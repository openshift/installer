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
package privateca

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

func (r *CaPool) validate() error {

	if err := dcl.Required(r, "name"); err != nil {
		return err
	}
	if err := dcl.Required(r, "tier"); err != nil {
		return err
	}
	if err := dcl.RequiredParameter(r.Project, "Project"); err != nil {
		return err
	}
	if err := dcl.RequiredParameter(r.Location, "Location"); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(r.IssuancePolicy) {
		if err := r.IssuancePolicy.validate(); err != nil {
			return err
		}
	}
	if !dcl.IsEmptyValueIndirect(r.PublishingOptions) {
		if err := r.PublishingOptions.validate(); err != nil {
			return err
		}
	}
	return nil
}
func (r *CaPoolIssuancePolicy) validate() error {
	if !dcl.IsEmptyValueIndirect(r.AllowedIssuanceModes) {
		if err := r.AllowedIssuanceModes.validate(); err != nil {
			return err
		}
	}
	if !dcl.IsEmptyValueIndirect(r.BaselineValues) {
		if err := r.BaselineValues.validate(); err != nil {
			return err
		}
	}
	if !dcl.IsEmptyValueIndirect(r.IdentityConstraints) {
		if err := r.IdentityConstraints.validate(); err != nil {
			return err
		}
	}
	if !dcl.IsEmptyValueIndirect(r.PassthroughExtensions) {
		if err := r.PassthroughExtensions.validate(); err != nil {
			return err
		}
	}
	return nil
}
func (r *CaPoolIssuancePolicyAllowedKeyTypes) validate() error {
	if err := dcl.ValidateAtMostOneOfFieldsSet([]string{"Rsa", "EllipticCurve"}, r.Rsa, r.EllipticCurve); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(r.Rsa) {
		if err := r.Rsa.validate(); err != nil {
			return err
		}
	}
	if !dcl.IsEmptyValueIndirect(r.EllipticCurve) {
		if err := r.EllipticCurve.validate(); err != nil {
			return err
		}
	}
	return nil
}
func (r *CaPoolIssuancePolicyAllowedKeyTypesRsa) validate() error {
	return nil
}
func (r *CaPoolIssuancePolicyAllowedKeyTypesEllipticCurve) validate() error {
	return nil
}
func (r *CaPoolIssuancePolicyAllowedIssuanceModes) validate() error {
	return nil
}
func (r *CaPoolIssuancePolicyBaselineValues) validate() error {
	if !dcl.IsEmptyValueIndirect(r.KeyUsage) {
		if err := r.KeyUsage.validate(); err != nil {
			return err
		}
	}
	if !dcl.IsEmptyValueIndirect(r.CaOptions) {
		if err := r.CaOptions.validate(); err != nil {
			return err
		}
	}
	return nil
}
func (r *CaPoolIssuancePolicyBaselineValuesKeyUsage) validate() error {
	if !dcl.IsEmptyValueIndirect(r.BaseKeyUsage) {
		if err := r.BaseKeyUsage.validate(); err != nil {
			return err
		}
	}
	if !dcl.IsEmptyValueIndirect(r.ExtendedKeyUsage) {
		if err := r.ExtendedKeyUsage.validate(); err != nil {
			return err
		}
	}
	return nil
}
func (r *CaPoolIssuancePolicyBaselineValuesKeyUsageBaseKeyUsage) validate() error {
	return nil
}
func (r *CaPoolIssuancePolicyBaselineValuesKeyUsageExtendedKeyUsage) validate() error {
	return nil
}
func (r *CaPoolIssuancePolicyBaselineValuesKeyUsageUnknownExtendedKeyUsages) validate() error {
	if err := dcl.Required(r, "objectIdPath"); err != nil {
		return err
	}
	return nil
}
func (r *CaPoolIssuancePolicyBaselineValuesCaOptions) validate() error {
	return nil
}
func (r *CaPoolIssuancePolicyBaselineValuesPolicyIds) validate() error {
	if err := dcl.Required(r, "objectIdPath"); err != nil {
		return err
	}
	return nil
}
func (r *CaPoolIssuancePolicyBaselineValuesAdditionalExtensions) validate() error {
	if err := dcl.Required(r, "objectId"); err != nil {
		return err
	}
	if err := dcl.Required(r, "value"); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(r.ObjectId) {
		if err := r.ObjectId.validate(); err != nil {
			return err
		}
	}
	return nil
}
func (r *CaPoolIssuancePolicyBaselineValuesAdditionalExtensionsObjectId) validate() error {
	if err := dcl.Required(r, "objectIdPath"); err != nil {
		return err
	}
	return nil
}
func (r *CaPoolIssuancePolicyIdentityConstraints) validate() error {
	if err := dcl.Required(r, "allowSubjectPassthrough"); err != nil {
		return err
	}
	if err := dcl.Required(r, "allowSubjectAltNamesPassthrough"); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(r.CelExpression) {
		if err := r.CelExpression.validate(); err != nil {
			return err
		}
	}
	return nil
}
func (r *CaPoolIssuancePolicyIdentityConstraintsCelExpression) validate() error {
	return nil
}
func (r *CaPoolIssuancePolicyPassthroughExtensions) validate() error {
	return nil
}
func (r *CaPoolIssuancePolicyPassthroughExtensionsAdditionalExtensions) validate() error {
	if err := dcl.Required(r, "objectIdPath"); err != nil {
		return err
	}
	return nil
}
func (r *CaPoolPublishingOptions) validate() error {
	return nil
}
func (r *CaPool) basePath() string {
	params := map[string]interface{}{}
	return dcl.Nprintf("https://privateca.googleapis.com/v1/", params)
}

func (r *CaPool) getURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"project":  dcl.ValueOrEmptyString(nr.Project),
		"location": dcl.ValueOrEmptyString(nr.Location),
		"name":     dcl.ValueOrEmptyString(nr.Name),
	}
	return dcl.URL("projects/{{project}}/locations/{{location}}/caPools/{{name}}", nr.basePath(), userBasePath, params), nil
}

func (r *CaPool) listURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"project":  dcl.ValueOrEmptyString(nr.Project),
		"location": dcl.ValueOrEmptyString(nr.Location),
	}
	return dcl.URL("projects/{{project}}/locations/{{location}}/caPools", nr.basePath(), userBasePath, params), nil

}

func (r *CaPool) createURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"project":  dcl.ValueOrEmptyString(nr.Project),
		"location": dcl.ValueOrEmptyString(nr.Location),
		"name":     dcl.ValueOrEmptyString(nr.Name),
	}
	return dcl.URL("projects/{{project}}/locations/{{location}}/caPools?caPoolId={{name}}", nr.basePath(), userBasePath, params), nil

}

func (r *CaPool) deleteURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"project":  dcl.ValueOrEmptyString(nr.Project),
		"location": dcl.ValueOrEmptyString(nr.Location),
		"name":     dcl.ValueOrEmptyString(nr.Name),
	}
	return dcl.URL("projects/{{project}}/locations/{{location}}/caPools/{{name}}", nr.basePath(), userBasePath, params), nil
}

// caPoolApiOperation represents a mutable operation in the underlying REST
// API such as Create, Update, or Delete.
type caPoolApiOperation interface {
	do(context.Context, *CaPool, *Client) error
}

// newUpdateCaPoolUpdateCaPoolRequest creates a request for an
// CaPool resource's UpdateCaPool update type by filling in the update
// fields based on the intended state of the resource.
func newUpdateCaPoolUpdateCaPoolRequest(ctx context.Context, f *CaPool, c *Client) (map[string]interface{}, error) {
	req := map[string]interface{}{}
	res := f
	_ = res

	if v, err := expandCaPoolIssuancePolicy(c, f.IssuancePolicy, res); err != nil {
		return nil, fmt.Errorf("error expanding IssuancePolicy into issuancePolicy: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		req["issuancePolicy"] = v
	}
	if v, err := expandCaPoolPublishingOptions(c, f.PublishingOptions, res); err != nil {
		return nil, fmt.Errorf("error expanding PublishingOptions into publishingOptions: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		req["publishingOptions"] = v
	}
	if v := f.Labels; !dcl.IsEmptyValueIndirect(v) {
		req["labels"] = v
	}
	req["name"] = fmt.Sprintf("projects/%s/locations/%s/caPools/%s", *f.Project, *f.Location, *f.Name)

	return req, nil
}

// marshalUpdateCaPoolUpdateCaPoolRequest converts the update into
// the final JSON request body.
func marshalUpdateCaPoolUpdateCaPoolRequest(c *Client, m map[string]interface{}) ([]byte, error) {

	return json.Marshal(m)
}

type updateCaPoolUpdateCaPoolOperation struct {
	// If the update operation has the REQUIRES_APPLY_OPTIONS trait, this will be populated.
	// Usually it will be nil - this is to prevent us from accidentally depending on apply
	// options, which should usually be unnecessary.
	ApplyOptions []dcl.ApplyOption
	FieldDiffs   []*dcl.FieldDiff
}

// do creates a request and sends it to the appropriate URL. In most operations,
// do will transcribe a subset of the resource into a request object and send a
// PUT request to a single URL.

func (op *updateCaPoolUpdateCaPoolOperation) do(ctx context.Context, r *CaPool, c *Client) error {
	_, err := c.GetCaPool(ctx, r)
	if err != nil {
		return err
	}

	u, err := r.updateURL(c.Config.BasePath, "UpdateCaPool")
	if err != nil {
		return err
	}
	mask := dcl.TopLevelUpdateMask(op.FieldDiffs)
	u, err = dcl.AddQueryParams(u, map[string]string{"updateMask": mask})
	if err != nil {
		return err
	}

	req, err := newUpdateCaPoolUpdateCaPoolRequest(ctx, r, c)
	if err != nil {
		return err
	}

	c.Config.Logger.InfoWithContextf(ctx, "Created update: %#v", req)
	body, err := marshalUpdateCaPoolUpdateCaPoolRequest(c, req)
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

func (c *Client) listCaPoolRaw(ctx context.Context, r *CaPool, pageToken string, pageSize int32) ([]byte, error) {
	u, err := r.urlNormalized().listURL(c.Config.BasePath)
	if err != nil {
		return nil, err
	}

	m := make(map[string]string)
	if pageToken != "" {
		m["pageToken"] = pageToken
	}

	if pageSize != CaPoolMaxPage {
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

type listCaPoolOperation struct {
	CaPools []map[string]interface{} `json:"caPools"`
	Token   string                   `json:"nextPageToken"`
}

func (c *Client) listCaPool(ctx context.Context, r *CaPool, pageToken string, pageSize int32) ([]*CaPool, string, error) {
	b, err := c.listCaPoolRaw(ctx, r, pageToken, pageSize)
	if err != nil {
		return nil, "", err
	}

	var m listCaPoolOperation
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, "", err
	}

	var l []*CaPool
	for _, v := range m.CaPools {
		res, err := unmarshalMapCaPool(v, c, r)
		if err != nil {
			return nil, m.Token, err
		}
		res.Project = r.Project
		res.Location = r.Location
		l = append(l, res)
	}

	return l, m.Token, nil
}

func (c *Client) deleteAllCaPool(ctx context.Context, f func(*CaPool) bool, resources []*CaPool) error {
	var errors []string
	for _, res := range resources {
		if f(res) {
			// We do not want deleteAll to fail on a deletion or else it will stop deleting other resources.
			err := c.DeleteCaPool(ctx, res)
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

type deleteCaPoolOperation struct{}

func (op *deleteCaPoolOperation) do(ctx context.Context, r *CaPool, c *Client) error {
	r, err := c.GetCaPool(ctx, r)
	if err != nil {
		if dcl.IsNotFound(err) {
			c.Config.Logger.InfoWithContextf(ctx, "CaPool not found, returning. Original error: %v", err)
			return nil
		}
		c.Config.Logger.WarningWithContextf(ctx, "GetCaPool checking for existence. error: %v", err)
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
		_, err := c.GetCaPool(ctx, r)
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
type createCaPoolOperation struct {
	response map[string]interface{}
}

func (op *createCaPoolOperation) FirstResponse() (map[string]interface{}, bool) {
	return op.response, len(op.response) > 0
}

func (op *createCaPoolOperation) do(ctx context.Context, r *CaPool, c *Client) error {
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

	if _, err := c.GetCaPool(ctx, r); err != nil {
		c.Config.Logger.WarningWithContextf(ctx, "get returned error: %v", err)
		return err
	}

	return nil
}

func (c *Client) getCaPoolRaw(ctx context.Context, r *CaPool) ([]byte, error) {

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

func (c *Client) caPoolDiffsForRawDesired(ctx context.Context, rawDesired *CaPool, opts ...dcl.ApplyOption) (initial, desired *CaPool, diffs []*dcl.FieldDiff, err error) {
	c.Config.Logger.InfoWithContext(ctx, "Fetching initial state...")
	// First, let us see if the user provided a state hint.  If they did, we will start fetching based on that.
	var fetchState *CaPool
	if sh := dcl.FetchStateHint(opts); sh != nil {
		if r, ok := sh.(*CaPool); !ok {
			c.Config.Logger.WarningWithContextf(ctx, "Initial state hint was of the wrong type; expected CaPool, got %T", sh)
		} else {
			fetchState = r
		}
	}
	if fetchState == nil {
		fetchState = rawDesired
	}

	// 1.2: Retrieval of raw initial state from API
	rawInitial, err := c.GetCaPool(ctx, fetchState)
	if rawInitial == nil {
		if !dcl.IsNotFound(err) {
			c.Config.Logger.WarningWithContextf(ctx, "Failed to retrieve whether a CaPool resource already exists: %s", err)
			return nil, nil, nil, fmt.Errorf("failed to retrieve CaPool resource: %v", err)
		}
		c.Config.Logger.InfoWithContext(ctx, "Found that CaPool resource did not exist.")
		// Perform canonicalization to pick up defaults.
		desired, err = canonicalizeCaPoolDesiredState(rawDesired, rawInitial)
		return nil, desired, nil, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Found initial state for CaPool: %v", rawInitial)
	c.Config.Logger.InfoWithContextf(ctx, "Initial desired state for CaPool: %v", rawDesired)

	// The Get call applies postReadExtract and so the result may contain fields that are not part of API version.
	if err := extractCaPoolFields(rawInitial); err != nil {
		return nil, nil, nil, err
	}

	// 1.3: Canonicalize raw initial state into initial state.
	initial, err = canonicalizeCaPoolInitialState(rawInitial, rawDesired)
	if err != nil {
		return nil, nil, nil, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Canonicalized initial state for CaPool: %v", initial)

	// 1.4: Canonicalize raw desired state into desired state.
	desired, err = canonicalizeCaPoolDesiredState(rawDesired, rawInitial, opts...)
	if err != nil {
		return nil, nil, nil, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Canonicalized desired state for CaPool: %v", desired)

	// 2.1: Comparison of initial and desired state.
	diffs, err = diffCaPool(c, desired, initial, opts...)
	return initial, desired, diffs, err
}

func canonicalizeCaPoolInitialState(rawInitial, rawDesired *CaPool) (*CaPool, error) {
	// TODO(magic-modules-eng): write canonicalizer once relevant traits are added.
	return rawInitial, nil
}

/*
* Canonicalizers
*
* These are responsible for converting either a user-specified config or a
* GCP API response to a standard format that can be used for difference checking.
* */

func canonicalizeCaPoolDesiredState(rawDesired, rawInitial *CaPool, opts ...dcl.ApplyOption) (*CaPool, error) {

	if rawInitial == nil {
		// Since the initial state is empty, the desired state is all we have.
		// We canonicalize the remaining nested objects with nil to pick up defaults.
		rawDesired.IssuancePolicy = canonicalizeCaPoolIssuancePolicy(rawDesired.IssuancePolicy, nil, opts...)
		rawDesired.PublishingOptions = canonicalizeCaPoolPublishingOptions(rawDesired.PublishingOptions, nil, opts...)

		return rawDesired, nil
	}
	canonicalDesired := &CaPool{}
	if dcl.PartialSelfLinkToSelfLink(rawDesired.Name, rawInitial.Name) {
		canonicalDesired.Name = rawInitial.Name
	} else {
		canonicalDesired.Name = rawDesired.Name
	}
	if dcl.IsZeroValue(rawDesired.Tier) || (dcl.IsEmptyValueIndirect(rawDesired.Tier) && dcl.IsEmptyValueIndirect(rawInitial.Tier)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		canonicalDesired.Tier = rawInitial.Tier
	} else {
		canonicalDesired.Tier = rawDesired.Tier
	}
	canonicalDesired.IssuancePolicy = canonicalizeCaPoolIssuancePolicy(rawDesired.IssuancePolicy, rawInitial.IssuancePolicy, opts...)
	canonicalDesired.PublishingOptions = canonicalizeCaPoolPublishingOptions(rawDesired.PublishingOptions, rawInitial.PublishingOptions, opts...)
	if dcl.IsZeroValue(rawDesired.Labels) || (dcl.IsEmptyValueIndirect(rawDesired.Labels) && dcl.IsEmptyValueIndirect(rawInitial.Labels)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		canonicalDesired.Labels = rawInitial.Labels
	} else {
		canonicalDesired.Labels = rawDesired.Labels
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

func canonicalizeCaPoolNewState(c *Client, rawNew, rawDesired *CaPool) (*CaPool, error) {

	if dcl.IsEmptyValueIndirect(rawNew.Name) && dcl.IsEmptyValueIndirect(rawDesired.Name) {
		rawNew.Name = rawDesired.Name
	} else {
		if dcl.PartialSelfLinkToSelfLink(rawDesired.Name, rawNew.Name) {
			rawNew.Name = rawDesired.Name
		}
	}

	if dcl.IsEmptyValueIndirect(rawNew.Tier) && dcl.IsEmptyValueIndirect(rawDesired.Tier) {
		rawNew.Tier = rawDesired.Tier
	} else {
	}

	if dcl.IsEmptyValueIndirect(rawNew.IssuancePolicy) && dcl.IsEmptyValueIndirect(rawDesired.IssuancePolicy) {
		rawNew.IssuancePolicy = rawDesired.IssuancePolicy
	} else {
		rawNew.IssuancePolicy = canonicalizeNewCaPoolIssuancePolicy(c, rawDesired.IssuancePolicy, rawNew.IssuancePolicy)
	}

	if dcl.IsEmptyValueIndirect(rawNew.PublishingOptions) && dcl.IsEmptyValueIndirect(rawDesired.PublishingOptions) {
		rawNew.PublishingOptions = rawDesired.PublishingOptions
	} else {
		rawNew.PublishingOptions = canonicalizeNewCaPoolPublishingOptions(c, rawDesired.PublishingOptions, rawNew.PublishingOptions)
	}

	if dcl.IsEmptyValueIndirect(rawNew.Labels) && dcl.IsEmptyValueIndirect(rawDesired.Labels) {
		rawNew.Labels = rawDesired.Labels
	} else {
	}

	rawNew.Project = rawDesired.Project

	rawNew.Location = rawDesired.Location

	return rawNew, nil
}

func canonicalizeCaPoolIssuancePolicy(des, initial *CaPoolIssuancePolicy, opts ...dcl.ApplyOption) *CaPoolIssuancePolicy {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &CaPoolIssuancePolicy{}

	cDes.AllowedKeyTypes = canonicalizeCaPoolIssuancePolicyAllowedKeyTypesSlice(des.AllowedKeyTypes, initial.AllowedKeyTypes, opts...)
	if dcl.StringCanonicalize(des.MaximumLifetime, initial.MaximumLifetime) || dcl.IsZeroValue(des.MaximumLifetime) {
		cDes.MaximumLifetime = initial.MaximumLifetime
	} else {
		cDes.MaximumLifetime = des.MaximumLifetime
	}
	cDes.AllowedIssuanceModes = canonicalizeCaPoolIssuancePolicyAllowedIssuanceModes(des.AllowedIssuanceModes, initial.AllowedIssuanceModes, opts...)
	cDes.BaselineValues = canonicalizeCaPoolIssuancePolicyBaselineValues(des.BaselineValues, initial.BaselineValues, opts...)
	cDes.IdentityConstraints = canonicalizeCaPoolIssuancePolicyIdentityConstraints(des.IdentityConstraints, initial.IdentityConstraints, opts...)
	cDes.PassthroughExtensions = canonicalizeCaPoolIssuancePolicyPassthroughExtensions(des.PassthroughExtensions, initial.PassthroughExtensions, opts...)

	return cDes
}

func canonicalizeCaPoolIssuancePolicySlice(des, initial []CaPoolIssuancePolicy, opts ...dcl.ApplyOption) []CaPoolIssuancePolicy {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]CaPoolIssuancePolicy, 0, len(des))
		for _, d := range des {
			cd := canonicalizeCaPoolIssuancePolicy(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]CaPoolIssuancePolicy, 0, len(des))
	for i, d := range des {
		cd := canonicalizeCaPoolIssuancePolicy(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewCaPoolIssuancePolicy(c *Client, des, nw *CaPoolIssuancePolicy) *CaPoolIssuancePolicy {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for CaPoolIssuancePolicy while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	nw.AllowedKeyTypes = canonicalizeNewCaPoolIssuancePolicyAllowedKeyTypesSlice(c, des.AllowedKeyTypes, nw.AllowedKeyTypes)
	if dcl.StringCanonicalize(des.MaximumLifetime, nw.MaximumLifetime) {
		nw.MaximumLifetime = des.MaximumLifetime
	}
	nw.AllowedIssuanceModes = canonicalizeNewCaPoolIssuancePolicyAllowedIssuanceModes(c, des.AllowedIssuanceModes, nw.AllowedIssuanceModes)
	nw.BaselineValues = canonicalizeNewCaPoolIssuancePolicyBaselineValues(c, des.BaselineValues, nw.BaselineValues)
	nw.IdentityConstraints = canonicalizeNewCaPoolIssuancePolicyIdentityConstraints(c, des.IdentityConstraints, nw.IdentityConstraints)
	nw.PassthroughExtensions = canonicalizeNewCaPoolIssuancePolicyPassthroughExtensions(c, des.PassthroughExtensions, nw.PassthroughExtensions)

	return nw
}

func canonicalizeNewCaPoolIssuancePolicySet(c *Client, des, nw []CaPoolIssuancePolicy) []CaPoolIssuancePolicy {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []CaPoolIssuancePolicy
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareCaPoolIssuancePolicyNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewCaPoolIssuancePolicy(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewCaPoolIssuancePolicySlice(c *Client, des, nw []CaPoolIssuancePolicy) []CaPoolIssuancePolicy {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []CaPoolIssuancePolicy
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewCaPoolIssuancePolicy(c, &d, &n))
	}

	return items
}

func canonicalizeCaPoolIssuancePolicyAllowedKeyTypes(des, initial *CaPoolIssuancePolicyAllowedKeyTypes, opts ...dcl.ApplyOption) *CaPoolIssuancePolicyAllowedKeyTypes {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if des.Rsa != nil || (initial != nil && initial.Rsa != nil) {
		// Check if anything else is set.
		if dcl.AnySet(des.EllipticCurve) {
			des.Rsa = nil
			if initial != nil {
				initial.Rsa = nil
			}
		}
	}

	if des.EllipticCurve != nil || (initial != nil && initial.EllipticCurve != nil) {
		// Check if anything else is set.
		if dcl.AnySet(des.Rsa) {
			des.EllipticCurve = nil
			if initial != nil {
				initial.EllipticCurve = nil
			}
		}
	}

	if initial == nil {
		return des
	}

	cDes := &CaPoolIssuancePolicyAllowedKeyTypes{}

	cDes.Rsa = canonicalizeCaPoolIssuancePolicyAllowedKeyTypesRsa(des.Rsa, initial.Rsa, opts...)
	cDes.EllipticCurve = canonicalizeCaPoolIssuancePolicyAllowedKeyTypesEllipticCurve(des.EllipticCurve, initial.EllipticCurve, opts...)

	return cDes
}

func canonicalizeCaPoolIssuancePolicyAllowedKeyTypesSlice(des, initial []CaPoolIssuancePolicyAllowedKeyTypes, opts ...dcl.ApplyOption) []CaPoolIssuancePolicyAllowedKeyTypes {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]CaPoolIssuancePolicyAllowedKeyTypes, 0, len(des))
		for _, d := range des {
			cd := canonicalizeCaPoolIssuancePolicyAllowedKeyTypes(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]CaPoolIssuancePolicyAllowedKeyTypes, 0, len(des))
	for i, d := range des {
		cd := canonicalizeCaPoolIssuancePolicyAllowedKeyTypes(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewCaPoolIssuancePolicyAllowedKeyTypes(c *Client, des, nw *CaPoolIssuancePolicyAllowedKeyTypes) *CaPoolIssuancePolicyAllowedKeyTypes {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for CaPoolIssuancePolicyAllowedKeyTypes while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	nw.Rsa = canonicalizeNewCaPoolIssuancePolicyAllowedKeyTypesRsa(c, des.Rsa, nw.Rsa)
	nw.EllipticCurve = canonicalizeNewCaPoolIssuancePolicyAllowedKeyTypesEllipticCurve(c, des.EllipticCurve, nw.EllipticCurve)

	return nw
}

func canonicalizeNewCaPoolIssuancePolicyAllowedKeyTypesSet(c *Client, des, nw []CaPoolIssuancePolicyAllowedKeyTypes) []CaPoolIssuancePolicyAllowedKeyTypes {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []CaPoolIssuancePolicyAllowedKeyTypes
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareCaPoolIssuancePolicyAllowedKeyTypesNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewCaPoolIssuancePolicyAllowedKeyTypes(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewCaPoolIssuancePolicyAllowedKeyTypesSlice(c *Client, des, nw []CaPoolIssuancePolicyAllowedKeyTypes) []CaPoolIssuancePolicyAllowedKeyTypes {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []CaPoolIssuancePolicyAllowedKeyTypes
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewCaPoolIssuancePolicyAllowedKeyTypes(c, &d, &n))
	}

	return items
}

func canonicalizeCaPoolIssuancePolicyAllowedKeyTypesRsa(des, initial *CaPoolIssuancePolicyAllowedKeyTypesRsa, opts ...dcl.ApplyOption) *CaPoolIssuancePolicyAllowedKeyTypesRsa {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &CaPoolIssuancePolicyAllowedKeyTypesRsa{}

	if dcl.IsZeroValue(des.MinModulusSize) || (dcl.IsEmptyValueIndirect(des.MinModulusSize) && dcl.IsEmptyValueIndirect(initial.MinModulusSize)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.MinModulusSize = initial.MinModulusSize
	} else {
		cDes.MinModulusSize = des.MinModulusSize
	}
	if dcl.IsZeroValue(des.MaxModulusSize) || (dcl.IsEmptyValueIndirect(des.MaxModulusSize) && dcl.IsEmptyValueIndirect(initial.MaxModulusSize)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.MaxModulusSize = initial.MaxModulusSize
	} else {
		cDes.MaxModulusSize = des.MaxModulusSize
	}

	return cDes
}

func canonicalizeCaPoolIssuancePolicyAllowedKeyTypesRsaSlice(des, initial []CaPoolIssuancePolicyAllowedKeyTypesRsa, opts ...dcl.ApplyOption) []CaPoolIssuancePolicyAllowedKeyTypesRsa {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]CaPoolIssuancePolicyAllowedKeyTypesRsa, 0, len(des))
		for _, d := range des {
			cd := canonicalizeCaPoolIssuancePolicyAllowedKeyTypesRsa(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]CaPoolIssuancePolicyAllowedKeyTypesRsa, 0, len(des))
	for i, d := range des {
		cd := canonicalizeCaPoolIssuancePolicyAllowedKeyTypesRsa(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewCaPoolIssuancePolicyAllowedKeyTypesRsa(c *Client, des, nw *CaPoolIssuancePolicyAllowedKeyTypesRsa) *CaPoolIssuancePolicyAllowedKeyTypesRsa {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for CaPoolIssuancePolicyAllowedKeyTypesRsa while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	return nw
}

func canonicalizeNewCaPoolIssuancePolicyAllowedKeyTypesRsaSet(c *Client, des, nw []CaPoolIssuancePolicyAllowedKeyTypesRsa) []CaPoolIssuancePolicyAllowedKeyTypesRsa {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []CaPoolIssuancePolicyAllowedKeyTypesRsa
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareCaPoolIssuancePolicyAllowedKeyTypesRsaNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewCaPoolIssuancePolicyAllowedKeyTypesRsa(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewCaPoolIssuancePolicyAllowedKeyTypesRsaSlice(c *Client, des, nw []CaPoolIssuancePolicyAllowedKeyTypesRsa) []CaPoolIssuancePolicyAllowedKeyTypesRsa {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []CaPoolIssuancePolicyAllowedKeyTypesRsa
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewCaPoolIssuancePolicyAllowedKeyTypesRsa(c, &d, &n))
	}

	return items
}

func canonicalizeCaPoolIssuancePolicyAllowedKeyTypesEllipticCurve(des, initial *CaPoolIssuancePolicyAllowedKeyTypesEllipticCurve, opts ...dcl.ApplyOption) *CaPoolIssuancePolicyAllowedKeyTypesEllipticCurve {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &CaPoolIssuancePolicyAllowedKeyTypesEllipticCurve{}

	if dcl.IsZeroValue(des.SignatureAlgorithm) || (dcl.IsEmptyValueIndirect(des.SignatureAlgorithm) && dcl.IsEmptyValueIndirect(initial.SignatureAlgorithm)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.SignatureAlgorithm = initial.SignatureAlgorithm
	} else {
		cDes.SignatureAlgorithm = des.SignatureAlgorithm
	}

	return cDes
}

func canonicalizeCaPoolIssuancePolicyAllowedKeyTypesEllipticCurveSlice(des, initial []CaPoolIssuancePolicyAllowedKeyTypesEllipticCurve, opts ...dcl.ApplyOption) []CaPoolIssuancePolicyAllowedKeyTypesEllipticCurve {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]CaPoolIssuancePolicyAllowedKeyTypesEllipticCurve, 0, len(des))
		for _, d := range des {
			cd := canonicalizeCaPoolIssuancePolicyAllowedKeyTypesEllipticCurve(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]CaPoolIssuancePolicyAllowedKeyTypesEllipticCurve, 0, len(des))
	for i, d := range des {
		cd := canonicalizeCaPoolIssuancePolicyAllowedKeyTypesEllipticCurve(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewCaPoolIssuancePolicyAllowedKeyTypesEllipticCurve(c *Client, des, nw *CaPoolIssuancePolicyAllowedKeyTypesEllipticCurve) *CaPoolIssuancePolicyAllowedKeyTypesEllipticCurve {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for CaPoolIssuancePolicyAllowedKeyTypesEllipticCurve while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	return nw
}

func canonicalizeNewCaPoolIssuancePolicyAllowedKeyTypesEllipticCurveSet(c *Client, des, nw []CaPoolIssuancePolicyAllowedKeyTypesEllipticCurve) []CaPoolIssuancePolicyAllowedKeyTypesEllipticCurve {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []CaPoolIssuancePolicyAllowedKeyTypesEllipticCurve
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareCaPoolIssuancePolicyAllowedKeyTypesEllipticCurveNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewCaPoolIssuancePolicyAllowedKeyTypesEllipticCurve(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewCaPoolIssuancePolicyAllowedKeyTypesEllipticCurveSlice(c *Client, des, nw []CaPoolIssuancePolicyAllowedKeyTypesEllipticCurve) []CaPoolIssuancePolicyAllowedKeyTypesEllipticCurve {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []CaPoolIssuancePolicyAllowedKeyTypesEllipticCurve
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewCaPoolIssuancePolicyAllowedKeyTypesEllipticCurve(c, &d, &n))
	}

	return items
}

func canonicalizeCaPoolIssuancePolicyAllowedIssuanceModes(des, initial *CaPoolIssuancePolicyAllowedIssuanceModes, opts ...dcl.ApplyOption) *CaPoolIssuancePolicyAllowedIssuanceModes {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &CaPoolIssuancePolicyAllowedIssuanceModes{}

	if dcl.BoolCanonicalize(des.AllowCsrBasedIssuance, initial.AllowCsrBasedIssuance) || dcl.IsZeroValue(des.AllowCsrBasedIssuance) {
		cDes.AllowCsrBasedIssuance = initial.AllowCsrBasedIssuance
	} else {
		cDes.AllowCsrBasedIssuance = des.AllowCsrBasedIssuance
	}
	if dcl.BoolCanonicalize(des.AllowConfigBasedIssuance, initial.AllowConfigBasedIssuance) || dcl.IsZeroValue(des.AllowConfigBasedIssuance) {
		cDes.AllowConfigBasedIssuance = initial.AllowConfigBasedIssuance
	} else {
		cDes.AllowConfigBasedIssuance = des.AllowConfigBasedIssuance
	}

	return cDes
}

func canonicalizeCaPoolIssuancePolicyAllowedIssuanceModesSlice(des, initial []CaPoolIssuancePolicyAllowedIssuanceModes, opts ...dcl.ApplyOption) []CaPoolIssuancePolicyAllowedIssuanceModes {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]CaPoolIssuancePolicyAllowedIssuanceModes, 0, len(des))
		for _, d := range des {
			cd := canonicalizeCaPoolIssuancePolicyAllowedIssuanceModes(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]CaPoolIssuancePolicyAllowedIssuanceModes, 0, len(des))
	for i, d := range des {
		cd := canonicalizeCaPoolIssuancePolicyAllowedIssuanceModes(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewCaPoolIssuancePolicyAllowedIssuanceModes(c *Client, des, nw *CaPoolIssuancePolicyAllowedIssuanceModes) *CaPoolIssuancePolicyAllowedIssuanceModes {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for CaPoolIssuancePolicyAllowedIssuanceModes while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.BoolCanonicalize(des.AllowCsrBasedIssuance, nw.AllowCsrBasedIssuance) {
		nw.AllowCsrBasedIssuance = des.AllowCsrBasedIssuance
	}
	if dcl.BoolCanonicalize(des.AllowConfigBasedIssuance, nw.AllowConfigBasedIssuance) {
		nw.AllowConfigBasedIssuance = des.AllowConfigBasedIssuance
	}

	return nw
}

func canonicalizeNewCaPoolIssuancePolicyAllowedIssuanceModesSet(c *Client, des, nw []CaPoolIssuancePolicyAllowedIssuanceModes) []CaPoolIssuancePolicyAllowedIssuanceModes {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []CaPoolIssuancePolicyAllowedIssuanceModes
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareCaPoolIssuancePolicyAllowedIssuanceModesNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewCaPoolIssuancePolicyAllowedIssuanceModes(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewCaPoolIssuancePolicyAllowedIssuanceModesSlice(c *Client, des, nw []CaPoolIssuancePolicyAllowedIssuanceModes) []CaPoolIssuancePolicyAllowedIssuanceModes {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []CaPoolIssuancePolicyAllowedIssuanceModes
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewCaPoolIssuancePolicyAllowedIssuanceModes(c, &d, &n))
	}

	return items
}

func canonicalizeCaPoolIssuancePolicyBaselineValues(des, initial *CaPoolIssuancePolicyBaselineValues, opts ...dcl.ApplyOption) *CaPoolIssuancePolicyBaselineValues {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &CaPoolIssuancePolicyBaselineValues{}

	cDes.KeyUsage = canonicalizeCaPoolIssuancePolicyBaselineValuesKeyUsage(des.KeyUsage, initial.KeyUsage, opts...)
	cDes.CaOptions = canonicalizeCaPoolIssuancePolicyBaselineValuesCaOptions(des.CaOptions, initial.CaOptions, opts...)
	cDes.PolicyIds = canonicalizeCaPoolIssuancePolicyBaselineValuesPolicyIdsSlice(des.PolicyIds, initial.PolicyIds, opts...)
	if dcl.StringArrayCanonicalize(des.AiaOcspServers, initial.AiaOcspServers) {
		cDes.AiaOcspServers = initial.AiaOcspServers
	} else {
		cDes.AiaOcspServers = des.AiaOcspServers
	}
	cDes.AdditionalExtensions = canonicalizeCaPoolIssuancePolicyBaselineValuesAdditionalExtensionsSlice(des.AdditionalExtensions, initial.AdditionalExtensions, opts...)

	return cDes
}

func canonicalizeCaPoolIssuancePolicyBaselineValuesSlice(des, initial []CaPoolIssuancePolicyBaselineValues, opts ...dcl.ApplyOption) []CaPoolIssuancePolicyBaselineValues {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]CaPoolIssuancePolicyBaselineValues, 0, len(des))
		for _, d := range des {
			cd := canonicalizeCaPoolIssuancePolicyBaselineValues(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]CaPoolIssuancePolicyBaselineValues, 0, len(des))
	for i, d := range des {
		cd := canonicalizeCaPoolIssuancePolicyBaselineValues(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewCaPoolIssuancePolicyBaselineValues(c *Client, des, nw *CaPoolIssuancePolicyBaselineValues) *CaPoolIssuancePolicyBaselineValues {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for CaPoolIssuancePolicyBaselineValues while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	nw.KeyUsage = canonicalizeNewCaPoolIssuancePolicyBaselineValuesKeyUsage(c, des.KeyUsage, nw.KeyUsage)
	nw.CaOptions = canonicalizeNewCaPoolIssuancePolicyBaselineValuesCaOptions(c, des.CaOptions, nw.CaOptions)
	nw.PolicyIds = canonicalizeNewCaPoolIssuancePolicyBaselineValuesPolicyIdsSlice(c, des.PolicyIds, nw.PolicyIds)
	if dcl.StringArrayCanonicalize(des.AiaOcspServers, nw.AiaOcspServers) {
		nw.AiaOcspServers = des.AiaOcspServers
	}
	nw.AdditionalExtensions = canonicalizeNewCaPoolIssuancePolicyBaselineValuesAdditionalExtensionsSlice(c, des.AdditionalExtensions, nw.AdditionalExtensions)

	return nw
}

func canonicalizeNewCaPoolIssuancePolicyBaselineValuesSet(c *Client, des, nw []CaPoolIssuancePolicyBaselineValues) []CaPoolIssuancePolicyBaselineValues {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []CaPoolIssuancePolicyBaselineValues
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareCaPoolIssuancePolicyBaselineValuesNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewCaPoolIssuancePolicyBaselineValues(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewCaPoolIssuancePolicyBaselineValuesSlice(c *Client, des, nw []CaPoolIssuancePolicyBaselineValues) []CaPoolIssuancePolicyBaselineValues {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []CaPoolIssuancePolicyBaselineValues
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewCaPoolIssuancePolicyBaselineValues(c, &d, &n))
	}

	return items
}

func canonicalizeCaPoolIssuancePolicyBaselineValuesKeyUsage(des, initial *CaPoolIssuancePolicyBaselineValuesKeyUsage, opts ...dcl.ApplyOption) *CaPoolIssuancePolicyBaselineValuesKeyUsage {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &CaPoolIssuancePolicyBaselineValuesKeyUsage{}

	cDes.BaseKeyUsage = canonicalizeCaPoolIssuancePolicyBaselineValuesKeyUsageBaseKeyUsage(des.BaseKeyUsage, initial.BaseKeyUsage, opts...)
	cDes.ExtendedKeyUsage = canonicalizeCaPoolIssuancePolicyBaselineValuesKeyUsageExtendedKeyUsage(des.ExtendedKeyUsage, initial.ExtendedKeyUsage, opts...)
	cDes.UnknownExtendedKeyUsages = canonicalizeCaPoolIssuancePolicyBaselineValuesKeyUsageUnknownExtendedKeyUsagesSlice(des.UnknownExtendedKeyUsages, initial.UnknownExtendedKeyUsages, opts...)

	return cDes
}

func canonicalizeCaPoolIssuancePolicyBaselineValuesKeyUsageSlice(des, initial []CaPoolIssuancePolicyBaselineValuesKeyUsage, opts ...dcl.ApplyOption) []CaPoolIssuancePolicyBaselineValuesKeyUsage {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]CaPoolIssuancePolicyBaselineValuesKeyUsage, 0, len(des))
		for _, d := range des {
			cd := canonicalizeCaPoolIssuancePolicyBaselineValuesKeyUsage(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]CaPoolIssuancePolicyBaselineValuesKeyUsage, 0, len(des))
	for i, d := range des {
		cd := canonicalizeCaPoolIssuancePolicyBaselineValuesKeyUsage(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewCaPoolIssuancePolicyBaselineValuesKeyUsage(c *Client, des, nw *CaPoolIssuancePolicyBaselineValuesKeyUsage) *CaPoolIssuancePolicyBaselineValuesKeyUsage {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for CaPoolIssuancePolicyBaselineValuesKeyUsage while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	nw.BaseKeyUsage = canonicalizeNewCaPoolIssuancePolicyBaselineValuesKeyUsageBaseKeyUsage(c, des.BaseKeyUsage, nw.BaseKeyUsage)
	nw.ExtendedKeyUsage = canonicalizeNewCaPoolIssuancePolicyBaselineValuesKeyUsageExtendedKeyUsage(c, des.ExtendedKeyUsage, nw.ExtendedKeyUsage)
	nw.UnknownExtendedKeyUsages = canonicalizeNewCaPoolIssuancePolicyBaselineValuesKeyUsageUnknownExtendedKeyUsagesSlice(c, des.UnknownExtendedKeyUsages, nw.UnknownExtendedKeyUsages)

	return nw
}

func canonicalizeNewCaPoolIssuancePolicyBaselineValuesKeyUsageSet(c *Client, des, nw []CaPoolIssuancePolicyBaselineValuesKeyUsage) []CaPoolIssuancePolicyBaselineValuesKeyUsage {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []CaPoolIssuancePolicyBaselineValuesKeyUsage
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareCaPoolIssuancePolicyBaselineValuesKeyUsageNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewCaPoolIssuancePolicyBaselineValuesKeyUsage(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewCaPoolIssuancePolicyBaselineValuesKeyUsageSlice(c *Client, des, nw []CaPoolIssuancePolicyBaselineValuesKeyUsage) []CaPoolIssuancePolicyBaselineValuesKeyUsage {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []CaPoolIssuancePolicyBaselineValuesKeyUsage
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewCaPoolIssuancePolicyBaselineValuesKeyUsage(c, &d, &n))
	}

	return items
}

func canonicalizeCaPoolIssuancePolicyBaselineValuesKeyUsageBaseKeyUsage(des, initial *CaPoolIssuancePolicyBaselineValuesKeyUsageBaseKeyUsage, opts ...dcl.ApplyOption) *CaPoolIssuancePolicyBaselineValuesKeyUsageBaseKeyUsage {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &CaPoolIssuancePolicyBaselineValuesKeyUsageBaseKeyUsage{}

	if dcl.BoolCanonicalize(des.DigitalSignature, initial.DigitalSignature) || dcl.IsZeroValue(des.DigitalSignature) {
		cDes.DigitalSignature = initial.DigitalSignature
	} else {
		cDes.DigitalSignature = des.DigitalSignature
	}
	if dcl.BoolCanonicalize(des.ContentCommitment, initial.ContentCommitment) || dcl.IsZeroValue(des.ContentCommitment) {
		cDes.ContentCommitment = initial.ContentCommitment
	} else {
		cDes.ContentCommitment = des.ContentCommitment
	}
	if dcl.BoolCanonicalize(des.KeyEncipherment, initial.KeyEncipherment) || dcl.IsZeroValue(des.KeyEncipherment) {
		cDes.KeyEncipherment = initial.KeyEncipherment
	} else {
		cDes.KeyEncipherment = des.KeyEncipherment
	}
	if dcl.BoolCanonicalize(des.DataEncipherment, initial.DataEncipherment) || dcl.IsZeroValue(des.DataEncipherment) {
		cDes.DataEncipherment = initial.DataEncipherment
	} else {
		cDes.DataEncipherment = des.DataEncipherment
	}
	if dcl.BoolCanonicalize(des.KeyAgreement, initial.KeyAgreement) || dcl.IsZeroValue(des.KeyAgreement) {
		cDes.KeyAgreement = initial.KeyAgreement
	} else {
		cDes.KeyAgreement = des.KeyAgreement
	}
	if dcl.BoolCanonicalize(des.CertSign, initial.CertSign) || dcl.IsZeroValue(des.CertSign) {
		cDes.CertSign = initial.CertSign
	} else {
		cDes.CertSign = des.CertSign
	}
	if dcl.BoolCanonicalize(des.CrlSign, initial.CrlSign) || dcl.IsZeroValue(des.CrlSign) {
		cDes.CrlSign = initial.CrlSign
	} else {
		cDes.CrlSign = des.CrlSign
	}
	if dcl.BoolCanonicalize(des.EncipherOnly, initial.EncipherOnly) || dcl.IsZeroValue(des.EncipherOnly) {
		cDes.EncipherOnly = initial.EncipherOnly
	} else {
		cDes.EncipherOnly = des.EncipherOnly
	}
	if dcl.BoolCanonicalize(des.DecipherOnly, initial.DecipherOnly) || dcl.IsZeroValue(des.DecipherOnly) {
		cDes.DecipherOnly = initial.DecipherOnly
	} else {
		cDes.DecipherOnly = des.DecipherOnly
	}

	return cDes
}

func canonicalizeCaPoolIssuancePolicyBaselineValuesKeyUsageBaseKeyUsageSlice(des, initial []CaPoolIssuancePolicyBaselineValuesKeyUsageBaseKeyUsage, opts ...dcl.ApplyOption) []CaPoolIssuancePolicyBaselineValuesKeyUsageBaseKeyUsage {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]CaPoolIssuancePolicyBaselineValuesKeyUsageBaseKeyUsage, 0, len(des))
		for _, d := range des {
			cd := canonicalizeCaPoolIssuancePolicyBaselineValuesKeyUsageBaseKeyUsage(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]CaPoolIssuancePolicyBaselineValuesKeyUsageBaseKeyUsage, 0, len(des))
	for i, d := range des {
		cd := canonicalizeCaPoolIssuancePolicyBaselineValuesKeyUsageBaseKeyUsage(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewCaPoolIssuancePolicyBaselineValuesKeyUsageBaseKeyUsage(c *Client, des, nw *CaPoolIssuancePolicyBaselineValuesKeyUsageBaseKeyUsage) *CaPoolIssuancePolicyBaselineValuesKeyUsageBaseKeyUsage {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for CaPoolIssuancePolicyBaselineValuesKeyUsageBaseKeyUsage while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.BoolCanonicalize(des.DigitalSignature, nw.DigitalSignature) {
		nw.DigitalSignature = des.DigitalSignature
	}
	if dcl.BoolCanonicalize(des.ContentCommitment, nw.ContentCommitment) {
		nw.ContentCommitment = des.ContentCommitment
	}
	if dcl.BoolCanonicalize(des.KeyEncipherment, nw.KeyEncipherment) {
		nw.KeyEncipherment = des.KeyEncipherment
	}
	if dcl.BoolCanonicalize(des.DataEncipherment, nw.DataEncipherment) {
		nw.DataEncipherment = des.DataEncipherment
	}
	if dcl.BoolCanonicalize(des.KeyAgreement, nw.KeyAgreement) {
		nw.KeyAgreement = des.KeyAgreement
	}
	if dcl.BoolCanonicalize(des.CertSign, nw.CertSign) {
		nw.CertSign = des.CertSign
	}
	if dcl.BoolCanonicalize(des.CrlSign, nw.CrlSign) {
		nw.CrlSign = des.CrlSign
	}
	if dcl.BoolCanonicalize(des.EncipherOnly, nw.EncipherOnly) {
		nw.EncipherOnly = des.EncipherOnly
	}
	if dcl.BoolCanonicalize(des.DecipherOnly, nw.DecipherOnly) {
		nw.DecipherOnly = des.DecipherOnly
	}

	return nw
}

func canonicalizeNewCaPoolIssuancePolicyBaselineValuesKeyUsageBaseKeyUsageSet(c *Client, des, nw []CaPoolIssuancePolicyBaselineValuesKeyUsageBaseKeyUsage) []CaPoolIssuancePolicyBaselineValuesKeyUsageBaseKeyUsage {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []CaPoolIssuancePolicyBaselineValuesKeyUsageBaseKeyUsage
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareCaPoolIssuancePolicyBaselineValuesKeyUsageBaseKeyUsageNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewCaPoolIssuancePolicyBaselineValuesKeyUsageBaseKeyUsage(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewCaPoolIssuancePolicyBaselineValuesKeyUsageBaseKeyUsageSlice(c *Client, des, nw []CaPoolIssuancePolicyBaselineValuesKeyUsageBaseKeyUsage) []CaPoolIssuancePolicyBaselineValuesKeyUsageBaseKeyUsage {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []CaPoolIssuancePolicyBaselineValuesKeyUsageBaseKeyUsage
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewCaPoolIssuancePolicyBaselineValuesKeyUsageBaseKeyUsage(c, &d, &n))
	}

	return items
}

func canonicalizeCaPoolIssuancePolicyBaselineValuesKeyUsageExtendedKeyUsage(des, initial *CaPoolIssuancePolicyBaselineValuesKeyUsageExtendedKeyUsage, opts ...dcl.ApplyOption) *CaPoolIssuancePolicyBaselineValuesKeyUsageExtendedKeyUsage {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &CaPoolIssuancePolicyBaselineValuesKeyUsageExtendedKeyUsage{}

	if dcl.BoolCanonicalize(des.ServerAuth, initial.ServerAuth) || dcl.IsZeroValue(des.ServerAuth) {
		cDes.ServerAuth = initial.ServerAuth
	} else {
		cDes.ServerAuth = des.ServerAuth
	}
	if dcl.BoolCanonicalize(des.ClientAuth, initial.ClientAuth) || dcl.IsZeroValue(des.ClientAuth) {
		cDes.ClientAuth = initial.ClientAuth
	} else {
		cDes.ClientAuth = des.ClientAuth
	}
	if dcl.BoolCanonicalize(des.CodeSigning, initial.CodeSigning) || dcl.IsZeroValue(des.CodeSigning) {
		cDes.CodeSigning = initial.CodeSigning
	} else {
		cDes.CodeSigning = des.CodeSigning
	}
	if dcl.BoolCanonicalize(des.EmailProtection, initial.EmailProtection) || dcl.IsZeroValue(des.EmailProtection) {
		cDes.EmailProtection = initial.EmailProtection
	} else {
		cDes.EmailProtection = des.EmailProtection
	}
	if dcl.BoolCanonicalize(des.TimeStamping, initial.TimeStamping) || dcl.IsZeroValue(des.TimeStamping) {
		cDes.TimeStamping = initial.TimeStamping
	} else {
		cDes.TimeStamping = des.TimeStamping
	}
	if dcl.BoolCanonicalize(des.OcspSigning, initial.OcspSigning) || dcl.IsZeroValue(des.OcspSigning) {
		cDes.OcspSigning = initial.OcspSigning
	} else {
		cDes.OcspSigning = des.OcspSigning
	}

	return cDes
}

func canonicalizeCaPoolIssuancePolicyBaselineValuesKeyUsageExtendedKeyUsageSlice(des, initial []CaPoolIssuancePolicyBaselineValuesKeyUsageExtendedKeyUsage, opts ...dcl.ApplyOption) []CaPoolIssuancePolicyBaselineValuesKeyUsageExtendedKeyUsage {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]CaPoolIssuancePolicyBaselineValuesKeyUsageExtendedKeyUsage, 0, len(des))
		for _, d := range des {
			cd := canonicalizeCaPoolIssuancePolicyBaselineValuesKeyUsageExtendedKeyUsage(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]CaPoolIssuancePolicyBaselineValuesKeyUsageExtendedKeyUsage, 0, len(des))
	for i, d := range des {
		cd := canonicalizeCaPoolIssuancePolicyBaselineValuesKeyUsageExtendedKeyUsage(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewCaPoolIssuancePolicyBaselineValuesKeyUsageExtendedKeyUsage(c *Client, des, nw *CaPoolIssuancePolicyBaselineValuesKeyUsageExtendedKeyUsage) *CaPoolIssuancePolicyBaselineValuesKeyUsageExtendedKeyUsage {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for CaPoolIssuancePolicyBaselineValuesKeyUsageExtendedKeyUsage while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.BoolCanonicalize(des.ServerAuth, nw.ServerAuth) {
		nw.ServerAuth = des.ServerAuth
	}
	if dcl.BoolCanonicalize(des.ClientAuth, nw.ClientAuth) {
		nw.ClientAuth = des.ClientAuth
	}
	if dcl.BoolCanonicalize(des.CodeSigning, nw.CodeSigning) {
		nw.CodeSigning = des.CodeSigning
	}
	if dcl.BoolCanonicalize(des.EmailProtection, nw.EmailProtection) {
		nw.EmailProtection = des.EmailProtection
	}
	if dcl.BoolCanonicalize(des.TimeStamping, nw.TimeStamping) {
		nw.TimeStamping = des.TimeStamping
	}
	if dcl.BoolCanonicalize(des.OcspSigning, nw.OcspSigning) {
		nw.OcspSigning = des.OcspSigning
	}

	return nw
}

func canonicalizeNewCaPoolIssuancePolicyBaselineValuesKeyUsageExtendedKeyUsageSet(c *Client, des, nw []CaPoolIssuancePolicyBaselineValuesKeyUsageExtendedKeyUsage) []CaPoolIssuancePolicyBaselineValuesKeyUsageExtendedKeyUsage {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []CaPoolIssuancePolicyBaselineValuesKeyUsageExtendedKeyUsage
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareCaPoolIssuancePolicyBaselineValuesKeyUsageExtendedKeyUsageNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewCaPoolIssuancePolicyBaselineValuesKeyUsageExtendedKeyUsage(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewCaPoolIssuancePolicyBaselineValuesKeyUsageExtendedKeyUsageSlice(c *Client, des, nw []CaPoolIssuancePolicyBaselineValuesKeyUsageExtendedKeyUsage) []CaPoolIssuancePolicyBaselineValuesKeyUsageExtendedKeyUsage {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []CaPoolIssuancePolicyBaselineValuesKeyUsageExtendedKeyUsage
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewCaPoolIssuancePolicyBaselineValuesKeyUsageExtendedKeyUsage(c, &d, &n))
	}

	return items
}

func canonicalizeCaPoolIssuancePolicyBaselineValuesKeyUsageUnknownExtendedKeyUsages(des, initial *CaPoolIssuancePolicyBaselineValuesKeyUsageUnknownExtendedKeyUsages, opts ...dcl.ApplyOption) *CaPoolIssuancePolicyBaselineValuesKeyUsageUnknownExtendedKeyUsages {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &CaPoolIssuancePolicyBaselineValuesKeyUsageUnknownExtendedKeyUsages{}

	if dcl.IsZeroValue(des.ObjectIdPath) || (dcl.IsEmptyValueIndirect(des.ObjectIdPath) && dcl.IsEmptyValueIndirect(initial.ObjectIdPath)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.ObjectIdPath = initial.ObjectIdPath
	} else {
		cDes.ObjectIdPath = des.ObjectIdPath
	}

	return cDes
}

func canonicalizeCaPoolIssuancePolicyBaselineValuesKeyUsageUnknownExtendedKeyUsagesSlice(des, initial []CaPoolIssuancePolicyBaselineValuesKeyUsageUnknownExtendedKeyUsages, opts ...dcl.ApplyOption) []CaPoolIssuancePolicyBaselineValuesKeyUsageUnknownExtendedKeyUsages {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]CaPoolIssuancePolicyBaselineValuesKeyUsageUnknownExtendedKeyUsages, 0, len(des))
		for _, d := range des {
			cd := canonicalizeCaPoolIssuancePolicyBaselineValuesKeyUsageUnknownExtendedKeyUsages(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]CaPoolIssuancePolicyBaselineValuesKeyUsageUnknownExtendedKeyUsages, 0, len(des))
	for i, d := range des {
		cd := canonicalizeCaPoolIssuancePolicyBaselineValuesKeyUsageUnknownExtendedKeyUsages(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewCaPoolIssuancePolicyBaselineValuesKeyUsageUnknownExtendedKeyUsages(c *Client, des, nw *CaPoolIssuancePolicyBaselineValuesKeyUsageUnknownExtendedKeyUsages) *CaPoolIssuancePolicyBaselineValuesKeyUsageUnknownExtendedKeyUsages {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for CaPoolIssuancePolicyBaselineValuesKeyUsageUnknownExtendedKeyUsages while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	return nw
}

func canonicalizeNewCaPoolIssuancePolicyBaselineValuesKeyUsageUnknownExtendedKeyUsagesSet(c *Client, des, nw []CaPoolIssuancePolicyBaselineValuesKeyUsageUnknownExtendedKeyUsages) []CaPoolIssuancePolicyBaselineValuesKeyUsageUnknownExtendedKeyUsages {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []CaPoolIssuancePolicyBaselineValuesKeyUsageUnknownExtendedKeyUsages
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareCaPoolIssuancePolicyBaselineValuesKeyUsageUnknownExtendedKeyUsagesNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewCaPoolIssuancePolicyBaselineValuesKeyUsageUnknownExtendedKeyUsages(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewCaPoolIssuancePolicyBaselineValuesKeyUsageUnknownExtendedKeyUsagesSlice(c *Client, des, nw []CaPoolIssuancePolicyBaselineValuesKeyUsageUnknownExtendedKeyUsages) []CaPoolIssuancePolicyBaselineValuesKeyUsageUnknownExtendedKeyUsages {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []CaPoolIssuancePolicyBaselineValuesKeyUsageUnknownExtendedKeyUsages
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewCaPoolIssuancePolicyBaselineValuesKeyUsageUnknownExtendedKeyUsages(c, &d, &n))
	}

	return items
}

func canonicalizeCaPoolIssuancePolicyBaselineValuesCaOptions(des, initial *CaPoolIssuancePolicyBaselineValuesCaOptions, opts ...dcl.ApplyOption) *CaPoolIssuancePolicyBaselineValuesCaOptions {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &CaPoolIssuancePolicyBaselineValuesCaOptions{}

	if dcl.BoolCanonicalize(des.IsCa, initial.IsCa) || dcl.IsZeroValue(des.IsCa) {
		cDes.IsCa = initial.IsCa
	} else {
		cDes.IsCa = des.IsCa
	}
	if dcl.IsZeroValue(des.MaxIssuerPathLength) || (dcl.IsEmptyValueIndirect(des.MaxIssuerPathLength) && dcl.IsEmptyValueIndirect(initial.MaxIssuerPathLength)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.MaxIssuerPathLength = initial.MaxIssuerPathLength
	} else {
		cDes.MaxIssuerPathLength = des.MaxIssuerPathLength
	}
	if dcl.BoolCanonicalize(des.ZeroMaxIssuerPathLength, initial.ZeroMaxIssuerPathLength) || dcl.IsZeroValue(des.ZeroMaxIssuerPathLength) {
		cDes.ZeroMaxIssuerPathLength = initial.ZeroMaxIssuerPathLength
	} else {
		cDes.ZeroMaxIssuerPathLength = des.ZeroMaxIssuerPathLength
	}

	return cDes
}

func canonicalizeCaPoolIssuancePolicyBaselineValuesCaOptionsSlice(des, initial []CaPoolIssuancePolicyBaselineValuesCaOptions, opts ...dcl.ApplyOption) []CaPoolIssuancePolicyBaselineValuesCaOptions {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]CaPoolIssuancePolicyBaselineValuesCaOptions, 0, len(des))
		for _, d := range des {
			cd := canonicalizeCaPoolIssuancePolicyBaselineValuesCaOptions(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]CaPoolIssuancePolicyBaselineValuesCaOptions, 0, len(des))
	for i, d := range des {
		cd := canonicalizeCaPoolIssuancePolicyBaselineValuesCaOptions(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewCaPoolIssuancePolicyBaselineValuesCaOptions(c *Client, des, nw *CaPoolIssuancePolicyBaselineValuesCaOptions) *CaPoolIssuancePolicyBaselineValuesCaOptions {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for CaPoolIssuancePolicyBaselineValuesCaOptions while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.BoolCanonicalize(des.IsCa, nw.IsCa) {
		nw.IsCa = des.IsCa
	}
	if dcl.BoolCanonicalize(des.ZeroMaxIssuerPathLength, nw.ZeroMaxIssuerPathLength) {
		nw.ZeroMaxIssuerPathLength = des.ZeroMaxIssuerPathLength
	}

	return nw
}

func canonicalizeNewCaPoolIssuancePolicyBaselineValuesCaOptionsSet(c *Client, des, nw []CaPoolIssuancePolicyBaselineValuesCaOptions) []CaPoolIssuancePolicyBaselineValuesCaOptions {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []CaPoolIssuancePolicyBaselineValuesCaOptions
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareCaPoolIssuancePolicyBaselineValuesCaOptionsNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewCaPoolIssuancePolicyBaselineValuesCaOptions(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewCaPoolIssuancePolicyBaselineValuesCaOptionsSlice(c *Client, des, nw []CaPoolIssuancePolicyBaselineValuesCaOptions) []CaPoolIssuancePolicyBaselineValuesCaOptions {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []CaPoolIssuancePolicyBaselineValuesCaOptions
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewCaPoolIssuancePolicyBaselineValuesCaOptions(c, &d, &n))
	}

	return items
}

func canonicalizeCaPoolIssuancePolicyBaselineValuesPolicyIds(des, initial *CaPoolIssuancePolicyBaselineValuesPolicyIds, opts ...dcl.ApplyOption) *CaPoolIssuancePolicyBaselineValuesPolicyIds {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &CaPoolIssuancePolicyBaselineValuesPolicyIds{}

	if dcl.IsZeroValue(des.ObjectIdPath) || (dcl.IsEmptyValueIndirect(des.ObjectIdPath) && dcl.IsEmptyValueIndirect(initial.ObjectIdPath)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.ObjectIdPath = initial.ObjectIdPath
	} else {
		cDes.ObjectIdPath = des.ObjectIdPath
	}

	return cDes
}

func canonicalizeCaPoolIssuancePolicyBaselineValuesPolicyIdsSlice(des, initial []CaPoolIssuancePolicyBaselineValuesPolicyIds, opts ...dcl.ApplyOption) []CaPoolIssuancePolicyBaselineValuesPolicyIds {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]CaPoolIssuancePolicyBaselineValuesPolicyIds, 0, len(des))
		for _, d := range des {
			cd := canonicalizeCaPoolIssuancePolicyBaselineValuesPolicyIds(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]CaPoolIssuancePolicyBaselineValuesPolicyIds, 0, len(des))
	for i, d := range des {
		cd := canonicalizeCaPoolIssuancePolicyBaselineValuesPolicyIds(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewCaPoolIssuancePolicyBaselineValuesPolicyIds(c *Client, des, nw *CaPoolIssuancePolicyBaselineValuesPolicyIds) *CaPoolIssuancePolicyBaselineValuesPolicyIds {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for CaPoolIssuancePolicyBaselineValuesPolicyIds while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	return nw
}

func canonicalizeNewCaPoolIssuancePolicyBaselineValuesPolicyIdsSet(c *Client, des, nw []CaPoolIssuancePolicyBaselineValuesPolicyIds) []CaPoolIssuancePolicyBaselineValuesPolicyIds {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []CaPoolIssuancePolicyBaselineValuesPolicyIds
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareCaPoolIssuancePolicyBaselineValuesPolicyIdsNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewCaPoolIssuancePolicyBaselineValuesPolicyIds(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewCaPoolIssuancePolicyBaselineValuesPolicyIdsSlice(c *Client, des, nw []CaPoolIssuancePolicyBaselineValuesPolicyIds) []CaPoolIssuancePolicyBaselineValuesPolicyIds {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []CaPoolIssuancePolicyBaselineValuesPolicyIds
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewCaPoolIssuancePolicyBaselineValuesPolicyIds(c, &d, &n))
	}

	return items
}

func canonicalizeCaPoolIssuancePolicyBaselineValuesAdditionalExtensions(des, initial *CaPoolIssuancePolicyBaselineValuesAdditionalExtensions, opts ...dcl.ApplyOption) *CaPoolIssuancePolicyBaselineValuesAdditionalExtensions {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &CaPoolIssuancePolicyBaselineValuesAdditionalExtensions{}

	cDes.ObjectId = canonicalizeCaPoolIssuancePolicyBaselineValuesAdditionalExtensionsObjectId(des.ObjectId, initial.ObjectId, opts...)
	if dcl.BoolCanonicalize(des.Critical, initial.Critical) || dcl.IsZeroValue(des.Critical) {
		cDes.Critical = initial.Critical
	} else {
		cDes.Critical = des.Critical
	}
	if dcl.StringCanonicalize(des.Value, initial.Value) || dcl.IsZeroValue(des.Value) {
		cDes.Value = initial.Value
	} else {
		cDes.Value = des.Value
	}

	return cDes
}

func canonicalizeCaPoolIssuancePolicyBaselineValuesAdditionalExtensionsSlice(des, initial []CaPoolIssuancePolicyBaselineValuesAdditionalExtensions, opts ...dcl.ApplyOption) []CaPoolIssuancePolicyBaselineValuesAdditionalExtensions {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]CaPoolIssuancePolicyBaselineValuesAdditionalExtensions, 0, len(des))
		for _, d := range des {
			cd := canonicalizeCaPoolIssuancePolicyBaselineValuesAdditionalExtensions(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]CaPoolIssuancePolicyBaselineValuesAdditionalExtensions, 0, len(des))
	for i, d := range des {
		cd := canonicalizeCaPoolIssuancePolicyBaselineValuesAdditionalExtensions(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewCaPoolIssuancePolicyBaselineValuesAdditionalExtensions(c *Client, des, nw *CaPoolIssuancePolicyBaselineValuesAdditionalExtensions) *CaPoolIssuancePolicyBaselineValuesAdditionalExtensions {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for CaPoolIssuancePolicyBaselineValuesAdditionalExtensions while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	nw.ObjectId = canonicalizeNewCaPoolIssuancePolicyBaselineValuesAdditionalExtensionsObjectId(c, des.ObjectId, nw.ObjectId)
	if dcl.BoolCanonicalize(des.Critical, nw.Critical) {
		nw.Critical = des.Critical
	}
	if dcl.StringCanonicalize(des.Value, nw.Value) {
		nw.Value = des.Value
	}

	return nw
}

func canonicalizeNewCaPoolIssuancePolicyBaselineValuesAdditionalExtensionsSet(c *Client, des, nw []CaPoolIssuancePolicyBaselineValuesAdditionalExtensions) []CaPoolIssuancePolicyBaselineValuesAdditionalExtensions {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []CaPoolIssuancePolicyBaselineValuesAdditionalExtensions
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareCaPoolIssuancePolicyBaselineValuesAdditionalExtensionsNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewCaPoolIssuancePolicyBaselineValuesAdditionalExtensions(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewCaPoolIssuancePolicyBaselineValuesAdditionalExtensionsSlice(c *Client, des, nw []CaPoolIssuancePolicyBaselineValuesAdditionalExtensions) []CaPoolIssuancePolicyBaselineValuesAdditionalExtensions {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []CaPoolIssuancePolicyBaselineValuesAdditionalExtensions
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewCaPoolIssuancePolicyBaselineValuesAdditionalExtensions(c, &d, &n))
	}

	return items
}

func canonicalizeCaPoolIssuancePolicyBaselineValuesAdditionalExtensionsObjectId(des, initial *CaPoolIssuancePolicyBaselineValuesAdditionalExtensionsObjectId, opts ...dcl.ApplyOption) *CaPoolIssuancePolicyBaselineValuesAdditionalExtensionsObjectId {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &CaPoolIssuancePolicyBaselineValuesAdditionalExtensionsObjectId{}

	if dcl.IsZeroValue(des.ObjectIdPath) || (dcl.IsEmptyValueIndirect(des.ObjectIdPath) && dcl.IsEmptyValueIndirect(initial.ObjectIdPath)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.ObjectIdPath = initial.ObjectIdPath
	} else {
		cDes.ObjectIdPath = des.ObjectIdPath
	}

	return cDes
}

func canonicalizeCaPoolIssuancePolicyBaselineValuesAdditionalExtensionsObjectIdSlice(des, initial []CaPoolIssuancePolicyBaselineValuesAdditionalExtensionsObjectId, opts ...dcl.ApplyOption) []CaPoolIssuancePolicyBaselineValuesAdditionalExtensionsObjectId {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]CaPoolIssuancePolicyBaselineValuesAdditionalExtensionsObjectId, 0, len(des))
		for _, d := range des {
			cd := canonicalizeCaPoolIssuancePolicyBaselineValuesAdditionalExtensionsObjectId(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]CaPoolIssuancePolicyBaselineValuesAdditionalExtensionsObjectId, 0, len(des))
	for i, d := range des {
		cd := canonicalizeCaPoolIssuancePolicyBaselineValuesAdditionalExtensionsObjectId(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewCaPoolIssuancePolicyBaselineValuesAdditionalExtensionsObjectId(c *Client, des, nw *CaPoolIssuancePolicyBaselineValuesAdditionalExtensionsObjectId) *CaPoolIssuancePolicyBaselineValuesAdditionalExtensionsObjectId {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for CaPoolIssuancePolicyBaselineValuesAdditionalExtensionsObjectId while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	return nw
}

func canonicalizeNewCaPoolIssuancePolicyBaselineValuesAdditionalExtensionsObjectIdSet(c *Client, des, nw []CaPoolIssuancePolicyBaselineValuesAdditionalExtensionsObjectId) []CaPoolIssuancePolicyBaselineValuesAdditionalExtensionsObjectId {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []CaPoolIssuancePolicyBaselineValuesAdditionalExtensionsObjectId
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareCaPoolIssuancePolicyBaselineValuesAdditionalExtensionsObjectIdNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewCaPoolIssuancePolicyBaselineValuesAdditionalExtensionsObjectId(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewCaPoolIssuancePolicyBaselineValuesAdditionalExtensionsObjectIdSlice(c *Client, des, nw []CaPoolIssuancePolicyBaselineValuesAdditionalExtensionsObjectId) []CaPoolIssuancePolicyBaselineValuesAdditionalExtensionsObjectId {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []CaPoolIssuancePolicyBaselineValuesAdditionalExtensionsObjectId
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewCaPoolIssuancePolicyBaselineValuesAdditionalExtensionsObjectId(c, &d, &n))
	}

	return items
}

func canonicalizeCaPoolIssuancePolicyIdentityConstraints(des, initial *CaPoolIssuancePolicyIdentityConstraints, opts ...dcl.ApplyOption) *CaPoolIssuancePolicyIdentityConstraints {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &CaPoolIssuancePolicyIdentityConstraints{}

	cDes.CelExpression = canonicalizeCaPoolIssuancePolicyIdentityConstraintsCelExpression(des.CelExpression, initial.CelExpression, opts...)
	if dcl.BoolCanonicalize(des.AllowSubjectPassthrough, initial.AllowSubjectPassthrough) || dcl.IsZeroValue(des.AllowSubjectPassthrough) {
		cDes.AllowSubjectPassthrough = initial.AllowSubjectPassthrough
	} else {
		cDes.AllowSubjectPassthrough = des.AllowSubjectPassthrough
	}
	if dcl.BoolCanonicalize(des.AllowSubjectAltNamesPassthrough, initial.AllowSubjectAltNamesPassthrough) || dcl.IsZeroValue(des.AllowSubjectAltNamesPassthrough) {
		cDes.AllowSubjectAltNamesPassthrough = initial.AllowSubjectAltNamesPassthrough
	} else {
		cDes.AllowSubjectAltNamesPassthrough = des.AllowSubjectAltNamesPassthrough
	}

	return cDes
}

func canonicalizeCaPoolIssuancePolicyIdentityConstraintsSlice(des, initial []CaPoolIssuancePolicyIdentityConstraints, opts ...dcl.ApplyOption) []CaPoolIssuancePolicyIdentityConstraints {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]CaPoolIssuancePolicyIdentityConstraints, 0, len(des))
		for _, d := range des {
			cd := canonicalizeCaPoolIssuancePolicyIdentityConstraints(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]CaPoolIssuancePolicyIdentityConstraints, 0, len(des))
	for i, d := range des {
		cd := canonicalizeCaPoolIssuancePolicyIdentityConstraints(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewCaPoolIssuancePolicyIdentityConstraints(c *Client, des, nw *CaPoolIssuancePolicyIdentityConstraints) *CaPoolIssuancePolicyIdentityConstraints {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for CaPoolIssuancePolicyIdentityConstraints while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	nw.CelExpression = canonicalizeNewCaPoolIssuancePolicyIdentityConstraintsCelExpression(c, des.CelExpression, nw.CelExpression)
	if dcl.BoolCanonicalize(des.AllowSubjectPassthrough, nw.AllowSubjectPassthrough) {
		nw.AllowSubjectPassthrough = des.AllowSubjectPassthrough
	}
	if dcl.BoolCanonicalize(des.AllowSubjectAltNamesPassthrough, nw.AllowSubjectAltNamesPassthrough) {
		nw.AllowSubjectAltNamesPassthrough = des.AllowSubjectAltNamesPassthrough
	}

	return nw
}

func canonicalizeNewCaPoolIssuancePolicyIdentityConstraintsSet(c *Client, des, nw []CaPoolIssuancePolicyIdentityConstraints) []CaPoolIssuancePolicyIdentityConstraints {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []CaPoolIssuancePolicyIdentityConstraints
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareCaPoolIssuancePolicyIdentityConstraintsNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewCaPoolIssuancePolicyIdentityConstraints(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewCaPoolIssuancePolicyIdentityConstraintsSlice(c *Client, des, nw []CaPoolIssuancePolicyIdentityConstraints) []CaPoolIssuancePolicyIdentityConstraints {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []CaPoolIssuancePolicyIdentityConstraints
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewCaPoolIssuancePolicyIdentityConstraints(c, &d, &n))
	}

	return items
}

func canonicalizeCaPoolIssuancePolicyIdentityConstraintsCelExpression(des, initial *CaPoolIssuancePolicyIdentityConstraintsCelExpression, opts ...dcl.ApplyOption) *CaPoolIssuancePolicyIdentityConstraintsCelExpression {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &CaPoolIssuancePolicyIdentityConstraintsCelExpression{}

	if dcl.StringCanonicalize(des.Expression, initial.Expression) || dcl.IsZeroValue(des.Expression) {
		cDes.Expression = initial.Expression
	} else {
		cDes.Expression = des.Expression
	}
	if dcl.StringCanonicalize(des.Title, initial.Title) || dcl.IsZeroValue(des.Title) {
		cDes.Title = initial.Title
	} else {
		cDes.Title = des.Title
	}
	if dcl.StringCanonicalize(des.Description, initial.Description) || dcl.IsZeroValue(des.Description) {
		cDes.Description = initial.Description
	} else {
		cDes.Description = des.Description
	}
	if dcl.StringCanonicalize(des.Location, initial.Location) || dcl.IsZeroValue(des.Location) {
		cDes.Location = initial.Location
	} else {
		cDes.Location = des.Location
	}

	return cDes
}

func canonicalizeCaPoolIssuancePolicyIdentityConstraintsCelExpressionSlice(des, initial []CaPoolIssuancePolicyIdentityConstraintsCelExpression, opts ...dcl.ApplyOption) []CaPoolIssuancePolicyIdentityConstraintsCelExpression {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]CaPoolIssuancePolicyIdentityConstraintsCelExpression, 0, len(des))
		for _, d := range des {
			cd := canonicalizeCaPoolIssuancePolicyIdentityConstraintsCelExpression(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]CaPoolIssuancePolicyIdentityConstraintsCelExpression, 0, len(des))
	for i, d := range des {
		cd := canonicalizeCaPoolIssuancePolicyIdentityConstraintsCelExpression(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewCaPoolIssuancePolicyIdentityConstraintsCelExpression(c *Client, des, nw *CaPoolIssuancePolicyIdentityConstraintsCelExpression) *CaPoolIssuancePolicyIdentityConstraintsCelExpression {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for CaPoolIssuancePolicyIdentityConstraintsCelExpression while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringCanonicalize(des.Expression, nw.Expression) {
		nw.Expression = des.Expression
	}
	if dcl.StringCanonicalize(des.Title, nw.Title) {
		nw.Title = des.Title
	}
	if dcl.StringCanonicalize(des.Description, nw.Description) {
		nw.Description = des.Description
	}
	if dcl.StringCanonicalize(des.Location, nw.Location) {
		nw.Location = des.Location
	}

	return nw
}

func canonicalizeNewCaPoolIssuancePolicyIdentityConstraintsCelExpressionSet(c *Client, des, nw []CaPoolIssuancePolicyIdentityConstraintsCelExpression) []CaPoolIssuancePolicyIdentityConstraintsCelExpression {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []CaPoolIssuancePolicyIdentityConstraintsCelExpression
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareCaPoolIssuancePolicyIdentityConstraintsCelExpressionNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewCaPoolIssuancePolicyIdentityConstraintsCelExpression(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewCaPoolIssuancePolicyIdentityConstraintsCelExpressionSlice(c *Client, des, nw []CaPoolIssuancePolicyIdentityConstraintsCelExpression) []CaPoolIssuancePolicyIdentityConstraintsCelExpression {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []CaPoolIssuancePolicyIdentityConstraintsCelExpression
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewCaPoolIssuancePolicyIdentityConstraintsCelExpression(c, &d, &n))
	}

	return items
}

func canonicalizeCaPoolIssuancePolicyPassthroughExtensions(des, initial *CaPoolIssuancePolicyPassthroughExtensions, opts ...dcl.ApplyOption) *CaPoolIssuancePolicyPassthroughExtensions {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &CaPoolIssuancePolicyPassthroughExtensions{}

	if dcl.IsZeroValue(des.KnownExtensions) || (dcl.IsEmptyValueIndirect(des.KnownExtensions) && dcl.IsEmptyValueIndirect(initial.KnownExtensions)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.KnownExtensions = initial.KnownExtensions
	} else {
		cDes.KnownExtensions = des.KnownExtensions
	}
	cDes.AdditionalExtensions = canonicalizeCaPoolIssuancePolicyPassthroughExtensionsAdditionalExtensionsSlice(des.AdditionalExtensions, initial.AdditionalExtensions, opts...)

	return cDes
}

func canonicalizeCaPoolIssuancePolicyPassthroughExtensionsSlice(des, initial []CaPoolIssuancePolicyPassthroughExtensions, opts ...dcl.ApplyOption) []CaPoolIssuancePolicyPassthroughExtensions {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]CaPoolIssuancePolicyPassthroughExtensions, 0, len(des))
		for _, d := range des {
			cd := canonicalizeCaPoolIssuancePolicyPassthroughExtensions(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]CaPoolIssuancePolicyPassthroughExtensions, 0, len(des))
	for i, d := range des {
		cd := canonicalizeCaPoolIssuancePolicyPassthroughExtensions(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewCaPoolIssuancePolicyPassthroughExtensions(c *Client, des, nw *CaPoolIssuancePolicyPassthroughExtensions) *CaPoolIssuancePolicyPassthroughExtensions {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for CaPoolIssuancePolicyPassthroughExtensions while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	nw.AdditionalExtensions = canonicalizeNewCaPoolIssuancePolicyPassthroughExtensionsAdditionalExtensionsSlice(c, des.AdditionalExtensions, nw.AdditionalExtensions)

	return nw
}

func canonicalizeNewCaPoolIssuancePolicyPassthroughExtensionsSet(c *Client, des, nw []CaPoolIssuancePolicyPassthroughExtensions) []CaPoolIssuancePolicyPassthroughExtensions {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []CaPoolIssuancePolicyPassthroughExtensions
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareCaPoolIssuancePolicyPassthroughExtensionsNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewCaPoolIssuancePolicyPassthroughExtensions(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewCaPoolIssuancePolicyPassthroughExtensionsSlice(c *Client, des, nw []CaPoolIssuancePolicyPassthroughExtensions) []CaPoolIssuancePolicyPassthroughExtensions {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []CaPoolIssuancePolicyPassthroughExtensions
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewCaPoolIssuancePolicyPassthroughExtensions(c, &d, &n))
	}

	return items
}

func canonicalizeCaPoolIssuancePolicyPassthroughExtensionsAdditionalExtensions(des, initial *CaPoolIssuancePolicyPassthroughExtensionsAdditionalExtensions, opts ...dcl.ApplyOption) *CaPoolIssuancePolicyPassthroughExtensionsAdditionalExtensions {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &CaPoolIssuancePolicyPassthroughExtensionsAdditionalExtensions{}

	if dcl.IsZeroValue(des.ObjectIdPath) || (dcl.IsEmptyValueIndirect(des.ObjectIdPath) && dcl.IsEmptyValueIndirect(initial.ObjectIdPath)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.ObjectIdPath = initial.ObjectIdPath
	} else {
		cDes.ObjectIdPath = des.ObjectIdPath
	}

	return cDes
}

func canonicalizeCaPoolIssuancePolicyPassthroughExtensionsAdditionalExtensionsSlice(des, initial []CaPoolIssuancePolicyPassthroughExtensionsAdditionalExtensions, opts ...dcl.ApplyOption) []CaPoolIssuancePolicyPassthroughExtensionsAdditionalExtensions {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]CaPoolIssuancePolicyPassthroughExtensionsAdditionalExtensions, 0, len(des))
		for _, d := range des {
			cd := canonicalizeCaPoolIssuancePolicyPassthroughExtensionsAdditionalExtensions(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]CaPoolIssuancePolicyPassthroughExtensionsAdditionalExtensions, 0, len(des))
	for i, d := range des {
		cd := canonicalizeCaPoolIssuancePolicyPassthroughExtensionsAdditionalExtensions(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewCaPoolIssuancePolicyPassthroughExtensionsAdditionalExtensions(c *Client, des, nw *CaPoolIssuancePolicyPassthroughExtensionsAdditionalExtensions) *CaPoolIssuancePolicyPassthroughExtensionsAdditionalExtensions {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for CaPoolIssuancePolicyPassthroughExtensionsAdditionalExtensions while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	return nw
}

func canonicalizeNewCaPoolIssuancePolicyPassthroughExtensionsAdditionalExtensionsSet(c *Client, des, nw []CaPoolIssuancePolicyPassthroughExtensionsAdditionalExtensions) []CaPoolIssuancePolicyPassthroughExtensionsAdditionalExtensions {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []CaPoolIssuancePolicyPassthroughExtensionsAdditionalExtensions
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareCaPoolIssuancePolicyPassthroughExtensionsAdditionalExtensionsNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewCaPoolIssuancePolicyPassthroughExtensionsAdditionalExtensions(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewCaPoolIssuancePolicyPassthroughExtensionsAdditionalExtensionsSlice(c *Client, des, nw []CaPoolIssuancePolicyPassthroughExtensionsAdditionalExtensions) []CaPoolIssuancePolicyPassthroughExtensionsAdditionalExtensions {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []CaPoolIssuancePolicyPassthroughExtensionsAdditionalExtensions
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewCaPoolIssuancePolicyPassthroughExtensionsAdditionalExtensions(c, &d, &n))
	}

	return items
}

func canonicalizeCaPoolPublishingOptions(des, initial *CaPoolPublishingOptions, opts ...dcl.ApplyOption) *CaPoolPublishingOptions {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &CaPoolPublishingOptions{}

	if dcl.BoolCanonicalize(des.PublishCaCert, initial.PublishCaCert) || dcl.IsZeroValue(des.PublishCaCert) {
		cDes.PublishCaCert = initial.PublishCaCert
	} else {
		cDes.PublishCaCert = des.PublishCaCert
	}
	if dcl.BoolCanonicalize(des.PublishCrl, initial.PublishCrl) || dcl.IsZeroValue(des.PublishCrl) {
		cDes.PublishCrl = initial.PublishCrl
	} else {
		cDes.PublishCrl = des.PublishCrl
	}

	return cDes
}

func canonicalizeCaPoolPublishingOptionsSlice(des, initial []CaPoolPublishingOptions, opts ...dcl.ApplyOption) []CaPoolPublishingOptions {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]CaPoolPublishingOptions, 0, len(des))
		for _, d := range des {
			cd := canonicalizeCaPoolPublishingOptions(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]CaPoolPublishingOptions, 0, len(des))
	for i, d := range des {
		cd := canonicalizeCaPoolPublishingOptions(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewCaPoolPublishingOptions(c *Client, des, nw *CaPoolPublishingOptions) *CaPoolPublishingOptions {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for CaPoolPublishingOptions while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.BoolCanonicalize(des.PublishCaCert, nw.PublishCaCert) {
		nw.PublishCaCert = des.PublishCaCert
	}
	if dcl.BoolCanonicalize(des.PublishCrl, nw.PublishCrl) {
		nw.PublishCrl = des.PublishCrl
	}

	return nw
}

func canonicalizeNewCaPoolPublishingOptionsSet(c *Client, des, nw []CaPoolPublishingOptions) []CaPoolPublishingOptions {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []CaPoolPublishingOptions
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareCaPoolPublishingOptionsNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewCaPoolPublishingOptions(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewCaPoolPublishingOptionsSlice(c *Client, des, nw []CaPoolPublishingOptions) []CaPoolPublishingOptions {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []CaPoolPublishingOptions
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewCaPoolPublishingOptions(c, &d, &n))
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
func diffCaPool(c *Client, desired, actual *CaPool, opts ...dcl.ApplyOption) ([]*dcl.FieldDiff, error) {
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

	if ds, err := dcl.Diff(desired.Tier, actual.Tier, dcl.DiffInfo{Type: "EnumType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Tier")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.IssuancePolicy, actual.IssuancePolicy, dcl.DiffInfo{ObjectFunction: compareCaPoolIssuancePolicyNewStyle, EmptyObject: EmptyCaPoolIssuancePolicy, OperationSelector: dcl.TriggersOperation("updateCaPoolUpdateCaPoolOperation")}, fn.AddNest("IssuancePolicy")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.PublishingOptions, actual.PublishingOptions, dcl.DiffInfo{ObjectFunction: compareCaPoolPublishingOptionsNewStyle, EmptyObject: EmptyCaPoolPublishingOptions, OperationSelector: dcl.TriggersOperation("updateCaPoolUpdateCaPoolOperation")}, fn.AddNest("PublishingOptions")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Labels, actual.Labels, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateCaPoolUpdateCaPoolOperation")}, fn.AddNest("Labels")); len(ds) != 0 || err != nil {
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
func compareCaPoolIssuancePolicyNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*CaPoolIssuancePolicy)
	if !ok {
		desiredNotPointer, ok := d.(CaPoolIssuancePolicy)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CaPoolIssuancePolicy or *CaPoolIssuancePolicy", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*CaPoolIssuancePolicy)
	if !ok {
		actualNotPointer, ok := a.(CaPoolIssuancePolicy)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CaPoolIssuancePolicy", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.AllowedKeyTypes, actual.AllowedKeyTypes, dcl.DiffInfo{ObjectFunction: compareCaPoolIssuancePolicyAllowedKeyTypesNewStyle, EmptyObject: EmptyCaPoolIssuancePolicyAllowedKeyTypes, OperationSelector: dcl.TriggersOperation("updateCaPoolUpdateCaPoolOperation")}, fn.AddNest("AllowedKeyTypes")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.MaximumLifetime, actual.MaximumLifetime, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateCaPoolUpdateCaPoolOperation")}, fn.AddNest("MaximumLifetime")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.AllowedIssuanceModes, actual.AllowedIssuanceModes, dcl.DiffInfo{ObjectFunction: compareCaPoolIssuancePolicyAllowedIssuanceModesNewStyle, EmptyObject: EmptyCaPoolIssuancePolicyAllowedIssuanceModes, OperationSelector: dcl.TriggersOperation("updateCaPoolUpdateCaPoolOperation")}, fn.AddNest("AllowedIssuanceModes")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.BaselineValues, actual.BaselineValues, dcl.DiffInfo{ObjectFunction: compareCaPoolIssuancePolicyBaselineValuesNewStyle, EmptyObject: EmptyCaPoolIssuancePolicyBaselineValues, OperationSelector: dcl.TriggersOperation("updateCaPoolUpdateCaPoolOperation")}, fn.AddNest("BaselineValues")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.IdentityConstraints, actual.IdentityConstraints, dcl.DiffInfo{ObjectFunction: compareCaPoolIssuancePolicyIdentityConstraintsNewStyle, EmptyObject: EmptyCaPoolIssuancePolicyIdentityConstraints, OperationSelector: dcl.TriggersOperation("updateCaPoolUpdateCaPoolOperation")}, fn.AddNest("IdentityConstraints")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.PassthroughExtensions, actual.PassthroughExtensions, dcl.DiffInfo{ObjectFunction: compareCaPoolIssuancePolicyPassthroughExtensionsNewStyle, EmptyObject: EmptyCaPoolIssuancePolicyPassthroughExtensions, OperationSelector: dcl.TriggersOperation("updateCaPoolUpdateCaPoolOperation")}, fn.AddNest("PassthroughExtensions")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareCaPoolIssuancePolicyAllowedKeyTypesNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*CaPoolIssuancePolicyAllowedKeyTypes)
	if !ok {
		desiredNotPointer, ok := d.(CaPoolIssuancePolicyAllowedKeyTypes)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CaPoolIssuancePolicyAllowedKeyTypes or *CaPoolIssuancePolicyAllowedKeyTypes", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*CaPoolIssuancePolicyAllowedKeyTypes)
	if !ok {
		actualNotPointer, ok := a.(CaPoolIssuancePolicyAllowedKeyTypes)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CaPoolIssuancePolicyAllowedKeyTypes", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.Rsa, actual.Rsa, dcl.DiffInfo{ObjectFunction: compareCaPoolIssuancePolicyAllowedKeyTypesRsaNewStyle, EmptyObject: EmptyCaPoolIssuancePolicyAllowedKeyTypesRsa, OperationSelector: dcl.TriggersOperation("updateCaPoolUpdateCaPoolOperation")}, fn.AddNest("Rsa")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.EllipticCurve, actual.EllipticCurve, dcl.DiffInfo{ObjectFunction: compareCaPoolIssuancePolicyAllowedKeyTypesEllipticCurveNewStyle, EmptyObject: EmptyCaPoolIssuancePolicyAllowedKeyTypesEllipticCurve, OperationSelector: dcl.TriggersOperation("updateCaPoolUpdateCaPoolOperation")}, fn.AddNest("EllipticCurve")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareCaPoolIssuancePolicyAllowedKeyTypesRsaNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*CaPoolIssuancePolicyAllowedKeyTypesRsa)
	if !ok {
		desiredNotPointer, ok := d.(CaPoolIssuancePolicyAllowedKeyTypesRsa)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CaPoolIssuancePolicyAllowedKeyTypesRsa or *CaPoolIssuancePolicyAllowedKeyTypesRsa", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*CaPoolIssuancePolicyAllowedKeyTypesRsa)
	if !ok {
		actualNotPointer, ok := a.(CaPoolIssuancePolicyAllowedKeyTypesRsa)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CaPoolIssuancePolicyAllowedKeyTypesRsa", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.MinModulusSize, actual.MinModulusSize, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateCaPoolUpdateCaPoolOperation")}, fn.AddNest("MinModulusSize")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.MaxModulusSize, actual.MaxModulusSize, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateCaPoolUpdateCaPoolOperation")}, fn.AddNest("MaxModulusSize")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareCaPoolIssuancePolicyAllowedKeyTypesEllipticCurveNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*CaPoolIssuancePolicyAllowedKeyTypesEllipticCurve)
	if !ok {
		desiredNotPointer, ok := d.(CaPoolIssuancePolicyAllowedKeyTypesEllipticCurve)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CaPoolIssuancePolicyAllowedKeyTypesEllipticCurve or *CaPoolIssuancePolicyAllowedKeyTypesEllipticCurve", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*CaPoolIssuancePolicyAllowedKeyTypesEllipticCurve)
	if !ok {
		actualNotPointer, ok := a.(CaPoolIssuancePolicyAllowedKeyTypesEllipticCurve)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CaPoolIssuancePolicyAllowedKeyTypesEllipticCurve", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.SignatureAlgorithm, actual.SignatureAlgorithm, dcl.DiffInfo{Type: "EnumType", OperationSelector: dcl.TriggersOperation("updateCaPoolUpdateCaPoolOperation")}, fn.AddNest("SignatureAlgorithm")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareCaPoolIssuancePolicyAllowedIssuanceModesNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*CaPoolIssuancePolicyAllowedIssuanceModes)
	if !ok {
		desiredNotPointer, ok := d.(CaPoolIssuancePolicyAllowedIssuanceModes)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CaPoolIssuancePolicyAllowedIssuanceModes or *CaPoolIssuancePolicyAllowedIssuanceModes", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*CaPoolIssuancePolicyAllowedIssuanceModes)
	if !ok {
		actualNotPointer, ok := a.(CaPoolIssuancePolicyAllowedIssuanceModes)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CaPoolIssuancePolicyAllowedIssuanceModes", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.AllowCsrBasedIssuance, actual.AllowCsrBasedIssuance, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateCaPoolUpdateCaPoolOperation")}, fn.AddNest("AllowCsrBasedIssuance")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.AllowConfigBasedIssuance, actual.AllowConfigBasedIssuance, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateCaPoolUpdateCaPoolOperation")}, fn.AddNest("AllowConfigBasedIssuance")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareCaPoolIssuancePolicyBaselineValuesNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*CaPoolIssuancePolicyBaselineValues)
	if !ok {
		desiredNotPointer, ok := d.(CaPoolIssuancePolicyBaselineValues)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CaPoolIssuancePolicyBaselineValues or *CaPoolIssuancePolicyBaselineValues", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*CaPoolIssuancePolicyBaselineValues)
	if !ok {
		actualNotPointer, ok := a.(CaPoolIssuancePolicyBaselineValues)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CaPoolIssuancePolicyBaselineValues", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.KeyUsage, actual.KeyUsage, dcl.DiffInfo{ObjectFunction: compareCaPoolIssuancePolicyBaselineValuesKeyUsageNewStyle, EmptyObject: EmptyCaPoolIssuancePolicyBaselineValuesKeyUsage, OperationSelector: dcl.TriggersOperation("updateCaPoolUpdateCaPoolOperation")}, fn.AddNest("KeyUsage")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.CaOptions, actual.CaOptions, dcl.DiffInfo{ObjectFunction: compareCaPoolIssuancePolicyBaselineValuesCaOptionsNewStyle, EmptyObject: EmptyCaPoolIssuancePolicyBaselineValuesCaOptions, OperationSelector: dcl.TriggersOperation("updateCaPoolUpdateCaPoolOperation")}, fn.AddNest("CaOptions")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.PolicyIds, actual.PolicyIds, dcl.DiffInfo{ObjectFunction: compareCaPoolIssuancePolicyBaselineValuesPolicyIdsNewStyle, EmptyObject: EmptyCaPoolIssuancePolicyBaselineValuesPolicyIds, OperationSelector: dcl.TriggersOperation("updateCaPoolUpdateCaPoolOperation")}, fn.AddNest("PolicyIds")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.AiaOcspServers, actual.AiaOcspServers, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateCaPoolUpdateCaPoolOperation")}, fn.AddNest("AiaOcspServers")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.AdditionalExtensions, actual.AdditionalExtensions, dcl.DiffInfo{ObjectFunction: compareCaPoolIssuancePolicyBaselineValuesAdditionalExtensionsNewStyle, EmptyObject: EmptyCaPoolIssuancePolicyBaselineValuesAdditionalExtensions, OperationSelector: dcl.TriggersOperation("updateCaPoolUpdateCaPoolOperation")}, fn.AddNest("AdditionalExtensions")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareCaPoolIssuancePolicyBaselineValuesKeyUsageNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*CaPoolIssuancePolicyBaselineValuesKeyUsage)
	if !ok {
		desiredNotPointer, ok := d.(CaPoolIssuancePolicyBaselineValuesKeyUsage)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CaPoolIssuancePolicyBaselineValuesKeyUsage or *CaPoolIssuancePolicyBaselineValuesKeyUsage", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*CaPoolIssuancePolicyBaselineValuesKeyUsage)
	if !ok {
		actualNotPointer, ok := a.(CaPoolIssuancePolicyBaselineValuesKeyUsage)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CaPoolIssuancePolicyBaselineValuesKeyUsage", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.BaseKeyUsage, actual.BaseKeyUsage, dcl.DiffInfo{ObjectFunction: compareCaPoolIssuancePolicyBaselineValuesKeyUsageBaseKeyUsageNewStyle, EmptyObject: EmptyCaPoolIssuancePolicyBaselineValuesKeyUsageBaseKeyUsage, OperationSelector: dcl.TriggersOperation("updateCaPoolUpdateCaPoolOperation")}, fn.AddNest("BaseKeyUsage")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.ExtendedKeyUsage, actual.ExtendedKeyUsage, dcl.DiffInfo{ObjectFunction: compareCaPoolIssuancePolicyBaselineValuesKeyUsageExtendedKeyUsageNewStyle, EmptyObject: EmptyCaPoolIssuancePolicyBaselineValuesKeyUsageExtendedKeyUsage, OperationSelector: dcl.TriggersOperation("updateCaPoolUpdateCaPoolOperation")}, fn.AddNest("ExtendedKeyUsage")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.UnknownExtendedKeyUsages, actual.UnknownExtendedKeyUsages, dcl.DiffInfo{ObjectFunction: compareCaPoolIssuancePolicyBaselineValuesKeyUsageUnknownExtendedKeyUsagesNewStyle, EmptyObject: EmptyCaPoolIssuancePolicyBaselineValuesKeyUsageUnknownExtendedKeyUsages, OperationSelector: dcl.TriggersOperation("updateCaPoolUpdateCaPoolOperation")}, fn.AddNest("UnknownExtendedKeyUsages")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareCaPoolIssuancePolicyBaselineValuesKeyUsageBaseKeyUsageNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*CaPoolIssuancePolicyBaselineValuesKeyUsageBaseKeyUsage)
	if !ok {
		desiredNotPointer, ok := d.(CaPoolIssuancePolicyBaselineValuesKeyUsageBaseKeyUsage)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CaPoolIssuancePolicyBaselineValuesKeyUsageBaseKeyUsage or *CaPoolIssuancePolicyBaselineValuesKeyUsageBaseKeyUsage", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*CaPoolIssuancePolicyBaselineValuesKeyUsageBaseKeyUsage)
	if !ok {
		actualNotPointer, ok := a.(CaPoolIssuancePolicyBaselineValuesKeyUsageBaseKeyUsage)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CaPoolIssuancePolicyBaselineValuesKeyUsageBaseKeyUsage", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.DigitalSignature, actual.DigitalSignature, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateCaPoolUpdateCaPoolOperation")}, fn.AddNest("DigitalSignature")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.ContentCommitment, actual.ContentCommitment, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateCaPoolUpdateCaPoolOperation")}, fn.AddNest("ContentCommitment")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.KeyEncipherment, actual.KeyEncipherment, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateCaPoolUpdateCaPoolOperation")}, fn.AddNest("KeyEncipherment")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.DataEncipherment, actual.DataEncipherment, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateCaPoolUpdateCaPoolOperation")}, fn.AddNest("DataEncipherment")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.KeyAgreement, actual.KeyAgreement, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateCaPoolUpdateCaPoolOperation")}, fn.AddNest("KeyAgreement")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.CertSign, actual.CertSign, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateCaPoolUpdateCaPoolOperation")}, fn.AddNest("CertSign")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.CrlSign, actual.CrlSign, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateCaPoolUpdateCaPoolOperation")}, fn.AddNest("CrlSign")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.EncipherOnly, actual.EncipherOnly, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateCaPoolUpdateCaPoolOperation")}, fn.AddNest("EncipherOnly")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.DecipherOnly, actual.DecipherOnly, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateCaPoolUpdateCaPoolOperation")}, fn.AddNest("DecipherOnly")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareCaPoolIssuancePolicyBaselineValuesKeyUsageExtendedKeyUsageNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*CaPoolIssuancePolicyBaselineValuesKeyUsageExtendedKeyUsage)
	if !ok {
		desiredNotPointer, ok := d.(CaPoolIssuancePolicyBaselineValuesKeyUsageExtendedKeyUsage)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CaPoolIssuancePolicyBaselineValuesKeyUsageExtendedKeyUsage or *CaPoolIssuancePolicyBaselineValuesKeyUsageExtendedKeyUsage", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*CaPoolIssuancePolicyBaselineValuesKeyUsageExtendedKeyUsage)
	if !ok {
		actualNotPointer, ok := a.(CaPoolIssuancePolicyBaselineValuesKeyUsageExtendedKeyUsage)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CaPoolIssuancePolicyBaselineValuesKeyUsageExtendedKeyUsage", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.ServerAuth, actual.ServerAuth, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateCaPoolUpdateCaPoolOperation")}, fn.AddNest("ServerAuth")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.ClientAuth, actual.ClientAuth, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateCaPoolUpdateCaPoolOperation")}, fn.AddNest("ClientAuth")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.CodeSigning, actual.CodeSigning, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateCaPoolUpdateCaPoolOperation")}, fn.AddNest("CodeSigning")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.EmailProtection, actual.EmailProtection, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateCaPoolUpdateCaPoolOperation")}, fn.AddNest("EmailProtection")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.TimeStamping, actual.TimeStamping, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateCaPoolUpdateCaPoolOperation")}, fn.AddNest("TimeStamping")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.OcspSigning, actual.OcspSigning, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateCaPoolUpdateCaPoolOperation")}, fn.AddNest("OcspSigning")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareCaPoolIssuancePolicyBaselineValuesKeyUsageUnknownExtendedKeyUsagesNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*CaPoolIssuancePolicyBaselineValuesKeyUsageUnknownExtendedKeyUsages)
	if !ok {
		desiredNotPointer, ok := d.(CaPoolIssuancePolicyBaselineValuesKeyUsageUnknownExtendedKeyUsages)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CaPoolIssuancePolicyBaselineValuesKeyUsageUnknownExtendedKeyUsages or *CaPoolIssuancePolicyBaselineValuesKeyUsageUnknownExtendedKeyUsages", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*CaPoolIssuancePolicyBaselineValuesKeyUsageUnknownExtendedKeyUsages)
	if !ok {
		actualNotPointer, ok := a.(CaPoolIssuancePolicyBaselineValuesKeyUsageUnknownExtendedKeyUsages)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CaPoolIssuancePolicyBaselineValuesKeyUsageUnknownExtendedKeyUsages", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.ObjectIdPath, actual.ObjectIdPath, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateCaPoolUpdateCaPoolOperation")}, fn.AddNest("ObjectIdPath")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareCaPoolIssuancePolicyBaselineValuesCaOptionsNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*CaPoolIssuancePolicyBaselineValuesCaOptions)
	if !ok {
		desiredNotPointer, ok := d.(CaPoolIssuancePolicyBaselineValuesCaOptions)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CaPoolIssuancePolicyBaselineValuesCaOptions or *CaPoolIssuancePolicyBaselineValuesCaOptions", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*CaPoolIssuancePolicyBaselineValuesCaOptions)
	if !ok {
		actualNotPointer, ok := a.(CaPoolIssuancePolicyBaselineValuesCaOptions)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CaPoolIssuancePolicyBaselineValuesCaOptions", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.IsCa, actual.IsCa, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateCaPoolUpdateCaPoolOperation")}, fn.AddNest("IsCa")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.MaxIssuerPathLength, actual.MaxIssuerPathLength, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateCaPoolUpdateCaPoolOperation")}, fn.AddNest("MaxIssuerPathLength")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.ZeroMaxIssuerPathLength, actual.ZeroMaxIssuerPathLength, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateCaPoolUpdateCaPoolOperation")}, fn.AddNest("ZeroMaxIssuerPathLength")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareCaPoolIssuancePolicyBaselineValuesPolicyIdsNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*CaPoolIssuancePolicyBaselineValuesPolicyIds)
	if !ok {
		desiredNotPointer, ok := d.(CaPoolIssuancePolicyBaselineValuesPolicyIds)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CaPoolIssuancePolicyBaselineValuesPolicyIds or *CaPoolIssuancePolicyBaselineValuesPolicyIds", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*CaPoolIssuancePolicyBaselineValuesPolicyIds)
	if !ok {
		actualNotPointer, ok := a.(CaPoolIssuancePolicyBaselineValuesPolicyIds)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CaPoolIssuancePolicyBaselineValuesPolicyIds", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.ObjectIdPath, actual.ObjectIdPath, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateCaPoolUpdateCaPoolOperation")}, fn.AddNest("ObjectIdPath")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareCaPoolIssuancePolicyBaselineValuesAdditionalExtensionsNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*CaPoolIssuancePolicyBaselineValuesAdditionalExtensions)
	if !ok {
		desiredNotPointer, ok := d.(CaPoolIssuancePolicyBaselineValuesAdditionalExtensions)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CaPoolIssuancePolicyBaselineValuesAdditionalExtensions or *CaPoolIssuancePolicyBaselineValuesAdditionalExtensions", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*CaPoolIssuancePolicyBaselineValuesAdditionalExtensions)
	if !ok {
		actualNotPointer, ok := a.(CaPoolIssuancePolicyBaselineValuesAdditionalExtensions)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CaPoolIssuancePolicyBaselineValuesAdditionalExtensions", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.ObjectId, actual.ObjectId, dcl.DiffInfo{ObjectFunction: compareCaPoolIssuancePolicyBaselineValuesAdditionalExtensionsObjectIdNewStyle, EmptyObject: EmptyCaPoolIssuancePolicyBaselineValuesAdditionalExtensionsObjectId, OperationSelector: dcl.TriggersOperation("updateCaPoolUpdateCaPoolOperation")}, fn.AddNest("ObjectId")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Critical, actual.Critical, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateCaPoolUpdateCaPoolOperation")}, fn.AddNest("Critical")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Value, actual.Value, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateCaPoolUpdateCaPoolOperation")}, fn.AddNest("Value")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareCaPoolIssuancePolicyBaselineValuesAdditionalExtensionsObjectIdNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*CaPoolIssuancePolicyBaselineValuesAdditionalExtensionsObjectId)
	if !ok {
		desiredNotPointer, ok := d.(CaPoolIssuancePolicyBaselineValuesAdditionalExtensionsObjectId)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CaPoolIssuancePolicyBaselineValuesAdditionalExtensionsObjectId or *CaPoolIssuancePolicyBaselineValuesAdditionalExtensionsObjectId", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*CaPoolIssuancePolicyBaselineValuesAdditionalExtensionsObjectId)
	if !ok {
		actualNotPointer, ok := a.(CaPoolIssuancePolicyBaselineValuesAdditionalExtensionsObjectId)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CaPoolIssuancePolicyBaselineValuesAdditionalExtensionsObjectId", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.ObjectIdPath, actual.ObjectIdPath, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateCaPoolUpdateCaPoolOperation")}, fn.AddNest("ObjectIdPath")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareCaPoolIssuancePolicyIdentityConstraintsNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*CaPoolIssuancePolicyIdentityConstraints)
	if !ok {
		desiredNotPointer, ok := d.(CaPoolIssuancePolicyIdentityConstraints)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CaPoolIssuancePolicyIdentityConstraints or *CaPoolIssuancePolicyIdentityConstraints", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*CaPoolIssuancePolicyIdentityConstraints)
	if !ok {
		actualNotPointer, ok := a.(CaPoolIssuancePolicyIdentityConstraints)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CaPoolIssuancePolicyIdentityConstraints", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.CelExpression, actual.CelExpression, dcl.DiffInfo{ObjectFunction: compareCaPoolIssuancePolicyIdentityConstraintsCelExpressionNewStyle, EmptyObject: EmptyCaPoolIssuancePolicyIdentityConstraintsCelExpression, OperationSelector: dcl.TriggersOperation("updateCaPoolUpdateCaPoolOperation")}, fn.AddNest("CelExpression")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.AllowSubjectPassthrough, actual.AllowSubjectPassthrough, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateCaPoolUpdateCaPoolOperation")}, fn.AddNest("AllowSubjectPassthrough")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.AllowSubjectAltNamesPassthrough, actual.AllowSubjectAltNamesPassthrough, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateCaPoolUpdateCaPoolOperation")}, fn.AddNest("AllowSubjectAltNamesPassthrough")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareCaPoolIssuancePolicyIdentityConstraintsCelExpressionNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*CaPoolIssuancePolicyIdentityConstraintsCelExpression)
	if !ok {
		desiredNotPointer, ok := d.(CaPoolIssuancePolicyIdentityConstraintsCelExpression)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CaPoolIssuancePolicyIdentityConstraintsCelExpression or *CaPoolIssuancePolicyIdentityConstraintsCelExpression", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*CaPoolIssuancePolicyIdentityConstraintsCelExpression)
	if !ok {
		actualNotPointer, ok := a.(CaPoolIssuancePolicyIdentityConstraintsCelExpression)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CaPoolIssuancePolicyIdentityConstraintsCelExpression", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.Expression, actual.Expression, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateCaPoolUpdateCaPoolOperation")}, fn.AddNest("Expression")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Title, actual.Title, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateCaPoolUpdateCaPoolOperation")}, fn.AddNest("Title")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Description, actual.Description, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateCaPoolUpdateCaPoolOperation")}, fn.AddNest("Description")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Location, actual.Location, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateCaPoolUpdateCaPoolOperation")}, fn.AddNest("Location")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareCaPoolIssuancePolicyPassthroughExtensionsNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*CaPoolIssuancePolicyPassthroughExtensions)
	if !ok {
		desiredNotPointer, ok := d.(CaPoolIssuancePolicyPassthroughExtensions)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CaPoolIssuancePolicyPassthroughExtensions or *CaPoolIssuancePolicyPassthroughExtensions", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*CaPoolIssuancePolicyPassthroughExtensions)
	if !ok {
		actualNotPointer, ok := a.(CaPoolIssuancePolicyPassthroughExtensions)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CaPoolIssuancePolicyPassthroughExtensions", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.KnownExtensions, actual.KnownExtensions, dcl.DiffInfo{Type: "EnumType", OperationSelector: dcl.TriggersOperation("updateCaPoolUpdateCaPoolOperation")}, fn.AddNest("KnownExtensions")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.AdditionalExtensions, actual.AdditionalExtensions, dcl.DiffInfo{ObjectFunction: compareCaPoolIssuancePolicyPassthroughExtensionsAdditionalExtensionsNewStyle, EmptyObject: EmptyCaPoolIssuancePolicyPassthroughExtensionsAdditionalExtensions, OperationSelector: dcl.TriggersOperation("updateCaPoolUpdateCaPoolOperation")}, fn.AddNest("AdditionalExtensions")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareCaPoolIssuancePolicyPassthroughExtensionsAdditionalExtensionsNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*CaPoolIssuancePolicyPassthroughExtensionsAdditionalExtensions)
	if !ok {
		desiredNotPointer, ok := d.(CaPoolIssuancePolicyPassthroughExtensionsAdditionalExtensions)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CaPoolIssuancePolicyPassthroughExtensionsAdditionalExtensions or *CaPoolIssuancePolicyPassthroughExtensionsAdditionalExtensions", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*CaPoolIssuancePolicyPassthroughExtensionsAdditionalExtensions)
	if !ok {
		actualNotPointer, ok := a.(CaPoolIssuancePolicyPassthroughExtensionsAdditionalExtensions)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CaPoolIssuancePolicyPassthroughExtensionsAdditionalExtensions", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.ObjectIdPath, actual.ObjectIdPath, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateCaPoolUpdateCaPoolOperation")}, fn.AddNest("ObjectIdPath")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareCaPoolPublishingOptionsNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*CaPoolPublishingOptions)
	if !ok {
		desiredNotPointer, ok := d.(CaPoolPublishingOptions)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CaPoolPublishingOptions or *CaPoolPublishingOptions", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*CaPoolPublishingOptions)
	if !ok {
		actualNotPointer, ok := a.(CaPoolPublishingOptions)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CaPoolPublishingOptions", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.PublishCaCert, actual.PublishCaCert, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateCaPoolUpdateCaPoolOperation")}, fn.AddNest("PublishCaCert")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.PublishCrl, actual.PublishCrl, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateCaPoolUpdateCaPoolOperation")}, fn.AddNest("PublishCrl")); len(ds) != 0 || err != nil {
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
func (r *CaPool) urlNormalized() *CaPool {
	normalized := dcl.Copy(*r).(CaPool)
	normalized.Name = dcl.SelfLinkToName(r.Name)
	normalized.Project = dcl.SelfLinkToName(r.Project)
	normalized.Location = dcl.SelfLinkToName(r.Location)
	return &normalized
}

func (r *CaPool) updateURL(userBasePath, updateName string) (string, error) {
	nr := r.urlNormalized()
	if updateName == "UpdateCaPool" {
		fields := map[string]interface{}{
			"project":  dcl.ValueOrEmptyString(nr.Project),
			"location": dcl.ValueOrEmptyString(nr.Location),
			"name":     dcl.ValueOrEmptyString(nr.Name),
		}
		return dcl.URL("projects/{{project}}/locations/{{location}}/caPools/{{name}}", nr.basePath(), userBasePath, fields), nil

	}

	return "", fmt.Errorf("unknown update name: %s", updateName)
}

// marshal encodes the CaPool resource into JSON for a Create request, and
// performs transformations from the resource schema to the API schema if
// necessary.
func (r *CaPool) marshal(c *Client) ([]byte, error) {
	m, err := expandCaPool(c, r)
	if err != nil {
		return nil, fmt.Errorf("error marshalling CaPool: %w", err)
	}

	return json.Marshal(m)
}

// unmarshalCaPool decodes JSON responses into the CaPool resource schema.
func unmarshalCaPool(b []byte, c *Client, res *CaPool) (*CaPool, error) {
	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, err
	}
	return unmarshalMapCaPool(m, c, res)
}

func unmarshalMapCaPool(m map[string]interface{}, c *Client, res *CaPool) (*CaPool, error) {

	flattened := flattenCaPool(c, m, res)
	if flattened == nil {
		return nil, fmt.Errorf("attempted to flatten empty json object")
	}
	return flattened, nil
}

// expandCaPool expands CaPool into a JSON request object.
func expandCaPool(c *Client, f *CaPool) (map[string]interface{}, error) {
	m := make(map[string]interface{})
	res := f
	_ = res
	if v, err := dcl.DeriveField("projects/%s/locations/%s/caPools/%s", f.Name, dcl.SelfLinkToName(f.Project), dcl.SelfLinkToName(f.Location), dcl.SelfLinkToName(f.Name)); err != nil {
		return nil, fmt.Errorf("error expanding Name into name: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["name"] = v
	}
	if v := f.Tier; dcl.ValueShouldBeSent(v) {
		m["tier"] = v
	}
	if v, err := expandCaPoolIssuancePolicy(c, f.IssuancePolicy, res); err != nil {
		return nil, fmt.Errorf("error expanding IssuancePolicy into issuancePolicy: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["issuancePolicy"] = v
	}
	if v, err := expandCaPoolPublishingOptions(c, f.PublishingOptions, res); err != nil {
		return nil, fmt.Errorf("error expanding PublishingOptions into publishingOptions: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["publishingOptions"] = v
	}
	if v := f.Labels; dcl.ValueShouldBeSent(v) {
		m["labels"] = v
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

// flattenCaPool flattens CaPool from a JSON request object into the
// CaPool type.
func flattenCaPool(c *Client, i interface{}, res *CaPool) *CaPool {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}
	if len(m) == 0 {
		return nil
	}

	resultRes := &CaPool{}
	resultRes.Name = dcl.FlattenString(m["name"])
	resultRes.Tier = flattenCaPoolTierEnum(m["tier"])
	resultRes.IssuancePolicy = flattenCaPoolIssuancePolicy(c, m["issuancePolicy"], res)
	resultRes.PublishingOptions = flattenCaPoolPublishingOptions(c, m["publishingOptions"], res)
	resultRes.Labels = dcl.FlattenKeyValuePairs(m["labels"])
	resultRes.Project = dcl.FlattenString(m["project"])
	resultRes.Location = dcl.FlattenString(m["location"])

	return resultRes
}

// expandCaPoolIssuancePolicyMap expands the contents of CaPoolIssuancePolicy into a JSON
// request object.
func expandCaPoolIssuancePolicyMap(c *Client, f map[string]CaPoolIssuancePolicy, res *CaPool) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandCaPoolIssuancePolicy(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandCaPoolIssuancePolicySlice expands the contents of CaPoolIssuancePolicy into a JSON
// request object.
func expandCaPoolIssuancePolicySlice(c *Client, f []CaPoolIssuancePolicy, res *CaPool) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandCaPoolIssuancePolicy(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenCaPoolIssuancePolicyMap flattens the contents of CaPoolIssuancePolicy from a JSON
// response object.
func flattenCaPoolIssuancePolicyMap(c *Client, i interface{}, res *CaPool) map[string]CaPoolIssuancePolicy {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]CaPoolIssuancePolicy{}
	}

	if len(a) == 0 {
		return map[string]CaPoolIssuancePolicy{}
	}

	items := make(map[string]CaPoolIssuancePolicy)
	for k, item := range a {
		items[k] = *flattenCaPoolIssuancePolicy(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenCaPoolIssuancePolicySlice flattens the contents of CaPoolIssuancePolicy from a JSON
// response object.
func flattenCaPoolIssuancePolicySlice(c *Client, i interface{}, res *CaPool) []CaPoolIssuancePolicy {
	a, ok := i.([]interface{})
	if !ok {
		return []CaPoolIssuancePolicy{}
	}

	if len(a) == 0 {
		return []CaPoolIssuancePolicy{}
	}

	items := make([]CaPoolIssuancePolicy, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenCaPoolIssuancePolicy(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandCaPoolIssuancePolicy expands an instance of CaPoolIssuancePolicy into a JSON
// request object.
func expandCaPoolIssuancePolicy(c *Client, f *CaPoolIssuancePolicy, res *CaPool) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v, err := expandCaPoolIssuancePolicyAllowedKeyTypesSlice(c, f.AllowedKeyTypes, res); err != nil {
		return nil, fmt.Errorf("error expanding AllowedKeyTypes into allowedKeyTypes: %w", err)
	} else if v != nil {
		m["allowedKeyTypes"] = v
	}
	if v := f.MaximumLifetime; !dcl.IsEmptyValueIndirect(v) {
		m["maximumLifetime"] = v
	}
	if v, err := expandCaPoolIssuancePolicyAllowedIssuanceModes(c, f.AllowedIssuanceModes, res); err != nil {
		return nil, fmt.Errorf("error expanding AllowedIssuanceModes into allowedIssuanceModes: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["allowedIssuanceModes"] = v
	}
	if v, err := expandCaPoolIssuancePolicyBaselineValues(c, f.BaselineValues, res); err != nil {
		return nil, fmt.Errorf("error expanding BaselineValues into baselineValues: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["baselineValues"] = v
	}
	if v, err := expandCaPoolIssuancePolicyIdentityConstraints(c, f.IdentityConstraints, res); err != nil {
		return nil, fmt.Errorf("error expanding IdentityConstraints into identityConstraints: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["identityConstraints"] = v
	}
	if v, err := expandCaPoolIssuancePolicyPassthroughExtensions(c, f.PassthroughExtensions, res); err != nil {
		return nil, fmt.Errorf("error expanding PassthroughExtensions into passthroughExtensions: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["passthroughExtensions"] = v
	}

	return m, nil
}

// flattenCaPoolIssuancePolicy flattens an instance of CaPoolIssuancePolicy from a JSON
// response object.
func flattenCaPoolIssuancePolicy(c *Client, i interface{}, res *CaPool) *CaPoolIssuancePolicy {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &CaPoolIssuancePolicy{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyCaPoolIssuancePolicy
	}
	r.AllowedKeyTypes = flattenCaPoolIssuancePolicyAllowedKeyTypesSlice(c, m["allowedKeyTypes"], res)
	r.MaximumLifetime = dcl.FlattenString(m["maximumLifetime"])
	r.AllowedIssuanceModes = flattenCaPoolIssuancePolicyAllowedIssuanceModes(c, m["allowedIssuanceModes"], res)
	r.BaselineValues = flattenCaPoolIssuancePolicyBaselineValues(c, m["baselineValues"], res)
	r.IdentityConstraints = flattenCaPoolIssuancePolicyIdentityConstraints(c, m["identityConstraints"], res)
	r.PassthroughExtensions = flattenCaPoolIssuancePolicyPassthroughExtensions(c, m["passthroughExtensions"], res)

	return r
}

// expandCaPoolIssuancePolicyAllowedKeyTypesMap expands the contents of CaPoolIssuancePolicyAllowedKeyTypes into a JSON
// request object.
func expandCaPoolIssuancePolicyAllowedKeyTypesMap(c *Client, f map[string]CaPoolIssuancePolicyAllowedKeyTypes, res *CaPool) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandCaPoolIssuancePolicyAllowedKeyTypes(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandCaPoolIssuancePolicyAllowedKeyTypesSlice expands the contents of CaPoolIssuancePolicyAllowedKeyTypes into a JSON
// request object.
func expandCaPoolIssuancePolicyAllowedKeyTypesSlice(c *Client, f []CaPoolIssuancePolicyAllowedKeyTypes, res *CaPool) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandCaPoolIssuancePolicyAllowedKeyTypes(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenCaPoolIssuancePolicyAllowedKeyTypesMap flattens the contents of CaPoolIssuancePolicyAllowedKeyTypes from a JSON
// response object.
func flattenCaPoolIssuancePolicyAllowedKeyTypesMap(c *Client, i interface{}, res *CaPool) map[string]CaPoolIssuancePolicyAllowedKeyTypes {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]CaPoolIssuancePolicyAllowedKeyTypes{}
	}

	if len(a) == 0 {
		return map[string]CaPoolIssuancePolicyAllowedKeyTypes{}
	}

	items := make(map[string]CaPoolIssuancePolicyAllowedKeyTypes)
	for k, item := range a {
		items[k] = *flattenCaPoolIssuancePolicyAllowedKeyTypes(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenCaPoolIssuancePolicyAllowedKeyTypesSlice flattens the contents of CaPoolIssuancePolicyAllowedKeyTypes from a JSON
// response object.
func flattenCaPoolIssuancePolicyAllowedKeyTypesSlice(c *Client, i interface{}, res *CaPool) []CaPoolIssuancePolicyAllowedKeyTypes {
	a, ok := i.([]interface{})
	if !ok {
		return []CaPoolIssuancePolicyAllowedKeyTypes{}
	}

	if len(a) == 0 {
		return []CaPoolIssuancePolicyAllowedKeyTypes{}
	}

	items := make([]CaPoolIssuancePolicyAllowedKeyTypes, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenCaPoolIssuancePolicyAllowedKeyTypes(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandCaPoolIssuancePolicyAllowedKeyTypes expands an instance of CaPoolIssuancePolicyAllowedKeyTypes into a JSON
// request object.
func expandCaPoolIssuancePolicyAllowedKeyTypes(c *Client, f *CaPoolIssuancePolicyAllowedKeyTypes, res *CaPool) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v, err := expandCaPoolIssuancePolicyAllowedKeyTypesRsa(c, f.Rsa, res); err != nil {
		return nil, fmt.Errorf("error expanding Rsa into rsa: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["rsa"] = v
	}
	if v, err := expandCaPoolIssuancePolicyAllowedKeyTypesEllipticCurve(c, f.EllipticCurve, res); err != nil {
		return nil, fmt.Errorf("error expanding EllipticCurve into ellipticCurve: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["ellipticCurve"] = v
	}

	return m, nil
}

// flattenCaPoolIssuancePolicyAllowedKeyTypes flattens an instance of CaPoolIssuancePolicyAllowedKeyTypes from a JSON
// response object.
func flattenCaPoolIssuancePolicyAllowedKeyTypes(c *Client, i interface{}, res *CaPool) *CaPoolIssuancePolicyAllowedKeyTypes {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &CaPoolIssuancePolicyAllowedKeyTypes{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyCaPoolIssuancePolicyAllowedKeyTypes
	}
	r.Rsa = flattenCaPoolIssuancePolicyAllowedKeyTypesRsa(c, m["rsa"], res)
	r.EllipticCurve = flattenCaPoolIssuancePolicyAllowedKeyTypesEllipticCurve(c, m["ellipticCurve"], res)

	return r
}

// expandCaPoolIssuancePolicyAllowedKeyTypesRsaMap expands the contents of CaPoolIssuancePolicyAllowedKeyTypesRsa into a JSON
// request object.
func expandCaPoolIssuancePolicyAllowedKeyTypesRsaMap(c *Client, f map[string]CaPoolIssuancePolicyAllowedKeyTypesRsa, res *CaPool) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandCaPoolIssuancePolicyAllowedKeyTypesRsa(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandCaPoolIssuancePolicyAllowedKeyTypesRsaSlice expands the contents of CaPoolIssuancePolicyAllowedKeyTypesRsa into a JSON
// request object.
func expandCaPoolIssuancePolicyAllowedKeyTypesRsaSlice(c *Client, f []CaPoolIssuancePolicyAllowedKeyTypesRsa, res *CaPool) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandCaPoolIssuancePolicyAllowedKeyTypesRsa(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenCaPoolIssuancePolicyAllowedKeyTypesRsaMap flattens the contents of CaPoolIssuancePolicyAllowedKeyTypesRsa from a JSON
// response object.
func flattenCaPoolIssuancePolicyAllowedKeyTypesRsaMap(c *Client, i interface{}, res *CaPool) map[string]CaPoolIssuancePolicyAllowedKeyTypesRsa {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]CaPoolIssuancePolicyAllowedKeyTypesRsa{}
	}

	if len(a) == 0 {
		return map[string]CaPoolIssuancePolicyAllowedKeyTypesRsa{}
	}

	items := make(map[string]CaPoolIssuancePolicyAllowedKeyTypesRsa)
	for k, item := range a {
		items[k] = *flattenCaPoolIssuancePolicyAllowedKeyTypesRsa(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenCaPoolIssuancePolicyAllowedKeyTypesRsaSlice flattens the contents of CaPoolIssuancePolicyAllowedKeyTypesRsa from a JSON
// response object.
func flattenCaPoolIssuancePolicyAllowedKeyTypesRsaSlice(c *Client, i interface{}, res *CaPool) []CaPoolIssuancePolicyAllowedKeyTypesRsa {
	a, ok := i.([]interface{})
	if !ok {
		return []CaPoolIssuancePolicyAllowedKeyTypesRsa{}
	}

	if len(a) == 0 {
		return []CaPoolIssuancePolicyAllowedKeyTypesRsa{}
	}

	items := make([]CaPoolIssuancePolicyAllowedKeyTypesRsa, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenCaPoolIssuancePolicyAllowedKeyTypesRsa(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandCaPoolIssuancePolicyAllowedKeyTypesRsa expands an instance of CaPoolIssuancePolicyAllowedKeyTypesRsa into a JSON
// request object.
func expandCaPoolIssuancePolicyAllowedKeyTypesRsa(c *Client, f *CaPoolIssuancePolicyAllowedKeyTypesRsa, res *CaPool) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.MinModulusSize; !dcl.IsEmptyValueIndirect(v) {
		m["minModulusSize"] = v
	}
	if v := f.MaxModulusSize; !dcl.IsEmptyValueIndirect(v) {
		m["maxModulusSize"] = v
	}

	return m, nil
}

// flattenCaPoolIssuancePolicyAllowedKeyTypesRsa flattens an instance of CaPoolIssuancePolicyAllowedKeyTypesRsa from a JSON
// response object.
func flattenCaPoolIssuancePolicyAllowedKeyTypesRsa(c *Client, i interface{}, res *CaPool) *CaPoolIssuancePolicyAllowedKeyTypesRsa {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &CaPoolIssuancePolicyAllowedKeyTypesRsa{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyCaPoolIssuancePolicyAllowedKeyTypesRsa
	}
	r.MinModulusSize = dcl.FlattenInteger(m["minModulusSize"])
	r.MaxModulusSize = dcl.FlattenInteger(m["maxModulusSize"])

	return r
}

// expandCaPoolIssuancePolicyAllowedKeyTypesEllipticCurveMap expands the contents of CaPoolIssuancePolicyAllowedKeyTypesEllipticCurve into a JSON
// request object.
func expandCaPoolIssuancePolicyAllowedKeyTypesEllipticCurveMap(c *Client, f map[string]CaPoolIssuancePolicyAllowedKeyTypesEllipticCurve, res *CaPool) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandCaPoolIssuancePolicyAllowedKeyTypesEllipticCurve(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandCaPoolIssuancePolicyAllowedKeyTypesEllipticCurveSlice expands the contents of CaPoolIssuancePolicyAllowedKeyTypesEllipticCurve into a JSON
// request object.
func expandCaPoolIssuancePolicyAllowedKeyTypesEllipticCurveSlice(c *Client, f []CaPoolIssuancePolicyAllowedKeyTypesEllipticCurve, res *CaPool) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandCaPoolIssuancePolicyAllowedKeyTypesEllipticCurve(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenCaPoolIssuancePolicyAllowedKeyTypesEllipticCurveMap flattens the contents of CaPoolIssuancePolicyAllowedKeyTypesEllipticCurve from a JSON
// response object.
func flattenCaPoolIssuancePolicyAllowedKeyTypesEllipticCurveMap(c *Client, i interface{}, res *CaPool) map[string]CaPoolIssuancePolicyAllowedKeyTypesEllipticCurve {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]CaPoolIssuancePolicyAllowedKeyTypesEllipticCurve{}
	}

	if len(a) == 0 {
		return map[string]CaPoolIssuancePolicyAllowedKeyTypesEllipticCurve{}
	}

	items := make(map[string]CaPoolIssuancePolicyAllowedKeyTypesEllipticCurve)
	for k, item := range a {
		items[k] = *flattenCaPoolIssuancePolicyAllowedKeyTypesEllipticCurve(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenCaPoolIssuancePolicyAllowedKeyTypesEllipticCurveSlice flattens the contents of CaPoolIssuancePolicyAllowedKeyTypesEllipticCurve from a JSON
// response object.
func flattenCaPoolIssuancePolicyAllowedKeyTypesEllipticCurveSlice(c *Client, i interface{}, res *CaPool) []CaPoolIssuancePolicyAllowedKeyTypesEllipticCurve {
	a, ok := i.([]interface{})
	if !ok {
		return []CaPoolIssuancePolicyAllowedKeyTypesEllipticCurve{}
	}

	if len(a) == 0 {
		return []CaPoolIssuancePolicyAllowedKeyTypesEllipticCurve{}
	}

	items := make([]CaPoolIssuancePolicyAllowedKeyTypesEllipticCurve, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenCaPoolIssuancePolicyAllowedKeyTypesEllipticCurve(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandCaPoolIssuancePolicyAllowedKeyTypesEllipticCurve expands an instance of CaPoolIssuancePolicyAllowedKeyTypesEllipticCurve into a JSON
// request object.
func expandCaPoolIssuancePolicyAllowedKeyTypesEllipticCurve(c *Client, f *CaPoolIssuancePolicyAllowedKeyTypesEllipticCurve, res *CaPool) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.SignatureAlgorithm; !dcl.IsEmptyValueIndirect(v) {
		m["signatureAlgorithm"] = v
	}

	return m, nil
}

// flattenCaPoolIssuancePolicyAllowedKeyTypesEllipticCurve flattens an instance of CaPoolIssuancePolicyAllowedKeyTypesEllipticCurve from a JSON
// response object.
func flattenCaPoolIssuancePolicyAllowedKeyTypesEllipticCurve(c *Client, i interface{}, res *CaPool) *CaPoolIssuancePolicyAllowedKeyTypesEllipticCurve {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &CaPoolIssuancePolicyAllowedKeyTypesEllipticCurve{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyCaPoolIssuancePolicyAllowedKeyTypesEllipticCurve
	}
	r.SignatureAlgorithm = flattenCaPoolIssuancePolicyAllowedKeyTypesEllipticCurveSignatureAlgorithmEnum(m["signatureAlgorithm"])

	return r
}

// expandCaPoolIssuancePolicyAllowedIssuanceModesMap expands the contents of CaPoolIssuancePolicyAllowedIssuanceModes into a JSON
// request object.
func expandCaPoolIssuancePolicyAllowedIssuanceModesMap(c *Client, f map[string]CaPoolIssuancePolicyAllowedIssuanceModes, res *CaPool) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandCaPoolIssuancePolicyAllowedIssuanceModes(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandCaPoolIssuancePolicyAllowedIssuanceModesSlice expands the contents of CaPoolIssuancePolicyAllowedIssuanceModes into a JSON
// request object.
func expandCaPoolIssuancePolicyAllowedIssuanceModesSlice(c *Client, f []CaPoolIssuancePolicyAllowedIssuanceModes, res *CaPool) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandCaPoolIssuancePolicyAllowedIssuanceModes(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenCaPoolIssuancePolicyAllowedIssuanceModesMap flattens the contents of CaPoolIssuancePolicyAllowedIssuanceModes from a JSON
// response object.
func flattenCaPoolIssuancePolicyAllowedIssuanceModesMap(c *Client, i interface{}, res *CaPool) map[string]CaPoolIssuancePolicyAllowedIssuanceModes {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]CaPoolIssuancePolicyAllowedIssuanceModes{}
	}

	if len(a) == 0 {
		return map[string]CaPoolIssuancePolicyAllowedIssuanceModes{}
	}

	items := make(map[string]CaPoolIssuancePolicyAllowedIssuanceModes)
	for k, item := range a {
		items[k] = *flattenCaPoolIssuancePolicyAllowedIssuanceModes(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenCaPoolIssuancePolicyAllowedIssuanceModesSlice flattens the contents of CaPoolIssuancePolicyAllowedIssuanceModes from a JSON
// response object.
func flattenCaPoolIssuancePolicyAllowedIssuanceModesSlice(c *Client, i interface{}, res *CaPool) []CaPoolIssuancePolicyAllowedIssuanceModes {
	a, ok := i.([]interface{})
	if !ok {
		return []CaPoolIssuancePolicyAllowedIssuanceModes{}
	}

	if len(a) == 0 {
		return []CaPoolIssuancePolicyAllowedIssuanceModes{}
	}

	items := make([]CaPoolIssuancePolicyAllowedIssuanceModes, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenCaPoolIssuancePolicyAllowedIssuanceModes(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandCaPoolIssuancePolicyAllowedIssuanceModes expands an instance of CaPoolIssuancePolicyAllowedIssuanceModes into a JSON
// request object.
func expandCaPoolIssuancePolicyAllowedIssuanceModes(c *Client, f *CaPoolIssuancePolicyAllowedIssuanceModes, res *CaPool) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.AllowCsrBasedIssuance; !dcl.IsEmptyValueIndirect(v) {
		m["allowCsrBasedIssuance"] = v
	}
	if v := f.AllowConfigBasedIssuance; !dcl.IsEmptyValueIndirect(v) {
		m["allowConfigBasedIssuance"] = v
	}

	return m, nil
}

// flattenCaPoolIssuancePolicyAllowedIssuanceModes flattens an instance of CaPoolIssuancePolicyAllowedIssuanceModes from a JSON
// response object.
func flattenCaPoolIssuancePolicyAllowedIssuanceModes(c *Client, i interface{}, res *CaPool) *CaPoolIssuancePolicyAllowedIssuanceModes {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &CaPoolIssuancePolicyAllowedIssuanceModes{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyCaPoolIssuancePolicyAllowedIssuanceModes
	}
	r.AllowCsrBasedIssuance = dcl.FlattenBool(m["allowCsrBasedIssuance"])
	r.AllowConfigBasedIssuance = dcl.FlattenBool(m["allowConfigBasedIssuance"])

	return r
}

// expandCaPoolIssuancePolicyBaselineValuesMap expands the contents of CaPoolIssuancePolicyBaselineValues into a JSON
// request object.
func expandCaPoolIssuancePolicyBaselineValuesMap(c *Client, f map[string]CaPoolIssuancePolicyBaselineValues, res *CaPool) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandCaPoolIssuancePolicyBaselineValues(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandCaPoolIssuancePolicyBaselineValuesSlice expands the contents of CaPoolIssuancePolicyBaselineValues into a JSON
// request object.
func expandCaPoolIssuancePolicyBaselineValuesSlice(c *Client, f []CaPoolIssuancePolicyBaselineValues, res *CaPool) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandCaPoolIssuancePolicyBaselineValues(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenCaPoolIssuancePolicyBaselineValuesMap flattens the contents of CaPoolIssuancePolicyBaselineValues from a JSON
// response object.
func flattenCaPoolIssuancePolicyBaselineValuesMap(c *Client, i interface{}, res *CaPool) map[string]CaPoolIssuancePolicyBaselineValues {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]CaPoolIssuancePolicyBaselineValues{}
	}

	if len(a) == 0 {
		return map[string]CaPoolIssuancePolicyBaselineValues{}
	}

	items := make(map[string]CaPoolIssuancePolicyBaselineValues)
	for k, item := range a {
		items[k] = *flattenCaPoolIssuancePolicyBaselineValues(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenCaPoolIssuancePolicyBaselineValuesSlice flattens the contents of CaPoolIssuancePolicyBaselineValues from a JSON
// response object.
func flattenCaPoolIssuancePolicyBaselineValuesSlice(c *Client, i interface{}, res *CaPool) []CaPoolIssuancePolicyBaselineValues {
	a, ok := i.([]interface{})
	if !ok {
		return []CaPoolIssuancePolicyBaselineValues{}
	}

	if len(a) == 0 {
		return []CaPoolIssuancePolicyBaselineValues{}
	}

	items := make([]CaPoolIssuancePolicyBaselineValues, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenCaPoolIssuancePolicyBaselineValues(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandCaPoolIssuancePolicyBaselineValues expands an instance of CaPoolIssuancePolicyBaselineValues into a JSON
// request object.
func expandCaPoolIssuancePolicyBaselineValues(c *Client, f *CaPoolIssuancePolicyBaselineValues, res *CaPool) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v, err := expandCaPoolIssuancePolicyBaselineValuesKeyUsage(c, f.KeyUsage, res); err != nil {
		return nil, fmt.Errorf("error expanding KeyUsage into keyUsage: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["keyUsage"] = v
	}
	if v, err := expandCaPoolIssuancePolicyBaselineValuesCAOptions(c, f.CaOptions, res); err != nil {
		return nil, fmt.Errorf("error expanding CaOptions into caOptions: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["caOptions"] = v
	}
	if v, err := expandCaPoolIssuancePolicyBaselineValuesPolicyIdsSlice(c, f.PolicyIds, res); err != nil {
		return nil, fmt.Errorf("error expanding PolicyIds into policyIds: %w", err)
	} else if v != nil {
		m["policyIds"] = v
	}
	if v := f.AiaOcspServers; v != nil {
		m["aiaOcspServers"] = v
	}
	if v, err := expandCaPoolIssuancePolicyBaselineValuesAdditionalExtensionsSlice(c, f.AdditionalExtensions, res); err != nil {
		return nil, fmt.Errorf("error expanding AdditionalExtensions into additionalExtensions: %w", err)
	} else if v != nil {
		m["additionalExtensions"] = v
	}

	return m, nil
}

// flattenCaPoolIssuancePolicyBaselineValues flattens an instance of CaPoolIssuancePolicyBaselineValues from a JSON
// response object.
func flattenCaPoolIssuancePolicyBaselineValues(c *Client, i interface{}, res *CaPool) *CaPoolIssuancePolicyBaselineValues {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &CaPoolIssuancePolicyBaselineValues{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyCaPoolIssuancePolicyBaselineValues
	}
	r.KeyUsage = flattenCaPoolIssuancePolicyBaselineValuesKeyUsage(c, m["keyUsage"], res)
	r.CaOptions = flattenCaPoolIssuancePolicyBaselineValuesCAOptions(c, m["caOptions"], res)
	r.PolicyIds = flattenCaPoolIssuancePolicyBaselineValuesPolicyIdsSlice(c, m["policyIds"], res)
	r.AiaOcspServers = dcl.FlattenStringSlice(m["aiaOcspServers"])
	r.AdditionalExtensions = flattenCaPoolIssuancePolicyBaselineValuesAdditionalExtensionsSlice(c, m["additionalExtensions"], res)

	return r
}

// expandCaPoolIssuancePolicyBaselineValuesKeyUsageMap expands the contents of CaPoolIssuancePolicyBaselineValuesKeyUsage into a JSON
// request object.
func expandCaPoolIssuancePolicyBaselineValuesKeyUsageMap(c *Client, f map[string]CaPoolIssuancePolicyBaselineValuesKeyUsage, res *CaPool) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandCaPoolIssuancePolicyBaselineValuesKeyUsage(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandCaPoolIssuancePolicyBaselineValuesKeyUsageSlice expands the contents of CaPoolIssuancePolicyBaselineValuesKeyUsage into a JSON
// request object.
func expandCaPoolIssuancePolicyBaselineValuesKeyUsageSlice(c *Client, f []CaPoolIssuancePolicyBaselineValuesKeyUsage, res *CaPool) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandCaPoolIssuancePolicyBaselineValuesKeyUsage(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenCaPoolIssuancePolicyBaselineValuesKeyUsageMap flattens the contents of CaPoolIssuancePolicyBaselineValuesKeyUsage from a JSON
// response object.
func flattenCaPoolIssuancePolicyBaselineValuesKeyUsageMap(c *Client, i interface{}, res *CaPool) map[string]CaPoolIssuancePolicyBaselineValuesKeyUsage {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]CaPoolIssuancePolicyBaselineValuesKeyUsage{}
	}

	if len(a) == 0 {
		return map[string]CaPoolIssuancePolicyBaselineValuesKeyUsage{}
	}

	items := make(map[string]CaPoolIssuancePolicyBaselineValuesKeyUsage)
	for k, item := range a {
		items[k] = *flattenCaPoolIssuancePolicyBaselineValuesKeyUsage(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenCaPoolIssuancePolicyBaselineValuesKeyUsageSlice flattens the contents of CaPoolIssuancePolicyBaselineValuesKeyUsage from a JSON
// response object.
func flattenCaPoolIssuancePolicyBaselineValuesKeyUsageSlice(c *Client, i interface{}, res *CaPool) []CaPoolIssuancePolicyBaselineValuesKeyUsage {
	a, ok := i.([]interface{})
	if !ok {
		return []CaPoolIssuancePolicyBaselineValuesKeyUsage{}
	}

	if len(a) == 0 {
		return []CaPoolIssuancePolicyBaselineValuesKeyUsage{}
	}

	items := make([]CaPoolIssuancePolicyBaselineValuesKeyUsage, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenCaPoolIssuancePolicyBaselineValuesKeyUsage(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandCaPoolIssuancePolicyBaselineValuesKeyUsage expands an instance of CaPoolIssuancePolicyBaselineValuesKeyUsage into a JSON
// request object.
func expandCaPoolIssuancePolicyBaselineValuesKeyUsage(c *Client, f *CaPoolIssuancePolicyBaselineValuesKeyUsage, res *CaPool) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v, err := expandCaPoolIssuancePolicyBaselineValuesKeyUsageBaseKeyUsage(c, f.BaseKeyUsage, res); err != nil {
		return nil, fmt.Errorf("error expanding BaseKeyUsage into baseKeyUsage: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["baseKeyUsage"] = v
	}
	if v, err := expandCaPoolIssuancePolicyBaselineValuesKeyUsageExtendedKeyUsage(c, f.ExtendedKeyUsage, res); err != nil {
		return nil, fmt.Errorf("error expanding ExtendedKeyUsage into extendedKeyUsage: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["extendedKeyUsage"] = v
	}
	if v, err := expandCaPoolIssuancePolicyBaselineValuesKeyUsageUnknownExtendedKeyUsagesSlice(c, f.UnknownExtendedKeyUsages, res); err != nil {
		return nil, fmt.Errorf("error expanding UnknownExtendedKeyUsages into unknownExtendedKeyUsages: %w", err)
	} else if v != nil {
		m["unknownExtendedKeyUsages"] = v
	}

	return m, nil
}

// flattenCaPoolIssuancePolicyBaselineValuesKeyUsage flattens an instance of CaPoolIssuancePolicyBaselineValuesKeyUsage from a JSON
// response object.
func flattenCaPoolIssuancePolicyBaselineValuesKeyUsage(c *Client, i interface{}, res *CaPool) *CaPoolIssuancePolicyBaselineValuesKeyUsage {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &CaPoolIssuancePolicyBaselineValuesKeyUsage{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyCaPoolIssuancePolicyBaselineValuesKeyUsage
	}
	r.BaseKeyUsage = flattenCaPoolIssuancePolicyBaselineValuesKeyUsageBaseKeyUsage(c, m["baseKeyUsage"], res)
	r.ExtendedKeyUsage = flattenCaPoolIssuancePolicyBaselineValuesKeyUsageExtendedKeyUsage(c, m["extendedKeyUsage"], res)
	r.UnknownExtendedKeyUsages = flattenCaPoolIssuancePolicyBaselineValuesKeyUsageUnknownExtendedKeyUsagesSlice(c, m["unknownExtendedKeyUsages"], res)

	return r
}

// expandCaPoolIssuancePolicyBaselineValuesKeyUsageBaseKeyUsageMap expands the contents of CaPoolIssuancePolicyBaselineValuesKeyUsageBaseKeyUsage into a JSON
// request object.
func expandCaPoolIssuancePolicyBaselineValuesKeyUsageBaseKeyUsageMap(c *Client, f map[string]CaPoolIssuancePolicyBaselineValuesKeyUsageBaseKeyUsage, res *CaPool) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandCaPoolIssuancePolicyBaselineValuesKeyUsageBaseKeyUsage(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandCaPoolIssuancePolicyBaselineValuesKeyUsageBaseKeyUsageSlice expands the contents of CaPoolIssuancePolicyBaselineValuesKeyUsageBaseKeyUsage into a JSON
// request object.
func expandCaPoolIssuancePolicyBaselineValuesKeyUsageBaseKeyUsageSlice(c *Client, f []CaPoolIssuancePolicyBaselineValuesKeyUsageBaseKeyUsage, res *CaPool) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandCaPoolIssuancePolicyBaselineValuesKeyUsageBaseKeyUsage(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenCaPoolIssuancePolicyBaselineValuesKeyUsageBaseKeyUsageMap flattens the contents of CaPoolIssuancePolicyBaselineValuesKeyUsageBaseKeyUsage from a JSON
// response object.
func flattenCaPoolIssuancePolicyBaselineValuesKeyUsageBaseKeyUsageMap(c *Client, i interface{}, res *CaPool) map[string]CaPoolIssuancePolicyBaselineValuesKeyUsageBaseKeyUsage {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]CaPoolIssuancePolicyBaselineValuesKeyUsageBaseKeyUsage{}
	}

	if len(a) == 0 {
		return map[string]CaPoolIssuancePolicyBaselineValuesKeyUsageBaseKeyUsage{}
	}

	items := make(map[string]CaPoolIssuancePolicyBaselineValuesKeyUsageBaseKeyUsage)
	for k, item := range a {
		items[k] = *flattenCaPoolIssuancePolicyBaselineValuesKeyUsageBaseKeyUsage(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenCaPoolIssuancePolicyBaselineValuesKeyUsageBaseKeyUsageSlice flattens the contents of CaPoolIssuancePolicyBaselineValuesKeyUsageBaseKeyUsage from a JSON
// response object.
func flattenCaPoolIssuancePolicyBaselineValuesKeyUsageBaseKeyUsageSlice(c *Client, i interface{}, res *CaPool) []CaPoolIssuancePolicyBaselineValuesKeyUsageBaseKeyUsage {
	a, ok := i.([]interface{})
	if !ok {
		return []CaPoolIssuancePolicyBaselineValuesKeyUsageBaseKeyUsage{}
	}

	if len(a) == 0 {
		return []CaPoolIssuancePolicyBaselineValuesKeyUsageBaseKeyUsage{}
	}

	items := make([]CaPoolIssuancePolicyBaselineValuesKeyUsageBaseKeyUsage, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenCaPoolIssuancePolicyBaselineValuesKeyUsageBaseKeyUsage(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandCaPoolIssuancePolicyBaselineValuesKeyUsageBaseKeyUsage expands an instance of CaPoolIssuancePolicyBaselineValuesKeyUsageBaseKeyUsage into a JSON
// request object.
func expandCaPoolIssuancePolicyBaselineValuesKeyUsageBaseKeyUsage(c *Client, f *CaPoolIssuancePolicyBaselineValuesKeyUsageBaseKeyUsage, res *CaPool) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.DigitalSignature; !dcl.IsEmptyValueIndirect(v) {
		m["digitalSignature"] = v
	}
	if v := f.ContentCommitment; !dcl.IsEmptyValueIndirect(v) {
		m["contentCommitment"] = v
	}
	if v := f.KeyEncipherment; !dcl.IsEmptyValueIndirect(v) {
		m["keyEncipherment"] = v
	}
	if v := f.DataEncipherment; !dcl.IsEmptyValueIndirect(v) {
		m["dataEncipherment"] = v
	}
	if v := f.KeyAgreement; !dcl.IsEmptyValueIndirect(v) {
		m["keyAgreement"] = v
	}
	if v := f.CertSign; !dcl.IsEmptyValueIndirect(v) {
		m["certSign"] = v
	}
	if v := f.CrlSign; !dcl.IsEmptyValueIndirect(v) {
		m["crlSign"] = v
	}
	if v := f.EncipherOnly; !dcl.IsEmptyValueIndirect(v) {
		m["encipherOnly"] = v
	}
	if v := f.DecipherOnly; !dcl.IsEmptyValueIndirect(v) {
		m["decipherOnly"] = v
	}

	return m, nil
}

// flattenCaPoolIssuancePolicyBaselineValuesKeyUsageBaseKeyUsage flattens an instance of CaPoolIssuancePolicyBaselineValuesKeyUsageBaseKeyUsage from a JSON
// response object.
func flattenCaPoolIssuancePolicyBaselineValuesKeyUsageBaseKeyUsage(c *Client, i interface{}, res *CaPool) *CaPoolIssuancePolicyBaselineValuesKeyUsageBaseKeyUsage {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &CaPoolIssuancePolicyBaselineValuesKeyUsageBaseKeyUsage{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyCaPoolIssuancePolicyBaselineValuesKeyUsageBaseKeyUsage
	}
	r.DigitalSignature = dcl.FlattenBool(m["digitalSignature"])
	r.ContentCommitment = dcl.FlattenBool(m["contentCommitment"])
	r.KeyEncipherment = dcl.FlattenBool(m["keyEncipherment"])
	r.DataEncipherment = dcl.FlattenBool(m["dataEncipherment"])
	r.KeyAgreement = dcl.FlattenBool(m["keyAgreement"])
	r.CertSign = dcl.FlattenBool(m["certSign"])
	r.CrlSign = dcl.FlattenBool(m["crlSign"])
	r.EncipherOnly = dcl.FlattenBool(m["encipherOnly"])
	r.DecipherOnly = dcl.FlattenBool(m["decipherOnly"])

	return r
}

// expandCaPoolIssuancePolicyBaselineValuesKeyUsageExtendedKeyUsageMap expands the contents of CaPoolIssuancePolicyBaselineValuesKeyUsageExtendedKeyUsage into a JSON
// request object.
func expandCaPoolIssuancePolicyBaselineValuesKeyUsageExtendedKeyUsageMap(c *Client, f map[string]CaPoolIssuancePolicyBaselineValuesKeyUsageExtendedKeyUsage, res *CaPool) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandCaPoolIssuancePolicyBaselineValuesKeyUsageExtendedKeyUsage(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandCaPoolIssuancePolicyBaselineValuesKeyUsageExtendedKeyUsageSlice expands the contents of CaPoolIssuancePolicyBaselineValuesKeyUsageExtendedKeyUsage into a JSON
// request object.
func expandCaPoolIssuancePolicyBaselineValuesKeyUsageExtendedKeyUsageSlice(c *Client, f []CaPoolIssuancePolicyBaselineValuesKeyUsageExtendedKeyUsage, res *CaPool) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandCaPoolIssuancePolicyBaselineValuesKeyUsageExtendedKeyUsage(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenCaPoolIssuancePolicyBaselineValuesKeyUsageExtendedKeyUsageMap flattens the contents of CaPoolIssuancePolicyBaselineValuesKeyUsageExtendedKeyUsage from a JSON
// response object.
func flattenCaPoolIssuancePolicyBaselineValuesKeyUsageExtendedKeyUsageMap(c *Client, i interface{}, res *CaPool) map[string]CaPoolIssuancePolicyBaselineValuesKeyUsageExtendedKeyUsage {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]CaPoolIssuancePolicyBaselineValuesKeyUsageExtendedKeyUsage{}
	}

	if len(a) == 0 {
		return map[string]CaPoolIssuancePolicyBaselineValuesKeyUsageExtendedKeyUsage{}
	}

	items := make(map[string]CaPoolIssuancePolicyBaselineValuesKeyUsageExtendedKeyUsage)
	for k, item := range a {
		items[k] = *flattenCaPoolIssuancePolicyBaselineValuesKeyUsageExtendedKeyUsage(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenCaPoolIssuancePolicyBaselineValuesKeyUsageExtendedKeyUsageSlice flattens the contents of CaPoolIssuancePolicyBaselineValuesKeyUsageExtendedKeyUsage from a JSON
// response object.
func flattenCaPoolIssuancePolicyBaselineValuesKeyUsageExtendedKeyUsageSlice(c *Client, i interface{}, res *CaPool) []CaPoolIssuancePolicyBaselineValuesKeyUsageExtendedKeyUsage {
	a, ok := i.([]interface{})
	if !ok {
		return []CaPoolIssuancePolicyBaselineValuesKeyUsageExtendedKeyUsage{}
	}

	if len(a) == 0 {
		return []CaPoolIssuancePolicyBaselineValuesKeyUsageExtendedKeyUsage{}
	}

	items := make([]CaPoolIssuancePolicyBaselineValuesKeyUsageExtendedKeyUsage, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenCaPoolIssuancePolicyBaselineValuesKeyUsageExtendedKeyUsage(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandCaPoolIssuancePolicyBaselineValuesKeyUsageExtendedKeyUsage expands an instance of CaPoolIssuancePolicyBaselineValuesKeyUsageExtendedKeyUsage into a JSON
// request object.
func expandCaPoolIssuancePolicyBaselineValuesKeyUsageExtendedKeyUsage(c *Client, f *CaPoolIssuancePolicyBaselineValuesKeyUsageExtendedKeyUsage, res *CaPool) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.ServerAuth; !dcl.IsEmptyValueIndirect(v) {
		m["serverAuth"] = v
	}
	if v := f.ClientAuth; !dcl.IsEmptyValueIndirect(v) {
		m["clientAuth"] = v
	}
	if v := f.CodeSigning; !dcl.IsEmptyValueIndirect(v) {
		m["codeSigning"] = v
	}
	if v := f.EmailProtection; !dcl.IsEmptyValueIndirect(v) {
		m["emailProtection"] = v
	}
	if v := f.TimeStamping; !dcl.IsEmptyValueIndirect(v) {
		m["timeStamping"] = v
	}
	if v := f.OcspSigning; !dcl.IsEmptyValueIndirect(v) {
		m["ocspSigning"] = v
	}

	return m, nil
}

// flattenCaPoolIssuancePolicyBaselineValuesKeyUsageExtendedKeyUsage flattens an instance of CaPoolIssuancePolicyBaselineValuesKeyUsageExtendedKeyUsage from a JSON
// response object.
func flattenCaPoolIssuancePolicyBaselineValuesKeyUsageExtendedKeyUsage(c *Client, i interface{}, res *CaPool) *CaPoolIssuancePolicyBaselineValuesKeyUsageExtendedKeyUsage {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &CaPoolIssuancePolicyBaselineValuesKeyUsageExtendedKeyUsage{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyCaPoolIssuancePolicyBaselineValuesKeyUsageExtendedKeyUsage
	}
	r.ServerAuth = dcl.FlattenBool(m["serverAuth"])
	r.ClientAuth = dcl.FlattenBool(m["clientAuth"])
	r.CodeSigning = dcl.FlattenBool(m["codeSigning"])
	r.EmailProtection = dcl.FlattenBool(m["emailProtection"])
	r.TimeStamping = dcl.FlattenBool(m["timeStamping"])
	r.OcspSigning = dcl.FlattenBool(m["ocspSigning"])

	return r
}

// expandCaPoolIssuancePolicyBaselineValuesKeyUsageUnknownExtendedKeyUsagesMap expands the contents of CaPoolIssuancePolicyBaselineValuesKeyUsageUnknownExtendedKeyUsages into a JSON
// request object.
func expandCaPoolIssuancePolicyBaselineValuesKeyUsageUnknownExtendedKeyUsagesMap(c *Client, f map[string]CaPoolIssuancePolicyBaselineValuesKeyUsageUnknownExtendedKeyUsages, res *CaPool) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandCaPoolIssuancePolicyBaselineValuesKeyUsageUnknownExtendedKeyUsages(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandCaPoolIssuancePolicyBaselineValuesKeyUsageUnknownExtendedKeyUsagesSlice expands the contents of CaPoolIssuancePolicyBaselineValuesKeyUsageUnknownExtendedKeyUsages into a JSON
// request object.
func expandCaPoolIssuancePolicyBaselineValuesKeyUsageUnknownExtendedKeyUsagesSlice(c *Client, f []CaPoolIssuancePolicyBaselineValuesKeyUsageUnknownExtendedKeyUsages, res *CaPool) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandCaPoolIssuancePolicyBaselineValuesKeyUsageUnknownExtendedKeyUsages(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenCaPoolIssuancePolicyBaselineValuesKeyUsageUnknownExtendedKeyUsagesMap flattens the contents of CaPoolIssuancePolicyBaselineValuesKeyUsageUnknownExtendedKeyUsages from a JSON
// response object.
func flattenCaPoolIssuancePolicyBaselineValuesKeyUsageUnknownExtendedKeyUsagesMap(c *Client, i interface{}, res *CaPool) map[string]CaPoolIssuancePolicyBaselineValuesKeyUsageUnknownExtendedKeyUsages {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]CaPoolIssuancePolicyBaselineValuesKeyUsageUnknownExtendedKeyUsages{}
	}

	if len(a) == 0 {
		return map[string]CaPoolIssuancePolicyBaselineValuesKeyUsageUnknownExtendedKeyUsages{}
	}

	items := make(map[string]CaPoolIssuancePolicyBaselineValuesKeyUsageUnknownExtendedKeyUsages)
	for k, item := range a {
		items[k] = *flattenCaPoolIssuancePolicyBaselineValuesKeyUsageUnknownExtendedKeyUsages(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenCaPoolIssuancePolicyBaselineValuesKeyUsageUnknownExtendedKeyUsagesSlice flattens the contents of CaPoolIssuancePolicyBaselineValuesKeyUsageUnknownExtendedKeyUsages from a JSON
// response object.
func flattenCaPoolIssuancePolicyBaselineValuesKeyUsageUnknownExtendedKeyUsagesSlice(c *Client, i interface{}, res *CaPool) []CaPoolIssuancePolicyBaselineValuesKeyUsageUnknownExtendedKeyUsages {
	a, ok := i.([]interface{})
	if !ok {
		return []CaPoolIssuancePolicyBaselineValuesKeyUsageUnknownExtendedKeyUsages{}
	}

	if len(a) == 0 {
		return []CaPoolIssuancePolicyBaselineValuesKeyUsageUnknownExtendedKeyUsages{}
	}

	items := make([]CaPoolIssuancePolicyBaselineValuesKeyUsageUnknownExtendedKeyUsages, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenCaPoolIssuancePolicyBaselineValuesKeyUsageUnknownExtendedKeyUsages(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandCaPoolIssuancePolicyBaselineValuesKeyUsageUnknownExtendedKeyUsages expands an instance of CaPoolIssuancePolicyBaselineValuesKeyUsageUnknownExtendedKeyUsages into a JSON
// request object.
func expandCaPoolIssuancePolicyBaselineValuesKeyUsageUnknownExtendedKeyUsages(c *Client, f *CaPoolIssuancePolicyBaselineValuesKeyUsageUnknownExtendedKeyUsages, res *CaPool) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.ObjectIdPath; v != nil {
		m["objectIdPath"] = v
	}

	return m, nil
}

// flattenCaPoolIssuancePolicyBaselineValuesKeyUsageUnknownExtendedKeyUsages flattens an instance of CaPoolIssuancePolicyBaselineValuesKeyUsageUnknownExtendedKeyUsages from a JSON
// response object.
func flattenCaPoolIssuancePolicyBaselineValuesKeyUsageUnknownExtendedKeyUsages(c *Client, i interface{}, res *CaPool) *CaPoolIssuancePolicyBaselineValuesKeyUsageUnknownExtendedKeyUsages {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &CaPoolIssuancePolicyBaselineValuesKeyUsageUnknownExtendedKeyUsages{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyCaPoolIssuancePolicyBaselineValuesKeyUsageUnknownExtendedKeyUsages
	}
	r.ObjectIdPath = dcl.FlattenIntSlice(m["objectIdPath"])

	return r
}

// expandCaPoolIssuancePolicyBaselineValuesCaOptionsMap expands the contents of CaPoolIssuancePolicyBaselineValuesCaOptions into a JSON
// request object.
func expandCaPoolIssuancePolicyBaselineValuesCaOptionsMap(c *Client, f map[string]CaPoolIssuancePolicyBaselineValuesCaOptions, res *CaPool) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandCaPoolIssuancePolicyBaselineValuesCaOptions(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandCaPoolIssuancePolicyBaselineValuesCaOptionsSlice expands the contents of CaPoolIssuancePolicyBaselineValuesCaOptions into a JSON
// request object.
func expandCaPoolIssuancePolicyBaselineValuesCaOptionsSlice(c *Client, f []CaPoolIssuancePolicyBaselineValuesCaOptions, res *CaPool) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandCaPoolIssuancePolicyBaselineValuesCaOptions(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenCaPoolIssuancePolicyBaselineValuesCaOptionsMap flattens the contents of CaPoolIssuancePolicyBaselineValuesCaOptions from a JSON
// response object.
func flattenCaPoolIssuancePolicyBaselineValuesCaOptionsMap(c *Client, i interface{}, res *CaPool) map[string]CaPoolIssuancePolicyBaselineValuesCaOptions {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]CaPoolIssuancePolicyBaselineValuesCaOptions{}
	}

	if len(a) == 0 {
		return map[string]CaPoolIssuancePolicyBaselineValuesCaOptions{}
	}

	items := make(map[string]CaPoolIssuancePolicyBaselineValuesCaOptions)
	for k, item := range a {
		items[k] = *flattenCaPoolIssuancePolicyBaselineValuesCaOptions(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenCaPoolIssuancePolicyBaselineValuesCaOptionsSlice flattens the contents of CaPoolIssuancePolicyBaselineValuesCaOptions from a JSON
// response object.
func flattenCaPoolIssuancePolicyBaselineValuesCaOptionsSlice(c *Client, i interface{}, res *CaPool) []CaPoolIssuancePolicyBaselineValuesCaOptions {
	a, ok := i.([]interface{})
	if !ok {
		return []CaPoolIssuancePolicyBaselineValuesCaOptions{}
	}

	if len(a) == 0 {
		return []CaPoolIssuancePolicyBaselineValuesCaOptions{}
	}

	items := make([]CaPoolIssuancePolicyBaselineValuesCaOptions, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenCaPoolIssuancePolicyBaselineValuesCaOptions(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandCaPoolIssuancePolicyBaselineValuesCaOptions expands an instance of CaPoolIssuancePolicyBaselineValuesCaOptions into a JSON
// request object.
func expandCaPoolIssuancePolicyBaselineValuesCaOptions(c *Client, f *CaPoolIssuancePolicyBaselineValuesCaOptions, res *CaPool) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.IsCa; !dcl.IsEmptyValueIndirect(v) {
		m["isCa"] = v
	}
	if v := f.MaxIssuerPathLength; !dcl.IsEmptyValueIndirect(v) {
		m["maxIssuerPathLength"] = v
	}
	if v := f.ZeroMaxIssuerPathLength; !dcl.IsEmptyValueIndirect(v) {
		m["zeroMaxIssuerPathLength"] = v
	}

	return m, nil
}

// flattenCaPoolIssuancePolicyBaselineValuesCaOptions flattens an instance of CaPoolIssuancePolicyBaselineValuesCaOptions from a JSON
// response object.
func flattenCaPoolIssuancePolicyBaselineValuesCaOptions(c *Client, i interface{}, res *CaPool) *CaPoolIssuancePolicyBaselineValuesCaOptions {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &CaPoolIssuancePolicyBaselineValuesCaOptions{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyCaPoolIssuancePolicyBaselineValuesCaOptions
	}
	r.IsCa = dcl.FlattenBool(m["isCa"])
	r.MaxIssuerPathLength = dcl.FlattenInteger(m["maxIssuerPathLength"])
	r.ZeroMaxIssuerPathLength = dcl.FlattenBool(m["zeroMaxIssuerPathLength"])

	return r
}

// expandCaPoolIssuancePolicyBaselineValuesPolicyIdsMap expands the contents of CaPoolIssuancePolicyBaselineValuesPolicyIds into a JSON
// request object.
func expandCaPoolIssuancePolicyBaselineValuesPolicyIdsMap(c *Client, f map[string]CaPoolIssuancePolicyBaselineValuesPolicyIds, res *CaPool) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandCaPoolIssuancePolicyBaselineValuesPolicyIds(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandCaPoolIssuancePolicyBaselineValuesPolicyIdsSlice expands the contents of CaPoolIssuancePolicyBaselineValuesPolicyIds into a JSON
// request object.
func expandCaPoolIssuancePolicyBaselineValuesPolicyIdsSlice(c *Client, f []CaPoolIssuancePolicyBaselineValuesPolicyIds, res *CaPool) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandCaPoolIssuancePolicyBaselineValuesPolicyIds(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenCaPoolIssuancePolicyBaselineValuesPolicyIdsMap flattens the contents of CaPoolIssuancePolicyBaselineValuesPolicyIds from a JSON
// response object.
func flattenCaPoolIssuancePolicyBaselineValuesPolicyIdsMap(c *Client, i interface{}, res *CaPool) map[string]CaPoolIssuancePolicyBaselineValuesPolicyIds {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]CaPoolIssuancePolicyBaselineValuesPolicyIds{}
	}

	if len(a) == 0 {
		return map[string]CaPoolIssuancePolicyBaselineValuesPolicyIds{}
	}

	items := make(map[string]CaPoolIssuancePolicyBaselineValuesPolicyIds)
	for k, item := range a {
		items[k] = *flattenCaPoolIssuancePolicyBaselineValuesPolicyIds(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenCaPoolIssuancePolicyBaselineValuesPolicyIdsSlice flattens the contents of CaPoolIssuancePolicyBaselineValuesPolicyIds from a JSON
// response object.
func flattenCaPoolIssuancePolicyBaselineValuesPolicyIdsSlice(c *Client, i interface{}, res *CaPool) []CaPoolIssuancePolicyBaselineValuesPolicyIds {
	a, ok := i.([]interface{})
	if !ok {
		return []CaPoolIssuancePolicyBaselineValuesPolicyIds{}
	}

	if len(a) == 0 {
		return []CaPoolIssuancePolicyBaselineValuesPolicyIds{}
	}

	items := make([]CaPoolIssuancePolicyBaselineValuesPolicyIds, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenCaPoolIssuancePolicyBaselineValuesPolicyIds(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandCaPoolIssuancePolicyBaselineValuesPolicyIds expands an instance of CaPoolIssuancePolicyBaselineValuesPolicyIds into a JSON
// request object.
func expandCaPoolIssuancePolicyBaselineValuesPolicyIds(c *Client, f *CaPoolIssuancePolicyBaselineValuesPolicyIds, res *CaPool) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.ObjectIdPath; v != nil {
		m["objectIdPath"] = v
	}

	return m, nil
}

// flattenCaPoolIssuancePolicyBaselineValuesPolicyIds flattens an instance of CaPoolIssuancePolicyBaselineValuesPolicyIds from a JSON
// response object.
func flattenCaPoolIssuancePolicyBaselineValuesPolicyIds(c *Client, i interface{}, res *CaPool) *CaPoolIssuancePolicyBaselineValuesPolicyIds {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &CaPoolIssuancePolicyBaselineValuesPolicyIds{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyCaPoolIssuancePolicyBaselineValuesPolicyIds
	}
	r.ObjectIdPath = dcl.FlattenIntSlice(m["objectIdPath"])

	return r
}

// expandCaPoolIssuancePolicyBaselineValuesAdditionalExtensionsMap expands the contents of CaPoolIssuancePolicyBaselineValuesAdditionalExtensions into a JSON
// request object.
func expandCaPoolIssuancePolicyBaselineValuesAdditionalExtensionsMap(c *Client, f map[string]CaPoolIssuancePolicyBaselineValuesAdditionalExtensions, res *CaPool) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandCaPoolIssuancePolicyBaselineValuesAdditionalExtensions(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandCaPoolIssuancePolicyBaselineValuesAdditionalExtensionsSlice expands the contents of CaPoolIssuancePolicyBaselineValuesAdditionalExtensions into a JSON
// request object.
func expandCaPoolIssuancePolicyBaselineValuesAdditionalExtensionsSlice(c *Client, f []CaPoolIssuancePolicyBaselineValuesAdditionalExtensions, res *CaPool) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandCaPoolIssuancePolicyBaselineValuesAdditionalExtensions(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenCaPoolIssuancePolicyBaselineValuesAdditionalExtensionsMap flattens the contents of CaPoolIssuancePolicyBaselineValuesAdditionalExtensions from a JSON
// response object.
func flattenCaPoolIssuancePolicyBaselineValuesAdditionalExtensionsMap(c *Client, i interface{}, res *CaPool) map[string]CaPoolIssuancePolicyBaselineValuesAdditionalExtensions {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]CaPoolIssuancePolicyBaselineValuesAdditionalExtensions{}
	}

	if len(a) == 0 {
		return map[string]CaPoolIssuancePolicyBaselineValuesAdditionalExtensions{}
	}

	items := make(map[string]CaPoolIssuancePolicyBaselineValuesAdditionalExtensions)
	for k, item := range a {
		items[k] = *flattenCaPoolIssuancePolicyBaselineValuesAdditionalExtensions(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenCaPoolIssuancePolicyBaselineValuesAdditionalExtensionsSlice flattens the contents of CaPoolIssuancePolicyBaselineValuesAdditionalExtensions from a JSON
// response object.
func flattenCaPoolIssuancePolicyBaselineValuesAdditionalExtensionsSlice(c *Client, i interface{}, res *CaPool) []CaPoolIssuancePolicyBaselineValuesAdditionalExtensions {
	a, ok := i.([]interface{})
	if !ok {
		return []CaPoolIssuancePolicyBaselineValuesAdditionalExtensions{}
	}

	if len(a) == 0 {
		return []CaPoolIssuancePolicyBaselineValuesAdditionalExtensions{}
	}

	items := make([]CaPoolIssuancePolicyBaselineValuesAdditionalExtensions, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenCaPoolIssuancePolicyBaselineValuesAdditionalExtensions(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandCaPoolIssuancePolicyBaselineValuesAdditionalExtensions expands an instance of CaPoolIssuancePolicyBaselineValuesAdditionalExtensions into a JSON
// request object.
func expandCaPoolIssuancePolicyBaselineValuesAdditionalExtensions(c *Client, f *CaPoolIssuancePolicyBaselineValuesAdditionalExtensions, res *CaPool) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v, err := expandCaPoolIssuancePolicyBaselineValuesAdditionalExtensionsObjectId(c, f.ObjectId, res); err != nil {
		return nil, fmt.Errorf("error expanding ObjectId into objectId: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["objectId"] = v
	}
	if v := f.Critical; !dcl.IsEmptyValueIndirect(v) {
		m["critical"] = v
	}
	if v := f.Value; !dcl.IsEmptyValueIndirect(v) {
		m["value"] = v
	}

	return m, nil
}

// flattenCaPoolIssuancePolicyBaselineValuesAdditionalExtensions flattens an instance of CaPoolIssuancePolicyBaselineValuesAdditionalExtensions from a JSON
// response object.
func flattenCaPoolIssuancePolicyBaselineValuesAdditionalExtensions(c *Client, i interface{}, res *CaPool) *CaPoolIssuancePolicyBaselineValuesAdditionalExtensions {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &CaPoolIssuancePolicyBaselineValuesAdditionalExtensions{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyCaPoolIssuancePolicyBaselineValuesAdditionalExtensions
	}
	r.ObjectId = flattenCaPoolIssuancePolicyBaselineValuesAdditionalExtensionsObjectId(c, m["objectId"], res)
	r.Critical = dcl.FlattenBool(m["critical"])
	r.Value = dcl.FlattenString(m["value"])

	return r
}

// expandCaPoolIssuancePolicyBaselineValuesAdditionalExtensionsObjectIdMap expands the contents of CaPoolIssuancePolicyBaselineValuesAdditionalExtensionsObjectId into a JSON
// request object.
func expandCaPoolIssuancePolicyBaselineValuesAdditionalExtensionsObjectIdMap(c *Client, f map[string]CaPoolIssuancePolicyBaselineValuesAdditionalExtensionsObjectId, res *CaPool) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandCaPoolIssuancePolicyBaselineValuesAdditionalExtensionsObjectId(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandCaPoolIssuancePolicyBaselineValuesAdditionalExtensionsObjectIdSlice expands the contents of CaPoolIssuancePolicyBaselineValuesAdditionalExtensionsObjectId into a JSON
// request object.
func expandCaPoolIssuancePolicyBaselineValuesAdditionalExtensionsObjectIdSlice(c *Client, f []CaPoolIssuancePolicyBaselineValuesAdditionalExtensionsObjectId, res *CaPool) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandCaPoolIssuancePolicyBaselineValuesAdditionalExtensionsObjectId(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenCaPoolIssuancePolicyBaselineValuesAdditionalExtensionsObjectIdMap flattens the contents of CaPoolIssuancePolicyBaselineValuesAdditionalExtensionsObjectId from a JSON
// response object.
func flattenCaPoolIssuancePolicyBaselineValuesAdditionalExtensionsObjectIdMap(c *Client, i interface{}, res *CaPool) map[string]CaPoolIssuancePolicyBaselineValuesAdditionalExtensionsObjectId {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]CaPoolIssuancePolicyBaselineValuesAdditionalExtensionsObjectId{}
	}

	if len(a) == 0 {
		return map[string]CaPoolIssuancePolicyBaselineValuesAdditionalExtensionsObjectId{}
	}

	items := make(map[string]CaPoolIssuancePolicyBaselineValuesAdditionalExtensionsObjectId)
	for k, item := range a {
		items[k] = *flattenCaPoolIssuancePolicyBaselineValuesAdditionalExtensionsObjectId(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenCaPoolIssuancePolicyBaselineValuesAdditionalExtensionsObjectIdSlice flattens the contents of CaPoolIssuancePolicyBaselineValuesAdditionalExtensionsObjectId from a JSON
// response object.
func flattenCaPoolIssuancePolicyBaselineValuesAdditionalExtensionsObjectIdSlice(c *Client, i interface{}, res *CaPool) []CaPoolIssuancePolicyBaselineValuesAdditionalExtensionsObjectId {
	a, ok := i.([]interface{})
	if !ok {
		return []CaPoolIssuancePolicyBaselineValuesAdditionalExtensionsObjectId{}
	}

	if len(a) == 0 {
		return []CaPoolIssuancePolicyBaselineValuesAdditionalExtensionsObjectId{}
	}

	items := make([]CaPoolIssuancePolicyBaselineValuesAdditionalExtensionsObjectId, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenCaPoolIssuancePolicyBaselineValuesAdditionalExtensionsObjectId(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandCaPoolIssuancePolicyBaselineValuesAdditionalExtensionsObjectId expands an instance of CaPoolIssuancePolicyBaselineValuesAdditionalExtensionsObjectId into a JSON
// request object.
func expandCaPoolIssuancePolicyBaselineValuesAdditionalExtensionsObjectId(c *Client, f *CaPoolIssuancePolicyBaselineValuesAdditionalExtensionsObjectId, res *CaPool) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.ObjectIdPath; v != nil {
		m["objectIdPath"] = v
	}

	return m, nil
}

// flattenCaPoolIssuancePolicyBaselineValuesAdditionalExtensionsObjectId flattens an instance of CaPoolIssuancePolicyBaselineValuesAdditionalExtensionsObjectId from a JSON
// response object.
func flattenCaPoolIssuancePolicyBaselineValuesAdditionalExtensionsObjectId(c *Client, i interface{}, res *CaPool) *CaPoolIssuancePolicyBaselineValuesAdditionalExtensionsObjectId {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &CaPoolIssuancePolicyBaselineValuesAdditionalExtensionsObjectId{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyCaPoolIssuancePolicyBaselineValuesAdditionalExtensionsObjectId
	}
	r.ObjectIdPath = dcl.FlattenIntSlice(m["objectIdPath"])

	return r
}

// expandCaPoolIssuancePolicyIdentityConstraintsMap expands the contents of CaPoolIssuancePolicyIdentityConstraints into a JSON
// request object.
func expandCaPoolIssuancePolicyIdentityConstraintsMap(c *Client, f map[string]CaPoolIssuancePolicyIdentityConstraints, res *CaPool) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandCaPoolIssuancePolicyIdentityConstraints(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandCaPoolIssuancePolicyIdentityConstraintsSlice expands the contents of CaPoolIssuancePolicyIdentityConstraints into a JSON
// request object.
func expandCaPoolIssuancePolicyIdentityConstraintsSlice(c *Client, f []CaPoolIssuancePolicyIdentityConstraints, res *CaPool) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandCaPoolIssuancePolicyIdentityConstraints(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenCaPoolIssuancePolicyIdentityConstraintsMap flattens the contents of CaPoolIssuancePolicyIdentityConstraints from a JSON
// response object.
func flattenCaPoolIssuancePolicyIdentityConstraintsMap(c *Client, i interface{}, res *CaPool) map[string]CaPoolIssuancePolicyIdentityConstraints {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]CaPoolIssuancePolicyIdentityConstraints{}
	}

	if len(a) == 0 {
		return map[string]CaPoolIssuancePolicyIdentityConstraints{}
	}

	items := make(map[string]CaPoolIssuancePolicyIdentityConstraints)
	for k, item := range a {
		items[k] = *flattenCaPoolIssuancePolicyIdentityConstraints(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenCaPoolIssuancePolicyIdentityConstraintsSlice flattens the contents of CaPoolIssuancePolicyIdentityConstraints from a JSON
// response object.
func flattenCaPoolIssuancePolicyIdentityConstraintsSlice(c *Client, i interface{}, res *CaPool) []CaPoolIssuancePolicyIdentityConstraints {
	a, ok := i.([]interface{})
	if !ok {
		return []CaPoolIssuancePolicyIdentityConstraints{}
	}

	if len(a) == 0 {
		return []CaPoolIssuancePolicyIdentityConstraints{}
	}

	items := make([]CaPoolIssuancePolicyIdentityConstraints, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenCaPoolIssuancePolicyIdentityConstraints(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandCaPoolIssuancePolicyIdentityConstraints expands an instance of CaPoolIssuancePolicyIdentityConstraints into a JSON
// request object.
func expandCaPoolIssuancePolicyIdentityConstraints(c *Client, f *CaPoolIssuancePolicyIdentityConstraints, res *CaPool) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v, err := expandCaPoolIssuancePolicyIdentityConstraintsCelExpression(c, f.CelExpression, res); err != nil {
		return nil, fmt.Errorf("error expanding CelExpression into celExpression: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["celExpression"] = v
	}
	if v := f.AllowSubjectPassthrough; !dcl.IsEmptyValueIndirect(v) {
		m["allowSubjectPassthrough"] = v
	}
	if v := f.AllowSubjectAltNamesPassthrough; !dcl.IsEmptyValueIndirect(v) {
		m["allowSubjectAltNamesPassthrough"] = v
	}

	return m, nil
}

// flattenCaPoolIssuancePolicyIdentityConstraints flattens an instance of CaPoolIssuancePolicyIdentityConstraints from a JSON
// response object.
func flattenCaPoolIssuancePolicyIdentityConstraints(c *Client, i interface{}, res *CaPool) *CaPoolIssuancePolicyIdentityConstraints {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &CaPoolIssuancePolicyIdentityConstraints{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyCaPoolIssuancePolicyIdentityConstraints
	}
	r.CelExpression = flattenCaPoolIssuancePolicyIdentityConstraintsCelExpression(c, m["celExpression"], res)
	r.AllowSubjectPassthrough = dcl.FlattenBool(m["allowSubjectPassthrough"])
	r.AllowSubjectAltNamesPassthrough = dcl.FlattenBool(m["allowSubjectAltNamesPassthrough"])

	return r
}

// expandCaPoolIssuancePolicyIdentityConstraintsCelExpressionMap expands the contents of CaPoolIssuancePolicyIdentityConstraintsCelExpression into a JSON
// request object.
func expandCaPoolIssuancePolicyIdentityConstraintsCelExpressionMap(c *Client, f map[string]CaPoolIssuancePolicyIdentityConstraintsCelExpression, res *CaPool) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandCaPoolIssuancePolicyIdentityConstraintsCelExpression(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandCaPoolIssuancePolicyIdentityConstraintsCelExpressionSlice expands the contents of CaPoolIssuancePolicyIdentityConstraintsCelExpression into a JSON
// request object.
func expandCaPoolIssuancePolicyIdentityConstraintsCelExpressionSlice(c *Client, f []CaPoolIssuancePolicyIdentityConstraintsCelExpression, res *CaPool) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandCaPoolIssuancePolicyIdentityConstraintsCelExpression(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenCaPoolIssuancePolicyIdentityConstraintsCelExpressionMap flattens the contents of CaPoolIssuancePolicyIdentityConstraintsCelExpression from a JSON
// response object.
func flattenCaPoolIssuancePolicyIdentityConstraintsCelExpressionMap(c *Client, i interface{}, res *CaPool) map[string]CaPoolIssuancePolicyIdentityConstraintsCelExpression {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]CaPoolIssuancePolicyIdentityConstraintsCelExpression{}
	}

	if len(a) == 0 {
		return map[string]CaPoolIssuancePolicyIdentityConstraintsCelExpression{}
	}

	items := make(map[string]CaPoolIssuancePolicyIdentityConstraintsCelExpression)
	for k, item := range a {
		items[k] = *flattenCaPoolIssuancePolicyIdentityConstraintsCelExpression(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenCaPoolIssuancePolicyIdentityConstraintsCelExpressionSlice flattens the contents of CaPoolIssuancePolicyIdentityConstraintsCelExpression from a JSON
// response object.
func flattenCaPoolIssuancePolicyIdentityConstraintsCelExpressionSlice(c *Client, i interface{}, res *CaPool) []CaPoolIssuancePolicyIdentityConstraintsCelExpression {
	a, ok := i.([]interface{})
	if !ok {
		return []CaPoolIssuancePolicyIdentityConstraintsCelExpression{}
	}

	if len(a) == 0 {
		return []CaPoolIssuancePolicyIdentityConstraintsCelExpression{}
	}

	items := make([]CaPoolIssuancePolicyIdentityConstraintsCelExpression, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenCaPoolIssuancePolicyIdentityConstraintsCelExpression(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandCaPoolIssuancePolicyIdentityConstraintsCelExpression expands an instance of CaPoolIssuancePolicyIdentityConstraintsCelExpression into a JSON
// request object.
func expandCaPoolIssuancePolicyIdentityConstraintsCelExpression(c *Client, f *CaPoolIssuancePolicyIdentityConstraintsCelExpression, res *CaPool) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.Expression; !dcl.IsEmptyValueIndirect(v) {
		m["expression"] = v
	}
	if v := f.Title; !dcl.IsEmptyValueIndirect(v) {
		m["title"] = v
	}
	if v := f.Description; !dcl.IsEmptyValueIndirect(v) {
		m["description"] = v
	}
	if v := f.Location; !dcl.IsEmptyValueIndirect(v) {
		m["location"] = v
	}

	return m, nil
}

// flattenCaPoolIssuancePolicyIdentityConstraintsCelExpression flattens an instance of CaPoolIssuancePolicyIdentityConstraintsCelExpression from a JSON
// response object.
func flattenCaPoolIssuancePolicyIdentityConstraintsCelExpression(c *Client, i interface{}, res *CaPool) *CaPoolIssuancePolicyIdentityConstraintsCelExpression {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &CaPoolIssuancePolicyIdentityConstraintsCelExpression{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyCaPoolIssuancePolicyIdentityConstraintsCelExpression
	}
	r.Expression = dcl.FlattenString(m["expression"])
	r.Title = dcl.FlattenString(m["title"])
	r.Description = dcl.FlattenString(m["description"])
	r.Location = dcl.FlattenString(m["location"])

	return r
}

// expandCaPoolIssuancePolicyPassthroughExtensionsMap expands the contents of CaPoolIssuancePolicyPassthroughExtensions into a JSON
// request object.
func expandCaPoolIssuancePolicyPassthroughExtensionsMap(c *Client, f map[string]CaPoolIssuancePolicyPassthroughExtensions, res *CaPool) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandCaPoolIssuancePolicyPassthroughExtensions(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandCaPoolIssuancePolicyPassthroughExtensionsSlice expands the contents of CaPoolIssuancePolicyPassthroughExtensions into a JSON
// request object.
func expandCaPoolIssuancePolicyPassthroughExtensionsSlice(c *Client, f []CaPoolIssuancePolicyPassthroughExtensions, res *CaPool) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandCaPoolIssuancePolicyPassthroughExtensions(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenCaPoolIssuancePolicyPassthroughExtensionsMap flattens the contents of CaPoolIssuancePolicyPassthroughExtensions from a JSON
// response object.
func flattenCaPoolIssuancePolicyPassthroughExtensionsMap(c *Client, i interface{}, res *CaPool) map[string]CaPoolIssuancePolicyPassthroughExtensions {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]CaPoolIssuancePolicyPassthroughExtensions{}
	}

	if len(a) == 0 {
		return map[string]CaPoolIssuancePolicyPassthroughExtensions{}
	}

	items := make(map[string]CaPoolIssuancePolicyPassthroughExtensions)
	for k, item := range a {
		items[k] = *flattenCaPoolIssuancePolicyPassthroughExtensions(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenCaPoolIssuancePolicyPassthroughExtensionsSlice flattens the contents of CaPoolIssuancePolicyPassthroughExtensions from a JSON
// response object.
func flattenCaPoolIssuancePolicyPassthroughExtensionsSlice(c *Client, i interface{}, res *CaPool) []CaPoolIssuancePolicyPassthroughExtensions {
	a, ok := i.([]interface{})
	if !ok {
		return []CaPoolIssuancePolicyPassthroughExtensions{}
	}

	if len(a) == 0 {
		return []CaPoolIssuancePolicyPassthroughExtensions{}
	}

	items := make([]CaPoolIssuancePolicyPassthroughExtensions, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenCaPoolIssuancePolicyPassthroughExtensions(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandCaPoolIssuancePolicyPassthroughExtensions expands an instance of CaPoolIssuancePolicyPassthroughExtensions into a JSON
// request object.
func expandCaPoolIssuancePolicyPassthroughExtensions(c *Client, f *CaPoolIssuancePolicyPassthroughExtensions, res *CaPool) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.KnownExtensions; v != nil {
		m["knownExtensions"] = v
	}
	if v, err := expandCaPoolIssuancePolicyPassthroughExtensionsAdditionalExtensionsSlice(c, f.AdditionalExtensions, res); err != nil {
		return nil, fmt.Errorf("error expanding AdditionalExtensions into additionalExtensions: %w", err)
	} else if v != nil {
		m["additionalExtensions"] = v
	}

	return m, nil
}

// flattenCaPoolIssuancePolicyPassthroughExtensions flattens an instance of CaPoolIssuancePolicyPassthroughExtensions from a JSON
// response object.
func flattenCaPoolIssuancePolicyPassthroughExtensions(c *Client, i interface{}, res *CaPool) *CaPoolIssuancePolicyPassthroughExtensions {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &CaPoolIssuancePolicyPassthroughExtensions{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyCaPoolIssuancePolicyPassthroughExtensions
	}
	r.KnownExtensions = flattenCaPoolIssuancePolicyPassthroughExtensionsKnownExtensionsEnumSlice(c, m["knownExtensions"], res)
	r.AdditionalExtensions = flattenCaPoolIssuancePolicyPassthroughExtensionsAdditionalExtensionsSlice(c, m["additionalExtensions"], res)

	return r
}

// expandCaPoolIssuancePolicyPassthroughExtensionsAdditionalExtensionsMap expands the contents of CaPoolIssuancePolicyPassthroughExtensionsAdditionalExtensions into a JSON
// request object.
func expandCaPoolIssuancePolicyPassthroughExtensionsAdditionalExtensionsMap(c *Client, f map[string]CaPoolIssuancePolicyPassthroughExtensionsAdditionalExtensions, res *CaPool) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandCaPoolIssuancePolicyPassthroughExtensionsAdditionalExtensions(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandCaPoolIssuancePolicyPassthroughExtensionsAdditionalExtensionsSlice expands the contents of CaPoolIssuancePolicyPassthroughExtensionsAdditionalExtensions into a JSON
// request object.
func expandCaPoolIssuancePolicyPassthroughExtensionsAdditionalExtensionsSlice(c *Client, f []CaPoolIssuancePolicyPassthroughExtensionsAdditionalExtensions, res *CaPool) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandCaPoolIssuancePolicyPassthroughExtensionsAdditionalExtensions(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenCaPoolIssuancePolicyPassthroughExtensionsAdditionalExtensionsMap flattens the contents of CaPoolIssuancePolicyPassthroughExtensionsAdditionalExtensions from a JSON
// response object.
func flattenCaPoolIssuancePolicyPassthroughExtensionsAdditionalExtensionsMap(c *Client, i interface{}, res *CaPool) map[string]CaPoolIssuancePolicyPassthroughExtensionsAdditionalExtensions {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]CaPoolIssuancePolicyPassthroughExtensionsAdditionalExtensions{}
	}

	if len(a) == 0 {
		return map[string]CaPoolIssuancePolicyPassthroughExtensionsAdditionalExtensions{}
	}

	items := make(map[string]CaPoolIssuancePolicyPassthroughExtensionsAdditionalExtensions)
	for k, item := range a {
		items[k] = *flattenCaPoolIssuancePolicyPassthroughExtensionsAdditionalExtensions(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenCaPoolIssuancePolicyPassthroughExtensionsAdditionalExtensionsSlice flattens the contents of CaPoolIssuancePolicyPassthroughExtensionsAdditionalExtensions from a JSON
// response object.
func flattenCaPoolIssuancePolicyPassthroughExtensionsAdditionalExtensionsSlice(c *Client, i interface{}, res *CaPool) []CaPoolIssuancePolicyPassthroughExtensionsAdditionalExtensions {
	a, ok := i.([]interface{})
	if !ok {
		return []CaPoolIssuancePolicyPassthroughExtensionsAdditionalExtensions{}
	}

	if len(a) == 0 {
		return []CaPoolIssuancePolicyPassthroughExtensionsAdditionalExtensions{}
	}

	items := make([]CaPoolIssuancePolicyPassthroughExtensionsAdditionalExtensions, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenCaPoolIssuancePolicyPassthroughExtensionsAdditionalExtensions(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandCaPoolIssuancePolicyPassthroughExtensionsAdditionalExtensions expands an instance of CaPoolIssuancePolicyPassthroughExtensionsAdditionalExtensions into a JSON
// request object.
func expandCaPoolIssuancePolicyPassthroughExtensionsAdditionalExtensions(c *Client, f *CaPoolIssuancePolicyPassthroughExtensionsAdditionalExtensions, res *CaPool) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.ObjectIdPath; v != nil {
		m["objectIdPath"] = v
	}

	return m, nil
}

// flattenCaPoolIssuancePolicyPassthroughExtensionsAdditionalExtensions flattens an instance of CaPoolIssuancePolicyPassthroughExtensionsAdditionalExtensions from a JSON
// response object.
func flattenCaPoolIssuancePolicyPassthroughExtensionsAdditionalExtensions(c *Client, i interface{}, res *CaPool) *CaPoolIssuancePolicyPassthroughExtensionsAdditionalExtensions {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &CaPoolIssuancePolicyPassthroughExtensionsAdditionalExtensions{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyCaPoolIssuancePolicyPassthroughExtensionsAdditionalExtensions
	}
	r.ObjectIdPath = dcl.FlattenIntSlice(m["objectIdPath"])

	return r
}

// expandCaPoolPublishingOptionsMap expands the contents of CaPoolPublishingOptions into a JSON
// request object.
func expandCaPoolPublishingOptionsMap(c *Client, f map[string]CaPoolPublishingOptions, res *CaPool) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandCaPoolPublishingOptions(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandCaPoolPublishingOptionsSlice expands the contents of CaPoolPublishingOptions into a JSON
// request object.
func expandCaPoolPublishingOptionsSlice(c *Client, f []CaPoolPublishingOptions, res *CaPool) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandCaPoolPublishingOptions(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenCaPoolPublishingOptionsMap flattens the contents of CaPoolPublishingOptions from a JSON
// response object.
func flattenCaPoolPublishingOptionsMap(c *Client, i interface{}, res *CaPool) map[string]CaPoolPublishingOptions {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]CaPoolPublishingOptions{}
	}

	if len(a) == 0 {
		return map[string]CaPoolPublishingOptions{}
	}

	items := make(map[string]CaPoolPublishingOptions)
	for k, item := range a {
		items[k] = *flattenCaPoolPublishingOptions(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenCaPoolPublishingOptionsSlice flattens the contents of CaPoolPublishingOptions from a JSON
// response object.
func flattenCaPoolPublishingOptionsSlice(c *Client, i interface{}, res *CaPool) []CaPoolPublishingOptions {
	a, ok := i.([]interface{})
	if !ok {
		return []CaPoolPublishingOptions{}
	}

	if len(a) == 0 {
		return []CaPoolPublishingOptions{}
	}

	items := make([]CaPoolPublishingOptions, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenCaPoolPublishingOptions(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandCaPoolPublishingOptions expands an instance of CaPoolPublishingOptions into a JSON
// request object.
func expandCaPoolPublishingOptions(c *Client, f *CaPoolPublishingOptions, res *CaPool) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.PublishCaCert; !dcl.IsEmptyValueIndirect(v) {
		m["publishCaCert"] = v
	}
	if v := f.PublishCrl; !dcl.IsEmptyValueIndirect(v) {
		m["publishCrl"] = v
	}

	return m, nil
}

// flattenCaPoolPublishingOptions flattens an instance of CaPoolPublishingOptions from a JSON
// response object.
func flattenCaPoolPublishingOptions(c *Client, i interface{}, res *CaPool) *CaPoolPublishingOptions {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &CaPoolPublishingOptions{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyCaPoolPublishingOptions
	}
	r.PublishCaCert = dcl.FlattenBool(m["publishCaCert"])
	r.PublishCrl = dcl.FlattenBool(m["publishCrl"])

	return r
}

// flattenCaPoolTierEnumMap flattens the contents of CaPoolTierEnum from a JSON
// response object.
func flattenCaPoolTierEnumMap(c *Client, i interface{}, res *CaPool) map[string]CaPoolTierEnum {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]CaPoolTierEnum{}
	}

	if len(a) == 0 {
		return map[string]CaPoolTierEnum{}
	}

	items := make(map[string]CaPoolTierEnum)
	for k, item := range a {
		items[k] = *flattenCaPoolTierEnum(item.(interface{}))
	}

	return items
}

// flattenCaPoolTierEnumSlice flattens the contents of CaPoolTierEnum from a JSON
// response object.
func flattenCaPoolTierEnumSlice(c *Client, i interface{}, res *CaPool) []CaPoolTierEnum {
	a, ok := i.([]interface{})
	if !ok {
		return []CaPoolTierEnum{}
	}

	if len(a) == 0 {
		return []CaPoolTierEnum{}
	}

	items := make([]CaPoolTierEnum, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenCaPoolTierEnum(item.(interface{})))
	}

	return items
}

// flattenCaPoolTierEnum asserts that an interface is a string, and returns a
// pointer to a *CaPoolTierEnum with the same value as that string.
func flattenCaPoolTierEnum(i interface{}) *CaPoolTierEnum {
	s, ok := i.(string)
	if !ok {
		return nil
	}

	return CaPoolTierEnumRef(s)
}

// flattenCaPoolIssuancePolicyAllowedKeyTypesEllipticCurveSignatureAlgorithmEnumMap flattens the contents of CaPoolIssuancePolicyAllowedKeyTypesEllipticCurveSignatureAlgorithmEnum from a JSON
// response object.
func flattenCaPoolIssuancePolicyAllowedKeyTypesEllipticCurveSignatureAlgorithmEnumMap(c *Client, i interface{}, res *CaPool) map[string]CaPoolIssuancePolicyAllowedKeyTypesEllipticCurveSignatureAlgorithmEnum {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]CaPoolIssuancePolicyAllowedKeyTypesEllipticCurveSignatureAlgorithmEnum{}
	}

	if len(a) == 0 {
		return map[string]CaPoolIssuancePolicyAllowedKeyTypesEllipticCurveSignatureAlgorithmEnum{}
	}

	items := make(map[string]CaPoolIssuancePolicyAllowedKeyTypesEllipticCurveSignatureAlgorithmEnum)
	for k, item := range a {
		items[k] = *flattenCaPoolIssuancePolicyAllowedKeyTypesEllipticCurveSignatureAlgorithmEnum(item.(interface{}))
	}

	return items
}

// flattenCaPoolIssuancePolicyAllowedKeyTypesEllipticCurveSignatureAlgorithmEnumSlice flattens the contents of CaPoolIssuancePolicyAllowedKeyTypesEllipticCurveSignatureAlgorithmEnum from a JSON
// response object.
func flattenCaPoolIssuancePolicyAllowedKeyTypesEllipticCurveSignatureAlgorithmEnumSlice(c *Client, i interface{}, res *CaPool) []CaPoolIssuancePolicyAllowedKeyTypesEllipticCurveSignatureAlgorithmEnum {
	a, ok := i.([]interface{})
	if !ok {
		return []CaPoolIssuancePolicyAllowedKeyTypesEllipticCurveSignatureAlgorithmEnum{}
	}

	if len(a) == 0 {
		return []CaPoolIssuancePolicyAllowedKeyTypesEllipticCurveSignatureAlgorithmEnum{}
	}

	items := make([]CaPoolIssuancePolicyAllowedKeyTypesEllipticCurveSignatureAlgorithmEnum, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenCaPoolIssuancePolicyAllowedKeyTypesEllipticCurveSignatureAlgorithmEnum(item.(interface{})))
	}

	return items
}

// flattenCaPoolIssuancePolicyAllowedKeyTypesEllipticCurveSignatureAlgorithmEnum asserts that an interface is a string, and returns a
// pointer to a *CaPoolIssuancePolicyAllowedKeyTypesEllipticCurveSignatureAlgorithmEnum with the same value as that string.
func flattenCaPoolIssuancePolicyAllowedKeyTypesEllipticCurveSignatureAlgorithmEnum(i interface{}) *CaPoolIssuancePolicyAllowedKeyTypesEllipticCurveSignatureAlgorithmEnum {
	s, ok := i.(string)
	if !ok {
		return nil
	}

	return CaPoolIssuancePolicyAllowedKeyTypesEllipticCurveSignatureAlgorithmEnumRef(s)
}

// flattenCaPoolIssuancePolicyPassthroughExtensionsKnownExtensionsEnumMap flattens the contents of CaPoolIssuancePolicyPassthroughExtensionsKnownExtensionsEnum from a JSON
// response object.
func flattenCaPoolIssuancePolicyPassthroughExtensionsKnownExtensionsEnumMap(c *Client, i interface{}, res *CaPool) map[string]CaPoolIssuancePolicyPassthroughExtensionsKnownExtensionsEnum {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]CaPoolIssuancePolicyPassthroughExtensionsKnownExtensionsEnum{}
	}

	if len(a) == 0 {
		return map[string]CaPoolIssuancePolicyPassthroughExtensionsKnownExtensionsEnum{}
	}

	items := make(map[string]CaPoolIssuancePolicyPassthroughExtensionsKnownExtensionsEnum)
	for k, item := range a {
		items[k] = *flattenCaPoolIssuancePolicyPassthroughExtensionsKnownExtensionsEnum(item.(interface{}))
	}

	return items
}

// flattenCaPoolIssuancePolicyPassthroughExtensionsKnownExtensionsEnumSlice flattens the contents of CaPoolIssuancePolicyPassthroughExtensionsKnownExtensionsEnum from a JSON
// response object.
func flattenCaPoolIssuancePolicyPassthroughExtensionsKnownExtensionsEnumSlice(c *Client, i interface{}, res *CaPool) []CaPoolIssuancePolicyPassthroughExtensionsKnownExtensionsEnum {
	a, ok := i.([]interface{})
	if !ok {
		return []CaPoolIssuancePolicyPassthroughExtensionsKnownExtensionsEnum{}
	}

	if len(a) == 0 {
		return []CaPoolIssuancePolicyPassthroughExtensionsKnownExtensionsEnum{}
	}

	items := make([]CaPoolIssuancePolicyPassthroughExtensionsKnownExtensionsEnum, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenCaPoolIssuancePolicyPassthroughExtensionsKnownExtensionsEnum(item.(interface{})))
	}

	return items
}

// flattenCaPoolIssuancePolicyPassthroughExtensionsKnownExtensionsEnum asserts that an interface is a string, and returns a
// pointer to a *CaPoolIssuancePolicyPassthroughExtensionsKnownExtensionsEnum with the same value as that string.
func flattenCaPoolIssuancePolicyPassthroughExtensionsKnownExtensionsEnum(i interface{}) *CaPoolIssuancePolicyPassthroughExtensionsKnownExtensionsEnum {
	s, ok := i.(string)
	if !ok {
		return nil
	}

	return CaPoolIssuancePolicyPassthroughExtensionsKnownExtensionsEnumRef(s)
}

// This function returns a matcher that checks whether a serialized resource matches this resource
// in its parameters (as defined by the fields in a Get, which definitionally define resource
// identity).  This is useful in extracting the element from a List call.
func (r *CaPool) matcher(c *Client) func([]byte) bool {
	return func(b []byte) bool {
		cr, err := unmarshalCaPool(b, c, r)
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

type caPoolDiff struct {
	// The diff should include one or the other of RequiresRecreate or UpdateOp.
	RequiresRecreate bool
	UpdateOp         caPoolApiOperation
	FieldName        string // used for error logging
}

func convertFieldDiffsToCaPoolDiffs(config *dcl.Config, fds []*dcl.FieldDiff, opts []dcl.ApplyOption) ([]caPoolDiff, error) {
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
	var diffs []caPoolDiff
	// For each operation name, create a caPoolDiff which contains the operation.
	for opName, fieldDiffs := range opNamesToFieldDiffs {
		// Use the first field diff's field name for logging required recreate error.
		diff := caPoolDiff{FieldName: fieldDiffs[0].FieldName}
		if opName == "Recreate" {
			diff.RequiresRecreate = true
		} else {
			apiOp, err := convertOpNameToCaPoolApiOperation(opName, fieldDiffs, opts...)
			if err != nil {
				return diffs, err
			}
			diff.UpdateOp = apiOp
		}
		diffs = append(diffs, diff)
	}
	return diffs, nil
}

func convertOpNameToCaPoolApiOperation(opName string, fieldDiffs []*dcl.FieldDiff, opts ...dcl.ApplyOption) (caPoolApiOperation, error) {
	switch opName {

	case "updateCaPoolUpdateCaPoolOperation":
		return &updateCaPoolUpdateCaPoolOperation{FieldDiffs: fieldDiffs}, nil

	default:
		return nil, fmt.Errorf("no such operation with name: %v", opName)
	}
}

func extractCaPoolFields(r *CaPool) error {
	vIssuancePolicy := r.IssuancePolicy
	if vIssuancePolicy == nil {
		// note: explicitly not the empty object.
		vIssuancePolicy = &CaPoolIssuancePolicy{}
	}
	if err := extractCaPoolIssuancePolicyFields(r, vIssuancePolicy); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vIssuancePolicy) {
		r.IssuancePolicy = vIssuancePolicy
	}
	vPublishingOptions := r.PublishingOptions
	if vPublishingOptions == nil {
		// note: explicitly not the empty object.
		vPublishingOptions = &CaPoolPublishingOptions{}
	}
	if err := extractCaPoolPublishingOptionsFields(r, vPublishingOptions); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vPublishingOptions) {
		r.PublishingOptions = vPublishingOptions
	}
	return nil
}
func extractCaPoolIssuancePolicyFields(r *CaPool, o *CaPoolIssuancePolicy) error {
	vAllowedIssuanceModes := o.AllowedIssuanceModes
	if vAllowedIssuanceModes == nil {
		// note: explicitly not the empty object.
		vAllowedIssuanceModes = &CaPoolIssuancePolicyAllowedIssuanceModes{}
	}
	if err := extractCaPoolIssuancePolicyAllowedIssuanceModesFields(r, vAllowedIssuanceModes); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vAllowedIssuanceModes) {
		o.AllowedIssuanceModes = vAllowedIssuanceModes
	}
	vBaselineValues := o.BaselineValues
	if vBaselineValues == nil {
		// note: explicitly not the empty object.
		vBaselineValues = &CaPoolIssuancePolicyBaselineValues{}
	}
	if err := extractCaPoolIssuancePolicyBaselineValuesFields(r, vBaselineValues); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vBaselineValues) {
		o.BaselineValues = vBaselineValues
	}
	vIdentityConstraints := o.IdentityConstraints
	if vIdentityConstraints == nil {
		// note: explicitly not the empty object.
		vIdentityConstraints = &CaPoolIssuancePolicyIdentityConstraints{}
	}
	if err := extractCaPoolIssuancePolicyIdentityConstraintsFields(r, vIdentityConstraints); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vIdentityConstraints) {
		o.IdentityConstraints = vIdentityConstraints
	}
	vPassthroughExtensions := o.PassthroughExtensions
	if vPassthroughExtensions == nil {
		// note: explicitly not the empty object.
		vPassthroughExtensions = &CaPoolIssuancePolicyPassthroughExtensions{}
	}
	if err := extractCaPoolIssuancePolicyPassthroughExtensionsFields(r, vPassthroughExtensions); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vPassthroughExtensions) {
		o.PassthroughExtensions = vPassthroughExtensions
	}
	return nil
}
func extractCaPoolIssuancePolicyAllowedKeyTypesFields(r *CaPool, o *CaPoolIssuancePolicyAllowedKeyTypes) error {
	vRsa := o.Rsa
	if vRsa == nil {
		// note: explicitly not the empty object.
		vRsa = &CaPoolIssuancePolicyAllowedKeyTypesRsa{}
	}
	if err := extractCaPoolIssuancePolicyAllowedKeyTypesRsaFields(r, vRsa); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vRsa) {
		o.Rsa = vRsa
	}
	vEllipticCurve := o.EllipticCurve
	if vEllipticCurve == nil {
		// note: explicitly not the empty object.
		vEllipticCurve = &CaPoolIssuancePolicyAllowedKeyTypesEllipticCurve{}
	}
	if err := extractCaPoolIssuancePolicyAllowedKeyTypesEllipticCurveFields(r, vEllipticCurve); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vEllipticCurve) {
		o.EllipticCurve = vEllipticCurve
	}
	return nil
}
func extractCaPoolIssuancePolicyAllowedKeyTypesRsaFields(r *CaPool, o *CaPoolIssuancePolicyAllowedKeyTypesRsa) error {
	return nil
}
func extractCaPoolIssuancePolicyAllowedKeyTypesEllipticCurveFields(r *CaPool, o *CaPoolIssuancePolicyAllowedKeyTypesEllipticCurve) error {
	return nil
}
func extractCaPoolIssuancePolicyAllowedIssuanceModesFields(r *CaPool, o *CaPoolIssuancePolicyAllowedIssuanceModes) error {
	return nil
}
func extractCaPoolIssuancePolicyBaselineValuesFields(r *CaPool, o *CaPoolIssuancePolicyBaselineValues) error {
	vKeyUsage := o.KeyUsage
	if vKeyUsage == nil {
		// note: explicitly not the empty object.
		vKeyUsage = &CaPoolIssuancePolicyBaselineValuesKeyUsage{}
	}
	if err := extractCaPoolIssuancePolicyBaselineValuesKeyUsageFields(r, vKeyUsage); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vKeyUsage) {
		o.KeyUsage = vKeyUsage
	}
	vCaOptions := o.CaOptions
	if vCaOptions == nil {
		// note: explicitly not the empty object.
		vCaOptions = &CaPoolIssuancePolicyBaselineValuesCaOptions{}
	}
	if err := extractCaPoolIssuancePolicyBaselineValuesCaOptionsFields(r, vCaOptions); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vCaOptions) {
		o.CaOptions = vCaOptions
	}
	return nil
}
func extractCaPoolIssuancePolicyBaselineValuesKeyUsageFields(r *CaPool, o *CaPoolIssuancePolicyBaselineValuesKeyUsage) error {
	vBaseKeyUsage := o.BaseKeyUsage
	if vBaseKeyUsage == nil {
		// note: explicitly not the empty object.
		vBaseKeyUsage = &CaPoolIssuancePolicyBaselineValuesKeyUsageBaseKeyUsage{}
	}
	if err := extractCaPoolIssuancePolicyBaselineValuesKeyUsageBaseKeyUsageFields(r, vBaseKeyUsage); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vBaseKeyUsage) {
		o.BaseKeyUsage = vBaseKeyUsage
	}
	vExtendedKeyUsage := o.ExtendedKeyUsage
	if vExtendedKeyUsage == nil {
		// note: explicitly not the empty object.
		vExtendedKeyUsage = &CaPoolIssuancePolicyBaselineValuesKeyUsageExtendedKeyUsage{}
	}
	if err := extractCaPoolIssuancePolicyBaselineValuesKeyUsageExtendedKeyUsageFields(r, vExtendedKeyUsage); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vExtendedKeyUsage) {
		o.ExtendedKeyUsage = vExtendedKeyUsage
	}
	return nil
}
func extractCaPoolIssuancePolicyBaselineValuesKeyUsageBaseKeyUsageFields(r *CaPool, o *CaPoolIssuancePolicyBaselineValuesKeyUsageBaseKeyUsage) error {
	return nil
}
func extractCaPoolIssuancePolicyBaselineValuesKeyUsageExtendedKeyUsageFields(r *CaPool, o *CaPoolIssuancePolicyBaselineValuesKeyUsageExtendedKeyUsage) error {
	return nil
}
func extractCaPoolIssuancePolicyBaselineValuesKeyUsageUnknownExtendedKeyUsagesFields(r *CaPool, o *CaPoolIssuancePolicyBaselineValuesKeyUsageUnknownExtendedKeyUsages) error {
	return nil
}
func extractCaPoolIssuancePolicyBaselineValuesCaOptionsFields(r *CaPool, o *CaPoolIssuancePolicyBaselineValuesCaOptions) error {
	return nil
}
func extractCaPoolIssuancePolicyBaselineValuesPolicyIdsFields(r *CaPool, o *CaPoolIssuancePolicyBaselineValuesPolicyIds) error {
	return nil
}
func extractCaPoolIssuancePolicyBaselineValuesAdditionalExtensionsFields(r *CaPool, o *CaPoolIssuancePolicyBaselineValuesAdditionalExtensions) error {
	vObjectId := o.ObjectId
	if vObjectId == nil {
		// note: explicitly not the empty object.
		vObjectId = &CaPoolIssuancePolicyBaselineValuesAdditionalExtensionsObjectId{}
	}
	if err := extractCaPoolIssuancePolicyBaselineValuesAdditionalExtensionsObjectIdFields(r, vObjectId); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vObjectId) {
		o.ObjectId = vObjectId
	}
	return nil
}
func extractCaPoolIssuancePolicyBaselineValuesAdditionalExtensionsObjectIdFields(r *CaPool, o *CaPoolIssuancePolicyBaselineValuesAdditionalExtensionsObjectId) error {
	return nil
}
func extractCaPoolIssuancePolicyIdentityConstraintsFields(r *CaPool, o *CaPoolIssuancePolicyIdentityConstraints) error {
	vCelExpression := o.CelExpression
	if vCelExpression == nil {
		// note: explicitly not the empty object.
		vCelExpression = &CaPoolIssuancePolicyIdentityConstraintsCelExpression{}
	}
	if err := extractCaPoolIssuancePolicyIdentityConstraintsCelExpressionFields(r, vCelExpression); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vCelExpression) {
		o.CelExpression = vCelExpression
	}
	return nil
}
func extractCaPoolIssuancePolicyIdentityConstraintsCelExpressionFields(r *CaPool, o *CaPoolIssuancePolicyIdentityConstraintsCelExpression) error {
	return nil
}
func extractCaPoolIssuancePolicyPassthroughExtensionsFields(r *CaPool, o *CaPoolIssuancePolicyPassthroughExtensions) error {
	return nil
}
func extractCaPoolIssuancePolicyPassthroughExtensionsAdditionalExtensionsFields(r *CaPool, o *CaPoolIssuancePolicyPassthroughExtensionsAdditionalExtensions) error {
	return nil
}
func extractCaPoolPublishingOptionsFields(r *CaPool, o *CaPoolPublishingOptions) error {
	return nil
}

func postReadExtractCaPoolFields(r *CaPool) error {
	vIssuancePolicy := r.IssuancePolicy
	if vIssuancePolicy == nil {
		// note: explicitly not the empty object.
		vIssuancePolicy = &CaPoolIssuancePolicy{}
	}
	if err := postReadExtractCaPoolIssuancePolicyFields(r, vIssuancePolicy); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vIssuancePolicy) {
		r.IssuancePolicy = vIssuancePolicy
	}
	vPublishingOptions := r.PublishingOptions
	if vPublishingOptions == nil {
		// note: explicitly not the empty object.
		vPublishingOptions = &CaPoolPublishingOptions{}
	}
	if err := postReadExtractCaPoolPublishingOptionsFields(r, vPublishingOptions); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vPublishingOptions) {
		r.PublishingOptions = vPublishingOptions
	}
	return nil
}
func postReadExtractCaPoolIssuancePolicyFields(r *CaPool, o *CaPoolIssuancePolicy) error {
	vAllowedIssuanceModes := o.AllowedIssuanceModes
	if vAllowedIssuanceModes == nil {
		// note: explicitly not the empty object.
		vAllowedIssuanceModes = &CaPoolIssuancePolicyAllowedIssuanceModes{}
	}
	if err := extractCaPoolIssuancePolicyAllowedIssuanceModesFields(r, vAllowedIssuanceModes); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vAllowedIssuanceModes) {
		o.AllowedIssuanceModes = vAllowedIssuanceModes
	}
	vBaselineValues := o.BaselineValues
	if vBaselineValues == nil {
		// note: explicitly not the empty object.
		vBaselineValues = &CaPoolIssuancePolicyBaselineValues{}
	}
	if err := extractCaPoolIssuancePolicyBaselineValuesFields(r, vBaselineValues); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vBaselineValues) {
		o.BaselineValues = vBaselineValues
	}
	vIdentityConstraints := o.IdentityConstraints
	if vIdentityConstraints == nil {
		// note: explicitly not the empty object.
		vIdentityConstraints = &CaPoolIssuancePolicyIdentityConstraints{}
	}
	if err := extractCaPoolIssuancePolicyIdentityConstraintsFields(r, vIdentityConstraints); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vIdentityConstraints) {
		o.IdentityConstraints = vIdentityConstraints
	}
	vPassthroughExtensions := o.PassthroughExtensions
	if vPassthroughExtensions == nil {
		// note: explicitly not the empty object.
		vPassthroughExtensions = &CaPoolIssuancePolicyPassthroughExtensions{}
	}
	if err := extractCaPoolIssuancePolicyPassthroughExtensionsFields(r, vPassthroughExtensions); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vPassthroughExtensions) {
		o.PassthroughExtensions = vPassthroughExtensions
	}
	return nil
}
func postReadExtractCaPoolIssuancePolicyAllowedKeyTypesFields(r *CaPool, o *CaPoolIssuancePolicyAllowedKeyTypes) error {
	vRsa := o.Rsa
	if vRsa == nil {
		// note: explicitly not the empty object.
		vRsa = &CaPoolIssuancePolicyAllowedKeyTypesRsa{}
	}
	if err := extractCaPoolIssuancePolicyAllowedKeyTypesRsaFields(r, vRsa); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vRsa) {
		o.Rsa = vRsa
	}
	vEllipticCurve := o.EllipticCurve
	if vEllipticCurve == nil {
		// note: explicitly not the empty object.
		vEllipticCurve = &CaPoolIssuancePolicyAllowedKeyTypesEllipticCurve{}
	}
	if err := extractCaPoolIssuancePolicyAllowedKeyTypesEllipticCurveFields(r, vEllipticCurve); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vEllipticCurve) {
		o.EllipticCurve = vEllipticCurve
	}
	return nil
}
func postReadExtractCaPoolIssuancePolicyAllowedKeyTypesRsaFields(r *CaPool, o *CaPoolIssuancePolicyAllowedKeyTypesRsa) error {
	return nil
}
func postReadExtractCaPoolIssuancePolicyAllowedKeyTypesEllipticCurveFields(r *CaPool, o *CaPoolIssuancePolicyAllowedKeyTypesEllipticCurve) error {
	return nil
}
func postReadExtractCaPoolIssuancePolicyAllowedIssuanceModesFields(r *CaPool, o *CaPoolIssuancePolicyAllowedIssuanceModes) error {
	return nil
}
func postReadExtractCaPoolIssuancePolicyBaselineValuesFields(r *CaPool, o *CaPoolIssuancePolicyBaselineValues) error {
	vKeyUsage := o.KeyUsage
	if vKeyUsage == nil {
		// note: explicitly not the empty object.
		vKeyUsage = &CaPoolIssuancePolicyBaselineValuesKeyUsage{}
	}
	if err := extractCaPoolIssuancePolicyBaselineValuesKeyUsageFields(r, vKeyUsage); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vKeyUsage) {
		o.KeyUsage = vKeyUsage
	}
	vCaOptions := o.CaOptions
	if vCaOptions == nil {
		// note: explicitly not the empty object.
		vCaOptions = &CaPoolIssuancePolicyBaselineValuesCaOptions{}
	}
	if err := extractCaPoolIssuancePolicyBaselineValuesCaOptionsFields(r, vCaOptions); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vCaOptions) {
		o.CaOptions = vCaOptions
	}
	return nil
}
func postReadExtractCaPoolIssuancePolicyBaselineValuesKeyUsageFields(r *CaPool, o *CaPoolIssuancePolicyBaselineValuesKeyUsage) error {
	vBaseKeyUsage := o.BaseKeyUsage
	if vBaseKeyUsage == nil {
		// note: explicitly not the empty object.
		vBaseKeyUsage = &CaPoolIssuancePolicyBaselineValuesKeyUsageBaseKeyUsage{}
	}
	if err := extractCaPoolIssuancePolicyBaselineValuesKeyUsageBaseKeyUsageFields(r, vBaseKeyUsage); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vBaseKeyUsage) {
		o.BaseKeyUsage = vBaseKeyUsage
	}
	vExtendedKeyUsage := o.ExtendedKeyUsage
	if vExtendedKeyUsage == nil {
		// note: explicitly not the empty object.
		vExtendedKeyUsage = &CaPoolIssuancePolicyBaselineValuesKeyUsageExtendedKeyUsage{}
	}
	if err := extractCaPoolIssuancePolicyBaselineValuesKeyUsageExtendedKeyUsageFields(r, vExtendedKeyUsage); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vExtendedKeyUsage) {
		o.ExtendedKeyUsage = vExtendedKeyUsage
	}
	return nil
}
func postReadExtractCaPoolIssuancePolicyBaselineValuesKeyUsageBaseKeyUsageFields(r *CaPool, o *CaPoolIssuancePolicyBaselineValuesKeyUsageBaseKeyUsage) error {
	return nil
}
func postReadExtractCaPoolIssuancePolicyBaselineValuesKeyUsageExtendedKeyUsageFields(r *CaPool, o *CaPoolIssuancePolicyBaselineValuesKeyUsageExtendedKeyUsage) error {
	return nil
}
func postReadExtractCaPoolIssuancePolicyBaselineValuesKeyUsageUnknownExtendedKeyUsagesFields(r *CaPool, o *CaPoolIssuancePolicyBaselineValuesKeyUsageUnknownExtendedKeyUsages) error {
	return nil
}
func postReadExtractCaPoolIssuancePolicyBaselineValuesCaOptionsFields(r *CaPool, o *CaPoolIssuancePolicyBaselineValuesCaOptions) error {
	return nil
}
func postReadExtractCaPoolIssuancePolicyBaselineValuesPolicyIdsFields(r *CaPool, o *CaPoolIssuancePolicyBaselineValuesPolicyIds) error {
	return nil
}
func postReadExtractCaPoolIssuancePolicyBaselineValuesAdditionalExtensionsFields(r *CaPool, o *CaPoolIssuancePolicyBaselineValuesAdditionalExtensions) error {
	vObjectId := o.ObjectId
	if vObjectId == nil {
		// note: explicitly not the empty object.
		vObjectId = &CaPoolIssuancePolicyBaselineValuesAdditionalExtensionsObjectId{}
	}
	if err := extractCaPoolIssuancePolicyBaselineValuesAdditionalExtensionsObjectIdFields(r, vObjectId); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vObjectId) {
		o.ObjectId = vObjectId
	}
	return nil
}
func postReadExtractCaPoolIssuancePolicyBaselineValuesAdditionalExtensionsObjectIdFields(r *CaPool, o *CaPoolIssuancePolicyBaselineValuesAdditionalExtensionsObjectId) error {
	return nil
}
func postReadExtractCaPoolIssuancePolicyIdentityConstraintsFields(r *CaPool, o *CaPoolIssuancePolicyIdentityConstraints) error {
	vCelExpression := o.CelExpression
	if vCelExpression == nil {
		// note: explicitly not the empty object.
		vCelExpression = &CaPoolIssuancePolicyIdentityConstraintsCelExpression{}
	}
	if err := extractCaPoolIssuancePolicyIdentityConstraintsCelExpressionFields(r, vCelExpression); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vCelExpression) {
		o.CelExpression = vCelExpression
	}
	return nil
}
func postReadExtractCaPoolIssuancePolicyIdentityConstraintsCelExpressionFields(r *CaPool, o *CaPoolIssuancePolicyIdentityConstraintsCelExpression) error {
	return nil
}
func postReadExtractCaPoolIssuancePolicyPassthroughExtensionsFields(r *CaPool, o *CaPoolIssuancePolicyPassthroughExtensions) error {
	return nil
}
func postReadExtractCaPoolIssuancePolicyPassthroughExtensionsAdditionalExtensionsFields(r *CaPool, o *CaPoolIssuancePolicyPassthroughExtensionsAdditionalExtensions) error {
	return nil
}
func postReadExtractCaPoolPublishingOptionsFields(r *CaPool, o *CaPoolPublishingOptions) error {
	return nil
}
