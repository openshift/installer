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

func (r *CertificateAuthority) validate() error {

	if err := dcl.Required(r, "name"); err != nil {
		return err
	}
	if err := dcl.Required(r, "type"); err != nil {
		return err
	}
	if err := dcl.Required(r, "config"); err != nil {
		return err
	}
	if err := dcl.Required(r, "lifetime"); err != nil {
		return err
	}
	if err := dcl.Required(r, "keySpec"); err != nil {
		return err
	}
	if err := dcl.RequiredParameter(r.Project, "Project"); err != nil {
		return err
	}
	if err := dcl.RequiredParameter(r.Location, "Location"); err != nil {
		return err
	}
	if err := dcl.RequiredParameter(r.CaPool, "CaPool"); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(r.Config) {
		if err := r.Config.validate(); err != nil {
			return err
		}
	}
	if !dcl.IsEmptyValueIndirect(r.KeySpec) {
		if err := r.KeySpec.validate(); err != nil {
			return err
		}
	}
	if !dcl.IsEmptyValueIndirect(r.SubordinateConfig) {
		if err := r.SubordinateConfig.validate(); err != nil {
			return err
		}
	}
	if !dcl.IsEmptyValueIndirect(r.AccessUrls) {
		if err := r.AccessUrls.validate(); err != nil {
			return err
		}
	}
	return nil
}
func (r *CertificateAuthorityConfig) validate() error {
	if err := dcl.Required(r, "subjectConfig"); err != nil {
		return err
	}
	if err := dcl.Required(r, "x509Config"); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(r.SubjectConfig) {
		if err := r.SubjectConfig.validate(); err != nil {
			return err
		}
	}
	if !dcl.IsEmptyValueIndirect(r.X509Config) {
		if err := r.X509Config.validate(); err != nil {
			return err
		}
	}
	if !dcl.IsEmptyValueIndirect(r.PublicKey) {
		if err := r.PublicKey.validate(); err != nil {
			return err
		}
	}
	return nil
}
func (r *CertificateAuthorityConfigSubjectConfig) validate() error {
	if err := dcl.Required(r, "subject"); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(r.Subject) {
		if err := r.Subject.validate(); err != nil {
			return err
		}
	}
	if !dcl.IsEmptyValueIndirect(r.SubjectAltName) {
		if err := r.SubjectAltName.validate(); err != nil {
			return err
		}
	}
	return nil
}
func (r *CertificateAuthorityConfigSubjectConfigSubject) validate() error {
	return nil
}
func (r *CertificateAuthorityConfigSubjectConfigSubjectAltName) validate() error {
	return nil
}
func (r *CertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSans) validate() error {
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
func (r *CertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSansObjectId) validate() error {
	if err := dcl.Required(r, "objectIdPath"); err != nil {
		return err
	}
	return nil
}
func (r *CertificateAuthorityConfigX509Config) validate() error {
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
func (r *CertificateAuthorityConfigX509ConfigKeyUsage) validate() error {
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
func (r *CertificateAuthorityConfigX509ConfigKeyUsageBaseKeyUsage) validate() error {
	return nil
}
func (r *CertificateAuthorityConfigX509ConfigKeyUsageExtendedKeyUsage) validate() error {
	return nil
}
func (r *CertificateAuthorityConfigX509ConfigKeyUsageUnknownExtendedKeyUsages) validate() error {
	if err := dcl.Required(r, "objectIdPath"); err != nil {
		return err
	}
	return nil
}
func (r *CertificateAuthorityConfigX509ConfigCaOptions) validate() error {
	return nil
}
func (r *CertificateAuthorityConfigX509ConfigPolicyIds) validate() error {
	if err := dcl.Required(r, "objectIdPath"); err != nil {
		return err
	}
	return nil
}
func (r *CertificateAuthorityConfigX509ConfigAdditionalExtensions) validate() error {
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
func (r *CertificateAuthorityConfigX509ConfigAdditionalExtensionsObjectId) validate() error {
	if err := dcl.Required(r, "objectIdPath"); err != nil {
		return err
	}
	return nil
}
func (r *CertificateAuthorityConfigPublicKey) validate() error {
	if err := dcl.Required(r, "key"); err != nil {
		return err
	}
	if err := dcl.Required(r, "format"); err != nil {
		return err
	}
	return nil
}
func (r *CertificateAuthorityKeySpec) validate() error {
	if err := dcl.ValidateExactlyOneOfFieldsSet([]string{"CloudKmsKeyVersion", "Algorithm"}, r.CloudKmsKeyVersion, r.Algorithm); err != nil {
		return err
	}
	return nil
}
func (r *CertificateAuthoritySubordinateConfig) validate() error {
	if err := dcl.ValidateExactlyOneOfFieldsSet([]string{"CertificateAuthority", "PemIssuerChain"}, r.CertificateAuthority, r.PemIssuerChain); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(r.PemIssuerChain) {
		if err := r.PemIssuerChain.validate(); err != nil {
			return err
		}
	}
	return nil
}
func (r *CertificateAuthoritySubordinateConfigPemIssuerChain) validate() error {
	if err := dcl.Required(r, "pemCertificates"); err != nil {
		return err
	}
	return nil
}
func (r *CertificateAuthorityCaCertificateDescriptions) validate() error {
	if !dcl.IsEmptyValueIndirect(r.SubjectDescription) {
		if err := r.SubjectDescription.validate(); err != nil {
			return err
		}
	}
	if !dcl.IsEmptyValueIndirect(r.X509Description) {
		if err := r.X509Description.validate(); err != nil {
			return err
		}
	}
	if !dcl.IsEmptyValueIndirect(r.PublicKey) {
		if err := r.PublicKey.validate(); err != nil {
			return err
		}
	}
	if !dcl.IsEmptyValueIndirect(r.SubjectKeyId) {
		if err := r.SubjectKeyId.validate(); err != nil {
			return err
		}
	}
	if !dcl.IsEmptyValueIndirect(r.AuthorityKeyId) {
		if err := r.AuthorityKeyId.validate(); err != nil {
			return err
		}
	}
	if !dcl.IsEmptyValueIndirect(r.CertFingerprint) {
		if err := r.CertFingerprint.validate(); err != nil {
			return err
		}
	}
	return nil
}
func (r *CertificateAuthorityCaCertificateDescriptionsSubjectDescription) validate() error {
	if !dcl.IsEmptyValueIndirect(r.Subject) {
		if err := r.Subject.validate(); err != nil {
			return err
		}
	}
	if !dcl.IsEmptyValueIndirect(r.SubjectAltName) {
		if err := r.SubjectAltName.validate(); err != nil {
			return err
		}
	}
	return nil
}
func (r *CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubject) validate() error {
	return nil
}
func (r *CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltName) validate() error {
	return nil
}
func (r *CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSans) validate() error {
	if err := dcl.Required(r, "objectId"); err != nil {
		return err
	}
	if err := dcl.Required(r, "critical"); err != nil {
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
func (r *CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSansObjectId) validate() error {
	if err := dcl.Required(r, "objectIdPath"); err != nil {
		return err
	}
	return nil
}
func (r *CertificateAuthorityCaCertificateDescriptionsX509Description) validate() error {
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
func (r *CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsage) validate() error {
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
func (r *CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageBaseKeyUsage) validate() error {
	return nil
}
func (r *CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageExtendedKeyUsage) validate() error {
	return nil
}
func (r *CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageUnknownExtendedKeyUsages) validate() error {
	if err := dcl.Required(r, "objectIdPath"); err != nil {
		return err
	}
	return nil
}
func (r *CertificateAuthorityCaCertificateDescriptionsX509DescriptionCaOptions) validate() error {
	return nil
}
func (r *CertificateAuthorityCaCertificateDescriptionsX509DescriptionPolicyIds) validate() error {
	if err := dcl.Required(r, "objectIdPath"); err != nil {
		return err
	}
	return nil
}
func (r *CertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensions) validate() error {
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
func (r *CertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensionsObjectId) validate() error {
	if err := dcl.Required(r, "objectIdPath"); err != nil {
		return err
	}
	return nil
}
func (r *CertificateAuthorityCaCertificateDescriptionsPublicKey) validate() error {
	if err := dcl.Required(r, "key"); err != nil {
		return err
	}
	if err := dcl.Required(r, "format"); err != nil {
		return err
	}
	return nil
}
func (r *CertificateAuthorityCaCertificateDescriptionsSubjectKeyId) validate() error {
	return nil
}
func (r *CertificateAuthorityCaCertificateDescriptionsAuthorityKeyId) validate() error {
	return nil
}
func (r *CertificateAuthorityCaCertificateDescriptionsCertFingerprint) validate() error {
	return nil
}
func (r *CertificateAuthorityAccessUrls) validate() error {
	return nil
}
func (r *CertificateAuthority) basePath() string {
	params := map[string]interface{}{}
	return dcl.Nprintf("https://privateca.googleapis.com/v1/", params)
}

func (r *CertificateAuthority) getURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"project":  dcl.ValueOrEmptyString(nr.Project),
		"location": dcl.ValueOrEmptyString(nr.Location),
		"caPool":   dcl.ValueOrEmptyString(nr.CaPool),
		"name":     dcl.ValueOrEmptyString(nr.Name),
	}
	return dcl.URL("projects/{{project}}/locations/{{location}}/caPools/{{caPool}}/certificateAuthorities/{{name}}", nr.basePath(), userBasePath, params), nil
}

func (r *CertificateAuthority) listURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"project":  dcl.ValueOrEmptyString(nr.Project),
		"location": dcl.ValueOrEmptyString(nr.Location),
		"caPool":   dcl.ValueOrEmptyString(nr.CaPool),
	}
	return dcl.URL("projects/{{project}}/locations/{{location}}/caPools/{{caPool}}/certificateAuthorities", nr.basePath(), userBasePath, params), nil

}

func (r *CertificateAuthority) createURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"project":  dcl.ValueOrEmptyString(nr.Project),
		"location": dcl.ValueOrEmptyString(nr.Location),
		"caPool":   dcl.ValueOrEmptyString(nr.CaPool),
		"name":     dcl.ValueOrEmptyString(nr.Name),
	}
	return dcl.URL("projects/{{project}}/locations/{{location}}/caPools/{{caPool}}/certificateAuthorities?certificateAuthorityId={{name}}", nr.basePath(), userBasePath, params), nil

}

func (r *CertificateAuthority) deleteURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"project":  dcl.ValueOrEmptyString(nr.Project),
		"location": dcl.ValueOrEmptyString(nr.Location),
		"caPool":   dcl.ValueOrEmptyString(nr.CaPool),
		"name":     dcl.ValueOrEmptyString(nr.Name),
	}
	return dcl.URL("projects/{{project}}/locations/{{location}}/caPools/{{caPool}}/certificateAuthorities/{{name}}?ignoreActiveCertificates=true", nr.basePath(), userBasePath, params), nil
}

// certificateAuthorityApiOperation represents a mutable operation in the underlying REST
// API such as Create, Update, or Delete.
type certificateAuthorityApiOperation interface {
	do(context.Context, *CertificateAuthority, *Client) error
}

// newUpdateCertificateAuthorityUpdateCertificateAuthorityRequest creates a request for an
// CertificateAuthority resource's UpdateCertificateAuthority update type by filling in the update
// fields based on the intended state of the resource.
func newUpdateCertificateAuthorityUpdateCertificateAuthorityRequest(ctx context.Context, f *CertificateAuthority, c *Client) (map[string]interface{}, error) {
	req := map[string]interface{}{}
	res := f
	_ = res

	if v := f.Labels; !dcl.IsEmptyValueIndirect(v) {
		req["labels"] = v
	}
	return req, nil
}

// marshalUpdateCertificateAuthorityUpdateCertificateAuthorityRequest converts the update into
// the final JSON request body.
func marshalUpdateCertificateAuthorityUpdateCertificateAuthorityRequest(c *Client, m map[string]interface{}) ([]byte, error) {

	return json.Marshal(m)
}

type updateCertificateAuthorityUpdateCertificateAuthorityOperation struct {
	// If the update operation has the REQUIRES_APPLY_OPTIONS trait, this will be populated.
	// Usually it will be nil - this is to prevent us from accidentally depending on apply
	// options, which should usually be unnecessary.
	ApplyOptions []dcl.ApplyOption
	FieldDiffs   []*dcl.FieldDiff
}

// do creates a request and sends it to the appropriate URL. In most operations,
// do will transcribe a subset of the resource into a request object and send a
// PUT request to a single URL.

func (op *updateCertificateAuthorityUpdateCertificateAuthorityOperation) do(ctx context.Context, r *CertificateAuthority, c *Client) error {
	_, err := c.GetCertificateAuthority(ctx, r)
	if err != nil {
		return err
	}

	u, err := r.updateURL(c.Config.BasePath, "UpdateCertificateAuthority")
	if err != nil {
		return err
	}
	mask := dcl.UpdateMask(op.FieldDiffs)
	u, err = dcl.AddQueryParams(u, map[string]string{"updateMask": mask})
	if err != nil {
		return err
	}

	req, err := newUpdateCertificateAuthorityUpdateCertificateAuthorityRequest(ctx, r, c)
	if err != nil {
		return err
	}

	c.Config.Logger.InfoWithContextf(ctx, "Created update: %#v", req)
	body, err := marshalUpdateCertificateAuthorityUpdateCertificateAuthorityRequest(c, req)
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

func (c *Client) listCertificateAuthorityRaw(ctx context.Context, r *CertificateAuthority, pageToken string, pageSize int32) ([]byte, error) {
	u, err := r.urlNormalized().listURL(c.Config.BasePath)
	if err != nil {
		return nil, err
	}

	m := make(map[string]string)
	if pageToken != "" {
		m["pageToken"] = pageToken
	}

	if pageSize != CertificateAuthorityMaxPage {
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

type listCertificateAuthorityOperation struct {
	CertificateAuthorities []map[string]interface{} `json:"certificateAuthorities"`
	Token                  string                   `json:"nextPageToken"`
}

func (c *Client) listCertificateAuthority(ctx context.Context, r *CertificateAuthority, pageToken string, pageSize int32) ([]*CertificateAuthority, string, error) {
	b, err := c.listCertificateAuthorityRaw(ctx, r, pageToken, pageSize)
	if err != nil {
		return nil, "", err
	}

	var m listCertificateAuthorityOperation
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, "", err
	}

	var l []*CertificateAuthority
	for _, v := range m.CertificateAuthorities {
		res, err := unmarshalMapCertificateAuthority(v, c, r)
		if err != nil {
			return nil, m.Token, err
		}
		res.Project = r.Project
		res.Location = r.Location
		res.CaPool = r.CaPool
		l = append(l, res)
	}

	return l, m.Token, nil
}

func (c *Client) deleteAllCertificateAuthority(ctx context.Context, f func(*CertificateAuthority) bool, resources []*CertificateAuthority) error {
	var errors []string
	for _, res := range resources {
		if f(res) {
			// We do not want deleteAll to fail on a deletion or else it will stop deleting other resources.
			err := c.DeleteCertificateAuthority(ctx, res)
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

type deleteCertificateAuthorityOperation struct{}

func (op *deleteCertificateAuthorityOperation) do(ctx context.Context, r *CertificateAuthority, c *Client) error {
	r, err := c.GetCertificateAuthority(ctx, r)
	if err != nil {
		if dcl.IsNotFound(err) {
			c.Config.Logger.InfoWithContextf(ctx, "CertificateAuthority not found, returning. Original error: %v", err)
			return nil
		}
		c.Config.Logger.WarningWithContextf(ctx, "GetCertificateAuthority checking for existence. error: %v", err)
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
	return nil
}

// Create operations are similar to Update operations, although they do not have
// specific request objects. The Create request object is the json encoding of
// the resource, which is modified by res.marshal to form the base request body.
type createCertificateAuthorityOperation struct {
	response map[string]interface{}
}

func (op *createCertificateAuthorityOperation) FirstResponse() (map[string]interface{}, bool) {
	return op.response, len(op.response) > 0
}

func (op *createCertificateAuthorityOperation) do(ctx context.Context, r *CertificateAuthority, c *Client) error {
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

	if _, err := c.GetCertificateAuthority(ctx, r); err != nil {
		c.Config.Logger.WarningWithContextf(ctx, "get returned error: %v", err)
		return err
	}

	err = r.waitForCertificateAuthorityEnabled(ctx, c)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) getCertificateAuthorityRaw(ctx context.Context, r *CertificateAuthority) ([]byte, error) {

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

func (c *Client) certificateAuthorityDiffsForRawDesired(ctx context.Context, rawDesired *CertificateAuthority, opts ...dcl.ApplyOption) (initial, desired *CertificateAuthority, diffs []*dcl.FieldDiff, err error) {
	c.Config.Logger.InfoWithContext(ctx, "Fetching initial state...")
	// First, let us see if the user provided a state hint.  If they did, we will start fetching based on that.
	var fetchState *CertificateAuthority
	if sh := dcl.FetchStateHint(opts); sh != nil {
		if r, ok := sh.(*CertificateAuthority); !ok {
			c.Config.Logger.WarningWithContextf(ctx, "Initial state hint was of the wrong type; expected CertificateAuthority, got %T", sh)
		} else {
			fetchState = r
		}
	}
	if fetchState == nil {
		fetchState = rawDesired
	}

	// 1.2: Retrieval of raw initial state from API
	rawInitial, err := c.GetCertificateAuthority(ctx, fetchState)
	if rawInitial == nil {
		if !dcl.IsNotFound(err) {
			c.Config.Logger.WarningWithContextf(ctx, "Failed to retrieve whether a CertificateAuthority resource already exists: %s", err)
			return nil, nil, nil, fmt.Errorf("failed to retrieve CertificateAuthority resource: %v", err)
		}
		c.Config.Logger.InfoWithContext(ctx, "Found that CertificateAuthority resource did not exist.")
		// Perform canonicalization to pick up defaults.
		desired, err = canonicalizeCertificateAuthorityDesiredState(rawDesired, rawInitial)
		return nil, desired, nil, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Found initial state for CertificateAuthority: %v", rawInitial)
	c.Config.Logger.InfoWithContextf(ctx, "Initial desired state for CertificateAuthority: %v", rawDesired)

	// The Get call applies postReadExtract and so the result may contain fields that are not part of API version.
	if err := extractCertificateAuthorityFields(rawInitial); err != nil {
		return nil, nil, nil, err
	}

	// 1.3: Canonicalize raw initial state into initial state.
	initial, err = canonicalizeCertificateAuthorityInitialState(rawInitial, rawDesired)
	if err != nil {
		return nil, nil, nil, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Canonicalized initial state for CertificateAuthority: %v", initial)

	// 1.4: Canonicalize raw desired state into desired state.
	desired, err = canonicalizeCertificateAuthorityDesiredState(rawDesired, rawInitial, opts...)
	if err != nil {
		return nil, nil, nil, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Canonicalized desired state for CertificateAuthority: %v", desired)

	// 2.1: Comparison of initial and desired state.
	diffs, err = diffCertificateAuthority(c, desired, initial, opts...)
	return initial, desired, diffs, err
}

func canonicalizeCertificateAuthorityInitialState(rawInitial, rawDesired *CertificateAuthority) (*CertificateAuthority, error) {
	// TODO(magic-modules-eng): write canonicalizer once relevant traits are added.
	return rawInitial, nil
}

/*
* Canonicalizers
*
* These are responsible for converting either a user-specified config or a
* GCP API response to a standard format that can be used for difference checking.
* */

func canonicalizeCertificateAuthorityDesiredState(rawDesired, rawInitial *CertificateAuthority, opts ...dcl.ApplyOption) (*CertificateAuthority, error) {

	if rawInitial == nil {
		// Since the initial state is empty, the desired state is all we have.
		// We canonicalize the remaining nested objects with nil to pick up defaults.
		rawDesired.Config = canonicalizeCertificateAuthorityConfig(rawDesired.Config, nil, opts...)
		rawDesired.KeySpec = canonicalizeCertificateAuthorityKeySpec(rawDesired.KeySpec, nil, opts...)
		rawDesired.SubordinateConfig = canonicalizeCertificateAuthoritySubordinateConfig(rawDesired.SubordinateConfig, nil, opts...)
		rawDesired.AccessUrls = canonicalizeCertificateAuthorityAccessUrls(rawDesired.AccessUrls, nil, opts...)

		return rawDesired, nil
	}
	canonicalDesired := &CertificateAuthority{}
	if dcl.PartialSelfLinkToSelfLink(rawDesired.Name, rawInitial.Name) {
		canonicalDesired.Name = rawInitial.Name
	} else {
		canonicalDesired.Name = rawDesired.Name
	}
	if dcl.IsZeroValue(rawDesired.Type) || (dcl.IsEmptyValueIndirect(rawDesired.Type) && dcl.IsEmptyValueIndirect(rawInitial.Type)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		canonicalDesired.Type = rawInitial.Type
	} else {
		canonicalDesired.Type = rawDesired.Type
	}
	canonicalDesired.Config = canonicalizeCertificateAuthorityConfig(rawDesired.Config, rawInitial.Config, opts...)
	if dcl.StringCanonicalize(rawDesired.Lifetime, rawInitial.Lifetime) {
		canonicalDesired.Lifetime = rawInitial.Lifetime
	} else {
		canonicalDesired.Lifetime = rawDesired.Lifetime
	}
	canonicalDesired.KeySpec = canonicalizeCertificateAuthorityKeySpec(rawDesired.KeySpec, rawInitial.KeySpec, opts...)
	if dcl.IsZeroValue(rawDesired.GcsBucket) || (dcl.IsEmptyValueIndirect(rawDesired.GcsBucket) && dcl.IsEmptyValueIndirect(rawInitial.GcsBucket)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		canonicalDesired.GcsBucket = rawInitial.GcsBucket
	} else {
		canonicalDesired.GcsBucket = rawDesired.GcsBucket
	}
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
	if dcl.NameToSelfLink(rawDesired.CaPool, rawInitial.CaPool) {
		canonicalDesired.CaPool = rawInitial.CaPool
	} else {
		canonicalDesired.CaPool = rawDesired.CaPool
	}
	return canonicalDesired, nil
}

func canonicalizeCertificateAuthorityNewState(c *Client, rawNew, rawDesired *CertificateAuthority) (*CertificateAuthority, error) {

	if dcl.IsEmptyValueIndirect(rawNew.Name) && dcl.IsEmptyValueIndirect(rawDesired.Name) {
		rawNew.Name = rawDesired.Name
	} else {
		if dcl.PartialSelfLinkToSelfLink(rawDesired.Name, rawNew.Name) {
			rawNew.Name = rawDesired.Name
		}
	}

	if dcl.IsEmptyValueIndirect(rawNew.Type) && dcl.IsEmptyValueIndirect(rawDesired.Type) {
		rawNew.Type = rawDesired.Type
	} else {
	}

	if dcl.IsEmptyValueIndirect(rawNew.Config) && dcl.IsEmptyValueIndirect(rawDesired.Config) {
		rawNew.Config = rawDesired.Config
	} else {
		rawNew.Config = canonicalizeNewCertificateAuthorityConfig(c, rawDesired.Config, rawNew.Config)
	}

	if dcl.IsEmptyValueIndirect(rawNew.Lifetime) && dcl.IsEmptyValueIndirect(rawDesired.Lifetime) {
		rawNew.Lifetime = rawDesired.Lifetime
	} else {
		if dcl.StringCanonicalize(rawDesired.Lifetime, rawNew.Lifetime) {
			rawNew.Lifetime = rawDesired.Lifetime
		}
	}

	if dcl.IsEmptyValueIndirect(rawNew.KeySpec) && dcl.IsEmptyValueIndirect(rawDesired.KeySpec) {
		rawNew.KeySpec = rawDesired.KeySpec
	} else {
		rawNew.KeySpec = canonicalizeNewCertificateAuthorityKeySpec(c, rawDesired.KeySpec, rawNew.KeySpec)
	}

	if dcl.IsEmptyValueIndirect(rawNew.SubordinateConfig) && dcl.IsEmptyValueIndirect(rawDesired.SubordinateConfig) {
		rawNew.SubordinateConfig = rawDesired.SubordinateConfig
	} else {
		rawNew.SubordinateConfig = canonicalizeNewCertificateAuthoritySubordinateConfig(c, rawDesired.SubordinateConfig, rawNew.SubordinateConfig)
	}

	if dcl.IsEmptyValueIndirect(rawNew.Tier) && dcl.IsEmptyValueIndirect(rawDesired.Tier) {
		rawNew.Tier = rawDesired.Tier
	} else {
	}

	if dcl.IsEmptyValueIndirect(rawNew.State) && dcl.IsEmptyValueIndirect(rawDesired.State) {
		rawNew.State = rawDesired.State
	} else {
	}

	if dcl.IsEmptyValueIndirect(rawNew.PemCaCertificates) && dcl.IsEmptyValueIndirect(rawDesired.PemCaCertificates) {
		rawNew.PemCaCertificates = rawDesired.PemCaCertificates
	} else {
		if dcl.StringArrayCanonicalize(rawDesired.PemCaCertificates, rawNew.PemCaCertificates) {
			rawNew.PemCaCertificates = rawDesired.PemCaCertificates
		}
	}

	if dcl.IsEmptyValueIndirect(rawNew.CaCertificateDescriptions) && dcl.IsEmptyValueIndirect(rawDesired.CaCertificateDescriptions) {
		rawNew.CaCertificateDescriptions = rawDesired.CaCertificateDescriptions
	} else {
		rawNew.CaCertificateDescriptions = canonicalizeNewCertificateAuthorityCaCertificateDescriptionsSlice(c, rawDesired.CaCertificateDescriptions, rawNew.CaCertificateDescriptions)
	}

	if dcl.IsEmptyValueIndirect(rawNew.GcsBucket) && dcl.IsEmptyValueIndirect(rawDesired.GcsBucket) {
		rawNew.GcsBucket = rawDesired.GcsBucket
	} else {
	}

	if dcl.IsEmptyValueIndirect(rawNew.AccessUrls) && dcl.IsEmptyValueIndirect(rawDesired.AccessUrls) {
		rawNew.AccessUrls = rawDesired.AccessUrls
	} else {
		rawNew.AccessUrls = canonicalizeNewCertificateAuthorityAccessUrls(c, rawDesired.AccessUrls, rawNew.AccessUrls)
	}

	if dcl.IsEmptyValueIndirect(rawNew.CreateTime) && dcl.IsEmptyValueIndirect(rawDesired.CreateTime) {
		rawNew.CreateTime = rawDesired.CreateTime
	} else {
	}

	if dcl.IsEmptyValueIndirect(rawNew.UpdateTime) && dcl.IsEmptyValueIndirect(rawDesired.UpdateTime) {
		rawNew.UpdateTime = rawDesired.UpdateTime
	} else {
	}

	if dcl.IsEmptyValueIndirect(rawNew.DeleteTime) && dcl.IsEmptyValueIndirect(rawDesired.DeleteTime) {
		rawNew.DeleteTime = rawDesired.DeleteTime
	} else {
	}

	if dcl.IsEmptyValueIndirect(rawNew.ExpireTime) && dcl.IsEmptyValueIndirect(rawDesired.ExpireTime) {
		rawNew.ExpireTime = rawDesired.ExpireTime
	} else {
	}

	if dcl.IsEmptyValueIndirect(rawNew.Labels) && dcl.IsEmptyValueIndirect(rawDesired.Labels) {
		rawNew.Labels = rawDesired.Labels
	} else {
	}

	rawNew.Project = rawDesired.Project

	rawNew.Location = rawDesired.Location

	rawNew.CaPool = rawDesired.CaPool

	return rawNew, nil
}

func canonicalizeCertificateAuthorityConfig(des, initial *CertificateAuthorityConfig, opts ...dcl.ApplyOption) *CertificateAuthorityConfig {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &CertificateAuthorityConfig{}

	cDes.SubjectConfig = canonicalizeCertificateAuthorityConfigSubjectConfig(des.SubjectConfig, initial.SubjectConfig, opts...)
	cDes.X509Config = canonicalizeCertificateAuthorityConfigX509Config(des.X509Config, initial.X509Config, opts...)

	return cDes
}

func canonicalizeCertificateAuthorityConfigSlice(des, initial []CertificateAuthorityConfig, opts ...dcl.ApplyOption) []CertificateAuthorityConfig {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]CertificateAuthorityConfig, 0, len(des))
		for _, d := range des {
			cd := canonicalizeCertificateAuthorityConfig(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]CertificateAuthorityConfig, 0, len(des))
	for i, d := range des {
		cd := canonicalizeCertificateAuthorityConfig(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewCertificateAuthorityConfig(c *Client, des, nw *CertificateAuthorityConfig) *CertificateAuthorityConfig {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for CertificateAuthorityConfig while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	nw.SubjectConfig = canonicalizeNewCertificateAuthorityConfigSubjectConfig(c, des.SubjectConfig, nw.SubjectConfig)
	nw.X509Config = canonicalizeNewCertificateAuthorityConfigX509Config(c, des.X509Config, nw.X509Config)
	nw.PublicKey = canonicalizeNewCertificateAuthorityConfigPublicKey(c, des.PublicKey, nw.PublicKey)

	return nw
}

func canonicalizeNewCertificateAuthorityConfigSet(c *Client, des, nw []CertificateAuthorityConfig) []CertificateAuthorityConfig {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []CertificateAuthorityConfig
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareCertificateAuthorityConfigNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewCertificateAuthorityConfig(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewCertificateAuthorityConfigSlice(c *Client, des, nw []CertificateAuthorityConfig) []CertificateAuthorityConfig {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []CertificateAuthorityConfig
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewCertificateAuthorityConfig(c, &d, &n))
	}

	return items
}

func canonicalizeCertificateAuthorityConfigSubjectConfig(des, initial *CertificateAuthorityConfigSubjectConfig, opts ...dcl.ApplyOption) *CertificateAuthorityConfigSubjectConfig {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &CertificateAuthorityConfigSubjectConfig{}

	cDes.Subject = canonicalizeCertificateAuthorityConfigSubjectConfigSubject(des.Subject, initial.Subject, opts...)
	cDes.SubjectAltName = canonicalizeCertificateAuthorityConfigSubjectConfigSubjectAltName(des.SubjectAltName, initial.SubjectAltName, opts...)

	return cDes
}

func canonicalizeCertificateAuthorityConfigSubjectConfigSlice(des, initial []CertificateAuthorityConfigSubjectConfig, opts ...dcl.ApplyOption) []CertificateAuthorityConfigSubjectConfig {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]CertificateAuthorityConfigSubjectConfig, 0, len(des))
		for _, d := range des {
			cd := canonicalizeCertificateAuthorityConfigSubjectConfig(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]CertificateAuthorityConfigSubjectConfig, 0, len(des))
	for i, d := range des {
		cd := canonicalizeCertificateAuthorityConfigSubjectConfig(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewCertificateAuthorityConfigSubjectConfig(c *Client, des, nw *CertificateAuthorityConfigSubjectConfig) *CertificateAuthorityConfigSubjectConfig {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for CertificateAuthorityConfigSubjectConfig while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	nw.Subject = canonicalizeNewCertificateAuthorityConfigSubjectConfigSubject(c, des.Subject, nw.Subject)
	nw.SubjectAltName = canonicalizeNewCertificateAuthorityConfigSubjectConfigSubjectAltName(c, des.SubjectAltName, nw.SubjectAltName)

	return nw
}

func canonicalizeNewCertificateAuthorityConfigSubjectConfigSet(c *Client, des, nw []CertificateAuthorityConfigSubjectConfig) []CertificateAuthorityConfigSubjectConfig {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []CertificateAuthorityConfigSubjectConfig
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareCertificateAuthorityConfigSubjectConfigNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewCertificateAuthorityConfigSubjectConfig(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewCertificateAuthorityConfigSubjectConfigSlice(c *Client, des, nw []CertificateAuthorityConfigSubjectConfig) []CertificateAuthorityConfigSubjectConfig {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []CertificateAuthorityConfigSubjectConfig
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewCertificateAuthorityConfigSubjectConfig(c, &d, &n))
	}

	return items
}

func canonicalizeCertificateAuthorityConfigSubjectConfigSubject(des, initial *CertificateAuthorityConfigSubjectConfigSubject, opts ...dcl.ApplyOption) *CertificateAuthorityConfigSubjectConfigSubject {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &CertificateAuthorityConfigSubjectConfigSubject{}

	if dcl.StringCanonicalize(des.CommonName, initial.CommonName) || dcl.IsZeroValue(des.CommonName) {
		cDes.CommonName = initial.CommonName
	} else {
		cDes.CommonName = des.CommonName
	}
	if dcl.StringCanonicalize(des.CountryCode, initial.CountryCode) || dcl.IsZeroValue(des.CountryCode) {
		cDes.CountryCode = initial.CountryCode
	} else {
		cDes.CountryCode = des.CountryCode
	}
	if dcl.StringCanonicalize(des.Organization, initial.Organization) || dcl.IsZeroValue(des.Organization) {
		cDes.Organization = initial.Organization
	} else {
		cDes.Organization = des.Organization
	}
	if dcl.StringCanonicalize(des.OrganizationalUnit, initial.OrganizationalUnit) || dcl.IsZeroValue(des.OrganizationalUnit) {
		cDes.OrganizationalUnit = initial.OrganizationalUnit
	} else {
		cDes.OrganizationalUnit = des.OrganizationalUnit
	}
	if dcl.StringCanonicalize(des.Locality, initial.Locality) || dcl.IsZeroValue(des.Locality) {
		cDes.Locality = initial.Locality
	} else {
		cDes.Locality = des.Locality
	}
	if dcl.StringCanonicalize(des.Province, initial.Province) || dcl.IsZeroValue(des.Province) {
		cDes.Province = initial.Province
	} else {
		cDes.Province = des.Province
	}
	if dcl.StringCanonicalize(des.StreetAddress, initial.StreetAddress) || dcl.IsZeroValue(des.StreetAddress) {
		cDes.StreetAddress = initial.StreetAddress
	} else {
		cDes.StreetAddress = des.StreetAddress
	}
	if dcl.StringCanonicalize(des.PostalCode, initial.PostalCode) || dcl.IsZeroValue(des.PostalCode) {
		cDes.PostalCode = initial.PostalCode
	} else {
		cDes.PostalCode = des.PostalCode
	}

	return cDes
}

func canonicalizeCertificateAuthorityConfigSubjectConfigSubjectSlice(des, initial []CertificateAuthorityConfigSubjectConfigSubject, opts ...dcl.ApplyOption) []CertificateAuthorityConfigSubjectConfigSubject {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]CertificateAuthorityConfigSubjectConfigSubject, 0, len(des))
		for _, d := range des {
			cd := canonicalizeCertificateAuthorityConfigSubjectConfigSubject(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]CertificateAuthorityConfigSubjectConfigSubject, 0, len(des))
	for i, d := range des {
		cd := canonicalizeCertificateAuthorityConfigSubjectConfigSubject(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewCertificateAuthorityConfigSubjectConfigSubject(c *Client, des, nw *CertificateAuthorityConfigSubjectConfigSubject) *CertificateAuthorityConfigSubjectConfigSubject {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for CertificateAuthorityConfigSubjectConfigSubject while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringCanonicalize(des.CommonName, nw.CommonName) {
		nw.CommonName = des.CommonName
	}
	if dcl.StringCanonicalize(des.CountryCode, nw.CountryCode) {
		nw.CountryCode = des.CountryCode
	}
	if dcl.StringCanonicalize(des.Organization, nw.Organization) {
		nw.Organization = des.Organization
	}
	if dcl.StringCanonicalize(des.OrganizationalUnit, nw.OrganizationalUnit) {
		nw.OrganizationalUnit = des.OrganizationalUnit
	}
	if dcl.StringCanonicalize(des.Locality, nw.Locality) {
		nw.Locality = des.Locality
	}
	if dcl.StringCanonicalize(des.Province, nw.Province) {
		nw.Province = des.Province
	}
	if dcl.StringCanonicalize(des.StreetAddress, nw.StreetAddress) {
		nw.StreetAddress = des.StreetAddress
	}
	if dcl.StringCanonicalize(des.PostalCode, nw.PostalCode) {
		nw.PostalCode = des.PostalCode
	}

	return nw
}

func canonicalizeNewCertificateAuthorityConfigSubjectConfigSubjectSet(c *Client, des, nw []CertificateAuthorityConfigSubjectConfigSubject) []CertificateAuthorityConfigSubjectConfigSubject {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []CertificateAuthorityConfigSubjectConfigSubject
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareCertificateAuthorityConfigSubjectConfigSubjectNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewCertificateAuthorityConfigSubjectConfigSubject(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewCertificateAuthorityConfigSubjectConfigSubjectSlice(c *Client, des, nw []CertificateAuthorityConfigSubjectConfigSubject) []CertificateAuthorityConfigSubjectConfigSubject {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []CertificateAuthorityConfigSubjectConfigSubject
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewCertificateAuthorityConfigSubjectConfigSubject(c, &d, &n))
	}

	return items
}

func canonicalizeCertificateAuthorityConfigSubjectConfigSubjectAltName(des, initial *CertificateAuthorityConfigSubjectConfigSubjectAltName, opts ...dcl.ApplyOption) *CertificateAuthorityConfigSubjectConfigSubjectAltName {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &CertificateAuthorityConfigSubjectConfigSubjectAltName{}

	if dcl.StringArrayCanonicalize(des.DnsNames, initial.DnsNames) {
		cDes.DnsNames = initial.DnsNames
	} else {
		cDes.DnsNames = des.DnsNames
	}
	if dcl.StringArrayCanonicalize(des.Uris, initial.Uris) {
		cDes.Uris = initial.Uris
	} else {
		cDes.Uris = des.Uris
	}
	if dcl.StringArrayCanonicalize(des.EmailAddresses, initial.EmailAddresses) {
		cDes.EmailAddresses = initial.EmailAddresses
	} else {
		cDes.EmailAddresses = des.EmailAddresses
	}
	if dcl.StringArrayCanonicalize(des.IPAddresses, initial.IPAddresses) {
		cDes.IPAddresses = initial.IPAddresses
	} else {
		cDes.IPAddresses = des.IPAddresses
	}
	cDes.CustomSans = canonicalizeCertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSansSlice(des.CustomSans, initial.CustomSans, opts...)

	return cDes
}

func canonicalizeCertificateAuthorityConfigSubjectConfigSubjectAltNameSlice(des, initial []CertificateAuthorityConfigSubjectConfigSubjectAltName, opts ...dcl.ApplyOption) []CertificateAuthorityConfigSubjectConfigSubjectAltName {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]CertificateAuthorityConfigSubjectConfigSubjectAltName, 0, len(des))
		for _, d := range des {
			cd := canonicalizeCertificateAuthorityConfigSubjectConfigSubjectAltName(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]CertificateAuthorityConfigSubjectConfigSubjectAltName, 0, len(des))
	for i, d := range des {
		cd := canonicalizeCertificateAuthorityConfigSubjectConfigSubjectAltName(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewCertificateAuthorityConfigSubjectConfigSubjectAltName(c *Client, des, nw *CertificateAuthorityConfigSubjectConfigSubjectAltName) *CertificateAuthorityConfigSubjectConfigSubjectAltName {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for CertificateAuthorityConfigSubjectConfigSubjectAltName while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringArrayCanonicalize(des.DnsNames, nw.DnsNames) {
		nw.DnsNames = des.DnsNames
	}
	if dcl.StringArrayCanonicalize(des.Uris, nw.Uris) {
		nw.Uris = des.Uris
	}
	if dcl.StringArrayCanonicalize(des.EmailAddresses, nw.EmailAddresses) {
		nw.EmailAddresses = des.EmailAddresses
	}
	if dcl.StringArrayCanonicalize(des.IPAddresses, nw.IPAddresses) {
		nw.IPAddresses = des.IPAddresses
	}
	nw.CustomSans = canonicalizeNewCertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSansSlice(c, des.CustomSans, nw.CustomSans)

	return nw
}

func canonicalizeNewCertificateAuthorityConfigSubjectConfigSubjectAltNameSet(c *Client, des, nw []CertificateAuthorityConfigSubjectConfigSubjectAltName) []CertificateAuthorityConfigSubjectConfigSubjectAltName {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []CertificateAuthorityConfigSubjectConfigSubjectAltName
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareCertificateAuthorityConfigSubjectConfigSubjectAltNameNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewCertificateAuthorityConfigSubjectConfigSubjectAltName(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewCertificateAuthorityConfigSubjectConfigSubjectAltNameSlice(c *Client, des, nw []CertificateAuthorityConfigSubjectConfigSubjectAltName) []CertificateAuthorityConfigSubjectConfigSubjectAltName {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []CertificateAuthorityConfigSubjectConfigSubjectAltName
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewCertificateAuthorityConfigSubjectConfigSubjectAltName(c, &d, &n))
	}

	return items
}

func canonicalizeCertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSans(des, initial *CertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSans, opts ...dcl.ApplyOption) *CertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSans {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &CertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSans{}

	cDes.ObjectId = canonicalizeCertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSansObjectId(des.ObjectId, initial.ObjectId, opts...)
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

func canonicalizeCertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSansSlice(des, initial []CertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSans, opts ...dcl.ApplyOption) []CertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSans {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]CertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSans, 0, len(des))
		for _, d := range des {
			cd := canonicalizeCertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSans(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]CertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSans, 0, len(des))
	for i, d := range des {
		cd := canonicalizeCertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSans(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewCertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSans(c *Client, des, nw *CertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSans) *CertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSans {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for CertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSans while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	nw.ObjectId = canonicalizeNewCertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSansObjectId(c, des.ObjectId, nw.ObjectId)
	if dcl.BoolCanonicalize(des.Critical, nw.Critical) {
		nw.Critical = des.Critical
	}
	if dcl.StringCanonicalize(des.Value, nw.Value) {
		nw.Value = des.Value
	}

	return nw
}

func canonicalizeNewCertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSansSet(c *Client, des, nw []CertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSans) []CertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSans {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []CertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSans
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareCertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSansNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewCertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSans(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewCertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSansSlice(c *Client, des, nw []CertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSans) []CertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSans {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []CertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSans
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewCertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSans(c, &d, &n))
	}

	return items
}

func canonicalizeCertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSansObjectId(des, initial *CertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSansObjectId, opts ...dcl.ApplyOption) *CertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSansObjectId {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &CertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSansObjectId{}

	if dcl.IsZeroValue(des.ObjectIdPath) || (dcl.IsEmptyValueIndirect(des.ObjectIdPath) && dcl.IsEmptyValueIndirect(initial.ObjectIdPath)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.ObjectIdPath = initial.ObjectIdPath
	} else {
		cDes.ObjectIdPath = des.ObjectIdPath
	}

	return cDes
}

func canonicalizeCertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSansObjectIdSlice(des, initial []CertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSansObjectId, opts ...dcl.ApplyOption) []CertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSansObjectId {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]CertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSansObjectId, 0, len(des))
		for _, d := range des {
			cd := canonicalizeCertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSansObjectId(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]CertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSansObjectId, 0, len(des))
	for i, d := range des {
		cd := canonicalizeCertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSansObjectId(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewCertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSansObjectId(c *Client, des, nw *CertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSansObjectId) *CertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSansObjectId {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for CertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSansObjectId while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	return nw
}

func canonicalizeNewCertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSansObjectIdSet(c *Client, des, nw []CertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSansObjectId) []CertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSansObjectId {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []CertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSansObjectId
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareCertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSansObjectIdNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewCertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSansObjectId(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewCertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSansObjectIdSlice(c *Client, des, nw []CertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSansObjectId) []CertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSansObjectId {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []CertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSansObjectId
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewCertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSansObjectId(c, &d, &n))
	}

	return items
}

func canonicalizeCertificateAuthorityConfigX509Config(des, initial *CertificateAuthorityConfigX509Config, opts ...dcl.ApplyOption) *CertificateAuthorityConfigX509Config {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &CertificateAuthorityConfigX509Config{}

	cDes.KeyUsage = canonicalizeCertificateAuthorityConfigX509ConfigKeyUsage(des.KeyUsage, initial.KeyUsage, opts...)
	cDes.CaOptions = canonicalizeCertificateAuthorityConfigX509ConfigCaOptions(des.CaOptions, initial.CaOptions, opts...)
	cDes.PolicyIds = canonicalizeCertificateAuthorityConfigX509ConfigPolicyIdsSlice(des.PolicyIds, initial.PolicyIds, opts...)
	cDes.AdditionalExtensions = canonicalizeCertificateAuthorityConfigX509ConfigAdditionalExtensionsSlice(des.AdditionalExtensions, initial.AdditionalExtensions, opts...)

	return cDes
}

func canonicalizeCertificateAuthorityConfigX509ConfigSlice(des, initial []CertificateAuthorityConfigX509Config, opts ...dcl.ApplyOption) []CertificateAuthorityConfigX509Config {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]CertificateAuthorityConfigX509Config, 0, len(des))
		for _, d := range des {
			cd := canonicalizeCertificateAuthorityConfigX509Config(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]CertificateAuthorityConfigX509Config, 0, len(des))
	for i, d := range des {
		cd := canonicalizeCertificateAuthorityConfigX509Config(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewCertificateAuthorityConfigX509Config(c *Client, des, nw *CertificateAuthorityConfigX509Config) *CertificateAuthorityConfigX509Config {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for CertificateAuthorityConfigX509Config while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	nw.KeyUsage = canonicalizeNewCertificateAuthorityConfigX509ConfigKeyUsage(c, des.KeyUsage, nw.KeyUsage)
	nw.CaOptions = canonicalizeNewCertificateAuthorityConfigX509ConfigCaOptions(c, des.CaOptions, nw.CaOptions)
	nw.PolicyIds = canonicalizeNewCertificateAuthorityConfigX509ConfigPolicyIdsSlice(c, des.PolicyIds, nw.PolicyIds)
	if dcl.StringArrayCanonicalize(des.AiaOcspServers, nw.AiaOcspServers) {
		nw.AiaOcspServers = des.AiaOcspServers
	}
	nw.AdditionalExtensions = canonicalizeNewCertificateAuthorityConfigX509ConfigAdditionalExtensionsSlice(c, des.AdditionalExtensions, nw.AdditionalExtensions)

	return nw
}

func canonicalizeNewCertificateAuthorityConfigX509ConfigSet(c *Client, des, nw []CertificateAuthorityConfigX509Config) []CertificateAuthorityConfigX509Config {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []CertificateAuthorityConfigX509Config
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareCertificateAuthorityConfigX509ConfigNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewCertificateAuthorityConfigX509Config(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewCertificateAuthorityConfigX509ConfigSlice(c *Client, des, nw []CertificateAuthorityConfigX509Config) []CertificateAuthorityConfigX509Config {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []CertificateAuthorityConfigX509Config
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewCertificateAuthorityConfigX509Config(c, &d, &n))
	}

	return items
}

func canonicalizeCertificateAuthorityConfigX509ConfigKeyUsage(des, initial *CertificateAuthorityConfigX509ConfigKeyUsage, opts ...dcl.ApplyOption) *CertificateAuthorityConfigX509ConfigKeyUsage {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &CertificateAuthorityConfigX509ConfigKeyUsage{}

	cDes.BaseKeyUsage = canonicalizeCertificateAuthorityConfigX509ConfigKeyUsageBaseKeyUsage(des.BaseKeyUsage, initial.BaseKeyUsage, opts...)
	cDes.ExtendedKeyUsage = canonicalizeCertificateAuthorityConfigX509ConfigKeyUsageExtendedKeyUsage(des.ExtendedKeyUsage, initial.ExtendedKeyUsage, opts...)
	cDes.UnknownExtendedKeyUsages = canonicalizeCertificateAuthorityConfigX509ConfigKeyUsageUnknownExtendedKeyUsagesSlice(des.UnknownExtendedKeyUsages, initial.UnknownExtendedKeyUsages, opts...)

	return cDes
}

func canonicalizeCertificateAuthorityConfigX509ConfigKeyUsageSlice(des, initial []CertificateAuthorityConfigX509ConfigKeyUsage, opts ...dcl.ApplyOption) []CertificateAuthorityConfigX509ConfigKeyUsage {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]CertificateAuthorityConfigX509ConfigKeyUsage, 0, len(des))
		for _, d := range des {
			cd := canonicalizeCertificateAuthorityConfigX509ConfigKeyUsage(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]CertificateAuthorityConfigX509ConfigKeyUsage, 0, len(des))
	for i, d := range des {
		cd := canonicalizeCertificateAuthorityConfigX509ConfigKeyUsage(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewCertificateAuthorityConfigX509ConfigKeyUsage(c *Client, des, nw *CertificateAuthorityConfigX509ConfigKeyUsage) *CertificateAuthorityConfigX509ConfigKeyUsage {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for CertificateAuthorityConfigX509ConfigKeyUsage while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	nw.BaseKeyUsage = canonicalizeNewCertificateAuthorityConfigX509ConfigKeyUsageBaseKeyUsage(c, des.BaseKeyUsage, nw.BaseKeyUsage)
	nw.ExtendedKeyUsage = canonicalizeNewCertificateAuthorityConfigX509ConfigKeyUsageExtendedKeyUsage(c, des.ExtendedKeyUsage, nw.ExtendedKeyUsage)
	nw.UnknownExtendedKeyUsages = canonicalizeNewCertificateAuthorityConfigX509ConfigKeyUsageUnknownExtendedKeyUsagesSlice(c, des.UnknownExtendedKeyUsages, nw.UnknownExtendedKeyUsages)

	return nw
}

func canonicalizeNewCertificateAuthorityConfigX509ConfigKeyUsageSet(c *Client, des, nw []CertificateAuthorityConfigX509ConfigKeyUsage) []CertificateAuthorityConfigX509ConfigKeyUsage {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []CertificateAuthorityConfigX509ConfigKeyUsage
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareCertificateAuthorityConfigX509ConfigKeyUsageNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewCertificateAuthorityConfigX509ConfigKeyUsage(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewCertificateAuthorityConfigX509ConfigKeyUsageSlice(c *Client, des, nw []CertificateAuthorityConfigX509ConfigKeyUsage) []CertificateAuthorityConfigX509ConfigKeyUsage {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []CertificateAuthorityConfigX509ConfigKeyUsage
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewCertificateAuthorityConfigX509ConfigKeyUsage(c, &d, &n))
	}

	return items
}

func canonicalizeCertificateAuthorityConfigX509ConfigKeyUsageBaseKeyUsage(des, initial *CertificateAuthorityConfigX509ConfigKeyUsageBaseKeyUsage, opts ...dcl.ApplyOption) *CertificateAuthorityConfigX509ConfigKeyUsageBaseKeyUsage {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &CertificateAuthorityConfigX509ConfigKeyUsageBaseKeyUsage{}

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

func canonicalizeCertificateAuthorityConfigX509ConfigKeyUsageBaseKeyUsageSlice(des, initial []CertificateAuthorityConfigX509ConfigKeyUsageBaseKeyUsage, opts ...dcl.ApplyOption) []CertificateAuthorityConfigX509ConfigKeyUsageBaseKeyUsage {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]CertificateAuthorityConfigX509ConfigKeyUsageBaseKeyUsage, 0, len(des))
		for _, d := range des {
			cd := canonicalizeCertificateAuthorityConfigX509ConfigKeyUsageBaseKeyUsage(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]CertificateAuthorityConfigX509ConfigKeyUsageBaseKeyUsage, 0, len(des))
	for i, d := range des {
		cd := canonicalizeCertificateAuthorityConfigX509ConfigKeyUsageBaseKeyUsage(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewCertificateAuthorityConfigX509ConfigKeyUsageBaseKeyUsage(c *Client, des, nw *CertificateAuthorityConfigX509ConfigKeyUsageBaseKeyUsage) *CertificateAuthorityConfigX509ConfigKeyUsageBaseKeyUsage {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for CertificateAuthorityConfigX509ConfigKeyUsageBaseKeyUsage while comparing non-nil desired to nil actual.  Returning desired object.")
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

func canonicalizeNewCertificateAuthorityConfigX509ConfigKeyUsageBaseKeyUsageSet(c *Client, des, nw []CertificateAuthorityConfigX509ConfigKeyUsageBaseKeyUsage) []CertificateAuthorityConfigX509ConfigKeyUsageBaseKeyUsage {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []CertificateAuthorityConfigX509ConfigKeyUsageBaseKeyUsage
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareCertificateAuthorityConfigX509ConfigKeyUsageBaseKeyUsageNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewCertificateAuthorityConfigX509ConfigKeyUsageBaseKeyUsage(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewCertificateAuthorityConfigX509ConfigKeyUsageBaseKeyUsageSlice(c *Client, des, nw []CertificateAuthorityConfigX509ConfigKeyUsageBaseKeyUsage) []CertificateAuthorityConfigX509ConfigKeyUsageBaseKeyUsage {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []CertificateAuthorityConfigX509ConfigKeyUsageBaseKeyUsage
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewCertificateAuthorityConfigX509ConfigKeyUsageBaseKeyUsage(c, &d, &n))
	}

	return items
}

func canonicalizeCertificateAuthorityConfigX509ConfigKeyUsageExtendedKeyUsage(des, initial *CertificateAuthorityConfigX509ConfigKeyUsageExtendedKeyUsage, opts ...dcl.ApplyOption) *CertificateAuthorityConfigX509ConfigKeyUsageExtendedKeyUsage {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &CertificateAuthorityConfigX509ConfigKeyUsageExtendedKeyUsage{}

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

func canonicalizeCertificateAuthorityConfigX509ConfigKeyUsageExtendedKeyUsageSlice(des, initial []CertificateAuthorityConfigX509ConfigKeyUsageExtendedKeyUsage, opts ...dcl.ApplyOption) []CertificateAuthorityConfigX509ConfigKeyUsageExtendedKeyUsage {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]CertificateAuthorityConfigX509ConfigKeyUsageExtendedKeyUsage, 0, len(des))
		for _, d := range des {
			cd := canonicalizeCertificateAuthorityConfigX509ConfigKeyUsageExtendedKeyUsage(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]CertificateAuthorityConfigX509ConfigKeyUsageExtendedKeyUsage, 0, len(des))
	for i, d := range des {
		cd := canonicalizeCertificateAuthorityConfigX509ConfigKeyUsageExtendedKeyUsage(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewCertificateAuthorityConfigX509ConfigKeyUsageExtendedKeyUsage(c *Client, des, nw *CertificateAuthorityConfigX509ConfigKeyUsageExtendedKeyUsage) *CertificateAuthorityConfigX509ConfigKeyUsageExtendedKeyUsage {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for CertificateAuthorityConfigX509ConfigKeyUsageExtendedKeyUsage while comparing non-nil desired to nil actual.  Returning desired object.")
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

func canonicalizeNewCertificateAuthorityConfigX509ConfigKeyUsageExtendedKeyUsageSet(c *Client, des, nw []CertificateAuthorityConfigX509ConfigKeyUsageExtendedKeyUsage) []CertificateAuthorityConfigX509ConfigKeyUsageExtendedKeyUsage {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []CertificateAuthorityConfigX509ConfigKeyUsageExtendedKeyUsage
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareCertificateAuthorityConfigX509ConfigKeyUsageExtendedKeyUsageNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewCertificateAuthorityConfigX509ConfigKeyUsageExtendedKeyUsage(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewCertificateAuthorityConfigX509ConfigKeyUsageExtendedKeyUsageSlice(c *Client, des, nw []CertificateAuthorityConfigX509ConfigKeyUsageExtendedKeyUsage) []CertificateAuthorityConfigX509ConfigKeyUsageExtendedKeyUsage {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []CertificateAuthorityConfigX509ConfigKeyUsageExtendedKeyUsage
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewCertificateAuthorityConfigX509ConfigKeyUsageExtendedKeyUsage(c, &d, &n))
	}

	return items
}

func canonicalizeCertificateAuthorityConfigX509ConfigKeyUsageUnknownExtendedKeyUsages(des, initial *CertificateAuthorityConfigX509ConfigKeyUsageUnknownExtendedKeyUsages, opts ...dcl.ApplyOption) *CertificateAuthorityConfigX509ConfigKeyUsageUnknownExtendedKeyUsages {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &CertificateAuthorityConfigX509ConfigKeyUsageUnknownExtendedKeyUsages{}

	if dcl.IsZeroValue(des.ObjectIdPath) || (dcl.IsEmptyValueIndirect(des.ObjectIdPath) && dcl.IsEmptyValueIndirect(initial.ObjectIdPath)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.ObjectIdPath = initial.ObjectIdPath
	} else {
		cDes.ObjectIdPath = des.ObjectIdPath
	}

	return cDes
}

func canonicalizeCertificateAuthorityConfigX509ConfigKeyUsageUnknownExtendedKeyUsagesSlice(des, initial []CertificateAuthorityConfigX509ConfigKeyUsageUnknownExtendedKeyUsages, opts ...dcl.ApplyOption) []CertificateAuthorityConfigX509ConfigKeyUsageUnknownExtendedKeyUsages {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]CertificateAuthorityConfigX509ConfigKeyUsageUnknownExtendedKeyUsages, 0, len(des))
		for _, d := range des {
			cd := canonicalizeCertificateAuthorityConfigX509ConfigKeyUsageUnknownExtendedKeyUsages(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]CertificateAuthorityConfigX509ConfigKeyUsageUnknownExtendedKeyUsages, 0, len(des))
	for i, d := range des {
		cd := canonicalizeCertificateAuthorityConfigX509ConfigKeyUsageUnknownExtendedKeyUsages(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewCertificateAuthorityConfigX509ConfigKeyUsageUnknownExtendedKeyUsages(c *Client, des, nw *CertificateAuthorityConfigX509ConfigKeyUsageUnknownExtendedKeyUsages) *CertificateAuthorityConfigX509ConfigKeyUsageUnknownExtendedKeyUsages {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for CertificateAuthorityConfigX509ConfigKeyUsageUnknownExtendedKeyUsages while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	return nw
}

func canonicalizeNewCertificateAuthorityConfigX509ConfigKeyUsageUnknownExtendedKeyUsagesSet(c *Client, des, nw []CertificateAuthorityConfigX509ConfigKeyUsageUnknownExtendedKeyUsages) []CertificateAuthorityConfigX509ConfigKeyUsageUnknownExtendedKeyUsages {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []CertificateAuthorityConfigX509ConfigKeyUsageUnknownExtendedKeyUsages
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareCertificateAuthorityConfigX509ConfigKeyUsageUnknownExtendedKeyUsagesNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewCertificateAuthorityConfigX509ConfigKeyUsageUnknownExtendedKeyUsages(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewCertificateAuthorityConfigX509ConfigKeyUsageUnknownExtendedKeyUsagesSlice(c *Client, des, nw []CertificateAuthorityConfigX509ConfigKeyUsageUnknownExtendedKeyUsages) []CertificateAuthorityConfigX509ConfigKeyUsageUnknownExtendedKeyUsages {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []CertificateAuthorityConfigX509ConfigKeyUsageUnknownExtendedKeyUsages
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewCertificateAuthorityConfigX509ConfigKeyUsageUnknownExtendedKeyUsages(c, &d, &n))
	}

	return items
}

func canonicalizeCertificateAuthorityConfigX509ConfigCaOptions(des, initial *CertificateAuthorityConfigX509ConfigCaOptions, opts ...dcl.ApplyOption) *CertificateAuthorityConfigX509ConfigCaOptions {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &CertificateAuthorityConfigX509ConfigCaOptions{}

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

func canonicalizeCertificateAuthorityConfigX509ConfigCaOptionsSlice(des, initial []CertificateAuthorityConfigX509ConfigCaOptions, opts ...dcl.ApplyOption) []CertificateAuthorityConfigX509ConfigCaOptions {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]CertificateAuthorityConfigX509ConfigCaOptions, 0, len(des))
		for _, d := range des {
			cd := canonicalizeCertificateAuthorityConfigX509ConfigCaOptions(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]CertificateAuthorityConfigX509ConfigCaOptions, 0, len(des))
	for i, d := range des {
		cd := canonicalizeCertificateAuthorityConfigX509ConfigCaOptions(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewCertificateAuthorityConfigX509ConfigCaOptions(c *Client, des, nw *CertificateAuthorityConfigX509ConfigCaOptions) *CertificateAuthorityConfigX509ConfigCaOptions {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for CertificateAuthorityConfigX509ConfigCaOptions while comparing non-nil desired to nil actual.  Returning desired object.")
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

func canonicalizeNewCertificateAuthorityConfigX509ConfigCaOptionsSet(c *Client, des, nw []CertificateAuthorityConfigX509ConfigCaOptions) []CertificateAuthorityConfigX509ConfigCaOptions {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []CertificateAuthorityConfigX509ConfigCaOptions
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareCertificateAuthorityConfigX509ConfigCaOptionsNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewCertificateAuthorityConfigX509ConfigCaOptions(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewCertificateAuthorityConfigX509ConfigCaOptionsSlice(c *Client, des, nw []CertificateAuthorityConfigX509ConfigCaOptions) []CertificateAuthorityConfigX509ConfigCaOptions {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []CertificateAuthorityConfigX509ConfigCaOptions
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewCertificateAuthorityConfigX509ConfigCaOptions(c, &d, &n))
	}

	return items
}

func canonicalizeCertificateAuthorityConfigX509ConfigPolicyIds(des, initial *CertificateAuthorityConfigX509ConfigPolicyIds, opts ...dcl.ApplyOption) *CertificateAuthorityConfigX509ConfigPolicyIds {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &CertificateAuthorityConfigX509ConfigPolicyIds{}

	if dcl.IsZeroValue(des.ObjectIdPath) || (dcl.IsEmptyValueIndirect(des.ObjectIdPath) && dcl.IsEmptyValueIndirect(initial.ObjectIdPath)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.ObjectIdPath = initial.ObjectIdPath
	} else {
		cDes.ObjectIdPath = des.ObjectIdPath
	}

	return cDes
}

func canonicalizeCertificateAuthorityConfigX509ConfigPolicyIdsSlice(des, initial []CertificateAuthorityConfigX509ConfigPolicyIds, opts ...dcl.ApplyOption) []CertificateAuthorityConfigX509ConfigPolicyIds {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]CertificateAuthorityConfigX509ConfigPolicyIds, 0, len(des))
		for _, d := range des {
			cd := canonicalizeCertificateAuthorityConfigX509ConfigPolicyIds(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]CertificateAuthorityConfigX509ConfigPolicyIds, 0, len(des))
	for i, d := range des {
		cd := canonicalizeCertificateAuthorityConfigX509ConfigPolicyIds(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewCertificateAuthorityConfigX509ConfigPolicyIds(c *Client, des, nw *CertificateAuthorityConfigX509ConfigPolicyIds) *CertificateAuthorityConfigX509ConfigPolicyIds {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for CertificateAuthorityConfigX509ConfigPolicyIds while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	return nw
}

func canonicalizeNewCertificateAuthorityConfigX509ConfigPolicyIdsSet(c *Client, des, nw []CertificateAuthorityConfigX509ConfigPolicyIds) []CertificateAuthorityConfigX509ConfigPolicyIds {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []CertificateAuthorityConfigX509ConfigPolicyIds
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareCertificateAuthorityConfigX509ConfigPolicyIdsNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewCertificateAuthorityConfigX509ConfigPolicyIds(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewCertificateAuthorityConfigX509ConfigPolicyIdsSlice(c *Client, des, nw []CertificateAuthorityConfigX509ConfigPolicyIds) []CertificateAuthorityConfigX509ConfigPolicyIds {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []CertificateAuthorityConfigX509ConfigPolicyIds
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewCertificateAuthorityConfigX509ConfigPolicyIds(c, &d, &n))
	}

	return items
}

func canonicalizeCertificateAuthorityConfigX509ConfigAdditionalExtensions(des, initial *CertificateAuthorityConfigX509ConfigAdditionalExtensions, opts ...dcl.ApplyOption) *CertificateAuthorityConfigX509ConfigAdditionalExtensions {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &CertificateAuthorityConfigX509ConfigAdditionalExtensions{}

	cDes.ObjectId = canonicalizeCertificateAuthorityConfigX509ConfigAdditionalExtensionsObjectId(des.ObjectId, initial.ObjectId, opts...)
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

func canonicalizeCertificateAuthorityConfigX509ConfigAdditionalExtensionsSlice(des, initial []CertificateAuthorityConfigX509ConfigAdditionalExtensions, opts ...dcl.ApplyOption) []CertificateAuthorityConfigX509ConfigAdditionalExtensions {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]CertificateAuthorityConfigX509ConfigAdditionalExtensions, 0, len(des))
		for _, d := range des {
			cd := canonicalizeCertificateAuthorityConfigX509ConfigAdditionalExtensions(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]CertificateAuthorityConfigX509ConfigAdditionalExtensions, 0, len(des))
	for i, d := range des {
		cd := canonicalizeCertificateAuthorityConfigX509ConfigAdditionalExtensions(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewCertificateAuthorityConfigX509ConfigAdditionalExtensions(c *Client, des, nw *CertificateAuthorityConfigX509ConfigAdditionalExtensions) *CertificateAuthorityConfigX509ConfigAdditionalExtensions {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for CertificateAuthorityConfigX509ConfigAdditionalExtensions while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	nw.ObjectId = canonicalizeNewCertificateAuthorityConfigX509ConfigAdditionalExtensionsObjectId(c, des.ObjectId, nw.ObjectId)
	if dcl.BoolCanonicalize(des.Critical, nw.Critical) {
		nw.Critical = des.Critical
	}
	if dcl.StringCanonicalize(des.Value, nw.Value) {
		nw.Value = des.Value
	}

	return nw
}

func canonicalizeNewCertificateAuthorityConfigX509ConfigAdditionalExtensionsSet(c *Client, des, nw []CertificateAuthorityConfigX509ConfigAdditionalExtensions) []CertificateAuthorityConfigX509ConfigAdditionalExtensions {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []CertificateAuthorityConfigX509ConfigAdditionalExtensions
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareCertificateAuthorityConfigX509ConfigAdditionalExtensionsNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewCertificateAuthorityConfigX509ConfigAdditionalExtensions(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewCertificateAuthorityConfigX509ConfigAdditionalExtensionsSlice(c *Client, des, nw []CertificateAuthorityConfigX509ConfigAdditionalExtensions) []CertificateAuthorityConfigX509ConfigAdditionalExtensions {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []CertificateAuthorityConfigX509ConfigAdditionalExtensions
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewCertificateAuthorityConfigX509ConfigAdditionalExtensions(c, &d, &n))
	}

	return items
}

func canonicalizeCertificateAuthorityConfigX509ConfigAdditionalExtensionsObjectId(des, initial *CertificateAuthorityConfigX509ConfigAdditionalExtensionsObjectId, opts ...dcl.ApplyOption) *CertificateAuthorityConfigX509ConfigAdditionalExtensionsObjectId {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &CertificateAuthorityConfigX509ConfigAdditionalExtensionsObjectId{}

	if dcl.IsZeroValue(des.ObjectIdPath) || (dcl.IsEmptyValueIndirect(des.ObjectIdPath) && dcl.IsEmptyValueIndirect(initial.ObjectIdPath)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.ObjectIdPath = initial.ObjectIdPath
	} else {
		cDes.ObjectIdPath = des.ObjectIdPath
	}

	return cDes
}

func canonicalizeCertificateAuthorityConfigX509ConfigAdditionalExtensionsObjectIdSlice(des, initial []CertificateAuthorityConfigX509ConfigAdditionalExtensionsObjectId, opts ...dcl.ApplyOption) []CertificateAuthorityConfigX509ConfigAdditionalExtensionsObjectId {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]CertificateAuthorityConfigX509ConfigAdditionalExtensionsObjectId, 0, len(des))
		for _, d := range des {
			cd := canonicalizeCertificateAuthorityConfigX509ConfigAdditionalExtensionsObjectId(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]CertificateAuthorityConfigX509ConfigAdditionalExtensionsObjectId, 0, len(des))
	for i, d := range des {
		cd := canonicalizeCertificateAuthorityConfigX509ConfigAdditionalExtensionsObjectId(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewCertificateAuthorityConfigX509ConfigAdditionalExtensionsObjectId(c *Client, des, nw *CertificateAuthorityConfigX509ConfigAdditionalExtensionsObjectId) *CertificateAuthorityConfigX509ConfigAdditionalExtensionsObjectId {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for CertificateAuthorityConfigX509ConfigAdditionalExtensionsObjectId while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	return nw
}

func canonicalizeNewCertificateAuthorityConfigX509ConfigAdditionalExtensionsObjectIdSet(c *Client, des, nw []CertificateAuthorityConfigX509ConfigAdditionalExtensionsObjectId) []CertificateAuthorityConfigX509ConfigAdditionalExtensionsObjectId {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []CertificateAuthorityConfigX509ConfigAdditionalExtensionsObjectId
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareCertificateAuthorityConfigX509ConfigAdditionalExtensionsObjectIdNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewCertificateAuthorityConfigX509ConfigAdditionalExtensionsObjectId(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewCertificateAuthorityConfigX509ConfigAdditionalExtensionsObjectIdSlice(c *Client, des, nw []CertificateAuthorityConfigX509ConfigAdditionalExtensionsObjectId) []CertificateAuthorityConfigX509ConfigAdditionalExtensionsObjectId {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []CertificateAuthorityConfigX509ConfigAdditionalExtensionsObjectId
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewCertificateAuthorityConfigX509ConfigAdditionalExtensionsObjectId(c, &d, &n))
	}

	return items
}

func canonicalizeCertificateAuthorityConfigPublicKey(des, initial *CertificateAuthorityConfigPublicKey, opts ...dcl.ApplyOption) *CertificateAuthorityConfigPublicKey {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &CertificateAuthorityConfigPublicKey{}

	if dcl.StringCanonicalize(des.Key, initial.Key) || dcl.IsZeroValue(des.Key) {
		cDes.Key = initial.Key
	} else {
		cDes.Key = des.Key
	}
	if dcl.IsZeroValue(des.Format) || (dcl.IsEmptyValueIndirect(des.Format) && dcl.IsEmptyValueIndirect(initial.Format)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.Format = initial.Format
	} else {
		cDes.Format = des.Format
	}

	return cDes
}

func canonicalizeCertificateAuthorityConfigPublicKeySlice(des, initial []CertificateAuthorityConfigPublicKey, opts ...dcl.ApplyOption) []CertificateAuthorityConfigPublicKey {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]CertificateAuthorityConfigPublicKey, 0, len(des))
		for _, d := range des {
			cd := canonicalizeCertificateAuthorityConfigPublicKey(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]CertificateAuthorityConfigPublicKey, 0, len(des))
	for i, d := range des {
		cd := canonicalizeCertificateAuthorityConfigPublicKey(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewCertificateAuthorityConfigPublicKey(c *Client, des, nw *CertificateAuthorityConfigPublicKey) *CertificateAuthorityConfigPublicKey {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for CertificateAuthorityConfigPublicKey while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringCanonicalize(des.Key, nw.Key) {
		nw.Key = des.Key
	}

	return nw
}

func canonicalizeNewCertificateAuthorityConfigPublicKeySet(c *Client, des, nw []CertificateAuthorityConfigPublicKey) []CertificateAuthorityConfigPublicKey {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []CertificateAuthorityConfigPublicKey
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareCertificateAuthorityConfigPublicKeyNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewCertificateAuthorityConfigPublicKey(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewCertificateAuthorityConfigPublicKeySlice(c *Client, des, nw []CertificateAuthorityConfigPublicKey) []CertificateAuthorityConfigPublicKey {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []CertificateAuthorityConfigPublicKey
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewCertificateAuthorityConfigPublicKey(c, &d, &n))
	}

	return items
}

func canonicalizeCertificateAuthorityKeySpec(des, initial *CertificateAuthorityKeySpec, opts ...dcl.ApplyOption) *CertificateAuthorityKeySpec {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if des.CloudKmsKeyVersion != nil || (initial != nil && initial.CloudKmsKeyVersion != nil) {
		// Check if anything else is set.
		if dcl.AnySet(des.Algorithm) {
			des.CloudKmsKeyVersion = nil
			if initial != nil {
				initial.CloudKmsKeyVersion = nil
			}
		}
	}

	if des.Algorithm != nil || (initial != nil && initial.Algorithm != nil) {
		// Check if anything else is set.
		if dcl.AnySet(des.CloudKmsKeyVersion) {
			des.Algorithm = nil
			if initial != nil {
				initial.Algorithm = nil
			}
		}
	}

	if initial == nil {
		return des
	}

	cDes := &CertificateAuthorityKeySpec{}

	if dcl.IsZeroValue(des.CloudKmsKeyVersion) || (dcl.IsEmptyValueIndirect(des.CloudKmsKeyVersion) && dcl.IsEmptyValueIndirect(initial.CloudKmsKeyVersion)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.CloudKmsKeyVersion = initial.CloudKmsKeyVersion
	} else {
		cDes.CloudKmsKeyVersion = des.CloudKmsKeyVersion
	}
	if dcl.IsZeroValue(des.Algorithm) || (dcl.IsEmptyValueIndirect(des.Algorithm) && dcl.IsEmptyValueIndirect(initial.Algorithm)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.Algorithm = initial.Algorithm
	} else {
		cDes.Algorithm = des.Algorithm
	}

	return cDes
}

func canonicalizeCertificateAuthorityKeySpecSlice(des, initial []CertificateAuthorityKeySpec, opts ...dcl.ApplyOption) []CertificateAuthorityKeySpec {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]CertificateAuthorityKeySpec, 0, len(des))
		for _, d := range des {
			cd := canonicalizeCertificateAuthorityKeySpec(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]CertificateAuthorityKeySpec, 0, len(des))
	for i, d := range des {
		cd := canonicalizeCertificateAuthorityKeySpec(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewCertificateAuthorityKeySpec(c *Client, des, nw *CertificateAuthorityKeySpec) *CertificateAuthorityKeySpec {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for CertificateAuthorityKeySpec while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	return nw
}

func canonicalizeNewCertificateAuthorityKeySpecSet(c *Client, des, nw []CertificateAuthorityKeySpec) []CertificateAuthorityKeySpec {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []CertificateAuthorityKeySpec
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareCertificateAuthorityKeySpecNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewCertificateAuthorityKeySpec(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewCertificateAuthorityKeySpecSlice(c *Client, des, nw []CertificateAuthorityKeySpec) []CertificateAuthorityKeySpec {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []CertificateAuthorityKeySpec
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewCertificateAuthorityKeySpec(c, &d, &n))
	}

	return items
}

func canonicalizeCertificateAuthoritySubordinateConfig(des, initial *CertificateAuthoritySubordinateConfig, opts ...dcl.ApplyOption) *CertificateAuthoritySubordinateConfig {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if des.CertificateAuthority != nil || (initial != nil && initial.CertificateAuthority != nil) {
		// Check if anything else is set.
		if dcl.AnySet(des.PemIssuerChain) {
			des.CertificateAuthority = nil
			if initial != nil {
				initial.CertificateAuthority = nil
			}
		}
	}

	if des.PemIssuerChain != nil || (initial != nil && initial.PemIssuerChain != nil) {
		// Check if anything else is set.
		if dcl.AnySet(des.CertificateAuthority) {
			des.PemIssuerChain = nil
			if initial != nil {
				initial.PemIssuerChain = nil
			}
		}
	}

	if initial == nil {
		return des
	}

	cDes := &CertificateAuthoritySubordinateConfig{}

	if dcl.IsZeroValue(des.CertificateAuthority) || (dcl.IsEmptyValueIndirect(des.CertificateAuthority) && dcl.IsEmptyValueIndirect(initial.CertificateAuthority)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.CertificateAuthority = initial.CertificateAuthority
	} else {
		cDes.CertificateAuthority = des.CertificateAuthority
	}
	cDes.PemIssuerChain = canonicalizeCertificateAuthoritySubordinateConfigPemIssuerChain(des.PemIssuerChain, initial.PemIssuerChain, opts...)

	return cDes
}

func canonicalizeCertificateAuthoritySubordinateConfigSlice(des, initial []CertificateAuthoritySubordinateConfig, opts ...dcl.ApplyOption) []CertificateAuthoritySubordinateConfig {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]CertificateAuthoritySubordinateConfig, 0, len(des))
		for _, d := range des {
			cd := canonicalizeCertificateAuthoritySubordinateConfig(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]CertificateAuthoritySubordinateConfig, 0, len(des))
	for i, d := range des {
		cd := canonicalizeCertificateAuthoritySubordinateConfig(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewCertificateAuthoritySubordinateConfig(c *Client, des, nw *CertificateAuthoritySubordinateConfig) *CertificateAuthoritySubordinateConfig {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for CertificateAuthoritySubordinateConfig while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	nw.PemIssuerChain = canonicalizeNewCertificateAuthoritySubordinateConfigPemIssuerChain(c, des.PemIssuerChain, nw.PemIssuerChain)

	return nw
}

func canonicalizeNewCertificateAuthoritySubordinateConfigSet(c *Client, des, nw []CertificateAuthoritySubordinateConfig) []CertificateAuthoritySubordinateConfig {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []CertificateAuthoritySubordinateConfig
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareCertificateAuthoritySubordinateConfigNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewCertificateAuthoritySubordinateConfig(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewCertificateAuthoritySubordinateConfigSlice(c *Client, des, nw []CertificateAuthoritySubordinateConfig) []CertificateAuthoritySubordinateConfig {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []CertificateAuthoritySubordinateConfig
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewCertificateAuthoritySubordinateConfig(c, &d, &n))
	}

	return items
}

func canonicalizeCertificateAuthoritySubordinateConfigPemIssuerChain(des, initial *CertificateAuthoritySubordinateConfigPemIssuerChain, opts ...dcl.ApplyOption) *CertificateAuthoritySubordinateConfigPemIssuerChain {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &CertificateAuthoritySubordinateConfigPemIssuerChain{}

	if dcl.StringArrayCanonicalize(des.PemCertificates, initial.PemCertificates) {
		cDes.PemCertificates = initial.PemCertificates
	} else {
		cDes.PemCertificates = des.PemCertificates
	}

	return cDes
}

func canonicalizeCertificateAuthoritySubordinateConfigPemIssuerChainSlice(des, initial []CertificateAuthoritySubordinateConfigPemIssuerChain, opts ...dcl.ApplyOption) []CertificateAuthoritySubordinateConfigPemIssuerChain {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]CertificateAuthoritySubordinateConfigPemIssuerChain, 0, len(des))
		for _, d := range des {
			cd := canonicalizeCertificateAuthoritySubordinateConfigPemIssuerChain(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]CertificateAuthoritySubordinateConfigPemIssuerChain, 0, len(des))
	for i, d := range des {
		cd := canonicalizeCertificateAuthoritySubordinateConfigPemIssuerChain(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewCertificateAuthoritySubordinateConfigPemIssuerChain(c *Client, des, nw *CertificateAuthoritySubordinateConfigPemIssuerChain) *CertificateAuthoritySubordinateConfigPemIssuerChain {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for CertificateAuthoritySubordinateConfigPemIssuerChain while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringArrayCanonicalize(des.PemCertificates, nw.PemCertificates) {
		nw.PemCertificates = des.PemCertificates
	}

	return nw
}

func canonicalizeNewCertificateAuthoritySubordinateConfigPemIssuerChainSet(c *Client, des, nw []CertificateAuthoritySubordinateConfigPemIssuerChain) []CertificateAuthoritySubordinateConfigPemIssuerChain {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []CertificateAuthoritySubordinateConfigPemIssuerChain
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareCertificateAuthoritySubordinateConfigPemIssuerChainNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewCertificateAuthoritySubordinateConfigPemIssuerChain(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewCertificateAuthoritySubordinateConfigPemIssuerChainSlice(c *Client, des, nw []CertificateAuthoritySubordinateConfigPemIssuerChain) []CertificateAuthoritySubordinateConfigPemIssuerChain {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []CertificateAuthoritySubordinateConfigPemIssuerChain
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewCertificateAuthoritySubordinateConfigPemIssuerChain(c, &d, &n))
	}

	return items
}

func canonicalizeCertificateAuthorityCaCertificateDescriptions(des, initial *CertificateAuthorityCaCertificateDescriptions, opts ...dcl.ApplyOption) *CertificateAuthorityCaCertificateDescriptions {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &CertificateAuthorityCaCertificateDescriptions{}

	cDes.SubjectDescription = canonicalizeCertificateAuthorityCaCertificateDescriptionsSubjectDescription(des.SubjectDescription, initial.SubjectDescription, opts...)
	cDes.X509Description = canonicalizeCertificateAuthorityCaCertificateDescriptionsX509Description(des.X509Description, initial.X509Description, opts...)
	cDes.PublicKey = canonicalizeCertificateAuthorityCaCertificateDescriptionsPublicKey(des.PublicKey, initial.PublicKey, opts...)
	cDes.SubjectKeyId = canonicalizeCertificateAuthorityCaCertificateDescriptionsSubjectKeyId(des.SubjectKeyId, initial.SubjectKeyId, opts...)
	cDes.AuthorityKeyId = canonicalizeCertificateAuthorityCaCertificateDescriptionsAuthorityKeyId(des.AuthorityKeyId, initial.AuthorityKeyId, opts...)
	if dcl.StringArrayCanonicalize(des.CrlDistributionPoints, initial.CrlDistributionPoints) {
		cDes.CrlDistributionPoints = initial.CrlDistributionPoints
	} else {
		cDes.CrlDistributionPoints = des.CrlDistributionPoints
	}
	if dcl.StringArrayCanonicalize(des.AiaIssuingCertificateUrls, initial.AiaIssuingCertificateUrls) {
		cDes.AiaIssuingCertificateUrls = initial.AiaIssuingCertificateUrls
	} else {
		cDes.AiaIssuingCertificateUrls = des.AiaIssuingCertificateUrls
	}
	cDes.CertFingerprint = canonicalizeCertificateAuthorityCaCertificateDescriptionsCertFingerprint(des.CertFingerprint, initial.CertFingerprint, opts...)

	return cDes
}

func canonicalizeCertificateAuthorityCaCertificateDescriptionsSlice(des, initial []CertificateAuthorityCaCertificateDescriptions, opts ...dcl.ApplyOption) []CertificateAuthorityCaCertificateDescriptions {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]CertificateAuthorityCaCertificateDescriptions, 0, len(des))
		for _, d := range des {
			cd := canonicalizeCertificateAuthorityCaCertificateDescriptions(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]CertificateAuthorityCaCertificateDescriptions, 0, len(des))
	for i, d := range des {
		cd := canonicalizeCertificateAuthorityCaCertificateDescriptions(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewCertificateAuthorityCaCertificateDescriptions(c *Client, des, nw *CertificateAuthorityCaCertificateDescriptions) *CertificateAuthorityCaCertificateDescriptions {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for CertificateAuthorityCaCertificateDescriptions while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	nw.SubjectDescription = canonicalizeNewCertificateAuthorityCaCertificateDescriptionsSubjectDescription(c, des.SubjectDescription, nw.SubjectDescription)
	nw.X509Description = canonicalizeNewCertificateAuthorityCaCertificateDescriptionsX509Description(c, des.X509Description, nw.X509Description)
	nw.PublicKey = canonicalizeNewCertificateAuthorityCaCertificateDescriptionsPublicKey(c, des.PublicKey, nw.PublicKey)
	nw.SubjectKeyId = canonicalizeNewCertificateAuthorityCaCertificateDescriptionsSubjectKeyId(c, des.SubjectKeyId, nw.SubjectKeyId)
	nw.AuthorityKeyId = canonicalizeNewCertificateAuthorityCaCertificateDescriptionsAuthorityKeyId(c, des.AuthorityKeyId, nw.AuthorityKeyId)
	if dcl.StringArrayCanonicalize(des.CrlDistributionPoints, nw.CrlDistributionPoints) {
		nw.CrlDistributionPoints = des.CrlDistributionPoints
	}
	if dcl.StringArrayCanonicalize(des.AiaIssuingCertificateUrls, nw.AiaIssuingCertificateUrls) {
		nw.AiaIssuingCertificateUrls = des.AiaIssuingCertificateUrls
	}
	nw.CertFingerprint = canonicalizeNewCertificateAuthorityCaCertificateDescriptionsCertFingerprint(c, des.CertFingerprint, nw.CertFingerprint)

	return nw
}

func canonicalizeNewCertificateAuthorityCaCertificateDescriptionsSet(c *Client, des, nw []CertificateAuthorityCaCertificateDescriptions) []CertificateAuthorityCaCertificateDescriptions {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []CertificateAuthorityCaCertificateDescriptions
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareCertificateAuthorityCaCertificateDescriptionsNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewCertificateAuthorityCaCertificateDescriptions(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewCertificateAuthorityCaCertificateDescriptionsSlice(c *Client, des, nw []CertificateAuthorityCaCertificateDescriptions) []CertificateAuthorityCaCertificateDescriptions {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []CertificateAuthorityCaCertificateDescriptions
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewCertificateAuthorityCaCertificateDescriptions(c, &d, &n))
	}

	return items
}

func canonicalizeCertificateAuthorityCaCertificateDescriptionsSubjectDescription(des, initial *CertificateAuthorityCaCertificateDescriptionsSubjectDescription, opts ...dcl.ApplyOption) *CertificateAuthorityCaCertificateDescriptionsSubjectDescription {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &CertificateAuthorityCaCertificateDescriptionsSubjectDescription{}

	cDes.Subject = canonicalizeCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubject(des.Subject, initial.Subject, opts...)
	cDes.SubjectAltName = canonicalizeCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltName(des.SubjectAltName, initial.SubjectAltName, opts...)
	if dcl.StringCanonicalize(des.HexSerialNumber, initial.HexSerialNumber) || dcl.IsZeroValue(des.HexSerialNumber) {
		cDes.HexSerialNumber = initial.HexSerialNumber
	} else {
		cDes.HexSerialNumber = des.HexSerialNumber
	}
	if dcl.StringCanonicalize(des.Lifetime, initial.Lifetime) || dcl.IsZeroValue(des.Lifetime) {
		cDes.Lifetime = initial.Lifetime
	} else {
		cDes.Lifetime = des.Lifetime
	}
	if dcl.IsZeroValue(des.NotBeforeTime) || (dcl.IsEmptyValueIndirect(des.NotBeforeTime) && dcl.IsEmptyValueIndirect(initial.NotBeforeTime)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.NotBeforeTime = initial.NotBeforeTime
	} else {
		cDes.NotBeforeTime = des.NotBeforeTime
	}
	if dcl.IsZeroValue(des.NotAfterTime) || (dcl.IsEmptyValueIndirect(des.NotAfterTime) && dcl.IsEmptyValueIndirect(initial.NotAfterTime)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.NotAfterTime = initial.NotAfterTime
	} else {
		cDes.NotAfterTime = des.NotAfterTime
	}

	return cDes
}

func canonicalizeCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSlice(des, initial []CertificateAuthorityCaCertificateDescriptionsSubjectDescription, opts ...dcl.ApplyOption) []CertificateAuthorityCaCertificateDescriptionsSubjectDescription {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]CertificateAuthorityCaCertificateDescriptionsSubjectDescription, 0, len(des))
		for _, d := range des {
			cd := canonicalizeCertificateAuthorityCaCertificateDescriptionsSubjectDescription(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]CertificateAuthorityCaCertificateDescriptionsSubjectDescription, 0, len(des))
	for i, d := range des {
		cd := canonicalizeCertificateAuthorityCaCertificateDescriptionsSubjectDescription(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewCertificateAuthorityCaCertificateDescriptionsSubjectDescription(c *Client, des, nw *CertificateAuthorityCaCertificateDescriptionsSubjectDescription) *CertificateAuthorityCaCertificateDescriptionsSubjectDescription {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for CertificateAuthorityCaCertificateDescriptionsSubjectDescription while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	nw.Subject = canonicalizeNewCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubject(c, des.Subject, nw.Subject)
	nw.SubjectAltName = canonicalizeNewCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltName(c, des.SubjectAltName, nw.SubjectAltName)
	if dcl.StringCanonicalize(des.HexSerialNumber, nw.HexSerialNumber) {
		nw.HexSerialNumber = des.HexSerialNumber
	}
	if dcl.StringCanonicalize(des.Lifetime, nw.Lifetime) {
		nw.Lifetime = des.Lifetime
	}

	return nw
}

func canonicalizeNewCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSet(c *Client, des, nw []CertificateAuthorityCaCertificateDescriptionsSubjectDescription) []CertificateAuthorityCaCertificateDescriptionsSubjectDescription {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []CertificateAuthorityCaCertificateDescriptionsSubjectDescription
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewCertificateAuthorityCaCertificateDescriptionsSubjectDescription(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSlice(c *Client, des, nw []CertificateAuthorityCaCertificateDescriptionsSubjectDescription) []CertificateAuthorityCaCertificateDescriptionsSubjectDescription {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []CertificateAuthorityCaCertificateDescriptionsSubjectDescription
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewCertificateAuthorityCaCertificateDescriptionsSubjectDescription(c, &d, &n))
	}

	return items
}

func canonicalizeCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubject(des, initial *CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubject, opts ...dcl.ApplyOption) *CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubject {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubject{}

	if dcl.StringCanonicalize(des.CommonName, initial.CommonName) || dcl.IsZeroValue(des.CommonName) {
		cDes.CommonName = initial.CommonName
	} else {
		cDes.CommonName = des.CommonName
	}
	if dcl.StringCanonicalize(des.CountryCode, initial.CountryCode) || dcl.IsZeroValue(des.CountryCode) {
		cDes.CountryCode = initial.CountryCode
	} else {
		cDes.CountryCode = des.CountryCode
	}
	if dcl.StringCanonicalize(des.Organization, initial.Organization) || dcl.IsZeroValue(des.Organization) {
		cDes.Organization = initial.Organization
	} else {
		cDes.Organization = des.Organization
	}
	if dcl.StringCanonicalize(des.OrganizationalUnit, initial.OrganizationalUnit) || dcl.IsZeroValue(des.OrganizationalUnit) {
		cDes.OrganizationalUnit = initial.OrganizationalUnit
	} else {
		cDes.OrganizationalUnit = des.OrganizationalUnit
	}
	if dcl.StringCanonicalize(des.Locality, initial.Locality) || dcl.IsZeroValue(des.Locality) {
		cDes.Locality = initial.Locality
	} else {
		cDes.Locality = des.Locality
	}
	if dcl.StringCanonicalize(des.Province, initial.Province) || dcl.IsZeroValue(des.Province) {
		cDes.Province = initial.Province
	} else {
		cDes.Province = des.Province
	}
	if dcl.StringCanonicalize(des.StreetAddress, initial.StreetAddress) || dcl.IsZeroValue(des.StreetAddress) {
		cDes.StreetAddress = initial.StreetAddress
	} else {
		cDes.StreetAddress = des.StreetAddress
	}
	if dcl.StringCanonicalize(des.PostalCode, initial.PostalCode) || dcl.IsZeroValue(des.PostalCode) {
		cDes.PostalCode = initial.PostalCode
	} else {
		cDes.PostalCode = des.PostalCode
	}

	return cDes
}

func canonicalizeCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectSlice(des, initial []CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubject, opts ...dcl.ApplyOption) []CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubject {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubject, 0, len(des))
		for _, d := range des {
			cd := canonicalizeCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubject(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubject, 0, len(des))
	for i, d := range des {
		cd := canonicalizeCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubject(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubject(c *Client, des, nw *CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubject) *CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubject {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubject while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringCanonicalize(des.CommonName, nw.CommonName) {
		nw.CommonName = des.CommonName
	}
	if dcl.StringCanonicalize(des.CountryCode, nw.CountryCode) {
		nw.CountryCode = des.CountryCode
	}
	if dcl.StringCanonicalize(des.Organization, nw.Organization) {
		nw.Organization = des.Organization
	}
	if dcl.StringCanonicalize(des.OrganizationalUnit, nw.OrganizationalUnit) {
		nw.OrganizationalUnit = des.OrganizationalUnit
	}
	if dcl.StringCanonicalize(des.Locality, nw.Locality) {
		nw.Locality = des.Locality
	}
	if dcl.StringCanonicalize(des.Province, nw.Province) {
		nw.Province = des.Province
	}
	if dcl.StringCanonicalize(des.StreetAddress, nw.StreetAddress) {
		nw.StreetAddress = des.StreetAddress
	}
	if dcl.StringCanonicalize(des.PostalCode, nw.PostalCode) {
		nw.PostalCode = des.PostalCode
	}

	return nw
}

func canonicalizeNewCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectSet(c *Client, des, nw []CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubject) []CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubject {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubject
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubject(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectSlice(c *Client, des, nw []CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubject) []CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubject {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubject
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubject(c, &d, &n))
	}

	return items
}

func canonicalizeCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltName(des, initial *CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltName, opts ...dcl.ApplyOption) *CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltName {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltName{}

	if dcl.StringArrayCanonicalize(des.DnsNames, initial.DnsNames) {
		cDes.DnsNames = initial.DnsNames
	} else {
		cDes.DnsNames = des.DnsNames
	}
	if dcl.StringArrayCanonicalize(des.Uris, initial.Uris) {
		cDes.Uris = initial.Uris
	} else {
		cDes.Uris = des.Uris
	}
	if dcl.StringArrayCanonicalize(des.EmailAddresses, initial.EmailAddresses) {
		cDes.EmailAddresses = initial.EmailAddresses
	} else {
		cDes.EmailAddresses = des.EmailAddresses
	}
	if dcl.StringArrayCanonicalize(des.IPAddresses, initial.IPAddresses) {
		cDes.IPAddresses = initial.IPAddresses
	} else {
		cDes.IPAddresses = des.IPAddresses
	}
	cDes.CustomSans = canonicalizeCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSansSlice(des.CustomSans, initial.CustomSans, opts...)

	return cDes
}

func canonicalizeCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameSlice(des, initial []CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltName, opts ...dcl.ApplyOption) []CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltName {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltName, 0, len(des))
		for _, d := range des {
			cd := canonicalizeCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltName(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltName, 0, len(des))
	for i, d := range des {
		cd := canonicalizeCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltName(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltName(c *Client, des, nw *CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltName) *CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltName {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltName while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringArrayCanonicalize(des.DnsNames, nw.DnsNames) {
		nw.DnsNames = des.DnsNames
	}
	if dcl.StringArrayCanonicalize(des.Uris, nw.Uris) {
		nw.Uris = des.Uris
	}
	if dcl.StringArrayCanonicalize(des.EmailAddresses, nw.EmailAddresses) {
		nw.EmailAddresses = des.EmailAddresses
	}
	if dcl.StringArrayCanonicalize(des.IPAddresses, nw.IPAddresses) {
		nw.IPAddresses = des.IPAddresses
	}
	nw.CustomSans = canonicalizeNewCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSansSlice(c, des.CustomSans, nw.CustomSans)

	return nw
}

func canonicalizeNewCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameSet(c *Client, des, nw []CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltName) []CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltName {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltName
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltName(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameSlice(c *Client, des, nw []CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltName) []CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltName {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltName
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltName(c, &d, &n))
	}

	return items
}

func canonicalizeCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSans(des, initial *CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSans, opts ...dcl.ApplyOption) *CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSans {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSans{}

	cDes.ObjectId = canonicalizeCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSansObjectId(des.ObjectId, initial.ObjectId, opts...)
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

func canonicalizeCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSansSlice(des, initial []CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSans, opts ...dcl.ApplyOption) []CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSans {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSans, 0, len(des))
		for _, d := range des {
			cd := canonicalizeCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSans(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSans, 0, len(des))
	for i, d := range des {
		cd := canonicalizeCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSans(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSans(c *Client, des, nw *CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSans) *CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSans {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSans while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	nw.ObjectId = canonicalizeNewCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSansObjectId(c, des.ObjectId, nw.ObjectId)
	if dcl.BoolCanonicalize(des.Critical, nw.Critical) {
		nw.Critical = des.Critical
	}
	if dcl.StringCanonicalize(des.Value, nw.Value) {
		nw.Value = des.Value
	}

	return nw
}

func canonicalizeNewCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSansSet(c *Client, des, nw []CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSans) []CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSans {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSans
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSansNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSans(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSansSlice(c *Client, des, nw []CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSans) []CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSans {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSans
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSans(c, &d, &n))
	}

	return items
}

func canonicalizeCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSansObjectId(des, initial *CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSansObjectId, opts ...dcl.ApplyOption) *CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSansObjectId {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSansObjectId{}

	if dcl.IsZeroValue(des.ObjectIdPath) || (dcl.IsEmptyValueIndirect(des.ObjectIdPath) && dcl.IsEmptyValueIndirect(initial.ObjectIdPath)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.ObjectIdPath = initial.ObjectIdPath
	} else {
		cDes.ObjectIdPath = des.ObjectIdPath
	}

	return cDes
}

func canonicalizeCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSansObjectIdSlice(des, initial []CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSansObjectId, opts ...dcl.ApplyOption) []CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSansObjectId {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSansObjectId, 0, len(des))
		for _, d := range des {
			cd := canonicalizeCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSansObjectId(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSansObjectId, 0, len(des))
	for i, d := range des {
		cd := canonicalizeCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSansObjectId(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSansObjectId(c *Client, des, nw *CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSansObjectId) *CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSansObjectId {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSansObjectId while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	return nw
}

func canonicalizeNewCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSansObjectIdSet(c *Client, des, nw []CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSansObjectId) []CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSansObjectId {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSansObjectId
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSansObjectIdNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSansObjectId(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSansObjectIdSlice(c *Client, des, nw []CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSansObjectId) []CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSansObjectId {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSansObjectId
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSansObjectId(c, &d, &n))
	}

	return items
}

func canonicalizeCertificateAuthorityCaCertificateDescriptionsX509Description(des, initial *CertificateAuthorityCaCertificateDescriptionsX509Description, opts ...dcl.ApplyOption) *CertificateAuthorityCaCertificateDescriptionsX509Description {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &CertificateAuthorityCaCertificateDescriptionsX509Description{}

	cDes.KeyUsage = canonicalizeCertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsage(des.KeyUsage, initial.KeyUsage, opts...)
	cDes.CaOptions = canonicalizeCertificateAuthorityCaCertificateDescriptionsX509DescriptionCaOptions(des.CaOptions, initial.CaOptions, opts...)
	cDes.PolicyIds = canonicalizeCertificateAuthorityCaCertificateDescriptionsX509DescriptionPolicyIdsSlice(des.PolicyIds, initial.PolicyIds, opts...)
	cDes.AdditionalExtensions = canonicalizeCertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensionsSlice(des.AdditionalExtensions, initial.AdditionalExtensions, opts...)

	return cDes
}

func canonicalizeCertificateAuthorityCaCertificateDescriptionsX509DescriptionSlice(des, initial []CertificateAuthorityCaCertificateDescriptionsX509Description, opts ...dcl.ApplyOption) []CertificateAuthorityCaCertificateDescriptionsX509Description {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]CertificateAuthorityCaCertificateDescriptionsX509Description, 0, len(des))
		for _, d := range des {
			cd := canonicalizeCertificateAuthorityCaCertificateDescriptionsX509Description(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]CertificateAuthorityCaCertificateDescriptionsX509Description, 0, len(des))
	for i, d := range des {
		cd := canonicalizeCertificateAuthorityCaCertificateDescriptionsX509Description(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewCertificateAuthorityCaCertificateDescriptionsX509Description(c *Client, des, nw *CertificateAuthorityCaCertificateDescriptionsX509Description) *CertificateAuthorityCaCertificateDescriptionsX509Description {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for CertificateAuthorityCaCertificateDescriptionsX509Description while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	nw.KeyUsage = canonicalizeNewCertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsage(c, des.KeyUsage, nw.KeyUsage)
	nw.CaOptions = canonicalizeNewCertificateAuthorityCaCertificateDescriptionsX509DescriptionCaOptions(c, des.CaOptions, nw.CaOptions)
	nw.PolicyIds = canonicalizeNewCertificateAuthorityCaCertificateDescriptionsX509DescriptionPolicyIdsSlice(c, des.PolicyIds, nw.PolicyIds)
	if dcl.StringArrayCanonicalize(des.AiaOcspServers, nw.AiaOcspServers) {
		nw.AiaOcspServers = des.AiaOcspServers
	}
	nw.AdditionalExtensions = canonicalizeNewCertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensionsSlice(c, des.AdditionalExtensions, nw.AdditionalExtensions)

	return nw
}

func canonicalizeNewCertificateAuthorityCaCertificateDescriptionsX509DescriptionSet(c *Client, des, nw []CertificateAuthorityCaCertificateDescriptionsX509Description) []CertificateAuthorityCaCertificateDescriptionsX509Description {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []CertificateAuthorityCaCertificateDescriptionsX509Description
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareCertificateAuthorityCaCertificateDescriptionsX509DescriptionNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewCertificateAuthorityCaCertificateDescriptionsX509Description(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewCertificateAuthorityCaCertificateDescriptionsX509DescriptionSlice(c *Client, des, nw []CertificateAuthorityCaCertificateDescriptionsX509Description) []CertificateAuthorityCaCertificateDescriptionsX509Description {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []CertificateAuthorityCaCertificateDescriptionsX509Description
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewCertificateAuthorityCaCertificateDescriptionsX509Description(c, &d, &n))
	}

	return items
}

func canonicalizeCertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsage(des, initial *CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsage, opts ...dcl.ApplyOption) *CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsage {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsage{}

	cDes.BaseKeyUsage = canonicalizeCertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageBaseKeyUsage(des.BaseKeyUsage, initial.BaseKeyUsage, opts...)
	cDes.ExtendedKeyUsage = canonicalizeCertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageExtendedKeyUsage(des.ExtendedKeyUsage, initial.ExtendedKeyUsage, opts...)
	cDes.UnknownExtendedKeyUsages = canonicalizeCertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageUnknownExtendedKeyUsagesSlice(des.UnknownExtendedKeyUsages, initial.UnknownExtendedKeyUsages, opts...)

	return cDes
}

func canonicalizeCertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageSlice(des, initial []CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsage, opts ...dcl.ApplyOption) []CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsage {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsage, 0, len(des))
		for _, d := range des {
			cd := canonicalizeCertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsage(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsage, 0, len(des))
	for i, d := range des {
		cd := canonicalizeCertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsage(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewCertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsage(c *Client, des, nw *CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsage) *CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsage {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsage while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	nw.BaseKeyUsage = canonicalizeNewCertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageBaseKeyUsage(c, des.BaseKeyUsage, nw.BaseKeyUsage)
	nw.ExtendedKeyUsage = canonicalizeNewCertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageExtendedKeyUsage(c, des.ExtendedKeyUsage, nw.ExtendedKeyUsage)
	nw.UnknownExtendedKeyUsages = canonicalizeNewCertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageUnknownExtendedKeyUsagesSlice(c, des.UnknownExtendedKeyUsages, nw.UnknownExtendedKeyUsages)

	return nw
}

func canonicalizeNewCertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageSet(c *Client, des, nw []CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsage) []CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsage {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsage
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareCertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewCertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsage(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewCertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageSlice(c *Client, des, nw []CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsage) []CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsage {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsage
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewCertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsage(c, &d, &n))
	}

	return items
}

func canonicalizeCertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageBaseKeyUsage(des, initial *CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageBaseKeyUsage, opts ...dcl.ApplyOption) *CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageBaseKeyUsage {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageBaseKeyUsage{}

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

func canonicalizeCertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageBaseKeyUsageSlice(des, initial []CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageBaseKeyUsage, opts ...dcl.ApplyOption) []CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageBaseKeyUsage {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageBaseKeyUsage, 0, len(des))
		for _, d := range des {
			cd := canonicalizeCertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageBaseKeyUsage(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageBaseKeyUsage, 0, len(des))
	for i, d := range des {
		cd := canonicalizeCertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageBaseKeyUsage(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewCertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageBaseKeyUsage(c *Client, des, nw *CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageBaseKeyUsage) *CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageBaseKeyUsage {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageBaseKeyUsage while comparing non-nil desired to nil actual.  Returning desired object.")
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

func canonicalizeNewCertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageBaseKeyUsageSet(c *Client, des, nw []CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageBaseKeyUsage) []CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageBaseKeyUsage {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageBaseKeyUsage
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareCertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageBaseKeyUsageNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewCertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageBaseKeyUsage(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewCertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageBaseKeyUsageSlice(c *Client, des, nw []CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageBaseKeyUsage) []CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageBaseKeyUsage {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageBaseKeyUsage
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewCertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageBaseKeyUsage(c, &d, &n))
	}

	return items
}

func canonicalizeCertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageExtendedKeyUsage(des, initial *CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageExtendedKeyUsage, opts ...dcl.ApplyOption) *CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageExtendedKeyUsage {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageExtendedKeyUsage{}

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

func canonicalizeCertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageExtendedKeyUsageSlice(des, initial []CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageExtendedKeyUsage, opts ...dcl.ApplyOption) []CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageExtendedKeyUsage {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageExtendedKeyUsage, 0, len(des))
		for _, d := range des {
			cd := canonicalizeCertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageExtendedKeyUsage(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageExtendedKeyUsage, 0, len(des))
	for i, d := range des {
		cd := canonicalizeCertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageExtendedKeyUsage(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewCertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageExtendedKeyUsage(c *Client, des, nw *CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageExtendedKeyUsage) *CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageExtendedKeyUsage {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageExtendedKeyUsage while comparing non-nil desired to nil actual.  Returning desired object.")
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

func canonicalizeNewCertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageExtendedKeyUsageSet(c *Client, des, nw []CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageExtendedKeyUsage) []CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageExtendedKeyUsage {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageExtendedKeyUsage
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareCertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageExtendedKeyUsageNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewCertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageExtendedKeyUsage(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewCertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageExtendedKeyUsageSlice(c *Client, des, nw []CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageExtendedKeyUsage) []CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageExtendedKeyUsage {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageExtendedKeyUsage
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewCertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageExtendedKeyUsage(c, &d, &n))
	}

	return items
}

func canonicalizeCertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageUnknownExtendedKeyUsages(des, initial *CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageUnknownExtendedKeyUsages, opts ...dcl.ApplyOption) *CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageUnknownExtendedKeyUsages {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageUnknownExtendedKeyUsages{}

	if dcl.IsZeroValue(des.ObjectIdPath) || (dcl.IsEmptyValueIndirect(des.ObjectIdPath) && dcl.IsEmptyValueIndirect(initial.ObjectIdPath)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.ObjectIdPath = initial.ObjectIdPath
	} else {
		cDes.ObjectIdPath = des.ObjectIdPath
	}

	return cDes
}

func canonicalizeCertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageUnknownExtendedKeyUsagesSlice(des, initial []CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageUnknownExtendedKeyUsages, opts ...dcl.ApplyOption) []CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageUnknownExtendedKeyUsages {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageUnknownExtendedKeyUsages, 0, len(des))
		for _, d := range des {
			cd := canonicalizeCertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageUnknownExtendedKeyUsages(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageUnknownExtendedKeyUsages, 0, len(des))
	for i, d := range des {
		cd := canonicalizeCertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageUnknownExtendedKeyUsages(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewCertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageUnknownExtendedKeyUsages(c *Client, des, nw *CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageUnknownExtendedKeyUsages) *CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageUnknownExtendedKeyUsages {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageUnknownExtendedKeyUsages while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	return nw
}

func canonicalizeNewCertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageUnknownExtendedKeyUsagesSet(c *Client, des, nw []CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageUnknownExtendedKeyUsages) []CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageUnknownExtendedKeyUsages {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageUnknownExtendedKeyUsages
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareCertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageUnknownExtendedKeyUsagesNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewCertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageUnknownExtendedKeyUsages(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewCertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageUnknownExtendedKeyUsagesSlice(c *Client, des, nw []CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageUnknownExtendedKeyUsages) []CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageUnknownExtendedKeyUsages {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageUnknownExtendedKeyUsages
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewCertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageUnknownExtendedKeyUsages(c, &d, &n))
	}

	return items
}

func canonicalizeCertificateAuthorityCaCertificateDescriptionsX509DescriptionCaOptions(des, initial *CertificateAuthorityCaCertificateDescriptionsX509DescriptionCaOptions, opts ...dcl.ApplyOption) *CertificateAuthorityCaCertificateDescriptionsX509DescriptionCaOptions {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &CertificateAuthorityCaCertificateDescriptionsX509DescriptionCaOptions{}

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

	return cDes
}

func canonicalizeCertificateAuthorityCaCertificateDescriptionsX509DescriptionCaOptionsSlice(des, initial []CertificateAuthorityCaCertificateDescriptionsX509DescriptionCaOptions, opts ...dcl.ApplyOption) []CertificateAuthorityCaCertificateDescriptionsX509DescriptionCaOptions {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]CertificateAuthorityCaCertificateDescriptionsX509DescriptionCaOptions, 0, len(des))
		for _, d := range des {
			cd := canonicalizeCertificateAuthorityCaCertificateDescriptionsX509DescriptionCaOptions(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]CertificateAuthorityCaCertificateDescriptionsX509DescriptionCaOptions, 0, len(des))
	for i, d := range des {
		cd := canonicalizeCertificateAuthorityCaCertificateDescriptionsX509DescriptionCaOptions(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewCertificateAuthorityCaCertificateDescriptionsX509DescriptionCaOptions(c *Client, des, nw *CertificateAuthorityCaCertificateDescriptionsX509DescriptionCaOptions) *CertificateAuthorityCaCertificateDescriptionsX509DescriptionCaOptions {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for CertificateAuthorityCaCertificateDescriptionsX509DescriptionCaOptions while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.BoolCanonicalize(des.IsCa, nw.IsCa) {
		nw.IsCa = des.IsCa
	}

	return nw
}

func canonicalizeNewCertificateAuthorityCaCertificateDescriptionsX509DescriptionCaOptionsSet(c *Client, des, nw []CertificateAuthorityCaCertificateDescriptionsX509DescriptionCaOptions) []CertificateAuthorityCaCertificateDescriptionsX509DescriptionCaOptions {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []CertificateAuthorityCaCertificateDescriptionsX509DescriptionCaOptions
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareCertificateAuthorityCaCertificateDescriptionsX509DescriptionCaOptionsNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewCertificateAuthorityCaCertificateDescriptionsX509DescriptionCaOptions(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewCertificateAuthorityCaCertificateDescriptionsX509DescriptionCaOptionsSlice(c *Client, des, nw []CertificateAuthorityCaCertificateDescriptionsX509DescriptionCaOptions) []CertificateAuthorityCaCertificateDescriptionsX509DescriptionCaOptions {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []CertificateAuthorityCaCertificateDescriptionsX509DescriptionCaOptions
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewCertificateAuthorityCaCertificateDescriptionsX509DescriptionCaOptions(c, &d, &n))
	}

	return items
}

func canonicalizeCertificateAuthorityCaCertificateDescriptionsX509DescriptionPolicyIds(des, initial *CertificateAuthorityCaCertificateDescriptionsX509DescriptionPolicyIds, opts ...dcl.ApplyOption) *CertificateAuthorityCaCertificateDescriptionsX509DescriptionPolicyIds {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &CertificateAuthorityCaCertificateDescriptionsX509DescriptionPolicyIds{}

	if dcl.IsZeroValue(des.ObjectIdPath) || (dcl.IsEmptyValueIndirect(des.ObjectIdPath) && dcl.IsEmptyValueIndirect(initial.ObjectIdPath)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.ObjectIdPath = initial.ObjectIdPath
	} else {
		cDes.ObjectIdPath = des.ObjectIdPath
	}

	return cDes
}

func canonicalizeCertificateAuthorityCaCertificateDescriptionsX509DescriptionPolicyIdsSlice(des, initial []CertificateAuthorityCaCertificateDescriptionsX509DescriptionPolicyIds, opts ...dcl.ApplyOption) []CertificateAuthorityCaCertificateDescriptionsX509DescriptionPolicyIds {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]CertificateAuthorityCaCertificateDescriptionsX509DescriptionPolicyIds, 0, len(des))
		for _, d := range des {
			cd := canonicalizeCertificateAuthorityCaCertificateDescriptionsX509DescriptionPolicyIds(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]CertificateAuthorityCaCertificateDescriptionsX509DescriptionPolicyIds, 0, len(des))
	for i, d := range des {
		cd := canonicalizeCertificateAuthorityCaCertificateDescriptionsX509DescriptionPolicyIds(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewCertificateAuthorityCaCertificateDescriptionsX509DescriptionPolicyIds(c *Client, des, nw *CertificateAuthorityCaCertificateDescriptionsX509DescriptionPolicyIds) *CertificateAuthorityCaCertificateDescriptionsX509DescriptionPolicyIds {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for CertificateAuthorityCaCertificateDescriptionsX509DescriptionPolicyIds while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	return nw
}

func canonicalizeNewCertificateAuthorityCaCertificateDescriptionsX509DescriptionPolicyIdsSet(c *Client, des, nw []CertificateAuthorityCaCertificateDescriptionsX509DescriptionPolicyIds) []CertificateAuthorityCaCertificateDescriptionsX509DescriptionPolicyIds {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []CertificateAuthorityCaCertificateDescriptionsX509DescriptionPolicyIds
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareCertificateAuthorityCaCertificateDescriptionsX509DescriptionPolicyIdsNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewCertificateAuthorityCaCertificateDescriptionsX509DescriptionPolicyIds(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewCertificateAuthorityCaCertificateDescriptionsX509DescriptionPolicyIdsSlice(c *Client, des, nw []CertificateAuthorityCaCertificateDescriptionsX509DescriptionPolicyIds) []CertificateAuthorityCaCertificateDescriptionsX509DescriptionPolicyIds {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []CertificateAuthorityCaCertificateDescriptionsX509DescriptionPolicyIds
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewCertificateAuthorityCaCertificateDescriptionsX509DescriptionPolicyIds(c, &d, &n))
	}

	return items
}

func canonicalizeCertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensions(des, initial *CertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensions, opts ...dcl.ApplyOption) *CertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensions {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &CertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensions{}

	cDes.ObjectId = canonicalizeCertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensionsObjectId(des.ObjectId, initial.ObjectId, opts...)
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

func canonicalizeCertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensionsSlice(des, initial []CertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensions, opts ...dcl.ApplyOption) []CertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensions {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]CertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensions, 0, len(des))
		for _, d := range des {
			cd := canonicalizeCertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensions(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]CertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensions, 0, len(des))
	for i, d := range des {
		cd := canonicalizeCertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensions(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewCertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensions(c *Client, des, nw *CertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensions) *CertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensions {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for CertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensions while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	nw.ObjectId = canonicalizeNewCertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensionsObjectId(c, des.ObjectId, nw.ObjectId)
	if dcl.BoolCanonicalize(des.Critical, nw.Critical) {
		nw.Critical = des.Critical
	}
	if dcl.StringCanonicalize(des.Value, nw.Value) {
		nw.Value = des.Value
	}

	return nw
}

func canonicalizeNewCertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensionsSet(c *Client, des, nw []CertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensions) []CertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensions {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []CertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensions
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareCertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensionsNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewCertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensions(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewCertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensionsSlice(c *Client, des, nw []CertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensions) []CertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensions {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []CertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensions
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewCertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensions(c, &d, &n))
	}

	return items
}

func canonicalizeCertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensionsObjectId(des, initial *CertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensionsObjectId, opts ...dcl.ApplyOption) *CertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensionsObjectId {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &CertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensionsObjectId{}

	if dcl.IsZeroValue(des.ObjectIdPath) || (dcl.IsEmptyValueIndirect(des.ObjectIdPath) && dcl.IsEmptyValueIndirect(initial.ObjectIdPath)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.ObjectIdPath = initial.ObjectIdPath
	} else {
		cDes.ObjectIdPath = des.ObjectIdPath
	}

	return cDes
}

func canonicalizeCertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensionsObjectIdSlice(des, initial []CertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensionsObjectId, opts ...dcl.ApplyOption) []CertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensionsObjectId {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]CertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensionsObjectId, 0, len(des))
		for _, d := range des {
			cd := canonicalizeCertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensionsObjectId(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]CertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensionsObjectId, 0, len(des))
	for i, d := range des {
		cd := canonicalizeCertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensionsObjectId(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewCertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensionsObjectId(c *Client, des, nw *CertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensionsObjectId) *CertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensionsObjectId {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for CertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensionsObjectId while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	return nw
}

func canonicalizeNewCertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensionsObjectIdSet(c *Client, des, nw []CertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensionsObjectId) []CertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensionsObjectId {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []CertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensionsObjectId
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareCertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensionsObjectIdNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewCertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensionsObjectId(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewCertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensionsObjectIdSlice(c *Client, des, nw []CertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensionsObjectId) []CertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensionsObjectId {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []CertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensionsObjectId
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewCertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensionsObjectId(c, &d, &n))
	}

	return items
}

func canonicalizeCertificateAuthorityCaCertificateDescriptionsPublicKey(des, initial *CertificateAuthorityCaCertificateDescriptionsPublicKey, opts ...dcl.ApplyOption) *CertificateAuthorityCaCertificateDescriptionsPublicKey {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &CertificateAuthorityCaCertificateDescriptionsPublicKey{}

	if dcl.StringCanonicalize(des.Key, initial.Key) || dcl.IsZeroValue(des.Key) {
		cDes.Key = initial.Key
	} else {
		cDes.Key = des.Key
	}
	if dcl.IsZeroValue(des.Format) || (dcl.IsEmptyValueIndirect(des.Format) && dcl.IsEmptyValueIndirect(initial.Format)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.Format = initial.Format
	} else {
		cDes.Format = des.Format
	}

	return cDes
}

func canonicalizeCertificateAuthorityCaCertificateDescriptionsPublicKeySlice(des, initial []CertificateAuthorityCaCertificateDescriptionsPublicKey, opts ...dcl.ApplyOption) []CertificateAuthorityCaCertificateDescriptionsPublicKey {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]CertificateAuthorityCaCertificateDescriptionsPublicKey, 0, len(des))
		for _, d := range des {
			cd := canonicalizeCertificateAuthorityCaCertificateDescriptionsPublicKey(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]CertificateAuthorityCaCertificateDescriptionsPublicKey, 0, len(des))
	for i, d := range des {
		cd := canonicalizeCertificateAuthorityCaCertificateDescriptionsPublicKey(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewCertificateAuthorityCaCertificateDescriptionsPublicKey(c *Client, des, nw *CertificateAuthorityCaCertificateDescriptionsPublicKey) *CertificateAuthorityCaCertificateDescriptionsPublicKey {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for CertificateAuthorityCaCertificateDescriptionsPublicKey while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringCanonicalize(des.Key, nw.Key) {
		nw.Key = des.Key
	}

	return nw
}

func canonicalizeNewCertificateAuthorityCaCertificateDescriptionsPublicKeySet(c *Client, des, nw []CertificateAuthorityCaCertificateDescriptionsPublicKey) []CertificateAuthorityCaCertificateDescriptionsPublicKey {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []CertificateAuthorityCaCertificateDescriptionsPublicKey
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareCertificateAuthorityCaCertificateDescriptionsPublicKeyNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewCertificateAuthorityCaCertificateDescriptionsPublicKey(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewCertificateAuthorityCaCertificateDescriptionsPublicKeySlice(c *Client, des, nw []CertificateAuthorityCaCertificateDescriptionsPublicKey) []CertificateAuthorityCaCertificateDescriptionsPublicKey {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []CertificateAuthorityCaCertificateDescriptionsPublicKey
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewCertificateAuthorityCaCertificateDescriptionsPublicKey(c, &d, &n))
	}

	return items
}

func canonicalizeCertificateAuthorityCaCertificateDescriptionsSubjectKeyId(des, initial *CertificateAuthorityCaCertificateDescriptionsSubjectKeyId, opts ...dcl.ApplyOption) *CertificateAuthorityCaCertificateDescriptionsSubjectKeyId {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &CertificateAuthorityCaCertificateDescriptionsSubjectKeyId{}

	if dcl.StringCanonicalize(des.KeyId, initial.KeyId) || dcl.IsZeroValue(des.KeyId) {
		cDes.KeyId = initial.KeyId
	} else {
		cDes.KeyId = des.KeyId
	}

	return cDes
}

func canonicalizeCertificateAuthorityCaCertificateDescriptionsSubjectKeyIdSlice(des, initial []CertificateAuthorityCaCertificateDescriptionsSubjectKeyId, opts ...dcl.ApplyOption) []CertificateAuthorityCaCertificateDescriptionsSubjectKeyId {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]CertificateAuthorityCaCertificateDescriptionsSubjectKeyId, 0, len(des))
		for _, d := range des {
			cd := canonicalizeCertificateAuthorityCaCertificateDescriptionsSubjectKeyId(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]CertificateAuthorityCaCertificateDescriptionsSubjectKeyId, 0, len(des))
	for i, d := range des {
		cd := canonicalizeCertificateAuthorityCaCertificateDescriptionsSubjectKeyId(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewCertificateAuthorityCaCertificateDescriptionsSubjectKeyId(c *Client, des, nw *CertificateAuthorityCaCertificateDescriptionsSubjectKeyId) *CertificateAuthorityCaCertificateDescriptionsSubjectKeyId {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for CertificateAuthorityCaCertificateDescriptionsSubjectKeyId while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringCanonicalize(des.KeyId, nw.KeyId) {
		nw.KeyId = des.KeyId
	}

	return nw
}

func canonicalizeNewCertificateAuthorityCaCertificateDescriptionsSubjectKeyIdSet(c *Client, des, nw []CertificateAuthorityCaCertificateDescriptionsSubjectKeyId) []CertificateAuthorityCaCertificateDescriptionsSubjectKeyId {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []CertificateAuthorityCaCertificateDescriptionsSubjectKeyId
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareCertificateAuthorityCaCertificateDescriptionsSubjectKeyIdNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewCertificateAuthorityCaCertificateDescriptionsSubjectKeyId(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewCertificateAuthorityCaCertificateDescriptionsSubjectKeyIdSlice(c *Client, des, nw []CertificateAuthorityCaCertificateDescriptionsSubjectKeyId) []CertificateAuthorityCaCertificateDescriptionsSubjectKeyId {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []CertificateAuthorityCaCertificateDescriptionsSubjectKeyId
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewCertificateAuthorityCaCertificateDescriptionsSubjectKeyId(c, &d, &n))
	}

	return items
}

func canonicalizeCertificateAuthorityCaCertificateDescriptionsAuthorityKeyId(des, initial *CertificateAuthorityCaCertificateDescriptionsAuthorityKeyId, opts ...dcl.ApplyOption) *CertificateAuthorityCaCertificateDescriptionsAuthorityKeyId {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &CertificateAuthorityCaCertificateDescriptionsAuthorityKeyId{}

	if dcl.StringCanonicalize(des.KeyId, initial.KeyId) || dcl.IsZeroValue(des.KeyId) {
		cDes.KeyId = initial.KeyId
	} else {
		cDes.KeyId = des.KeyId
	}

	return cDes
}

func canonicalizeCertificateAuthorityCaCertificateDescriptionsAuthorityKeyIdSlice(des, initial []CertificateAuthorityCaCertificateDescriptionsAuthorityKeyId, opts ...dcl.ApplyOption) []CertificateAuthorityCaCertificateDescriptionsAuthorityKeyId {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]CertificateAuthorityCaCertificateDescriptionsAuthorityKeyId, 0, len(des))
		for _, d := range des {
			cd := canonicalizeCertificateAuthorityCaCertificateDescriptionsAuthorityKeyId(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]CertificateAuthorityCaCertificateDescriptionsAuthorityKeyId, 0, len(des))
	for i, d := range des {
		cd := canonicalizeCertificateAuthorityCaCertificateDescriptionsAuthorityKeyId(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewCertificateAuthorityCaCertificateDescriptionsAuthorityKeyId(c *Client, des, nw *CertificateAuthorityCaCertificateDescriptionsAuthorityKeyId) *CertificateAuthorityCaCertificateDescriptionsAuthorityKeyId {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for CertificateAuthorityCaCertificateDescriptionsAuthorityKeyId while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringCanonicalize(des.KeyId, nw.KeyId) {
		nw.KeyId = des.KeyId
	}

	return nw
}

func canonicalizeNewCertificateAuthorityCaCertificateDescriptionsAuthorityKeyIdSet(c *Client, des, nw []CertificateAuthorityCaCertificateDescriptionsAuthorityKeyId) []CertificateAuthorityCaCertificateDescriptionsAuthorityKeyId {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []CertificateAuthorityCaCertificateDescriptionsAuthorityKeyId
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareCertificateAuthorityCaCertificateDescriptionsAuthorityKeyIdNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewCertificateAuthorityCaCertificateDescriptionsAuthorityKeyId(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewCertificateAuthorityCaCertificateDescriptionsAuthorityKeyIdSlice(c *Client, des, nw []CertificateAuthorityCaCertificateDescriptionsAuthorityKeyId) []CertificateAuthorityCaCertificateDescriptionsAuthorityKeyId {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []CertificateAuthorityCaCertificateDescriptionsAuthorityKeyId
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewCertificateAuthorityCaCertificateDescriptionsAuthorityKeyId(c, &d, &n))
	}

	return items
}

func canonicalizeCertificateAuthorityCaCertificateDescriptionsCertFingerprint(des, initial *CertificateAuthorityCaCertificateDescriptionsCertFingerprint, opts ...dcl.ApplyOption) *CertificateAuthorityCaCertificateDescriptionsCertFingerprint {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &CertificateAuthorityCaCertificateDescriptionsCertFingerprint{}

	if dcl.StringCanonicalize(des.Sha256Hash, initial.Sha256Hash) || dcl.IsZeroValue(des.Sha256Hash) {
		cDes.Sha256Hash = initial.Sha256Hash
	} else {
		cDes.Sha256Hash = des.Sha256Hash
	}

	return cDes
}

func canonicalizeCertificateAuthorityCaCertificateDescriptionsCertFingerprintSlice(des, initial []CertificateAuthorityCaCertificateDescriptionsCertFingerprint, opts ...dcl.ApplyOption) []CertificateAuthorityCaCertificateDescriptionsCertFingerprint {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]CertificateAuthorityCaCertificateDescriptionsCertFingerprint, 0, len(des))
		for _, d := range des {
			cd := canonicalizeCertificateAuthorityCaCertificateDescriptionsCertFingerprint(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]CertificateAuthorityCaCertificateDescriptionsCertFingerprint, 0, len(des))
	for i, d := range des {
		cd := canonicalizeCertificateAuthorityCaCertificateDescriptionsCertFingerprint(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewCertificateAuthorityCaCertificateDescriptionsCertFingerprint(c *Client, des, nw *CertificateAuthorityCaCertificateDescriptionsCertFingerprint) *CertificateAuthorityCaCertificateDescriptionsCertFingerprint {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for CertificateAuthorityCaCertificateDescriptionsCertFingerprint while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringCanonicalize(des.Sha256Hash, nw.Sha256Hash) {
		nw.Sha256Hash = des.Sha256Hash
	}

	return nw
}

func canonicalizeNewCertificateAuthorityCaCertificateDescriptionsCertFingerprintSet(c *Client, des, nw []CertificateAuthorityCaCertificateDescriptionsCertFingerprint) []CertificateAuthorityCaCertificateDescriptionsCertFingerprint {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []CertificateAuthorityCaCertificateDescriptionsCertFingerprint
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareCertificateAuthorityCaCertificateDescriptionsCertFingerprintNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewCertificateAuthorityCaCertificateDescriptionsCertFingerprint(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewCertificateAuthorityCaCertificateDescriptionsCertFingerprintSlice(c *Client, des, nw []CertificateAuthorityCaCertificateDescriptionsCertFingerprint) []CertificateAuthorityCaCertificateDescriptionsCertFingerprint {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []CertificateAuthorityCaCertificateDescriptionsCertFingerprint
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewCertificateAuthorityCaCertificateDescriptionsCertFingerprint(c, &d, &n))
	}

	return items
}

func canonicalizeCertificateAuthorityAccessUrls(des, initial *CertificateAuthorityAccessUrls, opts ...dcl.ApplyOption) *CertificateAuthorityAccessUrls {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &CertificateAuthorityAccessUrls{}

	if dcl.StringCanonicalize(des.CaCertificateAccessUrl, initial.CaCertificateAccessUrl) || dcl.IsZeroValue(des.CaCertificateAccessUrl) {
		cDes.CaCertificateAccessUrl = initial.CaCertificateAccessUrl
	} else {
		cDes.CaCertificateAccessUrl = des.CaCertificateAccessUrl
	}
	if dcl.StringArrayCanonicalize(des.CrlAccessUrls, initial.CrlAccessUrls) {
		cDes.CrlAccessUrls = initial.CrlAccessUrls
	} else {
		cDes.CrlAccessUrls = des.CrlAccessUrls
	}

	return cDes
}

func canonicalizeCertificateAuthorityAccessUrlsSlice(des, initial []CertificateAuthorityAccessUrls, opts ...dcl.ApplyOption) []CertificateAuthorityAccessUrls {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]CertificateAuthorityAccessUrls, 0, len(des))
		for _, d := range des {
			cd := canonicalizeCertificateAuthorityAccessUrls(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]CertificateAuthorityAccessUrls, 0, len(des))
	for i, d := range des {
		cd := canonicalizeCertificateAuthorityAccessUrls(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewCertificateAuthorityAccessUrls(c *Client, des, nw *CertificateAuthorityAccessUrls) *CertificateAuthorityAccessUrls {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for CertificateAuthorityAccessUrls while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringCanonicalize(des.CaCertificateAccessUrl, nw.CaCertificateAccessUrl) {
		nw.CaCertificateAccessUrl = des.CaCertificateAccessUrl
	}
	if dcl.StringArrayCanonicalize(des.CrlAccessUrls, nw.CrlAccessUrls) {
		nw.CrlAccessUrls = des.CrlAccessUrls
	}

	return nw
}

func canonicalizeNewCertificateAuthorityAccessUrlsSet(c *Client, des, nw []CertificateAuthorityAccessUrls) []CertificateAuthorityAccessUrls {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []CertificateAuthorityAccessUrls
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareCertificateAuthorityAccessUrlsNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewCertificateAuthorityAccessUrls(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewCertificateAuthorityAccessUrlsSlice(c *Client, des, nw []CertificateAuthorityAccessUrls) []CertificateAuthorityAccessUrls {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []CertificateAuthorityAccessUrls
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewCertificateAuthorityAccessUrls(c, &d, &n))
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
func diffCertificateAuthority(c *Client, desired, actual *CertificateAuthority, opts ...dcl.ApplyOption) ([]*dcl.FieldDiff, error) {
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

	if ds, err := dcl.Diff(desired.Type, actual.Type, dcl.DiffInfo{Type: "EnumType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Type")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Config, actual.Config, dcl.DiffInfo{ObjectFunction: compareCertificateAuthorityConfigNewStyle, EmptyObject: EmptyCertificateAuthorityConfig, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Config")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Lifetime, actual.Lifetime, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Lifetime")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.KeySpec, actual.KeySpec, dcl.DiffInfo{ObjectFunction: compareCertificateAuthorityKeySpecNewStyle, EmptyObject: EmptyCertificateAuthorityKeySpec, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("KeySpec")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.SubordinateConfig, actual.SubordinateConfig, dcl.DiffInfo{OutputOnly: true, ObjectFunction: compareCertificateAuthoritySubordinateConfigNewStyle, EmptyObject: EmptyCertificateAuthoritySubordinateConfig, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("SubordinateConfig")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Tier, actual.Tier, dcl.DiffInfo{OutputOnly: true, Type: "EnumType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Tier")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.State, actual.State, dcl.DiffInfo{OutputOnly: true, Type: "EnumType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("State")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.PemCaCertificates, actual.PemCaCertificates, dcl.DiffInfo{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("PemCaCertificates")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.CaCertificateDescriptions, actual.CaCertificateDescriptions, dcl.DiffInfo{OutputOnly: true, ObjectFunction: compareCertificateAuthorityCaCertificateDescriptionsNewStyle, EmptyObject: EmptyCertificateAuthorityCaCertificateDescriptions, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("CaCertificateDescriptions")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.GcsBucket, actual.GcsBucket, dcl.DiffInfo{Type: "ReferenceType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("GcsBucket")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.AccessUrls, actual.AccessUrls, dcl.DiffInfo{OutputOnly: true, ObjectFunction: compareCertificateAuthorityAccessUrlsNewStyle, EmptyObject: EmptyCertificateAuthorityAccessUrls, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("AccessUrls")); len(ds) != 0 || err != nil {
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

	if ds, err := dcl.Diff(desired.DeleteTime, actual.DeleteTime, dcl.DiffInfo{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("DeleteTime")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.ExpireTime, actual.ExpireTime, dcl.DiffInfo{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("ExpireTime")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Labels, actual.Labels, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateCertificateAuthorityUpdateCertificateAuthorityOperation")}, fn.AddNest("Labels")); len(ds) != 0 || err != nil {
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

	if ds, err := dcl.Diff(desired.CaPool, actual.CaPool, dcl.DiffInfo{Type: "ReferenceType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("CaPool")); len(ds) != 0 || err != nil {
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
func compareCertificateAuthorityConfigNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*CertificateAuthorityConfig)
	if !ok {
		desiredNotPointer, ok := d.(CertificateAuthorityConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CertificateAuthorityConfig or *CertificateAuthorityConfig", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*CertificateAuthorityConfig)
	if !ok {
		actualNotPointer, ok := a.(CertificateAuthorityConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CertificateAuthorityConfig", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.SubjectConfig, actual.SubjectConfig, dcl.DiffInfo{ObjectFunction: compareCertificateAuthorityConfigSubjectConfigNewStyle, EmptyObject: EmptyCertificateAuthorityConfigSubjectConfig, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("SubjectConfig")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.X509Config, actual.X509Config, dcl.DiffInfo{ObjectFunction: compareCertificateAuthorityConfigX509ConfigNewStyle, EmptyObject: EmptyCertificateAuthorityConfigX509Config, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("X509Config")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.PublicKey, actual.PublicKey, dcl.DiffInfo{OutputOnly: true, ObjectFunction: compareCertificateAuthorityConfigPublicKeyNewStyle, EmptyObject: EmptyCertificateAuthorityConfigPublicKey, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("PublicKey")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareCertificateAuthorityConfigSubjectConfigNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*CertificateAuthorityConfigSubjectConfig)
	if !ok {
		desiredNotPointer, ok := d.(CertificateAuthorityConfigSubjectConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CertificateAuthorityConfigSubjectConfig or *CertificateAuthorityConfigSubjectConfig", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*CertificateAuthorityConfigSubjectConfig)
	if !ok {
		actualNotPointer, ok := a.(CertificateAuthorityConfigSubjectConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CertificateAuthorityConfigSubjectConfig", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.Subject, actual.Subject, dcl.DiffInfo{ObjectFunction: compareCertificateAuthorityConfigSubjectConfigSubjectNewStyle, EmptyObject: EmptyCertificateAuthorityConfigSubjectConfigSubject, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Subject")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.SubjectAltName, actual.SubjectAltName, dcl.DiffInfo{ObjectFunction: compareCertificateAuthorityConfigSubjectConfigSubjectAltNameNewStyle, EmptyObject: EmptyCertificateAuthorityConfigSubjectConfigSubjectAltName, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("SubjectAltName")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareCertificateAuthorityConfigSubjectConfigSubjectNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*CertificateAuthorityConfigSubjectConfigSubject)
	if !ok {
		desiredNotPointer, ok := d.(CertificateAuthorityConfigSubjectConfigSubject)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CertificateAuthorityConfigSubjectConfigSubject or *CertificateAuthorityConfigSubjectConfigSubject", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*CertificateAuthorityConfigSubjectConfigSubject)
	if !ok {
		actualNotPointer, ok := a.(CertificateAuthorityConfigSubjectConfigSubject)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CertificateAuthorityConfigSubjectConfigSubject", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.CommonName, actual.CommonName, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("CommonName")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.CountryCode, actual.CountryCode, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("CountryCode")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Organization, actual.Organization, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Organization")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.OrganizationalUnit, actual.OrganizationalUnit, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("OrganizationalUnit")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Locality, actual.Locality, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Locality")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Province, actual.Province, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Province")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.StreetAddress, actual.StreetAddress, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("StreetAddress")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.PostalCode, actual.PostalCode, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("PostalCode")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareCertificateAuthorityConfigSubjectConfigSubjectAltNameNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*CertificateAuthorityConfigSubjectConfigSubjectAltName)
	if !ok {
		desiredNotPointer, ok := d.(CertificateAuthorityConfigSubjectConfigSubjectAltName)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CertificateAuthorityConfigSubjectConfigSubjectAltName or *CertificateAuthorityConfigSubjectConfigSubjectAltName", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*CertificateAuthorityConfigSubjectConfigSubjectAltName)
	if !ok {
		actualNotPointer, ok := a.(CertificateAuthorityConfigSubjectConfigSubjectAltName)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CertificateAuthorityConfigSubjectConfigSubjectAltName", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.DnsNames, actual.DnsNames, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("DnsNames")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Uris, actual.Uris, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Uris")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.EmailAddresses, actual.EmailAddresses, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("EmailAddresses")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.IPAddresses, actual.IPAddresses, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("IpAddresses")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.CustomSans, actual.CustomSans, dcl.DiffInfo{ObjectFunction: compareCertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSansNewStyle, EmptyObject: EmptyCertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSans, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("CustomSans")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareCertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSansNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*CertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSans)
	if !ok {
		desiredNotPointer, ok := d.(CertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSans)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSans or *CertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSans", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*CertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSans)
	if !ok {
		actualNotPointer, ok := a.(CertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSans)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSans", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.ObjectId, actual.ObjectId, dcl.DiffInfo{ObjectFunction: compareCertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSansObjectIdNewStyle, EmptyObject: EmptyCertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSansObjectId, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("ObjectId")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Critical, actual.Critical, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Critical")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Value, actual.Value, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Value")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareCertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSansObjectIdNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*CertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSansObjectId)
	if !ok {
		desiredNotPointer, ok := d.(CertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSansObjectId)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSansObjectId or *CertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSansObjectId", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*CertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSansObjectId)
	if !ok {
		actualNotPointer, ok := a.(CertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSansObjectId)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSansObjectId", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.ObjectIdPath, actual.ObjectIdPath, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("ObjectIdPath")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareCertificateAuthorityConfigX509ConfigNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*CertificateAuthorityConfigX509Config)
	if !ok {
		desiredNotPointer, ok := d.(CertificateAuthorityConfigX509Config)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CertificateAuthorityConfigX509Config or *CertificateAuthorityConfigX509Config", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*CertificateAuthorityConfigX509Config)
	if !ok {
		actualNotPointer, ok := a.(CertificateAuthorityConfigX509Config)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CertificateAuthorityConfigX509Config", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.KeyUsage, actual.KeyUsage, dcl.DiffInfo{ObjectFunction: compareCertificateAuthorityConfigX509ConfigKeyUsageNewStyle, EmptyObject: EmptyCertificateAuthorityConfigX509ConfigKeyUsage, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("KeyUsage")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.CaOptions, actual.CaOptions, dcl.DiffInfo{ObjectFunction: compareCertificateAuthorityConfigX509ConfigCaOptionsNewStyle, EmptyObject: EmptyCertificateAuthorityConfigX509ConfigCaOptions, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("CaOptions")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.PolicyIds, actual.PolicyIds, dcl.DiffInfo{ObjectFunction: compareCertificateAuthorityConfigX509ConfigPolicyIdsNewStyle, EmptyObject: EmptyCertificateAuthorityConfigX509ConfigPolicyIds, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("PolicyIds")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.AiaOcspServers, actual.AiaOcspServers, dcl.DiffInfo{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("AiaOcspServers")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.AdditionalExtensions, actual.AdditionalExtensions, dcl.DiffInfo{ObjectFunction: compareCertificateAuthorityConfigX509ConfigAdditionalExtensionsNewStyle, EmptyObject: EmptyCertificateAuthorityConfigX509ConfigAdditionalExtensions, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("AdditionalExtensions")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareCertificateAuthorityConfigX509ConfigKeyUsageNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*CertificateAuthorityConfigX509ConfigKeyUsage)
	if !ok {
		desiredNotPointer, ok := d.(CertificateAuthorityConfigX509ConfigKeyUsage)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CertificateAuthorityConfigX509ConfigKeyUsage or *CertificateAuthorityConfigX509ConfigKeyUsage", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*CertificateAuthorityConfigX509ConfigKeyUsage)
	if !ok {
		actualNotPointer, ok := a.(CertificateAuthorityConfigX509ConfigKeyUsage)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CertificateAuthorityConfigX509ConfigKeyUsage", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.BaseKeyUsage, actual.BaseKeyUsage, dcl.DiffInfo{ObjectFunction: compareCertificateAuthorityConfigX509ConfigKeyUsageBaseKeyUsageNewStyle, EmptyObject: EmptyCertificateAuthorityConfigX509ConfigKeyUsageBaseKeyUsage, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("BaseKeyUsage")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.ExtendedKeyUsage, actual.ExtendedKeyUsage, dcl.DiffInfo{ObjectFunction: compareCertificateAuthorityConfigX509ConfigKeyUsageExtendedKeyUsageNewStyle, EmptyObject: EmptyCertificateAuthorityConfigX509ConfigKeyUsageExtendedKeyUsage, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("ExtendedKeyUsage")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.UnknownExtendedKeyUsages, actual.UnknownExtendedKeyUsages, dcl.DiffInfo{ObjectFunction: compareCertificateAuthorityConfigX509ConfigKeyUsageUnknownExtendedKeyUsagesNewStyle, EmptyObject: EmptyCertificateAuthorityConfigX509ConfigKeyUsageUnknownExtendedKeyUsages, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("UnknownExtendedKeyUsages")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareCertificateAuthorityConfigX509ConfigKeyUsageBaseKeyUsageNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*CertificateAuthorityConfigX509ConfigKeyUsageBaseKeyUsage)
	if !ok {
		desiredNotPointer, ok := d.(CertificateAuthorityConfigX509ConfigKeyUsageBaseKeyUsage)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CertificateAuthorityConfigX509ConfigKeyUsageBaseKeyUsage or *CertificateAuthorityConfigX509ConfigKeyUsageBaseKeyUsage", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*CertificateAuthorityConfigX509ConfigKeyUsageBaseKeyUsage)
	if !ok {
		actualNotPointer, ok := a.(CertificateAuthorityConfigX509ConfigKeyUsageBaseKeyUsage)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CertificateAuthorityConfigX509ConfigKeyUsageBaseKeyUsage", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.DigitalSignature, actual.DigitalSignature, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("DigitalSignature")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.ContentCommitment, actual.ContentCommitment, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("ContentCommitment")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.KeyEncipherment, actual.KeyEncipherment, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("KeyEncipherment")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.DataEncipherment, actual.DataEncipherment, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("DataEncipherment")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.KeyAgreement, actual.KeyAgreement, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("KeyAgreement")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.CertSign, actual.CertSign, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("CertSign")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.CrlSign, actual.CrlSign, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("CrlSign")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.EncipherOnly, actual.EncipherOnly, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("EncipherOnly")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.DecipherOnly, actual.DecipherOnly, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("DecipherOnly")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareCertificateAuthorityConfigX509ConfigKeyUsageExtendedKeyUsageNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*CertificateAuthorityConfigX509ConfigKeyUsageExtendedKeyUsage)
	if !ok {
		desiredNotPointer, ok := d.(CertificateAuthorityConfigX509ConfigKeyUsageExtendedKeyUsage)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CertificateAuthorityConfigX509ConfigKeyUsageExtendedKeyUsage or *CertificateAuthorityConfigX509ConfigKeyUsageExtendedKeyUsage", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*CertificateAuthorityConfigX509ConfigKeyUsageExtendedKeyUsage)
	if !ok {
		actualNotPointer, ok := a.(CertificateAuthorityConfigX509ConfigKeyUsageExtendedKeyUsage)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CertificateAuthorityConfigX509ConfigKeyUsageExtendedKeyUsage", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.ServerAuth, actual.ServerAuth, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("ServerAuth")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.ClientAuth, actual.ClientAuth, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("ClientAuth")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.CodeSigning, actual.CodeSigning, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("CodeSigning")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.EmailProtection, actual.EmailProtection, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("EmailProtection")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.TimeStamping, actual.TimeStamping, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("TimeStamping")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.OcspSigning, actual.OcspSigning, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("OcspSigning")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareCertificateAuthorityConfigX509ConfigKeyUsageUnknownExtendedKeyUsagesNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*CertificateAuthorityConfigX509ConfigKeyUsageUnknownExtendedKeyUsages)
	if !ok {
		desiredNotPointer, ok := d.(CertificateAuthorityConfigX509ConfigKeyUsageUnknownExtendedKeyUsages)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CertificateAuthorityConfigX509ConfigKeyUsageUnknownExtendedKeyUsages or *CertificateAuthorityConfigX509ConfigKeyUsageUnknownExtendedKeyUsages", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*CertificateAuthorityConfigX509ConfigKeyUsageUnknownExtendedKeyUsages)
	if !ok {
		actualNotPointer, ok := a.(CertificateAuthorityConfigX509ConfigKeyUsageUnknownExtendedKeyUsages)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CertificateAuthorityConfigX509ConfigKeyUsageUnknownExtendedKeyUsages", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.ObjectIdPath, actual.ObjectIdPath, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("ObjectIdPath")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareCertificateAuthorityConfigX509ConfigCaOptionsNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*CertificateAuthorityConfigX509ConfigCaOptions)
	if !ok {
		desiredNotPointer, ok := d.(CertificateAuthorityConfigX509ConfigCaOptions)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CertificateAuthorityConfigX509ConfigCaOptions or *CertificateAuthorityConfigX509ConfigCaOptions", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*CertificateAuthorityConfigX509ConfigCaOptions)
	if !ok {
		actualNotPointer, ok := a.(CertificateAuthorityConfigX509ConfigCaOptions)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CertificateAuthorityConfigX509ConfigCaOptions", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.IsCa, actual.IsCa, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("IsCa")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.MaxIssuerPathLength, actual.MaxIssuerPathLength, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("MaxIssuerPathLength")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.ZeroMaxIssuerPathLength, actual.ZeroMaxIssuerPathLength, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("ZeroMaxIssuerPathLength")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareCertificateAuthorityConfigX509ConfigPolicyIdsNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*CertificateAuthorityConfigX509ConfigPolicyIds)
	if !ok {
		desiredNotPointer, ok := d.(CertificateAuthorityConfigX509ConfigPolicyIds)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CertificateAuthorityConfigX509ConfigPolicyIds or *CertificateAuthorityConfigX509ConfigPolicyIds", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*CertificateAuthorityConfigX509ConfigPolicyIds)
	if !ok {
		actualNotPointer, ok := a.(CertificateAuthorityConfigX509ConfigPolicyIds)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CertificateAuthorityConfigX509ConfigPolicyIds", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.ObjectIdPath, actual.ObjectIdPath, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("ObjectIdPath")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareCertificateAuthorityConfigX509ConfigAdditionalExtensionsNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*CertificateAuthorityConfigX509ConfigAdditionalExtensions)
	if !ok {
		desiredNotPointer, ok := d.(CertificateAuthorityConfigX509ConfigAdditionalExtensions)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CertificateAuthorityConfigX509ConfigAdditionalExtensions or *CertificateAuthorityConfigX509ConfigAdditionalExtensions", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*CertificateAuthorityConfigX509ConfigAdditionalExtensions)
	if !ok {
		actualNotPointer, ok := a.(CertificateAuthorityConfigX509ConfigAdditionalExtensions)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CertificateAuthorityConfigX509ConfigAdditionalExtensions", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.ObjectId, actual.ObjectId, dcl.DiffInfo{ObjectFunction: compareCertificateAuthorityConfigX509ConfigAdditionalExtensionsObjectIdNewStyle, EmptyObject: EmptyCertificateAuthorityConfigX509ConfigAdditionalExtensionsObjectId, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("ObjectId")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Critical, actual.Critical, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Critical")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Value, actual.Value, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Value")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareCertificateAuthorityConfigX509ConfigAdditionalExtensionsObjectIdNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*CertificateAuthorityConfigX509ConfigAdditionalExtensionsObjectId)
	if !ok {
		desiredNotPointer, ok := d.(CertificateAuthorityConfigX509ConfigAdditionalExtensionsObjectId)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CertificateAuthorityConfigX509ConfigAdditionalExtensionsObjectId or *CertificateAuthorityConfigX509ConfigAdditionalExtensionsObjectId", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*CertificateAuthorityConfigX509ConfigAdditionalExtensionsObjectId)
	if !ok {
		actualNotPointer, ok := a.(CertificateAuthorityConfigX509ConfigAdditionalExtensionsObjectId)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CertificateAuthorityConfigX509ConfigAdditionalExtensionsObjectId", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.ObjectIdPath, actual.ObjectIdPath, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("ObjectIdPath")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareCertificateAuthorityConfigPublicKeyNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*CertificateAuthorityConfigPublicKey)
	if !ok {
		desiredNotPointer, ok := d.(CertificateAuthorityConfigPublicKey)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CertificateAuthorityConfigPublicKey or *CertificateAuthorityConfigPublicKey", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*CertificateAuthorityConfigPublicKey)
	if !ok {
		actualNotPointer, ok := a.(CertificateAuthorityConfigPublicKey)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CertificateAuthorityConfigPublicKey", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.Key, actual.Key, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Key")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Format, actual.Format, dcl.DiffInfo{Type: "EnumType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Format")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareCertificateAuthorityKeySpecNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*CertificateAuthorityKeySpec)
	if !ok {
		desiredNotPointer, ok := d.(CertificateAuthorityKeySpec)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CertificateAuthorityKeySpec or *CertificateAuthorityKeySpec", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*CertificateAuthorityKeySpec)
	if !ok {
		actualNotPointer, ok := a.(CertificateAuthorityKeySpec)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CertificateAuthorityKeySpec", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.CloudKmsKeyVersion, actual.CloudKmsKeyVersion, dcl.DiffInfo{Type: "ReferenceType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("CloudKmsKeyVersion")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Algorithm, actual.Algorithm, dcl.DiffInfo{Type: "EnumType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Algorithm")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareCertificateAuthoritySubordinateConfigNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*CertificateAuthoritySubordinateConfig)
	if !ok {
		desiredNotPointer, ok := d.(CertificateAuthoritySubordinateConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CertificateAuthoritySubordinateConfig or *CertificateAuthoritySubordinateConfig", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*CertificateAuthoritySubordinateConfig)
	if !ok {
		actualNotPointer, ok := a.(CertificateAuthoritySubordinateConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CertificateAuthoritySubordinateConfig", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.CertificateAuthority, actual.CertificateAuthority, dcl.DiffInfo{Type: "ReferenceType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("CertificateAuthority")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.PemIssuerChain, actual.PemIssuerChain, dcl.DiffInfo{ObjectFunction: compareCertificateAuthoritySubordinateConfigPemIssuerChainNewStyle, EmptyObject: EmptyCertificateAuthoritySubordinateConfigPemIssuerChain, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("PemIssuerChain")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareCertificateAuthoritySubordinateConfigPemIssuerChainNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*CertificateAuthoritySubordinateConfigPemIssuerChain)
	if !ok {
		desiredNotPointer, ok := d.(CertificateAuthoritySubordinateConfigPemIssuerChain)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CertificateAuthoritySubordinateConfigPemIssuerChain or *CertificateAuthoritySubordinateConfigPemIssuerChain", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*CertificateAuthoritySubordinateConfigPemIssuerChain)
	if !ok {
		actualNotPointer, ok := a.(CertificateAuthoritySubordinateConfigPemIssuerChain)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CertificateAuthoritySubordinateConfigPemIssuerChain", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.PemCertificates, actual.PemCertificates, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("PemCertificates")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareCertificateAuthorityCaCertificateDescriptionsNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*CertificateAuthorityCaCertificateDescriptions)
	if !ok {
		desiredNotPointer, ok := d.(CertificateAuthorityCaCertificateDescriptions)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CertificateAuthorityCaCertificateDescriptions or *CertificateAuthorityCaCertificateDescriptions", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*CertificateAuthorityCaCertificateDescriptions)
	if !ok {
		actualNotPointer, ok := a.(CertificateAuthorityCaCertificateDescriptions)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CertificateAuthorityCaCertificateDescriptions", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.SubjectDescription, actual.SubjectDescription, dcl.DiffInfo{ObjectFunction: compareCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionNewStyle, EmptyObject: EmptyCertificateAuthorityCaCertificateDescriptionsSubjectDescription, OperationSelector: dcl.TriggersOperation("updateCertificateAuthorityUpdateCertificateAuthorityOperation")}, fn.AddNest("SubjectDescription")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.X509Description, actual.X509Description, dcl.DiffInfo{ObjectFunction: compareCertificateAuthorityCaCertificateDescriptionsX509DescriptionNewStyle, EmptyObject: EmptyCertificateAuthorityCaCertificateDescriptionsX509Description, OperationSelector: dcl.TriggersOperation("updateCertificateAuthorityUpdateCertificateAuthorityOperation")}, fn.AddNest("X509Description")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.PublicKey, actual.PublicKey, dcl.DiffInfo{ObjectFunction: compareCertificateAuthorityCaCertificateDescriptionsPublicKeyNewStyle, EmptyObject: EmptyCertificateAuthorityCaCertificateDescriptionsPublicKey, OperationSelector: dcl.TriggersOperation("updateCertificateAuthorityUpdateCertificateAuthorityOperation")}, fn.AddNest("PublicKey")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.SubjectKeyId, actual.SubjectKeyId, dcl.DiffInfo{ObjectFunction: compareCertificateAuthorityCaCertificateDescriptionsSubjectKeyIdNewStyle, EmptyObject: EmptyCertificateAuthorityCaCertificateDescriptionsSubjectKeyId, OperationSelector: dcl.TriggersOperation("updateCertificateAuthorityUpdateCertificateAuthorityOperation")}, fn.AddNest("SubjectKeyId")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.AuthorityKeyId, actual.AuthorityKeyId, dcl.DiffInfo{ObjectFunction: compareCertificateAuthorityCaCertificateDescriptionsAuthorityKeyIdNewStyle, EmptyObject: EmptyCertificateAuthorityCaCertificateDescriptionsAuthorityKeyId, OperationSelector: dcl.TriggersOperation("updateCertificateAuthorityUpdateCertificateAuthorityOperation")}, fn.AddNest("AuthorityKeyId")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.CrlDistributionPoints, actual.CrlDistributionPoints, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateCertificateAuthorityUpdateCertificateAuthorityOperation")}, fn.AddNest("CrlDistributionPoints")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.AiaIssuingCertificateUrls, actual.AiaIssuingCertificateUrls, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateCertificateAuthorityUpdateCertificateAuthorityOperation")}, fn.AddNest("AiaIssuingCertificateUrls")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.CertFingerprint, actual.CertFingerprint, dcl.DiffInfo{ObjectFunction: compareCertificateAuthorityCaCertificateDescriptionsCertFingerprintNewStyle, EmptyObject: EmptyCertificateAuthorityCaCertificateDescriptionsCertFingerprint, OperationSelector: dcl.TriggersOperation("updateCertificateAuthorityUpdateCertificateAuthorityOperation")}, fn.AddNest("CertFingerprint")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*CertificateAuthorityCaCertificateDescriptionsSubjectDescription)
	if !ok {
		desiredNotPointer, ok := d.(CertificateAuthorityCaCertificateDescriptionsSubjectDescription)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CertificateAuthorityCaCertificateDescriptionsSubjectDescription or *CertificateAuthorityCaCertificateDescriptionsSubjectDescription", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*CertificateAuthorityCaCertificateDescriptionsSubjectDescription)
	if !ok {
		actualNotPointer, ok := a.(CertificateAuthorityCaCertificateDescriptionsSubjectDescription)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CertificateAuthorityCaCertificateDescriptionsSubjectDescription", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.Subject, actual.Subject, dcl.DiffInfo{ObjectFunction: compareCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectNewStyle, EmptyObject: EmptyCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubject, OperationSelector: dcl.TriggersOperation("updateCertificateAuthorityUpdateCertificateAuthorityOperation")}, fn.AddNest("Subject")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.SubjectAltName, actual.SubjectAltName, dcl.DiffInfo{ObjectFunction: compareCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameNewStyle, EmptyObject: EmptyCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltName, OperationSelector: dcl.TriggersOperation("updateCertificateAuthorityUpdateCertificateAuthorityOperation")}, fn.AddNest("SubjectAltName")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.HexSerialNumber, actual.HexSerialNumber, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateCertificateAuthorityUpdateCertificateAuthorityOperation")}, fn.AddNest("HexSerialNumber")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Lifetime, actual.Lifetime, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateCertificateAuthorityUpdateCertificateAuthorityOperation")}, fn.AddNest("Lifetime")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.NotBeforeTime, actual.NotBeforeTime, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateCertificateAuthorityUpdateCertificateAuthorityOperation")}, fn.AddNest("NotBeforeTime")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.NotAfterTime, actual.NotAfterTime, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateCertificateAuthorityUpdateCertificateAuthorityOperation")}, fn.AddNest("NotAfterTime")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubject)
	if !ok {
		desiredNotPointer, ok := d.(CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubject)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubject or *CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubject", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubject)
	if !ok {
		actualNotPointer, ok := a.(CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubject)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubject", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.CommonName, actual.CommonName, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateCertificateAuthorityUpdateCertificateAuthorityOperation")}, fn.AddNest("CommonName")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.CountryCode, actual.CountryCode, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateCertificateAuthorityUpdateCertificateAuthorityOperation")}, fn.AddNest("CountryCode")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Organization, actual.Organization, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateCertificateAuthorityUpdateCertificateAuthorityOperation")}, fn.AddNest("Organization")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.OrganizationalUnit, actual.OrganizationalUnit, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateCertificateAuthorityUpdateCertificateAuthorityOperation")}, fn.AddNest("OrganizationalUnit")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Locality, actual.Locality, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateCertificateAuthorityUpdateCertificateAuthorityOperation")}, fn.AddNest("Locality")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Province, actual.Province, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateCertificateAuthorityUpdateCertificateAuthorityOperation")}, fn.AddNest("Province")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.StreetAddress, actual.StreetAddress, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateCertificateAuthorityUpdateCertificateAuthorityOperation")}, fn.AddNest("StreetAddress")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.PostalCode, actual.PostalCode, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateCertificateAuthorityUpdateCertificateAuthorityOperation")}, fn.AddNest("PostalCode")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltName)
	if !ok {
		desiredNotPointer, ok := d.(CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltName)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltName or *CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltName", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltName)
	if !ok {
		actualNotPointer, ok := a.(CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltName)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltName", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.DnsNames, actual.DnsNames, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateCertificateAuthorityUpdateCertificateAuthorityOperation")}, fn.AddNest("DnsNames")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Uris, actual.Uris, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateCertificateAuthorityUpdateCertificateAuthorityOperation")}, fn.AddNest("Uris")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.EmailAddresses, actual.EmailAddresses, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateCertificateAuthorityUpdateCertificateAuthorityOperation")}, fn.AddNest("EmailAddresses")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.IPAddresses, actual.IPAddresses, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateCertificateAuthorityUpdateCertificateAuthorityOperation")}, fn.AddNest("IpAddresses")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.CustomSans, actual.CustomSans, dcl.DiffInfo{ObjectFunction: compareCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSansNewStyle, EmptyObject: EmptyCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSans, OperationSelector: dcl.TriggersOperation("updateCertificateAuthorityUpdateCertificateAuthorityOperation")}, fn.AddNest("CustomSans")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSansNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSans)
	if !ok {
		desiredNotPointer, ok := d.(CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSans)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSans or *CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSans", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSans)
	if !ok {
		actualNotPointer, ok := a.(CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSans)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSans", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.ObjectId, actual.ObjectId, dcl.DiffInfo{ObjectFunction: compareCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSansObjectIdNewStyle, EmptyObject: EmptyCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSansObjectId, OperationSelector: dcl.TriggersOperation("updateCertificateAuthorityUpdateCertificateAuthorityOperation")}, fn.AddNest("ObjectId")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Critical, actual.Critical, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateCertificateAuthorityUpdateCertificateAuthorityOperation")}, fn.AddNest("Critical")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Value, actual.Value, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateCertificateAuthorityUpdateCertificateAuthorityOperation")}, fn.AddNest("Value")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSansObjectIdNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSansObjectId)
	if !ok {
		desiredNotPointer, ok := d.(CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSansObjectId)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSansObjectId or *CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSansObjectId", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSansObjectId)
	if !ok {
		actualNotPointer, ok := a.(CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSansObjectId)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSansObjectId", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.ObjectIdPath, actual.ObjectIdPath, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateCertificateAuthorityUpdateCertificateAuthorityOperation")}, fn.AddNest("ObjectIdPath")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareCertificateAuthorityCaCertificateDescriptionsX509DescriptionNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*CertificateAuthorityCaCertificateDescriptionsX509Description)
	if !ok {
		desiredNotPointer, ok := d.(CertificateAuthorityCaCertificateDescriptionsX509Description)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CertificateAuthorityCaCertificateDescriptionsX509Description or *CertificateAuthorityCaCertificateDescriptionsX509Description", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*CertificateAuthorityCaCertificateDescriptionsX509Description)
	if !ok {
		actualNotPointer, ok := a.(CertificateAuthorityCaCertificateDescriptionsX509Description)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CertificateAuthorityCaCertificateDescriptionsX509Description", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.KeyUsage, actual.KeyUsage, dcl.DiffInfo{ObjectFunction: compareCertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageNewStyle, EmptyObject: EmptyCertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsage, OperationSelector: dcl.TriggersOperation("updateCertificateAuthorityUpdateCertificateAuthorityOperation")}, fn.AddNest("KeyUsage")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.CaOptions, actual.CaOptions, dcl.DiffInfo{ObjectFunction: compareCertificateAuthorityCaCertificateDescriptionsX509DescriptionCaOptionsNewStyle, EmptyObject: EmptyCertificateAuthorityCaCertificateDescriptionsX509DescriptionCaOptions, OperationSelector: dcl.TriggersOperation("updateCertificateAuthorityUpdateCertificateAuthorityOperation")}, fn.AddNest("CaOptions")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.PolicyIds, actual.PolicyIds, dcl.DiffInfo{ObjectFunction: compareCertificateAuthorityCaCertificateDescriptionsX509DescriptionPolicyIdsNewStyle, EmptyObject: EmptyCertificateAuthorityCaCertificateDescriptionsX509DescriptionPolicyIds, OperationSelector: dcl.TriggersOperation("updateCertificateAuthorityUpdateCertificateAuthorityOperation")}, fn.AddNest("PolicyIds")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.AiaOcspServers, actual.AiaOcspServers, dcl.DiffInfo{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("AiaOcspServers")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.AdditionalExtensions, actual.AdditionalExtensions, dcl.DiffInfo{ObjectFunction: compareCertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensionsNewStyle, EmptyObject: EmptyCertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensions, OperationSelector: dcl.TriggersOperation("updateCertificateAuthorityUpdateCertificateAuthorityOperation")}, fn.AddNest("AdditionalExtensions")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareCertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsage)
	if !ok {
		desiredNotPointer, ok := d.(CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsage)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsage or *CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsage", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsage)
	if !ok {
		actualNotPointer, ok := a.(CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsage)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsage", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.BaseKeyUsage, actual.BaseKeyUsage, dcl.DiffInfo{ObjectFunction: compareCertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageBaseKeyUsageNewStyle, EmptyObject: EmptyCertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageBaseKeyUsage, OperationSelector: dcl.TriggersOperation("updateCertificateAuthorityUpdateCertificateAuthorityOperation")}, fn.AddNest("BaseKeyUsage")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.ExtendedKeyUsage, actual.ExtendedKeyUsage, dcl.DiffInfo{ObjectFunction: compareCertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageExtendedKeyUsageNewStyle, EmptyObject: EmptyCertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageExtendedKeyUsage, OperationSelector: dcl.TriggersOperation("updateCertificateAuthorityUpdateCertificateAuthorityOperation")}, fn.AddNest("ExtendedKeyUsage")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.UnknownExtendedKeyUsages, actual.UnknownExtendedKeyUsages, dcl.DiffInfo{ObjectFunction: compareCertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageUnknownExtendedKeyUsagesNewStyle, EmptyObject: EmptyCertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageUnknownExtendedKeyUsages, OperationSelector: dcl.TriggersOperation("updateCertificateAuthorityUpdateCertificateAuthorityOperation")}, fn.AddNest("UnknownExtendedKeyUsages")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareCertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageBaseKeyUsageNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageBaseKeyUsage)
	if !ok {
		desiredNotPointer, ok := d.(CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageBaseKeyUsage)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageBaseKeyUsage or *CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageBaseKeyUsage", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageBaseKeyUsage)
	if !ok {
		actualNotPointer, ok := a.(CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageBaseKeyUsage)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageBaseKeyUsage", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.DigitalSignature, actual.DigitalSignature, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateCertificateAuthorityUpdateCertificateAuthorityOperation")}, fn.AddNest("DigitalSignature")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.ContentCommitment, actual.ContentCommitment, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateCertificateAuthorityUpdateCertificateAuthorityOperation")}, fn.AddNest("ContentCommitment")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.KeyEncipherment, actual.KeyEncipherment, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateCertificateAuthorityUpdateCertificateAuthorityOperation")}, fn.AddNest("KeyEncipherment")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.DataEncipherment, actual.DataEncipherment, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateCertificateAuthorityUpdateCertificateAuthorityOperation")}, fn.AddNest("DataEncipherment")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.KeyAgreement, actual.KeyAgreement, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateCertificateAuthorityUpdateCertificateAuthorityOperation")}, fn.AddNest("KeyAgreement")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.CertSign, actual.CertSign, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateCertificateAuthorityUpdateCertificateAuthorityOperation")}, fn.AddNest("CertSign")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.CrlSign, actual.CrlSign, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateCertificateAuthorityUpdateCertificateAuthorityOperation")}, fn.AddNest("CrlSign")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.EncipherOnly, actual.EncipherOnly, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateCertificateAuthorityUpdateCertificateAuthorityOperation")}, fn.AddNest("EncipherOnly")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.DecipherOnly, actual.DecipherOnly, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateCertificateAuthorityUpdateCertificateAuthorityOperation")}, fn.AddNest("DecipherOnly")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareCertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageExtendedKeyUsageNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageExtendedKeyUsage)
	if !ok {
		desiredNotPointer, ok := d.(CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageExtendedKeyUsage)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageExtendedKeyUsage or *CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageExtendedKeyUsage", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageExtendedKeyUsage)
	if !ok {
		actualNotPointer, ok := a.(CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageExtendedKeyUsage)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageExtendedKeyUsage", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.ServerAuth, actual.ServerAuth, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateCertificateAuthorityUpdateCertificateAuthorityOperation")}, fn.AddNest("ServerAuth")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.ClientAuth, actual.ClientAuth, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateCertificateAuthorityUpdateCertificateAuthorityOperation")}, fn.AddNest("ClientAuth")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.CodeSigning, actual.CodeSigning, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateCertificateAuthorityUpdateCertificateAuthorityOperation")}, fn.AddNest("CodeSigning")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.EmailProtection, actual.EmailProtection, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateCertificateAuthorityUpdateCertificateAuthorityOperation")}, fn.AddNest("EmailProtection")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.TimeStamping, actual.TimeStamping, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateCertificateAuthorityUpdateCertificateAuthorityOperation")}, fn.AddNest("TimeStamping")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.OcspSigning, actual.OcspSigning, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateCertificateAuthorityUpdateCertificateAuthorityOperation")}, fn.AddNest("OcspSigning")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareCertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageUnknownExtendedKeyUsagesNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageUnknownExtendedKeyUsages)
	if !ok {
		desiredNotPointer, ok := d.(CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageUnknownExtendedKeyUsages)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageUnknownExtendedKeyUsages or *CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageUnknownExtendedKeyUsages", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageUnknownExtendedKeyUsages)
	if !ok {
		actualNotPointer, ok := a.(CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageUnknownExtendedKeyUsages)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageUnknownExtendedKeyUsages", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.ObjectIdPath, actual.ObjectIdPath, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateCertificateAuthorityUpdateCertificateAuthorityOperation")}, fn.AddNest("ObjectIdPath")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareCertificateAuthorityCaCertificateDescriptionsX509DescriptionCaOptionsNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*CertificateAuthorityCaCertificateDescriptionsX509DescriptionCaOptions)
	if !ok {
		desiredNotPointer, ok := d.(CertificateAuthorityCaCertificateDescriptionsX509DescriptionCaOptions)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CertificateAuthorityCaCertificateDescriptionsX509DescriptionCaOptions or *CertificateAuthorityCaCertificateDescriptionsX509DescriptionCaOptions", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*CertificateAuthorityCaCertificateDescriptionsX509DescriptionCaOptions)
	if !ok {
		actualNotPointer, ok := a.(CertificateAuthorityCaCertificateDescriptionsX509DescriptionCaOptions)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CertificateAuthorityCaCertificateDescriptionsX509DescriptionCaOptions", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.IsCa, actual.IsCa, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateCertificateAuthorityUpdateCertificateAuthorityOperation")}, fn.AddNest("IsCa")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.MaxIssuerPathLength, actual.MaxIssuerPathLength, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateCertificateAuthorityUpdateCertificateAuthorityOperation")}, fn.AddNest("MaxIssuerPathLength")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareCertificateAuthorityCaCertificateDescriptionsX509DescriptionPolicyIdsNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*CertificateAuthorityCaCertificateDescriptionsX509DescriptionPolicyIds)
	if !ok {
		desiredNotPointer, ok := d.(CertificateAuthorityCaCertificateDescriptionsX509DescriptionPolicyIds)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CertificateAuthorityCaCertificateDescriptionsX509DescriptionPolicyIds or *CertificateAuthorityCaCertificateDescriptionsX509DescriptionPolicyIds", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*CertificateAuthorityCaCertificateDescriptionsX509DescriptionPolicyIds)
	if !ok {
		actualNotPointer, ok := a.(CertificateAuthorityCaCertificateDescriptionsX509DescriptionPolicyIds)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CertificateAuthorityCaCertificateDescriptionsX509DescriptionPolicyIds", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.ObjectIdPath, actual.ObjectIdPath, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateCertificateAuthorityUpdateCertificateAuthorityOperation")}, fn.AddNest("ObjectIdPath")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareCertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensionsNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*CertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensions)
	if !ok {
		desiredNotPointer, ok := d.(CertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensions)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensions or *CertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensions", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*CertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensions)
	if !ok {
		actualNotPointer, ok := a.(CertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensions)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensions", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.ObjectId, actual.ObjectId, dcl.DiffInfo{ObjectFunction: compareCertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensionsObjectIdNewStyle, EmptyObject: EmptyCertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensionsObjectId, OperationSelector: dcl.TriggersOperation("updateCertificateAuthorityUpdateCertificateAuthorityOperation")}, fn.AddNest("ObjectId")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Critical, actual.Critical, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateCertificateAuthorityUpdateCertificateAuthorityOperation")}, fn.AddNest("Critical")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Value, actual.Value, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateCertificateAuthorityUpdateCertificateAuthorityOperation")}, fn.AddNest("Value")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareCertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensionsObjectIdNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*CertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensionsObjectId)
	if !ok {
		desiredNotPointer, ok := d.(CertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensionsObjectId)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensionsObjectId or *CertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensionsObjectId", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*CertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensionsObjectId)
	if !ok {
		actualNotPointer, ok := a.(CertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensionsObjectId)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensionsObjectId", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.ObjectIdPath, actual.ObjectIdPath, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateCertificateAuthorityUpdateCertificateAuthorityOperation")}, fn.AddNest("ObjectIdPath")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareCertificateAuthorityCaCertificateDescriptionsPublicKeyNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*CertificateAuthorityCaCertificateDescriptionsPublicKey)
	if !ok {
		desiredNotPointer, ok := d.(CertificateAuthorityCaCertificateDescriptionsPublicKey)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CertificateAuthorityCaCertificateDescriptionsPublicKey or *CertificateAuthorityCaCertificateDescriptionsPublicKey", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*CertificateAuthorityCaCertificateDescriptionsPublicKey)
	if !ok {
		actualNotPointer, ok := a.(CertificateAuthorityCaCertificateDescriptionsPublicKey)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CertificateAuthorityCaCertificateDescriptionsPublicKey", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.Key, actual.Key, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateCertificateAuthorityUpdateCertificateAuthorityOperation")}, fn.AddNest("Key")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Format, actual.Format, dcl.DiffInfo{Type: "EnumType", OperationSelector: dcl.TriggersOperation("updateCertificateAuthorityUpdateCertificateAuthorityOperation")}, fn.AddNest("Format")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareCertificateAuthorityCaCertificateDescriptionsSubjectKeyIdNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*CertificateAuthorityCaCertificateDescriptionsSubjectKeyId)
	if !ok {
		desiredNotPointer, ok := d.(CertificateAuthorityCaCertificateDescriptionsSubjectKeyId)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CertificateAuthorityCaCertificateDescriptionsSubjectKeyId or *CertificateAuthorityCaCertificateDescriptionsSubjectKeyId", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*CertificateAuthorityCaCertificateDescriptionsSubjectKeyId)
	if !ok {
		actualNotPointer, ok := a.(CertificateAuthorityCaCertificateDescriptionsSubjectKeyId)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CertificateAuthorityCaCertificateDescriptionsSubjectKeyId", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.KeyId, actual.KeyId, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateCertificateAuthorityUpdateCertificateAuthorityOperation")}, fn.AddNest("KeyId")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareCertificateAuthorityCaCertificateDescriptionsAuthorityKeyIdNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*CertificateAuthorityCaCertificateDescriptionsAuthorityKeyId)
	if !ok {
		desiredNotPointer, ok := d.(CertificateAuthorityCaCertificateDescriptionsAuthorityKeyId)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CertificateAuthorityCaCertificateDescriptionsAuthorityKeyId or *CertificateAuthorityCaCertificateDescriptionsAuthorityKeyId", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*CertificateAuthorityCaCertificateDescriptionsAuthorityKeyId)
	if !ok {
		actualNotPointer, ok := a.(CertificateAuthorityCaCertificateDescriptionsAuthorityKeyId)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CertificateAuthorityCaCertificateDescriptionsAuthorityKeyId", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.KeyId, actual.KeyId, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateCertificateAuthorityUpdateCertificateAuthorityOperation")}, fn.AddNest("KeyId")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareCertificateAuthorityCaCertificateDescriptionsCertFingerprintNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*CertificateAuthorityCaCertificateDescriptionsCertFingerprint)
	if !ok {
		desiredNotPointer, ok := d.(CertificateAuthorityCaCertificateDescriptionsCertFingerprint)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CertificateAuthorityCaCertificateDescriptionsCertFingerprint or *CertificateAuthorityCaCertificateDescriptionsCertFingerprint", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*CertificateAuthorityCaCertificateDescriptionsCertFingerprint)
	if !ok {
		actualNotPointer, ok := a.(CertificateAuthorityCaCertificateDescriptionsCertFingerprint)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CertificateAuthorityCaCertificateDescriptionsCertFingerprint", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.Sha256Hash, actual.Sha256Hash, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateCertificateAuthorityUpdateCertificateAuthorityOperation")}, fn.AddNest("Sha256Hash")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareCertificateAuthorityAccessUrlsNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*CertificateAuthorityAccessUrls)
	if !ok {
		desiredNotPointer, ok := d.(CertificateAuthorityAccessUrls)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CertificateAuthorityAccessUrls or *CertificateAuthorityAccessUrls", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*CertificateAuthorityAccessUrls)
	if !ok {
		actualNotPointer, ok := a.(CertificateAuthorityAccessUrls)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CertificateAuthorityAccessUrls", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.CaCertificateAccessUrl, actual.CaCertificateAccessUrl, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("CaCertificateAccessUrl")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.CrlAccessUrls, actual.CrlAccessUrls, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("CrlAccessUrls")); len(ds) != 0 || err != nil {
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
func (r *CertificateAuthority) urlNormalized() *CertificateAuthority {
	normalized := dcl.Copy(*r).(CertificateAuthority)
	normalized.Name = dcl.SelfLinkToName(r.Name)
	normalized.Lifetime = dcl.SelfLinkToName(r.Lifetime)
	normalized.GcsBucket = dcl.SelfLinkToName(r.GcsBucket)
	normalized.Project = dcl.SelfLinkToName(r.Project)
	normalized.Location = dcl.SelfLinkToName(r.Location)
	normalized.CaPool = dcl.SelfLinkToName(r.CaPool)
	return &normalized
}

func (r *CertificateAuthority) updateURL(userBasePath, updateName string) (string, error) {
	nr := r.urlNormalized()
	if updateName == "UpdateCertificateAuthority" {
		fields := map[string]interface{}{
			"project":  dcl.ValueOrEmptyString(nr.Project),
			"location": dcl.ValueOrEmptyString(nr.Location),
			"caPool":   dcl.ValueOrEmptyString(nr.CaPool),
			"name":     dcl.ValueOrEmptyString(nr.Name),
		}
		return dcl.URL("projects/{{project}}/locations/{{location}}/caPools/{{caPool}}/certificateAuthorities/{{name}}", nr.basePath(), userBasePath, fields), nil

	}

	return "", fmt.Errorf("unknown update name: %s", updateName)
}

// marshal encodes the CertificateAuthority resource into JSON for a Create request, and
// performs transformations from the resource schema to the API schema if
// necessary.
func (r *CertificateAuthority) marshal(c *Client) ([]byte, error) {
	m, err := expandCertificateAuthority(c, r)
	if err != nil {
		return nil, fmt.Errorf("error marshalling CertificateAuthority: %w", err)
	}

	return json.Marshal(m)
}

// unmarshalCertificateAuthority decodes JSON responses into the CertificateAuthority resource schema.
func unmarshalCertificateAuthority(b []byte, c *Client, res *CertificateAuthority) (*CertificateAuthority, error) {
	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, err
	}
	return unmarshalMapCertificateAuthority(m, c, res)
}

func unmarshalMapCertificateAuthority(m map[string]interface{}, c *Client, res *CertificateAuthority) (*CertificateAuthority, error) {

	flattened := flattenCertificateAuthority(c, m, res)
	if flattened == nil {
		return nil, fmt.Errorf("attempted to flatten empty json object")
	}
	return flattened, nil
}

// expandCertificateAuthority expands CertificateAuthority into a JSON request object.
func expandCertificateAuthority(c *Client, f *CertificateAuthority) (map[string]interface{}, error) {
	m := make(map[string]interface{})
	res := f
	_ = res
	if v, err := dcl.DeriveField("projects/%s/locations/%s/caPools/%s/certificateAuthorities/%s", f.Name, dcl.SelfLinkToName(f.Project), dcl.SelfLinkToName(f.Location), dcl.SelfLinkToName(f.CaPool), dcl.SelfLinkToName(f.Name)); err != nil {
		return nil, fmt.Errorf("error expanding Name into name: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["name"] = v
	}
	if v := f.Type; dcl.ValueShouldBeSent(v) {
		m["type"] = v
	}
	if v, err := expandCertificateAuthorityConfig(c, f.Config, res); err != nil {
		return nil, fmt.Errorf("error expanding Config into config: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["config"] = v
	}
	if v := f.Lifetime; dcl.ValueShouldBeSent(v) {
		m["lifetime"] = v
	}
	if v, err := expandCertificateAuthorityKeySpec(c, f.KeySpec, res); err != nil {
		return nil, fmt.Errorf("error expanding KeySpec into keySpec: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["keySpec"] = v
	}
	if v := f.GcsBucket; dcl.ValueShouldBeSent(v) {
		m["gcsBucket"] = v
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
	if v, err := dcl.EmptyValue(); err != nil {
		return nil, fmt.Errorf("error expanding CaPool into caPool: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["caPool"] = v
	}

	return m, nil
}

// flattenCertificateAuthority flattens CertificateAuthority from a JSON request object into the
// CertificateAuthority type.
func flattenCertificateAuthority(c *Client, i interface{}, res *CertificateAuthority) *CertificateAuthority {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}
	if len(m) == 0 {
		return nil
	}

	resultRes := &CertificateAuthority{}
	resultRes.Name = dcl.FlattenString(m["name"])
	resultRes.Type = flattenCertificateAuthorityTypeEnum(m["type"])
	resultRes.Config = flattenCertificateAuthorityConfig(c, m["config"], res)
	resultRes.Lifetime = dcl.FlattenString(m["lifetime"])
	resultRes.KeySpec = flattenCertificateAuthorityKeySpec(c, m["keySpec"], res)
	resultRes.SubordinateConfig = flattenCertificateAuthoritySubordinateConfig(c, m["subordinateConfig"], res)
	resultRes.Tier = flattenCertificateAuthorityTierEnum(m["tier"])
	resultRes.State = flattenCertificateAuthorityStateEnum(m["state"])
	resultRes.PemCaCertificates = dcl.FlattenStringSlice(m["pemCaCertificates"])
	resultRes.CaCertificateDescriptions = flattenCertificateAuthorityCaCertificateDescriptionsSlice(c, m["caCertificateDescriptions"], res)
	resultRes.GcsBucket = dcl.FlattenString(m["gcsBucket"])
	resultRes.AccessUrls = flattenCertificateAuthorityAccessUrls(c, m["accessUrls"], res)
	resultRes.CreateTime = dcl.FlattenString(m["createTime"])
	resultRes.UpdateTime = dcl.FlattenString(m["updateTime"])
	resultRes.DeleteTime = dcl.FlattenString(m["deleteTime"])
	resultRes.ExpireTime = dcl.FlattenString(m["expireTime"])
	resultRes.Labels = dcl.FlattenKeyValuePairs(m["labels"])
	resultRes.Project = dcl.FlattenString(m["project"])
	resultRes.Location = dcl.FlattenString(m["location"])
	resultRes.CaPool = dcl.FlattenString(m["caPool"])

	return resultRes
}

// expandCertificateAuthorityConfigMap expands the contents of CertificateAuthorityConfig into a JSON
// request object.
func expandCertificateAuthorityConfigMap(c *Client, f map[string]CertificateAuthorityConfig, res *CertificateAuthority) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandCertificateAuthorityConfig(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandCertificateAuthorityConfigSlice expands the contents of CertificateAuthorityConfig into a JSON
// request object.
func expandCertificateAuthorityConfigSlice(c *Client, f []CertificateAuthorityConfig, res *CertificateAuthority) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandCertificateAuthorityConfig(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenCertificateAuthorityConfigMap flattens the contents of CertificateAuthorityConfig from a JSON
// response object.
func flattenCertificateAuthorityConfigMap(c *Client, i interface{}, res *CertificateAuthority) map[string]CertificateAuthorityConfig {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]CertificateAuthorityConfig{}
	}

	if len(a) == 0 {
		return map[string]CertificateAuthorityConfig{}
	}

	items := make(map[string]CertificateAuthorityConfig)
	for k, item := range a {
		items[k] = *flattenCertificateAuthorityConfig(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenCertificateAuthorityConfigSlice flattens the contents of CertificateAuthorityConfig from a JSON
// response object.
func flattenCertificateAuthorityConfigSlice(c *Client, i interface{}, res *CertificateAuthority) []CertificateAuthorityConfig {
	a, ok := i.([]interface{})
	if !ok {
		return []CertificateAuthorityConfig{}
	}

	if len(a) == 0 {
		return []CertificateAuthorityConfig{}
	}

	items := make([]CertificateAuthorityConfig, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenCertificateAuthorityConfig(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandCertificateAuthorityConfig expands an instance of CertificateAuthorityConfig into a JSON
// request object.
func expandCertificateAuthorityConfig(c *Client, f *CertificateAuthorityConfig, res *CertificateAuthority) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v, err := expandCertificateAuthorityConfigSubjectConfig(c, f.SubjectConfig, res); err != nil {
		return nil, fmt.Errorf("error expanding SubjectConfig into subjectConfig: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["subjectConfig"] = v
	}
	if v, err := expandCertificateAuthorityConfigX509Config(c, f.X509Config, res); err != nil {
		return nil, fmt.Errorf("error expanding X509Config into x509Config: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["x509Config"] = v
	}

	return m, nil
}

// flattenCertificateAuthorityConfig flattens an instance of CertificateAuthorityConfig from a JSON
// response object.
func flattenCertificateAuthorityConfig(c *Client, i interface{}, res *CertificateAuthority) *CertificateAuthorityConfig {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &CertificateAuthorityConfig{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyCertificateAuthorityConfig
	}
	r.SubjectConfig = flattenCertificateAuthorityConfigSubjectConfig(c, m["subjectConfig"], res)
	r.X509Config = flattenCertificateAuthorityConfigX509Config(c, m["x509Config"], res)
	r.PublicKey = flattenCertificateAuthorityConfigPublicKey(c, m["publicKey"], res)

	return r
}

// expandCertificateAuthorityConfigSubjectConfigMap expands the contents of CertificateAuthorityConfigSubjectConfig into a JSON
// request object.
func expandCertificateAuthorityConfigSubjectConfigMap(c *Client, f map[string]CertificateAuthorityConfigSubjectConfig, res *CertificateAuthority) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandCertificateAuthorityConfigSubjectConfig(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandCertificateAuthorityConfigSubjectConfigSlice expands the contents of CertificateAuthorityConfigSubjectConfig into a JSON
// request object.
func expandCertificateAuthorityConfigSubjectConfigSlice(c *Client, f []CertificateAuthorityConfigSubjectConfig, res *CertificateAuthority) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandCertificateAuthorityConfigSubjectConfig(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenCertificateAuthorityConfigSubjectConfigMap flattens the contents of CertificateAuthorityConfigSubjectConfig from a JSON
// response object.
func flattenCertificateAuthorityConfigSubjectConfigMap(c *Client, i interface{}, res *CertificateAuthority) map[string]CertificateAuthorityConfigSubjectConfig {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]CertificateAuthorityConfigSubjectConfig{}
	}

	if len(a) == 0 {
		return map[string]CertificateAuthorityConfigSubjectConfig{}
	}

	items := make(map[string]CertificateAuthorityConfigSubjectConfig)
	for k, item := range a {
		items[k] = *flattenCertificateAuthorityConfigSubjectConfig(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenCertificateAuthorityConfigSubjectConfigSlice flattens the contents of CertificateAuthorityConfigSubjectConfig from a JSON
// response object.
func flattenCertificateAuthorityConfigSubjectConfigSlice(c *Client, i interface{}, res *CertificateAuthority) []CertificateAuthorityConfigSubjectConfig {
	a, ok := i.([]interface{})
	if !ok {
		return []CertificateAuthorityConfigSubjectConfig{}
	}

	if len(a) == 0 {
		return []CertificateAuthorityConfigSubjectConfig{}
	}

	items := make([]CertificateAuthorityConfigSubjectConfig, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenCertificateAuthorityConfigSubjectConfig(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandCertificateAuthorityConfigSubjectConfig expands an instance of CertificateAuthorityConfigSubjectConfig into a JSON
// request object.
func expandCertificateAuthorityConfigSubjectConfig(c *Client, f *CertificateAuthorityConfigSubjectConfig, res *CertificateAuthority) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v, err := expandCertificateAuthorityConfigSubjectConfigSubject(c, f.Subject, res); err != nil {
		return nil, fmt.Errorf("error expanding Subject into subject: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["subject"] = v
	}
	if v, err := expandCertificateAuthorityConfigSubjectConfigSubjectAltName(c, f.SubjectAltName, res); err != nil {
		return nil, fmt.Errorf("error expanding SubjectAltName into subjectAltName: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["subjectAltName"] = v
	}

	return m, nil
}

// flattenCertificateAuthorityConfigSubjectConfig flattens an instance of CertificateAuthorityConfigSubjectConfig from a JSON
// response object.
func flattenCertificateAuthorityConfigSubjectConfig(c *Client, i interface{}, res *CertificateAuthority) *CertificateAuthorityConfigSubjectConfig {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &CertificateAuthorityConfigSubjectConfig{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyCertificateAuthorityConfigSubjectConfig
	}
	r.Subject = flattenCertificateAuthorityConfigSubjectConfigSubject(c, m["subject"], res)
	r.SubjectAltName = flattenCertificateAuthorityConfigSubjectConfigSubjectAltName(c, m["subjectAltName"], res)

	return r
}

// expandCertificateAuthorityConfigSubjectConfigSubjectMap expands the contents of CertificateAuthorityConfigSubjectConfigSubject into a JSON
// request object.
func expandCertificateAuthorityConfigSubjectConfigSubjectMap(c *Client, f map[string]CertificateAuthorityConfigSubjectConfigSubject, res *CertificateAuthority) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandCertificateAuthorityConfigSubjectConfigSubject(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandCertificateAuthorityConfigSubjectConfigSubjectSlice expands the contents of CertificateAuthorityConfigSubjectConfigSubject into a JSON
// request object.
func expandCertificateAuthorityConfigSubjectConfigSubjectSlice(c *Client, f []CertificateAuthorityConfigSubjectConfigSubject, res *CertificateAuthority) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandCertificateAuthorityConfigSubjectConfigSubject(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenCertificateAuthorityConfigSubjectConfigSubjectMap flattens the contents of CertificateAuthorityConfigSubjectConfigSubject from a JSON
// response object.
func flattenCertificateAuthorityConfigSubjectConfigSubjectMap(c *Client, i interface{}, res *CertificateAuthority) map[string]CertificateAuthorityConfigSubjectConfigSubject {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]CertificateAuthorityConfigSubjectConfigSubject{}
	}

	if len(a) == 0 {
		return map[string]CertificateAuthorityConfigSubjectConfigSubject{}
	}

	items := make(map[string]CertificateAuthorityConfigSubjectConfigSubject)
	for k, item := range a {
		items[k] = *flattenCertificateAuthorityConfigSubjectConfigSubject(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenCertificateAuthorityConfigSubjectConfigSubjectSlice flattens the contents of CertificateAuthorityConfigSubjectConfigSubject from a JSON
// response object.
func flattenCertificateAuthorityConfigSubjectConfigSubjectSlice(c *Client, i interface{}, res *CertificateAuthority) []CertificateAuthorityConfigSubjectConfigSubject {
	a, ok := i.([]interface{})
	if !ok {
		return []CertificateAuthorityConfigSubjectConfigSubject{}
	}

	if len(a) == 0 {
		return []CertificateAuthorityConfigSubjectConfigSubject{}
	}

	items := make([]CertificateAuthorityConfigSubjectConfigSubject, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenCertificateAuthorityConfigSubjectConfigSubject(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandCertificateAuthorityConfigSubjectConfigSubject expands an instance of CertificateAuthorityConfigSubjectConfigSubject into a JSON
// request object.
func expandCertificateAuthorityConfigSubjectConfigSubject(c *Client, f *CertificateAuthorityConfigSubjectConfigSubject, res *CertificateAuthority) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.CommonName; !dcl.IsEmptyValueIndirect(v) {
		m["commonName"] = v
	}
	if v := f.CountryCode; !dcl.IsEmptyValueIndirect(v) {
		m["countryCode"] = v
	}
	if v := f.Organization; !dcl.IsEmptyValueIndirect(v) {
		m["organization"] = v
	}
	if v := f.OrganizationalUnit; !dcl.IsEmptyValueIndirect(v) {
		m["organizationalUnit"] = v
	}
	if v := f.Locality; !dcl.IsEmptyValueIndirect(v) {
		m["locality"] = v
	}
	if v := f.Province; !dcl.IsEmptyValueIndirect(v) {
		m["province"] = v
	}
	if v := f.StreetAddress; !dcl.IsEmptyValueIndirect(v) {
		m["streetAddress"] = v
	}
	if v := f.PostalCode; !dcl.IsEmptyValueIndirect(v) {
		m["postalCode"] = v
	}

	return m, nil
}

// flattenCertificateAuthorityConfigSubjectConfigSubject flattens an instance of CertificateAuthorityConfigSubjectConfigSubject from a JSON
// response object.
func flattenCertificateAuthorityConfigSubjectConfigSubject(c *Client, i interface{}, res *CertificateAuthority) *CertificateAuthorityConfigSubjectConfigSubject {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &CertificateAuthorityConfigSubjectConfigSubject{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyCertificateAuthorityConfigSubjectConfigSubject
	}
	r.CommonName = dcl.FlattenString(m["commonName"])
	r.CountryCode = dcl.FlattenString(m["countryCode"])
	r.Organization = dcl.FlattenString(m["organization"])
	r.OrganizationalUnit = dcl.FlattenString(m["organizationalUnit"])
	r.Locality = dcl.FlattenString(m["locality"])
	r.Province = dcl.FlattenString(m["province"])
	r.StreetAddress = dcl.FlattenString(m["streetAddress"])
	r.PostalCode = dcl.FlattenString(m["postalCode"])

	return r
}

// expandCertificateAuthorityConfigSubjectConfigSubjectAltNameMap expands the contents of CertificateAuthorityConfigSubjectConfigSubjectAltName into a JSON
// request object.
func expandCertificateAuthorityConfigSubjectConfigSubjectAltNameMap(c *Client, f map[string]CertificateAuthorityConfigSubjectConfigSubjectAltName, res *CertificateAuthority) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandCertificateAuthorityConfigSubjectConfigSubjectAltName(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandCertificateAuthorityConfigSubjectConfigSubjectAltNameSlice expands the contents of CertificateAuthorityConfigSubjectConfigSubjectAltName into a JSON
// request object.
func expandCertificateAuthorityConfigSubjectConfigSubjectAltNameSlice(c *Client, f []CertificateAuthorityConfigSubjectConfigSubjectAltName, res *CertificateAuthority) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandCertificateAuthorityConfigSubjectConfigSubjectAltName(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenCertificateAuthorityConfigSubjectConfigSubjectAltNameMap flattens the contents of CertificateAuthorityConfigSubjectConfigSubjectAltName from a JSON
// response object.
func flattenCertificateAuthorityConfigSubjectConfigSubjectAltNameMap(c *Client, i interface{}, res *CertificateAuthority) map[string]CertificateAuthorityConfigSubjectConfigSubjectAltName {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]CertificateAuthorityConfigSubjectConfigSubjectAltName{}
	}

	if len(a) == 0 {
		return map[string]CertificateAuthorityConfigSubjectConfigSubjectAltName{}
	}

	items := make(map[string]CertificateAuthorityConfigSubjectConfigSubjectAltName)
	for k, item := range a {
		items[k] = *flattenCertificateAuthorityConfigSubjectConfigSubjectAltName(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenCertificateAuthorityConfigSubjectConfigSubjectAltNameSlice flattens the contents of CertificateAuthorityConfigSubjectConfigSubjectAltName from a JSON
// response object.
func flattenCertificateAuthorityConfigSubjectConfigSubjectAltNameSlice(c *Client, i interface{}, res *CertificateAuthority) []CertificateAuthorityConfigSubjectConfigSubjectAltName {
	a, ok := i.([]interface{})
	if !ok {
		return []CertificateAuthorityConfigSubjectConfigSubjectAltName{}
	}

	if len(a) == 0 {
		return []CertificateAuthorityConfigSubjectConfigSubjectAltName{}
	}

	items := make([]CertificateAuthorityConfigSubjectConfigSubjectAltName, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenCertificateAuthorityConfigSubjectConfigSubjectAltName(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandCertificateAuthorityConfigSubjectConfigSubjectAltName expands an instance of CertificateAuthorityConfigSubjectConfigSubjectAltName into a JSON
// request object.
func expandCertificateAuthorityConfigSubjectConfigSubjectAltName(c *Client, f *CertificateAuthorityConfigSubjectConfigSubjectAltName, res *CertificateAuthority) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.DnsNames; v != nil {
		m["dnsNames"] = v
	}
	if v := f.Uris; v != nil {
		m["uris"] = v
	}
	if v := f.EmailAddresses; v != nil {
		m["emailAddresses"] = v
	}
	if v := f.IPAddresses; v != nil {
		m["ipAddresses"] = v
	}
	if v, err := expandCertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSansSlice(c, f.CustomSans, res); err != nil {
		return nil, fmt.Errorf("error expanding CustomSans into customSans: %w", err)
	} else if v != nil {
		m["customSans"] = v
	}

	return m, nil
}

// flattenCertificateAuthorityConfigSubjectConfigSubjectAltName flattens an instance of CertificateAuthorityConfigSubjectConfigSubjectAltName from a JSON
// response object.
func flattenCertificateAuthorityConfigSubjectConfigSubjectAltName(c *Client, i interface{}, res *CertificateAuthority) *CertificateAuthorityConfigSubjectConfigSubjectAltName {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &CertificateAuthorityConfigSubjectConfigSubjectAltName{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyCertificateAuthorityConfigSubjectConfigSubjectAltName
	}
	r.DnsNames = dcl.FlattenStringSlice(m["dnsNames"])
	r.Uris = dcl.FlattenStringSlice(m["uris"])
	r.EmailAddresses = dcl.FlattenStringSlice(m["emailAddresses"])
	r.IPAddresses = dcl.FlattenStringSlice(m["ipAddresses"])
	r.CustomSans = flattenCertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSansSlice(c, m["customSans"], res)

	return r
}

// expandCertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSansMap expands the contents of CertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSans into a JSON
// request object.
func expandCertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSansMap(c *Client, f map[string]CertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSans, res *CertificateAuthority) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandCertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSans(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandCertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSansSlice expands the contents of CertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSans into a JSON
// request object.
func expandCertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSansSlice(c *Client, f []CertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSans, res *CertificateAuthority) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandCertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSans(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenCertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSansMap flattens the contents of CertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSans from a JSON
// response object.
func flattenCertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSansMap(c *Client, i interface{}, res *CertificateAuthority) map[string]CertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSans {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]CertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSans{}
	}

	if len(a) == 0 {
		return map[string]CertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSans{}
	}

	items := make(map[string]CertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSans)
	for k, item := range a {
		items[k] = *flattenCertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSans(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenCertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSansSlice flattens the contents of CertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSans from a JSON
// response object.
func flattenCertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSansSlice(c *Client, i interface{}, res *CertificateAuthority) []CertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSans {
	a, ok := i.([]interface{})
	if !ok {
		return []CertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSans{}
	}

	if len(a) == 0 {
		return []CertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSans{}
	}

	items := make([]CertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSans, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenCertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSans(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandCertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSans expands an instance of CertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSans into a JSON
// request object.
func expandCertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSans(c *Client, f *CertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSans, res *CertificateAuthority) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v, err := expandCertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSansObjectId(c, f.ObjectId, res); err != nil {
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

// flattenCertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSans flattens an instance of CertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSans from a JSON
// response object.
func flattenCertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSans(c *Client, i interface{}, res *CertificateAuthority) *CertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSans {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &CertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSans{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyCertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSans
	}
	r.ObjectId = flattenCertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSansObjectId(c, m["objectId"], res)
	r.Critical = dcl.FlattenBool(m["critical"])
	r.Value = dcl.FlattenString(m["value"])

	return r
}

// expandCertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSansObjectIdMap expands the contents of CertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSansObjectId into a JSON
// request object.
func expandCertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSansObjectIdMap(c *Client, f map[string]CertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSansObjectId, res *CertificateAuthority) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandCertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSansObjectId(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandCertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSansObjectIdSlice expands the contents of CertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSansObjectId into a JSON
// request object.
func expandCertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSansObjectIdSlice(c *Client, f []CertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSansObjectId, res *CertificateAuthority) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandCertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSansObjectId(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenCertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSansObjectIdMap flattens the contents of CertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSansObjectId from a JSON
// response object.
func flattenCertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSansObjectIdMap(c *Client, i interface{}, res *CertificateAuthority) map[string]CertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSansObjectId {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]CertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSansObjectId{}
	}

	if len(a) == 0 {
		return map[string]CertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSansObjectId{}
	}

	items := make(map[string]CertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSansObjectId)
	for k, item := range a {
		items[k] = *flattenCertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSansObjectId(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenCertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSansObjectIdSlice flattens the contents of CertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSansObjectId from a JSON
// response object.
func flattenCertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSansObjectIdSlice(c *Client, i interface{}, res *CertificateAuthority) []CertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSansObjectId {
	a, ok := i.([]interface{})
	if !ok {
		return []CertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSansObjectId{}
	}

	if len(a) == 0 {
		return []CertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSansObjectId{}
	}

	items := make([]CertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSansObjectId, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenCertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSansObjectId(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandCertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSansObjectId expands an instance of CertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSansObjectId into a JSON
// request object.
func expandCertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSansObjectId(c *Client, f *CertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSansObjectId, res *CertificateAuthority) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.ObjectIdPath; v != nil {
		m["objectIdPath"] = v
	}

	return m, nil
}

// flattenCertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSansObjectId flattens an instance of CertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSansObjectId from a JSON
// response object.
func flattenCertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSansObjectId(c *Client, i interface{}, res *CertificateAuthority) *CertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSansObjectId {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &CertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSansObjectId{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyCertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSansObjectId
	}
	r.ObjectIdPath = dcl.FlattenIntSlice(m["objectIdPath"])

	return r
}

// expandCertificateAuthorityConfigX509ConfigMap expands the contents of CertificateAuthorityConfigX509Config into a JSON
// request object.
func expandCertificateAuthorityConfigX509ConfigMap(c *Client, f map[string]CertificateAuthorityConfigX509Config, res *CertificateAuthority) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandCertificateAuthorityConfigX509Config(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandCertificateAuthorityConfigX509ConfigSlice expands the contents of CertificateAuthorityConfigX509Config into a JSON
// request object.
func expandCertificateAuthorityConfigX509ConfigSlice(c *Client, f []CertificateAuthorityConfigX509Config, res *CertificateAuthority) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandCertificateAuthorityConfigX509Config(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenCertificateAuthorityConfigX509ConfigMap flattens the contents of CertificateAuthorityConfigX509Config from a JSON
// response object.
func flattenCertificateAuthorityConfigX509ConfigMap(c *Client, i interface{}, res *CertificateAuthority) map[string]CertificateAuthorityConfigX509Config {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]CertificateAuthorityConfigX509Config{}
	}

	if len(a) == 0 {
		return map[string]CertificateAuthorityConfigX509Config{}
	}

	items := make(map[string]CertificateAuthorityConfigX509Config)
	for k, item := range a {
		items[k] = *flattenCertificateAuthorityConfigX509Config(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenCertificateAuthorityConfigX509ConfigSlice flattens the contents of CertificateAuthorityConfigX509Config from a JSON
// response object.
func flattenCertificateAuthorityConfigX509ConfigSlice(c *Client, i interface{}, res *CertificateAuthority) []CertificateAuthorityConfigX509Config {
	a, ok := i.([]interface{})
	if !ok {
		return []CertificateAuthorityConfigX509Config{}
	}

	if len(a) == 0 {
		return []CertificateAuthorityConfigX509Config{}
	}

	items := make([]CertificateAuthorityConfigX509Config, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenCertificateAuthorityConfigX509Config(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandCertificateAuthorityConfigX509Config expands an instance of CertificateAuthorityConfigX509Config into a JSON
// request object.
func expandCertificateAuthorityConfigX509Config(c *Client, f *CertificateAuthorityConfigX509Config, res *CertificateAuthority) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v, err := expandCertificateAuthorityConfigX509ConfigKeyUsage(c, f.KeyUsage, res); err != nil {
		return nil, fmt.Errorf("error expanding KeyUsage into keyUsage: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["keyUsage"] = v
	}
	if v, err := expandCertificateAuthorityConfigX509ConfigCAOptions(c, f.CaOptions, res); err != nil {
		return nil, fmt.Errorf("error expanding CaOptions into caOptions: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["caOptions"] = v
	}
	if v, err := expandCertificateAuthorityConfigX509ConfigPolicyIdsSlice(c, f.PolicyIds, res); err != nil {
		return nil, fmt.Errorf("error expanding PolicyIds into policyIds: %w", err)
	} else if v != nil {
		m["policyIds"] = v
	}
	if v, err := expandCertificateAuthorityConfigX509ConfigAdditionalExtensionsSlice(c, f.AdditionalExtensions, res); err != nil {
		return nil, fmt.Errorf("error expanding AdditionalExtensions into additionalExtensions: %w", err)
	} else if v != nil {
		m["additionalExtensions"] = v
	}

	return m, nil
}

// flattenCertificateAuthorityConfigX509Config flattens an instance of CertificateAuthorityConfigX509Config from a JSON
// response object.
func flattenCertificateAuthorityConfigX509Config(c *Client, i interface{}, res *CertificateAuthority) *CertificateAuthorityConfigX509Config {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &CertificateAuthorityConfigX509Config{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyCertificateAuthorityConfigX509Config
	}
	r.KeyUsage = flattenCertificateAuthorityConfigX509ConfigKeyUsage(c, m["keyUsage"], res)
	r.CaOptions = flattenCertificateAuthorityConfigX509ConfigCAOptions(c, m["caOptions"], res)
	r.PolicyIds = flattenCertificateAuthorityConfigX509ConfigPolicyIdsSlice(c, m["policyIds"], res)
	r.AiaOcspServers = dcl.FlattenStringSlice(m["aiaOcspServers"])
	r.AdditionalExtensions = flattenCertificateAuthorityConfigX509ConfigAdditionalExtensionsSlice(c, m["additionalExtensions"], res)

	return r
}

// expandCertificateAuthorityConfigX509ConfigKeyUsageMap expands the contents of CertificateAuthorityConfigX509ConfigKeyUsage into a JSON
// request object.
func expandCertificateAuthorityConfigX509ConfigKeyUsageMap(c *Client, f map[string]CertificateAuthorityConfigX509ConfigKeyUsage, res *CertificateAuthority) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandCertificateAuthorityConfigX509ConfigKeyUsage(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandCertificateAuthorityConfigX509ConfigKeyUsageSlice expands the contents of CertificateAuthorityConfigX509ConfigKeyUsage into a JSON
// request object.
func expandCertificateAuthorityConfigX509ConfigKeyUsageSlice(c *Client, f []CertificateAuthorityConfigX509ConfigKeyUsage, res *CertificateAuthority) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandCertificateAuthorityConfigX509ConfigKeyUsage(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenCertificateAuthorityConfigX509ConfigKeyUsageMap flattens the contents of CertificateAuthorityConfigX509ConfigKeyUsage from a JSON
// response object.
func flattenCertificateAuthorityConfigX509ConfigKeyUsageMap(c *Client, i interface{}, res *CertificateAuthority) map[string]CertificateAuthorityConfigX509ConfigKeyUsage {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]CertificateAuthorityConfigX509ConfigKeyUsage{}
	}

	if len(a) == 0 {
		return map[string]CertificateAuthorityConfigX509ConfigKeyUsage{}
	}

	items := make(map[string]CertificateAuthorityConfigX509ConfigKeyUsage)
	for k, item := range a {
		items[k] = *flattenCertificateAuthorityConfigX509ConfigKeyUsage(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenCertificateAuthorityConfigX509ConfigKeyUsageSlice flattens the contents of CertificateAuthorityConfigX509ConfigKeyUsage from a JSON
// response object.
func flattenCertificateAuthorityConfigX509ConfigKeyUsageSlice(c *Client, i interface{}, res *CertificateAuthority) []CertificateAuthorityConfigX509ConfigKeyUsage {
	a, ok := i.([]interface{})
	if !ok {
		return []CertificateAuthorityConfigX509ConfigKeyUsage{}
	}

	if len(a) == 0 {
		return []CertificateAuthorityConfigX509ConfigKeyUsage{}
	}

	items := make([]CertificateAuthorityConfigX509ConfigKeyUsage, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenCertificateAuthorityConfigX509ConfigKeyUsage(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandCertificateAuthorityConfigX509ConfigKeyUsage expands an instance of CertificateAuthorityConfigX509ConfigKeyUsage into a JSON
// request object.
func expandCertificateAuthorityConfigX509ConfigKeyUsage(c *Client, f *CertificateAuthorityConfigX509ConfigKeyUsage, res *CertificateAuthority) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v, err := expandCertificateAuthorityConfigX509ConfigKeyUsageBaseKeyUsage(c, f.BaseKeyUsage, res); err != nil {
		return nil, fmt.Errorf("error expanding BaseKeyUsage into baseKeyUsage: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["baseKeyUsage"] = v
	}
	if v, err := expandCertificateAuthorityConfigX509ConfigKeyUsageExtendedKeyUsage(c, f.ExtendedKeyUsage, res); err != nil {
		return nil, fmt.Errorf("error expanding ExtendedKeyUsage into extendedKeyUsage: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["extendedKeyUsage"] = v
	}
	if v, err := expandCertificateAuthorityConfigX509ConfigKeyUsageUnknownExtendedKeyUsagesSlice(c, f.UnknownExtendedKeyUsages, res); err != nil {
		return nil, fmt.Errorf("error expanding UnknownExtendedKeyUsages into unknownExtendedKeyUsages: %w", err)
	} else if v != nil {
		m["unknownExtendedKeyUsages"] = v
	}

	return m, nil
}

// flattenCertificateAuthorityConfigX509ConfigKeyUsage flattens an instance of CertificateAuthorityConfigX509ConfigKeyUsage from a JSON
// response object.
func flattenCertificateAuthorityConfigX509ConfigKeyUsage(c *Client, i interface{}, res *CertificateAuthority) *CertificateAuthorityConfigX509ConfigKeyUsage {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &CertificateAuthorityConfigX509ConfigKeyUsage{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyCertificateAuthorityConfigX509ConfigKeyUsage
	}
	r.BaseKeyUsage = flattenCertificateAuthorityConfigX509ConfigKeyUsageBaseKeyUsage(c, m["baseKeyUsage"], res)
	r.ExtendedKeyUsage = flattenCertificateAuthorityConfigX509ConfigKeyUsageExtendedKeyUsage(c, m["extendedKeyUsage"], res)
	r.UnknownExtendedKeyUsages = flattenCertificateAuthorityConfigX509ConfigKeyUsageUnknownExtendedKeyUsagesSlice(c, m["unknownExtendedKeyUsages"], res)

	return r
}

// expandCertificateAuthorityConfigX509ConfigKeyUsageBaseKeyUsageMap expands the contents of CertificateAuthorityConfigX509ConfigKeyUsageBaseKeyUsage into a JSON
// request object.
func expandCertificateAuthorityConfigX509ConfigKeyUsageBaseKeyUsageMap(c *Client, f map[string]CertificateAuthorityConfigX509ConfigKeyUsageBaseKeyUsage, res *CertificateAuthority) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandCertificateAuthorityConfigX509ConfigKeyUsageBaseKeyUsage(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandCertificateAuthorityConfigX509ConfigKeyUsageBaseKeyUsageSlice expands the contents of CertificateAuthorityConfigX509ConfigKeyUsageBaseKeyUsage into a JSON
// request object.
func expandCertificateAuthorityConfigX509ConfigKeyUsageBaseKeyUsageSlice(c *Client, f []CertificateAuthorityConfigX509ConfigKeyUsageBaseKeyUsage, res *CertificateAuthority) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandCertificateAuthorityConfigX509ConfigKeyUsageBaseKeyUsage(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenCertificateAuthorityConfigX509ConfigKeyUsageBaseKeyUsageMap flattens the contents of CertificateAuthorityConfigX509ConfigKeyUsageBaseKeyUsage from a JSON
// response object.
func flattenCertificateAuthorityConfigX509ConfigKeyUsageBaseKeyUsageMap(c *Client, i interface{}, res *CertificateAuthority) map[string]CertificateAuthorityConfigX509ConfigKeyUsageBaseKeyUsage {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]CertificateAuthorityConfigX509ConfigKeyUsageBaseKeyUsage{}
	}

	if len(a) == 0 {
		return map[string]CertificateAuthorityConfigX509ConfigKeyUsageBaseKeyUsage{}
	}

	items := make(map[string]CertificateAuthorityConfigX509ConfigKeyUsageBaseKeyUsage)
	for k, item := range a {
		items[k] = *flattenCertificateAuthorityConfigX509ConfigKeyUsageBaseKeyUsage(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenCertificateAuthorityConfigX509ConfigKeyUsageBaseKeyUsageSlice flattens the contents of CertificateAuthorityConfigX509ConfigKeyUsageBaseKeyUsage from a JSON
// response object.
func flattenCertificateAuthorityConfigX509ConfigKeyUsageBaseKeyUsageSlice(c *Client, i interface{}, res *CertificateAuthority) []CertificateAuthorityConfigX509ConfigKeyUsageBaseKeyUsage {
	a, ok := i.([]interface{})
	if !ok {
		return []CertificateAuthorityConfigX509ConfigKeyUsageBaseKeyUsage{}
	}

	if len(a) == 0 {
		return []CertificateAuthorityConfigX509ConfigKeyUsageBaseKeyUsage{}
	}

	items := make([]CertificateAuthorityConfigX509ConfigKeyUsageBaseKeyUsage, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenCertificateAuthorityConfigX509ConfigKeyUsageBaseKeyUsage(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandCertificateAuthorityConfigX509ConfigKeyUsageBaseKeyUsage expands an instance of CertificateAuthorityConfigX509ConfigKeyUsageBaseKeyUsage into a JSON
// request object.
func expandCertificateAuthorityConfigX509ConfigKeyUsageBaseKeyUsage(c *Client, f *CertificateAuthorityConfigX509ConfigKeyUsageBaseKeyUsage, res *CertificateAuthority) (map[string]interface{}, error) {
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

// flattenCertificateAuthorityConfigX509ConfigKeyUsageBaseKeyUsage flattens an instance of CertificateAuthorityConfigX509ConfigKeyUsageBaseKeyUsage from a JSON
// response object.
func flattenCertificateAuthorityConfigX509ConfigKeyUsageBaseKeyUsage(c *Client, i interface{}, res *CertificateAuthority) *CertificateAuthorityConfigX509ConfigKeyUsageBaseKeyUsage {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &CertificateAuthorityConfigX509ConfigKeyUsageBaseKeyUsage{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyCertificateAuthorityConfigX509ConfigKeyUsageBaseKeyUsage
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

// expandCertificateAuthorityConfigX509ConfigKeyUsageExtendedKeyUsageMap expands the contents of CertificateAuthorityConfigX509ConfigKeyUsageExtendedKeyUsage into a JSON
// request object.
func expandCertificateAuthorityConfigX509ConfigKeyUsageExtendedKeyUsageMap(c *Client, f map[string]CertificateAuthorityConfigX509ConfigKeyUsageExtendedKeyUsage, res *CertificateAuthority) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandCertificateAuthorityConfigX509ConfigKeyUsageExtendedKeyUsage(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandCertificateAuthorityConfigX509ConfigKeyUsageExtendedKeyUsageSlice expands the contents of CertificateAuthorityConfigX509ConfigKeyUsageExtendedKeyUsage into a JSON
// request object.
func expandCertificateAuthorityConfigX509ConfigKeyUsageExtendedKeyUsageSlice(c *Client, f []CertificateAuthorityConfigX509ConfigKeyUsageExtendedKeyUsage, res *CertificateAuthority) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandCertificateAuthorityConfigX509ConfigKeyUsageExtendedKeyUsage(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenCertificateAuthorityConfigX509ConfigKeyUsageExtendedKeyUsageMap flattens the contents of CertificateAuthorityConfigX509ConfigKeyUsageExtendedKeyUsage from a JSON
// response object.
func flattenCertificateAuthorityConfigX509ConfigKeyUsageExtendedKeyUsageMap(c *Client, i interface{}, res *CertificateAuthority) map[string]CertificateAuthorityConfigX509ConfigKeyUsageExtendedKeyUsage {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]CertificateAuthorityConfigX509ConfigKeyUsageExtendedKeyUsage{}
	}

	if len(a) == 0 {
		return map[string]CertificateAuthorityConfigX509ConfigKeyUsageExtendedKeyUsage{}
	}

	items := make(map[string]CertificateAuthorityConfigX509ConfigKeyUsageExtendedKeyUsage)
	for k, item := range a {
		items[k] = *flattenCertificateAuthorityConfigX509ConfigKeyUsageExtendedKeyUsage(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenCertificateAuthorityConfigX509ConfigKeyUsageExtendedKeyUsageSlice flattens the contents of CertificateAuthorityConfigX509ConfigKeyUsageExtendedKeyUsage from a JSON
// response object.
func flattenCertificateAuthorityConfigX509ConfigKeyUsageExtendedKeyUsageSlice(c *Client, i interface{}, res *CertificateAuthority) []CertificateAuthorityConfigX509ConfigKeyUsageExtendedKeyUsage {
	a, ok := i.([]interface{})
	if !ok {
		return []CertificateAuthorityConfigX509ConfigKeyUsageExtendedKeyUsage{}
	}

	if len(a) == 0 {
		return []CertificateAuthorityConfigX509ConfigKeyUsageExtendedKeyUsage{}
	}

	items := make([]CertificateAuthorityConfigX509ConfigKeyUsageExtendedKeyUsage, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenCertificateAuthorityConfigX509ConfigKeyUsageExtendedKeyUsage(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandCertificateAuthorityConfigX509ConfigKeyUsageExtendedKeyUsage expands an instance of CertificateAuthorityConfigX509ConfigKeyUsageExtendedKeyUsage into a JSON
// request object.
func expandCertificateAuthorityConfigX509ConfigKeyUsageExtendedKeyUsage(c *Client, f *CertificateAuthorityConfigX509ConfigKeyUsageExtendedKeyUsage, res *CertificateAuthority) (map[string]interface{}, error) {
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

// flattenCertificateAuthorityConfigX509ConfigKeyUsageExtendedKeyUsage flattens an instance of CertificateAuthorityConfigX509ConfigKeyUsageExtendedKeyUsage from a JSON
// response object.
func flattenCertificateAuthorityConfigX509ConfigKeyUsageExtendedKeyUsage(c *Client, i interface{}, res *CertificateAuthority) *CertificateAuthorityConfigX509ConfigKeyUsageExtendedKeyUsage {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &CertificateAuthorityConfigX509ConfigKeyUsageExtendedKeyUsage{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyCertificateAuthorityConfigX509ConfigKeyUsageExtendedKeyUsage
	}
	r.ServerAuth = dcl.FlattenBool(m["serverAuth"])
	r.ClientAuth = dcl.FlattenBool(m["clientAuth"])
	r.CodeSigning = dcl.FlattenBool(m["codeSigning"])
	r.EmailProtection = dcl.FlattenBool(m["emailProtection"])
	r.TimeStamping = dcl.FlattenBool(m["timeStamping"])
	r.OcspSigning = dcl.FlattenBool(m["ocspSigning"])

	return r
}

// expandCertificateAuthorityConfigX509ConfigKeyUsageUnknownExtendedKeyUsagesMap expands the contents of CertificateAuthorityConfigX509ConfigKeyUsageUnknownExtendedKeyUsages into a JSON
// request object.
func expandCertificateAuthorityConfigX509ConfigKeyUsageUnknownExtendedKeyUsagesMap(c *Client, f map[string]CertificateAuthorityConfigX509ConfigKeyUsageUnknownExtendedKeyUsages, res *CertificateAuthority) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandCertificateAuthorityConfigX509ConfigKeyUsageUnknownExtendedKeyUsages(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandCertificateAuthorityConfigX509ConfigKeyUsageUnknownExtendedKeyUsagesSlice expands the contents of CertificateAuthorityConfigX509ConfigKeyUsageUnknownExtendedKeyUsages into a JSON
// request object.
func expandCertificateAuthorityConfigX509ConfigKeyUsageUnknownExtendedKeyUsagesSlice(c *Client, f []CertificateAuthorityConfigX509ConfigKeyUsageUnknownExtendedKeyUsages, res *CertificateAuthority) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandCertificateAuthorityConfigX509ConfigKeyUsageUnknownExtendedKeyUsages(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenCertificateAuthorityConfigX509ConfigKeyUsageUnknownExtendedKeyUsagesMap flattens the contents of CertificateAuthorityConfigX509ConfigKeyUsageUnknownExtendedKeyUsages from a JSON
// response object.
func flattenCertificateAuthorityConfigX509ConfigKeyUsageUnknownExtendedKeyUsagesMap(c *Client, i interface{}, res *CertificateAuthority) map[string]CertificateAuthorityConfigX509ConfigKeyUsageUnknownExtendedKeyUsages {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]CertificateAuthorityConfigX509ConfigKeyUsageUnknownExtendedKeyUsages{}
	}

	if len(a) == 0 {
		return map[string]CertificateAuthorityConfigX509ConfigKeyUsageUnknownExtendedKeyUsages{}
	}

	items := make(map[string]CertificateAuthorityConfigX509ConfigKeyUsageUnknownExtendedKeyUsages)
	for k, item := range a {
		items[k] = *flattenCertificateAuthorityConfigX509ConfigKeyUsageUnknownExtendedKeyUsages(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenCertificateAuthorityConfigX509ConfigKeyUsageUnknownExtendedKeyUsagesSlice flattens the contents of CertificateAuthorityConfigX509ConfigKeyUsageUnknownExtendedKeyUsages from a JSON
// response object.
func flattenCertificateAuthorityConfigX509ConfigKeyUsageUnknownExtendedKeyUsagesSlice(c *Client, i interface{}, res *CertificateAuthority) []CertificateAuthorityConfigX509ConfigKeyUsageUnknownExtendedKeyUsages {
	a, ok := i.([]interface{})
	if !ok {
		return []CertificateAuthorityConfigX509ConfigKeyUsageUnknownExtendedKeyUsages{}
	}

	if len(a) == 0 {
		return []CertificateAuthorityConfigX509ConfigKeyUsageUnknownExtendedKeyUsages{}
	}

	items := make([]CertificateAuthorityConfigX509ConfigKeyUsageUnknownExtendedKeyUsages, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenCertificateAuthorityConfigX509ConfigKeyUsageUnknownExtendedKeyUsages(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandCertificateAuthorityConfigX509ConfigKeyUsageUnknownExtendedKeyUsages expands an instance of CertificateAuthorityConfigX509ConfigKeyUsageUnknownExtendedKeyUsages into a JSON
// request object.
func expandCertificateAuthorityConfigX509ConfigKeyUsageUnknownExtendedKeyUsages(c *Client, f *CertificateAuthorityConfigX509ConfigKeyUsageUnknownExtendedKeyUsages, res *CertificateAuthority) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.ObjectIdPath; v != nil {
		m["objectIdPath"] = v
	}

	return m, nil
}

// flattenCertificateAuthorityConfigX509ConfigKeyUsageUnknownExtendedKeyUsages flattens an instance of CertificateAuthorityConfigX509ConfigKeyUsageUnknownExtendedKeyUsages from a JSON
// response object.
func flattenCertificateAuthorityConfigX509ConfigKeyUsageUnknownExtendedKeyUsages(c *Client, i interface{}, res *CertificateAuthority) *CertificateAuthorityConfigX509ConfigKeyUsageUnknownExtendedKeyUsages {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &CertificateAuthorityConfigX509ConfigKeyUsageUnknownExtendedKeyUsages{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyCertificateAuthorityConfigX509ConfigKeyUsageUnknownExtendedKeyUsages
	}
	r.ObjectIdPath = dcl.FlattenIntSlice(m["objectIdPath"])

	return r
}

// expandCertificateAuthorityConfigX509ConfigCaOptionsMap expands the contents of CertificateAuthorityConfigX509ConfigCaOptions into a JSON
// request object.
func expandCertificateAuthorityConfigX509ConfigCaOptionsMap(c *Client, f map[string]CertificateAuthorityConfigX509ConfigCaOptions, res *CertificateAuthority) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandCertificateAuthorityConfigX509ConfigCaOptions(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandCertificateAuthorityConfigX509ConfigCaOptionsSlice expands the contents of CertificateAuthorityConfigX509ConfigCaOptions into a JSON
// request object.
func expandCertificateAuthorityConfigX509ConfigCaOptionsSlice(c *Client, f []CertificateAuthorityConfigX509ConfigCaOptions, res *CertificateAuthority) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandCertificateAuthorityConfigX509ConfigCaOptions(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenCertificateAuthorityConfigX509ConfigCaOptionsMap flattens the contents of CertificateAuthorityConfigX509ConfigCaOptions from a JSON
// response object.
func flattenCertificateAuthorityConfigX509ConfigCaOptionsMap(c *Client, i interface{}, res *CertificateAuthority) map[string]CertificateAuthorityConfigX509ConfigCaOptions {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]CertificateAuthorityConfigX509ConfigCaOptions{}
	}

	if len(a) == 0 {
		return map[string]CertificateAuthorityConfigX509ConfigCaOptions{}
	}

	items := make(map[string]CertificateAuthorityConfigX509ConfigCaOptions)
	for k, item := range a {
		items[k] = *flattenCertificateAuthorityConfigX509ConfigCaOptions(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenCertificateAuthorityConfigX509ConfigCaOptionsSlice flattens the contents of CertificateAuthorityConfigX509ConfigCaOptions from a JSON
// response object.
func flattenCertificateAuthorityConfigX509ConfigCaOptionsSlice(c *Client, i interface{}, res *CertificateAuthority) []CertificateAuthorityConfigX509ConfigCaOptions {
	a, ok := i.([]interface{})
	if !ok {
		return []CertificateAuthorityConfigX509ConfigCaOptions{}
	}

	if len(a) == 0 {
		return []CertificateAuthorityConfigX509ConfigCaOptions{}
	}

	items := make([]CertificateAuthorityConfigX509ConfigCaOptions, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenCertificateAuthorityConfigX509ConfigCaOptions(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandCertificateAuthorityConfigX509ConfigCaOptions expands an instance of CertificateAuthorityConfigX509ConfigCaOptions into a JSON
// request object.
func expandCertificateAuthorityConfigX509ConfigCaOptions(c *Client, f *CertificateAuthorityConfigX509ConfigCaOptions, res *CertificateAuthority) (map[string]interface{}, error) {
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

// flattenCertificateAuthorityConfigX509ConfigCaOptions flattens an instance of CertificateAuthorityConfigX509ConfigCaOptions from a JSON
// response object.
func flattenCertificateAuthorityConfigX509ConfigCaOptions(c *Client, i interface{}, res *CertificateAuthority) *CertificateAuthorityConfigX509ConfigCaOptions {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &CertificateAuthorityConfigX509ConfigCaOptions{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyCertificateAuthorityConfigX509ConfigCaOptions
	}
	r.IsCa = dcl.FlattenBool(m["isCa"])
	r.MaxIssuerPathLength = dcl.FlattenInteger(m["maxIssuerPathLength"])
	r.ZeroMaxIssuerPathLength = dcl.FlattenBool(m["zeroMaxIssuerPathLength"])

	return r
}

// expandCertificateAuthorityConfigX509ConfigPolicyIdsMap expands the contents of CertificateAuthorityConfigX509ConfigPolicyIds into a JSON
// request object.
func expandCertificateAuthorityConfigX509ConfigPolicyIdsMap(c *Client, f map[string]CertificateAuthorityConfigX509ConfigPolicyIds, res *CertificateAuthority) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandCertificateAuthorityConfigX509ConfigPolicyIds(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandCertificateAuthorityConfigX509ConfigPolicyIdsSlice expands the contents of CertificateAuthorityConfigX509ConfigPolicyIds into a JSON
// request object.
func expandCertificateAuthorityConfigX509ConfigPolicyIdsSlice(c *Client, f []CertificateAuthorityConfigX509ConfigPolicyIds, res *CertificateAuthority) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandCertificateAuthorityConfigX509ConfigPolicyIds(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenCertificateAuthorityConfigX509ConfigPolicyIdsMap flattens the contents of CertificateAuthorityConfigX509ConfigPolicyIds from a JSON
// response object.
func flattenCertificateAuthorityConfigX509ConfigPolicyIdsMap(c *Client, i interface{}, res *CertificateAuthority) map[string]CertificateAuthorityConfigX509ConfigPolicyIds {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]CertificateAuthorityConfigX509ConfigPolicyIds{}
	}

	if len(a) == 0 {
		return map[string]CertificateAuthorityConfigX509ConfigPolicyIds{}
	}

	items := make(map[string]CertificateAuthorityConfigX509ConfigPolicyIds)
	for k, item := range a {
		items[k] = *flattenCertificateAuthorityConfigX509ConfigPolicyIds(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenCertificateAuthorityConfigX509ConfigPolicyIdsSlice flattens the contents of CertificateAuthorityConfigX509ConfigPolicyIds from a JSON
// response object.
func flattenCertificateAuthorityConfigX509ConfigPolicyIdsSlice(c *Client, i interface{}, res *CertificateAuthority) []CertificateAuthorityConfigX509ConfigPolicyIds {
	a, ok := i.([]interface{})
	if !ok {
		return []CertificateAuthorityConfigX509ConfigPolicyIds{}
	}

	if len(a) == 0 {
		return []CertificateAuthorityConfigX509ConfigPolicyIds{}
	}

	items := make([]CertificateAuthorityConfigX509ConfigPolicyIds, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenCertificateAuthorityConfigX509ConfigPolicyIds(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandCertificateAuthorityConfigX509ConfigPolicyIds expands an instance of CertificateAuthorityConfigX509ConfigPolicyIds into a JSON
// request object.
func expandCertificateAuthorityConfigX509ConfigPolicyIds(c *Client, f *CertificateAuthorityConfigX509ConfigPolicyIds, res *CertificateAuthority) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.ObjectIdPath; v != nil {
		m["objectIdPath"] = v
	}

	return m, nil
}

// flattenCertificateAuthorityConfigX509ConfigPolicyIds flattens an instance of CertificateAuthorityConfigX509ConfigPolicyIds from a JSON
// response object.
func flattenCertificateAuthorityConfigX509ConfigPolicyIds(c *Client, i interface{}, res *CertificateAuthority) *CertificateAuthorityConfigX509ConfigPolicyIds {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &CertificateAuthorityConfigX509ConfigPolicyIds{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyCertificateAuthorityConfigX509ConfigPolicyIds
	}
	r.ObjectIdPath = dcl.FlattenIntSlice(m["objectIdPath"])

	return r
}

// expandCertificateAuthorityConfigX509ConfigAdditionalExtensionsMap expands the contents of CertificateAuthorityConfigX509ConfigAdditionalExtensions into a JSON
// request object.
func expandCertificateAuthorityConfigX509ConfigAdditionalExtensionsMap(c *Client, f map[string]CertificateAuthorityConfigX509ConfigAdditionalExtensions, res *CertificateAuthority) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandCertificateAuthorityConfigX509ConfigAdditionalExtensions(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandCertificateAuthorityConfigX509ConfigAdditionalExtensionsSlice expands the contents of CertificateAuthorityConfigX509ConfigAdditionalExtensions into a JSON
// request object.
func expandCertificateAuthorityConfigX509ConfigAdditionalExtensionsSlice(c *Client, f []CertificateAuthorityConfigX509ConfigAdditionalExtensions, res *CertificateAuthority) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandCertificateAuthorityConfigX509ConfigAdditionalExtensions(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenCertificateAuthorityConfigX509ConfigAdditionalExtensionsMap flattens the contents of CertificateAuthorityConfigX509ConfigAdditionalExtensions from a JSON
// response object.
func flattenCertificateAuthorityConfigX509ConfigAdditionalExtensionsMap(c *Client, i interface{}, res *CertificateAuthority) map[string]CertificateAuthorityConfigX509ConfigAdditionalExtensions {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]CertificateAuthorityConfigX509ConfigAdditionalExtensions{}
	}

	if len(a) == 0 {
		return map[string]CertificateAuthorityConfigX509ConfigAdditionalExtensions{}
	}

	items := make(map[string]CertificateAuthorityConfigX509ConfigAdditionalExtensions)
	for k, item := range a {
		items[k] = *flattenCertificateAuthorityConfigX509ConfigAdditionalExtensions(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenCertificateAuthorityConfigX509ConfigAdditionalExtensionsSlice flattens the contents of CertificateAuthorityConfigX509ConfigAdditionalExtensions from a JSON
// response object.
func flattenCertificateAuthorityConfigX509ConfigAdditionalExtensionsSlice(c *Client, i interface{}, res *CertificateAuthority) []CertificateAuthorityConfigX509ConfigAdditionalExtensions {
	a, ok := i.([]interface{})
	if !ok {
		return []CertificateAuthorityConfigX509ConfigAdditionalExtensions{}
	}

	if len(a) == 0 {
		return []CertificateAuthorityConfigX509ConfigAdditionalExtensions{}
	}

	items := make([]CertificateAuthorityConfigX509ConfigAdditionalExtensions, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenCertificateAuthorityConfigX509ConfigAdditionalExtensions(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandCertificateAuthorityConfigX509ConfigAdditionalExtensions expands an instance of CertificateAuthorityConfigX509ConfigAdditionalExtensions into a JSON
// request object.
func expandCertificateAuthorityConfigX509ConfigAdditionalExtensions(c *Client, f *CertificateAuthorityConfigX509ConfigAdditionalExtensions, res *CertificateAuthority) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v, err := expandCertificateAuthorityConfigX509ConfigAdditionalExtensionsObjectId(c, f.ObjectId, res); err != nil {
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

// flattenCertificateAuthorityConfigX509ConfigAdditionalExtensions flattens an instance of CertificateAuthorityConfigX509ConfigAdditionalExtensions from a JSON
// response object.
func flattenCertificateAuthorityConfigX509ConfigAdditionalExtensions(c *Client, i interface{}, res *CertificateAuthority) *CertificateAuthorityConfigX509ConfigAdditionalExtensions {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &CertificateAuthorityConfigX509ConfigAdditionalExtensions{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyCertificateAuthorityConfigX509ConfigAdditionalExtensions
	}
	r.ObjectId = flattenCertificateAuthorityConfigX509ConfigAdditionalExtensionsObjectId(c, m["objectId"], res)
	r.Critical = dcl.FlattenBool(m["critical"])
	r.Value = dcl.FlattenString(m["value"])

	return r
}

// expandCertificateAuthorityConfigX509ConfigAdditionalExtensionsObjectIdMap expands the contents of CertificateAuthorityConfigX509ConfigAdditionalExtensionsObjectId into a JSON
// request object.
func expandCertificateAuthorityConfigX509ConfigAdditionalExtensionsObjectIdMap(c *Client, f map[string]CertificateAuthorityConfigX509ConfigAdditionalExtensionsObjectId, res *CertificateAuthority) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandCertificateAuthorityConfigX509ConfigAdditionalExtensionsObjectId(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandCertificateAuthorityConfigX509ConfigAdditionalExtensionsObjectIdSlice expands the contents of CertificateAuthorityConfigX509ConfigAdditionalExtensionsObjectId into a JSON
// request object.
func expandCertificateAuthorityConfigX509ConfigAdditionalExtensionsObjectIdSlice(c *Client, f []CertificateAuthorityConfigX509ConfigAdditionalExtensionsObjectId, res *CertificateAuthority) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandCertificateAuthorityConfigX509ConfigAdditionalExtensionsObjectId(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenCertificateAuthorityConfigX509ConfigAdditionalExtensionsObjectIdMap flattens the contents of CertificateAuthorityConfigX509ConfigAdditionalExtensionsObjectId from a JSON
// response object.
func flattenCertificateAuthorityConfigX509ConfigAdditionalExtensionsObjectIdMap(c *Client, i interface{}, res *CertificateAuthority) map[string]CertificateAuthorityConfigX509ConfigAdditionalExtensionsObjectId {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]CertificateAuthorityConfigX509ConfigAdditionalExtensionsObjectId{}
	}

	if len(a) == 0 {
		return map[string]CertificateAuthorityConfigX509ConfigAdditionalExtensionsObjectId{}
	}

	items := make(map[string]CertificateAuthorityConfigX509ConfigAdditionalExtensionsObjectId)
	for k, item := range a {
		items[k] = *flattenCertificateAuthorityConfigX509ConfigAdditionalExtensionsObjectId(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenCertificateAuthorityConfigX509ConfigAdditionalExtensionsObjectIdSlice flattens the contents of CertificateAuthorityConfigX509ConfigAdditionalExtensionsObjectId from a JSON
// response object.
func flattenCertificateAuthorityConfigX509ConfigAdditionalExtensionsObjectIdSlice(c *Client, i interface{}, res *CertificateAuthority) []CertificateAuthorityConfigX509ConfigAdditionalExtensionsObjectId {
	a, ok := i.([]interface{})
	if !ok {
		return []CertificateAuthorityConfigX509ConfigAdditionalExtensionsObjectId{}
	}

	if len(a) == 0 {
		return []CertificateAuthorityConfigX509ConfigAdditionalExtensionsObjectId{}
	}

	items := make([]CertificateAuthorityConfigX509ConfigAdditionalExtensionsObjectId, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenCertificateAuthorityConfigX509ConfigAdditionalExtensionsObjectId(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandCertificateAuthorityConfigX509ConfigAdditionalExtensionsObjectId expands an instance of CertificateAuthorityConfigX509ConfigAdditionalExtensionsObjectId into a JSON
// request object.
func expandCertificateAuthorityConfigX509ConfigAdditionalExtensionsObjectId(c *Client, f *CertificateAuthorityConfigX509ConfigAdditionalExtensionsObjectId, res *CertificateAuthority) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.ObjectIdPath; v != nil {
		m["objectIdPath"] = v
	}

	return m, nil
}

// flattenCertificateAuthorityConfigX509ConfigAdditionalExtensionsObjectId flattens an instance of CertificateAuthorityConfigX509ConfigAdditionalExtensionsObjectId from a JSON
// response object.
func flattenCertificateAuthorityConfigX509ConfigAdditionalExtensionsObjectId(c *Client, i interface{}, res *CertificateAuthority) *CertificateAuthorityConfigX509ConfigAdditionalExtensionsObjectId {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &CertificateAuthorityConfigX509ConfigAdditionalExtensionsObjectId{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyCertificateAuthorityConfigX509ConfigAdditionalExtensionsObjectId
	}
	r.ObjectIdPath = dcl.FlattenIntSlice(m["objectIdPath"])

	return r
}

// expandCertificateAuthorityConfigPublicKeyMap expands the contents of CertificateAuthorityConfigPublicKey into a JSON
// request object.
func expandCertificateAuthorityConfigPublicKeyMap(c *Client, f map[string]CertificateAuthorityConfigPublicKey, res *CertificateAuthority) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandCertificateAuthorityConfigPublicKey(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandCertificateAuthorityConfigPublicKeySlice expands the contents of CertificateAuthorityConfigPublicKey into a JSON
// request object.
func expandCertificateAuthorityConfigPublicKeySlice(c *Client, f []CertificateAuthorityConfigPublicKey, res *CertificateAuthority) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandCertificateAuthorityConfigPublicKey(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenCertificateAuthorityConfigPublicKeyMap flattens the contents of CertificateAuthorityConfigPublicKey from a JSON
// response object.
func flattenCertificateAuthorityConfigPublicKeyMap(c *Client, i interface{}, res *CertificateAuthority) map[string]CertificateAuthorityConfigPublicKey {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]CertificateAuthorityConfigPublicKey{}
	}

	if len(a) == 0 {
		return map[string]CertificateAuthorityConfigPublicKey{}
	}

	items := make(map[string]CertificateAuthorityConfigPublicKey)
	for k, item := range a {
		items[k] = *flattenCertificateAuthorityConfigPublicKey(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenCertificateAuthorityConfigPublicKeySlice flattens the contents of CertificateAuthorityConfigPublicKey from a JSON
// response object.
func flattenCertificateAuthorityConfigPublicKeySlice(c *Client, i interface{}, res *CertificateAuthority) []CertificateAuthorityConfigPublicKey {
	a, ok := i.([]interface{})
	if !ok {
		return []CertificateAuthorityConfigPublicKey{}
	}

	if len(a) == 0 {
		return []CertificateAuthorityConfigPublicKey{}
	}

	items := make([]CertificateAuthorityConfigPublicKey, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenCertificateAuthorityConfigPublicKey(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandCertificateAuthorityConfigPublicKey expands an instance of CertificateAuthorityConfigPublicKey into a JSON
// request object.
func expandCertificateAuthorityConfigPublicKey(c *Client, f *CertificateAuthorityConfigPublicKey, res *CertificateAuthority) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.Key; !dcl.IsEmptyValueIndirect(v) {
		m["key"] = v
	}
	if v := f.Format; !dcl.IsEmptyValueIndirect(v) {
		m["format"] = v
	}

	return m, nil
}

// flattenCertificateAuthorityConfigPublicKey flattens an instance of CertificateAuthorityConfigPublicKey from a JSON
// response object.
func flattenCertificateAuthorityConfigPublicKey(c *Client, i interface{}, res *CertificateAuthority) *CertificateAuthorityConfigPublicKey {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &CertificateAuthorityConfigPublicKey{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyCertificateAuthorityConfigPublicKey
	}
	r.Key = dcl.FlattenString(m["key"])
	r.Format = flattenCertificateAuthorityConfigPublicKeyFormatEnum(m["format"])

	return r
}

// expandCertificateAuthorityKeySpecMap expands the contents of CertificateAuthorityKeySpec into a JSON
// request object.
func expandCertificateAuthorityKeySpecMap(c *Client, f map[string]CertificateAuthorityKeySpec, res *CertificateAuthority) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandCertificateAuthorityKeySpec(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandCertificateAuthorityKeySpecSlice expands the contents of CertificateAuthorityKeySpec into a JSON
// request object.
func expandCertificateAuthorityKeySpecSlice(c *Client, f []CertificateAuthorityKeySpec, res *CertificateAuthority) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandCertificateAuthorityKeySpec(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenCertificateAuthorityKeySpecMap flattens the contents of CertificateAuthorityKeySpec from a JSON
// response object.
func flattenCertificateAuthorityKeySpecMap(c *Client, i interface{}, res *CertificateAuthority) map[string]CertificateAuthorityKeySpec {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]CertificateAuthorityKeySpec{}
	}

	if len(a) == 0 {
		return map[string]CertificateAuthorityKeySpec{}
	}

	items := make(map[string]CertificateAuthorityKeySpec)
	for k, item := range a {
		items[k] = *flattenCertificateAuthorityKeySpec(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenCertificateAuthorityKeySpecSlice flattens the contents of CertificateAuthorityKeySpec from a JSON
// response object.
func flattenCertificateAuthorityKeySpecSlice(c *Client, i interface{}, res *CertificateAuthority) []CertificateAuthorityKeySpec {
	a, ok := i.([]interface{})
	if !ok {
		return []CertificateAuthorityKeySpec{}
	}

	if len(a) == 0 {
		return []CertificateAuthorityKeySpec{}
	}

	items := make([]CertificateAuthorityKeySpec, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenCertificateAuthorityKeySpec(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandCertificateAuthorityKeySpec expands an instance of CertificateAuthorityKeySpec into a JSON
// request object.
func expandCertificateAuthorityKeySpec(c *Client, f *CertificateAuthorityKeySpec, res *CertificateAuthority) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.CloudKmsKeyVersion; !dcl.IsEmptyValueIndirect(v) {
		m["cloudKmsKeyVersion"] = v
	}
	if v := f.Algorithm; !dcl.IsEmptyValueIndirect(v) {
		m["algorithm"] = v
	}

	return m, nil
}

// flattenCertificateAuthorityKeySpec flattens an instance of CertificateAuthorityKeySpec from a JSON
// response object.
func flattenCertificateAuthorityKeySpec(c *Client, i interface{}, res *CertificateAuthority) *CertificateAuthorityKeySpec {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &CertificateAuthorityKeySpec{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyCertificateAuthorityKeySpec
	}
	r.CloudKmsKeyVersion = dcl.FlattenString(m["cloudKmsKeyVersion"])
	r.Algorithm = flattenCertificateAuthorityKeySpecAlgorithmEnum(m["algorithm"])

	return r
}

// expandCertificateAuthoritySubordinateConfigMap expands the contents of CertificateAuthoritySubordinateConfig into a JSON
// request object.
func expandCertificateAuthoritySubordinateConfigMap(c *Client, f map[string]CertificateAuthoritySubordinateConfig, res *CertificateAuthority) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandCertificateAuthoritySubordinateConfig(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandCertificateAuthoritySubordinateConfigSlice expands the contents of CertificateAuthoritySubordinateConfig into a JSON
// request object.
func expandCertificateAuthoritySubordinateConfigSlice(c *Client, f []CertificateAuthoritySubordinateConfig, res *CertificateAuthority) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandCertificateAuthoritySubordinateConfig(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenCertificateAuthoritySubordinateConfigMap flattens the contents of CertificateAuthoritySubordinateConfig from a JSON
// response object.
func flattenCertificateAuthoritySubordinateConfigMap(c *Client, i interface{}, res *CertificateAuthority) map[string]CertificateAuthoritySubordinateConfig {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]CertificateAuthoritySubordinateConfig{}
	}

	if len(a) == 0 {
		return map[string]CertificateAuthoritySubordinateConfig{}
	}

	items := make(map[string]CertificateAuthoritySubordinateConfig)
	for k, item := range a {
		items[k] = *flattenCertificateAuthoritySubordinateConfig(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenCertificateAuthoritySubordinateConfigSlice flattens the contents of CertificateAuthoritySubordinateConfig from a JSON
// response object.
func flattenCertificateAuthoritySubordinateConfigSlice(c *Client, i interface{}, res *CertificateAuthority) []CertificateAuthoritySubordinateConfig {
	a, ok := i.([]interface{})
	if !ok {
		return []CertificateAuthoritySubordinateConfig{}
	}

	if len(a) == 0 {
		return []CertificateAuthoritySubordinateConfig{}
	}

	items := make([]CertificateAuthoritySubordinateConfig, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenCertificateAuthoritySubordinateConfig(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandCertificateAuthoritySubordinateConfig expands an instance of CertificateAuthoritySubordinateConfig into a JSON
// request object.
func expandCertificateAuthoritySubordinateConfig(c *Client, f *CertificateAuthoritySubordinateConfig, res *CertificateAuthority) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.CertificateAuthority; !dcl.IsEmptyValueIndirect(v) {
		m["certificateAuthority"] = v
	}
	if v, err := expandCertificateAuthoritySubordinateConfigPemIssuerChain(c, f.PemIssuerChain, res); err != nil {
		return nil, fmt.Errorf("error expanding PemIssuerChain into pemIssuerChain: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["pemIssuerChain"] = v
	}

	return m, nil
}

// flattenCertificateAuthoritySubordinateConfig flattens an instance of CertificateAuthoritySubordinateConfig from a JSON
// response object.
func flattenCertificateAuthoritySubordinateConfig(c *Client, i interface{}, res *CertificateAuthority) *CertificateAuthoritySubordinateConfig {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &CertificateAuthoritySubordinateConfig{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyCertificateAuthoritySubordinateConfig
	}
	r.CertificateAuthority = dcl.FlattenString(m["certificateAuthority"])
	r.PemIssuerChain = flattenCertificateAuthoritySubordinateConfigPemIssuerChain(c, m["pemIssuerChain"], res)

	return r
}

// expandCertificateAuthoritySubordinateConfigPemIssuerChainMap expands the contents of CertificateAuthoritySubordinateConfigPemIssuerChain into a JSON
// request object.
func expandCertificateAuthoritySubordinateConfigPemIssuerChainMap(c *Client, f map[string]CertificateAuthoritySubordinateConfigPemIssuerChain, res *CertificateAuthority) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandCertificateAuthoritySubordinateConfigPemIssuerChain(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandCertificateAuthoritySubordinateConfigPemIssuerChainSlice expands the contents of CertificateAuthoritySubordinateConfigPemIssuerChain into a JSON
// request object.
func expandCertificateAuthoritySubordinateConfigPemIssuerChainSlice(c *Client, f []CertificateAuthoritySubordinateConfigPemIssuerChain, res *CertificateAuthority) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandCertificateAuthoritySubordinateConfigPemIssuerChain(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenCertificateAuthoritySubordinateConfigPemIssuerChainMap flattens the contents of CertificateAuthoritySubordinateConfigPemIssuerChain from a JSON
// response object.
func flattenCertificateAuthoritySubordinateConfigPemIssuerChainMap(c *Client, i interface{}, res *CertificateAuthority) map[string]CertificateAuthoritySubordinateConfigPemIssuerChain {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]CertificateAuthoritySubordinateConfigPemIssuerChain{}
	}

	if len(a) == 0 {
		return map[string]CertificateAuthoritySubordinateConfigPemIssuerChain{}
	}

	items := make(map[string]CertificateAuthoritySubordinateConfigPemIssuerChain)
	for k, item := range a {
		items[k] = *flattenCertificateAuthoritySubordinateConfigPemIssuerChain(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenCertificateAuthoritySubordinateConfigPemIssuerChainSlice flattens the contents of CertificateAuthoritySubordinateConfigPemIssuerChain from a JSON
// response object.
func flattenCertificateAuthoritySubordinateConfigPemIssuerChainSlice(c *Client, i interface{}, res *CertificateAuthority) []CertificateAuthoritySubordinateConfigPemIssuerChain {
	a, ok := i.([]interface{})
	if !ok {
		return []CertificateAuthoritySubordinateConfigPemIssuerChain{}
	}

	if len(a) == 0 {
		return []CertificateAuthoritySubordinateConfigPemIssuerChain{}
	}

	items := make([]CertificateAuthoritySubordinateConfigPemIssuerChain, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenCertificateAuthoritySubordinateConfigPemIssuerChain(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandCertificateAuthoritySubordinateConfigPemIssuerChain expands an instance of CertificateAuthoritySubordinateConfigPemIssuerChain into a JSON
// request object.
func expandCertificateAuthoritySubordinateConfigPemIssuerChain(c *Client, f *CertificateAuthoritySubordinateConfigPemIssuerChain, res *CertificateAuthority) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.PemCertificates; v != nil {
		m["pemCertificates"] = v
	}

	return m, nil
}

// flattenCertificateAuthoritySubordinateConfigPemIssuerChain flattens an instance of CertificateAuthoritySubordinateConfigPemIssuerChain from a JSON
// response object.
func flattenCertificateAuthoritySubordinateConfigPemIssuerChain(c *Client, i interface{}, res *CertificateAuthority) *CertificateAuthoritySubordinateConfigPemIssuerChain {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &CertificateAuthoritySubordinateConfigPemIssuerChain{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyCertificateAuthoritySubordinateConfigPemIssuerChain
	}
	r.PemCertificates = dcl.FlattenStringSlice(m["pemCertificates"])

	return r
}

// expandCertificateAuthorityCaCertificateDescriptionsMap expands the contents of CertificateAuthorityCaCertificateDescriptions into a JSON
// request object.
func expandCertificateAuthorityCaCertificateDescriptionsMap(c *Client, f map[string]CertificateAuthorityCaCertificateDescriptions, res *CertificateAuthority) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandCertificateAuthorityCaCertificateDescriptions(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandCertificateAuthorityCaCertificateDescriptionsSlice expands the contents of CertificateAuthorityCaCertificateDescriptions into a JSON
// request object.
func expandCertificateAuthorityCaCertificateDescriptionsSlice(c *Client, f []CertificateAuthorityCaCertificateDescriptions, res *CertificateAuthority) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandCertificateAuthorityCaCertificateDescriptions(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenCertificateAuthorityCaCertificateDescriptionsMap flattens the contents of CertificateAuthorityCaCertificateDescriptions from a JSON
// response object.
func flattenCertificateAuthorityCaCertificateDescriptionsMap(c *Client, i interface{}, res *CertificateAuthority) map[string]CertificateAuthorityCaCertificateDescriptions {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]CertificateAuthorityCaCertificateDescriptions{}
	}

	if len(a) == 0 {
		return map[string]CertificateAuthorityCaCertificateDescriptions{}
	}

	items := make(map[string]CertificateAuthorityCaCertificateDescriptions)
	for k, item := range a {
		items[k] = *flattenCertificateAuthorityCaCertificateDescriptions(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenCertificateAuthorityCaCertificateDescriptionsSlice flattens the contents of CertificateAuthorityCaCertificateDescriptions from a JSON
// response object.
func flattenCertificateAuthorityCaCertificateDescriptionsSlice(c *Client, i interface{}, res *CertificateAuthority) []CertificateAuthorityCaCertificateDescriptions {
	a, ok := i.([]interface{})
	if !ok {
		return []CertificateAuthorityCaCertificateDescriptions{}
	}

	if len(a) == 0 {
		return []CertificateAuthorityCaCertificateDescriptions{}
	}

	items := make([]CertificateAuthorityCaCertificateDescriptions, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenCertificateAuthorityCaCertificateDescriptions(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandCertificateAuthorityCaCertificateDescriptions expands an instance of CertificateAuthorityCaCertificateDescriptions into a JSON
// request object.
func expandCertificateAuthorityCaCertificateDescriptions(c *Client, f *CertificateAuthorityCaCertificateDescriptions, res *CertificateAuthority) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v, err := expandCertificateAuthorityCaCertificateDescriptionsSubjectDescription(c, f.SubjectDescription, res); err != nil {
		return nil, fmt.Errorf("error expanding SubjectDescription into subjectDescription: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["subjectDescription"] = v
	}
	if v, err := expandCertificateAuthorityCaCertificateDescriptionsX509Description(c, f.X509Description, res); err != nil {
		return nil, fmt.Errorf("error expanding X509Description into x509Description: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["x509Description"] = v
	}
	if v, err := expandCertificateAuthorityCaCertificateDescriptionsPublicKey(c, f.PublicKey, res); err != nil {
		return nil, fmt.Errorf("error expanding PublicKey into publicKey: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["publicKey"] = v
	}
	if v, err := expandCertificateAuthorityCaCertificateDescriptionsSubjectKeyId(c, f.SubjectKeyId, res); err != nil {
		return nil, fmt.Errorf("error expanding SubjectKeyId into subjectKeyId: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["subjectKeyId"] = v
	}
	if v, err := expandCertificateAuthorityCaCertificateDescriptionsAuthorityKeyId(c, f.AuthorityKeyId, res); err != nil {
		return nil, fmt.Errorf("error expanding AuthorityKeyId into authorityKeyId: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["authorityKeyId"] = v
	}
	if v := f.CrlDistributionPoints; v != nil {
		m["crlDistributionPoints"] = v
	}
	if v := f.AiaIssuingCertificateUrls; v != nil {
		m["aiaIssuingCertificateUrls"] = v
	}
	if v, err := expandCertificateAuthorityCaCertificateDescriptionsCertFingerprint(c, f.CertFingerprint, res); err != nil {
		return nil, fmt.Errorf("error expanding CertFingerprint into certFingerprint: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["certFingerprint"] = v
	}

	return m, nil
}

// flattenCertificateAuthorityCaCertificateDescriptions flattens an instance of CertificateAuthorityCaCertificateDescriptions from a JSON
// response object.
func flattenCertificateAuthorityCaCertificateDescriptions(c *Client, i interface{}, res *CertificateAuthority) *CertificateAuthorityCaCertificateDescriptions {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &CertificateAuthorityCaCertificateDescriptions{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyCertificateAuthorityCaCertificateDescriptions
	}
	r.SubjectDescription = flattenCertificateAuthorityCaCertificateDescriptionsSubjectDescription(c, m["subjectDescription"], res)
	r.X509Description = flattenCertificateAuthorityCaCertificateDescriptionsX509Description(c, m["x509Description"], res)
	r.PublicKey = flattenCertificateAuthorityCaCertificateDescriptionsPublicKey(c, m["publicKey"], res)
	r.SubjectKeyId = flattenCertificateAuthorityCaCertificateDescriptionsSubjectKeyId(c, m["subjectKeyId"], res)
	r.AuthorityKeyId = flattenCertificateAuthorityCaCertificateDescriptionsAuthorityKeyId(c, m["authorityKeyId"], res)
	r.CrlDistributionPoints = dcl.FlattenStringSlice(m["crlDistributionPoints"])
	r.AiaIssuingCertificateUrls = dcl.FlattenStringSlice(m["aiaIssuingCertificateUrls"])
	r.CertFingerprint = flattenCertificateAuthorityCaCertificateDescriptionsCertFingerprint(c, m["certFingerprint"], res)

	return r
}

// expandCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionMap expands the contents of CertificateAuthorityCaCertificateDescriptionsSubjectDescription into a JSON
// request object.
func expandCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionMap(c *Client, f map[string]CertificateAuthorityCaCertificateDescriptionsSubjectDescription, res *CertificateAuthority) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandCertificateAuthorityCaCertificateDescriptionsSubjectDescription(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSlice expands the contents of CertificateAuthorityCaCertificateDescriptionsSubjectDescription into a JSON
// request object.
func expandCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSlice(c *Client, f []CertificateAuthorityCaCertificateDescriptionsSubjectDescription, res *CertificateAuthority) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandCertificateAuthorityCaCertificateDescriptionsSubjectDescription(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionMap flattens the contents of CertificateAuthorityCaCertificateDescriptionsSubjectDescription from a JSON
// response object.
func flattenCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionMap(c *Client, i interface{}, res *CertificateAuthority) map[string]CertificateAuthorityCaCertificateDescriptionsSubjectDescription {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]CertificateAuthorityCaCertificateDescriptionsSubjectDescription{}
	}

	if len(a) == 0 {
		return map[string]CertificateAuthorityCaCertificateDescriptionsSubjectDescription{}
	}

	items := make(map[string]CertificateAuthorityCaCertificateDescriptionsSubjectDescription)
	for k, item := range a {
		items[k] = *flattenCertificateAuthorityCaCertificateDescriptionsSubjectDescription(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSlice flattens the contents of CertificateAuthorityCaCertificateDescriptionsSubjectDescription from a JSON
// response object.
func flattenCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSlice(c *Client, i interface{}, res *CertificateAuthority) []CertificateAuthorityCaCertificateDescriptionsSubjectDescription {
	a, ok := i.([]interface{})
	if !ok {
		return []CertificateAuthorityCaCertificateDescriptionsSubjectDescription{}
	}

	if len(a) == 0 {
		return []CertificateAuthorityCaCertificateDescriptionsSubjectDescription{}
	}

	items := make([]CertificateAuthorityCaCertificateDescriptionsSubjectDescription, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenCertificateAuthorityCaCertificateDescriptionsSubjectDescription(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandCertificateAuthorityCaCertificateDescriptionsSubjectDescription expands an instance of CertificateAuthorityCaCertificateDescriptionsSubjectDescription into a JSON
// request object.
func expandCertificateAuthorityCaCertificateDescriptionsSubjectDescription(c *Client, f *CertificateAuthorityCaCertificateDescriptionsSubjectDescription, res *CertificateAuthority) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v, err := expandCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubject(c, f.Subject, res); err != nil {
		return nil, fmt.Errorf("error expanding Subject into subject: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["subject"] = v
	}
	if v, err := expandCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltName(c, f.SubjectAltName, res); err != nil {
		return nil, fmt.Errorf("error expanding SubjectAltName into subjectAltName: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["subjectAltName"] = v
	}
	if v := f.HexSerialNumber; !dcl.IsEmptyValueIndirect(v) {
		m["hexSerialNumber"] = v
	}
	if v := f.Lifetime; !dcl.IsEmptyValueIndirect(v) {
		m["lifetime"] = v
	}
	if v := f.NotBeforeTime; !dcl.IsEmptyValueIndirect(v) {
		m["notBeforeTime"] = v
	}
	if v := f.NotAfterTime; !dcl.IsEmptyValueIndirect(v) {
		m["notAfterTime"] = v
	}

	return m, nil
}

// flattenCertificateAuthorityCaCertificateDescriptionsSubjectDescription flattens an instance of CertificateAuthorityCaCertificateDescriptionsSubjectDescription from a JSON
// response object.
func flattenCertificateAuthorityCaCertificateDescriptionsSubjectDescription(c *Client, i interface{}, res *CertificateAuthority) *CertificateAuthorityCaCertificateDescriptionsSubjectDescription {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &CertificateAuthorityCaCertificateDescriptionsSubjectDescription{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyCertificateAuthorityCaCertificateDescriptionsSubjectDescription
	}
	r.Subject = flattenCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubject(c, m["subject"], res)
	r.SubjectAltName = flattenCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltName(c, m["subjectAltName"], res)
	r.HexSerialNumber = dcl.FlattenString(m["hexSerialNumber"])
	r.Lifetime = dcl.FlattenString(m["lifetime"])
	r.NotBeforeTime = dcl.FlattenString(m["notBeforeTime"])
	r.NotAfterTime = dcl.FlattenString(m["notAfterTime"])

	return r
}

// expandCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectMap expands the contents of CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubject into a JSON
// request object.
func expandCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectMap(c *Client, f map[string]CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubject, res *CertificateAuthority) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubject(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectSlice expands the contents of CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubject into a JSON
// request object.
func expandCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectSlice(c *Client, f []CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubject, res *CertificateAuthority) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubject(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectMap flattens the contents of CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubject from a JSON
// response object.
func flattenCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectMap(c *Client, i interface{}, res *CertificateAuthority) map[string]CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubject {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubject{}
	}

	if len(a) == 0 {
		return map[string]CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubject{}
	}

	items := make(map[string]CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubject)
	for k, item := range a {
		items[k] = *flattenCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubject(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectSlice flattens the contents of CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubject from a JSON
// response object.
func flattenCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectSlice(c *Client, i interface{}, res *CertificateAuthority) []CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubject {
	a, ok := i.([]interface{})
	if !ok {
		return []CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubject{}
	}

	if len(a) == 0 {
		return []CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubject{}
	}

	items := make([]CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubject, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubject(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubject expands an instance of CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubject into a JSON
// request object.
func expandCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubject(c *Client, f *CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubject, res *CertificateAuthority) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.CommonName; !dcl.IsEmptyValueIndirect(v) {
		m["commonName"] = v
	}
	if v := f.CountryCode; !dcl.IsEmptyValueIndirect(v) {
		m["countryCode"] = v
	}
	if v := f.Organization; !dcl.IsEmptyValueIndirect(v) {
		m["organization"] = v
	}
	if v := f.OrganizationalUnit; !dcl.IsEmptyValueIndirect(v) {
		m["organizationalUnit"] = v
	}
	if v := f.Locality; !dcl.IsEmptyValueIndirect(v) {
		m["locality"] = v
	}
	if v := f.Province; !dcl.IsEmptyValueIndirect(v) {
		m["province"] = v
	}
	if v := f.StreetAddress; !dcl.IsEmptyValueIndirect(v) {
		m["streetAddress"] = v
	}
	if v := f.PostalCode; !dcl.IsEmptyValueIndirect(v) {
		m["postalCode"] = v
	}

	return m, nil
}

// flattenCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubject flattens an instance of CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubject from a JSON
// response object.
func flattenCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubject(c *Client, i interface{}, res *CertificateAuthority) *CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubject {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubject{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubject
	}
	r.CommonName = dcl.FlattenString(m["commonName"])
	r.CountryCode = dcl.FlattenString(m["countryCode"])
	r.Organization = dcl.FlattenString(m["organization"])
	r.OrganizationalUnit = dcl.FlattenString(m["organizationalUnit"])
	r.Locality = dcl.FlattenString(m["locality"])
	r.Province = dcl.FlattenString(m["province"])
	r.StreetAddress = dcl.FlattenString(m["streetAddress"])
	r.PostalCode = dcl.FlattenString(m["postalCode"])

	return r
}

// expandCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameMap expands the contents of CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltName into a JSON
// request object.
func expandCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameMap(c *Client, f map[string]CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltName, res *CertificateAuthority) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltName(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameSlice expands the contents of CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltName into a JSON
// request object.
func expandCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameSlice(c *Client, f []CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltName, res *CertificateAuthority) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltName(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameMap flattens the contents of CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltName from a JSON
// response object.
func flattenCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameMap(c *Client, i interface{}, res *CertificateAuthority) map[string]CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltName {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltName{}
	}

	if len(a) == 0 {
		return map[string]CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltName{}
	}

	items := make(map[string]CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltName)
	for k, item := range a {
		items[k] = *flattenCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltName(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameSlice flattens the contents of CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltName from a JSON
// response object.
func flattenCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameSlice(c *Client, i interface{}, res *CertificateAuthority) []CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltName {
	a, ok := i.([]interface{})
	if !ok {
		return []CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltName{}
	}

	if len(a) == 0 {
		return []CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltName{}
	}

	items := make([]CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltName, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltName(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltName expands an instance of CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltName into a JSON
// request object.
func expandCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltName(c *Client, f *CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltName, res *CertificateAuthority) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.DnsNames; v != nil {
		m["dnsNames"] = v
	}
	if v := f.Uris; v != nil {
		m["uris"] = v
	}
	if v := f.EmailAddresses; v != nil {
		m["emailAddresses"] = v
	}
	if v := f.IPAddresses; v != nil {
		m["ipAddresses"] = v
	}
	if v, err := expandCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSansSlice(c, f.CustomSans, res); err != nil {
		return nil, fmt.Errorf("error expanding CustomSans into customSans: %w", err)
	} else if v != nil {
		m["customSans"] = v
	}

	return m, nil
}

// flattenCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltName flattens an instance of CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltName from a JSON
// response object.
func flattenCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltName(c *Client, i interface{}, res *CertificateAuthority) *CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltName {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltName{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltName
	}
	r.DnsNames = dcl.FlattenStringSlice(m["dnsNames"])
	r.Uris = dcl.FlattenStringSlice(m["uris"])
	r.EmailAddresses = dcl.FlattenStringSlice(m["emailAddresses"])
	r.IPAddresses = dcl.FlattenStringSlice(m["ipAddresses"])
	r.CustomSans = flattenCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSansSlice(c, m["customSans"], res)

	return r
}

// expandCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSansMap expands the contents of CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSans into a JSON
// request object.
func expandCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSansMap(c *Client, f map[string]CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSans, res *CertificateAuthority) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSans(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSansSlice expands the contents of CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSans into a JSON
// request object.
func expandCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSansSlice(c *Client, f []CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSans, res *CertificateAuthority) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSans(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSansMap flattens the contents of CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSans from a JSON
// response object.
func flattenCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSansMap(c *Client, i interface{}, res *CertificateAuthority) map[string]CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSans {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSans{}
	}

	if len(a) == 0 {
		return map[string]CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSans{}
	}

	items := make(map[string]CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSans)
	for k, item := range a {
		items[k] = *flattenCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSans(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSansSlice flattens the contents of CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSans from a JSON
// response object.
func flattenCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSansSlice(c *Client, i interface{}, res *CertificateAuthority) []CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSans {
	a, ok := i.([]interface{})
	if !ok {
		return []CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSans{}
	}

	if len(a) == 0 {
		return []CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSans{}
	}

	items := make([]CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSans, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSans(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSans expands an instance of CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSans into a JSON
// request object.
func expandCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSans(c *Client, f *CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSans, res *CertificateAuthority) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v, err := expandCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSansObjectId(c, f.ObjectId, res); err != nil {
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

// flattenCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSans flattens an instance of CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSans from a JSON
// response object.
func flattenCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSans(c *Client, i interface{}, res *CertificateAuthority) *CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSans {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSans{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSans
	}
	r.ObjectId = flattenCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSansObjectId(c, m["objectId"], res)
	r.Critical = dcl.FlattenBool(m["critical"])
	r.Value = dcl.FlattenString(m["value"])

	return r
}

// expandCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSansObjectIdMap expands the contents of CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSansObjectId into a JSON
// request object.
func expandCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSansObjectIdMap(c *Client, f map[string]CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSansObjectId, res *CertificateAuthority) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSansObjectId(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSansObjectIdSlice expands the contents of CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSansObjectId into a JSON
// request object.
func expandCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSansObjectIdSlice(c *Client, f []CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSansObjectId, res *CertificateAuthority) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSansObjectId(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSansObjectIdMap flattens the contents of CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSansObjectId from a JSON
// response object.
func flattenCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSansObjectIdMap(c *Client, i interface{}, res *CertificateAuthority) map[string]CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSansObjectId {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSansObjectId{}
	}

	if len(a) == 0 {
		return map[string]CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSansObjectId{}
	}

	items := make(map[string]CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSansObjectId)
	for k, item := range a {
		items[k] = *flattenCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSansObjectId(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSansObjectIdSlice flattens the contents of CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSansObjectId from a JSON
// response object.
func flattenCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSansObjectIdSlice(c *Client, i interface{}, res *CertificateAuthority) []CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSansObjectId {
	a, ok := i.([]interface{})
	if !ok {
		return []CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSansObjectId{}
	}

	if len(a) == 0 {
		return []CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSansObjectId{}
	}

	items := make([]CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSansObjectId, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSansObjectId(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSansObjectId expands an instance of CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSansObjectId into a JSON
// request object.
func expandCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSansObjectId(c *Client, f *CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSansObjectId, res *CertificateAuthority) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.ObjectIdPath; v != nil {
		m["objectIdPath"] = v
	}

	return m, nil
}

// flattenCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSansObjectId flattens an instance of CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSansObjectId from a JSON
// response object.
func flattenCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSansObjectId(c *Client, i interface{}, res *CertificateAuthority) *CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSansObjectId {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSansObjectId{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSansObjectId
	}
	r.ObjectIdPath = dcl.FlattenIntSlice(m["objectIdPath"])

	return r
}

// expandCertificateAuthorityCaCertificateDescriptionsX509DescriptionMap expands the contents of CertificateAuthorityCaCertificateDescriptionsX509Description into a JSON
// request object.
func expandCertificateAuthorityCaCertificateDescriptionsX509DescriptionMap(c *Client, f map[string]CertificateAuthorityCaCertificateDescriptionsX509Description, res *CertificateAuthority) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandCertificateAuthorityCaCertificateDescriptionsX509Description(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandCertificateAuthorityCaCertificateDescriptionsX509DescriptionSlice expands the contents of CertificateAuthorityCaCertificateDescriptionsX509Description into a JSON
// request object.
func expandCertificateAuthorityCaCertificateDescriptionsX509DescriptionSlice(c *Client, f []CertificateAuthorityCaCertificateDescriptionsX509Description, res *CertificateAuthority) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandCertificateAuthorityCaCertificateDescriptionsX509Description(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenCertificateAuthorityCaCertificateDescriptionsX509DescriptionMap flattens the contents of CertificateAuthorityCaCertificateDescriptionsX509Description from a JSON
// response object.
func flattenCertificateAuthorityCaCertificateDescriptionsX509DescriptionMap(c *Client, i interface{}, res *CertificateAuthority) map[string]CertificateAuthorityCaCertificateDescriptionsX509Description {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]CertificateAuthorityCaCertificateDescriptionsX509Description{}
	}

	if len(a) == 0 {
		return map[string]CertificateAuthorityCaCertificateDescriptionsX509Description{}
	}

	items := make(map[string]CertificateAuthorityCaCertificateDescriptionsX509Description)
	for k, item := range a {
		items[k] = *flattenCertificateAuthorityCaCertificateDescriptionsX509Description(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenCertificateAuthorityCaCertificateDescriptionsX509DescriptionSlice flattens the contents of CertificateAuthorityCaCertificateDescriptionsX509Description from a JSON
// response object.
func flattenCertificateAuthorityCaCertificateDescriptionsX509DescriptionSlice(c *Client, i interface{}, res *CertificateAuthority) []CertificateAuthorityCaCertificateDescriptionsX509Description {
	a, ok := i.([]interface{})
	if !ok {
		return []CertificateAuthorityCaCertificateDescriptionsX509Description{}
	}

	if len(a) == 0 {
		return []CertificateAuthorityCaCertificateDescriptionsX509Description{}
	}

	items := make([]CertificateAuthorityCaCertificateDescriptionsX509Description, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenCertificateAuthorityCaCertificateDescriptionsX509Description(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandCertificateAuthorityCaCertificateDescriptionsX509Description expands an instance of CertificateAuthorityCaCertificateDescriptionsX509Description into a JSON
// request object.
func expandCertificateAuthorityCaCertificateDescriptionsX509Description(c *Client, f *CertificateAuthorityCaCertificateDescriptionsX509Description, res *CertificateAuthority) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v, err := expandCertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsage(c, f.KeyUsage, res); err != nil {
		return nil, fmt.Errorf("error expanding KeyUsage into keyUsage: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["keyUsage"] = v
	}
	if v, err := expandCertificateAuthorityCaCertificateDescriptionsX509DescriptionCaOptions(c, f.CaOptions, res); err != nil {
		return nil, fmt.Errorf("error expanding CaOptions into caOptions: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["caOptions"] = v
	}
	if v, err := expandCertificateAuthorityCaCertificateDescriptionsX509DescriptionPolicyIdsSlice(c, f.PolicyIds, res); err != nil {
		return nil, fmt.Errorf("error expanding PolicyIds into policyIds: %w", err)
	} else if v != nil {
		m["policyIds"] = v
	}
	if v, err := expandCertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensionsSlice(c, f.AdditionalExtensions, res); err != nil {
		return nil, fmt.Errorf("error expanding AdditionalExtensions into additionalExtensions: %w", err)
	} else if v != nil {
		m["additionalExtensions"] = v
	}

	return m, nil
}

// flattenCertificateAuthorityCaCertificateDescriptionsX509Description flattens an instance of CertificateAuthorityCaCertificateDescriptionsX509Description from a JSON
// response object.
func flattenCertificateAuthorityCaCertificateDescriptionsX509Description(c *Client, i interface{}, res *CertificateAuthority) *CertificateAuthorityCaCertificateDescriptionsX509Description {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &CertificateAuthorityCaCertificateDescriptionsX509Description{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyCertificateAuthorityCaCertificateDescriptionsX509Description
	}
	r.KeyUsage = flattenCertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsage(c, m["keyUsage"], res)
	r.CaOptions = flattenCertificateAuthorityCaCertificateDescriptionsX509DescriptionCaOptions(c, m["caOptions"], res)
	r.PolicyIds = flattenCertificateAuthorityCaCertificateDescriptionsX509DescriptionPolicyIdsSlice(c, m["policyIds"], res)
	r.AiaOcspServers = dcl.FlattenStringSlice(m["aiaOcspServers"])
	r.AdditionalExtensions = flattenCertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensionsSlice(c, m["additionalExtensions"], res)

	return r
}

// expandCertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageMap expands the contents of CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsage into a JSON
// request object.
func expandCertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageMap(c *Client, f map[string]CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsage, res *CertificateAuthority) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandCertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsage(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandCertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageSlice expands the contents of CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsage into a JSON
// request object.
func expandCertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageSlice(c *Client, f []CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsage, res *CertificateAuthority) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandCertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsage(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenCertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageMap flattens the contents of CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsage from a JSON
// response object.
func flattenCertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageMap(c *Client, i interface{}, res *CertificateAuthority) map[string]CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsage {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsage{}
	}

	if len(a) == 0 {
		return map[string]CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsage{}
	}

	items := make(map[string]CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsage)
	for k, item := range a {
		items[k] = *flattenCertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsage(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenCertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageSlice flattens the contents of CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsage from a JSON
// response object.
func flattenCertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageSlice(c *Client, i interface{}, res *CertificateAuthority) []CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsage {
	a, ok := i.([]interface{})
	if !ok {
		return []CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsage{}
	}

	if len(a) == 0 {
		return []CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsage{}
	}

	items := make([]CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsage, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenCertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsage(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandCertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsage expands an instance of CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsage into a JSON
// request object.
func expandCertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsage(c *Client, f *CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsage, res *CertificateAuthority) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v, err := expandCertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageBaseKeyUsage(c, f.BaseKeyUsage, res); err != nil {
		return nil, fmt.Errorf("error expanding BaseKeyUsage into baseKeyUsage: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["baseKeyUsage"] = v
	}
	if v, err := expandCertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageExtendedKeyUsage(c, f.ExtendedKeyUsage, res); err != nil {
		return nil, fmt.Errorf("error expanding ExtendedKeyUsage into extendedKeyUsage: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["extendedKeyUsage"] = v
	}
	if v, err := expandCertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageUnknownExtendedKeyUsagesSlice(c, f.UnknownExtendedKeyUsages, res); err != nil {
		return nil, fmt.Errorf("error expanding UnknownExtendedKeyUsages into unknownExtendedKeyUsages: %w", err)
	} else if v != nil {
		m["unknownExtendedKeyUsages"] = v
	}

	return m, nil
}

// flattenCertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsage flattens an instance of CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsage from a JSON
// response object.
func flattenCertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsage(c *Client, i interface{}, res *CertificateAuthority) *CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsage {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsage{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyCertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsage
	}
	r.BaseKeyUsage = flattenCertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageBaseKeyUsage(c, m["baseKeyUsage"], res)
	r.ExtendedKeyUsage = flattenCertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageExtendedKeyUsage(c, m["extendedKeyUsage"], res)
	r.UnknownExtendedKeyUsages = flattenCertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageUnknownExtendedKeyUsagesSlice(c, m["unknownExtendedKeyUsages"], res)

	return r
}

// expandCertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageBaseKeyUsageMap expands the contents of CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageBaseKeyUsage into a JSON
// request object.
func expandCertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageBaseKeyUsageMap(c *Client, f map[string]CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageBaseKeyUsage, res *CertificateAuthority) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandCertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageBaseKeyUsage(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandCertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageBaseKeyUsageSlice expands the contents of CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageBaseKeyUsage into a JSON
// request object.
func expandCertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageBaseKeyUsageSlice(c *Client, f []CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageBaseKeyUsage, res *CertificateAuthority) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandCertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageBaseKeyUsage(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenCertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageBaseKeyUsageMap flattens the contents of CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageBaseKeyUsage from a JSON
// response object.
func flattenCertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageBaseKeyUsageMap(c *Client, i interface{}, res *CertificateAuthority) map[string]CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageBaseKeyUsage {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageBaseKeyUsage{}
	}

	if len(a) == 0 {
		return map[string]CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageBaseKeyUsage{}
	}

	items := make(map[string]CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageBaseKeyUsage)
	for k, item := range a {
		items[k] = *flattenCertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageBaseKeyUsage(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenCertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageBaseKeyUsageSlice flattens the contents of CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageBaseKeyUsage from a JSON
// response object.
func flattenCertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageBaseKeyUsageSlice(c *Client, i interface{}, res *CertificateAuthority) []CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageBaseKeyUsage {
	a, ok := i.([]interface{})
	if !ok {
		return []CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageBaseKeyUsage{}
	}

	if len(a) == 0 {
		return []CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageBaseKeyUsage{}
	}

	items := make([]CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageBaseKeyUsage, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenCertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageBaseKeyUsage(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandCertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageBaseKeyUsage expands an instance of CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageBaseKeyUsage into a JSON
// request object.
func expandCertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageBaseKeyUsage(c *Client, f *CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageBaseKeyUsage, res *CertificateAuthority) (map[string]interface{}, error) {
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

// flattenCertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageBaseKeyUsage flattens an instance of CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageBaseKeyUsage from a JSON
// response object.
func flattenCertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageBaseKeyUsage(c *Client, i interface{}, res *CertificateAuthority) *CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageBaseKeyUsage {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageBaseKeyUsage{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyCertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageBaseKeyUsage
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

// expandCertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageExtendedKeyUsageMap expands the contents of CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageExtendedKeyUsage into a JSON
// request object.
func expandCertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageExtendedKeyUsageMap(c *Client, f map[string]CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageExtendedKeyUsage, res *CertificateAuthority) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandCertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageExtendedKeyUsage(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandCertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageExtendedKeyUsageSlice expands the contents of CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageExtendedKeyUsage into a JSON
// request object.
func expandCertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageExtendedKeyUsageSlice(c *Client, f []CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageExtendedKeyUsage, res *CertificateAuthority) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandCertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageExtendedKeyUsage(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenCertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageExtendedKeyUsageMap flattens the contents of CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageExtendedKeyUsage from a JSON
// response object.
func flattenCertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageExtendedKeyUsageMap(c *Client, i interface{}, res *CertificateAuthority) map[string]CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageExtendedKeyUsage {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageExtendedKeyUsage{}
	}

	if len(a) == 0 {
		return map[string]CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageExtendedKeyUsage{}
	}

	items := make(map[string]CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageExtendedKeyUsage)
	for k, item := range a {
		items[k] = *flattenCertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageExtendedKeyUsage(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenCertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageExtendedKeyUsageSlice flattens the contents of CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageExtendedKeyUsage from a JSON
// response object.
func flattenCertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageExtendedKeyUsageSlice(c *Client, i interface{}, res *CertificateAuthority) []CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageExtendedKeyUsage {
	a, ok := i.([]interface{})
	if !ok {
		return []CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageExtendedKeyUsage{}
	}

	if len(a) == 0 {
		return []CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageExtendedKeyUsage{}
	}

	items := make([]CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageExtendedKeyUsage, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenCertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageExtendedKeyUsage(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandCertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageExtendedKeyUsage expands an instance of CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageExtendedKeyUsage into a JSON
// request object.
func expandCertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageExtendedKeyUsage(c *Client, f *CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageExtendedKeyUsage, res *CertificateAuthority) (map[string]interface{}, error) {
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

// flattenCertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageExtendedKeyUsage flattens an instance of CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageExtendedKeyUsage from a JSON
// response object.
func flattenCertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageExtendedKeyUsage(c *Client, i interface{}, res *CertificateAuthority) *CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageExtendedKeyUsage {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageExtendedKeyUsage{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyCertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageExtendedKeyUsage
	}
	r.ServerAuth = dcl.FlattenBool(m["serverAuth"])
	r.ClientAuth = dcl.FlattenBool(m["clientAuth"])
	r.CodeSigning = dcl.FlattenBool(m["codeSigning"])
	r.EmailProtection = dcl.FlattenBool(m["emailProtection"])
	r.TimeStamping = dcl.FlattenBool(m["timeStamping"])
	r.OcspSigning = dcl.FlattenBool(m["ocspSigning"])

	return r
}

// expandCertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageUnknownExtendedKeyUsagesMap expands the contents of CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageUnknownExtendedKeyUsages into a JSON
// request object.
func expandCertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageUnknownExtendedKeyUsagesMap(c *Client, f map[string]CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageUnknownExtendedKeyUsages, res *CertificateAuthority) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandCertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageUnknownExtendedKeyUsages(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandCertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageUnknownExtendedKeyUsagesSlice expands the contents of CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageUnknownExtendedKeyUsages into a JSON
// request object.
func expandCertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageUnknownExtendedKeyUsagesSlice(c *Client, f []CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageUnknownExtendedKeyUsages, res *CertificateAuthority) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandCertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageUnknownExtendedKeyUsages(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenCertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageUnknownExtendedKeyUsagesMap flattens the contents of CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageUnknownExtendedKeyUsages from a JSON
// response object.
func flattenCertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageUnknownExtendedKeyUsagesMap(c *Client, i interface{}, res *CertificateAuthority) map[string]CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageUnknownExtendedKeyUsages {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageUnknownExtendedKeyUsages{}
	}

	if len(a) == 0 {
		return map[string]CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageUnknownExtendedKeyUsages{}
	}

	items := make(map[string]CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageUnknownExtendedKeyUsages)
	for k, item := range a {
		items[k] = *flattenCertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageUnknownExtendedKeyUsages(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenCertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageUnknownExtendedKeyUsagesSlice flattens the contents of CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageUnknownExtendedKeyUsages from a JSON
// response object.
func flattenCertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageUnknownExtendedKeyUsagesSlice(c *Client, i interface{}, res *CertificateAuthority) []CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageUnknownExtendedKeyUsages {
	a, ok := i.([]interface{})
	if !ok {
		return []CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageUnknownExtendedKeyUsages{}
	}

	if len(a) == 0 {
		return []CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageUnknownExtendedKeyUsages{}
	}

	items := make([]CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageUnknownExtendedKeyUsages, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenCertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageUnknownExtendedKeyUsages(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandCertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageUnknownExtendedKeyUsages expands an instance of CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageUnknownExtendedKeyUsages into a JSON
// request object.
func expandCertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageUnknownExtendedKeyUsages(c *Client, f *CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageUnknownExtendedKeyUsages, res *CertificateAuthority) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.ObjectIdPath; v != nil {
		m["objectIdPath"] = v
	}

	return m, nil
}

// flattenCertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageUnknownExtendedKeyUsages flattens an instance of CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageUnknownExtendedKeyUsages from a JSON
// response object.
func flattenCertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageUnknownExtendedKeyUsages(c *Client, i interface{}, res *CertificateAuthority) *CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageUnknownExtendedKeyUsages {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageUnknownExtendedKeyUsages{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyCertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageUnknownExtendedKeyUsages
	}
	r.ObjectIdPath = dcl.FlattenIntSlice(m["objectIdPath"])

	return r
}

// expandCertificateAuthorityCaCertificateDescriptionsX509DescriptionCaOptionsMap expands the contents of CertificateAuthorityCaCertificateDescriptionsX509DescriptionCaOptions into a JSON
// request object.
func expandCertificateAuthorityCaCertificateDescriptionsX509DescriptionCaOptionsMap(c *Client, f map[string]CertificateAuthorityCaCertificateDescriptionsX509DescriptionCaOptions, res *CertificateAuthority) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandCertificateAuthorityCaCertificateDescriptionsX509DescriptionCaOptions(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandCertificateAuthorityCaCertificateDescriptionsX509DescriptionCaOptionsSlice expands the contents of CertificateAuthorityCaCertificateDescriptionsX509DescriptionCaOptions into a JSON
// request object.
func expandCertificateAuthorityCaCertificateDescriptionsX509DescriptionCaOptionsSlice(c *Client, f []CertificateAuthorityCaCertificateDescriptionsX509DescriptionCaOptions, res *CertificateAuthority) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandCertificateAuthorityCaCertificateDescriptionsX509DescriptionCaOptions(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenCertificateAuthorityCaCertificateDescriptionsX509DescriptionCaOptionsMap flattens the contents of CertificateAuthorityCaCertificateDescriptionsX509DescriptionCaOptions from a JSON
// response object.
func flattenCertificateAuthorityCaCertificateDescriptionsX509DescriptionCaOptionsMap(c *Client, i interface{}, res *CertificateAuthority) map[string]CertificateAuthorityCaCertificateDescriptionsX509DescriptionCaOptions {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]CertificateAuthorityCaCertificateDescriptionsX509DescriptionCaOptions{}
	}

	if len(a) == 0 {
		return map[string]CertificateAuthorityCaCertificateDescriptionsX509DescriptionCaOptions{}
	}

	items := make(map[string]CertificateAuthorityCaCertificateDescriptionsX509DescriptionCaOptions)
	for k, item := range a {
		items[k] = *flattenCertificateAuthorityCaCertificateDescriptionsX509DescriptionCaOptions(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenCertificateAuthorityCaCertificateDescriptionsX509DescriptionCaOptionsSlice flattens the contents of CertificateAuthorityCaCertificateDescriptionsX509DescriptionCaOptions from a JSON
// response object.
func flattenCertificateAuthorityCaCertificateDescriptionsX509DescriptionCaOptionsSlice(c *Client, i interface{}, res *CertificateAuthority) []CertificateAuthorityCaCertificateDescriptionsX509DescriptionCaOptions {
	a, ok := i.([]interface{})
	if !ok {
		return []CertificateAuthorityCaCertificateDescriptionsX509DescriptionCaOptions{}
	}

	if len(a) == 0 {
		return []CertificateAuthorityCaCertificateDescriptionsX509DescriptionCaOptions{}
	}

	items := make([]CertificateAuthorityCaCertificateDescriptionsX509DescriptionCaOptions, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenCertificateAuthorityCaCertificateDescriptionsX509DescriptionCaOptions(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandCertificateAuthorityCaCertificateDescriptionsX509DescriptionCaOptions expands an instance of CertificateAuthorityCaCertificateDescriptionsX509DescriptionCaOptions into a JSON
// request object.
func expandCertificateAuthorityCaCertificateDescriptionsX509DescriptionCaOptions(c *Client, f *CertificateAuthorityCaCertificateDescriptionsX509DescriptionCaOptions, res *CertificateAuthority) (map[string]interface{}, error) {
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

	return m, nil
}

// flattenCertificateAuthorityCaCertificateDescriptionsX509DescriptionCaOptions flattens an instance of CertificateAuthorityCaCertificateDescriptionsX509DescriptionCaOptions from a JSON
// response object.
func flattenCertificateAuthorityCaCertificateDescriptionsX509DescriptionCaOptions(c *Client, i interface{}, res *CertificateAuthority) *CertificateAuthorityCaCertificateDescriptionsX509DescriptionCaOptions {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &CertificateAuthorityCaCertificateDescriptionsX509DescriptionCaOptions{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyCertificateAuthorityCaCertificateDescriptionsX509DescriptionCaOptions
	}
	r.IsCa = dcl.FlattenBool(m["isCa"])
	r.MaxIssuerPathLength = dcl.FlattenInteger(m["maxIssuerPathLength"])

	return r
}

// expandCertificateAuthorityCaCertificateDescriptionsX509DescriptionPolicyIdsMap expands the contents of CertificateAuthorityCaCertificateDescriptionsX509DescriptionPolicyIds into a JSON
// request object.
func expandCertificateAuthorityCaCertificateDescriptionsX509DescriptionPolicyIdsMap(c *Client, f map[string]CertificateAuthorityCaCertificateDescriptionsX509DescriptionPolicyIds, res *CertificateAuthority) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandCertificateAuthorityCaCertificateDescriptionsX509DescriptionPolicyIds(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandCertificateAuthorityCaCertificateDescriptionsX509DescriptionPolicyIdsSlice expands the contents of CertificateAuthorityCaCertificateDescriptionsX509DescriptionPolicyIds into a JSON
// request object.
func expandCertificateAuthorityCaCertificateDescriptionsX509DescriptionPolicyIdsSlice(c *Client, f []CertificateAuthorityCaCertificateDescriptionsX509DescriptionPolicyIds, res *CertificateAuthority) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandCertificateAuthorityCaCertificateDescriptionsX509DescriptionPolicyIds(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenCertificateAuthorityCaCertificateDescriptionsX509DescriptionPolicyIdsMap flattens the contents of CertificateAuthorityCaCertificateDescriptionsX509DescriptionPolicyIds from a JSON
// response object.
func flattenCertificateAuthorityCaCertificateDescriptionsX509DescriptionPolicyIdsMap(c *Client, i interface{}, res *CertificateAuthority) map[string]CertificateAuthorityCaCertificateDescriptionsX509DescriptionPolicyIds {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]CertificateAuthorityCaCertificateDescriptionsX509DescriptionPolicyIds{}
	}

	if len(a) == 0 {
		return map[string]CertificateAuthorityCaCertificateDescriptionsX509DescriptionPolicyIds{}
	}

	items := make(map[string]CertificateAuthorityCaCertificateDescriptionsX509DescriptionPolicyIds)
	for k, item := range a {
		items[k] = *flattenCertificateAuthorityCaCertificateDescriptionsX509DescriptionPolicyIds(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenCertificateAuthorityCaCertificateDescriptionsX509DescriptionPolicyIdsSlice flattens the contents of CertificateAuthorityCaCertificateDescriptionsX509DescriptionPolicyIds from a JSON
// response object.
func flattenCertificateAuthorityCaCertificateDescriptionsX509DescriptionPolicyIdsSlice(c *Client, i interface{}, res *CertificateAuthority) []CertificateAuthorityCaCertificateDescriptionsX509DescriptionPolicyIds {
	a, ok := i.([]interface{})
	if !ok {
		return []CertificateAuthorityCaCertificateDescriptionsX509DescriptionPolicyIds{}
	}

	if len(a) == 0 {
		return []CertificateAuthorityCaCertificateDescriptionsX509DescriptionPolicyIds{}
	}

	items := make([]CertificateAuthorityCaCertificateDescriptionsX509DescriptionPolicyIds, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenCertificateAuthorityCaCertificateDescriptionsX509DescriptionPolicyIds(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandCertificateAuthorityCaCertificateDescriptionsX509DescriptionPolicyIds expands an instance of CertificateAuthorityCaCertificateDescriptionsX509DescriptionPolicyIds into a JSON
// request object.
func expandCertificateAuthorityCaCertificateDescriptionsX509DescriptionPolicyIds(c *Client, f *CertificateAuthorityCaCertificateDescriptionsX509DescriptionPolicyIds, res *CertificateAuthority) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.ObjectIdPath; v != nil {
		m["objectIdPath"] = v
	}

	return m, nil
}

// flattenCertificateAuthorityCaCertificateDescriptionsX509DescriptionPolicyIds flattens an instance of CertificateAuthorityCaCertificateDescriptionsX509DescriptionPolicyIds from a JSON
// response object.
func flattenCertificateAuthorityCaCertificateDescriptionsX509DescriptionPolicyIds(c *Client, i interface{}, res *CertificateAuthority) *CertificateAuthorityCaCertificateDescriptionsX509DescriptionPolicyIds {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &CertificateAuthorityCaCertificateDescriptionsX509DescriptionPolicyIds{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyCertificateAuthorityCaCertificateDescriptionsX509DescriptionPolicyIds
	}
	r.ObjectIdPath = dcl.FlattenIntSlice(m["objectIdPath"])

	return r
}

// expandCertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensionsMap expands the contents of CertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensions into a JSON
// request object.
func expandCertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensionsMap(c *Client, f map[string]CertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensions, res *CertificateAuthority) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandCertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensions(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandCertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensionsSlice expands the contents of CertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensions into a JSON
// request object.
func expandCertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensionsSlice(c *Client, f []CertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensions, res *CertificateAuthority) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandCertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensions(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenCertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensionsMap flattens the contents of CertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensions from a JSON
// response object.
func flattenCertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensionsMap(c *Client, i interface{}, res *CertificateAuthority) map[string]CertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensions {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]CertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensions{}
	}

	if len(a) == 0 {
		return map[string]CertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensions{}
	}

	items := make(map[string]CertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensions)
	for k, item := range a {
		items[k] = *flattenCertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensions(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenCertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensionsSlice flattens the contents of CertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensions from a JSON
// response object.
func flattenCertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensionsSlice(c *Client, i interface{}, res *CertificateAuthority) []CertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensions {
	a, ok := i.([]interface{})
	if !ok {
		return []CertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensions{}
	}

	if len(a) == 0 {
		return []CertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensions{}
	}

	items := make([]CertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensions, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenCertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensions(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandCertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensions expands an instance of CertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensions into a JSON
// request object.
func expandCertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensions(c *Client, f *CertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensions, res *CertificateAuthority) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v, err := expandCertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensionsObjectId(c, f.ObjectId, res); err != nil {
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

// flattenCertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensions flattens an instance of CertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensions from a JSON
// response object.
func flattenCertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensions(c *Client, i interface{}, res *CertificateAuthority) *CertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensions {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &CertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensions{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyCertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensions
	}
	r.ObjectId = flattenCertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensionsObjectId(c, m["objectId"], res)
	r.Critical = dcl.FlattenBool(m["critical"])
	r.Value = dcl.FlattenString(m["value"])

	return r
}

// expandCertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensionsObjectIdMap expands the contents of CertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensionsObjectId into a JSON
// request object.
func expandCertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensionsObjectIdMap(c *Client, f map[string]CertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensionsObjectId, res *CertificateAuthority) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandCertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensionsObjectId(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandCertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensionsObjectIdSlice expands the contents of CertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensionsObjectId into a JSON
// request object.
func expandCertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensionsObjectIdSlice(c *Client, f []CertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensionsObjectId, res *CertificateAuthority) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandCertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensionsObjectId(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenCertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensionsObjectIdMap flattens the contents of CertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensionsObjectId from a JSON
// response object.
func flattenCertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensionsObjectIdMap(c *Client, i interface{}, res *CertificateAuthority) map[string]CertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensionsObjectId {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]CertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensionsObjectId{}
	}

	if len(a) == 0 {
		return map[string]CertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensionsObjectId{}
	}

	items := make(map[string]CertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensionsObjectId)
	for k, item := range a {
		items[k] = *flattenCertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensionsObjectId(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenCertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensionsObjectIdSlice flattens the contents of CertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensionsObjectId from a JSON
// response object.
func flattenCertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensionsObjectIdSlice(c *Client, i interface{}, res *CertificateAuthority) []CertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensionsObjectId {
	a, ok := i.([]interface{})
	if !ok {
		return []CertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensionsObjectId{}
	}

	if len(a) == 0 {
		return []CertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensionsObjectId{}
	}

	items := make([]CertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensionsObjectId, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenCertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensionsObjectId(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandCertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensionsObjectId expands an instance of CertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensionsObjectId into a JSON
// request object.
func expandCertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensionsObjectId(c *Client, f *CertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensionsObjectId, res *CertificateAuthority) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.ObjectIdPath; v != nil {
		m["objectIdPath"] = v
	}

	return m, nil
}

// flattenCertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensionsObjectId flattens an instance of CertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensionsObjectId from a JSON
// response object.
func flattenCertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensionsObjectId(c *Client, i interface{}, res *CertificateAuthority) *CertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensionsObjectId {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &CertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensionsObjectId{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyCertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensionsObjectId
	}
	r.ObjectIdPath = dcl.FlattenIntSlice(m["objectIdPath"])

	return r
}

// expandCertificateAuthorityCaCertificateDescriptionsPublicKeyMap expands the contents of CertificateAuthorityCaCertificateDescriptionsPublicKey into a JSON
// request object.
func expandCertificateAuthorityCaCertificateDescriptionsPublicKeyMap(c *Client, f map[string]CertificateAuthorityCaCertificateDescriptionsPublicKey, res *CertificateAuthority) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandCertificateAuthorityCaCertificateDescriptionsPublicKey(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandCertificateAuthorityCaCertificateDescriptionsPublicKeySlice expands the contents of CertificateAuthorityCaCertificateDescriptionsPublicKey into a JSON
// request object.
func expandCertificateAuthorityCaCertificateDescriptionsPublicKeySlice(c *Client, f []CertificateAuthorityCaCertificateDescriptionsPublicKey, res *CertificateAuthority) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandCertificateAuthorityCaCertificateDescriptionsPublicKey(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenCertificateAuthorityCaCertificateDescriptionsPublicKeyMap flattens the contents of CertificateAuthorityCaCertificateDescriptionsPublicKey from a JSON
// response object.
func flattenCertificateAuthorityCaCertificateDescriptionsPublicKeyMap(c *Client, i interface{}, res *CertificateAuthority) map[string]CertificateAuthorityCaCertificateDescriptionsPublicKey {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]CertificateAuthorityCaCertificateDescriptionsPublicKey{}
	}

	if len(a) == 0 {
		return map[string]CertificateAuthorityCaCertificateDescriptionsPublicKey{}
	}

	items := make(map[string]CertificateAuthorityCaCertificateDescriptionsPublicKey)
	for k, item := range a {
		items[k] = *flattenCertificateAuthorityCaCertificateDescriptionsPublicKey(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenCertificateAuthorityCaCertificateDescriptionsPublicKeySlice flattens the contents of CertificateAuthorityCaCertificateDescriptionsPublicKey from a JSON
// response object.
func flattenCertificateAuthorityCaCertificateDescriptionsPublicKeySlice(c *Client, i interface{}, res *CertificateAuthority) []CertificateAuthorityCaCertificateDescriptionsPublicKey {
	a, ok := i.([]interface{})
	if !ok {
		return []CertificateAuthorityCaCertificateDescriptionsPublicKey{}
	}

	if len(a) == 0 {
		return []CertificateAuthorityCaCertificateDescriptionsPublicKey{}
	}

	items := make([]CertificateAuthorityCaCertificateDescriptionsPublicKey, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenCertificateAuthorityCaCertificateDescriptionsPublicKey(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandCertificateAuthorityCaCertificateDescriptionsPublicKey expands an instance of CertificateAuthorityCaCertificateDescriptionsPublicKey into a JSON
// request object.
func expandCertificateAuthorityCaCertificateDescriptionsPublicKey(c *Client, f *CertificateAuthorityCaCertificateDescriptionsPublicKey, res *CertificateAuthority) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.Key; !dcl.IsEmptyValueIndirect(v) {
		m["key"] = v
	}
	if v := f.Format; !dcl.IsEmptyValueIndirect(v) {
		m["format"] = v
	}

	return m, nil
}

// flattenCertificateAuthorityCaCertificateDescriptionsPublicKey flattens an instance of CertificateAuthorityCaCertificateDescriptionsPublicKey from a JSON
// response object.
func flattenCertificateAuthorityCaCertificateDescriptionsPublicKey(c *Client, i interface{}, res *CertificateAuthority) *CertificateAuthorityCaCertificateDescriptionsPublicKey {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &CertificateAuthorityCaCertificateDescriptionsPublicKey{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyCertificateAuthorityCaCertificateDescriptionsPublicKey
	}
	r.Key = dcl.FlattenString(m["key"])
	r.Format = flattenCertificateAuthorityCaCertificateDescriptionsPublicKeyFormatEnum(m["format"])

	return r
}

// expandCertificateAuthorityCaCertificateDescriptionsSubjectKeyIdMap expands the contents of CertificateAuthorityCaCertificateDescriptionsSubjectKeyId into a JSON
// request object.
func expandCertificateAuthorityCaCertificateDescriptionsSubjectKeyIdMap(c *Client, f map[string]CertificateAuthorityCaCertificateDescriptionsSubjectKeyId, res *CertificateAuthority) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandCertificateAuthorityCaCertificateDescriptionsSubjectKeyId(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandCertificateAuthorityCaCertificateDescriptionsSubjectKeyIdSlice expands the contents of CertificateAuthorityCaCertificateDescriptionsSubjectKeyId into a JSON
// request object.
func expandCertificateAuthorityCaCertificateDescriptionsSubjectKeyIdSlice(c *Client, f []CertificateAuthorityCaCertificateDescriptionsSubjectKeyId, res *CertificateAuthority) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandCertificateAuthorityCaCertificateDescriptionsSubjectKeyId(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenCertificateAuthorityCaCertificateDescriptionsSubjectKeyIdMap flattens the contents of CertificateAuthorityCaCertificateDescriptionsSubjectKeyId from a JSON
// response object.
func flattenCertificateAuthorityCaCertificateDescriptionsSubjectKeyIdMap(c *Client, i interface{}, res *CertificateAuthority) map[string]CertificateAuthorityCaCertificateDescriptionsSubjectKeyId {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]CertificateAuthorityCaCertificateDescriptionsSubjectKeyId{}
	}

	if len(a) == 0 {
		return map[string]CertificateAuthorityCaCertificateDescriptionsSubjectKeyId{}
	}

	items := make(map[string]CertificateAuthorityCaCertificateDescriptionsSubjectKeyId)
	for k, item := range a {
		items[k] = *flattenCertificateAuthorityCaCertificateDescriptionsSubjectKeyId(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenCertificateAuthorityCaCertificateDescriptionsSubjectKeyIdSlice flattens the contents of CertificateAuthorityCaCertificateDescriptionsSubjectKeyId from a JSON
// response object.
func flattenCertificateAuthorityCaCertificateDescriptionsSubjectKeyIdSlice(c *Client, i interface{}, res *CertificateAuthority) []CertificateAuthorityCaCertificateDescriptionsSubjectKeyId {
	a, ok := i.([]interface{})
	if !ok {
		return []CertificateAuthorityCaCertificateDescriptionsSubjectKeyId{}
	}

	if len(a) == 0 {
		return []CertificateAuthorityCaCertificateDescriptionsSubjectKeyId{}
	}

	items := make([]CertificateAuthorityCaCertificateDescriptionsSubjectKeyId, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenCertificateAuthorityCaCertificateDescriptionsSubjectKeyId(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandCertificateAuthorityCaCertificateDescriptionsSubjectKeyId expands an instance of CertificateAuthorityCaCertificateDescriptionsSubjectKeyId into a JSON
// request object.
func expandCertificateAuthorityCaCertificateDescriptionsSubjectKeyId(c *Client, f *CertificateAuthorityCaCertificateDescriptionsSubjectKeyId, res *CertificateAuthority) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.KeyId; !dcl.IsEmptyValueIndirect(v) {
		m["keyId"] = v
	}

	return m, nil
}

// flattenCertificateAuthorityCaCertificateDescriptionsSubjectKeyId flattens an instance of CertificateAuthorityCaCertificateDescriptionsSubjectKeyId from a JSON
// response object.
func flattenCertificateAuthorityCaCertificateDescriptionsSubjectKeyId(c *Client, i interface{}, res *CertificateAuthority) *CertificateAuthorityCaCertificateDescriptionsSubjectKeyId {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &CertificateAuthorityCaCertificateDescriptionsSubjectKeyId{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyCertificateAuthorityCaCertificateDescriptionsSubjectKeyId
	}
	r.KeyId = dcl.FlattenString(m["keyId"])

	return r
}

// expandCertificateAuthorityCaCertificateDescriptionsAuthorityKeyIdMap expands the contents of CertificateAuthorityCaCertificateDescriptionsAuthorityKeyId into a JSON
// request object.
func expandCertificateAuthorityCaCertificateDescriptionsAuthorityKeyIdMap(c *Client, f map[string]CertificateAuthorityCaCertificateDescriptionsAuthorityKeyId, res *CertificateAuthority) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandCertificateAuthorityCaCertificateDescriptionsAuthorityKeyId(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandCertificateAuthorityCaCertificateDescriptionsAuthorityKeyIdSlice expands the contents of CertificateAuthorityCaCertificateDescriptionsAuthorityKeyId into a JSON
// request object.
func expandCertificateAuthorityCaCertificateDescriptionsAuthorityKeyIdSlice(c *Client, f []CertificateAuthorityCaCertificateDescriptionsAuthorityKeyId, res *CertificateAuthority) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandCertificateAuthorityCaCertificateDescriptionsAuthorityKeyId(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenCertificateAuthorityCaCertificateDescriptionsAuthorityKeyIdMap flattens the contents of CertificateAuthorityCaCertificateDescriptionsAuthorityKeyId from a JSON
// response object.
func flattenCertificateAuthorityCaCertificateDescriptionsAuthorityKeyIdMap(c *Client, i interface{}, res *CertificateAuthority) map[string]CertificateAuthorityCaCertificateDescriptionsAuthorityKeyId {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]CertificateAuthorityCaCertificateDescriptionsAuthorityKeyId{}
	}

	if len(a) == 0 {
		return map[string]CertificateAuthorityCaCertificateDescriptionsAuthorityKeyId{}
	}

	items := make(map[string]CertificateAuthorityCaCertificateDescriptionsAuthorityKeyId)
	for k, item := range a {
		items[k] = *flattenCertificateAuthorityCaCertificateDescriptionsAuthorityKeyId(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenCertificateAuthorityCaCertificateDescriptionsAuthorityKeyIdSlice flattens the contents of CertificateAuthorityCaCertificateDescriptionsAuthorityKeyId from a JSON
// response object.
func flattenCertificateAuthorityCaCertificateDescriptionsAuthorityKeyIdSlice(c *Client, i interface{}, res *CertificateAuthority) []CertificateAuthorityCaCertificateDescriptionsAuthorityKeyId {
	a, ok := i.([]interface{})
	if !ok {
		return []CertificateAuthorityCaCertificateDescriptionsAuthorityKeyId{}
	}

	if len(a) == 0 {
		return []CertificateAuthorityCaCertificateDescriptionsAuthorityKeyId{}
	}

	items := make([]CertificateAuthorityCaCertificateDescriptionsAuthorityKeyId, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenCertificateAuthorityCaCertificateDescriptionsAuthorityKeyId(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandCertificateAuthorityCaCertificateDescriptionsAuthorityKeyId expands an instance of CertificateAuthorityCaCertificateDescriptionsAuthorityKeyId into a JSON
// request object.
func expandCertificateAuthorityCaCertificateDescriptionsAuthorityKeyId(c *Client, f *CertificateAuthorityCaCertificateDescriptionsAuthorityKeyId, res *CertificateAuthority) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.KeyId; !dcl.IsEmptyValueIndirect(v) {
		m["keyId"] = v
	}

	return m, nil
}

// flattenCertificateAuthorityCaCertificateDescriptionsAuthorityKeyId flattens an instance of CertificateAuthorityCaCertificateDescriptionsAuthorityKeyId from a JSON
// response object.
func flattenCertificateAuthorityCaCertificateDescriptionsAuthorityKeyId(c *Client, i interface{}, res *CertificateAuthority) *CertificateAuthorityCaCertificateDescriptionsAuthorityKeyId {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &CertificateAuthorityCaCertificateDescriptionsAuthorityKeyId{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyCertificateAuthorityCaCertificateDescriptionsAuthorityKeyId
	}
	r.KeyId = dcl.FlattenString(m["keyId"])

	return r
}

// expandCertificateAuthorityCaCertificateDescriptionsCertFingerprintMap expands the contents of CertificateAuthorityCaCertificateDescriptionsCertFingerprint into a JSON
// request object.
func expandCertificateAuthorityCaCertificateDescriptionsCertFingerprintMap(c *Client, f map[string]CertificateAuthorityCaCertificateDescriptionsCertFingerprint, res *CertificateAuthority) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandCertificateAuthorityCaCertificateDescriptionsCertFingerprint(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandCertificateAuthorityCaCertificateDescriptionsCertFingerprintSlice expands the contents of CertificateAuthorityCaCertificateDescriptionsCertFingerprint into a JSON
// request object.
func expandCertificateAuthorityCaCertificateDescriptionsCertFingerprintSlice(c *Client, f []CertificateAuthorityCaCertificateDescriptionsCertFingerprint, res *CertificateAuthority) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandCertificateAuthorityCaCertificateDescriptionsCertFingerprint(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenCertificateAuthorityCaCertificateDescriptionsCertFingerprintMap flattens the contents of CertificateAuthorityCaCertificateDescriptionsCertFingerprint from a JSON
// response object.
func flattenCertificateAuthorityCaCertificateDescriptionsCertFingerprintMap(c *Client, i interface{}, res *CertificateAuthority) map[string]CertificateAuthorityCaCertificateDescriptionsCertFingerprint {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]CertificateAuthorityCaCertificateDescriptionsCertFingerprint{}
	}

	if len(a) == 0 {
		return map[string]CertificateAuthorityCaCertificateDescriptionsCertFingerprint{}
	}

	items := make(map[string]CertificateAuthorityCaCertificateDescriptionsCertFingerprint)
	for k, item := range a {
		items[k] = *flattenCertificateAuthorityCaCertificateDescriptionsCertFingerprint(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenCertificateAuthorityCaCertificateDescriptionsCertFingerprintSlice flattens the contents of CertificateAuthorityCaCertificateDescriptionsCertFingerprint from a JSON
// response object.
func flattenCertificateAuthorityCaCertificateDescriptionsCertFingerprintSlice(c *Client, i interface{}, res *CertificateAuthority) []CertificateAuthorityCaCertificateDescriptionsCertFingerprint {
	a, ok := i.([]interface{})
	if !ok {
		return []CertificateAuthorityCaCertificateDescriptionsCertFingerprint{}
	}

	if len(a) == 0 {
		return []CertificateAuthorityCaCertificateDescriptionsCertFingerprint{}
	}

	items := make([]CertificateAuthorityCaCertificateDescriptionsCertFingerprint, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenCertificateAuthorityCaCertificateDescriptionsCertFingerprint(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandCertificateAuthorityCaCertificateDescriptionsCertFingerprint expands an instance of CertificateAuthorityCaCertificateDescriptionsCertFingerprint into a JSON
// request object.
func expandCertificateAuthorityCaCertificateDescriptionsCertFingerprint(c *Client, f *CertificateAuthorityCaCertificateDescriptionsCertFingerprint, res *CertificateAuthority) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.Sha256Hash; !dcl.IsEmptyValueIndirect(v) {
		m["sha256Hash"] = v
	}

	return m, nil
}

// flattenCertificateAuthorityCaCertificateDescriptionsCertFingerprint flattens an instance of CertificateAuthorityCaCertificateDescriptionsCertFingerprint from a JSON
// response object.
func flattenCertificateAuthorityCaCertificateDescriptionsCertFingerprint(c *Client, i interface{}, res *CertificateAuthority) *CertificateAuthorityCaCertificateDescriptionsCertFingerprint {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &CertificateAuthorityCaCertificateDescriptionsCertFingerprint{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyCertificateAuthorityCaCertificateDescriptionsCertFingerprint
	}
	r.Sha256Hash = dcl.FlattenString(m["sha256Hash"])

	return r
}

// expandCertificateAuthorityAccessUrlsMap expands the contents of CertificateAuthorityAccessUrls into a JSON
// request object.
func expandCertificateAuthorityAccessUrlsMap(c *Client, f map[string]CertificateAuthorityAccessUrls, res *CertificateAuthority) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandCertificateAuthorityAccessUrls(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandCertificateAuthorityAccessUrlsSlice expands the contents of CertificateAuthorityAccessUrls into a JSON
// request object.
func expandCertificateAuthorityAccessUrlsSlice(c *Client, f []CertificateAuthorityAccessUrls, res *CertificateAuthority) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandCertificateAuthorityAccessUrls(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenCertificateAuthorityAccessUrlsMap flattens the contents of CertificateAuthorityAccessUrls from a JSON
// response object.
func flattenCertificateAuthorityAccessUrlsMap(c *Client, i interface{}, res *CertificateAuthority) map[string]CertificateAuthorityAccessUrls {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]CertificateAuthorityAccessUrls{}
	}

	if len(a) == 0 {
		return map[string]CertificateAuthorityAccessUrls{}
	}

	items := make(map[string]CertificateAuthorityAccessUrls)
	for k, item := range a {
		items[k] = *flattenCertificateAuthorityAccessUrls(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenCertificateAuthorityAccessUrlsSlice flattens the contents of CertificateAuthorityAccessUrls from a JSON
// response object.
func flattenCertificateAuthorityAccessUrlsSlice(c *Client, i interface{}, res *CertificateAuthority) []CertificateAuthorityAccessUrls {
	a, ok := i.([]interface{})
	if !ok {
		return []CertificateAuthorityAccessUrls{}
	}

	if len(a) == 0 {
		return []CertificateAuthorityAccessUrls{}
	}

	items := make([]CertificateAuthorityAccessUrls, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenCertificateAuthorityAccessUrls(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandCertificateAuthorityAccessUrls expands an instance of CertificateAuthorityAccessUrls into a JSON
// request object.
func expandCertificateAuthorityAccessUrls(c *Client, f *CertificateAuthorityAccessUrls, res *CertificateAuthority) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.CaCertificateAccessUrl; !dcl.IsEmptyValueIndirect(v) {
		m["caCertificateAccessUrl"] = v
	}
	if v := f.CrlAccessUrls; v != nil {
		m["crlAccessUrls"] = v
	}

	return m, nil
}

// flattenCertificateAuthorityAccessUrls flattens an instance of CertificateAuthorityAccessUrls from a JSON
// response object.
func flattenCertificateAuthorityAccessUrls(c *Client, i interface{}, res *CertificateAuthority) *CertificateAuthorityAccessUrls {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &CertificateAuthorityAccessUrls{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyCertificateAuthorityAccessUrls
	}
	r.CaCertificateAccessUrl = dcl.FlattenString(m["caCertificateAccessUrl"])
	r.CrlAccessUrls = dcl.FlattenStringSlice(m["crlAccessUrls"])

	return r
}

// flattenCertificateAuthorityTypeEnumMap flattens the contents of CertificateAuthorityTypeEnum from a JSON
// response object.
func flattenCertificateAuthorityTypeEnumMap(c *Client, i interface{}, res *CertificateAuthority) map[string]CertificateAuthorityTypeEnum {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]CertificateAuthorityTypeEnum{}
	}

	if len(a) == 0 {
		return map[string]CertificateAuthorityTypeEnum{}
	}

	items := make(map[string]CertificateAuthorityTypeEnum)
	for k, item := range a {
		items[k] = *flattenCertificateAuthorityTypeEnum(item.(interface{}))
	}

	return items
}

// flattenCertificateAuthorityTypeEnumSlice flattens the contents of CertificateAuthorityTypeEnum from a JSON
// response object.
func flattenCertificateAuthorityTypeEnumSlice(c *Client, i interface{}, res *CertificateAuthority) []CertificateAuthorityTypeEnum {
	a, ok := i.([]interface{})
	if !ok {
		return []CertificateAuthorityTypeEnum{}
	}

	if len(a) == 0 {
		return []CertificateAuthorityTypeEnum{}
	}

	items := make([]CertificateAuthorityTypeEnum, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenCertificateAuthorityTypeEnum(item.(interface{})))
	}

	return items
}

// flattenCertificateAuthorityTypeEnum asserts that an interface is a string, and returns a
// pointer to a *CertificateAuthorityTypeEnum with the same value as that string.
func flattenCertificateAuthorityTypeEnum(i interface{}) *CertificateAuthorityTypeEnum {
	s, ok := i.(string)
	if !ok {
		return nil
	}

	return CertificateAuthorityTypeEnumRef(s)
}

// flattenCertificateAuthorityConfigPublicKeyFormatEnumMap flattens the contents of CertificateAuthorityConfigPublicKeyFormatEnum from a JSON
// response object.
func flattenCertificateAuthorityConfigPublicKeyFormatEnumMap(c *Client, i interface{}, res *CertificateAuthority) map[string]CertificateAuthorityConfigPublicKeyFormatEnum {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]CertificateAuthorityConfigPublicKeyFormatEnum{}
	}

	if len(a) == 0 {
		return map[string]CertificateAuthorityConfigPublicKeyFormatEnum{}
	}

	items := make(map[string]CertificateAuthorityConfigPublicKeyFormatEnum)
	for k, item := range a {
		items[k] = *flattenCertificateAuthorityConfigPublicKeyFormatEnum(item.(interface{}))
	}

	return items
}

// flattenCertificateAuthorityConfigPublicKeyFormatEnumSlice flattens the contents of CertificateAuthorityConfigPublicKeyFormatEnum from a JSON
// response object.
func flattenCertificateAuthorityConfigPublicKeyFormatEnumSlice(c *Client, i interface{}, res *CertificateAuthority) []CertificateAuthorityConfigPublicKeyFormatEnum {
	a, ok := i.([]interface{})
	if !ok {
		return []CertificateAuthorityConfigPublicKeyFormatEnum{}
	}

	if len(a) == 0 {
		return []CertificateAuthorityConfigPublicKeyFormatEnum{}
	}

	items := make([]CertificateAuthorityConfigPublicKeyFormatEnum, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenCertificateAuthorityConfigPublicKeyFormatEnum(item.(interface{})))
	}

	return items
}

// flattenCertificateAuthorityConfigPublicKeyFormatEnum asserts that an interface is a string, and returns a
// pointer to a *CertificateAuthorityConfigPublicKeyFormatEnum with the same value as that string.
func flattenCertificateAuthorityConfigPublicKeyFormatEnum(i interface{}) *CertificateAuthorityConfigPublicKeyFormatEnum {
	s, ok := i.(string)
	if !ok {
		return nil
	}

	return CertificateAuthorityConfigPublicKeyFormatEnumRef(s)
}

// flattenCertificateAuthorityKeySpecAlgorithmEnumMap flattens the contents of CertificateAuthorityKeySpecAlgorithmEnum from a JSON
// response object.
func flattenCertificateAuthorityKeySpecAlgorithmEnumMap(c *Client, i interface{}, res *CertificateAuthority) map[string]CertificateAuthorityKeySpecAlgorithmEnum {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]CertificateAuthorityKeySpecAlgorithmEnum{}
	}

	if len(a) == 0 {
		return map[string]CertificateAuthorityKeySpecAlgorithmEnum{}
	}

	items := make(map[string]CertificateAuthorityKeySpecAlgorithmEnum)
	for k, item := range a {
		items[k] = *flattenCertificateAuthorityKeySpecAlgorithmEnum(item.(interface{}))
	}

	return items
}

// flattenCertificateAuthorityKeySpecAlgorithmEnumSlice flattens the contents of CertificateAuthorityKeySpecAlgorithmEnum from a JSON
// response object.
func flattenCertificateAuthorityKeySpecAlgorithmEnumSlice(c *Client, i interface{}, res *CertificateAuthority) []CertificateAuthorityKeySpecAlgorithmEnum {
	a, ok := i.([]interface{})
	if !ok {
		return []CertificateAuthorityKeySpecAlgorithmEnum{}
	}

	if len(a) == 0 {
		return []CertificateAuthorityKeySpecAlgorithmEnum{}
	}

	items := make([]CertificateAuthorityKeySpecAlgorithmEnum, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenCertificateAuthorityKeySpecAlgorithmEnum(item.(interface{})))
	}

	return items
}

// flattenCertificateAuthorityKeySpecAlgorithmEnum asserts that an interface is a string, and returns a
// pointer to a *CertificateAuthorityKeySpecAlgorithmEnum with the same value as that string.
func flattenCertificateAuthorityKeySpecAlgorithmEnum(i interface{}) *CertificateAuthorityKeySpecAlgorithmEnum {
	s, ok := i.(string)
	if !ok {
		return nil
	}

	return CertificateAuthorityKeySpecAlgorithmEnumRef(s)
}

// flattenCertificateAuthorityTierEnumMap flattens the contents of CertificateAuthorityTierEnum from a JSON
// response object.
func flattenCertificateAuthorityTierEnumMap(c *Client, i interface{}, res *CertificateAuthority) map[string]CertificateAuthorityTierEnum {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]CertificateAuthorityTierEnum{}
	}

	if len(a) == 0 {
		return map[string]CertificateAuthorityTierEnum{}
	}

	items := make(map[string]CertificateAuthorityTierEnum)
	for k, item := range a {
		items[k] = *flattenCertificateAuthorityTierEnum(item.(interface{}))
	}

	return items
}

// flattenCertificateAuthorityTierEnumSlice flattens the contents of CertificateAuthorityTierEnum from a JSON
// response object.
func flattenCertificateAuthorityTierEnumSlice(c *Client, i interface{}, res *CertificateAuthority) []CertificateAuthorityTierEnum {
	a, ok := i.([]interface{})
	if !ok {
		return []CertificateAuthorityTierEnum{}
	}

	if len(a) == 0 {
		return []CertificateAuthorityTierEnum{}
	}

	items := make([]CertificateAuthorityTierEnum, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenCertificateAuthorityTierEnum(item.(interface{})))
	}

	return items
}

// flattenCertificateAuthorityTierEnum asserts that an interface is a string, and returns a
// pointer to a *CertificateAuthorityTierEnum with the same value as that string.
func flattenCertificateAuthorityTierEnum(i interface{}) *CertificateAuthorityTierEnum {
	s, ok := i.(string)
	if !ok {
		return nil
	}

	return CertificateAuthorityTierEnumRef(s)
}

// flattenCertificateAuthorityStateEnumMap flattens the contents of CertificateAuthorityStateEnum from a JSON
// response object.
func flattenCertificateAuthorityStateEnumMap(c *Client, i interface{}, res *CertificateAuthority) map[string]CertificateAuthorityStateEnum {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]CertificateAuthorityStateEnum{}
	}

	if len(a) == 0 {
		return map[string]CertificateAuthorityStateEnum{}
	}

	items := make(map[string]CertificateAuthorityStateEnum)
	for k, item := range a {
		items[k] = *flattenCertificateAuthorityStateEnum(item.(interface{}))
	}

	return items
}

// flattenCertificateAuthorityStateEnumSlice flattens the contents of CertificateAuthorityStateEnum from a JSON
// response object.
func flattenCertificateAuthorityStateEnumSlice(c *Client, i interface{}, res *CertificateAuthority) []CertificateAuthorityStateEnum {
	a, ok := i.([]interface{})
	if !ok {
		return []CertificateAuthorityStateEnum{}
	}

	if len(a) == 0 {
		return []CertificateAuthorityStateEnum{}
	}

	items := make([]CertificateAuthorityStateEnum, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenCertificateAuthorityStateEnum(item.(interface{})))
	}

	return items
}

// flattenCertificateAuthorityStateEnum asserts that an interface is a string, and returns a
// pointer to a *CertificateAuthorityStateEnum with the same value as that string.
func flattenCertificateAuthorityStateEnum(i interface{}) *CertificateAuthorityStateEnum {
	s, ok := i.(string)
	if !ok {
		return nil
	}

	return CertificateAuthorityStateEnumRef(s)
}

// flattenCertificateAuthorityCaCertificateDescriptionsPublicKeyFormatEnumMap flattens the contents of CertificateAuthorityCaCertificateDescriptionsPublicKeyFormatEnum from a JSON
// response object.
func flattenCertificateAuthorityCaCertificateDescriptionsPublicKeyFormatEnumMap(c *Client, i interface{}, res *CertificateAuthority) map[string]CertificateAuthorityCaCertificateDescriptionsPublicKeyFormatEnum {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]CertificateAuthorityCaCertificateDescriptionsPublicKeyFormatEnum{}
	}

	if len(a) == 0 {
		return map[string]CertificateAuthorityCaCertificateDescriptionsPublicKeyFormatEnum{}
	}

	items := make(map[string]CertificateAuthorityCaCertificateDescriptionsPublicKeyFormatEnum)
	for k, item := range a {
		items[k] = *flattenCertificateAuthorityCaCertificateDescriptionsPublicKeyFormatEnum(item.(interface{}))
	}

	return items
}

// flattenCertificateAuthorityCaCertificateDescriptionsPublicKeyFormatEnumSlice flattens the contents of CertificateAuthorityCaCertificateDescriptionsPublicKeyFormatEnum from a JSON
// response object.
func flattenCertificateAuthorityCaCertificateDescriptionsPublicKeyFormatEnumSlice(c *Client, i interface{}, res *CertificateAuthority) []CertificateAuthorityCaCertificateDescriptionsPublicKeyFormatEnum {
	a, ok := i.([]interface{})
	if !ok {
		return []CertificateAuthorityCaCertificateDescriptionsPublicKeyFormatEnum{}
	}

	if len(a) == 0 {
		return []CertificateAuthorityCaCertificateDescriptionsPublicKeyFormatEnum{}
	}

	items := make([]CertificateAuthorityCaCertificateDescriptionsPublicKeyFormatEnum, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenCertificateAuthorityCaCertificateDescriptionsPublicKeyFormatEnum(item.(interface{})))
	}

	return items
}

// flattenCertificateAuthorityCaCertificateDescriptionsPublicKeyFormatEnum asserts that an interface is a string, and returns a
// pointer to a *CertificateAuthorityCaCertificateDescriptionsPublicKeyFormatEnum with the same value as that string.
func flattenCertificateAuthorityCaCertificateDescriptionsPublicKeyFormatEnum(i interface{}) *CertificateAuthorityCaCertificateDescriptionsPublicKeyFormatEnum {
	s, ok := i.(string)
	if !ok {
		return nil
	}

	return CertificateAuthorityCaCertificateDescriptionsPublicKeyFormatEnumRef(s)
}

// This function returns a matcher that checks whether a serialized resource matches this resource
// in its parameters (as defined by the fields in a Get, which definitionally define resource
// identity).  This is useful in extracting the element from a List call.
func (r *CertificateAuthority) matcher(c *Client) func([]byte) bool {
	return func(b []byte) bool {
		cr, err := unmarshalCertificateAuthority(b, c, r)
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
		if nr.CaPool == nil && ncr.CaPool == nil {
			c.Config.Logger.Info("Both CaPool fields null - considering equal.")
		} else if nr.CaPool == nil || ncr.CaPool == nil {
			c.Config.Logger.Info("Only one CaPool field is null - considering unequal.")
			return false
		} else if *nr.CaPool != *ncr.CaPool {
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

type certificateAuthorityDiff struct {
	// The diff should include one or the other of RequiresRecreate or UpdateOp.
	RequiresRecreate bool
	UpdateOp         certificateAuthorityApiOperation
	FieldName        string // used for error logging
}

func convertFieldDiffsToCertificateAuthorityDiffs(config *dcl.Config, fds []*dcl.FieldDiff, opts []dcl.ApplyOption) ([]certificateAuthorityDiff, error) {
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
	var diffs []certificateAuthorityDiff
	// For each operation name, create a certificateAuthorityDiff which contains the operation.
	for opName, fieldDiffs := range opNamesToFieldDiffs {
		// Use the first field diff's field name for logging required recreate error.
		diff := certificateAuthorityDiff{FieldName: fieldDiffs[0].FieldName}
		if opName == "Recreate" {
			diff.RequiresRecreate = true
		} else {
			apiOp, err := convertOpNameToCertificateAuthorityApiOperation(opName, fieldDiffs, opts...)
			if err != nil {
				return diffs, err
			}
			diff.UpdateOp = apiOp
		}
		diffs = append(diffs, diff)
	}
	return diffs, nil
}

func convertOpNameToCertificateAuthorityApiOperation(opName string, fieldDiffs []*dcl.FieldDiff, opts ...dcl.ApplyOption) (certificateAuthorityApiOperation, error) {
	switch opName {

	case "updateCertificateAuthorityUpdateCertificateAuthorityOperation":
		return &updateCertificateAuthorityUpdateCertificateAuthorityOperation{FieldDiffs: fieldDiffs}, nil

	default:
		return nil, fmt.Errorf("no such operation with name: %v", opName)
	}
}

func extractCertificateAuthorityFields(r *CertificateAuthority) error {
	vConfig := r.Config
	if vConfig == nil {
		// note: explicitly not the empty object.
		vConfig = &CertificateAuthorityConfig{}
	}
	if err := extractCertificateAuthorityConfigFields(r, vConfig); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vConfig) {
		r.Config = vConfig
	}
	vKeySpec := r.KeySpec
	if vKeySpec == nil {
		// note: explicitly not the empty object.
		vKeySpec = &CertificateAuthorityKeySpec{}
	}
	if err := extractCertificateAuthorityKeySpecFields(r, vKeySpec); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vKeySpec) {
		r.KeySpec = vKeySpec
	}
	vSubordinateConfig := r.SubordinateConfig
	if vSubordinateConfig == nil {
		// note: explicitly not the empty object.
		vSubordinateConfig = &CertificateAuthoritySubordinateConfig{}
	}
	if err := extractCertificateAuthoritySubordinateConfigFields(r, vSubordinateConfig); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vSubordinateConfig) {
		r.SubordinateConfig = vSubordinateConfig
	}
	vAccessUrls := r.AccessUrls
	if vAccessUrls == nil {
		// note: explicitly not the empty object.
		vAccessUrls = &CertificateAuthorityAccessUrls{}
	}
	if err := extractCertificateAuthorityAccessUrlsFields(r, vAccessUrls); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vAccessUrls) {
		r.AccessUrls = vAccessUrls
	}
	return nil
}
func extractCertificateAuthorityConfigFields(r *CertificateAuthority, o *CertificateAuthorityConfig) error {
	vSubjectConfig := o.SubjectConfig
	if vSubjectConfig == nil {
		// note: explicitly not the empty object.
		vSubjectConfig = &CertificateAuthorityConfigSubjectConfig{}
	}
	if err := extractCertificateAuthorityConfigSubjectConfigFields(r, vSubjectConfig); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vSubjectConfig) {
		o.SubjectConfig = vSubjectConfig
	}
	vX509Config := o.X509Config
	if vX509Config == nil {
		// note: explicitly not the empty object.
		vX509Config = &CertificateAuthorityConfigX509Config{}
	}
	if err := extractCertificateAuthorityConfigX509ConfigFields(r, vX509Config); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vX509Config) {
		o.X509Config = vX509Config
	}
	vPublicKey := o.PublicKey
	if vPublicKey == nil {
		// note: explicitly not the empty object.
		vPublicKey = &CertificateAuthorityConfigPublicKey{}
	}
	if err := extractCertificateAuthorityConfigPublicKeyFields(r, vPublicKey); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vPublicKey) {
		o.PublicKey = vPublicKey
	}
	return nil
}
func extractCertificateAuthorityConfigSubjectConfigFields(r *CertificateAuthority, o *CertificateAuthorityConfigSubjectConfig) error {
	vSubject := o.Subject
	if vSubject == nil {
		// note: explicitly not the empty object.
		vSubject = &CertificateAuthorityConfigSubjectConfigSubject{}
	}
	if err := extractCertificateAuthorityConfigSubjectConfigSubjectFields(r, vSubject); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vSubject) {
		o.Subject = vSubject
	}
	vSubjectAltName := o.SubjectAltName
	if vSubjectAltName == nil {
		// note: explicitly not the empty object.
		vSubjectAltName = &CertificateAuthorityConfigSubjectConfigSubjectAltName{}
	}
	if err := extractCertificateAuthorityConfigSubjectConfigSubjectAltNameFields(r, vSubjectAltName); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vSubjectAltName) {
		o.SubjectAltName = vSubjectAltName
	}
	return nil
}
func extractCertificateAuthorityConfigSubjectConfigSubjectFields(r *CertificateAuthority, o *CertificateAuthorityConfigSubjectConfigSubject) error {
	return nil
}
func extractCertificateAuthorityConfigSubjectConfigSubjectAltNameFields(r *CertificateAuthority, o *CertificateAuthorityConfigSubjectConfigSubjectAltName) error {
	return nil
}
func extractCertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSansFields(r *CertificateAuthority, o *CertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSans) error {
	vObjectId := o.ObjectId
	if vObjectId == nil {
		// note: explicitly not the empty object.
		vObjectId = &CertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSansObjectId{}
	}
	if err := extractCertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSansObjectIdFields(r, vObjectId); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vObjectId) {
		o.ObjectId = vObjectId
	}
	return nil
}
func extractCertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSansObjectIdFields(r *CertificateAuthority, o *CertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSansObjectId) error {
	return nil
}
func extractCertificateAuthorityConfigX509ConfigFields(r *CertificateAuthority, o *CertificateAuthorityConfigX509Config) error {
	vKeyUsage := o.KeyUsage
	if vKeyUsage == nil {
		// note: explicitly not the empty object.
		vKeyUsage = &CertificateAuthorityConfigX509ConfigKeyUsage{}
	}
	if err := extractCertificateAuthorityConfigX509ConfigKeyUsageFields(r, vKeyUsage); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vKeyUsage) {
		o.KeyUsage = vKeyUsage
	}
	vCaOptions := o.CaOptions
	if vCaOptions == nil {
		// note: explicitly not the empty object.
		vCaOptions = &CertificateAuthorityConfigX509ConfigCaOptions{}
	}
	if err := extractCertificateAuthorityConfigX509ConfigCaOptionsFields(r, vCaOptions); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vCaOptions) {
		o.CaOptions = vCaOptions
	}
	return nil
}
func extractCertificateAuthorityConfigX509ConfigKeyUsageFields(r *CertificateAuthority, o *CertificateAuthorityConfigX509ConfigKeyUsage) error {
	vBaseKeyUsage := o.BaseKeyUsage
	if vBaseKeyUsage == nil {
		// note: explicitly not the empty object.
		vBaseKeyUsage = &CertificateAuthorityConfigX509ConfigKeyUsageBaseKeyUsage{}
	}
	if err := extractCertificateAuthorityConfigX509ConfigKeyUsageBaseKeyUsageFields(r, vBaseKeyUsage); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vBaseKeyUsage) {
		o.BaseKeyUsage = vBaseKeyUsage
	}
	vExtendedKeyUsage := o.ExtendedKeyUsage
	if vExtendedKeyUsage == nil {
		// note: explicitly not the empty object.
		vExtendedKeyUsage = &CertificateAuthorityConfigX509ConfigKeyUsageExtendedKeyUsage{}
	}
	if err := extractCertificateAuthorityConfigX509ConfigKeyUsageExtendedKeyUsageFields(r, vExtendedKeyUsage); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vExtendedKeyUsage) {
		o.ExtendedKeyUsage = vExtendedKeyUsage
	}
	return nil
}
func extractCertificateAuthorityConfigX509ConfigKeyUsageBaseKeyUsageFields(r *CertificateAuthority, o *CertificateAuthorityConfigX509ConfigKeyUsageBaseKeyUsage) error {
	return nil
}
func extractCertificateAuthorityConfigX509ConfigKeyUsageExtendedKeyUsageFields(r *CertificateAuthority, o *CertificateAuthorityConfigX509ConfigKeyUsageExtendedKeyUsage) error {
	return nil
}
func extractCertificateAuthorityConfigX509ConfigKeyUsageUnknownExtendedKeyUsagesFields(r *CertificateAuthority, o *CertificateAuthorityConfigX509ConfigKeyUsageUnknownExtendedKeyUsages) error {
	return nil
}
func extractCertificateAuthorityConfigX509ConfigCaOptionsFields(r *CertificateAuthority, o *CertificateAuthorityConfigX509ConfigCaOptions) error {
	return nil
}
func extractCertificateAuthorityConfigX509ConfigPolicyIdsFields(r *CertificateAuthority, o *CertificateAuthorityConfigX509ConfigPolicyIds) error {
	return nil
}
func extractCertificateAuthorityConfigX509ConfigAdditionalExtensionsFields(r *CertificateAuthority, o *CertificateAuthorityConfigX509ConfigAdditionalExtensions) error {
	vObjectId := o.ObjectId
	if vObjectId == nil {
		// note: explicitly not the empty object.
		vObjectId = &CertificateAuthorityConfigX509ConfigAdditionalExtensionsObjectId{}
	}
	if err := extractCertificateAuthorityConfigX509ConfigAdditionalExtensionsObjectIdFields(r, vObjectId); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vObjectId) {
		o.ObjectId = vObjectId
	}
	return nil
}
func extractCertificateAuthorityConfigX509ConfigAdditionalExtensionsObjectIdFields(r *CertificateAuthority, o *CertificateAuthorityConfigX509ConfigAdditionalExtensionsObjectId) error {
	return nil
}
func extractCertificateAuthorityConfigPublicKeyFields(r *CertificateAuthority, o *CertificateAuthorityConfigPublicKey) error {
	return nil
}
func extractCertificateAuthorityKeySpecFields(r *CertificateAuthority, o *CertificateAuthorityKeySpec) error {
	return nil
}
func extractCertificateAuthoritySubordinateConfigFields(r *CertificateAuthority, o *CertificateAuthoritySubordinateConfig) error {
	vPemIssuerChain := o.PemIssuerChain
	if vPemIssuerChain == nil {
		// note: explicitly not the empty object.
		vPemIssuerChain = &CertificateAuthoritySubordinateConfigPemIssuerChain{}
	}
	if err := extractCertificateAuthoritySubordinateConfigPemIssuerChainFields(r, vPemIssuerChain); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vPemIssuerChain) {
		o.PemIssuerChain = vPemIssuerChain
	}
	return nil
}
func extractCertificateAuthoritySubordinateConfigPemIssuerChainFields(r *CertificateAuthority, o *CertificateAuthoritySubordinateConfigPemIssuerChain) error {
	return nil
}
func extractCertificateAuthorityCaCertificateDescriptionsFields(r *CertificateAuthority, o *CertificateAuthorityCaCertificateDescriptions) error {
	vSubjectDescription := o.SubjectDescription
	if vSubjectDescription == nil {
		// note: explicitly not the empty object.
		vSubjectDescription = &CertificateAuthorityCaCertificateDescriptionsSubjectDescription{}
	}
	if err := extractCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionFields(r, vSubjectDescription); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vSubjectDescription) {
		o.SubjectDescription = vSubjectDescription
	}
	vX509Description := o.X509Description
	if vX509Description == nil {
		// note: explicitly not the empty object.
		vX509Description = &CertificateAuthorityCaCertificateDescriptionsX509Description{}
	}
	if err := extractCertificateAuthorityCaCertificateDescriptionsX509DescriptionFields(r, vX509Description); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vX509Description) {
		o.X509Description = vX509Description
	}
	vPublicKey := o.PublicKey
	if vPublicKey == nil {
		// note: explicitly not the empty object.
		vPublicKey = &CertificateAuthorityCaCertificateDescriptionsPublicKey{}
	}
	if err := extractCertificateAuthorityCaCertificateDescriptionsPublicKeyFields(r, vPublicKey); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vPublicKey) {
		o.PublicKey = vPublicKey
	}
	vSubjectKeyId := o.SubjectKeyId
	if vSubjectKeyId == nil {
		// note: explicitly not the empty object.
		vSubjectKeyId = &CertificateAuthorityCaCertificateDescriptionsSubjectKeyId{}
	}
	if err := extractCertificateAuthorityCaCertificateDescriptionsSubjectKeyIdFields(r, vSubjectKeyId); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vSubjectKeyId) {
		o.SubjectKeyId = vSubjectKeyId
	}
	vAuthorityKeyId := o.AuthorityKeyId
	if vAuthorityKeyId == nil {
		// note: explicitly not the empty object.
		vAuthorityKeyId = &CertificateAuthorityCaCertificateDescriptionsAuthorityKeyId{}
	}
	if err := extractCertificateAuthorityCaCertificateDescriptionsAuthorityKeyIdFields(r, vAuthorityKeyId); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vAuthorityKeyId) {
		o.AuthorityKeyId = vAuthorityKeyId
	}
	vCertFingerprint := o.CertFingerprint
	if vCertFingerprint == nil {
		// note: explicitly not the empty object.
		vCertFingerprint = &CertificateAuthorityCaCertificateDescriptionsCertFingerprint{}
	}
	if err := extractCertificateAuthorityCaCertificateDescriptionsCertFingerprintFields(r, vCertFingerprint); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vCertFingerprint) {
		o.CertFingerprint = vCertFingerprint
	}
	return nil
}
func extractCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionFields(r *CertificateAuthority, o *CertificateAuthorityCaCertificateDescriptionsSubjectDescription) error {
	vSubject := o.Subject
	if vSubject == nil {
		// note: explicitly not the empty object.
		vSubject = &CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubject{}
	}
	if err := extractCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectFields(r, vSubject); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vSubject) {
		o.Subject = vSubject
	}
	vSubjectAltName := o.SubjectAltName
	if vSubjectAltName == nil {
		// note: explicitly not the empty object.
		vSubjectAltName = &CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltName{}
	}
	if err := extractCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameFields(r, vSubjectAltName); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vSubjectAltName) {
		o.SubjectAltName = vSubjectAltName
	}
	return nil
}
func extractCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectFields(r *CertificateAuthority, o *CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubject) error {
	return nil
}
func extractCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameFields(r *CertificateAuthority, o *CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltName) error {
	return nil
}
func extractCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSansFields(r *CertificateAuthority, o *CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSans) error {
	vObjectId := o.ObjectId
	if vObjectId == nil {
		// note: explicitly not the empty object.
		vObjectId = &CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSansObjectId{}
	}
	if err := extractCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSansObjectIdFields(r, vObjectId); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vObjectId) {
		o.ObjectId = vObjectId
	}
	return nil
}
func extractCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSansObjectIdFields(r *CertificateAuthority, o *CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSansObjectId) error {
	return nil
}
func extractCertificateAuthorityCaCertificateDescriptionsX509DescriptionFields(r *CertificateAuthority, o *CertificateAuthorityCaCertificateDescriptionsX509Description) error {
	vKeyUsage := o.KeyUsage
	if vKeyUsage == nil {
		// note: explicitly not the empty object.
		vKeyUsage = &CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsage{}
	}
	if err := extractCertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageFields(r, vKeyUsage); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vKeyUsage) {
		o.KeyUsage = vKeyUsage
	}
	vCaOptions := o.CaOptions
	if vCaOptions == nil {
		// note: explicitly not the empty object.
		vCaOptions = &CertificateAuthorityCaCertificateDescriptionsX509DescriptionCaOptions{}
	}
	if err := extractCertificateAuthorityCaCertificateDescriptionsX509DescriptionCaOptionsFields(r, vCaOptions); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vCaOptions) {
		o.CaOptions = vCaOptions
	}
	return nil
}
func extractCertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageFields(r *CertificateAuthority, o *CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsage) error {
	vBaseKeyUsage := o.BaseKeyUsage
	if vBaseKeyUsage == nil {
		// note: explicitly not the empty object.
		vBaseKeyUsage = &CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageBaseKeyUsage{}
	}
	if err := extractCertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageBaseKeyUsageFields(r, vBaseKeyUsage); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vBaseKeyUsage) {
		o.BaseKeyUsage = vBaseKeyUsage
	}
	vExtendedKeyUsage := o.ExtendedKeyUsage
	if vExtendedKeyUsage == nil {
		// note: explicitly not the empty object.
		vExtendedKeyUsage = &CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageExtendedKeyUsage{}
	}
	if err := extractCertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageExtendedKeyUsageFields(r, vExtendedKeyUsage); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vExtendedKeyUsage) {
		o.ExtendedKeyUsage = vExtendedKeyUsage
	}
	return nil
}
func extractCertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageBaseKeyUsageFields(r *CertificateAuthority, o *CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageBaseKeyUsage) error {
	return nil
}
func extractCertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageExtendedKeyUsageFields(r *CertificateAuthority, o *CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageExtendedKeyUsage) error {
	return nil
}
func extractCertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageUnknownExtendedKeyUsagesFields(r *CertificateAuthority, o *CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageUnknownExtendedKeyUsages) error {
	return nil
}
func extractCertificateAuthorityCaCertificateDescriptionsX509DescriptionCaOptionsFields(r *CertificateAuthority, o *CertificateAuthorityCaCertificateDescriptionsX509DescriptionCaOptions) error {
	return nil
}
func extractCertificateAuthorityCaCertificateDescriptionsX509DescriptionPolicyIdsFields(r *CertificateAuthority, o *CertificateAuthorityCaCertificateDescriptionsX509DescriptionPolicyIds) error {
	return nil
}
func extractCertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensionsFields(r *CertificateAuthority, o *CertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensions) error {
	vObjectId := o.ObjectId
	if vObjectId == nil {
		// note: explicitly not the empty object.
		vObjectId = &CertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensionsObjectId{}
	}
	if err := extractCertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensionsObjectIdFields(r, vObjectId); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vObjectId) {
		o.ObjectId = vObjectId
	}
	return nil
}
func extractCertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensionsObjectIdFields(r *CertificateAuthority, o *CertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensionsObjectId) error {
	return nil
}
func extractCertificateAuthorityCaCertificateDescriptionsPublicKeyFields(r *CertificateAuthority, o *CertificateAuthorityCaCertificateDescriptionsPublicKey) error {
	return nil
}
func extractCertificateAuthorityCaCertificateDescriptionsSubjectKeyIdFields(r *CertificateAuthority, o *CertificateAuthorityCaCertificateDescriptionsSubjectKeyId) error {
	return nil
}
func extractCertificateAuthorityCaCertificateDescriptionsAuthorityKeyIdFields(r *CertificateAuthority, o *CertificateAuthorityCaCertificateDescriptionsAuthorityKeyId) error {
	return nil
}
func extractCertificateAuthorityCaCertificateDescriptionsCertFingerprintFields(r *CertificateAuthority, o *CertificateAuthorityCaCertificateDescriptionsCertFingerprint) error {
	return nil
}
func extractCertificateAuthorityAccessUrlsFields(r *CertificateAuthority, o *CertificateAuthorityAccessUrls) error {
	return nil
}

func postReadExtractCertificateAuthorityFields(r *CertificateAuthority) error {
	vConfig := r.Config
	if vConfig == nil {
		// note: explicitly not the empty object.
		vConfig = &CertificateAuthorityConfig{}
	}
	if err := postReadExtractCertificateAuthorityConfigFields(r, vConfig); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vConfig) {
		r.Config = vConfig
	}
	vKeySpec := r.KeySpec
	if vKeySpec == nil {
		// note: explicitly not the empty object.
		vKeySpec = &CertificateAuthorityKeySpec{}
	}
	if err := postReadExtractCertificateAuthorityKeySpecFields(r, vKeySpec); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vKeySpec) {
		r.KeySpec = vKeySpec
	}
	vSubordinateConfig := r.SubordinateConfig
	if vSubordinateConfig == nil {
		// note: explicitly not the empty object.
		vSubordinateConfig = &CertificateAuthoritySubordinateConfig{}
	}
	if err := postReadExtractCertificateAuthoritySubordinateConfigFields(r, vSubordinateConfig); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vSubordinateConfig) {
		r.SubordinateConfig = vSubordinateConfig
	}
	vAccessUrls := r.AccessUrls
	if vAccessUrls == nil {
		// note: explicitly not the empty object.
		vAccessUrls = &CertificateAuthorityAccessUrls{}
	}
	if err := postReadExtractCertificateAuthorityAccessUrlsFields(r, vAccessUrls); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vAccessUrls) {
		r.AccessUrls = vAccessUrls
	}
	return nil
}
func postReadExtractCertificateAuthorityConfigFields(r *CertificateAuthority, o *CertificateAuthorityConfig) error {
	vSubjectConfig := o.SubjectConfig
	if vSubjectConfig == nil {
		// note: explicitly not the empty object.
		vSubjectConfig = &CertificateAuthorityConfigSubjectConfig{}
	}
	if err := extractCertificateAuthorityConfigSubjectConfigFields(r, vSubjectConfig); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vSubjectConfig) {
		o.SubjectConfig = vSubjectConfig
	}
	vX509Config := o.X509Config
	if vX509Config == nil {
		// note: explicitly not the empty object.
		vX509Config = &CertificateAuthorityConfigX509Config{}
	}
	if err := extractCertificateAuthorityConfigX509ConfigFields(r, vX509Config); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vX509Config) {
		o.X509Config = vX509Config
	}
	vPublicKey := o.PublicKey
	if vPublicKey == nil {
		// note: explicitly not the empty object.
		vPublicKey = &CertificateAuthorityConfigPublicKey{}
	}
	if err := extractCertificateAuthorityConfigPublicKeyFields(r, vPublicKey); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vPublicKey) {
		o.PublicKey = vPublicKey
	}
	return nil
}
func postReadExtractCertificateAuthorityConfigSubjectConfigFields(r *CertificateAuthority, o *CertificateAuthorityConfigSubjectConfig) error {
	vSubject := o.Subject
	if vSubject == nil {
		// note: explicitly not the empty object.
		vSubject = &CertificateAuthorityConfigSubjectConfigSubject{}
	}
	if err := extractCertificateAuthorityConfigSubjectConfigSubjectFields(r, vSubject); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vSubject) {
		o.Subject = vSubject
	}
	vSubjectAltName := o.SubjectAltName
	if vSubjectAltName == nil {
		// note: explicitly not the empty object.
		vSubjectAltName = &CertificateAuthorityConfigSubjectConfigSubjectAltName{}
	}
	if err := extractCertificateAuthorityConfigSubjectConfigSubjectAltNameFields(r, vSubjectAltName); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vSubjectAltName) {
		o.SubjectAltName = vSubjectAltName
	}
	return nil
}
func postReadExtractCertificateAuthorityConfigSubjectConfigSubjectFields(r *CertificateAuthority, o *CertificateAuthorityConfigSubjectConfigSubject) error {
	return nil
}
func postReadExtractCertificateAuthorityConfigSubjectConfigSubjectAltNameFields(r *CertificateAuthority, o *CertificateAuthorityConfigSubjectConfigSubjectAltName) error {
	return nil
}
func postReadExtractCertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSansFields(r *CertificateAuthority, o *CertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSans) error {
	vObjectId := o.ObjectId
	if vObjectId == nil {
		// note: explicitly not the empty object.
		vObjectId = &CertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSansObjectId{}
	}
	if err := extractCertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSansObjectIdFields(r, vObjectId); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vObjectId) {
		o.ObjectId = vObjectId
	}
	return nil
}
func postReadExtractCertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSansObjectIdFields(r *CertificateAuthority, o *CertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSansObjectId) error {
	return nil
}
func postReadExtractCertificateAuthorityConfigX509ConfigFields(r *CertificateAuthority, o *CertificateAuthorityConfigX509Config) error {
	vKeyUsage := o.KeyUsage
	if vKeyUsage == nil {
		// note: explicitly not the empty object.
		vKeyUsage = &CertificateAuthorityConfigX509ConfigKeyUsage{}
	}
	if err := extractCertificateAuthorityConfigX509ConfigKeyUsageFields(r, vKeyUsage); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vKeyUsage) {
		o.KeyUsage = vKeyUsage
	}
	vCaOptions := o.CaOptions
	if vCaOptions == nil {
		// note: explicitly not the empty object.
		vCaOptions = &CertificateAuthorityConfigX509ConfigCaOptions{}
	}
	if err := extractCertificateAuthorityConfigX509ConfigCaOptionsFields(r, vCaOptions); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vCaOptions) {
		o.CaOptions = vCaOptions
	}
	return nil
}
func postReadExtractCertificateAuthorityConfigX509ConfigKeyUsageFields(r *CertificateAuthority, o *CertificateAuthorityConfigX509ConfigKeyUsage) error {
	vBaseKeyUsage := o.BaseKeyUsage
	if vBaseKeyUsage == nil {
		// note: explicitly not the empty object.
		vBaseKeyUsage = &CertificateAuthorityConfigX509ConfigKeyUsageBaseKeyUsage{}
	}
	if err := extractCertificateAuthorityConfigX509ConfigKeyUsageBaseKeyUsageFields(r, vBaseKeyUsage); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vBaseKeyUsage) {
		o.BaseKeyUsage = vBaseKeyUsage
	}
	vExtendedKeyUsage := o.ExtendedKeyUsage
	if vExtendedKeyUsage == nil {
		// note: explicitly not the empty object.
		vExtendedKeyUsage = &CertificateAuthorityConfigX509ConfigKeyUsageExtendedKeyUsage{}
	}
	if err := extractCertificateAuthorityConfigX509ConfigKeyUsageExtendedKeyUsageFields(r, vExtendedKeyUsage); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vExtendedKeyUsage) {
		o.ExtendedKeyUsage = vExtendedKeyUsage
	}
	return nil
}
func postReadExtractCertificateAuthorityConfigX509ConfigKeyUsageBaseKeyUsageFields(r *CertificateAuthority, o *CertificateAuthorityConfigX509ConfigKeyUsageBaseKeyUsage) error {
	return nil
}
func postReadExtractCertificateAuthorityConfigX509ConfigKeyUsageExtendedKeyUsageFields(r *CertificateAuthority, o *CertificateAuthorityConfigX509ConfigKeyUsageExtendedKeyUsage) error {
	return nil
}
func postReadExtractCertificateAuthorityConfigX509ConfigKeyUsageUnknownExtendedKeyUsagesFields(r *CertificateAuthority, o *CertificateAuthorityConfigX509ConfigKeyUsageUnknownExtendedKeyUsages) error {
	return nil
}
func postReadExtractCertificateAuthorityConfigX509ConfigCaOptionsFields(r *CertificateAuthority, o *CertificateAuthorityConfigX509ConfigCaOptions) error {
	return nil
}
func postReadExtractCertificateAuthorityConfigX509ConfigPolicyIdsFields(r *CertificateAuthority, o *CertificateAuthorityConfigX509ConfigPolicyIds) error {
	return nil
}
func postReadExtractCertificateAuthorityConfigX509ConfigAdditionalExtensionsFields(r *CertificateAuthority, o *CertificateAuthorityConfigX509ConfigAdditionalExtensions) error {
	vObjectId := o.ObjectId
	if vObjectId == nil {
		// note: explicitly not the empty object.
		vObjectId = &CertificateAuthorityConfigX509ConfigAdditionalExtensionsObjectId{}
	}
	if err := extractCertificateAuthorityConfigX509ConfigAdditionalExtensionsObjectIdFields(r, vObjectId); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vObjectId) {
		o.ObjectId = vObjectId
	}
	return nil
}
func postReadExtractCertificateAuthorityConfigX509ConfigAdditionalExtensionsObjectIdFields(r *CertificateAuthority, o *CertificateAuthorityConfigX509ConfigAdditionalExtensionsObjectId) error {
	return nil
}
func postReadExtractCertificateAuthorityConfigPublicKeyFields(r *CertificateAuthority, o *CertificateAuthorityConfigPublicKey) error {
	return nil
}
func postReadExtractCertificateAuthorityKeySpecFields(r *CertificateAuthority, o *CertificateAuthorityKeySpec) error {
	return nil
}
func postReadExtractCertificateAuthoritySubordinateConfigFields(r *CertificateAuthority, o *CertificateAuthoritySubordinateConfig) error {
	vPemIssuerChain := o.PemIssuerChain
	if vPemIssuerChain == nil {
		// note: explicitly not the empty object.
		vPemIssuerChain = &CertificateAuthoritySubordinateConfigPemIssuerChain{}
	}
	if err := extractCertificateAuthoritySubordinateConfigPemIssuerChainFields(r, vPemIssuerChain); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vPemIssuerChain) {
		o.PemIssuerChain = vPemIssuerChain
	}
	return nil
}
func postReadExtractCertificateAuthoritySubordinateConfigPemIssuerChainFields(r *CertificateAuthority, o *CertificateAuthoritySubordinateConfigPemIssuerChain) error {
	return nil
}
func postReadExtractCertificateAuthorityCaCertificateDescriptionsFields(r *CertificateAuthority, o *CertificateAuthorityCaCertificateDescriptions) error {
	vSubjectDescription := o.SubjectDescription
	if vSubjectDescription == nil {
		// note: explicitly not the empty object.
		vSubjectDescription = &CertificateAuthorityCaCertificateDescriptionsSubjectDescription{}
	}
	if err := extractCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionFields(r, vSubjectDescription); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vSubjectDescription) {
		o.SubjectDescription = vSubjectDescription
	}
	vX509Description := o.X509Description
	if vX509Description == nil {
		// note: explicitly not the empty object.
		vX509Description = &CertificateAuthorityCaCertificateDescriptionsX509Description{}
	}
	if err := extractCertificateAuthorityCaCertificateDescriptionsX509DescriptionFields(r, vX509Description); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vX509Description) {
		o.X509Description = vX509Description
	}
	vPublicKey := o.PublicKey
	if vPublicKey == nil {
		// note: explicitly not the empty object.
		vPublicKey = &CertificateAuthorityCaCertificateDescriptionsPublicKey{}
	}
	if err := extractCertificateAuthorityCaCertificateDescriptionsPublicKeyFields(r, vPublicKey); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vPublicKey) {
		o.PublicKey = vPublicKey
	}
	vSubjectKeyId := o.SubjectKeyId
	if vSubjectKeyId == nil {
		// note: explicitly not the empty object.
		vSubjectKeyId = &CertificateAuthorityCaCertificateDescriptionsSubjectKeyId{}
	}
	if err := extractCertificateAuthorityCaCertificateDescriptionsSubjectKeyIdFields(r, vSubjectKeyId); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vSubjectKeyId) {
		o.SubjectKeyId = vSubjectKeyId
	}
	vAuthorityKeyId := o.AuthorityKeyId
	if vAuthorityKeyId == nil {
		// note: explicitly not the empty object.
		vAuthorityKeyId = &CertificateAuthorityCaCertificateDescriptionsAuthorityKeyId{}
	}
	if err := extractCertificateAuthorityCaCertificateDescriptionsAuthorityKeyIdFields(r, vAuthorityKeyId); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vAuthorityKeyId) {
		o.AuthorityKeyId = vAuthorityKeyId
	}
	vCertFingerprint := o.CertFingerprint
	if vCertFingerprint == nil {
		// note: explicitly not the empty object.
		vCertFingerprint = &CertificateAuthorityCaCertificateDescriptionsCertFingerprint{}
	}
	if err := extractCertificateAuthorityCaCertificateDescriptionsCertFingerprintFields(r, vCertFingerprint); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vCertFingerprint) {
		o.CertFingerprint = vCertFingerprint
	}
	return nil
}
func postReadExtractCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionFields(r *CertificateAuthority, o *CertificateAuthorityCaCertificateDescriptionsSubjectDescription) error {
	vSubject := o.Subject
	if vSubject == nil {
		// note: explicitly not the empty object.
		vSubject = &CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubject{}
	}
	if err := extractCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectFields(r, vSubject); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vSubject) {
		o.Subject = vSubject
	}
	vSubjectAltName := o.SubjectAltName
	if vSubjectAltName == nil {
		// note: explicitly not the empty object.
		vSubjectAltName = &CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltName{}
	}
	if err := extractCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameFields(r, vSubjectAltName); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vSubjectAltName) {
		o.SubjectAltName = vSubjectAltName
	}
	return nil
}
func postReadExtractCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectFields(r *CertificateAuthority, o *CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubject) error {
	return nil
}
func postReadExtractCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameFields(r *CertificateAuthority, o *CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltName) error {
	return nil
}
func postReadExtractCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSansFields(r *CertificateAuthority, o *CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSans) error {
	vObjectId := o.ObjectId
	if vObjectId == nil {
		// note: explicitly not the empty object.
		vObjectId = &CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSansObjectId{}
	}
	if err := extractCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSansObjectIdFields(r, vObjectId); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vObjectId) {
		o.ObjectId = vObjectId
	}
	return nil
}
func postReadExtractCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSansObjectIdFields(r *CertificateAuthority, o *CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSansObjectId) error {
	return nil
}
func postReadExtractCertificateAuthorityCaCertificateDescriptionsX509DescriptionFields(r *CertificateAuthority, o *CertificateAuthorityCaCertificateDescriptionsX509Description) error {
	vKeyUsage := o.KeyUsage
	if vKeyUsage == nil {
		// note: explicitly not the empty object.
		vKeyUsage = &CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsage{}
	}
	if err := extractCertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageFields(r, vKeyUsage); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vKeyUsage) {
		o.KeyUsage = vKeyUsage
	}
	vCaOptions := o.CaOptions
	if vCaOptions == nil {
		// note: explicitly not the empty object.
		vCaOptions = &CertificateAuthorityCaCertificateDescriptionsX509DescriptionCaOptions{}
	}
	if err := extractCertificateAuthorityCaCertificateDescriptionsX509DescriptionCaOptionsFields(r, vCaOptions); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vCaOptions) {
		o.CaOptions = vCaOptions
	}
	return nil
}
func postReadExtractCertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageFields(r *CertificateAuthority, o *CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsage) error {
	vBaseKeyUsage := o.BaseKeyUsage
	if vBaseKeyUsage == nil {
		// note: explicitly not the empty object.
		vBaseKeyUsage = &CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageBaseKeyUsage{}
	}
	if err := extractCertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageBaseKeyUsageFields(r, vBaseKeyUsage); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vBaseKeyUsage) {
		o.BaseKeyUsage = vBaseKeyUsage
	}
	vExtendedKeyUsage := o.ExtendedKeyUsage
	if vExtendedKeyUsage == nil {
		// note: explicitly not the empty object.
		vExtendedKeyUsage = &CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageExtendedKeyUsage{}
	}
	if err := extractCertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageExtendedKeyUsageFields(r, vExtendedKeyUsage); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vExtendedKeyUsage) {
		o.ExtendedKeyUsage = vExtendedKeyUsage
	}
	return nil
}
func postReadExtractCertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageBaseKeyUsageFields(r *CertificateAuthority, o *CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageBaseKeyUsage) error {
	return nil
}
func postReadExtractCertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageExtendedKeyUsageFields(r *CertificateAuthority, o *CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageExtendedKeyUsage) error {
	return nil
}
func postReadExtractCertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageUnknownExtendedKeyUsagesFields(r *CertificateAuthority, o *CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageUnknownExtendedKeyUsages) error {
	return nil
}
func postReadExtractCertificateAuthorityCaCertificateDescriptionsX509DescriptionCaOptionsFields(r *CertificateAuthority, o *CertificateAuthorityCaCertificateDescriptionsX509DescriptionCaOptions) error {
	return nil
}
func postReadExtractCertificateAuthorityCaCertificateDescriptionsX509DescriptionPolicyIdsFields(r *CertificateAuthority, o *CertificateAuthorityCaCertificateDescriptionsX509DescriptionPolicyIds) error {
	return nil
}
func postReadExtractCertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensionsFields(r *CertificateAuthority, o *CertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensions) error {
	vObjectId := o.ObjectId
	if vObjectId == nil {
		// note: explicitly not the empty object.
		vObjectId = &CertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensionsObjectId{}
	}
	if err := extractCertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensionsObjectIdFields(r, vObjectId); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vObjectId) {
		o.ObjectId = vObjectId
	}
	return nil
}
func postReadExtractCertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensionsObjectIdFields(r *CertificateAuthority, o *CertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensionsObjectId) error {
	return nil
}
func postReadExtractCertificateAuthorityCaCertificateDescriptionsPublicKeyFields(r *CertificateAuthority, o *CertificateAuthorityCaCertificateDescriptionsPublicKey) error {
	return nil
}
func postReadExtractCertificateAuthorityCaCertificateDescriptionsSubjectKeyIdFields(r *CertificateAuthority, o *CertificateAuthorityCaCertificateDescriptionsSubjectKeyId) error {
	return nil
}
func postReadExtractCertificateAuthorityCaCertificateDescriptionsAuthorityKeyIdFields(r *CertificateAuthority, o *CertificateAuthorityCaCertificateDescriptionsAuthorityKeyId) error {
	return nil
}
func postReadExtractCertificateAuthorityCaCertificateDescriptionsCertFingerprintFields(r *CertificateAuthority, o *CertificateAuthorityCaCertificateDescriptionsCertFingerprint) error {
	return nil
}
func postReadExtractCertificateAuthorityAccessUrlsFields(r *CertificateAuthority, o *CertificateAuthorityAccessUrls) error {
	return nil
}
