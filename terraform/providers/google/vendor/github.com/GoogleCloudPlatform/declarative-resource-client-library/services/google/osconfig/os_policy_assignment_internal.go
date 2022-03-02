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
	"time"

	"github.com/GoogleCloudPlatform/declarative-resource-client-library/dcl"
	"github.com/GoogleCloudPlatform/declarative-resource-client-library/dcl/operations"
)

func (r *OSPolicyAssignment) validate() error {

	if err := dcl.Required(r, "name"); err != nil {
		return err
	}
	if err := dcl.Required(r, "osPolicies"); err != nil {
		return err
	}
	if err := dcl.Required(r, "instanceFilter"); err != nil {
		return err
	}
	if err := dcl.Required(r, "rollout"); err != nil {
		return err
	}
	if err := dcl.RequiredParameter(r.Project, "Project"); err != nil {
		return err
	}
	if err := dcl.RequiredParameter(r.Location, "Location"); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(r.InstanceFilter) {
		if err := r.InstanceFilter.validate(); err != nil {
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
func (r *OSPolicyAssignmentOSPolicies) validate() error {
	if err := dcl.Required(r, "id"); err != nil {
		return err
	}
	if err := dcl.Required(r, "mode"); err != nil {
		return err
	}
	if err := dcl.Required(r, "resourceGroups"); err != nil {
		return err
	}
	return nil
}
func (r *OSPolicyAssignmentOSPoliciesResourceGroups) validate() error {
	if err := dcl.Required(r, "resources"); err != nil {
		return err
	}
	return nil
}
func (r *OSPolicyAssignmentOSPoliciesResourceGroupsInventoryFilters) validate() error {
	if err := dcl.Required(r, "osShortName"); err != nil {
		return err
	}
	return nil
}
func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResources) validate() error {
	if err := dcl.Required(r, "id"); err != nil {
		return err
	}
	if err := dcl.ValidateExactlyOneOfFieldsSet([]string{"Pkg", "Repository", "Exec", "File"}, r.Pkg, r.Repository, r.Exec, r.File); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(r.Pkg) {
		if err := r.Pkg.validate(); err != nil {
			return err
		}
	}
	if !dcl.IsEmptyValueIndirect(r.Repository) {
		if err := r.Repository.validate(); err != nil {
			return err
		}
	}
	if !dcl.IsEmptyValueIndirect(r.Exec) {
		if err := r.Exec.validate(); err != nil {
			return err
		}
	}
	if !dcl.IsEmptyValueIndirect(r.File) {
		if err := r.File.validate(); err != nil {
			return err
		}
	}
	return nil
}
func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkg) validate() error {
	if err := dcl.Required(r, "desiredState"); err != nil {
		return err
	}
	if err := dcl.ValidateAtMostOneOfFieldsSet([]string{"Apt", "Deb", "Yum", "Zypper", "Rpm", "Googet", "Msi"}, r.Apt, r.Deb, r.Yum, r.Zypper, r.Rpm, r.Googet, r.Msi); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(r.Apt) {
		if err := r.Apt.validate(); err != nil {
			return err
		}
	}
	if !dcl.IsEmptyValueIndirect(r.Deb) {
		if err := r.Deb.validate(); err != nil {
			return err
		}
	}
	if !dcl.IsEmptyValueIndirect(r.Yum) {
		if err := r.Yum.validate(); err != nil {
			return err
		}
	}
	if !dcl.IsEmptyValueIndirect(r.Zypper) {
		if err := r.Zypper.validate(); err != nil {
			return err
		}
	}
	if !dcl.IsEmptyValueIndirect(r.Rpm) {
		if err := r.Rpm.validate(); err != nil {
			return err
		}
	}
	if !dcl.IsEmptyValueIndirect(r.Googet) {
		if err := r.Googet.validate(); err != nil {
			return err
		}
	}
	if !dcl.IsEmptyValueIndirect(r.Msi) {
		if err := r.Msi.validate(); err != nil {
			return err
		}
	}
	return nil
}
func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgApt) validate() error {
	if err := dcl.Required(r, "name"); err != nil {
		return err
	}
	return nil
}
func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDeb) validate() error {
	if err := dcl.Required(r, "source"); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(r.Source) {
		if err := r.Source.validate(); err != nil {
			return err
		}
	}
	return nil
}
func (r *OSPolicyAssignmentFile) validate() error {
	if err := dcl.ValidateAtMostOneOfFieldsSet([]string{"Remote", "Gcs", "LocalPath"}, r.Remote, r.Gcs, r.LocalPath); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(r.Remote) {
		if err := r.Remote.validate(); err != nil {
			return err
		}
	}
	if !dcl.IsEmptyValueIndirect(r.Gcs) {
		if err := r.Gcs.validate(); err != nil {
			return err
		}
	}
	return nil
}
func (r *OSPolicyAssignmentFileRemote) validate() error {
	if err := dcl.Required(r, "uri"); err != nil {
		return err
	}
	return nil
}
func (r *OSPolicyAssignmentFileGcs) validate() error {
	if err := dcl.Required(r, "bucket"); err != nil {
		return err
	}
	if err := dcl.Required(r, "object"); err != nil {
		return err
	}
	return nil
}
func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgYum) validate() error {
	if err := dcl.Required(r, "name"); err != nil {
		return err
	}
	return nil
}
func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgZypper) validate() error {
	if err := dcl.Required(r, "name"); err != nil {
		return err
	}
	return nil
}
func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpm) validate() error {
	if err := dcl.Required(r, "source"); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(r.Source) {
		if err := r.Source.validate(); err != nil {
			return err
		}
	}
	return nil
}
func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgGooget) validate() error {
	if err := dcl.Required(r, "name"); err != nil {
		return err
	}
	return nil
}
func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsi) validate() error {
	if err := dcl.Required(r, "source"); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(r.Source) {
		if err := r.Source.validate(); err != nil {
			return err
		}
	}
	return nil
}
func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepository) validate() error {
	if err := dcl.ValidateAtMostOneOfFieldsSet([]string{"Apt", "Yum", "Zypper", "Goo"}, r.Apt, r.Yum, r.Zypper, r.Goo); err != nil {
		return err
	}
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
	if !dcl.IsEmptyValueIndirect(r.Zypper) {
		if err := r.Zypper.validate(); err != nil {
			return err
		}
	}
	if !dcl.IsEmptyValueIndirect(r.Goo) {
		if err := r.Goo.validate(); err != nil {
			return err
		}
	}
	return nil
}
func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryApt) validate() error {
	if err := dcl.Required(r, "archiveType"); err != nil {
		return err
	}
	if err := dcl.Required(r, "uri"); err != nil {
		return err
	}
	if err := dcl.Required(r, "distribution"); err != nil {
		return err
	}
	if err := dcl.Required(r, "components"); err != nil {
		return err
	}
	return nil
}
func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryYum) validate() error {
	if err := dcl.Required(r, "id"); err != nil {
		return err
	}
	if err := dcl.Required(r, "baseUrl"); err != nil {
		return err
	}
	return nil
}
func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryZypper) validate() error {
	if err := dcl.Required(r, "id"); err != nil {
		return err
	}
	if err := dcl.Required(r, "baseUrl"); err != nil {
		return err
	}
	return nil
}
func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryGoo) validate() error {
	if err := dcl.Required(r, "name"); err != nil {
		return err
	}
	if err := dcl.Required(r, "url"); err != nil {
		return err
	}
	return nil
}
func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExec) validate() error {
	if err := dcl.Required(r, "validate"); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(r.Validate) {
		if err := r.Validate.validate(); err != nil {
			return err
		}
	}
	if !dcl.IsEmptyValueIndirect(r.Enforce) {
		if err := r.Enforce.validate(); err != nil {
			return err
		}
	}
	return nil
}
func (r *OSPolicyAssignmentExec) validate() error {
	if err := dcl.Required(r, "interpreter"); err != nil {
		return err
	}
	if err := dcl.ValidateAtMostOneOfFieldsSet([]string{"File", "Script"}, r.File, r.Script); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(r.File) {
		if err := r.File.validate(); err != nil {
			return err
		}
	}
	return nil
}
func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFile) validate() error {
	if err := dcl.Required(r, "path"); err != nil {
		return err
	}
	if err := dcl.Required(r, "state"); err != nil {
		return err
	}
	if err := dcl.ValidateAtMostOneOfFieldsSet([]string{"File", "Content"}, r.File, r.Content); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(r.File) {
		if err := r.File.validate(); err != nil {
			return err
		}
	}
	return nil
}
func (r *OSPolicyAssignmentInstanceFilter) validate() error {
	return nil
}
func (r *OSPolicyAssignmentInstanceFilterInclusionLabels) validate() error {
	return nil
}
func (r *OSPolicyAssignmentInstanceFilterExclusionLabels) validate() error {
	return nil
}
func (r *OSPolicyAssignmentInstanceFilterInventories) validate() error {
	if err := dcl.Required(r, "osShortName"); err != nil {
		return err
	}
	return nil
}
func (r *OSPolicyAssignmentRollout) validate() error {
	if err := dcl.Required(r, "disruptionBudget"); err != nil {
		return err
	}
	if err := dcl.Required(r, "minWaitDuration"); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(r.DisruptionBudget) {
		if err := r.DisruptionBudget.validate(); err != nil {
			return err
		}
	}
	return nil
}
func (r *OSPolicyAssignmentRolloutDisruptionBudget) validate() error {
	if err := dcl.ValidateExactlyOneOfFieldsSet([]string{"Fixed", "Percent"}, r.Fixed, r.Percent); err != nil {
		return err
	}
	return nil
}
func (r *OSPolicyAssignment) basePath() string {
	params := map[string]interface{}{}
	return dcl.Nprintf("https://osconfig.googleapis.com/v1", params)
}

func (r *OSPolicyAssignment) getURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"project":  dcl.ValueOrEmptyString(nr.Project),
		"location": dcl.ValueOrEmptyString(nr.Location),
		"name":     dcl.ValueOrEmptyString(nr.Name),
	}
	return dcl.URL("projects/{{project}}/locations/{{location}}/osPolicyAssignments/{{name}}", nr.basePath(), userBasePath, params), nil
}

func (r *OSPolicyAssignment) listURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"project":  dcl.ValueOrEmptyString(nr.Project),
		"location": dcl.ValueOrEmptyString(nr.Location),
	}
	return dcl.URL("projects/{{project}}/locations/{{location}}/osPolicyAssignments", nr.basePath(), userBasePath, params), nil

}

func (r *OSPolicyAssignment) createURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"project":  dcl.ValueOrEmptyString(nr.Project),
		"location": dcl.ValueOrEmptyString(nr.Location),
		"name":     dcl.ValueOrEmptyString(nr.Name),
	}
	return dcl.URL("projects/{{project}}/locations/{{location}}/osPolicyAssignments?osPolicyAssignmentId={{name}}", nr.basePath(), userBasePath, params), nil

}

func (r *OSPolicyAssignment) deleteURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"project":  dcl.ValueOrEmptyString(nr.Project),
		"location": dcl.ValueOrEmptyString(nr.Location),
		"name":     dcl.ValueOrEmptyString(nr.Name),
	}
	return dcl.URL("projects/{{project}}/locations/{{location}}/osPolicyAssignments/{{name}}", nr.basePath(), userBasePath, params), nil
}

// oSPolicyAssignmentApiOperation represents a mutable operation in the underlying REST
// API such as Create, Update, or Delete.
type oSPolicyAssignmentApiOperation interface {
	do(context.Context, *OSPolicyAssignment, *Client) error
}

// newUpdateOSPolicyAssignmentUpdateOSPolicyAssignmentRequest creates a request for an
// OSPolicyAssignment resource's UpdateOSPolicyAssignment update type by filling in the update
// fields based on the intended state of the resource.
func newUpdateOSPolicyAssignmentUpdateOSPolicyAssignmentRequest(ctx context.Context, f *OSPolicyAssignment, c *Client) (map[string]interface{}, error) {
	req := map[string]interface{}{}

	if v := f.Description; !dcl.IsEmptyValueIndirect(v) {
		req["description"] = v
	}
	if v, err := expandOSPolicyAssignmentOSPoliciesSlice(c, f.OSPolicies); err != nil {
		return nil, fmt.Errorf("error expanding OSPolicies into osPolicies: %w", err)
	} else if v != nil {
		req["osPolicies"] = v
	}
	if v, err := expandOSPolicyAssignmentInstanceFilter(c, f.InstanceFilter); err != nil {
		return nil, fmt.Errorf("error expanding InstanceFilter into instanceFilter: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		req["instanceFilter"] = v
	}
	if v, err := expandOSPolicyAssignmentRollout(c, f.Rollout); err != nil {
		return nil, fmt.Errorf("error expanding Rollout into rollout: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		req["rollout"] = v
	}
	b, err := c.getOSPolicyAssignmentRaw(ctx, f)
	if err != nil {
		return nil, err
	}
	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, err
	}
	rawEtag, err := dcl.GetMapEntry(
		m,
		[]string{"etag"},
	)
	if err != nil {
		c.Config.Logger.WarningWithContextf(ctx, "Failed to fetch from JSON Path: %v", err)
	} else {
		req["etag"] = rawEtag.(string)
	}
	return req, nil
}

// marshalUpdateOSPolicyAssignmentUpdateOSPolicyAssignmentRequest converts the update into
// the final JSON request body.
func marshalUpdateOSPolicyAssignmentUpdateOSPolicyAssignmentRequest(c *Client, m map[string]interface{}) ([]byte, error) {

	return json.Marshal(m)
}

type updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation struct {
	// If the update operation has the REQUIRES_APPLY_OPTIONS trait, this will be populated.
	// Usually it will be nil - this is to prevent us from accidentally depending on apply
	// options, which should usually be unnecessary.
	ApplyOptions []dcl.ApplyOption
	FieldDiffs   []*dcl.FieldDiff
}

// do creates a request and sends it to the appropriate URL. In most operations,
// do will transcribe a subset of the resource into a request object and send a
// PUT request to a single URL.

func (op *updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation) do(ctx context.Context, r *OSPolicyAssignment, c *Client) error {
	_, err := c.GetOSPolicyAssignment(ctx, r)
	if err != nil {
		return err
	}

	u, err := r.updateURL(c.Config.BasePath, "UpdateOSPolicyAssignment")
	if err != nil {
		return err
	}
	mask := dcl.TopLevelUpdateMask(op.FieldDiffs)
	u, err = dcl.AddQueryParams(u, map[string]string{"updateMask": mask})
	if err != nil {
		return err
	}

	req, err := newUpdateOSPolicyAssignmentUpdateOSPolicyAssignmentRequest(ctx, r, c)
	if err != nil {
		return err
	}

	c.Config.Logger.InfoWithContextf(ctx, "Created update: %#v", req)
	body, err := marshalUpdateOSPolicyAssignmentUpdateOSPolicyAssignmentRequest(c, req)
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

func (c *Client) listOSPolicyAssignmentRaw(ctx context.Context, r *OSPolicyAssignment, pageToken string, pageSize int32) ([]byte, error) {
	u, err := r.urlNormalized().listURL(c.Config.BasePath)
	if err != nil {
		return nil, err
	}

	m := make(map[string]string)
	if pageToken != "" {
		m["pageToken"] = pageToken
	}

	if pageSize != OSPolicyAssignmentMaxPage {
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

type listOSPolicyAssignmentOperation struct {
	OsPolicyAssignments []map[string]interface{} `json:"osPolicyAssignments"`
	Token               string                   `json:"nextPageToken"`
}

func (c *Client) listOSPolicyAssignment(ctx context.Context, r *OSPolicyAssignment, pageToken string, pageSize int32) ([]*OSPolicyAssignment, string, error) {
	b, err := c.listOSPolicyAssignmentRaw(ctx, r, pageToken, pageSize)
	if err != nil {
		return nil, "", err
	}

	var m listOSPolicyAssignmentOperation
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, "", err
	}

	var l []*OSPolicyAssignment
	for _, v := range m.OsPolicyAssignments {
		res, err := unmarshalMapOSPolicyAssignment(v, c)
		if err != nil {
			return nil, m.Token, err
		}
		res.Project = r.Project
		res.Location = r.Location
		l = append(l, res)
	}

	return l, m.Token, nil
}

func (c *Client) deleteAllOSPolicyAssignment(ctx context.Context, f func(*OSPolicyAssignment) bool, resources []*OSPolicyAssignment) error {
	var errors []string
	for _, res := range resources {
		if f(res) {
			// We do not want deleteAll to fail on a deletion or else it will stop deleting other resources.
			err := c.DeleteOSPolicyAssignment(ctx, res)
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

type deleteOSPolicyAssignmentOperation struct{}

func (op *deleteOSPolicyAssignmentOperation) do(ctx context.Context, r *OSPolicyAssignment, c *Client) error {
	r, err := c.GetOSPolicyAssignment(ctx, r)
	if err != nil {
		if dcl.IsNotFound(err) {
			c.Config.Logger.InfoWithContextf(ctx, "OSPolicyAssignment not found, returning. Original error: %v", err)
			return nil
		}
		c.Config.Logger.WarningWithContextf(ctx, "GetOSPolicyAssignment checking for existence. error: %v", err)
		return err
	}

	err = r.waitForNotReconciling(ctx, c)
	if err != nil {
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

	// we saw a race condition where for some successful delete operation, the Get calls returned resources for a short duration.
	// this is the reason we are adding retry to handle that case.
	maxRetry := 10
	for i := 1; i <= maxRetry; i++ {
		_, err = c.GetOSPolicyAssignment(ctx, r)
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
type createOSPolicyAssignmentOperation struct {
	response map[string]interface{}
}

func (op *createOSPolicyAssignmentOperation) FirstResponse() (map[string]interface{}, bool) {
	return op.response, len(op.response) > 0
}

func (op *createOSPolicyAssignmentOperation) do(ctx context.Context, r *OSPolicyAssignment, c *Client) error {
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

	if _, err := c.GetOSPolicyAssignment(ctx, r); err != nil {
		c.Config.Logger.WarningWithContextf(ctx, "get returned error: %v", err)
		return err
	}

	return nil
}

func (c *Client) getOSPolicyAssignmentRaw(ctx context.Context, r *OSPolicyAssignment) ([]byte, error) {

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

func (c *Client) oSPolicyAssignmentDiffsForRawDesired(ctx context.Context, rawDesired *OSPolicyAssignment, opts ...dcl.ApplyOption) (initial, desired *OSPolicyAssignment, diffs []*dcl.FieldDiff, err error) {
	c.Config.Logger.InfoWithContext(ctx, "Fetching initial state...")
	// First, let us see if the user provided a state hint.  If they did, we will start fetching based on that.
	var fetchState *OSPolicyAssignment
	if sh := dcl.FetchStateHint(opts); sh != nil {
		if r, ok := sh.(*OSPolicyAssignment); !ok {
			c.Config.Logger.WarningWithContextf(ctx, "Initial state hint was of the wrong type; expected OSPolicyAssignment, got %T", sh)
		} else {
			fetchState = r
		}
	}
	if fetchState == nil {
		fetchState = rawDesired
	}

	// 1.2: Retrieval of raw initial state from API
	rawInitial, err := c.GetOSPolicyAssignment(ctx, fetchState)
	if rawInitial == nil {
		if !dcl.IsNotFound(err) {
			c.Config.Logger.WarningWithContextf(ctx, "Failed to retrieve whether a OSPolicyAssignment resource already exists: %s", err)
			return nil, nil, nil, fmt.Errorf("failed to retrieve OSPolicyAssignment resource: %v", err)
		}
		c.Config.Logger.InfoWithContext(ctx, "Found that OSPolicyAssignment resource did not exist.")
		// Perform canonicalization to pick up defaults.
		desired, err = canonicalizeOSPolicyAssignmentDesiredState(rawDesired, rawInitial)
		return nil, desired, nil, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Found initial state for OSPolicyAssignment: %v", rawInitial)
	c.Config.Logger.InfoWithContextf(ctx, "Initial desired state for OSPolicyAssignment: %v", rawDesired)

	// 1.3: Canonicalize raw initial state into initial state.
	initial, err = canonicalizeOSPolicyAssignmentInitialState(rawInitial, rawDesired)
	if err != nil {
		return nil, nil, nil, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Canonicalized initial state for OSPolicyAssignment: %v", initial)

	// 1.4: Canonicalize raw desired state into desired state.
	desired, err = canonicalizeOSPolicyAssignmentDesiredState(rawDesired, rawInitial, opts...)
	if err != nil {
		return nil, nil, nil, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Canonicalized desired state for OSPolicyAssignment: %v", desired)

	// 2.1: Comparison of initial and desired state.
	diffs, err = diffOSPolicyAssignment(c, desired, initial, opts...)
	return initial, desired, diffs, err
}

func canonicalizeOSPolicyAssignmentInitialState(rawInitial, rawDesired *OSPolicyAssignment) (*OSPolicyAssignment, error) {
	// TODO(magic-modules-eng): write canonicalizer once relevant traits are added.
	return rawInitial, nil
}

/*
* Canonicalizers
*
* These are responsible for converting either a user-specified config or a
* GCP API response to a standard format that can be used for difference checking.
* */

func canonicalizeOSPolicyAssignmentDesiredState(rawDesired, rawInitial *OSPolicyAssignment, opts ...dcl.ApplyOption) (*OSPolicyAssignment, error) {

	if rawInitial == nil {
		// Since the initial state is empty, the desired state is all we have.
		// We canonicalize the remaining nested objects with nil to pick up defaults.
		rawDesired.InstanceFilter = canonicalizeOSPolicyAssignmentInstanceFilter(rawDesired.InstanceFilter, nil, opts...)
		rawDesired.Rollout = canonicalizeOSPolicyAssignmentRollout(rawDesired.Rollout, nil, opts...)

		return rawDesired, nil
	}
	canonicalDesired := &OSPolicyAssignment{}
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
	canonicalDesired.OSPolicies = canonicalizeOSPolicyAssignmentOSPoliciesSlice(rawDesired.OSPolicies, rawInitial.OSPolicies, opts...)
	canonicalDesired.InstanceFilter = canonicalizeOSPolicyAssignmentInstanceFilter(rawDesired.InstanceFilter, rawInitial.InstanceFilter, opts...)
	canonicalDesired.Rollout = canonicalizeOSPolicyAssignmentRollout(rawDesired.Rollout, rawInitial.Rollout, opts...)
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

func canonicalizeOSPolicyAssignmentNewState(c *Client, rawNew, rawDesired *OSPolicyAssignment) (*OSPolicyAssignment, error) {

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

	if dcl.IsNotReturnedByServer(rawNew.OSPolicies) && dcl.IsNotReturnedByServer(rawDesired.OSPolicies) {
		rawNew.OSPolicies = rawDesired.OSPolicies
	} else {
		rawNew.OSPolicies = canonicalizeNewOSPolicyAssignmentOSPoliciesSlice(c, rawDesired.OSPolicies, rawNew.OSPolicies)
	}

	if dcl.IsNotReturnedByServer(rawNew.InstanceFilter) && dcl.IsNotReturnedByServer(rawDesired.InstanceFilter) {
		rawNew.InstanceFilter = rawDesired.InstanceFilter
	} else {
		rawNew.InstanceFilter = canonicalizeNewOSPolicyAssignmentInstanceFilter(c, rawDesired.InstanceFilter, rawNew.InstanceFilter)
	}

	if dcl.IsNotReturnedByServer(rawNew.Rollout) && dcl.IsNotReturnedByServer(rawDesired.Rollout) {
		rawNew.Rollout = rawDesired.Rollout
	} else {
		rawNew.Rollout = canonicalizeNewOSPolicyAssignmentRollout(c, rawDesired.Rollout, rawNew.Rollout)
	}

	if dcl.IsNotReturnedByServer(rawNew.RevisionId) && dcl.IsNotReturnedByServer(rawDesired.RevisionId) {
		rawNew.RevisionId = rawDesired.RevisionId
	} else {
		if dcl.StringCanonicalize(rawDesired.RevisionId, rawNew.RevisionId) {
			rawNew.RevisionId = rawDesired.RevisionId
		}
	}

	if dcl.IsNotReturnedByServer(rawNew.RevisionCreateTime) && dcl.IsNotReturnedByServer(rawDesired.RevisionCreateTime) {
		rawNew.RevisionCreateTime = rawDesired.RevisionCreateTime
	} else {
	}

	if dcl.IsNotReturnedByServer(rawNew.Etag) && dcl.IsNotReturnedByServer(rawDesired.Etag) {
		rawNew.Etag = rawDesired.Etag
	} else {
		if dcl.StringCanonicalize(rawDesired.Etag, rawNew.Etag) {
			rawNew.Etag = rawDesired.Etag
		}
	}

	if dcl.IsNotReturnedByServer(rawNew.RolloutState) && dcl.IsNotReturnedByServer(rawDesired.RolloutState) {
		rawNew.RolloutState = rawDesired.RolloutState
	} else {
	}

	if dcl.IsNotReturnedByServer(rawNew.Baseline) && dcl.IsNotReturnedByServer(rawDesired.Baseline) {
		rawNew.Baseline = rawDesired.Baseline
	} else {
		if dcl.BoolCanonicalize(rawDesired.Baseline, rawNew.Baseline) {
			rawNew.Baseline = rawDesired.Baseline
		}
	}

	if dcl.IsNotReturnedByServer(rawNew.Deleted) && dcl.IsNotReturnedByServer(rawDesired.Deleted) {
		rawNew.Deleted = rawDesired.Deleted
	} else {
		if dcl.BoolCanonicalize(rawDesired.Deleted, rawNew.Deleted) {
			rawNew.Deleted = rawDesired.Deleted
		}
	}

	if dcl.IsNotReturnedByServer(rawNew.Reconciling) && dcl.IsNotReturnedByServer(rawDesired.Reconciling) {
		rawNew.Reconciling = rawDesired.Reconciling
	} else {
		if dcl.BoolCanonicalize(rawDesired.Reconciling, rawNew.Reconciling) {
			rawNew.Reconciling = rawDesired.Reconciling
		}
	}

	if dcl.IsNotReturnedByServer(rawNew.Uid) && dcl.IsNotReturnedByServer(rawDesired.Uid) {
		rawNew.Uid = rawDesired.Uid
	} else {
		if dcl.StringCanonicalize(rawDesired.Uid, rawNew.Uid) {
			rawNew.Uid = rawDesired.Uid
		}
	}

	rawNew.Project = rawDesired.Project

	rawNew.Location = rawDesired.Location

	return rawNew, nil
}

func canonicalizeOSPolicyAssignmentOSPolicies(des, initial *OSPolicyAssignmentOSPolicies, opts ...dcl.ApplyOption) *OSPolicyAssignmentOSPolicies {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &OSPolicyAssignmentOSPolicies{}

	if dcl.StringCanonicalize(des.Id, initial.Id) || dcl.IsZeroValue(des.Id) {
		cDes.Id = initial.Id
	} else {
		cDes.Id = des.Id
	}
	if dcl.StringCanonicalize(des.Description, initial.Description) || dcl.IsZeroValue(des.Description) {
		cDes.Description = initial.Description
	} else {
		cDes.Description = des.Description
	}
	if dcl.IsZeroValue(des.Mode) {
		cDes.Mode = initial.Mode
	} else {
		cDes.Mode = des.Mode
	}
	cDes.ResourceGroups = canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsSlice(des.ResourceGroups, initial.ResourceGroups, opts...)
	if dcl.BoolCanonicalize(des.AllowNoResourceGroupMatch, initial.AllowNoResourceGroupMatch) || dcl.IsZeroValue(des.AllowNoResourceGroupMatch) {
		cDes.AllowNoResourceGroupMatch = initial.AllowNoResourceGroupMatch
	} else {
		cDes.AllowNoResourceGroupMatch = des.AllowNoResourceGroupMatch
	}

	return cDes
}

func canonicalizeOSPolicyAssignmentOSPoliciesSlice(des, initial []OSPolicyAssignmentOSPolicies, opts ...dcl.ApplyOption) []OSPolicyAssignmentOSPolicies {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]OSPolicyAssignmentOSPolicies, 0, len(des))
		for _, d := range des {
			cd := canonicalizeOSPolicyAssignmentOSPolicies(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]OSPolicyAssignmentOSPolicies, 0, len(des))
	for i, d := range des {
		cd := canonicalizeOSPolicyAssignmentOSPolicies(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewOSPolicyAssignmentOSPolicies(c *Client, des, nw *OSPolicyAssignmentOSPolicies) *OSPolicyAssignmentOSPolicies {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsNotReturnedByServer(des) {
			c.Config.Logger.Info("Found explicitly empty value for OSPolicyAssignmentOSPolicies while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringCanonicalize(des.Id, nw.Id) {
		nw.Id = des.Id
	}
	if dcl.StringCanonicalize(des.Description, nw.Description) {
		nw.Description = des.Description
	}
	nw.ResourceGroups = canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsSlice(c, des.ResourceGroups, nw.ResourceGroups)
	if dcl.BoolCanonicalize(des.AllowNoResourceGroupMatch, nw.AllowNoResourceGroupMatch) {
		nw.AllowNoResourceGroupMatch = des.AllowNoResourceGroupMatch
	}

	return nw
}

func canonicalizeNewOSPolicyAssignmentOSPoliciesSet(c *Client, des, nw []OSPolicyAssignmentOSPolicies) []OSPolicyAssignmentOSPolicies {
	if des == nil {
		return nw
	}
	var reorderedNew []OSPolicyAssignmentOSPolicies
	for _, d := range des {
		matchedNew := -1
		for idx, n := range nw {
			if diffs, _ := compareOSPolicyAssignmentOSPoliciesNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
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

func canonicalizeNewOSPolicyAssignmentOSPoliciesSlice(c *Client, des, nw []OSPolicyAssignmentOSPolicies) []OSPolicyAssignmentOSPolicies {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []OSPolicyAssignmentOSPolicies
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewOSPolicyAssignmentOSPolicies(c, &d, &n))
	}

	return items
}

func canonicalizeOSPolicyAssignmentOSPoliciesResourceGroups(des, initial *OSPolicyAssignmentOSPoliciesResourceGroups, opts ...dcl.ApplyOption) *OSPolicyAssignmentOSPoliciesResourceGroups {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &OSPolicyAssignmentOSPoliciesResourceGroups{}

	cDes.InventoryFilters = canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsInventoryFiltersSlice(des.InventoryFilters, initial.InventoryFilters, opts...)
	cDes.Resources = canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesSlice(des.Resources, initial.Resources, opts...)

	return cDes
}

func canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsSlice(des, initial []OSPolicyAssignmentOSPoliciesResourceGroups, opts ...dcl.ApplyOption) []OSPolicyAssignmentOSPoliciesResourceGroups {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]OSPolicyAssignmentOSPoliciesResourceGroups, 0, len(des))
		for _, d := range des {
			cd := canonicalizeOSPolicyAssignmentOSPoliciesResourceGroups(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]OSPolicyAssignmentOSPoliciesResourceGroups, 0, len(des))
	for i, d := range des {
		cd := canonicalizeOSPolicyAssignmentOSPoliciesResourceGroups(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroups(c *Client, des, nw *OSPolicyAssignmentOSPoliciesResourceGroups) *OSPolicyAssignmentOSPoliciesResourceGroups {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsNotReturnedByServer(des) {
			c.Config.Logger.Info("Found explicitly empty value for OSPolicyAssignmentOSPoliciesResourceGroups while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	nw.InventoryFilters = canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsInventoryFiltersSlice(c, des.InventoryFilters, nw.InventoryFilters)
	nw.Resources = canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesSlice(c, des.Resources, nw.Resources)

	return nw
}

func canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsSet(c *Client, des, nw []OSPolicyAssignmentOSPoliciesResourceGroups) []OSPolicyAssignmentOSPoliciesResourceGroups {
	if des == nil {
		return nw
	}
	var reorderedNew []OSPolicyAssignmentOSPoliciesResourceGroups
	for _, d := range des {
		matchedNew := -1
		for idx, n := range nw {
			if diffs, _ := compareOSPolicyAssignmentOSPoliciesResourceGroupsNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
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

func canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsSlice(c *Client, des, nw []OSPolicyAssignmentOSPoliciesResourceGroups) []OSPolicyAssignmentOSPoliciesResourceGroups {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []OSPolicyAssignmentOSPoliciesResourceGroups
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroups(c, &d, &n))
	}

	return items
}

func canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsInventoryFilters(des, initial *OSPolicyAssignmentOSPoliciesResourceGroupsInventoryFilters, opts ...dcl.ApplyOption) *OSPolicyAssignmentOSPoliciesResourceGroupsInventoryFilters {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &OSPolicyAssignmentOSPoliciesResourceGroupsInventoryFilters{}

	if dcl.StringCanonicalize(des.OSShortName, initial.OSShortName) || dcl.IsZeroValue(des.OSShortName) {
		cDes.OSShortName = initial.OSShortName
	} else {
		cDes.OSShortName = des.OSShortName
	}
	if dcl.StringCanonicalize(des.OSVersion, initial.OSVersion) || dcl.IsZeroValue(des.OSVersion) {
		cDes.OSVersion = initial.OSVersion
	} else {
		cDes.OSVersion = des.OSVersion
	}

	return cDes
}

func canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsInventoryFiltersSlice(des, initial []OSPolicyAssignmentOSPoliciesResourceGroupsInventoryFilters, opts ...dcl.ApplyOption) []OSPolicyAssignmentOSPoliciesResourceGroupsInventoryFilters {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]OSPolicyAssignmentOSPoliciesResourceGroupsInventoryFilters, 0, len(des))
		for _, d := range des {
			cd := canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsInventoryFilters(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]OSPolicyAssignmentOSPoliciesResourceGroupsInventoryFilters, 0, len(des))
	for i, d := range des {
		cd := canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsInventoryFilters(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsInventoryFilters(c *Client, des, nw *OSPolicyAssignmentOSPoliciesResourceGroupsInventoryFilters) *OSPolicyAssignmentOSPoliciesResourceGroupsInventoryFilters {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsNotReturnedByServer(des) {
			c.Config.Logger.Info("Found explicitly empty value for OSPolicyAssignmentOSPoliciesResourceGroupsInventoryFilters while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringCanonicalize(des.OSShortName, nw.OSShortName) {
		nw.OSShortName = des.OSShortName
	}
	if dcl.StringCanonicalize(des.OSVersion, nw.OSVersion) {
		nw.OSVersion = des.OSVersion
	}

	return nw
}

func canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsInventoryFiltersSet(c *Client, des, nw []OSPolicyAssignmentOSPoliciesResourceGroupsInventoryFilters) []OSPolicyAssignmentOSPoliciesResourceGroupsInventoryFilters {
	if des == nil {
		return nw
	}
	var reorderedNew []OSPolicyAssignmentOSPoliciesResourceGroupsInventoryFilters
	for _, d := range des {
		matchedNew := -1
		for idx, n := range nw {
			if diffs, _ := compareOSPolicyAssignmentOSPoliciesResourceGroupsInventoryFiltersNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
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

func canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsInventoryFiltersSlice(c *Client, des, nw []OSPolicyAssignmentOSPoliciesResourceGroupsInventoryFilters) []OSPolicyAssignmentOSPoliciesResourceGroupsInventoryFilters {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []OSPolicyAssignmentOSPoliciesResourceGroupsInventoryFilters
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsInventoryFilters(c, &d, &n))
	}

	return items
}

func canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResources(des, initial *OSPolicyAssignmentOSPoliciesResourceGroupsResources, opts ...dcl.ApplyOption) *OSPolicyAssignmentOSPoliciesResourceGroupsResources {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if des.Pkg != nil || (initial != nil && initial.Pkg != nil) {
		// Check if anything else is set.
		if dcl.AnySet(des.Repository, des.Exec, des.File) {
			des.Pkg = nil
			if initial != nil {
				initial.Pkg = nil
			}
		}
	}

	if des.Repository != nil || (initial != nil && initial.Repository != nil) {
		// Check if anything else is set.
		if dcl.AnySet(des.Pkg, des.Exec, des.File) {
			des.Repository = nil
			if initial != nil {
				initial.Repository = nil
			}
		}
	}

	if des.Exec != nil || (initial != nil && initial.Exec != nil) {
		// Check if anything else is set.
		if dcl.AnySet(des.Pkg, des.Repository, des.File) {
			des.Exec = nil
			if initial != nil {
				initial.Exec = nil
			}
		}
	}

	if des.File != nil || (initial != nil && initial.File != nil) {
		// Check if anything else is set.
		if dcl.AnySet(des.Pkg, des.Repository, des.Exec) {
			des.File = nil
			if initial != nil {
				initial.File = nil
			}
		}
	}

	if initial == nil {
		return des
	}

	cDes := &OSPolicyAssignmentOSPoliciesResourceGroupsResources{}

	if dcl.StringCanonicalize(des.Id, initial.Id) || dcl.IsZeroValue(des.Id) {
		cDes.Id = initial.Id
	} else {
		cDes.Id = des.Id
	}
	cDes.Pkg = canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkg(des.Pkg, initial.Pkg, opts...)
	cDes.Repository = canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepository(des.Repository, initial.Repository, opts...)
	cDes.Exec = canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExec(des.Exec, initial.Exec, opts...)
	cDes.File = canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFile(des.File, initial.File, opts...)

	return cDes
}

func canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesSlice(des, initial []OSPolicyAssignmentOSPoliciesResourceGroupsResources, opts ...dcl.ApplyOption) []OSPolicyAssignmentOSPoliciesResourceGroupsResources {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]OSPolicyAssignmentOSPoliciesResourceGroupsResources, 0, len(des))
		for _, d := range des {
			cd := canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResources(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]OSPolicyAssignmentOSPoliciesResourceGroupsResources, 0, len(des))
	for i, d := range des {
		cd := canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResources(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResources(c *Client, des, nw *OSPolicyAssignmentOSPoliciesResourceGroupsResources) *OSPolicyAssignmentOSPoliciesResourceGroupsResources {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsNotReturnedByServer(des) {
			c.Config.Logger.Info("Found explicitly empty value for OSPolicyAssignmentOSPoliciesResourceGroupsResources while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringCanonicalize(des.Id, nw.Id) {
		nw.Id = des.Id
	}
	nw.Pkg = canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkg(c, des.Pkg, nw.Pkg)
	nw.Repository = canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepository(c, des.Repository, nw.Repository)
	nw.Exec = canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExec(c, des.Exec, nw.Exec)
	nw.File = canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFile(c, des.File, nw.File)

	return nw
}

func canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesSet(c *Client, des, nw []OSPolicyAssignmentOSPoliciesResourceGroupsResources) []OSPolicyAssignmentOSPoliciesResourceGroupsResources {
	if des == nil {
		return nw
	}
	var reorderedNew []OSPolicyAssignmentOSPoliciesResourceGroupsResources
	for _, d := range des {
		matchedNew := -1
		for idx, n := range nw {
			if diffs, _ := compareOSPolicyAssignmentOSPoliciesResourceGroupsResourcesNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
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

func canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesSlice(c *Client, des, nw []OSPolicyAssignmentOSPoliciesResourceGroupsResources) []OSPolicyAssignmentOSPoliciesResourceGroupsResources {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []OSPolicyAssignmentOSPoliciesResourceGroupsResources
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResources(c, &d, &n))
	}

	return items
}

func canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkg(des, initial *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkg, opts ...dcl.ApplyOption) *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkg {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if des.Apt != nil || (initial != nil && initial.Apt != nil) {
		// Check if anything else is set.
		if dcl.AnySet(des.Deb, des.Yum, des.Zypper, des.Rpm, des.Googet, des.Msi) {
			des.Apt = nil
			if initial != nil {
				initial.Apt = nil
			}
		}
	}

	if des.Deb != nil || (initial != nil && initial.Deb != nil) {
		// Check if anything else is set.
		if dcl.AnySet(des.Apt, des.Yum, des.Zypper, des.Rpm, des.Googet, des.Msi) {
			des.Deb = nil
			if initial != nil {
				initial.Deb = nil
			}
		}
	}

	if des.Yum != nil || (initial != nil && initial.Yum != nil) {
		// Check if anything else is set.
		if dcl.AnySet(des.Apt, des.Deb, des.Zypper, des.Rpm, des.Googet, des.Msi) {
			des.Yum = nil
			if initial != nil {
				initial.Yum = nil
			}
		}
	}

	if des.Zypper != nil || (initial != nil && initial.Zypper != nil) {
		// Check if anything else is set.
		if dcl.AnySet(des.Apt, des.Deb, des.Yum, des.Rpm, des.Googet, des.Msi) {
			des.Zypper = nil
			if initial != nil {
				initial.Zypper = nil
			}
		}
	}

	if des.Rpm != nil || (initial != nil && initial.Rpm != nil) {
		// Check if anything else is set.
		if dcl.AnySet(des.Apt, des.Deb, des.Yum, des.Zypper, des.Googet, des.Msi) {
			des.Rpm = nil
			if initial != nil {
				initial.Rpm = nil
			}
		}
	}

	if des.Googet != nil || (initial != nil && initial.Googet != nil) {
		// Check if anything else is set.
		if dcl.AnySet(des.Apt, des.Deb, des.Yum, des.Zypper, des.Rpm, des.Msi) {
			des.Googet = nil
			if initial != nil {
				initial.Googet = nil
			}
		}
	}

	if des.Msi != nil || (initial != nil && initial.Msi != nil) {
		// Check if anything else is set.
		if dcl.AnySet(des.Apt, des.Deb, des.Yum, des.Zypper, des.Rpm, des.Googet) {
			des.Msi = nil
			if initial != nil {
				initial.Msi = nil
			}
		}
	}

	if initial == nil {
		return des
	}

	cDes := &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkg{}

	if dcl.IsZeroValue(des.DesiredState) {
		cDes.DesiredState = initial.DesiredState
	} else {
		cDes.DesiredState = des.DesiredState
	}
	cDes.Apt = canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgApt(des.Apt, initial.Apt, opts...)
	cDes.Deb = canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDeb(des.Deb, initial.Deb, opts...)
	cDes.Yum = canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgYum(des.Yum, initial.Yum, opts...)
	cDes.Zypper = canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgZypper(des.Zypper, initial.Zypper, opts...)
	cDes.Rpm = canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpm(des.Rpm, initial.Rpm, opts...)
	cDes.Googet = canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgGooget(des.Googet, initial.Googet, opts...)
	cDes.Msi = canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsi(des.Msi, initial.Msi, opts...)

	return cDes
}

func canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgSlice(des, initial []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkg, opts ...dcl.ApplyOption) []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkg {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkg, 0, len(des))
		for _, d := range des {
			cd := canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkg(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkg, 0, len(des))
	for i, d := range des {
		cd := canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkg(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkg(c *Client, des, nw *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkg) *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkg {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsNotReturnedByServer(des) {
			c.Config.Logger.Info("Found explicitly empty value for OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkg while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	nw.Apt = canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgApt(c, des.Apt, nw.Apt)
	nw.Deb = canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDeb(c, des.Deb, nw.Deb)
	nw.Yum = canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgYum(c, des.Yum, nw.Yum)
	nw.Zypper = canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgZypper(c, des.Zypper, nw.Zypper)
	nw.Rpm = canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpm(c, des.Rpm, nw.Rpm)
	nw.Googet = canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgGooget(c, des.Googet, nw.Googet)
	nw.Msi = canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsi(c, des.Msi, nw.Msi)

	return nw
}

func canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgSet(c *Client, des, nw []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkg) []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkg {
	if des == nil {
		return nw
	}
	var reorderedNew []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkg
	for _, d := range des {
		matchedNew := -1
		for idx, n := range nw {
			if diffs, _ := compareOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
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

func canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgSlice(c *Client, des, nw []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkg) []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkg {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkg
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkg(c, &d, &n))
	}

	return items
}

func canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgApt(des, initial *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgApt, opts ...dcl.ApplyOption) *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgApt {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgApt{}

	if dcl.StringCanonicalize(des.Name, initial.Name) || dcl.IsZeroValue(des.Name) {
		cDes.Name = initial.Name
	} else {
		cDes.Name = des.Name
	}

	return cDes
}

func canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgAptSlice(des, initial []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgApt, opts ...dcl.ApplyOption) []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgApt {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgApt, 0, len(des))
		for _, d := range des {
			cd := canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgApt(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgApt, 0, len(des))
	for i, d := range des {
		cd := canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgApt(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgApt(c *Client, des, nw *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgApt) *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgApt {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsNotReturnedByServer(des) {
			c.Config.Logger.Info("Found explicitly empty value for OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgApt while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringCanonicalize(des.Name, nw.Name) {
		nw.Name = des.Name
	}

	return nw
}

func canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgAptSet(c *Client, des, nw []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgApt) []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgApt {
	if des == nil {
		return nw
	}
	var reorderedNew []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgApt
	for _, d := range des {
		matchedNew := -1
		for idx, n := range nw {
			if diffs, _ := compareOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgAptNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
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

func canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgAptSlice(c *Client, des, nw []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgApt) []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgApt {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgApt
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgApt(c, &d, &n))
	}

	return items
}

func canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDeb(des, initial *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDeb, opts ...dcl.ApplyOption) *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDeb {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDeb{}

	cDes.Source = canonicalizeOSPolicyAssignmentFile(des.Source, initial.Source, opts...)
	if dcl.BoolCanonicalize(des.PullDeps, initial.PullDeps) || dcl.IsZeroValue(des.PullDeps) {
		cDes.PullDeps = initial.PullDeps
	} else {
		cDes.PullDeps = des.PullDeps
	}

	return cDes
}

func canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSlice(des, initial []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDeb, opts ...dcl.ApplyOption) []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDeb {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDeb, 0, len(des))
		for _, d := range des {
			cd := canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDeb(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDeb, 0, len(des))
	for i, d := range des {
		cd := canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDeb(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDeb(c *Client, des, nw *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDeb) *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDeb {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsNotReturnedByServer(des) {
			c.Config.Logger.Info("Found explicitly empty value for OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDeb while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	nw.Source = canonicalizeNewOSPolicyAssignmentFile(c, des.Source, nw.Source)
	if dcl.BoolCanonicalize(des.PullDeps, nw.PullDeps) {
		nw.PullDeps = des.PullDeps
	}

	return nw
}

func canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSet(c *Client, des, nw []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDeb) []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDeb {
	if des == nil {
		return nw
	}
	var reorderedNew []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDeb
	for _, d := range des {
		matchedNew := -1
		for idx, n := range nw {
			if diffs, _ := compareOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
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

func canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSlice(c *Client, des, nw []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDeb) []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDeb {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDeb
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDeb(c, &d, &n))
	}

	return items
}

func canonicalizeOSPolicyAssignmentFile(des, initial *OSPolicyAssignmentFile, opts ...dcl.ApplyOption) *OSPolicyAssignmentFile {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if des.Remote != nil || (initial != nil && initial.Remote != nil) {
		// Check if anything else is set.
		if dcl.AnySet(des.Gcs, des.LocalPath) {
			des.Remote = nil
			if initial != nil {
				initial.Remote = nil
			}
		}
	}

	if des.Gcs != nil || (initial != nil && initial.Gcs != nil) {
		// Check if anything else is set.
		if dcl.AnySet(des.Remote, des.LocalPath) {
			des.Gcs = nil
			if initial != nil {
				initial.Gcs = nil
			}
		}
	}

	if des.LocalPath != nil || (initial != nil && initial.LocalPath != nil) {
		// Check if anything else is set.
		if dcl.AnySet(des.Remote, des.Gcs) {
			des.LocalPath = nil
			if initial != nil {
				initial.LocalPath = nil
			}
		}
	}

	if initial == nil {
		return des
	}

	cDes := &OSPolicyAssignmentFile{}

	cDes.Remote = canonicalizeOSPolicyAssignmentFileRemote(des.Remote, initial.Remote, opts...)
	cDes.Gcs = canonicalizeOSPolicyAssignmentFileGcs(des.Gcs, initial.Gcs, opts...)
	if dcl.StringCanonicalize(des.LocalPath, initial.LocalPath) || dcl.IsZeroValue(des.LocalPath) {
		cDes.LocalPath = initial.LocalPath
	} else {
		cDes.LocalPath = des.LocalPath
	}
	if dcl.BoolCanonicalize(des.AllowInsecure, initial.AllowInsecure) || dcl.IsZeroValue(des.AllowInsecure) {
		cDes.AllowInsecure = initial.AllowInsecure
	} else {
		cDes.AllowInsecure = des.AllowInsecure
	}

	return cDes
}

func canonicalizeOSPolicyAssignmentFileSlice(des, initial []OSPolicyAssignmentFile, opts ...dcl.ApplyOption) []OSPolicyAssignmentFile {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]OSPolicyAssignmentFile, 0, len(des))
		for _, d := range des {
			cd := canonicalizeOSPolicyAssignmentFile(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]OSPolicyAssignmentFile, 0, len(des))
	for i, d := range des {
		cd := canonicalizeOSPolicyAssignmentFile(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewOSPolicyAssignmentFile(c *Client, des, nw *OSPolicyAssignmentFile) *OSPolicyAssignmentFile {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsNotReturnedByServer(des) {
			c.Config.Logger.Info("Found explicitly empty value for OSPolicyAssignmentFile while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	nw.Remote = canonicalizeNewOSPolicyAssignmentFileRemote(c, des.Remote, nw.Remote)
	nw.Gcs = canonicalizeNewOSPolicyAssignmentFileGcs(c, des.Gcs, nw.Gcs)
	if dcl.StringCanonicalize(des.LocalPath, nw.LocalPath) {
		nw.LocalPath = des.LocalPath
	}
	if dcl.BoolCanonicalize(des.AllowInsecure, nw.AllowInsecure) {
		nw.AllowInsecure = des.AllowInsecure
	}

	return nw
}

func canonicalizeNewOSPolicyAssignmentFileSet(c *Client, des, nw []OSPolicyAssignmentFile) []OSPolicyAssignmentFile {
	if des == nil {
		return nw
	}
	var reorderedNew []OSPolicyAssignmentFile
	for _, d := range des {
		matchedNew := -1
		for idx, n := range nw {
			if diffs, _ := compareOSPolicyAssignmentFileNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
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

func canonicalizeNewOSPolicyAssignmentFileSlice(c *Client, des, nw []OSPolicyAssignmentFile) []OSPolicyAssignmentFile {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []OSPolicyAssignmentFile
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewOSPolicyAssignmentFile(c, &d, &n))
	}

	return items
}

func canonicalizeOSPolicyAssignmentFileRemote(des, initial *OSPolicyAssignmentFileRemote, opts ...dcl.ApplyOption) *OSPolicyAssignmentFileRemote {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &OSPolicyAssignmentFileRemote{}

	if dcl.StringCanonicalize(des.Uri, initial.Uri) || dcl.IsZeroValue(des.Uri) {
		cDes.Uri = initial.Uri
	} else {
		cDes.Uri = des.Uri
	}
	if dcl.StringCanonicalize(des.Sha256Checksum, initial.Sha256Checksum) || dcl.IsZeroValue(des.Sha256Checksum) {
		cDes.Sha256Checksum = initial.Sha256Checksum
	} else {
		cDes.Sha256Checksum = des.Sha256Checksum
	}

	return cDes
}

func canonicalizeOSPolicyAssignmentFileRemoteSlice(des, initial []OSPolicyAssignmentFileRemote, opts ...dcl.ApplyOption) []OSPolicyAssignmentFileRemote {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]OSPolicyAssignmentFileRemote, 0, len(des))
		for _, d := range des {
			cd := canonicalizeOSPolicyAssignmentFileRemote(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]OSPolicyAssignmentFileRemote, 0, len(des))
	for i, d := range des {
		cd := canonicalizeOSPolicyAssignmentFileRemote(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewOSPolicyAssignmentFileRemote(c *Client, des, nw *OSPolicyAssignmentFileRemote) *OSPolicyAssignmentFileRemote {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsNotReturnedByServer(des) {
			c.Config.Logger.Info("Found explicitly empty value for OSPolicyAssignmentFileRemote while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringCanonicalize(des.Uri, nw.Uri) {
		nw.Uri = des.Uri
	}
	if dcl.StringCanonicalize(des.Sha256Checksum, nw.Sha256Checksum) {
		nw.Sha256Checksum = des.Sha256Checksum
	}

	return nw
}

func canonicalizeNewOSPolicyAssignmentFileRemoteSet(c *Client, des, nw []OSPolicyAssignmentFileRemote) []OSPolicyAssignmentFileRemote {
	if des == nil {
		return nw
	}
	var reorderedNew []OSPolicyAssignmentFileRemote
	for _, d := range des {
		matchedNew := -1
		for idx, n := range nw {
			if diffs, _ := compareOSPolicyAssignmentFileRemoteNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
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

func canonicalizeNewOSPolicyAssignmentFileRemoteSlice(c *Client, des, nw []OSPolicyAssignmentFileRemote) []OSPolicyAssignmentFileRemote {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []OSPolicyAssignmentFileRemote
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewOSPolicyAssignmentFileRemote(c, &d, &n))
	}

	return items
}

func canonicalizeOSPolicyAssignmentFileGcs(des, initial *OSPolicyAssignmentFileGcs, opts ...dcl.ApplyOption) *OSPolicyAssignmentFileGcs {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &OSPolicyAssignmentFileGcs{}

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
	if dcl.IsZeroValue(des.Generation) {
		cDes.Generation = initial.Generation
	} else {
		cDes.Generation = des.Generation
	}

	return cDes
}

func canonicalizeOSPolicyAssignmentFileGcsSlice(des, initial []OSPolicyAssignmentFileGcs, opts ...dcl.ApplyOption) []OSPolicyAssignmentFileGcs {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]OSPolicyAssignmentFileGcs, 0, len(des))
		for _, d := range des {
			cd := canonicalizeOSPolicyAssignmentFileGcs(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]OSPolicyAssignmentFileGcs, 0, len(des))
	for i, d := range des {
		cd := canonicalizeOSPolicyAssignmentFileGcs(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewOSPolicyAssignmentFileGcs(c *Client, des, nw *OSPolicyAssignmentFileGcs) *OSPolicyAssignmentFileGcs {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsNotReturnedByServer(des) {
			c.Config.Logger.Info("Found explicitly empty value for OSPolicyAssignmentFileGcs while comparing non-nil desired to nil actual.  Returning desired object.")
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

func canonicalizeNewOSPolicyAssignmentFileGcsSet(c *Client, des, nw []OSPolicyAssignmentFileGcs) []OSPolicyAssignmentFileGcs {
	if des == nil {
		return nw
	}
	var reorderedNew []OSPolicyAssignmentFileGcs
	for _, d := range des {
		matchedNew := -1
		for idx, n := range nw {
			if diffs, _ := compareOSPolicyAssignmentFileGcsNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
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

func canonicalizeNewOSPolicyAssignmentFileGcsSlice(c *Client, des, nw []OSPolicyAssignmentFileGcs) []OSPolicyAssignmentFileGcs {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []OSPolicyAssignmentFileGcs
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewOSPolicyAssignmentFileGcs(c, &d, &n))
	}

	return items
}

func canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgYum(des, initial *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgYum, opts ...dcl.ApplyOption) *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgYum {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgYum{}

	if dcl.StringCanonicalize(des.Name, initial.Name) || dcl.IsZeroValue(des.Name) {
		cDes.Name = initial.Name
	} else {
		cDes.Name = des.Name
	}

	return cDes
}

func canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgYumSlice(des, initial []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgYum, opts ...dcl.ApplyOption) []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgYum {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgYum, 0, len(des))
		for _, d := range des {
			cd := canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgYum(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgYum, 0, len(des))
	for i, d := range des {
		cd := canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgYum(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgYum(c *Client, des, nw *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgYum) *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgYum {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsNotReturnedByServer(des) {
			c.Config.Logger.Info("Found explicitly empty value for OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgYum while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringCanonicalize(des.Name, nw.Name) {
		nw.Name = des.Name
	}

	return nw
}

func canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgYumSet(c *Client, des, nw []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgYum) []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgYum {
	if des == nil {
		return nw
	}
	var reorderedNew []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgYum
	for _, d := range des {
		matchedNew := -1
		for idx, n := range nw {
			if diffs, _ := compareOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgYumNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
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

func canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgYumSlice(c *Client, des, nw []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgYum) []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgYum {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgYum
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgYum(c, &d, &n))
	}

	return items
}

func canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgZypper(des, initial *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgZypper, opts ...dcl.ApplyOption) *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgZypper {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgZypper{}

	if dcl.StringCanonicalize(des.Name, initial.Name) || dcl.IsZeroValue(des.Name) {
		cDes.Name = initial.Name
	} else {
		cDes.Name = des.Name
	}

	return cDes
}

func canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgZypperSlice(des, initial []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgZypper, opts ...dcl.ApplyOption) []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgZypper {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgZypper, 0, len(des))
		for _, d := range des {
			cd := canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgZypper(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgZypper, 0, len(des))
	for i, d := range des {
		cd := canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgZypper(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgZypper(c *Client, des, nw *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgZypper) *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgZypper {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsNotReturnedByServer(des) {
			c.Config.Logger.Info("Found explicitly empty value for OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgZypper while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringCanonicalize(des.Name, nw.Name) {
		nw.Name = des.Name
	}

	return nw
}

func canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgZypperSet(c *Client, des, nw []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgZypper) []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgZypper {
	if des == nil {
		return nw
	}
	var reorderedNew []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgZypper
	for _, d := range des {
		matchedNew := -1
		for idx, n := range nw {
			if diffs, _ := compareOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgZypperNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
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

func canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgZypperSlice(c *Client, des, nw []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgZypper) []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgZypper {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgZypper
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgZypper(c, &d, &n))
	}

	return items
}

func canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpm(des, initial *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpm, opts ...dcl.ApplyOption) *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpm {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpm{}

	cDes.Source = canonicalizeOSPolicyAssignmentFile(des.Source, initial.Source, opts...)
	if dcl.BoolCanonicalize(des.PullDeps, initial.PullDeps) || dcl.IsZeroValue(des.PullDeps) {
		cDes.PullDeps = initial.PullDeps
	} else {
		cDes.PullDeps = des.PullDeps
	}

	return cDes
}

func canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSlice(des, initial []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpm, opts ...dcl.ApplyOption) []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpm {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpm, 0, len(des))
		for _, d := range des {
			cd := canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpm(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpm, 0, len(des))
	for i, d := range des {
		cd := canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpm(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpm(c *Client, des, nw *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpm) *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpm {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsNotReturnedByServer(des) {
			c.Config.Logger.Info("Found explicitly empty value for OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpm while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	nw.Source = canonicalizeNewOSPolicyAssignmentFile(c, des.Source, nw.Source)
	if dcl.BoolCanonicalize(des.PullDeps, nw.PullDeps) {
		nw.PullDeps = des.PullDeps
	}

	return nw
}

func canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSet(c *Client, des, nw []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpm) []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpm {
	if des == nil {
		return nw
	}
	var reorderedNew []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpm
	for _, d := range des {
		matchedNew := -1
		for idx, n := range nw {
			if diffs, _ := compareOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
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

func canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSlice(c *Client, des, nw []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpm) []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpm {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpm
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpm(c, &d, &n))
	}

	return items
}

func canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgGooget(des, initial *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgGooget, opts ...dcl.ApplyOption) *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgGooget {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgGooget{}

	if dcl.StringCanonicalize(des.Name, initial.Name) || dcl.IsZeroValue(des.Name) {
		cDes.Name = initial.Name
	} else {
		cDes.Name = des.Name
	}

	return cDes
}

func canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgGoogetSlice(des, initial []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgGooget, opts ...dcl.ApplyOption) []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgGooget {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgGooget, 0, len(des))
		for _, d := range des {
			cd := canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgGooget(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgGooget, 0, len(des))
	for i, d := range des {
		cd := canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgGooget(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgGooget(c *Client, des, nw *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgGooget) *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgGooget {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsNotReturnedByServer(des) {
			c.Config.Logger.Info("Found explicitly empty value for OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgGooget while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringCanonicalize(des.Name, nw.Name) {
		nw.Name = des.Name
	}

	return nw
}

func canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgGoogetSet(c *Client, des, nw []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgGooget) []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgGooget {
	if des == nil {
		return nw
	}
	var reorderedNew []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgGooget
	for _, d := range des {
		matchedNew := -1
		for idx, n := range nw {
			if diffs, _ := compareOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgGoogetNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
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

func canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgGoogetSlice(c *Client, des, nw []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgGooget) []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgGooget {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgGooget
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgGooget(c, &d, &n))
	}

	return items
}

func canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsi(des, initial *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsi, opts ...dcl.ApplyOption) *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsi {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsi{}

	cDes.Source = canonicalizeOSPolicyAssignmentFile(des.Source, initial.Source, opts...)
	if dcl.StringArrayCanonicalize(des.Properties, initial.Properties) || dcl.IsZeroValue(des.Properties) {
		cDes.Properties = initial.Properties
	} else {
		cDes.Properties = des.Properties
	}

	return cDes
}

func canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSlice(des, initial []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsi, opts ...dcl.ApplyOption) []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsi {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsi, 0, len(des))
		for _, d := range des {
			cd := canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsi(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsi, 0, len(des))
	for i, d := range des {
		cd := canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsi(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsi(c *Client, des, nw *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsi) *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsi {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsNotReturnedByServer(des) {
			c.Config.Logger.Info("Found explicitly empty value for OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsi while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	nw.Source = canonicalizeNewOSPolicyAssignmentFile(c, des.Source, nw.Source)
	if dcl.StringArrayCanonicalize(des.Properties, nw.Properties) {
		nw.Properties = des.Properties
	}

	return nw
}

func canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSet(c *Client, des, nw []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsi) []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsi {
	if des == nil {
		return nw
	}
	var reorderedNew []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsi
	for _, d := range des {
		matchedNew := -1
		for idx, n := range nw {
			if diffs, _ := compareOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
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

func canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSlice(c *Client, des, nw []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsi) []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsi {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsi
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsi(c, &d, &n))
	}

	return items
}

func canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepository(des, initial *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepository, opts ...dcl.ApplyOption) *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepository {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if des.Apt != nil || (initial != nil && initial.Apt != nil) {
		// Check if anything else is set.
		if dcl.AnySet(des.Yum, des.Zypper, des.Goo) {
			des.Apt = nil
			if initial != nil {
				initial.Apt = nil
			}
		}
	}

	if des.Yum != nil || (initial != nil && initial.Yum != nil) {
		// Check if anything else is set.
		if dcl.AnySet(des.Apt, des.Zypper, des.Goo) {
			des.Yum = nil
			if initial != nil {
				initial.Yum = nil
			}
		}
	}

	if des.Zypper != nil || (initial != nil && initial.Zypper != nil) {
		// Check if anything else is set.
		if dcl.AnySet(des.Apt, des.Yum, des.Goo) {
			des.Zypper = nil
			if initial != nil {
				initial.Zypper = nil
			}
		}
	}

	if des.Goo != nil || (initial != nil && initial.Goo != nil) {
		// Check if anything else is set.
		if dcl.AnySet(des.Apt, des.Yum, des.Zypper) {
			des.Goo = nil
			if initial != nil {
				initial.Goo = nil
			}
		}
	}

	if initial == nil {
		return des
	}

	cDes := &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepository{}

	cDes.Apt = canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryApt(des.Apt, initial.Apt, opts...)
	cDes.Yum = canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryYum(des.Yum, initial.Yum, opts...)
	cDes.Zypper = canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryZypper(des.Zypper, initial.Zypper, opts...)
	cDes.Goo = canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryGoo(des.Goo, initial.Goo, opts...)

	return cDes
}

func canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositorySlice(des, initial []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepository, opts ...dcl.ApplyOption) []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepository {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepository, 0, len(des))
		for _, d := range des {
			cd := canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepository(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepository, 0, len(des))
	for i, d := range des {
		cd := canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepository(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepository(c *Client, des, nw *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepository) *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepository {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsNotReturnedByServer(des) {
			c.Config.Logger.Info("Found explicitly empty value for OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepository while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	nw.Apt = canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryApt(c, des.Apt, nw.Apt)
	nw.Yum = canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryYum(c, des.Yum, nw.Yum)
	nw.Zypper = canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryZypper(c, des.Zypper, nw.Zypper)
	nw.Goo = canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryGoo(c, des.Goo, nw.Goo)

	return nw
}

func canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositorySet(c *Client, des, nw []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepository) []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepository {
	if des == nil {
		return nw
	}
	var reorderedNew []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepository
	for _, d := range des {
		matchedNew := -1
		for idx, n := range nw {
			if diffs, _ := compareOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
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

func canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositorySlice(c *Client, des, nw []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepository) []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepository {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepository
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepository(c, &d, &n))
	}

	return items
}

func canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryApt(des, initial *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryApt, opts ...dcl.ApplyOption) *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryApt {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryApt{}

	if dcl.IsZeroValue(des.ArchiveType) {
		cDes.ArchiveType = initial.ArchiveType
	} else {
		cDes.ArchiveType = des.ArchiveType
	}
	if dcl.StringCanonicalize(des.Uri, initial.Uri) || dcl.IsZeroValue(des.Uri) {
		cDes.Uri = initial.Uri
	} else {
		cDes.Uri = des.Uri
	}
	if dcl.StringCanonicalize(des.Distribution, initial.Distribution) || dcl.IsZeroValue(des.Distribution) {
		cDes.Distribution = initial.Distribution
	} else {
		cDes.Distribution = des.Distribution
	}
	if dcl.StringArrayCanonicalize(des.Components, initial.Components) || dcl.IsZeroValue(des.Components) {
		cDes.Components = initial.Components
	} else {
		cDes.Components = des.Components
	}
	if dcl.StringCanonicalize(des.GpgKey, initial.GpgKey) || dcl.IsZeroValue(des.GpgKey) {
		cDes.GpgKey = initial.GpgKey
	} else {
		cDes.GpgKey = des.GpgKey
	}

	return cDes
}

func canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryAptSlice(des, initial []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryApt, opts ...dcl.ApplyOption) []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryApt {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryApt, 0, len(des))
		for _, d := range des {
			cd := canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryApt(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryApt, 0, len(des))
	for i, d := range des {
		cd := canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryApt(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryApt(c *Client, des, nw *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryApt) *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryApt {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsNotReturnedByServer(des) {
			c.Config.Logger.Info("Found explicitly empty value for OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryApt while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringCanonicalize(des.Uri, nw.Uri) {
		nw.Uri = des.Uri
	}
	if dcl.StringCanonicalize(des.Distribution, nw.Distribution) {
		nw.Distribution = des.Distribution
	}
	if dcl.StringArrayCanonicalize(des.Components, nw.Components) {
		nw.Components = des.Components
	}
	if dcl.StringCanonicalize(des.GpgKey, nw.GpgKey) {
		nw.GpgKey = des.GpgKey
	}

	return nw
}

func canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryAptSet(c *Client, des, nw []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryApt) []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryApt {
	if des == nil {
		return nw
	}
	var reorderedNew []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryApt
	for _, d := range des {
		matchedNew := -1
		for idx, n := range nw {
			if diffs, _ := compareOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryAptNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
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

func canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryAptSlice(c *Client, des, nw []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryApt) []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryApt {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryApt
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryApt(c, &d, &n))
	}

	return items
}

func canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryYum(des, initial *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryYum, opts ...dcl.ApplyOption) *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryYum {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryYum{}

	if dcl.StringCanonicalize(des.Id, initial.Id) || dcl.IsZeroValue(des.Id) {
		cDes.Id = initial.Id
	} else {
		cDes.Id = des.Id
	}
	if dcl.StringCanonicalize(des.DisplayName, initial.DisplayName) || dcl.IsZeroValue(des.DisplayName) {
		cDes.DisplayName = initial.DisplayName
	} else {
		cDes.DisplayName = des.DisplayName
	}
	if dcl.StringCanonicalize(des.BaseUrl, initial.BaseUrl) || dcl.IsZeroValue(des.BaseUrl) {
		cDes.BaseUrl = initial.BaseUrl
	} else {
		cDes.BaseUrl = des.BaseUrl
	}
	if dcl.StringArrayCanonicalize(des.GpgKeys, initial.GpgKeys) || dcl.IsZeroValue(des.GpgKeys) {
		cDes.GpgKeys = initial.GpgKeys
	} else {
		cDes.GpgKeys = des.GpgKeys
	}

	return cDes
}

func canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryYumSlice(des, initial []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryYum, opts ...dcl.ApplyOption) []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryYum {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryYum, 0, len(des))
		for _, d := range des {
			cd := canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryYum(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryYum, 0, len(des))
	for i, d := range des {
		cd := canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryYum(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryYum(c *Client, des, nw *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryYum) *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryYum {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsNotReturnedByServer(des) {
			c.Config.Logger.Info("Found explicitly empty value for OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryYum while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringCanonicalize(des.Id, nw.Id) {
		nw.Id = des.Id
	}
	if dcl.StringCanonicalize(des.DisplayName, nw.DisplayName) {
		nw.DisplayName = des.DisplayName
	}
	if dcl.StringCanonicalize(des.BaseUrl, nw.BaseUrl) {
		nw.BaseUrl = des.BaseUrl
	}
	if dcl.StringArrayCanonicalize(des.GpgKeys, nw.GpgKeys) {
		nw.GpgKeys = des.GpgKeys
	}

	return nw
}

func canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryYumSet(c *Client, des, nw []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryYum) []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryYum {
	if des == nil {
		return nw
	}
	var reorderedNew []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryYum
	for _, d := range des {
		matchedNew := -1
		for idx, n := range nw {
			if diffs, _ := compareOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryYumNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
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

func canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryYumSlice(c *Client, des, nw []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryYum) []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryYum {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryYum
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryYum(c, &d, &n))
	}

	return items
}

func canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryZypper(des, initial *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryZypper, opts ...dcl.ApplyOption) *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryZypper {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryZypper{}

	if dcl.StringCanonicalize(des.Id, initial.Id) || dcl.IsZeroValue(des.Id) {
		cDes.Id = initial.Id
	} else {
		cDes.Id = des.Id
	}
	if dcl.StringCanonicalize(des.DisplayName, initial.DisplayName) || dcl.IsZeroValue(des.DisplayName) {
		cDes.DisplayName = initial.DisplayName
	} else {
		cDes.DisplayName = des.DisplayName
	}
	if dcl.StringCanonicalize(des.BaseUrl, initial.BaseUrl) || dcl.IsZeroValue(des.BaseUrl) {
		cDes.BaseUrl = initial.BaseUrl
	} else {
		cDes.BaseUrl = des.BaseUrl
	}
	if dcl.StringArrayCanonicalize(des.GpgKeys, initial.GpgKeys) || dcl.IsZeroValue(des.GpgKeys) {
		cDes.GpgKeys = initial.GpgKeys
	} else {
		cDes.GpgKeys = des.GpgKeys
	}

	return cDes
}

func canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryZypperSlice(des, initial []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryZypper, opts ...dcl.ApplyOption) []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryZypper {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryZypper, 0, len(des))
		for _, d := range des {
			cd := canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryZypper(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryZypper, 0, len(des))
	for i, d := range des {
		cd := canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryZypper(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryZypper(c *Client, des, nw *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryZypper) *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryZypper {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsNotReturnedByServer(des) {
			c.Config.Logger.Info("Found explicitly empty value for OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryZypper while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringCanonicalize(des.Id, nw.Id) {
		nw.Id = des.Id
	}
	if dcl.StringCanonicalize(des.DisplayName, nw.DisplayName) {
		nw.DisplayName = des.DisplayName
	}
	if dcl.StringCanonicalize(des.BaseUrl, nw.BaseUrl) {
		nw.BaseUrl = des.BaseUrl
	}
	if dcl.StringArrayCanonicalize(des.GpgKeys, nw.GpgKeys) {
		nw.GpgKeys = des.GpgKeys
	}

	return nw
}

func canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryZypperSet(c *Client, des, nw []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryZypper) []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryZypper {
	if des == nil {
		return nw
	}
	var reorderedNew []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryZypper
	for _, d := range des {
		matchedNew := -1
		for idx, n := range nw {
			if diffs, _ := compareOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryZypperNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
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

func canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryZypperSlice(c *Client, des, nw []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryZypper) []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryZypper {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryZypper
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryZypper(c, &d, &n))
	}

	return items
}

func canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryGoo(des, initial *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryGoo, opts ...dcl.ApplyOption) *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryGoo {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryGoo{}

	if dcl.StringCanonicalize(des.Name, initial.Name) || dcl.IsZeroValue(des.Name) {
		cDes.Name = initial.Name
	} else {
		cDes.Name = des.Name
	}
	if dcl.StringCanonicalize(des.Url, initial.Url) || dcl.IsZeroValue(des.Url) {
		cDes.Url = initial.Url
	} else {
		cDes.Url = des.Url
	}

	return cDes
}

func canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryGooSlice(des, initial []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryGoo, opts ...dcl.ApplyOption) []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryGoo {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryGoo, 0, len(des))
		for _, d := range des {
			cd := canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryGoo(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryGoo, 0, len(des))
	for i, d := range des {
		cd := canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryGoo(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryGoo(c *Client, des, nw *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryGoo) *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryGoo {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsNotReturnedByServer(des) {
			c.Config.Logger.Info("Found explicitly empty value for OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryGoo while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringCanonicalize(des.Name, nw.Name) {
		nw.Name = des.Name
	}
	if dcl.StringCanonicalize(des.Url, nw.Url) {
		nw.Url = des.Url
	}

	return nw
}

func canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryGooSet(c *Client, des, nw []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryGoo) []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryGoo {
	if des == nil {
		return nw
	}
	var reorderedNew []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryGoo
	for _, d := range des {
		matchedNew := -1
		for idx, n := range nw {
			if diffs, _ := compareOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryGooNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
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

func canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryGooSlice(c *Client, des, nw []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryGoo) []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryGoo {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryGoo
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryGoo(c, &d, &n))
	}

	return items
}

func canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExec(des, initial *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExec, opts ...dcl.ApplyOption) *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExec {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExec{}

	cDes.Validate = canonicalizeOSPolicyAssignmentExec(des.Validate, initial.Validate, opts...)
	cDes.Enforce = canonicalizeOSPolicyAssignmentExec(des.Enforce, initial.Enforce, opts...)

	return cDes
}

func canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecSlice(des, initial []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExec, opts ...dcl.ApplyOption) []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExec {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExec, 0, len(des))
		for _, d := range des {
			cd := canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExec(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExec, 0, len(des))
	for i, d := range des {
		cd := canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExec(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExec(c *Client, des, nw *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExec) *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExec {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsNotReturnedByServer(des) {
			c.Config.Logger.Info("Found explicitly empty value for OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExec while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	nw.Validate = canonicalizeNewOSPolicyAssignmentExec(c, des.Validate, nw.Validate)
	nw.Enforce = canonicalizeNewOSPolicyAssignmentExec(c, des.Enforce, nw.Enforce)

	return nw
}

func canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecSet(c *Client, des, nw []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExec) []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExec {
	if des == nil {
		return nw
	}
	var reorderedNew []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExec
	for _, d := range des {
		matchedNew := -1
		for idx, n := range nw {
			if diffs, _ := compareOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
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

func canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecSlice(c *Client, des, nw []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExec) []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExec {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExec
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExec(c, &d, &n))
	}

	return items
}

func canonicalizeOSPolicyAssignmentExec(des, initial *OSPolicyAssignmentExec, opts ...dcl.ApplyOption) *OSPolicyAssignmentExec {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if des.File != nil || (initial != nil && initial.File != nil) {
		// Check if anything else is set.
		if dcl.AnySet(des.Script) {
			des.File = nil
			if initial != nil {
				initial.File = nil
			}
		}
	}

	if des.Script != nil || (initial != nil && initial.Script != nil) {
		// Check if anything else is set.
		if dcl.AnySet(des.File) {
			des.Script = nil
			if initial != nil {
				initial.Script = nil
			}
		}
	}

	if initial == nil {
		return des
	}

	cDes := &OSPolicyAssignmentExec{}

	cDes.File = canonicalizeOSPolicyAssignmentFile(des.File, initial.File, opts...)
	if dcl.StringCanonicalize(des.Script, initial.Script) || dcl.IsZeroValue(des.Script) {
		cDes.Script = initial.Script
	} else {
		cDes.Script = des.Script
	}
	if dcl.StringArrayCanonicalize(des.Args, initial.Args) || dcl.IsZeroValue(des.Args) {
		cDes.Args = initial.Args
	} else {
		cDes.Args = des.Args
	}
	if dcl.IsZeroValue(des.Interpreter) {
		cDes.Interpreter = initial.Interpreter
	} else {
		cDes.Interpreter = des.Interpreter
	}
	if dcl.StringCanonicalize(des.OutputFilePath, initial.OutputFilePath) || dcl.IsZeroValue(des.OutputFilePath) {
		cDes.OutputFilePath = initial.OutputFilePath
	} else {
		cDes.OutputFilePath = des.OutputFilePath
	}

	return cDes
}

func canonicalizeOSPolicyAssignmentExecSlice(des, initial []OSPolicyAssignmentExec, opts ...dcl.ApplyOption) []OSPolicyAssignmentExec {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]OSPolicyAssignmentExec, 0, len(des))
		for _, d := range des {
			cd := canonicalizeOSPolicyAssignmentExec(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]OSPolicyAssignmentExec, 0, len(des))
	for i, d := range des {
		cd := canonicalizeOSPolicyAssignmentExec(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewOSPolicyAssignmentExec(c *Client, des, nw *OSPolicyAssignmentExec) *OSPolicyAssignmentExec {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsNotReturnedByServer(des) {
			c.Config.Logger.Info("Found explicitly empty value for OSPolicyAssignmentExec while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	nw.File = canonicalizeNewOSPolicyAssignmentFile(c, des.File, nw.File)
	if dcl.StringCanonicalize(des.Script, nw.Script) {
		nw.Script = des.Script
	}
	if dcl.StringArrayCanonicalize(des.Args, nw.Args) {
		nw.Args = des.Args
	}
	if dcl.StringCanonicalize(des.OutputFilePath, nw.OutputFilePath) {
		nw.OutputFilePath = des.OutputFilePath
	}

	return nw
}

func canonicalizeNewOSPolicyAssignmentExecSet(c *Client, des, nw []OSPolicyAssignmentExec) []OSPolicyAssignmentExec {
	if des == nil {
		return nw
	}
	var reorderedNew []OSPolicyAssignmentExec
	for _, d := range des {
		matchedNew := -1
		for idx, n := range nw {
			if diffs, _ := compareOSPolicyAssignmentExecNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
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

func canonicalizeNewOSPolicyAssignmentExecSlice(c *Client, des, nw []OSPolicyAssignmentExec) []OSPolicyAssignmentExec {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []OSPolicyAssignmentExec
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewOSPolicyAssignmentExec(c, &d, &n))
	}

	return items
}

func canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFile(des, initial *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFile, opts ...dcl.ApplyOption) *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFile {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if des.File != nil || (initial != nil && initial.File != nil) {
		// Check if anything else is set.
		if dcl.AnySet(des.Content) {
			des.File = nil
			if initial != nil {
				initial.File = nil
			}
		}
	}

	if des.Content != nil || (initial != nil && initial.Content != nil) {
		// Check if anything else is set.
		if dcl.AnySet(des.File) {
			des.Content = nil
			if initial != nil {
				initial.Content = nil
			}
		}
	}

	if initial == nil {
		return des
	}

	cDes := &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFile{}

	cDes.File = canonicalizeOSPolicyAssignmentFile(des.File, initial.File, opts...)
	if dcl.StringCanonicalize(des.Content, initial.Content) || dcl.IsZeroValue(des.Content) {
		cDes.Content = initial.Content
	} else {
		cDes.Content = des.Content
	}
	if dcl.StringCanonicalize(des.Path, initial.Path) || dcl.IsZeroValue(des.Path) {
		cDes.Path = initial.Path
	} else {
		cDes.Path = des.Path
	}
	if dcl.IsZeroValue(des.State) {
		cDes.State = initial.State
	} else {
		cDes.State = des.State
	}

	return cDes
}

func canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileSlice(des, initial []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFile, opts ...dcl.ApplyOption) []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFile {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFile, 0, len(des))
		for _, d := range des {
			cd := canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFile(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFile, 0, len(des))
	for i, d := range des {
		cd := canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFile(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFile(c *Client, des, nw *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFile) *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFile {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsNotReturnedByServer(des) {
			c.Config.Logger.Info("Found explicitly empty value for OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFile while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	nw.File = canonicalizeNewOSPolicyAssignmentFile(c, des.File, nw.File)
	if dcl.StringCanonicalize(des.Content, nw.Content) {
		nw.Content = des.Content
	}
	if dcl.StringCanonicalize(des.Path, nw.Path) {
		nw.Path = des.Path
	}
	if dcl.StringCanonicalize(des.Permissions, nw.Permissions) {
		nw.Permissions = des.Permissions
	}

	return nw
}

func canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileSet(c *Client, des, nw []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFile) []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFile {
	if des == nil {
		return nw
	}
	var reorderedNew []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFile
	for _, d := range des {
		matchedNew := -1
		for idx, n := range nw {
			if diffs, _ := compareOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
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

func canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileSlice(c *Client, des, nw []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFile) []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFile {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFile
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFile(c, &d, &n))
	}

	return items
}

func canonicalizeOSPolicyAssignmentInstanceFilter(des, initial *OSPolicyAssignmentInstanceFilter, opts ...dcl.ApplyOption) *OSPolicyAssignmentInstanceFilter {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &OSPolicyAssignmentInstanceFilter{}

	if dcl.BoolCanonicalize(des.All, initial.All) || dcl.IsZeroValue(des.All) {
		cDes.All = initial.All
	} else {
		cDes.All = des.All
	}
	cDes.InclusionLabels = canonicalizeOSPolicyAssignmentInstanceFilterInclusionLabelsSlice(des.InclusionLabels, initial.InclusionLabels, opts...)
	cDes.ExclusionLabels = canonicalizeOSPolicyAssignmentInstanceFilterExclusionLabelsSlice(des.ExclusionLabels, initial.ExclusionLabels, opts...)
	cDes.Inventories = canonicalizeOSPolicyAssignmentInstanceFilterInventoriesSlice(des.Inventories, initial.Inventories, opts...)

	return cDes
}

func canonicalizeOSPolicyAssignmentInstanceFilterSlice(des, initial []OSPolicyAssignmentInstanceFilter, opts ...dcl.ApplyOption) []OSPolicyAssignmentInstanceFilter {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]OSPolicyAssignmentInstanceFilter, 0, len(des))
		for _, d := range des {
			cd := canonicalizeOSPolicyAssignmentInstanceFilter(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]OSPolicyAssignmentInstanceFilter, 0, len(des))
	for i, d := range des {
		cd := canonicalizeOSPolicyAssignmentInstanceFilter(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewOSPolicyAssignmentInstanceFilter(c *Client, des, nw *OSPolicyAssignmentInstanceFilter) *OSPolicyAssignmentInstanceFilter {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsNotReturnedByServer(des) {
			c.Config.Logger.Info("Found explicitly empty value for OSPolicyAssignmentInstanceFilter while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.BoolCanonicalize(des.All, nw.All) {
		nw.All = des.All
	}
	nw.InclusionLabels = canonicalizeNewOSPolicyAssignmentInstanceFilterInclusionLabelsSlice(c, des.InclusionLabels, nw.InclusionLabels)
	nw.ExclusionLabels = canonicalizeNewOSPolicyAssignmentInstanceFilterExclusionLabelsSlice(c, des.ExclusionLabels, nw.ExclusionLabels)
	nw.Inventories = canonicalizeNewOSPolicyAssignmentInstanceFilterInventoriesSlice(c, des.Inventories, nw.Inventories)

	return nw
}

func canonicalizeNewOSPolicyAssignmentInstanceFilterSet(c *Client, des, nw []OSPolicyAssignmentInstanceFilter) []OSPolicyAssignmentInstanceFilter {
	if des == nil {
		return nw
	}
	var reorderedNew []OSPolicyAssignmentInstanceFilter
	for _, d := range des {
		matchedNew := -1
		for idx, n := range nw {
			if diffs, _ := compareOSPolicyAssignmentInstanceFilterNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
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

func canonicalizeNewOSPolicyAssignmentInstanceFilterSlice(c *Client, des, nw []OSPolicyAssignmentInstanceFilter) []OSPolicyAssignmentInstanceFilter {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []OSPolicyAssignmentInstanceFilter
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewOSPolicyAssignmentInstanceFilter(c, &d, &n))
	}

	return items
}

func canonicalizeOSPolicyAssignmentInstanceFilterInclusionLabels(des, initial *OSPolicyAssignmentInstanceFilterInclusionLabels, opts ...dcl.ApplyOption) *OSPolicyAssignmentInstanceFilterInclusionLabels {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &OSPolicyAssignmentInstanceFilterInclusionLabels{}

	if dcl.IsZeroValue(des.Labels) {
		cDes.Labels = initial.Labels
	} else {
		cDes.Labels = des.Labels
	}

	return cDes
}

func canonicalizeOSPolicyAssignmentInstanceFilterInclusionLabelsSlice(des, initial []OSPolicyAssignmentInstanceFilterInclusionLabels, opts ...dcl.ApplyOption) []OSPolicyAssignmentInstanceFilterInclusionLabels {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]OSPolicyAssignmentInstanceFilterInclusionLabels, 0, len(des))
		for _, d := range des {
			cd := canonicalizeOSPolicyAssignmentInstanceFilterInclusionLabels(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]OSPolicyAssignmentInstanceFilterInclusionLabels, 0, len(des))
	for i, d := range des {
		cd := canonicalizeOSPolicyAssignmentInstanceFilterInclusionLabels(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewOSPolicyAssignmentInstanceFilterInclusionLabels(c *Client, des, nw *OSPolicyAssignmentInstanceFilterInclusionLabels) *OSPolicyAssignmentInstanceFilterInclusionLabels {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsNotReturnedByServer(des) {
			c.Config.Logger.Info("Found explicitly empty value for OSPolicyAssignmentInstanceFilterInclusionLabels while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	return nw
}

func canonicalizeNewOSPolicyAssignmentInstanceFilterInclusionLabelsSet(c *Client, des, nw []OSPolicyAssignmentInstanceFilterInclusionLabels) []OSPolicyAssignmentInstanceFilterInclusionLabels {
	if des == nil {
		return nw
	}
	var reorderedNew []OSPolicyAssignmentInstanceFilterInclusionLabels
	for _, d := range des {
		matchedNew := -1
		for idx, n := range nw {
			if diffs, _ := compareOSPolicyAssignmentInstanceFilterInclusionLabelsNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
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

func canonicalizeNewOSPolicyAssignmentInstanceFilterInclusionLabelsSlice(c *Client, des, nw []OSPolicyAssignmentInstanceFilterInclusionLabels) []OSPolicyAssignmentInstanceFilterInclusionLabels {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []OSPolicyAssignmentInstanceFilterInclusionLabels
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewOSPolicyAssignmentInstanceFilterInclusionLabels(c, &d, &n))
	}

	return items
}

func canonicalizeOSPolicyAssignmentInstanceFilterExclusionLabels(des, initial *OSPolicyAssignmentInstanceFilterExclusionLabels, opts ...dcl.ApplyOption) *OSPolicyAssignmentInstanceFilterExclusionLabels {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &OSPolicyAssignmentInstanceFilterExclusionLabels{}

	if dcl.IsZeroValue(des.Labels) {
		cDes.Labels = initial.Labels
	} else {
		cDes.Labels = des.Labels
	}

	return cDes
}

func canonicalizeOSPolicyAssignmentInstanceFilterExclusionLabelsSlice(des, initial []OSPolicyAssignmentInstanceFilterExclusionLabels, opts ...dcl.ApplyOption) []OSPolicyAssignmentInstanceFilterExclusionLabels {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]OSPolicyAssignmentInstanceFilterExclusionLabels, 0, len(des))
		for _, d := range des {
			cd := canonicalizeOSPolicyAssignmentInstanceFilterExclusionLabels(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]OSPolicyAssignmentInstanceFilterExclusionLabels, 0, len(des))
	for i, d := range des {
		cd := canonicalizeOSPolicyAssignmentInstanceFilterExclusionLabels(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewOSPolicyAssignmentInstanceFilterExclusionLabels(c *Client, des, nw *OSPolicyAssignmentInstanceFilterExclusionLabels) *OSPolicyAssignmentInstanceFilterExclusionLabels {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsNotReturnedByServer(des) {
			c.Config.Logger.Info("Found explicitly empty value for OSPolicyAssignmentInstanceFilterExclusionLabels while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	return nw
}

func canonicalizeNewOSPolicyAssignmentInstanceFilterExclusionLabelsSet(c *Client, des, nw []OSPolicyAssignmentInstanceFilterExclusionLabels) []OSPolicyAssignmentInstanceFilterExclusionLabels {
	if des == nil {
		return nw
	}
	var reorderedNew []OSPolicyAssignmentInstanceFilterExclusionLabels
	for _, d := range des {
		matchedNew := -1
		for idx, n := range nw {
			if diffs, _ := compareOSPolicyAssignmentInstanceFilterExclusionLabelsNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
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

func canonicalizeNewOSPolicyAssignmentInstanceFilterExclusionLabelsSlice(c *Client, des, nw []OSPolicyAssignmentInstanceFilterExclusionLabels) []OSPolicyAssignmentInstanceFilterExclusionLabels {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []OSPolicyAssignmentInstanceFilterExclusionLabels
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewOSPolicyAssignmentInstanceFilterExclusionLabels(c, &d, &n))
	}

	return items
}

func canonicalizeOSPolicyAssignmentInstanceFilterInventories(des, initial *OSPolicyAssignmentInstanceFilterInventories, opts ...dcl.ApplyOption) *OSPolicyAssignmentInstanceFilterInventories {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &OSPolicyAssignmentInstanceFilterInventories{}

	if dcl.StringCanonicalize(des.OSShortName, initial.OSShortName) || dcl.IsZeroValue(des.OSShortName) {
		cDes.OSShortName = initial.OSShortName
	} else {
		cDes.OSShortName = des.OSShortName
	}
	if dcl.StringCanonicalize(des.OSVersion, initial.OSVersion) || dcl.IsZeroValue(des.OSVersion) {
		cDes.OSVersion = initial.OSVersion
	} else {
		cDes.OSVersion = des.OSVersion
	}

	return cDes
}

func canonicalizeOSPolicyAssignmentInstanceFilterInventoriesSlice(des, initial []OSPolicyAssignmentInstanceFilterInventories, opts ...dcl.ApplyOption) []OSPolicyAssignmentInstanceFilterInventories {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]OSPolicyAssignmentInstanceFilterInventories, 0, len(des))
		for _, d := range des {
			cd := canonicalizeOSPolicyAssignmentInstanceFilterInventories(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]OSPolicyAssignmentInstanceFilterInventories, 0, len(des))
	for i, d := range des {
		cd := canonicalizeOSPolicyAssignmentInstanceFilterInventories(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewOSPolicyAssignmentInstanceFilterInventories(c *Client, des, nw *OSPolicyAssignmentInstanceFilterInventories) *OSPolicyAssignmentInstanceFilterInventories {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsNotReturnedByServer(des) {
			c.Config.Logger.Info("Found explicitly empty value for OSPolicyAssignmentInstanceFilterInventories while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringCanonicalize(des.OSShortName, nw.OSShortName) {
		nw.OSShortName = des.OSShortName
	}
	if dcl.StringCanonicalize(des.OSVersion, nw.OSVersion) {
		nw.OSVersion = des.OSVersion
	}

	return nw
}

func canonicalizeNewOSPolicyAssignmentInstanceFilterInventoriesSet(c *Client, des, nw []OSPolicyAssignmentInstanceFilterInventories) []OSPolicyAssignmentInstanceFilterInventories {
	if des == nil {
		return nw
	}
	var reorderedNew []OSPolicyAssignmentInstanceFilterInventories
	for _, d := range des {
		matchedNew := -1
		for idx, n := range nw {
			if diffs, _ := compareOSPolicyAssignmentInstanceFilterInventoriesNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
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

func canonicalizeNewOSPolicyAssignmentInstanceFilterInventoriesSlice(c *Client, des, nw []OSPolicyAssignmentInstanceFilterInventories) []OSPolicyAssignmentInstanceFilterInventories {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []OSPolicyAssignmentInstanceFilterInventories
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewOSPolicyAssignmentInstanceFilterInventories(c, &d, &n))
	}

	return items
}

func canonicalizeOSPolicyAssignmentRollout(des, initial *OSPolicyAssignmentRollout, opts ...dcl.ApplyOption) *OSPolicyAssignmentRollout {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &OSPolicyAssignmentRollout{}

	cDes.DisruptionBudget = canonicalizeOSPolicyAssignmentRolloutDisruptionBudget(des.DisruptionBudget, initial.DisruptionBudget, opts...)
	if canonicalizeOSPolicyAssignmentRolloutMinWaitDuration(des.MinWaitDuration, initial.MinWaitDuration) || dcl.IsZeroValue(des.MinWaitDuration) {
		cDes.MinWaitDuration = initial.MinWaitDuration
	} else {
		cDes.MinWaitDuration = des.MinWaitDuration
	}

	return cDes
}

func canonicalizeOSPolicyAssignmentRolloutSlice(des, initial []OSPolicyAssignmentRollout, opts ...dcl.ApplyOption) []OSPolicyAssignmentRollout {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]OSPolicyAssignmentRollout, 0, len(des))
		for _, d := range des {
			cd := canonicalizeOSPolicyAssignmentRollout(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]OSPolicyAssignmentRollout, 0, len(des))
	for i, d := range des {
		cd := canonicalizeOSPolicyAssignmentRollout(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewOSPolicyAssignmentRollout(c *Client, des, nw *OSPolicyAssignmentRollout) *OSPolicyAssignmentRollout {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsNotReturnedByServer(des) {
			c.Config.Logger.Info("Found explicitly empty value for OSPolicyAssignmentRollout while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	nw.DisruptionBudget = canonicalizeNewOSPolicyAssignmentRolloutDisruptionBudget(c, des.DisruptionBudget, nw.DisruptionBudget)
	if canonicalizeOSPolicyAssignmentRolloutMinWaitDuration(des.MinWaitDuration, nw.MinWaitDuration) {
		nw.MinWaitDuration = des.MinWaitDuration
	}

	return nw
}

func canonicalizeNewOSPolicyAssignmentRolloutSet(c *Client, des, nw []OSPolicyAssignmentRollout) []OSPolicyAssignmentRollout {
	if des == nil {
		return nw
	}
	var reorderedNew []OSPolicyAssignmentRollout
	for _, d := range des {
		matchedNew := -1
		for idx, n := range nw {
			if diffs, _ := compareOSPolicyAssignmentRolloutNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
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

func canonicalizeNewOSPolicyAssignmentRolloutSlice(c *Client, des, nw []OSPolicyAssignmentRollout) []OSPolicyAssignmentRollout {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []OSPolicyAssignmentRollout
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewOSPolicyAssignmentRollout(c, &d, &n))
	}

	return items
}

func canonicalizeOSPolicyAssignmentRolloutDisruptionBudget(des, initial *OSPolicyAssignmentRolloutDisruptionBudget, opts ...dcl.ApplyOption) *OSPolicyAssignmentRolloutDisruptionBudget {
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

	cDes := &OSPolicyAssignmentRolloutDisruptionBudget{}

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

func canonicalizeOSPolicyAssignmentRolloutDisruptionBudgetSlice(des, initial []OSPolicyAssignmentRolloutDisruptionBudget, opts ...dcl.ApplyOption) []OSPolicyAssignmentRolloutDisruptionBudget {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]OSPolicyAssignmentRolloutDisruptionBudget, 0, len(des))
		for _, d := range des {
			cd := canonicalizeOSPolicyAssignmentRolloutDisruptionBudget(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]OSPolicyAssignmentRolloutDisruptionBudget, 0, len(des))
	for i, d := range des {
		cd := canonicalizeOSPolicyAssignmentRolloutDisruptionBudget(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewOSPolicyAssignmentRolloutDisruptionBudget(c *Client, des, nw *OSPolicyAssignmentRolloutDisruptionBudget) *OSPolicyAssignmentRolloutDisruptionBudget {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsNotReturnedByServer(des) {
			c.Config.Logger.Info("Found explicitly empty value for OSPolicyAssignmentRolloutDisruptionBudget while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	return nw
}

func canonicalizeNewOSPolicyAssignmentRolloutDisruptionBudgetSet(c *Client, des, nw []OSPolicyAssignmentRolloutDisruptionBudget) []OSPolicyAssignmentRolloutDisruptionBudget {
	if des == nil {
		return nw
	}
	var reorderedNew []OSPolicyAssignmentRolloutDisruptionBudget
	for _, d := range des {
		matchedNew := -1
		for idx, n := range nw {
			if diffs, _ := compareOSPolicyAssignmentRolloutDisruptionBudgetNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
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

func canonicalizeNewOSPolicyAssignmentRolloutDisruptionBudgetSlice(c *Client, des, nw []OSPolicyAssignmentRolloutDisruptionBudget) []OSPolicyAssignmentRolloutDisruptionBudget {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []OSPolicyAssignmentRolloutDisruptionBudget
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewOSPolicyAssignmentRolloutDisruptionBudget(c, &d, &n))
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
func diffOSPolicyAssignment(c *Client, desired, actual *OSPolicyAssignment, opts ...dcl.ApplyOption) ([]*dcl.FieldDiff, error) {
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

	if ds, err := dcl.Diff(desired.Description, actual.Description, dcl.Info{OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("Description")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.OSPolicies, actual.OSPolicies, dcl.Info{ObjectFunction: compareOSPolicyAssignmentOSPoliciesNewStyle, EmptyObject: EmptyOSPolicyAssignmentOSPolicies, OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("OsPolicies")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.InstanceFilter, actual.InstanceFilter, dcl.Info{ObjectFunction: compareOSPolicyAssignmentInstanceFilterNewStyle, EmptyObject: EmptyOSPolicyAssignmentInstanceFilter, OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("InstanceFilter")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Rollout, actual.Rollout, dcl.Info{ObjectFunction: compareOSPolicyAssignmentRolloutNewStyle, EmptyObject: EmptyOSPolicyAssignmentRollout, OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("Rollout")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.RevisionId, actual.RevisionId, dcl.Info{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("RevisionId")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.RevisionCreateTime, actual.RevisionCreateTime, dcl.Info{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("RevisionCreateTime")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Etag, actual.Etag, dcl.Info{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Etag")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.RolloutState, actual.RolloutState, dcl.Info{OutputOnly: true, Type: "EnumType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("RolloutState")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Baseline, actual.Baseline, dcl.Info{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Baseline")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Deleted, actual.Deleted, dcl.Info{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Deleted")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Reconciling, actual.Reconciling, dcl.Info{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Reconciling")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Uid, actual.Uid, dcl.Info{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Uid")); len(ds) != 0 || err != nil {
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
func compareOSPolicyAssignmentOSPoliciesNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*OSPolicyAssignmentOSPolicies)
	if !ok {
		desiredNotPointer, ok := d.(OSPolicyAssignmentOSPolicies)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a OSPolicyAssignmentOSPolicies or *OSPolicyAssignmentOSPolicies", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*OSPolicyAssignmentOSPolicies)
	if !ok {
		actualNotPointer, ok := a.(OSPolicyAssignmentOSPolicies)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a OSPolicyAssignmentOSPolicies", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.Id, actual.Id, dcl.Info{OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("Id")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Description, actual.Description, dcl.Info{OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("Description")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Mode, actual.Mode, dcl.Info{Type: "EnumType", OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("Mode")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.ResourceGroups, actual.ResourceGroups, dcl.Info{ObjectFunction: compareOSPolicyAssignmentOSPoliciesResourceGroupsNewStyle, EmptyObject: EmptyOSPolicyAssignmentOSPoliciesResourceGroups, OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("ResourceGroups")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.AllowNoResourceGroupMatch, actual.AllowNoResourceGroupMatch, dcl.Info{OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("AllowNoResourceGroupMatch")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareOSPolicyAssignmentOSPoliciesResourceGroupsNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*OSPolicyAssignmentOSPoliciesResourceGroups)
	if !ok {
		desiredNotPointer, ok := d.(OSPolicyAssignmentOSPoliciesResourceGroups)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a OSPolicyAssignmentOSPoliciesResourceGroups or *OSPolicyAssignmentOSPoliciesResourceGroups", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*OSPolicyAssignmentOSPoliciesResourceGroups)
	if !ok {
		actualNotPointer, ok := a.(OSPolicyAssignmentOSPoliciesResourceGroups)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a OSPolicyAssignmentOSPoliciesResourceGroups", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.InventoryFilters, actual.InventoryFilters, dcl.Info{ObjectFunction: compareOSPolicyAssignmentOSPoliciesResourceGroupsInventoryFiltersNewStyle, EmptyObject: EmptyOSPolicyAssignmentOSPoliciesResourceGroupsInventoryFilters, OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("InventoryFilters")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Resources, actual.Resources, dcl.Info{ObjectFunction: compareOSPolicyAssignmentOSPoliciesResourceGroupsResourcesNewStyle, EmptyObject: EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResources, OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("Resources")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareOSPolicyAssignmentOSPoliciesResourceGroupsInventoryFiltersNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*OSPolicyAssignmentOSPoliciesResourceGroupsInventoryFilters)
	if !ok {
		desiredNotPointer, ok := d.(OSPolicyAssignmentOSPoliciesResourceGroupsInventoryFilters)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a OSPolicyAssignmentOSPoliciesResourceGroupsInventoryFilters or *OSPolicyAssignmentOSPoliciesResourceGroupsInventoryFilters", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*OSPolicyAssignmentOSPoliciesResourceGroupsInventoryFilters)
	if !ok {
		actualNotPointer, ok := a.(OSPolicyAssignmentOSPoliciesResourceGroupsInventoryFilters)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a OSPolicyAssignmentOSPoliciesResourceGroupsInventoryFilters", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.OSShortName, actual.OSShortName, dcl.Info{OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("OsShortName")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.OSVersion, actual.OSVersion, dcl.Info{OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("OsVersion")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareOSPolicyAssignmentOSPoliciesResourceGroupsResourcesNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*OSPolicyAssignmentOSPoliciesResourceGroupsResources)
	if !ok {
		desiredNotPointer, ok := d.(OSPolicyAssignmentOSPoliciesResourceGroupsResources)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a OSPolicyAssignmentOSPoliciesResourceGroupsResources or *OSPolicyAssignmentOSPoliciesResourceGroupsResources", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*OSPolicyAssignmentOSPoliciesResourceGroupsResources)
	if !ok {
		actualNotPointer, ok := a.(OSPolicyAssignmentOSPoliciesResourceGroupsResources)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a OSPolicyAssignmentOSPoliciesResourceGroupsResources", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.Id, actual.Id, dcl.Info{OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("Id")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Pkg, actual.Pkg, dcl.Info{ObjectFunction: compareOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgNewStyle, EmptyObject: EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkg, OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("Pkg")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Repository, actual.Repository, dcl.Info{ObjectFunction: compareOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryNewStyle, EmptyObject: EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepository, OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("Repository")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Exec, actual.Exec, dcl.Info{ObjectFunction: compareOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecNewStyle, EmptyObject: EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExec, OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("Exec")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.File, actual.File, dcl.Info{ObjectFunction: compareOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileNewStyle, EmptyObject: EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFile, OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("File")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkg)
	if !ok {
		desiredNotPointer, ok := d.(OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkg)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkg or *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkg", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkg)
	if !ok {
		actualNotPointer, ok := a.(OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkg)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkg", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.DesiredState, actual.DesiredState, dcl.Info{Type: "EnumType", OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("DesiredState")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Apt, actual.Apt, dcl.Info{ObjectFunction: compareOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgAptNewStyle, EmptyObject: EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgApt, OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("Apt")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Deb, actual.Deb, dcl.Info{ObjectFunction: compareOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebNewStyle, EmptyObject: EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDeb, OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("Deb")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Yum, actual.Yum, dcl.Info{ObjectFunction: compareOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgYumNewStyle, EmptyObject: EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgYum, OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("Yum")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Zypper, actual.Zypper, dcl.Info{ObjectFunction: compareOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgZypperNewStyle, EmptyObject: EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgZypper, OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("Zypper")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Rpm, actual.Rpm, dcl.Info{ObjectFunction: compareOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmNewStyle, EmptyObject: EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpm, OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("Rpm")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Googet, actual.Googet, dcl.Info{ObjectFunction: compareOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgGoogetNewStyle, EmptyObject: EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgGooget, OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("Googet")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Msi, actual.Msi, dcl.Info{ObjectFunction: compareOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiNewStyle, EmptyObject: EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsi, OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("Msi")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgAptNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgApt)
	if !ok {
		desiredNotPointer, ok := d.(OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgApt)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgApt or *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgApt", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgApt)
	if !ok {
		actualNotPointer, ok := a.(OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgApt)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgApt", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.Name, actual.Name, dcl.Info{OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("Name")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDeb)
	if !ok {
		desiredNotPointer, ok := d.(OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDeb)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDeb or *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDeb", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDeb)
	if !ok {
		actualNotPointer, ok := a.(OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDeb)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDeb", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.Source, actual.Source, dcl.Info{ObjectFunction: compareOSPolicyAssignmentFileNewStyle, EmptyObject: EmptyOSPolicyAssignmentFile, OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("Source")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.PullDeps, actual.PullDeps, dcl.Info{OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("PullDeps")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareOSPolicyAssignmentFileNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*OSPolicyAssignmentFile)
	if !ok {
		desiredNotPointer, ok := d.(OSPolicyAssignmentFile)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a OSPolicyAssignmentFile or *OSPolicyAssignmentFile", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*OSPolicyAssignmentFile)
	if !ok {
		actualNotPointer, ok := a.(OSPolicyAssignmentFile)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a OSPolicyAssignmentFile", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.Remote, actual.Remote, dcl.Info{ObjectFunction: compareOSPolicyAssignmentFileRemoteNewStyle, EmptyObject: EmptyOSPolicyAssignmentFileRemote, OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("Remote")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Gcs, actual.Gcs, dcl.Info{ObjectFunction: compareOSPolicyAssignmentFileGcsNewStyle, EmptyObject: EmptyOSPolicyAssignmentFileGcs, OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("Gcs")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.LocalPath, actual.LocalPath, dcl.Info{OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("LocalPath")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.AllowInsecure, actual.AllowInsecure, dcl.Info{OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("AllowInsecure")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareOSPolicyAssignmentFileRemoteNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*OSPolicyAssignmentFileRemote)
	if !ok {
		desiredNotPointer, ok := d.(OSPolicyAssignmentFileRemote)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a OSPolicyAssignmentFileRemote or *OSPolicyAssignmentFileRemote", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*OSPolicyAssignmentFileRemote)
	if !ok {
		actualNotPointer, ok := a.(OSPolicyAssignmentFileRemote)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a OSPolicyAssignmentFileRemote", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.Uri, actual.Uri, dcl.Info{OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("Uri")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Sha256Checksum, actual.Sha256Checksum, dcl.Info{OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("Sha256Checksum")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareOSPolicyAssignmentFileGcsNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*OSPolicyAssignmentFileGcs)
	if !ok {
		desiredNotPointer, ok := d.(OSPolicyAssignmentFileGcs)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a OSPolicyAssignmentFileGcs or *OSPolicyAssignmentFileGcs", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*OSPolicyAssignmentFileGcs)
	if !ok {
		actualNotPointer, ok := a.(OSPolicyAssignmentFileGcs)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a OSPolicyAssignmentFileGcs", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.Bucket, actual.Bucket, dcl.Info{OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("Bucket")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Object, actual.Object, dcl.Info{OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("Object")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Generation, actual.Generation, dcl.Info{OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("Generation")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgYumNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgYum)
	if !ok {
		desiredNotPointer, ok := d.(OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgYum)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgYum or *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgYum", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgYum)
	if !ok {
		actualNotPointer, ok := a.(OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgYum)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgYum", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.Name, actual.Name, dcl.Info{OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("Name")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgZypperNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgZypper)
	if !ok {
		desiredNotPointer, ok := d.(OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgZypper)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgZypper or *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgZypper", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgZypper)
	if !ok {
		actualNotPointer, ok := a.(OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgZypper)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgZypper", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.Name, actual.Name, dcl.Info{OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("Name")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpm)
	if !ok {
		desiredNotPointer, ok := d.(OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpm)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpm or *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpm", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpm)
	if !ok {
		actualNotPointer, ok := a.(OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpm)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpm", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.Source, actual.Source, dcl.Info{ObjectFunction: compareOSPolicyAssignmentFileNewStyle, EmptyObject: EmptyOSPolicyAssignmentFile, OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("Source")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.PullDeps, actual.PullDeps, dcl.Info{OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("PullDeps")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgGoogetNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgGooget)
	if !ok {
		desiredNotPointer, ok := d.(OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgGooget)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgGooget or *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgGooget", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgGooget)
	if !ok {
		actualNotPointer, ok := a.(OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgGooget)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgGooget", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.Name, actual.Name, dcl.Info{OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("Name")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsi)
	if !ok {
		desiredNotPointer, ok := d.(OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsi)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsi or *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsi", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsi)
	if !ok {
		actualNotPointer, ok := a.(OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsi)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsi", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.Source, actual.Source, dcl.Info{ObjectFunction: compareOSPolicyAssignmentFileNewStyle, EmptyObject: EmptyOSPolicyAssignmentFile, OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("Source")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Properties, actual.Properties, dcl.Info{OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("Properties")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepository)
	if !ok {
		desiredNotPointer, ok := d.(OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepository)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepository or *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepository", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepository)
	if !ok {
		actualNotPointer, ok := a.(OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepository)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepository", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.Apt, actual.Apt, dcl.Info{ObjectFunction: compareOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryAptNewStyle, EmptyObject: EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryApt, OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("Apt")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Yum, actual.Yum, dcl.Info{ObjectFunction: compareOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryYumNewStyle, EmptyObject: EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryYum, OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("Yum")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Zypper, actual.Zypper, dcl.Info{ObjectFunction: compareOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryZypperNewStyle, EmptyObject: EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryZypper, OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("Zypper")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Goo, actual.Goo, dcl.Info{ObjectFunction: compareOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryGooNewStyle, EmptyObject: EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryGoo, OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("Goo")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryAptNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryApt)
	if !ok {
		desiredNotPointer, ok := d.(OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryApt)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryApt or *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryApt", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryApt)
	if !ok {
		actualNotPointer, ok := a.(OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryApt)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryApt", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.ArchiveType, actual.ArchiveType, dcl.Info{Type: "EnumType", OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("ArchiveType")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Uri, actual.Uri, dcl.Info{OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("Uri")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Distribution, actual.Distribution, dcl.Info{OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("Distribution")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Components, actual.Components, dcl.Info{OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("Components")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.GpgKey, actual.GpgKey, dcl.Info{OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("GpgKey")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryYumNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryYum)
	if !ok {
		desiredNotPointer, ok := d.(OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryYum)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryYum or *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryYum", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryYum)
	if !ok {
		actualNotPointer, ok := a.(OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryYum)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryYum", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.Id, actual.Id, dcl.Info{OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("Id")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.DisplayName, actual.DisplayName, dcl.Info{OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("DisplayName")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.BaseUrl, actual.BaseUrl, dcl.Info{OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("BaseUrl")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.GpgKeys, actual.GpgKeys, dcl.Info{OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("GpgKeys")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryZypperNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryZypper)
	if !ok {
		desiredNotPointer, ok := d.(OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryZypper)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryZypper or *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryZypper", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryZypper)
	if !ok {
		actualNotPointer, ok := a.(OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryZypper)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryZypper", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.Id, actual.Id, dcl.Info{OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("Id")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.DisplayName, actual.DisplayName, dcl.Info{OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("DisplayName")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.BaseUrl, actual.BaseUrl, dcl.Info{OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("BaseUrl")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.GpgKeys, actual.GpgKeys, dcl.Info{OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("GpgKeys")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryGooNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryGoo)
	if !ok {
		desiredNotPointer, ok := d.(OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryGoo)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryGoo or *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryGoo", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryGoo)
	if !ok {
		actualNotPointer, ok := a.(OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryGoo)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryGoo", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.Name, actual.Name, dcl.Info{OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("Name")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Url, actual.Url, dcl.Info{OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("Url")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExec)
	if !ok {
		desiredNotPointer, ok := d.(OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExec)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExec or *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExec", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExec)
	if !ok {
		actualNotPointer, ok := a.(OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExec)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExec", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.Validate, actual.Validate, dcl.Info{ObjectFunction: compareOSPolicyAssignmentExecNewStyle, EmptyObject: EmptyOSPolicyAssignmentExec, OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("Validate")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Enforce, actual.Enforce, dcl.Info{ObjectFunction: compareOSPolicyAssignmentExecNewStyle, EmptyObject: EmptyOSPolicyAssignmentExec, OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("Enforce")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareOSPolicyAssignmentExecNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*OSPolicyAssignmentExec)
	if !ok {
		desiredNotPointer, ok := d.(OSPolicyAssignmentExec)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a OSPolicyAssignmentExec or *OSPolicyAssignmentExec", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*OSPolicyAssignmentExec)
	if !ok {
		actualNotPointer, ok := a.(OSPolicyAssignmentExec)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a OSPolicyAssignmentExec", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.File, actual.File, dcl.Info{ObjectFunction: compareOSPolicyAssignmentFileNewStyle, EmptyObject: EmptyOSPolicyAssignmentFile, OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("File")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Script, actual.Script, dcl.Info{OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("Script")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Args, actual.Args, dcl.Info{OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("Args")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Interpreter, actual.Interpreter, dcl.Info{Type: "EnumType", OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("Interpreter")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.OutputFilePath, actual.OutputFilePath, dcl.Info{OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("OutputFilePath")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFile)
	if !ok {
		desiredNotPointer, ok := d.(OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFile)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFile or *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFile", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFile)
	if !ok {
		actualNotPointer, ok := a.(OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFile)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFile", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.File, actual.File, dcl.Info{ObjectFunction: compareOSPolicyAssignmentFileNewStyle, EmptyObject: EmptyOSPolicyAssignmentFile, OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("File")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Content, actual.Content, dcl.Info{OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("Content")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Path, actual.Path, dcl.Info{OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("Path")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.State, actual.State, dcl.Info{Type: "EnumType", OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("State")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Permissions, actual.Permissions, dcl.Info{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Permissions")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareOSPolicyAssignmentInstanceFilterNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*OSPolicyAssignmentInstanceFilter)
	if !ok {
		desiredNotPointer, ok := d.(OSPolicyAssignmentInstanceFilter)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a OSPolicyAssignmentInstanceFilter or *OSPolicyAssignmentInstanceFilter", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*OSPolicyAssignmentInstanceFilter)
	if !ok {
		actualNotPointer, ok := a.(OSPolicyAssignmentInstanceFilter)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a OSPolicyAssignmentInstanceFilter", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.All, actual.All, dcl.Info{OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("All")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.InclusionLabels, actual.InclusionLabels, dcl.Info{ObjectFunction: compareOSPolicyAssignmentInstanceFilterInclusionLabelsNewStyle, EmptyObject: EmptyOSPolicyAssignmentInstanceFilterInclusionLabels, OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("InclusionLabels")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.ExclusionLabels, actual.ExclusionLabels, dcl.Info{ObjectFunction: compareOSPolicyAssignmentInstanceFilterExclusionLabelsNewStyle, EmptyObject: EmptyOSPolicyAssignmentInstanceFilterExclusionLabels, OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("ExclusionLabels")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Inventories, actual.Inventories, dcl.Info{ObjectFunction: compareOSPolicyAssignmentInstanceFilterInventoriesNewStyle, EmptyObject: EmptyOSPolicyAssignmentInstanceFilterInventories, OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("Inventories")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareOSPolicyAssignmentInstanceFilterInclusionLabelsNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*OSPolicyAssignmentInstanceFilterInclusionLabels)
	if !ok {
		desiredNotPointer, ok := d.(OSPolicyAssignmentInstanceFilterInclusionLabels)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a OSPolicyAssignmentInstanceFilterInclusionLabels or *OSPolicyAssignmentInstanceFilterInclusionLabels", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*OSPolicyAssignmentInstanceFilterInclusionLabels)
	if !ok {
		actualNotPointer, ok := a.(OSPolicyAssignmentInstanceFilterInclusionLabels)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a OSPolicyAssignmentInstanceFilterInclusionLabels", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.Labels, actual.Labels, dcl.Info{OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("Labels")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareOSPolicyAssignmentInstanceFilterExclusionLabelsNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*OSPolicyAssignmentInstanceFilterExclusionLabels)
	if !ok {
		desiredNotPointer, ok := d.(OSPolicyAssignmentInstanceFilterExclusionLabels)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a OSPolicyAssignmentInstanceFilterExclusionLabels or *OSPolicyAssignmentInstanceFilterExclusionLabels", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*OSPolicyAssignmentInstanceFilterExclusionLabels)
	if !ok {
		actualNotPointer, ok := a.(OSPolicyAssignmentInstanceFilterExclusionLabels)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a OSPolicyAssignmentInstanceFilterExclusionLabels", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.Labels, actual.Labels, dcl.Info{OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("Labels")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareOSPolicyAssignmentInstanceFilterInventoriesNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*OSPolicyAssignmentInstanceFilterInventories)
	if !ok {
		desiredNotPointer, ok := d.(OSPolicyAssignmentInstanceFilterInventories)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a OSPolicyAssignmentInstanceFilterInventories or *OSPolicyAssignmentInstanceFilterInventories", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*OSPolicyAssignmentInstanceFilterInventories)
	if !ok {
		actualNotPointer, ok := a.(OSPolicyAssignmentInstanceFilterInventories)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a OSPolicyAssignmentInstanceFilterInventories", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.OSShortName, actual.OSShortName, dcl.Info{OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("OsShortName")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.OSVersion, actual.OSVersion, dcl.Info{OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("OsVersion")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareOSPolicyAssignmentRolloutNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*OSPolicyAssignmentRollout)
	if !ok {
		desiredNotPointer, ok := d.(OSPolicyAssignmentRollout)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a OSPolicyAssignmentRollout or *OSPolicyAssignmentRollout", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*OSPolicyAssignmentRollout)
	if !ok {
		actualNotPointer, ok := a.(OSPolicyAssignmentRollout)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a OSPolicyAssignmentRollout", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.DisruptionBudget, actual.DisruptionBudget, dcl.Info{ObjectFunction: compareOSPolicyAssignmentRolloutDisruptionBudgetNewStyle, EmptyObject: EmptyOSPolicyAssignmentRolloutDisruptionBudget, OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("DisruptionBudget")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.MinWaitDuration, actual.MinWaitDuration, dcl.Info{CustomDiff: canonicalizeOSPolicyAssignmentRolloutMinWaitDuration, OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("MinWaitDuration")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareOSPolicyAssignmentRolloutDisruptionBudgetNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*OSPolicyAssignmentRolloutDisruptionBudget)
	if !ok {
		desiredNotPointer, ok := d.(OSPolicyAssignmentRolloutDisruptionBudget)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a OSPolicyAssignmentRolloutDisruptionBudget or *OSPolicyAssignmentRolloutDisruptionBudget", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*OSPolicyAssignmentRolloutDisruptionBudget)
	if !ok {
		actualNotPointer, ok := a.(OSPolicyAssignmentRolloutDisruptionBudget)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a OSPolicyAssignmentRolloutDisruptionBudget", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.Fixed, actual.Fixed, dcl.Info{OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("Fixed")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Percent, actual.Percent, dcl.Info{OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("Percent")); len(ds) != 0 || err != nil {
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
func (r *OSPolicyAssignment) urlNormalized() *OSPolicyAssignment {
	normalized := dcl.Copy(*r).(OSPolicyAssignment)
	normalized.Name = dcl.SelfLinkToName(r.Name)
	normalized.Description = dcl.SelfLinkToName(r.Description)
	normalized.RevisionId = dcl.SelfLinkToName(r.RevisionId)
	normalized.Etag = dcl.SelfLinkToName(r.Etag)
	normalized.Uid = dcl.SelfLinkToName(r.Uid)
	normalized.Project = dcl.SelfLinkToName(r.Project)
	normalized.Location = dcl.SelfLinkToName(r.Location)
	return &normalized
}

func (r *OSPolicyAssignment) updateURL(userBasePath, updateName string) (string, error) {
	nr := r.urlNormalized()
	if updateName == "UpdateOSPolicyAssignment" {
		fields := map[string]interface{}{
			"project":  dcl.ValueOrEmptyString(nr.Project),
			"location": dcl.ValueOrEmptyString(nr.Location),
			"name":     dcl.ValueOrEmptyString(nr.Name),
		}
		return dcl.URL("projects/{{project}}/locations/{{location}}/osPolicyAssignments/{{name}}", nr.basePath(), userBasePath, fields), nil

	}

	return "", fmt.Errorf("unknown update name: %s", updateName)
}

// marshal encodes the OSPolicyAssignment resource into JSON for a Create request, and
// performs transformations from the resource schema to the API schema if
// necessary.
func (r *OSPolicyAssignment) marshal(c *Client) ([]byte, error) {
	m, err := expandOSPolicyAssignment(c, r)
	if err != nil {
		return nil, fmt.Errorf("error marshalling OSPolicyAssignment: %w", err)
	}

	return json.Marshal(m)
}

// unmarshalOSPolicyAssignment decodes JSON responses into the OSPolicyAssignment resource schema.
func unmarshalOSPolicyAssignment(b []byte, c *Client) (*OSPolicyAssignment, error) {
	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, err
	}
	return unmarshalMapOSPolicyAssignment(m, c)
}

func unmarshalMapOSPolicyAssignment(m map[string]interface{}, c *Client) (*OSPolicyAssignment, error) {

	flattened := flattenOSPolicyAssignment(c, m)
	if flattened == nil {
		return nil, fmt.Errorf("attempted to flatten empty json object")
	}
	return flattened, nil
}

// expandOSPolicyAssignment expands OSPolicyAssignment into a JSON request object.
func expandOSPolicyAssignment(c *Client, f *OSPolicyAssignment) (map[string]interface{}, error) {
	m := make(map[string]interface{})
	if v, err := dcl.DeriveField("projects/%s/locations/%s/osPolicyAssignments/%s", f.Name, dcl.SelfLinkToName(f.Project), dcl.SelfLinkToName(f.Location), dcl.SelfLinkToName(f.Name)); err != nil {
		return nil, fmt.Errorf("error expanding Name into name: %w", err)
	} else if v != nil {
		m["name"] = v
	}
	if v := f.Description; dcl.ValueShouldBeSent(v) {
		m["description"] = v
	}
	if v, err := expandOSPolicyAssignmentOSPoliciesSlice(c, f.OSPolicies); err != nil {
		return nil, fmt.Errorf("error expanding OSPolicies into osPolicies: %w", err)
	} else {
		m["osPolicies"] = v
	}
	if v, err := expandOSPolicyAssignmentInstanceFilter(c, f.InstanceFilter); err != nil {
		return nil, fmt.Errorf("error expanding InstanceFilter into instanceFilter: %w", err)
	} else if v != nil {
		m["instanceFilter"] = v
	}
	if v, err := expandOSPolicyAssignmentRollout(c, f.Rollout); err != nil {
		return nil, fmt.Errorf("error expanding Rollout into rollout: %w", err)
	} else if v != nil {
		m["rollout"] = v
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

// flattenOSPolicyAssignment flattens OSPolicyAssignment from a JSON request object into the
// OSPolicyAssignment type.
func flattenOSPolicyAssignment(c *Client, i interface{}) *OSPolicyAssignment {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}
	if len(m) == 0 {
		return nil
	}

	res := &OSPolicyAssignment{}
	res.Name = dcl.FlattenString(m["name"])
	res.Description = dcl.FlattenString(m["description"])
	res.OSPolicies = flattenOSPolicyAssignmentOSPoliciesSlice(c, m["osPolicies"])
	res.InstanceFilter = flattenOSPolicyAssignmentInstanceFilter(c, m["instanceFilter"])
	res.Rollout = flattenOSPolicyAssignmentRollout(c, m["rollout"])
	res.RevisionId = dcl.FlattenString(m["revisionId"])
	res.RevisionCreateTime = dcl.FlattenString(m["revisionCreateTime"])
	res.Etag = dcl.FlattenString(m["etag"])
	res.RolloutState = flattenOSPolicyAssignmentRolloutStateEnum(m["rolloutState"])
	res.Baseline = dcl.FlattenBool(m["baseline"])
	res.Deleted = dcl.FlattenBool(m["deleted"])
	res.Reconciling = dcl.FlattenBool(m["reconciling"])
	res.Uid = dcl.FlattenString(m["uid"])
	res.Project = dcl.FlattenString(m["project"])
	res.Location = dcl.FlattenString(m["location"])

	return res
}

// expandOSPolicyAssignmentOSPoliciesMap expands the contents of OSPolicyAssignmentOSPolicies into a JSON
// request object.
func expandOSPolicyAssignmentOSPoliciesMap(c *Client, f map[string]OSPolicyAssignmentOSPolicies) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandOSPolicyAssignmentOSPolicies(c, &item)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandOSPolicyAssignmentOSPoliciesSlice expands the contents of OSPolicyAssignmentOSPolicies into a JSON
// request object.
func expandOSPolicyAssignmentOSPoliciesSlice(c *Client, f []OSPolicyAssignmentOSPolicies) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandOSPolicyAssignmentOSPolicies(c, &item)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenOSPolicyAssignmentOSPoliciesMap flattens the contents of OSPolicyAssignmentOSPolicies from a JSON
// response object.
func flattenOSPolicyAssignmentOSPoliciesMap(c *Client, i interface{}) map[string]OSPolicyAssignmentOSPolicies {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]OSPolicyAssignmentOSPolicies{}
	}

	if len(a) == 0 {
		return map[string]OSPolicyAssignmentOSPolicies{}
	}

	items := make(map[string]OSPolicyAssignmentOSPolicies)
	for k, item := range a {
		items[k] = *flattenOSPolicyAssignmentOSPolicies(c, item.(map[string]interface{}))
	}

	return items
}

// flattenOSPolicyAssignmentOSPoliciesSlice flattens the contents of OSPolicyAssignmentOSPolicies from a JSON
// response object.
func flattenOSPolicyAssignmentOSPoliciesSlice(c *Client, i interface{}) []OSPolicyAssignmentOSPolicies {
	a, ok := i.([]interface{})
	if !ok {
		return []OSPolicyAssignmentOSPolicies{}
	}

	if len(a) == 0 {
		return []OSPolicyAssignmentOSPolicies{}
	}

	items := make([]OSPolicyAssignmentOSPolicies, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenOSPolicyAssignmentOSPolicies(c, item.(map[string]interface{})))
	}

	return items
}

// expandOSPolicyAssignmentOSPolicies expands an instance of OSPolicyAssignmentOSPolicies into a JSON
// request object.
func expandOSPolicyAssignmentOSPolicies(c *Client, f *OSPolicyAssignmentOSPolicies) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.Id; !dcl.IsEmptyValueIndirect(v) {
		m["id"] = v
	}
	if v := f.Description; !dcl.IsEmptyValueIndirect(v) {
		m["description"] = v
	}
	if v := f.Mode; !dcl.IsEmptyValueIndirect(v) {
		m["mode"] = v
	}
	if v, err := expandOSPolicyAssignmentOSPoliciesResourceGroupsSlice(c, f.ResourceGroups); err != nil {
		return nil, fmt.Errorf("error expanding ResourceGroups into resourceGroups: %w", err)
	} else if v != nil {
		m["resourceGroups"] = v
	}
	if v := f.AllowNoResourceGroupMatch; !dcl.IsEmptyValueIndirect(v) {
		m["allowNoResourceGroupMatch"] = v
	}

	return m, nil
}

// flattenOSPolicyAssignmentOSPolicies flattens an instance of OSPolicyAssignmentOSPolicies from a JSON
// response object.
func flattenOSPolicyAssignmentOSPolicies(c *Client, i interface{}) *OSPolicyAssignmentOSPolicies {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &OSPolicyAssignmentOSPolicies{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyOSPolicyAssignmentOSPolicies
	}
	r.Id = dcl.FlattenString(m["id"])
	r.Description = dcl.FlattenString(m["description"])
	r.Mode = flattenOSPolicyAssignmentOSPoliciesModeEnum(m["mode"])
	r.ResourceGroups = flattenOSPolicyAssignmentOSPoliciesResourceGroupsSlice(c, m["resourceGroups"])
	r.AllowNoResourceGroupMatch = dcl.FlattenBool(m["allowNoResourceGroupMatch"])

	return r
}

// expandOSPolicyAssignmentOSPoliciesResourceGroupsMap expands the contents of OSPolicyAssignmentOSPoliciesResourceGroups into a JSON
// request object.
func expandOSPolicyAssignmentOSPoliciesResourceGroupsMap(c *Client, f map[string]OSPolicyAssignmentOSPoliciesResourceGroups) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandOSPolicyAssignmentOSPoliciesResourceGroups(c, &item)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandOSPolicyAssignmentOSPoliciesResourceGroupsSlice expands the contents of OSPolicyAssignmentOSPoliciesResourceGroups into a JSON
// request object.
func expandOSPolicyAssignmentOSPoliciesResourceGroupsSlice(c *Client, f []OSPolicyAssignmentOSPoliciesResourceGroups) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandOSPolicyAssignmentOSPoliciesResourceGroups(c, &item)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenOSPolicyAssignmentOSPoliciesResourceGroupsMap flattens the contents of OSPolicyAssignmentOSPoliciesResourceGroups from a JSON
// response object.
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsMap(c *Client, i interface{}) map[string]OSPolicyAssignmentOSPoliciesResourceGroups {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]OSPolicyAssignmentOSPoliciesResourceGroups{}
	}

	if len(a) == 0 {
		return map[string]OSPolicyAssignmentOSPoliciesResourceGroups{}
	}

	items := make(map[string]OSPolicyAssignmentOSPoliciesResourceGroups)
	for k, item := range a {
		items[k] = *flattenOSPolicyAssignmentOSPoliciesResourceGroups(c, item.(map[string]interface{}))
	}

	return items
}

// flattenOSPolicyAssignmentOSPoliciesResourceGroupsSlice flattens the contents of OSPolicyAssignmentOSPoliciesResourceGroups from a JSON
// response object.
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsSlice(c *Client, i interface{}) []OSPolicyAssignmentOSPoliciesResourceGroups {
	a, ok := i.([]interface{})
	if !ok {
		return []OSPolicyAssignmentOSPoliciesResourceGroups{}
	}

	if len(a) == 0 {
		return []OSPolicyAssignmentOSPoliciesResourceGroups{}
	}

	items := make([]OSPolicyAssignmentOSPoliciesResourceGroups, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenOSPolicyAssignmentOSPoliciesResourceGroups(c, item.(map[string]interface{})))
	}

	return items
}

// expandOSPolicyAssignmentOSPoliciesResourceGroups expands an instance of OSPolicyAssignmentOSPoliciesResourceGroups into a JSON
// request object.
func expandOSPolicyAssignmentOSPoliciesResourceGroups(c *Client, f *OSPolicyAssignmentOSPoliciesResourceGroups) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v, err := expandOSPolicyAssignmentOSPoliciesResourceGroupsInventoryFiltersSlice(c, f.InventoryFilters); err != nil {
		return nil, fmt.Errorf("error expanding InventoryFilters into inventoryFilters: %w", err)
	} else if v != nil {
		m["inventoryFilters"] = v
	}
	if v, err := expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesSlice(c, f.Resources); err != nil {
		return nil, fmt.Errorf("error expanding Resources into resources: %w", err)
	} else if v != nil {
		m["resources"] = v
	}

	return m, nil
}

// flattenOSPolicyAssignmentOSPoliciesResourceGroups flattens an instance of OSPolicyAssignmentOSPoliciesResourceGroups from a JSON
// response object.
func flattenOSPolicyAssignmentOSPoliciesResourceGroups(c *Client, i interface{}) *OSPolicyAssignmentOSPoliciesResourceGroups {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &OSPolicyAssignmentOSPoliciesResourceGroups{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyOSPolicyAssignmentOSPoliciesResourceGroups
	}
	r.InventoryFilters = flattenOSPolicyAssignmentOSPoliciesResourceGroupsInventoryFiltersSlice(c, m["inventoryFilters"])
	r.Resources = flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesSlice(c, m["resources"])

	return r
}

// expandOSPolicyAssignmentOSPoliciesResourceGroupsInventoryFiltersMap expands the contents of OSPolicyAssignmentOSPoliciesResourceGroupsInventoryFilters into a JSON
// request object.
func expandOSPolicyAssignmentOSPoliciesResourceGroupsInventoryFiltersMap(c *Client, f map[string]OSPolicyAssignmentOSPoliciesResourceGroupsInventoryFilters) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandOSPolicyAssignmentOSPoliciesResourceGroupsInventoryFilters(c, &item)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandOSPolicyAssignmentOSPoliciesResourceGroupsInventoryFiltersSlice expands the contents of OSPolicyAssignmentOSPoliciesResourceGroupsInventoryFilters into a JSON
// request object.
func expandOSPolicyAssignmentOSPoliciesResourceGroupsInventoryFiltersSlice(c *Client, f []OSPolicyAssignmentOSPoliciesResourceGroupsInventoryFilters) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandOSPolicyAssignmentOSPoliciesResourceGroupsInventoryFilters(c, &item)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenOSPolicyAssignmentOSPoliciesResourceGroupsInventoryFiltersMap flattens the contents of OSPolicyAssignmentOSPoliciesResourceGroupsInventoryFilters from a JSON
// response object.
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsInventoryFiltersMap(c *Client, i interface{}) map[string]OSPolicyAssignmentOSPoliciesResourceGroupsInventoryFilters {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]OSPolicyAssignmentOSPoliciesResourceGroupsInventoryFilters{}
	}

	if len(a) == 0 {
		return map[string]OSPolicyAssignmentOSPoliciesResourceGroupsInventoryFilters{}
	}

	items := make(map[string]OSPolicyAssignmentOSPoliciesResourceGroupsInventoryFilters)
	for k, item := range a {
		items[k] = *flattenOSPolicyAssignmentOSPoliciesResourceGroupsInventoryFilters(c, item.(map[string]interface{}))
	}

	return items
}

// flattenOSPolicyAssignmentOSPoliciesResourceGroupsInventoryFiltersSlice flattens the contents of OSPolicyAssignmentOSPoliciesResourceGroupsInventoryFilters from a JSON
// response object.
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsInventoryFiltersSlice(c *Client, i interface{}) []OSPolicyAssignmentOSPoliciesResourceGroupsInventoryFilters {
	a, ok := i.([]interface{})
	if !ok {
		return []OSPolicyAssignmentOSPoliciesResourceGroupsInventoryFilters{}
	}

	if len(a) == 0 {
		return []OSPolicyAssignmentOSPoliciesResourceGroupsInventoryFilters{}
	}

	items := make([]OSPolicyAssignmentOSPoliciesResourceGroupsInventoryFilters, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenOSPolicyAssignmentOSPoliciesResourceGroupsInventoryFilters(c, item.(map[string]interface{})))
	}

	return items
}

// expandOSPolicyAssignmentOSPoliciesResourceGroupsInventoryFilters expands an instance of OSPolicyAssignmentOSPoliciesResourceGroupsInventoryFilters into a JSON
// request object.
func expandOSPolicyAssignmentOSPoliciesResourceGroupsInventoryFilters(c *Client, f *OSPolicyAssignmentOSPoliciesResourceGroupsInventoryFilters) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.OSShortName; !dcl.IsEmptyValueIndirect(v) {
		m["osShortName"] = v
	}
	if v := f.OSVersion; !dcl.IsEmptyValueIndirect(v) {
		m["osVersion"] = v
	}

	return m, nil
}

// flattenOSPolicyAssignmentOSPoliciesResourceGroupsInventoryFilters flattens an instance of OSPolicyAssignmentOSPoliciesResourceGroupsInventoryFilters from a JSON
// response object.
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsInventoryFilters(c *Client, i interface{}) *OSPolicyAssignmentOSPoliciesResourceGroupsInventoryFilters {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &OSPolicyAssignmentOSPoliciesResourceGroupsInventoryFilters{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyOSPolicyAssignmentOSPoliciesResourceGroupsInventoryFilters
	}
	r.OSShortName = dcl.FlattenString(m["osShortName"])
	r.OSVersion = dcl.FlattenString(m["osVersion"])

	return r
}

// expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesMap expands the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResources into a JSON
// request object.
func expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesMap(c *Client, f map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResources) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandOSPolicyAssignmentOSPoliciesResourceGroupsResources(c, &item)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesSlice expands the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResources into a JSON
// request object.
func expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesSlice(c *Client, f []OSPolicyAssignmentOSPoliciesResourceGroupsResources) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandOSPolicyAssignmentOSPoliciesResourceGroupsResources(c, &item)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesMap flattens the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResources from a JSON
// response object.
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesMap(c *Client, i interface{}) map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResources {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResources{}
	}

	if len(a) == 0 {
		return map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResources{}
	}

	items := make(map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResources)
	for k, item := range a {
		items[k] = *flattenOSPolicyAssignmentOSPoliciesResourceGroupsResources(c, item.(map[string]interface{}))
	}

	return items
}

// flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesSlice flattens the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResources from a JSON
// response object.
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesSlice(c *Client, i interface{}) []OSPolicyAssignmentOSPoliciesResourceGroupsResources {
	a, ok := i.([]interface{})
	if !ok {
		return []OSPolicyAssignmentOSPoliciesResourceGroupsResources{}
	}

	if len(a) == 0 {
		return []OSPolicyAssignmentOSPoliciesResourceGroupsResources{}
	}

	items := make([]OSPolicyAssignmentOSPoliciesResourceGroupsResources, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenOSPolicyAssignmentOSPoliciesResourceGroupsResources(c, item.(map[string]interface{})))
	}

	return items
}

// expandOSPolicyAssignmentOSPoliciesResourceGroupsResources expands an instance of OSPolicyAssignmentOSPoliciesResourceGroupsResources into a JSON
// request object.
func expandOSPolicyAssignmentOSPoliciesResourceGroupsResources(c *Client, f *OSPolicyAssignmentOSPoliciesResourceGroupsResources) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.Id; !dcl.IsEmptyValueIndirect(v) {
		m["id"] = v
	}
	if v, err := expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkg(c, f.Pkg); err != nil {
		return nil, fmt.Errorf("error expanding Pkg into pkg: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["pkg"] = v
	}
	if v, err := expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepository(c, f.Repository); err != nil {
		return nil, fmt.Errorf("error expanding Repository into repository: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["repository"] = v
	}
	if v, err := expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExec(c, f.Exec); err != nil {
		return nil, fmt.Errorf("error expanding Exec into exec: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["exec"] = v
	}
	if v, err := expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFile(c, f.File); err != nil {
		return nil, fmt.Errorf("error expanding File into file: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["file"] = v
	}

	return m, nil
}

// flattenOSPolicyAssignmentOSPoliciesResourceGroupsResources flattens an instance of OSPolicyAssignmentOSPoliciesResourceGroupsResources from a JSON
// response object.
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsResources(c *Client, i interface{}) *OSPolicyAssignmentOSPoliciesResourceGroupsResources {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &OSPolicyAssignmentOSPoliciesResourceGroupsResources{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResources
	}
	r.Id = dcl.FlattenString(m["id"])
	r.Pkg = flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkg(c, m["pkg"])
	r.Repository = flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepository(c, m["repository"])
	r.Exec = flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExec(c, m["exec"])
	r.File = flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFile(c, m["file"])

	return r
}

// expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMap expands the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkg into a JSON
// request object.
func expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMap(c *Client, f map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkg) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkg(c, &item)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgSlice expands the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkg into a JSON
// request object.
func expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgSlice(c *Client, f []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkg) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkg(c, &item)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMap flattens the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkg from a JSON
// response object.
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMap(c *Client, i interface{}) map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkg {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkg{}
	}

	if len(a) == 0 {
		return map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkg{}
	}

	items := make(map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkg)
	for k, item := range a {
		items[k] = *flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkg(c, item.(map[string]interface{}))
	}

	return items
}

// flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgSlice flattens the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkg from a JSON
// response object.
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgSlice(c *Client, i interface{}) []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkg {
	a, ok := i.([]interface{})
	if !ok {
		return []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkg{}
	}

	if len(a) == 0 {
		return []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkg{}
	}

	items := make([]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkg, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkg(c, item.(map[string]interface{})))
	}

	return items
}

// expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkg expands an instance of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkg into a JSON
// request object.
func expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkg(c *Client, f *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkg) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.DesiredState; !dcl.IsEmptyValueIndirect(v) {
		m["desiredState"] = v
	}
	if v, err := expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgApt(c, f.Apt); err != nil {
		return nil, fmt.Errorf("error expanding Apt into apt: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["apt"] = v
	}
	if v, err := expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDeb(c, f.Deb); err != nil {
		return nil, fmt.Errorf("error expanding Deb into deb: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["deb"] = v
	}
	if v, err := expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgYum(c, f.Yum); err != nil {
		return nil, fmt.Errorf("error expanding Yum into yum: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["yum"] = v
	}
	if v, err := expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgZypper(c, f.Zypper); err != nil {
		return nil, fmt.Errorf("error expanding Zypper into zypper: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["zypper"] = v
	}
	if v, err := expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpm(c, f.Rpm); err != nil {
		return nil, fmt.Errorf("error expanding Rpm into rpm: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["rpm"] = v
	}
	if v, err := expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgGooget(c, f.Googet); err != nil {
		return nil, fmt.Errorf("error expanding Googet into googet: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["googet"] = v
	}
	if v, err := expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsi(c, f.Msi); err != nil {
		return nil, fmt.Errorf("error expanding Msi into msi: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["msi"] = v
	}

	return m, nil
}

// flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkg flattens an instance of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkg from a JSON
// response object.
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkg(c *Client, i interface{}) *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkg {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkg{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkg
	}
	r.DesiredState = flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDesiredStateEnum(m["desiredState"])
	r.Apt = flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgApt(c, m["apt"])
	r.Deb = flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDeb(c, m["deb"])
	r.Yum = flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgYum(c, m["yum"])
	r.Zypper = flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgZypper(c, m["zypper"])
	r.Rpm = flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpm(c, m["rpm"])
	r.Googet = flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgGooget(c, m["googet"])
	r.Msi = flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsi(c, m["msi"])

	return r
}

// expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgAptMap expands the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgApt into a JSON
// request object.
func expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgAptMap(c *Client, f map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgApt) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgApt(c, &item)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgAptSlice expands the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgApt into a JSON
// request object.
func expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgAptSlice(c *Client, f []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgApt) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgApt(c, &item)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgAptMap flattens the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgApt from a JSON
// response object.
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgAptMap(c *Client, i interface{}) map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgApt {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgApt{}
	}

	if len(a) == 0 {
		return map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgApt{}
	}

	items := make(map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgApt)
	for k, item := range a {
		items[k] = *flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgApt(c, item.(map[string]interface{}))
	}

	return items
}

// flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgAptSlice flattens the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgApt from a JSON
// response object.
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgAptSlice(c *Client, i interface{}) []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgApt {
	a, ok := i.([]interface{})
	if !ok {
		return []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgApt{}
	}

	if len(a) == 0 {
		return []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgApt{}
	}

	items := make([]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgApt, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgApt(c, item.(map[string]interface{})))
	}

	return items
}

// expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgApt expands an instance of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgApt into a JSON
// request object.
func expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgApt(c *Client, f *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgApt) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.Name; !dcl.IsEmptyValueIndirect(v) {
		m["name"] = v
	}

	return m, nil
}

// flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgApt flattens an instance of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgApt from a JSON
// response object.
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgApt(c *Client, i interface{}) *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgApt {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgApt{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgApt
	}
	r.Name = dcl.FlattenString(m["name"])

	return r
}

// expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebMap expands the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDeb into a JSON
// request object.
func expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebMap(c *Client, f map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDeb) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDeb(c, &item)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSlice expands the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDeb into a JSON
// request object.
func expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSlice(c *Client, f []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDeb) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDeb(c, &item)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebMap flattens the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDeb from a JSON
// response object.
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebMap(c *Client, i interface{}) map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDeb {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDeb{}
	}

	if len(a) == 0 {
		return map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDeb{}
	}

	items := make(map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDeb)
	for k, item := range a {
		items[k] = *flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDeb(c, item.(map[string]interface{}))
	}

	return items
}

// flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSlice flattens the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDeb from a JSON
// response object.
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSlice(c *Client, i interface{}) []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDeb {
	a, ok := i.([]interface{})
	if !ok {
		return []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDeb{}
	}

	if len(a) == 0 {
		return []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDeb{}
	}

	items := make([]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDeb, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDeb(c, item.(map[string]interface{})))
	}

	return items
}

// expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDeb expands an instance of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDeb into a JSON
// request object.
func expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDeb(c *Client, f *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDeb) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v, err := expandOSPolicyAssignmentFile(c, f.Source); err != nil {
		return nil, fmt.Errorf("error expanding Source into source: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["source"] = v
	}
	if v := f.PullDeps; !dcl.IsEmptyValueIndirect(v) {
		m["pullDeps"] = v
	}

	return m, nil
}

// flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDeb flattens an instance of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDeb from a JSON
// response object.
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDeb(c *Client, i interface{}) *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDeb {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDeb{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDeb
	}
	r.Source = flattenOSPolicyAssignmentFile(c, m["source"])
	r.PullDeps = dcl.FlattenBool(m["pullDeps"])

	return r
}

// expandOSPolicyAssignmentFileMap expands the contents of OSPolicyAssignmentFile into a JSON
// request object.
func expandOSPolicyAssignmentFileMap(c *Client, f map[string]OSPolicyAssignmentFile) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandOSPolicyAssignmentFile(c, &item)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandOSPolicyAssignmentFileSlice expands the contents of OSPolicyAssignmentFile into a JSON
// request object.
func expandOSPolicyAssignmentFileSlice(c *Client, f []OSPolicyAssignmentFile) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandOSPolicyAssignmentFile(c, &item)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenOSPolicyAssignmentFileMap flattens the contents of OSPolicyAssignmentFile from a JSON
// response object.
func flattenOSPolicyAssignmentFileMap(c *Client, i interface{}) map[string]OSPolicyAssignmentFile {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]OSPolicyAssignmentFile{}
	}

	if len(a) == 0 {
		return map[string]OSPolicyAssignmentFile{}
	}

	items := make(map[string]OSPolicyAssignmentFile)
	for k, item := range a {
		items[k] = *flattenOSPolicyAssignmentFile(c, item.(map[string]interface{}))
	}

	return items
}

// flattenOSPolicyAssignmentFileSlice flattens the contents of OSPolicyAssignmentFile from a JSON
// response object.
func flattenOSPolicyAssignmentFileSlice(c *Client, i interface{}) []OSPolicyAssignmentFile {
	a, ok := i.([]interface{})
	if !ok {
		return []OSPolicyAssignmentFile{}
	}

	if len(a) == 0 {
		return []OSPolicyAssignmentFile{}
	}

	items := make([]OSPolicyAssignmentFile, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenOSPolicyAssignmentFile(c, item.(map[string]interface{})))
	}

	return items
}

// expandOSPolicyAssignmentFile expands an instance of OSPolicyAssignmentFile into a JSON
// request object.
func expandOSPolicyAssignmentFile(c *Client, f *OSPolicyAssignmentFile) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v, err := expandOSPolicyAssignmentFileRemote(c, f.Remote); err != nil {
		return nil, fmt.Errorf("error expanding Remote into remote: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["remote"] = v
	}
	if v, err := expandOSPolicyAssignmentFileGcs(c, f.Gcs); err != nil {
		return nil, fmt.Errorf("error expanding Gcs into gcs: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["gcs"] = v
	}
	if v := f.LocalPath; !dcl.IsEmptyValueIndirect(v) {
		m["localPath"] = v
	}
	if v := f.AllowInsecure; !dcl.IsEmptyValueIndirect(v) {
		m["allowInsecure"] = v
	}

	return m, nil
}

// flattenOSPolicyAssignmentFile flattens an instance of OSPolicyAssignmentFile from a JSON
// response object.
func flattenOSPolicyAssignmentFile(c *Client, i interface{}) *OSPolicyAssignmentFile {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &OSPolicyAssignmentFile{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyOSPolicyAssignmentFile
	}
	r.Remote = flattenOSPolicyAssignmentFileRemote(c, m["remote"])
	r.Gcs = flattenOSPolicyAssignmentFileGcs(c, m["gcs"])
	r.LocalPath = dcl.FlattenString(m["localPath"])
	r.AllowInsecure = dcl.FlattenBool(m["allowInsecure"])

	return r
}

// expandOSPolicyAssignmentFileRemoteMap expands the contents of OSPolicyAssignmentFileRemote into a JSON
// request object.
func expandOSPolicyAssignmentFileRemoteMap(c *Client, f map[string]OSPolicyAssignmentFileRemote) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandOSPolicyAssignmentFileRemote(c, &item)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandOSPolicyAssignmentFileRemoteSlice expands the contents of OSPolicyAssignmentFileRemote into a JSON
// request object.
func expandOSPolicyAssignmentFileRemoteSlice(c *Client, f []OSPolicyAssignmentFileRemote) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandOSPolicyAssignmentFileRemote(c, &item)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenOSPolicyAssignmentFileRemoteMap flattens the contents of OSPolicyAssignmentFileRemote from a JSON
// response object.
func flattenOSPolicyAssignmentFileRemoteMap(c *Client, i interface{}) map[string]OSPolicyAssignmentFileRemote {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]OSPolicyAssignmentFileRemote{}
	}

	if len(a) == 0 {
		return map[string]OSPolicyAssignmentFileRemote{}
	}

	items := make(map[string]OSPolicyAssignmentFileRemote)
	for k, item := range a {
		items[k] = *flattenOSPolicyAssignmentFileRemote(c, item.(map[string]interface{}))
	}

	return items
}

// flattenOSPolicyAssignmentFileRemoteSlice flattens the contents of OSPolicyAssignmentFileRemote from a JSON
// response object.
func flattenOSPolicyAssignmentFileRemoteSlice(c *Client, i interface{}) []OSPolicyAssignmentFileRemote {
	a, ok := i.([]interface{})
	if !ok {
		return []OSPolicyAssignmentFileRemote{}
	}

	if len(a) == 0 {
		return []OSPolicyAssignmentFileRemote{}
	}

	items := make([]OSPolicyAssignmentFileRemote, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenOSPolicyAssignmentFileRemote(c, item.(map[string]interface{})))
	}

	return items
}

// expandOSPolicyAssignmentFileRemote expands an instance of OSPolicyAssignmentFileRemote into a JSON
// request object.
func expandOSPolicyAssignmentFileRemote(c *Client, f *OSPolicyAssignmentFileRemote) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.Uri; !dcl.IsEmptyValueIndirect(v) {
		m["uri"] = v
	}
	if v := f.Sha256Checksum; !dcl.IsEmptyValueIndirect(v) {
		m["sha256Checksum"] = v
	}

	return m, nil
}

// flattenOSPolicyAssignmentFileRemote flattens an instance of OSPolicyAssignmentFileRemote from a JSON
// response object.
func flattenOSPolicyAssignmentFileRemote(c *Client, i interface{}) *OSPolicyAssignmentFileRemote {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &OSPolicyAssignmentFileRemote{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyOSPolicyAssignmentFileRemote
	}
	r.Uri = dcl.FlattenString(m["uri"])
	r.Sha256Checksum = dcl.FlattenString(m["sha256Checksum"])

	return r
}

// expandOSPolicyAssignmentFileGcsMap expands the contents of OSPolicyAssignmentFileGcs into a JSON
// request object.
func expandOSPolicyAssignmentFileGcsMap(c *Client, f map[string]OSPolicyAssignmentFileGcs) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandOSPolicyAssignmentFileGcs(c, &item)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandOSPolicyAssignmentFileGcsSlice expands the contents of OSPolicyAssignmentFileGcs into a JSON
// request object.
func expandOSPolicyAssignmentFileGcsSlice(c *Client, f []OSPolicyAssignmentFileGcs) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandOSPolicyAssignmentFileGcs(c, &item)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenOSPolicyAssignmentFileGcsMap flattens the contents of OSPolicyAssignmentFileGcs from a JSON
// response object.
func flattenOSPolicyAssignmentFileGcsMap(c *Client, i interface{}) map[string]OSPolicyAssignmentFileGcs {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]OSPolicyAssignmentFileGcs{}
	}

	if len(a) == 0 {
		return map[string]OSPolicyAssignmentFileGcs{}
	}

	items := make(map[string]OSPolicyAssignmentFileGcs)
	for k, item := range a {
		items[k] = *flattenOSPolicyAssignmentFileGcs(c, item.(map[string]interface{}))
	}

	return items
}

// flattenOSPolicyAssignmentFileGcsSlice flattens the contents of OSPolicyAssignmentFileGcs from a JSON
// response object.
func flattenOSPolicyAssignmentFileGcsSlice(c *Client, i interface{}) []OSPolicyAssignmentFileGcs {
	a, ok := i.([]interface{})
	if !ok {
		return []OSPolicyAssignmentFileGcs{}
	}

	if len(a) == 0 {
		return []OSPolicyAssignmentFileGcs{}
	}

	items := make([]OSPolicyAssignmentFileGcs, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenOSPolicyAssignmentFileGcs(c, item.(map[string]interface{})))
	}

	return items
}

// expandOSPolicyAssignmentFileGcs expands an instance of OSPolicyAssignmentFileGcs into a JSON
// request object.
func expandOSPolicyAssignmentFileGcs(c *Client, f *OSPolicyAssignmentFileGcs) (map[string]interface{}, error) {
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
	if v := f.Generation; !dcl.IsEmptyValueIndirect(v) {
		m["generation"] = v
	}

	return m, nil
}

// flattenOSPolicyAssignmentFileGcs flattens an instance of OSPolicyAssignmentFileGcs from a JSON
// response object.
func flattenOSPolicyAssignmentFileGcs(c *Client, i interface{}) *OSPolicyAssignmentFileGcs {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &OSPolicyAssignmentFileGcs{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyOSPolicyAssignmentFileGcs
	}
	r.Bucket = dcl.FlattenString(m["bucket"])
	r.Object = dcl.FlattenString(m["object"])
	r.Generation = dcl.FlattenInteger(m["generation"])

	return r
}

// expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgYumMap expands the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgYum into a JSON
// request object.
func expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgYumMap(c *Client, f map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgYum) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgYum(c, &item)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgYumSlice expands the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgYum into a JSON
// request object.
func expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgYumSlice(c *Client, f []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgYum) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgYum(c, &item)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgYumMap flattens the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgYum from a JSON
// response object.
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgYumMap(c *Client, i interface{}) map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgYum {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgYum{}
	}

	if len(a) == 0 {
		return map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgYum{}
	}

	items := make(map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgYum)
	for k, item := range a {
		items[k] = *flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgYum(c, item.(map[string]interface{}))
	}

	return items
}

// flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgYumSlice flattens the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgYum from a JSON
// response object.
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgYumSlice(c *Client, i interface{}) []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgYum {
	a, ok := i.([]interface{})
	if !ok {
		return []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgYum{}
	}

	if len(a) == 0 {
		return []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgYum{}
	}

	items := make([]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgYum, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgYum(c, item.(map[string]interface{})))
	}

	return items
}

// expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgYum expands an instance of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgYum into a JSON
// request object.
func expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgYum(c *Client, f *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgYum) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.Name; !dcl.IsEmptyValueIndirect(v) {
		m["name"] = v
	}

	return m, nil
}

// flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgYum flattens an instance of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgYum from a JSON
// response object.
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgYum(c *Client, i interface{}) *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgYum {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgYum{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgYum
	}
	r.Name = dcl.FlattenString(m["name"])

	return r
}

// expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgZypperMap expands the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgZypper into a JSON
// request object.
func expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgZypperMap(c *Client, f map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgZypper) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgZypper(c, &item)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgZypperSlice expands the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgZypper into a JSON
// request object.
func expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgZypperSlice(c *Client, f []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgZypper) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgZypper(c, &item)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgZypperMap flattens the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgZypper from a JSON
// response object.
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgZypperMap(c *Client, i interface{}) map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgZypper {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgZypper{}
	}

	if len(a) == 0 {
		return map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgZypper{}
	}

	items := make(map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgZypper)
	for k, item := range a {
		items[k] = *flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgZypper(c, item.(map[string]interface{}))
	}

	return items
}

// flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgZypperSlice flattens the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgZypper from a JSON
// response object.
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgZypperSlice(c *Client, i interface{}) []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgZypper {
	a, ok := i.([]interface{})
	if !ok {
		return []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgZypper{}
	}

	if len(a) == 0 {
		return []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgZypper{}
	}

	items := make([]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgZypper, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgZypper(c, item.(map[string]interface{})))
	}

	return items
}

// expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgZypper expands an instance of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgZypper into a JSON
// request object.
func expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgZypper(c *Client, f *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgZypper) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.Name; !dcl.IsEmptyValueIndirect(v) {
		m["name"] = v
	}

	return m, nil
}

// flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgZypper flattens an instance of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgZypper from a JSON
// response object.
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgZypper(c *Client, i interface{}) *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgZypper {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgZypper{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgZypper
	}
	r.Name = dcl.FlattenString(m["name"])

	return r
}

// expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmMap expands the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpm into a JSON
// request object.
func expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmMap(c *Client, f map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpm) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpm(c, &item)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSlice expands the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpm into a JSON
// request object.
func expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSlice(c *Client, f []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpm) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpm(c, &item)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmMap flattens the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpm from a JSON
// response object.
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmMap(c *Client, i interface{}) map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpm {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpm{}
	}

	if len(a) == 0 {
		return map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpm{}
	}

	items := make(map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpm)
	for k, item := range a {
		items[k] = *flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpm(c, item.(map[string]interface{}))
	}

	return items
}

// flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSlice flattens the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpm from a JSON
// response object.
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSlice(c *Client, i interface{}) []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpm {
	a, ok := i.([]interface{})
	if !ok {
		return []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpm{}
	}

	if len(a) == 0 {
		return []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpm{}
	}

	items := make([]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpm, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpm(c, item.(map[string]interface{})))
	}

	return items
}

// expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpm expands an instance of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpm into a JSON
// request object.
func expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpm(c *Client, f *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpm) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v, err := expandOSPolicyAssignmentFile(c, f.Source); err != nil {
		return nil, fmt.Errorf("error expanding Source into source: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["source"] = v
	}
	if v := f.PullDeps; !dcl.IsEmptyValueIndirect(v) {
		m["pullDeps"] = v
	}

	return m, nil
}

// flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpm flattens an instance of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpm from a JSON
// response object.
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpm(c *Client, i interface{}) *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpm {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpm{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpm
	}
	r.Source = flattenOSPolicyAssignmentFile(c, m["source"])
	r.PullDeps = dcl.FlattenBool(m["pullDeps"])

	return r
}

// expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgGoogetMap expands the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgGooget into a JSON
// request object.
func expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgGoogetMap(c *Client, f map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgGooget) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgGooget(c, &item)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgGoogetSlice expands the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgGooget into a JSON
// request object.
func expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgGoogetSlice(c *Client, f []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgGooget) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgGooget(c, &item)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgGoogetMap flattens the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgGooget from a JSON
// response object.
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgGoogetMap(c *Client, i interface{}) map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgGooget {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgGooget{}
	}

	if len(a) == 0 {
		return map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgGooget{}
	}

	items := make(map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgGooget)
	for k, item := range a {
		items[k] = *flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgGooget(c, item.(map[string]interface{}))
	}

	return items
}

// flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgGoogetSlice flattens the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgGooget from a JSON
// response object.
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgGoogetSlice(c *Client, i interface{}) []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgGooget {
	a, ok := i.([]interface{})
	if !ok {
		return []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgGooget{}
	}

	if len(a) == 0 {
		return []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgGooget{}
	}

	items := make([]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgGooget, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgGooget(c, item.(map[string]interface{})))
	}

	return items
}

// expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgGooget expands an instance of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgGooget into a JSON
// request object.
func expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgGooget(c *Client, f *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgGooget) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.Name; !dcl.IsEmptyValueIndirect(v) {
		m["name"] = v
	}

	return m, nil
}

// flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgGooget flattens an instance of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgGooget from a JSON
// response object.
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgGooget(c *Client, i interface{}) *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgGooget {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgGooget{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgGooget
	}
	r.Name = dcl.FlattenString(m["name"])

	return r
}

// expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiMap expands the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsi into a JSON
// request object.
func expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiMap(c *Client, f map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsi) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsi(c, &item)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSlice expands the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsi into a JSON
// request object.
func expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSlice(c *Client, f []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsi) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsi(c, &item)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiMap flattens the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsi from a JSON
// response object.
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiMap(c *Client, i interface{}) map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsi {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsi{}
	}

	if len(a) == 0 {
		return map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsi{}
	}

	items := make(map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsi)
	for k, item := range a {
		items[k] = *flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsi(c, item.(map[string]interface{}))
	}

	return items
}

// flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSlice flattens the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsi from a JSON
// response object.
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSlice(c *Client, i interface{}) []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsi {
	a, ok := i.([]interface{})
	if !ok {
		return []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsi{}
	}

	if len(a) == 0 {
		return []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsi{}
	}

	items := make([]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsi, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsi(c, item.(map[string]interface{})))
	}

	return items
}

// expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsi expands an instance of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsi into a JSON
// request object.
func expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsi(c *Client, f *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsi) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v, err := expandOSPolicyAssignmentFile(c, f.Source); err != nil {
		return nil, fmt.Errorf("error expanding Source into source: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["source"] = v
	}
	if v := f.Properties; v != nil {
		m["properties"] = v
	}

	return m, nil
}

// flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsi flattens an instance of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsi from a JSON
// response object.
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsi(c *Client, i interface{}) *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsi {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsi{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsi
	}
	r.Source = flattenOSPolicyAssignmentFile(c, m["source"])
	r.Properties = dcl.FlattenStringSlice(m["properties"])

	return r
}

// expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryMap expands the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepository into a JSON
// request object.
func expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryMap(c *Client, f map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepository) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepository(c, &item)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositorySlice expands the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepository into a JSON
// request object.
func expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositorySlice(c *Client, f []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepository) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepository(c, &item)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryMap flattens the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepository from a JSON
// response object.
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryMap(c *Client, i interface{}) map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepository {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepository{}
	}

	if len(a) == 0 {
		return map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepository{}
	}

	items := make(map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepository)
	for k, item := range a {
		items[k] = *flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepository(c, item.(map[string]interface{}))
	}

	return items
}

// flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositorySlice flattens the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepository from a JSON
// response object.
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositorySlice(c *Client, i interface{}) []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepository {
	a, ok := i.([]interface{})
	if !ok {
		return []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepository{}
	}

	if len(a) == 0 {
		return []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepository{}
	}

	items := make([]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepository, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepository(c, item.(map[string]interface{})))
	}

	return items
}

// expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepository expands an instance of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepository into a JSON
// request object.
func expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepository(c *Client, f *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepository) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v, err := expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryApt(c, f.Apt); err != nil {
		return nil, fmt.Errorf("error expanding Apt into apt: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["apt"] = v
	}
	if v, err := expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryYum(c, f.Yum); err != nil {
		return nil, fmt.Errorf("error expanding Yum into yum: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["yum"] = v
	}
	if v, err := expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryZypper(c, f.Zypper); err != nil {
		return nil, fmt.Errorf("error expanding Zypper into zypper: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["zypper"] = v
	}
	if v, err := expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryGoo(c, f.Goo); err != nil {
		return nil, fmt.Errorf("error expanding Goo into goo: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["goo"] = v
	}

	return m, nil
}

// flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepository flattens an instance of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepository from a JSON
// response object.
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepository(c *Client, i interface{}) *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepository {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepository{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepository
	}
	r.Apt = flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryApt(c, m["apt"])
	r.Yum = flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryYum(c, m["yum"])
	r.Zypper = flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryZypper(c, m["zypper"])
	r.Goo = flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryGoo(c, m["goo"])

	return r
}

// expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryAptMap expands the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryApt into a JSON
// request object.
func expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryAptMap(c *Client, f map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryApt) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryApt(c, &item)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryAptSlice expands the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryApt into a JSON
// request object.
func expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryAptSlice(c *Client, f []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryApt) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryApt(c, &item)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryAptMap flattens the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryApt from a JSON
// response object.
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryAptMap(c *Client, i interface{}) map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryApt {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryApt{}
	}

	if len(a) == 0 {
		return map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryApt{}
	}

	items := make(map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryApt)
	for k, item := range a {
		items[k] = *flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryApt(c, item.(map[string]interface{}))
	}

	return items
}

// flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryAptSlice flattens the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryApt from a JSON
// response object.
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryAptSlice(c *Client, i interface{}) []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryApt {
	a, ok := i.([]interface{})
	if !ok {
		return []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryApt{}
	}

	if len(a) == 0 {
		return []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryApt{}
	}

	items := make([]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryApt, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryApt(c, item.(map[string]interface{})))
	}

	return items
}

// expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryApt expands an instance of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryApt into a JSON
// request object.
func expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryApt(c *Client, f *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryApt) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.ArchiveType; !dcl.IsEmptyValueIndirect(v) {
		m["archiveType"] = v
	}
	if v := f.Uri; !dcl.IsEmptyValueIndirect(v) {
		m["uri"] = v
	}
	if v := f.Distribution; !dcl.IsEmptyValueIndirect(v) {
		m["distribution"] = v
	}
	if v := f.Components; v != nil {
		m["components"] = v
	}
	if v := f.GpgKey; !dcl.IsEmptyValueIndirect(v) {
		m["gpgKey"] = v
	}

	return m, nil
}

// flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryApt flattens an instance of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryApt from a JSON
// response object.
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryApt(c *Client, i interface{}) *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryApt {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryApt{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryApt
	}
	r.ArchiveType = flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryAptArchiveTypeEnum(m["archiveType"])
	r.Uri = dcl.FlattenString(m["uri"])
	r.Distribution = dcl.FlattenString(m["distribution"])
	r.Components = dcl.FlattenStringSlice(m["components"])
	r.GpgKey = dcl.FlattenString(m["gpgKey"])

	return r
}

// expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryYumMap expands the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryYum into a JSON
// request object.
func expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryYumMap(c *Client, f map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryYum) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryYum(c, &item)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryYumSlice expands the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryYum into a JSON
// request object.
func expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryYumSlice(c *Client, f []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryYum) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryYum(c, &item)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryYumMap flattens the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryYum from a JSON
// response object.
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryYumMap(c *Client, i interface{}) map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryYum {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryYum{}
	}

	if len(a) == 0 {
		return map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryYum{}
	}

	items := make(map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryYum)
	for k, item := range a {
		items[k] = *flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryYum(c, item.(map[string]interface{}))
	}

	return items
}

// flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryYumSlice flattens the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryYum from a JSON
// response object.
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryYumSlice(c *Client, i interface{}) []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryYum {
	a, ok := i.([]interface{})
	if !ok {
		return []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryYum{}
	}

	if len(a) == 0 {
		return []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryYum{}
	}

	items := make([]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryYum, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryYum(c, item.(map[string]interface{})))
	}

	return items
}

// expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryYum expands an instance of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryYum into a JSON
// request object.
func expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryYum(c *Client, f *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryYum) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.Id; !dcl.IsEmptyValueIndirect(v) {
		m["id"] = v
	}
	if v := f.DisplayName; !dcl.IsEmptyValueIndirect(v) {
		m["displayName"] = v
	}
	if v := f.BaseUrl; !dcl.IsEmptyValueIndirect(v) {
		m["baseUrl"] = v
	}
	if v := f.GpgKeys; v != nil {
		m["gpgKeys"] = v
	}

	return m, nil
}

// flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryYum flattens an instance of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryYum from a JSON
// response object.
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryYum(c *Client, i interface{}) *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryYum {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryYum{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryYum
	}
	r.Id = dcl.FlattenString(m["id"])
	r.DisplayName = dcl.FlattenString(m["displayName"])
	r.BaseUrl = dcl.FlattenString(m["baseUrl"])
	r.GpgKeys = dcl.FlattenStringSlice(m["gpgKeys"])

	return r
}

// expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryZypperMap expands the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryZypper into a JSON
// request object.
func expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryZypperMap(c *Client, f map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryZypper) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryZypper(c, &item)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryZypperSlice expands the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryZypper into a JSON
// request object.
func expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryZypperSlice(c *Client, f []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryZypper) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryZypper(c, &item)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryZypperMap flattens the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryZypper from a JSON
// response object.
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryZypperMap(c *Client, i interface{}) map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryZypper {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryZypper{}
	}

	if len(a) == 0 {
		return map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryZypper{}
	}

	items := make(map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryZypper)
	for k, item := range a {
		items[k] = *flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryZypper(c, item.(map[string]interface{}))
	}

	return items
}

// flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryZypperSlice flattens the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryZypper from a JSON
// response object.
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryZypperSlice(c *Client, i interface{}) []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryZypper {
	a, ok := i.([]interface{})
	if !ok {
		return []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryZypper{}
	}

	if len(a) == 0 {
		return []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryZypper{}
	}

	items := make([]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryZypper, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryZypper(c, item.(map[string]interface{})))
	}

	return items
}

// expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryZypper expands an instance of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryZypper into a JSON
// request object.
func expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryZypper(c *Client, f *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryZypper) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.Id; !dcl.IsEmptyValueIndirect(v) {
		m["id"] = v
	}
	if v := f.DisplayName; !dcl.IsEmptyValueIndirect(v) {
		m["displayName"] = v
	}
	if v := f.BaseUrl; !dcl.IsEmptyValueIndirect(v) {
		m["baseUrl"] = v
	}
	if v := f.GpgKeys; v != nil {
		m["gpgKeys"] = v
	}

	return m, nil
}

// flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryZypper flattens an instance of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryZypper from a JSON
// response object.
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryZypper(c *Client, i interface{}) *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryZypper {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryZypper{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryZypper
	}
	r.Id = dcl.FlattenString(m["id"])
	r.DisplayName = dcl.FlattenString(m["displayName"])
	r.BaseUrl = dcl.FlattenString(m["baseUrl"])
	r.GpgKeys = dcl.FlattenStringSlice(m["gpgKeys"])

	return r
}

// expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryGooMap expands the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryGoo into a JSON
// request object.
func expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryGooMap(c *Client, f map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryGoo) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryGoo(c, &item)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryGooSlice expands the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryGoo into a JSON
// request object.
func expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryGooSlice(c *Client, f []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryGoo) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryGoo(c, &item)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryGooMap flattens the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryGoo from a JSON
// response object.
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryGooMap(c *Client, i interface{}) map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryGoo {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryGoo{}
	}

	if len(a) == 0 {
		return map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryGoo{}
	}

	items := make(map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryGoo)
	for k, item := range a {
		items[k] = *flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryGoo(c, item.(map[string]interface{}))
	}

	return items
}

// flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryGooSlice flattens the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryGoo from a JSON
// response object.
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryGooSlice(c *Client, i interface{}) []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryGoo {
	a, ok := i.([]interface{})
	if !ok {
		return []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryGoo{}
	}

	if len(a) == 0 {
		return []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryGoo{}
	}

	items := make([]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryGoo, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryGoo(c, item.(map[string]interface{})))
	}

	return items
}

// expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryGoo expands an instance of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryGoo into a JSON
// request object.
func expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryGoo(c *Client, f *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryGoo) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.Name; !dcl.IsEmptyValueIndirect(v) {
		m["name"] = v
	}
	if v := f.Url; !dcl.IsEmptyValueIndirect(v) {
		m["url"] = v
	}

	return m, nil
}

// flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryGoo flattens an instance of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryGoo from a JSON
// response object.
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryGoo(c *Client, i interface{}) *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryGoo {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryGoo{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryGoo
	}
	r.Name = dcl.FlattenString(m["name"])
	r.Url = dcl.FlattenString(m["url"])

	return r
}

// expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecMap expands the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExec into a JSON
// request object.
func expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecMap(c *Client, f map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExec) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExec(c, &item)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecSlice expands the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExec into a JSON
// request object.
func expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecSlice(c *Client, f []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExec) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExec(c, &item)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecMap flattens the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExec from a JSON
// response object.
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecMap(c *Client, i interface{}) map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExec {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExec{}
	}

	if len(a) == 0 {
		return map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExec{}
	}

	items := make(map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExec)
	for k, item := range a {
		items[k] = *flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExec(c, item.(map[string]interface{}))
	}

	return items
}

// flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecSlice flattens the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExec from a JSON
// response object.
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecSlice(c *Client, i interface{}) []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExec {
	a, ok := i.([]interface{})
	if !ok {
		return []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExec{}
	}

	if len(a) == 0 {
		return []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExec{}
	}

	items := make([]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExec, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExec(c, item.(map[string]interface{})))
	}

	return items
}

// expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExec expands an instance of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExec into a JSON
// request object.
func expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExec(c *Client, f *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExec) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v, err := expandOSPolicyAssignmentExec(c, f.Validate); err != nil {
		return nil, fmt.Errorf("error expanding Validate into validate: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["validate"] = v
	}
	if v, err := expandOSPolicyAssignmentExec(c, f.Enforce); err != nil {
		return nil, fmt.Errorf("error expanding Enforce into enforce: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["enforce"] = v
	}

	return m, nil
}

// flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExec flattens an instance of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExec from a JSON
// response object.
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExec(c *Client, i interface{}) *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExec {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExec{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExec
	}
	r.Validate = flattenOSPolicyAssignmentExec(c, m["validate"])
	r.Enforce = flattenOSPolicyAssignmentExec(c, m["enforce"])

	return r
}

// expandOSPolicyAssignmentExecMap expands the contents of OSPolicyAssignmentExec into a JSON
// request object.
func expandOSPolicyAssignmentExecMap(c *Client, f map[string]OSPolicyAssignmentExec) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandOSPolicyAssignmentExec(c, &item)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandOSPolicyAssignmentExecSlice expands the contents of OSPolicyAssignmentExec into a JSON
// request object.
func expandOSPolicyAssignmentExecSlice(c *Client, f []OSPolicyAssignmentExec) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandOSPolicyAssignmentExec(c, &item)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenOSPolicyAssignmentExecMap flattens the contents of OSPolicyAssignmentExec from a JSON
// response object.
func flattenOSPolicyAssignmentExecMap(c *Client, i interface{}) map[string]OSPolicyAssignmentExec {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]OSPolicyAssignmentExec{}
	}

	if len(a) == 0 {
		return map[string]OSPolicyAssignmentExec{}
	}

	items := make(map[string]OSPolicyAssignmentExec)
	for k, item := range a {
		items[k] = *flattenOSPolicyAssignmentExec(c, item.(map[string]interface{}))
	}

	return items
}

// flattenOSPolicyAssignmentExecSlice flattens the contents of OSPolicyAssignmentExec from a JSON
// response object.
func flattenOSPolicyAssignmentExecSlice(c *Client, i interface{}) []OSPolicyAssignmentExec {
	a, ok := i.([]interface{})
	if !ok {
		return []OSPolicyAssignmentExec{}
	}

	if len(a) == 0 {
		return []OSPolicyAssignmentExec{}
	}

	items := make([]OSPolicyAssignmentExec, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenOSPolicyAssignmentExec(c, item.(map[string]interface{})))
	}

	return items
}

// expandOSPolicyAssignmentExec expands an instance of OSPolicyAssignmentExec into a JSON
// request object.
func expandOSPolicyAssignmentExec(c *Client, f *OSPolicyAssignmentExec) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v, err := expandOSPolicyAssignmentFile(c, f.File); err != nil {
		return nil, fmt.Errorf("error expanding File into file: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["file"] = v
	}
	if v := f.Script; !dcl.IsEmptyValueIndirect(v) {
		m["script"] = v
	}
	if v := f.Args; v != nil {
		m["args"] = v
	}
	if v := f.Interpreter; !dcl.IsEmptyValueIndirect(v) {
		m["interpreter"] = v
	}
	if v := f.OutputFilePath; !dcl.IsEmptyValueIndirect(v) {
		m["outputFilePath"] = v
	}

	return m, nil
}

// flattenOSPolicyAssignmentExec flattens an instance of OSPolicyAssignmentExec from a JSON
// response object.
func flattenOSPolicyAssignmentExec(c *Client, i interface{}) *OSPolicyAssignmentExec {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &OSPolicyAssignmentExec{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyOSPolicyAssignmentExec
	}
	r.File = flattenOSPolicyAssignmentFile(c, m["file"])
	r.Script = dcl.FlattenString(m["script"])
	r.Args = dcl.FlattenStringSlice(m["args"])
	r.Interpreter = flattenOSPolicyAssignmentExecInterpreterEnum(m["interpreter"])
	r.OutputFilePath = dcl.FlattenString(m["outputFilePath"])

	return r
}

// expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileMap expands the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFile into a JSON
// request object.
func expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileMap(c *Client, f map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFile) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFile(c, &item)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileSlice expands the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFile into a JSON
// request object.
func expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileSlice(c *Client, f []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFile) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFile(c, &item)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileMap flattens the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFile from a JSON
// response object.
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileMap(c *Client, i interface{}) map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFile {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFile{}
	}

	if len(a) == 0 {
		return map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFile{}
	}

	items := make(map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFile)
	for k, item := range a {
		items[k] = *flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFile(c, item.(map[string]interface{}))
	}

	return items
}

// flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileSlice flattens the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFile from a JSON
// response object.
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileSlice(c *Client, i interface{}) []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFile {
	a, ok := i.([]interface{})
	if !ok {
		return []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFile{}
	}

	if len(a) == 0 {
		return []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFile{}
	}

	items := make([]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFile, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFile(c, item.(map[string]interface{})))
	}

	return items
}

// expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFile expands an instance of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFile into a JSON
// request object.
func expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFile(c *Client, f *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFile) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v, err := expandOSPolicyAssignmentFile(c, f.File); err != nil {
		return nil, fmt.Errorf("error expanding File into file: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["file"] = v
	}
	if v := f.Content; !dcl.IsEmptyValueIndirect(v) {
		m["content"] = v
	}
	if v := f.Path; !dcl.IsEmptyValueIndirect(v) {
		m["path"] = v
	}
	if v := f.State; !dcl.IsEmptyValueIndirect(v) {
		m["state"] = v
	}

	return m, nil
}

// flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFile flattens an instance of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFile from a JSON
// response object.
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFile(c *Client, i interface{}) *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFile {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFile{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFile
	}
	r.File = flattenOSPolicyAssignmentFile(c, m["file"])
	r.Content = dcl.FlattenString(m["content"])
	r.Path = dcl.FlattenString(m["path"])
	r.State = flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileStateEnum(m["state"])
	r.Permissions = dcl.FlattenString(m["permissions"])

	return r
}

// expandOSPolicyAssignmentInstanceFilterMap expands the contents of OSPolicyAssignmentInstanceFilter into a JSON
// request object.
func expandOSPolicyAssignmentInstanceFilterMap(c *Client, f map[string]OSPolicyAssignmentInstanceFilter) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandOSPolicyAssignmentInstanceFilter(c, &item)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandOSPolicyAssignmentInstanceFilterSlice expands the contents of OSPolicyAssignmentInstanceFilter into a JSON
// request object.
func expandOSPolicyAssignmentInstanceFilterSlice(c *Client, f []OSPolicyAssignmentInstanceFilter) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandOSPolicyAssignmentInstanceFilter(c, &item)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenOSPolicyAssignmentInstanceFilterMap flattens the contents of OSPolicyAssignmentInstanceFilter from a JSON
// response object.
func flattenOSPolicyAssignmentInstanceFilterMap(c *Client, i interface{}) map[string]OSPolicyAssignmentInstanceFilter {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]OSPolicyAssignmentInstanceFilter{}
	}

	if len(a) == 0 {
		return map[string]OSPolicyAssignmentInstanceFilter{}
	}

	items := make(map[string]OSPolicyAssignmentInstanceFilter)
	for k, item := range a {
		items[k] = *flattenOSPolicyAssignmentInstanceFilter(c, item.(map[string]interface{}))
	}

	return items
}

// flattenOSPolicyAssignmentInstanceFilterSlice flattens the contents of OSPolicyAssignmentInstanceFilter from a JSON
// response object.
func flattenOSPolicyAssignmentInstanceFilterSlice(c *Client, i interface{}) []OSPolicyAssignmentInstanceFilter {
	a, ok := i.([]interface{})
	if !ok {
		return []OSPolicyAssignmentInstanceFilter{}
	}

	if len(a) == 0 {
		return []OSPolicyAssignmentInstanceFilter{}
	}

	items := make([]OSPolicyAssignmentInstanceFilter, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenOSPolicyAssignmentInstanceFilter(c, item.(map[string]interface{})))
	}

	return items
}

// expandOSPolicyAssignmentInstanceFilter expands an instance of OSPolicyAssignmentInstanceFilter into a JSON
// request object.
func expandOSPolicyAssignmentInstanceFilter(c *Client, f *OSPolicyAssignmentInstanceFilter) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.All; v != nil {
		m["all"] = v
	}
	if v, err := expandOSPolicyAssignmentInstanceFilterInclusionLabelsSlice(c, f.InclusionLabels); err != nil {
		return nil, fmt.Errorf("error expanding InclusionLabels into inclusionLabels: %w", err)
	} else if v != nil {
		m["inclusionLabels"] = v
	}
	if v, err := expandOSPolicyAssignmentInstanceFilterExclusionLabelsSlice(c, f.ExclusionLabels); err != nil {
		return nil, fmt.Errorf("error expanding ExclusionLabels into exclusionLabels: %w", err)
	} else if v != nil {
		m["exclusionLabels"] = v
	}
	if v, err := expandOSPolicyAssignmentInstanceFilterInventoriesSlice(c, f.Inventories); err != nil {
		return nil, fmt.Errorf("error expanding Inventories into inventories: %w", err)
	} else if v != nil {
		m["inventories"] = v
	}

	return m, nil
}

// flattenOSPolicyAssignmentInstanceFilter flattens an instance of OSPolicyAssignmentInstanceFilter from a JSON
// response object.
func flattenOSPolicyAssignmentInstanceFilter(c *Client, i interface{}) *OSPolicyAssignmentInstanceFilter {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &OSPolicyAssignmentInstanceFilter{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyOSPolicyAssignmentInstanceFilter
	}
	r.All = dcl.FlattenBool(m["all"])
	r.InclusionLabels = flattenOSPolicyAssignmentInstanceFilterInclusionLabelsSlice(c, m["inclusionLabels"])
	r.ExclusionLabels = flattenOSPolicyAssignmentInstanceFilterExclusionLabelsSlice(c, m["exclusionLabels"])
	r.Inventories = flattenOSPolicyAssignmentInstanceFilterInventoriesSlice(c, m["inventories"])

	return r
}

// expandOSPolicyAssignmentInstanceFilterInclusionLabelsMap expands the contents of OSPolicyAssignmentInstanceFilterInclusionLabels into a JSON
// request object.
func expandOSPolicyAssignmentInstanceFilterInclusionLabelsMap(c *Client, f map[string]OSPolicyAssignmentInstanceFilterInclusionLabels) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandOSPolicyAssignmentInstanceFilterInclusionLabels(c, &item)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandOSPolicyAssignmentInstanceFilterInclusionLabelsSlice expands the contents of OSPolicyAssignmentInstanceFilterInclusionLabels into a JSON
// request object.
func expandOSPolicyAssignmentInstanceFilterInclusionLabelsSlice(c *Client, f []OSPolicyAssignmentInstanceFilterInclusionLabels) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandOSPolicyAssignmentInstanceFilterInclusionLabels(c, &item)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenOSPolicyAssignmentInstanceFilterInclusionLabelsMap flattens the contents of OSPolicyAssignmentInstanceFilterInclusionLabels from a JSON
// response object.
func flattenOSPolicyAssignmentInstanceFilterInclusionLabelsMap(c *Client, i interface{}) map[string]OSPolicyAssignmentInstanceFilterInclusionLabels {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]OSPolicyAssignmentInstanceFilterInclusionLabels{}
	}

	if len(a) == 0 {
		return map[string]OSPolicyAssignmentInstanceFilterInclusionLabels{}
	}

	items := make(map[string]OSPolicyAssignmentInstanceFilterInclusionLabels)
	for k, item := range a {
		items[k] = *flattenOSPolicyAssignmentInstanceFilterInclusionLabels(c, item.(map[string]interface{}))
	}

	return items
}

// flattenOSPolicyAssignmentInstanceFilterInclusionLabelsSlice flattens the contents of OSPolicyAssignmentInstanceFilterInclusionLabels from a JSON
// response object.
func flattenOSPolicyAssignmentInstanceFilterInclusionLabelsSlice(c *Client, i interface{}) []OSPolicyAssignmentInstanceFilterInclusionLabels {
	a, ok := i.([]interface{})
	if !ok {
		return []OSPolicyAssignmentInstanceFilterInclusionLabels{}
	}

	if len(a) == 0 {
		return []OSPolicyAssignmentInstanceFilterInclusionLabels{}
	}

	items := make([]OSPolicyAssignmentInstanceFilterInclusionLabels, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenOSPolicyAssignmentInstanceFilterInclusionLabels(c, item.(map[string]interface{})))
	}

	return items
}

// expandOSPolicyAssignmentInstanceFilterInclusionLabels expands an instance of OSPolicyAssignmentInstanceFilterInclusionLabels into a JSON
// request object.
func expandOSPolicyAssignmentInstanceFilterInclusionLabels(c *Client, f *OSPolicyAssignmentInstanceFilterInclusionLabels) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.Labels; !dcl.IsEmptyValueIndirect(v) {
		m["labels"] = v
	}

	return m, nil
}

// flattenOSPolicyAssignmentInstanceFilterInclusionLabels flattens an instance of OSPolicyAssignmentInstanceFilterInclusionLabels from a JSON
// response object.
func flattenOSPolicyAssignmentInstanceFilterInclusionLabels(c *Client, i interface{}) *OSPolicyAssignmentInstanceFilterInclusionLabels {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &OSPolicyAssignmentInstanceFilterInclusionLabels{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyOSPolicyAssignmentInstanceFilterInclusionLabels
	}
	r.Labels = dcl.FlattenKeyValuePairs(m["labels"])

	return r
}

// expandOSPolicyAssignmentInstanceFilterExclusionLabelsMap expands the contents of OSPolicyAssignmentInstanceFilterExclusionLabels into a JSON
// request object.
func expandOSPolicyAssignmentInstanceFilterExclusionLabelsMap(c *Client, f map[string]OSPolicyAssignmentInstanceFilterExclusionLabels) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandOSPolicyAssignmentInstanceFilterExclusionLabels(c, &item)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandOSPolicyAssignmentInstanceFilterExclusionLabelsSlice expands the contents of OSPolicyAssignmentInstanceFilterExclusionLabels into a JSON
// request object.
func expandOSPolicyAssignmentInstanceFilterExclusionLabelsSlice(c *Client, f []OSPolicyAssignmentInstanceFilterExclusionLabels) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandOSPolicyAssignmentInstanceFilterExclusionLabels(c, &item)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenOSPolicyAssignmentInstanceFilterExclusionLabelsMap flattens the contents of OSPolicyAssignmentInstanceFilterExclusionLabels from a JSON
// response object.
func flattenOSPolicyAssignmentInstanceFilterExclusionLabelsMap(c *Client, i interface{}) map[string]OSPolicyAssignmentInstanceFilterExclusionLabels {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]OSPolicyAssignmentInstanceFilterExclusionLabels{}
	}

	if len(a) == 0 {
		return map[string]OSPolicyAssignmentInstanceFilterExclusionLabels{}
	}

	items := make(map[string]OSPolicyAssignmentInstanceFilterExclusionLabels)
	for k, item := range a {
		items[k] = *flattenOSPolicyAssignmentInstanceFilterExclusionLabels(c, item.(map[string]interface{}))
	}

	return items
}

// flattenOSPolicyAssignmentInstanceFilterExclusionLabelsSlice flattens the contents of OSPolicyAssignmentInstanceFilterExclusionLabels from a JSON
// response object.
func flattenOSPolicyAssignmentInstanceFilterExclusionLabelsSlice(c *Client, i interface{}) []OSPolicyAssignmentInstanceFilterExclusionLabels {
	a, ok := i.([]interface{})
	if !ok {
		return []OSPolicyAssignmentInstanceFilterExclusionLabels{}
	}

	if len(a) == 0 {
		return []OSPolicyAssignmentInstanceFilterExclusionLabels{}
	}

	items := make([]OSPolicyAssignmentInstanceFilterExclusionLabels, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenOSPolicyAssignmentInstanceFilterExclusionLabels(c, item.(map[string]interface{})))
	}

	return items
}

// expandOSPolicyAssignmentInstanceFilterExclusionLabels expands an instance of OSPolicyAssignmentInstanceFilterExclusionLabels into a JSON
// request object.
func expandOSPolicyAssignmentInstanceFilterExclusionLabels(c *Client, f *OSPolicyAssignmentInstanceFilterExclusionLabels) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.Labels; !dcl.IsEmptyValueIndirect(v) {
		m["labels"] = v
	}

	return m, nil
}

// flattenOSPolicyAssignmentInstanceFilterExclusionLabels flattens an instance of OSPolicyAssignmentInstanceFilterExclusionLabels from a JSON
// response object.
func flattenOSPolicyAssignmentInstanceFilterExclusionLabels(c *Client, i interface{}) *OSPolicyAssignmentInstanceFilterExclusionLabels {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &OSPolicyAssignmentInstanceFilterExclusionLabels{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyOSPolicyAssignmentInstanceFilterExclusionLabels
	}
	r.Labels = dcl.FlattenKeyValuePairs(m["labels"])

	return r
}

// expandOSPolicyAssignmentInstanceFilterInventoriesMap expands the contents of OSPolicyAssignmentInstanceFilterInventories into a JSON
// request object.
func expandOSPolicyAssignmentInstanceFilterInventoriesMap(c *Client, f map[string]OSPolicyAssignmentInstanceFilterInventories) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandOSPolicyAssignmentInstanceFilterInventories(c, &item)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandOSPolicyAssignmentInstanceFilterInventoriesSlice expands the contents of OSPolicyAssignmentInstanceFilterInventories into a JSON
// request object.
func expandOSPolicyAssignmentInstanceFilterInventoriesSlice(c *Client, f []OSPolicyAssignmentInstanceFilterInventories) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandOSPolicyAssignmentInstanceFilterInventories(c, &item)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenOSPolicyAssignmentInstanceFilterInventoriesMap flattens the contents of OSPolicyAssignmentInstanceFilterInventories from a JSON
// response object.
func flattenOSPolicyAssignmentInstanceFilterInventoriesMap(c *Client, i interface{}) map[string]OSPolicyAssignmentInstanceFilterInventories {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]OSPolicyAssignmentInstanceFilterInventories{}
	}

	if len(a) == 0 {
		return map[string]OSPolicyAssignmentInstanceFilterInventories{}
	}

	items := make(map[string]OSPolicyAssignmentInstanceFilterInventories)
	for k, item := range a {
		items[k] = *flattenOSPolicyAssignmentInstanceFilterInventories(c, item.(map[string]interface{}))
	}

	return items
}

// flattenOSPolicyAssignmentInstanceFilterInventoriesSlice flattens the contents of OSPolicyAssignmentInstanceFilterInventories from a JSON
// response object.
func flattenOSPolicyAssignmentInstanceFilterInventoriesSlice(c *Client, i interface{}) []OSPolicyAssignmentInstanceFilterInventories {
	a, ok := i.([]interface{})
	if !ok {
		return []OSPolicyAssignmentInstanceFilterInventories{}
	}

	if len(a) == 0 {
		return []OSPolicyAssignmentInstanceFilterInventories{}
	}

	items := make([]OSPolicyAssignmentInstanceFilterInventories, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenOSPolicyAssignmentInstanceFilterInventories(c, item.(map[string]interface{})))
	}

	return items
}

// expandOSPolicyAssignmentInstanceFilterInventories expands an instance of OSPolicyAssignmentInstanceFilterInventories into a JSON
// request object.
func expandOSPolicyAssignmentInstanceFilterInventories(c *Client, f *OSPolicyAssignmentInstanceFilterInventories) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.OSShortName; !dcl.IsEmptyValueIndirect(v) {
		m["osShortName"] = v
	}
	if v := f.OSVersion; !dcl.IsEmptyValueIndirect(v) {
		m["osVersion"] = v
	}

	return m, nil
}

// flattenOSPolicyAssignmentInstanceFilterInventories flattens an instance of OSPolicyAssignmentInstanceFilterInventories from a JSON
// response object.
func flattenOSPolicyAssignmentInstanceFilterInventories(c *Client, i interface{}) *OSPolicyAssignmentInstanceFilterInventories {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &OSPolicyAssignmentInstanceFilterInventories{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyOSPolicyAssignmentInstanceFilterInventories
	}
	r.OSShortName = dcl.FlattenString(m["osShortName"])
	r.OSVersion = dcl.FlattenString(m["osVersion"])

	return r
}

// expandOSPolicyAssignmentRolloutMap expands the contents of OSPolicyAssignmentRollout into a JSON
// request object.
func expandOSPolicyAssignmentRolloutMap(c *Client, f map[string]OSPolicyAssignmentRollout) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandOSPolicyAssignmentRollout(c, &item)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandOSPolicyAssignmentRolloutSlice expands the contents of OSPolicyAssignmentRollout into a JSON
// request object.
func expandOSPolicyAssignmentRolloutSlice(c *Client, f []OSPolicyAssignmentRollout) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandOSPolicyAssignmentRollout(c, &item)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenOSPolicyAssignmentRolloutMap flattens the contents of OSPolicyAssignmentRollout from a JSON
// response object.
func flattenOSPolicyAssignmentRolloutMap(c *Client, i interface{}) map[string]OSPolicyAssignmentRollout {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]OSPolicyAssignmentRollout{}
	}

	if len(a) == 0 {
		return map[string]OSPolicyAssignmentRollout{}
	}

	items := make(map[string]OSPolicyAssignmentRollout)
	for k, item := range a {
		items[k] = *flattenOSPolicyAssignmentRollout(c, item.(map[string]interface{}))
	}

	return items
}

// flattenOSPolicyAssignmentRolloutSlice flattens the contents of OSPolicyAssignmentRollout from a JSON
// response object.
func flattenOSPolicyAssignmentRolloutSlice(c *Client, i interface{}) []OSPolicyAssignmentRollout {
	a, ok := i.([]interface{})
	if !ok {
		return []OSPolicyAssignmentRollout{}
	}

	if len(a) == 0 {
		return []OSPolicyAssignmentRollout{}
	}

	items := make([]OSPolicyAssignmentRollout, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenOSPolicyAssignmentRollout(c, item.(map[string]interface{})))
	}

	return items
}

// expandOSPolicyAssignmentRollout expands an instance of OSPolicyAssignmentRollout into a JSON
// request object.
func expandOSPolicyAssignmentRollout(c *Client, f *OSPolicyAssignmentRollout) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v, err := expandOSPolicyAssignmentRolloutDisruptionBudget(c, f.DisruptionBudget); err != nil {
		return nil, fmt.Errorf("error expanding DisruptionBudget into disruptionBudget: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["disruptionBudget"] = v
	}
	if v := f.MinWaitDuration; !dcl.IsEmptyValueIndirect(v) {
		m["minWaitDuration"] = v
	}

	return m, nil
}

// flattenOSPolicyAssignmentRollout flattens an instance of OSPolicyAssignmentRollout from a JSON
// response object.
func flattenOSPolicyAssignmentRollout(c *Client, i interface{}) *OSPolicyAssignmentRollout {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &OSPolicyAssignmentRollout{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyOSPolicyAssignmentRollout
	}
	r.DisruptionBudget = flattenOSPolicyAssignmentRolloutDisruptionBudget(c, m["disruptionBudget"])
	r.MinWaitDuration = dcl.FlattenString(m["minWaitDuration"])

	return r
}

// expandOSPolicyAssignmentRolloutDisruptionBudgetMap expands the contents of OSPolicyAssignmentRolloutDisruptionBudget into a JSON
// request object.
func expandOSPolicyAssignmentRolloutDisruptionBudgetMap(c *Client, f map[string]OSPolicyAssignmentRolloutDisruptionBudget) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandOSPolicyAssignmentRolloutDisruptionBudget(c, &item)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandOSPolicyAssignmentRolloutDisruptionBudgetSlice expands the contents of OSPolicyAssignmentRolloutDisruptionBudget into a JSON
// request object.
func expandOSPolicyAssignmentRolloutDisruptionBudgetSlice(c *Client, f []OSPolicyAssignmentRolloutDisruptionBudget) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandOSPolicyAssignmentRolloutDisruptionBudget(c, &item)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenOSPolicyAssignmentRolloutDisruptionBudgetMap flattens the contents of OSPolicyAssignmentRolloutDisruptionBudget from a JSON
// response object.
func flattenOSPolicyAssignmentRolloutDisruptionBudgetMap(c *Client, i interface{}) map[string]OSPolicyAssignmentRolloutDisruptionBudget {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]OSPolicyAssignmentRolloutDisruptionBudget{}
	}

	if len(a) == 0 {
		return map[string]OSPolicyAssignmentRolloutDisruptionBudget{}
	}

	items := make(map[string]OSPolicyAssignmentRolloutDisruptionBudget)
	for k, item := range a {
		items[k] = *flattenOSPolicyAssignmentRolloutDisruptionBudget(c, item.(map[string]interface{}))
	}

	return items
}

// flattenOSPolicyAssignmentRolloutDisruptionBudgetSlice flattens the contents of OSPolicyAssignmentRolloutDisruptionBudget from a JSON
// response object.
func flattenOSPolicyAssignmentRolloutDisruptionBudgetSlice(c *Client, i interface{}) []OSPolicyAssignmentRolloutDisruptionBudget {
	a, ok := i.([]interface{})
	if !ok {
		return []OSPolicyAssignmentRolloutDisruptionBudget{}
	}

	if len(a) == 0 {
		return []OSPolicyAssignmentRolloutDisruptionBudget{}
	}

	items := make([]OSPolicyAssignmentRolloutDisruptionBudget, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenOSPolicyAssignmentRolloutDisruptionBudget(c, item.(map[string]interface{})))
	}

	return items
}

// expandOSPolicyAssignmentRolloutDisruptionBudget expands an instance of OSPolicyAssignmentRolloutDisruptionBudget into a JSON
// request object.
func expandOSPolicyAssignmentRolloutDisruptionBudget(c *Client, f *OSPolicyAssignmentRolloutDisruptionBudget) (map[string]interface{}, error) {
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

// flattenOSPolicyAssignmentRolloutDisruptionBudget flattens an instance of OSPolicyAssignmentRolloutDisruptionBudget from a JSON
// response object.
func flattenOSPolicyAssignmentRolloutDisruptionBudget(c *Client, i interface{}) *OSPolicyAssignmentRolloutDisruptionBudget {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &OSPolicyAssignmentRolloutDisruptionBudget{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyOSPolicyAssignmentRolloutDisruptionBudget
	}
	r.Fixed = dcl.FlattenInteger(m["fixed"])
	r.Percent = dcl.FlattenInteger(m["percent"])

	return r
}

// flattenOSPolicyAssignmentOSPoliciesModeEnumMap flattens the contents of OSPolicyAssignmentOSPoliciesModeEnum from a JSON
// response object.
func flattenOSPolicyAssignmentOSPoliciesModeEnumMap(c *Client, i interface{}) map[string]OSPolicyAssignmentOSPoliciesModeEnum {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]OSPolicyAssignmentOSPoliciesModeEnum{}
	}

	if len(a) == 0 {
		return map[string]OSPolicyAssignmentOSPoliciesModeEnum{}
	}

	items := make(map[string]OSPolicyAssignmentOSPoliciesModeEnum)
	for k, item := range a {
		items[k] = *flattenOSPolicyAssignmentOSPoliciesModeEnum(item.(interface{}))
	}

	return items
}

// flattenOSPolicyAssignmentOSPoliciesModeEnumSlice flattens the contents of OSPolicyAssignmentOSPoliciesModeEnum from a JSON
// response object.
func flattenOSPolicyAssignmentOSPoliciesModeEnumSlice(c *Client, i interface{}) []OSPolicyAssignmentOSPoliciesModeEnum {
	a, ok := i.([]interface{})
	if !ok {
		return []OSPolicyAssignmentOSPoliciesModeEnum{}
	}

	if len(a) == 0 {
		return []OSPolicyAssignmentOSPoliciesModeEnum{}
	}

	items := make([]OSPolicyAssignmentOSPoliciesModeEnum, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenOSPolicyAssignmentOSPoliciesModeEnum(item.(interface{})))
	}

	return items
}

// flattenOSPolicyAssignmentOSPoliciesModeEnum asserts that an interface is a string, and returns a
// pointer to a *OSPolicyAssignmentOSPoliciesModeEnum with the same value as that string.
func flattenOSPolicyAssignmentOSPoliciesModeEnum(i interface{}) *OSPolicyAssignmentOSPoliciesModeEnum {
	s, ok := i.(string)
	if !ok {
		return OSPolicyAssignmentOSPoliciesModeEnumRef("")
	}

	return OSPolicyAssignmentOSPoliciesModeEnumRef(s)
}

// flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDesiredStateEnumMap flattens the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDesiredStateEnum from a JSON
// response object.
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDesiredStateEnumMap(c *Client, i interface{}) map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDesiredStateEnum {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDesiredStateEnum{}
	}

	if len(a) == 0 {
		return map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDesiredStateEnum{}
	}

	items := make(map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDesiredStateEnum)
	for k, item := range a {
		items[k] = *flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDesiredStateEnum(item.(interface{}))
	}

	return items
}

// flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDesiredStateEnumSlice flattens the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDesiredStateEnum from a JSON
// response object.
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDesiredStateEnumSlice(c *Client, i interface{}) []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDesiredStateEnum {
	a, ok := i.([]interface{})
	if !ok {
		return []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDesiredStateEnum{}
	}

	if len(a) == 0 {
		return []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDesiredStateEnum{}
	}

	items := make([]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDesiredStateEnum, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDesiredStateEnum(item.(interface{})))
	}

	return items
}

// flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDesiredStateEnum asserts that an interface is a string, and returns a
// pointer to a *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDesiredStateEnum with the same value as that string.
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDesiredStateEnum(i interface{}) *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDesiredStateEnum {
	s, ok := i.(string)
	if !ok {
		return OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDesiredStateEnumRef("")
	}

	return OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDesiredStateEnumRef(s)
}

// flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryAptArchiveTypeEnumMap flattens the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryAptArchiveTypeEnum from a JSON
// response object.
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryAptArchiveTypeEnumMap(c *Client, i interface{}) map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryAptArchiveTypeEnum {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryAptArchiveTypeEnum{}
	}

	if len(a) == 0 {
		return map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryAptArchiveTypeEnum{}
	}

	items := make(map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryAptArchiveTypeEnum)
	for k, item := range a {
		items[k] = *flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryAptArchiveTypeEnum(item.(interface{}))
	}

	return items
}

// flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryAptArchiveTypeEnumSlice flattens the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryAptArchiveTypeEnum from a JSON
// response object.
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryAptArchiveTypeEnumSlice(c *Client, i interface{}) []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryAptArchiveTypeEnum {
	a, ok := i.([]interface{})
	if !ok {
		return []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryAptArchiveTypeEnum{}
	}

	if len(a) == 0 {
		return []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryAptArchiveTypeEnum{}
	}

	items := make([]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryAptArchiveTypeEnum, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryAptArchiveTypeEnum(item.(interface{})))
	}

	return items
}

// flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryAptArchiveTypeEnum asserts that an interface is a string, and returns a
// pointer to a *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryAptArchiveTypeEnum with the same value as that string.
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryAptArchiveTypeEnum(i interface{}) *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryAptArchiveTypeEnum {
	s, ok := i.(string)
	if !ok {
		return OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryAptArchiveTypeEnumRef("")
	}

	return OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryAptArchiveTypeEnumRef(s)
}

// flattenOSPolicyAssignmentExecInterpreterEnumMap flattens the contents of OSPolicyAssignmentExecInterpreterEnum from a JSON
// response object.
func flattenOSPolicyAssignmentExecInterpreterEnumMap(c *Client, i interface{}) map[string]OSPolicyAssignmentExecInterpreterEnum {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]OSPolicyAssignmentExecInterpreterEnum{}
	}

	if len(a) == 0 {
		return map[string]OSPolicyAssignmentExecInterpreterEnum{}
	}

	items := make(map[string]OSPolicyAssignmentExecInterpreterEnum)
	for k, item := range a {
		items[k] = *flattenOSPolicyAssignmentExecInterpreterEnum(item.(interface{}))
	}

	return items
}

// flattenOSPolicyAssignmentExecInterpreterEnumSlice flattens the contents of OSPolicyAssignmentExecInterpreterEnum from a JSON
// response object.
func flattenOSPolicyAssignmentExecInterpreterEnumSlice(c *Client, i interface{}) []OSPolicyAssignmentExecInterpreterEnum {
	a, ok := i.([]interface{})
	if !ok {
		return []OSPolicyAssignmentExecInterpreterEnum{}
	}

	if len(a) == 0 {
		return []OSPolicyAssignmentExecInterpreterEnum{}
	}

	items := make([]OSPolicyAssignmentExecInterpreterEnum, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenOSPolicyAssignmentExecInterpreterEnum(item.(interface{})))
	}

	return items
}

// flattenOSPolicyAssignmentExecInterpreterEnum asserts that an interface is a string, and returns a
// pointer to a *OSPolicyAssignmentExecInterpreterEnum with the same value as that string.
func flattenOSPolicyAssignmentExecInterpreterEnum(i interface{}) *OSPolicyAssignmentExecInterpreterEnum {
	s, ok := i.(string)
	if !ok {
		return OSPolicyAssignmentExecInterpreterEnumRef("")
	}

	return OSPolicyAssignmentExecInterpreterEnumRef(s)
}

// flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileStateEnumMap flattens the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileStateEnum from a JSON
// response object.
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileStateEnumMap(c *Client, i interface{}) map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileStateEnum {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileStateEnum{}
	}

	if len(a) == 0 {
		return map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileStateEnum{}
	}

	items := make(map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileStateEnum)
	for k, item := range a {
		items[k] = *flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileStateEnum(item.(interface{}))
	}

	return items
}

// flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileStateEnumSlice flattens the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileStateEnum from a JSON
// response object.
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileStateEnumSlice(c *Client, i interface{}) []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileStateEnum {
	a, ok := i.([]interface{})
	if !ok {
		return []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileStateEnum{}
	}

	if len(a) == 0 {
		return []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileStateEnum{}
	}

	items := make([]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileStateEnum, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileStateEnum(item.(interface{})))
	}

	return items
}

// flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileStateEnum asserts that an interface is a string, and returns a
// pointer to a *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileStateEnum with the same value as that string.
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileStateEnum(i interface{}) *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileStateEnum {
	s, ok := i.(string)
	if !ok {
		return OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileStateEnumRef("")
	}

	return OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileStateEnumRef(s)
}

// flattenOSPolicyAssignmentRolloutStateEnumMap flattens the contents of OSPolicyAssignmentRolloutStateEnum from a JSON
// response object.
func flattenOSPolicyAssignmentRolloutStateEnumMap(c *Client, i interface{}) map[string]OSPolicyAssignmentRolloutStateEnum {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]OSPolicyAssignmentRolloutStateEnum{}
	}

	if len(a) == 0 {
		return map[string]OSPolicyAssignmentRolloutStateEnum{}
	}

	items := make(map[string]OSPolicyAssignmentRolloutStateEnum)
	for k, item := range a {
		items[k] = *flattenOSPolicyAssignmentRolloutStateEnum(item.(interface{}))
	}

	return items
}

// flattenOSPolicyAssignmentRolloutStateEnumSlice flattens the contents of OSPolicyAssignmentRolloutStateEnum from a JSON
// response object.
func flattenOSPolicyAssignmentRolloutStateEnumSlice(c *Client, i interface{}) []OSPolicyAssignmentRolloutStateEnum {
	a, ok := i.([]interface{})
	if !ok {
		return []OSPolicyAssignmentRolloutStateEnum{}
	}

	if len(a) == 0 {
		return []OSPolicyAssignmentRolloutStateEnum{}
	}

	items := make([]OSPolicyAssignmentRolloutStateEnum, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenOSPolicyAssignmentRolloutStateEnum(item.(interface{})))
	}

	return items
}

// flattenOSPolicyAssignmentRolloutStateEnum asserts that an interface is a string, and returns a
// pointer to a *OSPolicyAssignmentRolloutStateEnum with the same value as that string.
func flattenOSPolicyAssignmentRolloutStateEnum(i interface{}) *OSPolicyAssignmentRolloutStateEnum {
	s, ok := i.(string)
	if !ok {
		return OSPolicyAssignmentRolloutStateEnumRef("")
	}

	return OSPolicyAssignmentRolloutStateEnumRef(s)
}

// This function returns a matcher that checks whether a serialized resource matches this resource
// in its parameters (as defined by the fields in a Get, which definitionally define resource
// identity).  This is useful in extracting the element from a List call.
func (r *OSPolicyAssignment) matcher(c *Client) func([]byte) bool {
	return func(b []byte) bool {
		cr, err := unmarshalOSPolicyAssignment(b, c)
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

type oSPolicyAssignmentDiff struct {
	// The diff should include one or the other of RequiresRecreate or UpdateOp.
	RequiresRecreate bool
	UpdateOp         oSPolicyAssignmentApiOperation
}

func convertFieldDiffsToOSPolicyAssignmentDiffs(config *dcl.Config, fds []*dcl.FieldDiff, opts []dcl.ApplyOption) ([]oSPolicyAssignmentDiff, error) {
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
	var diffs []oSPolicyAssignmentDiff
	// For each operation name, create a oSPolicyAssignmentDiff which contains the operation.
	for opName, fieldDiffs := range opNamesToFieldDiffs {
		diff := oSPolicyAssignmentDiff{}
		if opName == "Recreate" {
			diff.RequiresRecreate = true
		} else {
			apiOp, err := convertOpNameToOSPolicyAssignmentApiOperation(opName, fieldDiffs, opts...)
			if err != nil {
				return diffs, err
			}
			diff.UpdateOp = apiOp
		}
		diffs = append(diffs, diff)
	}
	return diffs, nil
}

func convertOpNameToOSPolicyAssignmentApiOperation(opName string, fieldDiffs []*dcl.FieldDiff, opts ...dcl.ApplyOption) (oSPolicyAssignmentApiOperation, error) {
	switch opName {

	case "updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation":
		return &updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation{FieldDiffs: fieldDiffs}, nil

	default:
		return nil, fmt.Errorf("no such operation with name: %v", opName)
	}
}

func extractOSPolicyAssignmentFields(r *OSPolicyAssignment) error {
	vInstanceFilter := r.InstanceFilter
	if vInstanceFilter == nil {
		// note: explicitly not the empty object.
		vInstanceFilter = &OSPolicyAssignmentInstanceFilter{}
	}
	if err := extractOSPolicyAssignmentInstanceFilterFields(r, vInstanceFilter); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vInstanceFilter) {
		r.InstanceFilter = vInstanceFilter
	}
	vRollout := r.Rollout
	if vRollout == nil {
		// note: explicitly not the empty object.
		vRollout = &OSPolicyAssignmentRollout{}
	}
	if err := extractOSPolicyAssignmentRolloutFields(r, vRollout); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vRollout) {
		r.Rollout = vRollout
	}
	return nil
}
func extractOSPolicyAssignmentOSPoliciesFields(r *OSPolicyAssignment, o *OSPolicyAssignmentOSPolicies) error {
	return nil
}
func extractOSPolicyAssignmentOSPoliciesResourceGroupsFields(r *OSPolicyAssignment, o *OSPolicyAssignmentOSPoliciesResourceGroups) error {
	return nil
}
func extractOSPolicyAssignmentOSPoliciesResourceGroupsInventoryFiltersFields(r *OSPolicyAssignment, o *OSPolicyAssignmentOSPoliciesResourceGroupsInventoryFilters) error {
	return nil
}
func extractOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFields(r *OSPolicyAssignment, o *OSPolicyAssignmentOSPoliciesResourceGroupsResources) error {
	vPkg := o.Pkg
	if vPkg == nil {
		// note: explicitly not the empty object.
		vPkg = &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkg{}
	}
	if err := extractOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgFields(r, vPkg); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vPkg) {
		o.Pkg = vPkg
	}
	vRepository := o.Repository
	if vRepository == nil {
		// note: explicitly not the empty object.
		vRepository = &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepository{}
	}
	if err := extractOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryFields(r, vRepository); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vRepository) {
		o.Repository = vRepository
	}
	vExec := o.Exec
	if vExec == nil {
		// note: explicitly not the empty object.
		vExec = &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExec{}
	}
	if err := extractOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecFields(r, vExec); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vExec) {
		o.Exec = vExec
	}
	vFile := o.File
	if vFile == nil {
		// note: explicitly not the empty object.
		vFile = &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFile{}
	}
	if err := extractOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFields(r, vFile); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vFile) {
		o.File = vFile
	}
	return nil
}
func extractOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgFields(r *OSPolicyAssignment, o *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkg) error {
	vApt := o.Apt
	if vApt == nil {
		// note: explicitly not the empty object.
		vApt = &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgApt{}
	}
	if err := extractOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgAptFields(r, vApt); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vApt) {
		o.Apt = vApt
	}
	vDeb := o.Deb
	if vDeb == nil {
		// note: explicitly not the empty object.
		vDeb = &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDeb{}
	}
	if err := extractOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebFields(r, vDeb); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vDeb) {
		o.Deb = vDeb
	}
	vYum := o.Yum
	if vYum == nil {
		// note: explicitly not the empty object.
		vYum = &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgYum{}
	}
	if err := extractOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgYumFields(r, vYum); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vYum) {
		o.Yum = vYum
	}
	vZypper := o.Zypper
	if vZypper == nil {
		// note: explicitly not the empty object.
		vZypper = &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgZypper{}
	}
	if err := extractOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgZypperFields(r, vZypper); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vZypper) {
		o.Zypper = vZypper
	}
	vRpm := o.Rpm
	if vRpm == nil {
		// note: explicitly not the empty object.
		vRpm = &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpm{}
	}
	if err := extractOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmFields(r, vRpm); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vRpm) {
		o.Rpm = vRpm
	}
	vGooget := o.Googet
	if vGooget == nil {
		// note: explicitly not the empty object.
		vGooget = &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgGooget{}
	}
	if err := extractOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgGoogetFields(r, vGooget); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vGooget) {
		o.Googet = vGooget
	}
	vMsi := o.Msi
	if vMsi == nil {
		// note: explicitly not the empty object.
		vMsi = &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsi{}
	}
	if err := extractOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiFields(r, vMsi); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vMsi) {
		o.Msi = vMsi
	}
	return nil
}
func extractOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgAptFields(r *OSPolicyAssignment, o *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgApt) error {
	return nil
}
func extractOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebFields(r *OSPolicyAssignment, o *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDeb) error {
	// *OSPolicyAssignmentFile is a reused type - that's not compatible with function extractors.

	return nil
}
func extractOSPolicyAssignmentFileFields(r *OSPolicyAssignment, o *OSPolicyAssignmentFile) error {
	vRemote := o.Remote
	if vRemote == nil {
		// note: explicitly not the empty object.
		vRemote = &OSPolicyAssignmentFileRemote{}
	}
	if err := extractOSPolicyAssignmentFileRemoteFields(r, vRemote); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vRemote) {
		o.Remote = vRemote
	}
	vGcs := o.Gcs
	if vGcs == nil {
		// note: explicitly not the empty object.
		vGcs = &OSPolicyAssignmentFileGcs{}
	}
	if err := extractOSPolicyAssignmentFileGcsFields(r, vGcs); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vGcs) {
		o.Gcs = vGcs
	}
	return nil
}
func extractOSPolicyAssignmentFileRemoteFields(r *OSPolicyAssignment, o *OSPolicyAssignmentFileRemote) error {
	return nil
}
func extractOSPolicyAssignmentFileGcsFields(r *OSPolicyAssignment, o *OSPolicyAssignmentFileGcs) error {
	return nil
}
func extractOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgYumFields(r *OSPolicyAssignment, o *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgYum) error {
	return nil
}
func extractOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgZypperFields(r *OSPolicyAssignment, o *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgZypper) error {
	return nil
}
func extractOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmFields(r *OSPolicyAssignment, o *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpm) error {
	// *OSPolicyAssignmentFile is a reused type - that's not compatible with function extractors.

	return nil
}
func extractOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgGoogetFields(r *OSPolicyAssignment, o *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgGooget) error {
	return nil
}
func extractOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiFields(r *OSPolicyAssignment, o *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsi) error {
	// *OSPolicyAssignmentFile is a reused type - that's not compatible with function extractors.

	return nil
}
func extractOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryFields(r *OSPolicyAssignment, o *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepository) error {
	vApt := o.Apt
	if vApt == nil {
		// note: explicitly not the empty object.
		vApt = &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryApt{}
	}
	if err := extractOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryAptFields(r, vApt); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vApt) {
		o.Apt = vApt
	}
	vYum := o.Yum
	if vYum == nil {
		// note: explicitly not the empty object.
		vYum = &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryYum{}
	}
	if err := extractOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryYumFields(r, vYum); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vYum) {
		o.Yum = vYum
	}
	vZypper := o.Zypper
	if vZypper == nil {
		// note: explicitly not the empty object.
		vZypper = &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryZypper{}
	}
	if err := extractOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryZypperFields(r, vZypper); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vZypper) {
		o.Zypper = vZypper
	}
	vGoo := o.Goo
	if vGoo == nil {
		// note: explicitly not the empty object.
		vGoo = &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryGoo{}
	}
	if err := extractOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryGooFields(r, vGoo); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vGoo) {
		o.Goo = vGoo
	}
	return nil
}
func extractOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryAptFields(r *OSPolicyAssignment, o *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryApt) error {
	return nil
}
func extractOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryYumFields(r *OSPolicyAssignment, o *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryYum) error {
	return nil
}
func extractOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryZypperFields(r *OSPolicyAssignment, o *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryZypper) error {
	return nil
}
func extractOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryGooFields(r *OSPolicyAssignment, o *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryGoo) error {
	return nil
}
func extractOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecFields(r *OSPolicyAssignment, o *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExec) error {
	// *OSPolicyAssignmentExec is a reused type - that's not compatible with function extractors.

	// *OSPolicyAssignmentExec is a reused type - that's not compatible with function extractors.

	return nil
}
func extractOSPolicyAssignmentExecFields(r *OSPolicyAssignment, o *OSPolicyAssignmentExec) error {
	// *OSPolicyAssignmentFile is a reused type - that's not compatible with function extractors.

	return nil
}
func extractOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFields(r *OSPolicyAssignment, o *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFile) error {
	// *OSPolicyAssignmentFile is a reused type - that's not compatible with function extractors.

	return nil
}
func extractOSPolicyAssignmentInstanceFilterFields(r *OSPolicyAssignment, o *OSPolicyAssignmentInstanceFilter) error {
	return nil
}
func extractOSPolicyAssignmentInstanceFilterInclusionLabelsFields(r *OSPolicyAssignment, o *OSPolicyAssignmentInstanceFilterInclusionLabels) error {
	return nil
}
func extractOSPolicyAssignmentInstanceFilterExclusionLabelsFields(r *OSPolicyAssignment, o *OSPolicyAssignmentInstanceFilterExclusionLabels) error {
	return nil
}
func extractOSPolicyAssignmentInstanceFilterInventoriesFields(r *OSPolicyAssignment, o *OSPolicyAssignmentInstanceFilterInventories) error {
	return nil
}
func extractOSPolicyAssignmentRolloutFields(r *OSPolicyAssignment, o *OSPolicyAssignmentRollout) error {
	vDisruptionBudget := o.DisruptionBudget
	if vDisruptionBudget == nil {
		// note: explicitly not the empty object.
		vDisruptionBudget = &OSPolicyAssignmentRolloutDisruptionBudget{}
	}
	if err := extractOSPolicyAssignmentRolloutDisruptionBudgetFields(r, vDisruptionBudget); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vDisruptionBudget) {
		o.DisruptionBudget = vDisruptionBudget
	}
	return nil
}
func extractOSPolicyAssignmentRolloutDisruptionBudgetFields(r *OSPolicyAssignment, o *OSPolicyAssignmentRolloutDisruptionBudget) error {
	return nil
}

func postReadExtractOSPolicyAssignmentFields(r *OSPolicyAssignment) error {
	vInstanceFilter := r.InstanceFilter
	if vInstanceFilter == nil {
		// note: explicitly not the empty object.
		vInstanceFilter = &OSPolicyAssignmentInstanceFilter{}
	}
	if err := postReadExtractOSPolicyAssignmentInstanceFilterFields(r, vInstanceFilter); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vInstanceFilter) {
		r.InstanceFilter = vInstanceFilter
	}
	vRollout := r.Rollout
	if vRollout == nil {
		// note: explicitly not the empty object.
		vRollout = &OSPolicyAssignmentRollout{}
	}
	if err := postReadExtractOSPolicyAssignmentRolloutFields(r, vRollout); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vRollout) {
		r.Rollout = vRollout
	}
	return nil
}
func postReadExtractOSPolicyAssignmentOSPoliciesFields(r *OSPolicyAssignment, o *OSPolicyAssignmentOSPolicies) error {
	return nil
}
func postReadExtractOSPolicyAssignmentOSPoliciesResourceGroupsFields(r *OSPolicyAssignment, o *OSPolicyAssignmentOSPoliciesResourceGroups) error {
	return nil
}
func postReadExtractOSPolicyAssignmentOSPoliciesResourceGroupsInventoryFiltersFields(r *OSPolicyAssignment, o *OSPolicyAssignmentOSPoliciesResourceGroupsInventoryFilters) error {
	return nil
}
func postReadExtractOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFields(r *OSPolicyAssignment, o *OSPolicyAssignmentOSPoliciesResourceGroupsResources) error {
	vPkg := o.Pkg
	if vPkg == nil {
		// note: explicitly not the empty object.
		vPkg = &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkg{}
	}
	if err := extractOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgFields(r, vPkg); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vPkg) {
		o.Pkg = vPkg
	}
	vRepository := o.Repository
	if vRepository == nil {
		// note: explicitly not the empty object.
		vRepository = &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepository{}
	}
	if err := extractOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryFields(r, vRepository); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vRepository) {
		o.Repository = vRepository
	}
	vExec := o.Exec
	if vExec == nil {
		// note: explicitly not the empty object.
		vExec = &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExec{}
	}
	if err := extractOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecFields(r, vExec); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vExec) {
		o.Exec = vExec
	}
	vFile := o.File
	if vFile == nil {
		// note: explicitly not the empty object.
		vFile = &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFile{}
	}
	if err := extractOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFields(r, vFile); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vFile) {
		o.File = vFile
	}
	return nil
}
func postReadExtractOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgFields(r *OSPolicyAssignment, o *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkg) error {
	vApt := o.Apt
	if vApt == nil {
		// note: explicitly not the empty object.
		vApt = &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgApt{}
	}
	if err := extractOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgAptFields(r, vApt); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vApt) {
		o.Apt = vApt
	}
	vDeb := o.Deb
	if vDeb == nil {
		// note: explicitly not the empty object.
		vDeb = &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDeb{}
	}
	if err := extractOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebFields(r, vDeb); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vDeb) {
		o.Deb = vDeb
	}
	vYum := o.Yum
	if vYum == nil {
		// note: explicitly not the empty object.
		vYum = &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgYum{}
	}
	if err := extractOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgYumFields(r, vYum); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vYum) {
		o.Yum = vYum
	}
	vZypper := o.Zypper
	if vZypper == nil {
		// note: explicitly not the empty object.
		vZypper = &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgZypper{}
	}
	if err := extractOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgZypperFields(r, vZypper); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vZypper) {
		o.Zypper = vZypper
	}
	vRpm := o.Rpm
	if vRpm == nil {
		// note: explicitly not the empty object.
		vRpm = &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpm{}
	}
	if err := extractOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmFields(r, vRpm); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vRpm) {
		o.Rpm = vRpm
	}
	vGooget := o.Googet
	if vGooget == nil {
		// note: explicitly not the empty object.
		vGooget = &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgGooget{}
	}
	if err := extractOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgGoogetFields(r, vGooget); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vGooget) {
		o.Googet = vGooget
	}
	vMsi := o.Msi
	if vMsi == nil {
		// note: explicitly not the empty object.
		vMsi = &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsi{}
	}
	if err := extractOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiFields(r, vMsi); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vMsi) {
		o.Msi = vMsi
	}
	return nil
}
func postReadExtractOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgAptFields(r *OSPolicyAssignment, o *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgApt) error {
	return nil
}
func postReadExtractOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebFields(r *OSPolicyAssignment, o *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDeb) error {
	// *OSPolicyAssignmentFile is a reused type - that's not compatible with function extractors.

	return nil
}
func postReadExtractOSPolicyAssignmentFileFields(r *OSPolicyAssignment, o *OSPolicyAssignmentFile) error {
	vRemote := o.Remote
	if vRemote == nil {
		// note: explicitly not the empty object.
		vRemote = &OSPolicyAssignmentFileRemote{}
	}
	if err := extractOSPolicyAssignmentFileRemoteFields(r, vRemote); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vRemote) {
		o.Remote = vRemote
	}
	vGcs := o.Gcs
	if vGcs == nil {
		// note: explicitly not the empty object.
		vGcs = &OSPolicyAssignmentFileGcs{}
	}
	if err := extractOSPolicyAssignmentFileGcsFields(r, vGcs); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vGcs) {
		o.Gcs = vGcs
	}
	return nil
}
func postReadExtractOSPolicyAssignmentFileRemoteFields(r *OSPolicyAssignment, o *OSPolicyAssignmentFileRemote) error {
	return nil
}
func postReadExtractOSPolicyAssignmentFileGcsFields(r *OSPolicyAssignment, o *OSPolicyAssignmentFileGcs) error {
	return nil
}
func postReadExtractOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgYumFields(r *OSPolicyAssignment, o *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgYum) error {
	return nil
}
func postReadExtractOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgZypperFields(r *OSPolicyAssignment, o *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgZypper) error {
	return nil
}
func postReadExtractOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmFields(r *OSPolicyAssignment, o *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpm) error {
	// *OSPolicyAssignmentFile is a reused type - that's not compatible with function extractors.

	return nil
}
func postReadExtractOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgGoogetFields(r *OSPolicyAssignment, o *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgGooget) error {
	return nil
}
func postReadExtractOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiFields(r *OSPolicyAssignment, o *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsi) error {
	// *OSPolicyAssignmentFile is a reused type - that's not compatible with function extractors.

	return nil
}
func postReadExtractOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryFields(r *OSPolicyAssignment, o *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepository) error {
	vApt := o.Apt
	if vApt == nil {
		// note: explicitly not the empty object.
		vApt = &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryApt{}
	}
	if err := extractOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryAptFields(r, vApt); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vApt) {
		o.Apt = vApt
	}
	vYum := o.Yum
	if vYum == nil {
		// note: explicitly not the empty object.
		vYum = &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryYum{}
	}
	if err := extractOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryYumFields(r, vYum); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vYum) {
		o.Yum = vYum
	}
	vZypper := o.Zypper
	if vZypper == nil {
		// note: explicitly not the empty object.
		vZypper = &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryZypper{}
	}
	if err := extractOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryZypperFields(r, vZypper); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vZypper) {
		o.Zypper = vZypper
	}
	vGoo := o.Goo
	if vGoo == nil {
		// note: explicitly not the empty object.
		vGoo = &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryGoo{}
	}
	if err := extractOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryGooFields(r, vGoo); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vGoo) {
		o.Goo = vGoo
	}
	return nil
}
func postReadExtractOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryAptFields(r *OSPolicyAssignment, o *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryApt) error {
	return nil
}
func postReadExtractOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryYumFields(r *OSPolicyAssignment, o *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryYum) error {
	return nil
}
func postReadExtractOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryZypperFields(r *OSPolicyAssignment, o *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryZypper) error {
	return nil
}
func postReadExtractOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryGooFields(r *OSPolicyAssignment, o *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryGoo) error {
	return nil
}
func postReadExtractOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecFields(r *OSPolicyAssignment, o *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExec) error {
	// *OSPolicyAssignmentExec is a reused type - that's not compatible with function extractors.

	// *OSPolicyAssignmentExec is a reused type - that's not compatible with function extractors.

	return nil
}
func postReadExtractOSPolicyAssignmentExecFields(r *OSPolicyAssignment, o *OSPolicyAssignmentExec) error {
	// *OSPolicyAssignmentFile is a reused type - that's not compatible with function extractors.

	return nil
}
func postReadExtractOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFields(r *OSPolicyAssignment, o *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFile) error {
	// *OSPolicyAssignmentFile is a reused type - that's not compatible with function extractors.

	return nil
}
func postReadExtractOSPolicyAssignmentInstanceFilterFields(r *OSPolicyAssignment, o *OSPolicyAssignmentInstanceFilter) error {
	return nil
}
func postReadExtractOSPolicyAssignmentInstanceFilterInclusionLabelsFields(r *OSPolicyAssignment, o *OSPolicyAssignmentInstanceFilterInclusionLabels) error {
	return nil
}
func postReadExtractOSPolicyAssignmentInstanceFilterExclusionLabelsFields(r *OSPolicyAssignment, o *OSPolicyAssignmentInstanceFilterExclusionLabels) error {
	return nil
}
func postReadExtractOSPolicyAssignmentInstanceFilterInventoriesFields(r *OSPolicyAssignment, o *OSPolicyAssignmentInstanceFilterInventories) error {
	return nil
}
func postReadExtractOSPolicyAssignmentRolloutFields(r *OSPolicyAssignment, o *OSPolicyAssignmentRollout) error {
	vDisruptionBudget := o.DisruptionBudget
	if vDisruptionBudget == nil {
		// note: explicitly not the empty object.
		vDisruptionBudget = &OSPolicyAssignmentRolloutDisruptionBudget{}
	}
	if err := extractOSPolicyAssignmentRolloutDisruptionBudgetFields(r, vDisruptionBudget); err != nil {
		return err
	}
	if !dcl.IsNotReturnedByServer(vDisruptionBudget) {
		o.DisruptionBudget = vDisruptionBudget
	}
	return nil
}
func postReadExtractOSPolicyAssignmentRolloutDisruptionBudgetFields(r *OSPolicyAssignment, o *OSPolicyAssignmentRolloutDisruptionBudget) error {
	return nil
}
