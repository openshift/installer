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
func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSource) validate() error {
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
func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceRemote) validate() error {
	if err := dcl.Required(r, "uri"); err != nil {
		return err
	}
	return nil
}
func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceGcs) validate() error {
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
func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSource) validate() error {
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
func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceRemote) validate() error {
	if err := dcl.Required(r, "uri"); err != nil {
		return err
	}
	return nil
}
func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceGcs) validate() error {
	if err := dcl.Required(r, "bucket"); err != nil {
		return err
	}
	if err := dcl.Required(r, "object"); err != nil {
		return err
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
func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSource) validate() error {
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
func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceRemote) validate() error {
	if err := dcl.Required(r, "uri"); err != nil {
		return err
	}
	return nil
}
func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceGcs) validate() error {
	if err := dcl.Required(r, "bucket"); err != nil {
		return err
	}
	if err := dcl.Required(r, "object"); err != nil {
		return err
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
func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidate) validate() error {
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
func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFile) validate() error {
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
func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileRemote) validate() error {
	if err := dcl.Required(r, "uri"); err != nil {
		return err
	}
	return nil
}
func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileGcs) validate() error {
	if err := dcl.Required(r, "bucket"); err != nil {
		return err
	}
	if err := dcl.Required(r, "object"); err != nil {
		return err
	}
	return nil
}
func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforce) validate() error {
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
func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFile) validate() error {
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
func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileRemote) validate() error {
	if err := dcl.Required(r, "uri"); err != nil {
		return err
	}
	return nil
}
func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileGcs) validate() error {
	if err := dcl.Required(r, "bucket"); err != nil {
		return err
	}
	if err := dcl.Required(r, "object"); err != nil {
		return err
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
func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFile) validate() error {
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
func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileRemote) validate() error {
	if err := dcl.Required(r, "uri"); err != nil {
		return err
	}
	return nil
}
func (r *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileGcs) validate() error {
	if err := dcl.Required(r, "bucket"); err != nil {
		return err
	}
	if err := dcl.Required(r, "object"); err != nil {
		return err
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
	res := f
	_ = res

	if v := f.Description; !dcl.IsEmptyValueIndirect(v) {
		req["description"] = v
	}
	if v, err := expandOSPolicyAssignmentOSPoliciesSlice(c, f.OSPolicies, res); err != nil {
		return nil, fmt.Errorf("error expanding OSPolicies into osPolicies: %w", err)
	} else if v != nil {
		req["osPolicies"] = v
	}
	if v, err := expandOSPolicyAssignmentInstanceFilter(c, f.InstanceFilter, res); err != nil {
		return nil, fmt.Errorf("error expanding InstanceFilter into instanceFilter: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		req["instanceFilter"] = v
	}
	if v, err := expandOSPolicyAssignmentRollout(c, f.Rollout, res); err != nil {
		return nil, fmt.Errorf("error expanding Rollout into rollout: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		req["rollout"] = v
	}
	if v := f.SkipAwaitRollout; !dcl.IsEmptyValueIndirect(v) {
		req["skipAwaitRollout"] = v
	}
	return req, nil
}

// marshalUpdateOSPolicyAssignmentUpdateOSPolicyAssignmentRequest converts the update into
// the final JSON request body.
func marshalUpdateOSPolicyAssignmentUpdateOSPolicyAssignmentRequest(c *Client, m map[string]interface{}) ([]byte, error) {

	dcl.MoveMapEntry(
		m,
		[]string{"skipAwaitRollout"},
		[]string{},
	)
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
		res, err := unmarshalMapOSPolicyAssignment(v, c, r)
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

// Create operations are similar to Update operations, although they do not have
// specific request objects. The Create request object is the json encoding of
// the resource, which is modified by res.marshal to form the base request body.
type createOSPolicyAssignmentOperation struct {
	response map[string]interface{}
}

func (op *createOSPolicyAssignmentOperation) FirstResponse() (map[string]interface{}, bool) {
	return op.response, len(op.response) > 0
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

	// The Get call applies postReadExtract and so the result may contain fields that are not part of API version.
	if err := extractOSPolicyAssignmentFields(rawInitial); err != nil {
		return nil, nil, nil, err
	}

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
	if dcl.BoolCanonicalize(rawDesired.SkipAwaitRollout, rawInitial.SkipAwaitRollout) {
		canonicalDesired.SkipAwaitRollout = rawInitial.SkipAwaitRollout
	} else {
		canonicalDesired.SkipAwaitRollout = rawDesired.SkipAwaitRollout
	}
	return canonicalDesired, nil
}

func canonicalizeOSPolicyAssignmentNewState(c *Client, rawNew, rawDesired *OSPolicyAssignment) (*OSPolicyAssignment, error) {

	if dcl.IsEmptyValueIndirect(rawNew.Name) && dcl.IsEmptyValueIndirect(rawDesired.Name) {
		rawNew.Name = rawDesired.Name
	} else {
		if dcl.PartialSelfLinkToSelfLink(rawDesired.Name, rawNew.Name) {
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

	if dcl.IsEmptyValueIndirect(rawNew.OSPolicies) && dcl.IsEmptyValueIndirect(rawDesired.OSPolicies) {
		rawNew.OSPolicies = rawDesired.OSPolicies
	} else {
		rawNew.OSPolicies = canonicalizeNewOSPolicyAssignmentOSPoliciesSlice(c, rawDesired.OSPolicies, rawNew.OSPolicies)
	}

	if dcl.IsEmptyValueIndirect(rawNew.InstanceFilter) && dcl.IsEmptyValueIndirect(rawDesired.InstanceFilter) {
		rawNew.InstanceFilter = rawDesired.InstanceFilter
	} else {
		rawNew.InstanceFilter = canonicalizeNewOSPolicyAssignmentInstanceFilter(c, rawDesired.InstanceFilter, rawNew.InstanceFilter)
	}

	if dcl.IsEmptyValueIndirect(rawNew.Rollout) && dcl.IsEmptyValueIndirect(rawDesired.Rollout) {
		rawNew.Rollout = rawDesired.Rollout
	} else {
		rawNew.Rollout = canonicalizeNewOSPolicyAssignmentRollout(c, rawDesired.Rollout, rawNew.Rollout)
	}

	if dcl.IsEmptyValueIndirect(rawNew.RevisionId) && dcl.IsEmptyValueIndirect(rawDesired.RevisionId) {
		rawNew.RevisionId = rawDesired.RevisionId
	} else {
		if dcl.StringCanonicalize(rawDesired.RevisionId, rawNew.RevisionId) {
			rawNew.RevisionId = rawDesired.RevisionId
		}
	}

	if dcl.IsEmptyValueIndirect(rawNew.RevisionCreateTime) && dcl.IsEmptyValueIndirect(rawDesired.RevisionCreateTime) {
		rawNew.RevisionCreateTime = rawDesired.RevisionCreateTime
	} else {
	}

	if dcl.IsEmptyValueIndirect(rawNew.Etag) && dcl.IsEmptyValueIndirect(rawDesired.Etag) {
		rawNew.Etag = rawDesired.Etag
	} else {
		if dcl.StringCanonicalize(rawDesired.Etag, rawNew.Etag) {
			rawNew.Etag = rawDesired.Etag
		}
	}

	if dcl.IsEmptyValueIndirect(rawNew.RolloutState) && dcl.IsEmptyValueIndirect(rawDesired.RolloutState) {
		rawNew.RolloutState = rawDesired.RolloutState
	} else {
	}

	if dcl.IsEmptyValueIndirect(rawNew.Baseline) && dcl.IsEmptyValueIndirect(rawDesired.Baseline) {
		rawNew.Baseline = rawDesired.Baseline
	} else {
		if dcl.BoolCanonicalize(rawDesired.Baseline, rawNew.Baseline) {
			rawNew.Baseline = rawDesired.Baseline
		}
	}

	if dcl.IsEmptyValueIndirect(rawNew.Deleted) && dcl.IsEmptyValueIndirect(rawDesired.Deleted) {
		rawNew.Deleted = rawDesired.Deleted
	} else {
		if dcl.BoolCanonicalize(rawDesired.Deleted, rawNew.Deleted) {
			rawNew.Deleted = rawDesired.Deleted
		}
	}

	if dcl.IsEmptyValueIndirect(rawNew.Reconciling) && dcl.IsEmptyValueIndirect(rawDesired.Reconciling) {
		rawNew.Reconciling = rawDesired.Reconciling
	} else {
		if dcl.BoolCanonicalize(rawDesired.Reconciling, rawNew.Reconciling) {
			rawNew.Reconciling = rawDesired.Reconciling
		}
	}

	if dcl.IsEmptyValueIndirect(rawNew.Uid) && dcl.IsEmptyValueIndirect(rawDesired.Uid) {
		rawNew.Uid = rawDesired.Uid
	} else {
		if dcl.StringCanonicalize(rawDesired.Uid, rawNew.Uid) {
			rawNew.Uid = rawDesired.Uid
		}
	}

	rawNew.Project = rawDesired.Project

	rawNew.Location = rawDesired.Location

	if dcl.IsEmptyValueIndirect(rawNew.SkipAwaitRollout) && dcl.IsEmptyValueIndirect(rawDesired.SkipAwaitRollout) {
		rawNew.SkipAwaitRollout = rawDesired.SkipAwaitRollout
	} else {
		rawNew.SkipAwaitRollout = rawDesired.SkipAwaitRollout
	}

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
	if dcl.IsZeroValue(des.Mode) || (dcl.IsEmptyValueIndirect(des.Mode) && dcl.IsEmptyValueIndirect(initial.Mode)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
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
		if dcl.IsEmptyValueIndirect(des) {
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

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []OSPolicyAssignmentOSPolicies
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareOSPolicyAssignmentOSPoliciesNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewOSPolicyAssignmentOSPolicies(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
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
		if dcl.IsEmptyValueIndirect(des) {
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

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []OSPolicyAssignmentOSPoliciesResourceGroups
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareOSPolicyAssignmentOSPoliciesResourceGroupsNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroups(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
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
		if dcl.IsEmptyValueIndirect(des) {
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

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []OSPolicyAssignmentOSPoliciesResourceGroupsInventoryFilters
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareOSPolicyAssignmentOSPoliciesResourceGroupsInventoryFiltersNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsInventoryFilters(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
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
		if dcl.IsEmptyValueIndirect(des) {
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

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []OSPolicyAssignmentOSPoliciesResourceGroupsResources
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareOSPolicyAssignmentOSPoliciesResourceGroupsResourcesNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResources(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
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

	if dcl.IsZeroValue(des.DesiredState) || (dcl.IsEmptyValueIndirect(des.DesiredState) && dcl.IsEmptyValueIndirect(initial.DesiredState)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
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
	if dcl.IsEmptyValueIndirect(des) {
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
		if dcl.IsEmptyValueIndirect(des) {
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

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkg
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkg(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
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
	if dcl.IsEmptyValueIndirect(des) {
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
		if dcl.IsEmptyValueIndirect(des) {
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

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgApt
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgAptNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgApt(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
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

	cDes.Source = canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSource(des.Source, initial.Source, opts...)
	if dcl.BoolCanonicalize(des.PullDeps, initial.PullDeps) || dcl.IsZeroValue(des.PullDeps) {
		cDes.PullDeps = initial.PullDeps
	} else {
		cDes.PullDeps = des.PullDeps
	}

	return cDes
}

func canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSlice(des, initial []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDeb, opts ...dcl.ApplyOption) []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDeb {
	if dcl.IsEmptyValueIndirect(des) {
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
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDeb while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	nw.Source = canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSource(c, des.Source, nw.Source)
	if dcl.BoolCanonicalize(des.PullDeps, nw.PullDeps) {
		nw.PullDeps = des.PullDeps
	}

	return nw
}

func canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSet(c *Client, des, nw []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDeb) []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDeb {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDeb
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDeb(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
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

func canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSource(des, initial *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSource, opts ...dcl.ApplyOption) *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSource {
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

	cDes := &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSource{}

	cDes.Remote = canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceRemote(des.Remote, initial.Remote, opts...)
	cDes.Gcs = canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceGcs(des.Gcs, initial.Gcs, opts...)
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

func canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceSlice(des, initial []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSource, opts ...dcl.ApplyOption) []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSource {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSource, 0, len(des))
		for _, d := range des {
			cd := canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSource(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSource, 0, len(des))
	for i, d := range des {
		cd := canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSource(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSource(c *Client, des, nw *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSource) *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSource {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSource while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	nw.Remote = canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceRemote(c, des.Remote, nw.Remote)
	nw.Gcs = canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceGcs(c, des.Gcs, nw.Gcs)
	if dcl.StringCanonicalize(des.LocalPath, nw.LocalPath) {
		nw.LocalPath = des.LocalPath
	}
	if dcl.BoolCanonicalize(des.AllowInsecure, nw.AllowInsecure) {
		nw.AllowInsecure = des.AllowInsecure
	}

	return nw
}

func canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceSet(c *Client, des, nw []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSource) []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSource {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSource
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSource(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceSlice(c *Client, des, nw []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSource) []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSource {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSource
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSource(c, &d, &n))
	}

	return items
}

func canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceRemote(des, initial *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceRemote, opts ...dcl.ApplyOption) *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceRemote {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceRemote{}

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

func canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceRemoteSlice(des, initial []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceRemote, opts ...dcl.ApplyOption) []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceRemote {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceRemote, 0, len(des))
		for _, d := range des {
			cd := canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceRemote(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceRemote, 0, len(des))
	for i, d := range des {
		cd := canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceRemote(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceRemote(c *Client, des, nw *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceRemote) *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceRemote {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceRemote while comparing non-nil desired to nil actual.  Returning desired object.")
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

func canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceRemoteSet(c *Client, des, nw []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceRemote) []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceRemote {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceRemote
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceRemoteNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceRemote(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceRemoteSlice(c *Client, des, nw []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceRemote) []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceRemote {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceRemote
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceRemote(c, &d, &n))
	}

	return items
}

func canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceGcs(des, initial *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceGcs, opts ...dcl.ApplyOption) *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceGcs {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceGcs{}

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
	if dcl.IsZeroValue(des.Generation) || (dcl.IsEmptyValueIndirect(des.Generation) && dcl.IsEmptyValueIndirect(initial.Generation)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.Generation = initial.Generation
	} else {
		cDes.Generation = des.Generation
	}

	return cDes
}

func canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceGcsSlice(des, initial []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceGcs, opts ...dcl.ApplyOption) []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceGcs {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceGcs, 0, len(des))
		for _, d := range des {
			cd := canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceGcs(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceGcs, 0, len(des))
	for i, d := range des {
		cd := canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceGcs(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceGcs(c *Client, des, nw *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceGcs) *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceGcs {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceGcs while comparing non-nil desired to nil actual.  Returning desired object.")
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

func canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceGcsSet(c *Client, des, nw []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceGcs) []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceGcs {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceGcs
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceGcsNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceGcs(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceGcsSlice(c *Client, des, nw []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceGcs) []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceGcs {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceGcs
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceGcs(c, &d, &n))
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
	if dcl.IsEmptyValueIndirect(des) {
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
		if dcl.IsEmptyValueIndirect(des) {
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

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgYum
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgYumNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgYum(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
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
	if dcl.IsEmptyValueIndirect(des) {
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
		if dcl.IsEmptyValueIndirect(des) {
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

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgZypper
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgZypperNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgZypper(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
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

	cDes.Source = canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSource(des.Source, initial.Source, opts...)
	if dcl.BoolCanonicalize(des.PullDeps, initial.PullDeps) || dcl.IsZeroValue(des.PullDeps) {
		cDes.PullDeps = initial.PullDeps
	} else {
		cDes.PullDeps = des.PullDeps
	}

	return cDes
}

func canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSlice(des, initial []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpm, opts ...dcl.ApplyOption) []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpm {
	if dcl.IsEmptyValueIndirect(des) {
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
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpm while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	nw.Source = canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSource(c, des.Source, nw.Source)
	if dcl.BoolCanonicalize(des.PullDeps, nw.PullDeps) {
		nw.PullDeps = des.PullDeps
	}

	return nw
}

func canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSet(c *Client, des, nw []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpm) []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpm {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpm
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpm(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
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

func canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSource(des, initial *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSource, opts ...dcl.ApplyOption) *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSource {
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

	cDes := &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSource{}

	cDes.Remote = canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceRemote(des.Remote, initial.Remote, opts...)
	cDes.Gcs = canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceGcs(des.Gcs, initial.Gcs, opts...)
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

func canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceSlice(des, initial []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSource, opts ...dcl.ApplyOption) []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSource {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSource, 0, len(des))
		for _, d := range des {
			cd := canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSource(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSource, 0, len(des))
	for i, d := range des {
		cd := canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSource(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSource(c *Client, des, nw *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSource) *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSource {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSource while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	nw.Remote = canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceRemote(c, des.Remote, nw.Remote)
	nw.Gcs = canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceGcs(c, des.Gcs, nw.Gcs)
	if dcl.StringCanonicalize(des.LocalPath, nw.LocalPath) {
		nw.LocalPath = des.LocalPath
	}
	if dcl.BoolCanonicalize(des.AllowInsecure, nw.AllowInsecure) {
		nw.AllowInsecure = des.AllowInsecure
	}

	return nw
}

func canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceSet(c *Client, des, nw []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSource) []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSource {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSource
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSource(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceSlice(c *Client, des, nw []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSource) []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSource {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSource
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSource(c, &d, &n))
	}

	return items
}

func canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceRemote(des, initial *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceRemote, opts ...dcl.ApplyOption) *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceRemote {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceRemote{}

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

func canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceRemoteSlice(des, initial []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceRemote, opts ...dcl.ApplyOption) []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceRemote {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceRemote, 0, len(des))
		for _, d := range des {
			cd := canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceRemote(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceRemote, 0, len(des))
	for i, d := range des {
		cd := canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceRemote(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceRemote(c *Client, des, nw *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceRemote) *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceRemote {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceRemote while comparing non-nil desired to nil actual.  Returning desired object.")
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

func canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceRemoteSet(c *Client, des, nw []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceRemote) []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceRemote {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceRemote
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceRemoteNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceRemote(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceRemoteSlice(c *Client, des, nw []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceRemote) []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceRemote {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceRemote
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceRemote(c, &d, &n))
	}

	return items
}

func canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceGcs(des, initial *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceGcs, opts ...dcl.ApplyOption) *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceGcs {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceGcs{}

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
	if dcl.IsZeroValue(des.Generation) || (dcl.IsEmptyValueIndirect(des.Generation) && dcl.IsEmptyValueIndirect(initial.Generation)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.Generation = initial.Generation
	} else {
		cDes.Generation = des.Generation
	}

	return cDes
}

func canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceGcsSlice(des, initial []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceGcs, opts ...dcl.ApplyOption) []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceGcs {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceGcs, 0, len(des))
		for _, d := range des {
			cd := canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceGcs(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceGcs, 0, len(des))
	for i, d := range des {
		cd := canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceGcs(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceGcs(c *Client, des, nw *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceGcs) *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceGcs {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceGcs while comparing non-nil desired to nil actual.  Returning desired object.")
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

func canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceGcsSet(c *Client, des, nw []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceGcs) []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceGcs {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceGcs
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceGcsNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceGcs(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceGcsSlice(c *Client, des, nw []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceGcs) []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceGcs {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceGcs
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceGcs(c, &d, &n))
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
	if dcl.IsEmptyValueIndirect(des) {
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
		if dcl.IsEmptyValueIndirect(des) {
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

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgGooget
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgGoogetNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgGooget(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
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

	cDes.Source = canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSource(des.Source, initial.Source, opts...)
	if dcl.StringArrayCanonicalize(des.Properties, initial.Properties) {
		cDes.Properties = initial.Properties
	} else {
		cDes.Properties = des.Properties
	}

	return cDes
}

func canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSlice(des, initial []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsi, opts ...dcl.ApplyOption) []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsi {
	if dcl.IsEmptyValueIndirect(des) {
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
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsi while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	nw.Source = canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSource(c, des.Source, nw.Source)
	if dcl.StringArrayCanonicalize(des.Properties, nw.Properties) {
		nw.Properties = des.Properties
	}

	return nw
}

func canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSet(c *Client, des, nw []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsi) []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsi {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsi
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsi(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
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

func canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSource(des, initial *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSource, opts ...dcl.ApplyOption) *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSource {
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

	cDes := &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSource{}

	cDes.Remote = canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceRemote(des.Remote, initial.Remote, opts...)
	cDes.Gcs = canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceGcs(des.Gcs, initial.Gcs, opts...)
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

func canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceSlice(des, initial []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSource, opts ...dcl.ApplyOption) []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSource {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSource, 0, len(des))
		for _, d := range des {
			cd := canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSource(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSource, 0, len(des))
	for i, d := range des {
		cd := canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSource(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSource(c *Client, des, nw *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSource) *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSource {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSource while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	nw.Remote = canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceRemote(c, des.Remote, nw.Remote)
	nw.Gcs = canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceGcs(c, des.Gcs, nw.Gcs)
	if dcl.StringCanonicalize(des.LocalPath, nw.LocalPath) {
		nw.LocalPath = des.LocalPath
	}
	if dcl.BoolCanonicalize(des.AllowInsecure, nw.AllowInsecure) {
		nw.AllowInsecure = des.AllowInsecure
	}

	return nw
}

func canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceSet(c *Client, des, nw []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSource) []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSource {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSource
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSource(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceSlice(c *Client, des, nw []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSource) []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSource {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSource
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSource(c, &d, &n))
	}

	return items
}

func canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceRemote(des, initial *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceRemote, opts ...dcl.ApplyOption) *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceRemote {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceRemote{}

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

func canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceRemoteSlice(des, initial []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceRemote, opts ...dcl.ApplyOption) []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceRemote {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceRemote, 0, len(des))
		for _, d := range des {
			cd := canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceRemote(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceRemote, 0, len(des))
	for i, d := range des {
		cd := canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceRemote(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceRemote(c *Client, des, nw *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceRemote) *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceRemote {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceRemote while comparing non-nil desired to nil actual.  Returning desired object.")
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

func canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceRemoteSet(c *Client, des, nw []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceRemote) []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceRemote {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceRemote
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceRemoteNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceRemote(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceRemoteSlice(c *Client, des, nw []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceRemote) []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceRemote {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceRemote
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceRemote(c, &d, &n))
	}

	return items
}

func canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceGcs(des, initial *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceGcs, opts ...dcl.ApplyOption) *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceGcs {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceGcs{}

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
	if dcl.IsZeroValue(des.Generation) || (dcl.IsEmptyValueIndirect(des.Generation) && dcl.IsEmptyValueIndirect(initial.Generation)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.Generation = initial.Generation
	} else {
		cDes.Generation = des.Generation
	}

	return cDes
}

func canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceGcsSlice(des, initial []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceGcs, opts ...dcl.ApplyOption) []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceGcs {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceGcs, 0, len(des))
		for _, d := range des {
			cd := canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceGcs(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceGcs, 0, len(des))
	for i, d := range des {
		cd := canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceGcs(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceGcs(c *Client, des, nw *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceGcs) *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceGcs {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceGcs while comparing non-nil desired to nil actual.  Returning desired object.")
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

func canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceGcsSet(c *Client, des, nw []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceGcs) []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceGcs {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceGcs
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceGcsNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceGcs(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceGcsSlice(c *Client, des, nw []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceGcs) []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceGcs {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceGcs
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceGcs(c, &d, &n))
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
	if dcl.IsEmptyValueIndirect(des) {
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
		if dcl.IsEmptyValueIndirect(des) {
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

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepository
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepository(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
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

	if dcl.IsZeroValue(des.ArchiveType) || (dcl.IsEmptyValueIndirect(des.ArchiveType) && dcl.IsEmptyValueIndirect(initial.ArchiveType)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
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
	if dcl.StringArrayCanonicalize(des.Components, initial.Components) {
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
	if dcl.IsEmptyValueIndirect(des) {
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
		if dcl.IsEmptyValueIndirect(des) {
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

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryApt
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryAptNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryApt(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
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
	if dcl.StringArrayCanonicalize(des.GpgKeys, initial.GpgKeys) {
		cDes.GpgKeys = initial.GpgKeys
	} else {
		cDes.GpgKeys = des.GpgKeys
	}

	return cDes
}

func canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryYumSlice(des, initial []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryYum, opts ...dcl.ApplyOption) []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryYum {
	if dcl.IsEmptyValueIndirect(des) {
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
		if dcl.IsEmptyValueIndirect(des) {
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

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryYum
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryYumNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryYum(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
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
	if dcl.StringArrayCanonicalize(des.GpgKeys, initial.GpgKeys) {
		cDes.GpgKeys = initial.GpgKeys
	} else {
		cDes.GpgKeys = des.GpgKeys
	}

	return cDes
}

func canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryZypperSlice(des, initial []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryZypper, opts ...dcl.ApplyOption) []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryZypper {
	if dcl.IsEmptyValueIndirect(des) {
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
		if dcl.IsEmptyValueIndirect(des) {
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

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryZypper
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryZypperNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryZypper(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
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
	if dcl.IsEmptyValueIndirect(des) {
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
		if dcl.IsEmptyValueIndirect(des) {
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

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryGoo
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryGooNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryGoo(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
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

	cDes.Validate = canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidate(des.Validate, initial.Validate, opts...)
	cDes.Enforce = canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforce(des.Enforce, initial.Enforce, opts...)

	return cDes
}

func canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecSlice(des, initial []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExec, opts ...dcl.ApplyOption) []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExec {
	if dcl.IsEmptyValueIndirect(des) {
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
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExec while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	nw.Validate = canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidate(c, des.Validate, nw.Validate)
	nw.Enforce = canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforce(c, des.Enforce, nw.Enforce)

	return nw
}

func canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecSet(c *Client, des, nw []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExec) []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExec {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExec
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExec(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
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

func canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidate(des, initial *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidate, opts ...dcl.ApplyOption) *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidate {
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

	cDes := &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidate{}

	cDes.File = canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFile(des.File, initial.File, opts...)
	if dcl.StringCanonicalize(des.Script, initial.Script) || dcl.IsZeroValue(des.Script) {
		cDes.Script = initial.Script
	} else {
		cDes.Script = des.Script
	}
	if dcl.StringArrayCanonicalize(des.Args, initial.Args) {
		cDes.Args = initial.Args
	} else {
		cDes.Args = des.Args
	}
	if dcl.IsZeroValue(des.Interpreter) || (dcl.IsEmptyValueIndirect(des.Interpreter) && dcl.IsEmptyValueIndirect(initial.Interpreter)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
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

func canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateSlice(des, initial []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidate, opts ...dcl.ApplyOption) []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidate {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidate, 0, len(des))
		for _, d := range des {
			cd := canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidate(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidate, 0, len(des))
	for i, d := range des {
		cd := canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidate(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidate(c *Client, des, nw *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidate) *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidate {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidate while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	nw.File = canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFile(c, des.File, nw.File)
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

func canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateSet(c *Client, des, nw []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidate) []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidate {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidate
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidate(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateSlice(c *Client, des, nw []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidate) []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidate {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidate
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidate(c, &d, &n))
	}

	return items
}

func canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFile(des, initial *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFile, opts ...dcl.ApplyOption) *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFile {
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

	cDes := &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFile{}

	cDes.Remote = canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileRemote(des.Remote, initial.Remote, opts...)
	cDes.Gcs = canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileGcs(des.Gcs, initial.Gcs, opts...)
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

func canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileSlice(des, initial []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFile, opts ...dcl.ApplyOption) []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFile {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFile, 0, len(des))
		for _, d := range des {
			cd := canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFile(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFile, 0, len(des))
	for i, d := range des {
		cd := canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFile(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFile(c *Client, des, nw *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFile) *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFile {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFile while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	nw.Remote = canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileRemote(c, des.Remote, nw.Remote)
	nw.Gcs = canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileGcs(c, des.Gcs, nw.Gcs)
	if dcl.StringCanonicalize(des.LocalPath, nw.LocalPath) {
		nw.LocalPath = des.LocalPath
	}
	if dcl.BoolCanonicalize(des.AllowInsecure, nw.AllowInsecure) {
		nw.AllowInsecure = des.AllowInsecure
	}

	return nw
}

func canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileSet(c *Client, des, nw []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFile) []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFile {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFile
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFile(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileSlice(c *Client, des, nw []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFile) []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFile {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFile
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFile(c, &d, &n))
	}

	return items
}

func canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileRemote(des, initial *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileRemote, opts ...dcl.ApplyOption) *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileRemote {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileRemote{}

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

func canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileRemoteSlice(des, initial []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileRemote, opts ...dcl.ApplyOption) []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileRemote {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileRemote, 0, len(des))
		for _, d := range des {
			cd := canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileRemote(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileRemote, 0, len(des))
	for i, d := range des {
		cd := canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileRemote(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileRemote(c *Client, des, nw *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileRemote) *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileRemote {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileRemote while comparing non-nil desired to nil actual.  Returning desired object.")
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

func canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileRemoteSet(c *Client, des, nw []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileRemote) []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileRemote {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileRemote
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileRemoteNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileRemote(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileRemoteSlice(c *Client, des, nw []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileRemote) []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileRemote {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileRemote
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileRemote(c, &d, &n))
	}

	return items
}

func canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileGcs(des, initial *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileGcs, opts ...dcl.ApplyOption) *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileGcs {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileGcs{}

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
	if dcl.IsZeroValue(des.Generation) || (dcl.IsEmptyValueIndirect(des.Generation) && dcl.IsEmptyValueIndirect(initial.Generation)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.Generation = initial.Generation
	} else {
		cDes.Generation = des.Generation
	}

	return cDes
}

func canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileGcsSlice(des, initial []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileGcs, opts ...dcl.ApplyOption) []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileGcs {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileGcs, 0, len(des))
		for _, d := range des {
			cd := canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileGcs(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileGcs, 0, len(des))
	for i, d := range des {
		cd := canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileGcs(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileGcs(c *Client, des, nw *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileGcs) *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileGcs {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileGcs while comparing non-nil desired to nil actual.  Returning desired object.")
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

func canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileGcsSet(c *Client, des, nw []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileGcs) []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileGcs {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileGcs
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileGcsNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileGcs(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileGcsSlice(c *Client, des, nw []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileGcs) []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileGcs {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileGcs
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileGcs(c, &d, &n))
	}

	return items
}

func canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforce(des, initial *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforce, opts ...dcl.ApplyOption) *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforce {
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

	cDes := &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforce{}

	cDes.File = canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFile(des.File, initial.File, opts...)
	if dcl.StringCanonicalize(des.Script, initial.Script) || dcl.IsZeroValue(des.Script) {
		cDes.Script = initial.Script
	} else {
		cDes.Script = des.Script
	}
	if dcl.StringArrayCanonicalize(des.Args, initial.Args) {
		cDes.Args = initial.Args
	} else {
		cDes.Args = des.Args
	}
	if dcl.IsZeroValue(des.Interpreter) || (dcl.IsEmptyValueIndirect(des.Interpreter) && dcl.IsEmptyValueIndirect(initial.Interpreter)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
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

func canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceSlice(des, initial []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforce, opts ...dcl.ApplyOption) []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforce {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforce, 0, len(des))
		for _, d := range des {
			cd := canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforce(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforce, 0, len(des))
	for i, d := range des {
		cd := canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforce(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforce(c *Client, des, nw *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforce) *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforce {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforce while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	nw.File = canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFile(c, des.File, nw.File)
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

func canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceSet(c *Client, des, nw []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforce) []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforce {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforce
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforce(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceSlice(c *Client, des, nw []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforce) []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforce {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforce
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforce(c, &d, &n))
	}

	return items
}

func canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFile(des, initial *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFile, opts ...dcl.ApplyOption) *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFile {
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

	cDes := &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFile{}

	cDes.Remote = canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileRemote(des.Remote, initial.Remote, opts...)
	cDes.Gcs = canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileGcs(des.Gcs, initial.Gcs, opts...)
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

func canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileSlice(des, initial []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFile, opts ...dcl.ApplyOption) []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFile {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFile, 0, len(des))
		for _, d := range des {
			cd := canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFile(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFile, 0, len(des))
	for i, d := range des {
		cd := canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFile(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFile(c *Client, des, nw *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFile) *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFile {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFile while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	nw.Remote = canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileRemote(c, des.Remote, nw.Remote)
	nw.Gcs = canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileGcs(c, des.Gcs, nw.Gcs)
	if dcl.StringCanonicalize(des.LocalPath, nw.LocalPath) {
		nw.LocalPath = des.LocalPath
	}
	if dcl.BoolCanonicalize(des.AllowInsecure, nw.AllowInsecure) {
		nw.AllowInsecure = des.AllowInsecure
	}

	return nw
}

func canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileSet(c *Client, des, nw []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFile) []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFile {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFile
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFile(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileSlice(c *Client, des, nw []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFile) []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFile {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFile
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFile(c, &d, &n))
	}

	return items
}

func canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileRemote(des, initial *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileRemote, opts ...dcl.ApplyOption) *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileRemote {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileRemote{}

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

func canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileRemoteSlice(des, initial []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileRemote, opts ...dcl.ApplyOption) []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileRemote {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileRemote, 0, len(des))
		for _, d := range des {
			cd := canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileRemote(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileRemote, 0, len(des))
	for i, d := range des {
		cd := canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileRemote(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileRemote(c *Client, des, nw *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileRemote) *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileRemote {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileRemote while comparing non-nil desired to nil actual.  Returning desired object.")
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

func canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileRemoteSet(c *Client, des, nw []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileRemote) []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileRemote {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileRemote
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileRemoteNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileRemote(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileRemoteSlice(c *Client, des, nw []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileRemote) []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileRemote {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileRemote
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileRemote(c, &d, &n))
	}

	return items
}

func canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileGcs(des, initial *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileGcs, opts ...dcl.ApplyOption) *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileGcs {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileGcs{}

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
	if dcl.IsZeroValue(des.Generation) || (dcl.IsEmptyValueIndirect(des.Generation) && dcl.IsEmptyValueIndirect(initial.Generation)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.Generation = initial.Generation
	} else {
		cDes.Generation = des.Generation
	}

	return cDes
}

func canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileGcsSlice(des, initial []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileGcs, opts ...dcl.ApplyOption) []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileGcs {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileGcs, 0, len(des))
		for _, d := range des {
			cd := canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileGcs(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileGcs, 0, len(des))
	for i, d := range des {
		cd := canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileGcs(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileGcs(c *Client, des, nw *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileGcs) *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileGcs {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileGcs while comparing non-nil desired to nil actual.  Returning desired object.")
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

func canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileGcsSet(c *Client, des, nw []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileGcs) []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileGcs {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileGcs
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileGcsNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileGcs(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileGcsSlice(c *Client, des, nw []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileGcs) []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileGcs {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileGcs
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileGcs(c, &d, &n))
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

	cDes.File = canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFile(des.File, initial.File, opts...)
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
	if dcl.IsZeroValue(des.State) || (dcl.IsEmptyValueIndirect(des.State) && dcl.IsEmptyValueIndirect(initial.State)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.State = initial.State
	} else {
		cDes.State = des.State
	}

	return cDes
}

func canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileSlice(des, initial []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFile, opts ...dcl.ApplyOption) []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFile {
	if dcl.IsEmptyValueIndirect(des) {
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
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFile while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	nw.File = canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFile(c, des.File, nw.File)
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

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFile
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFile(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
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

func canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFile(des, initial *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFile, opts ...dcl.ApplyOption) *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFile {
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

	cDes := &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFile{}

	cDes.Remote = canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileRemote(des.Remote, initial.Remote, opts...)
	cDes.Gcs = canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileGcs(des.Gcs, initial.Gcs, opts...)
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

func canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileSlice(des, initial []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFile, opts ...dcl.ApplyOption) []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFile {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFile, 0, len(des))
		for _, d := range des {
			cd := canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFile(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFile, 0, len(des))
	for i, d := range des {
		cd := canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFile(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFile(c *Client, des, nw *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFile) *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFile {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFile while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	nw.Remote = canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileRemote(c, des.Remote, nw.Remote)
	nw.Gcs = canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileGcs(c, des.Gcs, nw.Gcs)
	if dcl.StringCanonicalize(des.LocalPath, nw.LocalPath) {
		nw.LocalPath = des.LocalPath
	}
	if dcl.BoolCanonicalize(des.AllowInsecure, nw.AllowInsecure) {
		nw.AllowInsecure = des.AllowInsecure
	}

	return nw
}

func canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileSet(c *Client, des, nw []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFile) []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFile {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFile
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFile(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileSlice(c *Client, des, nw []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFile) []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFile {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFile
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFile(c, &d, &n))
	}

	return items
}

func canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileRemote(des, initial *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileRemote, opts ...dcl.ApplyOption) *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileRemote {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileRemote{}

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

func canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileRemoteSlice(des, initial []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileRemote, opts ...dcl.ApplyOption) []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileRemote {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileRemote, 0, len(des))
		for _, d := range des {
			cd := canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileRemote(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileRemote, 0, len(des))
	for i, d := range des {
		cd := canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileRemote(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileRemote(c *Client, des, nw *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileRemote) *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileRemote {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileRemote while comparing non-nil desired to nil actual.  Returning desired object.")
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

func canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileRemoteSet(c *Client, des, nw []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileRemote) []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileRemote {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileRemote
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileRemoteNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileRemote(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileRemoteSlice(c *Client, des, nw []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileRemote) []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileRemote {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileRemote
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileRemote(c, &d, &n))
	}

	return items
}

func canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileGcs(des, initial *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileGcs, opts ...dcl.ApplyOption) *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileGcs {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileGcs{}

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
	if dcl.IsZeroValue(des.Generation) || (dcl.IsEmptyValueIndirect(des.Generation) && dcl.IsEmptyValueIndirect(initial.Generation)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.Generation = initial.Generation
	} else {
		cDes.Generation = des.Generation
	}

	return cDes
}

func canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileGcsSlice(des, initial []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileGcs, opts ...dcl.ApplyOption) []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileGcs {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileGcs, 0, len(des))
		for _, d := range des {
			cd := canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileGcs(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileGcs, 0, len(des))
	for i, d := range des {
		cd := canonicalizeOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileGcs(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileGcs(c *Client, des, nw *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileGcs) *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileGcs {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileGcs while comparing non-nil desired to nil actual.  Returning desired object.")
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

func canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileGcsSet(c *Client, des, nw []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileGcs) []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileGcs {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileGcs
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileGcsNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileGcs(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileGcsSlice(c *Client, des, nw []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileGcs) []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileGcs {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileGcs
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileGcs(c, &d, &n))
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
	if dcl.IsEmptyValueIndirect(des) {
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
		if dcl.IsEmptyValueIndirect(des) {
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

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []OSPolicyAssignmentInstanceFilter
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareOSPolicyAssignmentInstanceFilterNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewOSPolicyAssignmentInstanceFilter(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
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

	if dcl.IsZeroValue(des.Labels) || (dcl.IsEmptyValueIndirect(des.Labels) && dcl.IsEmptyValueIndirect(initial.Labels)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
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
		if dcl.IsEmptyValueIndirect(des) {
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

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []OSPolicyAssignmentInstanceFilterInclusionLabels
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareOSPolicyAssignmentInstanceFilterInclusionLabelsNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewOSPolicyAssignmentInstanceFilterInclusionLabels(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
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

	if dcl.IsZeroValue(des.Labels) || (dcl.IsEmptyValueIndirect(des.Labels) && dcl.IsEmptyValueIndirect(initial.Labels)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
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
		if dcl.IsEmptyValueIndirect(des) {
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

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []OSPolicyAssignmentInstanceFilterExclusionLabels
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareOSPolicyAssignmentInstanceFilterExclusionLabelsNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewOSPolicyAssignmentInstanceFilterExclusionLabels(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
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
		if dcl.IsEmptyValueIndirect(des) {
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

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []OSPolicyAssignmentInstanceFilterInventories
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareOSPolicyAssignmentInstanceFilterInventoriesNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewOSPolicyAssignmentInstanceFilterInventories(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
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
	if dcl.IsEmptyValueIndirect(des) {
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
		if dcl.IsEmptyValueIndirect(des) {
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

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []OSPolicyAssignmentRollout
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareOSPolicyAssignmentRolloutNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewOSPolicyAssignmentRollout(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
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

	if dcl.IsZeroValue(des.Fixed) || (dcl.IsEmptyValueIndirect(des.Fixed) && dcl.IsEmptyValueIndirect(initial.Fixed)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.Fixed = initial.Fixed
	} else {
		cDes.Fixed = des.Fixed
	}
	if dcl.IsZeroValue(des.Percent) || (dcl.IsEmptyValueIndirect(des.Percent) && dcl.IsEmptyValueIndirect(initial.Percent)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.Percent = initial.Percent
	} else {
		cDes.Percent = des.Percent
	}

	return cDes
}

func canonicalizeOSPolicyAssignmentRolloutDisruptionBudgetSlice(des, initial []OSPolicyAssignmentRolloutDisruptionBudget, opts ...dcl.ApplyOption) []OSPolicyAssignmentRolloutDisruptionBudget {
	if dcl.IsEmptyValueIndirect(des) {
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
		if dcl.IsEmptyValueIndirect(des) {
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

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []OSPolicyAssignmentRolloutDisruptionBudget
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareOSPolicyAssignmentRolloutDisruptionBudgetNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewOSPolicyAssignmentRolloutDisruptionBudget(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
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
	if ds, err := dcl.Diff(desired.Name, actual.Name, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Name")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Description, actual.Description, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("Description")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.OSPolicies, actual.OSPolicies, dcl.DiffInfo{ObjectFunction: compareOSPolicyAssignmentOSPoliciesNewStyle, EmptyObject: EmptyOSPolicyAssignmentOSPolicies, OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("OsPolicies")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.InstanceFilter, actual.InstanceFilter, dcl.DiffInfo{ObjectFunction: compareOSPolicyAssignmentInstanceFilterNewStyle, EmptyObject: EmptyOSPolicyAssignmentInstanceFilter, OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("InstanceFilter")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Rollout, actual.Rollout, dcl.DiffInfo{ObjectFunction: compareOSPolicyAssignmentRolloutNewStyle, EmptyObject: EmptyOSPolicyAssignmentRollout, OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("Rollout")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.RevisionId, actual.RevisionId, dcl.DiffInfo{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("RevisionId")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.RevisionCreateTime, actual.RevisionCreateTime, dcl.DiffInfo{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("RevisionCreateTime")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Etag, actual.Etag, dcl.DiffInfo{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Etag")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.RolloutState, actual.RolloutState, dcl.DiffInfo{OutputOnly: true, Type: "EnumType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("RolloutState")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Baseline, actual.Baseline, dcl.DiffInfo{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Baseline")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Deleted, actual.Deleted, dcl.DiffInfo{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Deleted")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Reconciling, actual.Reconciling, dcl.DiffInfo{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Reconciling")); len(ds) != 0 || err != nil {
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

	if ds, err := dcl.Diff(desired.SkipAwaitRollout, actual.SkipAwaitRollout, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("SkipAwaitRollout")); len(ds) != 0 || err != nil {
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

	if ds, err := dcl.Diff(desired.Id, actual.Id, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("Id")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Description, actual.Description, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("Description")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Mode, actual.Mode, dcl.DiffInfo{Type: "EnumType", OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("Mode")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.ResourceGroups, actual.ResourceGroups, dcl.DiffInfo{ObjectFunction: compareOSPolicyAssignmentOSPoliciesResourceGroupsNewStyle, EmptyObject: EmptyOSPolicyAssignmentOSPoliciesResourceGroups, OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("ResourceGroups")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.AllowNoResourceGroupMatch, actual.AllowNoResourceGroupMatch, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("AllowNoResourceGroupMatch")); len(ds) != 0 || err != nil {
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

	if ds, err := dcl.Diff(desired.InventoryFilters, actual.InventoryFilters, dcl.DiffInfo{ObjectFunction: compareOSPolicyAssignmentOSPoliciesResourceGroupsInventoryFiltersNewStyle, EmptyObject: EmptyOSPolicyAssignmentOSPoliciesResourceGroupsInventoryFilters, OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("InventoryFilters")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Resources, actual.Resources, dcl.DiffInfo{ObjectFunction: compareOSPolicyAssignmentOSPoliciesResourceGroupsResourcesNewStyle, EmptyObject: EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResources, OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("Resources")); len(ds) != 0 || err != nil {
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

	if ds, err := dcl.Diff(desired.OSShortName, actual.OSShortName, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("OsShortName")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.OSVersion, actual.OSVersion, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("OsVersion")); len(ds) != 0 || err != nil {
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

	if ds, err := dcl.Diff(desired.Id, actual.Id, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("Id")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Pkg, actual.Pkg, dcl.DiffInfo{ObjectFunction: compareOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgNewStyle, EmptyObject: EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkg, OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("Pkg")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Repository, actual.Repository, dcl.DiffInfo{ObjectFunction: compareOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryNewStyle, EmptyObject: EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepository, OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("Repository")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Exec, actual.Exec, dcl.DiffInfo{ObjectFunction: compareOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecNewStyle, EmptyObject: EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExec, OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("Exec")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.File, actual.File, dcl.DiffInfo{ObjectFunction: compareOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileNewStyle, EmptyObject: EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFile, OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("File")); len(ds) != 0 || err != nil {
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

	if ds, err := dcl.Diff(desired.DesiredState, actual.DesiredState, dcl.DiffInfo{Type: "EnumType", OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("DesiredState")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Apt, actual.Apt, dcl.DiffInfo{ObjectFunction: compareOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgAptNewStyle, EmptyObject: EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgApt, OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("Apt")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Deb, actual.Deb, dcl.DiffInfo{ObjectFunction: compareOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebNewStyle, EmptyObject: EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDeb, OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("Deb")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Yum, actual.Yum, dcl.DiffInfo{ObjectFunction: compareOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgYumNewStyle, EmptyObject: EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgYum, OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("Yum")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Zypper, actual.Zypper, dcl.DiffInfo{ObjectFunction: compareOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgZypperNewStyle, EmptyObject: EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgZypper, OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("Zypper")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Rpm, actual.Rpm, dcl.DiffInfo{ObjectFunction: compareOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmNewStyle, EmptyObject: EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpm, OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("Rpm")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Googet, actual.Googet, dcl.DiffInfo{ObjectFunction: compareOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgGoogetNewStyle, EmptyObject: EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgGooget, OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("Googet")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Msi, actual.Msi, dcl.DiffInfo{ObjectFunction: compareOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiNewStyle, EmptyObject: EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsi, OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("Msi")); len(ds) != 0 || err != nil {
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

	if ds, err := dcl.Diff(desired.Name, actual.Name, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("Name")); len(ds) != 0 || err != nil {
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

	if ds, err := dcl.Diff(desired.Source, actual.Source, dcl.DiffInfo{ObjectFunction: compareOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceNewStyle, EmptyObject: EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSource, OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("Source")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.PullDeps, actual.PullDeps, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("PullDeps")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSource)
	if !ok {
		desiredNotPointer, ok := d.(OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSource)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSource or *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSource", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSource)
	if !ok {
		actualNotPointer, ok := a.(OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSource)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSource", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.Remote, actual.Remote, dcl.DiffInfo{ObjectFunction: compareOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceRemoteNewStyle, EmptyObject: EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceRemote, OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("Remote")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Gcs, actual.Gcs, dcl.DiffInfo{ObjectFunction: compareOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceGcsNewStyle, EmptyObject: EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceGcs, OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("Gcs")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.LocalPath, actual.LocalPath, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("LocalPath")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.AllowInsecure, actual.AllowInsecure, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("AllowInsecure")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceRemoteNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceRemote)
	if !ok {
		desiredNotPointer, ok := d.(OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceRemote)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceRemote or *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceRemote", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceRemote)
	if !ok {
		actualNotPointer, ok := a.(OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceRemote)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceRemote", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.Uri, actual.Uri, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("Uri")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Sha256Checksum, actual.Sha256Checksum, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("Sha256Checksum")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceGcsNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceGcs)
	if !ok {
		desiredNotPointer, ok := d.(OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceGcs)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceGcs or *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceGcs", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceGcs)
	if !ok {
		actualNotPointer, ok := a.(OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceGcs)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceGcs", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.Bucket, actual.Bucket, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("Bucket")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Object, actual.Object, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("Object")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Generation, actual.Generation, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("Generation")); len(ds) != 0 || err != nil {
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

	if ds, err := dcl.Diff(desired.Name, actual.Name, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("Name")); len(ds) != 0 || err != nil {
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

	if ds, err := dcl.Diff(desired.Name, actual.Name, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("Name")); len(ds) != 0 || err != nil {
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

	if ds, err := dcl.Diff(desired.Source, actual.Source, dcl.DiffInfo{ObjectFunction: compareOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceNewStyle, EmptyObject: EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSource, OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("Source")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.PullDeps, actual.PullDeps, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("PullDeps")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSource)
	if !ok {
		desiredNotPointer, ok := d.(OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSource)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSource or *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSource", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSource)
	if !ok {
		actualNotPointer, ok := a.(OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSource)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSource", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.Remote, actual.Remote, dcl.DiffInfo{ObjectFunction: compareOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceRemoteNewStyle, EmptyObject: EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceRemote, OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("Remote")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Gcs, actual.Gcs, dcl.DiffInfo{ObjectFunction: compareOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceGcsNewStyle, EmptyObject: EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceGcs, OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("Gcs")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.LocalPath, actual.LocalPath, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("LocalPath")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.AllowInsecure, actual.AllowInsecure, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("AllowInsecure")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceRemoteNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceRemote)
	if !ok {
		desiredNotPointer, ok := d.(OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceRemote)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceRemote or *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceRemote", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceRemote)
	if !ok {
		actualNotPointer, ok := a.(OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceRemote)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceRemote", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.Uri, actual.Uri, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("Uri")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Sha256Checksum, actual.Sha256Checksum, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("Sha256Checksum")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceGcsNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceGcs)
	if !ok {
		desiredNotPointer, ok := d.(OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceGcs)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceGcs or *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceGcs", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceGcs)
	if !ok {
		actualNotPointer, ok := a.(OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceGcs)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceGcs", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.Bucket, actual.Bucket, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("Bucket")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Object, actual.Object, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("Object")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Generation, actual.Generation, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("Generation")); len(ds) != 0 || err != nil {
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

	if ds, err := dcl.Diff(desired.Name, actual.Name, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("Name")); len(ds) != 0 || err != nil {
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

	if ds, err := dcl.Diff(desired.Source, actual.Source, dcl.DiffInfo{ObjectFunction: compareOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceNewStyle, EmptyObject: EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSource, OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("Source")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Properties, actual.Properties, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("Properties")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSource)
	if !ok {
		desiredNotPointer, ok := d.(OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSource)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSource or *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSource", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSource)
	if !ok {
		actualNotPointer, ok := a.(OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSource)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSource", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.Remote, actual.Remote, dcl.DiffInfo{ObjectFunction: compareOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceRemoteNewStyle, EmptyObject: EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceRemote, OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("Remote")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Gcs, actual.Gcs, dcl.DiffInfo{ObjectFunction: compareOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceGcsNewStyle, EmptyObject: EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceGcs, OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("Gcs")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.LocalPath, actual.LocalPath, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("LocalPath")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.AllowInsecure, actual.AllowInsecure, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("AllowInsecure")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceRemoteNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceRemote)
	if !ok {
		desiredNotPointer, ok := d.(OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceRemote)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceRemote or *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceRemote", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceRemote)
	if !ok {
		actualNotPointer, ok := a.(OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceRemote)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceRemote", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.Uri, actual.Uri, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("Uri")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Sha256Checksum, actual.Sha256Checksum, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("Sha256Checksum")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceGcsNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceGcs)
	if !ok {
		desiredNotPointer, ok := d.(OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceGcs)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceGcs or *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceGcs", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceGcs)
	if !ok {
		actualNotPointer, ok := a.(OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceGcs)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceGcs", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.Bucket, actual.Bucket, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("Bucket")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Object, actual.Object, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("Object")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Generation, actual.Generation, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("Generation")); len(ds) != 0 || err != nil {
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

	if ds, err := dcl.Diff(desired.Apt, actual.Apt, dcl.DiffInfo{ObjectFunction: compareOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryAptNewStyle, EmptyObject: EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryApt, OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("Apt")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Yum, actual.Yum, dcl.DiffInfo{ObjectFunction: compareOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryYumNewStyle, EmptyObject: EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryYum, OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("Yum")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Zypper, actual.Zypper, dcl.DiffInfo{ObjectFunction: compareOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryZypperNewStyle, EmptyObject: EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryZypper, OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("Zypper")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Goo, actual.Goo, dcl.DiffInfo{ObjectFunction: compareOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryGooNewStyle, EmptyObject: EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryGoo, OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("Goo")); len(ds) != 0 || err != nil {
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

	if ds, err := dcl.Diff(desired.ArchiveType, actual.ArchiveType, dcl.DiffInfo{Type: "EnumType", OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("ArchiveType")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Uri, actual.Uri, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("Uri")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Distribution, actual.Distribution, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("Distribution")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Components, actual.Components, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("Components")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.GpgKey, actual.GpgKey, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("GpgKey")); len(ds) != 0 || err != nil {
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

	if ds, err := dcl.Diff(desired.Id, actual.Id, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("Id")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.DisplayName, actual.DisplayName, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("DisplayName")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.BaseUrl, actual.BaseUrl, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("BaseUrl")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.GpgKeys, actual.GpgKeys, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("GpgKeys")); len(ds) != 0 || err != nil {
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

	if ds, err := dcl.Diff(desired.Id, actual.Id, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("Id")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.DisplayName, actual.DisplayName, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("DisplayName")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.BaseUrl, actual.BaseUrl, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("BaseUrl")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.GpgKeys, actual.GpgKeys, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("GpgKeys")); len(ds) != 0 || err != nil {
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

	if ds, err := dcl.Diff(desired.Name, actual.Name, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("Name")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Url, actual.Url, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("Url")); len(ds) != 0 || err != nil {
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

	if ds, err := dcl.Diff(desired.Validate, actual.Validate, dcl.DiffInfo{ObjectFunction: compareOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateNewStyle, EmptyObject: EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidate, OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("Validate")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Enforce, actual.Enforce, dcl.DiffInfo{ObjectFunction: compareOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceNewStyle, EmptyObject: EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforce, OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("Enforce")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidate)
	if !ok {
		desiredNotPointer, ok := d.(OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidate)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidate or *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidate", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidate)
	if !ok {
		actualNotPointer, ok := a.(OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidate)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidate", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.File, actual.File, dcl.DiffInfo{ObjectFunction: compareOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileNewStyle, EmptyObject: EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFile, OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("File")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Script, actual.Script, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("Script")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Args, actual.Args, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("Args")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Interpreter, actual.Interpreter, dcl.DiffInfo{Type: "EnumType", OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("Interpreter")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.OutputFilePath, actual.OutputFilePath, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("OutputFilePath")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFile)
	if !ok {
		desiredNotPointer, ok := d.(OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFile)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFile or *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFile", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFile)
	if !ok {
		actualNotPointer, ok := a.(OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFile)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFile", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.Remote, actual.Remote, dcl.DiffInfo{ObjectFunction: compareOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileRemoteNewStyle, EmptyObject: EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileRemote, OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("Remote")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Gcs, actual.Gcs, dcl.DiffInfo{ObjectFunction: compareOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileGcsNewStyle, EmptyObject: EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileGcs, OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("Gcs")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.LocalPath, actual.LocalPath, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("LocalPath")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.AllowInsecure, actual.AllowInsecure, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("AllowInsecure")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileRemoteNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileRemote)
	if !ok {
		desiredNotPointer, ok := d.(OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileRemote)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileRemote or *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileRemote", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileRemote)
	if !ok {
		actualNotPointer, ok := a.(OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileRemote)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileRemote", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.Uri, actual.Uri, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("Uri")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Sha256Checksum, actual.Sha256Checksum, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("Sha256Checksum")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileGcsNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileGcs)
	if !ok {
		desiredNotPointer, ok := d.(OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileGcs)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileGcs or *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileGcs", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileGcs)
	if !ok {
		actualNotPointer, ok := a.(OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileGcs)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileGcs", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.Bucket, actual.Bucket, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("Bucket")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Object, actual.Object, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("Object")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Generation, actual.Generation, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("Generation")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforce)
	if !ok {
		desiredNotPointer, ok := d.(OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforce)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforce or *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforce", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforce)
	if !ok {
		actualNotPointer, ok := a.(OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforce)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforce", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.File, actual.File, dcl.DiffInfo{ObjectFunction: compareOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileNewStyle, EmptyObject: EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFile, OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("File")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Script, actual.Script, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("Script")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Args, actual.Args, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("Args")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Interpreter, actual.Interpreter, dcl.DiffInfo{Type: "EnumType", OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("Interpreter")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.OutputFilePath, actual.OutputFilePath, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("OutputFilePath")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFile)
	if !ok {
		desiredNotPointer, ok := d.(OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFile)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFile or *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFile", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFile)
	if !ok {
		actualNotPointer, ok := a.(OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFile)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFile", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.Remote, actual.Remote, dcl.DiffInfo{ObjectFunction: compareOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileRemoteNewStyle, EmptyObject: EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileRemote, OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("Remote")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Gcs, actual.Gcs, dcl.DiffInfo{ObjectFunction: compareOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileGcsNewStyle, EmptyObject: EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileGcs, OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("Gcs")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.LocalPath, actual.LocalPath, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("LocalPath")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.AllowInsecure, actual.AllowInsecure, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("AllowInsecure")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileRemoteNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileRemote)
	if !ok {
		desiredNotPointer, ok := d.(OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileRemote)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileRemote or *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileRemote", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileRemote)
	if !ok {
		actualNotPointer, ok := a.(OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileRemote)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileRemote", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.Uri, actual.Uri, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("Uri")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Sha256Checksum, actual.Sha256Checksum, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("Sha256Checksum")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileGcsNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileGcs)
	if !ok {
		desiredNotPointer, ok := d.(OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileGcs)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileGcs or *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileGcs", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileGcs)
	if !ok {
		actualNotPointer, ok := a.(OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileGcs)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileGcs", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.Bucket, actual.Bucket, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("Bucket")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Object, actual.Object, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("Object")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Generation, actual.Generation, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("Generation")); len(ds) != 0 || err != nil {
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

	if ds, err := dcl.Diff(desired.File, actual.File, dcl.DiffInfo{ObjectFunction: compareOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileNewStyle, EmptyObject: EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFile, OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("File")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Content, actual.Content, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("Content")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Path, actual.Path, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("Path")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.State, actual.State, dcl.DiffInfo{Type: "EnumType", OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("State")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Permissions, actual.Permissions, dcl.DiffInfo{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Permissions")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFile)
	if !ok {
		desiredNotPointer, ok := d.(OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFile)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFile or *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFile", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFile)
	if !ok {
		actualNotPointer, ok := a.(OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFile)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFile", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.Remote, actual.Remote, dcl.DiffInfo{ObjectFunction: compareOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileRemoteNewStyle, EmptyObject: EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileRemote, OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("Remote")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Gcs, actual.Gcs, dcl.DiffInfo{ObjectFunction: compareOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileGcsNewStyle, EmptyObject: EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileGcs, OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("Gcs")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.LocalPath, actual.LocalPath, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("LocalPath")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.AllowInsecure, actual.AllowInsecure, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("AllowInsecure")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileRemoteNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileRemote)
	if !ok {
		desiredNotPointer, ok := d.(OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileRemote)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileRemote or *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileRemote", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileRemote)
	if !ok {
		actualNotPointer, ok := a.(OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileRemote)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileRemote", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.Uri, actual.Uri, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("Uri")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Sha256Checksum, actual.Sha256Checksum, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("Sha256Checksum")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileGcsNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileGcs)
	if !ok {
		desiredNotPointer, ok := d.(OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileGcs)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileGcs or *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileGcs", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileGcs)
	if !ok {
		actualNotPointer, ok := a.(OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileGcs)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileGcs", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.Bucket, actual.Bucket, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("Bucket")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Object, actual.Object, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("Object")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Generation, actual.Generation, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("Generation")); len(ds) != 0 || err != nil {
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

	if ds, err := dcl.Diff(desired.All, actual.All, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("All")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.InclusionLabels, actual.InclusionLabels, dcl.DiffInfo{ObjectFunction: compareOSPolicyAssignmentInstanceFilterInclusionLabelsNewStyle, EmptyObject: EmptyOSPolicyAssignmentInstanceFilterInclusionLabels, OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("InclusionLabels")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.ExclusionLabels, actual.ExclusionLabels, dcl.DiffInfo{ObjectFunction: compareOSPolicyAssignmentInstanceFilterExclusionLabelsNewStyle, EmptyObject: EmptyOSPolicyAssignmentInstanceFilterExclusionLabels, OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("ExclusionLabels")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Inventories, actual.Inventories, dcl.DiffInfo{ObjectFunction: compareOSPolicyAssignmentInstanceFilterInventoriesNewStyle, EmptyObject: EmptyOSPolicyAssignmentInstanceFilterInventories, OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("Inventories")); len(ds) != 0 || err != nil {
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

	if ds, err := dcl.Diff(desired.Labels, actual.Labels, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("Labels")); len(ds) != 0 || err != nil {
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

	if ds, err := dcl.Diff(desired.Labels, actual.Labels, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("Labels")); len(ds) != 0 || err != nil {
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

	if ds, err := dcl.Diff(desired.OSShortName, actual.OSShortName, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("OsShortName")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.OSVersion, actual.OSVersion, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("OsVersion")); len(ds) != 0 || err != nil {
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

	if ds, err := dcl.Diff(desired.DisruptionBudget, actual.DisruptionBudget, dcl.DiffInfo{ObjectFunction: compareOSPolicyAssignmentRolloutDisruptionBudgetNewStyle, EmptyObject: EmptyOSPolicyAssignmentRolloutDisruptionBudget, OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("DisruptionBudget")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.MinWaitDuration, actual.MinWaitDuration, dcl.DiffInfo{CustomDiff: canonicalizeOSPolicyAssignmentRolloutMinWaitDuration, OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("MinWaitDuration")); len(ds) != 0 || err != nil {
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

	if ds, err := dcl.Diff(desired.Fixed, actual.Fixed, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("Fixed")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Percent, actual.Percent, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateOSPolicyAssignmentUpdateOSPolicyAssignmentOperation")}, fn.AddNest("Percent")); len(ds) != 0 || err != nil {
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
	dcl.MoveMapEntry(
		m,
		[]string{"skipAwaitRollout"},
		[]string{},
	)

	return json.Marshal(m)
}

// unmarshalOSPolicyAssignment decodes JSON responses into the OSPolicyAssignment resource schema.
func unmarshalOSPolicyAssignment(b []byte, c *Client, res *OSPolicyAssignment) (*OSPolicyAssignment, error) {
	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, err
	}
	return unmarshalMapOSPolicyAssignment(m, c, res)
}

func unmarshalMapOSPolicyAssignment(m map[string]interface{}, c *Client, res *OSPolicyAssignment) (*OSPolicyAssignment, error) {

	flattened := flattenOSPolicyAssignment(c, m, res)
	if flattened == nil {
		return nil, fmt.Errorf("attempted to flatten empty json object")
	}
	return flattened, nil
}

// expandOSPolicyAssignment expands OSPolicyAssignment into a JSON request object.
func expandOSPolicyAssignment(c *Client, f *OSPolicyAssignment) (map[string]interface{}, error) {
	m := make(map[string]interface{})
	res := f
	_ = res
	if v, err := dcl.DeriveField("projects/%s/locations/%s/osPolicyAssignments/%s", f.Name, dcl.SelfLinkToName(f.Project), dcl.SelfLinkToName(f.Location), dcl.SelfLinkToName(f.Name)); err != nil {
		return nil, fmt.Errorf("error expanding Name into name: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["name"] = v
	}
	if v := f.Description; dcl.ValueShouldBeSent(v) {
		m["description"] = v
	}
	if v, err := expandOSPolicyAssignmentOSPoliciesSlice(c, f.OSPolicies, res); err != nil {
		return nil, fmt.Errorf("error expanding OSPolicies into osPolicies: %w", err)
	} else if v != nil {
		m["osPolicies"] = v
	}
	if v, err := expandOSPolicyAssignmentInstanceFilter(c, f.InstanceFilter, res); err != nil {
		return nil, fmt.Errorf("error expanding InstanceFilter into instanceFilter: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["instanceFilter"] = v
	}
	if v, err := expandOSPolicyAssignmentRollout(c, f.Rollout, res); err != nil {
		return nil, fmt.Errorf("error expanding Rollout into rollout: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["rollout"] = v
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
	if v := f.SkipAwaitRollout; dcl.ValueShouldBeSent(v) {
		m["skipAwaitRollout"] = v
	}

	return m, nil
}

// flattenOSPolicyAssignment flattens OSPolicyAssignment from a JSON request object into the
// OSPolicyAssignment type.
func flattenOSPolicyAssignment(c *Client, i interface{}, res *OSPolicyAssignment) *OSPolicyAssignment {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}
	if len(m) == 0 {
		return nil
	}

	resultRes := &OSPolicyAssignment{}
	resultRes.Name = dcl.FlattenString(m["name"])
	resultRes.Description = dcl.FlattenString(m["description"])
	resultRes.OSPolicies = flattenOSPolicyAssignmentOSPoliciesSlice(c, m["osPolicies"], res)
	resultRes.InstanceFilter = flattenOSPolicyAssignmentInstanceFilter(c, m["instanceFilter"], res)
	resultRes.Rollout = flattenOSPolicyAssignmentRollout(c, m["rollout"], res)
	resultRes.RevisionId = dcl.FlattenString(m["revisionId"])
	resultRes.RevisionCreateTime = dcl.FlattenString(m["revisionCreateTime"])
	resultRes.Etag = dcl.FlattenString(m["etag"])
	resultRes.RolloutState = flattenOSPolicyAssignmentRolloutStateEnum(m["rolloutState"])
	resultRes.Baseline = dcl.FlattenBool(m["baseline"])
	resultRes.Deleted = dcl.FlattenBool(m["deleted"])
	resultRes.Reconciling = dcl.FlattenBool(m["reconciling"])
	resultRes.Uid = dcl.FlattenString(m["uid"])
	resultRes.Project = dcl.FlattenString(m["project"])
	resultRes.Location = dcl.FlattenString(m["location"])
	resultRes.SkipAwaitRollout = dcl.FlattenBool(m["skipAwaitRollout"])

	return resultRes
}

// expandOSPolicyAssignmentOSPoliciesMap expands the contents of OSPolicyAssignmentOSPolicies into a JSON
// request object.
func expandOSPolicyAssignmentOSPoliciesMap(c *Client, f map[string]OSPolicyAssignmentOSPolicies, res *OSPolicyAssignment) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandOSPolicyAssignmentOSPolicies(c, &item, res)
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
func expandOSPolicyAssignmentOSPoliciesSlice(c *Client, f []OSPolicyAssignmentOSPolicies, res *OSPolicyAssignment) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandOSPolicyAssignmentOSPolicies(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenOSPolicyAssignmentOSPoliciesMap flattens the contents of OSPolicyAssignmentOSPolicies from a JSON
// response object.
func flattenOSPolicyAssignmentOSPoliciesMap(c *Client, i interface{}, res *OSPolicyAssignment) map[string]OSPolicyAssignmentOSPolicies {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]OSPolicyAssignmentOSPolicies{}
	}

	if len(a) == 0 {
		return map[string]OSPolicyAssignmentOSPolicies{}
	}

	items := make(map[string]OSPolicyAssignmentOSPolicies)
	for k, item := range a {
		items[k] = *flattenOSPolicyAssignmentOSPolicies(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenOSPolicyAssignmentOSPoliciesSlice flattens the contents of OSPolicyAssignmentOSPolicies from a JSON
// response object.
func flattenOSPolicyAssignmentOSPoliciesSlice(c *Client, i interface{}, res *OSPolicyAssignment) []OSPolicyAssignmentOSPolicies {
	a, ok := i.([]interface{})
	if !ok {
		return []OSPolicyAssignmentOSPolicies{}
	}

	if len(a) == 0 {
		return []OSPolicyAssignmentOSPolicies{}
	}

	items := make([]OSPolicyAssignmentOSPolicies, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenOSPolicyAssignmentOSPolicies(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandOSPolicyAssignmentOSPolicies expands an instance of OSPolicyAssignmentOSPolicies into a JSON
// request object.
func expandOSPolicyAssignmentOSPolicies(c *Client, f *OSPolicyAssignmentOSPolicies, res *OSPolicyAssignment) (map[string]interface{}, error) {
	if f == nil {
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
	if v, err := expandOSPolicyAssignmentOSPoliciesResourceGroupsSlice(c, f.ResourceGroups, res); err != nil {
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
func flattenOSPolicyAssignmentOSPolicies(c *Client, i interface{}, res *OSPolicyAssignment) *OSPolicyAssignmentOSPolicies {
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
	r.ResourceGroups = flattenOSPolicyAssignmentOSPoliciesResourceGroupsSlice(c, m["resourceGroups"], res)
	r.AllowNoResourceGroupMatch = dcl.FlattenBool(m["allowNoResourceGroupMatch"])

	return r
}

// expandOSPolicyAssignmentOSPoliciesResourceGroupsMap expands the contents of OSPolicyAssignmentOSPoliciesResourceGroups into a JSON
// request object.
func expandOSPolicyAssignmentOSPoliciesResourceGroupsMap(c *Client, f map[string]OSPolicyAssignmentOSPoliciesResourceGroups, res *OSPolicyAssignment) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandOSPolicyAssignmentOSPoliciesResourceGroups(c, &item, res)
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
func expandOSPolicyAssignmentOSPoliciesResourceGroupsSlice(c *Client, f []OSPolicyAssignmentOSPoliciesResourceGroups, res *OSPolicyAssignment) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandOSPolicyAssignmentOSPoliciesResourceGroups(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenOSPolicyAssignmentOSPoliciesResourceGroupsMap flattens the contents of OSPolicyAssignmentOSPoliciesResourceGroups from a JSON
// response object.
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsMap(c *Client, i interface{}, res *OSPolicyAssignment) map[string]OSPolicyAssignmentOSPoliciesResourceGroups {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]OSPolicyAssignmentOSPoliciesResourceGroups{}
	}

	if len(a) == 0 {
		return map[string]OSPolicyAssignmentOSPoliciesResourceGroups{}
	}

	items := make(map[string]OSPolicyAssignmentOSPoliciesResourceGroups)
	for k, item := range a {
		items[k] = *flattenOSPolicyAssignmentOSPoliciesResourceGroups(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenOSPolicyAssignmentOSPoliciesResourceGroupsSlice flattens the contents of OSPolicyAssignmentOSPoliciesResourceGroups from a JSON
// response object.
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsSlice(c *Client, i interface{}, res *OSPolicyAssignment) []OSPolicyAssignmentOSPoliciesResourceGroups {
	a, ok := i.([]interface{})
	if !ok {
		return []OSPolicyAssignmentOSPoliciesResourceGroups{}
	}

	if len(a) == 0 {
		return []OSPolicyAssignmentOSPoliciesResourceGroups{}
	}

	items := make([]OSPolicyAssignmentOSPoliciesResourceGroups, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenOSPolicyAssignmentOSPoliciesResourceGroups(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandOSPolicyAssignmentOSPoliciesResourceGroups expands an instance of OSPolicyAssignmentOSPoliciesResourceGroups into a JSON
// request object.
func expandOSPolicyAssignmentOSPoliciesResourceGroups(c *Client, f *OSPolicyAssignmentOSPoliciesResourceGroups, res *OSPolicyAssignment) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v, err := expandOSPolicyAssignmentOSPoliciesResourceGroupsInventoryFiltersSlice(c, f.InventoryFilters, res); err != nil {
		return nil, fmt.Errorf("error expanding InventoryFilters into inventoryFilters: %w", err)
	} else if v != nil {
		m["inventoryFilters"] = v
	}
	if v, err := expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesSlice(c, f.Resources, res); err != nil {
		return nil, fmt.Errorf("error expanding Resources into resources: %w", err)
	} else if v != nil {
		m["resources"] = v
	}

	return m, nil
}

// flattenOSPolicyAssignmentOSPoliciesResourceGroups flattens an instance of OSPolicyAssignmentOSPoliciesResourceGroups from a JSON
// response object.
func flattenOSPolicyAssignmentOSPoliciesResourceGroups(c *Client, i interface{}, res *OSPolicyAssignment) *OSPolicyAssignmentOSPoliciesResourceGroups {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &OSPolicyAssignmentOSPoliciesResourceGroups{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyOSPolicyAssignmentOSPoliciesResourceGroups
	}
	r.InventoryFilters = flattenOSPolicyAssignmentOSPoliciesResourceGroupsInventoryFiltersSlice(c, m["inventoryFilters"], res)
	r.Resources = flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesSlice(c, m["resources"], res)

	return r
}

// expandOSPolicyAssignmentOSPoliciesResourceGroupsInventoryFiltersMap expands the contents of OSPolicyAssignmentOSPoliciesResourceGroupsInventoryFilters into a JSON
// request object.
func expandOSPolicyAssignmentOSPoliciesResourceGroupsInventoryFiltersMap(c *Client, f map[string]OSPolicyAssignmentOSPoliciesResourceGroupsInventoryFilters, res *OSPolicyAssignment) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandOSPolicyAssignmentOSPoliciesResourceGroupsInventoryFilters(c, &item, res)
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
func expandOSPolicyAssignmentOSPoliciesResourceGroupsInventoryFiltersSlice(c *Client, f []OSPolicyAssignmentOSPoliciesResourceGroupsInventoryFilters, res *OSPolicyAssignment) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandOSPolicyAssignmentOSPoliciesResourceGroupsInventoryFilters(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenOSPolicyAssignmentOSPoliciesResourceGroupsInventoryFiltersMap flattens the contents of OSPolicyAssignmentOSPoliciesResourceGroupsInventoryFilters from a JSON
// response object.
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsInventoryFiltersMap(c *Client, i interface{}, res *OSPolicyAssignment) map[string]OSPolicyAssignmentOSPoliciesResourceGroupsInventoryFilters {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]OSPolicyAssignmentOSPoliciesResourceGroupsInventoryFilters{}
	}

	if len(a) == 0 {
		return map[string]OSPolicyAssignmentOSPoliciesResourceGroupsInventoryFilters{}
	}

	items := make(map[string]OSPolicyAssignmentOSPoliciesResourceGroupsInventoryFilters)
	for k, item := range a {
		items[k] = *flattenOSPolicyAssignmentOSPoliciesResourceGroupsInventoryFilters(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenOSPolicyAssignmentOSPoliciesResourceGroupsInventoryFiltersSlice flattens the contents of OSPolicyAssignmentOSPoliciesResourceGroupsInventoryFilters from a JSON
// response object.
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsInventoryFiltersSlice(c *Client, i interface{}, res *OSPolicyAssignment) []OSPolicyAssignmentOSPoliciesResourceGroupsInventoryFilters {
	a, ok := i.([]interface{})
	if !ok {
		return []OSPolicyAssignmentOSPoliciesResourceGroupsInventoryFilters{}
	}

	if len(a) == 0 {
		return []OSPolicyAssignmentOSPoliciesResourceGroupsInventoryFilters{}
	}

	items := make([]OSPolicyAssignmentOSPoliciesResourceGroupsInventoryFilters, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenOSPolicyAssignmentOSPoliciesResourceGroupsInventoryFilters(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandOSPolicyAssignmentOSPoliciesResourceGroupsInventoryFilters expands an instance of OSPolicyAssignmentOSPoliciesResourceGroupsInventoryFilters into a JSON
// request object.
func expandOSPolicyAssignmentOSPoliciesResourceGroupsInventoryFilters(c *Client, f *OSPolicyAssignmentOSPoliciesResourceGroupsInventoryFilters, res *OSPolicyAssignment) (map[string]interface{}, error) {
	if f == nil {
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
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsInventoryFilters(c *Client, i interface{}, res *OSPolicyAssignment) *OSPolicyAssignmentOSPoliciesResourceGroupsInventoryFilters {
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
func expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesMap(c *Client, f map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResources, res *OSPolicyAssignment) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandOSPolicyAssignmentOSPoliciesResourceGroupsResources(c, &item, res)
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
func expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesSlice(c *Client, f []OSPolicyAssignmentOSPoliciesResourceGroupsResources, res *OSPolicyAssignment) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandOSPolicyAssignmentOSPoliciesResourceGroupsResources(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesMap flattens the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResources from a JSON
// response object.
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesMap(c *Client, i interface{}, res *OSPolicyAssignment) map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResources {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResources{}
	}

	if len(a) == 0 {
		return map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResources{}
	}

	items := make(map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResources)
	for k, item := range a {
		items[k] = *flattenOSPolicyAssignmentOSPoliciesResourceGroupsResources(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesSlice flattens the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResources from a JSON
// response object.
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesSlice(c *Client, i interface{}, res *OSPolicyAssignment) []OSPolicyAssignmentOSPoliciesResourceGroupsResources {
	a, ok := i.([]interface{})
	if !ok {
		return []OSPolicyAssignmentOSPoliciesResourceGroupsResources{}
	}

	if len(a) == 0 {
		return []OSPolicyAssignmentOSPoliciesResourceGroupsResources{}
	}

	items := make([]OSPolicyAssignmentOSPoliciesResourceGroupsResources, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenOSPolicyAssignmentOSPoliciesResourceGroupsResources(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandOSPolicyAssignmentOSPoliciesResourceGroupsResources expands an instance of OSPolicyAssignmentOSPoliciesResourceGroupsResources into a JSON
// request object.
func expandOSPolicyAssignmentOSPoliciesResourceGroupsResources(c *Client, f *OSPolicyAssignmentOSPoliciesResourceGroupsResources, res *OSPolicyAssignment) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.Id; !dcl.IsEmptyValueIndirect(v) {
		m["id"] = v
	}
	if v, err := expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkg(c, f.Pkg, res); err != nil {
		return nil, fmt.Errorf("error expanding Pkg into pkg: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["pkg"] = v
	}
	if v, err := expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepository(c, f.Repository, res); err != nil {
		return nil, fmt.Errorf("error expanding Repository into repository: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["repository"] = v
	}
	if v, err := expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExec(c, f.Exec, res); err != nil {
		return nil, fmt.Errorf("error expanding Exec into exec: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["exec"] = v
	}
	if v, err := expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFile(c, f.File, res); err != nil {
		return nil, fmt.Errorf("error expanding File into file: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["file"] = v
	}

	return m, nil
}

// flattenOSPolicyAssignmentOSPoliciesResourceGroupsResources flattens an instance of OSPolicyAssignmentOSPoliciesResourceGroupsResources from a JSON
// response object.
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsResources(c *Client, i interface{}, res *OSPolicyAssignment) *OSPolicyAssignmentOSPoliciesResourceGroupsResources {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &OSPolicyAssignmentOSPoliciesResourceGroupsResources{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResources
	}
	r.Id = dcl.FlattenString(m["id"])
	r.Pkg = flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkg(c, m["pkg"], res)
	r.Repository = flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepository(c, m["repository"], res)
	r.Exec = flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExec(c, m["exec"], res)
	r.File = flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFile(c, m["file"], res)

	return r
}

// expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMap expands the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkg into a JSON
// request object.
func expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMap(c *Client, f map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkg, res *OSPolicyAssignment) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkg(c, &item, res)
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
func expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgSlice(c *Client, f []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkg, res *OSPolicyAssignment) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkg(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMap flattens the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkg from a JSON
// response object.
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMap(c *Client, i interface{}, res *OSPolicyAssignment) map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkg {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkg{}
	}

	if len(a) == 0 {
		return map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkg{}
	}

	items := make(map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkg)
	for k, item := range a {
		items[k] = *flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkg(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgSlice flattens the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkg from a JSON
// response object.
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgSlice(c *Client, i interface{}, res *OSPolicyAssignment) []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkg {
	a, ok := i.([]interface{})
	if !ok {
		return []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkg{}
	}

	if len(a) == 0 {
		return []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkg{}
	}

	items := make([]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkg, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkg(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkg expands an instance of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkg into a JSON
// request object.
func expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkg(c *Client, f *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkg, res *OSPolicyAssignment) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.DesiredState; !dcl.IsEmptyValueIndirect(v) {
		m["desiredState"] = v
	}
	if v, err := expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgApt(c, f.Apt, res); err != nil {
		return nil, fmt.Errorf("error expanding Apt into apt: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["apt"] = v
	}
	if v, err := expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDeb(c, f.Deb, res); err != nil {
		return nil, fmt.Errorf("error expanding Deb into deb: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["deb"] = v
	}
	if v, err := expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgYum(c, f.Yum, res); err != nil {
		return nil, fmt.Errorf("error expanding Yum into yum: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["yum"] = v
	}
	if v, err := expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgZypper(c, f.Zypper, res); err != nil {
		return nil, fmt.Errorf("error expanding Zypper into zypper: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["zypper"] = v
	}
	if v, err := expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpm(c, f.Rpm, res); err != nil {
		return nil, fmt.Errorf("error expanding Rpm into rpm: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["rpm"] = v
	}
	if v, err := expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgGooget(c, f.Googet, res); err != nil {
		return nil, fmt.Errorf("error expanding Googet into googet: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["googet"] = v
	}
	if v, err := expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsi(c, f.Msi, res); err != nil {
		return nil, fmt.Errorf("error expanding Msi into msi: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["msi"] = v
	}

	return m, nil
}

// flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkg flattens an instance of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkg from a JSON
// response object.
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkg(c *Client, i interface{}, res *OSPolicyAssignment) *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkg {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkg{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkg
	}
	r.DesiredState = flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDesiredStateEnum(m["desiredState"])
	r.Apt = flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgApt(c, m["apt"], res)
	r.Deb = flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDeb(c, m["deb"], res)
	r.Yum = flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgYum(c, m["yum"], res)
	r.Zypper = flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgZypper(c, m["zypper"], res)
	r.Rpm = flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpm(c, m["rpm"], res)
	r.Googet = flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgGooget(c, m["googet"], res)
	r.Msi = flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsi(c, m["msi"], res)

	return r
}

// expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgAptMap expands the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgApt into a JSON
// request object.
func expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgAptMap(c *Client, f map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgApt, res *OSPolicyAssignment) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgApt(c, &item, res)
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
func expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgAptSlice(c *Client, f []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgApt, res *OSPolicyAssignment) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgApt(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgAptMap flattens the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgApt from a JSON
// response object.
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgAptMap(c *Client, i interface{}, res *OSPolicyAssignment) map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgApt {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgApt{}
	}

	if len(a) == 0 {
		return map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgApt{}
	}

	items := make(map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgApt)
	for k, item := range a {
		items[k] = *flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgApt(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgAptSlice flattens the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgApt from a JSON
// response object.
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgAptSlice(c *Client, i interface{}, res *OSPolicyAssignment) []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgApt {
	a, ok := i.([]interface{})
	if !ok {
		return []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgApt{}
	}

	if len(a) == 0 {
		return []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgApt{}
	}

	items := make([]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgApt, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgApt(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgApt expands an instance of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgApt into a JSON
// request object.
func expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgApt(c *Client, f *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgApt, res *OSPolicyAssignment) (map[string]interface{}, error) {
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
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgApt(c *Client, i interface{}, res *OSPolicyAssignment) *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgApt {
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
func expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebMap(c *Client, f map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDeb, res *OSPolicyAssignment) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDeb(c, &item, res)
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
func expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSlice(c *Client, f []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDeb, res *OSPolicyAssignment) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDeb(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebMap flattens the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDeb from a JSON
// response object.
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebMap(c *Client, i interface{}, res *OSPolicyAssignment) map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDeb {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDeb{}
	}

	if len(a) == 0 {
		return map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDeb{}
	}

	items := make(map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDeb)
	for k, item := range a {
		items[k] = *flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDeb(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSlice flattens the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDeb from a JSON
// response object.
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSlice(c *Client, i interface{}, res *OSPolicyAssignment) []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDeb {
	a, ok := i.([]interface{})
	if !ok {
		return []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDeb{}
	}

	if len(a) == 0 {
		return []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDeb{}
	}

	items := make([]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDeb, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDeb(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDeb expands an instance of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDeb into a JSON
// request object.
func expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDeb(c *Client, f *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDeb, res *OSPolicyAssignment) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v, err := expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSource(c, f.Source, res); err != nil {
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
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDeb(c *Client, i interface{}, res *OSPolicyAssignment) *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDeb {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDeb{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDeb
	}
	r.Source = flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSource(c, m["source"], res)
	r.PullDeps = dcl.FlattenBool(m["pullDeps"])

	return r
}

// expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceMap expands the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSource into a JSON
// request object.
func expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceMap(c *Client, f map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSource, res *OSPolicyAssignment) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSource(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceSlice expands the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSource into a JSON
// request object.
func expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceSlice(c *Client, f []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSource, res *OSPolicyAssignment) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSource(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceMap flattens the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSource from a JSON
// response object.
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceMap(c *Client, i interface{}, res *OSPolicyAssignment) map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSource {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSource{}
	}

	if len(a) == 0 {
		return map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSource{}
	}

	items := make(map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSource)
	for k, item := range a {
		items[k] = *flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSource(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceSlice flattens the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSource from a JSON
// response object.
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceSlice(c *Client, i interface{}, res *OSPolicyAssignment) []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSource {
	a, ok := i.([]interface{})
	if !ok {
		return []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSource{}
	}

	if len(a) == 0 {
		return []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSource{}
	}

	items := make([]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSource, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSource(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSource expands an instance of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSource into a JSON
// request object.
func expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSource(c *Client, f *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSource, res *OSPolicyAssignment) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v, err := expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceRemote(c, f.Remote, res); err != nil {
		return nil, fmt.Errorf("error expanding Remote into remote: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["remote"] = v
	}
	if v, err := expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceGcs(c, f.Gcs, res); err != nil {
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

// flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSource flattens an instance of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSource from a JSON
// response object.
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSource(c *Client, i interface{}, res *OSPolicyAssignment) *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSource {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSource{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSource
	}
	r.Remote = flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceRemote(c, m["remote"], res)
	r.Gcs = flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceGcs(c, m["gcs"], res)
	r.LocalPath = dcl.FlattenString(m["localPath"])
	r.AllowInsecure = dcl.FlattenBool(m["allowInsecure"])

	return r
}

// expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceRemoteMap expands the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceRemote into a JSON
// request object.
func expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceRemoteMap(c *Client, f map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceRemote, res *OSPolicyAssignment) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceRemote(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceRemoteSlice expands the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceRemote into a JSON
// request object.
func expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceRemoteSlice(c *Client, f []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceRemote, res *OSPolicyAssignment) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceRemote(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceRemoteMap flattens the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceRemote from a JSON
// response object.
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceRemoteMap(c *Client, i interface{}, res *OSPolicyAssignment) map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceRemote {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceRemote{}
	}

	if len(a) == 0 {
		return map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceRemote{}
	}

	items := make(map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceRemote)
	for k, item := range a {
		items[k] = *flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceRemote(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceRemoteSlice flattens the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceRemote from a JSON
// response object.
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceRemoteSlice(c *Client, i interface{}, res *OSPolicyAssignment) []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceRemote {
	a, ok := i.([]interface{})
	if !ok {
		return []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceRemote{}
	}

	if len(a) == 0 {
		return []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceRemote{}
	}

	items := make([]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceRemote, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceRemote(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceRemote expands an instance of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceRemote into a JSON
// request object.
func expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceRemote(c *Client, f *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceRemote, res *OSPolicyAssignment) (map[string]interface{}, error) {
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

// flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceRemote flattens an instance of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceRemote from a JSON
// response object.
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceRemote(c *Client, i interface{}, res *OSPolicyAssignment) *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceRemote {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceRemote{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceRemote
	}
	r.Uri = dcl.FlattenString(m["uri"])
	r.Sha256Checksum = dcl.FlattenString(m["sha256Checksum"])

	return r
}

// expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceGcsMap expands the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceGcs into a JSON
// request object.
func expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceGcsMap(c *Client, f map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceGcs, res *OSPolicyAssignment) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceGcs(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceGcsSlice expands the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceGcs into a JSON
// request object.
func expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceGcsSlice(c *Client, f []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceGcs, res *OSPolicyAssignment) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceGcs(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceGcsMap flattens the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceGcs from a JSON
// response object.
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceGcsMap(c *Client, i interface{}, res *OSPolicyAssignment) map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceGcs {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceGcs{}
	}

	if len(a) == 0 {
		return map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceGcs{}
	}

	items := make(map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceGcs)
	for k, item := range a {
		items[k] = *flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceGcs(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceGcsSlice flattens the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceGcs from a JSON
// response object.
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceGcsSlice(c *Client, i interface{}, res *OSPolicyAssignment) []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceGcs {
	a, ok := i.([]interface{})
	if !ok {
		return []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceGcs{}
	}

	if len(a) == 0 {
		return []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceGcs{}
	}

	items := make([]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceGcs, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceGcs(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceGcs expands an instance of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceGcs into a JSON
// request object.
func expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceGcs(c *Client, f *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceGcs, res *OSPolicyAssignment) (map[string]interface{}, error) {
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

// flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceGcs flattens an instance of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceGcs from a JSON
// response object.
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceGcs(c *Client, i interface{}, res *OSPolicyAssignment) *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceGcs {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceGcs{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceGcs
	}
	r.Bucket = dcl.FlattenString(m["bucket"])
	r.Object = dcl.FlattenString(m["object"])
	r.Generation = dcl.FlattenInteger(m["generation"])

	return r
}

// expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgYumMap expands the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgYum into a JSON
// request object.
func expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgYumMap(c *Client, f map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgYum, res *OSPolicyAssignment) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgYum(c, &item, res)
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
func expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgYumSlice(c *Client, f []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgYum, res *OSPolicyAssignment) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgYum(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgYumMap flattens the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgYum from a JSON
// response object.
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgYumMap(c *Client, i interface{}, res *OSPolicyAssignment) map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgYum {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgYum{}
	}

	if len(a) == 0 {
		return map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgYum{}
	}

	items := make(map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgYum)
	for k, item := range a {
		items[k] = *flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgYum(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgYumSlice flattens the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgYum from a JSON
// response object.
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgYumSlice(c *Client, i interface{}, res *OSPolicyAssignment) []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgYum {
	a, ok := i.([]interface{})
	if !ok {
		return []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgYum{}
	}

	if len(a) == 0 {
		return []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgYum{}
	}

	items := make([]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgYum, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgYum(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgYum expands an instance of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgYum into a JSON
// request object.
func expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgYum(c *Client, f *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgYum, res *OSPolicyAssignment) (map[string]interface{}, error) {
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
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgYum(c *Client, i interface{}, res *OSPolicyAssignment) *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgYum {
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
func expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgZypperMap(c *Client, f map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgZypper, res *OSPolicyAssignment) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgZypper(c, &item, res)
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
func expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgZypperSlice(c *Client, f []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgZypper, res *OSPolicyAssignment) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgZypper(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgZypperMap flattens the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgZypper from a JSON
// response object.
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgZypperMap(c *Client, i interface{}, res *OSPolicyAssignment) map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgZypper {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgZypper{}
	}

	if len(a) == 0 {
		return map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgZypper{}
	}

	items := make(map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgZypper)
	for k, item := range a {
		items[k] = *flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgZypper(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgZypperSlice flattens the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgZypper from a JSON
// response object.
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgZypperSlice(c *Client, i interface{}, res *OSPolicyAssignment) []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgZypper {
	a, ok := i.([]interface{})
	if !ok {
		return []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgZypper{}
	}

	if len(a) == 0 {
		return []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgZypper{}
	}

	items := make([]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgZypper, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgZypper(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgZypper expands an instance of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgZypper into a JSON
// request object.
func expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgZypper(c *Client, f *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgZypper, res *OSPolicyAssignment) (map[string]interface{}, error) {
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
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgZypper(c *Client, i interface{}, res *OSPolicyAssignment) *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgZypper {
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
func expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmMap(c *Client, f map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpm, res *OSPolicyAssignment) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpm(c, &item, res)
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
func expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSlice(c *Client, f []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpm, res *OSPolicyAssignment) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpm(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmMap flattens the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpm from a JSON
// response object.
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmMap(c *Client, i interface{}, res *OSPolicyAssignment) map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpm {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpm{}
	}

	if len(a) == 0 {
		return map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpm{}
	}

	items := make(map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpm)
	for k, item := range a {
		items[k] = *flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpm(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSlice flattens the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpm from a JSON
// response object.
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSlice(c *Client, i interface{}, res *OSPolicyAssignment) []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpm {
	a, ok := i.([]interface{})
	if !ok {
		return []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpm{}
	}

	if len(a) == 0 {
		return []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpm{}
	}

	items := make([]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpm, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpm(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpm expands an instance of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpm into a JSON
// request object.
func expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpm(c *Client, f *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpm, res *OSPolicyAssignment) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v, err := expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSource(c, f.Source, res); err != nil {
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
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpm(c *Client, i interface{}, res *OSPolicyAssignment) *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpm {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpm{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpm
	}
	r.Source = flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSource(c, m["source"], res)
	r.PullDeps = dcl.FlattenBool(m["pullDeps"])

	return r
}

// expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceMap expands the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSource into a JSON
// request object.
func expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceMap(c *Client, f map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSource, res *OSPolicyAssignment) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSource(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceSlice expands the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSource into a JSON
// request object.
func expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceSlice(c *Client, f []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSource, res *OSPolicyAssignment) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSource(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceMap flattens the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSource from a JSON
// response object.
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceMap(c *Client, i interface{}, res *OSPolicyAssignment) map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSource {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSource{}
	}

	if len(a) == 0 {
		return map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSource{}
	}

	items := make(map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSource)
	for k, item := range a {
		items[k] = *flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSource(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceSlice flattens the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSource from a JSON
// response object.
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceSlice(c *Client, i interface{}, res *OSPolicyAssignment) []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSource {
	a, ok := i.([]interface{})
	if !ok {
		return []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSource{}
	}

	if len(a) == 0 {
		return []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSource{}
	}

	items := make([]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSource, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSource(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSource expands an instance of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSource into a JSON
// request object.
func expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSource(c *Client, f *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSource, res *OSPolicyAssignment) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v, err := expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceRemote(c, f.Remote, res); err != nil {
		return nil, fmt.Errorf("error expanding Remote into remote: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["remote"] = v
	}
	if v, err := expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceGcs(c, f.Gcs, res); err != nil {
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

// flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSource flattens an instance of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSource from a JSON
// response object.
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSource(c *Client, i interface{}, res *OSPolicyAssignment) *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSource {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSource{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSource
	}
	r.Remote = flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceRemote(c, m["remote"], res)
	r.Gcs = flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceGcs(c, m["gcs"], res)
	r.LocalPath = dcl.FlattenString(m["localPath"])
	r.AllowInsecure = dcl.FlattenBool(m["allowInsecure"])

	return r
}

// expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceRemoteMap expands the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceRemote into a JSON
// request object.
func expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceRemoteMap(c *Client, f map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceRemote, res *OSPolicyAssignment) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceRemote(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceRemoteSlice expands the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceRemote into a JSON
// request object.
func expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceRemoteSlice(c *Client, f []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceRemote, res *OSPolicyAssignment) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceRemote(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceRemoteMap flattens the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceRemote from a JSON
// response object.
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceRemoteMap(c *Client, i interface{}, res *OSPolicyAssignment) map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceRemote {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceRemote{}
	}

	if len(a) == 0 {
		return map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceRemote{}
	}

	items := make(map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceRemote)
	for k, item := range a {
		items[k] = *flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceRemote(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceRemoteSlice flattens the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceRemote from a JSON
// response object.
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceRemoteSlice(c *Client, i interface{}, res *OSPolicyAssignment) []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceRemote {
	a, ok := i.([]interface{})
	if !ok {
		return []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceRemote{}
	}

	if len(a) == 0 {
		return []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceRemote{}
	}

	items := make([]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceRemote, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceRemote(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceRemote expands an instance of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceRemote into a JSON
// request object.
func expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceRemote(c *Client, f *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceRemote, res *OSPolicyAssignment) (map[string]interface{}, error) {
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

// flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceRemote flattens an instance of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceRemote from a JSON
// response object.
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceRemote(c *Client, i interface{}, res *OSPolicyAssignment) *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceRemote {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceRemote{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceRemote
	}
	r.Uri = dcl.FlattenString(m["uri"])
	r.Sha256Checksum = dcl.FlattenString(m["sha256Checksum"])

	return r
}

// expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceGcsMap expands the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceGcs into a JSON
// request object.
func expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceGcsMap(c *Client, f map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceGcs, res *OSPolicyAssignment) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceGcs(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceGcsSlice expands the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceGcs into a JSON
// request object.
func expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceGcsSlice(c *Client, f []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceGcs, res *OSPolicyAssignment) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceGcs(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceGcsMap flattens the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceGcs from a JSON
// response object.
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceGcsMap(c *Client, i interface{}, res *OSPolicyAssignment) map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceGcs {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceGcs{}
	}

	if len(a) == 0 {
		return map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceGcs{}
	}

	items := make(map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceGcs)
	for k, item := range a {
		items[k] = *flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceGcs(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceGcsSlice flattens the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceGcs from a JSON
// response object.
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceGcsSlice(c *Client, i interface{}, res *OSPolicyAssignment) []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceGcs {
	a, ok := i.([]interface{})
	if !ok {
		return []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceGcs{}
	}

	if len(a) == 0 {
		return []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceGcs{}
	}

	items := make([]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceGcs, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceGcs(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceGcs expands an instance of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceGcs into a JSON
// request object.
func expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceGcs(c *Client, f *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceGcs, res *OSPolicyAssignment) (map[string]interface{}, error) {
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

// flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceGcs flattens an instance of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceGcs from a JSON
// response object.
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceGcs(c *Client, i interface{}, res *OSPolicyAssignment) *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceGcs {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceGcs{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceGcs
	}
	r.Bucket = dcl.FlattenString(m["bucket"])
	r.Object = dcl.FlattenString(m["object"])
	r.Generation = dcl.FlattenInteger(m["generation"])

	return r
}

// expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgGoogetMap expands the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgGooget into a JSON
// request object.
func expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgGoogetMap(c *Client, f map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgGooget, res *OSPolicyAssignment) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgGooget(c, &item, res)
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
func expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgGoogetSlice(c *Client, f []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgGooget, res *OSPolicyAssignment) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgGooget(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgGoogetMap flattens the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgGooget from a JSON
// response object.
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgGoogetMap(c *Client, i interface{}, res *OSPolicyAssignment) map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgGooget {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgGooget{}
	}

	if len(a) == 0 {
		return map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgGooget{}
	}

	items := make(map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgGooget)
	for k, item := range a {
		items[k] = *flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgGooget(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgGoogetSlice flattens the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgGooget from a JSON
// response object.
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgGoogetSlice(c *Client, i interface{}, res *OSPolicyAssignment) []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgGooget {
	a, ok := i.([]interface{})
	if !ok {
		return []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgGooget{}
	}

	if len(a) == 0 {
		return []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgGooget{}
	}

	items := make([]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgGooget, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgGooget(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgGooget expands an instance of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgGooget into a JSON
// request object.
func expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgGooget(c *Client, f *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgGooget, res *OSPolicyAssignment) (map[string]interface{}, error) {
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
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgGooget(c *Client, i interface{}, res *OSPolicyAssignment) *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgGooget {
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
func expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiMap(c *Client, f map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsi, res *OSPolicyAssignment) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsi(c, &item, res)
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
func expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSlice(c *Client, f []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsi, res *OSPolicyAssignment) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsi(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiMap flattens the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsi from a JSON
// response object.
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiMap(c *Client, i interface{}, res *OSPolicyAssignment) map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsi {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsi{}
	}

	if len(a) == 0 {
		return map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsi{}
	}

	items := make(map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsi)
	for k, item := range a {
		items[k] = *flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsi(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSlice flattens the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsi from a JSON
// response object.
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSlice(c *Client, i interface{}, res *OSPolicyAssignment) []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsi {
	a, ok := i.([]interface{})
	if !ok {
		return []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsi{}
	}

	if len(a) == 0 {
		return []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsi{}
	}

	items := make([]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsi, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsi(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsi expands an instance of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsi into a JSON
// request object.
func expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsi(c *Client, f *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsi, res *OSPolicyAssignment) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v, err := expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSource(c, f.Source, res); err != nil {
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
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsi(c *Client, i interface{}, res *OSPolicyAssignment) *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsi {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsi{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsi
	}
	r.Source = flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSource(c, m["source"], res)
	r.Properties = dcl.FlattenStringSlice(m["properties"])

	return r
}

// expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceMap expands the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSource into a JSON
// request object.
func expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceMap(c *Client, f map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSource, res *OSPolicyAssignment) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSource(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceSlice expands the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSource into a JSON
// request object.
func expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceSlice(c *Client, f []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSource, res *OSPolicyAssignment) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSource(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceMap flattens the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSource from a JSON
// response object.
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceMap(c *Client, i interface{}, res *OSPolicyAssignment) map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSource {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSource{}
	}

	if len(a) == 0 {
		return map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSource{}
	}

	items := make(map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSource)
	for k, item := range a {
		items[k] = *flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSource(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceSlice flattens the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSource from a JSON
// response object.
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceSlice(c *Client, i interface{}, res *OSPolicyAssignment) []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSource {
	a, ok := i.([]interface{})
	if !ok {
		return []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSource{}
	}

	if len(a) == 0 {
		return []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSource{}
	}

	items := make([]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSource, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSource(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSource expands an instance of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSource into a JSON
// request object.
func expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSource(c *Client, f *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSource, res *OSPolicyAssignment) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v, err := expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceRemote(c, f.Remote, res); err != nil {
		return nil, fmt.Errorf("error expanding Remote into remote: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["remote"] = v
	}
	if v, err := expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceGcs(c, f.Gcs, res); err != nil {
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

// flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSource flattens an instance of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSource from a JSON
// response object.
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSource(c *Client, i interface{}, res *OSPolicyAssignment) *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSource {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSource{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSource
	}
	r.Remote = flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceRemote(c, m["remote"], res)
	r.Gcs = flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceGcs(c, m["gcs"], res)
	r.LocalPath = dcl.FlattenString(m["localPath"])
	r.AllowInsecure = dcl.FlattenBool(m["allowInsecure"])

	return r
}

// expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceRemoteMap expands the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceRemote into a JSON
// request object.
func expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceRemoteMap(c *Client, f map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceRemote, res *OSPolicyAssignment) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceRemote(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceRemoteSlice expands the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceRemote into a JSON
// request object.
func expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceRemoteSlice(c *Client, f []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceRemote, res *OSPolicyAssignment) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceRemote(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceRemoteMap flattens the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceRemote from a JSON
// response object.
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceRemoteMap(c *Client, i interface{}, res *OSPolicyAssignment) map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceRemote {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceRemote{}
	}

	if len(a) == 0 {
		return map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceRemote{}
	}

	items := make(map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceRemote)
	for k, item := range a {
		items[k] = *flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceRemote(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceRemoteSlice flattens the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceRemote from a JSON
// response object.
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceRemoteSlice(c *Client, i interface{}, res *OSPolicyAssignment) []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceRemote {
	a, ok := i.([]interface{})
	if !ok {
		return []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceRemote{}
	}

	if len(a) == 0 {
		return []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceRemote{}
	}

	items := make([]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceRemote, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceRemote(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceRemote expands an instance of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceRemote into a JSON
// request object.
func expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceRemote(c *Client, f *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceRemote, res *OSPolicyAssignment) (map[string]interface{}, error) {
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

// flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceRemote flattens an instance of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceRemote from a JSON
// response object.
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceRemote(c *Client, i interface{}, res *OSPolicyAssignment) *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceRemote {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceRemote{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceRemote
	}
	r.Uri = dcl.FlattenString(m["uri"])
	r.Sha256Checksum = dcl.FlattenString(m["sha256Checksum"])

	return r
}

// expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceGcsMap expands the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceGcs into a JSON
// request object.
func expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceGcsMap(c *Client, f map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceGcs, res *OSPolicyAssignment) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceGcs(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceGcsSlice expands the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceGcs into a JSON
// request object.
func expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceGcsSlice(c *Client, f []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceGcs, res *OSPolicyAssignment) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceGcs(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceGcsMap flattens the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceGcs from a JSON
// response object.
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceGcsMap(c *Client, i interface{}, res *OSPolicyAssignment) map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceGcs {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceGcs{}
	}

	if len(a) == 0 {
		return map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceGcs{}
	}

	items := make(map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceGcs)
	for k, item := range a {
		items[k] = *flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceGcs(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceGcsSlice flattens the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceGcs from a JSON
// response object.
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceGcsSlice(c *Client, i interface{}, res *OSPolicyAssignment) []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceGcs {
	a, ok := i.([]interface{})
	if !ok {
		return []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceGcs{}
	}

	if len(a) == 0 {
		return []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceGcs{}
	}

	items := make([]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceGcs, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceGcs(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceGcs expands an instance of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceGcs into a JSON
// request object.
func expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceGcs(c *Client, f *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceGcs, res *OSPolicyAssignment) (map[string]interface{}, error) {
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

// flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceGcs flattens an instance of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceGcs from a JSON
// response object.
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceGcs(c *Client, i interface{}, res *OSPolicyAssignment) *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceGcs {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceGcs{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceGcs
	}
	r.Bucket = dcl.FlattenString(m["bucket"])
	r.Object = dcl.FlattenString(m["object"])
	r.Generation = dcl.FlattenInteger(m["generation"])

	return r
}

// expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryMap expands the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepository into a JSON
// request object.
func expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryMap(c *Client, f map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepository, res *OSPolicyAssignment) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepository(c, &item, res)
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
func expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositorySlice(c *Client, f []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepository, res *OSPolicyAssignment) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepository(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryMap flattens the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepository from a JSON
// response object.
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryMap(c *Client, i interface{}, res *OSPolicyAssignment) map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepository {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepository{}
	}

	if len(a) == 0 {
		return map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepository{}
	}

	items := make(map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepository)
	for k, item := range a {
		items[k] = *flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepository(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositorySlice flattens the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepository from a JSON
// response object.
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositorySlice(c *Client, i interface{}, res *OSPolicyAssignment) []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepository {
	a, ok := i.([]interface{})
	if !ok {
		return []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepository{}
	}

	if len(a) == 0 {
		return []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepository{}
	}

	items := make([]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepository, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepository(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepository expands an instance of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepository into a JSON
// request object.
func expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepository(c *Client, f *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepository, res *OSPolicyAssignment) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v, err := expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryApt(c, f.Apt, res); err != nil {
		return nil, fmt.Errorf("error expanding Apt into apt: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["apt"] = v
	}
	if v, err := expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryYum(c, f.Yum, res); err != nil {
		return nil, fmt.Errorf("error expanding Yum into yum: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["yum"] = v
	}
	if v, err := expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryZypper(c, f.Zypper, res); err != nil {
		return nil, fmt.Errorf("error expanding Zypper into zypper: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["zypper"] = v
	}
	if v, err := expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryGoo(c, f.Goo, res); err != nil {
		return nil, fmt.Errorf("error expanding Goo into goo: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["goo"] = v
	}

	return m, nil
}

// flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepository flattens an instance of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepository from a JSON
// response object.
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepository(c *Client, i interface{}, res *OSPolicyAssignment) *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepository {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepository{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepository
	}
	r.Apt = flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryApt(c, m["apt"], res)
	r.Yum = flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryYum(c, m["yum"], res)
	r.Zypper = flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryZypper(c, m["zypper"], res)
	r.Goo = flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryGoo(c, m["goo"], res)

	return r
}

// expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryAptMap expands the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryApt into a JSON
// request object.
func expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryAptMap(c *Client, f map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryApt, res *OSPolicyAssignment) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryApt(c, &item, res)
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
func expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryAptSlice(c *Client, f []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryApt, res *OSPolicyAssignment) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryApt(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryAptMap flattens the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryApt from a JSON
// response object.
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryAptMap(c *Client, i interface{}, res *OSPolicyAssignment) map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryApt {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryApt{}
	}

	if len(a) == 0 {
		return map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryApt{}
	}

	items := make(map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryApt)
	for k, item := range a {
		items[k] = *flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryApt(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryAptSlice flattens the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryApt from a JSON
// response object.
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryAptSlice(c *Client, i interface{}, res *OSPolicyAssignment) []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryApt {
	a, ok := i.([]interface{})
	if !ok {
		return []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryApt{}
	}

	if len(a) == 0 {
		return []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryApt{}
	}

	items := make([]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryApt, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryApt(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryApt expands an instance of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryApt into a JSON
// request object.
func expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryApt(c *Client, f *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryApt, res *OSPolicyAssignment) (map[string]interface{}, error) {
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
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryApt(c *Client, i interface{}, res *OSPolicyAssignment) *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryApt {
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
func expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryYumMap(c *Client, f map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryYum, res *OSPolicyAssignment) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryYum(c, &item, res)
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
func expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryYumSlice(c *Client, f []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryYum, res *OSPolicyAssignment) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryYum(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryYumMap flattens the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryYum from a JSON
// response object.
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryYumMap(c *Client, i interface{}, res *OSPolicyAssignment) map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryYum {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryYum{}
	}

	if len(a) == 0 {
		return map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryYum{}
	}

	items := make(map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryYum)
	for k, item := range a {
		items[k] = *flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryYum(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryYumSlice flattens the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryYum from a JSON
// response object.
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryYumSlice(c *Client, i interface{}, res *OSPolicyAssignment) []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryYum {
	a, ok := i.([]interface{})
	if !ok {
		return []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryYum{}
	}

	if len(a) == 0 {
		return []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryYum{}
	}

	items := make([]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryYum, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryYum(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryYum expands an instance of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryYum into a JSON
// request object.
func expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryYum(c *Client, f *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryYum, res *OSPolicyAssignment) (map[string]interface{}, error) {
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
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryYum(c *Client, i interface{}, res *OSPolicyAssignment) *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryYum {
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
func expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryZypperMap(c *Client, f map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryZypper, res *OSPolicyAssignment) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryZypper(c, &item, res)
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
func expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryZypperSlice(c *Client, f []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryZypper, res *OSPolicyAssignment) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryZypper(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryZypperMap flattens the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryZypper from a JSON
// response object.
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryZypperMap(c *Client, i interface{}, res *OSPolicyAssignment) map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryZypper {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryZypper{}
	}

	if len(a) == 0 {
		return map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryZypper{}
	}

	items := make(map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryZypper)
	for k, item := range a {
		items[k] = *flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryZypper(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryZypperSlice flattens the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryZypper from a JSON
// response object.
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryZypperSlice(c *Client, i interface{}, res *OSPolicyAssignment) []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryZypper {
	a, ok := i.([]interface{})
	if !ok {
		return []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryZypper{}
	}

	if len(a) == 0 {
		return []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryZypper{}
	}

	items := make([]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryZypper, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryZypper(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryZypper expands an instance of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryZypper into a JSON
// request object.
func expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryZypper(c *Client, f *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryZypper, res *OSPolicyAssignment) (map[string]interface{}, error) {
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
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryZypper(c *Client, i interface{}, res *OSPolicyAssignment) *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryZypper {
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
func expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryGooMap(c *Client, f map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryGoo, res *OSPolicyAssignment) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryGoo(c, &item, res)
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
func expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryGooSlice(c *Client, f []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryGoo, res *OSPolicyAssignment) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryGoo(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryGooMap flattens the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryGoo from a JSON
// response object.
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryGooMap(c *Client, i interface{}, res *OSPolicyAssignment) map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryGoo {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryGoo{}
	}

	if len(a) == 0 {
		return map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryGoo{}
	}

	items := make(map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryGoo)
	for k, item := range a {
		items[k] = *flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryGoo(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryGooSlice flattens the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryGoo from a JSON
// response object.
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryGooSlice(c *Client, i interface{}, res *OSPolicyAssignment) []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryGoo {
	a, ok := i.([]interface{})
	if !ok {
		return []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryGoo{}
	}

	if len(a) == 0 {
		return []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryGoo{}
	}

	items := make([]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryGoo, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryGoo(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryGoo expands an instance of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryGoo into a JSON
// request object.
func expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryGoo(c *Client, f *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryGoo, res *OSPolicyAssignment) (map[string]interface{}, error) {
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
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryGoo(c *Client, i interface{}, res *OSPolicyAssignment) *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryGoo {
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
func expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecMap(c *Client, f map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExec, res *OSPolicyAssignment) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExec(c, &item, res)
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
func expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecSlice(c *Client, f []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExec, res *OSPolicyAssignment) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExec(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecMap flattens the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExec from a JSON
// response object.
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecMap(c *Client, i interface{}, res *OSPolicyAssignment) map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExec {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExec{}
	}

	if len(a) == 0 {
		return map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExec{}
	}

	items := make(map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExec)
	for k, item := range a {
		items[k] = *flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExec(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecSlice flattens the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExec from a JSON
// response object.
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecSlice(c *Client, i interface{}, res *OSPolicyAssignment) []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExec {
	a, ok := i.([]interface{})
	if !ok {
		return []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExec{}
	}

	if len(a) == 0 {
		return []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExec{}
	}

	items := make([]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExec, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExec(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExec expands an instance of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExec into a JSON
// request object.
func expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExec(c *Client, f *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExec, res *OSPolicyAssignment) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v, err := expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidate(c, f.Validate, res); err != nil {
		return nil, fmt.Errorf("error expanding Validate into validate: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["validate"] = v
	}
	if v, err := expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforce(c, f.Enforce, res); err != nil {
		return nil, fmt.Errorf("error expanding Enforce into enforce: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["enforce"] = v
	}

	return m, nil
}

// flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExec flattens an instance of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExec from a JSON
// response object.
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExec(c *Client, i interface{}, res *OSPolicyAssignment) *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExec {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExec{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExec
	}
	r.Validate = flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidate(c, m["validate"], res)
	r.Enforce = flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforce(c, m["enforce"], res)

	return r
}

// expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateMap expands the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidate into a JSON
// request object.
func expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateMap(c *Client, f map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidate, res *OSPolicyAssignment) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidate(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateSlice expands the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidate into a JSON
// request object.
func expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateSlice(c *Client, f []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidate, res *OSPolicyAssignment) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidate(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateMap flattens the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidate from a JSON
// response object.
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateMap(c *Client, i interface{}, res *OSPolicyAssignment) map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidate {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidate{}
	}

	if len(a) == 0 {
		return map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidate{}
	}

	items := make(map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidate)
	for k, item := range a {
		items[k] = *flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidate(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateSlice flattens the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidate from a JSON
// response object.
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateSlice(c *Client, i interface{}, res *OSPolicyAssignment) []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidate {
	a, ok := i.([]interface{})
	if !ok {
		return []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidate{}
	}

	if len(a) == 0 {
		return []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidate{}
	}

	items := make([]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidate, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidate(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidate expands an instance of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidate into a JSON
// request object.
func expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidate(c *Client, f *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidate, res *OSPolicyAssignment) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v, err := expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFile(c, f.File, res); err != nil {
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

// flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidate flattens an instance of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidate from a JSON
// response object.
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidate(c *Client, i interface{}, res *OSPolicyAssignment) *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidate {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidate{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidate
	}
	r.File = flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFile(c, m["file"], res)
	r.Script = dcl.FlattenString(m["script"])
	r.Args = dcl.FlattenStringSlice(m["args"])
	r.Interpreter = flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateInterpreterEnum(m["interpreter"])
	r.OutputFilePath = dcl.FlattenString(m["outputFilePath"])

	return r
}

// expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileMap expands the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFile into a JSON
// request object.
func expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileMap(c *Client, f map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFile, res *OSPolicyAssignment) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFile(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileSlice expands the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFile into a JSON
// request object.
func expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileSlice(c *Client, f []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFile, res *OSPolicyAssignment) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFile(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileMap flattens the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFile from a JSON
// response object.
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileMap(c *Client, i interface{}, res *OSPolicyAssignment) map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFile {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFile{}
	}

	if len(a) == 0 {
		return map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFile{}
	}

	items := make(map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFile)
	for k, item := range a {
		items[k] = *flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFile(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileSlice flattens the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFile from a JSON
// response object.
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileSlice(c *Client, i interface{}, res *OSPolicyAssignment) []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFile {
	a, ok := i.([]interface{})
	if !ok {
		return []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFile{}
	}

	if len(a) == 0 {
		return []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFile{}
	}

	items := make([]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFile, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFile(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFile expands an instance of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFile into a JSON
// request object.
func expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFile(c *Client, f *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFile, res *OSPolicyAssignment) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v, err := expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileRemote(c, f.Remote, res); err != nil {
		return nil, fmt.Errorf("error expanding Remote into remote: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["remote"] = v
	}
	if v, err := expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileGcs(c, f.Gcs, res); err != nil {
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

// flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFile flattens an instance of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFile from a JSON
// response object.
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFile(c *Client, i interface{}, res *OSPolicyAssignment) *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFile {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFile{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFile
	}
	r.Remote = flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileRemote(c, m["remote"], res)
	r.Gcs = flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileGcs(c, m["gcs"], res)
	r.LocalPath = dcl.FlattenString(m["localPath"])
	r.AllowInsecure = dcl.FlattenBool(m["allowInsecure"])

	return r
}

// expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileRemoteMap expands the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileRemote into a JSON
// request object.
func expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileRemoteMap(c *Client, f map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileRemote, res *OSPolicyAssignment) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileRemote(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileRemoteSlice expands the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileRemote into a JSON
// request object.
func expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileRemoteSlice(c *Client, f []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileRemote, res *OSPolicyAssignment) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileRemote(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileRemoteMap flattens the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileRemote from a JSON
// response object.
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileRemoteMap(c *Client, i interface{}, res *OSPolicyAssignment) map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileRemote {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileRemote{}
	}

	if len(a) == 0 {
		return map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileRemote{}
	}

	items := make(map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileRemote)
	for k, item := range a {
		items[k] = *flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileRemote(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileRemoteSlice flattens the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileRemote from a JSON
// response object.
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileRemoteSlice(c *Client, i interface{}, res *OSPolicyAssignment) []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileRemote {
	a, ok := i.([]interface{})
	if !ok {
		return []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileRemote{}
	}

	if len(a) == 0 {
		return []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileRemote{}
	}

	items := make([]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileRemote, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileRemote(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileRemote expands an instance of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileRemote into a JSON
// request object.
func expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileRemote(c *Client, f *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileRemote, res *OSPolicyAssignment) (map[string]interface{}, error) {
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

// flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileRemote flattens an instance of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileRemote from a JSON
// response object.
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileRemote(c *Client, i interface{}, res *OSPolicyAssignment) *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileRemote {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileRemote{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileRemote
	}
	r.Uri = dcl.FlattenString(m["uri"])
	r.Sha256Checksum = dcl.FlattenString(m["sha256Checksum"])

	return r
}

// expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileGcsMap expands the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileGcs into a JSON
// request object.
func expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileGcsMap(c *Client, f map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileGcs, res *OSPolicyAssignment) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileGcs(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileGcsSlice expands the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileGcs into a JSON
// request object.
func expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileGcsSlice(c *Client, f []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileGcs, res *OSPolicyAssignment) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileGcs(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileGcsMap flattens the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileGcs from a JSON
// response object.
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileGcsMap(c *Client, i interface{}, res *OSPolicyAssignment) map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileGcs {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileGcs{}
	}

	if len(a) == 0 {
		return map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileGcs{}
	}

	items := make(map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileGcs)
	for k, item := range a {
		items[k] = *flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileGcs(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileGcsSlice flattens the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileGcs from a JSON
// response object.
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileGcsSlice(c *Client, i interface{}, res *OSPolicyAssignment) []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileGcs {
	a, ok := i.([]interface{})
	if !ok {
		return []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileGcs{}
	}

	if len(a) == 0 {
		return []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileGcs{}
	}

	items := make([]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileGcs, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileGcs(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileGcs expands an instance of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileGcs into a JSON
// request object.
func expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileGcs(c *Client, f *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileGcs, res *OSPolicyAssignment) (map[string]interface{}, error) {
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

// flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileGcs flattens an instance of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileGcs from a JSON
// response object.
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileGcs(c *Client, i interface{}, res *OSPolicyAssignment) *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileGcs {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileGcs{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileGcs
	}
	r.Bucket = dcl.FlattenString(m["bucket"])
	r.Object = dcl.FlattenString(m["object"])
	r.Generation = dcl.FlattenInteger(m["generation"])

	return r
}

// expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceMap expands the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforce into a JSON
// request object.
func expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceMap(c *Client, f map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforce, res *OSPolicyAssignment) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforce(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceSlice expands the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforce into a JSON
// request object.
func expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceSlice(c *Client, f []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforce, res *OSPolicyAssignment) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforce(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceMap flattens the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforce from a JSON
// response object.
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceMap(c *Client, i interface{}, res *OSPolicyAssignment) map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforce {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforce{}
	}

	if len(a) == 0 {
		return map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforce{}
	}

	items := make(map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforce)
	for k, item := range a {
		items[k] = *flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforce(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceSlice flattens the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforce from a JSON
// response object.
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceSlice(c *Client, i interface{}, res *OSPolicyAssignment) []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforce {
	a, ok := i.([]interface{})
	if !ok {
		return []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforce{}
	}

	if len(a) == 0 {
		return []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforce{}
	}

	items := make([]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforce, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforce(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforce expands an instance of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforce into a JSON
// request object.
func expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforce(c *Client, f *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforce, res *OSPolicyAssignment) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v, err := expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFile(c, f.File, res); err != nil {
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

// flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforce flattens an instance of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforce from a JSON
// response object.
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforce(c *Client, i interface{}, res *OSPolicyAssignment) *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforce {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforce{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforce
	}
	r.File = flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFile(c, m["file"], res)
	r.Script = dcl.FlattenString(m["script"])
	r.Args = dcl.FlattenStringSlice(m["args"])
	r.Interpreter = flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceInterpreterEnum(m["interpreter"])
	r.OutputFilePath = dcl.FlattenString(m["outputFilePath"])

	return r
}

// expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileMap expands the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFile into a JSON
// request object.
func expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileMap(c *Client, f map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFile, res *OSPolicyAssignment) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFile(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileSlice expands the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFile into a JSON
// request object.
func expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileSlice(c *Client, f []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFile, res *OSPolicyAssignment) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFile(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileMap flattens the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFile from a JSON
// response object.
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileMap(c *Client, i interface{}, res *OSPolicyAssignment) map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFile {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFile{}
	}

	if len(a) == 0 {
		return map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFile{}
	}

	items := make(map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFile)
	for k, item := range a {
		items[k] = *flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFile(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileSlice flattens the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFile from a JSON
// response object.
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileSlice(c *Client, i interface{}, res *OSPolicyAssignment) []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFile {
	a, ok := i.([]interface{})
	if !ok {
		return []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFile{}
	}

	if len(a) == 0 {
		return []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFile{}
	}

	items := make([]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFile, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFile(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFile expands an instance of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFile into a JSON
// request object.
func expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFile(c *Client, f *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFile, res *OSPolicyAssignment) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v, err := expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileRemote(c, f.Remote, res); err != nil {
		return nil, fmt.Errorf("error expanding Remote into remote: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["remote"] = v
	}
	if v, err := expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileGcs(c, f.Gcs, res); err != nil {
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

// flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFile flattens an instance of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFile from a JSON
// response object.
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFile(c *Client, i interface{}, res *OSPolicyAssignment) *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFile {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFile{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFile
	}
	r.Remote = flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileRemote(c, m["remote"], res)
	r.Gcs = flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileGcs(c, m["gcs"], res)
	r.LocalPath = dcl.FlattenString(m["localPath"])
	r.AllowInsecure = dcl.FlattenBool(m["allowInsecure"])

	return r
}

// expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileRemoteMap expands the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileRemote into a JSON
// request object.
func expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileRemoteMap(c *Client, f map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileRemote, res *OSPolicyAssignment) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileRemote(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileRemoteSlice expands the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileRemote into a JSON
// request object.
func expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileRemoteSlice(c *Client, f []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileRemote, res *OSPolicyAssignment) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileRemote(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileRemoteMap flattens the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileRemote from a JSON
// response object.
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileRemoteMap(c *Client, i interface{}, res *OSPolicyAssignment) map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileRemote {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileRemote{}
	}

	if len(a) == 0 {
		return map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileRemote{}
	}

	items := make(map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileRemote)
	for k, item := range a {
		items[k] = *flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileRemote(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileRemoteSlice flattens the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileRemote from a JSON
// response object.
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileRemoteSlice(c *Client, i interface{}, res *OSPolicyAssignment) []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileRemote {
	a, ok := i.([]interface{})
	if !ok {
		return []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileRemote{}
	}

	if len(a) == 0 {
		return []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileRemote{}
	}

	items := make([]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileRemote, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileRemote(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileRemote expands an instance of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileRemote into a JSON
// request object.
func expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileRemote(c *Client, f *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileRemote, res *OSPolicyAssignment) (map[string]interface{}, error) {
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

// flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileRemote flattens an instance of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileRemote from a JSON
// response object.
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileRemote(c *Client, i interface{}, res *OSPolicyAssignment) *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileRemote {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileRemote{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileRemote
	}
	r.Uri = dcl.FlattenString(m["uri"])
	r.Sha256Checksum = dcl.FlattenString(m["sha256Checksum"])

	return r
}

// expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileGcsMap expands the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileGcs into a JSON
// request object.
func expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileGcsMap(c *Client, f map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileGcs, res *OSPolicyAssignment) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileGcs(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileGcsSlice expands the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileGcs into a JSON
// request object.
func expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileGcsSlice(c *Client, f []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileGcs, res *OSPolicyAssignment) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileGcs(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileGcsMap flattens the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileGcs from a JSON
// response object.
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileGcsMap(c *Client, i interface{}, res *OSPolicyAssignment) map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileGcs {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileGcs{}
	}

	if len(a) == 0 {
		return map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileGcs{}
	}

	items := make(map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileGcs)
	for k, item := range a {
		items[k] = *flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileGcs(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileGcsSlice flattens the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileGcs from a JSON
// response object.
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileGcsSlice(c *Client, i interface{}, res *OSPolicyAssignment) []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileGcs {
	a, ok := i.([]interface{})
	if !ok {
		return []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileGcs{}
	}

	if len(a) == 0 {
		return []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileGcs{}
	}

	items := make([]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileGcs, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileGcs(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileGcs expands an instance of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileGcs into a JSON
// request object.
func expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileGcs(c *Client, f *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileGcs, res *OSPolicyAssignment) (map[string]interface{}, error) {
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

// flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileGcs flattens an instance of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileGcs from a JSON
// response object.
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileGcs(c *Client, i interface{}, res *OSPolicyAssignment) *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileGcs {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileGcs{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileGcs
	}
	r.Bucket = dcl.FlattenString(m["bucket"])
	r.Object = dcl.FlattenString(m["object"])
	r.Generation = dcl.FlattenInteger(m["generation"])

	return r
}

// expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileMap expands the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFile into a JSON
// request object.
func expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileMap(c *Client, f map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFile, res *OSPolicyAssignment) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFile(c, &item, res)
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
func expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileSlice(c *Client, f []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFile, res *OSPolicyAssignment) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFile(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileMap flattens the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFile from a JSON
// response object.
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileMap(c *Client, i interface{}, res *OSPolicyAssignment) map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFile {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFile{}
	}

	if len(a) == 0 {
		return map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFile{}
	}

	items := make(map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFile)
	for k, item := range a {
		items[k] = *flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFile(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileSlice flattens the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFile from a JSON
// response object.
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileSlice(c *Client, i interface{}, res *OSPolicyAssignment) []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFile {
	a, ok := i.([]interface{})
	if !ok {
		return []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFile{}
	}

	if len(a) == 0 {
		return []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFile{}
	}

	items := make([]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFile, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFile(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFile expands an instance of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFile into a JSON
// request object.
func expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFile(c *Client, f *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFile, res *OSPolicyAssignment) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v, err := expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFile(c, f.File, res); err != nil {
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
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFile(c *Client, i interface{}, res *OSPolicyAssignment) *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFile {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFile{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFile
	}
	r.File = flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFile(c, m["file"], res)
	r.Content = dcl.FlattenString(m["content"])
	r.Path = dcl.FlattenString(m["path"])
	r.State = flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileStateEnum(m["state"])
	r.Permissions = dcl.FlattenString(m["permissions"])

	return r
}

// expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileMap expands the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFile into a JSON
// request object.
func expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileMap(c *Client, f map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFile, res *OSPolicyAssignment) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFile(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileSlice expands the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFile into a JSON
// request object.
func expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileSlice(c *Client, f []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFile, res *OSPolicyAssignment) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFile(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileMap flattens the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFile from a JSON
// response object.
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileMap(c *Client, i interface{}, res *OSPolicyAssignment) map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFile {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFile{}
	}

	if len(a) == 0 {
		return map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFile{}
	}

	items := make(map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFile)
	for k, item := range a {
		items[k] = *flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFile(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileSlice flattens the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFile from a JSON
// response object.
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileSlice(c *Client, i interface{}, res *OSPolicyAssignment) []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFile {
	a, ok := i.([]interface{})
	if !ok {
		return []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFile{}
	}

	if len(a) == 0 {
		return []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFile{}
	}

	items := make([]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFile, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFile(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFile expands an instance of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFile into a JSON
// request object.
func expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFile(c *Client, f *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFile, res *OSPolicyAssignment) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v, err := expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileRemote(c, f.Remote, res); err != nil {
		return nil, fmt.Errorf("error expanding Remote into remote: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["remote"] = v
	}
	if v, err := expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileGcs(c, f.Gcs, res); err != nil {
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

// flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFile flattens an instance of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFile from a JSON
// response object.
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFile(c *Client, i interface{}, res *OSPolicyAssignment) *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFile {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFile{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFile
	}
	r.Remote = flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileRemote(c, m["remote"], res)
	r.Gcs = flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileGcs(c, m["gcs"], res)
	r.LocalPath = dcl.FlattenString(m["localPath"])
	r.AllowInsecure = dcl.FlattenBool(m["allowInsecure"])

	return r
}

// expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileRemoteMap expands the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileRemote into a JSON
// request object.
func expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileRemoteMap(c *Client, f map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileRemote, res *OSPolicyAssignment) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileRemote(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileRemoteSlice expands the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileRemote into a JSON
// request object.
func expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileRemoteSlice(c *Client, f []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileRemote, res *OSPolicyAssignment) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileRemote(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileRemoteMap flattens the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileRemote from a JSON
// response object.
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileRemoteMap(c *Client, i interface{}, res *OSPolicyAssignment) map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileRemote {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileRemote{}
	}

	if len(a) == 0 {
		return map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileRemote{}
	}

	items := make(map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileRemote)
	for k, item := range a {
		items[k] = *flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileRemote(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileRemoteSlice flattens the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileRemote from a JSON
// response object.
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileRemoteSlice(c *Client, i interface{}, res *OSPolicyAssignment) []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileRemote {
	a, ok := i.([]interface{})
	if !ok {
		return []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileRemote{}
	}

	if len(a) == 0 {
		return []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileRemote{}
	}

	items := make([]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileRemote, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileRemote(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileRemote expands an instance of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileRemote into a JSON
// request object.
func expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileRemote(c *Client, f *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileRemote, res *OSPolicyAssignment) (map[string]interface{}, error) {
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

// flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileRemote flattens an instance of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileRemote from a JSON
// response object.
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileRemote(c *Client, i interface{}, res *OSPolicyAssignment) *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileRemote {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileRemote{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileRemote
	}
	r.Uri = dcl.FlattenString(m["uri"])
	r.Sha256Checksum = dcl.FlattenString(m["sha256Checksum"])

	return r
}

// expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileGcsMap expands the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileGcs into a JSON
// request object.
func expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileGcsMap(c *Client, f map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileGcs, res *OSPolicyAssignment) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileGcs(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileGcsSlice expands the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileGcs into a JSON
// request object.
func expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileGcsSlice(c *Client, f []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileGcs, res *OSPolicyAssignment) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileGcs(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileGcsMap flattens the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileGcs from a JSON
// response object.
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileGcsMap(c *Client, i interface{}, res *OSPolicyAssignment) map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileGcs {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileGcs{}
	}

	if len(a) == 0 {
		return map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileGcs{}
	}

	items := make(map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileGcs)
	for k, item := range a {
		items[k] = *flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileGcs(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileGcsSlice flattens the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileGcs from a JSON
// response object.
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileGcsSlice(c *Client, i interface{}, res *OSPolicyAssignment) []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileGcs {
	a, ok := i.([]interface{})
	if !ok {
		return []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileGcs{}
	}

	if len(a) == 0 {
		return []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileGcs{}
	}

	items := make([]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileGcs, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileGcs(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileGcs expands an instance of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileGcs into a JSON
// request object.
func expandOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileGcs(c *Client, f *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileGcs, res *OSPolicyAssignment) (map[string]interface{}, error) {
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

// flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileGcs flattens an instance of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileGcs from a JSON
// response object.
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileGcs(c *Client, i interface{}, res *OSPolicyAssignment) *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileGcs {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileGcs{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileGcs
	}
	r.Bucket = dcl.FlattenString(m["bucket"])
	r.Object = dcl.FlattenString(m["object"])
	r.Generation = dcl.FlattenInteger(m["generation"])

	return r
}

// expandOSPolicyAssignmentInstanceFilterMap expands the contents of OSPolicyAssignmentInstanceFilter into a JSON
// request object.
func expandOSPolicyAssignmentInstanceFilterMap(c *Client, f map[string]OSPolicyAssignmentInstanceFilter, res *OSPolicyAssignment) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandOSPolicyAssignmentInstanceFilter(c, &item, res)
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
func expandOSPolicyAssignmentInstanceFilterSlice(c *Client, f []OSPolicyAssignmentInstanceFilter, res *OSPolicyAssignment) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandOSPolicyAssignmentInstanceFilter(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenOSPolicyAssignmentInstanceFilterMap flattens the contents of OSPolicyAssignmentInstanceFilter from a JSON
// response object.
func flattenOSPolicyAssignmentInstanceFilterMap(c *Client, i interface{}, res *OSPolicyAssignment) map[string]OSPolicyAssignmentInstanceFilter {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]OSPolicyAssignmentInstanceFilter{}
	}

	if len(a) == 0 {
		return map[string]OSPolicyAssignmentInstanceFilter{}
	}

	items := make(map[string]OSPolicyAssignmentInstanceFilter)
	for k, item := range a {
		items[k] = *flattenOSPolicyAssignmentInstanceFilter(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenOSPolicyAssignmentInstanceFilterSlice flattens the contents of OSPolicyAssignmentInstanceFilter from a JSON
// response object.
func flattenOSPolicyAssignmentInstanceFilterSlice(c *Client, i interface{}, res *OSPolicyAssignment) []OSPolicyAssignmentInstanceFilter {
	a, ok := i.([]interface{})
	if !ok {
		return []OSPolicyAssignmentInstanceFilter{}
	}

	if len(a) == 0 {
		return []OSPolicyAssignmentInstanceFilter{}
	}

	items := make([]OSPolicyAssignmentInstanceFilter, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenOSPolicyAssignmentInstanceFilter(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandOSPolicyAssignmentInstanceFilter expands an instance of OSPolicyAssignmentInstanceFilter into a JSON
// request object.
func expandOSPolicyAssignmentInstanceFilter(c *Client, f *OSPolicyAssignmentInstanceFilter, res *OSPolicyAssignment) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.All; v != nil {
		m["all"] = v
	}
	if v, err := expandOSPolicyAssignmentInstanceFilterInclusionLabelsSlice(c, f.InclusionLabels, res); err != nil {
		return nil, fmt.Errorf("error expanding InclusionLabels into inclusionLabels: %w", err)
	} else if v != nil {
		m["inclusionLabels"] = v
	}
	if v, err := expandOSPolicyAssignmentInstanceFilterExclusionLabelsSlice(c, f.ExclusionLabels, res); err != nil {
		return nil, fmt.Errorf("error expanding ExclusionLabels into exclusionLabels: %w", err)
	} else if v != nil {
		m["exclusionLabels"] = v
	}
	if v, err := expandOSPolicyAssignmentInstanceFilterInventoriesSlice(c, f.Inventories, res); err != nil {
		return nil, fmt.Errorf("error expanding Inventories into inventories: %w", err)
	} else if v != nil {
		m["inventories"] = v
	}

	return m, nil
}

// flattenOSPolicyAssignmentInstanceFilter flattens an instance of OSPolicyAssignmentInstanceFilter from a JSON
// response object.
func flattenOSPolicyAssignmentInstanceFilter(c *Client, i interface{}, res *OSPolicyAssignment) *OSPolicyAssignmentInstanceFilter {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &OSPolicyAssignmentInstanceFilter{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyOSPolicyAssignmentInstanceFilter
	}
	r.All = dcl.FlattenBool(m["all"])
	r.InclusionLabels = flattenOSPolicyAssignmentInstanceFilterInclusionLabelsSlice(c, m["inclusionLabels"], res)
	r.ExclusionLabels = flattenOSPolicyAssignmentInstanceFilterExclusionLabelsSlice(c, m["exclusionLabels"], res)
	r.Inventories = flattenOSPolicyAssignmentInstanceFilterInventoriesSlice(c, m["inventories"], res)

	return r
}

// expandOSPolicyAssignmentInstanceFilterInclusionLabelsMap expands the contents of OSPolicyAssignmentInstanceFilterInclusionLabels into a JSON
// request object.
func expandOSPolicyAssignmentInstanceFilterInclusionLabelsMap(c *Client, f map[string]OSPolicyAssignmentInstanceFilterInclusionLabels, res *OSPolicyAssignment) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandOSPolicyAssignmentInstanceFilterInclusionLabels(c, &item, res)
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
func expandOSPolicyAssignmentInstanceFilterInclusionLabelsSlice(c *Client, f []OSPolicyAssignmentInstanceFilterInclusionLabels, res *OSPolicyAssignment) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandOSPolicyAssignmentInstanceFilterInclusionLabels(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenOSPolicyAssignmentInstanceFilterInclusionLabelsMap flattens the contents of OSPolicyAssignmentInstanceFilterInclusionLabels from a JSON
// response object.
func flattenOSPolicyAssignmentInstanceFilterInclusionLabelsMap(c *Client, i interface{}, res *OSPolicyAssignment) map[string]OSPolicyAssignmentInstanceFilterInclusionLabels {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]OSPolicyAssignmentInstanceFilterInclusionLabels{}
	}

	if len(a) == 0 {
		return map[string]OSPolicyAssignmentInstanceFilterInclusionLabels{}
	}

	items := make(map[string]OSPolicyAssignmentInstanceFilterInclusionLabels)
	for k, item := range a {
		items[k] = *flattenOSPolicyAssignmentInstanceFilterInclusionLabels(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenOSPolicyAssignmentInstanceFilterInclusionLabelsSlice flattens the contents of OSPolicyAssignmentInstanceFilterInclusionLabels from a JSON
// response object.
func flattenOSPolicyAssignmentInstanceFilterInclusionLabelsSlice(c *Client, i interface{}, res *OSPolicyAssignment) []OSPolicyAssignmentInstanceFilterInclusionLabels {
	a, ok := i.([]interface{})
	if !ok {
		return []OSPolicyAssignmentInstanceFilterInclusionLabels{}
	}

	if len(a) == 0 {
		return []OSPolicyAssignmentInstanceFilterInclusionLabels{}
	}

	items := make([]OSPolicyAssignmentInstanceFilterInclusionLabels, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenOSPolicyAssignmentInstanceFilterInclusionLabels(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandOSPolicyAssignmentInstanceFilterInclusionLabels expands an instance of OSPolicyAssignmentInstanceFilterInclusionLabels into a JSON
// request object.
func expandOSPolicyAssignmentInstanceFilterInclusionLabels(c *Client, f *OSPolicyAssignmentInstanceFilterInclusionLabels, res *OSPolicyAssignment) (map[string]interface{}, error) {
	if f == nil {
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
func flattenOSPolicyAssignmentInstanceFilterInclusionLabels(c *Client, i interface{}, res *OSPolicyAssignment) *OSPolicyAssignmentInstanceFilterInclusionLabels {
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
func expandOSPolicyAssignmentInstanceFilterExclusionLabelsMap(c *Client, f map[string]OSPolicyAssignmentInstanceFilterExclusionLabels, res *OSPolicyAssignment) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandOSPolicyAssignmentInstanceFilterExclusionLabels(c, &item, res)
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
func expandOSPolicyAssignmentInstanceFilterExclusionLabelsSlice(c *Client, f []OSPolicyAssignmentInstanceFilterExclusionLabels, res *OSPolicyAssignment) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandOSPolicyAssignmentInstanceFilterExclusionLabels(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenOSPolicyAssignmentInstanceFilterExclusionLabelsMap flattens the contents of OSPolicyAssignmentInstanceFilterExclusionLabels from a JSON
// response object.
func flattenOSPolicyAssignmentInstanceFilterExclusionLabelsMap(c *Client, i interface{}, res *OSPolicyAssignment) map[string]OSPolicyAssignmentInstanceFilterExclusionLabels {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]OSPolicyAssignmentInstanceFilterExclusionLabels{}
	}

	if len(a) == 0 {
		return map[string]OSPolicyAssignmentInstanceFilterExclusionLabels{}
	}

	items := make(map[string]OSPolicyAssignmentInstanceFilterExclusionLabels)
	for k, item := range a {
		items[k] = *flattenOSPolicyAssignmentInstanceFilterExclusionLabels(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenOSPolicyAssignmentInstanceFilterExclusionLabelsSlice flattens the contents of OSPolicyAssignmentInstanceFilterExclusionLabels from a JSON
// response object.
func flattenOSPolicyAssignmentInstanceFilterExclusionLabelsSlice(c *Client, i interface{}, res *OSPolicyAssignment) []OSPolicyAssignmentInstanceFilterExclusionLabels {
	a, ok := i.([]interface{})
	if !ok {
		return []OSPolicyAssignmentInstanceFilterExclusionLabels{}
	}

	if len(a) == 0 {
		return []OSPolicyAssignmentInstanceFilterExclusionLabels{}
	}

	items := make([]OSPolicyAssignmentInstanceFilterExclusionLabels, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenOSPolicyAssignmentInstanceFilterExclusionLabels(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandOSPolicyAssignmentInstanceFilterExclusionLabels expands an instance of OSPolicyAssignmentInstanceFilterExclusionLabels into a JSON
// request object.
func expandOSPolicyAssignmentInstanceFilterExclusionLabels(c *Client, f *OSPolicyAssignmentInstanceFilterExclusionLabels, res *OSPolicyAssignment) (map[string]interface{}, error) {
	if f == nil {
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
func flattenOSPolicyAssignmentInstanceFilterExclusionLabels(c *Client, i interface{}, res *OSPolicyAssignment) *OSPolicyAssignmentInstanceFilterExclusionLabels {
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
func expandOSPolicyAssignmentInstanceFilterInventoriesMap(c *Client, f map[string]OSPolicyAssignmentInstanceFilterInventories, res *OSPolicyAssignment) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandOSPolicyAssignmentInstanceFilterInventories(c, &item, res)
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
func expandOSPolicyAssignmentInstanceFilterInventoriesSlice(c *Client, f []OSPolicyAssignmentInstanceFilterInventories, res *OSPolicyAssignment) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandOSPolicyAssignmentInstanceFilterInventories(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenOSPolicyAssignmentInstanceFilterInventoriesMap flattens the contents of OSPolicyAssignmentInstanceFilterInventories from a JSON
// response object.
func flattenOSPolicyAssignmentInstanceFilterInventoriesMap(c *Client, i interface{}, res *OSPolicyAssignment) map[string]OSPolicyAssignmentInstanceFilterInventories {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]OSPolicyAssignmentInstanceFilterInventories{}
	}

	if len(a) == 0 {
		return map[string]OSPolicyAssignmentInstanceFilterInventories{}
	}

	items := make(map[string]OSPolicyAssignmentInstanceFilterInventories)
	for k, item := range a {
		items[k] = *flattenOSPolicyAssignmentInstanceFilterInventories(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenOSPolicyAssignmentInstanceFilterInventoriesSlice flattens the contents of OSPolicyAssignmentInstanceFilterInventories from a JSON
// response object.
func flattenOSPolicyAssignmentInstanceFilterInventoriesSlice(c *Client, i interface{}, res *OSPolicyAssignment) []OSPolicyAssignmentInstanceFilterInventories {
	a, ok := i.([]interface{})
	if !ok {
		return []OSPolicyAssignmentInstanceFilterInventories{}
	}

	if len(a) == 0 {
		return []OSPolicyAssignmentInstanceFilterInventories{}
	}

	items := make([]OSPolicyAssignmentInstanceFilterInventories, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenOSPolicyAssignmentInstanceFilterInventories(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandOSPolicyAssignmentInstanceFilterInventories expands an instance of OSPolicyAssignmentInstanceFilterInventories into a JSON
// request object.
func expandOSPolicyAssignmentInstanceFilterInventories(c *Client, f *OSPolicyAssignmentInstanceFilterInventories, res *OSPolicyAssignment) (map[string]interface{}, error) {
	if f == nil {
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
func flattenOSPolicyAssignmentInstanceFilterInventories(c *Client, i interface{}, res *OSPolicyAssignment) *OSPolicyAssignmentInstanceFilterInventories {
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
func expandOSPolicyAssignmentRolloutMap(c *Client, f map[string]OSPolicyAssignmentRollout, res *OSPolicyAssignment) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandOSPolicyAssignmentRollout(c, &item, res)
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
func expandOSPolicyAssignmentRolloutSlice(c *Client, f []OSPolicyAssignmentRollout, res *OSPolicyAssignment) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandOSPolicyAssignmentRollout(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenOSPolicyAssignmentRolloutMap flattens the contents of OSPolicyAssignmentRollout from a JSON
// response object.
func flattenOSPolicyAssignmentRolloutMap(c *Client, i interface{}, res *OSPolicyAssignment) map[string]OSPolicyAssignmentRollout {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]OSPolicyAssignmentRollout{}
	}

	if len(a) == 0 {
		return map[string]OSPolicyAssignmentRollout{}
	}

	items := make(map[string]OSPolicyAssignmentRollout)
	for k, item := range a {
		items[k] = *flattenOSPolicyAssignmentRollout(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenOSPolicyAssignmentRolloutSlice flattens the contents of OSPolicyAssignmentRollout from a JSON
// response object.
func flattenOSPolicyAssignmentRolloutSlice(c *Client, i interface{}, res *OSPolicyAssignment) []OSPolicyAssignmentRollout {
	a, ok := i.([]interface{})
	if !ok {
		return []OSPolicyAssignmentRollout{}
	}

	if len(a) == 0 {
		return []OSPolicyAssignmentRollout{}
	}

	items := make([]OSPolicyAssignmentRollout, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenOSPolicyAssignmentRollout(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandOSPolicyAssignmentRollout expands an instance of OSPolicyAssignmentRollout into a JSON
// request object.
func expandOSPolicyAssignmentRollout(c *Client, f *OSPolicyAssignmentRollout, res *OSPolicyAssignment) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v, err := expandOSPolicyAssignmentRolloutDisruptionBudget(c, f.DisruptionBudget, res); err != nil {
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
func flattenOSPolicyAssignmentRollout(c *Client, i interface{}, res *OSPolicyAssignment) *OSPolicyAssignmentRollout {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &OSPolicyAssignmentRollout{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyOSPolicyAssignmentRollout
	}
	r.DisruptionBudget = flattenOSPolicyAssignmentRolloutDisruptionBudget(c, m["disruptionBudget"], res)
	r.MinWaitDuration = dcl.FlattenString(m["minWaitDuration"])

	return r
}

// expandOSPolicyAssignmentRolloutDisruptionBudgetMap expands the contents of OSPolicyAssignmentRolloutDisruptionBudget into a JSON
// request object.
func expandOSPolicyAssignmentRolloutDisruptionBudgetMap(c *Client, f map[string]OSPolicyAssignmentRolloutDisruptionBudget, res *OSPolicyAssignment) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandOSPolicyAssignmentRolloutDisruptionBudget(c, &item, res)
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
func expandOSPolicyAssignmentRolloutDisruptionBudgetSlice(c *Client, f []OSPolicyAssignmentRolloutDisruptionBudget, res *OSPolicyAssignment) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandOSPolicyAssignmentRolloutDisruptionBudget(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenOSPolicyAssignmentRolloutDisruptionBudgetMap flattens the contents of OSPolicyAssignmentRolloutDisruptionBudget from a JSON
// response object.
func flattenOSPolicyAssignmentRolloutDisruptionBudgetMap(c *Client, i interface{}, res *OSPolicyAssignment) map[string]OSPolicyAssignmentRolloutDisruptionBudget {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]OSPolicyAssignmentRolloutDisruptionBudget{}
	}

	if len(a) == 0 {
		return map[string]OSPolicyAssignmentRolloutDisruptionBudget{}
	}

	items := make(map[string]OSPolicyAssignmentRolloutDisruptionBudget)
	for k, item := range a {
		items[k] = *flattenOSPolicyAssignmentRolloutDisruptionBudget(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenOSPolicyAssignmentRolloutDisruptionBudgetSlice flattens the contents of OSPolicyAssignmentRolloutDisruptionBudget from a JSON
// response object.
func flattenOSPolicyAssignmentRolloutDisruptionBudgetSlice(c *Client, i interface{}, res *OSPolicyAssignment) []OSPolicyAssignmentRolloutDisruptionBudget {
	a, ok := i.([]interface{})
	if !ok {
		return []OSPolicyAssignmentRolloutDisruptionBudget{}
	}

	if len(a) == 0 {
		return []OSPolicyAssignmentRolloutDisruptionBudget{}
	}

	items := make([]OSPolicyAssignmentRolloutDisruptionBudget, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenOSPolicyAssignmentRolloutDisruptionBudget(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandOSPolicyAssignmentRolloutDisruptionBudget expands an instance of OSPolicyAssignmentRolloutDisruptionBudget into a JSON
// request object.
func expandOSPolicyAssignmentRolloutDisruptionBudget(c *Client, f *OSPolicyAssignmentRolloutDisruptionBudget, res *OSPolicyAssignment) (map[string]interface{}, error) {
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
func flattenOSPolicyAssignmentRolloutDisruptionBudget(c *Client, i interface{}, res *OSPolicyAssignment) *OSPolicyAssignmentRolloutDisruptionBudget {
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
func flattenOSPolicyAssignmentOSPoliciesModeEnumMap(c *Client, i interface{}, res *OSPolicyAssignment) map[string]OSPolicyAssignmentOSPoliciesModeEnum {
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
func flattenOSPolicyAssignmentOSPoliciesModeEnumSlice(c *Client, i interface{}, res *OSPolicyAssignment) []OSPolicyAssignmentOSPoliciesModeEnum {
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
		return nil
	}

	return OSPolicyAssignmentOSPoliciesModeEnumRef(s)
}

// flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDesiredStateEnumMap flattens the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDesiredStateEnum from a JSON
// response object.
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDesiredStateEnumMap(c *Client, i interface{}, res *OSPolicyAssignment) map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDesiredStateEnum {
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
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDesiredStateEnumSlice(c *Client, i interface{}, res *OSPolicyAssignment) []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDesiredStateEnum {
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
		return nil
	}

	return OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDesiredStateEnumRef(s)
}

// flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryAptArchiveTypeEnumMap flattens the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryAptArchiveTypeEnum from a JSON
// response object.
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryAptArchiveTypeEnumMap(c *Client, i interface{}, res *OSPolicyAssignment) map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryAptArchiveTypeEnum {
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
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryAptArchiveTypeEnumSlice(c *Client, i interface{}, res *OSPolicyAssignment) []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryAptArchiveTypeEnum {
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
		return nil
	}

	return OSPolicyAssignmentOSPoliciesResourceGroupsResourcesRepositoryAptArchiveTypeEnumRef(s)
}

// flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateInterpreterEnumMap flattens the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateInterpreterEnum from a JSON
// response object.
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateInterpreterEnumMap(c *Client, i interface{}, res *OSPolicyAssignment) map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateInterpreterEnum {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateInterpreterEnum{}
	}

	if len(a) == 0 {
		return map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateInterpreterEnum{}
	}

	items := make(map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateInterpreterEnum)
	for k, item := range a {
		items[k] = *flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateInterpreterEnum(item.(interface{}))
	}

	return items
}

// flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateInterpreterEnumSlice flattens the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateInterpreterEnum from a JSON
// response object.
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateInterpreterEnumSlice(c *Client, i interface{}, res *OSPolicyAssignment) []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateInterpreterEnum {
	a, ok := i.([]interface{})
	if !ok {
		return []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateInterpreterEnum{}
	}

	if len(a) == 0 {
		return []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateInterpreterEnum{}
	}

	items := make([]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateInterpreterEnum, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateInterpreterEnum(item.(interface{})))
	}

	return items
}

// flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateInterpreterEnum asserts that an interface is a string, and returns a
// pointer to a *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateInterpreterEnum with the same value as that string.
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateInterpreterEnum(i interface{}) *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateInterpreterEnum {
	s, ok := i.(string)
	if !ok {
		return nil
	}

	return OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateInterpreterEnumRef(s)
}

// flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceInterpreterEnumMap flattens the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceInterpreterEnum from a JSON
// response object.
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceInterpreterEnumMap(c *Client, i interface{}, res *OSPolicyAssignment) map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceInterpreterEnum {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceInterpreterEnum{}
	}

	if len(a) == 0 {
		return map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceInterpreterEnum{}
	}

	items := make(map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceInterpreterEnum)
	for k, item := range a {
		items[k] = *flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceInterpreterEnum(item.(interface{}))
	}

	return items
}

// flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceInterpreterEnumSlice flattens the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceInterpreterEnum from a JSON
// response object.
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceInterpreterEnumSlice(c *Client, i interface{}, res *OSPolicyAssignment) []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceInterpreterEnum {
	a, ok := i.([]interface{})
	if !ok {
		return []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceInterpreterEnum{}
	}

	if len(a) == 0 {
		return []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceInterpreterEnum{}
	}

	items := make([]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceInterpreterEnum, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceInterpreterEnum(item.(interface{})))
	}

	return items
}

// flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceInterpreterEnum asserts that an interface is a string, and returns a
// pointer to a *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceInterpreterEnum with the same value as that string.
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceInterpreterEnum(i interface{}) *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceInterpreterEnum {
	s, ok := i.(string)
	if !ok {
		return nil
	}

	return OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceInterpreterEnumRef(s)
}

// flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileStateEnumMap flattens the contents of OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileStateEnum from a JSON
// response object.
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileStateEnumMap(c *Client, i interface{}, res *OSPolicyAssignment) map[string]OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileStateEnum {
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
func flattenOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileStateEnumSlice(c *Client, i interface{}, res *OSPolicyAssignment) []OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileStateEnum {
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
		return nil
	}

	return OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileStateEnumRef(s)
}

// flattenOSPolicyAssignmentRolloutStateEnumMap flattens the contents of OSPolicyAssignmentRolloutStateEnum from a JSON
// response object.
func flattenOSPolicyAssignmentRolloutStateEnumMap(c *Client, i interface{}, res *OSPolicyAssignment) map[string]OSPolicyAssignmentRolloutStateEnum {
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
func flattenOSPolicyAssignmentRolloutStateEnumSlice(c *Client, i interface{}, res *OSPolicyAssignment) []OSPolicyAssignmentRolloutStateEnum {
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
		return nil
	}

	return OSPolicyAssignmentRolloutStateEnumRef(s)
}

// This function returns a matcher that checks whether a serialized resource matches this resource
// in its parameters (as defined by the fields in a Get, which definitionally define resource
// identity).  This is useful in extracting the element from a List call.
func (r *OSPolicyAssignment) matcher(c *Client) func([]byte) bool {
	return func(b []byte) bool {
		cr, err := unmarshalOSPolicyAssignment(b, c, r)
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
	FieldName        string // used for error logging
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
		// Use the first field diff's field name for logging required recreate error.
		diff := oSPolicyAssignmentDiff{FieldName: fieldDiffs[0].FieldName}
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
	if !dcl.IsEmptyValueIndirect(vInstanceFilter) {
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
	if !dcl.IsEmptyValueIndirect(vRollout) {
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
	if !dcl.IsEmptyValueIndirect(vPkg) {
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
	if !dcl.IsEmptyValueIndirect(vRepository) {
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
	if !dcl.IsEmptyValueIndirect(vExec) {
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
	if !dcl.IsEmptyValueIndirect(vFile) {
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
	if !dcl.IsEmptyValueIndirect(vApt) {
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
	if !dcl.IsEmptyValueIndirect(vDeb) {
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
	if !dcl.IsEmptyValueIndirect(vYum) {
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
	if !dcl.IsEmptyValueIndirect(vZypper) {
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
	if !dcl.IsEmptyValueIndirect(vRpm) {
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
	if !dcl.IsEmptyValueIndirect(vGooget) {
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
	if !dcl.IsEmptyValueIndirect(vMsi) {
		o.Msi = vMsi
	}
	return nil
}
func extractOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgAptFields(r *OSPolicyAssignment, o *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgApt) error {
	return nil
}
func extractOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebFields(r *OSPolicyAssignment, o *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDeb) error {
	vSource := o.Source
	if vSource == nil {
		// note: explicitly not the empty object.
		vSource = &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSource{}
	}
	if err := extractOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceFields(r, vSource); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vSource) {
		o.Source = vSource
	}
	return nil
}
func extractOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceFields(r *OSPolicyAssignment, o *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSource) error {
	vRemote := o.Remote
	if vRemote == nil {
		// note: explicitly not the empty object.
		vRemote = &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceRemote{}
	}
	if err := extractOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceRemoteFields(r, vRemote); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vRemote) {
		o.Remote = vRemote
	}
	vGcs := o.Gcs
	if vGcs == nil {
		// note: explicitly not the empty object.
		vGcs = &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceGcs{}
	}
	if err := extractOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceGcsFields(r, vGcs); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vGcs) {
		o.Gcs = vGcs
	}
	return nil
}
func extractOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceRemoteFields(r *OSPolicyAssignment, o *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceRemote) error {
	return nil
}
func extractOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceGcsFields(r *OSPolicyAssignment, o *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceGcs) error {
	return nil
}
func extractOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgYumFields(r *OSPolicyAssignment, o *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgYum) error {
	return nil
}
func extractOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgZypperFields(r *OSPolicyAssignment, o *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgZypper) error {
	return nil
}
func extractOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmFields(r *OSPolicyAssignment, o *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpm) error {
	vSource := o.Source
	if vSource == nil {
		// note: explicitly not the empty object.
		vSource = &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSource{}
	}
	if err := extractOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceFields(r, vSource); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vSource) {
		o.Source = vSource
	}
	return nil
}
func extractOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceFields(r *OSPolicyAssignment, o *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSource) error {
	vRemote := o.Remote
	if vRemote == nil {
		// note: explicitly not the empty object.
		vRemote = &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceRemote{}
	}
	if err := extractOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceRemoteFields(r, vRemote); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vRemote) {
		o.Remote = vRemote
	}
	vGcs := o.Gcs
	if vGcs == nil {
		// note: explicitly not the empty object.
		vGcs = &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceGcs{}
	}
	if err := extractOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceGcsFields(r, vGcs); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vGcs) {
		o.Gcs = vGcs
	}
	return nil
}
func extractOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceRemoteFields(r *OSPolicyAssignment, o *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceRemote) error {
	return nil
}
func extractOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceGcsFields(r *OSPolicyAssignment, o *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceGcs) error {
	return nil
}
func extractOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgGoogetFields(r *OSPolicyAssignment, o *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgGooget) error {
	return nil
}
func extractOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiFields(r *OSPolicyAssignment, o *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsi) error {
	vSource := o.Source
	if vSource == nil {
		// note: explicitly not the empty object.
		vSource = &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSource{}
	}
	if err := extractOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceFields(r, vSource); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vSource) {
		o.Source = vSource
	}
	return nil
}
func extractOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceFields(r *OSPolicyAssignment, o *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSource) error {
	vRemote := o.Remote
	if vRemote == nil {
		// note: explicitly not the empty object.
		vRemote = &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceRemote{}
	}
	if err := extractOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceRemoteFields(r, vRemote); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vRemote) {
		o.Remote = vRemote
	}
	vGcs := o.Gcs
	if vGcs == nil {
		// note: explicitly not the empty object.
		vGcs = &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceGcs{}
	}
	if err := extractOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceGcsFields(r, vGcs); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vGcs) {
		o.Gcs = vGcs
	}
	return nil
}
func extractOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceRemoteFields(r *OSPolicyAssignment, o *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceRemote) error {
	return nil
}
func extractOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceGcsFields(r *OSPolicyAssignment, o *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceGcs) error {
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
	if !dcl.IsEmptyValueIndirect(vApt) {
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
	if !dcl.IsEmptyValueIndirect(vYum) {
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
	if !dcl.IsEmptyValueIndirect(vZypper) {
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
	if !dcl.IsEmptyValueIndirect(vGoo) {
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
	vValidate := o.Validate
	if vValidate == nil {
		// note: explicitly not the empty object.
		vValidate = &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidate{}
	}
	if err := extractOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFields(r, vValidate); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vValidate) {
		o.Validate = vValidate
	}
	vEnforce := o.Enforce
	if vEnforce == nil {
		// note: explicitly not the empty object.
		vEnforce = &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforce{}
	}
	if err := extractOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFields(r, vEnforce); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vEnforce) {
		o.Enforce = vEnforce
	}
	return nil
}
func extractOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFields(r *OSPolicyAssignment, o *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidate) error {
	vFile := o.File
	if vFile == nil {
		// note: explicitly not the empty object.
		vFile = &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFile{}
	}
	if err := extractOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileFields(r, vFile); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vFile) {
		o.File = vFile
	}
	return nil
}
func extractOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileFields(r *OSPolicyAssignment, o *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFile) error {
	vRemote := o.Remote
	if vRemote == nil {
		// note: explicitly not the empty object.
		vRemote = &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileRemote{}
	}
	if err := extractOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileRemoteFields(r, vRemote); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vRemote) {
		o.Remote = vRemote
	}
	vGcs := o.Gcs
	if vGcs == nil {
		// note: explicitly not the empty object.
		vGcs = &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileGcs{}
	}
	if err := extractOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileGcsFields(r, vGcs); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vGcs) {
		o.Gcs = vGcs
	}
	return nil
}
func extractOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileRemoteFields(r *OSPolicyAssignment, o *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileRemote) error {
	return nil
}
func extractOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileGcsFields(r *OSPolicyAssignment, o *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileGcs) error {
	return nil
}
func extractOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFields(r *OSPolicyAssignment, o *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforce) error {
	vFile := o.File
	if vFile == nil {
		// note: explicitly not the empty object.
		vFile = &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFile{}
	}
	if err := extractOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileFields(r, vFile); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vFile) {
		o.File = vFile
	}
	return nil
}
func extractOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileFields(r *OSPolicyAssignment, o *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFile) error {
	vRemote := o.Remote
	if vRemote == nil {
		// note: explicitly not the empty object.
		vRemote = &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileRemote{}
	}
	if err := extractOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileRemoteFields(r, vRemote); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vRemote) {
		o.Remote = vRemote
	}
	vGcs := o.Gcs
	if vGcs == nil {
		// note: explicitly not the empty object.
		vGcs = &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileGcs{}
	}
	if err := extractOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileGcsFields(r, vGcs); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vGcs) {
		o.Gcs = vGcs
	}
	return nil
}
func extractOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileRemoteFields(r *OSPolicyAssignment, o *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileRemote) error {
	return nil
}
func extractOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileGcsFields(r *OSPolicyAssignment, o *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileGcs) error {
	return nil
}
func extractOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFields(r *OSPolicyAssignment, o *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFile) error {
	vFile := o.File
	if vFile == nil {
		// note: explicitly not the empty object.
		vFile = &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFile{}
	}
	if err := extractOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileFields(r, vFile); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vFile) {
		o.File = vFile
	}
	return nil
}
func extractOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileFields(r *OSPolicyAssignment, o *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFile) error {
	vRemote := o.Remote
	if vRemote == nil {
		// note: explicitly not the empty object.
		vRemote = &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileRemote{}
	}
	if err := extractOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileRemoteFields(r, vRemote); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vRemote) {
		o.Remote = vRemote
	}
	vGcs := o.Gcs
	if vGcs == nil {
		// note: explicitly not the empty object.
		vGcs = &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileGcs{}
	}
	if err := extractOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileGcsFields(r, vGcs); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vGcs) {
		o.Gcs = vGcs
	}
	return nil
}
func extractOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileRemoteFields(r *OSPolicyAssignment, o *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileRemote) error {
	return nil
}
func extractOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileGcsFields(r *OSPolicyAssignment, o *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileGcs) error {
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
	if !dcl.IsEmptyValueIndirect(vDisruptionBudget) {
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
	if !dcl.IsEmptyValueIndirect(vInstanceFilter) {
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
	if !dcl.IsEmptyValueIndirect(vRollout) {
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
	if !dcl.IsEmptyValueIndirect(vPkg) {
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
	if !dcl.IsEmptyValueIndirect(vRepository) {
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
	if !dcl.IsEmptyValueIndirect(vExec) {
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
	if !dcl.IsEmptyValueIndirect(vFile) {
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
	if !dcl.IsEmptyValueIndirect(vApt) {
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
	if !dcl.IsEmptyValueIndirect(vDeb) {
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
	if !dcl.IsEmptyValueIndirect(vYum) {
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
	if !dcl.IsEmptyValueIndirect(vZypper) {
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
	if !dcl.IsEmptyValueIndirect(vRpm) {
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
	if !dcl.IsEmptyValueIndirect(vGooget) {
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
	if !dcl.IsEmptyValueIndirect(vMsi) {
		o.Msi = vMsi
	}
	return nil
}
func postReadExtractOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgAptFields(r *OSPolicyAssignment, o *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgApt) error {
	return nil
}
func postReadExtractOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebFields(r *OSPolicyAssignment, o *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDeb) error {
	vSource := o.Source
	if vSource == nil {
		// note: explicitly not the empty object.
		vSource = &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSource{}
	}
	if err := extractOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceFields(r, vSource); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vSource) {
		o.Source = vSource
	}
	return nil
}
func postReadExtractOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceFields(r *OSPolicyAssignment, o *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSource) error {
	vRemote := o.Remote
	if vRemote == nil {
		// note: explicitly not the empty object.
		vRemote = &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceRemote{}
	}
	if err := extractOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceRemoteFields(r, vRemote); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vRemote) {
		o.Remote = vRemote
	}
	vGcs := o.Gcs
	if vGcs == nil {
		// note: explicitly not the empty object.
		vGcs = &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceGcs{}
	}
	if err := extractOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceGcsFields(r, vGcs); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vGcs) {
		o.Gcs = vGcs
	}
	return nil
}
func postReadExtractOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceRemoteFields(r *OSPolicyAssignment, o *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceRemote) error {
	return nil
}
func postReadExtractOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceGcsFields(r *OSPolicyAssignment, o *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgDebSourceGcs) error {
	return nil
}
func postReadExtractOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgYumFields(r *OSPolicyAssignment, o *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgYum) error {
	return nil
}
func postReadExtractOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgZypperFields(r *OSPolicyAssignment, o *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgZypper) error {
	return nil
}
func postReadExtractOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmFields(r *OSPolicyAssignment, o *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpm) error {
	vSource := o.Source
	if vSource == nil {
		// note: explicitly not the empty object.
		vSource = &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSource{}
	}
	if err := extractOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceFields(r, vSource); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vSource) {
		o.Source = vSource
	}
	return nil
}
func postReadExtractOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceFields(r *OSPolicyAssignment, o *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSource) error {
	vRemote := o.Remote
	if vRemote == nil {
		// note: explicitly not the empty object.
		vRemote = &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceRemote{}
	}
	if err := extractOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceRemoteFields(r, vRemote); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vRemote) {
		o.Remote = vRemote
	}
	vGcs := o.Gcs
	if vGcs == nil {
		// note: explicitly not the empty object.
		vGcs = &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceGcs{}
	}
	if err := extractOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceGcsFields(r, vGcs); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vGcs) {
		o.Gcs = vGcs
	}
	return nil
}
func postReadExtractOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceRemoteFields(r *OSPolicyAssignment, o *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceRemote) error {
	return nil
}
func postReadExtractOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceGcsFields(r *OSPolicyAssignment, o *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgRpmSourceGcs) error {
	return nil
}
func postReadExtractOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgGoogetFields(r *OSPolicyAssignment, o *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgGooget) error {
	return nil
}
func postReadExtractOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiFields(r *OSPolicyAssignment, o *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsi) error {
	vSource := o.Source
	if vSource == nil {
		// note: explicitly not the empty object.
		vSource = &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSource{}
	}
	if err := extractOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceFields(r, vSource); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vSource) {
		o.Source = vSource
	}
	return nil
}
func postReadExtractOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceFields(r *OSPolicyAssignment, o *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSource) error {
	vRemote := o.Remote
	if vRemote == nil {
		// note: explicitly not the empty object.
		vRemote = &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceRemote{}
	}
	if err := extractOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceRemoteFields(r, vRemote); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vRemote) {
		o.Remote = vRemote
	}
	vGcs := o.Gcs
	if vGcs == nil {
		// note: explicitly not the empty object.
		vGcs = &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceGcs{}
	}
	if err := extractOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceGcsFields(r, vGcs); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vGcs) {
		o.Gcs = vGcs
	}
	return nil
}
func postReadExtractOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceRemoteFields(r *OSPolicyAssignment, o *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceRemote) error {
	return nil
}
func postReadExtractOSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceGcsFields(r *OSPolicyAssignment, o *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesPkgMsiSourceGcs) error {
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
	if !dcl.IsEmptyValueIndirect(vApt) {
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
	if !dcl.IsEmptyValueIndirect(vYum) {
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
	if !dcl.IsEmptyValueIndirect(vZypper) {
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
	if !dcl.IsEmptyValueIndirect(vGoo) {
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
	vValidate := o.Validate
	if vValidate == nil {
		// note: explicitly not the empty object.
		vValidate = &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidate{}
	}
	if err := extractOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFields(r, vValidate); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vValidate) {
		o.Validate = vValidate
	}
	vEnforce := o.Enforce
	if vEnforce == nil {
		// note: explicitly not the empty object.
		vEnforce = &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforce{}
	}
	if err := extractOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFields(r, vEnforce); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vEnforce) {
		o.Enforce = vEnforce
	}
	return nil
}
func postReadExtractOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFields(r *OSPolicyAssignment, o *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidate) error {
	vFile := o.File
	if vFile == nil {
		// note: explicitly not the empty object.
		vFile = &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFile{}
	}
	if err := extractOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileFields(r, vFile); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vFile) {
		o.File = vFile
	}
	return nil
}
func postReadExtractOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileFields(r *OSPolicyAssignment, o *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFile) error {
	vRemote := o.Remote
	if vRemote == nil {
		// note: explicitly not the empty object.
		vRemote = &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileRemote{}
	}
	if err := extractOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileRemoteFields(r, vRemote); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vRemote) {
		o.Remote = vRemote
	}
	vGcs := o.Gcs
	if vGcs == nil {
		// note: explicitly not the empty object.
		vGcs = &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileGcs{}
	}
	if err := extractOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileGcsFields(r, vGcs); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vGcs) {
		o.Gcs = vGcs
	}
	return nil
}
func postReadExtractOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileRemoteFields(r *OSPolicyAssignment, o *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileRemote) error {
	return nil
}
func postReadExtractOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileGcsFields(r *OSPolicyAssignment, o *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecValidateFileGcs) error {
	return nil
}
func postReadExtractOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFields(r *OSPolicyAssignment, o *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforce) error {
	vFile := o.File
	if vFile == nil {
		// note: explicitly not the empty object.
		vFile = &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFile{}
	}
	if err := extractOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileFields(r, vFile); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vFile) {
		o.File = vFile
	}
	return nil
}
func postReadExtractOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileFields(r *OSPolicyAssignment, o *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFile) error {
	vRemote := o.Remote
	if vRemote == nil {
		// note: explicitly not the empty object.
		vRemote = &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileRemote{}
	}
	if err := extractOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileRemoteFields(r, vRemote); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vRemote) {
		o.Remote = vRemote
	}
	vGcs := o.Gcs
	if vGcs == nil {
		// note: explicitly not the empty object.
		vGcs = &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileGcs{}
	}
	if err := extractOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileGcsFields(r, vGcs); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vGcs) {
		o.Gcs = vGcs
	}
	return nil
}
func postReadExtractOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileRemoteFields(r *OSPolicyAssignment, o *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileRemote) error {
	return nil
}
func postReadExtractOSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileGcsFields(r *OSPolicyAssignment, o *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesExecEnforceFileGcs) error {
	return nil
}
func postReadExtractOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFields(r *OSPolicyAssignment, o *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFile) error {
	vFile := o.File
	if vFile == nil {
		// note: explicitly not the empty object.
		vFile = &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFile{}
	}
	if err := extractOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileFields(r, vFile); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vFile) {
		o.File = vFile
	}
	return nil
}
func postReadExtractOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileFields(r *OSPolicyAssignment, o *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFile) error {
	vRemote := o.Remote
	if vRemote == nil {
		// note: explicitly not the empty object.
		vRemote = &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileRemote{}
	}
	if err := extractOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileRemoteFields(r, vRemote); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vRemote) {
		o.Remote = vRemote
	}
	vGcs := o.Gcs
	if vGcs == nil {
		// note: explicitly not the empty object.
		vGcs = &OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileGcs{}
	}
	if err := extractOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileGcsFields(r, vGcs); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vGcs) {
		o.Gcs = vGcs
	}
	return nil
}
func postReadExtractOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileRemoteFields(r *OSPolicyAssignment, o *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileRemote) error {
	return nil
}
func postReadExtractOSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileGcsFields(r *OSPolicyAssignment, o *OSPolicyAssignmentOSPoliciesResourceGroupsResourcesFileFileGcs) error {
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
	if !dcl.IsEmptyValueIndirect(vDisruptionBudget) {
		o.DisruptionBudget = vDisruptionBudget
	}
	return nil
}
func postReadExtractOSPolicyAssignmentRolloutDisruptionBudgetFields(r *OSPolicyAssignment, o *OSPolicyAssignmentRolloutDisruptionBudget) error {
	return nil
}
