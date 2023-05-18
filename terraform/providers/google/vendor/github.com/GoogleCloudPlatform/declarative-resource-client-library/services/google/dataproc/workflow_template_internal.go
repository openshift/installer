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
package dataproc

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

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
func (r *WorkflowTemplatePlacementManagedClusterConfig) validate() error {
	if !dcl.IsEmptyValueIndirect(r.GceClusterConfig) {
		if err := r.GceClusterConfig.validate(); err != nil {
			return err
		}
	}
	if !dcl.IsEmptyValueIndirect(r.MasterConfig) {
		if err := r.MasterConfig.validate(); err != nil {
			return err
		}
	}
	if !dcl.IsEmptyValueIndirect(r.WorkerConfig) {
		if err := r.WorkerConfig.validate(); err != nil {
			return err
		}
	}
	if !dcl.IsEmptyValueIndirect(r.SecondaryWorkerConfig) {
		if err := r.SecondaryWorkerConfig.validate(); err != nil {
			return err
		}
	}
	if !dcl.IsEmptyValueIndirect(r.SoftwareConfig) {
		if err := r.SoftwareConfig.validate(); err != nil {
			return err
		}
	}
	if !dcl.IsEmptyValueIndirect(r.EncryptionConfig) {
		if err := r.EncryptionConfig.validate(); err != nil {
			return err
		}
	}
	if !dcl.IsEmptyValueIndirect(r.AutoscalingConfig) {
		if err := r.AutoscalingConfig.validate(); err != nil {
			return err
		}
	}
	if !dcl.IsEmptyValueIndirect(r.SecurityConfig) {
		if err := r.SecurityConfig.validate(); err != nil {
			return err
		}
	}
	if !dcl.IsEmptyValueIndirect(r.LifecycleConfig) {
		if err := r.LifecycleConfig.validate(); err != nil {
			return err
		}
	}
	if !dcl.IsEmptyValueIndirect(r.EndpointConfig) {
		if err := r.EndpointConfig.validate(); err != nil {
			return err
		}
	}
	return nil
}
func (r *WorkflowTemplatePlacementManagedClusterConfigGceClusterConfig) validate() error {
	if !dcl.IsEmptyValueIndirect(r.ReservationAffinity) {
		if err := r.ReservationAffinity.validate(); err != nil {
			return err
		}
	}
	if !dcl.IsEmptyValueIndirect(r.NodeGroupAffinity) {
		if err := r.NodeGroupAffinity.validate(); err != nil {
			return err
		}
	}
	if !dcl.IsEmptyValueIndirect(r.ShieldedInstanceConfig) {
		if err := r.ShieldedInstanceConfig.validate(); err != nil {
			return err
		}
	}
	return nil
}
func (r *WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigReservationAffinity) validate() error {
	return nil
}
func (r *WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigNodeGroupAffinity) validate() error {
	if err := dcl.Required(r, "nodeGroup"); err != nil {
		return err
	}
	return nil
}
func (r *WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigShieldedInstanceConfig) validate() error {
	return nil
}
func (r *WorkflowTemplatePlacementManagedClusterConfigMasterConfig) validate() error {
	if !dcl.IsEmptyValueIndirect(r.DiskConfig) {
		if err := r.DiskConfig.validate(); err != nil {
			return err
		}
	}
	if !dcl.IsEmptyValueIndirect(r.ManagedGroupConfig) {
		if err := r.ManagedGroupConfig.validate(); err != nil {
			return err
		}
	}
	return nil
}
func (r *WorkflowTemplatePlacementManagedClusterConfigMasterConfigDiskConfig) validate() error {
	return nil
}
func (r *WorkflowTemplatePlacementManagedClusterConfigMasterConfigManagedGroupConfig) validate() error {
	return nil
}
func (r *WorkflowTemplatePlacementManagedClusterConfigMasterConfigAccelerators) validate() error {
	return nil
}
func (r *WorkflowTemplatePlacementManagedClusterConfigWorkerConfig) validate() error {
	if !dcl.IsEmptyValueIndirect(r.DiskConfig) {
		if err := r.DiskConfig.validate(); err != nil {
			return err
		}
	}
	if !dcl.IsEmptyValueIndirect(r.ManagedGroupConfig) {
		if err := r.ManagedGroupConfig.validate(); err != nil {
			return err
		}
	}
	return nil
}
func (r *WorkflowTemplatePlacementManagedClusterConfigWorkerConfigDiskConfig) validate() error {
	return nil
}
func (r *WorkflowTemplatePlacementManagedClusterConfigWorkerConfigManagedGroupConfig) validate() error {
	return nil
}
func (r *WorkflowTemplatePlacementManagedClusterConfigWorkerConfigAccelerators) validate() error {
	return nil
}
func (r *WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfig) validate() error {
	if !dcl.IsEmptyValueIndirect(r.DiskConfig) {
		if err := r.DiskConfig.validate(); err != nil {
			return err
		}
	}
	if !dcl.IsEmptyValueIndirect(r.ManagedGroupConfig) {
		if err := r.ManagedGroupConfig.validate(); err != nil {
			return err
		}
	}
	return nil
}
func (r *WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigDiskConfig) validate() error {
	return nil
}
func (r *WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigManagedGroupConfig) validate() error {
	return nil
}
func (r *WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigAccelerators) validate() error {
	return nil
}
func (r *WorkflowTemplatePlacementManagedClusterConfigSoftwareConfig) validate() error {
	return nil
}
func (r *WorkflowTemplatePlacementManagedClusterConfigInitializationActions) validate() error {
	return nil
}
func (r *WorkflowTemplatePlacementManagedClusterConfigEncryptionConfig) validate() error {
	return nil
}
func (r *WorkflowTemplatePlacementManagedClusterConfigAutoscalingConfig) validate() error {
	return nil
}
func (r *WorkflowTemplatePlacementManagedClusterConfigSecurityConfig) validate() error {
	if !dcl.IsEmptyValueIndirect(r.KerberosConfig) {
		if err := r.KerberosConfig.validate(); err != nil {
			return err
		}
	}
	return nil
}
func (r *WorkflowTemplatePlacementManagedClusterConfigSecurityConfigKerberosConfig) validate() error {
	return nil
}
func (r *WorkflowTemplatePlacementManagedClusterConfigLifecycleConfig) validate() error {
	return nil
}
func (r *WorkflowTemplatePlacementManagedClusterConfigEndpointConfig) validate() error {
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
		res, err := unmarshalMapWorkflowTemplate(v, c, r)
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

	// We saw a race condition where for some successful delete operation, the Get calls returned resources for a short duration.
	// This is the reason we are adding retry to handle that case.
	retriesRemaining := 10
	dcl.Do(ctx, func(ctx context.Context) (*dcl.RetryDetails, error) {
		_, err := c.GetWorkflowTemplate(ctx, r)
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

	// The Get call applies postReadExtract and so the result may contain fields that are not part of API version.
	if err := extractWorkflowTemplateFields(rawInitial); err != nil {
		return nil, nil, nil, err
	}

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
	if dcl.IsZeroValue(rawDesired.Labels) || (dcl.IsEmptyValueIndirect(rawDesired.Labels) && dcl.IsEmptyValueIndirect(rawInitial.Labels)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
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

	if dcl.IsEmptyValueIndirect(rawNew.Version) && dcl.IsEmptyValueIndirect(rawDesired.Version) {
		rawNew.Version = rawDesired.Version
	} else {
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

	if dcl.IsEmptyValueIndirect(rawNew.Placement) && dcl.IsEmptyValueIndirect(rawDesired.Placement) {
		rawNew.Placement = rawDesired.Placement
	} else {
		rawNew.Placement = canonicalizeNewWorkflowTemplatePlacement(c, rawDesired.Placement, rawNew.Placement)
	}

	if dcl.IsEmptyValueIndirect(rawNew.Jobs) && dcl.IsEmptyValueIndirect(rawDesired.Jobs) {
		rawNew.Jobs = rawDesired.Jobs
	} else {
		rawNew.Jobs = canonicalizeNewWorkflowTemplateJobsSlice(c, rawDesired.Jobs, rawNew.Jobs)
	}

	if dcl.IsEmptyValueIndirect(rawNew.Parameters) && dcl.IsEmptyValueIndirect(rawDesired.Parameters) {
		rawNew.Parameters = rawDesired.Parameters
	} else {
		rawNew.Parameters = canonicalizeNewWorkflowTemplateParametersSlice(c, rawDesired.Parameters, rawNew.Parameters)
	}

	if dcl.IsEmptyValueIndirect(rawNew.DagTimeout) && dcl.IsEmptyValueIndirect(rawDesired.DagTimeout) {
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
	if dcl.IsEmptyValueIndirect(des) {
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
		if dcl.IsEmptyValueIndirect(des) {
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

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []WorkflowTemplatePlacement
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareWorkflowTemplatePlacementNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewWorkflowTemplatePlacement(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
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
	cDes.Config = canonicalizeWorkflowTemplatePlacementManagedClusterConfig(des.Config, initial.Config, opts...)
	if dcl.IsZeroValue(des.Labels) || (dcl.IsEmptyValueIndirect(des.Labels) && dcl.IsEmptyValueIndirect(initial.Labels)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.Labels = initial.Labels
	} else {
		cDes.Labels = des.Labels
	}

	return cDes
}

func canonicalizeWorkflowTemplatePlacementManagedClusterSlice(des, initial []WorkflowTemplatePlacementManagedCluster, opts ...dcl.ApplyOption) []WorkflowTemplatePlacementManagedCluster {
	if dcl.IsEmptyValueIndirect(des) {
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
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for WorkflowTemplatePlacementManagedCluster while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringCanonicalize(des.ClusterName, nw.ClusterName) {
		nw.ClusterName = des.ClusterName
	}
	nw.Config = canonicalizeNewWorkflowTemplatePlacementManagedClusterConfig(c, des.Config, nw.Config)

	return nw
}

func canonicalizeNewWorkflowTemplatePlacementManagedClusterSet(c *Client, des, nw []WorkflowTemplatePlacementManagedCluster) []WorkflowTemplatePlacementManagedCluster {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []WorkflowTemplatePlacementManagedCluster
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareWorkflowTemplatePlacementManagedClusterNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewWorkflowTemplatePlacementManagedCluster(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
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

func canonicalizeWorkflowTemplatePlacementManagedClusterConfig(des, initial *WorkflowTemplatePlacementManagedClusterConfig, opts ...dcl.ApplyOption) *WorkflowTemplatePlacementManagedClusterConfig {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &WorkflowTemplatePlacementManagedClusterConfig{}

	if dcl.IsZeroValue(des.StagingBucket) || (dcl.IsEmptyValueIndirect(des.StagingBucket) && dcl.IsEmptyValueIndirect(initial.StagingBucket)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.StagingBucket = initial.StagingBucket
	} else {
		cDes.StagingBucket = des.StagingBucket
	}
	if dcl.IsZeroValue(des.TempBucket) || (dcl.IsEmptyValueIndirect(des.TempBucket) && dcl.IsEmptyValueIndirect(initial.TempBucket)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.TempBucket = initial.TempBucket
	} else {
		cDes.TempBucket = des.TempBucket
	}
	cDes.GceClusterConfig = canonicalizeWorkflowTemplatePlacementManagedClusterConfigGceClusterConfig(des.GceClusterConfig, initial.GceClusterConfig, opts...)
	cDes.MasterConfig = canonicalizeWorkflowTemplatePlacementManagedClusterConfigMasterConfig(des.MasterConfig, initial.MasterConfig, opts...)
	cDes.WorkerConfig = canonicalizeWorkflowTemplatePlacementManagedClusterConfigWorkerConfig(des.WorkerConfig, initial.WorkerConfig, opts...)
	cDes.SecondaryWorkerConfig = canonicalizeWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfig(des.SecondaryWorkerConfig, initial.SecondaryWorkerConfig, opts...)
	cDes.SoftwareConfig = canonicalizeWorkflowTemplatePlacementManagedClusterConfigSoftwareConfig(des.SoftwareConfig, initial.SoftwareConfig, opts...)
	cDes.InitializationActions = canonicalizeWorkflowTemplatePlacementManagedClusterConfigInitializationActionsSlice(des.InitializationActions, initial.InitializationActions, opts...)
	cDes.EncryptionConfig = canonicalizeWorkflowTemplatePlacementManagedClusterConfigEncryptionConfig(des.EncryptionConfig, initial.EncryptionConfig, opts...)
	cDes.AutoscalingConfig = canonicalizeWorkflowTemplatePlacementManagedClusterConfigAutoscalingConfig(des.AutoscalingConfig, initial.AutoscalingConfig, opts...)
	cDes.SecurityConfig = canonicalizeWorkflowTemplatePlacementManagedClusterConfigSecurityConfig(des.SecurityConfig, initial.SecurityConfig, opts...)
	cDes.LifecycleConfig = canonicalizeWorkflowTemplatePlacementManagedClusterConfigLifecycleConfig(des.LifecycleConfig, initial.LifecycleConfig, opts...)
	cDes.EndpointConfig = canonicalizeWorkflowTemplatePlacementManagedClusterConfigEndpointConfig(des.EndpointConfig, initial.EndpointConfig, opts...)

	return cDes
}

func canonicalizeWorkflowTemplatePlacementManagedClusterConfigSlice(des, initial []WorkflowTemplatePlacementManagedClusterConfig, opts ...dcl.ApplyOption) []WorkflowTemplatePlacementManagedClusterConfig {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]WorkflowTemplatePlacementManagedClusterConfig, 0, len(des))
		for _, d := range des {
			cd := canonicalizeWorkflowTemplatePlacementManagedClusterConfig(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]WorkflowTemplatePlacementManagedClusterConfig, 0, len(des))
	for i, d := range des {
		cd := canonicalizeWorkflowTemplatePlacementManagedClusterConfig(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewWorkflowTemplatePlacementManagedClusterConfig(c *Client, des, nw *WorkflowTemplatePlacementManagedClusterConfig) *WorkflowTemplatePlacementManagedClusterConfig {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for WorkflowTemplatePlacementManagedClusterConfig while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	nw.GceClusterConfig = canonicalizeNewWorkflowTemplatePlacementManagedClusterConfigGceClusterConfig(c, des.GceClusterConfig, nw.GceClusterConfig)
	nw.MasterConfig = canonicalizeNewWorkflowTemplatePlacementManagedClusterConfigMasterConfig(c, des.MasterConfig, nw.MasterConfig)
	nw.WorkerConfig = canonicalizeNewWorkflowTemplatePlacementManagedClusterConfigWorkerConfig(c, des.WorkerConfig, nw.WorkerConfig)
	nw.SecondaryWorkerConfig = canonicalizeNewWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfig(c, des.SecondaryWorkerConfig, nw.SecondaryWorkerConfig)
	nw.SoftwareConfig = canonicalizeNewWorkflowTemplatePlacementManagedClusterConfigSoftwareConfig(c, des.SoftwareConfig, nw.SoftwareConfig)
	nw.InitializationActions = canonicalizeNewWorkflowTemplatePlacementManagedClusterConfigInitializationActionsSlice(c, des.InitializationActions, nw.InitializationActions)
	nw.EncryptionConfig = canonicalizeNewWorkflowTemplatePlacementManagedClusterConfigEncryptionConfig(c, des.EncryptionConfig, nw.EncryptionConfig)
	nw.AutoscalingConfig = canonicalizeNewWorkflowTemplatePlacementManagedClusterConfigAutoscalingConfig(c, des.AutoscalingConfig, nw.AutoscalingConfig)
	nw.SecurityConfig = canonicalizeNewWorkflowTemplatePlacementManagedClusterConfigSecurityConfig(c, des.SecurityConfig, nw.SecurityConfig)
	nw.LifecycleConfig = canonicalizeNewWorkflowTemplatePlacementManagedClusterConfigLifecycleConfig(c, des.LifecycleConfig, nw.LifecycleConfig)
	nw.EndpointConfig = canonicalizeNewWorkflowTemplatePlacementManagedClusterConfigEndpointConfig(c, des.EndpointConfig, nw.EndpointConfig)

	return nw
}

func canonicalizeNewWorkflowTemplatePlacementManagedClusterConfigSet(c *Client, des, nw []WorkflowTemplatePlacementManagedClusterConfig) []WorkflowTemplatePlacementManagedClusterConfig {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []WorkflowTemplatePlacementManagedClusterConfig
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareWorkflowTemplatePlacementManagedClusterConfigNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewWorkflowTemplatePlacementManagedClusterConfig(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewWorkflowTemplatePlacementManagedClusterConfigSlice(c *Client, des, nw []WorkflowTemplatePlacementManagedClusterConfig) []WorkflowTemplatePlacementManagedClusterConfig {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []WorkflowTemplatePlacementManagedClusterConfig
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewWorkflowTemplatePlacementManagedClusterConfig(c, &d, &n))
	}

	return items
}

func canonicalizeWorkflowTemplatePlacementManagedClusterConfigGceClusterConfig(des, initial *WorkflowTemplatePlacementManagedClusterConfigGceClusterConfig, opts ...dcl.ApplyOption) *WorkflowTemplatePlacementManagedClusterConfigGceClusterConfig {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &WorkflowTemplatePlacementManagedClusterConfigGceClusterConfig{}

	if dcl.StringCanonicalize(des.Zone, initial.Zone) || dcl.IsZeroValue(des.Zone) {
		cDes.Zone = initial.Zone
	} else {
		cDes.Zone = des.Zone
	}
	if dcl.IsZeroValue(des.Network) || (dcl.IsEmptyValueIndirect(des.Network) && dcl.IsEmptyValueIndirect(initial.Network)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.Network = initial.Network
	} else {
		cDes.Network = des.Network
	}
	if dcl.IsZeroValue(des.Subnetwork) || (dcl.IsEmptyValueIndirect(des.Subnetwork) && dcl.IsEmptyValueIndirect(initial.Subnetwork)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.Subnetwork = initial.Subnetwork
	} else {
		cDes.Subnetwork = des.Subnetwork
	}
	if dcl.BoolCanonicalize(des.InternalIPOnly, initial.InternalIPOnly) || dcl.IsZeroValue(des.InternalIPOnly) {
		cDes.InternalIPOnly = initial.InternalIPOnly
	} else {
		cDes.InternalIPOnly = des.InternalIPOnly
	}
	if dcl.IsZeroValue(des.PrivateIPv6GoogleAccess) || (dcl.IsEmptyValueIndirect(des.PrivateIPv6GoogleAccess) && dcl.IsEmptyValueIndirect(initial.PrivateIPv6GoogleAccess)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.PrivateIPv6GoogleAccess = initial.PrivateIPv6GoogleAccess
	} else {
		cDes.PrivateIPv6GoogleAccess = des.PrivateIPv6GoogleAccess
	}
	if dcl.IsZeroValue(des.ServiceAccount) || (dcl.IsEmptyValueIndirect(des.ServiceAccount) && dcl.IsEmptyValueIndirect(initial.ServiceAccount)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.ServiceAccount = initial.ServiceAccount
	} else {
		cDes.ServiceAccount = des.ServiceAccount
	}
	if dcl.StringArrayCanonicalize(des.ServiceAccountScopes, initial.ServiceAccountScopes) {
		cDes.ServiceAccountScopes = initial.ServiceAccountScopes
	} else {
		cDes.ServiceAccountScopes = des.ServiceAccountScopes
	}
	if dcl.StringArrayCanonicalize(des.Tags, initial.Tags) {
		cDes.Tags = initial.Tags
	} else {
		cDes.Tags = des.Tags
	}
	if dcl.IsZeroValue(des.Metadata) || (dcl.IsEmptyValueIndirect(des.Metadata) && dcl.IsEmptyValueIndirect(initial.Metadata)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.Metadata = initial.Metadata
	} else {
		cDes.Metadata = des.Metadata
	}
	cDes.ReservationAffinity = canonicalizeWorkflowTemplatePlacementManagedClusterConfigGceClusterConfigReservationAffinity(des.ReservationAffinity, initial.ReservationAffinity, opts...)
	cDes.NodeGroupAffinity = canonicalizeWorkflowTemplatePlacementManagedClusterConfigGceClusterConfigNodeGroupAffinity(des.NodeGroupAffinity, initial.NodeGroupAffinity, opts...)
	cDes.ShieldedInstanceConfig = canonicalizeWorkflowTemplatePlacementManagedClusterConfigGceClusterConfigShieldedInstanceConfig(des.ShieldedInstanceConfig, initial.ShieldedInstanceConfig, opts...)

	return cDes
}

func canonicalizeWorkflowTemplatePlacementManagedClusterConfigGceClusterConfigSlice(des, initial []WorkflowTemplatePlacementManagedClusterConfigGceClusterConfig, opts ...dcl.ApplyOption) []WorkflowTemplatePlacementManagedClusterConfigGceClusterConfig {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]WorkflowTemplatePlacementManagedClusterConfigGceClusterConfig, 0, len(des))
		for _, d := range des {
			cd := canonicalizeWorkflowTemplatePlacementManagedClusterConfigGceClusterConfig(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]WorkflowTemplatePlacementManagedClusterConfigGceClusterConfig, 0, len(des))
	for i, d := range des {
		cd := canonicalizeWorkflowTemplatePlacementManagedClusterConfigGceClusterConfig(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewWorkflowTemplatePlacementManagedClusterConfigGceClusterConfig(c *Client, des, nw *WorkflowTemplatePlacementManagedClusterConfigGceClusterConfig) *WorkflowTemplatePlacementManagedClusterConfigGceClusterConfig {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for WorkflowTemplatePlacementManagedClusterConfigGceClusterConfig while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringCanonicalize(des.Zone, nw.Zone) {
		nw.Zone = des.Zone
	}
	if dcl.BoolCanonicalize(des.InternalIPOnly, nw.InternalIPOnly) {
		nw.InternalIPOnly = des.InternalIPOnly
	}
	if dcl.StringArrayCanonicalize(des.ServiceAccountScopes, nw.ServiceAccountScopes) {
		nw.ServiceAccountScopes = des.ServiceAccountScopes
	}
	if dcl.StringArrayCanonicalize(des.Tags, nw.Tags) {
		nw.Tags = des.Tags
	}
	nw.ReservationAffinity = canonicalizeNewWorkflowTemplatePlacementManagedClusterConfigGceClusterConfigReservationAffinity(c, des.ReservationAffinity, nw.ReservationAffinity)
	nw.NodeGroupAffinity = canonicalizeNewWorkflowTemplatePlacementManagedClusterConfigGceClusterConfigNodeGroupAffinity(c, des.NodeGroupAffinity, nw.NodeGroupAffinity)
	nw.ShieldedInstanceConfig = canonicalizeNewWorkflowTemplatePlacementManagedClusterConfigGceClusterConfigShieldedInstanceConfig(c, des.ShieldedInstanceConfig, nw.ShieldedInstanceConfig)

	return nw
}

func canonicalizeNewWorkflowTemplatePlacementManagedClusterConfigGceClusterConfigSet(c *Client, des, nw []WorkflowTemplatePlacementManagedClusterConfigGceClusterConfig) []WorkflowTemplatePlacementManagedClusterConfigGceClusterConfig {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []WorkflowTemplatePlacementManagedClusterConfigGceClusterConfig
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareWorkflowTemplatePlacementManagedClusterConfigGceClusterConfigNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewWorkflowTemplatePlacementManagedClusterConfigGceClusterConfig(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewWorkflowTemplatePlacementManagedClusterConfigGceClusterConfigSlice(c *Client, des, nw []WorkflowTemplatePlacementManagedClusterConfigGceClusterConfig) []WorkflowTemplatePlacementManagedClusterConfigGceClusterConfig {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []WorkflowTemplatePlacementManagedClusterConfigGceClusterConfig
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewWorkflowTemplatePlacementManagedClusterConfigGceClusterConfig(c, &d, &n))
	}

	return items
}

func canonicalizeWorkflowTemplatePlacementManagedClusterConfigGceClusterConfigReservationAffinity(des, initial *WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigReservationAffinity, opts ...dcl.ApplyOption) *WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigReservationAffinity {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigReservationAffinity{}

	if dcl.IsZeroValue(des.ConsumeReservationType) || (dcl.IsEmptyValueIndirect(des.ConsumeReservationType) && dcl.IsEmptyValueIndirect(initial.ConsumeReservationType)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.ConsumeReservationType = initial.ConsumeReservationType
	} else {
		cDes.ConsumeReservationType = des.ConsumeReservationType
	}
	if dcl.StringCanonicalize(des.Key, initial.Key) || dcl.IsZeroValue(des.Key) {
		cDes.Key = initial.Key
	} else {
		cDes.Key = des.Key
	}
	if dcl.StringArrayCanonicalize(des.Values, initial.Values) {
		cDes.Values = initial.Values
	} else {
		cDes.Values = des.Values
	}

	return cDes
}

func canonicalizeWorkflowTemplatePlacementManagedClusterConfigGceClusterConfigReservationAffinitySlice(des, initial []WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigReservationAffinity, opts ...dcl.ApplyOption) []WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigReservationAffinity {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigReservationAffinity, 0, len(des))
		for _, d := range des {
			cd := canonicalizeWorkflowTemplatePlacementManagedClusterConfigGceClusterConfigReservationAffinity(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigReservationAffinity, 0, len(des))
	for i, d := range des {
		cd := canonicalizeWorkflowTemplatePlacementManagedClusterConfigGceClusterConfigReservationAffinity(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewWorkflowTemplatePlacementManagedClusterConfigGceClusterConfigReservationAffinity(c *Client, des, nw *WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigReservationAffinity) *WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigReservationAffinity {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigReservationAffinity while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringCanonicalize(des.Key, nw.Key) {
		nw.Key = des.Key
	}
	if dcl.StringArrayCanonicalize(des.Values, nw.Values) {
		nw.Values = des.Values
	}

	return nw
}

func canonicalizeNewWorkflowTemplatePlacementManagedClusterConfigGceClusterConfigReservationAffinitySet(c *Client, des, nw []WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigReservationAffinity) []WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigReservationAffinity {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigReservationAffinity
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareWorkflowTemplatePlacementManagedClusterConfigGceClusterConfigReservationAffinityNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewWorkflowTemplatePlacementManagedClusterConfigGceClusterConfigReservationAffinity(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewWorkflowTemplatePlacementManagedClusterConfigGceClusterConfigReservationAffinitySlice(c *Client, des, nw []WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigReservationAffinity) []WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigReservationAffinity {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigReservationAffinity
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewWorkflowTemplatePlacementManagedClusterConfigGceClusterConfigReservationAffinity(c, &d, &n))
	}

	return items
}

func canonicalizeWorkflowTemplatePlacementManagedClusterConfigGceClusterConfigNodeGroupAffinity(des, initial *WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigNodeGroupAffinity, opts ...dcl.ApplyOption) *WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigNodeGroupAffinity {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigNodeGroupAffinity{}

	if dcl.IsZeroValue(des.NodeGroup) || (dcl.IsEmptyValueIndirect(des.NodeGroup) && dcl.IsEmptyValueIndirect(initial.NodeGroup)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.NodeGroup = initial.NodeGroup
	} else {
		cDes.NodeGroup = des.NodeGroup
	}

	return cDes
}

func canonicalizeWorkflowTemplatePlacementManagedClusterConfigGceClusterConfigNodeGroupAffinitySlice(des, initial []WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigNodeGroupAffinity, opts ...dcl.ApplyOption) []WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigNodeGroupAffinity {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigNodeGroupAffinity, 0, len(des))
		for _, d := range des {
			cd := canonicalizeWorkflowTemplatePlacementManagedClusterConfigGceClusterConfigNodeGroupAffinity(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigNodeGroupAffinity, 0, len(des))
	for i, d := range des {
		cd := canonicalizeWorkflowTemplatePlacementManagedClusterConfigGceClusterConfigNodeGroupAffinity(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewWorkflowTemplatePlacementManagedClusterConfigGceClusterConfigNodeGroupAffinity(c *Client, des, nw *WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigNodeGroupAffinity) *WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigNodeGroupAffinity {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigNodeGroupAffinity while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	return nw
}

func canonicalizeNewWorkflowTemplatePlacementManagedClusterConfigGceClusterConfigNodeGroupAffinitySet(c *Client, des, nw []WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigNodeGroupAffinity) []WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigNodeGroupAffinity {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigNodeGroupAffinity
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareWorkflowTemplatePlacementManagedClusterConfigGceClusterConfigNodeGroupAffinityNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewWorkflowTemplatePlacementManagedClusterConfigGceClusterConfigNodeGroupAffinity(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewWorkflowTemplatePlacementManagedClusterConfigGceClusterConfigNodeGroupAffinitySlice(c *Client, des, nw []WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigNodeGroupAffinity) []WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigNodeGroupAffinity {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigNodeGroupAffinity
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewWorkflowTemplatePlacementManagedClusterConfigGceClusterConfigNodeGroupAffinity(c, &d, &n))
	}

	return items
}

func canonicalizeWorkflowTemplatePlacementManagedClusterConfigGceClusterConfigShieldedInstanceConfig(des, initial *WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigShieldedInstanceConfig, opts ...dcl.ApplyOption) *WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigShieldedInstanceConfig {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigShieldedInstanceConfig{}

	if dcl.BoolCanonicalize(des.EnableSecureBoot, initial.EnableSecureBoot) || dcl.IsZeroValue(des.EnableSecureBoot) {
		cDes.EnableSecureBoot = initial.EnableSecureBoot
	} else {
		cDes.EnableSecureBoot = des.EnableSecureBoot
	}
	if dcl.BoolCanonicalize(des.EnableVtpm, initial.EnableVtpm) || dcl.IsZeroValue(des.EnableVtpm) {
		cDes.EnableVtpm = initial.EnableVtpm
	} else {
		cDes.EnableVtpm = des.EnableVtpm
	}
	if dcl.BoolCanonicalize(des.EnableIntegrityMonitoring, initial.EnableIntegrityMonitoring) || dcl.IsZeroValue(des.EnableIntegrityMonitoring) {
		cDes.EnableIntegrityMonitoring = initial.EnableIntegrityMonitoring
	} else {
		cDes.EnableIntegrityMonitoring = des.EnableIntegrityMonitoring
	}

	return cDes
}

func canonicalizeWorkflowTemplatePlacementManagedClusterConfigGceClusterConfigShieldedInstanceConfigSlice(des, initial []WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigShieldedInstanceConfig, opts ...dcl.ApplyOption) []WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigShieldedInstanceConfig {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigShieldedInstanceConfig, 0, len(des))
		for _, d := range des {
			cd := canonicalizeWorkflowTemplatePlacementManagedClusterConfigGceClusterConfigShieldedInstanceConfig(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigShieldedInstanceConfig, 0, len(des))
	for i, d := range des {
		cd := canonicalizeWorkflowTemplatePlacementManagedClusterConfigGceClusterConfigShieldedInstanceConfig(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewWorkflowTemplatePlacementManagedClusterConfigGceClusterConfigShieldedInstanceConfig(c *Client, des, nw *WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigShieldedInstanceConfig) *WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigShieldedInstanceConfig {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigShieldedInstanceConfig while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.BoolCanonicalize(des.EnableSecureBoot, nw.EnableSecureBoot) {
		nw.EnableSecureBoot = des.EnableSecureBoot
	}
	if dcl.BoolCanonicalize(des.EnableVtpm, nw.EnableVtpm) {
		nw.EnableVtpm = des.EnableVtpm
	}
	if dcl.BoolCanonicalize(des.EnableIntegrityMonitoring, nw.EnableIntegrityMonitoring) {
		nw.EnableIntegrityMonitoring = des.EnableIntegrityMonitoring
	}

	return nw
}

func canonicalizeNewWorkflowTemplatePlacementManagedClusterConfigGceClusterConfigShieldedInstanceConfigSet(c *Client, des, nw []WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigShieldedInstanceConfig) []WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigShieldedInstanceConfig {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigShieldedInstanceConfig
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareWorkflowTemplatePlacementManagedClusterConfigGceClusterConfigShieldedInstanceConfigNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewWorkflowTemplatePlacementManagedClusterConfigGceClusterConfigShieldedInstanceConfig(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewWorkflowTemplatePlacementManagedClusterConfigGceClusterConfigShieldedInstanceConfigSlice(c *Client, des, nw []WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigShieldedInstanceConfig) []WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigShieldedInstanceConfig {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigShieldedInstanceConfig
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewWorkflowTemplatePlacementManagedClusterConfigGceClusterConfigShieldedInstanceConfig(c, &d, &n))
	}

	return items
}

func canonicalizeWorkflowTemplatePlacementManagedClusterConfigMasterConfig(des, initial *WorkflowTemplatePlacementManagedClusterConfigMasterConfig, opts ...dcl.ApplyOption) *WorkflowTemplatePlacementManagedClusterConfigMasterConfig {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &WorkflowTemplatePlacementManagedClusterConfigMasterConfig{}

	if dcl.IsZeroValue(des.NumInstances) || (dcl.IsEmptyValueIndirect(des.NumInstances) && dcl.IsEmptyValueIndirect(initial.NumInstances)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.NumInstances = initial.NumInstances
	} else {
		cDes.NumInstances = des.NumInstances
	}
	if dcl.IsZeroValue(des.Image) || (dcl.IsEmptyValueIndirect(des.Image) && dcl.IsEmptyValueIndirect(initial.Image)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.Image = initial.Image
	} else {
		cDes.Image = des.Image
	}
	if dcl.StringCanonicalize(des.MachineType, initial.MachineType) || dcl.IsZeroValue(des.MachineType) {
		cDes.MachineType = initial.MachineType
	} else {
		cDes.MachineType = des.MachineType
	}
	cDes.DiskConfig = canonicalizeWorkflowTemplatePlacementManagedClusterConfigMasterConfigDiskConfig(des.DiskConfig, initial.DiskConfig, opts...)
	if dcl.IsZeroValue(des.Preemptibility) || (dcl.IsEmptyValueIndirect(des.Preemptibility) && dcl.IsEmptyValueIndirect(initial.Preemptibility)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.Preemptibility = initial.Preemptibility
	} else {
		cDes.Preemptibility = des.Preemptibility
	}
	cDes.Accelerators = canonicalizeWorkflowTemplatePlacementManagedClusterConfigMasterConfigAcceleratorsSlice(des.Accelerators, initial.Accelerators, opts...)
	if dcl.StringCanonicalize(des.MinCpuPlatform, initial.MinCpuPlatform) || dcl.IsZeroValue(des.MinCpuPlatform) {
		cDes.MinCpuPlatform = initial.MinCpuPlatform
	} else {
		cDes.MinCpuPlatform = des.MinCpuPlatform
	}

	return cDes
}

func canonicalizeWorkflowTemplatePlacementManagedClusterConfigMasterConfigSlice(des, initial []WorkflowTemplatePlacementManagedClusterConfigMasterConfig, opts ...dcl.ApplyOption) []WorkflowTemplatePlacementManagedClusterConfigMasterConfig {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]WorkflowTemplatePlacementManagedClusterConfigMasterConfig, 0, len(des))
		for _, d := range des {
			cd := canonicalizeWorkflowTemplatePlacementManagedClusterConfigMasterConfig(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]WorkflowTemplatePlacementManagedClusterConfigMasterConfig, 0, len(des))
	for i, d := range des {
		cd := canonicalizeWorkflowTemplatePlacementManagedClusterConfigMasterConfig(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewWorkflowTemplatePlacementManagedClusterConfigMasterConfig(c *Client, des, nw *WorkflowTemplatePlacementManagedClusterConfigMasterConfig) *WorkflowTemplatePlacementManagedClusterConfigMasterConfig {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for WorkflowTemplatePlacementManagedClusterConfigMasterConfig while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringArrayCanonicalize(des.InstanceNames, nw.InstanceNames) {
		nw.InstanceNames = des.InstanceNames
	}
	if dcl.StringCanonicalize(des.MachineType, nw.MachineType) {
		nw.MachineType = des.MachineType
	}
	nw.DiskConfig = canonicalizeNewWorkflowTemplatePlacementManagedClusterConfigMasterConfigDiskConfig(c, des.DiskConfig, nw.DiskConfig)
	if dcl.BoolCanonicalize(des.IsPreemptible, nw.IsPreemptible) {
		nw.IsPreemptible = des.IsPreemptible
	}
	nw.ManagedGroupConfig = canonicalizeNewWorkflowTemplatePlacementManagedClusterConfigMasterConfigManagedGroupConfig(c, des.ManagedGroupConfig, nw.ManagedGroupConfig)
	nw.Accelerators = canonicalizeNewWorkflowTemplatePlacementManagedClusterConfigMasterConfigAcceleratorsSlice(c, des.Accelerators, nw.Accelerators)
	if dcl.StringCanonicalize(des.MinCpuPlatform, nw.MinCpuPlatform) {
		nw.MinCpuPlatform = des.MinCpuPlatform
	}

	return nw
}

func canonicalizeNewWorkflowTemplatePlacementManagedClusterConfigMasterConfigSet(c *Client, des, nw []WorkflowTemplatePlacementManagedClusterConfigMasterConfig) []WorkflowTemplatePlacementManagedClusterConfigMasterConfig {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []WorkflowTemplatePlacementManagedClusterConfigMasterConfig
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareWorkflowTemplatePlacementManagedClusterConfigMasterConfigNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewWorkflowTemplatePlacementManagedClusterConfigMasterConfig(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewWorkflowTemplatePlacementManagedClusterConfigMasterConfigSlice(c *Client, des, nw []WorkflowTemplatePlacementManagedClusterConfigMasterConfig) []WorkflowTemplatePlacementManagedClusterConfigMasterConfig {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []WorkflowTemplatePlacementManagedClusterConfigMasterConfig
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewWorkflowTemplatePlacementManagedClusterConfigMasterConfig(c, &d, &n))
	}

	return items
}

func canonicalizeWorkflowTemplatePlacementManagedClusterConfigMasterConfigDiskConfig(des, initial *WorkflowTemplatePlacementManagedClusterConfigMasterConfigDiskConfig, opts ...dcl.ApplyOption) *WorkflowTemplatePlacementManagedClusterConfigMasterConfigDiskConfig {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &WorkflowTemplatePlacementManagedClusterConfigMasterConfigDiskConfig{}

	if dcl.StringCanonicalize(des.BootDiskType, initial.BootDiskType) || dcl.IsZeroValue(des.BootDiskType) {
		cDes.BootDiskType = initial.BootDiskType
	} else {
		cDes.BootDiskType = des.BootDiskType
	}
	if dcl.IsZeroValue(des.BootDiskSizeGb) || (dcl.IsEmptyValueIndirect(des.BootDiskSizeGb) && dcl.IsEmptyValueIndirect(initial.BootDiskSizeGb)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.BootDiskSizeGb = initial.BootDiskSizeGb
	} else {
		cDes.BootDiskSizeGb = des.BootDiskSizeGb
	}
	if dcl.IsZeroValue(des.NumLocalSsds) || (dcl.IsEmptyValueIndirect(des.NumLocalSsds) && dcl.IsEmptyValueIndirect(initial.NumLocalSsds)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.NumLocalSsds = initial.NumLocalSsds
	} else {
		cDes.NumLocalSsds = des.NumLocalSsds
	}

	return cDes
}

func canonicalizeWorkflowTemplatePlacementManagedClusterConfigMasterConfigDiskConfigSlice(des, initial []WorkflowTemplatePlacementManagedClusterConfigMasterConfigDiskConfig, opts ...dcl.ApplyOption) []WorkflowTemplatePlacementManagedClusterConfigMasterConfigDiskConfig {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]WorkflowTemplatePlacementManagedClusterConfigMasterConfigDiskConfig, 0, len(des))
		for _, d := range des {
			cd := canonicalizeWorkflowTemplatePlacementManagedClusterConfigMasterConfigDiskConfig(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]WorkflowTemplatePlacementManagedClusterConfigMasterConfigDiskConfig, 0, len(des))
	for i, d := range des {
		cd := canonicalizeWorkflowTemplatePlacementManagedClusterConfigMasterConfigDiskConfig(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewWorkflowTemplatePlacementManagedClusterConfigMasterConfigDiskConfig(c *Client, des, nw *WorkflowTemplatePlacementManagedClusterConfigMasterConfigDiskConfig) *WorkflowTemplatePlacementManagedClusterConfigMasterConfigDiskConfig {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for WorkflowTemplatePlacementManagedClusterConfigMasterConfigDiskConfig while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringCanonicalize(des.BootDiskType, nw.BootDiskType) {
		nw.BootDiskType = des.BootDiskType
	}

	return nw
}

func canonicalizeNewWorkflowTemplatePlacementManagedClusterConfigMasterConfigDiskConfigSet(c *Client, des, nw []WorkflowTemplatePlacementManagedClusterConfigMasterConfigDiskConfig) []WorkflowTemplatePlacementManagedClusterConfigMasterConfigDiskConfig {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []WorkflowTemplatePlacementManagedClusterConfigMasterConfigDiskConfig
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareWorkflowTemplatePlacementManagedClusterConfigMasterConfigDiskConfigNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewWorkflowTemplatePlacementManagedClusterConfigMasterConfigDiskConfig(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewWorkflowTemplatePlacementManagedClusterConfigMasterConfigDiskConfigSlice(c *Client, des, nw []WorkflowTemplatePlacementManagedClusterConfigMasterConfigDiskConfig) []WorkflowTemplatePlacementManagedClusterConfigMasterConfigDiskConfig {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []WorkflowTemplatePlacementManagedClusterConfigMasterConfigDiskConfig
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewWorkflowTemplatePlacementManagedClusterConfigMasterConfigDiskConfig(c, &d, &n))
	}

	return items
}

func canonicalizeWorkflowTemplatePlacementManagedClusterConfigMasterConfigManagedGroupConfig(des, initial *WorkflowTemplatePlacementManagedClusterConfigMasterConfigManagedGroupConfig, opts ...dcl.ApplyOption) *WorkflowTemplatePlacementManagedClusterConfigMasterConfigManagedGroupConfig {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &WorkflowTemplatePlacementManagedClusterConfigMasterConfigManagedGroupConfig{}

	return cDes
}

func canonicalizeWorkflowTemplatePlacementManagedClusterConfigMasterConfigManagedGroupConfigSlice(des, initial []WorkflowTemplatePlacementManagedClusterConfigMasterConfigManagedGroupConfig, opts ...dcl.ApplyOption) []WorkflowTemplatePlacementManagedClusterConfigMasterConfigManagedGroupConfig {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]WorkflowTemplatePlacementManagedClusterConfigMasterConfigManagedGroupConfig, 0, len(des))
		for _, d := range des {
			cd := canonicalizeWorkflowTemplatePlacementManagedClusterConfigMasterConfigManagedGroupConfig(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]WorkflowTemplatePlacementManagedClusterConfigMasterConfigManagedGroupConfig, 0, len(des))
	for i, d := range des {
		cd := canonicalizeWorkflowTemplatePlacementManagedClusterConfigMasterConfigManagedGroupConfig(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewWorkflowTemplatePlacementManagedClusterConfigMasterConfigManagedGroupConfig(c *Client, des, nw *WorkflowTemplatePlacementManagedClusterConfigMasterConfigManagedGroupConfig) *WorkflowTemplatePlacementManagedClusterConfigMasterConfigManagedGroupConfig {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for WorkflowTemplatePlacementManagedClusterConfigMasterConfigManagedGroupConfig while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringCanonicalize(des.InstanceTemplateName, nw.InstanceTemplateName) {
		nw.InstanceTemplateName = des.InstanceTemplateName
	}
	if dcl.StringCanonicalize(des.InstanceGroupManagerName, nw.InstanceGroupManagerName) {
		nw.InstanceGroupManagerName = des.InstanceGroupManagerName
	}

	return nw
}

func canonicalizeNewWorkflowTemplatePlacementManagedClusterConfigMasterConfigManagedGroupConfigSet(c *Client, des, nw []WorkflowTemplatePlacementManagedClusterConfigMasterConfigManagedGroupConfig) []WorkflowTemplatePlacementManagedClusterConfigMasterConfigManagedGroupConfig {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []WorkflowTemplatePlacementManagedClusterConfigMasterConfigManagedGroupConfig
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareWorkflowTemplatePlacementManagedClusterConfigMasterConfigManagedGroupConfigNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewWorkflowTemplatePlacementManagedClusterConfigMasterConfigManagedGroupConfig(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewWorkflowTemplatePlacementManagedClusterConfigMasterConfigManagedGroupConfigSlice(c *Client, des, nw []WorkflowTemplatePlacementManagedClusterConfigMasterConfigManagedGroupConfig) []WorkflowTemplatePlacementManagedClusterConfigMasterConfigManagedGroupConfig {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []WorkflowTemplatePlacementManagedClusterConfigMasterConfigManagedGroupConfig
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewWorkflowTemplatePlacementManagedClusterConfigMasterConfigManagedGroupConfig(c, &d, &n))
	}

	return items
}

func canonicalizeWorkflowTemplatePlacementManagedClusterConfigMasterConfigAccelerators(des, initial *WorkflowTemplatePlacementManagedClusterConfigMasterConfigAccelerators, opts ...dcl.ApplyOption) *WorkflowTemplatePlacementManagedClusterConfigMasterConfigAccelerators {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &WorkflowTemplatePlacementManagedClusterConfigMasterConfigAccelerators{}

	if dcl.StringCanonicalize(des.AcceleratorType, initial.AcceleratorType) || dcl.IsZeroValue(des.AcceleratorType) {
		cDes.AcceleratorType = initial.AcceleratorType
	} else {
		cDes.AcceleratorType = des.AcceleratorType
	}
	if dcl.IsZeroValue(des.AcceleratorCount) || (dcl.IsEmptyValueIndirect(des.AcceleratorCount) && dcl.IsEmptyValueIndirect(initial.AcceleratorCount)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.AcceleratorCount = initial.AcceleratorCount
	} else {
		cDes.AcceleratorCount = des.AcceleratorCount
	}

	return cDes
}

func canonicalizeWorkflowTemplatePlacementManagedClusterConfigMasterConfigAcceleratorsSlice(des, initial []WorkflowTemplatePlacementManagedClusterConfigMasterConfigAccelerators, opts ...dcl.ApplyOption) []WorkflowTemplatePlacementManagedClusterConfigMasterConfigAccelerators {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]WorkflowTemplatePlacementManagedClusterConfigMasterConfigAccelerators, 0, len(des))
		for _, d := range des {
			cd := canonicalizeWorkflowTemplatePlacementManagedClusterConfigMasterConfigAccelerators(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]WorkflowTemplatePlacementManagedClusterConfigMasterConfigAccelerators, 0, len(des))
	for i, d := range des {
		cd := canonicalizeWorkflowTemplatePlacementManagedClusterConfigMasterConfigAccelerators(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewWorkflowTemplatePlacementManagedClusterConfigMasterConfigAccelerators(c *Client, des, nw *WorkflowTemplatePlacementManagedClusterConfigMasterConfigAccelerators) *WorkflowTemplatePlacementManagedClusterConfigMasterConfigAccelerators {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for WorkflowTemplatePlacementManagedClusterConfigMasterConfigAccelerators while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringCanonicalize(des.AcceleratorType, nw.AcceleratorType) {
		nw.AcceleratorType = des.AcceleratorType
	}

	return nw
}

func canonicalizeNewWorkflowTemplatePlacementManagedClusterConfigMasterConfigAcceleratorsSet(c *Client, des, nw []WorkflowTemplatePlacementManagedClusterConfigMasterConfigAccelerators) []WorkflowTemplatePlacementManagedClusterConfigMasterConfigAccelerators {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []WorkflowTemplatePlacementManagedClusterConfigMasterConfigAccelerators
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareWorkflowTemplatePlacementManagedClusterConfigMasterConfigAcceleratorsNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewWorkflowTemplatePlacementManagedClusterConfigMasterConfigAccelerators(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewWorkflowTemplatePlacementManagedClusterConfigMasterConfigAcceleratorsSlice(c *Client, des, nw []WorkflowTemplatePlacementManagedClusterConfigMasterConfigAccelerators) []WorkflowTemplatePlacementManagedClusterConfigMasterConfigAccelerators {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []WorkflowTemplatePlacementManagedClusterConfigMasterConfigAccelerators
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewWorkflowTemplatePlacementManagedClusterConfigMasterConfigAccelerators(c, &d, &n))
	}

	return items
}

func canonicalizeWorkflowTemplatePlacementManagedClusterConfigWorkerConfig(des, initial *WorkflowTemplatePlacementManagedClusterConfigWorkerConfig, opts ...dcl.ApplyOption) *WorkflowTemplatePlacementManagedClusterConfigWorkerConfig {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &WorkflowTemplatePlacementManagedClusterConfigWorkerConfig{}

	if dcl.IsZeroValue(des.NumInstances) || (dcl.IsEmptyValueIndirect(des.NumInstances) && dcl.IsEmptyValueIndirect(initial.NumInstances)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.NumInstances = initial.NumInstances
	} else {
		cDes.NumInstances = des.NumInstances
	}
	if dcl.IsZeroValue(des.Image) || (dcl.IsEmptyValueIndirect(des.Image) && dcl.IsEmptyValueIndirect(initial.Image)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.Image = initial.Image
	} else {
		cDes.Image = des.Image
	}
	if dcl.StringCanonicalize(des.MachineType, initial.MachineType) || dcl.IsZeroValue(des.MachineType) {
		cDes.MachineType = initial.MachineType
	} else {
		cDes.MachineType = des.MachineType
	}
	cDes.DiskConfig = canonicalizeWorkflowTemplatePlacementManagedClusterConfigWorkerConfigDiskConfig(des.DiskConfig, initial.DiskConfig, opts...)
	if dcl.IsZeroValue(des.Preemptibility) || (dcl.IsEmptyValueIndirect(des.Preemptibility) && dcl.IsEmptyValueIndirect(initial.Preemptibility)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.Preemptibility = initial.Preemptibility
	} else {
		cDes.Preemptibility = des.Preemptibility
	}
	cDes.Accelerators = canonicalizeWorkflowTemplatePlacementManagedClusterConfigWorkerConfigAcceleratorsSlice(des.Accelerators, initial.Accelerators, opts...)
	if dcl.StringCanonicalize(des.MinCpuPlatform, initial.MinCpuPlatform) || dcl.IsZeroValue(des.MinCpuPlatform) {
		cDes.MinCpuPlatform = initial.MinCpuPlatform
	} else {
		cDes.MinCpuPlatform = des.MinCpuPlatform
	}

	return cDes
}

func canonicalizeWorkflowTemplatePlacementManagedClusterConfigWorkerConfigSlice(des, initial []WorkflowTemplatePlacementManagedClusterConfigWorkerConfig, opts ...dcl.ApplyOption) []WorkflowTemplatePlacementManagedClusterConfigWorkerConfig {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]WorkflowTemplatePlacementManagedClusterConfigWorkerConfig, 0, len(des))
		for _, d := range des {
			cd := canonicalizeWorkflowTemplatePlacementManagedClusterConfigWorkerConfig(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]WorkflowTemplatePlacementManagedClusterConfigWorkerConfig, 0, len(des))
	for i, d := range des {
		cd := canonicalizeWorkflowTemplatePlacementManagedClusterConfigWorkerConfig(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewWorkflowTemplatePlacementManagedClusterConfigWorkerConfig(c *Client, des, nw *WorkflowTemplatePlacementManagedClusterConfigWorkerConfig) *WorkflowTemplatePlacementManagedClusterConfigWorkerConfig {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for WorkflowTemplatePlacementManagedClusterConfigWorkerConfig while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringArrayCanonicalize(des.InstanceNames, nw.InstanceNames) {
		nw.InstanceNames = des.InstanceNames
	}
	if dcl.StringCanonicalize(des.MachineType, nw.MachineType) {
		nw.MachineType = des.MachineType
	}
	nw.DiskConfig = canonicalizeNewWorkflowTemplatePlacementManagedClusterConfigWorkerConfigDiskConfig(c, des.DiskConfig, nw.DiskConfig)
	if dcl.BoolCanonicalize(des.IsPreemptible, nw.IsPreemptible) {
		nw.IsPreemptible = des.IsPreemptible
	}
	nw.ManagedGroupConfig = canonicalizeNewWorkflowTemplatePlacementManagedClusterConfigWorkerConfigManagedGroupConfig(c, des.ManagedGroupConfig, nw.ManagedGroupConfig)
	nw.Accelerators = canonicalizeNewWorkflowTemplatePlacementManagedClusterConfigWorkerConfigAcceleratorsSlice(c, des.Accelerators, nw.Accelerators)
	if dcl.StringCanonicalize(des.MinCpuPlatform, nw.MinCpuPlatform) {
		nw.MinCpuPlatform = des.MinCpuPlatform
	}

	return nw
}

func canonicalizeNewWorkflowTemplatePlacementManagedClusterConfigWorkerConfigSet(c *Client, des, nw []WorkflowTemplatePlacementManagedClusterConfigWorkerConfig) []WorkflowTemplatePlacementManagedClusterConfigWorkerConfig {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []WorkflowTemplatePlacementManagedClusterConfigWorkerConfig
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareWorkflowTemplatePlacementManagedClusterConfigWorkerConfigNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewWorkflowTemplatePlacementManagedClusterConfigWorkerConfig(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewWorkflowTemplatePlacementManagedClusterConfigWorkerConfigSlice(c *Client, des, nw []WorkflowTemplatePlacementManagedClusterConfigWorkerConfig) []WorkflowTemplatePlacementManagedClusterConfigWorkerConfig {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []WorkflowTemplatePlacementManagedClusterConfigWorkerConfig
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewWorkflowTemplatePlacementManagedClusterConfigWorkerConfig(c, &d, &n))
	}

	return items
}

func canonicalizeWorkflowTemplatePlacementManagedClusterConfigWorkerConfigDiskConfig(des, initial *WorkflowTemplatePlacementManagedClusterConfigWorkerConfigDiskConfig, opts ...dcl.ApplyOption) *WorkflowTemplatePlacementManagedClusterConfigWorkerConfigDiskConfig {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &WorkflowTemplatePlacementManagedClusterConfigWorkerConfigDiskConfig{}

	if dcl.StringCanonicalize(des.BootDiskType, initial.BootDiskType) || dcl.IsZeroValue(des.BootDiskType) {
		cDes.BootDiskType = initial.BootDiskType
	} else {
		cDes.BootDiskType = des.BootDiskType
	}
	if dcl.IsZeroValue(des.BootDiskSizeGb) || (dcl.IsEmptyValueIndirect(des.BootDiskSizeGb) && dcl.IsEmptyValueIndirect(initial.BootDiskSizeGb)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.BootDiskSizeGb = initial.BootDiskSizeGb
	} else {
		cDes.BootDiskSizeGb = des.BootDiskSizeGb
	}
	if dcl.IsZeroValue(des.NumLocalSsds) || (dcl.IsEmptyValueIndirect(des.NumLocalSsds) && dcl.IsEmptyValueIndirect(initial.NumLocalSsds)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.NumLocalSsds = initial.NumLocalSsds
	} else {
		cDes.NumLocalSsds = des.NumLocalSsds
	}

	return cDes
}

func canonicalizeWorkflowTemplatePlacementManagedClusterConfigWorkerConfigDiskConfigSlice(des, initial []WorkflowTemplatePlacementManagedClusterConfigWorkerConfigDiskConfig, opts ...dcl.ApplyOption) []WorkflowTemplatePlacementManagedClusterConfigWorkerConfigDiskConfig {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]WorkflowTemplatePlacementManagedClusterConfigWorkerConfigDiskConfig, 0, len(des))
		for _, d := range des {
			cd := canonicalizeWorkflowTemplatePlacementManagedClusterConfigWorkerConfigDiskConfig(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]WorkflowTemplatePlacementManagedClusterConfigWorkerConfigDiskConfig, 0, len(des))
	for i, d := range des {
		cd := canonicalizeWorkflowTemplatePlacementManagedClusterConfigWorkerConfigDiskConfig(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewWorkflowTemplatePlacementManagedClusterConfigWorkerConfigDiskConfig(c *Client, des, nw *WorkflowTemplatePlacementManagedClusterConfigWorkerConfigDiskConfig) *WorkflowTemplatePlacementManagedClusterConfigWorkerConfigDiskConfig {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for WorkflowTemplatePlacementManagedClusterConfigWorkerConfigDiskConfig while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringCanonicalize(des.BootDiskType, nw.BootDiskType) {
		nw.BootDiskType = des.BootDiskType
	}

	return nw
}

func canonicalizeNewWorkflowTemplatePlacementManagedClusterConfigWorkerConfigDiskConfigSet(c *Client, des, nw []WorkflowTemplatePlacementManagedClusterConfigWorkerConfigDiskConfig) []WorkflowTemplatePlacementManagedClusterConfigWorkerConfigDiskConfig {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []WorkflowTemplatePlacementManagedClusterConfigWorkerConfigDiskConfig
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareWorkflowTemplatePlacementManagedClusterConfigWorkerConfigDiskConfigNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewWorkflowTemplatePlacementManagedClusterConfigWorkerConfigDiskConfig(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewWorkflowTemplatePlacementManagedClusterConfigWorkerConfigDiskConfigSlice(c *Client, des, nw []WorkflowTemplatePlacementManagedClusterConfigWorkerConfigDiskConfig) []WorkflowTemplatePlacementManagedClusterConfigWorkerConfigDiskConfig {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []WorkflowTemplatePlacementManagedClusterConfigWorkerConfigDiskConfig
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewWorkflowTemplatePlacementManagedClusterConfigWorkerConfigDiskConfig(c, &d, &n))
	}

	return items
}

func canonicalizeWorkflowTemplatePlacementManagedClusterConfigWorkerConfigManagedGroupConfig(des, initial *WorkflowTemplatePlacementManagedClusterConfigWorkerConfigManagedGroupConfig, opts ...dcl.ApplyOption) *WorkflowTemplatePlacementManagedClusterConfigWorkerConfigManagedGroupConfig {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &WorkflowTemplatePlacementManagedClusterConfigWorkerConfigManagedGroupConfig{}

	return cDes
}

func canonicalizeWorkflowTemplatePlacementManagedClusterConfigWorkerConfigManagedGroupConfigSlice(des, initial []WorkflowTemplatePlacementManagedClusterConfigWorkerConfigManagedGroupConfig, opts ...dcl.ApplyOption) []WorkflowTemplatePlacementManagedClusterConfigWorkerConfigManagedGroupConfig {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]WorkflowTemplatePlacementManagedClusterConfigWorkerConfigManagedGroupConfig, 0, len(des))
		for _, d := range des {
			cd := canonicalizeWorkflowTemplatePlacementManagedClusterConfigWorkerConfigManagedGroupConfig(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]WorkflowTemplatePlacementManagedClusterConfigWorkerConfigManagedGroupConfig, 0, len(des))
	for i, d := range des {
		cd := canonicalizeWorkflowTemplatePlacementManagedClusterConfigWorkerConfigManagedGroupConfig(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewWorkflowTemplatePlacementManagedClusterConfigWorkerConfigManagedGroupConfig(c *Client, des, nw *WorkflowTemplatePlacementManagedClusterConfigWorkerConfigManagedGroupConfig) *WorkflowTemplatePlacementManagedClusterConfigWorkerConfigManagedGroupConfig {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for WorkflowTemplatePlacementManagedClusterConfigWorkerConfigManagedGroupConfig while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringCanonicalize(des.InstanceTemplateName, nw.InstanceTemplateName) {
		nw.InstanceTemplateName = des.InstanceTemplateName
	}
	if dcl.StringCanonicalize(des.InstanceGroupManagerName, nw.InstanceGroupManagerName) {
		nw.InstanceGroupManagerName = des.InstanceGroupManagerName
	}

	return nw
}

func canonicalizeNewWorkflowTemplatePlacementManagedClusterConfigWorkerConfigManagedGroupConfigSet(c *Client, des, nw []WorkflowTemplatePlacementManagedClusterConfigWorkerConfigManagedGroupConfig) []WorkflowTemplatePlacementManagedClusterConfigWorkerConfigManagedGroupConfig {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []WorkflowTemplatePlacementManagedClusterConfigWorkerConfigManagedGroupConfig
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareWorkflowTemplatePlacementManagedClusterConfigWorkerConfigManagedGroupConfigNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewWorkflowTemplatePlacementManagedClusterConfigWorkerConfigManagedGroupConfig(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewWorkflowTemplatePlacementManagedClusterConfigWorkerConfigManagedGroupConfigSlice(c *Client, des, nw []WorkflowTemplatePlacementManagedClusterConfigWorkerConfigManagedGroupConfig) []WorkflowTemplatePlacementManagedClusterConfigWorkerConfigManagedGroupConfig {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []WorkflowTemplatePlacementManagedClusterConfigWorkerConfigManagedGroupConfig
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewWorkflowTemplatePlacementManagedClusterConfigWorkerConfigManagedGroupConfig(c, &d, &n))
	}

	return items
}

func canonicalizeWorkflowTemplatePlacementManagedClusterConfigWorkerConfigAccelerators(des, initial *WorkflowTemplatePlacementManagedClusterConfigWorkerConfigAccelerators, opts ...dcl.ApplyOption) *WorkflowTemplatePlacementManagedClusterConfigWorkerConfigAccelerators {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &WorkflowTemplatePlacementManagedClusterConfigWorkerConfigAccelerators{}

	if dcl.StringCanonicalize(des.AcceleratorType, initial.AcceleratorType) || dcl.IsZeroValue(des.AcceleratorType) {
		cDes.AcceleratorType = initial.AcceleratorType
	} else {
		cDes.AcceleratorType = des.AcceleratorType
	}
	if dcl.IsZeroValue(des.AcceleratorCount) || (dcl.IsEmptyValueIndirect(des.AcceleratorCount) && dcl.IsEmptyValueIndirect(initial.AcceleratorCount)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.AcceleratorCount = initial.AcceleratorCount
	} else {
		cDes.AcceleratorCount = des.AcceleratorCount
	}

	return cDes
}

func canonicalizeWorkflowTemplatePlacementManagedClusterConfigWorkerConfigAcceleratorsSlice(des, initial []WorkflowTemplatePlacementManagedClusterConfigWorkerConfigAccelerators, opts ...dcl.ApplyOption) []WorkflowTemplatePlacementManagedClusterConfigWorkerConfigAccelerators {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]WorkflowTemplatePlacementManagedClusterConfigWorkerConfigAccelerators, 0, len(des))
		for _, d := range des {
			cd := canonicalizeWorkflowTemplatePlacementManagedClusterConfigWorkerConfigAccelerators(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]WorkflowTemplatePlacementManagedClusterConfigWorkerConfigAccelerators, 0, len(des))
	for i, d := range des {
		cd := canonicalizeWorkflowTemplatePlacementManagedClusterConfigWorkerConfigAccelerators(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewWorkflowTemplatePlacementManagedClusterConfigWorkerConfigAccelerators(c *Client, des, nw *WorkflowTemplatePlacementManagedClusterConfigWorkerConfigAccelerators) *WorkflowTemplatePlacementManagedClusterConfigWorkerConfigAccelerators {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for WorkflowTemplatePlacementManagedClusterConfigWorkerConfigAccelerators while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringCanonicalize(des.AcceleratorType, nw.AcceleratorType) {
		nw.AcceleratorType = des.AcceleratorType
	}

	return nw
}

func canonicalizeNewWorkflowTemplatePlacementManagedClusterConfigWorkerConfigAcceleratorsSet(c *Client, des, nw []WorkflowTemplatePlacementManagedClusterConfigWorkerConfigAccelerators) []WorkflowTemplatePlacementManagedClusterConfigWorkerConfigAccelerators {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []WorkflowTemplatePlacementManagedClusterConfigWorkerConfigAccelerators
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareWorkflowTemplatePlacementManagedClusterConfigWorkerConfigAcceleratorsNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewWorkflowTemplatePlacementManagedClusterConfigWorkerConfigAccelerators(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewWorkflowTemplatePlacementManagedClusterConfigWorkerConfigAcceleratorsSlice(c *Client, des, nw []WorkflowTemplatePlacementManagedClusterConfigWorkerConfigAccelerators) []WorkflowTemplatePlacementManagedClusterConfigWorkerConfigAccelerators {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []WorkflowTemplatePlacementManagedClusterConfigWorkerConfigAccelerators
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewWorkflowTemplatePlacementManagedClusterConfigWorkerConfigAccelerators(c, &d, &n))
	}

	return items
}

func canonicalizeWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfig(des, initial *WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfig, opts ...dcl.ApplyOption) *WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfig {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfig{}

	if dcl.IsZeroValue(des.NumInstances) || (dcl.IsEmptyValueIndirect(des.NumInstances) && dcl.IsEmptyValueIndirect(initial.NumInstances)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.NumInstances = initial.NumInstances
	} else {
		cDes.NumInstances = des.NumInstances
	}
	if dcl.IsZeroValue(des.Image) || (dcl.IsEmptyValueIndirect(des.Image) && dcl.IsEmptyValueIndirect(initial.Image)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.Image = initial.Image
	} else {
		cDes.Image = des.Image
	}
	if dcl.StringCanonicalize(des.MachineType, initial.MachineType) || dcl.IsZeroValue(des.MachineType) {
		cDes.MachineType = initial.MachineType
	} else {
		cDes.MachineType = des.MachineType
	}
	cDes.DiskConfig = canonicalizeWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigDiskConfig(des.DiskConfig, initial.DiskConfig, opts...)
	if dcl.IsZeroValue(des.Preemptibility) || (dcl.IsEmptyValueIndirect(des.Preemptibility) && dcl.IsEmptyValueIndirect(initial.Preemptibility)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.Preemptibility = initial.Preemptibility
	} else {
		cDes.Preemptibility = des.Preemptibility
	}
	cDes.Accelerators = canonicalizeWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigAcceleratorsSlice(des.Accelerators, initial.Accelerators, opts...)
	if dcl.StringCanonicalize(des.MinCpuPlatform, initial.MinCpuPlatform) || dcl.IsZeroValue(des.MinCpuPlatform) {
		cDes.MinCpuPlatform = initial.MinCpuPlatform
	} else {
		cDes.MinCpuPlatform = des.MinCpuPlatform
	}

	return cDes
}

func canonicalizeWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigSlice(des, initial []WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfig, opts ...dcl.ApplyOption) []WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfig {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfig, 0, len(des))
		for _, d := range des {
			cd := canonicalizeWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfig(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfig, 0, len(des))
	for i, d := range des {
		cd := canonicalizeWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfig(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfig(c *Client, des, nw *WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfig) *WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfig {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfig while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringArrayCanonicalize(des.InstanceNames, nw.InstanceNames) {
		nw.InstanceNames = des.InstanceNames
	}
	if dcl.StringCanonicalize(des.MachineType, nw.MachineType) {
		nw.MachineType = des.MachineType
	}
	nw.DiskConfig = canonicalizeNewWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigDiskConfig(c, des.DiskConfig, nw.DiskConfig)
	if dcl.BoolCanonicalize(des.IsPreemptible, nw.IsPreemptible) {
		nw.IsPreemptible = des.IsPreemptible
	}
	nw.ManagedGroupConfig = canonicalizeNewWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigManagedGroupConfig(c, des.ManagedGroupConfig, nw.ManagedGroupConfig)
	nw.Accelerators = canonicalizeNewWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigAcceleratorsSlice(c, des.Accelerators, nw.Accelerators)
	if dcl.StringCanonicalize(des.MinCpuPlatform, nw.MinCpuPlatform) {
		nw.MinCpuPlatform = des.MinCpuPlatform
	}

	return nw
}

func canonicalizeNewWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigSet(c *Client, des, nw []WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfig) []WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfig {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfig
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfig(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigSlice(c *Client, des, nw []WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfig) []WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfig {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfig
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfig(c, &d, &n))
	}

	return items
}

func canonicalizeWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigDiskConfig(des, initial *WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigDiskConfig, opts ...dcl.ApplyOption) *WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigDiskConfig {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigDiskConfig{}

	if dcl.StringCanonicalize(des.BootDiskType, initial.BootDiskType) || dcl.IsZeroValue(des.BootDiskType) {
		cDes.BootDiskType = initial.BootDiskType
	} else {
		cDes.BootDiskType = des.BootDiskType
	}
	if dcl.IsZeroValue(des.BootDiskSizeGb) || (dcl.IsEmptyValueIndirect(des.BootDiskSizeGb) && dcl.IsEmptyValueIndirect(initial.BootDiskSizeGb)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.BootDiskSizeGb = initial.BootDiskSizeGb
	} else {
		cDes.BootDiskSizeGb = des.BootDiskSizeGb
	}
	if dcl.IsZeroValue(des.NumLocalSsds) || (dcl.IsEmptyValueIndirect(des.NumLocalSsds) && dcl.IsEmptyValueIndirect(initial.NumLocalSsds)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.NumLocalSsds = initial.NumLocalSsds
	} else {
		cDes.NumLocalSsds = des.NumLocalSsds
	}

	return cDes
}

func canonicalizeWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigDiskConfigSlice(des, initial []WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigDiskConfig, opts ...dcl.ApplyOption) []WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigDiskConfig {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigDiskConfig, 0, len(des))
		for _, d := range des {
			cd := canonicalizeWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigDiskConfig(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigDiskConfig, 0, len(des))
	for i, d := range des {
		cd := canonicalizeWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigDiskConfig(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigDiskConfig(c *Client, des, nw *WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigDiskConfig) *WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigDiskConfig {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigDiskConfig while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringCanonicalize(des.BootDiskType, nw.BootDiskType) {
		nw.BootDiskType = des.BootDiskType
	}

	return nw
}

func canonicalizeNewWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigDiskConfigSet(c *Client, des, nw []WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigDiskConfig) []WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigDiskConfig {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigDiskConfig
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigDiskConfigNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigDiskConfig(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigDiskConfigSlice(c *Client, des, nw []WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigDiskConfig) []WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigDiskConfig {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigDiskConfig
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigDiskConfig(c, &d, &n))
	}

	return items
}

func canonicalizeWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigManagedGroupConfig(des, initial *WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigManagedGroupConfig, opts ...dcl.ApplyOption) *WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigManagedGroupConfig {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigManagedGroupConfig{}

	return cDes
}

func canonicalizeWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigManagedGroupConfigSlice(des, initial []WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigManagedGroupConfig, opts ...dcl.ApplyOption) []WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigManagedGroupConfig {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigManagedGroupConfig, 0, len(des))
		for _, d := range des {
			cd := canonicalizeWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigManagedGroupConfig(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigManagedGroupConfig, 0, len(des))
	for i, d := range des {
		cd := canonicalizeWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigManagedGroupConfig(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigManagedGroupConfig(c *Client, des, nw *WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigManagedGroupConfig) *WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigManagedGroupConfig {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigManagedGroupConfig while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringCanonicalize(des.InstanceTemplateName, nw.InstanceTemplateName) {
		nw.InstanceTemplateName = des.InstanceTemplateName
	}
	if dcl.StringCanonicalize(des.InstanceGroupManagerName, nw.InstanceGroupManagerName) {
		nw.InstanceGroupManagerName = des.InstanceGroupManagerName
	}

	return nw
}

func canonicalizeNewWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigManagedGroupConfigSet(c *Client, des, nw []WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigManagedGroupConfig) []WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigManagedGroupConfig {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigManagedGroupConfig
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigManagedGroupConfigNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigManagedGroupConfig(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigManagedGroupConfigSlice(c *Client, des, nw []WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigManagedGroupConfig) []WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigManagedGroupConfig {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigManagedGroupConfig
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigManagedGroupConfig(c, &d, &n))
	}

	return items
}

func canonicalizeWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigAccelerators(des, initial *WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigAccelerators, opts ...dcl.ApplyOption) *WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigAccelerators {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigAccelerators{}

	if dcl.StringCanonicalize(des.AcceleratorType, initial.AcceleratorType) || dcl.IsZeroValue(des.AcceleratorType) {
		cDes.AcceleratorType = initial.AcceleratorType
	} else {
		cDes.AcceleratorType = des.AcceleratorType
	}
	if dcl.IsZeroValue(des.AcceleratorCount) || (dcl.IsEmptyValueIndirect(des.AcceleratorCount) && dcl.IsEmptyValueIndirect(initial.AcceleratorCount)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.AcceleratorCount = initial.AcceleratorCount
	} else {
		cDes.AcceleratorCount = des.AcceleratorCount
	}

	return cDes
}

func canonicalizeWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigAcceleratorsSlice(des, initial []WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigAccelerators, opts ...dcl.ApplyOption) []WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigAccelerators {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigAccelerators, 0, len(des))
		for _, d := range des {
			cd := canonicalizeWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigAccelerators(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigAccelerators, 0, len(des))
	for i, d := range des {
		cd := canonicalizeWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigAccelerators(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigAccelerators(c *Client, des, nw *WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigAccelerators) *WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigAccelerators {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigAccelerators while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringCanonicalize(des.AcceleratorType, nw.AcceleratorType) {
		nw.AcceleratorType = des.AcceleratorType
	}

	return nw
}

func canonicalizeNewWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigAcceleratorsSet(c *Client, des, nw []WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigAccelerators) []WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigAccelerators {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigAccelerators
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigAcceleratorsNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigAccelerators(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigAcceleratorsSlice(c *Client, des, nw []WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigAccelerators) []WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigAccelerators {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigAccelerators
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigAccelerators(c, &d, &n))
	}

	return items
}

func canonicalizeWorkflowTemplatePlacementManagedClusterConfigSoftwareConfig(des, initial *WorkflowTemplatePlacementManagedClusterConfigSoftwareConfig, opts ...dcl.ApplyOption) *WorkflowTemplatePlacementManagedClusterConfigSoftwareConfig {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &WorkflowTemplatePlacementManagedClusterConfigSoftwareConfig{}

	if dcl.StringCanonicalize(des.ImageVersion, initial.ImageVersion) || dcl.IsZeroValue(des.ImageVersion) {
		cDes.ImageVersion = initial.ImageVersion
	} else {
		cDes.ImageVersion = des.ImageVersion
	}
	if dcl.IsZeroValue(des.Properties) || (dcl.IsEmptyValueIndirect(des.Properties) && dcl.IsEmptyValueIndirect(initial.Properties)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.Properties = initial.Properties
	} else {
		cDes.Properties = des.Properties
	}
	if dcl.IsZeroValue(des.OptionalComponents) || (dcl.IsEmptyValueIndirect(des.OptionalComponents) && dcl.IsEmptyValueIndirect(initial.OptionalComponents)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.OptionalComponents = initial.OptionalComponents
	} else {
		cDes.OptionalComponents = des.OptionalComponents
	}

	return cDes
}

func canonicalizeWorkflowTemplatePlacementManagedClusterConfigSoftwareConfigSlice(des, initial []WorkflowTemplatePlacementManagedClusterConfigSoftwareConfig, opts ...dcl.ApplyOption) []WorkflowTemplatePlacementManagedClusterConfigSoftwareConfig {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]WorkflowTemplatePlacementManagedClusterConfigSoftwareConfig, 0, len(des))
		for _, d := range des {
			cd := canonicalizeWorkflowTemplatePlacementManagedClusterConfigSoftwareConfig(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]WorkflowTemplatePlacementManagedClusterConfigSoftwareConfig, 0, len(des))
	for i, d := range des {
		cd := canonicalizeWorkflowTemplatePlacementManagedClusterConfigSoftwareConfig(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewWorkflowTemplatePlacementManagedClusterConfigSoftwareConfig(c *Client, des, nw *WorkflowTemplatePlacementManagedClusterConfigSoftwareConfig) *WorkflowTemplatePlacementManagedClusterConfigSoftwareConfig {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for WorkflowTemplatePlacementManagedClusterConfigSoftwareConfig while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringCanonicalize(des.ImageVersion, nw.ImageVersion) {
		nw.ImageVersion = des.ImageVersion
	}

	return nw
}

func canonicalizeNewWorkflowTemplatePlacementManagedClusterConfigSoftwareConfigSet(c *Client, des, nw []WorkflowTemplatePlacementManagedClusterConfigSoftwareConfig) []WorkflowTemplatePlacementManagedClusterConfigSoftwareConfig {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []WorkflowTemplatePlacementManagedClusterConfigSoftwareConfig
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareWorkflowTemplatePlacementManagedClusterConfigSoftwareConfigNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewWorkflowTemplatePlacementManagedClusterConfigSoftwareConfig(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewWorkflowTemplatePlacementManagedClusterConfigSoftwareConfigSlice(c *Client, des, nw []WorkflowTemplatePlacementManagedClusterConfigSoftwareConfig) []WorkflowTemplatePlacementManagedClusterConfigSoftwareConfig {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []WorkflowTemplatePlacementManagedClusterConfigSoftwareConfig
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewWorkflowTemplatePlacementManagedClusterConfigSoftwareConfig(c, &d, &n))
	}

	return items
}

func canonicalizeWorkflowTemplatePlacementManagedClusterConfigInitializationActions(des, initial *WorkflowTemplatePlacementManagedClusterConfigInitializationActions, opts ...dcl.ApplyOption) *WorkflowTemplatePlacementManagedClusterConfigInitializationActions {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &WorkflowTemplatePlacementManagedClusterConfigInitializationActions{}

	if dcl.StringCanonicalize(des.ExecutableFile, initial.ExecutableFile) || dcl.IsZeroValue(des.ExecutableFile) {
		cDes.ExecutableFile = initial.ExecutableFile
	} else {
		cDes.ExecutableFile = des.ExecutableFile
	}
	if dcl.StringCanonicalize(des.ExecutionTimeout, initial.ExecutionTimeout) || dcl.IsZeroValue(des.ExecutionTimeout) {
		cDes.ExecutionTimeout = initial.ExecutionTimeout
	} else {
		cDes.ExecutionTimeout = des.ExecutionTimeout
	}

	return cDes
}

func canonicalizeWorkflowTemplatePlacementManagedClusterConfigInitializationActionsSlice(des, initial []WorkflowTemplatePlacementManagedClusterConfigInitializationActions, opts ...dcl.ApplyOption) []WorkflowTemplatePlacementManagedClusterConfigInitializationActions {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]WorkflowTemplatePlacementManagedClusterConfigInitializationActions, 0, len(des))
		for _, d := range des {
			cd := canonicalizeWorkflowTemplatePlacementManagedClusterConfigInitializationActions(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]WorkflowTemplatePlacementManagedClusterConfigInitializationActions, 0, len(des))
	for i, d := range des {
		cd := canonicalizeWorkflowTemplatePlacementManagedClusterConfigInitializationActions(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewWorkflowTemplatePlacementManagedClusterConfigInitializationActions(c *Client, des, nw *WorkflowTemplatePlacementManagedClusterConfigInitializationActions) *WorkflowTemplatePlacementManagedClusterConfigInitializationActions {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for WorkflowTemplatePlacementManagedClusterConfigInitializationActions while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringCanonicalize(des.ExecutableFile, nw.ExecutableFile) {
		nw.ExecutableFile = des.ExecutableFile
	}
	if dcl.StringCanonicalize(des.ExecutionTimeout, nw.ExecutionTimeout) {
		nw.ExecutionTimeout = des.ExecutionTimeout
	}

	return nw
}

func canonicalizeNewWorkflowTemplatePlacementManagedClusterConfigInitializationActionsSet(c *Client, des, nw []WorkflowTemplatePlacementManagedClusterConfigInitializationActions) []WorkflowTemplatePlacementManagedClusterConfigInitializationActions {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []WorkflowTemplatePlacementManagedClusterConfigInitializationActions
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareWorkflowTemplatePlacementManagedClusterConfigInitializationActionsNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewWorkflowTemplatePlacementManagedClusterConfigInitializationActions(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewWorkflowTemplatePlacementManagedClusterConfigInitializationActionsSlice(c *Client, des, nw []WorkflowTemplatePlacementManagedClusterConfigInitializationActions) []WorkflowTemplatePlacementManagedClusterConfigInitializationActions {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []WorkflowTemplatePlacementManagedClusterConfigInitializationActions
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewWorkflowTemplatePlacementManagedClusterConfigInitializationActions(c, &d, &n))
	}

	return items
}

func canonicalizeWorkflowTemplatePlacementManagedClusterConfigEncryptionConfig(des, initial *WorkflowTemplatePlacementManagedClusterConfigEncryptionConfig, opts ...dcl.ApplyOption) *WorkflowTemplatePlacementManagedClusterConfigEncryptionConfig {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &WorkflowTemplatePlacementManagedClusterConfigEncryptionConfig{}

	if dcl.IsZeroValue(des.GcePdKmsKeyName) || (dcl.IsEmptyValueIndirect(des.GcePdKmsKeyName) && dcl.IsEmptyValueIndirect(initial.GcePdKmsKeyName)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.GcePdKmsKeyName = initial.GcePdKmsKeyName
	} else {
		cDes.GcePdKmsKeyName = des.GcePdKmsKeyName
	}

	return cDes
}

func canonicalizeWorkflowTemplatePlacementManagedClusterConfigEncryptionConfigSlice(des, initial []WorkflowTemplatePlacementManagedClusterConfigEncryptionConfig, opts ...dcl.ApplyOption) []WorkflowTemplatePlacementManagedClusterConfigEncryptionConfig {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]WorkflowTemplatePlacementManagedClusterConfigEncryptionConfig, 0, len(des))
		for _, d := range des {
			cd := canonicalizeWorkflowTemplatePlacementManagedClusterConfigEncryptionConfig(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]WorkflowTemplatePlacementManagedClusterConfigEncryptionConfig, 0, len(des))
	for i, d := range des {
		cd := canonicalizeWorkflowTemplatePlacementManagedClusterConfigEncryptionConfig(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewWorkflowTemplatePlacementManagedClusterConfigEncryptionConfig(c *Client, des, nw *WorkflowTemplatePlacementManagedClusterConfigEncryptionConfig) *WorkflowTemplatePlacementManagedClusterConfigEncryptionConfig {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for WorkflowTemplatePlacementManagedClusterConfigEncryptionConfig while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	return nw
}

func canonicalizeNewWorkflowTemplatePlacementManagedClusterConfigEncryptionConfigSet(c *Client, des, nw []WorkflowTemplatePlacementManagedClusterConfigEncryptionConfig) []WorkflowTemplatePlacementManagedClusterConfigEncryptionConfig {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []WorkflowTemplatePlacementManagedClusterConfigEncryptionConfig
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareWorkflowTemplatePlacementManagedClusterConfigEncryptionConfigNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewWorkflowTemplatePlacementManagedClusterConfigEncryptionConfig(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewWorkflowTemplatePlacementManagedClusterConfigEncryptionConfigSlice(c *Client, des, nw []WorkflowTemplatePlacementManagedClusterConfigEncryptionConfig) []WorkflowTemplatePlacementManagedClusterConfigEncryptionConfig {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []WorkflowTemplatePlacementManagedClusterConfigEncryptionConfig
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewWorkflowTemplatePlacementManagedClusterConfigEncryptionConfig(c, &d, &n))
	}

	return items
}

func canonicalizeWorkflowTemplatePlacementManagedClusterConfigAutoscalingConfig(des, initial *WorkflowTemplatePlacementManagedClusterConfigAutoscalingConfig, opts ...dcl.ApplyOption) *WorkflowTemplatePlacementManagedClusterConfigAutoscalingConfig {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &WorkflowTemplatePlacementManagedClusterConfigAutoscalingConfig{}

	if dcl.IsZeroValue(des.Policy) || (dcl.IsEmptyValueIndirect(des.Policy) && dcl.IsEmptyValueIndirect(initial.Policy)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.Policy = initial.Policy
	} else {
		cDes.Policy = des.Policy
	}

	return cDes
}

func canonicalizeWorkflowTemplatePlacementManagedClusterConfigAutoscalingConfigSlice(des, initial []WorkflowTemplatePlacementManagedClusterConfigAutoscalingConfig, opts ...dcl.ApplyOption) []WorkflowTemplatePlacementManagedClusterConfigAutoscalingConfig {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]WorkflowTemplatePlacementManagedClusterConfigAutoscalingConfig, 0, len(des))
		for _, d := range des {
			cd := canonicalizeWorkflowTemplatePlacementManagedClusterConfigAutoscalingConfig(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]WorkflowTemplatePlacementManagedClusterConfigAutoscalingConfig, 0, len(des))
	for i, d := range des {
		cd := canonicalizeWorkflowTemplatePlacementManagedClusterConfigAutoscalingConfig(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewWorkflowTemplatePlacementManagedClusterConfigAutoscalingConfig(c *Client, des, nw *WorkflowTemplatePlacementManagedClusterConfigAutoscalingConfig) *WorkflowTemplatePlacementManagedClusterConfigAutoscalingConfig {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for WorkflowTemplatePlacementManagedClusterConfigAutoscalingConfig while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	return nw
}

func canonicalizeNewWorkflowTemplatePlacementManagedClusterConfigAutoscalingConfigSet(c *Client, des, nw []WorkflowTemplatePlacementManagedClusterConfigAutoscalingConfig) []WorkflowTemplatePlacementManagedClusterConfigAutoscalingConfig {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []WorkflowTemplatePlacementManagedClusterConfigAutoscalingConfig
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareWorkflowTemplatePlacementManagedClusterConfigAutoscalingConfigNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewWorkflowTemplatePlacementManagedClusterConfigAutoscalingConfig(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewWorkflowTemplatePlacementManagedClusterConfigAutoscalingConfigSlice(c *Client, des, nw []WorkflowTemplatePlacementManagedClusterConfigAutoscalingConfig) []WorkflowTemplatePlacementManagedClusterConfigAutoscalingConfig {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []WorkflowTemplatePlacementManagedClusterConfigAutoscalingConfig
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewWorkflowTemplatePlacementManagedClusterConfigAutoscalingConfig(c, &d, &n))
	}

	return items
}

func canonicalizeWorkflowTemplatePlacementManagedClusterConfigSecurityConfig(des, initial *WorkflowTemplatePlacementManagedClusterConfigSecurityConfig, opts ...dcl.ApplyOption) *WorkflowTemplatePlacementManagedClusterConfigSecurityConfig {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &WorkflowTemplatePlacementManagedClusterConfigSecurityConfig{}

	cDes.KerberosConfig = canonicalizeWorkflowTemplatePlacementManagedClusterConfigSecurityConfigKerberosConfig(des.KerberosConfig, initial.KerberosConfig, opts...)

	return cDes
}

func canonicalizeWorkflowTemplatePlacementManagedClusterConfigSecurityConfigSlice(des, initial []WorkflowTemplatePlacementManagedClusterConfigSecurityConfig, opts ...dcl.ApplyOption) []WorkflowTemplatePlacementManagedClusterConfigSecurityConfig {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]WorkflowTemplatePlacementManagedClusterConfigSecurityConfig, 0, len(des))
		for _, d := range des {
			cd := canonicalizeWorkflowTemplatePlacementManagedClusterConfigSecurityConfig(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]WorkflowTemplatePlacementManagedClusterConfigSecurityConfig, 0, len(des))
	for i, d := range des {
		cd := canonicalizeWorkflowTemplatePlacementManagedClusterConfigSecurityConfig(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewWorkflowTemplatePlacementManagedClusterConfigSecurityConfig(c *Client, des, nw *WorkflowTemplatePlacementManagedClusterConfigSecurityConfig) *WorkflowTemplatePlacementManagedClusterConfigSecurityConfig {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for WorkflowTemplatePlacementManagedClusterConfigSecurityConfig while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	nw.KerberosConfig = canonicalizeNewWorkflowTemplatePlacementManagedClusterConfigSecurityConfigKerberosConfig(c, des.KerberosConfig, nw.KerberosConfig)

	return nw
}

func canonicalizeNewWorkflowTemplatePlacementManagedClusterConfigSecurityConfigSet(c *Client, des, nw []WorkflowTemplatePlacementManagedClusterConfigSecurityConfig) []WorkflowTemplatePlacementManagedClusterConfigSecurityConfig {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []WorkflowTemplatePlacementManagedClusterConfigSecurityConfig
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareWorkflowTemplatePlacementManagedClusterConfigSecurityConfigNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewWorkflowTemplatePlacementManagedClusterConfigSecurityConfig(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewWorkflowTemplatePlacementManagedClusterConfigSecurityConfigSlice(c *Client, des, nw []WorkflowTemplatePlacementManagedClusterConfigSecurityConfig) []WorkflowTemplatePlacementManagedClusterConfigSecurityConfig {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []WorkflowTemplatePlacementManagedClusterConfigSecurityConfig
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewWorkflowTemplatePlacementManagedClusterConfigSecurityConfig(c, &d, &n))
	}

	return items
}

func canonicalizeWorkflowTemplatePlacementManagedClusterConfigSecurityConfigKerberosConfig(des, initial *WorkflowTemplatePlacementManagedClusterConfigSecurityConfigKerberosConfig, opts ...dcl.ApplyOption) *WorkflowTemplatePlacementManagedClusterConfigSecurityConfigKerberosConfig {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &WorkflowTemplatePlacementManagedClusterConfigSecurityConfigKerberosConfig{}

	if dcl.BoolCanonicalize(des.EnableKerberos, initial.EnableKerberos) || dcl.IsZeroValue(des.EnableKerberos) {
		cDes.EnableKerberos = initial.EnableKerberos
	} else {
		cDes.EnableKerberos = des.EnableKerberos
	}
	if dcl.StringCanonicalize(des.RootPrincipalPassword, initial.RootPrincipalPassword) || dcl.IsZeroValue(des.RootPrincipalPassword) {
		cDes.RootPrincipalPassword = initial.RootPrincipalPassword
	} else {
		cDes.RootPrincipalPassword = des.RootPrincipalPassword
	}
	if dcl.IsZeroValue(des.KmsKey) || (dcl.IsEmptyValueIndirect(des.KmsKey) && dcl.IsEmptyValueIndirect(initial.KmsKey)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.KmsKey = initial.KmsKey
	} else {
		cDes.KmsKey = des.KmsKey
	}
	if dcl.StringCanonicalize(des.Keystore, initial.Keystore) || dcl.IsZeroValue(des.Keystore) {
		cDes.Keystore = initial.Keystore
	} else {
		cDes.Keystore = des.Keystore
	}
	if dcl.StringCanonicalize(des.Truststore, initial.Truststore) || dcl.IsZeroValue(des.Truststore) {
		cDes.Truststore = initial.Truststore
	} else {
		cDes.Truststore = des.Truststore
	}
	if dcl.StringCanonicalize(des.KeystorePassword, initial.KeystorePassword) || dcl.IsZeroValue(des.KeystorePassword) {
		cDes.KeystorePassword = initial.KeystorePassword
	} else {
		cDes.KeystorePassword = des.KeystorePassword
	}
	if dcl.StringCanonicalize(des.KeyPassword, initial.KeyPassword) || dcl.IsZeroValue(des.KeyPassword) {
		cDes.KeyPassword = initial.KeyPassword
	} else {
		cDes.KeyPassword = des.KeyPassword
	}
	if dcl.StringCanonicalize(des.TruststorePassword, initial.TruststorePassword) || dcl.IsZeroValue(des.TruststorePassword) {
		cDes.TruststorePassword = initial.TruststorePassword
	} else {
		cDes.TruststorePassword = des.TruststorePassword
	}
	if dcl.StringCanonicalize(des.CrossRealmTrustRealm, initial.CrossRealmTrustRealm) || dcl.IsZeroValue(des.CrossRealmTrustRealm) {
		cDes.CrossRealmTrustRealm = initial.CrossRealmTrustRealm
	} else {
		cDes.CrossRealmTrustRealm = des.CrossRealmTrustRealm
	}
	if dcl.StringCanonicalize(des.CrossRealmTrustKdc, initial.CrossRealmTrustKdc) || dcl.IsZeroValue(des.CrossRealmTrustKdc) {
		cDes.CrossRealmTrustKdc = initial.CrossRealmTrustKdc
	} else {
		cDes.CrossRealmTrustKdc = des.CrossRealmTrustKdc
	}
	if dcl.StringCanonicalize(des.CrossRealmTrustAdminServer, initial.CrossRealmTrustAdminServer) || dcl.IsZeroValue(des.CrossRealmTrustAdminServer) {
		cDes.CrossRealmTrustAdminServer = initial.CrossRealmTrustAdminServer
	} else {
		cDes.CrossRealmTrustAdminServer = des.CrossRealmTrustAdminServer
	}
	if dcl.StringCanonicalize(des.CrossRealmTrustSharedPassword, initial.CrossRealmTrustSharedPassword) || dcl.IsZeroValue(des.CrossRealmTrustSharedPassword) {
		cDes.CrossRealmTrustSharedPassword = initial.CrossRealmTrustSharedPassword
	} else {
		cDes.CrossRealmTrustSharedPassword = des.CrossRealmTrustSharedPassword
	}
	if dcl.StringCanonicalize(des.KdcDbKey, initial.KdcDbKey) || dcl.IsZeroValue(des.KdcDbKey) {
		cDes.KdcDbKey = initial.KdcDbKey
	} else {
		cDes.KdcDbKey = des.KdcDbKey
	}
	if dcl.IsZeroValue(des.TgtLifetimeHours) || (dcl.IsEmptyValueIndirect(des.TgtLifetimeHours) && dcl.IsEmptyValueIndirect(initial.TgtLifetimeHours)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.TgtLifetimeHours = initial.TgtLifetimeHours
	} else {
		cDes.TgtLifetimeHours = des.TgtLifetimeHours
	}
	if dcl.StringCanonicalize(des.Realm, initial.Realm) || dcl.IsZeroValue(des.Realm) {
		cDes.Realm = initial.Realm
	} else {
		cDes.Realm = des.Realm
	}

	return cDes
}

func canonicalizeWorkflowTemplatePlacementManagedClusterConfigSecurityConfigKerberosConfigSlice(des, initial []WorkflowTemplatePlacementManagedClusterConfigSecurityConfigKerberosConfig, opts ...dcl.ApplyOption) []WorkflowTemplatePlacementManagedClusterConfigSecurityConfigKerberosConfig {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]WorkflowTemplatePlacementManagedClusterConfigSecurityConfigKerberosConfig, 0, len(des))
		for _, d := range des {
			cd := canonicalizeWorkflowTemplatePlacementManagedClusterConfigSecurityConfigKerberosConfig(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]WorkflowTemplatePlacementManagedClusterConfigSecurityConfigKerberosConfig, 0, len(des))
	for i, d := range des {
		cd := canonicalizeWorkflowTemplatePlacementManagedClusterConfigSecurityConfigKerberosConfig(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewWorkflowTemplatePlacementManagedClusterConfigSecurityConfigKerberosConfig(c *Client, des, nw *WorkflowTemplatePlacementManagedClusterConfigSecurityConfigKerberosConfig) *WorkflowTemplatePlacementManagedClusterConfigSecurityConfigKerberosConfig {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for WorkflowTemplatePlacementManagedClusterConfigSecurityConfigKerberosConfig while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.BoolCanonicalize(des.EnableKerberos, nw.EnableKerberos) {
		nw.EnableKerberos = des.EnableKerberos
	}
	if dcl.StringCanonicalize(des.RootPrincipalPassword, nw.RootPrincipalPassword) {
		nw.RootPrincipalPassword = des.RootPrincipalPassword
	}
	if dcl.StringCanonicalize(des.Keystore, nw.Keystore) {
		nw.Keystore = des.Keystore
	}
	if dcl.StringCanonicalize(des.Truststore, nw.Truststore) {
		nw.Truststore = des.Truststore
	}
	if dcl.StringCanonicalize(des.KeystorePassword, nw.KeystorePassword) {
		nw.KeystorePassword = des.KeystorePassword
	}
	if dcl.StringCanonicalize(des.KeyPassword, nw.KeyPassword) {
		nw.KeyPassword = des.KeyPassword
	}
	if dcl.StringCanonicalize(des.TruststorePassword, nw.TruststorePassword) {
		nw.TruststorePassword = des.TruststorePassword
	}
	if dcl.StringCanonicalize(des.CrossRealmTrustRealm, nw.CrossRealmTrustRealm) {
		nw.CrossRealmTrustRealm = des.CrossRealmTrustRealm
	}
	if dcl.StringCanonicalize(des.CrossRealmTrustKdc, nw.CrossRealmTrustKdc) {
		nw.CrossRealmTrustKdc = des.CrossRealmTrustKdc
	}
	if dcl.StringCanonicalize(des.CrossRealmTrustAdminServer, nw.CrossRealmTrustAdminServer) {
		nw.CrossRealmTrustAdminServer = des.CrossRealmTrustAdminServer
	}
	if dcl.StringCanonicalize(des.CrossRealmTrustSharedPassword, nw.CrossRealmTrustSharedPassword) {
		nw.CrossRealmTrustSharedPassword = des.CrossRealmTrustSharedPassword
	}
	if dcl.StringCanonicalize(des.KdcDbKey, nw.KdcDbKey) {
		nw.KdcDbKey = des.KdcDbKey
	}
	if dcl.StringCanonicalize(des.Realm, nw.Realm) {
		nw.Realm = des.Realm
	}

	return nw
}

func canonicalizeNewWorkflowTemplatePlacementManagedClusterConfigSecurityConfigKerberosConfigSet(c *Client, des, nw []WorkflowTemplatePlacementManagedClusterConfigSecurityConfigKerberosConfig) []WorkflowTemplatePlacementManagedClusterConfigSecurityConfigKerberosConfig {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []WorkflowTemplatePlacementManagedClusterConfigSecurityConfigKerberosConfig
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareWorkflowTemplatePlacementManagedClusterConfigSecurityConfigKerberosConfigNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewWorkflowTemplatePlacementManagedClusterConfigSecurityConfigKerberosConfig(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewWorkflowTemplatePlacementManagedClusterConfigSecurityConfigKerberosConfigSlice(c *Client, des, nw []WorkflowTemplatePlacementManagedClusterConfigSecurityConfigKerberosConfig) []WorkflowTemplatePlacementManagedClusterConfigSecurityConfigKerberosConfig {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []WorkflowTemplatePlacementManagedClusterConfigSecurityConfigKerberosConfig
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewWorkflowTemplatePlacementManagedClusterConfigSecurityConfigKerberosConfig(c, &d, &n))
	}

	return items
}

func canonicalizeWorkflowTemplatePlacementManagedClusterConfigLifecycleConfig(des, initial *WorkflowTemplatePlacementManagedClusterConfigLifecycleConfig, opts ...dcl.ApplyOption) *WorkflowTemplatePlacementManagedClusterConfigLifecycleConfig {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &WorkflowTemplatePlacementManagedClusterConfigLifecycleConfig{}

	if dcl.StringCanonicalize(des.IdleDeleteTtl, initial.IdleDeleteTtl) || dcl.IsZeroValue(des.IdleDeleteTtl) {
		cDes.IdleDeleteTtl = initial.IdleDeleteTtl
	} else {
		cDes.IdleDeleteTtl = des.IdleDeleteTtl
	}
	if dcl.IsZeroValue(des.AutoDeleteTime) || (dcl.IsEmptyValueIndirect(des.AutoDeleteTime) && dcl.IsEmptyValueIndirect(initial.AutoDeleteTime)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.AutoDeleteTime = initial.AutoDeleteTime
	} else {
		cDes.AutoDeleteTime = des.AutoDeleteTime
	}
	if dcl.StringCanonicalize(des.AutoDeleteTtl, initial.AutoDeleteTtl) || dcl.IsZeroValue(des.AutoDeleteTtl) {
		cDes.AutoDeleteTtl = initial.AutoDeleteTtl
	} else {
		cDes.AutoDeleteTtl = des.AutoDeleteTtl
	}

	return cDes
}

func canonicalizeWorkflowTemplatePlacementManagedClusterConfigLifecycleConfigSlice(des, initial []WorkflowTemplatePlacementManagedClusterConfigLifecycleConfig, opts ...dcl.ApplyOption) []WorkflowTemplatePlacementManagedClusterConfigLifecycleConfig {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]WorkflowTemplatePlacementManagedClusterConfigLifecycleConfig, 0, len(des))
		for _, d := range des {
			cd := canonicalizeWorkflowTemplatePlacementManagedClusterConfigLifecycleConfig(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]WorkflowTemplatePlacementManagedClusterConfigLifecycleConfig, 0, len(des))
	for i, d := range des {
		cd := canonicalizeWorkflowTemplatePlacementManagedClusterConfigLifecycleConfig(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewWorkflowTemplatePlacementManagedClusterConfigLifecycleConfig(c *Client, des, nw *WorkflowTemplatePlacementManagedClusterConfigLifecycleConfig) *WorkflowTemplatePlacementManagedClusterConfigLifecycleConfig {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for WorkflowTemplatePlacementManagedClusterConfigLifecycleConfig while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringCanonicalize(des.IdleDeleteTtl, nw.IdleDeleteTtl) {
		nw.IdleDeleteTtl = des.IdleDeleteTtl
	}
	if dcl.StringCanonicalize(des.AutoDeleteTtl, nw.AutoDeleteTtl) {
		nw.AutoDeleteTtl = des.AutoDeleteTtl
	}

	return nw
}

func canonicalizeNewWorkflowTemplatePlacementManagedClusterConfigLifecycleConfigSet(c *Client, des, nw []WorkflowTemplatePlacementManagedClusterConfigLifecycleConfig) []WorkflowTemplatePlacementManagedClusterConfigLifecycleConfig {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []WorkflowTemplatePlacementManagedClusterConfigLifecycleConfig
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareWorkflowTemplatePlacementManagedClusterConfigLifecycleConfigNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewWorkflowTemplatePlacementManagedClusterConfigLifecycleConfig(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewWorkflowTemplatePlacementManagedClusterConfigLifecycleConfigSlice(c *Client, des, nw []WorkflowTemplatePlacementManagedClusterConfigLifecycleConfig) []WorkflowTemplatePlacementManagedClusterConfigLifecycleConfig {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []WorkflowTemplatePlacementManagedClusterConfigLifecycleConfig
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewWorkflowTemplatePlacementManagedClusterConfigLifecycleConfig(c, &d, &n))
	}

	return items
}

func canonicalizeWorkflowTemplatePlacementManagedClusterConfigEndpointConfig(des, initial *WorkflowTemplatePlacementManagedClusterConfigEndpointConfig, opts ...dcl.ApplyOption) *WorkflowTemplatePlacementManagedClusterConfigEndpointConfig {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &WorkflowTemplatePlacementManagedClusterConfigEndpointConfig{}

	if dcl.BoolCanonicalize(des.EnableHttpPortAccess, initial.EnableHttpPortAccess) || dcl.IsZeroValue(des.EnableHttpPortAccess) {
		cDes.EnableHttpPortAccess = initial.EnableHttpPortAccess
	} else {
		cDes.EnableHttpPortAccess = des.EnableHttpPortAccess
	}

	return cDes
}

func canonicalizeWorkflowTemplatePlacementManagedClusterConfigEndpointConfigSlice(des, initial []WorkflowTemplatePlacementManagedClusterConfigEndpointConfig, opts ...dcl.ApplyOption) []WorkflowTemplatePlacementManagedClusterConfigEndpointConfig {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]WorkflowTemplatePlacementManagedClusterConfigEndpointConfig, 0, len(des))
		for _, d := range des {
			cd := canonicalizeWorkflowTemplatePlacementManagedClusterConfigEndpointConfig(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]WorkflowTemplatePlacementManagedClusterConfigEndpointConfig, 0, len(des))
	for i, d := range des {
		cd := canonicalizeWorkflowTemplatePlacementManagedClusterConfigEndpointConfig(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewWorkflowTemplatePlacementManagedClusterConfigEndpointConfig(c *Client, des, nw *WorkflowTemplatePlacementManagedClusterConfigEndpointConfig) *WorkflowTemplatePlacementManagedClusterConfigEndpointConfig {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for WorkflowTemplatePlacementManagedClusterConfigEndpointConfig while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.BoolCanonicalize(des.EnableHttpPortAccess, nw.EnableHttpPortAccess) {
		nw.EnableHttpPortAccess = des.EnableHttpPortAccess
	}

	return nw
}

func canonicalizeNewWorkflowTemplatePlacementManagedClusterConfigEndpointConfigSet(c *Client, des, nw []WorkflowTemplatePlacementManagedClusterConfigEndpointConfig) []WorkflowTemplatePlacementManagedClusterConfigEndpointConfig {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []WorkflowTemplatePlacementManagedClusterConfigEndpointConfig
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareWorkflowTemplatePlacementManagedClusterConfigEndpointConfigNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewWorkflowTemplatePlacementManagedClusterConfigEndpointConfig(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewWorkflowTemplatePlacementManagedClusterConfigEndpointConfigSlice(c *Client, des, nw []WorkflowTemplatePlacementManagedClusterConfigEndpointConfig) []WorkflowTemplatePlacementManagedClusterConfigEndpointConfig {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []WorkflowTemplatePlacementManagedClusterConfigEndpointConfig
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewWorkflowTemplatePlacementManagedClusterConfigEndpointConfig(c, &d, &n))
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
	if dcl.IsZeroValue(des.ClusterLabels) || (dcl.IsEmptyValueIndirect(des.ClusterLabels) && dcl.IsEmptyValueIndirect(initial.ClusterLabels)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.ClusterLabels = initial.ClusterLabels
	} else {
		cDes.ClusterLabels = des.ClusterLabels
	}

	return cDes
}

func canonicalizeWorkflowTemplatePlacementClusterSelectorSlice(des, initial []WorkflowTemplatePlacementClusterSelector, opts ...dcl.ApplyOption) []WorkflowTemplatePlacementClusterSelector {
	if dcl.IsEmptyValueIndirect(des) {
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
		if dcl.IsEmptyValueIndirect(des) {
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

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []WorkflowTemplatePlacementClusterSelector
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareWorkflowTemplatePlacementClusterSelectorNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewWorkflowTemplatePlacementClusterSelector(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
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
	if dcl.IsZeroValue(des.Labels) || (dcl.IsEmptyValueIndirect(des.Labels) && dcl.IsEmptyValueIndirect(initial.Labels)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.Labels = initial.Labels
	} else {
		cDes.Labels = des.Labels
	}
	cDes.Scheduling = canonicalizeWorkflowTemplateJobsScheduling(des.Scheduling, initial.Scheduling, opts...)
	if dcl.StringArrayCanonicalize(des.PrerequisiteStepIds, initial.PrerequisiteStepIds) {
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
		if dcl.IsEmptyValueIndirect(des) {
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

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []WorkflowTemplateJobs
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareWorkflowTemplateJobsNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewWorkflowTemplateJobs(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
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
	if dcl.StringArrayCanonicalize(des.Args, initial.Args) {
		cDes.Args = initial.Args
	} else {
		cDes.Args = des.Args
	}
	if dcl.StringArrayCanonicalize(des.JarFileUris, initial.JarFileUris) {
		cDes.JarFileUris = initial.JarFileUris
	} else {
		cDes.JarFileUris = des.JarFileUris
	}
	if dcl.StringArrayCanonicalize(des.FileUris, initial.FileUris) {
		cDes.FileUris = initial.FileUris
	} else {
		cDes.FileUris = des.FileUris
	}
	if dcl.StringArrayCanonicalize(des.ArchiveUris, initial.ArchiveUris) {
		cDes.ArchiveUris = initial.ArchiveUris
	} else {
		cDes.ArchiveUris = des.ArchiveUris
	}
	if dcl.IsZeroValue(des.Properties) || (dcl.IsEmptyValueIndirect(des.Properties) && dcl.IsEmptyValueIndirect(initial.Properties)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.Properties = initial.Properties
	} else {
		cDes.Properties = des.Properties
	}
	cDes.LoggingConfig = canonicalizeWorkflowTemplateJobsHadoopJobLoggingConfig(des.LoggingConfig, initial.LoggingConfig, opts...)

	return cDes
}

func canonicalizeWorkflowTemplateJobsHadoopJobSlice(des, initial []WorkflowTemplateJobsHadoopJob, opts ...dcl.ApplyOption) []WorkflowTemplateJobsHadoopJob {
	if dcl.IsEmptyValueIndirect(des) {
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
		if dcl.IsEmptyValueIndirect(des) {
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

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []WorkflowTemplateJobsHadoopJob
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareWorkflowTemplateJobsHadoopJobNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewWorkflowTemplateJobsHadoopJob(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
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

	if dcl.IsZeroValue(des.DriverLogLevels) || (dcl.IsEmptyValueIndirect(des.DriverLogLevels) && dcl.IsEmptyValueIndirect(initial.DriverLogLevels)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.DriverLogLevels = initial.DriverLogLevels
	} else {
		cDes.DriverLogLevels = des.DriverLogLevels
	}

	return cDes
}

func canonicalizeWorkflowTemplateJobsHadoopJobLoggingConfigSlice(des, initial []WorkflowTemplateJobsHadoopJobLoggingConfig, opts ...dcl.ApplyOption) []WorkflowTemplateJobsHadoopJobLoggingConfig {
	if dcl.IsEmptyValueIndirect(des) {
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
		if dcl.IsEmptyValueIndirect(des) {
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

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []WorkflowTemplateJobsHadoopJobLoggingConfig
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareWorkflowTemplateJobsHadoopJobLoggingConfigNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewWorkflowTemplateJobsHadoopJobLoggingConfig(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
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
	if dcl.StringArrayCanonicalize(des.Args, initial.Args) {
		cDes.Args = initial.Args
	} else {
		cDes.Args = des.Args
	}
	if dcl.StringArrayCanonicalize(des.JarFileUris, initial.JarFileUris) {
		cDes.JarFileUris = initial.JarFileUris
	} else {
		cDes.JarFileUris = des.JarFileUris
	}
	if dcl.StringArrayCanonicalize(des.FileUris, initial.FileUris) {
		cDes.FileUris = initial.FileUris
	} else {
		cDes.FileUris = des.FileUris
	}
	if dcl.StringArrayCanonicalize(des.ArchiveUris, initial.ArchiveUris) {
		cDes.ArchiveUris = initial.ArchiveUris
	} else {
		cDes.ArchiveUris = des.ArchiveUris
	}
	if dcl.IsZeroValue(des.Properties) || (dcl.IsEmptyValueIndirect(des.Properties) && dcl.IsEmptyValueIndirect(initial.Properties)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.Properties = initial.Properties
	} else {
		cDes.Properties = des.Properties
	}
	cDes.LoggingConfig = canonicalizeWorkflowTemplateJobsSparkJobLoggingConfig(des.LoggingConfig, initial.LoggingConfig, opts...)

	return cDes
}

func canonicalizeWorkflowTemplateJobsSparkJobSlice(des, initial []WorkflowTemplateJobsSparkJob, opts ...dcl.ApplyOption) []WorkflowTemplateJobsSparkJob {
	if dcl.IsEmptyValueIndirect(des) {
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
		if dcl.IsEmptyValueIndirect(des) {
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

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []WorkflowTemplateJobsSparkJob
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareWorkflowTemplateJobsSparkJobNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewWorkflowTemplateJobsSparkJob(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
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

	if dcl.IsZeroValue(des.DriverLogLevels) || (dcl.IsEmptyValueIndirect(des.DriverLogLevels) && dcl.IsEmptyValueIndirect(initial.DriverLogLevels)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.DriverLogLevels = initial.DriverLogLevels
	} else {
		cDes.DriverLogLevels = des.DriverLogLevels
	}

	return cDes
}

func canonicalizeWorkflowTemplateJobsSparkJobLoggingConfigSlice(des, initial []WorkflowTemplateJobsSparkJobLoggingConfig, opts ...dcl.ApplyOption) []WorkflowTemplateJobsSparkJobLoggingConfig {
	if dcl.IsEmptyValueIndirect(des) {
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
		if dcl.IsEmptyValueIndirect(des) {
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

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []WorkflowTemplateJobsSparkJobLoggingConfig
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareWorkflowTemplateJobsSparkJobLoggingConfigNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewWorkflowTemplateJobsSparkJobLoggingConfig(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
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
	if dcl.StringArrayCanonicalize(des.Args, initial.Args) {
		cDes.Args = initial.Args
	} else {
		cDes.Args = des.Args
	}
	if dcl.StringArrayCanonicalize(des.PythonFileUris, initial.PythonFileUris) {
		cDes.PythonFileUris = initial.PythonFileUris
	} else {
		cDes.PythonFileUris = des.PythonFileUris
	}
	if dcl.StringArrayCanonicalize(des.JarFileUris, initial.JarFileUris) {
		cDes.JarFileUris = initial.JarFileUris
	} else {
		cDes.JarFileUris = des.JarFileUris
	}
	if dcl.StringArrayCanonicalize(des.FileUris, initial.FileUris) {
		cDes.FileUris = initial.FileUris
	} else {
		cDes.FileUris = des.FileUris
	}
	if dcl.StringArrayCanonicalize(des.ArchiveUris, initial.ArchiveUris) {
		cDes.ArchiveUris = initial.ArchiveUris
	} else {
		cDes.ArchiveUris = des.ArchiveUris
	}
	if dcl.IsZeroValue(des.Properties) || (dcl.IsEmptyValueIndirect(des.Properties) && dcl.IsEmptyValueIndirect(initial.Properties)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.Properties = initial.Properties
	} else {
		cDes.Properties = des.Properties
	}
	cDes.LoggingConfig = canonicalizeWorkflowTemplateJobsPysparkJobLoggingConfig(des.LoggingConfig, initial.LoggingConfig, opts...)

	return cDes
}

func canonicalizeWorkflowTemplateJobsPysparkJobSlice(des, initial []WorkflowTemplateJobsPysparkJob, opts ...dcl.ApplyOption) []WorkflowTemplateJobsPysparkJob {
	if dcl.IsEmptyValueIndirect(des) {
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
		if dcl.IsEmptyValueIndirect(des) {
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

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []WorkflowTemplateJobsPysparkJob
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareWorkflowTemplateJobsPysparkJobNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewWorkflowTemplateJobsPysparkJob(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
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

	if dcl.IsZeroValue(des.DriverLogLevels) || (dcl.IsEmptyValueIndirect(des.DriverLogLevels) && dcl.IsEmptyValueIndirect(initial.DriverLogLevels)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.DriverLogLevels = initial.DriverLogLevels
	} else {
		cDes.DriverLogLevels = des.DriverLogLevels
	}

	return cDes
}

func canonicalizeWorkflowTemplateJobsPysparkJobLoggingConfigSlice(des, initial []WorkflowTemplateJobsPysparkJobLoggingConfig, opts ...dcl.ApplyOption) []WorkflowTemplateJobsPysparkJobLoggingConfig {
	if dcl.IsEmptyValueIndirect(des) {
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
		if dcl.IsEmptyValueIndirect(des) {
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

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []WorkflowTemplateJobsPysparkJobLoggingConfig
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareWorkflowTemplateJobsPysparkJobLoggingConfigNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewWorkflowTemplateJobsPysparkJobLoggingConfig(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
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
	if dcl.IsZeroValue(des.ScriptVariables) || (dcl.IsEmptyValueIndirect(des.ScriptVariables) && dcl.IsEmptyValueIndirect(initial.ScriptVariables)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.ScriptVariables = initial.ScriptVariables
	} else {
		cDes.ScriptVariables = des.ScriptVariables
	}
	if dcl.IsZeroValue(des.Properties) || (dcl.IsEmptyValueIndirect(des.Properties) && dcl.IsEmptyValueIndirect(initial.Properties)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.Properties = initial.Properties
	} else {
		cDes.Properties = des.Properties
	}
	if dcl.StringArrayCanonicalize(des.JarFileUris, initial.JarFileUris) {
		cDes.JarFileUris = initial.JarFileUris
	} else {
		cDes.JarFileUris = des.JarFileUris
	}

	return cDes
}

func canonicalizeWorkflowTemplateJobsHiveJobSlice(des, initial []WorkflowTemplateJobsHiveJob, opts ...dcl.ApplyOption) []WorkflowTemplateJobsHiveJob {
	if dcl.IsEmptyValueIndirect(des) {
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
		if dcl.IsEmptyValueIndirect(des) {
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

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []WorkflowTemplateJobsHiveJob
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareWorkflowTemplateJobsHiveJobNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewWorkflowTemplateJobsHiveJob(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
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

	if dcl.StringArrayCanonicalize(des.Queries, initial.Queries) {
		cDes.Queries = initial.Queries
	} else {
		cDes.Queries = des.Queries
	}

	return cDes
}

func canonicalizeWorkflowTemplateJobsHiveJobQueryListSlice(des, initial []WorkflowTemplateJobsHiveJobQueryList, opts ...dcl.ApplyOption) []WorkflowTemplateJobsHiveJobQueryList {
	if dcl.IsEmptyValueIndirect(des) {
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
		if dcl.IsEmptyValueIndirect(des) {
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

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []WorkflowTemplateJobsHiveJobQueryList
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareWorkflowTemplateJobsHiveJobQueryListNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewWorkflowTemplateJobsHiveJobQueryList(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
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
	if dcl.IsZeroValue(des.ScriptVariables) || (dcl.IsEmptyValueIndirect(des.ScriptVariables) && dcl.IsEmptyValueIndirect(initial.ScriptVariables)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.ScriptVariables = initial.ScriptVariables
	} else {
		cDes.ScriptVariables = des.ScriptVariables
	}
	if dcl.IsZeroValue(des.Properties) || (dcl.IsEmptyValueIndirect(des.Properties) && dcl.IsEmptyValueIndirect(initial.Properties)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.Properties = initial.Properties
	} else {
		cDes.Properties = des.Properties
	}
	if dcl.StringArrayCanonicalize(des.JarFileUris, initial.JarFileUris) {
		cDes.JarFileUris = initial.JarFileUris
	} else {
		cDes.JarFileUris = des.JarFileUris
	}
	cDes.LoggingConfig = canonicalizeWorkflowTemplateJobsPigJobLoggingConfig(des.LoggingConfig, initial.LoggingConfig, opts...)

	return cDes
}

func canonicalizeWorkflowTemplateJobsPigJobSlice(des, initial []WorkflowTemplateJobsPigJob, opts ...dcl.ApplyOption) []WorkflowTemplateJobsPigJob {
	if dcl.IsEmptyValueIndirect(des) {
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
		if dcl.IsEmptyValueIndirect(des) {
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

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []WorkflowTemplateJobsPigJob
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareWorkflowTemplateJobsPigJobNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewWorkflowTemplateJobsPigJob(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
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

	if dcl.StringArrayCanonicalize(des.Queries, initial.Queries) {
		cDes.Queries = initial.Queries
	} else {
		cDes.Queries = des.Queries
	}

	return cDes
}

func canonicalizeWorkflowTemplateJobsPigJobQueryListSlice(des, initial []WorkflowTemplateJobsPigJobQueryList, opts ...dcl.ApplyOption) []WorkflowTemplateJobsPigJobQueryList {
	if dcl.IsEmptyValueIndirect(des) {
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
		if dcl.IsEmptyValueIndirect(des) {
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

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []WorkflowTemplateJobsPigJobQueryList
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareWorkflowTemplateJobsPigJobQueryListNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewWorkflowTemplateJobsPigJobQueryList(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
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

	if dcl.IsZeroValue(des.DriverLogLevels) || (dcl.IsEmptyValueIndirect(des.DriverLogLevels) && dcl.IsEmptyValueIndirect(initial.DriverLogLevels)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.DriverLogLevels = initial.DriverLogLevels
	} else {
		cDes.DriverLogLevels = des.DriverLogLevels
	}

	return cDes
}

func canonicalizeWorkflowTemplateJobsPigJobLoggingConfigSlice(des, initial []WorkflowTemplateJobsPigJobLoggingConfig, opts ...dcl.ApplyOption) []WorkflowTemplateJobsPigJobLoggingConfig {
	if dcl.IsEmptyValueIndirect(des) {
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
		if dcl.IsEmptyValueIndirect(des) {
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

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []WorkflowTemplateJobsPigJobLoggingConfig
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareWorkflowTemplateJobsPigJobLoggingConfigNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewWorkflowTemplateJobsPigJobLoggingConfig(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
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
	if dcl.StringArrayCanonicalize(des.Args, initial.Args) {
		cDes.Args = initial.Args
	} else {
		cDes.Args = des.Args
	}
	if dcl.StringArrayCanonicalize(des.FileUris, initial.FileUris) {
		cDes.FileUris = initial.FileUris
	} else {
		cDes.FileUris = des.FileUris
	}
	if dcl.StringArrayCanonicalize(des.ArchiveUris, initial.ArchiveUris) {
		cDes.ArchiveUris = initial.ArchiveUris
	} else {
		cDes.ArchiveUris = des.ArchiveUris
	}
	if dcl.IsZeroValue(des.Properties) || (dcl.IsEmptyValueIndirect(des.Properties) && dcl.IsEmptyValueIndirect(initial.Properties)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.Properties = initial.Properties
	} else {
		cDes.Properties = des.Properties
	}
	cDes.LoggingConfig = canonicalizeWorkflowTemplateJobsSparkRJobLoggingConfig(des.LoggingConfig, initial.LoggingConfig, opts...)

	return cDes
}

func canonicalizeWorkflowTemplateJobsSparkRJobSlice(des, initial []WorkflowTemplateJobsSparkRJob, opts ...dcl.ApplyOption) []WorkflowTemplateJobsSparkRJob {
	if dcl.IsEmptyValueIndirect(des) {
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
		if dcl.IsEmptyValueIndirect(des) {
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

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []WorkflowTemplateJobsSparkRJob
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareWorkflowTemplateJobsSparkRJobNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewWorkflowTemplateJobsSparkRJob(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
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

	if dcl.IsZeroValue(des.DriverLogLevels) || (dcl.IsEmptyValueIndirect(des.DriverLogLevels) && dcl.IsEmptyValueIndirect(initial.DriverLogLevels)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.DriverLogLevels = initial.DriverLogLevels
	} else {
		cDes.DriverLogLevels = des.DriverLogLevels
	}

	return cDes
}

func canonicalizeWorkflowTemplateJobsSparkRJobLoggingConfigSlice(des, initial []WorkflowTemplateJobsSparkRJobLoggingConfig, opts ...dcl.ApplyOption) []WorkflowTemplateJobsSparkRJobLoggingConfig {
	if dcl.IsEmptyValueIndirect(des) {
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
		if dcl.IsEmptyValueIndirect(des) {
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

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []WorkflowTemplateJobsSparkRJobLoggingConfig
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareWorkflowTemplateJobsSparkRJobLoggingConfigNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewWorkflowTemplateJobsSparkRJobLoggingConfig(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
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
	if dcl.IsZeroValue(des.ScriptVariables) || (dcl.IsEmptyValueIndirect(des.ScriptVariables) && dcl.IsEmptyValueIndirect(initial.ScriptVariables)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.ScriptVariables = initial.ScriptVariables
	} else {
		cDes.ScriptVariables = des.ScriptVariables
	}
	if dcl.IsZeroValue(des.Properties) || (dcl.IsEmptyValueIndirect(des.Properties) && dcl.IsEmptyValueIndirect(initial.Properties)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.Properties = initial.Properties
	} else {
		cDes.Properties = des.Properties
	}
	if dcl.StringArrayCanonicalize(des.JarFileUris, initial.JarFileUris) {
		cDes.JarFileUris = initial.JarFileUris
	} else {
		cDes.JarFileUris = des.JarFileUris
	}
	cDes.LoggingConfig = canonicalizeWorkflowTemplateJobsSparkSqlJobLoggingConfig(des.LoggingConfig, initial.LoggingConfig, opts...)

	return cDes
}

func canonicalizeWorkflowTemplateJobsSparkSqlJobSlice(des, initial []WorkflowTemplateJobsSparkSqlJob, opts ...dcl.ApplyOption) []WorkflowTemplateJobsSparkSqlJob {
	if dcl.IsEmptyValueIndirect(des) {
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
		if dcl.IsEmptyValueIndirect(des) {
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

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []WorkflowTemplateJobsSparkSqlJob
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareWorkflowTemplateJobsSparkSqlJobNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewWorkflowTemplateJobsSparkSqlJob(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
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

	if dcl.StringArrayCanonicalize(des.Queries, initial.Queries) {
		cDes.Queries = initial.Queries
	} else {
		cDes.Queries = des.Queries
	}

	return cDes
}

func canonicalizeWorkflowTemplateJobsSparkSqlJobQueryListSlice(des, initial []WorkflowTemplateJobsSparkSqlJobQueryList, opts ...dcl.ApplyOption) []WorkflowTemplateJobsSparkSqlJobQueryList {
	if dcl.IsEmptyValueIndirect(des) {
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
		if dcl.IsEmptyValueIndirect(des) {
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

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []WorkflowTemplateJobsSparkSqlJobQueryList
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareWorkflowTemplateJobsSparkSqlJobQueryListNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewWorkflowTemplateJobsSparkSqlJobQueryList(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
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

	if dcl.IsZeroValue(des.DriverLogLevels) || (dcl.IsEmptyValueIndirect(des.DriverLogLevels) && dcl.IsEmptyValueIndirect(initial.DriverLogLevels)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.DriverLogLevels = initial.DriverLogLevels
	} else {
		cDes.DriverLogLevels = des.DriverLogLevels
	}

	return cDes
}

func canonicalizeWorkflowTemplateJobsSparkSqlJobLoggingConfigSlice(des, initial []WorkflowTemplateJobsSparkSqlJobLoggingConfig, opts ...dcl.ApplyOption) []WorkflowTemplateJobsSparkSqlJobLoggingConfig {
	if dcl.IsEmptyValueIndirect(des) {
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
		if dcl.IsEmptyValueIndirect(des) {
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

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []WorkflowTemplateJobsSparkSqlJobLoggingConfig
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareWorkflowTemplateJobsSparkSqlJobLoggingConfigNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewWorkflowTemplateJobsSparkSqlJobLoggingConfig(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
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
	if dcl.StringArrayCanonicalize(des.ClientTags, initial.ClientTags) {
		cDes.ClientTags = initial.ClientTags
	} else {
		cDes.ClientTags = des.ClientTags
	}
	if dcl.IsZeroValue(des.Properties) || (dcl.IsEmptyValueIndirect(des.Properties) && dcl.IsEmptyValueIndirect(initial.Properties)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.Properties = initial.Properties
	} else {
		cDes.Properties = des.Properties
	}
	cDes.LoggingConfig = canonicalizeWorkflowTemplateJobsPrestoJobLoggingConfig(des.LoggingConfig, initial.LoggingConfig, opts...)

	return cDes
}

func canonicalizeWorkflowTemplateJobsPrestoJobSlice(des, initial []WorkflowTemplateJobsPrestoJob, opts ...dcl.ApplyOption) []WorkflowTemplateJobsPrestoJob {
	if dcl.IsEmptyValueIndirect(des) {
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
		if dcl.IsEmptyValueIndirect(des) {
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

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []WorkflowTemplateJobsPrestoJob
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareWorkflowTemplateJobsPrestoJobNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewWorkflowTemplateJobsPrestoJob(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
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

	if dcl.StringArrayCanonicalize(des.Queries, initial.Queries) {
		cDes.Queries = initial.Queries
	} else {
		cDes.Queries = des.Queries
	}

	return cDes
}

func canonicalizeWorkflowTemplateJobsPrestoJobQueryListSlice(des, initial []WorkflowTemplateJobsPrestoJobQueryList, opts ...dcl.ApplyOption) []WorkflowTemplateJobsPrestoJobQueryList {
	if dcl.IsEmptyValueIndirect(des) {
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
		if dcl.IsEmptyValueIndirect(des) {
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

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []WorkflowTemplateJobsPrestoJobQueryList
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareWorkflowTemplateJobsPrestoJobQueryListNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewWorkflowTemplateJobsPrestoJobQueryList(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
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

	if dcl.IsZeroValue(des.DriverLogLevels) || (dcl.IsEmptyValueIndirect(des.DriverLogLevels) && dcl.IsEmptyValueIndirect(initial.DriverLogLevels)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.DriverLogLevels = initial.DriverLogLevels
	} else {
		cDes.DriverLogLevels = des.DriverLogLevels
	}

	return cDes
}

func canonicalizeWorkflowTemplateJobsPrestoJobLoggingConfigSlice(des, initial []WorkflowTemplateJobsPrestoJobLoggingConfig, opts ...dcl.ApplyOption) []WorkflowTemplateJobsPrestoJobLoggingConfig {
	if dcl.IsEmptyValueIndirect(des) {
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
		if dcl.IsEmptyValueIndirect(des) {
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

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []WorkflowTemplateJobsPrestoJobLoggingConfig
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareWorkflowTemplateJobsPrestoJobLoggingConfigNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewWorkflowTemplateJobsPrestoJobLoggingConfig(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
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

	if dcl.IsZeroValue(des.MaxFailuresPerHour) || (dcl.IsEmptyValueIndirect(des.MaxFailuresPerHour) && dcl.IsEmptyValueIndirect(initial.MaxFailuresPerHour)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.MaxFailuresPerHour = initial.MaxFailuresPerHour
	} else {
		cDes.MaxFailuresPerHour = des.MaxFailuresPerHour
	}
	if dcl.IsZeroValue(des.MaxFailuresTotal) || (dcl.IsEmptyValueIndirect(des.MaxFailuresTotal) && dcl.IsEmptyValueIndirect(initial.MaxFailuresTotal)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.MaxFailuresTotal = initial.MaxFailuresTotal
	} else {
		cDes.MaxFailuresTotal = des.MaxFailuresTotal
	}

	return cDes
}

func canonicalizeWorkflowTemplateJobsSchedulingSlice(des, initial []WorkflowTemplateJobsScheduling, opts ...dcl.ApplyOption) []WorkflowTemplateJobsScheduling {
	if dcl.IsEmptyValueIndirect(des) {
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
		if dcl.IsEmptyValueIndirect(des) {
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

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []WorkflowTemplateJobsScheduling
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareWorkflowTemplateJobsSchedulingNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewWorkflowTemplateJobsScheduling(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
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
	if dcl.StringArrayCanonicalize(des.Fields, initial.Fields) {
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
		if dcl.IsEmptyValueIndirect(des) {
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

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []WorkflowTemplateParameters
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareWorkflowTemplateParametersNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewWorkflowTemplateParameters(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
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
	if dcl.IsEmptyValueIndirect(des) {
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
		if dcl.IsEmptyValueIndirect(des) {
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

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []WorkflowTemplateParametersValidation
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareWorkflowTemplateParametersValidationNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewWorkflowTemplateParametersValidation(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
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

	if dcl.StringArrayCanonicalize(des.Regexes, initial.Regexes) {
		cDes.Regexes = initial.Regexes
	} else {
		cDes.Regexes = des.Regexes
	}

	return cDes
}

func canonicalizeWorkflowTemplateParametersValidationRegexSlice(des, initial []WorkflowTemplateParametersValidationRegex, opts ...dcl.ApplyOption) []WorkflowTemplateParametersValidationRegex {
	if dcl.IsEmptyValueIndirect(des) {
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
		if dcl.IsEmptyValueIndirect(des) {
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

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []WorkflowTemplateParametersValidationRegex
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareWorkflowTemplateParametersValidationRegexNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewWorkflowTemplateParametersValidationRegex(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
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

	if dcl.StringArrayCanonicalize(des.Values, initial.Values) {
		cDes.Values = initial.Values
	} else {
		cDes.Values = des.Values
	}

	return cDes
}

func canonicalizeWorkflowTemplateParametersValidationValuesSlice(des, initial []WorkflowTemplateParametersValidationValues, opts ...dcl.ApplyOption) []WorkflowTemplateParametersValidationValues {
	if dcl.IsEmptyValueIndirect(des) {
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
		if dcl.IsEmptyValueIndirect(des) {
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

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []WorkflowTemplateParametersValidationValues
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareWorkflowTemplateParametersValidationValuesNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewWorkflowTemplateParametersValidationValues(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
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
	if ds, err := dcl.Diff(desired.Name, actual.Name, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Name")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Version, actual.Version, dcl.DiffInfo{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Version")); len(ds) != 0 || err != nil {
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

	if ds, err := dcl.Diff(desired.Labels, actual.Labels, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Labels")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Placement, actual.Placement, dcl.DiffInfo{ObjectFunction: compareWorkflowTemplatePlacementNewStyle, EmptyObject: EmptyWorkflowTemplatePlacement, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Placement")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Jobs, actual.Jobs, dcl.DiffInfo{ObjectFunction: compareWorkflowTemplateJobsNewStyle, EmptyObject: EmptyWorkflowTemplateJobs, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Jobs")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Parameters, actual.Parameters, dcl.DiffInfo{ObjectFunction: compareWorkflowTemplateParametersNewStyle, EmptyObject: EmptyWorkflowTemplateParameters, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Parameters")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.DagTimeout, actual.DagTimeout, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("DagTimeout")); len(ds) != 0 || err != nil {
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

	if ds, err := dcl.Diff(desired.ManagedCluster, actual.ManagedCluster, dcl.DiffInfo{ObjectFunction: compareWorkflowTemplatePlacementManagedClusterNewStyle, EmptyObject: EmptyWorkflowTemplatePlacementManagedCluster, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("ManagedCluster")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.ClusterSelector, actual.ClusterSelector, dcl.DiffInfo{ObjectFunction: compareWorkflowTemplatePlacementClusterSelectorNewStyle, EmptyObject: EmptyWorkflowTemplatePlacementClusterSelector, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("ClusterSelector")); len(ds) != 0 || err != nil {
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

	if ds, err := dcl.Diff(desired.ClusterName, actual.ClusterName, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("ClusterName")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Config, actual.Config, dcl.DiffInfo{ObjectFunction: compareWorkflowTemplatePlacementManagedClusterConfigNewStyle, EmptyObject: EmptyWorkflowTemplatePlacementManagedClusterConfig, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Config")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Labels, actual.Labels, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Labels")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareWorkflowTemplatePlacementManagedClusterConfigNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*WorkflowTemplatePlacementManagedClusterConfig)
	if !ok {
		desiredNotPointer, ok := d.(WorkflowTemplatePlacementManagedClusterConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a WorkflowTemplatePlacementManagedClusterConfig or *WorkflowTemplatePlacementManagedClusterConfig", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*WorkflowTemplatePlacementManagedClusterConfig)
	if !ok {
		actualNotPointer, ok := a.(WorkflowTemplatePlacementManagedClusterConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a WorkflowTemplatePlacementManagedClusterConfig", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.StagingBucket, actual.StagingBucket, dcl.DiffInfo{Type: "ReferenceType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("ConfigBucket")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.TempBucket, actual.TempBucket, dcl.DiffInfo{Type: "ReferenceType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("TempBucket")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.GceClusterConfig, actual.GceClusterConfig, dcl.DiffInfo{ObjectFunction: compareWorkflowTemplatePlacementManagedClusterConfigGceClusterConfigNewStyle, EmptyObject: EmptyWorkflowTemplatePlacementManagedClusterConfigGceClusterConfig, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("GceClusterConfig")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.MasterConfig, actual.MasterConfig, dcl.DiffInfo{ServerDefault: true, ObjectFunction: compareWorkflowTemplatePlacementManagedClusterConfigMasterConfigNewStyle, EmptyObject: EmptyWorkflowTemplatePlacementManagedClusterConfigMasterConfig, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("MasterConfig")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.WorkerConfig, actual.WorkerConfig, dcl.DiffInfo{ServerDefault: true, ObjectFunction: compareWorkflowTemplatePlacementManagedClusterConfigWorkerConfigNewStyle, EmptyObject: EmptyWorkflowTemplatePlacementManagedClusterConfigWorkerConfig, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("WorkerConfig")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.SecondaryWorkerConfig, actual.SecondaryWorkerConfig, dcl.DiffInfo{ServerDefault: true, ObjectFunction: compareWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigNewStyle, EmptyObject: EmptyWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfig, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("SecondaryWorkerConfig")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.SoftwareConfig, actual.SoftwareConfig, dcl.DiffInfo{ObjectFunction: compareWorkflowTemplatePlacementManagedClusterConfigSoftwareConfigNewStyle, EmptyObject: EmptyWorkflowTemplatePlacementManagedClusterConfigSoftwareConfig, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("SoftwareConfig")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.InitializationActions, actual.InitializationActions, dcl.DiffInfo{ObjectFunction: compareWorkflowTemplatePlacementManagedClusterConfigInitializationActionsNewStyle, EmptyObject: EmptyWorkflowTemplatePlacementManagedClusterConfigInitializationActions, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("InitializationActions")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.EncryptionConfig, actual.EncryptionConfig, dcl.DiffInfo{ObjectFunction: compareWorkflowTemplatePlacementManagedClusterConfigEncryptionConfigNewStyle, EmptyObject: EmptyWorkflowTemplatePlacementManagedClusterConfigEncryptionConfig, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("EncryptionConfig")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.AutoscalingConfig, actual.AutoscalingConfig, dcl.DiffInfo{ObjectFunction: compareWorkflowTemplatePlacementManagedClusterConfigAutoscalingConfigNewStyle, EmptyObject: EmptyWorkflowTemplatePlacementManagedClusterConfigAutoscalingConfig, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("AutoscalingConfig")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.SecurityConfig, actual.SecurityConfig, dcl.DiffInfo{ObjectFunction: compareWorkflowTemplatePlacementManagedClusterConfigSecurityConfigNewStyle, EmptyObject: EmptyWorkflowTemplatePlacementManagedClusterConfigSecurityConfig, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("SecurityConfig")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.LifecycleConfig, actual.LifecycleConfig, dcl.DiffInfo{ObjectFunction: compareWorkflowTemplatePlacementManagedClusterConfigLifecycleConfigNewStyle, EmptyObject: EmptyWorkflowTemplatePlacementManagedClusterConfigLifecycleConfig, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("LifecycleConfig")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.EndpointConfig, actual.EndpointConfig, dcl.DiffInfo{ObjectFunction: compareWorkflowTemplatePlacementManagedClusterConfigEndpointConfigNewStyle, EmptyObject: EmptyWorkflowTemplatePlacementManagedClusterConfigEndpointConfig, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("EndpointConfig")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareWorkflowTemplatePlacementManagedClusterConfigGceClusterConfigNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*WorkflowTemplatePlacementManagedClusterConfigGceClusterConfig)
	if !ok {
		desiredNotPointer, ok := d.(WorkflowTemplatePlacementManagedClusterConfigGceClusterConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a WorkflowTemplatePlacementManagedClusterConfigGceClusterConfig or *WorkflowTemplatePlacementManagedClusterConfigGceClusterConfig", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*WorkflowTemplatePlacementManagedClusterConfigGceClusterConfig)
	if !ok {
		actualNotPointer, ok := a.(WorkflowTemplatePlacementManagedClusterConfigGceClusterConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a WorkflowTemplatePlacementManagedClusterConfigGceClusterConfig", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.Zone, actual.Zone, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("ZoneUri")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Network, actual.Network, dcl.DiffInfo{Type: "ReferenceType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("NetworkUri")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Subnetwork, actual.Subnetwork, dcl.DiffInfo{Type: "ReferenceType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("SubnetworkUri")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.InternalIPOnly, actual.InternalIPOnly, dcl.DiffInfo{ServerDefault: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("InternalIpOnly")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.PrivateIPv6GoogleAccess, actual.PrivateIPv6GoogleAccess, dcl.DiffInfo{Type: "EnumType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("PrivateIpv6GoogleAccess")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.ServiceAccount, actual.ServiceAccount, dcl.DiffInfo{Type: "ReferenceType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("ServiceAccount")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.ServiceAccountScopes, actual.ServiceAccountScopes, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("ServiceAccountScopes")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Tags, actual.Tags, dcl.DiffInfo{Type: "Set", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Tags")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Metadata, actual.Metadata, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Metadata")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.ReservationAffinity, actual.ReservationAffinity, dcl.DiffInfo{ObjectFunction: compareWorkflowTemplatePlacementManagedClusterConfigGceClusterConfigReservationAffinityNewStyle, EmptyObject: EmptyWorkflowTemplatePlacementManagedClusterConfigGceClusterConfigReservationAffinity, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("ReservationAffinity")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.NodeGroupAffinity, actual.NodeGroupAffinity, dcl.DiffInfo{ObjectFunction: compareWorkflowTemplatePlacementManagedClusterConfigGceClusterConfigNodeGroupAffinityNewStyle, EmptyObject: EmptyWorkflowTemplatePlacementManagedClusterConfigGceClusterConfigNodeGroupAffinity, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("NodeGroupAffinity")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.ShieldedInstanceConfig, actual.ShieldedInstanceConfig, dcl.DiffInfo{ObjectFunction: compareWorkflowTemplatePlacementManagedClusterConfigGceClusterConfigShieldedInstanceConfigNewStyle, EmptyObject: EmptyWorkflowTemplatePlacementManagedClusterConfigGceClusterConfigShieldedInstanceConfig, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("ShieldedInstanceConfig")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareWorkflowTemplatePlacementManagedClusterConfigGceClusterConfigReservationAffinityNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigReservationAffinity)
	if !ok {
		desiredNotPointer, ok := d.(WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigReservationAffinity)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigReservationAffinity or *WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigReservationAffinity", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigReservationAffinity)
	if !ok {
		actualNotPointer, ok := a.(WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigReservationAffinity)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigReservationAffinity", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.ConsumeReservationType, actual.ConsumeReservationType, dcl.DiffInfo{Type: "EnumType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("ConsumeReservationType")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Key, actual.Key, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Key")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Values, actual.Values, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Values")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareWorkflowTemplatePlacementManagedClusterConfigGceClusterConfigNodeGroupAffinityNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigNodeGroupAffinity)
	if !ok {
		desiredNotPointer, ok := d.(WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigNodeGroupAffinity)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigNodeGroupAffinity or *WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigNodeGroupAffinity", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigNodeGroupAffinity)
	if !ok {
		actualNotPointer, ok := a.(WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigNodeGroupAffinity)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigNodeGroupAffinity", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.NodeGroup, actual.NodeGroup, dcl.DiffInfo{Type: "ReferenceType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("NodeGroupUri")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareWorkflowTemplatePlacementManagedClusterConfigGceClusterConfigShieldedInstanceConfigNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigShieldedInstanceConfig)
	if !ok {
		desiredNotPointer, ok := d.(WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigShieldedInstanceConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigShieldedInstanceConfig or *WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigShieldedInstanceConfig", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigShieldedInstanceConfig)
	if !ok {
		actualNotPointer, ok := a.(WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigShieldedInstanceConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigShieldedInstanceConfig", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.EnableSecureBoot, actual.EnableSecureBoot, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("EnableSecureBoot")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.EnableVtpm, actual.EnableVtpm, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("EnableVtpm")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.EnableIntegrityMonitoring, actual.EnableIntegrityMonitoring, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("EnableIntegrityMonitoring")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareWorkflowTemplatePlacementManagedClusterConfigMasterConfigNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*WorkflowTemplatePlacementManagedClusterConfigMasterConfig)
	if !ok {
		desiredNotPointer, ok := d.(WorkflowTemplatePlacementManagedClusterConfigMasterConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a WorkflowTemplatePlacementManagedClusterConfigMasterConfig or *WorkflowTemplatePlacementManagedClusterConfigMasterConfig", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*WorkflowTemplatePlacementManagedClusterConfigMasterConfig)
	if !ok {
		actualNotPointer, ok := a.(WorkflowTemplatePlacementManagedClusterConfigMasterConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a WorkflowTemplatePlacementManagedClusterConfigMasterConfig", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.NumInstances, actual.NumInstances, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("NumInstances")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.InstanceNames, actual.InstanceNames, dcl.DiffInfo{OutputOnly: true, ServerDefault: true, Type: "ReferenceType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("InstanceNames")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Image, actual.Image, dcl.DiffInfo{Type: "ReferenceType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("ImageUri")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.MachineType, actual.MachineType, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("MachineTypeUri")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.DiskConfig, actual.DiskConfig, dcl.DiffInfo{ServerDefault: true, ObjectFunction: compareWorkflowTemplatePlacementManagedClusterConfigMasterConfigDiskConfigNewStyle, EmptyObject: EmptyWorkflowTemplatePlacementManagedClusterConfigMasterConfigDiskConfig, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("DiskConfig")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.IsPreemptible, actual.IsPreemptible, dcl.DiffInfo{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("IsPreemptible")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Preemptibility, actual.Preemptibility, dcl.DiffInfo{Type: "EnumType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Preemptibility")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.ManagedGroupConfig, actual.ManagedGroupConfig, dcl.DiffInfo{OutputOnly: true, ServerDefault: true, ObjectFunction: compareWorkflowTemplatePlacementManagedClusterConfigMasterConfigManagedGroupConfigNewStyle, EmptyObject: EmptyWorkflowTemplatePlacementManagedClusterConfigMasterConfigManagedGroupConfig, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("ManagedGroupConfig")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Accelerators, actual.Accelerators, dcl.DiffInfo{ServerDefault: true, ObjectFunction: compareWorkflowTemplatePlacementManagedClusterConfigMasterConfigAcceleratorsNewStyle, EmptyObject: EmptyWorkflowTemplatePlacementManagedClusterConfigMasterConfigAccelerators, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Accelerators")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.MinCpuPlatform, actual.MinCpuPlatform, dcl.DiffInfo{ServerDefault: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("MinCpuPlatform")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareWorkflowTemplatePlacementManagedClusterConfigMasterConfigDiskConfigNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*WorkflowTemplatePlacementManagedClusterConfigMasterConfigDiskConfig)
	if !ok {
		desiredNotPointer, ok := d.(WorkflowTemplatePlacementManagedClusterConfigMasterConfigDiskConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a WorkflowTemplatePlacementManagedClusterConfigMasterConfigDiskConfig or *WorkflowTemplatePlacementManagedClusterConfigMasterConfigDiskConfig", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*WorkflowTemplatePlacementManagedClusterConfigMasterConfigDiskConfig)
	if !ok {
		actualNotPointer, ok := a.(WorkflowTemplatePlacementManagedClusterConfigMasterConfigDiskConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a WorkflowTemplatePlacementManagedClusterConfigMasterConfigDiskConfig", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.BootDiskType, actual.BootDiskType, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("BootDiskType")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.BootDiskSizeGb, actual.BootDiskSizeGb, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("BootDiskSizeGb")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.NumLocalSsds, actual.NumLocalSsds, dcl.DiffInfo{ServerDefault: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("NumLocalSsds")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareWorkflowTemplatePlacementManagedClusterConfigMasterConfigManagedGroupConfigNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*WorkflowTemplatePlacementManagedClusterConfigMasterConfigManagedGroupConfig)
	if !ok {
		desiredNotPointer, ok := d.(WorkflowTemplatePlacementManagedClusterConfigMasterConfigManagedGroupConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a WorkflowTemplatePlacementManagedClusterConfigMasterConfigManagedGroupConfig or *WorkflowTemplatePlacementManagedClusterConfigMasterConfigManagedGroupConfig", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*WorkflowTemplatePlacementManagedClusterConfigMasterConfigManagedGroupConfig)
	if !ok {
		actualNotPointer, ok := a.(WorkflowTemplatePlacementManagedClusterConfigMasterConfigManagedGroupConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a WorkflowTemplatePlacementManagedClusterConfigMasterConfigManagedGroupConfig", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.InstanceTemplateName, actual.InstanceTemplateName, dcl.DiffInfo{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("InstanceTemplateName")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.InstanceGroupManagerName, actual.InstanceGroupManagerName, dcl.DiffInfo{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("InstanceGroupManagerName")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareWorkflowTemplatePlacementManagedClusterConfigMasterConfigAcceleratorsNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*WorkflowTemplatePlacementManagedClusterConfigMasterConfigAccelerators)
	if !ok {
		desiredNotPointer, ok := d.(WorkflowTemplatePlacementManagedClusterConfigMasterConfigAccelerators)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a WorkflowTemplatePlacementManagedClusterConfigMasterConfigAccelerators or *WorkflowTemplatePlacementManagedClusterConfigMasterConfigAccelerators", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*WorkflowTemplatePlacementManagedClusterConfigMasterConfigAccelerators)
	if !ok {
		actualNotPointer, ok := a.(WorkflowTemplatePlacementManagedClusterConfigMasterConfigAccelerators)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a WorkflowTemplatePlacementManagedClusterConfigMasterConfigAccelerators", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.AcceleratorType, actual.AcceleratorType, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("AcceleratorTypeUri")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.AcceleratorCount, actual.AcceleratorCount, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("AcceleratorCount")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareWorkflowTemplatePlacementManagedClusterConfigWorkerConfigNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*WorkflowTemplatePlacementManagedClusterConfigWorkerConfig)
	if !ok {
		desiredNotPointer, ok := d.(WorkflowTemplatePlacementManagedClusterConfigWorkerConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a WorkflowTemplatePlacementManagedClusterConfigWorkerConfig or *WorkflowTemplatePlacementManagedClusterConfigWorkerConfig", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*WorkflowTemplatePlacementManagedClusterConfigWorkerConfig)
	if !ok {
		actualNotPointer, ok := a.(WorkflowTemplatePlacementManagedClusterConfigWorkerConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a WorkflowTemplatePlacementManagedClusterConfigWorkerConfig", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.NumInstances, actual.NumInstances, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("NumInstances")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.InstanceNames, actual.InstanceNames, dcl.DiffInfo{OutputOnly: true, ServerDefault: true, Type: "ReferenceType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("InstanceNames")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Image, actual.Image, dcl.DiffInfo{Type: "ReferenceType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("ImageUri")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.MachineType, actual.MachineType, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("MachineTypeUri")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.DiskConfig, actual.DiskConfig, dcl.DiffInfo{ServerDefault: true, ObjectFunction: compareWorkflowTemplatePlacementManagedClusterConfigWorkerConfigDiskConfigNewStyle, EmptyObject: EmptyWorkflowTemplatePlacementManagedClusterConfigWorkerConfigDiskConfig, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("DiskConfig")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.IsPreemptible, actual.IsPreemptible, dcl.DiffInfo{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("IsPreemptible")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Preemptibility, actual.Preemptibility, dcl.DiffInfo{Type: "EnumType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Preemptibility")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.ManagedGroupConfig, actual.ManagedGroupConfig, dcl.DiffInfo{OutputOnly: true, ServerDefault: true, ObjectFunction: compareWorkflowTemplatePlacementManagedClusterConfigWorkerConfigManagedGroupConfigNewStyle, EmptyObject: EmptyWorkflowTemplatePlacementManagedClusterConfigWorkerConfigManagedGroupConfig, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("ManagedGroupConfig")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Accelerators, actual.Accelerators, dcl.DiffInfo{ServerDefault: true, ObjectFunction: compareWorkflowTemplatePlacementManagedClusterConfigWorkerConfigAcceleratorsNewStyle, EmptyObject: EmptyWorkflowTemplatePlacementManagedClusterConfigWorkerConfigAccelerators, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Accelerators")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.MinCpuPlatform, actual.MinCpuPlatform, dcl.DiffInfo{ServerDefault: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("MinCpuPlatform")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareWorkflowTemplatePlacementManagedClusterConfigWorkerConfigDiskConfigNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*WorkflowTemplatePlacementManagedClusterConfigWorkerConfigDiskConfig)
	if !ok {
		desiredNotPointer, ok := d.(WorkflowTemplatePlacementManagedClusterConfigWorkerConfigDiskConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a WorkflowTemplatePlacementManagedClusterConfigWorkerConfigDiskConfig or *WorkflowTemplatePlacementManagedClusterConfigWorkerConfigDiskConfig", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*WorkflowTemplatePlacementManagedClusterConfigWorkerConfigDiskConfig)
	if !ok {
		actualNotPointer, ok := a.(WorkflowTemplatePlacementManagedClusterConfigWorkerConfigDiskConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a WorkflowTemplatePlacementManagedClusterConfigWorkerConfigDiskConfig", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.BootDiskType, actual.BootDiskType, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("BootDiskType")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.BootDiskSizeGb, actual.BootDiskSizeGb, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("BootDiskSizeGb")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.NumLocalSsds, actual.NumLocalSsds, dcl.DiffInfo{ServerDefault: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("NumLocalSsds")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareWorkflowTemplatePlacementManagedClusterConfigWorkerConfigManagedGroupConfigNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*WorkflowTemplatePlacementManagedClusterConfigWorkerConfigManagedGroupConfig)
	if !ok {
		desiredNotPointer, ok := d.(WorkflowTemplatePlacementManagedClusterConfigWorkerConfigManagedGroupConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a WorkflowTemplatePlacementManagedClusterConfigWorkerConfigManagedGroupConfig or *WorkflowTemplatePlacementManagedClusterConfigWorkerConfigManagedGroupConfig", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*WorkflowTemplatePlacementManagedClusterConfigWorkerConfigManagedGroupConfig)
	if !ok {
		actualNotPointer, ok := a.(WorkflowTemplatePlacementManagedClusterConfigWorkerConfigManagedGroupConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a WorkflowTemplatePlacementManagedClusterConfigWorkerConfigManagedGroupConfig", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.InstanceTemplateName, actual.InstanceTemplateName, dcl.DiffInfo{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("InstanceTemplateName")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.InstanceGroupManagerName, actual.InstanceGroupManagerName, dcl.DiffInfo{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("InstanceGroupManagerName")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareWorkflowTemplatePlacementManagedClusterConfigWorkerConfigAcceleratorsNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*WorkflowTemplatePlacementManagedClusterConfigWorkerConfigAccelerators)
	if !ok {
		desiredNotPointer, ok := d.(WorkflowTemplatePlacementManagedClusterConfigWorkerConfigAccelerators)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a WorkflowTemplatePlacementManagedClusterConfigWorkerConfigAccelerators or *WorkflowTemplatePlacementManagedClusterConfigWorkerConfigAccelerators", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*WorkflowTemplatePlacementManagedClusterConfigWorkerConfigAccelerators)
	if !ok {
		actualNotPointer, ok := a.(WorkflowTemplatePlacementManagedClusterConfigWorkerConfigAccelerators)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a WorkflowTemplatePlacementManagedClusterConfigWorkerConfigAccelerators", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.AcceleratorType, actual.AcceleratorType, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("AcceleratorTypeUri")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.AcceleratorCount, actual.AcceleratorCount, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("AcceleratorCount")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfig)
	if !ok {
		desiredNotPointer, ok := d.(WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfig or *WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfig", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfig)
	if !ok {
		actualNotPointer, ok := a.(WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfig", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.NumInstances, actual.NumInstances, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("NumInstances")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.InstanceNames, actual.InstanceNames, dcl.DiffInfo{OutputOnly: true, ServerDefault: true, Type: "ReferenceType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("InstanceNames")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Image, actual.Image, dcl.DiffInfo{Type: "ReferenceType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("ImageUri")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.MachineType, actual.MachineType, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("MachineTypeUri")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.DiskConfig, actual.DiskConfig, dcl.DiffInfo{ServerDefault: true, ObjectFunction: compareWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigDiskConfigNewStyle, EmptyObject: EmptyWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigDiskConfig, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("DiskConfig")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.IsPreemptible, actual.IsPreemptible, dcl.DiffInfo{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("IsPreemptible")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Preemptibility, actual.Preemptibility, dcl.DiffInfo{Type: "EnumType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Preemptibility")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.ManagedGroupConfig, actual.ManagedGroupConfig, dcl.DiffInfo{OutputOnly: true, ServerDefault: true, ObjectFunction: compareWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigManagedGroupConfigNewStyle, EmptyObject: EmptyWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigManagedGroupConfig, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("ManagedGroupConfig")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Accelerators, actual.Accelerators, dcl.DiffInfo{ServerDefault: true, ObjectFunction: compareWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigAcceleratorsNewStyle, EmptyObject: EmptyWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigAccelerators, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Accelerators")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.MinCpuPlatform, actual.MinCpuPlatform, dcl.DiffInfo{ServerDefault: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("MinCpuPlatform")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigDiskConfigNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigDiskConfig)
	if !ok {
		desiredNotPointer, ok := d.(WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigDiskConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigDiskConfig or *WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigDiskConfig", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigDiskConfig)
	if !ok {
		actualNotPointer, ok := a.(WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigDiskConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigDiskConfig", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.BootDiskType, actual.BootDiskType, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("BootDiskType")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.BootDiskSizeGb, actual.BootDiskSizeGb, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("BootDiskSizeGb")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.NumLocalSsds, actual.NumLocalSsds, dcl.DiffInfo{ServerDefault: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("NumLocalSsds")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigManagedGroupConfigNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigManagedGroupConfig)
	if !ok {
		desiredNotPointer, ok := d.(WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigManagedGroupConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigManagedGroupConfig or *WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigManagedGroupConfig", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigManagedGroupConfig)
	if !ok {
		actualNotPointer, ok := a.(WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigManagedGroupConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigManagedGroupConfig", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.InstanceTemplateName, actual.InstanceTemplateName, dcl.DiffInfo{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("InstanceTemplateName")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.InstanceGroupManagerName, actual.InstanceGroupManagerName, dcl.DiffInfo{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("InstanceGroupManagerName")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigAcceleratorsNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigAccelerators)
	if !ok {
		desiredNotPointer, ok := d.(WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigAccelerators)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigAccelerators or *WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigAccelerators", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigAccelerators)
	if !ok {
		actualNotPointer, ok := a.(WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigAccelerators)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigAccelerators", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.AcceleratorType, actual.AcceleratorType, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("AcceleratorTypeUri")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.AcceleratorCount, actual.AcceleratorCount, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("AcceleratorCount")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareWorkflowTemplatePlacementManagedClusterConfigSoftwareConfigNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*WorkflowTemplatePlacementManagedClusterConfigSoftwareConfig)
	if !ok {
		desiredNotPointer, ok := d.(WorkflowTemplatePlacementManagedClusterConfigSoftwareConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a WorkflowTemplatePlacementManagedClusterConfigSoftwareConfig or *WorkflowTemplatePlacementManagedClusterConfigSoftwareConfig", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*WorkflowTemplatePlacementManagedClusterConfigSoftwareConfig)
	if !ok {
		actualNotPointer, ok := a.(WorkflowTemplatePlacementManagedClusterConfigSoftwareConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a WorkflowTemplatePlacementManagedClusterConfigSoftwareConfig", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.ImageVersion, actual.ImageVersion, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("ImageVersion")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Properties, actual.Properties, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Properties")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.OptionalComponents, actual.OptionalComponents, dcl.DiffInfo{Type: "EnumType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("OptionalComponents")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareWorkflowTemplatePlacementManagedClusterConfigInitializationActionsNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*WorkflowTemplatePlacementManagedClusterConfigInitializationActions)
	if !ok {
		desiredNotPointer, ok := d.(WorkflowTemplatePlacementManagedClusterConfigInitializationActions)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a WorkflowTemplatePlacementManagedClusterConfigInitializationActions or *WorkflowTemplatePlacementManagedClusterConfigInitializationActions", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*WorkflowTemplatePlacementManagedClusterConfigInitializationActions)
	if !ok {
		actualNotPointer, ok := a.(WorkflowTemplatePlacementManagedClusterConfigInitializationActions)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a WorkflowTemplatePlacementManagedClusterConfigInitializationActions", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.ExecutableFile, actual.ExecutableFile, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("ExecutableFile")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.ExecutionTimeout, actual.ExecutionTimeout, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("ExecutionTimeout")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareWorkflowTemplatePlacementManagedClusterConfigEncryptionConfigNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*WorkflowTemplatePlacementManagedClusterConfigEncryptionConfig)
	if !ok {
		desiredNotPointer, ok := d.(WorkflowTemplatePlacementManagedClusterConfigEncryptionConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a WorkflowTemplatePlacementManagedClusterConfigEncryptionConfig or *WorkflowTemplatePlacementManagedClusterConfigEncryptionConfig", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*WorkflowTemplatePlacementManagedClusterConfigEncryptionConfig)
	if !ok {
		actualNotPointer, ok := a.(WorkflowTemplatePlacementManagedClusterConfigEncryptionConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a WorkflowTemplatePlacementManagedClusterConfigEncryptionConfig", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.GcePdKmsKeyName, actual.GcePdKmsKeyName, dcl.DiffInfo{Type: "ReferenceType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("GcePdKmsKeyName")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareWorkflowTemplatePlacementManagedClusterConfigAutoscalingConfigNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*WorkflowTemplatePlacementManagedClusterConfigAutoscalingConfig)
	if !ok {
		desiredNotPointer, ok := d.(WorkflowTemplatePlacementManagedClusterConfigAutoscalingConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a WorkflowTemplatePlacementManagedClusterConfigAutoscalingConfig or *WorkflowTemplatePlacementManagedClusterConfigAutoscalingConfig", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*WorkflowTemplatePlacementManagedClusterConfigAutoscalingConfig)
	if !ok {
		actualNotPointer, ok := a.(WorkflowTemplatePlacementManagedClusterConfigAutoscalingConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a WorkflowTemplatePlacementManagedClusterConfigAutoscalingConfig", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.Policy, actual.Policy, dcl.DiffInfo{Type: "ReferenceType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("PolicyUri")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareWorkflowTemplatePlacementManagedClusterConfigSecurityConfigNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*WorkflowTemplatePlacementManagedClusterConfigSecurityConfig)
	if !ok {
		desiredNotPointer, ok := d.(WorkflowTemplatePlacementManagedClusterConfigSecurityConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a WorkflowTemplatePlacementManagedClusterConfigSecurityConfig or *WorkflowTemplatePlacementManagedClusterConfigSecurityConfig", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*WorkflowTemplatePlacementManagedClusterConfigSecurityConfig)
	if !ok {
		actualNotPointer, ok := a.(WorkflowTemplatePlacementManagedClusterConfigSecurityConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a WorkflowTemplatePlacementManagedClusterConfigSecurityConfig", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.KerberosConfig, actual.KerberosConfig, dcl.DiffInfo{ObjectFunction: compareWorkflowTemplatePlacementManagedClusterConfigSecurityConfigKerberosConfigNewStyle, EmptyObject: EmptyWorkflowTemplatePlacementManagedClusterConfigSecurityConfigKerberosConfig, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("KerberosConfig")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareWorkflowTemplatePlacementManagedClusterConfigSecurityConfigKerberosConfigNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*WorkflowTemplatePlacementManagedClusterConfigSecurityConfigKerberosConfig)
	if !ok {
		desiredNotPointer, ok := d.(WorkflowTemplatePlacementManagedClusterConfigSecurityConfigKerberosConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a WorkflowTemplatePlacementManagedClusterConfigSecurityConfigKerberosConfig or *WorkflowTemplatePlacementManagedClusterConfigSecurityConfigKerberosConfig", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*WorkflowTemplatePlacementManagedClusterConfigSecurityConfigKerberosConfig)
	if !ok {
		actualNotPointer, ok := a.(WorkflowTemplatePlacementManagedClusterConfigSecurityConfigKerberosConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a WorkflowTemplatePlacementManagedClusterConfigSecurityConfigKerberosConfig", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.EnableKerberos, actual.EnableKerberos, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("EnableKerberos")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.RootPrincipalPassword, actual.RootPrincipalPassword, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("RootPrincipalPasswordUri")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.KmsKey, actual.KmsKey, dcl.DiffInfo{Type: "ReferenceType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("KmsKeyUri")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Keystore, actual.Keystore, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("KeystoreUri")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Truststore, actual.Truststore, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("TruststoreUri")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.KeystorePassword, actual.KeystorePassword, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("KeystorePasswordUri")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.KeyPassword, actual.KeyPassword, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("KeyPasswordUri")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.TruststorePassword, actual.TruststorePassword, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("TruststorePasswordUri")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.CrossRealmTrustRealm, actual.CrossRealmTrustRealm, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("CrossRealmTrustRealm")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.CrossRealmTrustKdc, actual.CrossRealmTrustKdc, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("CrossRealmTrustKdc")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.CrossRealmTrustAdminServer, actual.CrossRealmTrustAdminServer, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("CrossRealmTrustAdminServer")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.CrossRealmTrustSharedPassword, actual.CrossRealmTrustSharedPassword, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("CrossRealmTrustSharedPasswordUri")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.KdcDbKey, actual.KdcDbKey, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("KdcDbKeyUri")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.TgtLifetimeHours, actual.TgtLifetimeHours, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("TgtLifetimeHours")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Realm, actual.Realm, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Realm")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareWorkflowTemplatePlacementManagedClusterConfigLifecycleConfigNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*WorkflowTemplatePlacementManagedClusterConfigLifecycleConfig)
	if !ok {
		desiredNotPointer, ok := d.(WorkflowTemplatePlacementManagedClusterConfigLifecycleConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a WorkflowTemplatePlacementManagedClusterConfigLifecycleConfig or *WorkflowTemplatePlacementManagedClusterConfigLifecycleConfig", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*WorkflowTemplatePlacementManagedClusterConfigLifecycleConfig)
	if !ok {
		actualNotPointer, ok := a.(WorkflowTemplatePlacementManagedClusterConfigLifecycleConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a WorkflowTemplatePlacementManagedClusterConfigLifecycleConfig", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.IdleDeleteTtl, actual.IdleDeleteTtl, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("IdleDeleteTtl")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.AutoDeleteTime, actual.AutoDeleteTime, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("AutoDeleteTime")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.AutoDeleteTtl, actual.AutoDeleteTtl, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("AutoDeleteTtl")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.IdleStartTime, actual.IdleStartTime, dcl.DiffInfo{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("IdleStartTime")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareWorkflowTemplatePlacementManagedClusterConfigEndpointConfigNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*WorkflowTemplatePlacementManagedClusterConfigEndpointConfig)
	if !ok {
		desiredNotPointer, ok := d.(WorkflowTemplatePlacementManagedClusterConfigEndpointConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a WorkflowTemplatePlacementManagedClusterConfigEndpointConfig or *WorkflowTemplatePlacementManagedClusterConfigEndpointConfig", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*WorkflowTemplatePlacementManagedClusterConfigEndpointConfig)
	if !ok {
		actualNotPointer, ok := a.(WorkflowTemplatePlacementManagedClusterConfigEndpointConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a WorkflowTemplatePlacementManagedClusterConfigEndpointConfig", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.HttpPorts, actual.HttpPorts, dcl.DiffInfo{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("HttpPorts")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.EnableHttpPortAccess, actual.EnableHttpPortAccess, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("EnableHttpPortAccess")); len(ds) != 0 || err != nil {
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

	if ds, err := dcl.Diff(desired.Zone, actual.Zone, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Zone")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.ClusterLabels, actual.ClusterLabels, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("ClusterLabels")); len(ds) != 0 || err != nil {
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

	if ds, err := dcl.Diff(desired.StepId, actual.StepId, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("StepId")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.HadoopJob, actual.HadoopJob, dcl.DiffInfo{ObjectFunction: compareWorkflowTemplateJobsHadoopJobNewStyle, EmptyObject: EmptyWorkflowTemplateJobsHadoopJob, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("HadoopJob")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.SparkJob, actual.SparkJob, dcl.DiffInfo{ObjectFunction: compareWorkflowTemplateJobsSparkJobNewStyle, EmptyObject: EmptyWorkflowTemplateJobsSparkJob, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("SparkJob")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.PysparkJob, actual.PysparkJob, dcl.DiffInfo{ObjectFunction: compareWorkflowTemplateJobsPysparkJobNewStyle, EmptyObject: EmptyWorkflowTemplateJobsPysparkJob, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("PysparkJob")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.HiveJob, actual.HiveJob, dcl.DiffInfo{ObjectFunction: compareWorkflowTemplateJobsHiveJobNewStyle, EmptyObject: EmptyWorkflowTemplateJobsHiveJob, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("HiveJob")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.PigJob, actual.PigJob, dcl.DiffInfo{ObjectFunction: compareWorkflowTemplateJobsPigJobNewStyle, EmptyObject: EmptyWorkflowTemplateJobsPigJob, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("PigJob")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.SparkRJob, actual.SparkRJob, dcl.DiffInfo{ObjectFunction: compareWorkflowTemplateJobsSparkRJobNewStyle, EmptyObject: EmptyWorkflowTemplateJobsSparkRJob, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("SparkRJob")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.SparkSqlJob, actual.SparkSqlJob, dcl.DiffInfo{ObjectFunction: compareWorkflowTemplateJobsSparkSqlJobNewStyle, EmptyObject: EmptyWorkflowTemplateJobsSparkSqlJob, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("SparkSqlJob")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.PrestoJob, actual.PrestoJob, dcl.DiffInfo{ObjectFunction: compareWorkflowTemplateJobsPrestoJobNewStyle, EmptyObject: EmptyWorkflowTemplateJobsPrestoJob, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("PrestoJob")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Labels, actual.Labels, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Labels")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Scheduling, actual.Scheduling, dcl.DiffInfo{ObjectFunction: compareWorkflowTemplateJobsSchedulingNewStyle, EmptyObject: EmptyWorkflowTemplateJobsScheduling, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Scheduling")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.PrerequisiteStepIds, actual.PrerequisiteStepIds, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("PrerequisiteStepIds")); len(ds) != 0 || err != nil {
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

	if ds, err := dcl.Diff(desired.MainJarFileUri, actual.MainJarFileUri, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("MainJarFileUri")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.MainClass, actual.MainClass, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("MainClass")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Args, actual.Args, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Args")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.JarFileUris, actual.JarFileUris, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("JarFileUris")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.FileUris, actual.FileUris, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("FileUris")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.ArchiveUris, actual.ArchiveUris, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("ArchiveUris")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Properties, actual.Properties, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Properties")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.LoggingConfig, actual.LoggingConfig, dcl.DiffInfo{ObjectFunction: compareWorkflowTemplateJobsHadoopJobLoggingConfigNewStyle, EmptyObject: EmptyWorkflowTemplateJobsHadoopJobLoggingConfig, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("LoggingConfig")); len(ds) != 0 || err != nil {
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

	if ds, err := dcl.Diff(desired.DriverLogLevels, actual.DriverLogLevels, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("DriverLogLevels")); len(ds) != 0 || err != nil {
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

	if ds, err := dcl.Diff(desired.MainJarFileUri, actual.MainJarFileUri, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("MainJarFileUri")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.MainClass, actual.MainClass, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("MainClass")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Args, actual.Args, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Args")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.JarFileUris, actual.JarFileUris, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("JarFileUris")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.FileUris, actual.FileUris, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("FileUris")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.ArchiveUris, actual.ArchiveUris, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("ArchiveUris")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Properties, actual.Properties, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Properties")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.LoggingConfig, actual.LoggingConfig, dcl.DiffInfo{ObjectFunction: compareWorkflowTemplateJobsSparkJobLoggingConfigNewStyle, EmptyObject: EmptyWorkflowTemplateJobsSparkJobLoggingConfig, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("LoggingConfig")); len(ds) != 0 || err != nil {
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

	if ds, err := dcl.Diff(desired.DriverLogLevels, actual.DriverLogLevels, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("DriverLogLevels")); len(ds) != 0 || err != nil {
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

	if ds, err := dcl.Diff(desired.MainPythonFileUri, actual.MainPythonFileUri, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("MainPythonFileUri")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Args, actual.Args, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Args")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.PythonFileUris, actual.PythonFileUris, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("PythonFileUris")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.JarFileUris, actual.JarFileUris, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("JarFileUris")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.FileUris, actual.FileUris, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("FileUris")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.ArchiveUris, actual.ArchiveUris, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("ArchiveUris")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Properties, actual.Properties, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Properties")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.LoggingConfig, actual.LoggingConfig, dcl.DiffInfo{ObjectFunction: compareWorkflowTemplateJobsPysparkJobLoggingConfigNewStyle, EmptyObject: EmptyWorkflowTemplateJobsPysparkJobLoggingConfig, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("LoggingConfig")); len(ds) != 0 || err != nil {
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

	if ds, err := dcl.Diff(desired.DriverLogLevels, actual.DriverLogLevels, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("DriverLogLevels")); len(ds) != 0 || err != nil {
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

	if ds, err := dcl.Diff(desired.QueryFileUri, actual.QueryFileUri, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("QueryFileUri")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.QueryList, actual.QueryList, dcl.DiffInfo{ObjectFunction: compareWorkflowTemplateJobsHiveJobQueryListNewStyle, EmptyObject: EmptyWorkflowTemplateJobsHiveJobQueryList, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("QueryList")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.ContinueOnFailure, actual.ContinueOnFailure, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("ContinueOnFailure")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.ScriptVariables, actual.ScriptVariables, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("ScriptVariables")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Properties, actual.Properties, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Properties")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.JarFileUris, actual.JarFileUris, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("JarFileUris")); len(ds) != 0 || err != nil {
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

	if ds, err := dcl.Diff(desired.Queries, actual.Queries, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Queries")); len(ds) != 0 || err != nil {
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

	if ds, err := dcl.Diff(desired.QueryFileUri, actual.QueryFileUri, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("QueryFileUri")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.QueryList, actual.QueryList, dcl.DiffInfo{ObjectFunction: compareWorkflowTemplateJobsPigJobQueryListNewStyle, EmptyObject: EmptyWorkflowTemplateJobsPigJobQueryList, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("QueryList")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.ContinueOnFailure, actual.ContinueOnFailure, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("ContinueOnFailure")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.ScriptVariables, actual.ScriptVariables, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("ScriptVariables")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Properties, actual.Properties, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Properties")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.JarFileUris, actual.JarFileUris, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("JarFileUris")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.LoggingConfig, actual.LoggingConfig, dcl.DiffInfo{ObjectFunction: compareWorkflowTemplateJobsPigJobLoggingConfigNewStyle, EmptyObject: EmptyWorkflowTemplateJobsPigJobLoggingConfig, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("LoggingConfig")); len(ds) != 0 || err != nil {
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

	if ds, err := dcl.Diff(desired.Queries, actual.Queries, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Queries")); len(ds) != 0 || err != nil {
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

	if ds, err := dcl.Diff(desired.DriverLogLevels, actual.DriverLogLevels, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("DriverLogLevels")); len(ds) != 0 || err != nil {
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

	if ds, err := dcl.Diff(desired.MainRFileUri, actual.MainRFileUri, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("MainRFileUri")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Args, actual.Args, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Args")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.FileUris, actual.FileUris, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("FileUris")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.ArchiveUris, actual.ArchiveUris, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("ArchiveUris")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Properties, actual.Properties, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Properties")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.LoggingConfig, actual.LoggingConfig, dcl.DiffInfo{ObjectFunction: compareWorkflowTemplateJobsSparkRJobLoggingConfigNewStyle, EmptyObject: EmptyWorkflowTemplateJobsSparkRJobLoggingConfig, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("LoggingConfig")); len(ds) != 0 || err != nil {
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

	if ds, err := dcl.Diff(desired.DriverLogLevels, actual.DriverLogLevels, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("DriverLogLevels")); len(ds) != 0 || err != nil {
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

	if ds, err := dcl.Diff(desired.QueryFileUri, actual.QueryFileUri, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("QueryFileUri")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.QueryList, actual.QueryList, dcl.DiffInfo{ObjectFunction: compareWorkflowTemplateJobsSparkSqlJobQueryListNewStyle, EmptyObject: EmptyWorkflowTemplateJobsSparkSqlJobQueryList, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("QueryList")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.ScriptVariables, actual.ScriptVariables, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("ScriptVariables")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Properties, actual.Properties, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Properties")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.JarFileUris, actual.JarFileUris, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("JarFileUris")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.LoggingConfig, actual.LoggingConfig, dcl.DiffInfo{ObjectFunction: compareWorkflowTemplateJobsSparkSqlJobLoggingConfigNewStyle, EmptyObject: EmptyWorkflowTemplateJobsSparkSqlJobLoggingConfig, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("LoggingConfig")); len(ds) != 0 || err != nil {
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

	if ds, err := dcl.Diff(desired.Queries, actual.Queries, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Queries")); len(ds) != 0 || err != nil {
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

	if ds, err := dcl.Diff(desired.DriverLogLevels, actual.DriverLogLevels, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("DriverLogLevels")); len(ds) != 0 || err != nil {
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

	if ds, err := dcl.Diff(desired.QueryFileUri, actual.QueryFileUri, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("QueryFileUri")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.QueryList, actual.QueryList, dcl.DiffInfo{ObjectFunction: compareWorkflowTemplateJobsPrestoJobQueryListNewStyle, EmptyObject: EmptyWorkflowTemplateJobsPrestoJobQueryList, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("QueryList")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.ContinueOnFailure, actual.ContinueOnFailure, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("ContinueOnFailure")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.OutputFormat, actual.OutputFormat, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("OutputFormat")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.ClientTags, actual.ClientTags, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("ClientTags")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Properties, actual.Properties, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Properties")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.LoggingConfig, actual.LoggingConfig, dcl.DiffInfo{ObjectFunction: compareWorkflowTemplateJobsPrestoJobLoggingConfigNewStyle, EmptyObject: EmptyWorkflowTemplateJobsPrestoJobLoggingConfig, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("LoggingConfig")); len(ds) != 0 || err != nil {
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

	if ds, err := dcl.Diff(desired.Queries, actual.Queries, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Queries")); len(ds) != 0 || err != nil {
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

	if ds, err := dcl.Diff(desired.DriverLogLevels, actual.DriverLogLevels, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("DriverLogLevels")); len(ds) != 0 || err != nil {
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

	if ds, err := dcl.Diff(desired.MaxFailuresPerHour, actual.MaxFailuresPerHour, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("MaxFailuresPerHour")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.MaxFailuresTotal, actual.MaxFailuresTotal, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("MaxFailuresTotal")); len(ds) != 0 || err != nil {
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

	if ds, err := dcl.Diff(desired.Name, actual.Name, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Name")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Fields, actual.Fields, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Fields")); len(ds) != 0 || err != nil {
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

	if ds, err := dcl.Diff(desired.Validation, actual.Validation, dcl.DiffInfo{ObjectFunction: compareWorkflowTemplateParametersValidationNewStyle, EmptyObject: EmptyWorkflowTemplateParametersValidation, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Validation")); len(ds) != 0 || err != nil {
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

	if ds, err := dcl.Diff(desired.Regex, actual.Regex, dcl.DiffInfo{ObjectFunction: compareWorkflowTemplateParametersValidationRegexNewStyle, EmptyObject: EmptyWorkflowTemplateParametersValidationRegex, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Regex")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Values, actual.Values, dcl.DiffInfo{ObjectFunction: compareWorkflowTemplateParametersValidationValuesNewStyle, EmptyObject: EmptyWorkflowTemplateParametersValidationValues, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Values")); len(ds) != 0 || err != nil {
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

	if ds, err := dcl.Diff(desired.Regexes, actual.Regexes, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Regexes")); len(ds) != 0 || err != nil {
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

	if ds, err := dcl.Diff(desired.Values, actual.Values, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Values")); len(ds) != 0 || err != nil {
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
func unmarshalWorkflowTemplate(b []byte, c *Client, res *WorkflowTemplate) (*WorkflowTemplate, error) {
	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, err
	}
	return unmarshalMapWorkflowTemplate(m, c, res)
}

func unmarshalMapWorkflowTemplate(m map[string]interface{}, c *Client, res *WorkflowTemplate) (*WorkflowTemplate, error) {

	flattened := flattenWorkflowTemplate(c, m, res)
	if flattened == nil {
		return nil, fmt.Errorf("attempted to flatten empty json object")
	}
	return flattened, nil
}

// expandWorkflowTemplate expands WorkflowTemplate into a JSON request object.
func expandWorkflowTemplate(c *Client, f *WorkflowTemplate) (map[string]interface{}, error) {
	m := make(map[string]interface{})
	res := f
	_ = res
	if v, err := dcl.EmptyValue(); err != nil {
		return nil, fmt.Errorf("error expanding Name into name: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["name"] = v
	}
	if v := f.Labels; dcl.ValueShouldBeSent(v) {
		m["labels"] = v
	}
	if v, err := expandWorkflowTemplatePlacement(c, f.Placement, res); err != nil {
		return nil, fmt.Errorf("error expanding Placement into placement: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["placement"] = v
	}
	if v, err := expandWorkflowTemplateJobsSlice(c, f.Jobs, res); err != nil {
		return nil, fmt.Errorf("error expanding Jobs into jobs: %w", err)
	} else if v != nil {
		m["jobs"] = v
	}
	if v, err := expandWorkflowTemplateParametersSlice(c, f.Parameters, res); err != nil {
		return nil, fmt.Errorf("error expanding Parameters into parameters: %w", err)
	} else if v != nil {
		m["parameters"] = v
	}
	if v := f.DagTimeout; dcl.ValueShouldBeSent(v) {
		m["dagTimeout"] = v
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

// flattenWorkflowTemplate flattens WorkflowTemplate from a JSON request object into the
// WorkflowTemplate type.
func flattenWorkflowTemplate(c *Client, i interface{}, res *WorkflowTemplate) *WorkflowTemplate {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}
	if len(m) == 0 {
		return nil
	}

	resultRes := &WorkflowTemplate{}
	resultRes.Name = dcl.FlattenString(m["name"])
	resultRes.Version = dcl.FlattenInteger(m["version"])
	resultRes.CreateTime = dcl.FlattenString(m["createTime"])
	resultRes.UpdateTime = dcl.FlattenString(m["updateTime"])
	resultRes.Labels = dcl.FlattenKeyValuePairs(m["labels"])
	resultRes.Placement = flattenWorkflowTemplatePlacement(c, m["placement"], res)
	resultRes.Jobs = flattenWorkflowTemplateJobsSlice(c, m["jobs"], res)
	resultRes.Parameters = flattenWorkflowTemplateParametersSlice(c, m["parameters"], res)
	resultRes.DagTimeout = dcl.FlattenString(m["dagTimeout"])
	resultRes.Project = dcl.FlattenString(m["project"])
	resultRes.Location = dcl.FlattenString(m["location"])

	return resultRes
}

// expandWorkflowTemplatePlacementMap expands the contents of WorkflowTemplatePlacement into a JSON
// request object.
func expandWorkflowTemplatePlacementMap(c *Client, f map[string]WorkflowTemplatePlacement, res *WorkflowTemplate) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandWorkflowTemplatePlacement(c, &item, res)
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
func expandWorkflowTemplatePlacementSlice(c *Client, f []WorkflowTemplatePlacement, res *WorkflowTemplate) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandWorkflowTemplatePlacement(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenWorkflowTemplatePlacementMap flattens the contents of WorkflowTemplatePlacement from a JSON
// response object.
func flattenWorkflowTemplatePlacementMap(c *Client, i interface{}, res *WorkflowTemplate) map[string]WorkflowTemplatePlacement {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]WorkflowTemplatePlacement{}
	}

	if len(a) == 0 {
		return map[string]WorkflowTemplatePlacement{}
	}

	items := make(map[string]WorkflowTemplatePlacement)
	for k, item := range a {
		items[k] = *flattenWorkflowTemplatePlacement(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenWorkflowTemplatePlacementSlice flattens the contents of WorkflowTemplatePlacement from a JSON
// response object.
func flattenWorkflowTemplatePlacementSlice(c *Client, i interface{}, res *WorkflowTemplate) []WorkflowTemplatePlacement {
	a, ok := i.([]interface{})
	if !ok {
		return []WorkflowTemplatePlacement{}
	}

	if len(a) == 0 {
		return []WorkflowTemplatePlacement{}
	}

	items := make([]WorkflowTemplatePlacement, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenWorkflowTemplatePlacement(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandWorkflowTemplatePlacement expands an instance of WorkflowTemplatePlacement into a JSON
// request object.
func expandWorkflowTemplatePlacement(c *Client, f *WorkflowTemplatePlacement, res *WorkflowTemplate) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v, err := expandWorkflowTemplatePlacementManagedCluster(c, f.ManagedCluster, res); err != nil {
		return nil, fmt.Errorf("error expanding ManagedCluster into managedCluster: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["managedCluster"] = v
	}
	if v, err := expandWorkflowTemplatePlacementClusterSelector(c, f.ClusterSelector, res); err != nil {
		return nil, fmt.Errorf("error expanding ClusterSelector into clusterSelector: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["clusterSelector"] = v
	}

	return m, nil
}

// flattenWorkflowTemplatePlacement flattens an instance of WorkflowTemplatePlacement from a JSON
// response object.
func flattenWorkflowTemplatePlacement(c *Client, i interface{}, res *WorkflowTemplate) *WorkflowTemplatePlacement {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &WorkflowTemplatePlacement{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyWorkflowTemplatePlacement
	}
	r.ManagedCluster = flattenWorkflowTemplatePlacementManagedCluster(c, m["managedCluster"], res)
	r.ClusterSelector = flattenWorkflowTemplatePlacementClusterSelector(c, m["clusterSelector"], res)

	return r
}

// expandWorkflowTemplatePlacementManagedClusterMap expands the contents of WorkflowTemplatePlacementManagedCluster into a JSON
// request object.
func expandWorkflowTemplatePlacementManagedClusterMap(c *Client, f map[string]WorkflowTemplatePlacementManagedCluster, res *WorkflowTemplate) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandWorkflowTemplatePlacementManagedCluster(c, &item, res)
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
func expandWorkflowTemplatePlacementManagedClusterSlice(c *Client, f []WorkflowTemplatePlacementManagedCluster, res *WorkflowTemplate) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandWorkflowTemplatePlacementManagedCluster(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenWorkflowTemplatePlacementManagedClusterMap flattens the contents of WorkflowTemplatePlacementManagedCluster from a JSON
// response object.
func flattenWorkflowTemplatePlacementManagedClusterMap(c *Client, i interface{}, res *WorkflowTemplate) map[string]WorkflowTemplatePlacementManagedCluster {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]WorkflowTemplatePlacementManagedCluster{}
	}

	if len(a) == 0 {
		return map[string]WorkflowTemplatePlacementManagedCluster{}
	}

	items := make(map[string]WorkflowTemplatePlacementManagedCluster)
	for k, item := range a {
		items[k] = *flattenWorkflowTemplatePlacementManagedCluster(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenWorkflowTemplatePlacementManagedClusterSlice flattens the contents of WorkflowTemplatePlacementManagedCluster from a JSON
// response object.
func flattenWorkflowTemplatePlacementManagedClusterSlice(c *Client, i interface{}, res *WorkflowTemplate) []WorkflowTemplatePlacementManagedCluster {
	a, ok := i.([]interface{})
	if !ok {
		return []WorkflowTemplatePlacementManagedCluster{}
	}

	if len(a) == 0 {
		return []WorkflowTemplatePlacementManagedCluster{}
	}

	items := make([]WorkflowTemplatePlacementManagedCluster, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenWorkflowTemplatePlacementManagedCluster(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandWorkflowTemplatePlacementManagedCluster expands an instance of WorkflowTemplatePlacementManagedCluster into a JSON
// request object.
func expandWorkflowTemplatePlacementManagedCluster(c *Client, f *WorkflowTemplatePlacementManagedCluster, res *WorkflowTemplate) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.ClusterName; !dcl.IsEmptyValueIndirect(v) {
		m["clusterName"] = v
	}
	if v, err := expandWorkflowTemplatePlacementManagedClusterConfig(c, f.Config, res); err != nil {
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
func flattenWorkflowTemplatePlacementManagedCluster(c *Client, i interface{}, res *WorkflowTemplate) *WorkflowTemplatePlacementManagedCluster {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &WorkflowTemplatePlacementManagedCluster{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyWorkflowTemplatePlacementManagedCluster
	}
	r.ClusterName = dcl.FlattenString(m["clusterName"])
	r.Config = flattenWorkflowTemplatePlacementManagedClusterConfig(c, m["config"], res)
	r.Labels = dcl.FlattenKeyValuePairs(m["labels"])

	return r
}

// expandWorkflowTemplatePlacementManagedClusterConfigMap expands the contents of WorkflowTemplatePlacementManagedClusterConfig into a JSON
// request object.
func expandWorkflowTemplatePlacementManagedClusterConfigMap(c *Client, f map[string]WorkflowTemplatePlacementManagedClusterConfig, res *WorkflowTemplate) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandWorkflowTemplatePlacementManagedClusterConfig(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandWorkflowTemplatePlacementManagedClusterConfigSlice expands the contents of WorkflowTemplatePlacementManagedClusterConfig into a JSON
// request object.
func expandWorkflowTemplatePlacementManagedClusterConfigSlice(c *Client, f []WorkflowTemplatePlacementManagedClusterConfig, res *WorkflowTemplate) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandWorkflowTemplatePlacementManagedClusterConfig(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenWorkflowTemplatePlacementManagedClusterConfigMap flattens the contents of WorkflowTemplatePlacementManagedClusterConfig from a JSON
// response object.
func flattenWorkflowTemplatePlacementManagedClusterConfigMap(c *Client, i interface{}, res *WorkflowTemplate) map[string]WorkflowTemplatePlacementManagedClusterConfig {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]WorkflowTemplatePlacementManagedClusterConfig{}
	}

	if len(a) == 0 {
		return map[string]WorkflowTemplatePlacementManagedClusterConfig{}
	}

	items := make(map[string]WorkflowTemplatePlacementManagedClusterConfig)
	for k, item := range a {
		items[k] = *flattenWorkflowTemplatePlacementManagedClusterConfig(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenWorkflowTemplatePlacementManagedClusterConfigSlice flattens the contents of WorkflowTemplatePlacementManagedClusterConfig from a JSON
// response object.
func flattenWorkflowTemplatePlacementManagedClusterConfigSlice(c *Client, i interface{}, res *WorkflowTemplate) []WorkflowTemplatePlacementManagedClusterConfig {
	a, ok := i.([]interface{})
	if !ok {
		return []WorkflowTemplatePlacementManagedClusterConfig{}
	}

	if len(a) == 0 {
		return []WorkflowTemplatePlacementManagedClusterConfig{}
	}

	items := make([]WorkflowTemplatePlacementManagedClusterConfig, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenWorkflowTemplatePlacementManagedClusterConfig(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandWorkflowTemplatePlacementManagedClusterConfig expands an instance of WorkflowTemplatePlacementManagedClusterConfig into a JSON
// request object.
func expandWorkflowTemplatePlacementManagedClusterConfig(c *Client, f *WorkflowTemplatePlacementManagedClusterConfig, res *WorkflowTemplate) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.StagingBucket; !dcl.IsEmptyValueIndirect(v) {
		m["configBucket"] = v
	}
	if v := f.TempBucket; !dcl.IsEmptyValueIndirect(v) {
		m["tempBucket"] = v
	}
	if v, err := expandWorkflowTemplatePlacementManagedClusterConfigGceClusterConfig(c, f.GceClusterConfig, res); err != nil {
		return nil, fmt.Errorf("error expanding GceClusterConfig into gceClusterConfig: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["gceClusterConfig"] = v
	}
	if v, err := expandWorkflowTemplatePlacementManagedClusterConfigMasterConfig(c, f.MasterConfig, res); err != nil {
		return nil, fmt.Errorf("error expanding MasterConfig into masterConfig: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["masterConfig"] = v
	}
	if v, err := expandWorkflowTemplatePlacementManagedClusterConfigWorkerConfig(c, f.WorkerConfig, res); err != nil {
		return nil, fmt.Errorf("error expanding WorkerConfig into workerConfig: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["workerConfig"] = v
	}
	if v, err := expandWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfig(c, f.SecondaryWorkerConfig, res); err != nil {
		return nil, fmt.Errorf("error expanding SecondaryWorkerConfig into secondaryWorkerConfig: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["secondaryWorkerConfig"] = v
	}
	if v, err := expandWorkflowTemplatePlacementManagedClusterConfigSoftwareConfig(c, f.SoftwareConfig, res); err != nil {
		return nil, fmt.Errorf("error expanding SoftwareConfig into softwareConfig: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["softwareConfig"] = v
	}
	if v, err := expandWorkflowTemplatePlacementManagedClusterConfigInitializationActionsSlice(c, f.InitializationActions, res); err != nil {
		return nil, fmt.Errorf("error expanding InitializationActions into initializationActions: %w", err)
	} else if v != nil {
		m["initializationActions"] = v
	}
	if v, err := expandWorkflowTemplatePlacementManagedClusterConfigEncryptionConfig(c, f.EncryptionConfig, res); err != nil {
		return nil, fmt.Errorf("error expanding EncryptionConfig into encryptionConfig: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["encryptionConfig"] = v
	}
	if v, err := expandWorkflowTemplatePlacementManagedClusterConfigAutoscalingConfig(c, f.AutoscalingConfig, res); err != nil {
		return nil, fmt.Errorf("error expanding AutoscalingConfig into autoscalingConfig: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["autoscalingConfig"] = v
	}
	if v, err := expandWorkflowTemplatePlacementManagedClusterConfigSecurityConfig(c, f.SecurityConfig, res); err != nil {
		return nil, fmt.Errorf("error expanding SecurityConfig into securityConfig: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["securityConfig"] = v
	}
	if v, err := expandWorkflowTemplatePlacementManagedClusterConfigLifecycleConfig(c, f.LifecycleConfig, res); err != nil {
		return nil, fmt.Errorf("error expanding LifecycleConfig into lifecycleConfig: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["lifecycleConfig"] = v
	}
	if v, err := expandWorkflowTemplatePlacementManagedClusterConfigEndpointConfig(c, f.EndpointConfig, res); err != nil {
		return nil, fmt.Errorf("error expanding EndpointConfig into endpointConfig: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["endpointConfig"] = v
	}

	return m, nil
}

// flattenWorkflowTemplatePlacementManagedClusterConfig flattens an instance of WorkflowTemplatePlacementManagedClusterConfig from a JSON
// response object.
func flattenWorkflowTemplatePlacementManagedClusterConfig(c *Client, i interface{}, res *WorkflowTemplate) *WorkflowTemplatePlacementManagedClusterConfig {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &WorkflowTemplatePlacementManagedClusterConfig{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyWorkflowTemplatePlacementManagedClusterConfig
	}
	r.StagingBucket = dcl.FlattenString(m["configBucket"])
	r.TempBucket = dcl.FlattenString(m["tempBucket"])
	r.GceClusterConfig = flattenWorkflowTemplatePlacementManagedClusterConfigGceClusterConfig(c, m["gceClusterConfig"], res)
	r.MasterConfig = flattenWorkflowTemplatePlacementManagedClusterConfigMasterConfig(c, m["masterConfig"], res)
	r.WorkerConfig = flattenWorkflowTemplatePlacementManagedClusterConfigWorkerConfig(c, m["workerConfig"], res)
	r.SecondaryWorkerConfig = flattenWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfig(c, m["secondaryWorkerConfig"], res)
	r.SoftwareConfig = flattenWorkflowTemplatePlacementManagedClusterConfigSoftwareConfig(c, m["softwareConfig"], res)
	r.InitializationActions = flattenWorkflowTemplatePlacementManagedClusterConfigInitializationActionsSlice(c, m["initializationActions"], res)
	r.EncryptionConfig = flattenWorkflowTemplatePlacementManagedClusterConfigEncryptionConfig(c, m["encryptionConfig"], res)
	r.AutoscalingConfig = flattenWorkflowTemplatePlacementManagedClusterConfigAutoscalingConfig(c, m["autoscalingConfig"], res)
	r.SecurityConfig = flattenWorkflowTemplatePlacementManagedClusterConfigSecurityConfig(c, m["securityConfig"], res)
	r.LifecycleConfig = flattenWorkflowTemplatePlacementManagedClusterConfigLifecycleConfig(c, m["lifecycleConfig"], res)
	r.EndpointConfig = flattenWorkflowTemplatePlacementManagedClusterConfigEndpointConfig(c, m["endpointConfig"], res)

	return r
}

// expandWorkflowTemplatePlacementManagedClusterConfigGceClusterConfigMap expands the contents of WorkflowTemplatePlacementManagedClusterConfigGceClusterConfig into a JSON
// request object.
func expandWorkflowTemplatePlacementManagedClusterConfigGceClusterConfigMap(c *Client, f map[string]WorkflowTemplatePlacementManagedClusterConfigGceClusterConfig, res *WorkflowTemplate) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandWorkflowTemplatePlacementManagedClusterConfigGceClusterConfig(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandWorkflowTemplatePlacementManagedClusterConfigGceClusterConfigSlice expands the contents of WorkflowTemplatePlacementManagedClusterConfigGceClusterConfig into a JSON
// request object.
func expandWorkflowTemplatePlacementManagedClusterConfigGceClusterConfigSlice(c *Client, f []WorkflowTemplatePlacementManagedClusterConfigGceClusterConfig, res *WorkflowTemplate) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandWorkflowTemplatePlacementManagedClusterConfigGceClusterConfig(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenWorkflowTemplatePlacementManagedClusterConfigGceClusterConfigMap flattens the contents of WorkflowTemplatePlacementManagedClusterConfigGceClusterConfig from a JSON
// response object.
func flattenWorkflowTemplatePlacementManagedClusterConfigGceClusterConfigMap(c *Client, i interface{}, res *WorkflowTemplate) map[string]WorkflowTemplatePlacementManagedClusterConfigGceClusterConfig {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]WorkflowTemplatePlacementManagedClusterConfigGceClusterConfig{}
	}

	if len(a) == 0 {
		return map[string]WorkflowTemplatePlacementManagedClusterConfigGceClusterConfig{}
	}

	items := make(map[string]WorkflowTemplatePlacementManagedClusterConfigGceClusterConfig)
	for k, item := range a {
		items[k] = *flattenWorkflowTemplatePlacementManagedClusterConfigGceClusterConfig(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenWorkflowTemplatePlacementManagedClusterConfigGceClusterConfigSlice flattens the contents of WorkflowTemplatePlacementManagedClusterConfigGceClusterConfig from a JSON
// response object.
func flattenWorkflowTemplatePlacementManagedClusterConfigGceClusterConfigSlice(c *Client, i interface{}, res *WorkflowTemplate) []WorkflowTemplatePlacementManagedClusterConfigGceClusterConfig {
	a, ok := i.([]interface{})
	if !ok {
		return []WorkflowTemplatePlacementManagedClusterConfigGceClusterConfig{}
	}

	if len(a) == 0 {
		return []WorkflowTemplatePlacementManagedClusterConfigGceClusterConfig{}
	}

	items := make([]WorkflowTemplatePlacementManagedClusterConfigGceClusterConfig, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenWorkflowTemplatePlacementManagedClusterConfigGceClusterConfig(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandWorkflowTemplatePlacementManagedClusterConfigGceClusterConfig expands an instance of WorkflowTemplatePlacementManagedClusterConfigGceClusterConfig into a JSON
// request object.
func expandWorkflowTemplatePlacementManagedClusterConfigGceClusterConfig(c *Client, f *WorkflowTemplatePlacementManagedClusterConfigGceClusterConfig, res *WorkflowTemplate) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.Zone; !dcl.IsEmptyValueIndirect(v) {
		m["zoneUri"] = v
	}
	if v := f.Network; !dcl.IsEmptyValueIndirect(v) {
		m["networkUri"] = v
	}
	if v := f.Subnetwork; !dcl.IsEmptyValueIndirect(v) {
		m["subnetworkUri"] = v
	}
	if v := f.InternalIPOnly; !dcl.IsEmptyValueIndirect(v) {
		m["internalIpOnly"] = v
	}
	if v := f.PrivateIPv6GoogleAccess; !dcl.IsEmptyValueIndirect(v) {
		m["privateIpv6GoogleAccess"] = v
	}
	if v := f.ServiceAccount; !dcl.IsEmptyValueIndirect(v) {
		m["serviceAccount"] = v
	}
	if v := f.ServiceAccountScopes; v != nil {
		m["serviceAccountScopes"] = v
	}
	if v := f.Tags; v != nil {
		m["tags"] = v
	}
	if v := f.Metadata; !dcl.IsEmptyValueIndirect(v) {
		m["metadata"] = v
	}
	if v, err := expandWorkflowTemplatePlacementManagedClusterConfigGceClusterConfigReservationAffinity(c, f.ReservationAffinity, res); err != nil {
		return nil, fmt.Errorf("error expanding ReservationAffinity into reservationAffinity: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["reservationAffinity"] = v
	}
	if v, err := expandWorkflowTemplatePlacementManagedClusterConfigGceClusterConfigNodeGroupAffinity(c, f.NodeGroupAffinity, res); err != nil {
		return nil, fmt.Errorf("error expanding NodeGroupAffinity into nodeGroupAffinity: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["nodeGroupAffinity"] = v
	}
	if v, err := expandWorkflowTemplatePlacementManagedClusterConfigGceClusterConfigShieldedInstanceConfig(c, f.ShieldedInstanceConfig, res); err != nil {
		return nil, fmt.Errorf("error expanding ShieldedInstanceConfig into shieldedInstanceConfig: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["shieldedInstanceConfig"] = v
	}

	return m, nil
}

// flattenWorkflowTemplatePlacementManagedClusterConfigGceClusterConfig flattens an instance of WorkflowTemplatePlacementManagedClusterConfigGceClusterConfig from a JSON
// response object.
func flattenWorkflowTemplatePlacementManagedClusterConfigGceClusterConfig(c *Client, i interface{}, res *WorkflowTemplate) *WorkflowTemplatePlacementManagedClusterConfigGceClusterConfig {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &WorkflowTemplatePlacementManagedClusterConfigGceClusterConfig{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyWorkflowTemplatePlacementManagedClusterConfigGceClusterConfig
	}
	r.Zone = dcl.FlattenString(m["zoneUri"])
	r.Network = dcl.FlattenString(m["networkUri"])
	r.Subnetwork = dcl.FlattenString(m["subnetworkUri"])
	r.InternalIPOnly = dcl.FlattenBool(m["internalIpOnly"])
	r.PrivateIPv6GoogleAccess = flattenWorkflowTemplatePlacementManagedClusterConfigGceClusterConfigPrivateIPv6GoogleAccessEnum(m["privateIpv6GoogleAccess"])
	r.ServiceAccount = dcl.FlattenString(m["serviceAccount"])
	r.ServiceAccountScopes = dcl.FlattenStringSlice(m["serviceAccountScopes"])
	r.Tags = dcl.FlattenStringSlice(m["tags"])
	r.Metadata = dcl.FlattenKeyValuePairs(m["metadata"])
	r.ReservationAffinity = flattenWorkflowTemplatePlacementManagedClusterConfigGceClusterConfigReservationAffinity(c, m["reservationAffinity"], res)
	r.NodeGroupAffinity = flattenWorkflowTemplatePlacementManagedClusterConfigGceClusterConfigNodeGroupAffinity(c, m["nodeGroupAffinity"], res)
	r.ShieldedInstanceConfig = flattenWorkflowTemplatePlacementManagedClusterConfigGceClusterConfigShieldedInstanceConfig(c, m["shieldedInstanceConfig"], res)

	return r
}

// expandWorkflowTemplatePlacementManagedClusterConfigGceClusterConfigReservationAffinityMap expands the contents of WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigReservationAffinity into a JSON
// request object.
func expandWorkflowTemplatePlacementManagedClusterConfigGceClusterConfigReservationAffinityMap(c *Client, f map[string]WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigReservationAffinity, res *WorkflowTemplate) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandWorkflowTemplatePlacementManagedClusterConfigGceClusterConfigReservationAffinity(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandWorkflowTemplatePlacementManagedClusterConfigGceClusterConfigReservationAffinitySlice expands the contents of WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigReservationAffinity into a JSON
// request object.
func expandWorkflowTemplatePlacementManagedClusterConfigGceClusterConfigReservationAffinitySlice(c *Client, f []WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigReservationAffinity, res *WorkflowTemplate) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandWorkflowTemplatePlacementManagedClusterConfigGceClusterConfigReservationAffinity(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenWorkflowTemplatePlacementManagedClusterConfigGceClusterConfigReservationAffinityMap flattens the contents of WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigReservationAffinity from a JSON
// response object.
func flattenWorkflowTemplatePlacementManagedClusterConfigGceClusterConfigReservationAffinityMap(c *Client, i interface{}, res *WorkflowTemplate) map[string]WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigReservationAffinity {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigReservationAffinity{}
	}

	if len(a) == 0 {
		return map[string]WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigReservationAffinity{}
	}

	items := make(map[string]WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigReservationAffinity)
	for k, item := range a {
		items[k] = *flattenWorkflowTemplatePlacementManagedClusterConfigGceClusterConfigReservationAffinity(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenWorkflowTemplatePlacementManagedClusterConfigGceClusterConfigReservationAffinitySlice flattens the contents of WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigReservationAffinity from a JSON
// response object.
func flattenWorkflowTemplatePlacementManagedClusterConfigGceClusterConfigReservationAffinitySlice(c *Client, i interface{}, res *WorkflowTemplate) []WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigReservationAffinity {
	a, ok := i.([]interface{})
	if !ok {
		return []WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigReservationAffinity{}
	}

	if len(a) == 0 {
		return []WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigReservationAffinity{}
	}

	items := make([]WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigReservationAffinity, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenWorkflowTemplatePlacementManagedClusterConfigGceClusterConfigReservationAffinity(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandWorkflowTemplatePlacementManagedClusterConfigGceClusterConfigReservationAffinity expands an instance of WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigReservationAffinity into a JSON
// request object.
func expandWorkflowTemplatePlacementManagedClusterConfigGceClusterConfigReservationAffinity(c *Client, f *WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigReservationAffinity, res *WorkflowTemplate) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.ConsumeReservationType; !dcl.IsEmptyValueIndirect(v) {
		m["consumeReservationType"] = v
	}
	if v := f.Key; !dcl.IsEmptyValueIndirect(v) {
		m["key"] = v
	}
	if v := f.Values; v != nil {
		m["values"] = v
	}

	return m, nil
}

// flattenWorkflowTemplatePlacementManagedClusterConfigGceClusterConfigReservationAffinity flattens an instance of WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigReservationAffinity from a JSON
// response object.
func flattenWorkflowTemplatePlacementManagedClusterConfigGceClusterConfigReservationAffinity(c *Client, i interface{}, res *WorkflowTemplate) *WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigReservationAffinity {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigReservationAffinity{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyWorkflowTemplatePlacementManagedClusterConfigGceClusterConfigReservationAffinity
	}
	r.ConsumeReservationType = flattenWorkflowTemplatePlacementManagedClusterConfigGceClusterConfigReservationAffinityConsumeReservationTypeEnum(m["consumeReservationType"])
	r.Key = dcl.FlattenString(m["key"])
	r.Values = dcl.FlattenStringSlice(m["values"])

	return r
}

// expandWorkflowTemplatePlacementManagedClusterConfigGceClusterConfigNodeGroupAffinityMap expands the contents of WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigNodeGroupAffinity into a JSON
// request object.
func expandWorkflowTemplatePlacementManagedClusterConfigGceClusterConfigNodeGroupAffinityMap(c *Client, f map[string]WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigNodeGroupAffinity, res *WorkflowTemplate) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandWorkflowTemplatePlacementManagedClusterConfigGceClusterConfigNodeGroupAffinity(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandWorkflowTemplatePlacementManagedClusterConfigGceClusterConfigNodeGroupAffinitySlice expands the contents of WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigNodeGroupAffinity into a JSON
// request object.
func expandWorkflowTemplatePlacementManagedClusterConfigGceClusterConfigNodeGroupAffinitySlice(c *Client, f []WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigNodeGroupAffinity, res *WorkflowTemplate) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandWorkflowTemplatePlacementManagedClusterConfigGceClusterConfigNodeGroupAffinity(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenWorkflowTemplatePlacementManagedClusterConfigGceClusterConfigNodeGroupAffinityMap flattens the contents of WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigNodeGroupAffinity from a JSON
// response object.
func flattenWorkflowTemplatePlacementManagedClusterConfigGceClusterConfigNodeGroupAffinityMap(c *Client, i interface{}, res *WorkflowTemplate) map[string]WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigNodeGroupAffinity {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigNodeGroupAffinity{}
	}

	if len(a) == 0 {
		return map[string]WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigNodeGroupAffinity{}
	}

	items := make(map[string]WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigNodeGroupAffinity)
	for k, item := range a {
		items[k] = *flattenWorkflowTemplatePlacementManagedClusterConfigGceClusterConfigNodeGroupAffinity(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenWorkflowTemplatePlacementManagedClusterConfigGceClusterConfigNodeGroupAffinitySlice flattens the contents of WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigNodeGroupAffinity from a JSON
// response object.
func flattenWorkflowTemplatePlacementManagedClusterConfigGceClusterConfigNodeGroupAffinitySlice(c *Client, i interface{}, res *WorkflowTemplate) []WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigNodeGroupAffinity {
	a, ok := i.([]interface{})
	if !ok {
		return []WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigNodeGroupAffinity{}
	}

	if len(a) == 0 {
		return []WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigNodeGroupAffinity{}
	}

	items := make([]WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigNodeGroupAffinity, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenWorkflowTemplatePlacementManagedClusterConfigGceClusterConfigNodeGroupAffinity(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandWorkflowTemplatePlacementManagedClusterConfigGceClusterConfigNodeGroupAffinity expands an instance of WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigNodeGroupAffinity into a JSON
// request object.
func expandWorkflowTemplatePlacementManagedClusterConfigGceClusterConfigNodeGroupAffinity(c *Client, f *WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigNodeGroupAffinity, res *WorkflowTemplate) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.NodeGroup; !dcl.IsEmptyValueIndirect(v) {
		m["nodeGroupUri"] = v
	}

	return m, nil
}

// flattenWorkflowTemplatePlacementManagedClusterConfigGceClusterConfigNodeGroupAffinity flattens an instance of WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigNodeGroupAffinity from a JSON
// response object.
func flattenWorkflowTemplatePlacementManagedClusterConfigGceClusterConfigNodeGroupAffinity(c *Client, i interface{}, res *WorkflowTemplate) *WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigNodeGroupAffinity {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigNodeGroupAffinity{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyWorkflowTemplatePlacementManagedClusterConfigGceClusterConfigNodeGroupAffinity
	}
	r.NodeGroup = dcl.FlattenString(m["nodeGroupUri"])

	return r
}

// expandWorkflowTemplatePlacementManagedClusterConfigGceClusterConfigShieldedInstanceConfigMap expands the contents of WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigShieldedInstanceConfig into a JSON
// request object.
func expandWorkflowTemplatePlacementManagedClusterConfigGceClusterConfigShieldedInstanceConfigMap(c *Client, f map[string]WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigShieldedInstanceConfig, res *WorkflowTemplate) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandWorkflowTemplatePlacementManagedClusterConfigGceClusterConfigShieldedInstanceConfig(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandWorkflowTemplatePlacementManagedClusterConfigGceClusterConfigShieldedInstanceConfigSlice expands the contents of WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigShieldedInstanceConfig into a JSON
// request object.
func expandWorkflowTemplatePlacementManagedClusterConfigGceClusterConfigShieldedInstanceConfigSlice(c *Client, f []WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigShieldedInstanceConfig, res *WorkflowTemplate) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandWorkflowTemplatePlacementManagedClusterConfigGceClusterConfigShieldedInstanceConfig(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenWorkflowTemplatePlacementManagedClusterConfigGceClusterConfigShieldedInstanceConfigMap flattens the contents of WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigShieldedInstanceConfig from a JSON
// response object.
func flattenWorkflowTemplatePlacementManagedClusterConfigGceClusterConfigShieldedInstanceConfigMap(c *Client, i interface{}, res *WorkflowTemplate) map[string]WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigShieldedInstanceConfig {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigShieldedInstanceConfig{}
	}

	if len(a) == 0 {
		return map[string]WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigShieldedInstanceConfig{}
	}

	items := make(map[string]WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigShieldedInstanceConfig)
	for k, item := range a {
		items[k] = *flattenWorkflowTemplatePlacementManagedClusterConfigGceClusterConfigShieldedInstanceConfig(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenWorkflowTemplatePlacementManagedClusterConfigGceClusterConfigShieldedInstanceConfigSlice flattens the contents of WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigShieldedInstanceConfig from a JSON
// response object.
func flattenWorkflowTemplatePlacementManagedClusterConfigGceClusterConfigShieldedInstanceConfigSlice(c *Client, i interface{}, res *WorkflowTemplate) []WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigShieldedInstanceConfig {
	a, ok := i.([]interface{})
	if !ok {
		return []WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigShieldedInstanceConfig{}
	}

	if len(a) == 0 {
		return []WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigShieldedInstanceConfig{}
	}

	items := make([]WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigShieldedInstanceConfig, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenWorkflowTemplatePlacementManagedClusterConfigGceClusterConfigShieldedInstanceConfig(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandWorkflowTemplatePlacementManagedClusterConfigGceClusterConfigShieldedInstanceConfig expands an instance of WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigShieldedInstanceConfig into a JSON
// request object.
func expandWorkflowTemplatePlacementManagedClusterConfigGceClusterConfigShieldedInstanceConfig(c *Client, f *WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigShieldedInstanceConfig, res *WorkflowTemplate) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.EnableSecureBoot; !dcl.IsEmptyValueIndirect(v) {
		m["enableSecureBoot"] = v
	}
	if v := f.EnableVtpm; !dcl.IsEmptyValueIndirect(v) {
		m["enableVtpm"] = v
	}
	if v := f.EnableIntegrityMonitoring; !dcl.IsEmptyValueIndirect(v) {
		m["enableIntegrityMonitoring"] = v
	}

	return m, nil
}

// flattenWorkflowTemplatePlacementManagedClusterConfigGceClusterConfigShieldedInstanceConfig flattens an instance of WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigShieldedInstanceConfig from a JSON
// response object.
func flattenWorkflowTemplatePlacementManagedClusterConfigGceClusterConfigShieldedInstanceConfig(c *Client, i interface{}, res *WorkflowTemplate) *WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigShieldedInstanceConfig {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigShieldedInstanceConfig{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyWorkflowTemplatePlacementManagedClusterConfigGceClusterConfigShieldedInstanceConfig
	}
	r.EnableSecureBoot = dcl.FlattenBool(m["enableSecureBoot"])
	r.EnableVtpm = dcl.FlattenBool(m["enableVtpm"])
	r.EnableIntegrityMonitoring = dcl.FlattenBool(m["enableIntegrityMonitoring"])

	return r
}

// expandWorkflowTemplatePlacementManagedClusterConfigMasterConfigMap expands the contents of WorkflowTemplatePlacementManagedClusterConfigMasterConfig into a JSON
// request object.
func expandWorkflowTemplatePlacementManagedClusterConfigMasterConfigMap(c *Client, f map[string]WorkflowTemplatePlacementManagedClusterConfigMasterConfig, res *WorkflowTemplate) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandWorkflowTemplatePlacementManagedClusterConfigMasterConfig(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandWorkflowTemplatePlacementManagedClusterConfigMasterConfigSlice expands the contents of WorkflowTemplatePlacementManagedClusterConfigMasterConfig into a JSON
// request object.
func expandWorkflowTemplatePlacementManagedClusterConfigMasterConfigSlice(c *Client, f []WorkflowTemplatePlacementManagedClusterConfigMasterConfig, res *WorkflowTemplate) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandWorkflowTemplatePlacementManagedClusterConfigMasterConfig(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenWorkflowTemplatePlacementManagedClusterConfigMasterConfigMap flattens the contents of WorkflowTemplatePlacementManagedClusterConfigMasterConfig from a JSON
// response object.
func flattenWorkflowTemplatePlacementManagedClusterConfigMasterConfigMap(c *Client, i interface{}, res *WorkflowTemplate) map[string]WorkflowTemplatePlacementManagedClusterConfigMasterConfig {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]WorkflowTemplatePlacementManagedClusterConfigMasterConfig{}
	}

	if len(a) == 0 {
		return map[string]WorkflowTemplatePlacementManagedClusterConfigMasterConfig{}
	}

	items := make(map[string]WorkflowTemplatePlacementManagedClusterConfigMasterConfig)
	for k, item := range a {
		items[k] = *flattenWorkflowTemplatePlacementManagedClusterConfigMasterConfig(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenWorkflowTemplatePlacementManagedClusterConfigMasterConfigSlice flattens the contents of WorkflowTemplatePlacementManagedClusterConfigMasterConfig from a JSON
// response object.
func flattenWorkflowTemplatePlacementManagedClusterConfigMasterConfigSlice(c *Client, i interface{}, res *WorkflowTemplate) []WorkflowTemplatePlacementManagedClusterConfigMasterConfig {
	a, ok := i.([]interface{})
	if !ok {
		return []WorkflowTemplatePlacementManagedClusterConfigMasterConfig{}
	}

	if len(a) == 0 {
		return []WorkflowTemplatePlacementManagedClusterConfigMasterConfig{}
	}

	items := make([]WorkflowTemplatePlacementManagedClusterConfigMasterConfig, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenWorkflowTemplatePlacementManagedClusterConfigMasterConfig(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandWorkflowTemplatePlacementManagedClusterConfigMasterConfig expands an instance of WorkflowTemplatePlacementManagedClusterConfigMasterConfig into a JSON
// request object.
func expandWorkflowTemplatePlacementManagedClusterConfigMasterConfig(c *Client, f *WorkflowTemplatePlacementManagedClusterConfigMasterConfig, res *WorkflowTemplate) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.NumInstances; !dcl.IsEmptyValueIndirect(v) {
		m["numInstances"] = v
	}
	if v := f.Image; !dcl.IsEmptyValueIndirect(v) {
		m["imageUri"] = v
	}
	if v := f.MachineType; !dcl.IsEmptyValueIndirect(v) {
		m["machineTypeUri"] = v
	}
	if v, err := expandWorkflowTemplatePlacementManagedClusterConfigMasterConfigDiskConfig(c, f.DiskConfig, res); err != nil {
		return nil, fmt.Errorf("error expanding DiskConfig into diskConfig: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["diskConfig"] = v
	}
	if v := f.Preemptibility; !dcl.IsEmptyValueIndirect(v) {
		m["preemptibility"] = v
	}
	if v, err := expandWorkflowTemplatePlacementManagedClusterConfigMasterConfigAcceleratorsSlice(c, f.Accelerators, res); err != nil {
		return nil, fmt.Errorf("error expanding Accelerators into accelerators: %w", err)
	} else if v != nil {
		m["accelerators"] = v
	}
	if v := f.MinCpuPlatform; !dcl.IsEmptyValueIndirect(v) {
		m["minCpuPlatform"] = v
	}

	return m, nil
}

// flattenWorkflowTemplatePlacementManagedClusterConfigMasterConfig flattens an instance of WorkflowTemplatePlacementManagedClusterConfigMasterConfig from a JSON
// response object.
func flattenWorkflowTemplatePlacementManagedClusterConfigMasterConfig(c *Client, i interface{}, res *WorkflowTemplate) *WorkflowTemplatePlacementManagedClusterConfigMasterConfig {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &WorkflowTemplatePlacementManagedClusterConfigMasterConfig{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyWorkflowTemplatePlacementManagedClusterConfigMasterConfig
	}
	r.NumInstances = dcl.FlattenInteger(m["numInstances"])
	r.InstanceNames = dcl.FlattenStringSlice(m["instanceNames"])
	r.Image = dcl.FlattenString(m["imageUri"])
	r.MachineType = dcl.FlattenString(m["machineTypeUri"])
	r.DiskConfig = flattenWorkflowTemplatePlacementManagedClusterConfigMasterConfigDiskConfig(c, m["diskConfig"], res)
	r.IsPreemptible = dcl.FlattenBool(m["isPreemptible"])
	r.Preemptibility = flattenWorkflowTemplatePlacementManagedClusterConfigMasterConfigPreemptibilityEnum(m["preemptibility"])
	r.ManagedGroupConfig = flattenWorkflowTemplatePlacementManagedClusterConfigMasterConfigManagedGroupConfig(c, m["managedGroupConfig"], res)
	r.Accelerators = flattenWorkflowTemplatePlacementManagedClusterConfigMasterConfigAcceleratorsSlice(c, m["accelerators"], res)
	r.MinCpuPlatform = dcl.FlattenString(m["minCpuPlatform"])

	return r
}

// expandWorkflowTemplatePlacementManagedClusterConfigMasterConfigDiskConfigMap expands the contents of WorkflowTemplatePlacementManagedClusterConfigMasterConfigDiskConfig into a JSON
// request object.
func expandWorkflowTemplatePlacementManagedClusterConfigMasterConfigDiskConfigMap(c *Client, f map[string]WorkflowTemplatePlacementManagedClusterConfigMasterConfigDiskConfig, res *WorkflowTemplate) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandWorkflowTemplatePlacementManagedClusterConfigMasterConfigDiskConfig(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandWorkflowTemplatePlacementManagedClusterConfigMasterConfigDiskConfigSlice expands the contents of WorkflowTemplatePlacementManagedClusterConfigMasterConfigDiskConfig into a JSON
// request object.
func expandWorkflowTemplatePlacementManagedClusterConfigMasterConfigDiskConfigSlice(c *Client, f []WorkflowTemplatePlacementManagedClusterConfigMasterConfigDiskConfig, res *WorkflowTemplate) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandWorkflowTemplatePlacementManagedClusterConfigMasterConfigDiskConfig(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenWorkflowTemplatePlacementManagedClusterConfigMasterConfigDiskConfigMap flattens the contents of WorkflowTemplatePlacementManagedClusterConfigMasterConfigDiskConfig from a JSON
// response object.
func flattenWorkflowTemplatePlacementManagedClusterConfigMasterConfigDiskConfigMap(c *Client, i interface{}, res *WorkflowTemplate) map[string]WorkflowTemplatePlacementManagedClusterConfigMasterConfigDiskConfig {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]WorkflowTemplatePlacementManagedClusterConfigMasterConfigDiskConfig{}
	}

	if len(a) == 0 {
		return map[string]WorkflowTemplatePlacementManagedClusterConfigMasterConfigDiskConfig{}
	}

	items := make(map[string]WorkflowTemplatePlacementManagedClusterConfigMasterConfigDiskConfig)
	for k, item := range a {
		items[k] = *flattenWorkflowTemplatePlacementManagedClusterConfigMasterConfigDiskConfig(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenWorkflowTemplatePlacementManagedClusterConfigMasterConfigDiskConfigSlice flattens the contents of WorkflowTemplatePlacementManagedClusterConfigMasterConfigDiskConfig from a JSON
// response object.
func flattenWorkflowTemplatePlacementManagedClusterConfigMasterConfigDiskConfigSlice(c *Client, i interface{}, res *WorkflowTemplate) []WorkflowTemplatePlacementManagedClusterConfigMasterConfigDiskConfig {
	a, ok := i.([]interface{})
	if !ok {
		return []WorkflowTemplatePlacementManagedClusterConfigMasterConfigDiskConfig{}
	}

	if len(a) == 0 {
		return []WorkflowTemplatePlacementManagedClusterConfigMasterConfigDiskConfig{}
	}

	items := make([]WorkflowTemplatePlacementManagedClusterConfigMasterConfigDiskConfig, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenWorkflowTemplatePlacementManagedClusterConfigMasterConfigDiskConfig(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandWorkflowTemplatePlacementManagedClusterConfigMasterConfigDiskConfig expands an instance of WorkflowTemplatePlacementManagedClusterConfigMasterConfigDiskConfig into a JSON
// request object.
func expandWorkflowTemplatePlacementManagedClusterConfigMasterConfigDiskConfig(c *Client, f *WorkflowTemplatePlacementManagedClusterConfigMasterConfigDiskConfig, res *WorkflowTemplate) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.BootDiskType; !dcl.IsEmptyValueIndirect(v) {
		m["bootDiskType"] = v
	}
	if v := f.BootDiskSizeGb; !dcl.IsEmptyValueIndirect(v) {
		m["bootDiskSizeGb"] = v
	}
	if v := f.NumLocalSsds; !dcl.IsEmptyValueIndirect(v) {
		m["numLocalSsds"] = v
	}

	return m, nil
}

// flattenWorkflowTemplatePlacementManagedClusterConfigMasterConfigDiskConfig flattens an instance of WorkflowTemplatePlacementManagedClusterConfigMasterConfigDiskConfig from a JSON
// response object.
func flattenWorkflowTemplatePlacementManagedClusterConfigMasterConfigDiskConfig(c *Client, i interface{}, res *WorkflowTemplate) *WorkflowTemplatePlacementManagedClusterConfigMasterConfigDiskConfig {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &WorkflowTemplatePlacementManagedClusterConfigMasterConfigDiskConfig{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyWorkflowTemplatePlacementManagedClusterConfigMasterConfigDiskConfig
	}
	r.BootDiskType = dcl.FlattenString(m["bootDiskType"])
	r.BootDiskSizeGb = dcl.FlattenInteger(m["bootDiskSizeGb"])
	r.NumLocalSsds = dcl.FlattenInteger(m["numLocalSsds"])

	return r
}

// expandWorkflowTemplatePlacementManagedClusterConfigMasterConfigManagedGroupConfigMap expands the contents of WorkflowTemplatePlacementManagedClusterConfigMasterConfigManagedGroupConfig into a JSON
// request object.
func expandWorkflowTemplatePlacementManagedClusterConfigMasterConfigManagedGroupConfigMap(c *Client, f map[string]WorkflowTemplatePlacementManagedClusterConfigMasterConfigManagedGroupConfig, res *WorkflowTemplate) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandWorkflowTemplatePlacementManagedClusterConfigMasterConfigManagedGroupConfig(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandWorkflowTemplatePlacementManagedClusterConfigMasterConfigManagedGroupConfigSlice expands the contents of WorkflowTemplatePlacementManagedClusterConfigMasterConfigManagedGroupConfig into a JSON
// request object.
func expandWorkflowTemplatePlacementManagedClusterConfigMasterConfigManagedGroupConfigSlice(c *Client, f []WorkflowTemplatePlacementManagedClusterConfigMasterConfigManagedGroupConfig, res *WorkflowTemplate) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandWorkflowTemplatePlacementManagedClusterConfigMasterConfigManagedGroupConfig(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenWorkflowTemplatePlacementManagedClusterConfigMasterConfigManagedGroupConfigMap flattens the contents of WorkflowTemplatePlacementManagedClusterConfigMasterConfigManagedGroupConfig from a JSON
// response object.
func flattenWorkflowTemplatePlacementManagedClusterConfigMasterConfigManagedGroupConfigMap(c *Client, i interface{}, res *WorkflowTemplate) map[string]WorkflowTemplatePlacementManagedClusterConfigMasterConfigManagedGroupConfig {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]WorkflowTemplatePlacementManagedClusterConfigMasterConfigManagedGroupConfig{}
	}

	if len(a) == 0 {
		return map[string]WorkflowTemplatePlacementManagedClusterConfigMasterConfigManagedGroupConfig{}
	}

	items := make(map[string]WorkflowTemplatePlacementManagedClusterConfigMasterConfigManagedGroupConfig)
	for k, item := range a {
		items[k] = *flattenWorkflowTemplatePlacementManagedClusterConfigMasterConfigManagedGroupConfig(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenWorkflowTemplatePlacementManagedClusterConfigMasterConfigManagedGroupConfigSlice flattens the contents of WorkflowTemplatePlacementManagedClusterConfigMasterConfigManagedGroupConfig from a JSON
// response object.
func flattenWorkflowTemplatePlacementManagedClusterConfigMasterConfigManagedGroupConfigSlice(c *Client, i interface{}, res *WorkflowTemplate) []WorkflowTemplatePlacementManagedClusterConfigMasterConfigManagedGroupConfig {
	a, ok := i.([]interface{})
	if !ok {
		return []WorkflowTemplatePlacementManagedClusterConfigMasterConfigManagedGroupConfig{}
	}

	if len(a) == 0 {
		return []WorkflowTemplatePlacementManagedClusterConfigMasterConfigManagedGroupConfig{}
	}

	items := make([]WorkflowTemplatePlacementManagedClusterConfigMasterConfigManagedGroupConfig, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenWorkflowTemplatePlacementManagedClusterConfigMasterConfigManagedGroupConfig(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandWorkflowTemplatePlacementManagedClusterConfigMasterConfigManagedGroupConfig expands an instance of WorkflowTemplatePlacementManagedClusterConfigMasterConfigManagedGroupConfig into a JSON
// request object.
func expandWorkflowTemplatePlacementManagedClusterConfigMasterConfigManagedGroupConfig(c *Client, f *WorkflowTemplatePlacementManagedClusterConfigMasterConfigManagedGroupConfig, res *WorkflowTemplate) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})

	return m, nil
}

// flattenWorkflowTemplatePlacementManagedClusterConfigMasterConfigManagedGroupConfig flattens an instance of WorkflowTemplatePlacementManagedClusterConfigMasterConfigManagedGroupConfig from a JSON
// response object.
func flattenWorkflowTemplatePlacementManagedClusterConfigMasterConfigManagedGroupConfig(c *Client, i interface{}, res *WorkflowTemplate) *WorkflowTemplatePlacementManagedClusterConfigMasterConfigManagedGroupConfig {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &WorkflowTemplatePlacementManagedClusterConfigMasterConfigManagedGroupConfig{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyWorkflowTemplatePlacementManagedClusterConfigMasterConfigManagedGroupConfig
	}
	r.InstanceTemplateName = dcl.FlattenString(m["instanceTemplateName"])
	r.InstanceGroupManagerName = dcl.FlattenString(m["instanceGroupManagerName"])

	return r
}

// expandWorkflowTemplatePlacementManagedClusterConfigMasterConfigAcceleratorsMap expands the contents of WorkflowTemplatePlacementManagedClusterConfigMasterConfigAccelerators into a JSON
// request object.
func expandWorkflowTemplatePlacementManagedClusterConfigMasterConfigAcceleratorsMap(c *Client, f map[string]WorkflowTemplatePlacementManagedClusterConfigMasterConfigAccelerators, res *WorkflowTemplate) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandWorkflowTemplatePlacementManagedClusterConfigMasterConfigAccelerators(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandWorkflowTemplatePlacementManagedClusterConfigMasterConfigAcceleratorsSlice expands the contents of WorkflowTemplatePlacementManagedClusterConfigMasterConfigAccelerators into a JSON
// request object.
func expandWorkflowTemplatePlacementManagedClusterConfigMasterConfigAcceleratorsSlice(c *Client, f []WorkflowTemplatePlacementManagedClusterConfigMasterConfigAccelerators, res *WorkflowTemplate) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandWorkflowTemplatePlacementManagedClusterConfigMasterConfigAccelerators(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenWorkflowTemplatePlacementManagedClusterConfigMasterConfigAcceleratorsMap flattens the contents of WorkflowTemplatePlacementManagedClusterConfigMasterConfigAccelerators from a JSON
// response object.
func flattenWorkflowTemplatePlacementManagedClusterConfigMasterConfigAcceleratorsMap(c *Client, i interface{}, res *WorkflowTemplate) map[string]WorkflowTemplatePlacementManagedClusterConfigMasterConfigAccelerators {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]WorkflowTemplatePlacementManagedClusterConfigMasterConfigAccelerators{}
	}

	if len(a) == 0 {
		return map[string]WorkflowTemplatePlacementManagedClusterConfigMasterConfigAccelerators{}
	}

	items := make(map[string]WorkflowTemplatePlacementManagedClusterConfigMasterConfigAccelerators)
	for k, item := range a {
		items[k] = *flattenWorkflowTemplatePlacementManagedClusterConfigMasterConfigAccelerators(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenWorkflowTemplatePlacementManagedClusterConfigMasterConfigAcceleratorsSlice flattens the contents of WorkflowTemplatePlacementManagedClusterConfigMasterConfigAccelerators from a JSON
// response object.
func flattenWorkflowTemplatePlacementManagedClusterConfigMasterConfigAcceleratorsSlice(c *Client, i interface{}, res *WorkflowTemplate) []WorkflowTemplatePlacementManagedClusterConfigMasterConfigAccelerators {
	a, ok := i.([]interface{})
	if !ok {
		return []WorkflowTemplatePlacementManagedClusterConfigMasterConfigAccelerators{}
	}

	if len(a) == 0 {
		return []WorkflowTemplatePlacementManagedClusterConfigMasterConfigAccelerators{}
	}

	items := make([]WorkflowTemplatePlacementManagedClusterConfigMasterConfigAccelerators, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenWorkflowTemplatePlacementManagedClusterConfigMasterConfigAccelerators(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandWorkflowTemplatePlacementManagedClusterConfigMasterConfigAccelerators expands an instance of WorkflowTemplatePlacementManagedClusterConfigMasterConfigAccelerators into a JSON
// request object.
func expandWorkflowTemplatePlacementManagedClusterConfigMasterConfigAccelerators(c *Client, f *WorkflowTemplatePlacementManagedClusterConfigMasterConfigAccelerators, res *WorkflowTemplate) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.AcceleratorType; !dcl.IsEmptyValueIndirect(v) {
		m["acceleratorTypeUri"] = v
	}
	if v := f.AcceleratorCount; !dcl.IsEmptyValueIndirect(v) {
		m["acceleratorCount"] = v
	}

	return m, nil
}

// flattenWorkflowTemplatePlacementManagedClusterConfigMasterConfigAccelerators flattens an instance of WorkflowTemplatePlacementManagedClusterConfigMasterConfigAccelerators from a JSON
// response object.
func flattenWorkflowTemplatePlacementManagedClusterConfigMasterConfigAccelerators(c *Client, i interface{}, res *WorkflowTemplate) *WorkflowTemplatePlacementManagedClusterConfigMasterConfigAccelerators {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &WorkflowTemplatePlacementManagedClusterConfigMasterConfigAccelerators{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyWorkflowTemplatePlacementManagedClusterConfigMasterConfigAccelerators
	}
	r.AcceleratorType = dcl.FlattenString(m["acceleratorTypeUri"])
	r.AcceleratorCount = dcl.FlattenInteger(m["acceleratorCount"])

	return r
}

// expandWorkflowTemplatePlacementManagedClusterConfigWorkerConfigMap expands the contents of WorkflowTemplatePlacementManagedClusterConfigWorkerConfig into a JSON
// request object.
func expandWorkflowTemplatePlacementManagedClusterConfigWorkerConfigMap(c *Client, f map[string]WorkflowTemplatePlacementManagedClusterConfigWorkerConfig, res *WorkflowTemplate) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandWorkflowTemplatePlacementManagedClusterConfigWorkerConfig(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandWorkflowTemplatePlacementManagedClusterConfigWorkerConfigSlice expands the contents of WorkflowTemplatePlacementManagedClusterConfigWorkerConfig into a JSON
// request object.
func expandWorkflowTemplatePlacementManagedClusterConfigWorkerConfigSlice(c *Client, f []WorkflowTemplatePlacementManagedClusterConfigWorkerConfig, res *WorkflowTemplate) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandWorkflowTemplatePlacementManagedClusterConfigWorkerConfig(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenWorkflowTemplatePlacementManagedClusterConfigWorkerConfigMap flattens the contents of WorkflowTemplatePlacementManagedClusterConfigWorkerConfig from a JSON
// response object.
func flattenWorkflowTemplatePlacementManagedClusterConfigWorkerConfigMap(c *Client, i interface{}, res *WorkflowTemplate) map[string]WorkflowTemplatePlacementManagedClusterConfigWorkerConfig {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]WorkflowTemplatePlacementManagedClusterConfigWorkerConfig{}
	}

	if len(a) == 0 {
		return map[string]WorkflowTemplatePlacementManagedClusterConfigWorkerConfig{}
	}

	items := make(map[string]WorkflowTemplatePlacementManagedClusterConfigWorkerConfig)
	for k, item := range a {
		items[k] = *flattenWorkflowTemplatePlacementManagedClusterConfigWorkerConfig(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenWorkflowTemplatePlacementManagedClusterConfigWorkerConfigSlice flattens the contents of WorkflowTemplatePlacementManagedClusterConfigWorkerConfig from a JSON
// response object.
func flattenWorkflowTemplatePlacementManagedClusterConfigWorkerConfigSlice(c *Client, i interface{}, res *WorkflowTemplate) []WorkflowTemplatePlacementManagedClusterConfigWorkerConfig {
	a, ok := i.([]interface{})
	if !ok {
		return []WorkflowTemplatePlacementManagedClusterConfigWorkerConfig{}
	}

	if len(a) == 0 {
		return []WorkflowTemplatePlacementManagedClusterConfigWorkerConfig{}
	}

	items := make([]WorkflowTemplatePlacementManagedClusterConfigWorkerConfig, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenWorkflowTemplatePlacementManagedClusterConfigWorkerConfig(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandWorkflowTemplatePlacementManagedClusterConfigWorkerConfig expands an instance of WorkflowTemplatePlacementManagedClusterConfigWorkerConfig into a JSON
// request object.
func expandWorkflowTemplatePlacementManagedClusterConfigWorkerConfig(c *Client, f *WorkflowTemplatePlacementManagedClusterConfigWorkerConfig, res *WorkflowTemplate) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.NumInstances; !dcl.IsEmptyValueIndirect(v) {
		m["numInstances"] = v
	}
	if v := f.Image; !dcl.IsEmptyValueIndirect(v) {
		m["imageUri"] = v
	}
	if v := f.MachineType; !dcl.IsEmptyValueIndirect(v) {
		m["machineTypeUri"] = v
	}
	if v, err := expandWorkflowTemplatePlacementManagedClusterConfigWorkerConfigDiskConfig(c, f.DiskConfig, res); err != nil {
		return nil, fmt.Errorf("error expanding DiskConfig into diskConfig: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["diskConfig"] = v
	}
	if v := f.Preemptibility; !dcl.IsEmptyValueIndirect(v) {
		m["preemptibility"] = v
	}
	if v, err := expandWorkflowTemplatePlacementManagedClusterConfigWorkerConfigAcceleratorsSlice(c, f.Accelerators, res); err != nil {
		return nil, fmt.Errorf("error expanding Accelerators into accelerators: %w", err)
	} else if v != nil {
		m["accelerators"] = v
	}
	if v := f.MinCpuPlatform; !dcl.IsEmptyValueIndirect(v) {
		m["minCpuPlatform"] = v
	}

	return m, nil
}

// flattenWorkflowTemplatePlacementManagedClusterConfigWorkerConfig flattens an instance of WorkflowTemplatePlacementManagedClusterConfigWorkerConfig from a JSON
// response object.
func flattenWorkflowTemplatePlacementManagedClusterConfigWorkerConfig(c *Client, i interface{}, res *WorkflowTemplate) *WorkflowTemplatePlacementManagedClusterConfigWorkerConfig {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &WorkflowTemplatePlacementManagedClusterConfigWorkerConfig{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyWorkflowTemplatePlacementManagedClusterConfigWorkerConfig
	}
	r.NumInstances = dcl.FlattenInteger(m["numInstances"])
	r.InstanceNames = dcl.FlattenStringSlice(m["instanceNames"])
	r.Image = dcl.FlattenString(m["imageUri"])
	r.MachineType = dcl.FlattenString(m["machineTypeUri"])
	r.DiskConfig = flattenWorkflowTemplatePlacementManagedClusterConfigWorkerConfigDiskConfig(c, m["diskConfig"], res)
	r.IsPreemptible = dcl.FlattenBool(m["isPreemptible"])
	r.Preemptibility = flattenWorkflowTemplatePlacementManagedClusterConfigWorkerConfigPreemptibilityEnum(m["preemptibility"])
	r.ManagedGroupConfig = flattenWorkflowTemplatePlacementManagedClusterConfigWorkerConfigManagedGroupConfig(c, m["managedGroupConfig"], res)
	r.Accelerators = flattenWorkflowTemplatePlacementManagedClusterConfigWorkerConfigAcceleratorsSlice(c, m["accelerators"], res)
	r.MinCpuPlatform = dcl.FlattenString(m["minCpuPlatform"])

	return r
}

// expandWorkflowTemplatePlacementManagedClusterConfigWorkerConfigDiskConfigMap expands the contents of WorkflowTemplatePlacementManagedClusterConfigWorkerConfigDiskConfig into a JSON
// request object.
func expandWorkflowTemplatePlacementManagedClusterConfigWorkerConfigDiskConfigMap(c *Client, f map[string]WorkflowTemplatePlacementManagedClusterConfigWorkerConfigDiskConfig, res *WorkflowTemplate) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandWorkflowTemplatePlacementManagedClusterConfigWorkerConfigDiskConfig(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandWorkflowTemplatePlacementManagedClusterConfigWorkerConfigDiskConfigSlice expands the contents of WorkflowTemplatePlacementManagedClusterConfigWorkerConfigDiskConfig into a JSON
// request object.
func expandWorkflowTemplatePlacementManagedClusterConfigWorkerConfigDiskConfigSlice(c *Client, f []WorkflowTemplatePlacementManagedClusterConfigWorkerConfigDiskConfig, res *WorkflowTemplate) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandWorkflowTemplatePlacementManagedClusterConfigWorkerConfigDiskConfig(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenWorkflowTemplatePlacementManagedClusterConfigWorkerConfigDiskConfigMap flattens the contents of WorkflowTemplatePlacementManagedClusterConfigWorkerConfigDiskConfig from a JSON
// response object.
func flattenWorkflowTemplatePlacementManagedClusterConfigWorkerConfigDiskConfigMap(c *Client, i interface{}, res *WorkflowTemplate) map[string]WorkflowTemplatePlacementManagedClusterConfigWorkerConfigDiskConfig {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]WorkflowTemplatePlacementManagedClusterConfigWorkerConfigDiskConfig{}
	}

	if len(a) == 0 {
		return map[string]WorkflowTemplatePlacementManagedClusterConfigWorkerConfigDiskConfig{}
	}

	items := make(map[string]WorkflowTemplatePlacementManagedClusterConfigWorkerConfigDiskConfig)
	for k, item := range a {
		items[k] = *flattenWorkflowTemplatePlacementManagedClusterConfigWorkerConfigDiskConfig(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenWorkflowTemplatePlacementManagedClusterConfigWorkerConfigDiskConfigSlice flattens the contents of WorkflowTemplatePlacementManagedClusterConfigWorkerConfigDiskConfig from a JSON
// response object.
func flattenWorkflowTemplatePlacementManagedClusterConfigWorkerConfigDiskConfigSlice(c *Client, i interface{}, res *WorkflowTemplate) []WorkflowTemplatePlacementManagedClusterConfigWorkerConfigDiskConfig {
	a, ok := i.([]interface{})
	if !ok {
		return []WorkflowTemplatePlacementManagedClusterConfigWorkerConfigDiskConfig{}
	}

	if len(a) == 0 {
		return []WorkflowTemplatePlacementManagedClusterConfigWorkerConfigDiskConfig{}
	}

	items := make([]WorkflowTemplatePlacementManagedClusterConfigWorkerConfigDiskConfig, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenWorkflowTemplatePlacementManagedClusterConfigWorkerConfigDiskConfig(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandWorkflowTemplatePlacementManagedClusterConfigWorkerConfigDiskConfig expands an instance of WorkflowTemplatePlacementManagedClusterConfigWorkerConfigDiskConfig into a JSON
// request object.
func expandWorkflowTemplatePlacementManagedClusterConfigWorkerConfigDiskConfig(c *Client, f *WorkflowTemplatePlacementManagedClusterConfigWorkerConfigDiskConfig, res *WorkflowTemplate) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.BootDiskType; !dcl.IsEmptyValueIndirect(v) {
		m["bootDiskType"] = v
	}
	if v := f.BootDiskSizeGb; !dcl.IsEmptyValueIndirect(v) {
		m["bootDiskSizeGb"] = v
	}
	if v := f.NumLocalSsds; !dcl.IsEmptyValueIndirect(v) {
		m["numLocalSsds"] = v
	}

	return m, nil
}

// flattenWorkflowTemplatePlacementManagedClusterConfigWorkerConfigDiskConfig flattens an instance of WorkflowTemplatePlacementManagedClusterConfigWorkerConfigDiskConfig from a JSON
// response object.
func flattenWorkflowTemplatePlacementManagedClusterConfigWorkerConfigDiskConfig(c *Client, i interface{}, res *WorkflowTemplate) *WorkflowTemplatePlacementManagedClusterConfigWorkerConfigDiskConfig {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &WorkflowTemplatePlacementManagedClusterConfigWorkerConfigDiskConfig{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyWorkflowTemplatePlacementManagedClusterConfigWorkerConfigDiskConfig
	}
	r.BootDiskType = dcl.FlattenString(m["bootDiskType"])
	r.BootDiskSizeGb = dcl.FlattenInteger(m["bootDiskSizeGb"])
	r.NumLocalSsds = dcl.FlattenInteger(m["numLocalSsds"])

	return r
}

// expandWorkflowTemplatePlacementManagedClusterConfigWorkerConfigManagedGroupConfigMap expands the contents of WorkflowTemplatePlacementManagedClusterConfigWorkerConfigManagedGroupConfig into a JSON
// request object.
func expandWorkflowTemplatePlacementManagedClusterConfigWorkerConfigManagedGroupConfigMap(c *Client, f map[string]WorkflowTemplatePlacementManagedClusterConfigWorkerConfigManagedGroupConfig, res *WorkflowTemplate) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandWorkflowTemplatePlacementManagedClusterConfigWorkerConfigManagedGroupConfig(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandWorkflowTemplatePlacementManagedClusterConfigWorkerConfigManagedGroupConfigSlice expands the contents of WorkflowTemplatePlacementManagedClusterConfigWorkerConfigManagedGroupConfig into a JSON
// request object.
func expandWorkflowTemplatePlacementManagedClusterConfigWorkerConfigManagedGroupConfigSlice(c *Client, f []WorkflowTemplatePlacementManagedClusterConfigWorkerConfigManagedGroupConfig, res *WorkflowTemplate) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandWorkflowTemplatePlacementManagedClusterConfigWorkerConfigManagedGroupConfig(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenWorkflowTemplatePlacementManagedClusterConfigWorkerConfigManagedGroupConfigMap flattens the contents of WorkflowTemplatePlacementManagedClusterConfigWorkerConfigManagedGroupConfig from a JSON
// response object.
func flattenWorkflowTemplatePlacementManagedClusterConfigWorkerConfigManagedGroupConfigMap(c *Client, i interface{}, res *WorkflowTemplate) map[string]WorkflowTemplatePlacementManagedClusterConfigWorkerConfigManagedGroupConfig {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]WorkflowTemplatePlacementManagedClusterConfigWorkerConfigManagedGroupConfig{}
	}

	if len(a) == 0 {
		return map[string]WorkflowTemplatePlacementManagedClusterConfigWorkerConfigManagedGroupConfig{}
	}

	items := make(map[string]WorkflowTemplatePlacementManagedClusterConfigWorkerConfigManagedGroupConfig)
	for k, item := range a {
		items[k] = *flattenWorkflowTemplatePlacementManagedClusterConfigWorkerConfigManagedGroupConfig(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenWorkflowTemplatePlacementManagedClusterConfigWorkerConfigManagedGroupConfigSlice flattens the contents of WorkflowTemplatePlacementManagedClusterConfigWorkerConfigManagedGroupConfig from a JSON
// response object.
func flattenWorkflowTemplatePlacementManagedClusterConfigWorkerConfigManagedGroupConfigSlice(c *Client, i interface{}, res *WorkflowTemplate) []WorkflowTemplatePlacementManagedClusterConfigWorkerConfigManagedGroupConfig {
	a, ok := i.([]interface{})
	if !ok {
		return []WorkflowTemplatePlacementManagedClusterConfigWorkerConfigManagedGroupConfig{}
	}

	if len(a) == 0 {
		return []WorkflowTemplatePlacementManagedClusterConfigWorkerConfigManagedGroupConfig{}
	}

	items := make([]WorkflowTemplatePlacementManagedClusterConfigWorkerConfigManagedGroupConfig, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenWorkflowTemplatePlacementManagedClusterConfigWorkerConfigManagedGroupConfig(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandWorkflowTemplatePlacementManagedClusterConfigWorkerConfigManagedGroupConfig expands an instance of WorkflowTemplatePlacementManagedClusterConfigWorkerConfigManagedGroupConfig into a JSON
// request object.
func expandWorkflowTemplatePlacementManagedClusterConfigWorkerConfigManagedGroupConfig(c *Client, f *WorkflowTemplatePlacementManagedClusterConfigWorkerConfigManagedGroupConfig, res *WorkflowTemplate) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})

	return m, nil
}

// flattenWorkflowTemplatePlacementManagedClusterConfigWorkerConfigManagedGroupConfig flattens an instance of WorkflowTemplatePlacementManagedClusterConfigWorkerConfigManagedGroupConfig from a JSON
// response object.
func flattenWorkflowTemplatePlacementManagedClusterConfigWorkerConfigManagedGroupConfig(c *Client, i interface{}, res *WorkflowTemplate) *WorkflowTemplatePlacementManagedClusterConfigWorkerConfigManagedGroupConfig {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &WorkflowTemplatePlacementManagedClusterConfigWorkerConfigManagedGroupConfig{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyWorkflowTemplatePlacementManagedClusterConfigWorkerConfigManagedGroupConfig
	}
	r.InstanceTemplateName = dcl.FlattenString(m["instanceTemplateName"])
	r.InstanceGroupManagerName = dcl.FlattenString(m["instanceGroupManagerName"])

	return r
}

// expandWorkflowTemplatePlacementManagedClusterConfigWorkerConfigAcceleratorsMap expands the contents of WorkflowTemplatePlacementManagedClusterConfigWorkerConfigAccelerators into a JSON
// request object.
func expandWorkflowTemplatePlacementManagedClusterConfigWorkerConfigAcceleratorsMap(c *Client, f map[string]WorkflowTemplatePlacementManagedClusterConfigWorkerConfigAccelerators, res *WorkflowTemplate) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandWorkflowTemplatePlacementManagedClusterConfigWorkerConfigAccelerators(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandWorkflowTemplatePlacementManagedClusterConfigWorkerConfigAcceleratorsSlice expands the contents of WorkflowTemplatePlacementManagedClusterConfigWorkerConfigAccelerators into a JSON
// request object.
func expandWorkflowTemplatePlacementManagedClusterConfigWorkerConfigAcceleratorsSlice(c *Client, f []WorkflowTemplatePlacementManagedClusterConfigWorkerConfigAccelerators, res *WorkflowTemplate) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandWorkflowTemplatePlacementManagedClusterConfigWorkerConfigAccelerators(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenWorkflowTemplatePlacementManagedClusterConfigWorkerConfigAcceleratorsMap flattens the contents of WorkflowTemplatePlacementManagedClusterConfigWorkerConfigAccelerators from a JSON
// response object.
func flattenWorkflowTemplatePlacementManagedClusterConfigWorkerConfigAcceleratorsMap(c *Client, i interface{}, res *WorkflowTemplate) map[string]WorkflowTemplatePlacementManagedClusterConfigWorkerConfigAccelerators {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]WorkflowTemplatePlacementManagedClusterConfigWorkerConfigAccelerators{}
	}

	if len(a) == 0 {
		return map[string]WorkflowTemplatePlacementManagedClusterConfigWorkerConfigAccelerators{}
	}

	items := make(map[string]WorkflowTemplatePlacementManagedClusterConfigWorkerConfigAccelerators)
	for k, item := range a {
		items[k] = *flattenWorkflowTemplatePlacementManagedClusterConfigWorkerConfigAccelerators(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenWorkflowTemplatePlacementManagedClusterConfigWorkerConfigAcceleratorsSlice flattens the contents of WorkflowTemplatePlacementManagedClusterConfigWorkerConfigAccelerators from a JSON
// response object.
func flattenWorkflowTemplatePlacementManagedClusterConfigWorkerConfigAcceleratorsSlice(c *Client, i interface{}, res *WorkflowTemplate) []WorkflowTemplatePlacementManagedClusterConfigWorkerConfigAccelerators {
	a, ok := i.([]interface{})
	if !ok {
		return []WorkflowTemplatePlacementManagedClusterConfigWorkerConfigAccelerators{}
	}

	if len(a) == 0 {
		return []WorkflowTemplatePlacementManagedClusterConfigWorkerConfigAccelerators{}
	}

	items := make([]WorkflowTemplatePlacementManagedClusterConfigWorkerConfigAccelerators, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenWorkflowTemplatePlacementManagedClusterConfigWorkerConfigAccelerators(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandWorkflowTemplatePlacementManagedClusterConfigWorkerConfigAccelerators expands an instance of WorkflowTemplatePlacementManagedClusterConfigWorkerConfigAccelerators into a JSON
// request object.
func expandWorkflowTemplatePlacementManagedClusterConfigWorkerConfigAccelerators(c *Client, f *WorkflowTemplatePlacementManagedClusterConfigWorkerConfigAccelerators, res *WorkflowTemplate) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.AcceleratorType; !dcl.IsEmptyValueIndirect(v) {
		m["acceleratorTypeUri"] = v
	}
	if v := f.AcceleratorCount; !dcl.IsEmptyValueIndirect(v) {
		m["acceleratorCount"] = v
	}

	return m, nil
}

// flattenWorkflowTemplatePlacementManagedClusterConfigWorkerConfigAccelerators flattens an instance of WorkflowTemplatePlacementManagedClusterConfigWorkerConfigAccelerators from a JSON
// response object.
func flattenWorkflowTemplatePlacementManagedClusterConfigWorkerConfigAccelerators(c *Client, i interface{}, res *WorkflowTemplate) *WorkflowTemplatePlacementManagedClusterConfigWorkerConfigAccelerators {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &WorkflowTemplatePlacementManagedClusterConfigWorkerConfigAccelerators{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyWorkflowTemplatePlacementManagedClusterConfigWorkerConfigAccelerators
	}
	r.AcceleratorType = dcl.FlattenString(m["acceleratorTypeUri"])
	r.AcceleratorCount = dcl.FlattenInteger(m["acceleratorCount"])

	return r
}

// expandWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigMap expands the contents of WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfig into a JSON
// request object.
func expandWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigMap(c *Client, f map[string]WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfig, res *WorkflowTemplate) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfig(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigSlice expands the contents of WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfig into a JSON
// request object.
func expandWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigSlice(c *Client, f []WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfig, res *WorkflowTemplate) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfig(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigMap flattens the contents of WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfig from a JSON
// response object.
func flattenWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigMap(c *Client, i interface{}, res *WorkflowTemplate) map[string]WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfig {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfig{}
	}

	if len(a) == 0 {
		return map[string]WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfig{}
	}

	items := make(map[string]WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfig)
	for k, item := range a {
		items[k] = *flattenWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfig(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigSlice flattens the contents of WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfig from a JSON
// response object.
func flattenWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigSlice(c *Client, i interface{}, res *WorkflowTemplate) []WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfig {
	a, ok := i.([]interface{})
	if !ok {
		return []WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfig{}
	}

	if len(a) == 0 {
		return []WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfig{}
	}

	items := make([]WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfig, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfig(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfig expands an instance of WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfig into a JSON
// request object.
func expandWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfig(c *Client, f *WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfig, res *WorkflowTemplate) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.NumInstances; !dcl.IsEmptyValueIndirect(v) {
		m["numInstances"] = v
	}
	if v := f.Image; !dcl.IsEmptyValueIndirect(v) {
		m["imageUri"] = v
	}
	if v := f.MachineType; !dcl.IsEmptyValueIndirect(v) {
		m["machineTypeUri"] = v
	}
	if v, err := expandWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigDiskConfig(c, f.DiskConfig, res); err != nil {
		return nil, fmt.Errorf("error expanding DiskConfig into diskConfig: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["diskConfig"] = v
	}
	if v := f.Preemptibility; !dcl.IsEmptyValueIndirect(v) {
		m["preemptibility"] = v
	}
	if v, err := expandWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigAcceleratorsSlice(c, f.Accelerators, res); err != nil {
		return nil, fmt.Errorf("error expanding Accelerators into accelerators: %w", err)
	} else if v != nil {
		m["accelerators"] = v
	}
	if v := f.MinCpuPlatform; !dcl.IsEmptyValueIndirect(v) {
		m["minCpuPlatform"] = v
	}

	return m, nil
}

// flattenWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfig flattens an instance of WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfig from a JSON
// response object.
func flattenWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfig(c *Client, i interface{}, res *WorkflowTemplate) *WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfig {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfig{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfig
	}
	r.NumInstances = dcl.FlattenInteger(m["numInstances"])
	r.InstanceNames = dcl.FlattenStringSlice(m["instanceNames"])
	r.Image = dcl.FlattenString(m["imageUri"])
	r.MachineType = dcl.FlattenString(m["machineTypeUri"])
	r.DiskConfig = flattenWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigDiskConfig(c, m["diskConfig"], res)
	r.IsPreemptible = dcl.FlattenBool(m["isPreemptible"])
	r.Preemptibility = flattenWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigPreemptibilityEnum(m["preemptibility"])
	r.ManagedGroupConfig = flattenWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigManagedGroupConfig(c, m["managedGroupConfig"], res)
	r.Accelerators = flattenWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigAcceleratorsSlice(c, m["accelerators"], res)
	r.MinCpuPlatform = dcl.FlattenString(m["minCpuPlatform"])

	return r
}

// expandWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigDiskConfigMap expands the contents of WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigDiskConfig into a JSON
// request object.
func expandWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigDiskConfigMap(c *Client, f map[string]WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigDiskConfig, res *WorkflowTemplate) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigDiskConfig(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigDiskConfigSlice expands the contents of WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigDiskConfig into a JSON
// request object.
func expandWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigDiskConfigSlice(c *Client, f []WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigDiskConfig, res *WorkflowTemplate) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigDiskConfig(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigDiskConfigMap flattens the contents of WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigDiskConfig from a JSON
// response object.
func flattenWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigDiskConfigMap(c *Client, i interface{}, res *WorkflowTemplate) map[string]WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigDiskConfig {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigDiskConfig{}
	}

	if len(a) == 0 {
		return map[string]WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigDiskConfig{}
	}

	items := make(map[string]WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigDiskConfig)
	for k, item := range a {
		items[k] = *flattenWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigDiskConfig(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigDiskConfigSlice flattens the contents of WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigDiskConfig from a JSON
// response object.
func flattenWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigDiskConfigSlice(c *Client, i interface{}, res *WorkflowTemplate) []WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigDiskConfig {
	a, ok := i.([]interface{})
	if !ok {
		return []WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigDiskConfig{}
	}

	if len(a) == 0 {
		return []WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigDiskConfig{}
	}

	items := make([]WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigDiskConfig, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigDiskConfig(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigDiskConfig expands an instance of WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigDiskConfig into a JSON
// request object.
func expandWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigDiskConfig(c *Client, f *WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigDiskConfig, res *WorkflowTemplate) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.BootDiskType; !dcl.IsEmptyValueIndirect(v) {
		m["bootDiskType"] = v
	}
	if v := f.BootDiskSizeGb; !dcl.IsEmptyValueIndirect(v) {
		m["bootDiskSizeGb"] = v
	}
	if v := f.NumLocalSsds; !dcl.IsEmptyValueIndirect(v) {
		m["numLocalSsds"] = v
	}

	return m, nil
}

// flattenWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigDiskConfig flattens an instance of WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigDiskConfig from a JSON
// response object.
func flattenWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigDiskConfig(c *Client, i interface{}, res *WorkflowTemplate) *WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigDiskConfig {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigDiskConfig{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigDiskConfig
	}
	r.BootDiskType = dcl.FlattenString(m["bootDiskType"])
	r.BootDiskSizeGb = dcl.FlattenInteger(m["bootDiskSizeGb"])
	r.NumLocalSsds = dcl.FlattenInteger(m["numLocalSsds"])

	return r
}

// expandWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigManagedGroupConfigMap expands the contents of WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigManagedGroupConfig into a JSON
// request object.
func expandWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigManagedGroupConfigMap(c *Client, f map[string]WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigManagedGroupConfig, res *WorkflowTemplate) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigManagedGroupConfig(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigManagedGroupConfigSlice expands the contents of WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigManagedGroupConfig into a JSON
// request object.
func expandWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigManagedGroupConfigSlice(c *Client, f []WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigManagedGroupConfig, res *WorkflowTemplate) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigManagedGroupConfig(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigManagedGroupConfigMap flattens the contents of WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigManagedGroupConfig from a JSON
// response object.
func flattenWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigManagedGroupConfigMap(c *Client, i interface{}, res *WorkflowTemplate) map[string]WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigManagedGroupConfig {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigManagedGroupConfig{}
	}

	if len(a) == 0 {
		return map[string]WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigManagedGroupConfig{}
	}

	items := make(map[string]WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigManagedGroupConfig)
	for k, item := range a {
		items[k] = *flattenWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigManagedGroupConfig(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigManagedGroupConfigSlice flattens the contents of WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigManagedGroupConfig from a JSON
// response object.
func flattenWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigManagedGroupConfigSlice(c *Client, i interface{}, res *WorkflowTemplate) []WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigManagedGroupConfig {
	a, ok := i.([]interface{})
	if !ok {
		return []WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigManagedGroupConfig{}
	}

	if len(a) == 0 {
		return []WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigManagedGroupConfig{}
	}

	items := make([]WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigManagedGroupConfig, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigManagedGroupConfig(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigManagedGroupConfig expands an instance of WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigManagedGroupConfig into a JSON
// request object.
func expandWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigManagedGroupConfig(c *Client, f *WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigManagedGroupConfig, res *WorkflowTemplate) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})

	return m, nil
}

// flattenWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigManagedGroupConfig flattens an instance of WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigManagedGroupConfig from a JSON
// response object.
func flattenWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigManagedGroupConfig(c *Client, i interface{}, res *WorkflowTemplate) *WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigManagedGroupConfig {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigManagedGroupConfig{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigManagedGroupConfig
	}
	r.InstanceTemplateName = dcl.FlattenString(m["instanceTemplateName"])
	r.InstanceGroupManagerName = dcl.FlattenString(m["instanceGroupManagerName"])

	return r
}

// expandWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigAcceleratorsMap expands the contents of WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigAccelerators into a JSON
// request object.
func expandWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigAcceleratorsMap(c *Client, f map[string]WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigAccelerators, res *WorkflowTemplate) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigAccelerators(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigAcceleratorsSlice expands the contents of WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigAccelerators into a JSON
// request object.
func expandWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigAcceleratorsSlice(c *Client, f []WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigAccelerators, res *WorkflowTemplate) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigAccelerators(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigAcceleratorsMap flattens the contents of WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigAccelerators from a JSON
// response object.
func flattenWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigAcceleratorsMap(c *Client, i interface{}, res *WorkflowTemplate) map[string]WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigAccelerators {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigAccelerators{}
	}

	if len(a) == 0 {
		return map[string]WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigAccelerators{}
	}

	items := make(map[string]WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigAccelerators)
	for k, item := range a {
		items[k] = *flattenWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigAccelerators(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigAcceleratorsSlice flattens the contents of WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigAccelerators from a JSON
// response object.
func flattenWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigAcceleratorsSlice(c *Client, i interface{}, res *WorkflowTemplate) []WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigAccelerators {
	a, ok := i.([]interface{})
	if !ok {
		return []WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigAccelerators{}
	}

	if len(a) == 0 {
		return []WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigAccelerators{}
	}

	items := make([]WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigAccelerators, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigAccelerators(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigAccelerators expands an instance of WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigAccelerators into a JSON
// request object.
func expandWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigAccelerators(c *Client, f *WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigAccelerators, res *WorkflowTemplate) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.AcceleratorType; !dcl.IsEmptyValueIndirect(v) {
		m["acceleratorTypeUri"] = v
	}
	if v := f.AcceleratorCount; !dcl.IsEmptyValueIndirect(v) {
		m["acceleratorCount"] = v
	}

	return m, nil
}

// flattenWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigAccelerators flattens an instance of WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigAccelerators from a JSON
// response object.
func flattenWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigAccelerators(c *Client, i interface{}, res *WorkflowTemplate) *WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigAccelerators {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigAccelerators{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigAccelerators
	}
	r.AcceleratorType = dcl.FlattenString(m["acceleratorTypeUri"])
	r.AcceleratorCount = dcl.FlattenInteger(m["acceleratorCount"])

	return r
}

// expandWorkflowTemplatePlacementManagedClusterConfigSoftwareConfigMap expands the contents of WorkflowTemplatePlacementManagedClusterConfigSoftwareConfig into a JSON
// request object.
func expandWorkflowTemplatePlacementManagedClusterConfigSoftwareConfigMap(c *Client, f map[string]WorkflowTemplatePlacementManagedClusterConfigSoftwareConfig, res *WorkflowTemplate) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandWorkflowTemplatePlacementManagedClusterConfigSoftwareConfig(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandWorkflowTemplatePlacementManagedClusterConfigSoftwareConfigSlice expands the contents of WorkflowTemplatePlacementManagedClusterConfigSoftwareConfig into a JSON
// request object.
func expandWorkflowTemplatePlacementManagedClusterConfigSoftwareConfigSlice(c *Client, f []WorkflowTemplatePlacementManagedClusterConfigSoftwareConfig, res *WorkflowTemplate) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandWorkflowTemplatePlacementManagedClusterConfigSoftwareConfig(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenWorkflowTemplatePlacementManagedClusterConfigSoftwareConfigMap flattens the contents of WorkflowTemplatePlacementManagedClusterConfigSoftwareConfig from a JSON
// response object.
func flattenWorkflowTemplatePlacementManagedClusterConfigSoftwareConfigMap(c *Client, i interface{}, res *WorkflowTemplate) map[string]WorkflowTemplatePlacementManagedClusterConfigSoftwareConfig {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]WorkflowTemplatePlacementManagedClusterConfigSoftwareConfig{}
	}

	if len(a) == 0 {
		return map[string]WorkflowTemplatePlacementManagedClusterConfigSoftwareConfig{}
	}

	items := make(map[string]WorkflowTemplatePlacementManagedClusterConfigSoftwareConfig)
	for k, item := range a {
		items[k] = *flattenWorkflowTemplatePlacementManagedClusterConfigSoftwareConfig(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenWorkflowTemplatePlacementManagedClusterConfigSoftwareConfigSlice flattens the contents of WorkflowTemplatePlacementManagedClusterConfigSoftwareConfig from a JSON
// response object.
func flattenWorkflowTemplatePlacementManagedClusterConfigSoftwareConfigSlice(c *Client, i interface{}, res *WorkflowTemplate) []WorkflowTemplatePlacementManagedClusterConfigSoftwareConfig {
	a, ok := i.([]interface{})
	if !ok {
		return []WorkflowTemplatePlacementManagedClusterConfigSoftwareConfig{}
	}

	if len(a) == 0 {
		return []WorkflowTemplatePlacementManagedClusterConfigSoftwareConfig{}
	}

	items := make([]WorkflowTemplatePlacementManagedClusterConfigSoftwareConfig, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenWorkflowTemplatePlacementManagedClusterConfigSoftwareConfig(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandWorkflowTemplatePlacementManagedClusterConfigSoftwareConfig expands an instance of WorkflowTemplatePlacementManagedClusterConfigSoftwareConfig into a JSON
// request object.
func expandWorkflowTemplatePlacementManagedClusterConfigSoftwareConfig(c *Client, f *WorkflowTemplatePlacementManagedClusterConfigSoftwareConfig, res *WorkflowTemplate) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.ImageVersion; !dcl.IsEmptyValueIndirect(v) {
		m["imageVersion"] = v
	}
	if v := f.Properties; !dcl.IsEmptyValueIndirect(v) {
		m["properties"] = v
	}
	if v := f.OptionalComponents; v != nil {
		m["optionalComponents"] = v
	}

	return m, nil
}

// flattenWorkflowTemplatePlacementManagedClusterConfigSoftwareConfig flattens an instance of WorkflowTemplatePlacementManagedClusterConfigSoftwareConfig from a JSON
// response object.
func flattenWorkflowTemplatePlacementManagedClusterConfigSoftwareConfig(c *Client, i interface{}, res *WorkflowTemplate) *WorkflowTemplatePlacementManagedClusterConfigSoftwareConfig {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &WorkflowTemplatePlacementManagedClusterConfigSoftwareConfig{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyWorkflowTemplatePlacementManagedClusterConfigSoftwareConfig
	}
	r.ImageVersion = dcl.FlattenString(m["imageVersion"])
	r.Properties = dcl.FlattenKeyValuePairs(m["properties"])
	r.OptionalComponents = flattenWorkflowTemplatePlacementManagedClusterConfigSoftwareConfigOptionalComponentsEnumSlice(c, m["optionalComponents"], res)

	return r
}

// expandWorkflowTemplatePlacementManagedClusterConfigInitializationActionsMap expands the contents of WorkflowTemplatePlacementManagedClusterConfigInitializationActions into a JSON
// request object.
func expandWorkflowTemplatePlacementManagedClusterConfigInitializationActionsMap(c *Client, f map[string]WorkflowTemplatePlacementManagedClusterConfigInitializationActions, res *WorkflowTemplate) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandWorkflowTemplatePlacementManagedClusterConfigInitializationActions(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandWorkflowTemplatePlacementManagedClusterConfigInitializationActionsSlice expands the contents of WorkflowTemplatePlacementManagedClusterConfigInitializationActions into a JSON
// request object.
func expandWorkflowTemplatePlacementManagedClusterConfigInitializationActionsSlice(c *Client, f []WorkflowTemplatePlacementManagedClusterConfigInitializationActions, res *WorkflowTemplate) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandWorkflowTemplatePlacementManagedClusterConfigInitializationActions(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenWorkflowTemplatePlacementManagedClusterConfigInitializationActionsMap flattens the contents of WorkflowTemplatePlacementManagedClusterConfigInitializationActions from a JSON
// response object.
func flattenWorkflowTemplatePlacementManagedClusterConfigInitializationActionsMap(c *Client, i interface{}, res *WorkflowTemplate) map[string]WorkflowTemplatePlacementManagedClusterConfigInitializationActions {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]WorkflowTemplatePlacementManagedClusterConfigInitializationActions{}
	}

	if len(a) == 0 {
		return map[string]WorkflowTemplatePlacementManagedClusterConfigInitializationActions{}
	}

	items := make(map[string]WorkflowTemplatePlacementManagedClusterConfigInitializationActions)
	for k, item := range a {
		items[k] = *flattenWorkflowTemplatePlacementManagedClusterConfigInitializationActions(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenWorkflowTemplatePlacementManagedClusterConfigInitializationActionsSlice flattens the contents of WorkflowTemplatePlacementManagedClusterConfigInitializationActions from a JSON
// response object.
func flattenWorkflowTemplatePlacementManagedClusterConfigInitializationActionsSlice(c *Client, i interface{}, res *WorkflowTemplate) []WorkflowTemplatePlacementManagedClusterConfigInitializationActions {
	a, ok := i.([]interface{})
	if !ok {
		return []WorkflowTemplatePlacementManagedClusterConfigInitializationActions{}
	}

	if len(a) == 0 {
		return []WorkflowTemplatePlacementManagedClusterConfigInitializationActions{}
	}

	items := make([]WorkflowTemplatePlacementManagedClusterConfigInitializationActions, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenWorkflowTemplatePlacementManagedClusterConfigInitializationActions(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandWorkflowTemplatePlacementManagedClusterConfigInitializationActions expands an instance of WorkflowTemplatePlacementManagedClusterConfigInitializationActions into a JSON
// request object.
func expandWorkflowTemplatePlacementManagedClusterConfigInitializationActions(c *Client, f *WorkflowTemplatePlacementManagedClusterConfigInitializationActions, res *WorkflowTemplate) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.ExecutableFile; !dcl.IsEmptyValueIndirect(v) {
		m["executableFile"] = v
	}
	if v := f.ExecutionTimeout; !dcl.IsEmptyValueIndirect(v) {
		m["executionTimeout"] = v
	}

	return m, nil
}

// flattenWorkflowTemplatePlacementManagedClusterConfigInitializationActions flattens an instance of WorkflowTemplatePlacementManagedClusterConfigInitializationActions from a JSON
// response object.
func flattenWorkflowTemplatePlacementManagedClusterConfigInitializationActions(c *Client, i interface{}, res *WorkflowTemplate) *WorkflowTemplatePlacementManagedClusterConfigInitializationActions {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &WorkflowTemplatePlacementManagedClusterConfigInitializationActions{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyWorkflowTemplatePlacementManagedClusterConfigInitializationActions
	}
	r.ExecutableFile = dcl.FlattenString(m["executableFile"])
	r.ExecutionTimeout = dcl.FlattenString(m["executionTimeout"])

	return r
}

// expandWorkflowTemplatePlacementManagedClusterConfigEncryptionConfigMap expands the contents of WorkflowTemplatePlacementManagedClusterConfigEncryptionConfig into a JSON
// request object.
func expandWorkflowTemplatePlacementManagedClusterConfigEncryptionConfigMap(c *Client, f map[string]WorkflowTemplatePlacementManagedClusterConfigEncryptionConfig, res *WorkflowTemplate) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandWorkflowTemplatePlacementManagedClusterConfigEncryptionConfig(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandWorkflowTemplatePlacementManagedClusterConfigEncryptionConfigSlice expands the contents of WorkflowTemplatePlacementManagedClusterConfigEncryptionConfig into a JSON
// request object.
func expandWorkflowTemplatePlacementManagedClusterConfigEncryptionConfigSlice(c *Client, f []WorkflowTemplatePlacementManagedClusterConfigEncryptionConfig, res *WorkflowTemplate) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandWorkflowTemplatePlacementManagedClusterConfigEncryptionConfig(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenWorkflowTemplatePlacementManagedClusterConfigEncryptionConfigMap flattens the contents of WorkflowTemplatePlacementManagedClusterConfigEncryptionConfig from a JSON
// response object.
func flattenWorkflowTemplatePlacementManagedClusterConfigEncryptionConfigMap(c *Client, i interface{}, res *WorkflowTemplate) map[string]WorkflowTemplatePlacementManagedClusterConfigEncryptionConfig {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]WorkflowTemplatePlacementManagedClusterConfigEncryptionConfig{}
	}

	if len(a) == 0 {
		return map[string]WorkflowTemplatePlacementManagedClusterConfigEncryptionConfig{}
	}

	items := make(map[string]WorkflowTemplatePlacementManagedClusterConfigEncryptionConfig)
	for k, item := range a {
		items[k] = *flattenWorkflowTemplatePlacementManagedClusterConfigEncryptionConfig(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenWorkflowTemplatePlacementManagedClusterConfigEncryptionConfigSlice flattens the contents of WorkflowTemplatePlacementManagedClusterConfigEncryptionConfig from a JSON
// response object.
func flattenWorkflowTemplatePlacementManagedClusterConfigEncryptionConfigSlice(c *Client, i interface{}, res *WorkflowTemplate) []WorkflowTemplatePlacementManagedClusterConfigEncryptionConfig {
	a, ok := i.([]interface{})
	if !ok {
		return []WorkflowTemplatePlacementManagedClusterConfigEncryptionConfig{}
	}

	if len(a) == 0 {
		return []WorkflowTemplatePlacementManagedClusterConfigEncryptionConfig{}
	}

	items := make([]WorkflowTemplatePlacementManagedClusterConfigEncryptionConfig, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenWorkflowTemplatePlacementManagedClusterConfigEncryptionConfig(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandWorkflowTemplatePlacementManagedClusterConfigEncryptionConfig expands an instance of WorkflowTemplatePlacementManagedClusterConfigEncryptionConfig into a JSON
// request object.
func expandWorkflowTemplatePlacementManagedClusterConfigEncryptionConfig(c *Client, f *WorkflowTemplatePlacementManagedClusterConfigEncryptionConfig, res *WorkflowTemplate) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.GcePdKmsKeyName; !dcl.IsEmptyValueIndirect(v) {
		m["gcePdKmsKeyName"] = v
	}

	return m, nil
}

// flattenWorkflowTemplatePlacementManagedClusterConfigEncryptionConfig flattens an instance of WorkflowTemplatePlacementManagedClusterConfigEncryptionConfig from a JSON
// response object.
func flattenWorkflowTemplatePlacementManagedClusterConfigEncryptionConfig(c *Client, i interface{}, res *WorkflowTemplate) *WorkflowTemplatePlacementManagedClusterConfigEncryptionConfig {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &WorkflowTemplatePlacementManagedClusterConfigEncryptionConfig{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyWorkflowTemplatePlacementManagedClusterConfigEncryptionConfig
	}
	r.GcePdKmsKeyName = dcl.FlattenString(m["gcePdKmsKeyName"])

	return r
}

// expandWorkflowTemplatePlacementManagedClusterConfigAutoscalingConfigMap expands the contents of WorkflowTemplatePlacementManagedClusterConfigAutoscalingConfig into a JSON
// request object.
func expandWorkflowTemplatePlacementManagedClusterConfigAutoscalingConfigMap(c *Client, f map[string]WorkflowTemplatePlacementManagedClusterConfigAutoscalingConfig, res *WorkflowTemplate) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandWorkflowTemplatePlacementManagedClusterConfigAutoscalingConfig(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandWorkflowTemplatePlacementManagedClusterConfigAutoscalingConfigSlice expands the contents of WorkflowTemplatePlacementManagedClusterConfigAutoscalingConfig into a JSON
// request object.
func expandWorkflowTemplatePlacementManagedClusterConfigAutoscalingConfigSlice(c *Client, f []WorkflowTemplatePlacementManagedClusterConfigAutoscalingConfig, res *WorkflowTemplate) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandWorkflowTemplatePlacementManagedClusterConfigAutoscalingConfig(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenWorkflowTemplatePlacementManagedClusterConfigAutoscalingConfigMap flattens the contents of WorkflowTemplatePlacementManagedClusterConfigAutoscalingConfig from a JSON
// response object.
func flattenWorkflowTemplatePlacementManagedClusterConfigAutoscalingConfigMap(c *Client, i interface{}, res *WorkflowTemplate) map[string]WorkflowTemplatePlacementManagedClusterConfigAutoscalingConfig {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]WorkflowTemplatePlacementManagedClusterConfigAutoscalingConfig{}
	}

	if len(a) == 0 {
		return map[string]WorkflowTemplatePlacementManagedClusterConfigAutoscalingConfig{}
	}

	items := make(map[string]WorkflowTemplatePlacementManagedClusterConfigAutoscalingConfig)
	for k, item := range a {
		items[k] = *flattenWorkflowTemplatePlacementManagedClusterConfigAutoscalingConfig(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenWorkflowTemplatePlacementManagedClusterConfigAutoscalingConfigSlice flattens the contents of WorkflowTemplatePlacementManagedClusterConfigAutoscalingConfig from a JSON
// response object.
func flattenWorkflowTemplatePlacementManagedClusterConfigAutoscalingConfigSlice(c *Client, i interface{}, res *WorkflowTemplate) []WorkflowTemplatePlacementManagedClusterConfigAutoscalingConfig {
	a, ok := i.([]interface{})
	if !ok {
		return []WorkflowTemplatePlacementManagedClusterConfigAutoscalingConfig{}
	}

	if len(a) == 0 {
		return []WorkflowTemplatePlacementManagedClusterConfigAutoscalingConfig{}
	}

	items := make([]WorkflowTemplatePlacementManagedClusterConfigAutoscalingConfig, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenWorkflowTemplatePlacementManagedClusterConfigAutoscalingConfig(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandWorkflowTemplatePlacementManagedClusterConfigAutoscalingConfig expands an instance of WorkflowTemplatePlacementManagedClusterConfigAutoscalingConfig into a JSON
// request object.
func expandWorkflowTemplatePlacementManagedClusterConfigAutoscalingConfig(c *Client, f *WorkflowTemplatePlacementManagedClusterConfigAutoscalingConfig, res *WorkflowTemplate) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.Policy; !dcl.IsEmptyValueIndirect(v) {
		m["policyUri"] = v
	}

	return m, nil
}

// flattenWorkflowTemplatePlacementManagedClusterConfigAutoscalingConfig flattens an instance of WorkflowTemplatePlacementManagedClusterConfigAutoscalingConfig from a JSON
// response object.
func flattenWorkflowTemplatePlacementManagedClusterConfigAutoscalingConfig(c *Client, i interface{}, res *WorkflowTemplate) *WorkflowTemplatePlacementManagedClusterConfigAutoscalingConfig {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &WorkflowTemplatePlacementManagedClusterConfigAutoscalingConfig{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyWorkflowTemplatePlacementManagedClusterConfigAutoscalingConfig
	}
	r.Policy = dcl.FlattenString(m["policyUri"])

	return r
}

// expandWorkflowTemplatePlacementManagedClusterConfigSecurityConfigMap expands the contents of WorkflowTemplatePlacementManagedClusterConfigSecurityConfig into a JSON
// request object.
func expandWorkflowTemplatePlacementManagedClusterConfigSecurityConfigMap(c *Client, f map[string]WorkflowTemplatePlacementManagedClusterConfigSecurityConfig, res *WorkflowTemplate) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandWorkflowTemplatePlacementManagedClusterConfigSecurityConfig(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandWorkflowTemplatePlacementManagedClusterConfigSecurityConfigSlice expands the contents of WorkflowTemplatePlacementManagedClusterConfigSecurityConfig into a JSON
// request object.
func expandWorkflowTemplatePlacementManagedClusterConfigSecurityConfigSlice(c *Client, f []WorkflowTemplatePlacementManagedClusterConfigSecurityConfig, res *WorkflowTemplate) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandWorkflowTemplatePlacementManagedClusterConfigSecurityConfig(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenWorkflowTemplatePlacementManagedClusterConfigSecurityConfigMap flattens the contents of WorkflowTemplatePlacementManagedClusterConfigSecurityConfig from a JSON
// response object.
func flattenWorkflowTemplatePlacementManagedClusterConfigSecurityConfigMap(c *Client, i interface{}, res *WorkflowTemplate) map[string]WorkflowTemplatePlacementManagedClusterConfigSecurityConfig {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]WorkflowTemplatePlacementManagedClusterConfigSecurityConfig{}
	}

	if len(a) == 0 {
		return map[string]WorkflowTemplatePlacementManagedClusterConfigSecurityConfig{}
	}

	items := make(map[string]WorkflowTemplatePlacementManagedClusterConfigSecurityConfig)
	for k, item := range a {
		items[k] = *flattenWorkflowTemplatePlacementManagedClusterConfigSecurityConfig(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenWorkflowTemplatePlacementManagedClusterConfigSecurityConfigSlice flattens the contents of WorkflowTemplatePlacementManagedClusterConfigSecurityConfig from a JSON
// response object.
func flattenWorkflowTemplatePlacementManagedClusterConfigSecurityConfigSlice(c *Client, i interface{}, res *WorkflowTemplate) []WorkflowTemplatePlacementManagedClusterConfigSecurityConfig {
	a, ok := i.([]interface{})
	if !ok {
		return []WorkflowTemplatePlacementManagedClusterConfigSecurityConfig{}
	}

	if len(a) == 0 {
		return []WorkflowTemplatePlacementManagedClusterConfigSecurityConfig{}
	}

	items := make([]WorkflowTemplatePlacementManagedClusterConfigSecurityConfig, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenWorkflowTemplatePlacementManagedClusterConfigSecurityConfig(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandWorkflowTemplatePlacementManagedClusterConfigSecurityConfig expands an instance of WorkflowTemplatePlacementManagedClusterConfigSecurityConfig into a JSON
// request object.
func expandWorkflowTemplatePlacementManagedClusterConfigSecurityConfig(c *Client, f *WorkflowTemplatePlacementManagedClusterConfigSecurityConfig, res *WorkflowTemplate) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v, err := expandWorkflowTemplatePlacementManagedClusterConfigSecurityConfigKerberosConfig(c, f.KerberosConfig, res); err != nil {
		return nil, fmt.Errorf("error expanding KerberosConfig into kerberosConfig: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["kerberosConfig"] = v
	}

	return m, nil
}

// flattenWorkflowTemplatePlacementManagedClusterConfigSecurityConfig flattens an instance of WorkflowTemplatePlacementManagedClusterConfigSecurityConfig from a JSON
// response object.
func flattenWorkflowTemplatePlacementManagedClusterConfigSecurityConfig(c *Client, i interface{}, res *WorkflowTemplate) *WorkflowTemplatePlacementManagedClusterConfigSecurityConfig {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &WorkflowTemplatePlacementManagedClusterConfigSecurityConfig{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyWorkflowTemplatePlacementManagedClusterConfigSecurityConfig
	}
	r.KerberosConfig = flattenWorkflowTemplatePlacementManagedClusterConfigSecurityConfigKerberosConfig(c, m["kerberosConfig"], res)

	return r
}

// expandWorkflowTemplatePlacementManagedClusterConfigSecurityConfigKerberosConfigMap expands the contents of WorkflowTemplatePlacementManagedClusterConfigSecurityConfigKerberosConfig into a JSON
// request object.
func expandWorkflowTemplatePlacementManagedClusterConfigSecurityConfigKerberosConfigMap(c *Client, f map[string]WorkflowTemplatePlacementManagedClusterConfigSecurityConfigKerberosConfig, res *WorkflowTemplate) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandWorkflowTemplatePlacementManagedClusterConfigSecurityConfigKerberosConfig(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandWorkflowTemplatePlacementManagedClusterConfigSecurityConfigKerberosConfigSlice expands the contents of WorkflowTemplatePlacementManagedClusterConfigSecurityConfigKerberosConfig into a JSON
// request object.
func expandWorkflowTemplatePlacementManagedClusterConfigSecurityConfigKerberosConfigSlice(c *Client, f []WorkflowTemplatePlacementManagedClusterConfigSecurityConfigKerberosConfig, res *WorkflowTemplate) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandWorkflowTemplatePlacementManagedClusterConfigSecurityConfigKerberosConfig(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenWorkflowTemplatePlacementManagedClusterConfigSecurityConfigKerberosConfigMap flattens the contents of WorkflowTemplatePlacementManagedClusterConfigSecurityConfigKerberosConfig from a JSON
// response object.
func flattenWorkflowTemplatePlacementManagedClusterConfigSecurityConfigKerberosConfigMap(c *Client, i interface{}, res *WorkflowTemplate) map[string]WorkflowTemplatePlacementManagedClusterConfigSecurityConfigKerberosConfig {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]WorkflowTemplatePlacementManagedClusterConfigSecurityConfigKerberosConfig{}
	}

	if len(a) == 0 {
		return map[string]WorkflowTemplatePlacementManagedClusterConfigSecurityConfigKerberosConfig{}
	}

	items := make(map[string]WorkflowTemplatePlacementManagedClusterConfigSecurityConfigKerberosConfig)
	for k, item := range a {
		items[k] = *flattenWorkflowTemplatePlacementManagedClusterConfigSecurityConfigKerberosConfig(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenWorkflowTemplatePlacementManagedClusterConfigSecurityConfigKerberosConfigSlice flattens the contents of WorkflowTemplatePlacementManagedClusterConfigSecurityConfigKerberosConfig from a JSON
// response object.
func flattenWorkflowTemplatePlacementManagedClusterConfigSecurityConfigKerberosConfigSlice(c *Client, i interface{}, res *WorkflowTemplate) []WorkflowTemplatePlacementManagedClusterConfigSecurityConfigKerberosConfig {
	a, ok := i.([]interface{})
	if !ok {
		return []WorkflowTemplatePlacementManagedClusterConfigSecurityConfigKerberosConfig{}
	}

	if len(a) == 0 {
		return []WorkflowTemplatePlacementManagedClusterConfigSecurityConfigKerberosConfig{}
	}

	items := make([]WorkflowTemplatePlacementManagedClusterConfigSecurityConfigKerberosConfig, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenWorkflowTemplatePlacementManagedClusterConfigSecurityConfigKerberosConfig(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandWorkflowTemplatePlacementManagedClusterConfigSecurityConfigKerberosConfig expands an instance of WorkflowTemplatePlacementManagedClusterConfigSecurityConfigKerberosConfig into a JSON
// request object.
func expandWorkflowTemplatePlacementManagedClusterConfigSecurityConfigKerberosConfig(c *Client, f *WorkflowTemplatePlacementManagedClusterConfigSecurityConfigKerberosConfig, res *WorkflowTemplate) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.EnableKerberos; !dcl.IsEmptyValueIndirect(v) {
		m["enableKerberos"] = v
	}
	if v := f.RootPrincipalPassword; !dcl.IsEmptyValueIndirect(v) {
		m["rootPrincipalPasswordUri"] = v
	}
	if v := f.KmsKey; !dcl.IsEmptyValueIndirect(v) {
		m["kmsKeyUri"] = v
	}
	if v := f.Keystore; !dcl.IsEmptyValueIndirect(v) {
		m["keystoreUri"] = v
	}
	if v := f.Truststore; !dcl.IsEmptyValueIndirect(v) {
		m["truststoreUri"] = v
	}
	if v := f.KeystorePassword; !dcl.IsEmptyValueIndirect(v) {
		m["keystorePasswordUri"] = v
	}
	if v := f.KeyPassword; !dcl.IsEmptyValueIndirect(v) {
		m["keyPasswordUri"] = v
	}
	if v := f.TruststorePassword; !dcl.IsEmptyValueIndirect(v) {
		m["truststorePasswordUri"] = v
	}
	if v := f.CrossRealmTrustRealm; !dcl.IsEmptyValueIndirect(v) {
		m["crossRealmTrustRealm"] = v
	}
	if v := f.CrossRealmTrustKdc; !dcl.IsEmptyValueIndirect(v) {
		m["crossRealmTrustKdc"] = v
	}
	if v := f.CrossRealmTrustAdminServer; !dcl.IsEmptyValueIndirect(v) {
		m["crossRealmTrustAdminServer"] = v
	}
	if v := f.CrossRealmTrustSharedPassword; !dcl.IsEmptyValueIndirect(v) {
		m["crossRealmTrustSharedPasswordUri"] = v
	}
	if v := f.KdcDbKey; !dcl.IsEmptyValueIndirect(v) {
		m["kdcDbKeyUri"] = v
	}
	if v := f.TgtLifetimeHours; !dcl.IsEmptyValueIndirect(v) {
		m["tgtLifetimeHours"] = v
	}
	if v := f.Realm; !dcl.IsEmptyValueIndirect(v) {
		m["realm"] = v
	}

	return m, nil
}

// flattenWorkflowTemplatePlacementManagedClusterConfigSecurityConfigKerberosConfig flattens an instance of WorkflowTemplatePlacementManagedClusterConfigSecurityConfigKerberosConfig from a JSON
// response object.
func flattenWorkflowTemplatePlacementManagedClusterConfigSecurityConfigKerberosConfig(c *Client, i interface{}, res *WorkflowTemplate) *WorkflowTemplatePlacementManagedClusterConfigSecurityConfigKerberosConfig {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &WorkflowTemplatePlacementManagedClusterConfigSecurityConfigKerberosConfig{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyWorkflowTemplatePlacementManagedClusterConfigSecurityConfigKerberosConfig
	}
	r.EnableKerberos = dcl.FlattenBool(m["enableKerberos"])
	r.RootPrincipalPassword = dcl.FlattenString(m["rootPrincipalPasswordUri"])
	r.KmsKey = dcl.FlattenString(m["kmsKeyUri"])
	r.Keystore = dcl.FlattenString(m["keystoreUri"])
	r.Truststore = dcl.FlattenString(m["truststoreUri"])
	r.KeystorePassword = dcl.FlattenString(m["keystorePasswordUri"])
	r.KeyPassword = dcl.FlattenString(m["keyPasswordUri"])
	r.TruststorePassword = dcl.FlattenString(m["truststorePasswordUri"])
	r.CrossRealmTrustRealm = dcl.FlattenString(m["crossRealmTrustRealm"])
	r.CrossRealmTrustKdc = dcl.FlattenString(m["crossRealmTrustKdc"])
	r.CrossRealmTrustAdminServer = dcl.FlattenString(m["crossRealmTrustAdminServer"])
	r.CrossRealmTrustSharedPassword = dcl.FlattenString(m["crossRealmTrustSharedPasswordUri"])
	r.KdcDbKey = dcl.FlattenString(m["kdcDbKeyUri"])
	r.TgtLifetimeHours = dcl.FlattenInteger(m["tgtLifetimeHours"])
	r.Realm = dcl.FlattenString(m["realm"])

	return r
}

// expandWorkflowTemplatePlacementManagedClusterConfigLifecycleConfigMap expands the contents of WorkflowTemplatePlacementManagedClusterConfigLifecycleConfig into a JSON
// request object.
func expandWorkflowTemplatePlacementManagedClusterConfigLifecycleConfigMap(c *Client, f map[string]WorkflowTemplatePlacementManagedClusterConfigLifecycleConfig, res *WorkflowTemplate) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandWorkflowTemplatePlacementManagedClusterConfigLifecycleConfig(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandWorkflowTemplatePlacementManagedClusterConfigLifecycleConfigSlice expands the contents of WorkflowTemplatePlacementManagedClusterConfigLifecycleConfig into a JSON
// request object.
func expandWorkflowTemplatePlacementManagedClusterConfigLifecycleConfigSlice(c *Client, f []WorkflowTemplatePlacementManagedClusterConfigLifecycleConfig, res *WorkflowTemplate) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandWorkflowTemplatePlacementManagedClusterConfigLifecycleConfig(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenWorkflowTemplatePlacementManagedClusterConfigLifecycleConfigMap flattens the contents of WorkflowTemplatePlacementManagedClusterConfigLifecycleConfig from a JSON
// response object.
func flattenWorkflowTemplatePlacementManagedClusterConfigLifecycleConfigMap(c *Client, i interface{}, res *WorkflowTemplate) map[string]WorkflowTemplatePlacementManagedClusterConfigLifecycleConfig {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]WorkflowTemplatePlacementManagedClusterConfigLifecycleConfig{}
	}

	if len(a) == 0 {
		return map[string]WorkflowTemplatePlacementManagedClusterConfigLifecycleConfig{}
	}

	items := make(map[string]WorkflowTemplatePlacementManagedClusterConfigLifecycleConfig)
	for k, item := range a {
		items[k] = *flattenWorkflowTemplatePlacementManagedClusterConfigLifecycleConfig(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenWorkflowTemplatePlacementManagedClusterConfigLifecycleConfigSlice flattens the contents of WorkflowTemplatePlacementManagedClusterConfigLifecycleConfig from a JSON
// response object.
func flattenWorkflowTemplatePlacementManagedClusterConfigLifecycleConfigSlice(c *Client, i interface{}, res *WorkflowTemplate) []WorkflowTemplatePlacementManagedClusterConfigLifecycleConfig {
	a, ok := i.([]interface{})
	if !ok {
		return []WorkflowTemplatePlacementManagedClusterConfigLifecycleConfig{}
	}

	if len(a) == 0 {
		return []WorkflowTemplatePlacementManagedClusterConfigLifecycleConfig{}
	}

	items := make([]WorkflowTemplatePlacementManagedClusterConfigLifecycleConfig, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenWorkflowTemplatePlacementManagedClusterConfigLifecycleConfig(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandWorkflowTemplatePlacementManagedClusterConfigLifecycleConfig expands an instance of WorkflowTemplatePlacementManagedClusterConfigLifecycleConfig into a JSON
// request object.
func expandWorkflowTemplatePlacementManagedClusterConfigLifecycleConfig(c *Client, f *WorkflowTemplatePlacementManagedClusterConfigLifecycleConfig, res *WorkflowTemplate) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.IdleDeleteTtl; !dcl.IsEmptyValueIndirect(v) {
		m["idleDeleteTtl"] = v
	}
	if v := f.AutoDeleteTime; !dcl.IsEmptyValueIndirect(v) {
		m["autoDeleteTime"] = v
	}
	if v := f.AutoDeleteTtl; !dcl.IsEmptyValueIndirect(v) {
		m["autoDeleteTtl"] = v
	}

	return m, nil
}

// flattenWorkflowTemplatePlacementManagedClusterConfigLifecycleConfig flattens an instance of WorkflowTemplatePlacementManagedClusterConfigLifecycleConfig from a JSON
// response object.
func flattenWorkflowTemplatePlacementManagedClusterConfigLifecycleConfig(c *Client, i interface{}, res *WorkflowTemplate) *WorkflowTemplatePlacementManagedClusterConfigLifecycleConfig {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &WorkflowTemplatePlacementManagedClusterConfigLifecycleConfig{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyWorkflowTemplatePlacementManagedClusterConfigLifecycleConfig
	}
	r.IdleDeleteTtl = dcl.FlattenString(m["idleDeleteTtl"])
	r.AutoDeleteTime = dcl.FlattenString(m["autoDeleteTime"])
	r.AutoDeleteTtl = dcl.FlattenString(m["autoDeleteTtl"])
	r.IdleStartTime = dcl.FlattenString(m["idleStartTime"])

	return r
}

// expandWorkflowTemplatePlacementManagedClusterConfigEndpointConfigMap expands the contents of WorkflowTemplatePlacementManagedClusterConfigEndpointConfig into a JSON
// request object.
func expandWorkflowTemplatePlacementManagedClusterConfigEndpointConfigMap(c *Client, f map[string]WorkflowTemplatePlacementManagedClusterConfigEndpointConfig, res *WorkflowTemplate) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandWorkflowTemplatePlacementManagedClusterConfigEndpointConfig(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandWorkflowTemplatePlacementManagedClusterConfigEndpointConfigSlice expands the contents of WorkflowTemplatePlacementManagedClusterConfigEndpointConfig into a JSON
// request object.
func expandWorkflowTemplatePlacementManagedClusterConfigEndpointConfigSlice(c *Client, f []WorkflowTemplatePlacementManagedClusterConfigEndpointConfig, res *WorkflowTemplate) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandWorkflowTemplatePlacementManagedClusterConfigEndpointConfig(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenWorkflowTemplatePlacementManagedClusterConfigEndpointConfigMap flattens the contents of WorkflowTemplatePlacementManagedClusterConfigEndpointConfig from a JSON
// response object.
func flattenWorkflowTemplatePlacementManagedClusterConfigEndpointConfigMap(c *Client, i interface{}, res *WorkflowTemplate) map[string]WorkflowTemplatePlacementManagedClusterConfigEndpointConfig {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]WorkflowTemplatePlacementManagedClusterConfigEndpointConfig{}
	}

	if len(a) == 0 {
		return map[string]WorkflowTemplatePlacementManagedClusterConfigEndpointConfig{}
	}

	items := make(map[string]WorkflowTemplatePlacementManagedClusterConfigEndpointConfig)
	for k, item := range a {
		items[k] = *flattenWorkflowTemplatePlacementManagedClusterConfigEndpointConfig(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenWorkflowTemplatePlacementManagedClusterConfigEndpointConfigSlice flattens the contents of WorkflowTemplatePlacementManagedClusterConfigEndpointConfig from a JSON
// response object.
func flattenWorkflowTemplatePlacementManagedClusterConfigEndpointConfigSlice(c *Client, i interface{}, res *WorkflowTemplate) []WorkflowTemplatePlacementManagedClusterConfigEndpointConfig {
	a, ok := i.([]interface{})
	if !ok {
		return []WorkflowTemplatePlacementManagedClusterConfigEndpointConfig{}
	}

	if len(a) == 0 {
		return []WorkflowTemplatePlacementManagedClusterConfigEndpointConfig{}
	}

	items := make([]WorkflowTemplatePlacementManagedClusterConfigEndpointConfig, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenWorkflowTemplatePlacementManagedClusterConfigEndpointConfig(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandWorkflowTemplatePlacementManagedClusterConfigEndpointConfig expands an instance of WorkflowTemplatePlacementManagedClusterConfigEndpointConfig into a JSON
// request object.
func expandWorkflowTemplatePlacementManagedClusterConfigEndpointConfig(c *Client, f *WorkflowTemplatePlacementManagedClusterConfigEndpointConfig, res *WorkflowTemplate) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.EnableHttpPortAccess; !dcl.IsEmptyValueIndirect(v) {
		m["enableHttpPortAccess"] = v
	}

	return m, nil
}

// flattenWorkflowTemplatePlacementManagedClusterConfigEndpointConfig flattens an instance of WorkflowTemplatePlacementManagedClusterConfigEndpointConfig from a JSON
// response object.
func flattenWorkflowTemplatePlacementManagedClusterConfigEndpointConfig(c *Client, i interface{}, res *WorkflowTemplate) *WorkflowTemplatePlacementManagedClusterConfigEndpointConfig {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &WorkflowTemplatePlacementManagedClusterConfigEndpointConfig{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyWorkflowTemplatePlacementManagedClusterConfigEndpointConfig
	}
	r.HttpPorts = dcl.FlattenKeyValuePairs(m["httpPorts"])
	r.EnableHttpPortAccess = dcl.FlattenBool(m["enableHttpPortAccess"])

	return r
}

// expandWorkflowTemplatePlacementClusterSelectorMap expands the contents of WorkflowTemplatePlacementClusterSelector into a JSON
// request object.
func expandWorkflowTemplatePlacementClusterSelectorMap(c *Client, f map[string]WorkflowTemplatePlacementClusterSelector, res *WorkflowTemplate) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandWorkflowTemplatePlacementClusterSelector(c, &item, res)
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
func expandWorkflowTemplatePlacementClusterSelectorSlice(c *Client, f []WorkflowTemplatePlacementClusterSelector, res *WorkflowTemplate) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandWorkflowTemplatePlacementClusterSelector(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenWorkflowTemplatePlacementClusterSelectorMap flattens the contents of WorkflowTemplatePlacementClusterSelector from a JSON
// response object.
func flattenWorkflowTemplatePlacementClusterSelectorMap(c *Client, i interface{}, res *WorkflowTemplate) map[string]WorkflowTemplatePlacementClusterSelector {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]WorkflowTemplatePlacementClusterSelector{}
	}

	if len(a) == 0 {
		return map[string]WorkflowTemplatePlacementClusterSelector{}
	}

	items := make(map[string]WorkflowTemplatePlacementClusterSelector)
	for k, item := range a {
		items[k] = *flattenWorkflowTemplatePlacementClusterSelector(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenWorkflowTemplatePlacementClusterSelectorSlice flattens the contents of WorkflowTemplatePlacementClusterSelector from a JSON
// response object.
func flattenWorkflowTemplatePlacementClusterSelectorSlice(c *Client, i interface{}, res *WorkflowTemplate) []WorkflowTemplatePlacementClusterSelector {
	a, ok := i.([]interface{})
	if !ok {
		return []WorkflowTemplatePlacementClusterSelector{}
	}

	if len(a) == 0 {
		return []WorkflowTemplatePlacementClusterSelector{}
	}

	items := make([]WorkflowTemplatePlacementClusterSelector, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenWorkflowTemplatePlacementClusterSelector(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandWorkflowTemplatePlacementClusterSelector expands an instance of WorkflowTemplatePlacementClusterSelector into a JSON
// request object.
func expandWorkflowTemplatePlacementClusterSelector(c *Client, f *WorkflowTemplatePlacementClusterSelector, res *WorkflowTemplate) (map[string]interface{}, error) {
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
func flattenWorkflowTemplatePlacementClusterSelector(c *Client, i interface{}, res *WorkflowTemplate) *WorkflowTemplatePlacementClusterSelector {
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
func expandWorkflowTemplateJobsMap(c *Client, f map[string]WorkflowTemplateJobs, res *WorkflowTemplate) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandWorkflowTemplateJobs(c, &item, res)
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
func expandWorkflowTemplateJobsSlice(c *Client, f []WorkflowTemplateJobs, res *WorkflowTemplate) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandWorkflowTemplateJobs(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenWorkflowTemplateJobsMap flattens the contents of WorkflowTemplateJobs from a JSON
// response object.
func flattenWorkflowTemplateJobsMap(c *Client, i interface{}, res *WorkflowTemplate) map[string]WorkflowTemplateJobs {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]WorkflowTemplateJobs{}
	}

	if len(a) == 0 {
		return map[string]WorkflowTemplateJobs{}
	}

	items := make(map[string]WorkflowTemplateJobs)
	for k, item := range a {
		items[k] = *flattenWorkflowTemplateJobs(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenWorkflowTemplateJobsSlice flattens the contents of WorkflowTemplateJobs from a JSON
// response object.
func flattenWorkflowTemplateJobsSlice(c *Client, i interface{}, res *WorkflowTemplate) []WorkflowTemplateJobs {
	a, ok := i.([]interface{})
	if !ok {
		return []WorkflowTemplateJobs{}
	}

	if len(a) == 0 {
		return []WorkflowTemplateJobs{}
	}

	items := make([]WorkflowTemplateJobs, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenWorkflowTemplateJobs(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandWorkflowTemplateJobs expands an instance of WorkflowTemplateJobs into a JSON
// request object.
func expandWorkflowTemplateJobs(c *Client, f *WorkflowTemplateJobs, res *WorkflowTemplate) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.StepId; !dcl.IsEmptyValueIndirect(v) {
		m["stepId"] = v
	}
	if v, err := expandWorkflowTemplateJobsHadoopJob(c, f.HadoopJob, res); err != nil {
		return nil, fmt.Errorf("error expanding HadoopJob into hadoopJob: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["hadoopJob"] = v
	}
	if v, err := expandWorkflowTemplateJobsSparkJob(c, f.SparkJob, res); err != nil {
		return nil, fmt.Errorf("error expanding SparkJob into sparkJob: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["sparkJob"] = v
	}
	if v, err := expandWorkflowTemplateJobsPysparkJob(c, f.PysparkJob, res); err != nil {
		return nil, fmt.Errorf("error expanding PysparkJob into pysparkJob: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["pysparkJob"] = v
	}
	if v, err := expandWorkflowTemplateJobsHiveJob(c, f.HiveJob, res); err != nil {
		return nil, fmt.Errorf("error expanding HiveJob into hiveJob: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["hiveJob"] = v
	}
	if v, err := expandWorkflowTemplateJobsPigJob(c, f.PigJob, res); err != nil {
		return nil, fmt.Errorf("error expanding PigJob into pigJob: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["pigJob"] = v
	}
	if v, err := expandWorkflowTemplateJobsSparkRJob(c, f.SparkRJob, res); err != nil {
		return nil, fmt.Errorf("error expanding SparkRJob into sparkRJob: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["sparkRJob"] = v
	}
	if v, err := expandWorkflowTemplateJobsSparkSqlJob(c, f.SparkSqlJob, res); err != nil {
		return nil, fmt.Errorf("error expanding SparkSqlJob into sparkSqlJob: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["sparkSqlJob"] = v
	}
	if v, err := expandWorkflowTemplateJobsPrestoJob(c, f.PrestoJob, res); err != nil {
		return nil, fmt.Errorf("error expanding PrestoJob into prestoJob: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["prestoJob"] = v
	}
	if v := f.Labels; !dcl.IsEmptyValueIndirect(v) {
		m["labels"] = v
	}
	if v, err := expandWorkflowTemplateJobsScheduling(c, f.Scheduling, res); err != nil {
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
func flattenWorkflowTemplateJobs(c *Client, i interface{}, res *WorkflowTemplate) *WorkflowTemplateJobs {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &WorkflowTemplateJobs{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyWorkflowTemplateJobs
	}
	r.StepId = dcl.FlattenString(m["stepId"])
	r.HadoopJob = flattenWorkflowTemplateJobsHadoopJob(c, m["hadoopJob"], res)
	r.SparkJob = flattenWorkflowTemplateJobsSparkJob(c, m["sparkJob"], res)
	r.PysparkJob = flattenWorkflowTemplateJobsPysparkJob(c, m["pysparkJob"], res)
	r.HiveJob = flattenWorkflowTemplateJobsHiveJob(c, m["hiveJob"], res)
	r.PigJob = flattenWorkflowTemplateJobsPigJob(c, m["pigJob"], res)
	r.SparkRJob = flattenWorkflowTemplateJobsSparkRJob(c, m["sparkRJob"], res)
	r.SparkSqlJob = flattenWorkflowTemplateJobsSparkSqlJob(c, m["sparkSqlJob"], res)
	r.PrestoJob = flattenWorkflowTemplateJobsPrestoJob(c, m["prestoJob"], res)
	r.Labels = dcl.FlattenKeyValuePairs(m["labels"])
	r.Scheduling = flattenWorkflowTemplateJobsScheduling(c, m["scheduling"], res)
	r.PrerequisiteStepIds = dcl.FlattenStringSlice(m["prerequisiteStepIds"])

	return r
}

// expandWorkflowTemplateJobsHadoopJobMap expands the contents of WorkflowTemplateJobsHadoopJob into a JSON
// request object.
func expandWorkflowTemplateJobsHadoopJobMap(c *Client, f map[string]WorkflowTemplateJobsHadoopJob, res *WorkflowTemplate) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandWorkflowTemplateJobsHadoopJob(c, &item, res)
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
func expandWorkflowTemplateJobsHadoopJobSlice(c *Client, f []WorkflowTemplateJobsHadoopJob, res *WorkflowTemplate) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandWorkflowTemplateJobsHadoopJob(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenWorkflowTemplateJobsHadoopJobMap flattens the contents of WorkflowTemplateJobsHadoopJob from a JSON
// response object.
func flattenWorkflowTemplateJobsHadoopJobMap(c *Client, i interface{}, res *WorkflowTemplate) map[string]WorkflowTemplateJobsHadoopJob {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]WorkflowTemplateJobsHadoopJob{}
	}

	if len(a) == 0 {
		return map[string]WorkflowTemplateJobsHadoopJob{}
	}

	items := make(map[string]WorkflowTemplateJobsHadoopJob)
	for k, item := range a {
		items[k] = *flattenWorkflowTemplateJobsHadoopJob(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenWorkflowTemplateJobsHadoopJobSlice flattens the contents of WorkflowTemplateJobsHadoopJob from a JSON
// response object.
func flattenWorkflowTemplateJobsHadoopJobSlice(c *Client, i interface{}, res *WorkflowTemplate) []WorkflowTemplateJobsHadoopJob {
	a, ok := i.([]interface{})
	if !ok {
		return []WorkflowTemplateJobsHadoopJob{}
	}

	if len(a) == 0 {
		return []WorkflowTemplateJobsHadoopJob{}
	}

	items := make([]WorkflowTemplateJobsHadoopJob, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenWorkflowTemplateJobsHadoopJob(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandWorkflowTemplateJobsHadoopJob expands an instance of WorkflowTemplateJobsHadoopJob into a JSON
// request object.
func expandWorkflowTemplateJobsHadoopJob(c *Client, f *WorkflowTemplateJobsHadoopJob, res *WorkflowTemplate) (map[string]interface{}, error) {
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
	if v, err := expandWorkflowTemplateJobsHadoopJobLoggingConfig(c, f.LoggingConfig, res); err != nil {
		return nil, fmt.Errorf("error expanding LoggingConfig into loggingConfig: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["loggingConfig"] = v
	}

	return m, nil
}

// flattenWorkflowTemplateJobsHadoopJob flattens an instance of WorkflowTemplateJobsHadoopJob from a JSON
// response object.
func flattenWorkflowTemplateJobsHadoopJob(c *Client, i interface{}, res *WorkflowTemplate) *WorkflowTemplateJobsHadoopJob {
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
	r.LoggingConfig = flattenWorkflowTemplateJobsHadoopJobLoggingConfig(c, m["loggingConfig"], res)

	return r
}

// expandWorkflowTemplateJobsHadoopJobLoggingConfigMap expands the contents of WorkflowTemplateJobsHadoopJobLoggingConfig into a JSON
// request object.
func expandWorkflowTemplateJobsHadoopJobLoggingConfigMap(c *Client, f map[string]WorkflowTemplateJobsHadoopJobLoggingConfig, res *WorkflowTemplate) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandWorkflowTemplateJobsHadoopJobLoggingConfig(c, &item, res)
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
func expandWorkflowTemplateJobsHadoopJobLoggingConfigSlice(c *Client, f []WorkflowTemplateJobsHadoopJobLoggingConfig, res *WorkflowTemplate) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandWorkflowTemplateJobsHadoopJobLoggingConfig(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenWorkflowTemplateJobsHadoopJobLoggingConfigMap flattens the contents of WorkflowTemplateJobsHadoopJobLoggingConfig from a JSON
// response object.
func flattenWorkflowTemplateJobsHadoopJobLoggingConfigMap(c *Client, i interface{}, res *WorkflowTemplate) map[string]WorkflowTemplateJobsHadoopJobLoggingConfig {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]WorkflowTemplateJobsHadoopJobLoggingConfig{}
	}

	if len(a) == 0 {
		return map[string]WorkflowTemplateJobsHadoopJobLoggingConfig{}
	}

	items := make(map[string]WorkflowTemplateJobsHadoopJobLoggingConfig)
	for k, item := range a {
		items[k] = *flattenWorkflowTemplateJobsHadoopJobLoggingConfig(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenWorkflowTemplateJobsHadoopJobLoggingConfigSlice flattens the contents of WorkflowTemplateJobsHadoopJobLoggingConfig from a JSON
// response object.
func flattenWorkflowTemplateJobsHadoopJobLoggingConfigSlice(c *Client, i interface{}, res *WorkflowTemplate) []WorkflowTemplateJobsHadoopJobLoggingConfig {
	a, ok := i.([]interface{})
	if !ok {
		return []WorkflowTemplateJobsHadoopJobLoggingConfig{}
	}

	if len(a) == 0 {
		return []WorkflowTemplateJobsHadoopJobLoggingConfig{}
	}

	items := make([]WorkflowTemplateJobsHadoopJobLoggingConfig, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenWorkflowTemplateJobsHadoopJobLoggingConfig(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandWorkflowTemplateJobsHadoopJobLoggingConfig expands an instance of WorkflowTemplateJobsHadoopJobLoggingConfig into a JSON
// request object.
func expandWorkflowTemplateJobsHadoopJobLoggingConfig(c *Client, f *WorkflowTemplateJobsHadoopJobLoggingConfig, res *WorkflowTemplate) (map[string]interface{}, error) {
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
func flattenWorkflowTemplateJobsHadoopJobLoggingConfig(c *Client, i interface{}, res *WorkflowTemplate) *WorkflowTemplateJobsHadoopJobLoggingConfig {
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
func expandWorkflowTemplateJobsSparkJobMap(c *Client, f map[string]WorkflowTemplateJobsSparkJob, res *WorkflowTemplate) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandWorkflowTemplateJobsSparkJob(c, &item, res)
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
func expandWorkflowTemplateJobsSparkJobSlice(c *Client, f []WorkflowTemplateJobsSparkJob, res *WorkflowTemplate) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandWorkflowTemplateJobsSparkJob(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenWorkflowTemplateJobsSparkJobMap flattens the contents of WorkflowTemplateJobsSparkJob from a JSON
// response object.
func flattenWorkflowTemplateJobsSparkJobMap(c *Client, i interface{}, res *WorkflowTemplate) map[string]WorkflowTemplateJobsSparkJob {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]WorkflowTemplateJobsSparkJob{}
	}

	if len(a) == 0 {
		return map[string]WorkflowTemplateJobsSparkJob{}
	}

	items := make(map[string]WorkflowTemplateJobsSparkJob)
	for k, item := range a {
		items[k] = *flattenWorkflowTemplateJobsSparkJob(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenWorkflowTemplateJobsSparkJobSlice flattens the contents of WorkflowTemplateJobsSparkJob from a JSON
// response object.
func flattenWorkflowTemplateJobsSparkJobSlice(c *Client, i interface{}, res *WorkflowTemplate) []WorkflowTemplateJobsSparkJob {
	a, ok := i.([]interface{})
	if !ok {
		return []WorkflowTemplateJobsSparkJob{}
	}

	if len(a) == 0 {
		return []WorkflowTemplateJobsSparkJob{}
	}

	items := make([]WorkflowTemplateJobsSparkJob, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenWorkflowTemplateJobsSparkJob(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandWorkflowTemplateJobsSparkJob expands an instance of WorkflowTemplateJobsSparkJob into a JSON
// request object.
func expandWorkflowTemplateJobsSparkJob(c *Client, f *WorkflowTemplateJobsSparkJob, res *WorkflowTemplate) (map[string]interface{}, error) {
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
	if v, err := expandWorkflowTemplateJobsSparkJobLoggingConfig(c, f.LoggingConfig, res); err != nil {
		return nil, fmt.Errorf("error expanding LoggingConfig into loggingConfig: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["loggingConfig"] = v
	}

	return m, nil
}

// flattenWorkflowTemplateJobsSparkJob flattens an instance of WorkflowTemplateJobsSparkJob from a JSON
// response object.
func flattenWorkflowTemplateJobsSparkJob(c *Client, i interface{}, res *WorkflowTemplate) *WorkflowTemplateJobsSparkJob {
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
	r.LoggingConfig = flattenWorkflowTemplateJobsSparkJobLoggingConfig(c, m["loggingConfig"], res)

	return r
}

// expandWorkflowTemplateJobsSparkJobLoggingConfigMap expands the contents of WorkflowTemplateJobsSparkJobLoggingConfig into a JSON
// request object.
func expandWorkflowTemplateJobsSparkJobLoggingConfigMap(c *Client, f map[string]WorkflowTemplateJobsSparkJobLoggingConfig, res *WorkflowTemplate) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandWorkflowTemplateJobsSparkJobLoggingConfig(c, &item, res)
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
func expandWorkflowTemplateJobsSparkJobLoggingConfigSlice(c *Client, f []WorkflowTemplateJobsSparkJobLoggingConfig, res *WorkflowTemplate) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandWorkflowTemplateJobsSparkJobLoggingConfig(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenWorkflowTemplateJobsSparkJobLoggingConfigMap flattens the contents of WorkflowTemplateJobsSparkJobLoggingConfig from a JSON
// response object.
func flattenWorkflowTemplateJobsSparkJobLoggingConfigMap(c *Client, i interface{}, res *WorkflowTemplate) map[string]WorkflowTemplateJobsSparkJobLoggingConfig {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]WorkflowTemplateJobsSparkJobLoggingConfig{}
	}

	if len(a) == 0 {
		return map[string]WorkflowTemplateJobsSparkJobLoggingConfig{}
	}

	items := make(map[string]WorkflowTemplateJobsSparkJobLoggingConfig)
	for k, item := range a {
		items[k] = *flattenWorkflowTemplateJobsSparkJobLoggingConfig(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenWorkflowTemplateJobsSparkJobLoggingConfigSlice flattens the contents of WorkflowTemplateJobsSparkJobLoggingConfig from a JSON
// response object.
func flattenWorkflowTemplateJobsSparkJobLoggingConfigSlice(c *Client, i interface{}, res *WorkflowTemplate) []WorkflowTemplateJobsSparkJobLoggingConfig {
	a, ok := i.([]interface{})
	if !ok {
		return []WorkflowTemplateJobsSparkJobLoggingConfig{}
	}

	if len(a) == 0 {
		return []WorkflowTemplateJobsSparkJobLoggingConfig{}
	}

	items := make([]WorkflowTemplateJobsSparkJobLoggingConfig, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenWorkflowTemplateJobsSparkJobLoggingConfig(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandWorkflowTemplateJobsSparkJobLoggingConfig expands an instance of WorkflowTemplateJobsSparkJobLoggingConfig into a JSON
// request object.
func expandWorkflowTemplateJobsSparkJobLoggingConfig(c *Client, f *WorkflowTemplateJobsSparkJobLoggingConfig, res *WorkflowTemplate) (map[string]interface{}, error) {
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
func flattenWorkflowTemplateJobsSparkJobLoggingConfig(c *Client, i interface{}, res *WorkflowTemplate) *WorkflowTemplateJobsSparkJobLoggingConfig {
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
func expandWorkflowTemplateJobsPysparkJobMap(c *Client, f map[string]WorkflowTemplateJobsPysparkJob, res *WorkflowTemplate) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandWorkflowTemplateJobsPysparkJob(c, &item, res)
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
func expandWorkflowTemplateJobsPysparkJobSlice(c *Client, f []WorkflowTemplateJobsPysparkJob, res *WorkflowTemplate) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandWorkflowTemplateJobsPysparkJob(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenWorkflowTemplateJobsPysparkJobMap flattens the contents of WorkflowTemplateJobsPysparkJob from a JSON
// response object.
func flattenWorkflowTemplateJobsPysparkJobMap(c *Client, i interface{}, res *WorkflowTemplate) map[string]WorkflowTemplateJobsPysparkJob {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]WorkflowTemplateJobsPysparkJob{}
	}

	if len(a) == 0 {
		return map[string]WorkflowTemplateJobsPysparkJob{}
	}

	items := make(map[string]WorkflowTemplateJobsPysparkJob)
	for k, item := range a {
		items[k] = *flattenWorkflowTemplateJobsPysparkJob(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenWorkflowTemplateJobsPysparkJobSlice flattens the contents of WorkflowTemplateJobsPysparkJob from a JSON
// response object.
func flattenWorkflowTemplateJobsPysparkJobSlice(c *Client, i interface{}, res *WorkflowTemplate) []WorkflowTemplateJobsPysparkJob {
	a, ok := i.([]interface{})
	if !ok {
		return []WorkflowTemplateJobsPysparkJob{}
	}

	if len(a) == 0 {
		return []WorkflowTemplateJobsPysparkJob{}
	}

	items := make([]WorkflowTemplateJobsPysparkJob, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenWorkflowTemplateJobsPysparkJob(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandWorkflowTemplateJobsPysparkJob expands an instance of WorkflowTemplateJobsPysparkJob into a JSON
// request object.
func expandWorkflowTemplateJobsPysparkJob(c *Client, f *WorkflowTemplateJobsPysparkJob, res *WorkflowTemplate) (map[string]interface{}, error) {
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
	if v, err := expandWorkflowTemplateJobsPysparkJobLoggingConfig(c, f.LoggingConfig, res); err != nil {
		return nil, fmt.Errorf("error expanding LoggingConfig into loggingConfig: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["loggingConfig"] = v
	}

	return m, nil
}

// flattenWorkflowTemplateJobsPysparkJob flattens an instance of WorkflowTemplateJobsPysparkJob from a JSON
// response object.
func flattenWorkflowTemplateJobsPysparkJob(c *Client, i interface{}, res *WorkflowTemplate) *WorkflowTemplateJobsPysparkJob {
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
	r.LoggingConfig = flattenWorkflowTemplateJobsPysparkJobLoggingConfig(c, m["loggingConfig"], res)

	return r
}

// expandWorkflowTemplateJobsPysparkJobLoggingConfigMap expands the contents of WorkflowTemplateJobsPysparkJobLoggingConfig into a JSON
// request object.
func expandWorkflowTemplateJobsPysparkJobLoggingConfigMap(c *Client, f map[string]WorkflowTemplateJobsPysparkJobLoggingConfig, res *WorkflowTemplate) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandWorkflowTemplateJobsPysparkJobLoggingConfig(c, &item, res)
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
func expandWorkflowTemplateJobsPysparkJobLoggingConfigSlice(c *Client, f []WorkflowTemplateJobsPysparkJobLoggingConfig, res *WorkflowTemplate) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandWorkflowTemplateJobsPysparkJobLoggingConfig(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenWorkflowTemplateJobsPysparkJobLoggingConfigMap flattens the contents of WorkflowTemplateJobsPysparkJobLoggingConfig from a JSON
// response object.
func flattenWorkflowTemplateJobsPysparkJobLoggingConfigMap(c *Client, i interface{}, res *WorkflowTemplate) map[string]WorkflowTemplateJobsPysparkJobLoggingConfig {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]WorkflowTemplateJobsPysparkJobLoggingConfig{}
	}

	if len(a) == 0 {
		return map[string]WorkflowTemplateJobsPysparkJobLoggingConfig{}
	}

	items := make(map[string]WorkflowTemplateJobsPysparkJobLoggingConfig)
	for k, item := range a {
		items[k] = *flattenWorkflowTemplateJobsPysparkJobLoggingConfig(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenWorkflowTemplateJobsPysparkJobLoggingConfigSlice flattens the contents of WorkflowTemplateJobsPysparkJobLoggingConfig from a JSON
// response object.
func flattenWorkflowTemplateJobsPysparkJobLoggingConfigSlice(c *Client, i interface{}, res *WorkflowTemplate) []WorkflowTemplateJobsPysparkJobLoggingConfig {
	a, ok := i.([]interface{})
	if !ok {
		return []WorkflowTemplateJobsPysparkJobLoggingConfig{}
	}

	if len(a) == 0 {
		return []WorkflowTemplateJobsPysparkJobLoggingConfig{}
	}

	items := make([]WorkflowTemplateJobsPysparkJobLoggingConfig, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenWorkflowTemplateJobsPysparkJobLoggingConfig(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandWorkflowTemplateJobsPysparkJobLoggingConfig expands an instance of WorkflowTemplateJobsPysparkJobLoggingConfig into a JSON
// request object.
func expandWorkflowTemplateJobsPysparkJobLoggingConfig(c *Client, f *WorkflowTemplateJobsPysparkJobLoggingConfig, res *WorkflowTemplate) (map[string]interface{}, error) {
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
func flattenWorkflowTemplateJobsPysparkJobLoggingConfig(c *Client, i interface{}, res *WorkflowTemplate) *WorkflowTemplateJobsPysparkJobLoggingConfig {
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
func expandWorkflowTemplateJobsHiveJobMap(c *Client, f map[string]WorkflowTemplateJobsHiveJob, res *WorkflowTemplate) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandWorkflowTemplateJobsHiveJob(c, &item, res)
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
func expandWorkflowTemplateJobsHiveJobSlice(c *Client, f []WorkflowTemplateJobsHiveJob, res *WorkflowTemplate) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandWorkflowTemplateJobsHiveJob(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenWorkflowTemplateJobsHiveJobMap flattens the contents of WorkflowTemplateJobsHiveJob from a JSON
// response object.
func flattenWorkflowTemplateJobsHiveJobMap(c *Client, i interface{}, res *WorkflowTemplate) map[string]WorkflowTemplateJobsHiveJob {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]WorkflowTemplateJobsHiveJob{}
	}

	if len(a) == 0 {
		return map[string]WorkflowTemplateJobsHiveJob{}
	}

	items := make(map[string]WorkflowTemplateJobsHiveJob)
	for k, item := range a {
		items[k] = *flattenWorkflowTemplateJobsHiveJob(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenWorkflowTemplateJobsHiveJobSlice flattens the contents of WorkflowTemplateJobsHiveJob from a JSON
// response object.
func flattenWorkflowTemplateJobsHiveJobSlice(c *Client, i interface{}, res *WorkflowTemplate) []WorkflowTemplateJobsHiveJob {
	a, ok := i.([]interface{})
	if !ok {
		return []WorkflowTemplateJobsHiveJob{}
	}

	if len(a) == 0 {
		return []WorkflowTemplateJobsHiveJob{}
	}

	items := make([]WorkflowTemplateJobsHiveJob, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenWorkflowTemplateJobsHiveJob(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandWorkflowTemplateJobsHiveJob expands an instance of WorkflowTemplateJobsHiveJob into a JSON
// request object.
func expandWorkflowTemplateJobsHiveJob(c *Client, f *WorkflowTemplateJobsHiveJob, res *WorkflowTemplate) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.QueryFileUri; !dcl.IsEmptyValueIndirect(v) {
		m["queryFileUri"] = v
	}
	if v, err := expandWorkflowTemplateJobsHiveJobQueryList(c, f.QueryList, res); err != nil {
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
func flattenWorkflowTemplateJobsHiveJob(c *Client, i interface{}, res *WorkflowTemplate) *WorkflowTemplateJobsHiveJob {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &WorkflowTemplateJobsHiveJob{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyWorkflowTemplateJobsHiveJob
	}
	r.QueryFileUri = dcl.FlattenString(m["queryFileUri"])
	r.QueryList = flattenWorkflowTemplateJobsHiveJobQueryList(c, m["queryList"], res)
	r.ContinueOnFailure = dcl.FlattenBool(m["continueOnFailure"])
	r.ScriptVariables = dcl.FlattenKeyValuePairs(m["scriptVariables"])
	r.Properties = dcl.FlattenKeyValuePairs(m["properties"])
	r.JarFileUris = dcl.FlattenStringSlice(m["jarFileUris"])

	return r
}

// expandWorkflowTemplateJobsHiveJobQueryListMap expands the contents of WorkflowTemplateJobsHiveJobQueryList into a JSON
// request object.
func expandWorkflowTemplateJobsHiveJobQueryListMap(c *Client, f map[string]WorkflowTemplateJobsHiveJobQueryList, res *WorkflowTemplate) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandWorkflowTemplateJobsHiveJobQueryList(c, &item, res)
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
func expandWorkflowTemplateJobsHiveJobQueryListSlice(c *Client, f []WorkflowTemplateJobsHiveJobQueryList, res *WorkflowTemplate) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandWorkflowTemplateJobsHiveJobQueryList(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenWorkflowTemplateJobsHiveJobQueryListMap flattens the contents of WorkflowTemplateJobsHiveJobQueryList from a JSON
// response object.
func flattenWorkflowTemplateJobsHiveJobQueryListMap(c *Client, i interface{}, res *WorkflowTemplate) map[string]WorkflowTemplateJobsHiveJobQueryList {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]WorkflowTemplateJobsHiveJobQueryList{}
	}

	if len(a) == 0 {
		return map[string]WorkflowTemplateJobsHiveJobQueryList{}
	}

	items := make(map[string]WorkflowTemplateJobsHiveJobQueryList)
	for k, item := range a {
		items[k] = *flattenWorkflowTemplateJobsHiveJobQueryList(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenWorkflowTemplateJobsHiveJobQueryListSlice flattens the contents of WorkflowTemplateJobsHiveJobQueryList from a JSON
// response object.
func flattenWorkflowTemplateJobsHiveJobQueryListSlice(c *Client, i interface{}, res *WorkflowTemplate) []WorkflowTemplateJobsHiveJobQueryList {
	a, ok := i.([]interface{})
	if !ok {
		return []WorkflowTemplateJobsHiveJobQueryList{}
	}

	if len(a) == 0 {
		return []WorkflowTemplateJobsHiveJobQueryList{}
	}

	items := make([]WorkflowTemplateJobsHiveJobQueryList, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenWorkflowTemplateJobsHiveJobQueryList(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandWorkflowTemplateJobsHiveJobQueryList expands an instance of WorkflowTemplateJobsHiveJobQueryList into a JSON
// request object.
func expandWorkflowTemplateJobsHiveJobQueryList(c *Client, f *WorkflowTemplateJobsHiveJobQueryList, res *WorkflowTemplate) (map[string]interface{}, error) {
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
func flattenWorkflowTemplateJobsHiveJobQueryList(c *Client, i interface{}, res *WorkflowTemplate) *WorkflowTemplateJobsHiveJobQueryList {
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
func expandWorkflowTemplateJobsPigJobMap(c *Client, f map[string]WorkflowTemplateJobsPigJob, res *WorkflowTemplate) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandWorkflowTemplateJobsPigJob(c, &item, res)
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
func expandWorkflowTemplateJobsPigJobSlice(c *Client, f []WorkflowTemplateJobsPigJob, res *WorkflowTemplate) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandWorkflowTemplateJobsPigJob(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenWorkflowTemplateJobsPigJobMap flattens the contents of WorkflowTemplateJobsPigJob from a JSON
// response object.
func flattenWorkflowTemplateJobsPigJobMap(c *Client, i interface{}, res *WorkflowTemplate) map[string]WorkflowTemplateJobsPigJob {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]WorkflowTemplateJobsPigJob{}
	}

	if len(a) == 0 {
		return map[string]WorkflowTemplateJobsPigJob{}
	}

	items := make(map[string]WorkflowTemplateJobsPigJob)
	for k, item := range a {
		items[k] = *flattenWorkflowTemplateJobsPigJob(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenWorkflowTemplateJobsPigJobSlice flattens the contents of WorkflowTemplateJobsPigJob from a JSON
// response object.
func flattenWorkflowTemplateJobsPigJobSlice(c *Client, i interface{}, res *WorkflowTemplate) []WorkflowTemplateJobsPigJob {
	a, ok := i.([]interface{})
	if !ok {
		return []WorkflowTemplateJobsPigJob{}
	}

	if len(a) == 0 {
		return []WorkflowTemplateJobsPigJob{}
	}

	items := make([]WorkflowTemplateJobsPigJob, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenWorkflowTemplateJobsPigJob(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandWorkflowTemplateJobsPigJob expands an instance of WorkflowTemplateJobsPigJob into a JSON
// request object.
func expandWorkflowTemplateJobsPigJob(c *Client, f *WorkflowTemplateJobsPigJob, res *WorkflowTemplate) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.QueryFileUri; !dcl.IsEmptyValueIndirect(v) {
		m["queryFileUri"] = v
	}
	if v, err := expandWorkflowTemplateJobsPigJobQueryList(c, f.QueryList, res); err != nil {
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
	if v, err := expandWorkflowTemplateJobsPigJobLoggingConfig(c, f.LoggingConfig, res); err != nil {
		return nil, fmt.Errorf("error expanding LoggingConfig into loggingConfig: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["loggingConfig"] = v
	}

	return m, nil
}

// flattenWorkflowTemplateJobsPigJob flattens an instance of WorkflowTemplateJobsPigJob from a JSON
// response object.
func flattenWorkflowTemplateJobsPigJob(c *Client, i interface{}, res *WorkflowTemplate) *WorkflowTemplateJobsPigJob {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &WorkflowTemplateJobsPigJob{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyWorkflowTemplateJobsPigJob
	}
	r.QueryFileUri = dcl.FlattenString(m["queryFileUri"])
	r.QueryList = flattenWorkflowTemplateJobsPigJobQueryList(c, m["queryList"], res)
	r.ContinueOnFailure = dcl.FlattenBool(m["continueOnFailure"])
	r.ScriptVariables = dcl.FlattenKeyValuePairs(m["scriptVariables"])
	r.Properties = dcl.FlattenKeyValuePairs(m["properties"])
	r.JarFileUris = dcl.FlattenStringSlice(m["jarFileUris"])
	r.LoggingConfig = flattenWorkflowTemplateJobsPigJobLoggingConfig(c, m["loggingConfig"], res)

	return r
}

// expandWorkflowTemplateJobsPigJobQueryListMap expands the contents of WorkflowTemplateJobsPigJobQueryList into a JSON
// request object.
func expandWorkflowTemplateJobsPigJobQueryListMap(c *Client, f map[string]WorkflowTemplateJobsPigJobQueryList, res *WorkflowTemplate) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandWorkflowTemplateJobsPigJobQueryList(c, &item, res)
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
func expandWorkflowTemplateJobsPigJobQueryListSlice(c *Client, f []WorkflowTemplateJobsPigJobQueryList, res *WorkflowTemplate) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandWorkflowTemplateJobsPigJobQueryList(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenWorkflowTemplateJobsPigJobQueryListMap flattens the contents of WorkflowTemplateJobsPigJobQueryList from a JSON
// response object.
func flattenWorkflowTemplateJobsPigJobQueryListMap(c *Client, i interface{}, res *WorkflowTemplate) map[string]WorkflowTemplateJobsPigJobQueryList {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]WorkflowTemplateJobsPigJobQueryList{}
	}

	if len(a) == 0 {
		return map[string]WorkflowTemplateJobsPigJobQueryList{}
	}

	items := make(map[string]WorkflowTemplateJobsPigJobQueryList)
	for k, item := range a {
		items[k] = *flattenWorkflowTemplateJobsPigJobQueryList(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenWorkflowTemplateJobsPigJobQueryListSlice flattens the contents of WorkflowTemplateJobsPigJobQueryList from a JSON
// response object.
func flattenWorkflowTemplateJobsPigJobQueryListSlice(c *Client, i interface{}, res *WorkflowTemplate) []WorkflowTemplateJobsPigJobQueryList {
	a, ok := i.([]interface{})
	if !ok {
		return []WorkflowTemplateJobsPigJobQueryList{}
	}

	if len(a) == 0 {
		return []WorkflowTemplateJobsPigJobQueryList{}
	}

	items := make([]WorkflowTemplateJobsPigJobQueryList, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenWorkflowTemplateJobsPigJobQueryList(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandWorkflowTemplateJobsPigJobQueryList expands an instance of WorkflowTemplateJobsPigJobQueryList into a JSON
// request object.
func expandWorkflowTemplateJobsPigJobQueryList(c *Client, f *WorkflowTemplateJobsPigJobQueryList, res *WorkflowTemplate) (map[string]interface{}, error) {
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
func flattenWorkflowTemplateJobsPigJobQueryList(c *Client, i interface{}, res *WorkflowTemplate) *WorkflowTemplateJobsPigJobQueryList {
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
func expandWorkflowTemplateJobsPigJobLoggingConfigMap(c *Client, f map[string]WorkflowTemplateJobsPigJobLoggingConfig, res *WorkflowTemplate) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandWorkflowTemplateJobsPigJobLoggingConfig(c, &item, res)
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
func expandWorkflowTemplateJobsPigJobLoggingConfigSlice(c *Client, f []WorkflowTemplateJobsPigJobLoggingConfig, res *WorkflowTemplate) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandWorkflowTemplateJobsPigJobLoggingConfig(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenWorkflowTemplateJobsPigJobLoggingConfigMap flattens the contents of WorkflowTemplateJobsPigJobLoggingConfig from a JSON
// response object.
func flattenWorkflowTemplateJobsPigJobLoggingConfigMap(c *Client, i interface{}, res *WorkflowTemplate) map[string]WorkflowTemplateJobsPigJobLoggingConfig {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]WorkflowTemplateJobsPigJobLoggingConfig{}
	}

	if len(a) == 0 {
		return map[string]WorkflowTemplateJobsPigJobLoggingConfig{}
	}

	items := make(map[string]WorkflowTemplateJobsPigJobLoggingConfig)
	for k, item := range a {
		items[k] = *flattenWorkflowTemplateJobsPigJobLoggingConfig(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenWorkflowTemplateJobsPigJobLoggingConfigSlice flattens the contents of WorkflowTemplateJobsPigJobLoggingConfig from a JSON
// response object.
func flattenWorkflowTemplateJobsPigJobLoggingConfigSlice(c *Client, i interface{}, res *WorkflowTemplate) []WorkflowTemplateJobsPigJobLoggingConfig {
	a, ok := i.([]interface{})
	if !ok {
		return []WorkflowTemplateJobsPigJobLoggingConfig{}
	}

	if len(a) == 0 {
		return []WorkflowTemplateJobsPigJobLoggingConfig{}
	}

	items := make([]WorkflowTemplateJobsPigJobLoggingConfig, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenWorkflowTemplateJobsPigJobLoggingConfig(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandWorkflowTemplateJobsPigJobLoggingConfig expands an instance of WorkflowTemplateJobsPigJobLoggingConfig into a JSON
// request object.
func expandWorkflowTemplateJobsPigJobLoggingConfig(c *Client, f *WorkflowTemplateJobsPigJobLoggingConfig, res *WorkflowTemplate) (map[string]interface{}, error) {
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
func flattenWorkflowTemplateJobsPigJobLoggingConfig(c *Client, i interface{}, res *WorkflowTemplate) *WorkflowTemplateJobsPigJobLoggingConfig {
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
func expandWorkflowTemplateJobsSparkRJobMap(c *Client, f map[string]WorkflowTemplateJobsSparkRJob, res *WorkflowTemplate) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandWorkflowTemplateJobsSparkRJob(c, &item, res)
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
func expandWorkflowTemplateJobsSparkRJobSlice(c *Client, f []WorkflowTemplateJobsSparkRJob, res *WorkflowTemplate) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandWorkflowTemplateJobsSparkRJob(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenWorkflowTemplateJobsSparkRJobMap flattens the contents of WorkflowTemplateJobsSparkRJob from a JSON
// response object.
func flattenWorkflowTemplateJobsSparkRJobMap(c *Client, i interface{}, res *WorkflowTemplate) map[string]WorkflowTemplateJobsSparkRJob {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]WorkflowTemplateJobsSparkRJob{}
	}

	if len(a) == 0 {
		return map[string]WorkflowTemplateJobsSparkRJob{}
	}

	items := make(map[string]WorkflowTemplateJobsSparkRJob)
	for k, item := range a {
		items[k] = *flattenWorkflowTemplateJobsSparkRJob(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenWorkflowTemplateJobsSparkRJobSlice flattens the contents of WorkflowTemplateJobsSparkRJob from a JSON
// response object.
func flattenWorkflowTemplateJobsSparkRJobSlice(c *Client, i interface{}, res *WorkflowTemplate) []WorkflowTemplateJobsSparkRJob {
	a, ok := i.([]interface{})
	if !ok {
		return []WorkflowTemplateJobsSparkRJob{}
	}

	if len(a) == 0 {
		return []WorkflowTemplateJobsSparkRJob{}
	}

	items := make([]WorkflowTemplateJobsSparkRJob, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenWorkflowTemplateJobsSparkRJob(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandWorkflowTemplateJobsSparkRJob expands an instance of WorkflowTemplateJobsSparkRJob into a JSON
// request object.
func expandWorkflowTemplateJobsSparkRJob(c *Client, f *WorkflowTemplateJobsSparkRJob, res *WorkflowTemplate) (map[string]interface{}, error) {
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
	if v, err := expandWorkflowTemplateJobsSparkRJobLoggingConfig(c, f.LoggingConfig, res); err != nil {
		return nil, fmt.Errorf("error expanding LoggingConfig into loggingConfig: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["loggingConfig"] = v
	}

	return m, nil
}

// flattenWorkflowTemplateJobsSparkRJob flattens an instance of WorkflowTemplateJobsSparkRJob from a JSON
// response object.
func flattenWorkflowTemplateJobsSparkRJob(c *Client, i interface{}, res *WorkflowTemplate) *WorkflowTemplateJobsSparkRJob {
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
	r.LoggingConfig = flattenWorkflowTemplateJobsSparkRJobLoggingConfig(c, m["loggingConfig"], res)

	return r
}

// expandWorkflowTemplateJobsSparkRJobLoggingConfigMap expands the contents of WorkflowTemplateJobsSparkRJobLoggingConfig into a JSON
// request object.
func expandWorkflowTemplateJobsSparkRJobLoggingConfigMap(c *Client, f map[string]WorkflowTemplateJobsSparkRJobLoggingConfig, res *WorkflowTemplate) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandWorkflowTemplateJobsSparkRJobLoggingConfig(c, &item, res)
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
func expandWorkflowTemplateJobsSparkRJobLoggingConfigSlice(c *Client, f []WorkflowTemplateJobsSparkRJobLoggingConfig, res *WorkflowTemplate) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandWorkflowTemplateJobsSparkRJobLoggingConfig(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenWorkflowTemplateJobsSparkRJobLoggingConfigMap flattens the contents of WorkflowTemplateJobsSparkRJobLoggingConfig from a JSON
// response object.
func flattenWorkflowTemplateJobsSparkRJobLoggingConfigMap(c *Client, i interface{}, res *WorkflowTemplate) map[string]WorkflowTemplateJobsSparkRJobLoggingConfig {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]WorkflowTemplateJobsSparkRJobLoggingConfig{}
	}

	if len(a) == 0 {
		return map[string]WorkflowTemplateJobsSparkRJobLoggingConfig{}
	}

	items := make(map[string]WorkflowTemplateJobsSparkRJobLoggingConfig)
	for k, item := range a {
		items[k] = *flattenWorkflowTemplateJobsSparkRJobLoggingConfig(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenWorkflowTemplateJobsSparkRJobLoggingConfigSlice flattens the contents of WorkflowTemplateJobsSparkRJobLoggingConfig from a JSON
// response object.
func flattenWorkflowTemplateJobsSparkRJobLoggingConfigSlice(c *Client, i interface{}, res *WorkflowTemplate) []WorkflowTemplateJobsSparkRJobLoggingConfig {
	a, ok := i.([]interface{})
	if !ok {
		return []WorkflowTemplateJobsSparkRJobLoggingConfig{}
	}

	if len(a) == 0 {
		return []WorkflowTemplateJobsSparkRJobLoggingConfig{}
	}

	items := make([]WorkflowTemplateJobsSparkRJobLoggingConfig, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenWorkflowTemplateJobsSparkRJobLoggingConfig(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandWorkflowTemplateJobsSparkRJobLoggingConfig expands an instance of WorkflowTemplateJobsSparkRJobLoggingConfig into a JSON
// request object.
func expandWorkflowTemplateJobsSparkRJobLoggingConfig(c *Client, f *WorkflowTemplateJobsSparkRJobLoggingConfig, res *WorkflowTemplate) (map[string]interface{}, error) {
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
func flattenWorkflowTemplateJobsSparkRJobLoggingConfig(c *Client, i interface{}, res *WorkflowTemplate) *WorkflowTemplateJobsSparkRJobLoggingConfig {
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
func expandWorkflowTemplateJobsSparkSqlJobMap(c *Client, f map[string]WorkflowTemplateJobsSparkSqlJob, res *WorkflowTemplate) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandWorkflowTemplateJobsSparkSqlJob(c, &item, res)
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
func expandWorkflowTemplateJobsSparkSqlJobSlice(c *Client, f []WorkflowTemplateJobsSparkSqlJob, res *WorkflowTemplate) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandWorkflowTemplateJobsSparkSqlJob(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenWorkflowTemplateJobsSparkSqlJobMap flattens the contents of WorkflowTemplateJobsSparkSqlJob from a JSON
// response object.
func flattenWorkflowTemplateJobsSparkSqlJobMap(c *Client, i interface{}, res *WorkflowTemplate) map[string]WorkflowTemplateJobsSparkSqlJob {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]WorkflowTemplateJobsSparkSqlJob{}
	}

	if len(a) == 0 {
		return map[string]WorkflowTemplateJobsSparkSqlJob{}
	}

	items := make(map[string]WorkflowTemplateJobsSparkSqlJob)
	for k, item := range a {
		items[k] = *flattenWorkflowTemplateJobsSparkSqlJob(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenWorkflowTemplateJobsSparkSqlJobSlice flattens the contents of WorkflowTemplateJobsSparkSqlJob from a JSON
// response object.
func flattenWorkflowTemplateJobsSparkSqlJobSlice(c *Client, i interface{}, res *WorkflowTemplate) []WorkflowTemplateJobsSparkSqlJob {
	a, ok := i.([]interface{})
	if !ok {
		return []WorkflowTemplateJobsSparkSqlJob{}
	}

	if len(a) == 0 {
		return []WorkflowTemplateJobsSparkSqlJob{}
	}

	items := make([]WorkflowTemplateJobsSparkSqlJob, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenWorkflowTemplateJobsSparkSqlJob(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandWorkflowTemplateJobsSparkSqlJob expands an instance of WorkflowTemplateJobsSparkSqlJob into a JSON
// request object.
func expandWorkflowTemplateJobsSparkSqlJob(c *Client, f *WorkflowTemplateJobsSparkSqlJob, res *WorkflowTemplate) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.QueryFileUri; !dcl.IsEmptyValueIndirect(v) {
		m["queryFileUri"] = v
	}
	if v, err := expandWorkflowTemplateJobsSparkSqlJobQueryList(c, f.QueryList, res); err != nil {
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
	if v, err := expandWorkflowTemplateJobsSparkSqlJobLoggingConfig(c, f.LoggingConfig, res); err != nil {
		return nil, fmt.Errorf("error expanding LoggingConfig into loggingConfig: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["loggingConfig"] = v
	}

	return m, nil
}

// flattenWorkflowTemplateJobsSparkSqlJob flattens an instance of WorkflowTemplateJobsSparkSqlJob from a JSON
// response object.
func flattenWorkflowTemplateJobsSparkSqlJob(c *Client, i interface{}, res *WorkflowTemplate) *WorkflowTemplateJobsSparkSqlJob {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &WorkflowTemplateJobsSparkSqlJob{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyWorkflowTemplateJobsSparkSqlJob
	}
	r.QueryFileUri = dcl.FlattenString(m["queryFileUri"])
	r.QueryList = flattenWorkflowTemplateJobsSparkSqlJobQueryList(c, m["queryList"], res)
	r.ScriptVariables = dcl.FlattenKeyValuePairs(m["scriptVariables"])
	r.Properties = dcl.FlattenKeyValuePairs(m["properties"])
	r.JarFileUris = dcl.FlattenStringSlice(m["jarFileUris"])
	r.LoggingConfig = flattenWorkflowTemplateJobsSparkSqlJobLoggingConfig(c, m["loggingConfig"], res)

	return r
}

// expandWorkflowTemplateJobsSparkSqlJobQueryListMap expands the contents of WorkflowTemplateJobsSparkSqlJobQueryList into a JSON
// request object.
func expandWorkflowTemplateJobsSparkSqlJobQueryListMap(c *Client, f map[string]WorkflowTemplateJobsSparkSqlJobQueryList, res *WorkflowTemplate) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandWorkflowTemplateJobsSparkSqlJobQueryList(c, &item, res)
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
func expandWorkflowTemplateJobsSparkSqlJobQueryListSlice(c *Client, f []WorkflowTemplateJobsSparkSqlJobQueryList, res *WorkflowTemplate) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandWorkflowTemplateJobsSparkSqlJobQueryList(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenWorkflowTemplateJobsSparkSqlJobQueryListMap flattens the contents of WorkflowTemplateJobsSparkSqlJobQueryList from a JSON
// response object.
func flattenWorkflowTemplateJobsSparkSqlJobQueryListMap(c *Client, i interface{}, res *WorkflowTemplate) map[string]WorkflowTemplateJobsSparkSqlJobQueryList {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]WorkflowTemplateJobsSparkSqlJobQueryList{}
	}

	if len(a) == 0 {
		return map[string]WorkflowTemplateJobsSparkSqlJobQueryList{}
	}

	items := make(map[string]WorkflowTemplateJobsSparkSqlJobQueryList)
	for k, item := range a {
		items[k] = *flattenWorkflowTemplateJobsSparkSqlJobQueryList(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenWorkflowTemplateJobsSparkSqlJobQueryListSlice flattens the contents of WorkflowTemplateJobsSparkSqlJobQueryList from a JSON
// response object.
func flattenWorkflowTemplateJobsSparkSqlJobQueryListSlice(c *Client, i interface{}, res *WorkflowTemplate) []WorkflowTemplateJobsSparkSqlJobQueryList {
	a, ok := i.([]interface{})
	if !ok {
		return []WorkflowTemplateJobsSparkSqlJobQueryList{}
	}

	if len(a) == 0 {
		return []WorkflowTemplateJobsSparkSqlJobQueryList{}
	}

	items := make([]WorkflowTemplateJobsSparkSqlJobQueryList, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenWorkflowTemplateJobsSparkSqlJobQueryList(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandWorkflowTemplateJobsSparkSqlJobQueryList expands an instance of WorkflowTemplateJobsSparkSqlJobQueryList into a JSON
// request object.
func expandWorkflowTemplateJobsSparkSqlJobQueryList(c *Client, f *WorkflowTemplateJobsSparkSqlJobQueryList, res *WorkflowTemplate) (map[string]interface{}, error) {
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
func flattenWorkflowTemplateJobsSparkSqlJobQueryList(c *Client, i interface{}, res *WorkflowTemplate) *WorkflowTemplateJobsSparkSqlJobQueryList {
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
func expandWorkflowTemplateJobsSparkSqlJobLoggingConfigMap(c *Client, f map[string]WorkflowTemplateJobsSparkSqlJobLoggingConfig, res *WorkflowTemplate) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandWorkflowTemplateJobsSparkSqlJobLoggingConfig(c, &item, res)
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
func expandWorkflowTemplateJobsSparkSqlJobLoggingConfigSlice(c *Client, f []WorkflowTemplateJobsSparkSqlJobLoggingConfig, res *WorkflowTemplate) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandWorkflowTemplateJobsSparkSqlJobLoggingConfig(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenWorkflowTemplateJobsSparkSqlJobLoggingConfigMap flattens the contents of WorkflowTemplateJobsSparkSqlJobLoggingConfig from a JSON
// response object.
func flattenWorkflowTemplateJobsSparkSqlJobLoggingConfigMap(c *Client, i interface{}, res *WorkflowTemplate) map[string]WorkflowTemplateJobsSparkSqlJobLoggingConfig {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]WorkflowTemplateJobsSparkSqlJobLoggingConfig{}
	}

	if len(a) == 0 {
		return map[string]WorkflowTemplateJobsSparkSqlJobLoggingConfig{}
	}

	items := make(map[string]WorkflowTemplateJobsSparkSqlJobLoggingConfig)
	for k, item := range a {
		items[k] = *flattenWorkflowTemplateJobsSparkSqlJobLoggingConfig(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenWorkflowTemplateJobsSparkSqlJobLoggingConfigSlice flattens the contents of WorkflowTemplateJobsSparkSqlJobLoggingConfig from a JSON
// response object.
func flattenWorkflowTemplateJobsSparkSqlJobLoggingConfigSlice(c *Client, i interface{}, res *WorkflowTemplate) []WorkflowTemplateJobsSparkSqlJobLoggingConfig {
	a, ok := i.([]interface{})
	if !ok {
		return []WorkflowTemplateJobsSparkSqlJobLoggingConfig{}
	}

	if len(a) == 0 {
		return []WorkflowTemplateJobsSparkSqlJobLoggingConfig{}
	}

	items := make([]WorkflowTemplateJobsSparkSqlJobLoggingConfig, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenWorkflowTemplateJobsSparkSqlJobLoggingConfig(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandWorkflowTemplateJobsSparkSqlJobLoggingConfig expands an instance of WorkflowTemplateJobsSparkSqlJobLoggingConfig into a JSON
// request object.
func expandWorkflowTemplateJobsSparkSqlJobLoggingConfig(c *Client, f *WorkflowTemplateJobsSparkSqlJobLoggingConfig, res *WorkflowTemplate) (map[string]interface{}, error) {
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
func flattenWorkflowTemplateJobsSparkSqlJobLoggingConfig(c *Client, i interface{}, res *WorkflowTemplate) *WorkflowTemplateJobsSparkSqlJobLoggingConfig {
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
func expandWorkflowTemplateJobsPrestoJobMap(c *Client, f map[string]WorkflowTemplateJobsPrestoJob, res *WorkflowTemplate) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandWorkflowTemplateJobsPrestoJob(c, &item, res)
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
func expandWorkflowTemplateJobsPrestoJobSlice(c *Client, f []WorkflowTemplateJobsPrestoJob, res *WorkflowTemplate) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandWorkflowTemplateJobsPrestoJob(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenWorkflowTemplateJobsPrestoJobMap flattens the contents of WorkflowTemplateJobsPrestoJob from a JSON
// response object.
func flattenWorkflowTemplateJobsPrestoJobMap(c *Client, i interface{}, res *WorkflowTemplate) map[string]WorkflowTemplateJobsPrestoJob {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]WorkflowTemplateJobsPrestoJob{}
	}

	if len(a) == 0 {
		return map[string]WorkflowTemplateJobsPrestoJob{}
	}

	items := make(map[string]WorkflowTemplateJobsPrestoJob)
	for k, item := range a {
		items[k] = *flattenWorkflowTemplateJobsPrestoJob(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenWorkflowTemplateJobsPrestoJobSlice flattens the contents of WorkflowTemplateJobsPrestoJob from a JSON
// response object.
func flattenWorkflowTemplateJobsPrestoJobSlice(c *Client, i interface{}, res *WorkflowTemplate) []WorkflowTemplateJobsPrestoJob {
	a, ok := i.([]interface{})
	if !ok {
		return []WorkflowTemplateJobsPrestoJob{}
	}

	if len(a) == 0 {
		return []WorkflowTemplateJobsPrestoJob{}
	}

	items := make([]WorkflowTemplateJobsPrestoJob, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenWorkflowTemplateJobsPrestoJob(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandWorkflowTemplateJobsPrestoJob expands an instance of WorkflowTemplateJobsPrestoJob into a JSON
// request object.
func expandWorkflowTemplateJobsPrestoJob(c *Client, f *WorkflowTemplateJobsPrestoJob, res *WorkflowTemplate) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.QueryFileUri; !dcl.IsEmptyValueIndirect(v) {
		m["queryFileUri"] = v
	}
	if v, err := expandWorkflowTemplateJobsPrestoJobQueryList(c, f.QueryList, res); err != nil {
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
	if v, err := expandWorkflowTemplateJobsPrestoJobLoggingConfig(c, f.LoggingConfig, res); err != nil {
		return nil, fmt.Errorf("error expanding LoggingConfig into loggingConfig: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["loggingConfig"] = v
	}

	return m, nil
}

// flattenWorkflowTemplateJobsPrestoJob flattens an instance of WorkflowTemplateJobsPrestoJob from a JSON
// response object.
func flattenWorkflowTemplateJobsPrestoJob(c *Client, i interface{}, res *WorkflowTemplate) *WorkflowTemplateJobsPrestoJob {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &WorkflowTemplateJobsPrestoJob{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyWorkflowTemplateJobsPrestoJob
	}
	r.QueryFileUri = dcl.FlattenString(m["queryFileUri"])
	r.QueryList = flattenWorkflowTemplateJobsPrestoJobQueryList(c, m["queryList"], res)
	r.ContinueOnFailure = dcl.FlattenBool(m["continueOnFailure"])
	r.OutputFormat = dcl.FlattenString(m["outputFormat"])
	r.ClientTags = dcl.FlattenStringSlice(m["clientTags"])
	r.Properties = dcl.FlattenKeyValuePairs(m["properties"])
	r.LoggingConfig = flattenWorkflowTemplateJobsPrestoJobLoggingConfig(c, m["loggingConfig"], res)

	return r
}

// expandWorkflowTemplateJobsPrestoJobQueryListMap expands the contents of WorkflowTemplateJobsPrestoJobQueryList into a JSON
// request object.
func expandWorkflowTemplateJobsPrestoJobQueryListMap(c *Client, f map[string]WorkflowTemplateJobsPrestoJobQueryList, res *WorkflowTemplate) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandWorkflowTemplateJobsPrestoJobQueryList(c, &item, res)
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
func expandWorkflowTemplateJobsPrestoJobQueryListSlice(c *Client, f []WorkflowTemplateJobsPrestoJobQueryList, res *WorkflowTemplate) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandWorkflowTemplateJobsPrestoJobQueryList(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenWorkflowTemplateJobsPrestoJobQueryListMap flattens the contents of WorkflowTemplateJobsPrestoJobQueryList from a JSON
// response object.
func flattenWorkflowTemplateJobsPrestoJobQueryListMap(c *Client, i interface{}, res *WorkflowTemplate) map[string]WorkflowTemplateJobsPrestoJobQueryList {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]WorkflowTemplateJobsPrestoJobQueryList{}
	}

	if len(a) == 0 {
		return map[string]WorkflowTemplateJobsPrestoJobQueryList{}
	}

	items := make(map[string]WorkflowTemplateJobsPrestoJobQueryList)
	for k, item := range a {
		items[k] = *flattenWorkflowTemplateJobsPrestoJobQueryList(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenWorkflowTemplateJobsPrestoJobQueryListSlice flattens the contents of WorkflowTemplateJobsPrestoJobQueryList from a JSON
// response object.
func flattenWorkflowTemplateJobsPrestoJobQueryListSlice(c *Client, i interface{}, res *WorkflowTemplate) []WorkflowTemplateJobsPrestoJobQueryList {
	a, ok := i.([]interface{})
	if !ok {
		return []WorkflowTemplateJobsPrestoJobQueryList{}
	}

	if len(a) == 0 {
		return []WorkflowTemplateJobsPrestoJobQueryList{}
	}

	items := make([]WorkflowTemplateJobsPrestoJobQueryList, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenWorkflowTemplateJobsPrestoJobQueryList(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandWorkflowTemplateJobsPrestoJobQueryList expands an instance of WorkflowTemplateJobsPrestoJobQueryList into a JSON
// request object.
func expandWorkflowTemplateJobsPrestoJobQueryList(c *Client, f *WorkflowTemplateJobsPrestoJobQueryList, res *WorkflowTemplate) (map[string]interface{}, error) {
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
func flattenWorkflowTemplateJobsPrestoJobQueryList(c *Client, i interface{}, res *WorkflowTemplate) *WorkflowTemplateJobsPrestoJobQueryList {
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
func expandWorkflowTemplateJobsPrestoJobLoggingConfigMap(c *Client, f map[string]WorkflowTemplateJobsPrestoJobLoggingConfig, res *WorkflowTemplate) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandWorkflowTemplateJobsPrestoJobLoggingConfig(c, &item, res)
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
func expandWorkflowTemplateJobsPrestoJobLoggingConfigSlice(c *Client, f []WorkflowTemplateJobsPrestoJobLoggingConfig, res *WorkflowTemplate) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandWorkflowTemplateJobsPrestoJobLoggingConfig(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenWorkflowTemplateJobsPrestoJobLoggingConfigMap flattens the contents of WorkflowTemplateJobsPrestoJobLoggingConfig from a JSON
// response object.
func flattenWorkflowTemplateJobsPrestoJobLoggingConfigMap(c *Client, i interface{}, res *WorkflowTemplate) map[string]WorkflowTemplateJobsPrestoJobLoggingConfig {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]WorkflowTemplateJobsPrestoJobLoggingConfig{}
	}

	if len(a) == 0 {
		return map[string]WorkflowTemplateJobsPrestoJobLoggingConfig{}
	}

	items := make(map[string]WorkflowTemplateJobsPrestoJobLoggingConfig)
	for k, item := range a {
		items[k] = *flattenWorkflowTemplateJobsPrestoJobLoggingConfig(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenWorkflowTemplateJobsPrestoJobLoggingConfigSlice flattens the contents of WorkflowTemplateJobsPrestoJobLoggingConfig from a JSON
// response object.
func flattenWorkflowTemplateJobsPrestoJobLoggingConfigSlice(c *Client, i interface{}, res *WorkflowTemplate) []WorkflowTemplateJobsPrestoJobLoggingConfig {
	a, ok := i.([]interface{})
	if !ok {
		return []WorkflowTemplateJobsPrestoJobLoggingConfig{}
	}

	if len(a) == 0 {
		return []WorkflowTemplateJobsPrestoJobLoggingConfig{}
	}

	items := make([]WorkflowTemplateJobsPrestoJobLoggingConfig, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenWorkflowTemplateJobsPrestoJobLoggingConfig(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandWorkflowTemplateJobsPrestoJobLoggingConfig expands an instance of WorkflowTemplateJobsPrestoJobLoggingConfig into a JSON
// request object.
func expandWorkflowTemplateJobsPrestoJobLoggingConfig(c *Client, f *WorkflowTemplateJobsPrestoJobLoggingConfig, res *WorkflowTemplate) (map[string]interface{}, error) {
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
func flattenWorkflowTemplateJobsPrestoJobLoggingConfig(c *Client, i interface{}, res *WorkflowTemplate) *WorkflowTemplateJobsPrestoJobLoggingConfig {
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
func expandWorkflowTemplateJobsSchedulingMap(c *Client, f map[string]WorkflowTemplateJobsScheduling, res *WorkflowTemplate) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandWorkflowTemplateJobsScheduling(c, &item, res)
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
func expandWorkflowTemplateJobsSchedulingSlice(c *Client, f []WorkflowTemplateJobsScheduling, res *WorkflowTemplate) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandWorkflowTemplateJobsScheduling(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenWorkflowTemplateJobsSchedulingMap flattens the contents of WorkflowTemplateJobsScheduling from a JSON
// response object.
func flattenWorkflowTemplateJobsSchedulingMap(c *Client, i interface{}, res *WorkflowTemplate) map[string]WorkflowTemplateJobsScheduling {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]WorkflowTemplateJobsScheduling{}
	}

	if len(a) == 0 {
		return map[string]WorkflowTemplateJobsScheduling{}
	}

	items := make(map[string]WorkflowTemplateJobsScheduling)
	for k, item := range a {
		items[k] = *flattenWorkflowTemplateJobsScheduling(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenWorkflowTemplateJobsSchedulingSlice flattens the contents of WorkflowTemplateJobsScheduling from a JSON
// response object.
func flattenWorkflowTemplateJobsSchedulingSlice(c *Client, i interface{}, res *WorkflowTemplate) []WorkflowTemplateJobsScheduling {
	a, ok := i.([]interface{})
	if !ok {
		return []WorkflowTemplateJobsScheduling{}
	}

	if len(a) == 0 {
		return []WorkflowTemplateJobsScheduling{}
	}

	items := make([]WorkflowTemplateJobsScheduling, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenWorkflowTemplateJobsScheduling(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandWorkflowTemplateJobsScheduling expands an instance of WorkflowTemplateJobsScheduling into a JSON
// request object.
func expandWorkflowTemplateJobsScheduling(c *Client, f *WorkflowTemplateJobsScheduling, res *WorkflowTemplate) (map[string]interface{}, error) {
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
func flattenWorkflowTemplateJobsScheduling(c *Client, i interface{}, res *WorkflowTemplate) *WorkflowTemplateJobsScheduling {
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
func expandWorkflowTemplateParametersMap(c *Client, f map[string]WorkflowTemplateParameters, res *WorkflowTemplate) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandWorkflowTemplateParameters(c, &item, res)
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
func expandWorkflowTemplateParametersSlice(c *Client, f []WorkflowTemplateParameters, res *WorkflowTemplate) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandWorkflowTemplateParameters(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenWorkflowTemplateParametersMap flattens the contents of WorkflowTemplateParameters from a JSON
// response object.
func flattenWorkflowTemplateParametersMap(c *Client, i interface{}, res *WorkflowTemplate) map[string]WorkflowTemplateParameters {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]WorkflowTemplateParameters{}
	}

	if len(a) == 0 {
		return map[string]WorkflowTemplateParameters{}
	}

	items := make(map[string]WorkflowTemplateParameters)
	for k, item := range a {
		items[k] = *flattenWorkflowTemplateParameters(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenWorkflowTemplateParametersSlice flattens the contents of WorkflowTemplateParameters from a JSON
// response object.
func flattenWorkflowTemplateParametersSlice(c *Client, i interface{}, res *WorkflowTemplate) []WorkflowTemplateParameters {
	a, ok := i.([]interface{})
	if !ok {
		return []WorkflowTemplateParameters{}
	}

	if len(a) == 0 {
		return []WorkflowTemplateParameters{}
	}

	items := make([]WorkflowTemplateParameters, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenWorkflowTemplateParameters(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandWorkflowTemplateParameters expands an instance of WorkflowTemplateParameters into a JSON
// request object.
func expandWorkflowTemplateParameters(c *Client, f *WorkflowTemplateParameters, res *WorkflowTemplate) (map[string]interface{}, error) {
	if f == nil {
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
	if v, err := expandWorkflowTemplateParametersValidation(c, f.Validation, res); err != nil {
		return nil, fmt.Errorf("error expanding Validation into validation: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["validation"] = v
	}

	return m, nil
}

// flattenWorkflowTemplateParameters flattens an instance of WorkflowTemplateParameters from a JSON
// response object.
func flattenWorkflowTemplateParameters(c *Client, i interface{}, res *WorkflowTemplate) *WorkflowTemplateParameters {
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
	r.Validation = flattenWorkflowTemplateParametersValidation(c, m["validation"], res)

	return r
}

// expandWorkflowTemplateParametersValidationMap expands the contents of WorkflowTemplateParametersValidation into a JSON
// request object.
func expandWorkflowTemplateParametersValidationMap(c *Client, f map[string]WorkflowTemplateParametersValidation, res *WorkflowTemplate) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandWorkflowTemplateParametersValidation(c, &item, res)
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
func expandWorkflowTemplateParametersValidationSlice(c *Client, f []WorkflowTemplateParametersValidation, res *WorkflowTemplate) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandWorkflowTemplateParametersValidation(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenWorkflowTemplateParametersValidationMap flattens the contents of WorkflowTemplateParametersValidation from a JSON
// response object.
func flattenWorkflowTemplateParametersValidationMap(c *Client, i interface{}, res *WorkflowTemplate) map[string]WorkflowTemplateParametersValidation {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]WorkflowTemplateParametersValidation{}
	}

	if len(a) == 0 {
		return map[string]WorkflowTemplateParametersValidation{}
	}

	items := make(map[string]WorkflowTemplateParametersValidation)
	for k, item := range a {
		items[k] = *flattenWorkflowTemplateParametersValidation(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenWorkflowTemplateParametersValidationSlice flattens the contents of WorkflowTemplateParametersValidation from a JSON
// response object.
func flattenWorkflowTemplateParametersValidationSlice(c *Client, i interface{}, res *WorkflowTemplate) []WorkflowTemplateParametersValidation {
	a, ok := i.([]interface{})
	if !ok {
		return []WorkflowTemplateParametersValidation{}
	}

	if len(a) == 0 {
		return []WorkflowTemplateParametersValidation{}
	}

	items := make([]WorkflowTemplateParametersValidation, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenWorkflowTemplateParametersValidation(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandWorkflowTemplateParametersValidation expands an instance of WorkflowTemplateParametersValidation into a JSON
// request object.
func expandWorkflowTemplateParametersValidation(c *Client, f *WorkflowTemplateParametersValidation, res *WorkflowTemplate) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v, err := expandWorkflowTemplateParametersValidationRegex(c, f.Regex, res); err != nil {
		return nil, fmt.Errorf("error expanding Regex into regex: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["regex"] = v
	}
	if v, err := expandWorkflowTemplateParametersValidationValues(c, f.Values, res); err != nil {
		return nil, fmt.Errorf("error expanding Values into values: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["values"] = v
	}

	return m, nil
}

// flattenWorkflowTemplateParametersValidation flattens an instance of WorkflowTemplateParametersValidation from a JSON
// response object.
func flattenWorkflowTemplateParametersValidation(c *Client, i interface{}, res *WorkflowTemplate) *WorkflowTemplateParametersValidation {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &WorkflowTemplateParametersValidation{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyWorkflowTemplateParametersValidation
	}
	r.Regex = flattenWorkflowTemplateParametersValidationRegex(c, m["regex"], res)
	r.Values = flattenWorkflowTemplateParametersValidationValues(c, m["values"], res)

	return r
}

// expandWorkflowTemplateParametersValidationRegexMap expands the contents of WorkflowTemplateParametersValidationRegex into a JSON
// request object.
func expandWorkflowTemplateParametersValidationRegexMap(c *Client, f map[string]WorkflowTemplateParametersValidationRegex, res *WorkflowTemplate) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandWorkflowTemplateParametersValidationRegex(c, &item, res)
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
func expandWorkflowTemplateParametersValidationRegexSlice(c *Client, f []WorkflowTemplateParametersValidationRegex, res *WorkflowTemplate) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandWorkflowTemplateParametersValidationRegex(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenWorkflowTemplateParametersValidationRegexMap flattens the contents of WorkflowTemplateParametersValidationRegex from a JSON
// response object.
func flattenWorkflowTemplateParametersValidationRegexMap(c *Client, i interface{}, res *WorkflowTemplate) map[string]WorkflowTemplateParametersValidationRegex {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]WorkflowTemplateParametersValidationRegex{}
	}

	if len(a) == 0 {
		return map[string]WorkflowTemplateParametersValidationRegex{}
	}

	items := make(map[string]WorkflowTemplateParametersValidationRegex)
	for k, item := range a {
		items[k] = *flattenWorkflowTemplateParametersValidationRegex(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenWorkflowTemplateParametersValidationRegexSlice flattens the contents of WorkflowTemplateParametersValidationRegex from a JSON
// response object.
func flattenWorkflowTemplateParametersValidationRegexSlice(c *Client, i interface{}, res *WorkflowTemplate) []WorkflowTemplateParametersValidationRegex {
	a, ok := i.([]interface{})
	if !ok {
		return []WorkflowTemplateParametersValidationRegex{}
	}

	if len(a) == 0 {
		return []WorkflowTemplateParametersValidationRegex{}
	}

	items := make([]WorkflowTemplateParametersValidationRegex, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenWorkflowTemplateParametersValidationRegex(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandWorkflowTemplateParametersValidationRegex expands an instance of WorkflowTemplateParametersValidationRegex into a JSON
// request object.
func expandWorkflowTemplateParametersValidationRegex(c *Client, f *WorkflowTemplateParametersValidationRegex, res *WorkflowTemplate) (map[string]interface{}, error) {
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
func flattenWorkflowTemplateParametersValidationRegex(c *Client, i interface{}, res *WorkflowTemplate) *WorkflowTemplateParametersValidationRegex {
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
func expandWorkflowTemplateParametersValidationValuesMap(c *Client, f map[string]WorkflowTemplateParametersValidationValues, res *WorkflowTemplate) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandWorkflowTemplateParametersValidationValues(c, &item, res)
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
func expandWorkflowTemplateParametersValidationValuesSlice(c *Client, f []WorkflowTemplateParametersValidationValues, res *WorkflowTemplate) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandWorkflowTemplateParametersValidationValues(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenWorkflowTemplateParametersValidationValuesMap flattens the contents of WorkflowTemplateParametersValidationValues from a JSON
// response object.
func flattenWorkflowTemplateParametersValidationValuesMap(c *Client, i interface{}, res *WorkflowTemplate) map[string]WorkflowTemplateParametersValidationValues {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]WorkflowTemplateParametersValidationValues{}
	}

	if len(a) == 0 {
		return map[string]WorkflowTemplateParametersValidationValues{}
	}

	items := make(map[string]WorkflowTemplateParametersValidationValues)
	for k, item := range a {
		items[k] = *flattenWorkflowTemplateParametersValidationValues(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenWorkflowTemplateParametersValidationValuesSlice flattens the contents of WorkflowTemplateParametersValidationValues from a JSON
// response object.
func flattenWorkflowTemplateParametersValidationValuesSlice(c *Client, i interface{}, res *WorkflowTemplate) []WorkflowTemplateParametersValidationValues {
	a, ok := i.([]interface{})
	if !ok {
		return []WorkflowTemplateParametersValidationValues{}
	}

	if len(a) == 0 {
		return []WorkflowTemplateParametersValidationValues{}
	}

	items := make([]WorkflowTemplateParametersValidationValues, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenWorkflowTemplateParametersValidationValues(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandWorkflowTemplateParametersValidationValues expands an instance of WorkflowTemplateParametersValidationValues into a JSON
// request object.
func expandWorkflowTemplateParametersValidationValues(c *Client, f *WorkflowTemplateParametersValidationValues, res *WorkflowTemplate) (map[string]interface{}, error) {
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
func flattenWorkflowTemplateParametersValidationValues(c *Client, i interface{}, res *WorkflowTemplate) *WorkflowTemplateParametersValidationValues {
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

// flattenWorkflowTemplatePlacementManagedClusterConfigGceClusterConfigPrivateIPv6GoogleAccessEnumMap flattens the contents of WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigPrivateIPv6GoogleAccessEnum from a JSON
// response object.
func flattenWorkflowTemplatePlacementManagedClusterConfigGceClusterConfigPrivateIPv6GoogleAccessEnumMap(c *Client, i interface{}, res *WorkflowTemplate) map[string]WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigPrivateIPv6GoogleAccessEnum {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigPrivateIPv6GoogleAccessEnum{}
	}

	if len(a) == 0 {
		return map[string]WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigPrivateIPv6GoogleAccessEnum{}
	}

	items := make(map[string]WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigPrivateIPv6GoogleAccessEnum)
	for k, item := range a {
		items[k] = *flattenWorkflowTemplatePlacementManagedClusterConfigGceClusterConfigPrivateIPv6GoogleAccessEnum(item.(interface{}))
	}

	return items
}

// flattenWorkflowTemplatePlacementManagedClusterConfigGceClusterConfigPrivateIPv6GoogleAccessEnumSlice flattens the contents of WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigPrivateIPv6GoogleAccessEnum from a JSON
// response object.
func flattenWorkflowTemplatePlacementManagedClusterConfigGceClusterConfigPrivateIPv6GoogleAccessEnumSlice(c *Client, i interface{}, res *WorkflowTemplate) []WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigPrivateIPv6GoogleAccessEnum {
	a, ok := i.([]interface{})
	if !ok {
		return []WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigPrivateIPv6GoogleAccessEnum{}
	}

	if len(a) == 0 {
		return []WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigPrivateIPv6GoogleAccessEnum{}
	}

	items := make([]WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigPrivateIPv6GoogleAccessEnum, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenWorkflowTemplatePlacementManagedClusterConfigGceClusterConfigPrivateIPv6GoogleAccessEnum(item.(interface{})))
	}

	return items
}

// flattenWorkflowTemplatePlacementManagedClusterConfigGceClusterConfigPrivateIPv6GoogleAccessEnum asserts that an interface is a string, and returns a
// pointer to a *WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigPrivateIPv6GoogleAccessEnum with the same value as that string.
func flattenWorkflowTemplatePlacementManagedClusterConfigGceClusterConfigPrivateIPv6GoogleAccessEnum(i interface{}) *WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigPrivateIPv6GoogleAccessEnum {
	s, ok := i.(string)
	if !ok {
		return nil
	}

	return WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigPrivateIPv6GoogleAccessEnumRef(s)
}

// flattenWorkflowTemplatePlacementManagedClusterConfigGceClusterConfigReservationAffinityConsumeReservationTypeEnumMap flattens the contents of WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigReservationAffinityConsumeReservationTypeEnum from a JSON
// response object.
func flattenWorkflowTemplatePlacementManagedClusterConfigGceClusterConfigReservationAffinityConsumeReservationTypeEnumMap(c *Client, i interface{}, res *WorkflowTemplate) map[string]WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigReservationAffinityConsumeReservationTypeEnum {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigReservationAffinityConsumeReservationTypeEnum{}
	}

	if len(a) == 0 {
		return map[string]WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigReservationAffinityConsumeReservationTypeEnum{}
	}

	items := make(map[string]WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigReservationAffinityConsumeReservationTypeEnum)
	for k, item := range a {
		items[k] = *flattenWorkflowTemplatePlacementManagedClusterConfigGceClusterConfigReservationAffinityConsumeReservationTypeEnum(item.(interface{}))
	}

	return items
}

// flattenWorkflowTemplatePlacementManagedClusterConfigGceClusterConfigReservationAffinityConsumeReservationTypeEnumSlice flattens the contents of WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigReservationAffinityConsumeReservationTypeEnum from a JSON
// response object.
func flattenWorkflowTemplatePlacementManagedClusterConfigGceClusterConfigReservationAffinityConsumeReservationTypeEnumSlice(c *Client, i interface{}, res *WorkflowTemplate) []WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigReservationAffinityConsumeReservationTypeEnum {
	a, ok := i.([]interface{})
	if !ok {
		return []WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigReservationAffinityConsumeReservationTypeEnum{}
	}

	if len(a) == 0 {
		return []WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigReservationAffinityConsumeReservationTypeEnum{}
	}

	items := make([]WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigReservationAffinityConsumeReservationTypeEnum, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenWorkflowTemplatePlacementManagedClusterConfigGceClusterConfigReservationAffinityConsumeReservationTypeEnum(item.(interface{})))
	}

	return items
}

// flattenWorkflowTemplatePlacementManagedClusterConfigGceClusterConfigReservationAffinityConsumeReservationTypeEnum asserts that an interface is a string, and returns a
// pointer to a *WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigReservationAffinityConsumeReservationTypeEnum with the same value as that string.
func flattenWorkflowTemplatePlacementManagedClusterConfigGceClusterConfigReservationAffinityConsumeReservationTypeEnum(i interface{}) *WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigReservationAffinityConsumeReservationTypeEnum {
	s, ok := i.(string)
	if !ok {
		return nil
	}

	return WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigReservationAffinityConsumeReservationTypeEnumRef(s)
}

// flattenWorkflowTemplatePlacementManagedClusterConfigMasterConfigPreemptibilityEnumMap flattens the contents of WorkflowTemplatePlacementManagedClusterConfigMasterConfigPreemptibilityEnum from a JSON
// response object.
func flattenWorkflowTemplatePlacementManagedClusterConfigMasterConfigPreemptibilityEnumMap(c *Client, i interface{}, res *WorkflowTemplate) map[string]WorkflowTemplatePlacementManagedClusterConfigMasterConfigPreemptibilityEnum {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]WorkflowTemplatePlacementManagedClusterConfigMasterConfigPreemptibilityEnum{}
	}

	if len(a) == 0 {
		return map[string]WorkflowTemplatePlacementManagedClusterConfigMasterConfigPreemptibilityEnum{}
	}

	items := make(map[string]WorkflowTemplatePlacementManagedClusterConfigMasterConfigPreemptibilityEnum)
	for k, item := range a {
		items[k] = *flattenWorkflowTemplatePlacementManagedClusterConfigMasterConfigPreemptibilityEnum(item.(interface{}))
	}

	return items
}

// flattenWorkflowTemplatePlacementManagedClusterConfigMasterConfigPreemptibilityEnumSlice flattens the contents of WorkflowTemplatePlacementManagedClusterConfigMasterConfigPreemptibilityEnum from a JSON
// response object.
func flattenWorkflowTemplatePlacementManagedClusterConfigMasterConfigPreemptibilityEnumSlice(c *Client, i interface{}, res *WorkflowTemplate) []WorkflowTemplatePlacementManagedClusterConfigMasterConfigPreemptibilityEnum {
	a, ok := i.([]interface{})
	if !ok {
		return []WorkflowTemplatePlacementManagedClusterConfigMasterConfigPreemptibilityEnum{}
	}

	if len(a) == 0 {
		return []WorkflowTemplatePlacementManagedClusterConfigMasterConfigPreemptibilityEnum{}
	}

	items := make([]WorkflowTemplatePlacementManagedClusterConfigMasterConfigPreemptibilityEnum, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenWorkflowTemplatePlacementManagedClusterConfigMasterConfigPreemptibilityEnum(item.(interface{})))
	}

	return items
}

// flattenWorkflowTemplatePlacementManagedClusterConfigMasterConfigPreemptibilityEnum asserts that an interface is a string, and returns a
// pointer to a *WorkflowTemplatePlacementManagedClusterConfigMasterConfigPreemptibilityEnum with the same value as that string.
func flattenWorkflowTemplatePlacementManagedClusterConfigMasterConfigPreemptibilityEnum(i interface{}) *WorkflowTemplatePlacementManagedClusterConfigMasterConfigPreemptibilityEnum {
	s, ok := i.(string)
	if !ok {
		return nil
	}

	return WorkflowTemplatePlacementManagedClusterConfigMasterConfigPreemptibilityEnumRef(s)
}

// flattenWorkflowTemplatePlacementManagedClusterConfigWorkerConfigPreemptibilityEnumMap flattens the contents of WorkflowTemplatePlacementManagedClusterConfigWorkerConfigPreemptibilityEnum from a JSON
// response object.
func flattenWorkflowTemplatePlacementManagedClusterConfigWorkerConfigPreemptibilityEnumMap(c *Client, i interface{}, res *WorkflowTemplate) map[string]WorkflowTemplatePlacementManagedClusterConfigWorkerConfigPreemptibilityEnum {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]WorkflowTemplatePlacementManagedClusterConfigWorkerConfigPreemptibilityEnum{}
	}

	if len(a) == 0 {
		return map[string]WorkflowTemplatePlacementManagedClusterConfigWorkerConfigPreemptibilityEnum{}
	}

	items := make(map[string]WorkflowTemplatePlacementManagedClusterConfigWorkerConfigPreemptibilityEnum)
	for k, item := range a {
		items[k] = *flattenWorkflowTemplatePlacementManagedClusterConfigWorkerConfigPreemptibilityEnum(item.(interface{}))
	}

	return items
}

// flattenWorkflowTemplatePlacementManagedClusterConfigWorkerConfigPreemptibilityEnumSlice flattens the contents of WorkflowTemplatePlacementManagedClusterConfigWorkerConfigPreemptibilityEnum from a JSON
// response object.
func flattenWorkflowTemplatePlacementManagedClusterConfigWorkerConfigPreemptibilityEnumSlice(c *Client, i interface{}, res *WorkflowTemplate) []WorkflowTemplatePlacementManagedClusterConfigWorkerConfigPreemptibilityEnum {
	a, ok := i.([]interface{})
	if !ok {
		return []WorkflowTemplatePlacementManagedClusterConfigWorkerConfigPreemptibilityEnum{}
	}

	if len(a) == 0 {
		return []WorkflowTemplatePlacementManagedClusterConfigWorkerConfigPreemptibilityEnum{}
	}

	items := make([]WorkflowTemplatePlacementManagedClusterConfigWorkerConfigPreemptibilityEnum, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenWorkflowTemplatePlacementManagedClusterConfigWorkerConfigPreemptibilityEnum(item.(interface{})))
	}

	return items
}

// flattenWorkflowTemplatePlacementManagedClusterConfigWorkerConfigPreemptibilityEnum asserts that an interface is a string, and returns a
// pointer to a *WorkflowTemplatePlacementManagedClusterConfigWorkerConfigPreemptibilityEnum with the same value as that string.
func flattenWorkflowTemplatePlacementManagedClusterConfigWorkerConfigPreemptibilityEnum(i interface{}) *WorkflowTemplatePlacementManagedClusterConfigWorkerConfigPreemptibilityEnum {
	s, ok := i.(string)
	if !ok {
		return nil
	}

	return WorkflowTemplatePlacementManagedClusterConfigWorkerConfigPreemptibilityEnumRef(s)
}

// flattenWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigPreemptibilityEnumMap flattens the contents of WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigPreemptibilityEnum from a JSON
// response object.
func flattenWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigPreemptibilityEnumMap(c *Client, i interface{}, res *WorkflowTemplate) map[string]WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigPreemptibilityEnum {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigPreemptibilityEnum{}
	}

	if len(a) == 0 {
		return map[string]WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigPreemptibilityEnum{}
	}

	items := make(map[string]WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigPreemptibilityEnum)
	for k, item := range a {
		items[k] = *flattenWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigPreemptibilityEnum(item.(interface{}))
	}

	return items
}

// flattenWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigPreemptibilityEnumSlice flattens the contents of WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigPreemptibilityEnum from a JSON
// response object.
func flattenWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigPreemptibilityEnumSlice(c *Client, i interface{}, res *WorkflowTemplate) []WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigPreemptibilityEnum {
	a, ok := i.([]interface{})
	if !ok {
		return []WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigPreemptibilityEnum{}
	}

	if len(a) == 0 {
		return []WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigPreemptibilityEnum{}
	}

	items := make([]WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigPreemptibilityEnum, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigPreemptibilityEnum(item.(interface{})))
	}

	return items
}

// flattenWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigPreemptibilityEnum asserts that an interface is a string, and returns a
// pointer to a *WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigPreemptibilityEnum with the same value as that string.
func flattenWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigPreemptibilityEnum(i interface{}) *WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigPreemptibilityEnum {
	s, ok := i.(string)
	if !ok {
		return nil
	}

	return WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigPreemptibilityEnumRef(s)
}

// flattenWorkflowTemplatePlacementManagedClusterConfigSoftwareConfigOptionalComponentsEnumMap flattens the contents of WorkflowTemplatePlacementManagedClusterConfigSoftwareConfigOptionalComponentsEnum from a JSON
// response object.
func flattenWorkflowTemplatePlacementManagedClusterConfigSoftwareConfigOptionalComponentsEnumMap(c *Client, i interface{}, res *WorkflowTemplate) map[string]WorkflowTemplatePlacementManagedClusterConfigSoftwareConfigOptionalComponentsEnum {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]WorkflowTemplatePlacementManagedClusterConfigSoftwareConfigOptionalComponentsEnum{}
	}

	if len(a) == 0 {
		return map[string]WorkflowTemplatePlacementManagedClusterConfigSoftwareConfigOptionalComponentsEnum{}
	}

	items := make(map[string]WorkflowTemplatePlacementManagedClusterConfigSoftwareConfigOptionalComponentsEnum)
	for k, item := range a {
		items[k] = *flattenWorkflowTemplatePlacementManagedClusterConfigSoftwareConfigOptionalComponentsEnum(item.(interface{}))
	}

	return items
}

// flattenWorkflowTemplatePlacementManagedClusterConfigSoftwareConfigOptionalComponentsEnumSlice flattens the contents of WorkflowTemplatePlacementManagedClusterConfigSoftwareConfigOptionalComponentsEnum from a JSON
// response object.
func flattenWorkflowTemplatePlacementManagedClusterConfigSoftwareConfigOptionalComponentsEnumSlice(c *Client, i interface{}, res *WorkflowTemplate) []WorkflowTemplatePlacementManagedClusterConfigSoftwareConfigOptionalComponentsEnum {
	a, ok := i.([]interface{})
	if !ok {
		return []WorkflowTemplatePlacementManagedClusterConfigSoftwareConfigOptionalComponentsEnum{}
	}

	if len(a) == 0 {
		return []WorkflowTemplatePlacementManagedClusterConfigSoftwareConfigOptionalComponentsEnum{}
	}

	items := make([]WorkflowTemplatePlacementManagedClusterConfigSoftwareConfigOptionalComponentsEnum, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenWorkflowTemplatePlacementManagedClusterConfigSoftwareConfigOptionalComponentsEnum(item.(interface{})))
	}

	return items
}

// flattenWorkflowTemplatePlacementManagedClusterConfigSoftwareConfigOptionalComponentsEnum asserts that an interface is a string, and returns a
// pointer to a *WorkflowTemplatePlacementManagedClusterConfigSoftwareConfigOptionalComponentsEnum with the same value as that string.
func flattenWorkflowTemplatePlacementManagedClusterConfigSoftwareConfigOptionalComponentsEnum(i interface{}) *WorkflowTemplatePlacementManagedClusterConfigSoftwareConfigOptionalComponentsEnum {
	s, ok := i.(string)
	if !ok {
		return nil
	}

	return WorkflowTemplatePlacementManagedClusterConfigSoftwareConfigOptionalComponentsEnumRef(s)
}

// This function returns a matcher that checks whether a serialized resource matches this resource
// in its parameters (as defined by the fields in a Get, which definitionally define resource
// identity).  This is useful in extracting the element from a List call.
func (r *WorkflowTemplate) matcher(c *Client) func([]byte) bool {
	return func(b []byte) bool {
		cr, err := unmarshalWorkflowTemplate(b, c, r)
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
	FieldName        string // used for error logging
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
		// Use the first field diff's field name for logging required recreate error.
		diff := workflowTemplateDiff{FieldName: fieldDiffs[0].FieldName}
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
	if !dcl.IsEmptyValueIndirect(vPlacement) {
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
	if !dcl.IsEmptyValueIndirect(vManagedCluster) {
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
	if !dcl.IsEmptyValueIndirect(vClusterSelector) {
		o.ClusterSelector = vClusterSelector
	}
	return nil
}
func extractWorkflowTemplatePlacementManagedClusterFields(r *WorkflowTemplate, o *WorkflowTemplatePlacementManagedCluster) error {
	vConfig := o.Config
	if vConfig == nil {
		// note: explicitly not the empty object.
		vConfig = &WorkflowTemplatePlacementManagedClusterConfig{}
	}
	if err := extractWorkflowTemplatePlacementManagedClusterConfigFields(r, vConfig); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vConfig) {
		o.Config = vConfig
	}
	return nil
}
func extractWorkflowTemplatePlacementManagedClusterConfigFields(r *WorkflowTemplate, o *WorkflowTemplatePlacementManagedClusterConfig) error {
	vGceClusterConfig := o.GceClusterConfig
	if vGceClusterConfig == nil {
		// note: explicitly not the empty object.
		vGceClusterConfig = &WorkflowTemplatePlacementManagedClusterConfigGceClusterConfig{}
	}
	if err := extractWorkflowTemplatePlacementManagedClusterConfigGceClusterConfigFields(r, vGceClusterConfig); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vGceClusterConfig) {
		o.GceClusterConfig = vGceClusterConfig
	}
	vMasterConfig := o.MasterConfig
	if vMasterConfig == nil {
		// note: explicitly not the empty object.
		vMasterConfig = &WorkflowTemplatePlacementManagedClusterConfigMasterConfig{}
	}
	if err := extractWorkflowTemplatePlacementManagedClusterConfigMasterConfigFields(r, vMasterConfig); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vMasterConfig) {
		o.MasterConfig = vMasterConfig
	}
	vWorkerConfig := o.WorkerConfig
	if vWorkerConfig == nil {
		// note: explicitly not the empty object.
		vWorkerConfig = &WorkflowTemplatePlacementManagedClusterConfigWorkerConfig{}
	}
	if err := extractWorkflowTemplatePlacementManagedClusterConfigWorkerConfigFields(r, vWorkerConfig); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vWorkerConfig) {
		o.WorkerConfig = vWorkerConfig
	}
	vSecondaryWorkerConfig := o.SecondaryWorkerConfig
	if vSecondaryWorkerConfig == nil {
		// note: explicitly not the empty object.
		vSecondaryWorkerConfig = &WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfig{}
	}
	if err := extractWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigFields(r, vSecondaryWorkerConfig); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vSecondaryWorkerConfig) {
		o.SecondaryWorkerConfig = vSecondaryWorkerConfig
	}
	vSoftwareConfig := o.SoftwareConfig
	if vSoftwareConfig == nil {
		// note: explicitly not the empty object.
		vSoftwareConfig = &WorkflowTemplatePlacementManagedClusterConfigSoftwareConfig{}
	}
	if err := extractWorkflowTemplatePlacementManagedClusterConfigSoftwareConfigFields(r, vSoftwareConfig); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vSoftwareConfig) {
		o.SoftwareConfig = vSoftwareConfig
	}
	vEncryptionConfig := o.EncryptionConfig
	if vEncryptionConfig == nil {
		// note: explicitly not the empty object.
		vEncryptionConfig = &WorkflowTemplatePlacementManagedClusterConfigEncryptionConfig{}
	}
	if err := extractWorkflowTemplatePlacementManagedClusterConfigEncryptionConfigFields(r, vEncryptionConfig); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vEncryptionConfig) {
		o.EncryptionConfig = vEncryptionConfig
	}
	vAutoscalingConfig := o.AutoscalingConfig
	if vAutoscalingConfig == nil {
		// note: explicitly not the empty object.
		vAutoscalingConfig = &WorkflowTemplatePlacementManagedClusterConfigAutoscalingConfig{}
	}
	if err := extractWorkflowTemplatePlacementManagedClusterConfigAutoscalingConfigFields(r, vAutoscalingConfig); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vAutoscalingConfig) {
		o.AutoscalingConfig = vAutoscalingConfig
	}
	vSecurityConfig := o.SecurityConfig
	if vSecurityConfig == nil {
		// note: explicitly not the empty object.
		vSecurityConfig = &WorkflowTemplatePlacementManagedClusterConfigSecurityConfig{}
	}
	if err := extractWorkflowTemplatePlacementManagedClusterConfigSecurityConfigFields(r, vSecurityConfig); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vSecurityConfig) {
		o.SecurityConfig = vSecurityConfig
	}
	vLifecycleConfig := o.LifecycleConfig
	if vLifecycleConfig == nil {
		// note: explicitly not the empty object.
		vLifecycleConfig = &WorkflowTemplatePlacementManagedClusterConfigLifecycleConfig{}
	}
	if err := extractWorkflowTemplatePlacementManagedClusterConfigLifecycleConfigFields(r, vLifecycleConfig); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vLifecycleConfig) {
		o.LifecycleConfig = vLifecycleConfig
	}
	vEndpointConfig := o.EndpointConfig
	if vEndpointConfig == nil {
		// note: explicitly not the empty object.
		vEndpointConfig = &WorkflowTemplatePlacementManagedClusterConfigEndpointConfig{}
	}
	if err := extractWorkflowTemplatePlacementManagedClusterConfigEndpointConfigFields(r, vEndpointConfig); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vEndpointConfig) {
		o.EndpointConfig = vEndpointConfig
	}
	return nil
}
func extractWorkflowTemplatePlacementManagedClusterConfigGceClusterConfigFields(r *WorkflowTemplate, o *WorkflowTemplatePlacementManagedClusterConfigGceClusterConfig) error {
	vReservationAffinity := o.ReservationAffinity
	if vReservationAffinity == nil {
		// note: explicitly not the empty object.
		vReservationAffinity = &WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigReservationAffinity{}
	}
	if err := extractWorkflowTemplatePlacementManagedClusterConfigGceClusterConfigReservationAffinityFields(r, vReservationAffinity); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vReservationAffinity) {
		o.ReservationAffinity = vReservationAffinity
	}
	vNodeGroupAffinity := o.NodeGroupAffinity
	if vNodeGroupAffinity == nil {
		// note: explicitly not the empty object.
		vNodeGroupAffinity = &WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigNodeGroupAffinity{}
	}
	if err := extractWorkflowTemplatePlacementManagedClusterConfigGceClusterConfigNodeGroupAffinityFields(r, vNodeGroupAffinity); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vNodeGroupAffinity) {
		o.NodeGroupAffinity = vNodeGroupAffinity
	}
	vShieldedInstanceConfig := o.ShieldedInstanceConfig
	if vShieldedInstanceConfig == nil {
		// note: explicitly not the empty object.
		vShieldedInstanceConfig = &WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigShieldedInstanceConfig{}
	}
	if err := extractWorkflowTemplatePlacementManagedClusterConfigGceClusterConfigShieldedInstanceConfigFields(r, vShieldedInstanceConfig); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vShieldedInstanceConfig) {
		o.ShieldedInstanceConfig = vShieldedInstanceConfig
	}
	return nil
}
func extractWorkflowTemplatePlacementManagedClusterConfigGceClusterConfigReservationAffinityFields(r *WorkflowTemplate, o *WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigReservationAffinity) error {
	return nil
}
func extractWorkflowTemplatePlacementManagedClusterConfigGceClusterConfigNodeGroupAffinityFields(r *WorkflowTemplate, o *WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigNodeGroupAffinity) error {
	return nil
}
func extractWorkflowTemplatePlacementManagedClusterConfigGceClusterConfigShieldedInstanceConfigFields(r *WorkflowTemplate, o *WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigShieldedInstanceConfig) error {
	return nil
}
func extractWorkflowTemplatePlacementManagedClusterConfigMasterConfigFields(r *WorkflowTemplate, o *WorkflowTemplatePlacementManagedClusterConfigMasterConfig) error {
	vDiskConfig := o.DiskConfig
	if vDiskConfig == nil {
		// note: explicitly not the empty object.
		vDiskConfig = &WorkflowTemplatePlacementManagedClusterConfigMasterConfigDiskConfig{}
	}
	if err := extractWorkflowTemplatePlacementManagedClusterConfigMasterConfigDiskConfigFields(r, vDiskConfig); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vDiskConfig) {
		o.DiskConfig = vDiskConfig
	}
	vManagedGroupConfig := o.ManagedGroupConfig
	if vManagedGroupConfig == nil {
		// note: explicitly not the empty object.
		vManagedGroupConfig = &WorkflowTemplatePlacementManagedClusterConfigMasterConfigManagedGroupConfig{}
	}
	if err := extractWorkflowTemplatePlacementManagedClusterConfigMasterConfigManagedGroupConfigFields(r, vManagedGroupConfig); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vManagedGroupConfig) {
		o.ManagedGroupConfig = vManagedGroupConfig
	}
	return nil
}
func extractWorkflowTemplatePlacementManagedClusterConfigMasterConfigDiskConfigFields(r *WorkflowTemplate, o *WorkflowTemplatePlacementManagedClusterConfigMasterConfigDiskConfig) error {
	return nil
}
func extractWorkflowTemplatePlacementManagedClusterConfigMasterConfigManagedGroupConfigFields(r *WorkflowTemplate, o *WorkflowTemplatePlacementManagedClusterConfigMasterConfigManagedGroupConfig) error {
	return nil
}
func extractWorkflowTemplatePlacementManagedClusterConfigMasterConfigAcceleratorsFields(r *WorkflowTemplate, o *WorkflowTemplatePlacementManagedClusterConfigMasterConfigAccelerators) error {
	return nil
}
func extractWorkflowTemplatePlacementManagedClusterConfigWorkerConfigFields(r *WorkflowTemplate, o *WorkflowTemplatePlacementManagedClusterConfigWorkerConfig) error {
	vDiskConfig := o.DiskConfig
	if vDiskConfig == nil {
		// note: explicitly not the empty object.
		vDiskConfig = &WorkflowTemplatePlacementManagedClusterConfigWorkerConfigDiskConfig{}
	}
	if err := extractWorkflowTemplatePlacementManagedClusterConfigWorkerConfigDiskConfigFields(r, vDiskConfig); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vDiskConfig) {
		o.DiskConfig = vDiskConfig
	}
	vManagedGroupConfig := o.ManagedGroupConfig
	if vManagedGroupConfig == nil {
		// note: explicitly not the empty object.
		vManagedGroupConfig = &WorkflowTemplatePlacementManagedClusterConfigWorkerConfigManagedGroupConfig{}
	}
	if err := extractWorkflowTemplatePlacementManagedClusterConfigWorkerConfigManagedGroupConfigFields(r, vManagedGroupConfig); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vManagedGroupConfig) {
		o.ManagedGroupConfig = vManagedGroupConfig
	}
	return nil
}
func extractWorkflowTemplatePlacementManagedClusterConfigWorkerConfigDiskConfigFields(r *WorkflowTemplate, o *WorkflowTemplatePlacementManagedClusterConfigWorkerConfigDiskConfig) error {
	return nil
}
func extractWorkflowTemplatePlacementManagedClusterConfigWorkerConfigManagedGroupConfigFields(r *WorkflowTemplate, o *WorkflowTemplatePlacementManagedClusterConfigWorkerConfigManagedGroupConfig) error {
	return nil
}
func extractWorkflowTemplatePlacementManagedClusterConfigWorkerConfigAcceleratorsFields(r *WorkflowTemplate, o *WorkflowTemplatePlacementManagedClusterConfigWorkerConfigAccelerators) error {
	return nil
}
func extractWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigFields(r *WorkflowTemplate, o *WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfig) error {
	vDiskConfig := o.DiskConfig
	if vDiskConfig == nil {
		// note: explicitly not the empty object.
		vDiskConfig = &WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigDiskConfig{}
	}
	if err := extractWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigDiskConfigFields(r, vDiskConfig); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vDiskConfig) {
		o.DiskConfig = vDiskConfig
	}
	vManagedGroupConfig := o.ManagedGroupConfig
	if vManagedGroupConfig == nil {
		// note: explicitly not the empty object.
		vManagedGroupConfig = &WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigManagedGroupConfig{}
	}
	if err := extractWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigManagedGroupConfigFields(r, vManagedGroupConfig); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vManagedGroupConfig) {
		o.ManagedGroupConfig = vManagedGroupConfig
	}
	return nil
}
func extractWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigDiskConfigFields(r *WorkflowTemplate, o *WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigDiskConfig) error {
	return nil
}
func extractWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigManagedGroupConfigFields(r *WorkflowTemplate, o *WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigManagedGroupConfig) error {
	return nil
}
func extractWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigAcceleratorsFields(r *WorkflowTemplate, o *WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigAccelerators) error {
	return nil
}
func extractWorkflowTemplatePlacementManagedClusterConfigSoftwareConfigFields(r *WorkflowTemplate, o *WorkflowTemplatePlacementManagedClusterConfigSoftwareConfig) error {
	return nil
}
func extractWorkflowTemplatePlacementManagedClusterConfigInitializationActionsFields(r *WorkflowTemplate, o *WorkflowTemplatePlacementManagedClusterConfigInitializationActions) error {
	return nil
}
func extractWorkflowTemplatePlacementManagedClusterConfigEncryptionConfigFields(r *WorkflowTemplate, o *WorkflowTemplatePlacementManagedClusterConfigEncryptionConfig) error {
	return nil
}
func extractWorkflowTemplatePlacementManagedClusterConfigAutoscalingConfigFields(r *WorkflowTemplate, o *WorkflowTemplatePlacementManagedClusterConfigAutoscalingConfig) error {
	return nil
}
func extractWorkflowTemplatePlacementManagedClusterConfigSecurityConfigFields(r *WorkflowTemplate, o *WorkflowTemplatePlacementManagedClusterConfigSecurityConfig) error {
	vKerberosConfig := o.KerberosConfig
	if vKerberosConfig == nil {
		// note: explicitly not the empty object.
		vKerberosConfig = &WorkflowTemplatePlacementManagedClusterConfigSecurityConfigKerberosConfig{}
	}
	if err := extractWorkflowTemplatePlacementManagedClusterConfigSecurityConfigKerberosConfigFields(r, vKerberosConfig); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vKerberosConfig) {
		o.KerberosConfig = vKerberosConfig
	}
	return nil
}
func extractWorkflowTemplatePlacementManagedClusterConfigSecurityConfigKerberosConfigFields(r *WorkflowTemplate, o *WorkflowTemplatePlacementManagedClusterConfigSecurityConfigKerberosConfig) error {
	return nil
}
func extractWorkflowTemplatePlacementManagedClusterConfigLifecycleConfigFields(r *WorkflowTemplate, o *WorkflowTemplatePlacementManagedClusterConfigLifecycleConfig) error {
	return nil
}
func extractWorkflowTemplatePlacementManagedClusterConfigEndpointConfigFields(r *WorkflowTemplate, o *WorkflowTemplatePlacementManagedClusterConfigEndpointConfig) error {
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
	if !dcl.IsEmptyValueIndirect(vHadoopJob) {
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
	if !dcl.IsEmptyValueIndirect(vSparkJob) {
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
	if !dcl.IsEmptyValueIndirect(vPysparkJob) {
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
	if !dcl.IsEmptyValueIndirect(vHiveJob) {
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
	if !dcl.IsEmptyValueIndirect(vPigJob) {
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
	if !dcl.IsEmptyValueIndirect(vSparkRJob) {
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
	if !dcl.IsEmptyValueIndirect(vSparkSqlJob) {
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
	if !dcl.IsEmptyValueIndirect(vPrestoJob) {
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
	if !dcl.IsEmptyValueIndirect(vScheduling) {
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
	if !dcl.IsEmptyValueIndirect(vLoggingConfig) {
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
	if !dcl.IsEmptyValueIndirect(vLoggingConfig) {
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
	if !dcl.IsEmptyValueIndirect(vLoggingConfig) {
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
	if !dcl.IsEmptyValueIndirect(vQueryList) {
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
	if !dcl.IsEmptyValueIndirect(vQueryList) {
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
	if !dcl.IsEmptyValueIndirect(vLoggingConfig) {
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
	if !dcl.IsEmptyValueIndirect(vLoggingConfig) {
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
	if !dcl.IsEmptyValueIndirect(vQueryList) {
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
	if !dcl.IsEmptyValueIndirect(vLoggingConfig) {
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
	if !dcl.IsEmptyValueIndirect(vQueryList) {
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
	if !dcl.IsEmptyValueIndirect(vLoggingConfig) {
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
	if !dcl.IsEmptyValueIndirect(vValidation) {
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
	if !dcl.IsEmptyValueIndirect(vRegex) {
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
	if !dcl.IsEmptyValueIndirect(vValues) {
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
	if !dcl.IsEmptyValueIndirect(vPlacement) {
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
	if !dcl.IsEmptyValueIndirect(vManagedCluster) {
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
	if !dcl.IsEmptyValueIndirect(vClusterSelector) {
		o.ClusterSelector = vClusterSelector
	}
	return nil
}
func postReadExtractWorkflowTemplatePlacementManagedClusterFields(r *WorkflowTemplate, o *WorkflowTemplatePlacementManagedCluster) error {
	vConfig := o.Config
	if vConfig == nil {
		// note: explicitly not the empty object.
		vConfig = &WorkflowTemplatePlacementManagedClusterConfig{}
	}
	if err := extractWorkflowTemplatePlacementManagedClusterConfigFields(r, vConfig); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vConfig) {
		o.Config = vConfig
	}
	return nil
}
func postReadExtractWorkflowTemplatePlacementManagedClusterConfigFields(r *WorkflowTemplate, o *WorkflowTemplatePlacementManagedClusterConfig) error {
	vGceClusterConfig := o.GceClusterConfig
	if vGceClusterConfig == nil {
		// note: explicitly not the empty object.
		vGceClusterConfig = &WorkflowTemplatePlacementManagedClusterConfigGceClusterConfig{}
	}
	if err := extractWorkflowTemplatePlacementManagedClusterConfigGceClusterConfigFields(r, vGceClusterConfig); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vGceClusterConfig) {
		o.GceClusterConfig = vGceClusterConfig
	}
	vMasterConfig := o.MasterConfig
	if vMasterConfig == nil {
		// note: explicitly not the empty object.
		vMasterConfig = &WorkflowTemplatePlacementManagedClusterConfigMasterConfig{}
	}
	if err := extractWorkflowTemplatePlacementManagedClusterConfigMasterConfigFields(r, vMasterConfig); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vMasterConfig) {
		o.MasterConfig = vMasterConfig
	}
	vWorkerConfig := o.WorkerConfig
	if vWorkerConfig == nil {
		// note: explicitly not the empty object.
		vWorkerConfig = &WorkflowTemplatePlacementManagedClusterConfigWorkerConfig{}
	}
	if err := extractWorkflowTemplatePlacementManagedClusterConfigWorkerConfigFields(r, vWorkerConfig); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vWorkerConfig) {
		o.WorkerConfig = vWorkerConfig
	}
	vSecondaryWorkerConfig := o.SecondaryWorkerConfig
	if vSecondaryWorkerConfig == nil {
		// note: explicitly not the empty object.
		vSecondaryWorkerConfig = &WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfig{}
	}
	if err := extractWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigFields(r, vSecondaryWorkerConfig); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vSecondaryWorkerConfig) {
		o.SecondaryWorkerConfig = vSecondaryWorkerConfig
	}
	vSoftwareConfig := o.SoftwareConfig
	if vSoftwareConfig == nil {
		// note: explicitly not the empty object.
		vSoftwareConfig = &WorkflowTemplatePlacementManagedClusterConfigSoftwareConfig{}
	}
	if err := extractWorkflowTemplatePlacementManagedClusterConfigSoftwareConfigFields(r, vSoftwareConfig); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vSoftwareConfig) {
		o.SoftwareConfig = vSoftwareConfig
	}
	vEncryptionConfig := o.EncryptionConfig
	if vEncryptionConfig == nil {
		// note: explicitly not the empty object.
		vEncryptionConfig = &WorkflowTemplatePlacementManagedClusterConfigEncryptionConfig{}
	}
	if err := extractWorkflowTemplatePlacementManagedClusterConfigEncryptionConfigFields(r, vEncryptionConfig); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vEncryptionConfig) {
		o.EncryptionConfig = vEncryptionConfig
	}
	vAutoscalingConfig := o.AutoscalingConfig
	if vAutoscalingConfig == nil {
		// note: explicitly not the empty object.
		vAutoscalingConfig = &WorkflowTemplatePlacementManagedClusterConfigAutoscalingConfig{}
	}
	if err := extractWorkflowTemplatePlacementManagedClusterConfigAutoscalingConfigFields(r, vAutoscalingConfig); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vAutoscalingConfig) {
		o.AutoscalingConfig = vAutoscalingConfig
	}
	vSecurityConfig := o.SecurityConfig
	if vSecurityConfig == nil {
		// note: explicitly not the empty object.
		vSecurityConfig = &WorkflowTemplatePlacementManagedClusterConfigSecurityConfig{}
	}
	if err := extractWorkflowTemplatePlacementManagedClusterConfigSecurityConfigFields(r, vSecurityConfig); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vSecurityConfig) {
		o.SecurityConfig = vSecurityConfig
	}
	vLifecycleConfig := o.LifecycleConfig
	if vLifecycleConfig == nil {
		// note: explicitly not the empty object.
		vLifecycleConfig = &WorkflowTemplatePlacementManagedClusterConfigLifecycleConfig{}
	}
	if err := extractWorkflowTemplatePlacementManagedClusterConfigLifecycleConfigFields(r, vLifecycleConfig); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vLifecycleConfig) {
		o.LifecycleConfig = vLifecycleConfig
	}
	vEndpointConfig := o.EndpointConfig
	if vEndpointConfig == nil {
		// note: explicitly not the empty object.
		vEndpointConfig = &WorkflowTemplatePlacementManagedClusterConfigEndpointConfig{}
	}
	if err := extractWorkflowTemplatePlacementManagedClusterConfigEndpointConfigFields(r, vEndpointConfig); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vEndpointConfig) {
		o.EndpointConfig = vEndpointConfig
	}
	return nil
}
func postReadExtractWorkflowTemplatePlacementManagedClusterConfigGceClusterConfigFields(r *WorkflowTemplate, o *WorkflowTemplatePlacementManagedClusterConfigGceClusterConfig) error {
	vReservationAffinity := o.ReservationAffinity
	if vReservationAffinity == nil {
		// note: explicitly not the empty object.
		vReservationAffinity = &WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigReservationAffinity{}
	}
	if err := extractWorkflowTemplatePlacementManagedClusterConfigGceClusterConfigReservationAffinityFields(r, vReservationAffinity); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vReservationAffinity) {
		o.ReservationAffinity = vReservationAffinity
	}
	vNodeGroupAffinity := o.NodeGroupAffinity
	if vNodeGroupAffinity == nil {
		// note: explicitly not the empty object.
		vNodeGroupAffinity = &WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigNodeGroupAffinity{}
	}
	if err := extractWorkflowTemplatePlacementManagedClusterConfigGceClusterConfigNodeGroupAffinityFields(r, vNodeGroupAffinity); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vNodeGroupAffinity) {
		o.NodeGroupAffinity = vNodeGroupAffinity
	}
	vShieldedInstanceConfig := o.ShieldedInstanceConfig
	if vShieldedInstanceConfig == nil {
		// note: explicitly not the empty object.
		vShieldedInstanceConfig = &WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigShieldedInstanceConfig{}
	}
	if err := extractWorkflowTemplatePlacementManagedClusterConfigGceClusterConfigShieldedInstanceConfigFields(r, vShieldedInstanceConfig); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vShieldedInstanceConfig) {
		o.ShieldedInstanceConfig = vShieldedInstanceConfig
	}
	return nil
}
func postReadExtractWorkflowTemplatePlacementManagedClusterConfigGceClusterConfigReservationAffinityFields(r *WorkflowTemplate, o *WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigReservationAffinity) error {
	return nil
}
func postReadExtractWorkflowTemplatePlacementManagedClusterConfigGceClusterConfigNodeGroupAffinityFields(r *WorkflowTemplate, o *WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigNodeGroupAffinity) error {
	return nil
}
func postReadExtractWorkflowTemplatePlacementManagedClusterConfigGceClusterConfigShieldedInstanceConfigFields(r *WorkflowTemplate, o *WorkflowTemplatePlacementManagedClusterConfigGceClusterConfigShieldedInstanceConfig) error {
	return nil
}
func postReadExtractWorkflowTemplatePlacementManagedClusterConfigMasterConfigFields(r *WorkflowTemplate, o *WorkflowTemplatePlacementManagedClusterConfigMasterConfig) error {
	vDiskConfig := o.DiskConfig
	if vDiskConfig == nil {
		// note: explicitly not the empty object.
		vDiskConfig = &WorkflowTemplatePlacementManagedClusterConfigMasterConfigDiskConfig{}
	}
	if err := extractWorkflowTemplatePlacementManagedClusterConfigMasterConfigDiskConfigFields(r, vDiskConfig); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vDiskConfig) {
		o.DiskConfig = vDiskConfig
	}
	vManagedGroupConfig := o.ManagedGroupConfig
	if vManagedGroupConfig == nil {
		// note: explicitly not the empty object.
		vManagedGroupConfig = &WorkflowTemplatePlacementManagedClusterConfigMasterConfigManagedGroupConfig{}
	}
	if err := extractWorkflowTemplatePlacementManagedClusterConfigMasterConfigManagedGroupConfigFields(r, vManagedGroupConfig); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vManagedGroupConfig) {
		o.ManagedGroupConfig = vManagedGroupConfig
	}
	return nil
}
func postReadExtractWorkflowTemplatePlacementManagedClusterConfigMasterConfigDiskConfigFields(r *WorkflowTemplate, o *WorkflowTemplatePlacementManagedClusterConfigMasterConfigDiskConfig) error {
	return nil
}
func postReadExtractWorkflowTemplatePlacementManagedClusterConfigMasterConfigManagedGroupConfigFields(r *WorkflowTemplate, o *WorkflowTemplatePlacementManagedClusterConfigMasterConfigManagedGroupConfig) error {
	return nil
}
func postReadExtractWorkflowTemplatePlacementManagedClusterConfigMasterConfigAcceleratorsFields(r *WorkflowTemplate, o *WorkflowTemplatePlacementManagedClusterConfigMasterConfigAccelerators) error {
	return nil
}
func postReadExtractWorkflowTemplatePlacementManagedClusterConfigWorkerConfigFields(r *WorkflowTemplate, o *WorkflowTemplatePlacementManagedClusterConfigWorkerConfig) error {
	vDiskConfig := o.DiskConfig
	if vDiskConfig == nil {
		// note: explicitly not the empty object.
		vDiskConfig = &WorkflowTemplatePlacementManagedClusterConfigWorkerConfigDiskConfig{}
	}
	if err := extractWorkflowTemplatePlacementManagedClusterConfigWorkerConfigDiskConfigFields(r, vDiskConfig); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vDiskConfig) {
		o.DiskConfig = vDiskConfig
	}
	vManagedGroupConfig := o.ManagedGroupConfig
	if vManagedGroupConfig == nil {
		// note: explicitly not the empty object.
		vManagedGroupConfig = &WorkflowTemplatePlacementManagedClusterConfigWorkerConfigManagedGroupConfig{}
	}
	if err := extractWorkflowTemplatePlacementManagedClusterConfigWorkerConfigManagedGroupConfigFields(r, vManagedGroupConfig); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vManagedGroupConfig) {
		o.ManagedGroupConfig = vManagedGroupConfig
	}
	return nil
}
func postReadExtractWorkflowTemplatePlacementManagedClusterConfigWorkerConfigDiskConfigFields(r *WorkflowTemplate, o *WorkflowTemplatePlacementManagedClusterConfigWorkerConfigDiskConfig) error {
	return nil
}
func postReadExtractWorkflowTemplatePlacementManagedClusterConfigWorkerConfigManagedGroupConfigFields(r *WorkflowTemplate, o *WorkflowTemplatePlacementManagedClusterConfigWorkerConfigManagedGroupConfig) error {
	return nil
}
func postReadExtractWorkflowTemplatePlacementManagedClusterConfigWorkerConfigAcceleratorsFields(r *WorkflowTemplate, o *WorkflowTemplatePlacementManagedClusterConfigWorkerConfigAccelerators) error {
	return nil
}
func postReadExtractWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigFields(r *WorkflowTemplate, o *WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfig) error {
	vDiskConfig := o.DiskConfig
	if vDiskConfig == nil {
		// note: explicitly not the empty object.
		vDiskConfig = &WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigDiskConfig{}
	}
	if err := extractWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigDiskConfigFields(r, vDiskConfig); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vDiskConfig) {
		o.DiskConfig = vDiskConfig
	}
	vManagedGroupConfig := o.ManagedGroupConfig
	if vManagedGroupConfig == nil {
		// note: explicitly not the empty object.
		vManagedGroupConfig = &WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigManagedGroupConfig{}
	}
	if err := extractWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigManagedGroupConfigFields(r, vManagedGroupConfig); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vManagedGroupConfig) {
		o.ManagedGroupConfig = vManagedGroupConfig
	}
	return nil
}
func postReadExtractWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigDiskConfigFields(r *WorkflowTemplate, o *WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigDiskConfig) error {
	return nil
}
func postReadExtractWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigManagedGroupConfigFields(r *WorkflowTemplate, o *WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigManagedGroupConfig) error {
	return nil
}
func postReadExtractWorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigAcceleratorsFields(r *WorkflowTemplate, o *WorkflowTemplatePlacementManagedClusterConfigSecondaryWorkerConfigAccelerators) error {
	return nil
}
func postReadExtractWorkflowTemplatePlacementManagedClusterConfigSoftwareConfigFields(r *WorkflowTemplate, o *WorkflowTemplatePlacementManagedClusterConfigSoftwareConfig) error {
	return nil
}
func postReadExtractWorkflowTemplatePlacementManagedClusterConfigInitializationActionsFields(r *WorkflowTemplate, o *WorkflowTemplatePlacementManagedClusterConfigInitializationActions) error {
	return nil
}
func postReadExtractWorkflowTemplatePlacementManagedClusterConfigEncryptionConfigFields(r *WorkflowTemplate, o *WorkflowTemplatePlacementManagedClusterConfigEncryptionConfig) error {
	return nil
}
func postReadExtractWorkflowTemplatePlacementManagedClusterConfigAutoscalingConfigFields(r *WorkflowTemplate, o *WorkflowTemplatePlacementManagedClusterConfigAutoscalingConfig) error {
	return nil
}
func postReadExtractWorkflowTemplatePlacementManagedClusterConfigSecurityConfigFields(r *WorkflowTemplate, o *WorkflowTemplatePlacementManagedClusterConfigSecurityConfig) error {
	vKerberosConfig := o.KerberosConfig
	if vKerberosConfig == nil {
		// note: explicitly not the empty object.
		vKerberosConfig = &WorkflowTemplatePlacementManagedClusterConfigSecurityConfigKerberosConfig{}
	}
	if err := extractWorkflowTemplatePlacementManagedClusterConfigSecurityConfigKerberosConfigFields(r, vKerberosConfig); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vKerberosConfig) {
		o.KerberosConfig = vKerberosConfig
	}
	return nil
}
func postReadExtractWorkflowTemplatePlacementManagedClusterConfigSecurityConfigKerberosConfigFields(r *WorkflowTemplate, o *WorkflowTemplatePlacementManagedClusterConfigSecurityConfigKerberosConfig) error {
	return nil
}
func postReadExtractWorkflowTemplatePlacementManagedClusterConfigLifecycleConfigFields(r *WorkflowTemplate, o *WorkflowTemplatePlacementManagedClusterConfigLifecycleConfig) error {
	return nil
}
func postReadExtractWorkflowTemplatePlacementManagedClusterConfigEndpointConfigFields(r *WorkflowTemplate, o *WorkflowTemplatePlacementManagedClusterConfigEndpointConfig) error {
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
	if !dcl.IsEmptyValueIndirect(vHadoopJob) {
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
	if !dcl.IsEmptyValueIndirect(vSparkJob) {
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
	if !dcl.IsEmptyValueIndirect(vPysparkJob) {
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
	if !dcl.IsEmptyValueIndirect(vHiveJob) {
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
	if !dcl.IsEmptyValueIndirect(vPigJob) {
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
	if !dcl.IsEmptyValueIndirect(vSparkRJob) {
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
	if !dcl.IsEmptyValueIndirect(vSparkSqlJob) {
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
	if !dcl.IsEmptyValueIndirect(vPrestoJob) {
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
	if !dcl.IsEmptyValueIndirect(vScheduling) {
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
	if !dcl.IsEmptyValueIndirect(vLoggingConfig) {
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
	if !dcl.IsEmptyValueIndirect(vLoggingConfig) {
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
	if !dcl.IsEmptyValueIndirect(vLoggingConfig) {
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
	if !dcl.IsEmptyValueIndirect(vQueryList) {
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
	if !dcl.IsEmptyValueIndirect(vQueryList) {
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
	if !dcl.IsEmptyValueIndirect(vLoggingConfig) {
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
	if !dcl.IsEmptyValueIndirect(vLoggingConfig) {
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
	if !dcl.IsEmptyValueIndirect(vQueryList) {
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
	if !dcl.IsEmptyValueIndirect(vLoggingConfig) {
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
	if !dcl.IsEmptyValueIndirect(vQueryList) {
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
	if !dcl.IsEmptyValueIndirect(vLoggingConfig) {
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
	if !dcl.IsEmptyValueIndirect(vValidation) {
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
	if !dcl.IsEmptyValueIndirect(vRegex) {
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
	if !dcl.IsEmptyValueIndirect(vValues) {
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
