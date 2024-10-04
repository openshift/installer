// Copyright 2024 Google LLC. All Rights Reserved.
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

func (r *Instance) validate() error {

	if err := dcl.Required(r, "zone"); err != nil {
		return err
	}
	if err := dcl.RequiredParameter(r.Project, "Project"); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(r.Scheduling) {
		if err := r.Scheduling.validate(); err != nil {
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
func (r *InstanceDisks) validate() error {
	if !dcl.IsEmptyValueIndirect(r.DiskEncryptionKey) {
		if err := r.DiskEncryptionKey.validate(); err != nil {
			return err
		}
	}
	if !dcl.IsEmptyValueIndirect(r.InitializeParams) {
		if err := r.InitializeParams.validate(); err != nil {
			return err
		}
	}
	return nil
}
func (r *InstanceDisksDiskEncryptionKey) validate() error {
	return nil
}
func (r *InstanceDisksInitializeParams) validate() error {
	if !dcl.IsEmptyValueIndirect(r.SourceImageEncryptionKey) {
		if err := r.SourceImageEncryptionKey.validate(); err != nil {
			return err
		}
	}
	return nil
}
func (r *InstanceDisksInitializeParamsSourceImageEncryptionKey) validate() error {
	return nil
}
func (r *InstanceGuestAccelerators) validate() error {
	return nil
}
func (r *InstanceNetworkInterfaces) validate() error {
	return nil
}
func (r *InstanceNetworkInterfacesAccessConfigs) validate() error {
	if err := dcl.Required(r, "name"); err != nil {
		return err
	}
	if err := dcl.Required(r, "type"); err != nil {
		return err
	}
	return nil
}
func (r *InstanceNetworkInterfacesIPv6AccessConfigs) validate() error {
	if err := dcl.Required(r, "name"); err != nil {
		return err
	}
	if err := dcl.Required(r, "type"); err != nil {
		return err
	}
	return nil
}
func (r *InstanceNetworkInterfacesAliasIPRanges) validate() error {
	return nil
}
func (r *InstanceScheduling) validate() error {
	return nil
}
func (r *InstanceServiceAccounts) validate() error {
	return nil
}
func (r *InstanceShieldedInstanceConfig) validate() error {
	return nil
}
func (r *Instance) basePath() string {
	params := map[string]interface{}{}
	return dcl.Nprintf("https://www.googleapis.com/compute/v1/", params)
}

func (r *Instance) getURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"project": dcl.ValueOrEmptyString(nr.Project),
		"zone":    dcl.ValueOrEmptyString(nr.Zone),
		"name":    dcl.ValueOrEmptyString(nr.Name),
	}
	return dcl.URL("projects/{{project}}/zones/{{zone}}/instances/{{name}}", nr.basePath(), userBasePath, params), nil
}

func (r *Instance) listURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"project": dcl.ValueOrEmptyString(nr.Project),
		"zone":    dcl.ValueOrEmptyString(nr.Zone),
	}
	return dcl.URL("projects/{{project}}/zones/{{zone}}/instances", nr.basePath(), userBasePath, params), nil

}

func (r *Instance) createURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"project": dcl.ValueOrEmptyString(nr.Project),
		"zone":    dcl.ValueOrEmptyString(nr.Zone),
	}
	return dcl.URL("projects/{{project}}/zones/{{zone}}/instances", nr.basePath(), userBasePath, params), nil

}

func (r *Instance) deleteURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"project": dcl.ValueOrEmptyString(nr.Project),
		"zone":    dcl.ValueOrEmptyString(nr.Zone),
		"name":    dcl.ValueOrEmptyString(nr.Name),
	}
	return dcl.URL("projects/{{project}}/zones/{{zone}}/instances/{{name}}", nr.basePath(), userBasePath, params), nil
}

func (r *Instance) SetPolicyURL(userBasePath string) string {
	nr := r.urlNormalized()
	fields := map[string]interface{}{
		"project": *nr.Project,
		"zone":    *nr.Zone,
		"name":    *nr.Name,
	}
	return dcl.URL("projects/{{project}}/zones/{{zone}}/instances/{{name}}/setIamPolicy", nr.basePath(), userBasePath, fields)
}

func (r *Instance) SetPolicyVerb() string {
	return "POST"
}

func (r *Instance) getPolicyURL(userBasePath string) string {
	nr := r.urlNormalized()
	fields := map[string]interface{}{
		"project": *nr.Project,
		"zone":    *nr.Zone,
		"name":    *nr.Name,
	}
	return dcl.URL("projects/{{project}}/zones/{{zone}}/instances/{{name}}/getIamPolicy", nr.basePath(), userBasePath, fields)
}

func (r *Instance) IAMPolicyVersion() int {
	return 3
}

// instanceApiOperation represents a mutable operation in the underlying REST
// API such as Create, Update, or Delete.
type instanceApiOperation interface {
	do(context.Context, *Instance, *Client) error
}

// newUpdateInstanceSetDeletionProtectionRequest creates a request for an
// Instance resource's setDeletionProtection update type by filling in the update
// fields based on the intended state of the resource.
func newUpdateInstanceSetDeletionProtectionRequest(ctx context.Context, f *Instance, c *Client) (map[string]interface{}, error) {
	req := map[string]interface{}{}
	res := f
	_ = res

	if v := f.DeletionProtection; !dcl.IsEmptyValueIndirect(v) {
		req["deletionProtection"] = v
	}
	return req, nil
}

// marshalUpdateInstanceSetDeletionProtectionRequest converts the update into
// the final JSON request body.
func marshalUpdateInstanceSetDeletionProtectionRequest(c *Client, m map[string]interface{}) ([]byte, error) {

	dcl.MoveMapEntry(
		m,
		[]string{"tags"},
		[]string{"tags", "items"},
	)
	return json.Marshal(m)
}

type updateInstanceSetDeletionProtectionOperation struct {
	// If the update operation has the REQUIRES_APPLY_OPTIONS trait, this will be populated.
	// Usually it will be nil - this is to prevent us from accidentally depending on apply
	// options, which should usually be unnecessary.
	ApplyOptions []dcl.ApplyOption
	FieldDiffs   []*dcl.FieldDiff
}

// do creates a request and sends it to the appropriate URL. In most operations,
// do will transcribe a subset of the resource into a request object and send a
// PUT request to a single URL.

func (op *updateInstanceSetDeletionProtectionOperation) do(ctx context.Context, r *Instance, c *Client) error {
	_, err := c.GetInstance(ctx, r)
	if err != nil {
		return err
	}

	u, err := r.updateURL(c.Config.BasePath, "setDeletionProtection")
	if err != nil {
		return err
	}

	req, err := newUpdateInstanceSetDeletionProtectionRequest(ctx, r, c)
	if err != nil {
		return err
	}

	c.Config.Logger.InfoWithContextf(ctx, "Created update: %#v", req)
	body, err := marshalUpdateInstanceSetDeletionProtectionRequest(c, req)
	if err != nil {
		return err
	}
	resp, err := dcl.SendRequest(ctx, c.Config, "POST", u, bytes.NewBuffer(body), c.Config.RetryProvider)
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

// newUpdateInstanceSetLabelsRequest creates a request for an
// Instance resource's setLabels update type by filling in the update
// fields based on the intended state of the resource.
func newUpdateInstanceSetLabelsRequest(ctx context.Context, f *Instance, c *Client) (map[string]interface{}, error) {
	req := map[string]interface{}{}
	res := f
	_ = res

	if v := f.Labels; !dcl.IsEmptyValueIndirect(v) {
		req["labels"] = v
	}
	b, err := c.getInstanceRaw(ctx, f)
	if err != nil {
		return nil, err
	}
	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, err
	}
	rawLabelFingerprint, err := dcl.GetMapEntry(
		m,
		[]string{"labelFingerprint"},
	)
	if err != nil {
		c.Config.Logger.WarningWithContextf(ctx, "Failed to fetch from JSON Path: %v", err)
	} else {
		req["labelFingerprint"] = rawLabelFingerprint.(string)
	}
	return req, nil
}

// marshalUpdateInstanceSetLabelsRequest converts the update into
// the final JSON request body.
func marshalUpdateInstanceSetLabelsRequest(c *Client, m map[string]interface{}) ([]byte, error) {

	dcl.MoveMapEntry(
		m,
		[]string{"tags"},
		[]string{"tags", "items"},
	)
	return json.Marshal(m)
}

type updateInstanceSetLabelsOperation struct {
	// If the update operation has the REQUIRES_APPLY_OPTIONS trait, this will be populated.
	// Usually it will be nil - this is to prevent us from accidentally depending on apply
	// options, which should usually be unnecessary.
	ApplyOptions []dcl.ApplyOption
	FieldDiffs   []*dcl.FieldDiff
}

// do creates a request and sends it to the appropriate URL. In most operations,
// do will transcribe a subset of the resource into a request object and send a
// PUT request to a single URL.

func (op *updateInstanceSetLabelsOperation) do(ctx context.Context, r *Instance, c *Client) error {
	_, err := c.GetInstance(ctx, r)
	if err != nil {
		return err
	}

	u, err := r.updateURL(c.Config.BasePath, "setLabels")
	if err != nil {
		return err
	}

	req, err := newUpdateInstanceSetLabelsRequest(ctx, r, c)
	if err != nil {
		return err
	}

	c.Config.Logger.InfoWithContextf(ctx, "Created update: %#v", req)
	body, err := marshalUpdateInstanceSetLabelsRequest(c, req)
	if err != nil {
		return err
	}
	resp, err := dcl.SendRequest(ctx, c.Config, "POST", u, bytes.NewBuffer(body), c.Config.RetryProvider)
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

// newUpdateInstanceSetMachineTypeRequest creates a request for an
// Instance resource's setMachineType update type by filling in the update
// fields based on the intended state of the resource.
func newUpdateInstanceSetMachineTypeRequest(ctx context.Context, f *Instance, c *Client) (map[string]interface{}, error) {
	req := map[string]interface{}{}
	res := f
	_ = res

	if v, err := dcl.DeriveFromPattern("zones/%s/machineTypes/%s", f.MachineType, dcl.SelfLinkToName(f.Zone), f.MachineType); err != nil {
		return nil, fmt.Errorf("error expanding MachineType into machineType: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		req["machineType"] = v
	}
	return req, nil
}

// marshalUpdateInstanceSetMachineTypeRequest converts the update into
// the final JSON request body.
func marshalUpdateInstanceSetMachineTypeRequest(c *Client, m map[string]interface{}) ([]byte, error) {

	dcl.MoveMapEntry(
		m,
		[]string{"tags"},
		[]string{"tags", "items"},
	)
	return json.Marshal(m)
}

type updateInstanceSetMachineTypeOperation struct {
	// If the update operation has the REQUIRES_APPLY_OPTIONS trait, this will be populated.
	// Usually it will be nil - this is to prevent us from accidentally depending on apply
	// options, which should usually be unnecessary.
	ApplyOptions []dcl.ApplyOption
	FieldDiffs   []*dcl.FieldDiff
}

// do creates a request and sends it to the appropriate URL. In most operations,
// do will transcribe a subset of the resource into a request object and send a
// PUT request to a single URL.

func (op *updateInstanceSetMachineTypeOperation) do(ctx context.Context, r *Instance, c *Client) error {
	_, err := c.GetInstance(ctx, r)
	if err != nil {
		return err
	}

	u, err := r.updateURL(c.Config.BasePath, "setMachineType")
	if err != nil {
		return err
	}

	req, err := newUpdateInstanceSetMachineTypeRequest(ctx, r, c)
	if err != nil {
		return err
	}

	c.Config.Logger.InfoWithContextf(ctx, "Created update: %#v", req)
	body, err := marshalUpdateInstanceSetMachineTypeRequest(c, req)
	if err != nil {
		return err
	}
	resp, err := dcl.SendRequest(ctx, c.Config, "POST", u, bytes.NewBuffer(body), c.Config.RetryProvider)
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

// newUpdateInstanceSetMetadataRequest creates a request for an
// Instance resource's setMetadata update type by filling in the update
// fields based on the intended state of the resource.
func newUpdateInstanceSetMetadataRequest(ctx context.Context, f *Instance, c *Client) (map[string]interface{}, error) {
	req := map[string]interface{}{}
	res := f
	_ = res

	if v, err := dcl.ListOfKeyValuesFromMap(f.Metadata, "key", "value"); err != nil {
		return nil, fmt.Errorf("error expanding Metadata into items: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		req["items"] = v
	}
	b, err := c.getInstanceRaw(ctx, f)
	if err != nil {
		return nil, err
	}
	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, err
	}
	rawFingerprint, err := dcl.GetMapEntry(
		m,
		[]string{"metadata", "fingerprint"},
	)
	if err != nil {
		c.Config.Logger.WarningWithContextf(ctx, "Failed to fetch from JSON Path: %v", err)
	} else {
		req["fingerprint"] = rawFingerprint.(string)
	}
	return req, nil
}

// marshalUpdateInstanceSetMetadataRequest converts the update into
// the final JSON request body.
func marshalUpdateInstanceSetMetadataRequest(c *Client, m map[string]interface{}) ([]byte, error) {

	dcl.MoveMapEntry(
		m,
		[]string{"tags"},
		[]string{"tags", "items"},
	)
	return json.Marshal(m)
}

type updateInstanceSetMetadataOperation struct {
	// If the update operation has the REQUIRES_APPLY_OPTIONS trait, this will be populated.
	// Usually it will be nil - this is to prevent us from accidentally depending on apply
	// options, which should usually be unnecessary.
	ApplyOptions []dcl.ApplyOption
	FieldDiffs   []*dcl.FieldDiff
}

// do creates a request and sends it to the appropriate URL. In most operations,
// do will transcribe a subset of the resource into a request object and send a
// PUT request to a single URL.

func (op *updateInstanceSetMetadataOperation) do(ctx context.Context, r *Instance, c *Client) error {
	_, err := c.GetInstance(ctx, r)
	if err != nil {
		return err
	}

	u, err := r.updateURL(c.Config.BasePath, "setMetadata")
	if err != nil {
		return err
	}

	req, err := newUpdateInstanceSetMetadataRequest(ctx, r, c)
	if err != nil {
		return err
	}

	c.Config.Logger.InfoWithContextf(ctx, "Created update: %#v", req)
	body, err := marshalUpdateInstanceSetMetadataRequest(c, req)
	if err != nil {
		return err
	}
	resp, err := dcl.SendRequest(ctx, c.Config, "POST", u, bytes.NewBuffer(body), c.Config.RetryProvider)
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

// newUpdateInstanceSetTagsRequest creates a request for an
// Instance resource's setTags update type by filling in the update
// fields based on the intended state of the resource.
func newUpdateInstanceSetTagsRequest(ctx context.Context, f *Instance, c *Client) (map[string]interface{}, error) {
	req := map[string]interface{}{}
	res := f
	_ = res

	if v := f.Tags; v != nil {
		req["items"] = v
	}
	b, err := c.getInstanceRaw(ctx, f)
	if err != nil {
		return nil, err
	}
	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, err
	}
	rawFingerprint, err := dcl.GetMapEntry(
		m,
		[]string{"tags", "fingerprint"},
	)
	if err != nil {
		c.Config.Logger.WarningWithContextf(ctx, "Failed to fetch from JSON Path: %v", err)
	} else {
		req["fingerprint"] = rawFingerprint.(string)
	}
	return req, nil
}

// marshalUpdateInstanceSetTagsRequest converts the update into
// the final JSON request body.
func marshalUpdateInstanceSetTagsRequest(c *Client, m map[string]interface{}) ([]byte, error) {

	dcl.MoveMapEntry(
		m,
		[]string{"tags"},
		[]string{"tags", "items"},
	)
	return json.Marshal(m)
}

type updateInstanceSetTagsOperation struct {
	// If the update operation has the REQUIRES_APPLY_OPTIONS trait, this will be populated.
	// Usually it will be nil - this is to prevent us from accidentally depending on apply
	// options, which should usually be unnecessary.
	ApplyOptions []dcl.ApplyOption
	FieldDiffs   []*dcl.FieldDiff
}

// do creates a request and sends it to the appropriate URL. In most operations,
// do will transcribe a subset of the resource into a request object and send a
// PUT request to a single URL.

func (op *updateInstanceSetTagsOperation) do(ctx context.Context, r *Instance, c *Client) error {
	_, err := c.GetInstance(ctx, r)
	if err != nil {
		return err
	}

	u, err := r.updateURL(c.Config.BasePath, "setTags")
	if err != nil {
		return err
	}

	req, err := newUpdateInstanceSetTagsRequest(ctx, r, c)
	if err != nil {
		return err
	}

	c.Config.Logger.InfoWithContextf(ctx, "Created update: %#v", req)
	body, err := marshalUpdateInstanceSetTagsRequest(c, req)
	if err != nil {
		return err
	}
	resp, err := dcl.SendRequest(ctx, c.Config, "POST", u, bytes.NewBuffer(body), c.Config.RetryProvider)
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

// newUpdateInstanceStartRequest creates a request for an
// Instance resource's start update type by filling in the update
// fields based on the intended state of the resource.
func newUpdateInstanceStartRequest(ctx context.Context, f *Instance, c *Client) (map[string]interface{}, error) {
	req := map[string]interface{}{}
	res := f
	_ = res

	return req, nil
}

// marshalUpdateInstanceStartRequest converts the update into
// the final JSON request body.
func marshalUpdateInstanceStartRequest(c *Client, m map[string]interface{}) ([]byte, error) {

	dcl.MoveMapEntry(
		m,
		[]string{"tags"},
		[]string{"tags", "items"},
	)
	return json.Marshal(m)
}

type updateInstanceStartOperation struct {
	// If the update operation has the REQUIRES_APPLY_OPTIONS trait, this will be populated.
	// Usually it will be nil - this is to prevent us from accidentally depending on apply
	// options, which should usually be unnecessary.
	ApplyOptions []dcl.ApplyOption
	FieldDiffs   []*dcl.FieldDiff
}

// do creates a request and sends it to the appropriate URL. In most operations,
// do will transcribe a subset of the resource into a request object and send a
// PUT request to a single URL.

func (op *updateInstanceStartOperation) do(ctx context.Context, r *Instance, c *Client) error {
	_, err := c.GetInstance(ctx, r)
	if err != nil {
		return err
	}

	u, err := r.updateURL(c.Config.BasePath, "start")
	if err != nil {
		return err
	}

	req, err := newUpdateInstanceStartRequest(ctx, r, c)
	if err != nil {
		return err
	}

	c.Config.Logger.InfoWithContextf(ctx, "Created update: %#v", req)
	body, err := marshalUpdateInstanceStartRequest(c, req)
	if err != nil {
		return err
	}
	resp, err := dcl.SendRequest(ctx, c.Config, "POST", u, bytes.NewBuffer(body), c.Config.RetryProvider)
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

// newUpdateInstanceStopRequest creates a request for an
// Instance resource's stop update type by filling in the update
// fields based on the intended state of the resource.
func newUpdateInstanceStopRequest(ctx context.Context, f *Instance, c *Client) (map[string]interface{}, error) {
	req := map[string]interface{}{}
	res := f
	_ = res

	return req, nil
}

// marshalUpdateInstanceStopRequest converts the update into
// the final JSON request body.
func marshalUpdateInstanceStopRequest(c *Client, m map[string]interface{}) ([]byte, error) {

	dcl.MoveMapEntry(
		m,
		[]string{"tags"},
		[]string{"tags", "items"},
	)
	return json.Marshal(m)
}

type updateInstanceStopOperation struct {
	// If the update operation has the REQUIRES_APPLY_OPTIONS trait, this will be populated.
	// Usually it will be nil - this is to prevent us from accidentally depending on apply
	// options, which should usually be unnecessary.
	ApplyOptions []dcl.ApplyOption
	FieldDiffs   []*dcl.FieldDiff
}

// do creates a request and sends it to the appropriate URL. In most operations,
// do will transcribe a subset of the resource into a request object and send a
// PUT request to a single URL.

func (op *updateInstanceStopOperation) do(ctx context.Context, r *Instance, c *Client) error {
	_, err := c.GetInstance(ctx, r)
	if err != nil {
		return err
	}

	u, err := r.updateURL(c.Config.BasePath, "stop")
	if err != nil {
		return err
	}

	req, err := newUpdateInstanceStopRequest(ctx, r, c)
	if err != nil {
		return err
	}

	c.Config.Logger.InfoWithContextf(ctx, "Created update: %#v", req)
	body, err := marshalUpdateInstanceStopRequest(c, req)
	if err != nil {
		return err
	}
	resp, err := dcl.SendRequest(ctx, c.Config, "POST", u, bytes.NewBuffer(body), c.Config.RetryProvider)
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

// newUpdateInstanceUpdateRequest creates a request for an
// Instance resource's update update type by filling in the update
// fields based on the intended state of the resource.
func newUpdateInstanceUpdateRequest(ctx context.Context, f *Instance, c *Client) (map[string]interface{}, error) {
	req := map[string]interface{}{}
	res := f
	_ = res

	return req, nil
}

// marshalUpdateInstanceUpdateRequest converts the update into
// the final JSON request body.
func marshalUpdateInstanceUpdateRequest(c *Client, m map[string]interface{}) ([]byte, error) {

	dcl.MoveMapEntry(
		m,
		[]string{"tags"},
		[]string{"tags", "items"},
	)
	return json.Marshal(m)
}

type updateInstanceUpdateOperation struct {
	// If the update operation has the REQUIRES_APPLY_OPTIONS trait, this will be populated.
	// Usually it will be nil - this is to prevent us from accidentally depending on apply
	// options, which should usually be unnecessary.
	ApplyOptions []dcl.ApplyOption
	FieldDiffs   []*dcl.FieldDiff
}

// do creates a request and sends it to the appropriate URL. In most operations,
// do will transcribe a subset of the resource into a request object and send a
// PUT request to a single URL.

func (op *updateInstanceUpdateOperation) do(ctx context.Context, r *Instance, c *Client) error {
	_, err := c.GetInstance(ctx, r)
	if err != nil {
		return err
	}

	u, err := r.updateURL(c.Config.BasePath, "update")
	if err != nil {
		return err
	}

	req, err := newUpdateInstanceUpdateRequest(ctx, r, c)
	if err != nil {
		return err
	}

	c.Config.Logger.InfoWithContextf(ctx, "Created update: %#v", req)
	body, err := marshalUpdateInstanceUpdateRequest(c, req)
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

// newUpdateInstanceUpdateShieldedInstanceConfigRequest creates a request for an
// Instance resource's updateShieldedInstanceConfig update type by filling in the update
// fields based on the intended state of the resource.
func newUpdateInstanceUpdateShieldedInstanceConfigRequest(ctx context.Context, f *Instance, c *Client) (map[string]interface{}, error) {
	req := map[string]interface{}{}
	res := f
	_ = res

	if v, err := expandInstanceShieldedInstanceConfig(c, f.ShieldedInstanceConfig, res); err != nil {
		return nil, fmt.Errorf("error expanding ShieldedInstanceConfig into shieldedInstanceConfig: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		req["shieldedInstanceConfig"] = v
	}
	return req, nil
}

// marshalUpdateInstanceUpdateShieldedInstanceConfigRequest converts the update into
// the final JSON request body.
func marshalUpdateInstanceUpdateShieldedInstanceConfigRequest(c *Client, m map[string]interface{}) ([]byte, error) {

	dcl.MoveMapEntry(
		m,
		[]string{"tags"},
		[]string{"tags", "items"},
	)
	return json.Marshal(m)
}

type updateInstanceUpdateShieldedInstanceConfigOperation struct {
	// If the update operation has the REQUIRES_APPLY_OPTIONS trait, this will be populated.
	// Usually it will be nil - this is to prevent us from accidentally depending on apply
	// options, which should usually be unnecessary.
	ApplyOptions []dcl.ApplyOption
	FieldDiffs   []*dcl.FieldDiff
}

// do creates a request and sends it to the appropriate URL. In most operations,
// do will transcribe a subset of the resource into a request object and send a
// PUT request to a single URL.

func (op *updateInstanceUpdateShieldedInstanceConfigOperation) do(ctx context.Context, r *Instance, c *Client) error {
	_, err := c.GetInstance(ctx, r)
	if err != nil {
		return err
	}

	u, err := r.updateURL(c.Config.BasePath, "updateShieldedInstanceConfig")
	if err != nil {
		return err
	}

	req, err := newUpdateInstanceUpdateShieldedInstanceConfigRequest(ctx, r, c)
	if err != nil {
		return err
	}

	c.Config.Logger.InfoWithContextf(ctx, "Created update: %#v", req)
	body, err := marshalUpdateInstanceUpdateShieldedInstanceConfigRequest(c, req)
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

func (c *Client) listInstanceRaw(ctx context.Context, r *Instance, pageToken string, pageSize int32) ([]byte, error) {
	u, err := r.urlNormalized().listURL(c.Config.BasePath)
	if err != nil {
		return nil, err
	}

	m := make(map[string]string)
	if pageToken != "" {
		m["pageToken"] = pageToken
	}

	if pageSize != InstanceMaxPage {
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

type listInstanceOperation struct {
	Items []map[string]interface{} `json:"items"`
	Token string                   `json:"nextPageToken"`
}

func (c *Client) listInstance(ctx context.Context, r *Instance, pageToken string, pageSize int32) ([]*Instance, string, error) {
	b, err := c.listInstanceRaw(ctx, r, pageToken, pageSize)
	if err != nil {
		return nil, "", err
	}

	var m listInstanceOperation
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, "", err
	}

	var l []*Instance
	for _, v := range m.Items {
		res, err := unmarshalMapInstance(v, c, r)
		if err != nil {
			return nil, m.Token, err
		}
		res.Project = r.Project
		res.Zone = r.Zone
		l = append(l, res)
	}

	return l, m.Token, nil
}

func (c *Client) deleteAllInstance(ctx context.Context, f func(*Instance) bool, resources []*Instance) error {
	var errors []string
	for _, res := range resources {
		if f(res) {
			// We do not want deleteAll to fail on a deletion or else it will stop deleting other resources.
			err := c.DeleteInstance(ctx, res)
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

type deleteInstanceOperation struct{}

func (op *deleteInstanceOperation) do(ctx context.Context, r *Instance, c *Client) error {
	r, err := c.GetInstance(ctx, r)
	if err != nil {
		if dcl.IsNotFound(err) {
			c.Config.Logger.InfoWithContextf(ctx, "Instance not found, returning. Original error: %v", err)
			return nil
		}
		c.Config.Logger.WarningWithContextf(ctx, "GetInstance checking for existence. error: %v", err)
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
		_, err := c.GetInstance(ctx, r)
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
type createInstanceOperation struct {
	response map[string]interface{}
}

func (op *createInstanceOperation) FirstResponse() (map[string]interface{}, bool) {
	return op.response, len(op.response) > 0
}

func (op *createInstanceOperation) do(ctx context.Context, r *Instance, c *Client) error {
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

	if _, err := c.GetInstance(ctx, r); err != nil {
		c.Config.Logger.WarningWithContextf(ctx, "get returned error: %v", err)
		return err
	}

	return nil
}

func (c *Client) getInstanceRaw(ctx context.Context, r *Instance) ([]byte, error) {

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

func (c *Client) instanceDiffsForRawDesired(ctx context.Context, rawDesired *Instance, opts ...dcl.ApplyOption) (initial, desired *Instance, diffs []*dcl.FieldDiff, err error) {
	c.Config.Logger.InfoWithContext(ctx, "Fetching initial state...")
	// First, let us see if the user provided a state hint.  If they did, we will start fetching based on that.
	var fetchState *Instance
	if sh := dcl.FetchStateHint(opts); sh != nil {
		if r, ok := sh.(*Instance); !ok {
			c.Config.Logger.WarningWithContextf(ctx, "Initial state hint was of the wrong type; expected Instance, got %T", sh)
		} else {
			fetchState = r
		}
	}
	if fetchState == nil {
		fetchState = rawDesired
	}

	// 1.2: Retrieval of raw initial state from API
	rawInitial, err := c.GetInstance(ctx, fetchState)
	if rawInitial == nil {
		if !dcl.IsNotFound(err) {
			c.Config.Logger.WarningWithContextf(ctx, "Failed to retrieve whether a Instance resource already exists: %s", err)
			return nil, nil, nil, fmt.Errorf("failed to retrieve Instance resource: %v", err)
		}
		c.Config.Logger.InfoWithContext(ctx, "Found that Instance resource did not exist.")
		// Perform canonicalization to pick up defaults.
		desired, err = canonicalizeInstanceDesiredState(rawDesired, rawInitial)
		return nil, desired, nil, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Found initial state for Instance: %v", rawInitial)
	c.Config.Logger.InfoWithContextf(ctx, "Initial desired state for Instance: %v", rawDesired)

	// The Get call applies postReadExtract and so the result may contain fields that are not part of API version.
	if err := extractInstanceFields(rawInitial); err != nil {
		return nil, nil, nil, err
	}

	// 1.3: Canonicalize raw initial state into initial state.
	initial, err = canonicalizeInstanceInitialState(rawInitial, rawDesired)
	if err != nil {
		return nil, nil, nil, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Canonicalized initial state for Instance: %v", initial)

	// 1.4: Canonicalize raw desired state into desired state.
	desired, err = canonicalizeInstanceDesiredState(rawDesired, rawInitial, opts...)
	if err != nil {
		return nil, nil, nil, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Canonicalized desired state for Instance: %v", desired)

	// 2.1: Comparison of initial and desired state.
	diffs, err = diffInstance(c, desired, initial, opts...)
	return initial, desired, diffs, err
}

func canonicalizeInstanceInitialState(rawInitial, rawDesired *Instance) (*Instance, error) {
	// TODO(magic-modules-eng): write canonicalizer once relevant traits are added.
	return rawInitial, nil
}

/*
* Canonicalizers
*
* These are responsible for converting either a user-specified config or a
* GCP API response to a standard format that can be used for difference checking.
* */

func canonicalizeInstanceDesiredState(rawDesired, rawInitial *Instance, opts ...dcl.ApplyOption) (*Instance, error) {

	if rawInitial == nil {
		// Since the initial state is empty, the desired state is all we have.
		// We canonicalize the remaining nested objects with nil to pick up defaults.
		rawDesired.Scheduling = canonicalizeInstanceScheduling(rawDesired.Scheduling, nil, opts...)
		rawDesired.ShieldedInstanceConfig = canonicalizeInstanceShieldedInstanceConfig(rawDesired.ShieldedInstanceConfig, nil, opts...)

		return rawDesired, nil
	}
	canonicalDesired := &Instance{}
	if dcl.BoolCanonicalize(rawDesired.CanIPForward, rawInitial.CanIPForward) {
		canonicalDesired.CanIPForward = rawInitial.CanIPForward
	} else {
		canonicalDesired.CanIPForward = rawDesired.CanIPForward
	}
	if dcl.BoolCanonicalize(rawDesired.DeletionProtection, rawInitial.DeletionProtection) {
		canonicalDesired.DeletionProtection = rawInitial.DeletionProtection
	} else {
		canonicalDesired.DeletionProtection = rawDesired.DeletionProtection
	}
	if dcl.StringCanonicalize(rawDesired.Description, rawInitial.Description) {
		canonicalDesired.Description = rawInitial.Description
	} else {
		canonicalDesired.Description = rawDesired.Description
	}
	canonicalDesired.Disks = canonicalizeInstanceDisksSlice(rawDesired.Disks, rawInitial.Disks, opts...)
	canonicalDesired.GuestAccelerators = canonicalizeInstanceGuestAcceleratorsSlice(rawDesired.GuestAccelerators, rawInitial.GuestAccelerators, opts...)
	if dcl.StringCanonicalize(rawDesired.Hostname, rawInitial.Hostname) {
		canonicalDesired.Hostname = rawInitial.Hostname
	} else {
		canonicalDesired.Hostname = rawDesired.Hostname
	}
	if dcl.IsZeroValue(rawDesired.Labels) || (dcl.IsEmptyValueIndirect(rawDesired.Labels) && dcl.IsEmptyValueIndirect(rawInitial.Labels)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		canonicalDesired.Labels = rawInitial.Labels
	} else {
		canonicalDesired.Labels = rawDesired.Labels
	}
	if dcl.IsZeroValue(rawDesired.Metadata) || (dcl.IsEmptyValueIndirect(rawDesired.Metadata) && dcl.IsEmptyValueIndirect(rawInitial.Metadata)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		canonicalDesired.Metadata = rawInitial.Metadata
	} else {
		canonicalDesired.Metadata = rawDesired.Metadata
	}
	if dcl.PartialSelfLinkToSelfLink(rawDesired.MachineType, rawInitial.MachineType) {
		canonicalDesired.MachineType = rawInitial.MachineType
	} else {
		canonicalDesired.MachineType = rawDesired.MachineType
	}
	if dcl.StringCanonicalize(rawDesired.MinCpuPlatform, rawInitial.MinCpuPlatform) {
		canonicalDesired.MinCpuPlatform = rawInitial.MinCpuPlatform
	} else {
		canonicalDesired.MinCpuPlatform = rawDesired.MinCpuPlatform
	}
	if dcl.StringCanonicalize(rawDesired.Name, rawInitial.Name) {
		canonicalDesired.Name = rawInitial.Name
	} else {
		canonicalDesired.Name = rawDesired.Name
	}
	canonicalDesired.NetworkInterfaces = canonicalizeInstanceNetworkInterfacesSlice(rawDesired.NetworkInterfaces, rawInitial.NetworkInterfaces, opts...)
	canonicalDesired.Scheduling = canonicalizeInstanceScheduling(rawDesired.Scheduling, rawInitial.Scheduling, opts...)
	canonicalDesired.ServiceAccounts = canonicalizeInstanceServiceAccountsSlice(rawDesired.ServiceAccounts, rawInitial.ServiceAccounts, opts...)
	canonicalDesired.ShieldedInstanceConfig = canonicalizeInstanceShieldedInstanceConfig(rawDesired.ShieldedInstanceConfig, rawInitial.ShieldedInstanceConfig, opts...)
	if dcl.IsZeroValue(rawDesired.Status) || (dcl.IsEmptyValueIndirect(rawDesired.Status) && dcl.IsEmptyValueIndirect(rawInitial.Status)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		canonicalDesired.Status = rawInitial.Status
	} else {
		canonicalDesired.Status = rawDesired.Status
	}
	if dcl.StringArrayCanonicalize(rawDesired.Tags, rawInitial.Tags) {
		canonicalDesired.Tags = rawInitial.Tags
	} else {
		canonicalDesired.Tags = rawDesired.Tags
	}
	if dcl.IsZeroValue(rawDesired.Zone) || (dcl.IsEmptyValueIndirect(rawDesired.Zone) && dcl.IsEmptyValueIndirect(rawInitial.Zone)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		canonicalDesired.Zone = rawInitial.Zone
	} else {
		canonicalDesired.Zone = rawDesired.Zone
	}
	if dcl.NameToSelfLink(rawDesired.Project, rawInitial.Project) {
		canonicalDesired.Project = rawInitial.Project
	} else {
		canonicalDesired.Project = rawDesired.Project
	}
	return canonicalDesired, nil
}

func canonicalizeInstanceNewState(c *Client, rawNew, rawDesired *Instance) (*Instance, error) {

	if dcl.IsEmptyValueIndirect(rawNew.CanIPForward) && dcl.IsEmptyValueIndirect(rawDesired.CanIPForward) {
		rawNew.CanIPForward = rawDesired.CanIPForward
	} else {
		if dcl.BoolCanonicalize(rawDesired.CanIPForward, rawNew.CanIPForward) {
			rawNew.CanIPForward = rawDesired.CanIPForward
		}
	}

	if dcl.IsEmptyValueIndirect(rawNew.CpuPlatform) && dcl.IsEmptyValueIndirect(rawDesired.CpuPlatform) {
		rawNew.CpuPlatform = rawDesired.CpuPlatform
	} else {
		if dcl.StringCanonicalize(rawDesired.CpuPlatform, rawNew.CpuPlatform) {
			rawNew.CpuPlatform = rawDesired.CpuPlatform
		}
	}

	if dcl.IsEmptyValueIndirect(rawNew.CreationTimestamp) && dcl.IsEmptyValueIndirect(rawDesired.CreationTimestamp) {
		rawNew.CreationTimestamp = rawDesired.CreationTimestamp
	} else {
		if dcl.StringCanonicalize(rawDesired.CreationTimestamp, rawNew.CreationTimestamp) {
			rawNew.CreationTimestamp = rawDesired.CreationTimestamp
		}
	}

	if dcl.IsEmptyValueIndirect(rawNew.DeletionProtection) && dcl.IsEmptyValueIndirect(rawDesired.DeletionProtection) {
		rawNew.DeletionProtection = rawDesired.DeletionProtection
	} else {
		if dcl.BoolCanonicalize(rawDesired.DeletionProtection, rawNew.DeletionProtection) {
			rawNew.DeletionProtection = rawDesired.DeletionProtection
		}
	}

	if dcl.IsEmptyValueIndirect(rawNew.Description) && dcl.IsEmptyValueIndirect(rawDesired.Description) {
		rawNew.Description = rawDesired.Description
	} else {
		if dcl.StringCanonicalize(rawDesired.Description, rawNew.Description) {
			rawNew.Description = rawDesired.Description
		}
	}

	rawNew.Disks = rawDesired.Disks

	if dcl.IsEmptyValueIndirect(rawNew.GuestAccelerators) && dcl.IsEmptyValueIndirect(rawDesired.GuestAccelerators) {
		rawNew.GuestAccelerators = rawDesired.GuestAccelerators
	} else {
		rawNew.GuestAccelerators = canonicalizeNewInstanceGuestAcceleratorsSlice(c, rawDesired.GuestAccelerators, rawNew.GuestAccelerators)
	}

	if dcl.IsEmptyValueIndirect(rawNew.Hostname) && dcl.IsEmptyValueIndirect(rawDesired.Hostname) {
		rawNew.Hostname = rawDesired.Hostname
	} else {
		if dcl.StringCanonicalize(rawDesired.Hostname, rawNew.Hostname) {
			rawNew.Hostname = rawDesired.Hostname
		}
	}

	if dcl.IsEmptyValueIndirect(rawNew.Id) && dcl.IsEmptyValueIndirect(rawDesired.Id) {
		rawNew.Id = rawDesired.Id
	} else {
		if dcl.StringCanonicalize(rawDesired.Id, rawNew.Id) {
			rawNew.Id = rawDesired.Id
		}
	}

	if dcl.IsEmptyValueIndirect(rawNew.Labels) && dcl.IsEmptyValueIndirect(rawDesired.Labels) {
		rawNew.Labels = rawDesired.Labels
	} else {
	}

	if dcl.IsEmptyValueIndirect(rawNew.Metadata) && dcl.IsEmptyValueIndirect(rawDesired.Metadata) {
		rawNew.Metadata = rawDesired.Metadata
	} else {
	}

	if dcl.IsEmptyValueIndirect(rawNew.MachineType) && dcl.IsEmptyValueIndirect(rawDesired.MachineType) {
		rawNew.MachineType = rawDesired.MachineType
	} else {
		if dcl.PartialSelfLinkToSelfLink(rawDesired.MachineType, rawNew.MachineType) {
			rawNew.MachineType = rawDesired.MachineType
		}
	}

	if dcl.IsEmptyValueIndirect(rawNew.MinCpuPlatform) && dcl.IsEmptyValueIndirect(rawDesired.MinCpuPlatform) {
		rawNew.MinCpuPlatform = rawDesired.MinCpuPlatform
	} else {
		if dcl.StringCanonicalize(rawDesired.MinCpuPlatform, rawNew.MinCpuPlatform) {
			rawNew.MinCpuPlatform = rawDesired.MinCpuPlatform
		}
	}

	if dcl.IsEmptyValueIndirect(rawNew.Name) && dcl.IsEmptyValueIndirect(rawDesired.Name) {
		rawNew.Name = rawDesired.Name
	} else {
		if dcl.StringCanonicalize(rawDesired.Name, rawNew.Name) {
			rawNew.Name = rawDesired.Name
		}
	}

	if dcl.IsEmptyValueIndirect(rawNew.NetworkInterfaces) && dcl.IsEmptyValueIndirect(rawDesired.NetworkInterfaces) {
		rawNew.NetworkInterfaces = rawDesired.NetworkInterfaces
	} else {
		rawNew.NetworkInterfaces = canonicalizeNewInstanceNetworkInterfacesSlice(c, rawDesired.NetworkInterfaces, rawNew.NetworkInterfaces)
	}

	if dcl.IsEmptyValueIndirect(rawNew.Scheduling) && dcl.IsEmptyValueIndirect(rawDesired.Scheduling) {
		rawNew.Scheduling = rawDesired.Scheduling
	} else {
		rawNew.Scheduling = canonicalizeNewInstanceScheduling(c, rawDesired.Scheduling, rawNew.Scheduling)
	}

	if dcl.IsEmptyValueIndirect(rawNew.ServiceAccounts) && dcl.IsEmptyValueIndirect(rawDesired.ServiceAccounts) {
		rawNew.ServiceAccounts = rawDesired.ServiceAccounts
	} else {
		rawNew.ServiceAccounts = canonicalizeNewInstanceServiceAccountsSlice(c, rawDesired.ServiceAccounts, rawNew.ServiceAccounts)
	}

	if dcl.IsEmptyValueIndirect(rawNew.ShieldedInstanceConfig) && dcl.IsEmptyValueIndirect(rawDesired.ShieldedInstanceConfig) {
		rawNew.ShieldedInstanceConfig = rawDesired.ShieldedInstanceConfig
	} else {
		rawNew.ShieldedInstanceConfig = canonicalizeNewInstanceShieldedInstanceConfig(c, rawDesired.ShieldedInstanceConfig, rawNew.ShieldedInstanceConfig)
	}

	if dcl.IsEmptyValueIndirect(rawNew.Status) && dcl.IsEmptyValueIndirect(rawDesired.Status) {
		rawNew.Status = rawDesired.Status
	} else {
	}

	if dcl.IsEmptyValueIndirect(rawNew.StatusMessage) && dcl.IsEmptyValueIndirect(rawDesired.StatusMessage) {
		rawNew.StatusMessage = rawDesired.StatusMessage
	} else {
		if dcl.StringCanonicalize(rawDesired.StatusMessage, rawNew.StatusMessage) {
			rawNew.StatusMessage = rawDesired.StatusMessage
		}
	}

	if dcl.IsEmptyValueIndirect(rawNew.Tags) && dcl.IsEmptyValueIndirect(rawDesired.Tags) {
		rawNew.Tags = rawDesired.Tags
	} else {
		if dcl.StringArrayCanonicalize(rawDesired.Tags, rawNew.Tags) {
			rawNew.Tags = rawDesired.Tags
		}
	}

	if dcl.IsEmptyValueIndirect(rawNew.Zone) && dcl.IsEmptyValueIndirect(rawDesired.Zone) {
		rawNew.Zone = rawDesired.Zone
	} else {
	}

	rawNew.Project = rawDesired.Project

	if dcl.IsEmptyValueIndirect(rawNew.SelfLink) && dcl.IsEmptyValueIndirect(rawDesired.SelfLink) {
		rawNew.SelfLink = rawDesired.SelfLink
	} else {
		if dcl.StringCanonicalize(rawDesired.SelfLink, rawNew.SelfLink) {
			rawNew.SelfLink = rawDesired.SelfLink
		}
	}

	return rawNew, nil
}

func canonicalizeInstanceDisks(des, initial *InstanceDisks, opts ...dcl.ApplyOption) *InstanceDisks {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &InstanceDisks{}

	if dcl.BoolCanonicalize(des.AutoDelete, initial.AutoDelete) || dcl.IsZeroValue(des.AutoDelete) {
		cDes.AutoDelete = initial.AutoDelete
	} else {
		cDes.AutoDelete = des.AutoDelete
	}
	if dcl.BoolCanonicalize(des.Boot, initial.Boot) || dcl.IsZeroValue(des.Boot) {
		cDes.Boot = initial.Boot
	} else {
		cDes.Boot = des.Boot
	}
	if dcl.StringCanonicalize(des.DeviceName, initial.DeviceName) || dcl.IsZeroValue(des.DeviceName) {
		cDes.DeviceName = initial.DeviceName
	} else {
		cDes.DeviceName = des.DeviceName
	}
	cDes.DiskEncryptionKey = canonicalizeInstanceDisksDiskEncryptionKey(des.DiskEncryptionKey, initial.DiskEncryptionKey, opts...)
	if dcl.IsZeroValue(des.Index) || (dcl.IsEmptyValueIndirect(des.Index) && dcl.IsEmptyValueIndirect(initial.Index)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.Index = initial.Index
	} else {
		cDes.Index = des.Index
	}
	cDes.InitializeParams = canonicalizeInstanceDisksInitializeParams(des.InitializeParams, initial.InitializeParams, opts...)
	if dcl.IsZeroValue(des.Interface) || (dcl.IsEmptyValueIndirect(des.Interface) && dcl.IsEmptyValueIndirect(initial.Interface)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.Interface = initial.Interface
	} else {
		cDes.Interface = des.Interface
	}
	if dcl.IsZeroValue(des.Mode) || (dcl.IsEmptyValueIndirect(des.Mode) && dcl.IsEmptyValueIndirect(initial.Mode)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.Mode = initial.Mode
	} else {
		cDes.Mode = des.Mode
	}
	if dcl.IsZeroValue(des.Source) || (dcl.IsEmptyValueIndirect(des.Source) && dcl.IsEmptyValueIndirect(initial.Source)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.Source = initial.Source
	} else {
		cDes.Source = des.Source
	}
	if dcl.IsZeroValue(des.Type) || (dcl.IsEmptyValueIndirect(des.Type) && dcl.IsEmptyValueIndirect(initial.Type)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.Type = initial.Type
	} else {
		cDes.Type = des.Type
	}

	return cDes
}

func canonicalizeInstanceDisksSlice(des, initial []InstanceDisks, opts ...dcl.ApplyOption) []InstanceDisks {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]InstanceDisks, 0, len(des))
		for _, d := range des {
			cd := canonicalizeInstanceDisks(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]InstanceDisks, 0, len(des))
	for i, d := range des {
		cd := canonicalizeInstanceDisks(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewInstanceDisks(c *Client, des, nw *InstanceDisks) *InstanceDisks {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for InstanceDisks while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.BoolCanonicalize(des.AutoDelete, nw.AutoDelete) {
		nw.AutoDelete = des.AutoDelete
	}
	if dcl.BoolCanonicalize(des.Boot, nw.Boot) {
		nw.Boot = des.Boot
	}
	if dcl.StringCanonicalize(des.DeviceName, nw.DeviceName) {
		nw.DeviceName = des.DeviceName
	}
	nw.DiskEncryptionKey = canonicalizeNewInstanceDisksDiskEncryptionKey(c, des.DiskEncryptionKey, nw.DiskEncryptionKey)
	nw.InitializeParams = des.InitializeParams

	return nw
}

func canonicalizeNewInstanceDisksSet(c *Client, des, nw []InstanceDisks) []InstanceDisks {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []InstanceDisks
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareInstanceDisksNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewInstanceDisks(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewInstanceDisksSlice(c *Client, des, nw []InstanceDisks) []InstanceDisks {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []InstanceDisks
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewInstanceDisks(c, &d, &n))
	}

	return items
}

func canonicalizeInstanceDisksDiskEncryptionKey(des, initial *InstanceDisksDiskEncryptionKey, opts ...dcl.ApplyOption) *InstanceDisksDiskEncryptionKey {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &InstanceDisksDiskEncryptionKey{}

	if dcl.StringCanonicalize(des.RawKey, initial.RawKey) || dcl.IsZeroValue(des.RawKey) {
		cDes.RawKey = initial.RawKey
	} else {
		cDes.RawKey = des.RawKey
	}
	if dcl.StringCanonicalize(des.RsaEncryptedKey, initial.RsaEncryptedKey) || dcl.IsZeroValue(des.RsaEncryptedKey) {
		cDes.RsaEncryptedKey = initial.RsaEncryptedKey
	} else {
		cDes.RsaEncryptedKey = des.RsaEncryptedKey
	}

	return cDes
}

func canonicalizeInstanceDisksDiskEncryptionKeySlice(des, initial []InstanceDisksDiskEncryptionKey, opts ...dcl.ApplyOption) []InstanceDisksDiskEncryptionKey {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]InstanceDisksDiskEncryptionKey, 0, len(des))
		for _, d := range des {
			cd := canonicalizeInstanceDisksDiskEncryptionKey(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]InstanceDisksDiskEncryptionKey, 0, len(des))
	for i, d := range des {
		cd := canonicalizeInstanceDisksDiskEncryptionKey(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewInstanceDisksDiskEncryptionKey(c *Client, des, nw *InstanceDisksDiskEncryptionKey) *InstanceDisksDiskEncryptionKey {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for InstanceDisksDiskEncryptionKey while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringCanonicalize(des.RawKey, nw.RawKey) {
		nw.RawKey = des.RawKey
	}
	if dcl.StringCanonicalize(des.RsaEncryptedKey, nw.RsaEncryptedKey) {
		nw.RsaEncryptedKey = des.RsaEncryptedKey
	}
	if dcl.StringCanonicalize(des.Sha256, nw.Sha256) {
		nw.Sha256 = des.Sha256
	}

	return nw
}

func canonicalizeNewInstanceDisksDiskEncryptionKeySet(c *Client, des, nw []InstanceDisksDiskEncryptionKey) []InstanceDisksDiskEncryptionKey {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []InstanceDisksDiskEncryptionKey
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareInstanceDisksDiskEncryptionKeyNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewInstanceDisksDiskEncryptionKey(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewInstanceDisksDiskEncryptionKeySlice(c *Client, des, nw []InstanceDisksDiskEncryptionKey) []InstanceDisksDiskEncryptionKey {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []InstanceDisksDiskEncryptionKey
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewInstanceDisksDiskEncryptionKey(c, &d, &n))
	}

	return items
}

func canonicalizeInstanceDisksInitializeParams(des, initial *InstanceDisksInitializeParams, opts ...dcl.ApplyOption) *InstanceDisksInitializeParams {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &InstanceDisksInitializeParams{}

	if dcl.StringCanonicalize(des.DiskName, initial.DiskName) || dcl.IsZeroValue(des.DiskName) {
		cDes.DiskName = initial.DiskName
	} else {
		cDes.DiskName = des.DiskName
	}
	if dcl.IsZeroValue(des.DiskSizeGb) || (dcl.IsEmptyValueIndirect(des.DiskSizeGb) && dcl.IsEmptyValueIndirect(initial.DiskSizeGb)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.DiskSizeGb = initial.DiskSizeGb
	} else {
		cDes.DiskSizeGb = des.DiskSizeGb
	}
	if dcl.IsZeroValue(des.DiskType) || (dcl.IsEmptyValueIndirect(des.DiskType) && dcl.IsEmptyValueIndirect(initial.DiskType)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.DiskType = initial.DiskType
	} else {
		cDes.DiskType = des.DiskType
	}
	if dcl.StringCanonicalize(des.SourceImage, initial.SourceImage) || dcl.IsZeroValue(des.SourceImage) {
		cDes.SourceImage = initial.SourceImage
	} else {
		cDes.SourceImage = des.SourceImage
	}
	cDes.SourceImageEncryptionKey = canonicalizeInstanceDisksInitializeParamsSourceImageEncryptionKey(des.SourceImageEncryptionKey, initial.SourceImageEncryptionKey, opts...)

	return cDes
}

func canonicalizeInstanceDisksInitializeParamsSlice(des, initial []InstanceDisksInitializeParams, opts ...dcl.ApplyOption) []InstanceDisksInitializeParams {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]InstanceDisksInitializeParams, 0, len(des))
		for _, d := range des {
			cd := canonicalizeInstanceDisksInitializeParams(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]InstanceDisksInitializeParams, 0, len(des))
	for i, d := range des {
		cd := canonicalizeInstanceDisksInitializeParams(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewInstanceDisksInitializeParams(c *Client, des, nw *InstanceDisksInitializeParams) *InstanceDisksInitializeParams {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for InstanceDisksInitializeParams while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringCanonicalize(des.DiskName, nw.DiskName) {
		nw.DiskName = des.DiskName
	}
	if dcl.StringCanonicalize(des.SourceImage, nw.SourceImage) {
		nw.SourceImage = des.SourceImage
	}
	nw.SourceImageEncryptionKey = canonicalizeNewInstanceDisksInitializeParamsSourceImageEncryptionKey(c, des.SourceImageEncryptionKey, nw.SourceImageEncryptionKey)

	return nw
}

func canonicalizeNewInstanceDisksInitializeParamsSet(c *Client, des, nw []InstanceDisksInitializeParams) []InstanceDisksInitializeParams {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []InstanceDisksInitializeParams
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareInstanceDisksInitializeParamsNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewInstanceDisksInitializeParams(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewInstanceDisksInitializeParamsSlice(c *Client, des, nw []InstanceDisksInitializeParams) []InstanceDisksInitializeParams {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []InstanceDisksInitializeParams
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewInstanceDisksInitializeParams(c, &d, &n))
	}

	return items
}

func canonicalizeInstanceDisksInitializeParamsSourceImageEncryptionKey(des, initial *InstanceDisksInitializeParamsSourceImageEncryptionKey, opts ...dcl.ApplyOption) *InstanceDisksInitializeParamsSourceImageEncryptionKey {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &InstanceDisksInitializeParamsSourceImageEncryptionKey{}

	if dcl.StringCanonicalize(des.RawKey, initial.RawKey) || dcl.IsZeroValue(des.RawKey) {
		cDes.RawKey = initial.RawKey
	} else {
		cDes.RawKey = des.RawKey
	}

	return cDes
}

func canonicalizeInstanceDisksInitializeParamsSourceImageEncryptionKeySlice(des, initial []InstanceDisksInitializeParamsSourceImageEncryptionKey, opts ...dcl.ApplyOption) []InstanceDisksInitializeParamsSourceImageEncryptionKey {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]InstanceDisksInitializeParamsSourceImageEncryptionKey, 0, len(des))
		for _, d := range des {
			cd := canonicalizeInstanceDisksInitializeParamsSourceImageEncryptionKey(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]InstanceDisksInitializeParamsSourceImageEncryptionKey, 0, len(des))
	for i, d := range des {
		cd := canonicalizeInstanceDisksInitializeParamsSourceImageEncryptionKey(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewInstanceDisksInitializeParamsSourceImageEncryptionKey(c *Client, des, nw *InstanceDisksInitializeParamsSourceImageEncryptionKey) *InstanceDisksInitializeParamsSourceImageEncryptionKey {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for InstanceDisksInitializeParamsSourceImageEncryptionKey while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringCanonicalize(des.RawKey, nw.RawKey) {
		nw.RawKey = des.RawKey
	}
	if dcl.StringCanonicalize(des.Sha256, nw.Sha256) {
		nw.Sha256 = des.Sha256
	}

	return nw
}

func canonicalizeNewInstanceDisksInitializeParamsSourceImageEncryptionKeySet(c *Client, des, nw []InstanceDisksInitializeParamsSourceImageEncryptionKey) []InstanceDisksInitializeParamsSourceImageEncryptionKey {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []InstanceDisksInitializeParamsSourceImageEncryptionKey
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareInstanceDisksInitializeParamsSourceImageEncryptionKeyNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewInstanceDisksInitializeParamsSourceImageEncryptionKey(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewInstanceDisksInitializeParamsSourceImageEncryptionKeySlice(c *Client, des, nw []InstanceDisksInitializeParamsSourceImageEncryptionKey) []InstanceDisksInitializeParamsSourceImageEncryptionKey {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []InstanceDisksInitializeParamsSourceImageEncryptionKey
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewInstanceDisksInitializeParamsSourceImageEncryptionKey(c, &d, &n))
	}

	return items
}

func canonicalizeInstanceGuestAccelerators(des, initial *InstanceGuestAccelerators, opts ...dcl.ApplyOption) *InstanceGuestAccelerators {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &InstanceGuestAccelerators{}

	if dcl.IsZeroValue(des.AcceleratorCount) || (dcl.IsEmptyValueIndirect(des.AcceleratorCount) && dcl.IsEmptyValueIndirect(initial.AcceleratorCount)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.AcceleratorCount = initial.AcceleratorCount
	} else {
		cDes.AcceleratorCount = des.AcceleratorCount
	}
	if dcl.StringCanonicalize(des.AcceleratorType, initial.AcceleratorType) || dcl.IsZeroValue(des.AcceleratorType) {
		cDes.AcceleratorType = initial.AcceleratorType
	} else {
		cDes.AcceleratorType = des.AcceleratorType
	}

	return cDes
}

func canonicalizeInstanceGuestAcceleratorsSlice(des, initial []InstanceGuestAccelerators, opts ...dcl.ApplyOption) []InstanceGuestAccelerators {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]InstanceGuestAccelerators, 0, len(des))
		for _, d := range des {
			cd := canonicalizeInstanceGuestAccelerators(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]InstanceGuestAccelerators, 0, len(des))
	for i, d := range des {
		cd := canonicalizeInstanceGuestAccelerators(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewInstanceGuestAccelerators(c *Client, des, nw *InstanceGuestAccelerators) *InstanceGuestAccelerators {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for InstanceGuestAccelerators while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringCanonicalize(des.AcceleratorType, nw.AcceleratorType) {
		nw.AcceleratorType = des.AcceleratorType
	}

	return nw
}

func canonicalizeNewInstanceGuestAcceleratorsSet(c *Client, des, nw []InstanceGuestAccelerators) []InstanceGuestAccelerators {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []InstanceGuestAccelerators
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareInstanceGuestAcceleratorsNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewInstanceGuestAccelerators(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewInstanceGuestAcceleratorsSlice(c *Client, des, nw []InstanceGuestAccelerators) []InstanceGuestAccelerators {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []InstanceGuestAccelerators
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewInstanceGuestAccelerators(c, &d, &n))
	}

	return items
}

func canonicalizeInstanceNetworkInterfaces(des, initial *InstanceNetworkInterfaces, opts ...dcl.ApplyOption) *InstanceNetworkInterfaces {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &InstanceNetworkInterfaces{}

	cDes.AccessConfigs = canonicalizeInstanceNetworkInterfacesAccessConfigsSlice(des.AccessConfigs, initial.AccessConfigs, opts...)
	cDes.IPv6AccessConfigs = canonicalizeInstanceNetworkInterfacesIPv6AccessConfigsSlice(des.IPv6AccessConfigs, initial.IPv6AccessConfigs, opts...)
	cDes.AliasIPRanges = canonicalizeInstanceNetworkInterfacesAliasIPRangesSlice(des.AliasIPRanges, initial.AliasIPRanges, opts...)
	if dcl.IsZeroValue(des.Network) || (dcl.IsEmptyValueIndirect(des.Network) && dcl.IsEmptyValueIndirect(initial.Network)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.Network = initial.Network
	} else {
		cDes.Network = des.Network
	}
	if dcl.StringCanonicalize(des.NetworkIP, initial.NetworkIP) || dcl.IsZeroValue(des.NetworkIP) {
		cDes.NetworkIP = initial.NetworkIP
	} else {
		cDes.NetworkIP = des.NetworkIP
	}
	if dcl.IsZeroValue(des.Subnetwork) || (dcl.IsEmptyValueIndirect(des.Subnetwork) && dcl.IsEmptyValueIndirect(initial.Subnetwork)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.Subnetwork = initial.Subnetwork
	} else {
		cDes.Subnetwork = des.Subnetwork
	}

	return cDes
}

func canonicalizeInstanceNetworkInterfacesSlice(des, initial []InstanceNetworkInterfaces, opts ...dcl.ApplyOption) []InstanceNetworkInterfaces {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]InstanceNetworkInterfaces, 0, len(des))
		for _, d := range des {
			cd := canonicalizeInstanceNetworkInterfaces(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]InstanceNetworkInterfaces, 0, len(des))
	for i, d := range des {
		cd := canonicalizeInstanceNetworkInterfaces(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewInstanceNetworkInterfaces(c *Client, des, nw *InstanceNetworkInterfaces) *InstanceNetworkInterfaces {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for InstanceNetworkInterfaces while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	nw.AccessConfigs = canonicalizeNewInstanceNetworkInterfacesAccessConfigsSlice(c, des.AccessConfigs, nw.AccessConfigs)
	nw.IPv6AccessConfigs = canonicalizeNewInstanceNetworkInterfacesIPv6AccessConfigsSlice(c, des.IPv6AccessConfigs, nw.IPv6AccessConfigs)
	nw.AliasIPRanges = canonicalizeNewInstanceNetworkInterfacesAliasIPRangesSlice(c, des.AliasIPRanges, nw.AliasIPRanges)
	if dcl.StringCanonicalize(des.Name, nw.Name) {
		nw.Name = des.Name
	}
	if dcl.StringCanonicalize(des.NetworkIP, nw.NetworkIP) {
		nw.NetworkIP = des.NetworkIP
	}

	return nw
}

func canonicalizeNewInstanceNetworkInterfacesSet(c *Client, des, nw []InstanceNetworkInterfaces) []InstanceNetworkInterfaces {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []InstanceNetworkInterfaces
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareInstanceNetworkInterfacesNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewInstanceNetworkInterfaces(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewInstanceNetworkInterfacesSlice(c *Client, des, nw []InstanceNetworkInterfaces) []InstanceNetworkInterfaces {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []InstanceNetworkInterfaces
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewInstanceNetworkInterfaces(c, &d, &n))
	}

	return items
}

func canonicalizeInstanceNetworkInterfacesAccessConfigs(des, initial *InstanceNetworkInterfacesAccessConfigs, opts ...dcl.ApplyOption) *InstanceNetworkInterfacesAccessConfigs {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &InstanceNetworkInterfacesAccessConfigs{}

	if dcl.StringCanonicalize(des.Name, initial.Name) || dcl.IsZeroValue(des.Name) {
		cDes.Name = initial.Name
	} else {
		cDes.Name = des.Name
	}
	if dcl.IsZeroValue(des.NatIP) || (dcl.IsEmptyValueIndirect(des.NatIP) && dcl.IsEmptyValueIndirect(initial.NatIP)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.NatIP = initial.NatIP
	} else {
		cDes.NatIP = des.NatIP
	}
	if dcl.BoolCanonicalize(des.SetPublicPtr, initial.SetPublicPtr) || dcl.IsZeroValue(des.SetPublicPtr) {
		cDes.SetPublicPtr = initial.SetPublicPtr
	} else {
		cDes.SetPublicPtr = des.SetPublicPtr
	}
	if dcl.StringCanonicalize(des.PublicPtrDomainName, initial.PublicPtrDomainName) || dcl.IsZeroValue(des.PublicPtrDomainName) {
		cDes.PublicPtrDomainName = initial.PublicPtrDomainName
	} else {
		cDes.PublicPtrDomainName = des.PublicPtrDomainName
	}
	if dcl.IsZeroValue(des.NetworkTier) || (dcl.IsEmptyValueIndirect(des.NetworkTier) && dcl.IsEmptyValueIndirect(initial.NetworkTier)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.NetworkTier = initial.NetworkTier
	} else {
		cDes.NetworkTier = des.NetworkTier
	}
	if dcl.IsZeroValue(des.Type) || (dcl.IsEmptyValueIndirect(des.Type) && dcl.IsEmptyValueIndirect(initial.Type)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.Type = initial.Type
	} else {
		cDes.Type = des.Type
	}

	return cDes
}

func canonicalizeInstanceNetworkInterfacesAccessConfigsSlice(des, initial []InstanceNetworkInterfacesAccessConfigs, opts ...dcl.ApplyOption) []InstanceNetworkInterfacesAccessConfigs {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]InstanceNetworkInterfacesAccessConfigs, 0, len(des))
		for _, d := range des {
			cd := canonicalizeInstanceNetworkInterfacesAccessConfigs(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]InstanceNetworkInterfacesAccessConfigs, 0, len(des))
	for i, d := range des {
		cd := canonicalizeInstanceNetworkInterfacesAccessConfigs(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewInstanceNetworkInterfacesAccessConfigs(c *Client, des, nw *InstanceNetworkInterfacesAccessConfigs) *InstanceNetworkInterfacesAccessConfigs {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for InstanceNetworkInterfacesAccessConfigs while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringCanonicalize(des.Name, nw.Name) {
		nw.Name = des.Name
	}
	if dcl.StringCanonicalize(des.ExternalIPv6, nw.ExternalIPv6) {
		nw.ExternalIPv6 = des.ExternalIPv6
	}
	if dcl.StringCanonicalize(des.ExternalIPv6PrefixLength, nw.ExternalIPv6PrefixLength) {
		nw.ExternalIPv6PrefixLength = des.ExternalIPv6PrefixLength
	}
	if dcl.BoolCanonicalize(des.SetPublicPtr, nw.SetPublicPtr) {
		nw.SetPublicPtr = des.SetPublicPtr
	}
	if dcl.StringCanonicalize(des.PublicPtrDomainName, nw.PublicPtrDomainName) {
		nw.PublicPtrDomainName = des.PublicPtrDomainName
	}

	return nw
}

func canonicalizeNewInstanceNetworkInterfacesAccessConfigsSet(c *Client, des, nw []InstanceNetworkInterfacesAccessConfigs) []InstanceNetworkInterfacesAccessConfigs {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []InstanceNetworkInterfacesAccessConfigs
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareInstanceNetworkInterfacesAccessConfigsNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewInstanceNetworkInterfacesAccessConfigs(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewInstanceNetworkInterfacesAccessConfigsSlice(c *Client, des, nw []InstanceNetworkInterfacesAccessConfigs) []InstanceNetworkInterfacesAccessConfigs {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []InstanceNetworkInterfacesAccessConfigs
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewInstanceNetworkInterfacesAccessConfigs(c, &d, &n))
	}

	return items
}

func canonicalizeInstanceNetworkInterfacesIPv6AccessConfigs(des, initial *InstanceNetworkInterfacesIPv6AccessConfigs, opts ...dcl.ApplyOption) *InstanceNetworkInterfacesIPv6AccessConfigs {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &InstanceNetworkInterfacesIPv6AccessConfigs{}

	if dcl.StringCanonicalize(des.Name, initial.Name) || dcl.IsZeroValue(des.Name) {
		cDes.Name = initial.Name
	} else {
		cDes.Name = des.Name
	}
	if dcl.IsZeroValue(des.NatIP) || (dcl.IsEmptyValueIndirect(des.NatIP) && dcl.IsEmptyValueIndirect(initial.NatIP)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.NatIP = initial.NatIP
	} else {
		cDes.NatIP = des.NatIP
	}
	if dcl.BoolCanonicalize(des.SetPublicPtr, initial.SetPublicPtr) || dcl.IsZeroValue(des.SetPublicPtr) {
		cDes.SetPublicPtr = initial.SetPublicPtr
	} else {
		cDes.SetPublicPtr = des.SetPublicPtr
	}
	if dcl.StringCanonicalize(des.PublicPtrDomainName, initial.PublicPtrDomainName) || dcl.IsZeroValue(des.PublicPtrDomainName) {
		cDes.PublicPtrDomainName = initial.PublicPtrDomainName
	} else {
		cDes.PublicPtrDomainName = des.PublicPtrDomainName
	}
	if dcl.IsZeroValue(des.NetworkTier) || (dcl.IsEmptyValueIndirect(des.NetworkTier) && dcl.IsEmptyValueIndirect(initial.NetworkTier)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.NetworkTier = initial.NetworkTier
	} else {
		cDes.NetworkTier = des.NetworkTier
	}
	if dcl.IsZeroValue(des.Type) || (dcl.IsEmptyValueIndirect(des.Type) && dcl.IsEmptyValueIndirect(initial.Type)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.Type = initial.Type
	} else {
		cDes.Type = des.Type
	}

	return cDes
}

func canonicalizeInstanceNetworkInterfacesIPv6AccessConfigsSlice(des, initial []InstanceNetworkInterfacesIPv6AccessConfigs, opts ...dcl.ApplyOption) []InstanceNetworkInterfacesIPv6AccessConfigs {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]InstanceNetworkInterfacesIPv6AccessConfigs, 0, len(des))
		for _, d := range des {
			cd := canonicalizeInstanceNetworkInterfacesIPv6AccessConfigs(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]InstanceNetworkInterfacesIPv6AccessConfigs, 0, len(des))
	for i, d := range des {
		cd := canonicalizeInstanceNetworkInterfacesIPv6AccessConfigs(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewInstanceNetworkInterfacesIPv6AccessConfigs(c *Client, des, nw *InstanceNetworkInterfacesIPv6AccessConfigs) *InstanceNetworkInterfacesIPv6AccessConfigs {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for InstanceNetworkInterfacesIPv6AccessConfigs while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringCanonicalize(des.Name, nw.Name) {
		nw.Name = des.Name
	}
	if dcl.StringCanonicalize(des.ExternalIPv6, nw.ExternalIPv6) {
		nw.ExternalIPv6 = des.ExternalIPv6
	}
	if dcl.StringCanonicalize(des.ExternalIPv6PrefixLength, nw.ExternalIPv6PrefixLength) {
		nw.ExternalIPv6PrefixLength = des.ExternalIPv6PrefixLength
	}
	if dcl.BoolCanonicalize(des.SetPublicPtr, nw.SetPublicPtr) {
		nw.SetPublicPtr = des.SetPublicPtr
	}
	if dcl.StringCanonicalize(des.PublicPtrDomainName, nw.PublicPtrDomainName) {
		nw.PublicPtrDomainName = des.PublicPtrDomainName
	}

	return nw
}

func canonicalizeNewInstanceNetworkInterfacesIPv6AccessConfigsSet(c *Client, des, nw []InstanceNetworkInterfacesIPv6AccessConfigs) []InstanceNetworkInterfacesIPv6AccessConfigs {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []InstanceNetworkInterfacesIPv6AccessConfigs
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareInstanceNetworkInterfacesIPv6AccessConfigsNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewInstanceNetworkInterfacesIPv6AccessConfigs(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewInstanceNetworkInterfacesIPv6AccessConfigsSlice(c *Client, des, nw []InstanceNetworkInterfacesIPv6AccessConfigs) []InstanceNetworkInterfacesIPv6AccessConfigs {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []InstanceNetworkInterfacesIPv6AccessConfigs
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewInstanceNetworkInterfacesIPv6AccessConfigs(c, &d, &n))
	}

	return items
}

func canonicalizeInstanceNetworkInterfacesAliasIPRanges(des, initial *InstanceNetworkInterfacesAliasIPRanges, opts ...dcl.ApplyOption) *InstanceNetworkInterfacesAliasIPRanges {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &InstanceNetworkInterfacesAliasIPRanges{}

	if dcl.StringCanonicalize(des.IPCidrRange, initial.IPCidrRange) || dcl.IsZeroValue(des.IPCidrRange) {
		cDes.IPCidrRange = initial.IPCidrRange
	} else {
		cDes.IPCidrRange = des.IPCidrRange
	}
	if dcl.StringCanonicalize(des.SubnetworkRangeName, initial.SubnetworkRangeName) || dcl.IsZeroValue(des.SubnetworkRangeName) {
		cDes.SubnetworkRangeName = initial.SubnetworkRangeName
	} else {
		cDes.SubnetworkRangeName = des.SubnetworkRangeName
	}

	return cDes
}

func canonicalizeInstanceNetworkInterfacesAliasIPRangesSlice(des, initial []InstanceNetworkInterfacesAliasIPRanges, opts ...dcl.ApplyOption) []InstanceNetworkInterfacesAliasIPRanges {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]InstanceNetworkInterfacesAliasIPRanges, 0, len(des))
		for _, d := range des {
			cd := canonicalizeInstanceNetworkInterfacesAliasIPRanges(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]InstanceNetworkInterfacesAliasIPRanges, 0, len(des))
	for i, d := range des {
		cd := canonicalizeInstanceNetworkInterfacesAliasIPRanges(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewInstanceNetworkInterfacesAliasIPRanges(c *Client, des, nw *InstanceNetworkInterfacesAliasIPRanges) *InstanceNetworkInterfacesAliasIPRanges {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for InstanceNetworkInterfacesAliasIPRanges while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringCanonicalize(des.IPCidrRange, nw.IPCidrRange) {
		nw.IPCidrRange = des.IPCidrRange
	}
	if dcl.StringCanonicalize(des.SubnetworkRangeName, nw.SubnetworkRangeName) {
		nw.SubnetworkRangeName = des.SubnetworkRangeName
	}

	return nw
}

func canonicalizeNewInstanceNetworkInterfacesAliasIPRangesSet(c *Client, des, nw []InstanceNetworkInterfacesAliasIPRanges) []InstanceNetworkInterfacesAliasIPRanges {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []InstanceNetworkInterfacesAliasIPRanges
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareInstanceNetworkInterfacesAliasIPRangesNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewInstanceNetworkInterfacesAliasIPRanges(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewInstanceNetworkInterfacesAliasIPRangesSlice(c *Client, des, nw []InstanceNetworkInterfacesAliasIPRanges) []InstanceNetworkInterfacesAliasIPRanges {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []InstanceNetworkInterfacesAliasIPRanges
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewInstanceNetworkInterfacesAliasIPRanges(c, &d, &n))
	}

	return items
}

func canonicalizeInstanceScheduling(des, initial *InstanceScheduling, opts ...dcl.ApplyOption) *InstanceScheduling {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &InstanceScheduling{}

	if dcl.BoolCanonicalize(des.AutomaticRestart, initial.AutomaticRestart) || dcl.IsZeroValue(des.AutomaticRestart) {
		cDes.AutomaticRestart = initial.AutomaticRestart
	} else {
		cDes.AutomaticRestart = des.AutomaticRestart
	}
	if dcl.StringCanonicalize(des.OnHostMaintenance, initial.OnHostMaintenance) || dcl.IsZeroValue(des.OnHostMaintenance) {
		cDes.OnHostMaintenance = initial.OnHostMaintenance
	} else {
		cDes.OnHostMaintenance = des.OnHostMaintenance
	}
	if dcl.BoolCanonicalize(des.Preemptible, initial.Preemptible) || dcl.IsZeroValue(des.Preemptible) {
		cDes.Preemptible = initial.Preemptible
	} else {
		cDes.Preemptible = des.Preemptible
	}

	return cDes
}

func canonicalizeInstanceSchedulingSlice(des, initial []InstanceScheduling, opts ...dcl.ApplyOption) []InstanceScheduling {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]InstanceScheduling, 0, len(des))
		for _, d := range des {
			cd := canonicalizeInstanceScheduling(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]InstanceScheduling, 0, len(des))
	for i, d := range des {
		cd := canonicalizeInstanceScheduling(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewInstanceScheduling(c *Client, des, nw *InstanceScheduling) *InstanceScheduling {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for InstanceScheduling while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.BoolCanonicalize(des.AutomaticRestart, nw.AutomaticRestart) {
		nw.AutomaticRestart = des.AutomaticRestart
	}
	if dcl.StringCanonicalize(des.OnHostMaintenance, nw.OnHostMaintenance) {
		nw.OnHostMaintenance = des.OnHostMaintenance
	}
	if dcl.BoolCanonicalize(des.Preemptible, nw.Preemptible) {
		nw.Preemptible = des.Preemptible
	}

	return nw
}

func canonicalizeNewInstanceSchedulingSet(c *Client, des, nw []InstanceScheduling) []InstanceScheduling {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []InstanceScheduling
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareInstanceSchedulingNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewInstanceScheduling(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewInstanceSchedulingSlice(c *Client, des, nw []InstanceScheduling) []InstanceScheduling {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []InstanceScheduling
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewInstanceScheduling(c, &d, &n))
	}

	return items
}

func canonicalizeInstanceServiceAccounts(des, initial *InstanceServiceAccounts, opts ...dcl.ApplyOption) *InstanceServiceAccounts {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &InstanceServiceAccounts{}

	if dcl.StringCanonicalize(des.Email, initial.Email) || dcl.IsZeroValue(des.Email) {
		cDes.Email = initial.Email
	} else {
		cDes.Email = des.Email
	}
	if dcl.StringArrayCanonicalize(des.Scopes, initial.Scopes) {
		cDes.Scopes = initial.Scopes
	} else {
		cDes.Scopes = des.Scopes
	}

	return cDes
}

func canonicalizeInstanceServiceAccountsSlice(des, initial []InstanceServiceAccounts, opts ...dcl.ApplyOption) []InstanceServiceAccounts {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]InstanceServiceAccounts, 0, len(des))
		for _, d := range des {
			cd := canonicalizeInstanceServiceAccounts(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]InstanceServiceAccounts, 0, len(des))
	for i, d := range des {
		cd := canonicalizeInstanceServiceAccounts(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewInstanceServiceAccounts(c *Client, des, nw *InstanceServiceAccounts) *InstanceServiceAccounts {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for InstanceServiceAccounts while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringCanonicalize(des.Email, nw.Email) {
		nw.Email = des.Email
	}
	if dcl.StringArrayCanonicalize(des.Scopes, nw.Scopes) {
		nw.Scopes = des.Scopes
	}

	return nw
}

func canonicalizeNewInstanceServiceAccountsSet(c *Client, des, nw []InstanceServiceAccounts) []InstanceServiceAccounts {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []InstanceServiceAccounts
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareInstanceServiceAccountsNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewInstanceServiceAccounts(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewInstanceServiceAccountsSlice(c *Client, des, nw []InstanceServiceAccounts) []InstanceServiceAccounts {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []InstanceServiceAccounts
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewInstanceServiceAccounts(c, &d, &n))
	}

	return items
}

func canonicalizeInstanceShieldedInstanceConfig(des, initial *InstanceShieldedInstanceConfig, opts ...dcl.ApplyOption) *InstanceShieldedInstanceConfig {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &InstanceShieldedInstanceConfig{}

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

func canonicalizeInstanceShieldedInstanceConfigSlice(des, initial []InstanceShieldedInstanceConfig, opts ...dcl.ApplyOption) []InstanceShieldedInstanceConfig {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]InstanceShieldedInstanceConfig, 0, len(des))
		for _, d := range des {
			cd := canonicalizeInstanceShieldedInstanceConfig(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]InstanceShieldedInstanceConfig, 0, len(des))
	for i, d := range des {
		cd := canonicalizeInstanceShieldedInstanceConfig(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewInstanceShieldedInstanceConfig(c *Client, des, nw *InstanceShieldedInstanceConfig) *InstanceShieldedInstanceConfig {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for InstanceShieldedInstanceConfig while comparing non-nil desired to nil actual.  Returning desired object.")
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

func canonicalizeNewInstanceShieldedInstanceConfigSet(c *Client, des, nw []InstanceShieldedInstanceConfig) []InstanceShieldedInstanceConfig {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []InstanceShieldedInstanceConfig
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareInstanceShieldedInstanceConfigNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewInstanceShieldedInstanceConfig(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewInstanceShieldedInstanceConfigSlice(c *Client, des, nw []InstanceShieldedInstanceConfig) []InstanceShieldedInstanceConfig {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []InstanceShieldedInstanceConfig
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewInstanceShieldedInstanceConfig(c, &d, &n))
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
func diffInstance(c *Client, desired, actual *Instance, opts ...dcl.ApplyOption) ([]*dcl.FieldDiff, error) {
	if desired == nil || actual == nil {
		return nil, fmt.Errorf("nil resource passed to diff - always a programming error: %#v, %#v", desired, actual)
	}

	c.Config.Logger.Infof("Diff function called with desired state: %v", desired)
	c.Config.Logger.Infof("Diff function called with actual state: %v", actual)

	var fn dcl.FieldName
	var newDiffs []*dcl.FieldDiff
	// New style diffs.
	if ds, err := dcl.Diff(desired.CanIPForward, actual.CanIPForward, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("CanIpForward")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.CpuPlatform, actual.CpuPlatform, dcl.DiffInfo{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("CpuPlatform")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.CreationTimestamp, actual.CreationTimestamp, dcl.DiffInfo{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("CreationTimestamp")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.DeletionProtection, actual.DeletionProtection, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateInstanceSetDeletionProtectionOperation")}, fn.AddNest("DeletionProtection")); len(ds) != 0 || err != nil {
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

	if ds, err := dcl.Diff(desired.Disks, actual.Disks, dcl.DiffInfo{Ignore: true, ObjectFunction: compareInstanceDisksNewStyle, EmptyObject: EmptyInstanceDisks, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Disks")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.GuestAccelerators, actual.GuestAccelerators, dcl.DiffInfo{ObjectFunction: compareInstanceGuestAcceleratorsNewStyle, EmptyObject: EmptyInstanceGuestAccelerators, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("GuestAccelerators")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Hostname, actual.Hostname, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Hostname")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Id, actual.Id, dcl.DiffInfo{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Id")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Labels, actual.Labels, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateInstanceSetLabelsOperation")}, fn.AddNest("Labels")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Metadata, actual.Metadata, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateInstanceSetMetadataOperation")}, fn.AddNest("Metadata")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.MachineType, actual.MachineType, dcl.DiffInfo{Type: "ReferenceType", OperationSelector: machineTypeOperations()}, fn.AddNest("MachineType")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.MinCpuPlatform, actual.MinCpuPlatform, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("MinCpuPlatform")); len(ds) != 0 || err != nil {
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

	if ds, err := dcl.Diff(desired.NetworkInterfaces, actual.NetworkInterfaces, dcl.DiffInfo{ObjectFunction: compareInstanceNetworkInterfacesNewStyle, EmptyObject: EmptyInstanceNetworkInterfaces, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("NetworkInterfaces")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Scheduling, actual.Scheduling, dcl.DiffInfo{ServerDefault: true, ObjectFunction: compareInstanceSchedulingNewStyle, EmptyObject: EmptyInstanceScheduling, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Scheduling")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.ServiceAccounts, actual.ServiceAccounts, dcl.DiffInfo{ObjectFunction: compareInstanceServiceAccountsNewStyle, EmptyObject: EmptyInstanceServiceAccounts, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("ServiceAccounts")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.ShieldedInstanceConfig, actual.ShieldedInstanceConfig, dcl.DiffInfo{ServerDefault: true, ObjectFunction: compareInstanceShieldedInstanceConfigNewStyle, EmptyObject: EmptyInstanceShieldedInstanceConfig, OperationSelector: dcl.TriggersOperation("updateInstanceUpdateShieldedInstanceConfigOperation")}, fn.AddNest("ShieldedInstanceConfig")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Status, actual.Status, dcl.DiffInfo{ServerDefault: true, Type: "EnumType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Status")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.StatusMessage, actual.StatusMessage, dcl.DiffInfo{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("StatusMessage")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Tags, actual.Tags, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateInstanceSetTagsOperation")}, fn.AddNest("Tags")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Zone, actual.Zone, dcl.DiffInfo{Type: "ReferenceType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Zone")); len(ds) != 0 || err != nil {
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

	if ds, err := dcl.Diff(desired.SelfLink, actual.SelfLink, dcl.DiffInfo{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("SelfLink")); len(ds) != 0 || err != nil {
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
func compareInstanceDisksNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*InstanceDisks)
	if !ok {
		desiredNotPointer, ok := d.(InstanceDisks)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a InstanceDisks or *InstanceDisks", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*InstanceDisks)
	if !ok {
		actualNotPointer, ok := a.(InstanceDisks)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a InstanceDisks", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.AutoDelete, actual.AutoDelete, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("AutoDelete")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Boot, actual.Boot, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Boot")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.DeviceName, actual.DeviceName, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("DeviceName")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.DiskEncryptionKey, actual.DiskEncryptionKey, dcl.DiffInfo{ObjectFunction: compareInstanceDisksDiskEncryptionKeyNewStyle, EmptyObject: EmptyInstanceDisksDiskEncryptionKey, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("DiskEncryptionKey")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Index, actual.Index, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Index")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.InitializeParams, actual.InitializeParams, dcl.DiffInfo{Ignore: true, ObjectFunction: compareInstanceDisksInitializeParamsNewStyle, EmptyObject: EmptyInstanceDisksInitializeParams, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("InitializeParams")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Interface, actual.Interface, dcl.DiffInfo{Type: "EnumType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Interface")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Mode, actual.Mode, dcl.DiffInfo{Type: "EnumType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Mode")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Source, actual.Source, dcl.DiffInfo{Type: "ReferenceType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Source")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Type, actual.Type, dcl.DiffInfo{Type: "EnumType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Type")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareInstanceDisksDiskEncryptionKeyNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*InstanceDisksDiskEncryptionKey)
	if !ok {
		desiredNotPointer, ok := d.(InstanceDisksDiskEncryptionKey)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a InstanceDisksDiskEncryptionKey or *InstanceDisksDiskEncryptionKey", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*InstanceDisksDiskEncryptionKey)
	if !ok {
		actualNotPointer, ok := a.(InstanceDisksDiskEncryptionKey)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a InstanceDisksDiskEncryptionKey", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.RawKey, actual.RawKey, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("RawKey")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.RsaEncryptedKey, actual.RsaEncryptedKey, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("RsaEncryptedKey")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Sha256, actual.Sha256, dcl.DiffInfo{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Sha256")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareInstanceDisksInitializeParamsNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*InstanceDisksInitializeParams)
	if !ok {
		desiredNotPointer, ok := d.(InstanceDisksInitializeParams)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a InstanceDisksInitializeParams or *InstanceDisksInitializeParams", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*InstanceDisksInitializeParams)
	if !ok {
		actualNotPointer, ok := a.(InstanceDisksInitializeParams)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a InstanceDisksInitializeParams", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.DiskName, actual.DiskName, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("DiskName")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.DiskSizeGb, actual.DiskSizeGb, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("DiskSizeGb")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.DiskType, actual.DiskType, dcl.DiffInfo{Type: "ReferenceType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("DiskType")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.SourceImage, actual.SourceImage, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("SourceImage")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.SourceImageEncryptionKey, actual.SourceImageEncryptionKey, dcl.DiffInfo{ObjectFunction: compareInstanceDisksInitializeParamsSourceImageEncryptionKeyNewStyle, EmptyObject: EmptyInstanceDisksInitializeParamsSourceImageEncryptionKey, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("SourceImageEncryptionKey")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareInstanceDisksInitializeParamsSourceImageEncryptionKeyNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*InstanceDisksInitializeParamsSourceImageEncryptionKey)
	if !ok {
		desiredNotPointer, ok := d.(InstanceDisksInitializeParamsSourceImageEncryptionKey)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a InstanceDisksInitializeParamsSourceImageEncryptionKey or *InstanceDisksInitializeParamsSourceImageEncryptionKey", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*InstanceDisksInitializeParamsSourceImageEncryptionKey)
	if !ok {
		actualNotPointer, ok := a.(InstanceDisksInitializeParamsSourceImageEncryptionKey)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a InstanceDisksInitializeParamsSourceImageEncryptionKey", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.RawKey, actual.RawKey, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("RawKey")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Sha256, actual.Sha256, dcl.DiffInfo{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Sha256")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareInstanceGuestAcceleratorsNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*InstanceGuestAccelerators)
	if !ok {
		desiredNotPointer, ok := d.(InstanceGuestAccelerators)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a InstanceGuestAccelerators or *InstanceGuestAccelerators", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*InstanceGuestAccelerators)
	if !ok {
		actualNotPointer, ok := a.(InstanceGuestAccelerators)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a InstanceGuestAccelerators", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.AcceleratorCount, actual.AcceleratorCount, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("AcceleratorCount")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.AcceleratorType, actual.AcceleratorType, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("AcceleratorType")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareInstanceNetworkInterfacesNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*InstanceNetworkInterfaces)
	if !ok {
		desiredNotPointer, ok := d.(InstanceNetworkInterfaces)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a InstanceNetworkInterfaces or *InstanceNetworkInterfaces", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*InstanceNetworkInterfaces)
	if !ok {
		actualNotPointer, ok := a.(InstanceNetworkInterfaces)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a InstanceNetworkInterfaces", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.AccessConfigs, actual.AccessConfigs, dcl.DiffInfo{ObjectFunction: compareInstanceNetworkInterfacesAccessConfigsNewStyle, EmptyObject: EmptyInstanceNetworkInterfacesAccessConfigs, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("AccessConfigs")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.IPv6AccessConfigs, actual.IPv6AccessConfigs, dcl.DiffInfo{ObjectFunction: compareInstanceNetworkInterfacesIPv6AccessConfigsNewStyle, EmptyObject: EmptyInstanceNetworkInterfacesIPv6AccessConfigs, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Ipv6AccessConfigs")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.AliasIPRanges, actual.AliasIPRanges, dcl.DiffInfo{ObjectFunction: compareInstanceNetworkInterfacesAliasIPRangesNewStyle, EmptyObject: EmptyInstanceNetworkInterfacesAliasIPRanges, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("AliasIPRanges")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Name, actual.Name, dcl.DiffInfo{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Name")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Network, actual.Network, dcl.DiffInfo{ServerDefault: true, Type: "ReferenceType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Network")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.NetworkIP, actual.NetworkIP, dcl.DiffInfo{ServerDefault: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("NetworkIP")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Subnetwork, actual.Subnetwork, dcl.DiffInfo{ServerDefault: true, Type: "ReferenceType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Subnetwork")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareInstanceNetworkInterfacesAccessConfigsNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*InstanceNetworkInterfacesAccessConfigs)
	if !ok {
		desiredNotPointer, ok := d.(InstanceNetworkInterfacesAccessConfigs)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a InstanceNetworkInterfacesAccessConfigs or *InstanceNetworkInterfacesAccessConfigs", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*InstanceNetworkInterfacesAccessConfigs)
	if !ok {
		actualNotPointer, ok := a.(InstanceNetworkInterfacesAccessConfigs)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a InstanceNetworkInterfacesAccessConfigs", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.Name, actual.Name, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Name")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.NatIP, actual.NatIP, dcl.DiffInfo{ServerDefault: true, Type: "ReferenceType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("NatIP")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.ExternalIPv6, actual.ExternalIPv6, dcl.DiffInfo{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("ExternalIPv6")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.ExternalIPv6PrefixLength, actual.ExternalIPv6PrefixLength, dcl.DiffInfo{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("ExternalIPv6PrefixLength")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.SetPublicPtr, actual.SetPublicPtr, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("SetPublicPtr")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.PublicPtrDomainName, actual.PublicPtrDomainName, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("PublicPtrDomainName")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.NetworkTier, actual.NetworkTier, dcl.DiffInfo{ServerDefault: true, Type: "EnumType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("NetworkTier")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Type, actual.Type, dcl.DiffInfo{Type: "EnumType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Type")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareInstanceNetworkInterfacesIPv6AccessConfigsNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*InstanceNetworkInterfacesIPv6AccessConfigs)
	if !ok {
		desiredNotPointer, ok := d.(InstanceNetworkInterfacesIPv6AccessConfigs)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a InstanceNetworkInterfacesIPv6AccessConfigs or *InstanceNetworkInterfacesIPv6AccessConfigs", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*InstanceNetworkInterfacesIPv6AccessConfigs)
	if !ok {
		actualNotPointer, ok := a.(InstanceNetworkInterfacesIPv6AccessConfigs)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a InstanceNetworkInterfacesIPv6AccessConfigs", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.Name, actual.Name, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Name")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.NatIP, actual.NatIP, dcl.DiffInfo{ServerDefault: true, Type: "ReferenceType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("NatIP")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.ExternalIPv6, actual.ExternalIPv6, dcl.DiffInfo{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("ExternalIPv6")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.ExternalIPv6PrefixLength, actual.ExternalIPv6PrefixLength, dcl.DiffInfo{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("ExternalIPv6PrefixLength")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.SetPublicPtr, actual.SetPublicPtr, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("SetPublicPtr")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.PublicPtrDomainName, actual.PublicPtrDomainName, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("PublicPtrDomainName")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.NetworkTier, actual.NetworkTier, dcl.DiffInfo{Type: "EnumType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("NetworkTier")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Type, actual.Type, dcl.DiffInfo{Type: "EnumType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Type")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareInstanceNetworkInterfacesAliasIPRangesNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*InstanceNetworkInterfacesAliasIPRanges)
	if !ok {
		desiredNotPointer, ok := d.(InstanceNetworkInterfacesAliasIPRanges)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a InstanceNetworkInterfacesAliasIPRanges or *InstanceNetworkInterfacesAliasIPRanges", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*InstanceNetworkInterfacesAliasIPRanges)
	if !ok {
		actualNotPointer, ok := a.(InstanceNetworkInterfacesAliasIPRanges)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a InstanceNetworkInterfacesAliasIPRanges", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.IPCidrRange, actual.IPCidrRange, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("IpCidrRange")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.SubnetworkRangeName, actual.SubnetworkRangeName, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("SubnetworkRangeName")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareInstanceSchedulingNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*InstanceScheduling)
	if !ok {
		desiredNotPointer, ok := d.(InstanceScheduling)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a InstanceScheduling or *InstanceScheduling", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*InstanceScheduling)
	if !ok {
		actualNotPointer, ok := a.(InstanceScheduling)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a InstanceScheduling", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.AutomaticRestart, actual.AutomaticRestart, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("AutomaticRestart")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.OnHostMaintenance, actual.OnHostMaintenance, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("OnHostMaintenance")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Preemptible, actual.Preemptible, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Preemptible")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareInstanceServiceAccountsNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*InstanceServiceAccounts)
	if !ok {
		desiredNotPointer, ok := d.(InstanceServiceAccounts)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a InstanceServiceAccounts or *InstanceServiceAccounts", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*InstanceServiceAccounts)
	if !ok {
		actualNotPointer, ok := a.(InstanceServiceAccounts)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a InstanceServiceAccounts", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.Email, actual.Email, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Email")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Scopes, actual.Scopes, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Scopes")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareInstanceShieldedInstanceConfigNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*InstanceShieldedInstanceConfig)
	if !ok {
		desiredNotPointer, ok := d.(InstanceShieldedInstanceConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a InstanceShieldedInstanceConfig or *InstanceShieldedInstanceConfig", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*InstanceShieldedInstanceConfig)
	if !ok {
		actualNotPointer, ok := a.(InstanceShieldedInstanceConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a InstanceShieldedInstanceConfig", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.EnableSecureBoot, actual.EnableSecureBoot, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateInstanceUpdateShieldedInstanceConfigOperation")}, fn.AddNest("EnableSecureBoot")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.EnableVtpm, actual.EnableVtpm, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateInstanceUpdateShieldedInstanceConfigOperation")}, fn.AddNest("EnableVtpm")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.EnableIntegrityMonitoring, actual.EnableIntegrityMonitoring, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateInstanceUpdateShieldedInstanceConfigOperation")}, fn.AddNest("EnableIntegrityMonitoring")); len(ds) != 0 || err != nil {
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
func (r *Instance) urlNormalized() *Instance {
	normalized := dcl.Copy(*r).(Instance)
	normalized.CpuPlatform = dcl.SelfLinkToName(r.CpuPlatform)
	normalized.CreationTimestamp = dcl.SelfLinkToName(r.CreationTimestamp)
	normalized.Description = dcl.SelfLinkToName(r.Description)
	normalized.Hostname = dcl.SelfLinkToName(r.Hostname)
	normalized.Id = dcl.SelfLinkToName(r.Id)
	normalized.MachineType = r.MachineType
	normalized.MinCpuPlatform = dcl.SelfLinkToName(r.MinCpuPlatform)
	normalized.Name = dcl.SelfLinkToName(r.Name)
	normalized.StatusMessage = dcl.SelfLinkToName(r.StatusMessage)
	normalized.Zone = dcl.SelfLinkToName(r.Zone)
	normalized.Project = dcl.SelfLinkToName(r.Project)
	normalized.SelfLink = dcl.SelfLinkToName(r.SelfLink)
	return &normalized
}

func (r *Instance) updateURL(userBasePath, updateName string) (string, error) {
	nr := r.urlNormalized()
	if updateName == "setDeletionProtection" {
		fields := map[string]interface{}{
			"project": dcl.ValueOrEmptyString(nr.Project),
			"zone":    dcl.ValueOrEmptyString(nr.Zone),
		}
		return dcl.URL("/projects/{{project}}/zones/{{zone}}/instances/{resourceId}/setDeletionProtection", nr.basePath(), userBasePath, fields), nil

	}
	if updateName == "setLabels" {
		fields := map[string]interface{}{
			"project": dcl.ValueOrEmptyString(nr.Project),
			"zone":    dcl.ValueOrEmptyString(nr.Zone),
			"name":    dcl.ValueOrEmptyString(nr.Name),
		}
		return dcl.URL("projects/{{project}}/zones/{{zone}}/instances/{{name}}/setLabels", nr.basePath(), userBasePath, fields), nil

	}
	if updateName == "setMachineType" {
		fields := map[string]interface{}{
			"project": dcl.ValueOrEmptyString(nr.Project),
			"zone":    dcl.ValueOrEmptyString(nr.Zone),
			"name":    dcl.ValueOrEmptyString(nr.Name),
		}
		return dcl.URL("projects/{{project}}/zones/{{zone}}/instances/{{name}}/setMachineType", nr.basePath(), userBasePath, fields), nil

	}
	if updateName == "setMetadata" {
		fields := map[string]interface{}{
			"project": dcl.ValueOrEmptyString(nr.Project),
			"zone":    dcl.ValueOrEmptyString(nr.Zone),
			"name":    dcl.ValueOrEmptyString(nr.Name),
		}
		return dcl.URL("projects/{{project}}/zones/{{zone}}/instances/{{name}}/setMetadata", nr.basePath(), userBasePath, fields), nil

	}
	if updateName == "setTags" {
		fields := map[string]interface{}{
			"project": dcl.ValueOrEmptyString(nr.Project),
			"zone":    dcl.ValueOrEmptyString(nr.Zone),
			"name":    dcl.ValueOrEmptyString(nr.Name),
		}
		return dcl.URL("projects/{{project}}/zones/{{zone}}/instances/{{name}}/setTags", nr.basePath(), userBasePath, fields), nil

	}
	if updateName == "start" {
		fields := map[string]interface{}{
			"project": dcl.ValueOrEmptyString(nr.Project),
			"zone":    dcl.ValueOrEmptyString(nr.Zone),
			"name":    dcl.ValueOrEmptyString(nr.Name),
		}
		return dcl.URL("projects/{{project}}/zones/{{zone}}/instances/{{name}}/start", nr.basePath(), userBasePath, fields), nil

	}
	if updateName == "stop" {
		fields := map[string]interface{}{
			"project": dcl.ValueOrEmptyString(nr.Project),
			"zone":    dcl.ValueOrEmptyString(nr.Zone),
			"name":    dcl.ValueOrEmptyString(nr.Name),
		}
		return dcl.URL("projects/{{project}}/zones/{{zone}}/instances/{{name}}/stop", nr.basePath(), userBasePath, fields), nil

	}
	if updateName == "update" {
		fields := map[string]interface{}{
			"project": dcl.ValueOrEmptyString(nr.Project),
			"zone":    dcl.ValueOrEmptyString(nr.Zone),
			"name":    dcl.ValueOrEmptyString(nr.Name),
		}
		return dcl.URL("projects/{{project}}/zones/{{zone}}/instances/{{name}}", nr.basePath(), userBasePath, fields), nil

	}
	if updateName == "updateShieldedInstanceConfig" {
		fields := map[string]interface{}{
			"project": dcl.ValueOrEmptyString(nr.Project),
			"name":    dcl.ValueOrEmptyString(nr.Name),
		}
		return dcl.URL("projects/{{project}}/instances/{{name}}/updateShieldedInstanceConfig", nr.basePath(), userBasePath, fields), nil

	}

	return "", fmt.Errorf("unknown update name: %s", updateName)
}

// marshal encodes the Instance resource into JSON for a Create request, and
// performs transformations from the resource schema to the API schema if
// necessary.
func (r *Instance) marshal(c *Client) ([]byte, error) {
	m, err := expandInstance(c, r)
	if err != nil {
		return nil, fmt.Errorf("error marshalling Instance: %w", err)
	}
	dcl.MoveMapEntry(
		m,
		[]string{"tags"},
		[]string{"tags", "items"},
	)

	return json.Marshal(m)
}

// unmarshalInstance decodes JSON responses into the Instance resource schema.
func unmarshalInstance(b []byte, c *Client, res *Instance) (*Instance, error) {
	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, err
	}
	return unmarshalMapInstance(m, c, res)
}

func unmarshalMapInstance(m map[string]interface{}, c *Client, res *Instance) (*Instance, error) {
	if v, err := dcl.MapFromListOfKeyValues(m, []string{"metadata", "items"}, "key", "value"); err != nil {
		return nil, err
	} else {
		dcl.PutMapEntry(
			m,
			[]string{"metadata"},
			v,
		)
	}
	dcl.MoveMapEntry(
		m,
		[]string{"tags", "items"},
		[]string{"tags"},
	)

	flattened := flattenInstance(c, m, res)
	if flattened == nil {
		return nil, fmt.Errorf("attempted to flatten empty json object")
	}
	return flattened, nil
}

// expandInstance expands Instance into a JSON request object.
func expandInstance(c *Client, f *Instance) (map[string]interface{}, error) {
	m := make(map[string]interface{})
	res := f
	_ = res
	if v := f.CanIPForward; dcl.ValueShouldBeSent(v) {
		m["canIpForward"] = v
	}
	if v := f.DeletionProtection; dcl.ValueShouldBeSent(v) {
		m["deletionProtection"] = v
	}
	if v := f.Description; dcl.ValueShouldBeSent(v) {
		m["description"] = v
	}
	if v, err := expandInstanceDisksSlice(c, f.Disks, res); err != nil {
		return nil, fmt.Errorf("error expanding Disks into disks: %w", err)
	} else if v != nil {
		m["disks"] = v
	}
	if v, err := expandInstanceGuestAcceleratorsSlice(c, f.GuestAccelerators, res); err != nil {
		return nil, fmt.Errorf("error expanding GuestAccelerators into guestAccelerators: %w", err)
	} else if v != nil {
		m["guestAccelerators"] = v
	}
	if v := f.Hostname; dcl.ValueShouldBeSent(v) {
		m["hostname"] = v
	}
	if v := f.Labels; dcl.ValueShouldBeSent(v) {
		m["labels"] = v
	}
	if v, err := dcl.ListOfKeyValuesFromMapInStruct(f.Metadata, "items", "key", "value"); err != nil {
		return nil, fmt.Errorf("error expanding Metadata into metadata: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["metadata"] = v
	}
	if v, err := dcl.DeriveFromPattern("zones/%s/machineTypes/%s", f.MachineType, dcl.SelfLinkToName(f.Zone), f.MachineType); err != nil {
		return nil, fmt.Errorf("error expanding MachineType into machineType: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["machineType"] = v
	}
	if v := f.MinCpuPlatform; dcl.ValueShouldBeSent(v) {
		m["minCpuPlatform"] = v
	}
	if v := f.Name; dcl.ValueShouldBeSent(v) {
		m["name"] = v
	}
	if v, err := expandInstanceNetworkInterfacesSlice(c, f.NetworkInterfaces, res); err != nil {
		return nil, fmt.Errorf("error expanding NetworkInterfaces into networkInterfaces: %w", err)
	} else if v != nil {
		m["networkInterfaces"] = v
	}
	if v, err := expandInstanceScheduling(c, f.Scheduling, res); err != nil {
		return nil, fmt.Errorf("error expanding Scheduling into scheduling: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["scheduling"] = v
	}
	if v, err := expandInstanceServiceAccountsSlice(c, f.ServiceAccounts, res); err != nil {
		return nil, fmt.Errorf("error expanding ServiceAccounts into serviceAccounts: %w", err)
	} else if v != nil {
		m["serviceAccounts"] = v
	}
	if v, err := expandInstanceShieldedInstanceConfig(c, f.ShieldedInstanceConfig, res); err != nil {
		return nil, fmt.Errorf("error expanding ShieldedInstanceConfig into shieldedInstanceConfig: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["shieldedInstanceConfig"] = v
	}
	if v := f.Status; dcl.ValueShouldBeSent(v) {
		m["status"] = v
	}
	if v := f.Tags; v != nil {
		m["tags"] = v
	}
	if v := f.Zone; dcl.ValueShouldBeSent(v) {
		m["zone"] = v
	}
	if v, err := dcl.EmptyValue(); err != nil {
		return nil, fmt.Errorf("error expanding Project into project: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["project"] = v
	}

	return m, nil
}

// flattenInstance flattens Instance from a JSON request object into the
// Instance type.
func flattenInstance(c *Client, i interface{}, res *Instance) *Instance {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}
	if len(m) == 0 {
		return nil
	}

	resultRes := &Instance{}
	resultRes.CanIPForward = dcl.FlattenBool(m["canIpForward"])
	resultRes.CpuPlatform = dcl.FlattenString(m["cpuPlatform"])
	resultRes.CreationTimestamp = dcl.FlattenString(m["creationTimestamp"])
	resultRes.DeletionProtection = dcl.FlattenBool(m["deletionProtection"])
	resultRes.Description = dcl.FlattenString(m["description"])
	resultRes.Disks = flattenInstanceDisksSlice(c, m["disks"], res)
	resultRes.GuestAccelerators = flattenInstanceGuestAcceleratorsSlice(c, m["guestAccelerators"], res)
	resultRes.Hostname = dcl.FlattenString(m["hostname"])
	resultRes.Id = dcl.FlattenString(m["id"])
	resultRes.Labels = dcl.FlattenKeyValuePairs(m["labels"])
	resultRes.Metadata = dcl.FlattenKeyValuePairs(m["metadata"])
	resultRes.MachineType = dcl.FlattenString(m["machineType"])
	resultRes.MinCpuPlatform = dcl.FlattenString(m["minCpuPlatform"])
	resultRes.Name = dcl.FlattenString(m["name"])
	resultRes.NetworkInterfaces = flattenInstanceNetworkInterfacesSlice(c, m["networkInterfaces"], res)
	resultRes.Scheduling = flattenInstanceScheduling(c, m["scheduling"], res)
	resultRes.ServiceAccounts = flattenInstanceServiceAccountsSlice(c, m["serviceAccounts"], res)
	resultRes.ShieldedInstanceConfig = flattenInstanceShieldedInstanceConfig(c, m["shieldedInstanceConfig"], res)
	resultRes.Status = flattenInstanceStatusEnum(m["status"])
	resultRes.StatusMessage = dcl.FlattenString(m["statusMessage"])
	resultRes.Tags = dcl.FlattenStringSlice(m["tags"])
	resultRes.Zone = dcl.FlattenString(m["zone"])
	resultRes.Project = dcl.FlattenString(m["project"])
	resultRes.SelfLink = dcl.FlattenString(m["selfLink"])

	return resultRes
}

// expandInstanceDisksMap expands the contents of InstanceDisks into a JSON
// request object.
func expandInstanceDisksMap(c *Client, f map[string]InstanceDisks, res *Instance) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandInstanceDisks(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandInstanceDisksSlice expands the contents of InstanceDisks into a JSON
// request object.
func expandInstanceDisksSlice(c *Client, f []InstanceDisks, res *Instance) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandInstanceDisks(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenInstanceDisksMap flattens the contents of InstanceDisks from a JSON
// response object.
func flattenInstanceDisksMap(c *Client, i interface{}, res *Instance) map[string]InstanceDisks {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]InstanceDisks{}
	}

	if len(a) == 0 {
		return map[string]InstanceDisks{}
	}

	items := make(map[string]InstanceDisks)
	for k, item := range a {
		items[k] = *flattenInstanceDisks(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenInstanceDisksSlice flattens the contents of InstanceDisks from a JSON
// response object.
func flattenInstanceDisksSlice(c *Client, i interface{}, res *Instance) []InstanceDisks {
	a, ok := i.([]interface{})
	if !ok {
		return []InstanceDisks{}
	}

	if len(a) == 0 {
		return []InstanceDisks{}
	}

	items := make([]InstanceDisks, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenInstanceDisks(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandInstanceDisks expands an instance of InstanceDisks into a JSON
// request object.
func expandInstanceDisks(c *Client, f *InstanceDisks, res *Instance) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.AutoDelete; !dcl.IsEmptyValueIndirect(v) {
		m["autoDelete"] = v
	}
	if v := f.Boot; !dcl.IsEmptyValueIndirect(v) {
		m["boot"] = v
	}
	if v := f.DeviceName; !dcl.IsEmptyValueIndirect(v) {
		m["deviceName"] = v
	}
	if v, err := expandInstanceDisksDiskEncryptionKey(c, f.DiskEncryptionKey, res); err != nil {
		return nil, fmt.Errorf("error expanding DiskEncryptionKey into diskEncryptionKey: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["diskEncryptionKey"] = v
	}
	if v := f.Index; !dcl.IsEmptyValueIndirect(v) {
		m["index"] = v
	}
	if v, err := expandInstanceDisksInitializeParams(c, f.InitializeParams, res); err != nil {
		return nil, fmt.Errorf("error expanding InitializeParams into initializeParams: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["initializeParams"] = v
	}
	if v := f.Interface; !dcl.IsEmptyValueIndirect(v) {
		m["interface"] = v
	}
	if v := f.Mode; !dcl.IsEmptyValueIndirect(v) {
		m["mode"] = v
	}
	if v := f.Source; !dcl.IsEmptyValueIndirect(v) {
		m["source"] = v
	}
	if v := f.Type; !dcl.IsEmptyValueIndirect(v) {
		m["type"] = v
	}

	return m, nil
}

// flattenInstanceDisks flattens an instance of InstanceDisks from a JSON
// response object.
func flattenInstanceDisks(c *Client, i interface{}, res *Instance) *InstanceDisks {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &InstanceDisks{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyInstanceDisks
	}
	r.AutoDelete = dcl.FlattenBool(m["autoDelete"])
	r.Boot = dcl.FlattenBool(m["boot"])
	r.DeviceName = dcl.FlattenString(m["deviceName"])
	r.DiskEncryptionKey = flattenInstanceDisksDiskEncryptionKey(c, m["diskEncryptionKey"], res)
	r.Index = dcl.FlattenInteger(m["index"])
	r.InitializeParams = flattenInstanceDisksInitializeParams(c, m["initializeParams"], res)
	r.Interface = flattenInstanceDisksInterfaceEnum(m["interface"])
	r.Mode = flattenInstanceDisksModeEnum(m["mode"])
	r.Source = dcl.FlattenString(m["source"])
	r.Type = flattenInstanceDisksTypeEnum(m["type"])

	return r
}

// expandInstanceDisksDiskEncryptionKeyMap expands the contents of InstanceDisksDiskEncryptionKey into a JSON
// request object.
func expandInstanceDisksDiskEncryptionKeyMap(c *Client, f map[string]InstanceDisksDiskEncryptionKey, res *Instance) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandInstanceDisksDiskEncryptionKey(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandInstanceDisksDiskEncryptionKeySlice expands the contents of InstanceDisksDiskEncryptionKey into a JSON
// request object.
func expandInstanceDisksDiskEncryptionKeySlice(c *Client, f []InstanceDisksDiskEncryptionKey, res *Instance) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandInstanceDisksDiskEncryptionKey(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenInstanceDisksDiskEncryptionKeyMap flattens the contents of InstanceDisksDiskEncryptionKey from a JSON
// response object.
func flattenInstanceDisksDiskEncryptionKeyMap(c *Client, i interface{}, res *Instance) map[string]InstanceDisksDiskEncryptionKey {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]InstanceDisksDiskEncryptionKey{}
	}

	if len(a) == 0 {
		return map[string]InstanceDisksDiskEncryptionKey{}
	}

	items := make(map[string]InstanceDisksDiskEncryptionKey)
	for k, item := range a {
		items[k] = *flattenInstanceDisksDiskEncryptionKey(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenInstanceDisksDiskEncryptionKeySlice flattens the contents of InstanceDisksDiskEncryptionKey from a JSON
// response object.
func flattenInstanceDisksDiskEncryptionKeySlice(c *Client, i interface{}, res *Instance) []InstanceDisksDiskEncryptionKey {
	a, ok := i.([]interface{})
	if !ok {
		return []InstanceDisksDiskEncryptionKey{}
	}

	if len(a) == 0 {
		return []InstanceDisksDiskEncryptionKey{}
	}

	items := make([]InstanceDisksDiskEncryptionKey, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenInstanceDisksDiskEncryptionKey(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandInstanceDisksDiskEncryptionKey expands an instance of InstanceDisksDiskEncryptionKey into a JSON
// request object.
func expandInstanceDisksDiskEncryptionKey(c *Client, f *InstanceDisksDiskEncryptionKey, res *Instance) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.RawKey; !dcl.IsEmptyValueIndirect(v) {
		m["rawKey"] = v
	}
	if v := f.RsaEncryptedKey; !dcl.IsEmptyValueIndirect(v) {
		m["rsaEncryptedKey"] = v
	}

	return m, nil
}

// flattenInstanceDisksDiskEncryptionKey flattens an instance of InstanceDisksDiskEncryptionKey from a JSON
// response object.
func flattenInstanceDisksDiskEncryptionKey(c *Client, i interface{}, res *Instance) *InstanceDisksDiskEncryptionKey {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &InstanceDisksDiskEncryptionKey{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyInstanceDisksDiskEncryptionKey
	}
	r.RawKey = dcl.FlattenString(m["rawKey"])
	r.RsaEncryptedKey = dcl.FlattenString(m["rsaEncryptedKey"])
	r.Sha256 = dcl.FlattenString(m["sha256"])

	return r
}

// expandInstanceDisksInitializeParamsMap expands the contents of InstanceDisksInitializeParams into a JSON
// request object.
func expandInstanceDisksInitializeParamsMap(c *Client, f map[string]InstanceDisksInitializeParams, res *Instance) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandInstanceDisksInitializeParams(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandInstanceDisksInitializeParamsSlice expands the contents of InstanceDisksInitializeParams into a JSON
// request object.
func expandInstanceDisksInitializeParamsSlice(c *Client, f []InstanceDisksInitializeParams, res *Instance) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandInstanceDisksInitializeParams(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenInstanceDisksInitializeParamsMap flattens the contents of InstanceDisksInitializeParams from a JSON
// response object.
func flattenInstanceDisksInitializeParamsMap(c *Client, i interface{}, res *Instance) map[string]InstanceDisksInitializeParams {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]InstanceDisksInitializeParams{}
	}

	if len(a) == 0 {
		return map[string]InstanceDisksInitializeParams{}
	}

	items := make(map[string]InstanceDisksInitializeParams)
	for k, item := range a {
		items[k] = *flattenInstanceDisksInitializeParams(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenInstanceDisksInitializeParamsSlice flattens the contents of InstanceDisksInitializeParams from a JSON
// response object.
func flattenInstanceDisksInitializeParamsSlice(c *Client, i interface{}, res *Instance) []InstanceDisksInitializeParams {
	a, ok := i.([]interface{})
	if !ok {
		return []InstanceDisksInitializeParams{}
	}

	if len(a) == 0 {
		return []InstanceDisksInitializeParams{}
	}

	items := make([]InstanceDisksInitializeParams, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenInstanceDisksInitializeParams(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandInstanceDisksInitializeParams expands an instance of InstanceDisksInitializeParams into a JSON
// request object.
func expandInstanceDisksInitializeParams(c *Client, f *InstanceDisksInitializeParams, res *Instance) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.DiskName; !dcl.IsEmptyValueIndirect(v) {
		m["diskName"] = v
	}
	if v := f.DiskSizeGb; !dcl.IsEmptyValueIndirect(v) {
		m["diskSizeGb"] = v
	}
	if v := f.DiskType; !dcl.IsEmptyValueIndirect(v) {
		m["diskType"] = v
	}
	if v := f.SourceImage; !dcl.IsEmptyValueIndirect(v) {
		m["sourceImage"] = v
	}
	if v, err := expandInstanceDisksInitializeParamsSourceImageEncryptionKey(c, f.SourceImageEncryptionKey, res); err != nil {
		return nil, fmt.Errorf("error expanding SourceImageEncryptionKey into sourceImageEncryptionKey: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["sourceImageEncryptionKey"] = v
	}

	return m, nil
}

// flattenInstanceDisksInitializeParams flattens an instance of InstanceDisksInitializeParams from a JSON
// response object.
func flattenInstanceDisksInitializeParams(c *Client, i interface{}, res *Instance) *InstanceDisksInitializeParams {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &InstanceDisksInitializeParams{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyInstanceDisksInitializeParams
	}
	r.DiskName = dcl.FlattenString(m["diskName"])
	r.DiskSizeGb = dcl.FlattenInteger(m["diskSizeGb"])
	r.DiskType = dcl.FlattenString(m["diskType"])
	r.SourceImage = dcl.FlattenString(m["sourceImage"])
	r.SourceImageEncryptionKey = flattenInstanceDisksInitializeParamsSourceImageEncryptionKey(c, m["sourceImageEncryptionKey"], res)

	return r
}

// expandInstanceDisksInitializeParamsSourceImageEncryptionKeyMap expands the contents of InstanceDisksInitializeParamsSourceImageEncryptionKey into a JSON
// request object.
func expandInstanceDisksInitializeParamsSourceImageEncryptionKeyMap(c *Client, f map[string]InstanceDisksInitializeParamsSourceImageEncryptionKey, res *Instance) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandInstanceDisksInitializeParamsSourceImageEncryptionKey(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandInstanceDisksInitializeParamsSourceImageEncryptionKeySlice expands the contents of InstanceDisksInitializeParamsSourceImageEncryptionKey into a JSON
// request object.
func expandInstanceDisksInitializeParamsSourceImageEncryptionKeySlice(c *Client, f []InstanceDisksInitializeParamsSourceImageEncryptionKey, res *Instance) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandInstanceDisksInitializeParamsSourceImageEncryptionKey(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenInstanceDisksInitializeParamsSourceImageEncryptionKeyMap flattens the contents of InstanceDisksInitializeParamsSourceImageEncryptionKey from a JSON
// response object.
func flattenInstanceDisksInitializeParamsSourceImageEncryptionKeyMap(c *Client, i interface{}, res *Instance) map[string]InstanceDisksInitializeParamsSourceImageEncryptionKey {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]InstanceDisksInitializeParamsSourceImageEncryptionKey{}
	}

	if len(a) == 0 {
		return map[string]InstanceDisksInitializeParamsSourceImageEncryptionKey{}
	}

	items := make(map[string]InstanceDisksInitializeParamsSourceImageEncryptionKey)
	for k, item := range a {
		items[k] = *flattenInstanceDisksInitializeParamsSourceImageEncryptionKey(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenInstanceDisksInitializeParamsSourceImageEncryptionKeySlice flattens the contents of InstanceDisksInitializeParamsSourceImageEncryptionKey from a JSON
// response object.
func flattenInstanceDisksInitializeParamsSourceImageEncryptionKeySlice(c *Client, i interface{}, res *Instance) []InstanceDisksInitializeParamsSourceImageEncryptionKey {
	a, ok := i.([]interface{})
	if !ok {
		return []InstanceDisksInitializeParamsSourceImageEncryptionKey{}
	}

	if len(a) == 0 {
		return []InstanceDisksInitializeParamsSourceImageEncryptionKey{}
	}

	items := make([]InstanceDisksInitializeParamsSourceImageEncryptionKey, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenInstanceDisksInitializeParamsSourceImageEncryptionKey(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandInstanceDisksInitializeParamsSourceImageEncryptionKey expands an instance of InstanceDisksInitializeParamsSourceImageEncryptionKey into a JSON
// request object.
func expandInstanceDisksInitializeParamsSourceImageEncryptionKey(c *Client, f *InstanceDisksInitializeParamsSourceImageEncryptionKey, res *Instance) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.RawKey; !dcl.IsEmptyValueIndirect(v) {
		m["rawKey"] = v
	}

	return m, nil
}

// flattenInstanceDisksInitializeParamsSourceImageEncryptionKey flattens an instance of InstanceDisksInitializeParamsSourceImageEncryptionKey from a JSON
// response object.
func flattenInstanceDisksInitializeParamsSourceImageEncryptionKey(c *Client, i interface{}, res *Instance) *InstanceDisksInitializeParamsSourceImageEncryptionKey {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &InstanceDisksInitializeParamsSourceImageEncryptionKey{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyInstanceDisksInitializeParamsSourceImageEncryptionKey
	}
	r.RawKey = dcl.FlattenString(m["rawKey"])
	r.Sha256 = dcl.FlattenString(m["sha256"])

	return r
}

// expandInstanceGuestAcceleratorsMap expands the contents of InstanceGuestAccelerators into a JSON
// request object.
func expandInstanceGuestAcceleratorsMap(c *Client, f map[string]InstanceGuestAccelerators, res *Instance) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandInstanceGuestAccelerators(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandInstanceGuestAcceleratorsSlice expands the contents of InstanceGuestAccelerators into a JSON
// request object.
func expandInstanceGuestAcceleratorsSlice(c *Client, f []InstanceGuestAccelerators, res *Instance) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandInstanceGuestAccelerators(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenInstanceGuestAcceleratorsMap flattens the contents of InstanceGuestAccelerators from a JSON
// response object.
func flattenInstanceGuestAcceleratorsMap(c *Client, i interface{}, res *Instance) map[string]InstanceGuestAccelerators {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]InstanceGuestAccelerators{}
	}

	if len(a) == 0 {
		return map[string]InstanceGuestAccelerators{}
	}

	items := make(map[string]InstanceGuestAccelerators)
	for k, item := range a {
		items[k] = *flattenInstanceGuestAccelerators(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenInstanceGuestAcceleratorsSlice flattens the contents of InstanceGuestAccelerators from a JSON
// response object.
func flattenInstanceGuestAcceleratorsSlice(c *Client, i interface{}, res *Instance) []InstanceGuestAccelerators {
	a, ok := i.([]interface{})
	if !ok {
		return []InstanceGuestAccelerators{}
	}

	if len(a) == 0 {
		return []InstanceGuestAccelerators{}
	}

	items := make([]InstanceGuestAccelerators, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenInstanceGuestAccelerators(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandInstanceGuestAccelerators expands an instance of InstanceGuestAccelerators into a JSON
// request object.
func expandInstanceGuestAccelerators(c *Client, f *InstanceGuestAccelerators, res *Instance) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.AcceleratorCount; !dcl.IsEmptyValueIndirect(v) {
		m["acceleratorCount"] = v
	}
	if v := f.AcceleratorType; !dcl.IsEmptyValueIndirect(v) {
		m["acceleratorType"] = v
	}

	return m, nil
}

// flattenInstanceGuestAccelerators flattens an instance of InstanceGuestAccelerators from a JSON
// response object.
func flattenInstanceGuestAccelerators(c *Client, i interface{}, res *Instance) *InstanceGuestAccelerators {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &InstanceGuestAccelerators{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyInstanceGuestAccelerators
	}
	r.AcceleratorCount = dcl.FlattenInteger(m["acceleratorCount"])
	r.AcceleratorType = dcl.FlattenString(m["acceleratorType"])

	return r
}

// expandInstanceNetworkInterfacesMap expands the contents of InstanceNetworkInterfaces into a JSON
// request object.
func expandInstanceNetworkInterfacesMap(c *Client, f map[string]InstanceNetworkInterfaces, res *Instance) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandInstanceNetworkInterfaces(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandInstanceNetworkInterfacesSlice expands the contents of InstanceNetworkInterfaces into a JSON
// request object.
func expandInstanceNetworkInterfacesSlice(c *Client, f []InstanceNetworkInterfaces, res *Instance) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandInstanceNetworkInterfaces(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenInstanceNetworkInterfacesMap flattens the contents of InstanceNetworkInterfaces from a JSON
// response object.
func flattenInstanceNetworkInterfacesMap(c *Client, i interface{}, res *Instance) map[string]InstanceNetworkInterfaces {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]InstanceNetworkInterfaces{}
	}

	if len(a) == 0 {
		return map[string]InstanceNetworkInterfaces{}
	}

	items := make(map[string]InstanceNetworkInterfaces)
	for k, item := range a {
		items[k] = *flattenInstanceNetworkInterfaces(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenInstanceNetworkInterfacesSlice flattens the contents of InstanceNetworkInterfaces from a JSON
// response object.
func flattenInstanceNetworkInterfacesSlice(c *Client, i interface{}, res *Instance) []InstanceNetworkInterfaces {
	a, ok := i.([]interface{})
	if !ok {
		return []InstanceNetworkInterfaces{}
	}

	if len(a) == 0 {
		return []InstanceNetworkInterfaces{}
	}

	items := make([]InstanceNetworkInterfaces, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenInstanceNetworkInterfaces(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandInstanceNetworkInterfaces expands an instance of InstanceNetworkInterfaces into a JSON
// request object.
func expandInstanceNetworkInterfaces(c *Client, f *InstanceNetworkInterfaces, res *Instance) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v, err := expandInstanceNetworkInterfacesAccessConfigsSlice(c, f.AccessConfigs, res); err != nil {
		return nil, fmt.Errorf("error expanding AccessConfigs into accessConfigs: %w", err)
	} else if v != nil {
		m["accessConfigs"] = v
	}
	if v, err := expandInstanceNetworkInterfacesIPv6AccessConfigsSlice(c, f.IPv6AccessConfigs, res); err != nil {
		return nil, fmt.Errorf("error expanding IPv6AccessConfigs into ipv6AccessConfigs: %w", err)
	} else if v != nil {
		m["ipv6AccessConfigs"] = v
	}
	if v, err := expandInstanceNetworkInterfacesAliasIPRangesSlice(c, f.AliasIPRanges, res); err != nil {
		return nil, fmt.Errorf("error expanding AliasIPRanges into aliasIPRanges: %w", err)
	} else if v != nil {
		m["aliasIPRanges"] = v
	}
	if v := f.Network; !dcl.IsEmptyValueIndirect(v) {
		m["network"] = v
	}
	if v := f.NetworkIP; !dcl.IsEmptyValueIndirect(v) {
		m["networkIP"] = v
	}
	if v := f.Subnetwork; !dcl.IsEmptyValueIndirect(v) {
		m["subnetwork"] = v
	}

	return m, nil
}

// flattenInstanceNetworkInterfaces flattens an instance of InstanceNetworkInterfaces from a JSON
// response object.
func flattenInstanceNetworkInterfaces(c *Client, i interface{}, res *Instance) *InstanceNetworkInterfaces {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &InstanceNetworkInterfaces{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyInstanceNetworkInterfaces
	}
	r.AccessConfigs = flattenInstanceNetworkInterfacesAccessConfigsSlice(c, m["accessConfigs"], res)
	r.IPv6AccessConfigs = flattenInstanceNetworkInterfacesIPv6AccessConfigsSlice(c, m["ipv6AccessConfigs"], res)
	r.AliasIPRanges = flattenInstanceNetworkInterfacesAliasIPRangesSlice(c, m["aliasIPRanges"], res)
	r.Name = dcl.FlattenString(m["name"])
	r.Network = dcl.FlattenString(m["network"])
	r.NetworkIP = dcl.FlattenString(m["networkIP"])
	r.Subnetwork = dcl.FlattenString(m["subnetwork"])

	return r
}

// expandInstanceNetworkInterfacesAccessConfigsMap expands the contents of InstanceNetworkInterfacesAccessConfigs into a JSON
// request object.
func expandInstanceNetworkInterfacesAccessConfigsMap(c *Client, f map[string]InstanceNetworkInterfacesAccessConfigs, res *Instance) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandInstanceNetworkInterfacesAccessConfigs(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandInstanceNetworkInterfacesAccessConfigsSlice expands the contents of InstanceNetworkInterfacesAccessConfigs into a JSON
// request object.
func expandInstanceNetworkInterfacesAccessConfigsSlice(c *Client, f []InstanceNetworkInterfacesAccessConfigs, res *Instance) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandInstanceNetworkInterfacesAccessConfigs(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenInstanceNetworkInterfacesAccessConfigsMap flattens the contents of InstanceNetworkInterfacesAccessConfigs from a JSON
// response object.
func flattenInstanceNetworkInterfacesAccessConfigsMap(c *Client, i interface{}, res *Instance) map[string]InstanceNetworkInterfacesAccessConfigs {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]InstanceNetworkInterfacesAccessConfigs{}
	}

	if len(a) == 0 {
		return map[string]InstanceNetworkInterfacesAccessConfigs{}
	}

	items := make(map[string]InstanceNetworkInterfacesAccessConfigs)
	for k, item := range a {
		items[k] = *flattenInstanceNetworkInterfacesAccessConfigs(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenInstanceNetworkInterfacesAccessConfigsSlice flattens the contents of InstanceNetworkInterfacesAccessConfigs from a JSON
// response object.
func flattenInstanceNetworkInterfacesAccessConfigsSlice(c *Client, i interface{}, res *Instance) []InstanceNetworkInterfacesAccessConfigs {
	a, ok := i.([]interface{})
	if !ok {
		return []InstanceNetworkInterfacesAccessConfigs{}
	}

	if len(a) == 0 {
		return []InstanceNetworkInterfacesAccessConfigs{}
	}

	items := make([]InstanceNetworkInterfacesAccessConfigs, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenInstanceNetworkInterfacesAccessConfigs(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandInstanceNetworkInterfacesAccessConfigs expands an instance of InstanceNetworkInterfacesAccessConfigs into a JSON
// request object.
func expandInstanceNetworkInterfacesAccessConfigs(c *Client, f *InstanceNetworkInterfacesAccessConfigs, res *Instance) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.Name; !dcl.IsEmptyValueIndirect(v) {
		m["name"] = v
	}
	if v := f.NatIP; !dcl.IsEmptyValueIndirect(v) {
		m["natIP"] = v
	}
	if v := f.SetPublicPtr; !dcl.IsEmptyValueIndirect(v) {
		m["setPublicPtr"] = v
	}
	if v := f.PublicPtrDomainName; !dcl.IsEmptyValueIndirect(v) {
		m["publicPtrDomainName"] = v
	}
	if v := f.NetworkTier; !dcl.IsEmptyValueIndirect(v) {
		m["networkTier"] = v
	}
	if v := f.Type; !dcl.IsEmptyValueIndirect(v) {
		m["type"] = v
	}

	return m, nil
}

// flattenInstanceNetworkInterfacesAccessConfigs flattens an instance of InstanceNetworkInterfacesAccessConfigs from a JSON
// response object.
func flattenInstanceNetworkInterfacesAccessConfigs(c *Client, i interface{}, res *Instance) *InstanceNetworkInterfacesAccessConfigs {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &InstanceNetworkInterfacesAccessConfigs{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyInstanceNetworkInterfacesAccessConfigs
	}
	r.Name = dcl.FlattenString(m["name"])
	r.NatIP = dcl.FlattenString(m["natIP"])
	r.ExternalIPv6 = dcl.FlattenString(m["externalIPv6"])
	r.ExternalIPv6PrefixLength = dcl.FlattenString(m["externalIPv6PrefixLength"])
	r.SetPublicPtr = dcl.FlattenBool(m["setPublicPtr"])
	r.PublicPtrDomainName = dcl.FlattenString(m["publicPtrDomainName"])
	r.NetworkTier = flattenInstanceNetworkInterfacesAccessConfigsNetworkTierEnum(m["networkTier"])
	r.Type = flattenInstanceNetworkInterfacesAccessConfigsTypeEnum(m["type"])

	return r
}

// expandInstanceNetworkInterfacesIPv6AccessConfigsMap expands the contents of InstanceNetworkInterfacesIPv6AccessConfigs into a JSON
// request object.
func expandInstanceNetworkInterfacesIPv6AccessConfigsMap(c *Client, f map[string]InstanceNetworkInterfacesIPv6AccessConfigs, res *Instance) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandInstanceNetworkInterfacesIPv6AccessConfigs(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandInstanceNetworkInterfacesIPv6AccessConfigsSlice expands the contents of InstanceNetworkInterfacesIPv6AccessConfigs into a JSON
// request object.
func expandInstanceNetworkInterfacesIPv6AccessConfigsSlice(c *Client, f []InstanceNetworkInterfacesIPv6AccessConfigs, res *Instance) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandInstanceNetworkInterfacesIPv6AccessConfigs(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenInstanceNetworkInterfacesIPv6AccessConfigsMap flattens the contents of InstanceNetworkInterfacesIPv6AccessConfigs from a JSON
// response object.
func flattenInstanceNetworkInterfacesIPv6AccessConfigsMap(c *Client, i interface{}, res *Instance) map[string]InstanceNetworkInterfacesIPv6AccessConfigs {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]InstanceNetworkInterfacesIPv6AccessConfigs{}
	}

	if len(a) == 0 {
		return map[string]InstanceNetworkInterfacesIPv6AccessConfigs{}
	}

	items := make(map[string]InstanceNetworkInterfacesIPv6AccessConfigs)
	for k, item := range a {
		items[k] = *flattenInstanceNetworkInterfacesIPv6AccessConfigs(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenInstanceNetworkInterfacesIPv6AccessConfigsSlice flattens the contents of InstanceNetworkInterfacesIPv6AccessConfigs from a JSON
// response object.
func flattenInstanceNetworkInterfacesIPv6AccessConfigsSlice(c *Client, i interface{}, res *Instance) []InstanceNetworkInterfacesIPv6AccessConfigs {
	a, ok := i.([]interface{})
	if !ok {
		return []InstanceNetworkInterfacesIPv6AccessConfigs{}
	}

	if len(a) == 0 {
		return []InstanceNetworkInterfacesIPv6AccessConfigs{}
	}

	items := make([]InstanceNetworkInterfacesIPv6AccessConfigs, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenInstanceNetworkInterfacesIPv6AccessConfigs(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandInstanceNetworkInterfacesIPv6AccessConfigs expands an instance of InstanceNetworkInterfacesIPv6AccessConfigs into a JSON
// request object.
func expandInstanceNetworkInterfacesIPv6AccessConfigs(c *Client, f *InstanceNetworkInterfacesIPv6AccessConfigs, res *Instance) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.Name; !dcl.IsEmptyValueIndirect(v) {
		m["name"] = v
	}
	if v := f.NatIP; !dcl.IsEmptyValueIndirect(v) {
		m["natIP"] = v
	}
	if v := f.SetPublicPtr; !dcl.IsEmptyValueIndirect(v) {
		m["setPublicPtr"] = v
	}
	if v := f.PublicPtrDomainName; !dcl.IsEmptyValueIndirect(v) {
		m["publicPtrDomainName"] = v
	}
	if v := f.NetworkTier; !dcl.IsEmptyValueIndirect(v) {
		m["networkTier"] = v
	}
	if v := f.Type; !dcl.IsEmptyValueIndirect(v) {
		m["type"] = v
	}

	return m, nil
}

// flattenInstanceNetworkInterfacesIPv6AccessConfigs flattens an instance of InstanceNetworkInterfacesIPv6AccessConfigs from a JSON
// response object.
func flattenInstanceNetworkInterfacesIPv6AccessConfigs(c *Client, i interface{}, res *Instance) *InstanceNetworkInterfacesIPv6AccessConfigs {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &InstanceNetworkInterfacesIPv6AccessConfigs{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyInstanceNetworkInterfacesIPv6AccessConfigs
	}
	r.Name = dcl.FlattenString(m["name"])
	r.NatIP = dcl.FlattenString(m["natIP"])
	r.ExternalIPv6 = dcl.FlattenString(m["externalIPv6"])
	r.ExternalIPv6PrefixLength = dcl.FlattenString(m["externalIPv6PrefixLength"])
	r.SetPublicPtr = dcl.FlattenBool(m["setPublicPtr"])
	r.PublicPtrDomainName = dcl.FlattenString(m["publicPtrDomainName"])
	r.NetworkTier = flattenInstanceNetworkInterfacesIPv6AccessConfigsNetworkTierEnum(m["networkTier"])
	r.Type = flattenInstanceNetworkInterfacesIPv6AccessConfigsTypeEnum(m["type"])

	return r
}

// expandInstanceNetworkInterfacesAliasIPRangesMap expands the contents of InstanceNetworkInterfacesAliasIPRanges into a JSON
// request object.
func expandInstanceNetworkInterfacesAliasIPRangesMap(c *Client, f map[string]InstanceNetworkInterfacesAliasIPRanges, res *Instance) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandInstanceNetworkInterfacesAliasIPRanges(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandInstanceNetworkInterfacesAliasIPRangesSlice expands the contents of InstanceNetworkInterfacesAliasIPRanges into a JSON
// request object.
func expandInstanceNetworkInterfacesAliasIPRangesSlice(c *Client, f []InstanceNetworkInterfacesAliasIPRanges, res *Instance) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandInstanceNetworkInterfacesAliasIPRanges(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenInstanceNetworkInterfacesAliasIPRangesMap flattens the contents of InstanceNetworkInterfacesAliasIPRanges from a JSON
// response object.
func flattenInstanceNetworkInterfacesAliasIPRangesMap(c *Client, i interface{}, res *Instance) map[string]InstanceNetworkInterfacesAliasIPRanges {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]InstanceNetworkInterfacesAliasIPRanges{}
	}

	if len(a) == 0 {
		return map[string]InstanceNetworkInterfacesAliasIPRanges{}
	}

	items := make(map[string]InstanceNetworkInterfacesAliasIPRanges)
	for k, item := range a {
		items[k] = *flattenInstanceNetworkInterfacesAliasIPRanges(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenInstanceNetworkInterfacesAliasIPRangesSlice flattens the contents of InstanceNetworkInterfacesAliasIPRanges from a JSON
// response object.
func flattenInstanceNetworkInterfacesAliasIPRangesSlice(c *Client, i interface{}, res *Instance) []InstanceNetworkInterfacesAliasIPRanges {
	a, ok := i.([]interface{})
	if !ok {
		return []InstanceNetworkInterfacesAliasIPRanges{}
	}

	if len(a) == 0 {
		return []InstanceNetworkInterfacesAliasIPRanges{}
	}

	items := make([]InstanceNetworkInterfacesAliasIPRanges, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenInstanceNetworkInterfacesAliasIPRanges(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandInstanceNetworkInterfacesAliasIPRanges expands an instance of InstanceNetworkInterfacesAliasIPRanges into a JSON
// request object.
func expandInstanceNetworkInterfacesAliasIPRanges(c *Client, f *InstanceNetworkInterfacesAliasIPRanges, res *Instance) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.IPCidrRange; !dcl.IsEmptyValueIndirect(v) {
		m["ipCidrRange"] = v
	}
	if v := f.SubnetworkRangeName; !dcl.IsEmptyValueIndirect(v) {
		m["subnetworkRangeName"] = v
	}

	return m, nil
}

// flattenInstanceNetworkInterfacesAliasIPRanges flattens an instance of InstanceNetworkInterfacesAliasIPRanges from a JSON
// response object.
func flattenInstanceNetworkInterfacesAliasIPRanges(c *Client, i interface{}, res *Instance) *InstanceNetworkInterfacesAliasIPRanges {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &InstanceNetworkInterfacesAliasIPRanges{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyInstanceNetworkInterfacesAliasIPRanges
	}
	r.IPCidrRange = dcl.FlattenString(m["ipCidrRange"])
	r.SubnetworkRangeName = dcl.FlattenString(m["subnetworkRangeName"])

	return r
}

// expandInstanceSchedulingMap expands the contents of InstanceScheduling into a JSON
// request object.
func expandInstanceSchedulingMap(c *Client, f map[string]InstanceScheduling, res *Instance) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandInstanceScheduling(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandInstanceSchedulingSlice expands the contents of InstanceScheduling into a JSON
// request object.
func expandInstanceSchedulingSlice(c *Client, f []InstanceScheduling, res *Instance) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandInstanceScheduling(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenInstanceSchedulingMap flattens the contents of InstanceScheduling from a JSON
// response object.
func flattenInstanceSchedulingMap(c *Client, i interface{}, res *Instance) map[string]InstanceScheduling {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]InstanceScheduling{}
	}

	if len(a) == 0 {
		return map[string]InstanceScheduling{}
	}

	items := make(map[string]InstanceScheduling)
	for k, item := range a {
		items[k] = *flattenInstanceScheduling(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenInstanceSchedulingSlice flattens the contents of InstanceScheduling from a JSON
// response object.
func flattenInstanceSchedulingSlice(c *Client, i interface{}, res *Instance) []InstanceScheduling {
	a, ok := i.([]interface{})
	if !ok {
		return []InstanceScheduling{}
	}

	if len(a) == 0 {
		return []InstanceScheduling{}
	}

	items := make([]InstanceScheduling, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenInstanceScheduling(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandInstanceScheduling expands an instance of InstanceScheduling into a JSON
// request object.
func expandInstanceScheduling(c *Client, f *InstanceScheduling, res *Instance) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.AutomaticRestart; !dcl.IsEmptyValueIndirect(v) {
		m["automaticRestart"] = v
	}
	if v := f.OnHostMaintenance; !dcl.IsEmptyValueIndirect(v) {
		m["onHostMaintenance"] = v
	}
	if v := f.Preemptible; !dcl.IsEmptyValueIndirect(v) {
		m["preemptible"] = v
	}

	return m, nil
}

// flattenInstanceScheduling flattens an instance of InstanceScheduling from a JSON
// response object.
func flattenInstanceScheduling(c *Client, i interface{}, res *Instance) *InstanceScheduling {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &InstanceScheduling{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyInstanceScheduling
	}
	r.AutomaticRestart = dcl.FlattenBool(m["automaticRestart"])
	r.OnHostMaintenance = dcl.FlattenString(m["onHostMaintenance"])
	r.Preemptible = dcl.FlattenBool(m["preemptible"])

	return r
}

// expandInstanceServiceAccountsMap expands the contents of InstanceServiceAccounts into a JSON
// request object.
func expandInstanceServiceAccountsMap(c *Client, f map[string]InstanceServiceAccounts, res *Instance) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandInstanceServiceAccounts(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandInstanceServiceAccountsSlice expands the contents of InstanceServiceAccounts into a JSON
// request object.
func expandInstanceServiceAccountsSlice(c *Client, f []InstanceServiceAccounts, res *Instance) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandInstanceServiceAccounts(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenInstanceServiceAccountsMap flattens the contents of InstanceServiceAccounts from a JSON
// response object.
func flattenInstanceServiceAccountsMap(c *Client, i interface{}, res *Instance) map[string]InstanceServiceAccounts {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]InstanceServiceAccounts{}
	}

	if len(a) == 0 {
		return map[string]InstanceServiceAccounts{}
	}

	items := make(map[string]InstanceServiceAccounts)
	for k, item := range a {
		items[k] = *flattenInstanceServiceAccounts(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenInstanceServiceAccountsSlice flattens the contents of InstanceServiceAccounts from a JSON
// response object.
func flattenInstanceServiceAccountsSlice(c *Client, i interface{}, res *Instance) []InstanceServiceAccounts {
	a, ok := i.([]interface{})
	if !ok {
		return []InstanceServiceAccounts{}
	}

	if len(a) == 0 {
		return []InstanceServiceAccounts{}
	}

	items := make([]InstanceServiceAccounts, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenInstanceServiceAccounts(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandInstanceServiceAccounts expands an instance of InstanceServiceAccounts into a JSON
// request object.
func expandInstanceServiceAccounts(c *Client, f *InstanceServiceAccounts, res *Instance) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.Email; !dcl.IsEmptyValueIndirect(v) {
		m["email"] = v
	}
	if v := f.Scopes; v != nil {
		m["scopes"] = v
	}

	return m, nil
}

// flattenInstanceServiceAccounts flattens an instance of InstanceServiceAccounts from a JSON
// response object.
func flattenInstanceServiceAccounts(c *Client, i interface{}, res *Instance) *InstanceServiceAccounts {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &InstanceServiceAccounts{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyInstanceServiceAccounts
	}
	r.Email = dcl.FlattenString(m["email"])
	r.Scopes = dcl.FlattenStringSlice(m["scopes"])

	return r
}

// expandInstanceShieldedInstanceConfigMap expands the contents of InstanceShieldedInstanceConfig into a JSON
// request object.
func expandInstanceShieldedInstanceConfigMap(c *Client, f map[string]InstanceShieldedInstanceConfig, res *Instance) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandInstanceShieldedInstanceConfig(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandInstanceShieldedInstanceConfigSlice expands the contents of InstanceShieldedInstanceConfig into a JSON
// request object.
func expandInstanceShieldedInstanceConfigSlice(c *Client, f []InstanceShieldedInstanceConfig, res *Instance) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandInstanceShieldedInstanceConfig(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenInstanceShieldedInstanceConfigMap flattens the contents of InstanceShieldedInstanceConfig from a JSON
// response object.
func flattenInstanceShieldedInstanceConfigMap(c *Client, i interface{}, res *Instance) map[string]InstanceShieldedInstanceConfig {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]InstanceShieldedInstanceConfig{}
	}

	if len(a) == 0 {
		return map[string]InstanceShieldedInstanceConfig{}
	}

	items := make(map[string]InstanceShieldedInstanceConfig)
	for k, item := range a {
		items[k] = *flattenInstanceShieldedInstanceConfig(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenInstanceShieldedInstanceConfigSlice flattens the contents of InstanceShieldedInstanceConfig from a JSON
// response object.
func flattenInstanceShieldedInstanceConfigSlice(c *Client, i interface{}, res *Instance) []InstanceShieldedInstanceConfig {
	a, ok := i.([]interface{})
	if !ok {
		return []InstanceShieldedInstanceConfig{}
	}

	if len(a) == 0 {
		return []InstanceShieldedInstanceConfig{}
	}

	items := make([]InstanceShieldedInstanceConfig, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenInstanceShieldedInstanceConfig(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandInstanceShieldedInstanceConfig expands an instance of InstanceShieldedInstanceConfig into a JSON
// request object.
func expandInstanceShieldedInstanceConfig(c *Client, f *InstanceShieldedInstanceConfig, res *Instance) (map[string]interface{}, error) {
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

// flattenInstanceShieldedInstanceConfig flattens an instance of InstanceShieldedInstanceConfig from a JSON
// response object.
func flattenInstanceShieldedInstanceConfig(c *Client, i interface{}, res *Instance) *InstanceShieldedInstanceConfig {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &InstanceShieldedInstanceConfig{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyInstanceShieldedInstanceConfig
	}
	r.EnableSecureBoot = dcl.FlattenBool(m["enableSecureBoot"])
	r.EnableVtpm = dcl.FlattenBool(m["enableVtpm"])
	r.EnableIntegrityMonitoring = dcl.FlattenBool(m["enableIntegrityMonitoring"])

	return r
}

// flattenInstanceDisksInterfaceEnumMap flattens the contents of InstanceDisksInterfaceEnum from a JSON
// response object.
func flattenInstanceDisksInterfaceEnumMap(c *Client, i interface{}, res *Instance) map[string]InstanceDisksInterfaceEnum {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]InstanceDisksInterfaceEnum{}
	}

	if len(a) == 0 {
		return map[string]InstanceDisksInterfaceEnum{}
	}

	items := make(map[string]InstanceDisksInterfaceEnum)
	for k, item := range a {
		items[k] = *flattenInstanceDisksInterfaceEnum(item.(interface{}))
	}

	return items
}

// flattenInstanceDisksInterfaceEnumSlice flattens the contents of InstanceDisksInterfaceEnum from a JSON
// response object.
func flattenInstanceDisksInterfaceEnumSlice(c *Client, i interface{}, res *Instance) []InstanceDisksInterfaceEnum {
	a, ok := i.([]interface{})
	if !ok {
		return []InstanceDisksInterfaceEnum{}
	}

	if len(a) == 0 {
		return []InstanceDisksInterfaceEnum{}
	}

	items := make([]InstanceDisksInterfaceEnum, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenInstanceDisksInterfaceEnum(item.(interface{})))
	}

	return items
}

// flattenInstanceDisksInterfaceEnum asserts that an interface is a string, and returns a
// pointer to a *InstanceDisksInterfaceEnum with the same value as that string.
func flattenInstanceDisksInterfaceEnum(i interface{}) *InstanceDisksInterfaceEnum {
	s, ok := i.(string)
	if !ok {
		return nil
	}

	return InstanceDisksInterfaceEnumRef(s)
}

// flattenInstanceDisksModeEnumMap flattens the contents of InstanceDisksModeEnum from a JSON
// response object.
func flattenInstanceDisksModeEnumMap(c *Client, i interface{}, res *Instance) map[string]InstanceDisksModeEnum {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]InstanceDisksModeEnum{}
	}

	if len(a) == 0 {
		return map[string]InstanceDisksModeEnum{}
	}

	items := make(map[string]InstanceDisksModeEnum)
	for k, item := range a {
		items[k] = *flattenInstanceDisksModeEnum(item.(interface{}))
	}

	return items
}

// flattenInstanceDisksModeEnumSlice flattens the contents of InstanceDisksModeEnum from a JSON
// response object.
func flattenInstanceDisksModeEnumSlice(c *Client, i interface{}, res *Instance) []InstanceDisksModeEnum {
	a, ok := i.([]interface{})
	if !ok {
		return []InstanceDisksModeEnum{}
	}

	if len(a) == 0 {
		return []InstanceDisksModeEnum{}
	}

	items := make([]InstanceDisksModeEnum, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenInstanceDisksModeEnum(item.(interface{})))
	}

	return items
}

// flattenInstanceDisksModeEnum asserts that an interface is a string, and returns a
// pointer to a *InstanceDisksModeEnum with the same value as that string.
func flattenInstanceDisksModeEnum(i interface{}) *InstanceDisksModeEnum {
	s, ok := i.(string)
	if !ok {
		return nil
	}

	return InstanceDisksModeEnumRef(s)
}

// flattenInstanceDisksTypeEnumMap flattens the contents of InstanceDisksTypeEnum from a JSON
// response object.
func flattenInstanceDisksTypeEnumMap(c *Client, i interface{}, res *Instance) map[string]InstanceDisksTypeEnum {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]InstanceDisksTypeEnum{}
	}

	if len(a) == 0 {
		return map[string]InstanceDisksTypeEnum{}
	}

	items := make(map[string]InstanceDisksTypeEnum)
	for k, item := range a {
		items[k] = *flattenInstanceDisksTypeEnum(item.(interface{}))
	}

	return items
}

// flattenInstanceDisksTypeEnumSlice flattens the contents of InstanceDisksTypeEnum from a JSON
// response object.
func flattenInstanceDisksTypeEnumSlice(c *Client, i interface{}, res *Instance) []InstanceDisksTypeEnum {
	a, ok := i.([]interface{})
	if !ok {
		return []InstanceDisksTypeEnum{}
	}

	if len(a) == 0 {
		return []InstanceDisksTypeEnum{}
	}

	items := make([]InstanceDisksTypeEnum, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenInstanceDisksTypeEnum(item.(interface{})))
	}

	return items
}

// flattenInstanceDisksTypeEnum asserts that an interface is a string, and returns a
// pointer to a *InstanceDisksTypeEnum with the same value as that string.
func flattenInstanceDisksTypeEnum(i interface{}) *InstanceDisksTypeEnum {
	s, ok := i.(string)
	if !ok {
		return nil
	}

	return InstanceDisksTypeEnumRef(s)
}

// flattenInstanceNetworkInterfacesAccessConfigsNetworkTierEnumMap flattens the contents of InstanceNetworkInterfacesAccessConfigsNetworkTierEnum from a JSON
// response object.
func flattenInstanceNetworkInterfacesAccessConfigsNetworkTierEnumMap(c *Client, i interface{}, res *Instance) map[string]InstanceNetworkInterfacesAccessConfigsNetworkTierEnum {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]InstanceNetworkInterfacesAccessConfigsNetworkTierEnum{}
	}

	if len(a) == 0 {
		return map[string]InstanceNetworkInterfacesAccessConfigsNetworkTierEnum{}
	}

	items := make(map[string]InstanceNetworkInterfacesAccessConfigsNetworkTierEnum)
	for k, item := range a {
		items[k] = *flattenInstanceNetworkInterfacesAccessConfigsNetworkTierEnum(item.(interface{}))
	}

	return items
}

// flattenInstanceNetworkInterfacesAccessConfigsNetworkTierEnumSlice flattens the contents of InstanceNetworkInterfacesAccessConfigsNetworkTierEnum from a JSON
// response object.
func flattenInstanceNetworkInterfacesAccessConfigsNetworkTierEnumSlice(c *Client, i interface{}, res *Instance) []InstanceNetworkInterfacesAccessConfigsNetworkTierEnum {
	a, ok := i.([]interface{})
	if !ok {
		return []InstanceNetworkInterfacesAccessConfigsNetworkTierEnum{}
	}

	if len(a) == 0 {
		return []InstanceNetworkInterfacesAccessConfigsNetworkTierEnum{}
	}

	items := make([]InstanceNetworkInterfacesAccessConfigsNetworkTierEnum, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenInstanceNetworkInterfacesAccessConfigsNetworkTierEnum(item.(interface{})))
	}

	return items
}

// flattenInstanceNetworkInterfacesAccessConfigsNetworkTierEnum asserts that an interface is a string, and returns a
// pointer to a *InstanceNetworkInterfacesAccessConfigsNetworkTierEnum with the same value as that string.
func flattenInstanceNetworkInterfacesAccessConfigsNetworkTierEnum(i interface{}) *InstanceNetworkInterfacesAccessConfigsNetworkTierEnum {
	s, ok := i.(string)
	if !ok {
		return nil
	}

	return InstanceNetworkInterfacesAccessConfigsNetworkTierEnumRef(s)
}

// flattenInstanceNetworkInterfacesAccessConfigsTypeEnumMap flattens the contents of InstanceNetworkInterfacesAccessConfigsTypeEnum from a JSON
// response object.
func flattenInstanceNetworkInterfacesAccessConfigsTypeEnumMap(c *Client, i interface{}, res *Instance) map[string]InstanceNetworkInterfacesAccessConfigsTypeEnum {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]InstanceNetworkInterfacesAccessConfigsTypeEnum{}
	}

	if len(a) == 0 {
		return map[string]InstanceNetworkInterfacesAccessConfigsTypeEnum{}
	}

	items := make(map[string]InstanceNetworkInterfacesAccessConfigsTypeEnum)
	for k, item := range a {
		items[k] = *flattenInstanceNetworkInterfacesAccessConfigsTypeEnum(item.(interface{}))
	}

	return items
}

// flattenInstanceNetworkInterfacesAccessConfigsTypeEnumSlice flattens the contents of InstanceNetworkInterfacesAccessConfigsTypeEnum from a JSON
// response object.
func flattenInstanceNetworkInterfacesAccessConfigsTypeEnumSlice(c *Client, i interface{}, res *Instance) []InstanceNetworkInterfacesAccessConfigsTypeEnum {
	a, ok := i.([]interface{})
	if !ok {
		return []InstanceNetworkInterfacesAccessConfigsTypeEnum{}
	}

	if len(a) == 0 {
		return []InstanceNetworkInterfacesAccessConfigsTypeEnum{}
	}

	items := make([]InstanceNetworkInterfacesAccessConfigsTypeEnum, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenInstanceNetworkInterfacesAccessConfigsTypeEnum(item.(interface{})))
	}

	return items
}

// flattenInstanceNetworkInterfacesAccessConfigsTypeEnum asserts that an interface is a string, and returns a
// pointer to a *InstanceNetworkInterfacesAccessConfigsTypeEnum with the same value as that string.
func flattenInstanceNetworkInterfacesAccessConfigsTypeEnum(i interface{}) *InstanceNetworkInterfacesAccessConfigsTypeEnum {
	s, ok := i.(string)
	if !ok {
		return nil
	}

	return InstanceNetworkInterfacesAccessConfigsTypeEnumRef(s)
}

// flattenInstanceNetworkInterfacesIPv6AccessConfigsNetworkTierEnumMap flattens the contents of InstanceNetworkInterfacesIPv6AccessConfigsNetworkTierEnum from a JSON
// response object.
func flattenInstanceNetworkInterfacesIPv6AccessConfigsNetworkTierEnumMap(c *Client, i interface{}, res *Instance) map[string]InstanceNetworkInterfacesIPv6AccessConfigsNetworkTierEnum {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]InstanceNetworkInterfacesIPv6AccessConfigsNetworkTierEnum{}
	}

	if len(a) == 0 {
		return map[string]InstanceNetworkInterfacesIPv6AccessConfigsNetworkTierEnum{}
	}

	items := make(map[string]InstanceNetworkInterfacesIPv6AccessConfigsNetworkTierEnum)
	for k, item := range a {
		items[k] = *flattenInstanceNetworkInterfacesIPv6AccessConfigsNetworkTierEnum(item.(interface{}))
	}

	return items
}

// flattenInstanceNetworkInterfacesIPv6AccessConfigsNetworkTierEnumSlice flattens the contents of InstanceNetworkInterfacesIPv6AccessConfigsNetworkTierEnum from a JSON
// response object.
func flattenInstanceNetworkInterfacesIPv6AccessConfigsNetworkTierEnumSlice(c *Client, i interface{}, res *Instance) []InstanceNetworkInterfacesIPv6AccessConfigsNetworkTierEnum {
	a, ok := i.([]interface{})
	if !ok {
		return []InstanceNetworkInterfacesIPv6AccessConfigsNetworkTierEnum{}
	}

	if len(a) == 0 {
		return []InstanceNetworkInterfacesIPv6AccessConfigsNetworkTierEnum{}
	}

	items := make([]InstanceNetworkInterfacesIPv6AccessConfigsNetworkTierEnum, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenInstanceNetworkInterfacesIPv6AccessConfigsNetworkTierEnum(item.(interface{})))
	}

	return items
}

// flattenInstanceNetworkInterfacesIPv6AccessConfigsNetworkTierEnum asserts that an interface is a string, and returns a
// pointer to a *InstanceNetworkInterfacesIPv6AccessConfigsNetworkTierEnum with the same value as that string.
func flattenInstanceNetworkInterfacesIPv6AccessConfigsNetworkTierEnum(i interface{}) *InstanceNetworkInterfacesIPv6AccessConfigsNetworkTierEnum {
	s, ok := i.(string)
	if !ok {
		return nil
	}

	return InstanceNetworkInterfacesIPv6AccessConfigsNetworkTierEnumRef(s)
}

// flattenInstanceNetworkInterfacesIPv6AccessConfigsTypeEnumMap flattens the contents of InstanceNetworkInterfacesIPv6AccessConfigsTypeEnum from a JSON
// response object.
func flattenInstanceNetworkInterfacesIPv6AccessConfigsTypeEnumMap(c *Client, i interface{}, res *Instance) map[string]InstanceNetworkInterfacesIPv6AccessConfigsTypeEnum {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]InstanceNetworkInterfacesIPv6AccessConfigsTypeEnum{}
	}

	if len(a) == 0 {
		return map[string]InstanceNetworkInterfacesIPv6AccessConfigsTypeEnum{}
	}

	items := make(map[string]InstanceNetworkInterfacesIPv6AccessConfigsTypeEnum)
	for k, item := range a {
		items[k] = *flattenInstanceNetworkInterfacesIPv6AccessConfigsTypeEnum(item.(interface{}))
	}

	return items
}

// flattenInstanceNetworkInterfacesIPv6AccessConfigsTypeEnumSlice flattens the contents of InstanceNetworkInterfacesIPv6AccessConfigsTypeEnum from a JSON
// response object.
func flattenInstanceNetworkInterfacesIPv6AccessConfigsTypeEnumSlice(c *Client, i interface{}, res *Instance) []InstanceNetworkInterfacesIPv6AccessConfigsTypeEnum {
	a, ok := i.([]interface{})
	if !ok {
		return []InstanceNetworkInterfacesIPv6AccessConfigsTypeEnum{}
	}

	if len(a) == 0 {
		return []InstanceNetworkInterfacesIPv6AccessConfigsTypeEnum{}
	}

	items := make([]InstanceNetworkInterfacesIPv6AccessConfigsTypeEnum, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenInstanceNetworkInterfacesIPv6AccessConfigsTypeEnum(item.(interface{})))
	}

	return items
}

// flattenInstanceNetworkInterfacesIPv6AccessConfigsTypeEnum asserts that an interface is a string, and returns a
// pointer to a *InstanceNetworkInterfacesIPv6AccessConfigsTypeEnum with the same value as that string.
func flattenInstanceNetworkInterfacesIPv6AccessConfigsTypeEnum(i interface{}) *InstanceNetworkInterfacesIPv6AccessConfigsTypeEnum {
	s, ok := i.(string)
	if !ok {
		return nil
	}

	return InstanceNetworkInterfacesIPv6AccessConfigsTypeEnumRef(s)
}

// flattenInstanceStatusEnumMap flattens the contents of InstanceStatusEnum from a JSON
// response object.
func flattenInstanceStatusEnumMap(c *Client, i interface{}, res *Instance) map[string]InstanceStatusEnum {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]InstanceStatusEnum{}
	}

	if len(a) == 0 {
		return map[string]InstanceStatusEnum{}
	}

	items := make(map[string]InstanceStatusEnum)
	for k, item := range a {
		items[k] = *flattenInstanceStatusEnum(item.(interface{}))
	}

	return items
}

// flattenInstanceStatusEnumSlice flattens the contents of InstanceStatusEnum from a JSON
// response object.
func flattenInstanceStatusEnumSlice(c *Client, i interface{}, res *Instance) []InstanceStatusEnum {
	a, ok := i.([]interface{})
	if !ok {
		return []InstanceStatusEnum{}
	}

	if len(a) == 0 {
		return []InstanceStatusEnum{}
	}

	items := make([]InstanceStatusEnum, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenInstanceStatusEnum(item.(interface{})))
	}

	return items
}

// flattenInstanceStatusEnum asserts that an interface is a string, and returns a
// pointer to a *InstanceStatusEnum with the same value as that string.
func flattenInstanceStatusEnum(i interface{}) *InstanceStatusEnum {
	s, ok := i.(string)
	if !ok {
		return nil
	}

	return InstanceStatusEnumRef(s)
}

// This function returns a matcher that checks whether a serialized resource matches this resource
// in its parameters (as defined by the fields in a Get, which definitionally define resource
// identity).  This is useful in extracting the element from a List call.
func (r *Instance) matcher(c *Client) func([]byte) bool {
	return func(b []byte) bool {
		cr, err := unmarshalInstance(b, c, r)
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
		if nr.Zone == nil && ncr.Zone == nil {
			c.Config.Logger.Info("Both Zone fields null - considering equal.")
		} else if nr.Zone == nil || ncr.Zone == nil {
			c.Config.Logger.Info("Only one Zone field is null - considering unequal.")
			return false
		} else if *nr.Zone != *ncr.Zone {
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

type instanceDiff struct {
	// The diff should include one or the other of RequiresRecreate or UpdateOp.
	RequiresRecreate bool
	UpdateOp         instanceApiOperation
	FieldName        string // used for error logging
}

func convertFieldDiffsToInstanceDiffs(config *dcl.Config, fds []*dcl.FieldDiff, opts []dcl.ApplyOption) ([]instanceDiff, error) {
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
	var diffs []instanceDiff
	// For each operation name, create a instanceDiff which contains the operation.
	for opName, fieldDiffs := range opNamesToFieldDiffs {
		// Use the first field diff's field name for logging required recreate error.
		diff := instanceDiff{FieldName: fieldDiffs[0].FieldName}
		if opName == "Recreate" {
			diff.RequiresRecreate = true
		} else {
			apiOp, err := convertOpNameToInstanceApiOperation(opName, fieldDiffs, opts...)
			if err != nil {
				return diffs, err
			}
			diff.UpdateOp = apiOp
		}
		diffs = append(diffs, diff)
	}
	return diffs, nil
}

func convertOpNameToInstanceApiOperation(opName string, fieldDiffs []*dcl.FieldDiff, opts ...dcl.ApplyOption) (instanceApiOperation, error) {
	switch opName {

	case "updateInstanceSetDeletionProtectionOperation":
		return &updateInstanceSetDeletionProtectionOperation{FieldDiffs: fieldDiffs}, nil

	case "updateInstanceSetLabelsOperation":
		return &updateInstanceSetLabelsOperation{FieldDiffs: fieldDiffs}, nil

	case "updateInstanceSetMachineTypeOperation":
		return &updateInstanceSetMachineTypeOperation{FieldDiffs: fieldDiffs}, nil

	case "updateInstanceSetMetadataOperation":
		return &updateInstanceSetMetadataOperation{FieldDiffs: fieldDiffs}, nil

	case "updateInstanceSetTagsOperation":
		return &updateInstanceSetTagsOperation{FieldDiffs: fieldDiffs}, nil

	case "updateInstanceStartOperation":
		return &updateInstanceStartOperation{FieldDiffs: fieldDiffs}, nil

	case "updateInstanceStopOperation":
		return &updateInstanceStopOperation{FieldDiffs: fieldDiffs}, nil

	case "updateInstanceUpdateOperation":
		return &updateInstanceUpdateOperation{FieldDiffs: fieldDiffs}, nil

	case "updateInstanceUpdateShieldedInstanceConfigOperation":
		return &updateInstanceUpdateShieldedInstanceConfigOperation{FieldDiffs: fieldDiffs}, nil

	default:
		return nil, fmt.Errorf("no such operation with name: %v", opName)
	}
}

func extractInstanceFields(r *Instance) error {
	vScheduling := r.Scheduling
	if vScheduling == nil {
		// note: explicitly not the empty object.
		vScheduling = &InstanceScheduling{}
	}
	if err := extractInstanceSchedulingFields(r, vScheduling); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vScheduling) {
		r.Scheduling = vScheduling
	}
	vShieldedInstanceConfig := r.ShieldedInstanceConfig
	if vShieldedInstanceConfig == nil {
		// note: explicitly not the empty object.
		vShieldedInstanceConfig = &InstanceShieldedInstanceConfig{}
	}
	if err := extractInstanceShieldedInstanceConfigFields(r, vShieldedInstanceConfig); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vShieldedInstanceConfig) {
		r.ShieldedInstanceConfig = vShieldedInstanceConfig
	}
	return nil
}
func extractInstanceDisksFields(r *Instance, o *InstanceDisks) error {
	vDiskEncryptionKey := o.DiskEncryptionKey
	if vDiskEncryptionKey == nil {
		// note: explicitly not the empty object.
		vDiskEncryptionKey = &InstanceDisksDiskEncryptionKey{}
	}
	if err := extractInstanceDisksDiskEncryptionKeyFields(r, vDiskEncryptionKey); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vDiskEncryptionKey) {
		o.DiskEncryptionKey = vDiskEncryptionKey
	}
	vInitializeParams := o.InitializeParams
	if vInitializeParams == nil {
		// note: explicitly not the empty object.
		vInitializeParams = &InstanceDisksInitializeParams{}
	}
	if err := extractInstanceDisksInitializeParamsFields(r, vInitializeParams); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vInitializeParams) {
		o.InitializeParams = vInitializeParams
	}
	return nil
}
func extractInstanceDisksDiskEncryptionKeyFields(r *Instance, o *InstanceDisksDiskEncryptionKey) error {
	return nil
}
func extractInstanceDisksInitializeParamsFields(r *Instance, o *InstanceDisksInitializeParams) error {
	vSourceImageEncryptionKey := o.SourceImageEncryptionKey
	if vSourceImageEncryptionKey == nil {
		// note: explicitly not the empty object.
		vSourceImageEncryptionKey = &InstanceDisksInitializeParamsSourceImageEncryptionKey{}
	}
	if err := extractInstanceDisksInitializeParamsSourceImageEncryptionKeyFields(r, vSourceImageEncryptionKey); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vSourceImageEncryptionKey) {
		o.SourceImageEncryptionKey = vSourceImageEncryptionKey
	}
	return nil
}
func extractInstanceDisksInitializeParamsSourceImageEncryptionKeyFields(r *Instance, o *InstanceDisksInitializeParamsSourceImageEncryptionKey) error {
	return nil
}
func extractInstanceGuestAcceleratorsFields(r *Instance, o *InstanceGuestAccelerators) error {
	return nil
}
func extractInstanceNetworkInterfacesFields(r *Instance, o *InstanceNetworkInterfaces) error {
	return nil
}
func extractInstanceNetworkInterfacesAccessConfigsFields(r *Instance, o *InstanceNetworkInterfacesAccessConfigs) error {
	return nil
}
func extractInstanceNetworkInterfacesIPv6AccessConfigsFields(r *Instance, o *InstanceNetworkInterfacesIPv6AccessConfigs) error {
	return nil
}
func extractInstanceNetworkInterfacesAliasIPRangesFields(r *Instance, o *InstanceNetworkInterfacesAliasIPRanges) error {
	return nil
}
func extractInstanceSchedulingFields(r *Instance, o *InstanceScheduling) error {
	return nil
}
func extractInstanceServiceAccountsFields(r *Instance, o *InstanceServiceAccounts) error {
	return nil
}
func extractInstanceShieldedInstanceConfigFields(r *Instance, o *InstanceShieldedInstanceConfig) error {
	return nil
}

func postReadExtractInstanceFields(r *Instance) error {
	vScheduling := r.Scheduling
	if vScheduling == nil {
		// note: explicitly not the empty object.
		vScheduling = &InstanceScheduling{}
	}
	if err := postReadExtractInstanceSchedulingFields(r, vScheduling); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vScheduling) {
		r.Scheduling = vScheduling
	}
	vShieldedInstanceConfig := r.ShieldedInstanceConfig
	if vShieldedInstanceConfig == nil {
		// note: explicitly not the empty object.
		vShieldedInstanceConfig = &InstanceShieldedInstanceConfig{}
	}
	if err := postReadExtractInstanceShieldedInstanceConfigFields(r, vShieldedInstanceConfig); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vShieldedInstanceConfig) {
		r.ShieldedInstanceConfig = vShieldedInstanceConfig
	}
	return nil
}
func postReadExtractInstanceDisksFields(r *Instance, o *InstanceDisks) error {
	vDiskEncryptionKey := o.DiskEncryptionKey
	if vDiskEncryptionKey == nil {
		// note: explicitly not the empty object.
		vDiskEncryptionKey = &InstanceDisksDiskEncryptionKey{}
	}
	if err := extractInstanceDisksDiskEncryptionKeyFields(r, vDiskEncryptionKey); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vDiskEncryptionKey) {
		o.DiskEncryptionKey = vDiskEncryptionKey
	}
	vInitializeParams := o.InitializeParams
	if vInitializeParams == nil {
		// note: explicitly not the empty object.
		vInitializeParams = &InstanceDisksInitializeParams{}
	}
	if err := extractInstanceDisksInitializeParamsFields(r, vInitializeParams); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vInitializeParams) {
		o.InitializeParams = vInitializeParams
	}
	return nil
}
func postReadExtractInstanceDisksDiskEncryptionKeyFields(r *Instance, o *InstanceDisksDiskEncryptionKey) error {
	return nil
}
func postReadExtractInstanceDisksInitializeParamsFields(r *Instance, o *InstanceDisksInitializeParams) error {
	vSourceImageEncryptionKey := o.SourceImageEncryptionKey
	if vSourceImageEncryptionKey == nil {
		// note: explicitly not the empty object.
		vSourceImageEncryptionKey = &InstanceDisksInitializeParamsSourceImageEncryptionKey{}
	}
	if err := extractInstanceDisksInitializeParamsSourceImageEncryptionKeyFields(r, vSourceImageEncryptionKey); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vSourceImageEncryptionKey) {
		o.SourceImageEncryptionKey = vSourceImageEncryptionKey
	}
	return nil
}
func postReadExtractInstanceDisksInitializeParamsSourceImageEncryptionKeyFields(r *Instance, o *InstanceDisksInitializeParamsSourceImageEncryptionKey) error {
	return nil
}
func postReadExtractInstanceGuestAcceleratorsFields(r *Instance, o *InstanceGuestAccelerators) error {
	return nil
}
func postReadExtractInstanceNetworkInterfacesFields(r *Instance, o *InstanceNetworkInterfaces) error {
	return nil
}
func postReadExtractInstanceNetworkInterfacesAccessConfigsFields(r *Instance, o *InstanceNetworkInterfacesAccessConfigs) error {
	return nil
}
func postReadExtractInstanceNetworkInterfacesIPv6AccessConfigsFields(r *Instance, o *InstanceNetworkInterfacesIPv6AccessConfigs) error {
	return nil
}
func postReadExtractInstanceNetworkInterfacesAliasIPRangesFields(r *Instance, o *InstanceNetworkInterfacesAliasIPRanges) error {
	return nil
}
func postReadExtractInstanceSchedulingFields(r *Instance, o *InstanceScheduling) error {
	return nil
}
func postReadExtractInstanceServiceAccountsFields(r *Instance, o *InstanceServiceAccounts) error {
	return nil
}
func postReadExtractInstanceShieldedInstanceConfigFields(r *Instance, o *InstanceShieldedInstanceConfig) error {
	return nil
}
