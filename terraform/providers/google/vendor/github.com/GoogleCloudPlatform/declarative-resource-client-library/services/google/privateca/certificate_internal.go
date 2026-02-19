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
)

func (r *Certificate) validate() error {

	if err := dcl.ValidateAtMostOneOfFieldsSet([]string{"PemCsr", "Config"}, r.PemCsr, r.Config); err != nil {
		return err
	}
	if err := dcl.Required(r, "name"); err != nil {
		return err
	}
	if err := dcl.Required(r, "lifetime"); err != nil {
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
	if !dcl.IsEmptyValueIndirect(r.RevocationDetails) {
		if err := r.RevocationDetails.validate(); err != nil {
			return err
		}
	}
	if !dcl.IsEmptyValueIndirect(r.CertificateDescription) {
		if err := r.CertificateDescription.validate(); err != nil {
			return err
		}
	}
	return nil
}
func (r *CertificateConfig) validate() error {
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
func (r *CertificateConfigSubjectConfig) validate() error {
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
func (r *CertificateConfigSubjectConfigSubject) validate() error {
	return nil
}
func (r *CertificateConfigSubjectConfigSubjectAltName) validate() error {
	return nil
}
func (r *CertificateConfigX509Config) validate() error {
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
func (r *CertificateConfigX509ConfigKeyUsage) validate() error {
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
func (r *CertificateConfigX509ConfigKeyUsageBaseKeyUsage) validate() error {
	return nil
}
func (r *CertificateConfigX509ConfigKeyUsageExtendedKeyUsage) validate() error {
	return nil
}
func (r *CertificateConfigX509ConfigKeyUsageUnknownExtendedKeyUsages) validate() error {
	if err := dcl.Required(r, "objectIdPath"); err != nil {
		return err
	}
	return nil
}
func (r *CertificateConfigX509ConfigCaOptions) validate() error {
	return nil
}
func (r *CertificateConfigX509ConfigPolicyIds) validate() error {
	if err := dcl.Required(r, "objectIdPath"); err != nil {
		return err
	}
	return nil
}
func (r *CertificateConfigX509ConfigAdditionalExtensions) validate() error {
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
func (r *CertificateConfigX509ConfigAdditionalExtensionsObjectId) validate() error {
	if err := dcl.Required(r, "objectIdPath"); err != nil {
		return err
	}
	return nil
}
func (r *CertificateConfigPublicKey) validate() error {
	if err := dcl.Required(r, "key"); err != nil {
		return err
	}
	if err := dcl.Required(r, "format"); err != nil {
		return err
	}
	return nil
}
func (r *CertificateRevocationDetails) validate() error {
	return nil
}
func (r *CertificateCertificateDescription) validate() error {
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
func (r *CertificateCertificateDescriptionSubjectDescription) validate() error {
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
func (r *CertificateCertificateDescriptionSubjectDescriptionSubject) validate() error {
	return nil
}
func (r *CertificateCertificateDescriptionSubjectDescriptionSubjectAltName) validate() error {
	return nil
}
func (r *CertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSans) validate() error {
	if !dcl.IsEmptyValueIndirect(r.ObjectId) {
		if err := r.ObjectId.validate(); err != nil {
			return err
		}
	}
	return nil
}
func (r *CertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSansObjectId) validate() error {
	return nil
}
func (r *CertificateCertificateDescriptionX509Description) validate() error {
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
func (r *CertificateCertificateDescriptionX509DescriptionKeyUsage) validate() error {
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
func (r *CertificateCertificateDescriptionX509DescriptionKeyUsageBaseKeyUsage) validate() error {
	return nil
}
func (r *CertificateCertificateDescriptionX509DescriptionKeyUsageExtendedKeyUsage) validate() error {
	return nil
}
func (r *CertificateCertificateDescriptionX509DescriptionKeyUsageUnknownExtendedKeyUsages) validate() error {
	return nil
}
func (r *CertificateCertificateDescriptionX509DescriptionCaOptions) validate() error {
	if err := dcl.ValidateAtMostOneOfFieldsSet([]string{"IsCa"}, r.IsCa); err != nil {
		return err
	}
	if err := dcl.ValidateAtMostOneOfFieldsSet([]string{"MaxIssuerPathLength"}, r.MaxIssuerPathLength); err != nil {
		return err
	}
	return nil
}
func (r *CertificateCertificateDescriptionX509DescriptionPolicyIds) validate() error {
	return nil
}
func (r *CertificateCertificateDescriptionX509DescriptionAdditionalExtensions) validate() error {
	if !dcl.IsEmptyValueIndirect(r.ObjectId) {
		if err := r.ObjectId.validate(); err != nil {
			return err
		}
	}
	return nil
}
func (r *CertificateCertificateDescriptionX509DescriptionAdditionalExtensionsObjectId) validate() error {
	return nil
}
func (r *CertificateCertificateDescriptionPublicKey) validate() error {
	return nil
}
func (r *CertificateCertificateDescriptionSubjectKeyId) validate() error {
	return nil
}
func (r *CertificateCertificateDescriptionAuthorityKeyId) validate() error {
	return nil
}
func (r *CertificateCertificateDescriptionCertFingerprint) validate() error {
	return nil
}
func (r *Certificate) basePath() string {
	params := map[string]interface{}{}
	return dcl.Nprintf("https://privateca.googleapis.com/v1/", params)
}

func (r *Certificate) getURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"project":  dcl.ValueOrEmptyString(nr.Project),
		"location": dcl.ValueOrEmptyString(nr.Location),
		"caPool":   dcl.ValueOrEmptyString(nr.CaPool),
		"name":     dcl.ValueOrEmptyString(nr.Name),
	}
	return dcl.URL("projects/{{project}}/locations/{{location}}/caPools/{{caPool}}/certificates/{{name}}", nr.basePath(), userBasePath, params), nil
}

func (r *Certificate) listURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"project":  dcl.ValueOrEmptyString(nr.Project),
		"location": dcl.ValueOrEmptyString(nr.Location),
		"caPool":   dcl.ValueOrEmptyString(nr.CaPool),
	}
	return dcl.URL("projects/{{project}}/locations/{{location}}/caPools/{{caPool}}/certificates", nr.basePath(), userBasePath, params), nil

}

func (r *Certificate) deleteURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"project":  dcl.ValueOrEmptyString(nr.Project),
		"location": dcl.ValueOrEmptyString(nr.Location),
		"caPool":   dcl.ValueOrEmptyString(nr.CaPool),
		"name":     dcl.ValueOrEmptyString(nr.Name),
	}
	return dcl.URL("projects/{{project}}/locations/{{location}}/caPools/{{caPool}}/certificates/{{name}}:revoke", nr.basePath(), userBasePath, params), nil
}

// certificateApiOperation represents a mutable operation in the underlying REST
// API such as Create, Update, or Delete.
type certificateApiOperation interface {
	do(context.Context, *Certificate, *Client) error
}

// newUpdateCertificateUpdateCertificateRequest creates a request for an
// Certificate resource's UpdateCertificate update type by filling in the update
// fields based on the intended state of the resource.
func newUpdateCertificateUpdateCertificateRequest(ctx context.Context, f *Certificate, c *Client) (map[string]interface{}, error) {
	req := map[string]interface{}{}
	res := f
	_ = res

	if v := f.Labels; !dcl.IsEmptyValueIndirect(v) {
		req["labels"] = v
	}
	return req, nil
}

// marshalUpdateCertificateUpdateCertificateRequest converts the update into
// the final JSON request body.
func marshalUpdateCertificateUpdateCertificateRequest(c *Client, m map[string]interface{}) ([]byte, error) {

	return json.Marshal(m)
}

type updateCertificateUpdateCertificateOperation struct {
	// If the update operation has the REQUIRES_APPLY_OPTIONS trait, this will be populated.
	// Usually it will be nil - this is to prevent us from accidentally depending on apply
	// options, which should usually be unnecessary.
	ApplyOptions []dcl.ApplyOption
	FieldDiffs   []*dcl.FieldDiff
}

// do creates a request and sends it to the appropriate URL. In most operations,
// do will transcribe a subset of the resource into a request object and send a
// PUT request to a single URL.

func (op *updateCertificateUpdateCertificateOperation) do(ctx context.Context, r *Certificate, c *Client) error {
	_, err := c.GetCertificate(ctx, r)
	if err != nil {
		return err
	}

	u, err := r.updateURL(c.Config.BasePath, "UpdateCertificate")
	if err != nil {
		return err
	}
	mask := dcl.UpdateMask(op.FieldDiffs)
	u, err = dcl.AddQueryParams(u, map[string]string{"updateMask": mask})
	if err != nil {
		return err
	}

	req, err := newUpdateCertificateUpdateCertificateRequest(ctx, r, c)
	if err != nil {
		return err
	}

	c.Config.Logger.InfoWithContextf(ctx, "Created update: %#v", req)
	body, err := marshalUpdateCertificateUpdateCertificateRequest(c, req)
	if err != nil {
		return err
	}
	_, err = dcl.SendRequest(ctx, c.Config, "PATCH", u, bytes.NewBuffer(body), c.Config.RetryProvider)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) listCertificateRaw(ctx context.Context, r *Certificate, pageToken string, pageSize int32) ([]byte, error) {
	u, err := r.urlNormalized().listURL(c.Config.BasePath)
	if err != nil {
		return nil, err
	}

	m := make(map[string]string)
	if pageToken != "" {
		m["pageToken"] = pageToken
	}

	if pageSize != CertificateMaxPage {
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

type listCertificateOperation struct {
	Certificates []map[string]interface{} `json:"certificates"`
	Token        string                   `json:"nextPageToken"`
}

func (c *Client) listCertificate(ctx context.Context, r *Certificate, pageToken string, pageSize int32) ([]*Certificate, string, error) {
	b, err := c.listCertificateRaw(ctx, r, pageToken, pageSize)
	if err != nil {
		return nil, "", err
	}

	var m listCertificateOperation
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, "", err
	}

	var l []*Certificate
	for _, v := range m.Certificates {
		res, err := unmarshalMapCertificate(v, c, r)
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

func (c *Client) deleteAllCertificate(ctx context.Context, f func(*Certificate) bool, resources []*Certificate) error {
	var errors []string
	for _, res := range resources {
		if f(res) {
			// We do not want deleteAll to fail on a deletion or else it will stop deleting other resources.
			err := c.DeleteCertificate(ctx, res)
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

type deleteCertificateOperation struct{}

func (op *deleteCertificateOperation) do(ctx context.Context, r *Certificate, c *Client) error {
	r, err := c.GetCertificate(ctx, r)
	if err != nil {
		if dcl.IsNotFound(err) {
			c.Config.Logger.InfoWithContextf(ctx, "Certificate not found, returning. Original error: %v", err)
			return nil
		}
		c.Config.Logger.WarningWithContextf(ctx, "GetCertificate checking for existence. error: %v", err)
		return err
	}

	u, err := r.deleteURL(c.Config.BasePath)
	if err != nil {
		return err
	}

	// Delete should never have a body
	body := &bytes.Buffer{}
	_, err = dcl.SendRequest(ctx, c.Config, "POST", u, body, c.Config.RetryProvider)
	if err != nil {
		return fmt.Errorf("failed to delete Certificate: %w", err)
	}
	return nil
}

// Create operations are similar to Update operations, although they do not have
// specific request objects. The Create request object is the json encoding of
// the resource, which is modified by res.marshal to form the base request body.
type createCertificateOperation struct {
	response map[string]interface{}
}

func (op *createCertificateOperation) FirstResponse() (map[string]interface{}, bool) {
	return op.response, len(op.response) > 0
}

func (op *createCertificateOperation) do(ctx context.Context, r *Certificate, c *Client) error {
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

	if _, err := c.GetCertificate(ctx, r); err != nil {
		c.Config.Logger.WarningWithContextf(ctx, "get returned error: %v", err)
		return err
	}

	return nil
}

func (c *Client) getCertificateRaw(ctx context.Context, r *Certificate) ([]byte, error) {

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

func (c *Client) certificateDiffsForRawDesired(ctx context.Context, rawDesired *Certificate, opts ...dcl.ApplyOption) (initial, desired *Certificate, diffs []*dcl.FieldDiff, err error) {
	c.Config.Logger.InfoWithContext(ctx, "Fetching initial state...")
	// First, let us see if the user provided a state hint.  If they did, we will start fetching based on that.
	var fetchState *Certificate
	if sh := dcl.FetchStateHint(opts); sh != nil {
		if r, ok := sh.(*Certificate); !ok {
			c.Config.Logger.WarningWithContextf(ctx, "Initial state hint was of the wrong type; expected Certificate, got %T", sh)
		} else {
			fetchState = r
		}
	}
	if fetchState == nil {
		fetchState = rawDesired
	}

	// 1.2: Retrieval of raw initial state from API
	rawInitial, err := c.GetCertificate(ctx, fetchState)
	if rawInitial == nil {
		if !dcl.IsNotFound(err) {
			c.Config.Logger.WarningWithContextf(ctx, "Failed to retrieve whether a Certificate resource already exists: %s", err)
			return nil, nil, nil, fmt.Errorf("failed to retrieve Certificate resource: %v", err)
		}
		c.Config.Logger.InfoWithContext(ctx, "Found that Certificate resource did not exist.")
		// Perform canonicalization to pick up defaults.
		desired, err = canonicalizeCertificateDesiredState(rawDesired, rawInitial)
		return nil, desired, nil, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Found initial state for Certificate: %v", rawInitial)
	c.Config.Logger.InfoWithContextf(ctx, "Initial desired state for Certificate: %v", rawDesired)

	// The Get call applies postReadExtract and so the result may contain fields that are not part of API version.
	if err := extractCertificateFields(rawInitial); err != nil {
		return nil, nil, nil, err
	}

	// 1.3: Canonicalize raw initial state into initial state.
	initial, err = canonicalizeCertificateInitialState(rawInitial, rawDesired)
	if err != nil {
		return nil, nil, nil, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Canonicalized initial state for Certificate: %v", initial)

	// 1.4: Canonicalize raw desired state into desired state.
	desired, err = canonicalizeCertificateDesiredState(rawDesired, rawInitial, opts...)
	if err != nil {
		return nil, nil, nil, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Canonicalized desired state for Certificate: %v", desired)

	// 2.1: Comparison of initial and desired state.
	diffs, err = diffCertificate(c, desired, initial, opts...)
	return initial, desired, diffs, err
}

func canonicalizeCertificateInitialState(rawInitial, rawDesired *Certificate) (*Certificate, error) {
	// TODO(magic-modules-eng): write canonicalizer once relevant traits are added.

	if !dcl.IsZeroValue(rawInitial.PemCsr) {
		// Check if anything else is set.
		if dcl.AnySet(rawInitial.Config) {
			rawInitial.PemCsr = dcl.String("")
		}
	}

	if !dcl.IsZeroValue(rawInitial.Config) {
		// Check if anything else is set.
		if dcl.AnySet(rawInitial.PemCsr) {
			rawInitial.Config = EmptyCertificateConfig
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

func canonicalizeCertificateDesiredState(rawDesired, rawInitial *Certificate, opts ...dcl.ApplyOption) (*Certificate, error) {

	if rawInitial == nil {
		// Since the initial state is empty, the desired state is all we have.
		// We canonicalize the remaining nested objects with nil to pick up defaults.
		rawDesired.Config = canonicalizeCertificateConfig(rawDesired.Config, nil, opts...)
		rawDesired.RevocationDetails = canonicalizeCertificateRevocationDetails(rawDesired.RevocationDetails, nil, opts...)
		rawDesired.CertificateDescription = canonicalizeCertificateCertificateDescription(rawDesired.CertificateDescription, nil, opts...)

		return rawDesired, nil
	}
	canonicalDesired := &Certificate{}
	if dcl.PartialSelfLinkToSelfLink(rawDesired.Name, rawInitial.Name) {
		canonicalDesired.Name = rawInitial.Name
	} else {
		canonicalDesired.Name = rawDesired.Name
	}
	if dcl.StringCanonicalize(rawDesired.PemCsr, rawInitial.PemCsr) {
		canonicalDesired.PemCsr = rawInitial.PemCsr
	} else {
		canonicalDesired.PemCsr = rawDesired.PemCsr
	}
	canonicalDesired.Config = canonicalizeCertificateConfig(rawDesired.Config, rawInitial.Config, opts...)
	if dcl.StringCanonicalize(rawDesired.Lifetime, rawInitial.Lifetime) {
		canonicalDesired.Lifetime = rawInitial.Lifetime
	} else {
		canonicalDesired.Lifetime = rawDesired.Lifetime
	}
	if dcl.PartialSelfLinkToSelfLink(rawDesired.CertificateTemplate, rawInitial.CertificateTemplate) {
		canonicalDesired.CertificateTemplate = rawInitial.CertificateTemplate
	} else {
		canonicalDesired.CertificateTemplate = rawDesired.CertificateTemplate
	}
	if dcl.IsZeroValue(rawDesired.SubjectMode) || (dcl.IsEmptyValueIndirect(rawDesired.SubjectMode) && dcl.IsEmptyValueIndirect(rawInitial.SubjectMode)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		canonicalDesired.SubjectMode = rawInitial.SubjectMode
	} else {
		canonicalDesired.SubjectMode = rawDesired.SubjectMode
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
	if dcl.NameToSelfLink(rawDesired.CertificateAuthority, rawInitial.CertificateAuthority) {
		canonicalDesired.CertificateAuthority = rawInitial.CertificateAuthority
	} else {
		canonicalDesired.CertificateAuthority = rawDesired.CertificateAuthority
	}

	if canonicalDesired.PemCsr != nil {
		// Check if anything else is set.
		if dcl.AnySet(rawDesired.Config) {
			canonicalDesired.PemCsr = dcl.String("")
		}
	}

	if canonicalDesired.Config != nil {
		// Check if anything else is set.
		if dcl.AnySet(rawDesired.PemCsr) {
			canonicalDesired.Config = EmptyCertificateConfig
		}
	}

	return canonicalDesired, nil
}

func canonicalizeCertificateNewState(c *Client, rawNew, rawDesired *Certificate) (*Certificate, error) {

	if dcl.IsEmptyValueIndirect(rawNew.Name) && dcl.IsEmptyValueIndirect(rawDesired.Name) {
		rawNew.Name = rawDesired.Name
	} else {
		if dcl.PartialSelfLinkToSelfLink(rawDesired.Name, rawNew.Name) {
			rawNew.Name = rawDesired.Name
		}
	}

	if dcl.IsEmptyValueIndirect(rawNew.PemCsr) && dcl.IsEmptyValueIndirect(rawDesired.PemCsr) {
		rawNew.PemCsr = rawDesired.PemCsr
	} else {
		if dcl.StringCanonicalize(rawDesired.PemCsr, rawNew.PemCsr) {
			rawNew.PemCsr = rawDesired.PemCsr
		}
	}

	if dcl.IsEmptyValueIndirect(rawNew.Config) && dcl.IsEmptyValueIndirect(rawDesired.Config) {
		rawNew.Config = rawDesired.Config
	} else {
		rawNew.Config = canonicalizeNewCertificateConfig(c, rawDesired.Config, rawNew.Config)
	}

	if dcl.IsEmptyValueIndirect(rawNew.IssuerCertificateAuthority) && dcl.IsEmptyValueIndirect(rawDesired.IssuerCertificateAuthority) {
		rawNew.IssuerCertificateAuthority = rawDesired.IssuerCertificateAuthority
	} else {
	}

	if dcl.IsEmptyValueIndirect(rawNew.Lifetime) && dcl.IsEmptyValueIndirect(rawDesired.Lifetime) {
		rawNew.Lifetime = rawDesired.Lifetime
	} else {
		if dcl.StringCanonicalize(rawDesired.Lifetime, rawNew.Lifetime) {
			rawNew.Lifetime = rawDesired.Lifetime
		}
	}

	if dcl.IsEmptyValueIndirect(rawNew.CertificateTemplate) && dcl.IsEmptyValueIndirect(rawDesired.CertificateTemplate) {
		rawNew.CertificateTemplate = rawDesired.CertificateTemplate
	} else {
		if dcl.PartialSelfLinkToSelfLink(rawDesired.CertificateTemplate, rawNew.CertificateTemplate) {
			rawNew.CertificateTemplate = rawDesired.CertificateTemplate
		}
	}

	if dcl.IsEmptyValueIndirect(rawNew.SubjectMode) && dcl.IsEmptyValueIndirect(rawDesired.SubjectMode) {
		rawNew.SubjectMode = rawDesired.SubjectMode
	} else {
	}

	if dcl.IsEmptyValueIndirect(rawNew.RevocationDetails) && dcl.IsEmptyValueIndirect(rawDesired.RevocationDetails) {
		rawNew.RevocationDetails = rawDesired.RevocationDetails
	} else {
		rawNew.RevocationDetails = canonicalizeNewCertificateRevocationDetails(c, rawDesired.RevocationDetails, rawNew.RevocationDetails)
	}

	if dcl.IsEmptyValueIndirect(rawNew.PemCertificate) && dcl.IsEmptyValueIndirect(rawDesired.PemCertificate) {
		rawNew.PemCertificate = rawDesired.PemCertificate
	} else {
		if dcl.StringCanonicalize(rawDesired.PemCertificate, rawNew.PemCertificate) {
			rawNew.PemCertificate = rawDesired.PemCertificate
		}
	}

	if dcl.IsEmptyValueIndirect(rawNew.CertificateDescription) && dcl.IsEmptyValueIndirect(rawDesired.CertificateDescription) {
		rawNew.CertificateDescription = rawDesired.CertificateDescription
	} else {
		rawNew.CertificateDescription = canonicalizeNewCertificateCertificateDescription(c, rawDesired.CertificateDescription, rawNew.CertificateDescription)
	}

	if dcl.IsEmptyValueIndirect(rawNew.PemCertificateChain) && dcl.IsEmptyValueIndirect(rawDesired.PemCertificateChain) {
		rawNew.PemCertificateChain = rawDesired.PemCertificateChain
	} else {
		if dcl.StringArrayCanonicalize(rawDesired.PemCertificateChain, rawNew.PemCertificateChain) {
			rawNew.PemCertificateChain = rawDesired.PemCertificateChain
		}
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

	rawNew.Project = rawDesired.Project

	rawNew.Location = rawDesired.Location

	rawNew.CaPool = rawDesired.CaPool

	rawNew.CertificateAuthority = rawDesired.CertificateAuthority

	return rawNew, nil
}

func canonicalizeCertificateConfig(des, initial *CertificateConfig, opts ...dcl.ApplyOption) *CertificateConfig {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &CertificateConfig{}

	cDes.SubjectConfig = canonicalizeCertificateConfigSubjectConfig(des.SubjectConfig, initial.SubjectConfig, opts...)
	cDes.X509Config = canonicalizeCertificateConfigX509Config(des.X509Config, initial.X509Config, opts...)
	cDes.PublicKey = canonicalizeCertificateConfigPublicKey(des.PublicKey, initial.PublicKey, opts...)

	return cDes
}

func canonicalizeCertificateConfigSlice(des, initial []CertificateConfig, opts ...dcl.ApplyOption) []CertificateConfig {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]CertificateConfig, 0, len(des))
		for _, d := range des {
			cd := canonicalizeCertificateConfig(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]CertificateConfig, 0, len(des))
	for i, d := range des {
		cd := canonicalizeCertificateConfig(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewCertificateConfig(c *Client, des, nw *CertificateConfig) *CertificateConfig {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for CertificateConfig while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	nw.SubjectConfig = canonicalizeNewCertificateConfigSubjectConfig(c, des.SubjectConfig, nw.SubjectConfig)
	nw.X509Config = canonicalizeNewCertificateConfigX509Config(c, des.X509Config, nw.X509Config)
	nw.PublicKey = canonicalizeNewCertificateConfigPublicKey(c, des.PublicKey, nw.PublicKey)

	return nw
}

func canonicalizeNewCertificateConfigSet(c *Client, des, nw []CertificateConfig) []CertificateConfig {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []CertificateConfig
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareCertificateConfigNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewCertificateConfig(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewCertificateConfigSlice(c *Client, des, nw []CertificateConfig) []CertificateConfig {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []CertificateConfig
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewCertificateConfig(c, &d, &n))
	}

	return items
}

func canonicalizeCertificateConfigSubjectConfig(des, initial *CertificateConfigSubjectConfig, opts ...dcl.ApplyOption) *CertificateConfigSubjectConfig {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &CertificateConfigSubjectConfig{}

	cDes.Subject = canonicalizeCertificateConfigSubjectConfigSubject(des.Subject, initial.Subject, opts...)
	cDes.SubjectAltName = canonicalizeCertificateConfigSubjectConfigSubjectAltName(des.SubjectAltName, initial.SubjectAltName, opts...)

	return cDes
}

func canonicalizeCertificateConfigSubjectConfigSlice(des, initial []CertificateConfigSubjectConfig, opts ...dcl.ApplyOption) []CertificateConfigSubjectConfig {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]CertificateConfigSubjectConfig, 0, len(des))
		for _, d := range des {
			cd := canonicalizeCertificateConfigSubjectConfig(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]CertificateConfigSubjectConfig, 0, len(des))
	for i, d := range des {
		cd := canonicalizeCertificateConfigSubjectConfig(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewCertificateConfigSubjectConfig(c *Client, des, nw *CertificateConfigSubjectConfig) *CertificateConfigSubjectConfig {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for CertificateConfigSubjectConfig while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	nw.Subject = canonicalizeNewCertificateConfigSubjectConfigSubject(c, des.Subject, nw.Subject)
	nw.SubjectAltName = canonicalizeNewCertificateConfigSubjectConfigSubjectAltName(c, des.SubjectAltName, nw.SubjectAltName)

	return nw
}

func canonicalizeNewCertificateConfigSubjectConfigSet(c *Client, des, nw []CertificateConfigSubjectConfig) []CertificateConfigSubjectConfig {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []CertificateConfigSubjectConfig
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareCertificateConfigSubjectConfigNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewCertificateConfigSubjectConfig(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewCertificateConfigSubjectConfigSlice(c *Client, des, nw []CertificateConfigSubjectConfig) []CertificateConfigSubjectConfig {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []CertificateConfigSubjectConfig
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewCertificateConfigSubjectConfig(c, &d, &n))
	}

	return items
}

func canonicalizeCertificateConfigSubjectConfigSubject(des, initial *CertificateConfigSubjectConfigSubject, opts ...dcl.ApplyOption) *CertificateConfigSubjectConfigSubject {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &CertificateConfigSubjectConfigSubject{}

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

func canonicalizeCertificateConfigSubjectConfigSubjectSlice(des, initial []CertificateConfigSubjectConfigSubject, opts ...dcl.ApplyOption) []CertificateConfigSubjectConfigSubject {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]CertificateConfigSubjectConfigSubject, 0, len(des))
		for _, d := range des {
			cd := canonicalizeCertificateConfigSubjectConfigSubject(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]CertificateConfigSubjectConfigSubject, 0, len(des))
	for i, d := range des {
		cd := canonicalizeCertificateConfigSubjectConfigSubject(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewCertificateConfigSubjectConfigSubject(c *Client, des, nw *CertificateConfigSubjectConfigSubject) *CertificateConfigSubjectConfigSubject {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for CertificateConfigSubjectConfigSubject while comparing non-nil desired to nil actual.  Returning desired object.")
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

func canonicalizeNewCertificateConfigSubjectConfigSubjectSet(c *Client, des, nw []CertificateConfigSubjectConfigSubject) []CertificateConfigSubjectConfigSubject {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []CertificateConfigSubjectConfigSubject
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareCertificateConfigSubjectConfigSubjectNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewCertificateConfigSubjectConfigSubject(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewCertificateConfigSubjectConfigSubjectSlice(c *Client, des, nw []CertificateConfigSubjectConfigSubject) []CertificateConfigSubjectConfigSubject {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []CertificateConfigSubjectConfigSubject
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewCertificateConfigSubjectConfigSubject(c, &d, &n))
	}

	return items
}

func canonicalizeCertificateConfigSubjectConfigSubjectAltName(des, initial *CertificateConfigSubjectConfigSubjectAltName, opts ...dcl.ApplyOption) *CertificateConfigSubjectConfigSubjectAltName {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &CertificateConfigSubjectConfigSubjectAltName{}

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

	return cDes
}

func canonicalizeCertificateConfigSubjectConfigSubjectAltNameSlice(des, initial []CertificateConfigSubjectConfigSubjectAltName, opts ...dcl.ApplyOption) []CertificateConfigSubjectConfigSubjectAltName {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]CertificateConfigSubjectConfigSubjectAltName, 0, len(des))
		for _, d := range des {
			cd := canonicalizeCertificateConfigSubjectConfigSubjectAltName(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]CertificateConfigSubjectConfigSubjectAltName, 0, len(des))
	for i, d := range des {
		cd := canonicalizeCertificateConfigSubjectConfigSubjectAltName(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewCertificateConfigSubjectConfigSubjectAltName(c *Client, des, nw *CertificateConfigSubjectConfigSubjectAltName) *CertificateConfigSubjectConfigSubjectAltName {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for CertificateConfigSubjectConfigSubjectAltName while comparing non-nil desired to nil actual.  Returning desired object.")
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

	return nw
}

func canonicalizeNewCertificateConfigSubjectConfigSubjectAltNameSet(c *Client, des, nw []CertificateConfigSubjectConfigSubjectAltName) []CertificateConfigSubjectConfigSubjectAltName {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []CertificateConfigSubjectConfigSubjectAltName
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareCertificateConfigSubjectConfigSubjectAltNameNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewCertificateConfigSubjectConfigSubjectAltName(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewCertificateConfigSubjectConfigSubjectAltNameSlice(c *Client, des, nw []CertificateConfigSubjectConfigSubjectAltName) []CertificateConfigSubjectConfigSubjectAltName {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []CertificateConfigSubjectConfigSubjectAltName
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewCertificateConfigSubjectConfigSubjectAltName(c, &d, &n))
	}

	return items
}

func canonicalizeCertificateConfigX509Config(des, initial *CertificateConfigX509Config, opts ...dcl.ApplyOption) *CertificateConfigX509Config {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &CertificateConfigX509Config{}

	cDes.KeyUsage = canonicalizeCertificateConfigX509ConfigKeyUsage(des.KeyUsage, initial.KeyUsage, opts...)
	cDes.CaOptions = canonicalizeCertificateConfigX509ConfigCaOptions(des.CaOptions, initial.CaOptions, opts...)
	cDes.PolicyIds = canonicalizeCertificateConfigX509ConfigPolicyIdsSlice(des.PolicyIds, initial.PolicyIds, opts...)
	if dcl.StringArrayCanonicalize(des.AiaOcspServers, initial.AiaOcspServers) {
		cDes.AiaOcspServers = initial.AiaOcspServers
	} else {
		cDes.AiaOcspServers = des.AiaOcspServers
	}
	cDes.AdditionalExtensions = canonicalizeCertificateConfigX509ConfigAdditionalExtensionsSlice(des.AdditionalExtensions, initial.AdditionalExtensions, opts...)

	return cDes
}

func canonicalizeCertificateConfigX509ConfigSlice(des, initial []CertificateConfigX509Config, opts ...dcl.ApplyOption) []CertificateConfigX509Config {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]CertificateConfigX509Config, 0, len(des))
		for _, d := range des {
			cd := canonicalizeCertificateConfigX509Config(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]CertificateConfigX509Config, 0, len(des))
	for i, d := range des {
		cd := canonicalizeCertificateConfigX509Config(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewCertificateConfigX509Config(c *Client, des, nw *CertificateConfigX509Config) *CertificateConfigX509Config {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for CertificateConfigX509Config while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	nw.KeyUsage = canonicalizeNewCertificateConfigX509ConfigKeyUsage(c, des.KeyUsage, nw.KeyUsage)
	nw.CaOptions = canonicalizeNewCertificateConfigX509ConfigCaOptions(c, des.CaOptions, nw.CaOptions)
	nw.PolicyIds = canonicalizeNewCertificateConfigX509ConfigPolicyIdsSlice(c, des.PolicyIds, nw.PolicyIds)
	if dcl.StringArrayCanonicalize(des.AiaOcspServers, nw.AiaOcspServers) {
		nw.AiaOcspServers = des.AiaOcspServers
	}
	nw.AdditionalExtensions = canonicalizeNewCertificateConfigX509ConfigAdditionalExtensionsSlice(c, des.AdditionalExtensions, nw.AdditionalExtensions)

	return nw
}

func canonicalizeNewCertificateConfigX509ConfigSet(c *Client, des, nw []CertificateConfigX509Config) []CertificateConfigX509Config {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []CertificateConfigX509Config
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareCertificateConfigX509ConfigNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewCertificateConfigX509Config(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewCertificateConfigX509ConfigSlice(c *Client, des, nw []CertificateConfigX509Config) []CertificateConfigX509Config {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []CertificateConfigX509Config
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewCertificateConfigX509Config(c, &d, &n))
	}

	return items
}

func canonicalizeCertificateConfigX509ConfigKeyUsage(des, initial *CertificateConfigX509ConfigKeyUsage, opts ...dcl.ApplyOption) *CertificateConfigX509ConfigKeyUsage {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &CertificateConfigX509ConfigKeyUsage{}

	cDes.BaseKeyUsage = canonicalizeCertificateConfigX509ConfigKeyUsageBaseKeyUsage(des.BaseKeyUsage, initial.BaseKeyUsage, opts...)
	cDes.ExtendedKeyUsage = canonicalizeCertificateConfigX509ConfigKeyUsageExtendedKeyUsage(des.ExtendedKeyUsage, initial.ExtendedKeyUsage, opts...)
	cDes.UnknownExtendedKeyUsages = canonicalizeCertificateConfigX509ConfigKeyUsageUnknownExtendedKeyUsagesSlice(des.UnknownExtendedKeyUsages, initial.UnknownExtendedKeyUsages, opts...)

	return cDes
}

func canonicalizeCertificateConfigX509ConfigKeyUsageSlice(des, initial []CertificateConfigX509ConfigKeyUsage, opts ...dcl.ApplyOption) []CertificateConfigX509ConfigKeyUsage {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]CertificateConfigX509ConfigKeyUsage, 0, len(des))
		for _, d := range des {
			cd := canonicalizeCertificateConfigX509ConfigKeyUsage(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]CertificateConfigX509ConfigKeyUsage, 0, len(des))
	for i, d := range des {
		cd := canonicalizeCertificateConfigX509ConfigKeyUsage(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewCertificateConfigX509ConfigKeyUsage(c *Client, des, nw *CertificateConfigX509ConfigKeyUsage) *CertificateConfigX509ConfigKeyUsage {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for CertificateConfigX509ConfigKeyUsage while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	nw.BaseKeyUsage = canonicalizeNewCertificateConfigX509ConfigKeyUsageBaseKeyUsage(c, des.BaseKeyUsage, nw.BaseKeyUsage)
	nw.ExtendedKeyUsage = canonicalizeNewCertificateConfigX509ConfigKeyUsageExtendedKeyUsage(c, des.ExtendedKeyUsage, nw.ExtendedKeyUsage)
	nw.UnknownExtendedKeyUsages = canonicalizeNewCertificateConfigX509ConfigKeyUsageUnknownExtendedKeyUsagesSlice(c, des.UnknownExtendedKeyUsages, nw.UnknownExtendedKeyUsages)

	return nw
}

func canonicalizeNewCertificateConfigX509ConfigKeyUsageSet(c *Client, des, nw []CertificateConfigX509ConfigKeyUsage) []CertificateConfigX509ConfigKeyUsage {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []CertificateConfigX509ConfigKeyUsage
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareCertificateConfigX509ConfigKeyUsageNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewCertificateConfigX509ConfigKeyUsage(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewCertificateConfigX509ConfigKeyUsageSlice(c *Client, des, nw []CertificateConfigX509ConfigKeyUsage) []CertificateConfigX509ConfigKeyUsage {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []CertificateConfigX509ConfigKeyUsage
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewCertificateConfigX509ConfigKeyUsage(c, &d, &n))
	}

	return items
}

func canonicalizeCertificateConfigX509ConfigKeyUsageBaseKeyUsage(des, initial *CertificateConfigX509ConfigKeyUsageBaseKeyUsage, opts ...dcl.ApplyOption) *CertificateConfigX509ConfigKeyUsageBaseKeyUsage {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &CertificateConfigX509ConfigKeyUsageBaseKeyUsage{}

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

func canonicalizeCertificateConfigX509ConfigKeyUsageBaseKeyUsageSlice(des, initial []CertificateConfigX509ConfigKeyUsageBaseKeyUsage, opts ...dcl.ApplyOption) []CertificateConfigX509ConfigKeyUsageBaseKeyUsage {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]CertificateConfigX509ConfigKeyUsageBaseKeyUsage, 0, len(des))
		for _, d := range des {
			cd := canonicalizeCertificateConfigX509ConfigKeyUsageBaseKeyUsage(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]CertificateConfigX509ConfigKeyUsageBaseKeyUsage, 0, len(des))
	for i, d := range des {
		cd := canonicalizeCertificateConfigX509ConfigKeyUsageBaseKeyUsage(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewCertificateConfigX509ConfigKeyUsageBaseKeyUsage(c *Client, des, nw *CertificateConfigX509ConfigKeyUsageBaseKeyUsage) *CertificateConfigX509ConfigKeyUsageBaseKeyUsage {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for CertificateConfigX509ConfigKeyUsageBaseKeyUsage while comparing non-nil desired to nil actual.  Returning desired object.")
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

func canonicalizeNewCertificateConfigX509ConfigKeyUsageBaseKeyUsageSet(c *Client, des, nw []CertificateConfigX509ConfigKeyUsageBaseKeyUsage) []CertificateConfigX509ConfigKeyUsageBaseKeyUsage {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []CertificateConfigX509ConfigKeyUsageBaseKeyUsage
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareCertificateConfigX509ConfigKeyUsageBaseKeyUsageNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewCertificateConfigX509ConfigKeyUsageBaseKeyUsage(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewCertificateConfigX509ConfigKeyUsageBaseKeyUsageSlice(c *Client, des, nw []CertificateConfigX509ConfigKeyUsageBaseKeyUsage) []CertificateConfigX509ConfigKeyUsageBaseKeyUsage {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []CertificateConfigX509ConfigKeyUsageBaseKeyUsage
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewCertificateConfigX509ConfigKeyUsageBaseKeyUsage(c, &d, &n))
	}

	return items
}

func canonicalizeCertificateConfigX509ConfigKeyUsageExtendedKeyUsage(des, initial *CertificateConfigX509ConfigKeyUsageExtendedKeyUsage, opts ...dcl.ApplyOption) *CertificateConfigX509ConfigKeyUsageExtendedKeyUsage {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &CertificateConfigX509ConfigKeyUsageExtendedKeyUsage{}

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

func canonicalizeCertificateConfigX509ConfigKeyUsageExtendedKeyUsageSlice(des, initial []CertificateConfigX509ConfigKeyUsageExtendedKeyUsage, opts ...dcl.ApplyOption) []CertificateConfigX509ConfigKeyUsageExtendedKeyUsage {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]CertificateConfigX509ConfigKeyUsageExtendedKeyUsage, 0, len(des))
		for _, d := range des {
			cd := canonicalizeCertificateConfigX509ConfigKeyUsageExtendedKeyUsage(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]CertificateConfigX509ConfigKeyUsageExtendedKeyUsage, 0, len(des))
	for i, d := range des {
		cd := canonicalizeCertificateConfigX509ConfigKeyUsageExtendedKeyUsage(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewCertificateConfigX509ConfigKeyUsageExtendedKeyUsage(c *Client, des, nw *CertificateConfigX509ConfigKeyUsageExtendedKeyUsage) *CertificateConfigX509ConfigKeyUsageExtendedKeyUsage {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for CertificateConfigX509ConfigKeyUsageExtendedKeyUsage while comparing non-nil desired to nil actual.  Returning desired object.")
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

func canonicalizeNewCertificateConfigX509ConfigKeyUsageExtendedKeyUsageSet(c *Client, des, nw []CertificateConfigX509ConfigKeyUsageExtendedKeyUsage) []CertificateConfigX509ConfigKeyUsageExtendedKeyUsage {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []CertificateConfigX509ConfigKeyUsageExtendedKeyUsage
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareCertificateConfigX509ConfigKeyUsageExtendedKeyUsageNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewCertificateConfigX509ConfigKeyUsageExtendedKeyUsage(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewCertificateConfigX509ConfigKeyUsageExtendedKeyUsageSlice(c *Client, des, nw []CertificateConfigX509ConfigKeyUsageExtendedKeyUsage) []CertificateConfigX509ConfigKeyUsageExtendedKeyUsage {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []CertificateConfigX509ConfigKeyUsageExtendedKeyUsage
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewCertificateConfigX509ConfigKeyUsageExtendedKeyUsage(c, &d, &n))
	}

	return items
}

func canonicalizeCertificateConfigX509ConfigKeyUsageUnknownExtendedKeyUsages(des, initial *CertificateConfigX509ConfigKeyUsageUnknownExtendedKeyUsages, opts ...dcl.ApplyOption) *CertificateConfigX509ConfigKeyUsageUnknownExtendedKeyUsages {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &CertificateConfigX509ConfigKeyUsageUnknownExtendedKeyUsages{}

	if dcl.IsZeroValue(des.ObjectIdPath) || (dcl.IsEmptyValueIndirect(des.ObjectIdPath) && dcl.IsEmptyValueIndirect(initial.ObjectIdPath)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.ObjectIdPath = initial.ObjectIdPath
	} else {
		cDes.ObjectIdPath = des.ObjectIdPath
	}

	return cDes
}

func canonicalizeCertificateConfigX509ConfigKeyUsageUnknownExtendedKeyUsagesSlice(des, initial []CertificateConfigX509ConfigKeyUsageUnknownExtendedKeyUsages, opts ...dcl.ApplyOption) []CertificateConfigX509ConfigKeyUsageUnknownExtendedKeyUsages {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]CertificateConfigX509ConfigKeyUsageUnknownExtendedKeyUsages, 0, len(des))
		for _, d := range des {
			cd := canonicalizeCertificateConfigX509ConfigKeyUsageUnknownExtendedKeyUsages(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]CertificateConfigX509ConfigKeyUsageUnknownExtendedKeyUsages, 0, len(des))
	for i, d := range des {
		cd := canonicalizeCertificateConfigX509ConfigKeyUsageUnknownExtendedKeyUsages(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewCertificateConfigX509ConfigKeyUsageUnknownExtendedKeyUsages(c *Client, des, nw *CertificateConfigX509ConfigKeyUsageUnknownExtendedKeyUsages) *CertificateConfigX509ConfigKeyUsageUnknownExtendedKeyUsages {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for CertificateConfigX509ConfigKeyUsageUnknownExtendedKeyUsages while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	return nw
}

func canonicalizeNewCertificateConfigX509ConfigKeyUsageUnknownExtendedKeyUsagesSet(c *Client, des, nw []CertificateConfigX509ConfigKeyUsageUnknownExtendedKeyUsages) []CertificateConfigX509ConfigKeyUsageUnknownExtendedKeyUsages {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []CertificateConfigX509ConfigKeyUsageUnknownExtendedKeyUsages
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareCertificateConfigX509ConfigKeyUsageUnknownExtendedKeyUsagesNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewCertificateConfigX509ConfigKeyUsageUnknownExtendedKeyUsages(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewCertificateConfigX509ConfigKeyUsageUnknownExtendedKeyUsagesSlice(c *Client, des, nw []CertificateConfigX509ConfigKeyUsageUnknownExtendedKeyUsages) []CertificateConfigX509ConfigKeyUsageUnknownExtendedKeyUsages {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []CertificateConfigX509ConfigKeyUsageUnknownExtendedKeyUsages
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewCertificateConfigX509ConfigKeyUsageUnknownExtendedKeyUsages(c, &d, &n))
	}

	return items
}

func canonicalizeCertificateConfigX509ConfigCaOptions(des, initial *CertificateConfigX509ConfigCaOptions, opts ...dcl.ApplyOption) *CertificateConfigX509ConfigCaOptions {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &CertificateConfigX509ConfigCaOptions{}

	if dcl.BoolCanonicalize(des.IsCa, initial.IsCa) || dcl.IsZeroValue(des.IsCa) {
		cDes.IsCa = initial.IsCa
	} else {
		cDes.IsCa = des.IsCa
	}
	if dcl.BoolCanonicalize(des.NonCa, initial.NonCa) || dcl.IsZeroValue(des.NonCa) {
		cDes.NonCa = initial.NonCa
	} else {
		cDes.NonCa = des.NonCa
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

func canonicalizeCertificateConfigX509ConfigCaOptionsSlice(des, initial []CertificateConfigX509ConfigCaOptions, opts ...dcl.ApplyOption) []CertificateConfigX509ConfigCaOptions {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]CertificateConfigX509ConfigCaOptions, 0, len(des))
		for _, d := range des {
			cd := canonicalizeCertificateConfigX509ConfigCaOptions(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]CertificateConfigX509ConfigCaOptions, 0, len(des))
	for i, d := range des {
		cd := canonicalizeCertificateConfigX509ConfigCaOptions(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewCertificateConfigX509ConfigCaOptions(c *Client, des, nw *CertificateConfigX509ConfigCaOptions) *CertificateConfigX509ConfigCaOptions {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for CertificateConfigX509ConfigCaOptions while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.BoolCanonicalize(des.IsCa, nw.IsCa) {
		nw.IsCa = des.IsCa
	}
	if dcl.BoolCanonicalize(des.NonCa, nw.NonCa) {
		nw.NonCa = des.NonCa
	}
	if dcl.BoolCanonicalize(des.ZeroMaxIssuerPathLength, nw.ZeroMaxIssuerPathLength) {
		nw.ZeroMaxIssuerPathLength = des.ZeroMaxIssuerPathLength
	}

	return nw
}

func canonicalizeNewCertificateConfigX509ConfigCaOptionsSet(c *Client, des, nw []CertificateConfigX509ConfigCaOptions) []CertificateConfigX509ConfigCaOptions {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []CertificateConfigX509ConfigCaOptions
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareCertificateConfigX509ConfigCaOptionsNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewCertificateConfigX509ConfigCaOptions(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewCertificateConfigX509ConfigCaOptionsSlice(c *Client, des, nw []CertificateConfigX509ConfigCaOptions) []CertificateConfigX509ConfigCaOptions {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []CertificateConfigX509ConfigCaOptions
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewCertificateConfigX509ConfigCaOptions(c, &d, &n))
	}

	return items
}

func canonicalizeCertificateConfigX509ConfigPolicyIds(des, initial *CertificateConfigX509ConfigPolicyIds, opts ...dcl.ApplyOption) *CertificateConfigX509ConfigPolicyIds {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &CertificateConfigX509ConfigPolicyIds{}

	if dcl.IsZeroValue(des.ObjectIdPath) || (dcl.IsEmptyValueIndirect(des.ObjectIdPath) && dcl.IsEmptyValueIndirect(initial.ObjectIdPath)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.ObjectIdPath = initial.ObjectIdPath
	} else {
		cDes.ObjectIdPath = des.ObjectIdPath
	}

	return cDes
}

func canonicalizeCertificateConfigX509ConfigPolicyIdsSlice(des, initial []CertificateConfigX509ConfigPolicyIds, opts ...dcl.ApplyOption) []CertificateConfigX509ConfigPolicyIds {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]CertificateConfigX509ConfigPolicyIds, 0, len(des))
		for _, d := range des {
			cd := canonicalizeCertificateConfigX509ConfigPolicyIds(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]CertificateConfigX509ConfigPolicyIds, 0, len(des))
	for i, d := range des {
		cd := canonicalizeCertificateConfigX509ConfigPolicyIds(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewCertificateConfigX509ConfigPolicyIds(c *Client, des, nw *CertificateConfigX509ConfigPolicyIds) *CertificateConfigX509ConfigPolicyIds {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for CertificateConfigX509ConfigPolicyIds while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	return nw
}

func canonicalizeNewCertificateConfigX509ConfigPolicyIdsSet(c *Client, des, nw []CertificateConfigX509ConfigPolicyIds) []CertificateConfigX509ConfigPolicyIds {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []CertificateConfigX509ConfigPolicyIds
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareCertificateConfigX509ConfigPolicyIdsNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewCertificateConfigX509ConfigPolicyIds(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewCertificateConfigX509ConfigPolicyIdsSlice(c *Client, des, nw []CertificateConfigX509ConfigPolicyIds) []CertificateConfigX509ConfigPolicyIds {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []CertificateConfigX509ConfigPolicyIds
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewCertificateConfigX509ConfigPolicyIds(c, &d, &n))
	}

	return items
}

func canonicalizeCertificateConfigX509ConfigAdditionalExtensions(des, initial *CertificateConfigX509ConfigAdditionalExtensions, opts ...dcl.ApplyOption) *CertificateConfigX509ConfigAdditionalExtensions {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &CertificateConfigX509ConfigAdditionalExtensions{}

	cDes.ObjectId = canonicalizeCertificateConfigX509ConfigAdditionalExtensionsObjectId(des.ObjectId, initial.ObjectId, opts...)
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

func canonicalizeCertificateConfigX509ConfigAdditionalExtensionsSlice(des, initial []CertificateConfigX509ConfigAdditionalExtensions, opts ...dcl.ApplyOption) []CertificateConfigX509ConfigAdditionalExtensions {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]CertificateConfigX509ConfigAdditionalExtensions, 0, len(des))
		for _, d := range des {
			cd := canonicalizeCertificateConfigX509ConfigAdditionalExtensions(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]CertificateConfigX509ConfigAdditionalExtensions, 0, len(des))
	for i, d := range des {
		cd := canonicalizeCertificateConfigX509ConfigAdditionalExtensions(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewCertificateConfigX509ConfigAdditionalExtensions(c *Client, des, nw *CertificateConfigX509ConfigAdditionalExtensions) *CertificateConfigX509ConfigAdditionalExtensions {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for CertificateConfigX509ConfigAdditionalExtensions while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	nw.ObjectId = canonicalizeNewCertificateConfigX509ConfigAdditionalExtensionsObjectId(c, des.ObjectId, nw.ObjectId)
	if dcl.BoolCanonicalize(des.Critical, nw.Critical) {
		nw.Critical = des.Critical
	}
	if dcl.StringCanonicalize(des.Value, nw.Value) {
		nw.Value = des.Value
	}

	return nw
}

func canonicalizeNewCertificateConfigX509ConfigAdditionalExtensionsSet(c *Client, des, nw []CertificateConfigX509ConfigAdditionalExtensions) []CertificateConfigX509ConfigAdditionalExtensions {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []CertificateConfigX509ConfigAdditionalExtensions
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareCertificateConfigX509ConfigAdditionalExtensionsNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewCertificateConfigX509ConfigAdditionalExtensions(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewCertificateConfigX509ConfigAdditionalExtensionsSlice(c *Client, des, nw []CertificateConfigX509ConfigAdditionalExtensions) []CertificateConfigX509ConfigAdditionalExtensions {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []CertificateConfigX509ConfigAdditionalExtensions
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewCertificateConfigX509ConfigAdditionalExtensions(c, &d, &n))
	}

	return items
}

func canonicalizeCertificateConfigX509ConfigAdditionalExtensionsObjectId(des, initial *CertificateConfigX509ConfigAdditionalExtensionsObjectId, opts ...dcl.ApplyOption) *CertificateConfigX509ConfigAdditionalExtensionsObjectId {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &CertificateConfigX509ConfigAdditionalExtensionsObjectId{}

	if dcl.IsZeroValue(des.ObjectIdPath) || (dcl.IsEmptyValueIndirect(des.ObjectIdPath) && dcl.IsEmptyValueIndirect(initial.ObjectIdPath)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.ObjectIdPath = initial.ObjectIdPath
	} else {
		cDes.ObjectIdPath = des.ObjectIdPath
	}

	return cDes
}

func canonicalizeCertificateConfigX509ConfigAdditionalExtensionsObjectIdSlice(des, initial []CertificateConfigX509ConfigAdditionalExtensionsObjectId, opts ...dcl.ApplyOption) []CertificateConfigX509ConfigAdditionalExtensionsObjectId {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]CertificateConfigX509ConfigAdditionalExtensionsObjectId, 0, len(des))
		for _, d := range des {
			cd := canonicalizeCertificateConfigX509ConfigAdditionalExtensionsObjectId(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]CertificateConfigX509ConfigAdditionalExtensionsObjectId, 0, len(des))
	for i, d := range des {
		cd := canonicalizeCertificateConfigX509ConfigAdditionalExtensionsObjectId(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewCertificateConfigX509ConfigAdditionalExtensionsObjectId(c *Client, des, nw *CertificateConfigX509ConfigAdditionalExtensionsObjectId) *CertificateConfigX509ConfigAdditionalExtensionsObjectId {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for CertificateConfigX509ConfigAdditionalExtensionsObjectId while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	return nw
}

func canonicalizeNewCertificateConfigX509ConfigAdditionalExtensionsObjectIdSet(c *Client, des, nw []CertificateConfigX509ConfigAdditionalExtensionsObjectId) []CertificateConfigX509ConfigAdditionalExtensionsObjectId {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []CertificateConfigX509ConfigAdditionalExtensionsObjectId
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareCertificateConfigX509ConfigAdditionalExtensionsObjectIdNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewCertificateConfigX509ConfigAdditionalExtensionsObjectId(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewCertificateConfigX509ConfigAdditionalExtensionsObjectIdSlice(c *Client, des, nw []CertificateConfigX509ConfigAdditionalExtensionsObjectId) []CertificateConfigX509ConfigAdditionalExtensionsObjectId {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []CertificateConfigX509ConfigAdditionalExtensionsObjectId
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewCertificateConfigX509ConfigAdditionalExtensionsObjectId(c, &d, &n))
	}

	return items
}

func canonicalizeCertificateConfigPublicKey(des, initial *CertificateConfigPublicKey, opts ...dcl.ApplyOption) *CertificateConfigPublicKey {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &CertificateConfigPublicKey{}

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

func canonicalizeCertificateConfigPublicKeySlice(des, initial []CertificateConfigPublicKey, opts ...dcl.ApplyOption) []CertificateConfigPublicKey {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]CertificateConfigPublicKey, 0, len(des))
		for _, d := range des {
			cd := canonicalizeCertificateConfigPublicKey(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]CertificateConfigPublicKey, 0, len(des))
	for i, d := range des {
		cd := canonicalizeCertificateConfigPublicKey(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewCertificateConfigPublicKey(c *Client, des, nw *CertificateConfigPublicKey) *CertificateConfigPublicKey {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for CertificateConfigPublicKey while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringCanonicalize(des.Key, nw.Key) {
		nw.Key = des.Key
	}

	return nw
}

func canonicalizeNewCertificateConfigPublicKeySet(c *Client, des, nw []CertificateConfigPublicKey) []CertificateConfigPublicKey {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []CertificateConfigPublicKey
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareCertificateConfigPublicKeyNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewCertificateConfigPublicKey(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewCertificateConfigPublicKeySlice(c *Client, des, nw []CertificateConfigPublicKey) []CertificateConfigPublicKey {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []CertificateConfigPublicKey
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewCertificateConfigPublicKey(c, &d, &n))
	}

	return items
}

func canonicalizeCertificateRevocationDetails(des, initial *CertificateRevocationDetails, opts ...dcl.ApplyOption) *CertificateRevocationDetails {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &CertificateRevocationDetails{}

	if dcl.IsZeroValue(des.RevocationState) || (dcl.IsEmptyValueIndirect(des.RevocationState) && dcl.IsEmptyValueIndirect(initial.RevocationState)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.RevocationState = initial.RevocationState
	} else {
		cDes.RevocationState = des.RevocationState
	}
	if dcl.IsZeroValue(des.RevocationTime) || (dcl.IsEmptyValueIndirect(des.RevocationTime) && dcl.IsEmptyValueIndirect(initial.RevocationTime)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.RevocationTime = initial.RevocationTime
	} else {
		cDes.RevocationTime = des.RevocationTime
	}

	return cDes
}

func canonicalizeCertificateRevocationDetailsSlice(des, initial []CertificateRevocationDetails, opts ...dcl.ApplyOption) []CertificateRevocationDetails {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]CertificateRevocationDetails, 0, len(des))
		for _, d := range des {
			cd := canonicalizeCertificateRevocationDetails(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]CertificateRevocationDetails, 0, len(des))
	for i, d := range des {
		cd := canonicalizeCertificateRevocationDetails(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewCertificateRevocationDetails(c *Client, des, nw *CertificateRevocationDetails) *CertificateRevocationDetails {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for CertificateRevocationDetails while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	return nw
}

func canonicalizeNewCertificateRevocationDetailsSet(c *Client, des, nw []CertificateRevocationDetails) []CertificateRevocationDetails {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []CertificateRevocationDetails
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareCertificateRevocationDetailsNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewCertificateRevocationDetails(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewCertificateRevocationDetailsSlice(c *Client, des, nw []CertificateRevocationDetails) []CertificateRevocationDetails {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []CertificateRevocationDetails
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewCertificateRevocationDetails(c, &d, &n))
	}

	return items
}

func canonicalizeCertificateCertificateDescription(des, initial *CertificateCertificateDescription, opts ...dcl.ApplyOption) *CertificateCertificateDescription {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &CertificateCertificateDescription{}

	cDes.SubjectDescription = canonicalizeCertificateCertificateDescriptionSubjectDescription(des.SubjectDescription, initial.SubjectDescription, opts...)
	cDes.X509Description = canonicalizeCertificateCertificateDescriptionX509Description(des.X509Description, initial.X509Description, opts...)
	cDes.PublicKey = canonicalizeCertificateCertificateDescriptionPublicKey(des.PublicKey, initial.PublicKey, opts...)
	cDes.SubjectKeyId = canonicalizeCertificateCertificateDescriptionSubjectKeyId(des.SubjectKeyId, initial.SubjectKeyId, opts...)
	cDes.AuthorityKeyId = canonicalizeCertificateCertificateDescriptionAuthorityKeyId(des.AuthorityKeyId, initial.AuthorityKeyId, opts...)
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
	cDes.CertFingerprint = canonicalizeCertificateCertificateDescriptionCertFingerprint(des.CertFingerprint, initial.CertFingerprint, opts...)

	return cDes
}

func canonicalizeCertificateCertificateDescriptionSlice(des, initial []CertificateCertificateDescription, opts ...dcl.ApplyOption) []CertificateCertificateDescription {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]CertificateCertificateDescription, 0, len(des))
		for _, d := range des {
			cd := canonicalizeCertificateCertificateDescription(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]CertificateCertificateDescription, 0, len(des))
	for i, d := range des {
		cd := canonicalizeCertificateCertificateDescription(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewCertificateCertificateDescription(c *Client, des, nw *CertificateCertificateDescription) *CertificateCertificateDescription {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for CertificateCertificateDescription while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	nw.SubjectDescription = canonicalizeNewCertificateCertificateDescriptionSubjectDescription(c, des.SubjectDescription, nw.SubjectDescription)
	nw.X509Description = canonicalizeNewCertificateCertificateDescriptionX509Description(c, des.X509Description, nw.X509Description)
	nw.PublicKey = canonicalizeNewCertificateCertificateDescriptionPublicKey(c, des.PublicKey, nw.PublicKey)
	nw.SubjectKeyId = canonicalizeNewCertificateCertificateDescriptionSubjectKeyId(c, des.SubjectKeyId, nw.SubjectKeyId)
	nw.AuthorityKeyId = canonicalizeNewCertificateCertificateDescriptionAuthorityKeyId(c, des.AuthorityKeyId, nw.AuthorityKeyId)
	if dcl.StringArrayCanonicalize(des.CrlDistributionPoints, nw.CrlDistributionPoints) {
		nw.CrlDistributionPoints = des.CrlDistributionPoints
	}
	if dcl.StringArrayCanonicalize(des.AiaIssuingCertificateUrls, nw.AiaIssuingCertificateUrls) {
		nw.AiaIssuingCertificateUrls = des.AiaIssuingCertificateUrls
	}
	nw.CertFingerprint = canonicalizeNewCertificateCertificateDescriptionCertFingerprint(c, des.CertFingerprint, nw.CertFingerprint)

	return nw
}

func canonicalizeNewCertificateCertificateDescriptionSet(c *Client, des, nw []CertificateCertificateDescription) []CertificateCertificateDescription {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []CertificateCertificateDescription
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareCertificateCertificateDescriptionNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewCertificateCertificateDescription(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewCertificateCertificateDescriptionSlice(c *Client, des, nw []CertificateCertificateDescription) []CertificateCertificateDescription {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []CertificateCertificateDescription
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewCertificateCertificateDescription(c, &d, &n))
	}

	return items
}

func canonicalizeCertificateCertificateDescriptionSubjectDescription(des, initial *CertificateCertificateDescriptionSubjectDescription, opts ...dcl.ApplyOption) *CertificateCertificateDescriptionSubjectDescription {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &CertificateCertificateDescriptionSubjectDescription{}

	cDes.Subject = canonicalizeCertificateCertificateDescriptionSubjectDescriptionSubject(des.Subject, initial.Subject, opts...)
	cDes.SubjectAltName = canonicalizeCertificateCertificateDescriptionSubjectDescriptionSubjectAltName(des.SubjectAltName, initial.SubjectAltName, opts...)
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

func canonicalizeCertificateCertificateDescriptionSubjectDescriptionSlice(des, initial []CertificateCertificateDescriptionSubjectDescription, opts ...dcl.ApplyOption) []CertificateCertificateDescriptionSubjectDescription {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]CertificateCertificateDescriptionSubjectDescription, 0, len(des))
		for _, d := range des {
			cd := canonicalizeCertificateCertificateDescriptionSubjectDescription(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]CertificateCertificateDescriptionSubjectDescription, 0, len(des))
	for i, d := range des {
		cd := canonicalizeCertificateCertificateDescriptionSubjectDescription(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewCertificateCertificateDescriptionSubjectDescription(c *Client, des, nw *CertificateCertificateDescriptionSubjectDescription) *CertificateCertificateDescriptionSubjectDescription {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for CertificateCertificateDescriptionSubjectDescription while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	nw.Subject = canonicalizeNewCertificateCertificateDescriptionSubjectDescriptionSubject(c, des.Subject, nw.Subject)
	nw.SubjectAltName = canonicalizeNewCertificateCertificateDescriptionSubjectDescriptionSubjectAltName(c, des.SubjectAltName, nw.SubjectAltName)
	if dcl.StringCanonicalize(des.HexSerialNumber, nw.HexSerialNumber) {
		nw.HexSerialNumber = des.HexSerialNumber
	}
	if dcl.StringCanonicalize(des.Lifetime, nw.Lifetime) {
		nw.Lifetime = des.Lifetime
	}

	return nw
}

func canonicalizeNewCertificateCertificateDescriptionSubjectDescriptionSet(c *Client, des, nw []CertificateCertificateDescriptionSubjectDescription) []CertificateCertificateDescriptionSubjectDescription {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []CertificateCertificateDescriptionSubjectDescription
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareCertificateCertificateDescriptionSubjectDescriptionNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewCertificateCertificateDescriptionSubjectDescription(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewCertificateCertificateDescriptionSubjectDescriptionSlice(c *Client, des, nw []CertificateCertificateDescriptionSubjectDescription) []CertificateCertificateDescriptionSubjectDescription {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []CertificateCertificateDescriptionSubjectDescription
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewCertificateCertificateDescriptionSubjectDescription(c, &d, &n))
	}

	return items
}

func canonicalizeCertificateCertificateDescriptionSubjectDescriptionSubject(des, initial *CertificateCertificateDescriptionSubjectDescriptionSubject, opts ...dcl.ApplyOption) *CertificateCertificateDescriptionSubjectDescriptionSubject {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &CertificateCertificateDescriptionSubjectDescriptionSubject{}

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

func canonicalizeCertificateCertificateDescriptionSubjectDescriptionSubjectSlice(des, initial []CertificateCertificateDescriptionSubjectDescriptionSubject, opts ...dcl.ApplyOption) []CertificateCertificateDescriptionSubjectDescriptionSubject {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]CertificateCertificateDescriptionSubjectDescriptionSubject, 0, len(des))
		for _, d := range des {
			cd := canonicalizeCertificateCertificateDescriptionSubjectDescriptionSubject(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]CertificateCertificateDescriptionSubjectDescriptionSubject, 0, len(des))
	for i, d := range des {
		cd := canonicalizeCertificateCertificateDescriptionSubjectDescriptionSubject(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewCertificateCertificateDescriptionSubjectDescriptionSubject(c *Client, des, nw *CertificateCertificateDescriptionSubjectDescriptionSubject) *CertificateCertificateDescriptionSubjectDescriptionSubject {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for CertificateCertificateDescriptionSubjectDescriptionSubject while comparing non-nil desired to nil actual.  Returning desired object.")
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

func canonicalizeNewCertificateCertificateDescriptionSubjectDescriptionSubjectSet(c *Client, des, nw []CertificateCertificateDescriptionSubjectDescriptionSubject) []CertificateCertificateDescriptionSubjectDescriptionSubject {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []CertificateCertificateDescriptionSubjectDescriptionSubject
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareCertificateCertificateDescriptionSubjectDescriptionSubjectNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewCertificateCertificateDescriptionSubjectDescriptionSubject(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewCertificateCertificateDescriptionSubjectDescriptionSubjectSlice(c *Client, des, nw []CertificateCertificateDescriptionSubjectDescriptionSubject) []CertificateCertificateDescriptionSubjectDescriptionSubject {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []CertificateCertificateDescriptionSubjectDescriptionSubject
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewCertificateCertificateDescriptionSubjectDescriptionSubject(c, &d, &n))
	}

	return items
}

func canonicalizeCertificateCertificateDescriptionSubjectDescriptionSubjectAltName(des, initial *CertificateCertificateDescriptionSubjectDescriptionSubjectAltName, opts ...dcl.ApplyOption) *CertificateCertificateDescriptionSubjectDescriptionSubjectAltName {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &CertificateCertificateDescriptionSubjectDescriptionSubjectAltName{}

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
	cDes.CustomSans = canonicalizeCertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSansSlice(des.CustomSans, initial.CustomSans, opts...)

	return cDes
}

func canonicalizeCertificateCertificateDescriptionSubjectDescriptionSubjectAltNameSlice(des, initial []CertificateCertificateDescriptionSubjectDescriptionSubjectAltName, opts ...dcl.ApplyOption) []CertificateCertificateDescriptionSubjectDescriptionSubjectAltName {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]CertificateCertificateDescriptionSubjectDescriptionSubjectAltName, 0, len(des))
		for _, d := range des {
			cd := canonicalizeCertificateCertificateDescriptionSubjectDescriptionSubjectAltName(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]CertificateCertificateDescriptionSubjectDescriptionSubjectAltName, 0, len(des))
	for i, d := range des {
		cd := canonicalizeCertificateCertificateDescriptionSubjectDescriptionSubjectAltName(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewCertificateCertificateDescriptionSubjectDescriptionSubjectAltName(c *Client, des, nw *CertificateCertificateDescriptionSubjectDescriptionSubjectAltName) *CertificateCertificateDescriptionSubjectDescriptionSubjectAltName {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for CertificateCertificateDescriptionSubjectDescriptionSubjectAltName while comparing non-nil desired to nil actual.  Returning desired object.")
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
	nw.CustomSans = canonicalizeNewCertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSansSlice(c, des.CustomSans, nw.CustomSans)

	return nw
}

func canonicalizeNewCertificateCertificateDescriptionSubjectDescriptionSubjectAltNameSet(c *Client, des, nw []CertificateCertificateDescriptionSubjectDescriptionSubjectAltName) []CertificateCertificateDescriptionSubjectDescriptionSubjectAltName {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []CertificateCertificateDescriptionSubjectDescriptionSubjectAltName
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareCertificateCertificateDescriptionSubjectDescriptionSubjectAltNameNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewCertificateCertificateDescriptionSubjectDescriptionSubjectAltName(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewCertificateCertificateDescriptionSubjectDescriptionSubjectAltNameSlice(c *Client, des, nw []CertificateCertificateDescriptionSubjectDescriptionSubjectAltName) []CertificateCertificateDescriptionSubjectDescriptionSubjectAltName {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []CertificateCertificateDescriptionSubjectDescriptionSubjectAltName
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewCertificateCertificateDescriptionSubjectDescriptionSubjectAltName(c, &d, &n))
	}

	return items
}

func canonicalizeCertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSans(des, initial *CertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSans, opts ...dcl.ApplyOption) *CertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSans {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &CertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSans{}

	cDes.ObjectId = canonicalizeCertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSansObjectId(des.ObjectId, initial.ObjectId, opts...)
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

func canonicalizeCertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSansSlice(des, initial []CertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSans, opts ...dcl.ApplyOption) []CertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSans {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]CertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSans, 0, len(des))
		for _, d := range des {
			cd := canonicalizeCertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSans(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]CertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSans, 0, len(des))
	for i, d := range des {
		cd := canonicalizeCertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSans(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewCertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSans(c *Client, des, nw *CertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSans) *CertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSans {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for CertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSans while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	nw.ObjectId = canonicalizeNewCertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSansObjectId(c, des.ObjectId, nw.ObjectId)
	if dcl.BoolCanonicalize(des.Critical, nw.Critical) {
		nw.Critical = des.Critical
	}
	if dcl.StringCanonicalize(des.Value, nw.Value) {
		nw.Value = des.Value
	}

	return nw
}

func canonicalizeNewCertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSansSet(c *Client, des, nw []CertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSans) []CertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSans {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []CertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSans
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareCertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSansNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewCertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSans(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewCertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSansSlice(c *Client, des, nw []CertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSans) []CertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSans {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []CertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSans
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewCertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSans(c, &d, &n))
	}

	return items
}

func canonicalizeCertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSansObjectId(des, initial *CertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSansObjectId, opts ...dcl.ApplyOption) *CertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSansObjectId {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &CertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSansObjectId{}

	if dcl.IsZeroValue(des.ObjectIdPath) || (dcl.IsEmptyValueIndirect(des.ObjectIdPath) && dcl.IsEmptyValueIndirect(initial.ObjectIdPath)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.ObjectIdPath = initial.ObjectIdPath
	} else {
		cDes.ObjectIdPath = des.ObjectIdPath
	}

	return cDes
}

func canonicalizeCertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSansObjectIdSlice(des, initial []CertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSansObjectId, opts ...dcl.ApplyOption) []CertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSansObjectId {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]CertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSansObjectId, 0, len(des))
		for _, d := range des {
			cd := canonicalizeCertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSansObjectId(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]CertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSansObjectId, 0, len(des))
	for i, d := range des {
		cd := canonicalizeCertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSansObjectId(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewCertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSansObjectId(c *Client, des, nw *CertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSansObjectId) *CertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSansObjectId {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for CertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSansObjectId while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	return nw
}

func canonicalizeNewCertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSansObjectIdSet(c *Client, des, nw []CertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSansObjectId) []CertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSansObjectId {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []CertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSansObjectId
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareCertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSansObjectIdNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewCertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSansObjectId(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewCertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSansObjectIdSlice(c *Client, des, nw []CertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSansObjectId) []CertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSansObjectId {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []CertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSansObjectId
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewCertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSansObjectId(c, &d, &n))
	}

	return items
}

func canonicalizeCertificateCertificateDescriptionX509Description(des, initial *CertificateCertificateDescriptionX509Description, opts ...dcl.ApplyOption) *CertificateCertificateDescriptionX509Description {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &CertificateCertificateDescriptionX509Description{}

	cDes.KeyUsage = canonicalizeCertificateCertificateDescriptionX509DescriptionKeyUsage(des.KeyUsage, initial.KeyUsage, opts...)
	cDes.CaOptions = canonicalizeCertificateCertificateDescriptionX509DescriptionCaOptions(des.CaOptions, initial.CaOptions, opts...)
	cDes.PolicyIds = canonicalizeCertificateCertificateDescriptionX509DescriptionPolicyIdsSlice(des.PolicyIds, initial.PolicyIds, opts...)
	if dcl.StringArrayCanonicalize(des.AiaOcspServers, initial.AiaOcspServers) {
		cDes.AiaOcspServers = initial.AiaOcspServers
	} else {
		cDes.AiaOcspServers = des.AiaOcspServers
	}
	cDes.AdditionalExtensions = canonicalizeCertificateCertificateDescriptionX509DescriptionAdditionalExtensionsSlice(des.AdditionalExtensions, initial.AdditionalExtensions, opts...)

	return cDes
}

func canonicalizeCertificateCertificateDescriptionX509DescriptionSlice(des, initial []CertificateCertificateDescriptionX509Description, opts ...dcl.ApplyOption) []CertificateCertificateDescriptionX509Description {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]CertificateCertificateDescriptionX509Description, 0, len(des))
		for _, d := range des {
			cd := canonicalizeCertificateCertificateDescriptionX509Description(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]CertificateCertificateDescriptionX509Description, 0, len(des))
	for i, d := range des {
		cd := canonicalizeCertificateCertificateDescriptionX509Description(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewCertificateCertificateDescriptionX509Description(c *Client, des, nw *CertificateCertificateDescriptionX509Description) *CertificateCertificateDescriptionX509Description {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for CertificateCertificateDescriptionX509Description while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	nw.KeyUsage = canonicalizeNewCertificateCertificateDescriptionX509DescriptionKeyUsage(c, des.KeyUsage, nw.KeyUsage)
	nw.CaOptions = canonicalizeNewCertificateCertificateDescriptionX509DescriptionCaOptions(c, des.CaOptions, nw.CaOptions)
	nw.PolicyIds = canonicalizeNewCertificateCertificateDescriptionX509DescriptionPolicyIdsSlice(c, des.PolicyIds, nw.PolicyIds)
	if dcl.StringArrayCanonicalize(des.AiaOcspServers, nw.AiaOcspServers) {
		nw.AiaOcspServers = des.AiaOcspServers
	}
	nw.AdditionalExtensions = canonicalizeNewCertificateCertificateDescriptionX509DescriptionAdditionalExtensionsSlice(c, des.AdditionalExtensions, nw.AdditionalExtensions)

	return nw
}

func canonicalizeNewCertificateCertificateDescriptionX509DescriptionSet(c *Client, des, nw []CertificateCertificateDescriptionX509Description) []CertificateCertificateDescriptionX509Description {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []CertificateCertificateDescriptionX509Description
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareCertificateCertificateDescriptionX509DescriptionNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewCertificateCertificateDescriptionX509Description(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewCertificateCertificateDescriptionX509DescriptionSlice(c *Client, des, nw []CertificateCertificateDescriptionX509Description) []CertificateCertificateDescriptionX509Description {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []CertificateCertificateDescriptionX509Description
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewCertificateCertificateDescriptionX509Description(c, &d, &n))
	}

	return items
}

func canonicalizeCertificateCertificateDescriptionX509DescriptionKeyUsage(des, initial *CertificateCertificateDescriptionX509DescriptionKeyUsage, opts ...dcl.ApplyOption) *CertificateCertificateDescriptionX509DescriptionKeyUsage {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &CertificateCertificateDescriptionX509DescriptionKeyUsage{}

	cDes.BaseKeyUsage = canonicalizeCertificateCertificateDescriptionX509DescriptionKeyUsageBaseKeyUsage(des.BaseKeyUsage, initial.BaseKeyUsage, opts...)
	cDes.ExtendedKeyUsage = canonicalizeCertificateCertificateDescriptionX509DescriptionKeyUsageExtendedKeyUsage(des.ExtendedKeyUsage, initial.ExtendedKeyUsage, opts...)
	cDes.UnknownExtendedKeyUsages = canonicalizeCertificateCertificateDescriptionX509DescriptionKeyUsageUnknownExtendedKeyUsagesSlice(des.UnknownExtendedKeyUsages, initial.UnknownExtendedKeyUsages, opts...)

	return cDes
}

func canonicalizeCertificateCertificateDescriptionX509DescriptionKeyUsageSlice(des, initial []CertificateCertificateDescriptionX509DescriptionKeyUsage, opts ...dcl.ApplyOption) []CertificateCertificateDescriptionX509DescriptionKeyUsage {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]CertificateCertificateDescriptionX509DescriptionKeyUsage, 0, len(des))
		for _, d := range des {
			cd := canonicalizeCertificateCertificateDescriptionX509DescriptionKeyUsage(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]CertificateCertificateDescriptionX509DescriptionKeyUsage, 0, len(des))
	for i, d := range des {
		cd := canonicalizeCertificateCertificateDescriptionX509DescriptionKeyUsage(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewCertificateCertificateDescriptionX509DescriptionKeyUsage(c *Client, des, nw *CertificateCertificateDescriptionX509DescriptionKeyUsage) *CertificateCertificateDescriptionX509DescriptionKeyUsage {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for CertificateCertificateDescriptionX509DescriptionKeyUsage while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	nw.BaseKeyUsage = canonicalizeNewCertificateCertificateDescriptionX509DescriptionKeyUsageBaseKeyUsage(c, des.BaseKeyUsage, nw.BaseKeyUsage)
	nw.ExtendedKeyUsage = canonicalizeNewCertificateCertificateDescriptionX509DescriptionKeyUsageExtendedKeyUsage(c, des.ExtendedKeyUsage, nw.ExtendedKeyUsage)
	nw.UnknownExtendedKeyUsages = canonicalizeNewCertificateCertificateDescriptionX509DescriptionKeyUsageUnknownExtendedKeyUsagesSlice(c, des.UnknownExtendedKeyUsages, nw.UnknownExtendedKeyUsages)

	return nw
}

func canonicalizeNewCertificateCertificateDescriptionX509DescriptionKeyUsageSet(c *Client, des, nw []CertificateCertificateDescriptionX509DescriptionKeyUsage) []CertificateCertificateDescriptionX509DescriptionKeyUsage {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []CertificateCertificateDescriptionX509DescriptionKeyUsage
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareCertificateCertificateDescriptionX509DescriptionKeyUsageNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewCertificateCertificateDescriptionX509DescriptionKeyUsage(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewCertificateCertificateDescriptionX509DescriptionKeyUsageSlice(c *Client, des, nw []CertificateCertificateDescriptionX509DescriptionKeyUsage) []CertificateCertificateDescriptionX509DescriptionKeyUsage {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []CertificateCertificateDescriptionX509DescriptionKeyUsage
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewCertificateCertificateDescriptionX509DescriptionKeyUsage(c, &d, &n))
	}

	return items
}

func canonicalizeCertificateCertificateDescriptionX509DescriptionKeyUsageBaseKeyUsage(des, initial *CertificateCertificateDescriptionX509DescriptionKeyUsageBaseKeyUsage, opts ...dcl.ApplyOption) *CertificateCertificateDescriptionX509DescriptionKeyUsageBaseKeyUsage {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &CertificateCertificateDescriptionX509DescriptionKeyUsageBaseKeyUsage{}

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

func canonicalizeCertificateCertificateDescriptionX509DescriptionKeyUsageBaseKeyUsageSlice(des, initial []CertificateCertificateDescriptionX509DescriptionKeyUsageBaseKeyUsage, opts ...dcl.ApplyOption) []CertificateCertificateDescriptionX509DescriptionKeyUsageBaseKeyUsage {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]CertificateCertificateDescriptionX509DescriptionKeyUsageBaseKeyUsage, 0, len(des))
		for _, d := range des {
			cd := canonicalizeCertificateCertificateDescriptionX509DescriptionKeyUsageBaseKeyUsage(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]CertificateCertificateDescriptionX509DescriptionKeyUsageBaseKeyUsage, 0, len(des))
	for i, d := range des {
		cd := canonicalizeCertificateCertificateDescriptionX509DescriptionKeyUsageBaseKeyUsage(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewCertificateCertificateDescriptionX509DescriptionKeyUsageBaseKeyUsage(c *Client, des, nw *CertificateCertificateDescriptionX509DescriptionKeyUsageBaseKeyUsage) *CertificateCertificateDescriptionX509DescriptionKeyUsageBaseKeyUsage {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for CertificateCertificateDescriptionX509DescriptionKeyUsageBaseKeyUsage while comparing non-nil desired to nil actual.  Returning desired object.")
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

func canonicalizeNewCertificateCertificateDescriptionX509DescriptionKeyUsageBaseKeyUsageSet(c *Client, des, nw []CertificateCertificateDescriptionX509DescriptionKeyUsageBaseKeyUsage) []CertificateCertificateDescriptionX509DescriptionKeyUsageBaseKeyUsage {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []CertificateCertificateDescriptionX509DescriptionKeyUsageBaseKeyUsage
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareCertificateCertificateDescriptionX509DescriptionKeyUsageBaseKeyUsageNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewCertificateCertificateDescriptionX509DescriptionKeyUsageBaseKeyUsage(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewCertificateCertificateDescriptionX509DescriptionKeyUsageBaseKeyUsageSlice(c *Client, des, nw []CertificateCertificateDescriptionX509DescriptionKeyUsageBaseKeyUsage) []CertificateCertificateDescriptionX509DescriptionKeyUsageBaseKeyUsage {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []CertificateCertificateDescriptionX509DescriptionKeyUsageBaseKeyUsage
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewCertificateCertificateDescriptionX509DescriptionKeyUsageBaseKeyUsage(c, &d, &n))
	}

	return items
}

func canonicalizeCertificateCertificateDescriptionX509DescriptionKeyUsageExtendedKeyUsage(des, initial *CertificateCertificateDescriptionX509DescriptionKeyUsageExtendedKeyUsage, opts ...dcl.ApplyOption) *CertificateCertificateDescriptionX509DescriptionKeyUsageExtendedKeyUsage {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &CertificateCertificateDescriptionX509DescriptionKeyUsageExtendedKeyUsage{}

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

func canonicalizeCertificateCertificateDescriptionX509DescriptionKeyUsageExtendedKeyUsageSlice(des, initial []CertificateCertificateDescriptionX509DescriptionKeyUsageExtendedKeyUsage, opts ...dcl.ApplyOption) []CertificateCertificateDescriptionX509DescriptionKeyUsageExtendedKeyUsage {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]CertificateCertificateDescriptionX509DescriptionKeyUsageExtendedKeyUsage, 0, len(des))
		for _, d := range des {
			cd := canonicalizeCertificateCertificateDescriptionX509DescriptionKeyUsageExtendedKeyUsage(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]CertificateCertificateDescriptionX509DescriptionKeyUsageExtendedKeyUsage, 0, len(des))
	for i, d := range des {
		cd := canonicalizeCertificateCertificateDescriptionX509DescriptionKeyUsageExtendedKeyUsage(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewCertificateCertificateDescriptionX509DescriptionKeyUsageExtendedKeyUsage(c *Client, des, nw *CertificateCertificateDescriptionX509DescriptionKeyUsageExtendedKeyUsage) *CertificateCertificateDescriptionX509DescriptionKeyUsageExtendedKeyUsage {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for CertificateCertificateDescriptionX509DescriptionKeyUsageExtendedKeyUsage while comparing non-nil desired to nil actual.  Returning desired object.")
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

func canonicalizeNewCertificateCertificateDescriptionX509DescriptionKeyUsageExtendedKeyUsageSet(c *Client, des, nw []CertificateCertificateDescriptionX509DescriptionKeyUsageExtendedKeyUsage) []CertificateCertificateDescriptionX509DescriptionKeyUsageExtendedKeyUsage {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []CertificateCertificateDescriptionX509DescriptionKeyUsageExtendedKeyUsage
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareCertificateCertificateDescriptionX509DescriptionKeyUsageExtendedKeyUsageNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewCertificateCertificateDescriptionX509DescriptionKeyUsageExtendedKeyUsage(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewCertificateCertificateDescriptionX509DescriptionKeyUsageExtendedKeyUsageSlice(c *Client, des, nw []CertificateCertificateDescriptionX509DescriptionKeyUsageExtendedKeyUsage) []CertificateCertificateDescriptionX509DescriptionKeyUsageExtendedKeyUsage {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []CertificateCertificateDescriptionX509DescriptionKeyUsageExtendedKeyUsage
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewCertificateCertificateDescriptionX509DescriptionKeyUsageExtendedKeyUsage(c, &d, &n))
	}

	return items
}

func canonicalizeCertificateCertificateDescriptionX509DescriptionKeyUsageUnknownExtendedKeyUsages(des, initial *CertificateCertificateDescriptionX509DescriptionKeyUsageUnknownExtendedKeyUsages, opts ...dcl.ApplyOption) *CertificateCertificateDescriptionX509DescriptionKeyUsageUnknownExtendedKeyUsages {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &CertificateCertificateDescriptionX509DescriptionKeyUsageUnknownExtendedKeyUsages{}

	if dcl.IsZeroValue(des.ObjectIdPath) || (dcl.IsEmptyValueIndirect(des.ObjectIdPath) && dcl.IsEmptyValueIndirect(initial.ObjectIdPath)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.ObjectIdPath = initial.ObjectIdPath
	} else {
		cDes.ObjectIdPath = des.ObjectIdPath
	}

	return cDes
}

func canonicalizeCertificateCertificateDescriptionX509DescriptionKeyUsageUnknownExtendedKeyUsagesSlice(des, initial []CertificateCertificateDescriptionX509DescriptionKeyUsageUnknownExtendedKeyUsages, opts ...dcl.ApplyOption) []CertificateCertificateDescriptionX509DescriptionKeyUsageUnknownExtendedKeyUsages {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]CertificateCertificateDescriptionX509DescriptionKeyUsageUnknownExtendedKeyUsages, 0, len(des))
		for _, d := range des {
			cd := canonicalizeCertificateCertificateDescriptionX509DescriptionKeyUsageUnknownExtendedKeyUsages(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]CertificateCertificateDescriptionX509DescriptionKeyUsageUnknownExtendedKeyUsages, 0, len(des))
	for i, d := range des {
		cd := canonicalizeCertificateCertificateDescriptionX509DescriptionKeyUsageUnknownExtendedKeyUsages(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewCertificateCertificateDescriptionX509DescriptionKeyUsageUnknownExtendedKeyUsages(c *Client, des, nw *CertificateCertificateDescriptionX509DescriptionKeyUsageUnknownExtendedKeyUsages) *CertificateCertificateDescriptionX509DescriptionKeyUsageUnknownExtendedKeyUsages {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for CertificateCertificateDescriptionX509DescriptionKeyUsageUnknownExtendedKeyUsages while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	return nw
}

func canonicalizeNewCertificateCertificateDescriptionX509DescriptionKeyUsageUnknownExtendedKeyUsagesSet(c *Client, des, nw []CertificateCertificateDescriptionX509DescriptionKeyUsageUnknownExtendedKeyUsages) []CertificateCertificateDescriptionX509DescriptionKeyUsageUnknownExtendedKeyUsages {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []CertificateCertificateDescriptionX509DescriptionKeyUsageUnknownExtendedKeyUsages
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareCertificateCertificateDescriptionX509DescriptionKeyUsageUnknownExtendedKeyUsagesNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewCertificateCertificateDescriptionX509DescriptionKeyUsageUnknownExtendedKeyUsages(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewCertificateCertificateDescriptionX509DescriptionKeyUsageUnknownExtendedKeyUsagesSlice(c *Client, des, nw []CertificateCertificateDescriptionX509DescriptionKeyUsageUnknownExtendedKeyUsages) []CertificateCertificateDescriptionX509DescriptionKeyUsageUnknownExtendedKeyUsages {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []CertificateCertificateDescriptionX509DescriptionKeyUsageUnknownExtendedKeyUsages
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewCertificateCertificateDescriptionX509DescriptionKeyUsageUnknownExtendedKeyUsages(c, &d, &n))
	}

	return items
}

func canonicalizeCertificateCertificateDescriptionX509DescriptionCaOptions(des, initial *CertificateCertificateDescriptionX509DescriptionCaOptions, opts ...dcl.ApplyOption) *CertificateCertificateDescriptionX509DescriptionCaOptions {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if des.IsCa != nil || (initial != nil && initial.IsCa != nil) {
		// Check if anything else is set.
		if dcl.AnySet() {
			des.IsCa = nil
			if initial != nil {
				initial.IsCa = nil
			}
		}
	}

	if des.MaxIssuerPathLength != nil || (initial != nil && initial.MaxIssuerPathLength != nil) {
		// Check if anything else is set.
		if dcl.AnySet() {
			des.MaxIssuerPathLength = nil
			if initial != nil {
				initial.MaxIssuerPathLength = nil
			}
		}
	}

	if initial == nil {
		return des
	}

	cDes := &CertificateCertificateDescriptionX509DescriptionCaOptions{}

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

func canonicalizeCertificateCertificateDescriptionX509DescriptionCaOptionsSlice(des, initial []CertificateCertificateDescriptionX509DescriptionCaOptions, opts ...dcl.ApplyOption) []CertificateCertificateDescriptionX509DescriptionCaOptions {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]CertificateCertificateDescriptionX509DescriptionCaOptions, 0, len(des))
		for _, d := range des {
			cd := canonicalizeCertificateCertificateDescriptionX509DescriptionCaOptions(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]CertificateCertificateDescriptionX509DescriptionCaOptions, 0, len(des))
	for i, d := range des {
		cd := canonicalizeCertificateCertificateDescriptionX509DescriptionCaOptions(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewCertificateCertificateDescriptionX509DescriptionCaOptions(c *Client, des, nw *CertificateCertificateDescriptionX509DescriptionCaOptions) *CertificateCertificateDescriptionX509DescriptionCaOptions {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for CertificateCertificateDescriptionX509DescriptionCaOptions while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.BoolCanonicalize(des.IsCa, nw.IsCa) {
		nw.IsCa = des.IsCa
	}

	return nw
}

func canonicalizeNewCertificateCertificateDescriptionX509DescriptionCaOptionsSet(c *Client, des, nw []CertificateCertificateDescriptionX509DescriptionCaOptions) []CertificateCertificateDescriptionX509DescriptionCaOptions {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []CertificateCertificateDescriptionX509DescriptionCaOptions
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareCertificateCertificateDescriptionX509DescriptionCaOptionsNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewCertificateCertificateDescriptionX509DescriptionCaOptions(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewCertificateCertificateDescriptionX509DescriptionCaOptionsSlice(c *Client, des, nw []CertificateCertificateDescriptionX509DescriptionCaOptions) []CertificateCertificateDescriptionX509DescriptionCaOptions {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []CertificateCertificateDescriptionX509DescriptionCaOptions
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewCertificateCertificateDescriptionX509DescriptionCaOptions(c, &d, &n))
	}

	return items
}

func canonicalizeCertificateCertificateDescriptionX509DescriptionPolicyIds(des, initial *CertificateCertificateDescriptionX509DescriptionPolicyIds, opts ...dcl.ApplyOption) *CertificateCertificateDescriptionX509DescriptionPolicyIds {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &CertificateCertificateDescriptionX509DescriptionPolicyIds{}

	if dcl.IsZeroValue(des.ObjectIdPath) || (dcl.IsEmptyValueIndirect(des.ObjectIdPath) && dcl.IsEmptyValueIndirect(initial.ObjectIdPath)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.ObjectIdPath = initial.ObjectIdPath
	} else {
		cDes.ObjectIdPath = des.ObjectIdPath
	}

	return cDes
}

func canonicalizeCertificateCertificateDescriptionX509DescriptionPolicyIdsSlice(des, initial []CertificateCertificateDescriptionX509DescriptionPolicyIds, opts ...dcl.ApplyOption) []CertificateCertificateDescriptionX509DescriptionPolicyIds {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]CertificateCertificateDescriptionX509DescriptionPolicyIds, 0, len(des))
		for _, d := range des {
			cd := canonicalizeCertificateCertificateDescriptionX509DescriptionPolicyIds(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]CertificateCertificateDescriptionX509DescriptionPolicyIds, 0, len(des))
	for i, d := range des {
		cd := canonicalizeCertificateCertificateDescriptionX509DescriptionPolicyIds(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewCertificateCertificateDescriptionX509DescriptionPolicyIds(c *Client, des, nw *CertificateCertificateDescriptionX509DescriptionPolicyIds) *CertificateCertificateDescriptionX509DescriptionPolicyIds {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for CertificateCertificateDescriptionX509DescriptionPolicyIds while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	return nw
}

func canonicalizeNewCertificateCertificateDescriptionX509DescriptionPolicyIdsSet(c *Client, des, nw []CertificateCertificateDescriptionX509DescriptionPolicyIds) []CertificateCertificateDescriptionX509DescriptionPolicyIds {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []CertificateCertificateDescriptionX509DescriptionPolicyIds
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareCertificateCertificateDescriptionX509DescriptionPolicyIdsNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewCertificateCertificateDescriptionX509DescriptionPolicyIds(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewCertificateCertificateDescriptionX509DescriptionPolicyIdsSlice(c *Client, des, nw []CertificateCertificateDescriptionX509DescriptionPolicyIds) []CertificateCertificateDescriptionX509DescriptionPolicyIds {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []CertificateCertificateDescriptionX509DescriptionPolicyIds
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewCertificateCertificateDescriptionX509DescriptionPolicyIds(c, &d, &n))
	}

	return items
}

func canonicalizeCertificateCertificateDescriptionX509DescriptionAdditionalExtensions(des, initial *CertificateCertificateDescriptionX509DescriptionAdditionalExtensions, opts ...dcl.ApplyOption) *CertificateCertificateDescriptionX509DescriptionAdditionalExtensions {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &CertificateCertificateDescriptionX509DescriptionAdditionalExtensions{}

	cDes.ObjectId = canonicalizeCertificateCertificateDescriptionX509DescriptionAdditionalExtensionsObjectId(des.ObjectId, initial.ObjectId, opts...)
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

func canonicalizeCertificateCertificateDescriptionX509DescriptionAdditionalExtensionsSlice(des, initial []CertificateCertificateDescriptionX509DescriptionAdditionalExtensions, opts ...dcl.ApplyOption) []CertificateCertificateDescriptionX509DescriptionAdditionalExtensions {
	if des == nil {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]CertificateCertificateDescriptionX509DescriptionAdditionalExtensions, 0, len(des))
		for _, d := range des {
			cd := canonicalizeCertificateCertificateDescriptionX509DescriptionAdditionalExtensions(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]CertificateCertificateDescriptionX509DescriptionAdditionalExtensions, 0, len(des))
	for i, d := range des {
		cd := canonicalizeCertificateCertificateDescriptionX509DescriptionAdditionalExtensions(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewCertificateCertificateDescriptionX509DescriptionAdditionalExtensions(c *Client, des, nw *CertificateCertificateDescriptionX509DescriptionAdditionalExtensions) *CertificateCertificateDescriptionX509DescriptionAdditionalExtensions {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for CertificateCertificateDescriptionX509DescriptionAdditionalExtensions while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	nw.ObjectId = canonicalizeNewCertificateCertificateDescriptionX509DescriptionAdditionalExtensionsObjectId(c, des.ObjectId, nw.ObjectId)
	if dcl.BoolCanonicalize(des.Critical, nw.Critical) {
		nw.Critical = des.Critical
	}
	if dcl.StringCanonicalize(des.Value, nw.Value) {
		nw.Value = des.Value
	}

	return nw
}

func canonicalizeNewCertificateCertificateDescriptionX509DescriptionAdditionalExtensionsSet(c *Client, des, nw []CertificateCertificateDescriptionX509DescriptionAdditionalExtensions) []CertificateCertificateDescriptionX509DescriptionAdditionalExtensions {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []CertificateCertificateDescriptionX509DescriptionAdditionalExtensions
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareCertificateCertificateDescriptionX509DescriptionAdditionalExtensionsNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewCertificateCertificateDescriptionX509DescriptionAdditionalExtensions(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewCertificateCertificateDescriptionX509DescriptionAdditionalExtensionsSlice(c *Client, des, nw []CertificateCertificateDescriptionX509DescriptionAdditionalExtensions) []CertificateCertificateDescriptionX509DescriptionAdditionalExtensions {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []CertificateCertificateDescriptionX509DescriptionAdditionalExtensions
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewCertificateCertificateDescriptionX509DescriptionAdditionalExtensions(c, &d, &n))
	}

	return items
}

func canonicalizeCertificateCertificateDescriptionX509DescriptionAdditionalExtensionsObjectId(des, initial *CertificateCertificateDescriptionX509DescriptionAdditionalExtensionsObjectId, opts ...dcl.ApplyOption) *CertificateCertificateDescriptionX509DescriptionAdditionalExtensionsObjectId {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &CertificateCertificateDescriptionX509DescriptionAdditionalExtensionsObjectId{}

	if dcl.IsZeroValue(des.ObjectIdPath) || (dcl.IsEmptyValueIndirect(des.ObjectIdPath) && dcl.IsEmptyValueIndirect(initial.ObjectIdPath)) {
		// Desired and initial values are equivalent, so set canonical desired value to initial value.
		cDes.ObjectIdPath = initial.ObjectIdPath
	} else {
		cDes.ObjectIdPath = des.ObjectIdPath
	}

	return cDes
}

func canonicalizeCertificateCertificateDescriptionX509DescriptionAdditionalExtensionsObjectIdSlice(des, initial []CertificateCertificateDescriptionX509DescriptionAdditionalExtensionsObjectId, opts ...dcl.ApplyOption) []CertificateCertificateDescriptionX509DescriptionAdditionalExtensionsObjectId {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]CertificateCertificateDescriptionX509DescriptionAdditionalExtensionsObjectId, 0, len(des))
		for _, d := range des {
			cd := canonicalizeCertificateCertificateDescriptionX509DescriptionAdditionalExtensionsObjectId(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]CertificateCertificateDescriptionX509DescriptionAdditionalExtensionsObjectId, 0, len(des))
	for i, d := range des {
		cd := canonicalizeCertificateCertificateDescriptionX509DescriptionAdditionalExtensionsObjectId(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewCertificateCertificateDescriptionX509DescriptionAdditionalExtensionsObjectId(c *Client, des, nw *CertificateCertificateDescriptionX509DescriptionAdditionalExtensionsObjectId) *CertificateCertificateDescriptionX509DescriptionAdditionalExtensionsObjectId {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for CertificateCertificateDescriptionX509DescriptionAdditionalExtensionsObjectId while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	return nw
}

func canonicalizeNewCertificateCertificateDescriptionX509DescriptionAdditionalExtensionsObjectIdSet(c *Client, des, nw []CertificateCertificateDescriptionX509DescriptionAdditionalExtensionsObjectId) []CertificateCertificateDescriptionX509DescriptionAdditionalExtensionsObjectId {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []CertificateCertificateDescriptionX509DescriptionAdditionalExtensionsObjectId
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareCertificateCertificateDescriptionX509DescriptionAdditionalExtensionsObjectIdNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewCertificateCertificateDescriptionX509DescriptionAdditionalExtensionsObjectId(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewCertificateCertificateDescriptionX509DescriptionAdditionalExtensionsObjectIdSlice(c *Client, des, nw []CertificateCertificateDescriptionX509DescriptionAdditionalExtensionsObjectId) []CertificateCertificateDescriptionX509DescriptionAdditionalExtensionsObjectId {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []CertificateCertificateDescriptionX509DescriptionAdditionalExtensionsObjectId
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewCertificateCertificateDescriptionX509DescriptionAdditionalExtensionsObjectId(c, &d, &n))
	}

	return items
}

func canonicalizeCertificateCertificateDescriptionPublicKey(des, initial *CertificateCertificateDescriptionPublicKey, opts ...dcl.ApplyOption) *CertificateCertificateDescriptionPublicKey {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &CertificateCertificateDescriptionPublicKey{}

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

func canonicalizeCertificateCertificateDescriptionPublicKeySlice(des, initial []CertificateCertificateDescriptionPublicKey, opts ...dcl.ApplyOption) []CertificateCertificateDescriptionPublicKey {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]CertificateCertificateDescriptionPublicKey, 0, len(des))
		for _, d := range des {
			cd := canonicalizeCertificateCertificateDescriptionPublicKey(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]CertificateCertificateDescriptionPublicKey, 0, len(des))
	for i, d := range des {
		cd := canonicalizeCertificateCertificateDescriptionPublicKey(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewCertificateCertificateDescriptionPublicKey(c *Client, des, nw *CertificateCertificateDescriptionPublicKey) *CertificateCertificateDescriptionPublicKey {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for CertificateCertificateDescriptionPublicKey while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringCanonicalize(des.Key, nw.Key) {
		nw.Key = des.Key
	}

	return nw
}

func canonicalizeNewCertificateCertificateDescriptionPublicKeySet(c *Client, des, nw []CertificateCertificateDescriptionPublicKey) []CertificateCertificateDescriptionPublicKey {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []CertificateCertificateDescriptionPublicKey
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareCertificateCertificateDescriptionPublicKeyNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewCertificateCertificateDescriptionPublicKey(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewCertificateCertificateDescriptionPublicKeySlice(c *Client, des, nw []CertificateCertificateDescriptionPublicKey) []CertificateCertificateDescriptionPublicKey {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []CertificateCertificateDescriptionPublicKey
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewCertificateCertificateDescriptionPublicKey(c, &d, &n))
	}

	return items
}

func canonicalizeCertificateCertificateDescriptionSubjectKeyId(des, initial *CertificateCertificateDescriptionSubjectKeyId, opts ...dcl.ApplyOption) *CertificateCertificateDescriptionSubjectKeyId {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &CertificateCertificateDescriptionSubjectKeyId{}

	if dcl.StringCanonicalize(des.KeyId, initial.KeyId) || dcl.IsZeroValue(des.KeyId) {
		cDes.KeyId = initial.KeyId
	} else {
		cDes.KeyId = des.KeyId
	}

	return cDes
}

func canonicalizeCertificateCertificateDescriptionSubjectKeyIdSlice(des, initial []CertificateCertificateDescriptionSubjectKeyId, opts ...dcl.ApplyOption) []CertificateCertificateDescriptionSubjectKeyId {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]CertificateCertificateDescriptionSubjectKeyId, 0, len(des))
		for _, d := range des {
			cd := canonicalizeCertificateCertificateDescriptionSubjectKeyId(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]CertificateCertificateDescriptionSubjectKeyId, 0, len(des))
	for i, d := range des {
		cd := canonicalizeCertificateCertificateDescriptionSubjectKeyId(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewCertificateCertificateDescriptionSubjectKeyId(c *Client, des, nw *CertificateCertificateDescriptionSubjectKeyId) *CertificateCertificateDescriptionSubjectKeyId {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for CertificateCertificateDescriptionSubjectKeyId while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringCanonicalize(des.KeyId, nw.KeyId) {
		nw.KeyId = des.KeyId
	}

	return nw
}

func canonicalizeNewCertificateCertificateDescriptionSubjectKeyIdSet(c *Client, des, nw []CertificateCertificateDescriptionSubjectKeyId) []CertificateCertificateDescriptionSubjectKeyId {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []CertificateCertificateDescriptionSubjectKeyId
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareCertificateCertificateDescriptionSubjectKeyIdNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewCertificateCertificateDescriptionSubjectKeyId(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewCertificateCertificateDescriptionSubjectKeyIdSlice(c *Client, des, nw []CertificateCertificateDescriptionSubjectKeyId) []CertificateCertificateDescriptionSubjectKeyId {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []CertificateCertificateDescriptionSubjectKeyId
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewCertificateCertificateDescriptionSubjectKeyId(c, &d, &n))
	}

	return items
}

func canonicalizeCertificateCertificateDescriptionAuthorityKeyId(des, initial *CertificateCertificateDescriptionAuthorityKeyId, opts ...dcl.ApplyOption) *CertificateCertificateDescriptionAuthorityKeyId {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &CertificateCertificateDescriptionAuthorityKeyId{}

	if dcl.StringCanonicalize(des.KeyId, initial.KeyId) || dcl.IsZeroValue(des.KeyId) {
		cDes.KeyId = initial.KeyId
	} else {
		cDes.KeyId = des.KeyId
	}

	return cDes
}

func canonicalizeCertificateCertificateDescriptionAuthorityKeyIdSlice(des, initial []CertificateCertificateDescriptionAuthorityKeyId, opts ...dcl.ApplyOption) []CertificateCertificateDescriptionAuthorityKeyId {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]CertificateCertificateDescriptionAuthorityKeyId, 0, len(des))
		for _, d := range des {
			cd := canonicalizeCertificateCertificateDescriptionAuthorityKeyId(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]CertificateCertificateDescriptionAuthorityKeyId, 0, len(des))
	for i, d := range des {
		cd := canonicalizeCertificateCertificateDescriptionAuthorityKeyId(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewCertificateCertificateDescriptionAuthorityKeyId(c *Client, des, nw *CertificateCertificateDescriptionAuthorityKeyId) *CertificateCertificateDescriptionAuthorityKeyId {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for CertificateCertificateDescriptionAuthorityKeyId while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringCanonicalize(des.KeyId, nw.KeyId) {
		nw.KeyId = des.KeyId
	}

	return nw
}

func canonicalizeNewCertificateCertificateDescriptionAuthorityKeyIdSet(c *Client, des, nw []CertificateCertificateDescriptionAuthorityKeyId) []CertificateCertificateDescriptionAuthorityKeyId {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []CertificateCertificateDescriptionAuthorityKeyId
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareCertificateCertificateDescriptionAuthorityKeyIdNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewCertificateCertificateDescriptionAuthorityKeyId(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewCertificateCertificateDescriptionAuthorityKeyIdSlice(c *Client, des, nw []CertificateCertificateDescriptionAuthorityKeyId) []CertificateCertificateDescriptionAuthorityKeyId {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []CertificateCertificateDescriptionAuthorityKeyId
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewCertificateCertificateDescriptionAuthorityKeyId(c, &d, &n))
	}

	return items
}

func canonicalizeCertificateCertificateDescriptionCertFingerprint(des, initial *CertificateCertificateDescriptionCertFingerprint, opts ...dcl.ApplyOption) *CertificateCertificateDescriptionCertFingerprint {
	if des == nil {
		return initial
	}
	if des.empty {
		return des
	}

	if initial == nil {
		return des
	}

	cDes := &CertificateCertificateDescriptionCertFingerprint{}

	if dcl.StringCanonicalize(des.Sha256Hash, initial.Sha256Hash) || dcl.IsZeroValue(des.Sha256Hash) {
		cDes.Sha256Hash = initial.Sha256Hash
	} else {
		cDes.Sha256Hash = des.Sha256Hash
	}

	return cDes
}

func canonicalizeCertificateCertificateDescriptionCertFingerprintSlice(des, initial []CertificateCertificateDescriptionCertFingerprint, opts ...dcl.ApplyOption) []CertificateCertificateDescriptionCertFingerprint {
	if dcl.IsEmptyValueIndirect(des) {
		return initial
	}

	if len(des) != len(initial) {

		items := make([]CertificateCertificateDescriptionCertFingerprint, 0, len(des))
		for _, d := range des {
			cd := canonicalizeCertificateCertificateDescriptionCertFingerprint(&d, nil, opts...)
			if cd != nil {
				items = append(items, *cd)
			}
		}
		return items
	}

	items := make([]CertificateCertificateDescriptionCertFingerprint, 0, len(des))
	for i, d := range des {
		cd := canonicalizeCertificateCertificateDescriptionCertFingerprint(&d, &initial[i], opts...)
		if cd != nil {
			items = append(items, *cd)
		}
	}
	return items

}

func canonicalizeNewCertificateCertificateDescriptionCertFingerprint(c *Client, des, nw *CertificateCertificateDescriptionCertFingerprint) *CertificateCertificateDescriptionCertFingerprint {

	if des == nil {
		return nw
	}

	if nw == nil {
		if dcl.IsEmptyValueIndirect(des) {
			c.Config.Logger.Info("Found explicitly empty value for CertificateCertificateDescriptionCertFingerprint while comparing non-nil desired to nil actual.  Returning desired object.")
			return des
		}
		return nil
	}

	if dcl.StringCanonicalize(des.Sha256Hash, nw.Sha256Hash) {
		nw.Sha256Hash = des.Sha256Hash
	}

	return nw
}

func canonicalizeNewCertificateCertificateDescriptionCertFingerprintSet(c *Client, des, nw []CertificateCertificateDescriptionCertFingerprint) []CertificateCertificateDescriptionCertFingerprint {
	if des == nil {
		return nw
	}

	// Find the elements in des that are also in nw and canonicalize them. Remove matched elements from nw.
	var items []CertificateCertificateDescriptionCertFingerprint
	for _, d := range des {
		matchedIndex := -1
		for i, n := range nw {
			if diffs, _ := compareCertificateCertificateDescriptionCertFingerprintNewStyle(&d, &n, dcl.FieldName{}); len(diffs) == 0 {
				matchedIndex = i
				break
			}
		}
		if matchedIndex != -1 {
			items = append(items, *canonicalizeNewCertificateCertificateDescriptionCertFingerprint(c, &d, &nw[matchedIndex]))
			nw = append(nw[:matchedIndex], nw[matchedIndex+1:]...)
		}
	}
	// Also include elements in nw that are not matched in des.
	items = append(items, nw...)

	return items
}

func canonicalizeNewCertificateCertificateDescriptionCertFingerprintSlice(c *Client, des, nw []CertificateCertificateDescriptionCertFingerprint) []CertificateCertificateDescriptionCertFingerprint {
	if des == nil {
		return nw
	}

	// Lengths are unequal. A diff will occur later, so we shouldn't canonicalize.
	// Return the original array.
	if len(des) != len(nw) {
		return nw
	}

	var items []CertificateCertificateDescriptionCertFingerprint
	for i, d := range des {
		n := nw[i]
		items = append(items, *canonicalizeNewCertificateCertificateDescriptionCertFingerprint(c, &d, &n))
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
func diffCertificate(c *Client, desired, actual *Certificate, opts ...dcl.ApplyOption) ([]*dcl.FieldDiff, error) {
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

	if ds, err := dcl.Diff(desired.PemCsr, actual.PemCsr, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("PemCsr")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Config, actual.Config, dcl.DiffInfo{ObjectFunction: compareCertificateConfigNewStyle, EmptyObject: EmptyCertificateConfig, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Config")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.IssuerCertificateAuthority, actual.IssuerCertificateAuthority, dcl.DiffInfo{OutputOnly: true, Type: "ReferenceType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("IssuerCertificateAuthority")); len(ds) != 0 || err != nil {
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

	if ds, err := dcl.Diff(desired.CertificateTemplate, actual.CertificateTemplate, dcl.DiffInfo{Type: "ReferenceType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("CertificateTemplate")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.SubjectMode, actual.SubjectMode, dcl.DiffInfo{Type: "EnumType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("SubjectMode")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.RevocationDetails, actual.RevocationDetails, dcl.DiffInfo{OutputOnly: true, ObjectFunction: compareCertificateRevocationDetailsNewStyle, EmptyObject: EmptyCertificateRevocationDetails, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("RevocationDetails")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.PemCertificate, actual.PemCertificate, dcl.DiffInfo{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("PemCertificate")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.CertificateDescription, actual.CertificateDescription, dcl.DiffInfo{OutputOnly: true, ObjectFunction: compareCertificateCertificateDescriptionNewStyle, EmptyObject: EmptyCertificateCertificateDescription, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("CertificateDescription")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		newDiffs = append(newDiffs, ds...)
	}

	if ds, err := dcl.Diff(desired.PemCertificateChain, actual.PemCertificateChain, dcl.DiffInfo{OutputOnly: true, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("PemCertificateChain")); len(ds) != 0 || err != nil {
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

	if ds, err := dcl.Diff(desired.Labels, actual.Labels, dcl.DiffInfo{OperationSelector: dcl.TriggersOperation("updateCertificateUpdateCertificateOperation")}, fn.AddNest("Labels")); len(ds) != 0 || err != nil {
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

	if ds, err := dcl.Diff(desired.CertificateAuthority, actual.CertificateAuthority, dcl.DiffInfo{Type: "ReferenceType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("CertificateAuthority")); len(ds) != 0 || err != nil {
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
func compareCertificateConfigNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*CertificateConfig)
	if !ok {
		desiredNotPointer, ok := d.(CertificateConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CertificateConfig or *CertificateConfig", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*CertificateConfig)
	if !ok {
		actualNotPointer, ok := a.(CertificateConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CertificateConfig", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.SubjectConfig, actual.SubjectConfig, dcl.DiffInfo{ObjectFunction: compareCertificateConfigSubjectConfigNewStyle, EmptyObject: EmptyCertificateConfigSubjectConfig, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("SubjectConfig")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.X509Config, actual.X509Config, dcl.DiffInfo{ObjectFunction: compareCertificateConfigX509ConfigNewStyle, EmptyObject: EmptyCertificateConfigX509Config, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("X509Config")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.PublicKey, actual.PublicKey, dcl.DiffInfo{ObjectFunction: compareCertificateConfigPublicKeyNewStyle, EmptyObject: EmptyCertificateConfigPublicKey, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("PublicKey")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareCertificateConfigSubjectConfigNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*CertificateConfigSubjectConfig)
	if !ok {
		desiredNotPointer, ok := d.(CertificateConfigSubjectConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CertificateConfigSubjectConfig or *CertificateConfigSubjectConfig", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*CertificateConfigSubjectConfig)
	if !ok {
		actualNotPointer, ok := a.(CertificateConfigSubjectConfig)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CertificateConfigSubjectConfig", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.Subject, actual.Subject, dcl.DiffInfo{ObjectFunction: compareCertificateConfigSubjectConfigSubjectNewStyle, EmptyObject: EmptyCertificateConfigSubjectConfigSubject, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Subject")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.SubjectAltName, actual.SubjectAltName, dcl.DiffInfo{ObjectFunction: compareCertificateConfigSubjectConfigSubjectAltNameNewStyle, EmptyObject: EmptyCertificateConfigSubjectConfigSubjectAltName, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("SubjectAltName")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareCertificateConfigSubjectConfigSubjectNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*CertificateConfigSubjectConfigSubject)
	if !ok {
		desiredNotPointer, ok := d.(CertificateConfigSubjectConfigSubject)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CertificateConfigSubjectConfigSubject or *CertificateConfigSubjectConfigSubject", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*CertificateConfigSubjectConfigSubject)
	if !ok {
		actualNotPointer, ok := a.(CertificateConfigSubjectConfigSubject)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CertificateConfigSubjectConfigSubject", a)
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

func compareCertificateConfigSubjectConfigSubjectAltNameNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*CertificateConfigSubjectConfigSubjectAltName)
	if !ok {
		desiredNotPointer, ok := d.(CertificateConfigSubjectConfigSubjectAltName)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CertificateConfigSubjectConfigSubjectAltName or *CertificateConfigSubjectConfigSubjectAltName", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*CertificateConfigSubjectConfigSubjectAltName)
	if !ok {
		actualNotPointer, ok := a.(CertificateConfigSubjectConfigSubjectAltName)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CertificateConfigSubjectConfigSubjectAltName", a)
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
	return diffs, nil
}

func compareCertificateConfigX509ConfigNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*CertificateConfigX509Config)
	if !ok {
		desiredNotPointer, ok := d.(CertificateConfigX509Config)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CertificateConfigX509Config or *CertificateConfigX509Config", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*CertificateConfigX509Config)
	if !ok {
		actualNotPointer, ok := a.(CertificateConfigX509Config)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CertificateConfigX509Config", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.KeyUsage, actual.KeyUsage, dcl.DiffInfo{ObjectFunction: compareCertificateConfigX509ConfigKeyUsageNewStyle, EmptyObject: EmptyCertificateConfigX509ConfigKeyUsage, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("KeyUsage")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.CaOptions, actual.CaOptions, dcl.DiffInfo{ObjectFunction: compareCertificateConfigX509ConfigCaOptionsNewStyle, EmptyObject: EmptyCertificateConfigX509ConfigCaOptions, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("CaOptions")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.PolicyIds, actual.PolicyIds, dcl.DiffInfo{ObjectFunction: compareCertificateConfigX509ConfigPolicyIdsNewStyle, EmptyObject: EmptyCertificateConfigX509ConfigPolicyIds, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("PolicyIds")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.AiaOcspServers, actual.AiaOcspServers, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("AiaOcspServers")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.AdditionalExtensions, actual.AdditionalExtensions, dcl.DiffInfo{ObjectFunction: compareCertificateConfigX509ConfigAdditionalExtensionsNewStyle, EmptyObject: EmptyCertificateConfigX509ConfigAdditionalExtensions, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("AdditionalExtensions")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareCertificateConfigX509ConfigKeyUsageNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*CertificateConfigX509ConfigKeyUsage)
	if !ok {
		desiredNotPointer, ok := d.(CertificateConfigX509ConfigKeyUsage)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CertificateConfigX509ConfigKeyUsage or *CertificateConfigX509ConfigKeyUsage", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*CertificateConfigX509ConfigKeyUsage)
	if !ok {
		actualNotPointer, ok := a.(CertificateConfigX509ConfigKeyUsage)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CertificateConfigX509ConfigKeyUsage", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.BaseKeyUsage, actual.BaseKeyUsage, dcl.DiffInfo{ObjectFunction: compareCertificateConfigX509ConfigKeyUsageBaseKeyUsageNewStyle, EmptyObject: EmptyCertificateConfigX509ConfigKeyUsageBaseKeyUsage, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("BaseKeyUsage")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.ExtendedKeyUsage, actual.ExtendedKeyUsage, dcl.DiffInfo{ObjectFunction: compareCertificateConfigX509ConfigKeyUsageExtendedKeyUsageNewStyle, EmptyObject: EmptyCertificateConfigX509ConfigKeyUsageExtendedKeyUsage, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("ExtendedKeyUsage")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.UnknownExtendedKeyUsages, actual.UnknownExtendedKeyUsages, dcl.DiffInfo{ObjectFunction: compareCertificateConfigX509ConfigKeyUsageUnknownExtendedKeyUsagesNewStyle, EmptyObject: EmptyCertificateConfigX509ConfigKeyUsageUnknownExtendedKeyUsages, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("UnknownExtendedKeyUsages")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareCertificateConfigX509ConfigKeyUsageBaseKeyUsageNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*CertificateConfigX509ConfigKeyUsageBaseKeyUsage)
	if !ok {
		desiredNotPointer, ok := d.(CertificateConfigX509ConfigKeyUsageBaseKeyUsage)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CertificateConfigX509ConfigKeyUsageBaseKeyUsage or *CertificateConfigX509ConfigKeyUsageBaseKeyUsage", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*CertificateConfigX509ConfigKeyUsageBaseKeyUsage)
	if !ok {
		actualNotPointer, ok := a.(CertificateConfigX509ConfigKeyUsageBaseKeyUsage)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CertificateConfigX509ConfigKeyUsageBaseKeyUsage", a)
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

func compareCertificateConfigX509ConfigKeyUsageExtendedKeyUsageNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*CertificateConfigX509ConfigKeyUsageExtendedKeyUsage)
	if !ok {
		desiredNotPointer, ok := d.(CertificateConfigX509ConfigKeyUsageExtendedKeyUsage)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CertificateConfigX509ConfigKeyUsageExtendedKeyUsage or *CertificateConfigX509ConfigKeyUsageExtendedKeyUsage", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*CertificateConfigX509ConfigKeyUsageExtendedKeyUsage)
	if !ok {
		actualNotPointer, ok := a.(CertificateConfigX509ConfigKeyUsageExtendedKeyUsage)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CertificateConfigX509ConfigKeyUsageExtendedKeyUsage", a)
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

func compareCertificateConfigX509ConfigKeyUsageUnknownExtendedKeyUsagesNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*CertificateConfigX509ConfigKeyUsageUnknownExtendedKeyUsages)
	if !ok {
		desiredNotPointer, ok := d.(CertificateConfigX509ConfigKeyUsageUnknownExtendedKeyUsages)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CertificateConfigX509ConfigKeyUsageUnknownExtendedKeyUsages or *CertificateConfigX509ConfigKeyUsageUnknownExtendedKeyUsages", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*CertificateConfigX509ConfigKeyUsageUnknownExtendedKeyUsages)
	if !ok {
		actualNotPointer, ok := a.(CertificateConfigX509ConfigKeyUsageUnknownExtendedKeyUsages)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CertificateConfigX509ConfigKeyUsageUnknownExtendedKeyUsages", a)
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

func compareCertificateConfigX509ConfigCaOptionsNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*CertificateConfigX509ConfigCaOptions)
	if !ok {
		desiredNotPointer, ok := d.(CertificateConfigX509ConfigCaOptions)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CertificateConfigX509ConfigCaOptions or *CertificateConfigX509ConfigCaOptions", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*CertificateConfigX509ConfigCaOptions)
	if !ok {
		actualNotPointer, ok := a.(CertificateConfigX509ConfigCaOptions)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CertificateConfigX509ConfigCaOptions", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.IsCa, actual.IsCa, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("IsCa")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.NonCa, actual.NonCa, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("NonCa")); len(ds) != 0 || err != nil {
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

func compareCertificateConfigX509ConfigPolicyIdsNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*CertificateConfigX509ConfigPolicyIds)
	if !ok {
		desiredNotPointer, ok := d.(CertificateConfigX509ConfigPolicyIds)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CertificateConfigX509ConfigPolicyIds or *CertificateConfigX509ConfigPolicyIds", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*CertificateConfigX509ConfigPolicyIds)
	if !ok {
		actualNotPointer, ok := a.(CertificateConfigX509ConfigPolicyIds)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CertificateConfigX509ConfigPolicyIds", a)
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

func compareCertificateConfigX509ConfigAdditionalExtensionsNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*CertificateConfigX509ConfigAdditionalExtensions)
	if !ok {
		desiredNotPointer, ok := d.(CertificateConfigX509ConfigAdditionalExtensions)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CertificateConfigX509ConfigAdditionalExtensions or *CertificateConfigX509ConfigAdditionalExtensions", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*CertificateConfigX509ConfigAdditionalExtensions)
	if !ok {
		actualNotPointer, ok := a.(CertificateConfigX509ConfigAdditionalExtensions)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CertificateConfigX509ConfigAdditionalExtensions", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.ObjectId, actual.ObjectId, dcl.DiffInfo{ObjectFunction: compareCertificateConfigX509ConfigAdditionalExtensionsObjectIdNewStyle, EmptyObject: EmptyCertificateConfigX509ConfigAdditionalExtensionsObjectId, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("ObjectId")); len(ds) != 0 || err != nil {
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

func compareCertificateConfigX509ConfigAdditionalExtensionsObjectIdNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*CertificateConfigX509ConfigAdditionalExtensionsObjectId)
	if !ok {
		desiredNotPointer, ok := d.(CertificateConfigX509ConfigAdditionalExtensionsObjectId)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CertificateConfigX509ConfigAdditionalExtensionsObjectId or *CertificateConfigX509ConfigAdditionalExtensionsObjectId", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*CertificateConfigX509ConfigAdditionalExtensionsObjectId)
	if !ok {
		actualNotPointer, ok := a.(CertificateConfigX509ConfigAdditionalExtensionsObjectId)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CertificateConfigX509ConfigAdditionalExtensionsObjectId", a)
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

func compareCertificateConfigPublicKeyNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*CertificateConfigPublicKey)
	if !ok {
		desiredNotPointer, ok := d.(CertificateConfigPublicKey)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CertificateConfigPublicKey or *CertificateConfigPublicKey", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*CertificateConfigPublicKey)
	if !ok {
		actualNotPointer, ok := a.(CertificateConfigPublicKey)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CertificateConfigPublicKey", a)
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

func compareCertificateRevocationDetailsNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*CertificateRevocationDetails)
	if !ok {
		desiredNotPointer, ok := d.(CertificateRevocationDetails)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CertificateRevocationDetails or *CertificateRevocationDetails", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*CertificateRevocationDetails)
	if !ok {
		actualNotPointer, ok := a.(CertificateRevocationDetails)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CertificateRevocationDetails", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.RevocationState, actual.RevocationState, dcl.DiffInfo{Type: "EnumType", OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("RevocationState")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.RevocationTime, actual.RevocationTime, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("RevocationTime")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareCertificateCertificateDescriptionNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*CertificateCertificateDescription)
	if !ok {
		desiredNotPointer, ok := d.(CertificateCertificateDescription)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CertificateCertificateDescription or *CertificateCertificateDescription", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*CertificateCertificateDescription)
	if !ok {
		actualNotPointer, ok := a.(CertificateCertificateDescription)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CertificateCertificateDescription", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.SubjectDescription, actual.SubjectDescription, dcl.DiffInfo{ObjectFunction: compareCertificateCertificateDescriptionSubjectDescriptionNewStyle, EmptyObject: EmptyCertificateCertificateDescriptionSubjectDescription, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("SubjectDescription")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.X509Description, actual.X509Description, dcl.DiffInfo{ObjectFunction: compareCertificateCertificateDescriptionX509DescriptionNewStyle, EmptyObject: EmptyCertificateCertificateDescriptionX509Description, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("X509Description")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.PublicKey, actual.PublicKey, dcl.DiffInfo{ObjectFunction: compareCertificateCertificateDescriptionPublicKeyNewStyle, EmptyObject: EmptyCertificateCertificateDescriptionPublicKey, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("PublicKey")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.SubjectKeyId, actual.SubjectKeyId, dcl.DiffInfo{ObjectFunction: compareCertificateCertificateDescriptionSubjectKeyIdNewStyle, EmptyObject: EmptyCertificateCertificateDescriptionSubjectKeyId, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("SubjectKeyId")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.AuthorityKeyId, actual.AuthorityKeyId, dcl.DiffInfo{ObjectFunction: compareCertificateCertificateDescriptionAuthorityKeyIdNewStyle, EmptyObject: EmptyCertificateCertificateDescriptionAuthorityKeyId, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("AuthorityKeyId")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.CrlDistributionPoints, actual.CrlDistributionPoints, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("CrlDistributionPoints")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.AiaIssuingCertificateUrls, actual.AiaIssuingCertificateUrls, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("AiaIssuingCertificateUrls")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.CertFingerprint, actual.CertFingerprint, dcl.DiffInfo{ObjectFunction: compareCertificateCertificateDescriptionCertFingerprintNewStyle, EmptyObject: EmptyCertificateCertificateDescriptionCertFingerprint, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("CertFingerprint")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareCertificateCertificateDescriptionSubjectDescriptionNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*CertificateCertificateDescriptionSubjectDescription)
	if !ok {
		desiredNotPointer, ok := d.(CertificateCertificateDescriptionSubjectDescription)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CertificateCertificateDescriptionSubjectDescription or *CertificateCertificateDescriptionSubjectDescription", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*CertificateCertificateDescriptionSubjectDescription)
	if !ok {
		actualNotPointer, ok := a.(CertificateCertificateDescriptionSubjectDescription)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CertificateCertificateDescriptionSubjectDescription", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.Subject, actual.Subject, dcl.DiffInfo{ObjectFunction: compareCertificateCertificateDescriptionSubjectDescriptionSubjectNewStyle, EmptyObject: EmptyCertificateCertificateDescriptionSubjectDescriptionSubject, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Subject")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.SubjectAltName, actual.SubjectAltName, dcl.DiffInfo{ObjectFunction: compareCertificateCertificateDescriptionSubjectDescriptionSubjectAltNameNewStyle, EmptyObject: EmptyCertificateCertificateDescriptionSubjectDescriptionSubjectAltName, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("SubjectAltName")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.HexSerialNumber, actual.HexSerialNumber, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("HexSerialNumber")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.Lifetime, actual.Lifetime, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Lifetime")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.NotBeforeTime, actual.NotBeforeTime, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("NotBeforeTime")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.NotAfterTime, actual.NotAfterTime, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("NotAfterTime")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareCertificateCertificateDescriptionSubjectDescriptionSubjectNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*CertificateCertificateDescriptionSubjectDescriptionSubject)
	if !ok {
		desiredNotPointer, ok := d.(CertificateCertificateDescriptionSubjectDescriptionSubject)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CertificateCertificateDescriptionSubjectDescriptionSubject or *CertificateCertificateDescriptionSubjectDescriptionSubject", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*CertificateCertificateDescriptionSubjectDescriptionSubject)
	if !ok {
		actualNotPointer, ok := a.(CertificateCertificateDescriptionSubjectDescriptionSubject)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CertificateCertificateDescriptionSubjectDescriptionSubject", a)
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

func compareCertificateCertificateDescriptionSubjectDescriptionSubjectAltNameNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*CertificateCertificateDescriptionSubjectDescriptionSubjectAltName)
	if !ok {
		desiredNotPointer, ok := d.(CertificateCertificateDescriptionSubjectDescriptionSubjectAltName)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CertificateCertificateDescriptionSubjectDescriptionSubjectAltName or *CertificateCertificateDescriptionSubjectDescriptionSubjectAltName", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*CertificateCertificateDescriptionSubjectDescriptionSubjectAltName)
	if !ok {
		actualNotPointer, ok := a.(CertificateCertificateDescriptionSubjectDescriptionSubjectAltName)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CertificateCertificateDescriptionSubjectDescriptionSubjectAltName", a)
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

	if ds, err := dcl.Diff(desired.CustomSans, actual.CustomSans, dcl.DiffInfo{ObjectFunction: compareCertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSansNewStyle, EmptyObject: EmptyCertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSans, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("CustomSans")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareCertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSansNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*CertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSans)
	if !ok {
		desiredNotPointer, ok := d.(CertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSans)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSans or *CertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSans", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*CertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSans)
	if !ok {
		actualNotPointer, ok := a.(CertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSans)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSans", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.ObjectId, actual.ObjectId, dcl.DiffInfo{ObjectFunction: compareCertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSansObjectIdNewStyle, EmptyObject: EmptyCertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSansObjectId, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("ObjectId")); len(ds) != 0 || err != nil {
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

func compareCertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSansObjectIdNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*CertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSansObjectId)
	if !ok {
		desiredNotPointer, ok := d.(CertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSansObjectId)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSansObjectId or *CertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSansObjectId", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*CertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSansObjectId)
	if !ok {
		actualNotPointer, ok := a.(CertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSansObjectId)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSansObjectId", a)
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

func compareCertificateCertificateDescriptionX509DescriptionNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*CertificateCertificateDescriptionX509Description)
	if !ok {
		desiredNotPointer, ok := d.(CertificateCertificateDescriptionX509Description)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CertificateCertificateDescriptionX509Description or *CertificateCertificateDescriptionX509Description", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*CertificateCertificateDescriptionX509Description)
	if !ok {
		actualNotPointer, ok := a.(CertificateCertificateDescriptionX509Description)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CertificateCertificateDescriptionX509Description", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.KeyUsage, actual.KeyUsage, dcl.DiffInfo{ObjectFunction: compareCertificateCertificateDescriptionX509DescriptionKeyUsageNewStyle, EmptyObject: EmptyCertificateCertificateDescriptionX509DescriptionKeyUsage, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("KeyUsage")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.CaOptions, actual.CaOptions, dcl.DiffInfo{ObjectFunction: compareCertificateCertificateDescriptionX509DescriptionCaOptionsNewStyle, EmptyObject: EmptyCertificateCertificateDescriptionX509DescriptionCaOptions, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("CaOptions")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.PolicyIds, actual.PolicyIds, dcl.DiffInfo{ObjectFunction: compareCertificateCertificateDescriptionX509DescriptionPolicyIdsNewStyle, EmptyObject: EmptyCertificateCertificateDescriptionX509DescriptionPolicyIds, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("PolicyIds")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.AiaOcspServers, actual.AiaOcspServers, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("AiaOcspServers")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.AdditionalExtensions, actual.AdditionalExtensions, dcl.DiffInfo{ObjectFunction: compareCertificateCertificateDescriptionX509DescriptionAdditionalExtensionsNewStyle, EmptyObject: EmptyCertificateCertificateDescriptionX509DescriptionAdditionalExtensions, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("AdditionalExtensions")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareCertificateCertificateDescriptionX509DescriptionKeyUsageNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*CertificateCertificateDescriptionX509DescriptionKeyUsage)
	if !ok {
		desiredNotPointer, ok := d.(CertificateCertificateDescriptionX509DescriptionKeyUsage)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CertificateCertificateDescriptionX509DescriptionKeyUsage or *CertificateCertificateDescriptionX509DescriptionKeyUsage", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*CertificateCertificateDescriptionX509DescriptionKeyUsage)
	if !ok {
		actualNotPointer, ok := a.(CertificateCertificateDescriptionX509DescriptionKeyUsage)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CertificateCertificateDescriptionX509DescriptionKeyUsage", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.BaseKeyUsage, actual.BaseKeyUsage, dcl.DiffInfo{ObjectFunction: compareCertificateCertificateDescriptionX509DescriptionKeyUsageBaseKeyUsageNewStyle, EmptyObject: EmptyCertificateCertificateDescriptionX509DescriptionKeyUsageBaseKeyUsage, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("BaseKeyUsage")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.ExtendedKeyUsage, actual.ExtendedKeyUsage, dcl.DiffInfo{ObjectFunction: compareCertificateCertificateDescriptionX509DescriptionKeyUsageExtendedKeyUsageNewStyle, EmptyObject: EmptyCertificateCertificateDescriptionX509DescriptionKeyUsageExtendedKeyUsage, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("ExtendedKeyUsage")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}

	if ds, err := dcl.Diff(desired.UnknownExtendedKeyUsages, actual.UnknownExtendedKeyUsages, dcl.DiffInfo{ObjectFunction: compareCertificateCertificateDescriptionX509DescriptionKeyUsageUnknownExtendedKeyUsagesNewStyle, EmptyObject: EmptyCertificateCertificateDescriptionX509DescriptionKeyUsageUnknownExtendedKeyUsages, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("UnknownExtendedKeyUsages")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareCertificateCertificateDescriptionX509DescriptionKeyUsageBaseKeyUsageNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*CertificateCertificateDescriptionX509DescriptionKeyUsageBaseKeyUsage)
	if !ok {
		desiredNotPointer, ok := d.(CertificateCertificateDescriptionX509DescriptionKeyUsageBaseKeyUsage)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CertificateCertificateDescriptionX509DescriptionKeyUsageBaseKeyUsage or *CertificateCertificateDescriptionX509DescriptionKeyUsageBaseKeyUsage", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*CertificateCertificateDescriptionX509DescriptionKeyUsageBaseKeyUsage)
	if !ok {
		actualNotPointer, ok := a.(CertificateCertificateDescriptionX509DescriptionKeyUsageBaseKeyUsage)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CertificateCertificateDescriptionX509DescriptionKeyUsageBaseKeyUsage", a)
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

func compareCertificateCertificateDescriptionX509DescriptionKeyUsageExtendedKeyUsageNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*CertificateCertificateDescriptionX509DescriptionKeyUsageExtendedKeyUsage)
	if !ok {
		desiredNotPointer, ok := d.(CertificateCertificateDescriptionX509DescriptionKeyUsageExtendedKeyUsage)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CertificateCertificateDescriptionX509DescriptionKeyUsageExtendedKeyUsage or *CertificateCertificateDescriptionX509DescriptionKeyUsageExtendedKeyUsage", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*CertificateCertificateDescriptionX509DescriptionKeyUsageExtendedKeyUsage)
	if !ok {
		actualNotPointer, ok := a.(CertificateCertificateDescriptionX509DescriptionKeyUsageExtendedKeyUsage)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CertificateCertificateDescriptionX509DescriptionKeyUsageExtendedKeyUsage", a)
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

func compareCertificateCertificateDescriptionX509DescriptionKeyUsageUnknownExtendedKeyUsagesNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*CertificateCertificateDescriptionX509DescriptionKeyUsageUnknownExtendedKeyUsages)
	if !ok {
		desiredNotPointer, ok := d.(CertificateCertificateDescriptionX509DescriptionKeyUsageUnknownExtendedKeyUsages)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CertificateCertificateDescriptionX509DescriptionKeyUsageUnknownExtendedKeyUsages or *CertificateCertificateDescriptionX509DescriptionKeyUsageUnknownExtendedKeyUsages", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*CertificateCertificateDescriptionX509DescriptionKeyUsageUnknownExtendedKeyUsages)
	if !ok {
		actualNotPointer, ok := a.(CertificateCertificateDescriptionX509DescriptionKeyUsageUnknownExtendedKeyUsages)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CertificateCertificateDescriptionX509DescriptionKeyUsageUnknownExtendedKeyUsages", a)
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

func compareCertificateCertificateDescriptionX509DescriptionCaOptionsNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*CertificateCertificateDescriptionX509DescriptionCaOptions)
	if !ok {
		desiredNotPointer, ok := d.(CertificateCertificateDescriptionX509DescriptionCaOptions)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CertificateCertificateDescriptionX509DescriptionCaOptions or *CertificateCertificateDescriptionX509DescriptionCaOptions", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*CertificateCertificateDescriptionX509DescriptionCaOptions)
	if !ok {
		actualNotPointer, ok := a.(CertificateCertificateDescriptionX509DescriptionCaOptions)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CertificateCertificateDescriptionX509DescriptionCaOptions", a)
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
	return diffs, nil
}

func compareCertificateCertificateDescriptionX509DescriptionPolicyIdsNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*CertificateCertificateDescriptionX509DescriptionPolicyIds)
	if !ok {
		desiredNotPointer, ok := d.(CertificateCertificateDescriptionX509DescriptionPolicyIds)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CertificateCertificateDescriptionX509DescriptionPolicyIds or *CertificateCertificateDescriptionX509DescriptionPolicyIds", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*CertificateCertificateDescriptionX509DescriptionPolicyIds)
	if !ok {
		actualNotPointer, ok := a.(CertificateCertificateDescriptionX509DescriptionPolicyIds)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CertificateCertificateDescriptionX509DescriptionPolicyIds", a)
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

func compareCertificateCertificateDescriptionX509DescriptionAdditionalExtensionsNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*CertificateCertificateDescriptionX509DescriptionAdditionalExtensions)
	if !ok {
		desiredNotPointer, ok := d.(CertificateCertificateDescriptionX509DescriptionAdditionalExtensions)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CertificateCertificateDescriptionX509DescriptionAdditionalExtensions or *CertificateCertificateDescriptionX509DescriptionAdditionalExtensions", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*CertificateCertificateDescriptionX509DescriptionAdditionalExtensions)
	if !ok {
		actualNotPointer, ok := a.(CertificateCertificateDescriptionX509DescriptionAdditionalExtensions)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CertificateCertificateDescriptionX509DescriptionAdditionalExtensions", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.ObjectId, actual.ObjectId, dcl.DiffInfo{ObjectFunction: compareCertificateCertificateDescriptionX509DescriptionAdditionalExtensionsObjectIdNewStyle, EmptyObject: EmptyCertificateCertificateDescriptionX509DescriptionAdditionalExtensionsObjectId, OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("ObjectId")); len(ds) != 0 || err != nil {
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

func compareCertificateCertificateDescriptionX509DescriptionAdditionalExtensionsObjectIdNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*CertificateCertificateDescriptionX509DescriptionAdditionalExtensionsObjectId)
	if !ok {
		desiredNotPointer, ok := d.(CertificateCertificateDescriptionX509DescriptionAdditionalExtensionsObjectId)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CertificateCertificateDescriptionX509DescriptionAdditionalExtensionsObjectId or *CertificateCertificateDescriptionX509DescriptionAdditionalExtensionsObjectId", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*CertificateCertificateDescriptionX509DescriptionAdditionalExtensionsObjectId)
	if !ok {
		actualNotPointer, ok := a.(CertificateCertificateDescriptionX509DescriptionAdditionalExtensionsObjectId)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CertificateCertificateDescriptionX509DescriptionAdditionalExtensionsObjectId", a)
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

func compareCertificateCertificateDescriptionPublicKeyNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*CertificateCertificateDescriptionPublicKey)
	if !ok {
		desiredNotPointer, ok := d.(CertificateCertificateDescriptionPublicKey)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CertificateCertificateDescriptionPublicKey or *CertificateCertificateDescriptionPublicKey", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*CertificateCertificateDescriptionPublicKey)
	if !ok {
		actualNotPointer, ok := a.(CertificateCertificateDescriptionPublicKey)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CertificateCertificateDescriptionPublicKey", a)
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

func compareCertificateCertificateDescriptionSubjectKeyIdNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*CertificateCertificateDescriptionSubjectKeyId)
	if !ok {
		desiredNotPointer, ok := d.(CertificateCertificateDescriptionSubjectKeyId)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CertificateCertificateDescriptionSubjectKeyId or *CertificateCertificateDescriptionSubjectKeyId", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*CertificateCertificateDescriptionSubjectKeyId)
	if !ok {
		actualNotPointer, ok := a.(CertificateCertificateDescriptionSubjectKeyId)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CertificateCertificateDescriptionSubjectKeyId", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.KeyId, actual.KeyId, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("KeyId")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareCertificateCertificateDescriptionAuthorityKeyIdNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*CertificateCertificateDescriptionAuthorityKeyId)
	if !ok {
		desiredNotPointer, ok := d.(CertificateCertificateDescriptionAuthorityKeyId)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CertificateCertificateDescriptionAuthorityKeyId or *CertificateCertificateDescriptionAuthorityKeyId", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*CertificateCertificateDescriptionAuthorityKeyId)
	if !ok {
		actualNotPointer, ok := a.(CertificateCertificateDescriptionAuthorityKeyId)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CertificateCertificateDescriptionAuthorityKeyId", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.KeyId, actual.KeyId, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("KeyId")); len(ds) != 0 || err != nil {
		if err != nil {
			return nil, err
		}
		diffs = append(diffs, ds...)
	}
	return diffs, nil
}

func compareCertificateCertificateDescriptionCertFingerprintNewStyle(d, a interface{}, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	var diffs []*dcl.FieldDiff

	desired, ok := d.(*CertificateCertificateDescriptionCertFingerprint)
	if !ok {
		desiredNotPointer, ok := d.(CertificateCertificateDescriptionCertFingerprint)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CertificateCertificateDescriptionCertFingerprint or *CertificateCertificateDescriptionCertFingerprint", d)
		}
		desired = &desiredNotPointer
	}
	actual, ok := a.(*CertificateCertificateDescriptionCertFingerprint)
	if !ok {
		actualNotPointer, ok := a.(CertificateCertificateDescriptionCertFingerprint)
		if !ok {
			return nil, fmt.Errorf("obj %v is not a CertificateCertificateDescriptionCertFingerprint", a)
		}
		actual = &actualNotPointer
	}

	if ds, err := dcl.Diff(desired.Sha256Hash, actual.Sha256Hash, dcl.DiffInfo{OperationSelector: dcl.RequiresRecreate()}, fn.AddNest("Sha256Hash")); len(ds) != 0 || err != nil {
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
func (r *Certificate) urlNormalized() *Certificate {
	normalized := dcl.Copy(*r).(Certificate)
	normalized.Name = dcl.SelfLinkToName(r.Name)
	normalized.PemCsr = dcl.SelfLinkToName(r.PemCsr)
	normalized.IssuerCertificateAuthority = dcl.SelfLinkToName(r.IssuerCertificateAuthority)
	normalized.Lifetime = dcl.SelfLinkToName(r.Lifetime)
	normalized.CertificateTemplate = dcl.SelfLinkToName(r.CertificateTemplate)
	normalized.PemCertificate = dcl.SelfLinkToName(r.PemCertificate)
	normalized.Project = dcl.SelfLinkToName(r.Project)
	normalized.Location = dcl.SelfLinkToName(r.Location)
	normalized.CaPool = dcl.SelfLinkToName(r.CaPool)
	normalized.CertificateAuthority = dcl.SelfLinkToName(r.CertificateAuthority)
	return &normalized
}

func (r *Certificate) updateURL(userBasePath, updateName string) (string, error) {
	nr := r.urlNormalized()
	if updateName == "UpdateCertificate" {
		fields := map[string]interface{}{
			"project":  dcl.ValueOrEmptyString(nr.Project),
			"location": dcl.ValueOrEmptyString(nr.Location),
			"caPool":   dcl.ValueOrEmptyString(nr.CaPool),
			"name":     dcl.ValueOrEmptyString(nr.Name),
		}
		return dcl.URL("projects/{{project}}/locations/{{location}}/caPools/{{caPool}}/certificates/{{name}}", nr.basePath(), userBasePath, fields), nil

	}

	return "", fmt.Errorf("unknown update name: %s", updateName)
}

// marshal encodes the Certificate resource into JSON for a Create request, and
// performs transformations from the resource schema to the API schema if
// necessary.
func (r *Certificate) marshal(c *Client) ([]byte, error) {
	m, err := expandCertificate(c, r)
	if err != nil {
		return nil, fmt.Errorf("error marshalling Certificate: %w", err)
	}

	return json.Marshal(m)
}

// unmarshalCertificate decodes JSON responses into the Certificate resource schema.
func unmarshalCertificate(b []byte, c *Client, res *Certificate) (*Certificate, error) {
	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, err
	}
	return unmarshalMapCertificate(m, c, res)
}

func unmarshalMapCertificate(m map[string]interface{}, c *Client, res *Certificate) (*Certificate, error) {

	flattened := flattenCertificate(c, m, res)
	if flattened == nil {
		return nil, fmt.Errorf("attempted to flatten empty json object")
	}
	return flattened, nil
}

// expandCertificate expands Certificate into a JSON request object.
func expandCertificate(c *Client, f *Certificate) (map[string]interface{}, error) {
	m := make(map[string]interface{})
	res := f
	_ = res
	if v, err := dcl.DeriveField("projects/%s/locations/%s/caPools/%s/certificates/%s", f.Name, dcl.SelfLinkToName(f.Project), dcl.SelfLinkToName(f.Location), dcl.SelfLinkToName(f.CaPool), dcl.SelfLinkToName(f.Name)); err != nil {
		return nil, fmt.Errorf("error expanding Name into name: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["name"] = v
	}
	if v := f.PemCsr; dcl.ValueShouldBeSent(v) {
		m["pemCsr"] = v
	}
	if v, err := expandCertificateConfig(c, f.Config, res); err != nil {
		return nil, fmt.Errorf("error expanding Config into config: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["config"] = v
	}
	if v := f.Lifetime; dcl.ValueShouldBeSent(v) {
		m["lifetime"] = v
	}
	if v, err := dcl.DeriveField("projects/%s/locations/%s/certificateTemplates/%s", f.CertificateTemplate, dcl.SelfLinkToName(f.Project), dcl.SelfLinkToName(f.Location), dcl.SelfLinkToName(f.CertificateTemplate)); err != nil {
		return nil, fmt.Errorf("error expanding CertificateTemplate into certificateTemplate: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["certificateTemplate"] = v
	}
	if v := f.SubjectMode; dcl.ValueShouldBeSent(v) {
		m["subjectMode"] = v
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
	if v, err := dcl.EmptyValue(); err != nil {
		return nil, fmt.Errorf("error expanding CertificateAuthority into certificateAuthority: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["certificateAuthority"] = v
	}

	return m, nil
}

// flattenCertificate flattens Certificate from a JSON request object into the
// Certificate type.
func flattenCertificate(c *Client, i interface{}, res *Certificate) *Certificate {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}
	if len(m) == 0 {
		return nil
	}

	resultRes := &Certificate{}
	resultRes.Name = dcl.FlattenString(m["name"])
	resultRes.PemCsr = dcl.FlattenString(m["pemCsr"])
	resultRes.Config = flattenCertificateConfig(c, m["config"], res)
	resultRes.IssuerCertificateAuthority = dcl.FlattenString(m["issuerCertificateAuthority"])
	resultRes.Lifetime = dcl.FlattenString(m["lifetime"])
	resultRes.CertificateTemplate = dcl.FlattenString(m["certificateTemplate"])
	resultRes.SubjectMode = flattenCertificateSubjectModeEnum(m["subjectMode"])
	resultRes.RevocationDetails = flattenCertificateRevocationDetails(c, m["revocationDetails"], res)
	resultRes.PemCertificate = dcl.FlattenString(m["pemCertificate"])
	resultRes.CertificateDescription = flattenCertificateCertificateDescription(c, m["certificateDescription"], res)
	resultRes.PemCertificateChain = dcl.FlattenStringSlice(m["pemCertificateChain"])
	resultRes.CreateTime = dcl.FlattenString(m["createTime"])
	resultRes.UpdateTime = dcl.FlattenString(m["updateTime"])
	resultRes.Labels = dcl.FlattenKeyValuePairs(m["labels"])
	resultRes.Project = dcl.FlattenString(m["project"])
	resultRes.Location = dcl.FlattenString(m["location"])
	resultRes.CaPool = dcl.FlattenString(m["caPool"])
	resultRes.CertificateAuthority = dcl.FlattenString(m["certificateAuthority"])

	return resultRes
}

// expandCertificateConfigMap expands the contents of CertificateConfig into a JSON
// request object.
func expandCertificateConfigMap(c *Client, f map[string]CertificateConfig, res *Certificate) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandCertificateConfig(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandCertificateConfigSlice expands the contents of CertificateConfig into a JSON
// request object.
func expandCertificateConfigSlice(c *Client, f []CertificateConfig, res *Certificate) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandCertificateConfig(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenCertificateConfigMap flattens the contents of CertificateConfig from a JSON
// response object.
func flattenCertificateConfigMap(c *Client, i interface{}, res *Certificate) map[string]CertificateConfig {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]CertificateConfig{}
	}

	if len(a) == 0 {
		return map[string]CertificateConfig{}
	}

	items := make(map[string]CertificateConfig)
	for k, item := range a {
		items[k] = *flattenCertificateConfig(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenCertificateConfigSlice flattens the contents of CertificateConfig from a JSON
// response object.
func flattenCertificateConfigSlice(c *Client, i interface{}, res *Certificate) []CertificateConfig {
	a, ok := i.([]interface{})
	if !ok {
		return []CertificateConfig{}
	}

	if len(a) == 0 {
		return []CertificateConfig{}
	}

	items := make([]CertificateConfig, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenCertificateConfig(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandCertificateConfig expands an instance of CertificateConfig into a JSON
// request object.
func expandCertificateConfig(c *Client, f *CertificateConfig, res *Certificate) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v, err := expandCertificateConfigSubjectConfig(c, f.SubjectConfig, res); err != nil {
		return nil, fmt.Errorf("error expanding SubjectConfig into subjectConfig: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["subjectConfig"] = v
	}
	if v, err := expandCertificateConfigX509Config(c, f.X509Config, res); err != nil {
		return nil, fmt.Errorf("error expanding X509Config into x509Config: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["x509Config"] = v
	}
	if v, err := expandCertificateConfigPublicKey(c, f.PublicKey, res); err != nil {
		return nil, fmt.Errorf("error expanding PublicKey into publicKey: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["publicKey"] = v
	}

	return m, nil
}

// flattenCertificateConfig flattens an instance of CertificateConfig from a JSON
// response object.
func flattenCertificateConfig(c *Client, i interface{}, res *Certificate) *CertificateConfig {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &CertificateConfig{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyCertificateConfig
	}
	r.SubjectConfig = flattenCertificateConfigSubjectConfig(c, m["subjectConfig"], res)
	r.X509Config = flattenCertificateConfigX509Config(c, m["x509Config"], res)
	r.PublicKey = flattenCertificateConfigPublicKey(c, m["publicKey"], res)

	return r
}

// expandCertificateConfigSubjectConfigMap expands the contents of CertificateConfigSubjectConfig into a JSON
// request object.
func expandCertificateConfigSubjectConfigMap(c *Client, f map[string]CertificateConfigSubjectConfig, res *Certificate) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandCertificateConfigSubjectConfig(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandCertificateConfigSubjectConfigSlice expands the contents of CertificateConfigSubjectConfig into a JSON
// request object.
func expandCertificateConfigSubjectConfigSlice(c *Client, f []CertificateConfigSubjectConfig, res *Certificate) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandCertificateConfigSubjectConfig(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenCertificateConfigSubjectConfigMap flattens the contents of CertificateConfigSubjectConfig from a JSON
// response object.
func flattenCertificateConfigSubjectConfigMap(c *Client, i interface{}, res *Certificate) map[string]CertificateConfigSubjectConfig {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]CertificateConfigSubjectConfig{}
	}

	if len(a) == 0 {
		return map[string]CertificateConfigSubjectConfig{}
	}

	items := make(map[string]CertificateConfigSubjectConfig)
	for k, item := range a {
		items[k] = *flattenCertificateConfigSubjectConfig(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenCertificateConfigSubjectConfigSlice flattens the contents of CertificateConfigSubjectConfig from a JSON
// response object.
func flattenCertificateConfigSubjectConfigSlice(c *Client, i interface{}, res *Certificate) []CertificateConfigSubjectConfig {
	a, ok := i.([]interface{})
	if !ok {
		return []CertificateConfigSubjectConfig{}
	}

	if len(a) == 0 {
		return []CertificateConfigSubjectConfig{}
	}

	items := make([]CertificateConfigSubjectConfig, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenCertificateConfigSubjectConfig(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandCertificateConfigSubjectConfig expands an instance of CertificateConfigSubjectConfig into a JSON
// request object.
func expandCertificateConfigSubjectConfig(c *Client, f *CertificateConfigSubjectConfig, res *Certificate) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v, err := expandCertificateConfigSubjectConfigSubject(c, f.Subject, res); err != nil {
		return nil, fmt.Errorf("error expanding Subject into subject: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["subject"] = v
	}
	if v, err := expandCertificateConfigSubjectConfigSubjectAltName(c, f.SubjectAltName, res); err != nil {
		return nil, fmt.Errorf("error expanding SubjectAltName into subjectAltName: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["subjectAltName"] = v
	}

	return m, nil
}

// flattenCertificateConfigSubjectConfig flattens an instance of CertificateConfigSubjectConfig from a JSON
// response object.
func flattenCertificateConfigSubjectConfig(c *Client, i interface{}, res *Certificate) *CertificateConfigSubjectConfig {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &CertificateConfigSubjectConfig{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyCertificateConfigSubjectConfig
	}
	r.Subject = flattenCertificateConfigSubjectConfigSubject(c, m["subject"], res)
	r.SubjectAltName = flattenCertificateConfigSubjectConfigSubjectAltName(c, m["subjectAltName"], res)

	return r
}

// expandCertificateConfigSubjectConfigSubjectMap expands the contents of CertificateConfigSubjectConfigSubject into a JSON
// request object.
func expandCertificateConfigSubjectConfigSubjectMap(c *Client, f map[string]CertificateConfigSubjectConfigSubject, res *Certificate) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandCertificateConfigSubjectConfigSubject(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandCertificateConfigSubjectConfigSubjectSlice expands the contents of CertificateConfigSubjectConfigSubject into a JSON
// request object.
func expandCertificateConfigSubjectConfigSubjectSlice(c *Client, f []CertificateConfigSubjectConfigSubject, res *Certificate) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandCertificateConfigSubjectConfigSubject(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenCertificateConfigSubjectConfigSubjectMap flattens the contents of CertificateConfigSubjectConfigSubject from a JSON
// response object.
func flattenCertificateConfigSubjectConfigSubjectMap(c *Client, i interface{}, res *Certificate) map[string]CertificateConfigSubjectConfigSubject {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]CertificateConfigSubjectConfigSubject{}
	}

	if len(a) == 0 {
		return map[string]CertificateConfigSubjectConfigSubject{}
	}

	items := make(map[string]CertificateConfigSubjectConfigSubject)
	for k, item := range a {
		items[k] = *flattenCertificateConfigSubjectConfigSubject(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenCertificateConfigSubjectConfigSubjectSlice flattens the contents of CertificateConfigSubjectConfigSubject from a JSON
// response object.
func flattenCertificateConfigSubjectConfigSubjectSlice(c *Client, i interface{}, res *Certificate) []CertificateConfigSubjectConfigSubject {
	a, ok := i.([]interface{})
	if !ok {
		return []CertificateConfigSubjectConfigSubject{}
	}

	if len(a) == 0 {
		return []CertificateConfigSubjectConfigSubject{}
	}

	items := make([]CertificateConfigSubjectConfigSubject, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenCertificateConfigSubjectConfigSubject(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandCertificateConfigSubjectConfigSubject expands an instance of CertificateConfigSubjectConfigSubject into a JSON
// request object.
func expandCertificateConfigSubjectConfigSubject(c *Client, f *CertificateConfigSubjectConfigSubject, res *Certificate) (map[string]interface{}, error) {
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

// flattenCertificateConfigSubjectConfigSubject flattens an instance of CertificateConfigSubjectConfigSubject from a JSON
// response object.
func flattenCertificateConfigSubjectConfigSubject(c *Client, i interface{}, res *Certificate) *CertificateConfigSubjectConfigSubject {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &CertificateConfigSubjectConfigSubject{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyCertificateConfigSubjectConfigSubject
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

// expandCertificateConfigSubjectConfigSubjectAltNameMap expands the contents of CertificateConfigSubjectConfigSubjectAltName into a JSON
// request object.
func expandCertificateConfigSubjectConfigSubjectAltNameMap(c *Client, f map[string]CertificateConfigSubjectConfigSubjectAltName, res *Certificate) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandCertificateConfigSubjectConfigSubjectAltName(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandCertificateConfigSubjectConfigSubjectAltNameSlice expands the contents of CertificateConfigSubjectConfigSubjectAltName into a JSON
// request object.
func expandCertificateConfigSubjectConfigSubjectAltNameSlice(c *Client, f []CertificateConfigSubjectConfigSubjectAltName, res *Certificate) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandCertificateConfigSubjectConfigSubjectAltName(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenCertificateConfigSubjectConfigSubjectAltNameMap flattens the contents of CertificateConfigSubjectConfigSubjectAltName from a JSON
// response object.
func flattenCertificateConfigSubjectConfigSubjectAltNameMap(c *Client, i interface{}, res *Certificate) map[string]CertificateConfigSubjectConfigSubjectAltName {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]CertificateConfigSubjectConfigSubjectAltName{}
	}

	if len(a) == 0 {
		return map[string]CertificateConfigSubjectConfigSubjectAltName{}
	}

	items := make(map[string]CertificateConfigSubjectConfigSubjectAltName)
	for k, item := range a {
		items[k] = *flattenCertificateConfigSubjectConfigSubjectAltName(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenCertificateConfigSubjectConfigSubjectAltNameSlice flattens the contents of CertificateConfigSubjectConfigSubjectAltName from a JSON
// response object.
func flattenCertificateConfigSubjectConfigSubjectAltNameSlice(c *Client, i interface{}, res *Certificate) []CertificateConfigSubjectConfigSubjectAltName {
	a, ok := i.([]interface{})
	if !ok {
		return []CertificateConfigSubjectConfigSubjectAltName{}
	}

	if len(a) == 0 {
		return []CertificateConfigSubjectConfigSubjectAltName{}
	}

	items := make([]CertificateConfigSubjectConfigSubjectAltName, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenCertificateConfigSubjectConfigSubjectAltName(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandCertificateConfigSubjectConfigSubjectAltName expands an instance of CertificateConfigSubjectConfigSubjectAltName into a JSON
// request object.
func expandCertificateConfigSubjectConfigSubjectAltName(c *Client, f *CertificateConfigSubjectConfigSubjectAltName, res *Certificate) (map[string]interface{}, error) {
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

	return m, nil
}

// flattenCertificateConfigSubjectConfigSubjectAltName flattens an instance of CertificateConfigSubjectConfigSubjectAltName from a JSON
// response object.
func flattenCertificateConfigSubjectConfigSubjectAltName(c *Client, i interface{}, res *Certificate) *CertificateConfigSubjectConfigSubjectAltName {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &CertificateConfigSubjectConfigSubjectAltName{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyCertificateConfigSubjectConfigSubjectAltName
	}
	r.DnsNames = dcl.FlattenStringSlice(m["dnsNames"])
	r.Uris = dcl.FlattenStringSlice(m["uris"])
	r.EmailAddresses = dcl.FlattenStringSlice(m["emailAddresses"])
	r.IPAddresses = dcl.FlattenStringSlice(m["ipAddresses"])

	return r
}

// expandCertificateConfigX509ConfigMap expands the contents of CertificateConfigX509Config into a JSON
// request object.
func expandCertificateConfigX509ConfigMap(c *Client, f map[string]CertificateConfigX509Config, res *Certificate) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandCertificateConfigX509Config(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandCertificateConfigX509ConfigSlice expands the contents of CertificateConfigX509Config into a JSON
// request object.
func expandCertificateConfigX509ConfigSlice(c *Client, f []CertificateConfigX509Config, res *Certificate) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandCertificateConfigX509Config(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenCertificateConfigX509ConfigMap flattens the contents of CertificateConfigX509Config from a JSON
// response object.
func flattenCertificateConfigX509ConfigMap(c *Client, i interface{}, res *Certificate) map[string]CertificateConfigX509Config {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]CertificateConfigX509Config{}
	}

	if len(a) == 0 {
		return map[string]CertificateConfigX509Config{}
	}

	items := make(map[string]CertificateConfigX509Config)
	for k, item := range a {
		items[k] = *flattenCertificateConfigX509Config(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenCertificateConfigX509ConfigSlice flattens the contents of CertificateConfigX509Config from a JSON
// response object.
func flattenCertificateConfigX509ConfigSlice(c *Client, i interface{}, res *Certificate) []CertificateConfigX509Config {
	a, ok := i.([]interface{})
	if !ok {
		return []CertificateConfigX509Config{}
	}

	if len(a) == 0 {
		return []CertificateConfigX509Config{}
	}

	items := make([]CertificateConfigX509Config, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenCertificateConfigX509Config(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandCertificateConfigX509Config expands an instance of CertificateConfigX509Config into a JSON
// request object.
func expandCertificateConfigX509Config(c *Client, f *CertificateConfigX509Config, res *Certificate) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v, err := expandCertificateConfigX509ConfigKeyUsage(c, f.KeyUsage, res); err != nil {
		return nil, fmt.Errorf("error expanding KeyUsage into keyUsage: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["keyUsage"] = v
	}
	if v, err := expandCertificateConfigX509ConfigCAOptions(c, f.CaOptions, res); err != nil {
		return nil, fmt.Errorf("error expanding CaOptions into caOptions: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["caOptions"] = v
	}
	if v, err := expandCertificateConfigX509ConfigPolicyIdsSlice(c, f.PolicyIds, res); err != nil {
		return nil, fmt.Errorf("error expanding PolicyIds into policyIds: %w", err)
	} else if v != nil {
		m["policyIds"] = v
	}
	if v := f.AiaOcspServers; v != nil {
		m["aiaOcspServers"] = v
	}
	if v, err := expandCertificateConfigX509ConfigAdditionalExtensionsSlice(c, f.AdditionalExtensions, res); err != nil {
		return nil, fmt.Errorf("error expanding AdditionalExtensions into additionalExtensions: %w", err)
	} else if v != nil {
		m["additionalExtensions"] = v
	}

	return m, nil
}

// flattenCertificateConfigX509Config flattens an instance of CertificateConfigX509Config from a JSON
// response object.
func flattenCertificateConfigX509Config(c *Client, i interface{}, res *Certificate) *CertificateConfigX509Config {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &CertificateConfigX509Config{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyCertificateConfigX509Config
	}
	r.KeyUsage = flattenCertificateConfigX509ConfigKeyUsage(c, m["keyUsage"], res)
	r.CaOptions = flattenCertificateConfigX509ConfigCAOptions(c, m["caOptions"], res)
	r.PolicyIds = flattenCertificateConfigX509ConfigPolicyIdsSlice(c, m["policyIds"], res)
	r.AiaOcspServers = dcl.FlattenStringSlice(m["aiaOcspServers"])
	r.AdditionalExtensions = flattenCertificateConfigX509ConfigAdditionalExtensionsSlice(c, m["additionalExtensions"], res)

	return r
}

// expandCertificateConfigX509ConfigKeyUsageMap expands the contents of CertificateConfigX509ConfigKeyUsage into a JSON
// request object.
func expandCertificateConfigX509ConfigKeyUsageMap(c *Client, f map[string]CertificateConfigX509ConfigKeyUsage, res *Certificate) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandCertificateConfigX509ConfigKeyUsage(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandCertificateConfigX509ConfigKeyUsageSlice expands the contents of CertificateConfigX509ConfigKeyUsage into a JSON
// request object.
func expandCertificateConfigX509ConfigKeyUsageSlice(c *Client, f []CertificateConfigX509ConfigKeyUsage, res *Certificate) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandCertificateConfigX509ConfigKeyUsage(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenCertificateConfigX509ConfigKeyUsageMap flattens the contents of CertificateConfigX509ConfigKeyUsage from a JSON
// response object.
func flattenCertificateConfigX509ConfigKeyUsageMap(c *Client, i interface{}, res *Certificate) map[string]CertificateConfigX509ConfigKeyUsage {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]CertificateConfigX509ConfigKeyUsage{}
	}

	if len(a) == 0 {
		return map[string]CertificateConfigX509ConfigKeyUsage{}
	}

	items := make(map[string]CertificateConfigX509ConfigKeyUsage)
	for k, item := range a {
		items[k] = *flattenCertificateConfigX509ConfigKeyUsage(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenCertificateConfigX509ConfigKeyUsageSlice flattens the contents of CertificateConfigX509ConfigKeyUsage from a JSON
// response object.
func flattenCertificateConfigX509ConfigKeyUsageSlice(c *Client, i interface{}, res *Certificate) []CertificateConfigX509ConfigKeyUsage {
	a, ok := i.([]interface{})
	if !ok {
		return []CertificateConfigX509ConfigKeyUsage{}
	}

	if len(a) == 0 {
		return []CertificateConfigX509ConfigKeyUsage{}
	}

	items := make([]CertificateConfigX509ConfigKeyUsage, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenCertificateConfigX509ConfigKeyUsage(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandCertificateConfigX509ConfigKeyUsage expands an instance of CertificateConfigX509ConfigKeyUsage into a JSON
// request object.
func expandCertificateConfigX509ConfigKeyUsage(c *Client, f *CertificateConfigX509ConfigKeyUsage, res *Certificate) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v, err := expandCertificateConfigX509ConfigKeyUsageBaseKeyUsage(c, f.BaseKeyUsage, res); err != nil {
		return nil, fmt.Errorf("error expanding BaseKeyUsage into baseKeyUsage: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["baseKeyUsage"] = v
	}
	if v, err := expandCertificateConfigX509ConfigKeyUsageExtendedKeyUsage(c, f.ExtendedKeyUsage, res); err != nil {
		return nil, fmt.Errorf("error expanding ExtendedKeyUsage into extendedKeyUsage: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["extendedKeyUsage"] = v
	}
	if v, err := expandCertificateConfigX509ConfigKeyUsageUnknownExtendedKeyUsagesSlice(c, f.UnknownExtendedKeyUsages, res); err != nil {
		return nil, fmt.Errorf("error expanding UnknownExtendedKeyUsages into unknownExtendedKeyUsages: %w", err)
	} else if v != nil {
		m["unknownExtendedKeyUsages"] = v
	}

	return m, nil
}

// flattenCertificateConfigX509ConfigKeyUsage flattens an instance of CertificateConfigX509ConfigKeyUsage from a JSON
// response object.
func flattenCertificateConfigX509ConfigKeyUsage(c *Client, i interface{}, res *Certificate) *CertificateConfigX509ConfigKeyUsage {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &CertificateConfigX509ConfigKeyUsage{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyCertificateConfigX509ConfigKeyUsage
	}
	r.BaseKeyUsage = flattenCertificateConfigX509ConfigKeyUsageBaseKeyUsage(c, m["baseKeyUsage"], res)
	r.ExtendedKeyUsage = flattenCertificateConfigX509ConfigKeyUsageExtendedKeyUsage(c, m["extendedKeyUsage"], res)
	r.UnknownExtendedKeyUsages = flattenCertificateConfigX509ConfigKeyUsageUnknownExtendedKeyUsagesSlice(c, m["unknownExtendedKeyUsages"], res)

	return r
}

// expandCertificateConfigX509ConfigKeyUsageBaseKeyUsageMap expands the contents of CertificateConfigX509ConfigKeyUsageBaseKeyUsage into a JSON
// request object.
func expandCertificateConfigX509ConfigKeyUsageBaseKeyUsageMap(c *Client, f map[string]CertificateConfigX509ConfigKeyUsageBaseKeyUsage, res *Certificate) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandCertificateConfigX509ConfigKeyUsageBaseKeyUsage(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandCertificateConfigX509ConfigKeyUsageBaseKeyUsageSlice expands the contents of CertificateConfigX509ConfigKeyUsageBaseKeyUsage into a JSON
// request object.
func expandCertificateConfigX509ConfigKeyUsageBaseKeyUsageSlice(c *Client, f []CertificateConfigX509ConfigKeyUsageBaseKeyUsage, res *Certificate) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandCertificateConfigX509ConfigKeyUsageBaseKeyUsage(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenCertificateConfigX509ConfigKeyUsageBaseKeyUsageMap flattens the contents of CertificateConfigX509ConfigKeyUsageBaseKeyUsage from a JSON
// response object.
func flattenCertificateConfigX509ConfigKeyUsageBaseKeyUsageMap(c *Client, i interface{}, res *Certificate) map[string]CertificateConfigX509ConfigKeyUsageBaseKeyUsage {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]CertificateConfigX509ConfigKeyUsageBaseKeyUsage{}
	}

	if len(a) == 0 {
		return map[string]CertificateConfigX509ConfigKeyUsageBaseKeyUsage{}
	}

	items := make(map[string]CertificateConfigX509ConfigKeyUsageBaseKeyUsage)
	for k, item := range a {
		items[k] = *flattenCertificateConfigX509ConfigKeyUsageBaseKeyUsage(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenCertificateConfigX509ConfigKeyUsageBaseKeyUsageSlice flattens the contents of CertificateConfigX509ConfigKeyUsageBaseKeyUsage from a JSON
// response object.
func flattenCertificateConfigX509ConfigKeyUsageBaseKeyUsageSlice(c *Client, i interface{}, res *Certificate) []CertificateConfigX509ConfigKeyUsageBaseKeyUsage {
	a, ok := i.([]interface{})
	if !ok {
		return []CertificateConfigX509ConfigKeyUsageBaseKeyUsage{}
	}

	if len(a) == 0 {
		return []CertificateConfigX509ConfigKeyUsageBaseKeyUsage{}
	}

	items := make([]CertificateConfigX509ConfigKeyUsageBaseKeyUsage, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenCertificateConfigX509ConfigKeyUsageBaseKeyUsage(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandCertificateConfigX509ConfigKeyUsageBaseKeyUsage expands an instance of CertificateConfigX509ConfigKeyUsageBaseKeyUsage into a JSON
// request object.
func expandCertificateConfigX509ConfigKeyUsageBaseKeyUsage(c *Client, f *CertificateConfigX509ConfigKeyUsageBaseKeyUsage, res *Certificate) (map[string]interface{}, error) {
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

// flattenCertificateConfigX509ConfigKeyUsageBaseKeyUsage flattens an instance of CertificateConfigX509ConfigKeyUsageBaseKeyUsage from a JSON
// response object.
func flattenCertificateConfigX509ConfigKeyUsageBaseKeyUsage(c *Client, i interface{}, res *Certificate) *CertificateConfigX509ConfigKeyUsageBaseKeyUsage {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &CertificateConfigX509ConfigKeyUsageBaseKeyUsage{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyCertificateConfigX509ConfigKeyUsageBaseKeyUsage
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

// expandCertificateConfigX509ConfigKeyUsageExtendedKeyUsageMap expands the contents of CertificateConfigX509ConfigKeyUsageExtendedKeyUsage into a JSON
// request object.
func expandCertificateConfigX509ConfigKeyUsageExtendedKeyUsageMap(c *Client, f map[string]CertificateConfigX509ConfigKeyUsageExtendedKeyUsage, res *Certificate) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandCertificateConfigX509ConfigKeyUsageExtendedKeyUsage(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandCertificateConfigX509ConfigKeyUsageExtendedKeyUsageSlice expands the contents of CertificateConfigX509ConfigKeyUsageExtendedKeyUsage into a JSON
// request object.
func expandCertificateConfigX509ConfigKeyUsageExtendedKeyUsageSlice(c *Client, f []CertificateConfigX509ConfigKeyUsageExtendedKeyUsage, res *Certificate) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandCertificateConfigX509ConfigKeyUsageExtendedKeyUsage(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenCertificateConfigX509ConfigKeyUsageExtendedKeyUsageMap flattens the contents of CertificateConfigX509ConfigKeyUsageExtendedKeyUsage from a JSON
// response object.
func flattenCertificateConfigX509ConfigKeyUsageExtendedKeyUsageMap(c *Client, i interface{}, res *Certificate) map[string]CertificateConfigX509ConfigKeyUsageExtendedKeyUsage {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]CertificateConfigX509ConfigKeyUsageExtendedKeyUsage{}
	}

	if len(a) == 0 {
		return map[string]CertificateConfigX509ConfigKeyUsageExtendedKeyUsage{}
	}

	items := make(map[string]CertificateConfigX509ConfigKeyUsageExtendedKeyUsage)
	for k, item := range a {
		items[k] = *flattenCertificateConfigX509ConfigKeyUsageExtendedKeyUsage(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenCertificateConfigX509ConfigKeyUsageExtendedKeyUsageSlice flattens the contents of CertificateConfigX509ConfigKeyUsageExtendedKeyUsage from a JSON
// response object.
func flattenCertificateConfigX509ConfigKeyUsageExtendedKeyUsageSlice(c *Client, i interface{}, res *Certificate) []CertificateConfigX509ConfigKeyUsageExtendedKeyUsage {
	a, ok := i.([]interface{})
	if !ok {
		return []CertificateConfigX509ConfigKeyUsageExtendedKeyUsage{}
	}

	if len(a) == 0 {
		return []CertificateConfigX509ConfigKeyUsageExtendedKeyUsage{}
	}

	items := make([]CertificateConfigX509ConfigKeyUsageExtendedKeyUsage, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenCertificateConfigX509ConfigKeyUsageExtendedKeyUsage(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandCertificateConfigX509ConfigKeyUsageExtendedKeyUsage expands an instance of CertificateConfigX509ConfigKeyUsageExtendedKeyUsage into a JSON
// request object.
func expandCertificateConfigX509ConfigKeyUsageExtendedKeyUsage(c *Client, f *CertificateConfigX509ConfigKeyUsageExtendedKeyUsage, res *Certificate) (map[string]interface{}, error) {
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

// flattenCertificateConfigX509ConfigKeyUsageExtendedKeyUsage flattens an instance of CertificateConfigX509ConfigKeyUsageExtendedKeyUsage from a JSON
// response object.
func flattenCertificateConfigX509ConfigKeyUsageExtendedKeyUsage(c *Client, i interface{}, res *Certificate) *CertificateConfigX509ConfigKeyUsageExtendedKeyUsage {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &CertificateConfigX509ConfigKeyUsageExtendedKeyUsage{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyCertificateConfigX509ConfigKeyUsageExtendedKeyUsage
	}
	r.ServerAuth = dcl.FlattenBool(m["serverAuth"])
	r.ClientAuth = dcl.FlattenBool(m["clientAuth"])
	r.CodeSigning = dcl.FlattenBool(m["codeSigning"])
	r.EmailProtection = dcl.FlattenBool(m["emailProtection"])
	r.TimeStamping = dcl.FlattenBool(m["timeStamping"])
	r.OcspSigning = dcl.FlattenBool(m["ocspSigning"])

	return r
}

// expandCertificateConfigX509ConfigKeyUsageUnknownExtendedKeyUsagesMap expands the contents of CertificateConfigX509ConfigKeyUsageUnknownExtendedKeyUsages into a JSON
// request object.
func expandCertificateConfigX509ConfigKeyUsageUnknownExtendedKeyUsagesMap(c *Client, f map[string]CertificateConfigX509ConfigKeyUsageUnknownExtendedKeyUsages, res *Certificate) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandCertificateConfigX509ConfigKeyUsageUnknownExtendedKeyUsages(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandCertificateConfigX509ConfigKeyUsageUnknownExtendedKeyUsagesSlice expands the contents of CertificateConfigX509ConfigKeyUsageUnknownExtendedKeyUsages into a JSON
// request object.
func expandCertificateConfigX509ConfigKeyUsageUnknownExtendedKeyUsagesSlice(c *Client, f []CertificateConfigX509ConfigKeyUsageUnknownExtendedKeyUsages, res *Certificate) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandCertificateConfigX509ConfigKeyUsageUnknownExtendedKeyUsages(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenCertificateConfigX509ConfigKeyUsageUnknownExtendedKeyUsagesMap flattens the contents of CertificateConfigX509ConfigKeyUsageUnknownExtendedKeyUsages from a JSON
// response object.
func flattenCertificateConfigX509ConfigKeyUsageUnknownExtendedKeyUsagesMap(c *Client, i interface{}, res *Certificate) map[string]CertificateConfigX509ConfigKeyUsageUnknownExtendedKeyUsages {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]CertificateConfigX509ConfigKeyUsageUnknownExtendedKeyUsages{}
	}

	if len(a) == 0 {
		return map[string]CertificateConfigX509ConfigKeyUsageUnknownExtendedKeyUsages{}
	}

	items := make(map[string]CertificateConfigX509ConfigKeyUsageUnknownExtendedKeyUsages)
	for k, item := range a {
		items[k] = *flattenCertificateConfigX509ConfigKeyUsageUnknownExtendedKeyUsages(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenCertificateConfigX509ConfigKeyUsageUnknownExtendedKeyUsagesSlice flattens the contents of CertificateConfigX509ConfigKeyUsageUnknownExtendedKeyUsages from a JSON
// response object.
func flattenCertificateConfigX509ConfigKeyUsageUnknownExtendedKeyUsagesSlice(c *Client, i interface{}, res *Certificate) []CertificateConfigX509ConfigKeyUsageUnknownExtendedKeyUsages {
	a, ok := i.([]interface{})
	if !ok {
		return []CertificateConfigX509ConfigKeyUsageUnknownExtendedKeyUsages{}
	}

	if len(a) == 0 {
		return []CertificateConfigX509ConfigKeyUsageUnknownExtendedKeyUsages{}
	}

	items := make([]CertificateConfigX509ConfigKeyUsageUnknownExtendedKeyUsages, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenCertificateConfigX509ConfigKeyUsageUnknownExtendedKeyUsages(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandCertificateConfigX509ConfigKeyUsageUnknownExtendedKeyUsages expands an instance of CertificateConfigX509ConfigKeyUsageUnknownExtendedKeyUsages into a JSON
// request object.
func expandCertificateConfigX509ConfigKeyUsageUnknownExtendedKeyUsages(c *Client, f *CertificateConfigX509ConfigKeyUsageUnknownExtendedKeyUsages, res *Certificate) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.ObjectIdPath; v != nil {
		m["objectIdPath"] = v
	}

	return m, nil
}

// flattenCertificateConfigX509ConfigKeyUsageUnknownExtendedKeyUsages flattens an instance of CertificateConfigX509ConfigKeyUsageUnknownExtendedKeyUsages from a JSON
// response object.
func flattenCertificateConfigX509ConfigKeyUsageUnknownExtendedKeyUsages(c *Client, i interface{}, res *Certificate) *CertificateConfigX509ConfigKeyUsageUnknownExtendedKeyUsages {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &CertificateConfigX509ConfigKeyUsageUnknownExtendedKeyUsages{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyCertificateConfigX509ConfigKeyUsageUnknownExtendedKeyUsages
	}
	r.ObjectIdPath = dcl.FlattenIntSlice(m["objectIdPath"])

	return r
}

// expandCertificateConfigX509ConfigCaOptionsMap expands the contents of CertificateConfigX509ConfigCaOptions into a JSON
// request object.
func expandCertificateConfigX509ConfigCaOptionsMap(c *Client, f map[string]CertificateConfigX509ConfigCaOptions, res *Certificate) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandCertificateConfigX509ConfigCaOptions(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandCertificateConfigX509ConfigCaOptionsSlice expands the contents of CertificateConfigX509ConfigCaOptions into a JSON
// request object.
func expandCertificateConfigX509ConfigCaOptionsSlice(c *Client, f []CertificateConfigX509ConfigCaOptions, res *Certificate) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandCertificateConfigX509ConfigCaOptions(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenCertificateConfigX509ConfigCaOptionsMap flattens the contents of CertificateConfigX509ConfigCaOptions from a JSON
// response object.
func flattenCertificateConfigX509ConfigCaOptionsMap(c *Client, i interface{}, res *Certificate) map[string]CertificateConfigX509ConfigCaOptions {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]CertificateConfigX509ConfigCaOptions{}
	}

	if len(a) == 0 {
		return map[string]CertificateConfigX509ConfigCaOptions{}
	}

	items := make(map[string]CertificateConfigX509ConfigCaOptions)
	for k, item := range a {
		items[k] = *flattenCertificateConfigX509ConfigCaOptions(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenCertificateConfigX509ConfigCaOptionsSlice flattens the contents of CertificateConfigX509ConfigCaOptions from a JSON
// response object.
func flattenCertificateConfigX509ConfigCaOptionsSlice(c *Client, i interface{}, res *Certificate) []CertificateConfigX509ConfigCaOptions {
	a, ok := i.([]interface{})
	if !ok {
		return []CertificateConfigX509ConfigCaOptions{}
	}

	if len(a) == 0 {
		return []CertificateConfigX509ConfigCaOptions{}
	}

	items := make([]CertificateConfigX509ConfigCaOptions, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenCertificateConfigX509ConfigCaOptions(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandCertificateConfigX509ConfigCaOptions expands an instance of CertificateConfigX509ConfigCaOptions into a JSON
// request object.
func expandCertificateConfigX509ConfigCaOptions(c *Client, f *CertificateConfigX509ConfigCaOptions, res *Certificate) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.IsCa; !dcl.IsEmptyValueIndirect(v) {
		m["isCa"] = v
	}
	if v := f.NonCa; !dcl.IsEmptyValueIndirect(v) {
		m["nonCa"] = v
	}
	if v := f.MaxIssuerPathLength; !dcl.IsEmptyValueIndirect(v) {
		m["maxIssuerPathLength"] = v
	}
	if v := f.ZeroMaxIssuerPathLength; !dcl.IsEmptyValueIndirect(v) {
		m["zeroMaxIssuerPathLength"] = v
	}

	return m, nil
}

// flattenCertificateConfigX509ConfigCaOptions flattens an instance of CertificateConfigX509ConfigCaOptions from a JSON
// response object.
func flattenCertificateConfigX509ConfigCaOptions(c *Client, i interface{}, res *Certificate) *CertificateConfigX509ConfigCaOptions {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &CertificateConfigX509ConfigCaOptions{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyCertificateConfigX509ConfigCaOptions
	}
	r.IsCa = dcl.FlattenBool(m["isCa"])
	r.NonCa = dcl.FlattenBool(m["nonCa"])
	r.MaxIssuerPathLength = dcl.FlattenInteger(m["maxIssuerPathLength"])
	r.ZeroMaxIssuerPathLength = dcl.FlattenBool(m["zeroMaxIssuerPathLength"])

	return r
}

// expandCertificateConfigX509ConfigPolicyIdsMap expands the contents of CertificateConfigX509ConfigPolicyIds into a JSON
// request object.
func expandCertificateConfigX509ConfigPolicyIdsMap(c *Client, f map[string]CertificateConfigX509ConfigPolicyIds, res *Certificate) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandCertificateConfigX509ConfigPolicyIds(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandCertificateConfigX509ConfigPolicyIdsSlice expands the contents of CertificateConfigX509ConfigPolicyIds into a JSON
// request object.
func expandCertificateConfigX509ConfigPolicyIdsSlice(c *Client, f []CertificateConfigX509ConfigPolicyIds, res *Certificate) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandCertificateConfigX509ConfigPolicyIds(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenCertificateConfigX509ConfigPolicyIdsMap flattens the contents of CertificateConfigX509ConfigPolicyIds from a JSON
// response object.
func flattenCertificateConfigX509ConfigPolicyIdsMap(c *Client, i interface{}, res *Certificate) map[string]CertificateConfigX509ConfigPolicyIds {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]CertificateConfigX509ConfigPolicyIds{}
	}

	if len(a) == 0 {
		return map[string]CertificateConfigX509ConfigPolicyIds{}
	}

	items := make(map[string]CertificateConfigX509ConfigPolicyIds)
	for k, item := range a {
		items[k] = *flattenCertificateConfigX509ConfigPolicyIds(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenCertificateConfigX509ConfigPolicyIdsSlice flattens the contents of CertificateConfigX509ConfigPolicyIds from a JSON
// response object.
func flattenCertificateConfigX509ConfigPolicyIdsSlice(c *Client, i interface{}, res *Certificate) []CertificateConfigX509ConfigPolicyIds {
	a, ok := i.([]interface{})
	if !ok {
		return []CertificateConfigX509ConfigPolicyIds{}
	}

	if len(a) == 0 {
		return []CertificateConfigX509ConfigPolicyIds{}
	}

	items := make([]CertificateConfigX509ConfigPolicyIds, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenCertificateConfigX509ConfigPolicyIds(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandCertificateConfigX509ConfigPolicyIds expands an instance of CertificateConfigX509ConfigPolicyIds into a JSON
// request object.
func expandCertificateConfigX509ConfigPolicyIds(c *Client, f *CertificateConfigX509ConfigPolicyIds, res *Certificate) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.ObjectIdPath; v != nil {
		m["objectIdPath"] = v
	}

	return m, nil
}

// flattenCertificateConfigX509ConfigPolicyIds flattens an instance of CertificateConfigX509ConfigPolicyIds from a JSON
// response object.
func flattenCertificateConfigX509ConfigPolicyIds(c *Client, i interface{}, res *Certificate) *CertificateConfigX509ConfigPolicyIds {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &CertificateConfigX509ConfigPolicyIds{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyCertificateConfigX509ConfigPolicyIds
	}
	r.ObjectIdPath = dcl.FlattenIntSlice(m["objectIdPath"])

	return r
}

// expandCertificateConfigX509ConfigAdditionalExtensionsMap expands the contents of CertificateConfigX509ConfigAdditionalExtensions into a JSON
// request object.
func expandCertificateConfigX509ConfigAdditionalExtensionsMap(c *Client, f map[string]CertificateConfigX509ConfigAdditionalExtensions, res *Certificate) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandCertificateConfigX509ConfigAdditionalExtensions(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandCertificateConfigX509ConfigAdditionalExtensionsSlice expands the contents of CertificateConfigX509ConfigAdditionalExtensions into a JSON
// request object.
func expandCertificateConfigX509ConfigAdditionalExtensionsSlice(c *Client, f []CertificateConfigX509ConfigAdditionalExtensions, res *Certificate) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandCertificateConfigX509ConfigAdditionalExtensions(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenCertificateConfigX509ConfigAdditionalExtensionsMap flattens the contents of CertificateConfigX509ConfigAdditionalExtensions from a JSON
// response object.
func flattenCertificateConfigX509ConfigAdditionalExtensionsMap(c *Client, i interface{}, res *Certificate) map[string]CertificateConfigX509ConfigAdditionalExtensions {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]CertificateConfigX509ConfigAdditionalExtensions{}
	}

	if len(a) == 0 {
		return map[string]CertificateConfigX509ConfigAdditionalExtensions{}
	}

	items := make(map[string]CertificateConfigX509ConfigAdditionalExtensions)
	for k, item := range a {
		items[k] = *flattenCertificateConfigX509ConfigAdditionalExtensions(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenCertificateConfigX509ConfigAdditionalExtensionsSlice flattens the contents of CertificateConfigX509ConfigAdditionalExtensions from a JSON
// response object.
func flattenCertificateConfigX509ConfigAdditionalExtensionsSlice(c *Client, i interface{}, res *Certificate) []CertificateConfigX509ConfigAdditionalExtensions {
	a, ok := i.([]interface{})
	if !ok {
		return []CertificateConfigX509ConfigAdditionalExtensions{}
	}

	if len(a) == 0 {
		return []CertificateConfigX509ConfigAdditionalExtensions{}
	}

	items := make([]CertificateConfigX509ConfigAdditionalExtensions, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenCertificateConfigX509ConfigAdditionalExtensions(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandCertificateConfigX509ConfigAdditionalExtensions expands an instance of CertificateConfigX509ConfigAdditionalExtensions into a JSON
// request object.
func expandCertificateConfigX509ConfigAdditionalExtensions(c *Client, f *CertificateConfigX509ConfigAdditionalExtensions, res *Certificate) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v, err := expandCertificateConfigX509ConfigAdditionalExtensionsObjectId(c, f.ObjectId, res); err != nil {
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

// flattenCertificateConfigX509ConfigAdditionalExtensions flattens an instance of CertificateConfigX509ConfigAdditionalExtensions from a JSON
// response object.
func flattenCertificateConfigX509ConfigAdditionalExtensions(c *Client, i interface{}, res *Certificate) *CertificateConfigX509ConfigAdditionalExtensions {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &CertificateConfigX509ConfigAdditionalExtensions{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyCertificateConfigX509ConfigAdditionalExtensions
	}
	r.ObjectId = flattenCertificateConfigX509ConfigAdditionalExtensionsObjectId(c, m["objectId"], res)
	r.Critical = dcl.FlattenBool(m["critical"])
	r.Value = dcl.FlattenString(m["value"])

	return r
}

// expandCertificateConfigX509ConfigAdditionalExtensionsObjectIdMap expands the contents of CertificateConfigX509ConfigAdditionalExtensionsObjectId into a JSON
// request object.
func expandCertificateConfigX509ConfigAdditionalExtensionsObjectIdMap(c *Client, f map[string]CertificateConfigX509ConfigAdditionalExtensionsObjectId, res *Certificate) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandCertificateConfigX509ConfigAdditionalExtensionsObjectId(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandCertificateConfigX509ConfigAdditionalExtensionsObjectIdSlice expands the contents of CertificateConfigX509ConfigAdditionalExtensionsObjectId into a JSON
// request object.
func expandCertificateConfigX509ConfigAdditionalExtensionsObjectIdSlice(c *Client, f []CertificateConfigX509ConfigAdditionalExtensionsObjectId, res *Certificate) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandCertificateConfigX509ConfigAdditionalExtensionsObjectId(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenCertificateConfigX509ConfigAdditionalExtensionsObjectIdMap flattens the contents of CertificateConfigX509ConfigAdditionalExtensionsObjectId from a JSON
// response object.
func flattenCertificateConfigX509ConfigAdditionalExtensionsObjectIdMap(c *Client, i interface{}, res *Certificate) map[string]CertificateConfigX509ConfigAdditionalExtensionsObjectId {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]CertificateConfigX509ConfigAdditionalExtensionsObjectId{}
	}

	if len(a) == 0 {
		return map[string]CertificateConfigX509ConfigAdditionalExtensionsObjectId{}
	}

	items := make(map[string]CertificateConfigX509ConfigAdditionalExtensionsObjectId)
	for k, item := range a {
		items[k] = *flattenCertificateConfigX509ConfigAdditionalExtensionsObjectId(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenCertificateConfigX509ConfigAdditionalExtensionsObjectIdSlice flattens the contents of CertificateConfigX509ConfigAdditionalExtensionsObjectId from a JSON
// response object.
func flattenCertificateConfigX509ConfigAdditionalExtensionsObjectIdSlice(c *Client, i interface{}, res *Certificate) []CertificateConfigX509ConfigAdditionalExtensionsObjectId {
	a, ok := i.([]interface{})
	if !ok {
		return []CertificateConfigX509ConfigAdditionalExtensionsObjectId{}
	}

	if len(a) == 0 {
		return []CertificateConfigX509ConfigAdditionalExtensionsObjectId{}
	}

	items := make([]CertificateConfigX509ConfigAdditionalExtensionsObjectId, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenCertificateConfigX509ConfigAdditionalExtensionsObjectId(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandCertificateConfigX509ConfigAdditionalExtensionsObjectId expands an instance of CertificateConfigX509ConfigAdditionalExtensionsObjectId into a JSON
// request object.
func expandCertificateConfigX509ConfigAdditionalExtensionsObjectId(c *Client, f *CertificateConfigX509ConfigAdditionalExtensionsObjectId, res *Certificate) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.ObjectIdPath; v != nil {
		m["objectIdPath"] = v
	}

	return m, nil
}

// flattenCertificateConfigX509ConfigAdditionalExtensionsObjectId flattens an instance of CertificateConfigX509ConfigAdditionalExtensionsObjectId from a JSON
// response object.
func flattenCertificateConfigX509ConfigAdditionalExtensionsObjectId(c *Client, i interface{}, res *Certificate) *CertificateConfigX509ConfigAdditionalExtensionsObjectId {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &CertificateConfigX509ConfigAdditionalExtensionsObjectId{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyCertificateConfigX509ConfigAdditionalExtensionsObjectId
	}
	r.ObjectIdPath = dcl.FlattenIntSlice(m["objectIdPath"])

	return r
}

// expandCertificateConfigPublicKeyMap expands the contents of CertificateConfigPublicKey into a JSON
// request object.
func expandCertificateConfigPublicKeyMap(c *Client, f map[string]CertificateConfigPublicKey, res *Certificate) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandCertificateConfigPublicKey(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandCertificateConfigPublicKeySlice expands the contents of CertificateConfigPublicKey into a JSON
// request object.
func expandCertificateConfigPublicKeySlice(c *Client, f []CertificateConfigPublicKey, res *Certificate) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandCertificateConfigPublicKey(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenCertificateConfigPublicKeyMap flattens the contents of CertificateConfigPublicKey from a JSON
// response object.
func flattenCertificateConfigPublicKeyMap(c *Client, i interface{}, res *Certificate) map[string]CertificateConfigPublicKey {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]CertificateConfigPublicKey{}
	}

	if len(a) == 0 {
		return map[string]CertificateConfigPublicKey{}
	}

	items := make(map[string]CertificateConfigPublicKey)
	for k, item := range a {
		items[k] = *flattenCertificateConfigPublicKey(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenCertificateConfigPublicKeySlice flattens the contents of CertificateConfigPublicKey from a JSON
// response object.
func flattenCertificateConfigPublicKeySlice(c *Client, i interface{}, res *Certificate) []CertificateConfigPublicKey {
	a, ok := i.([]interface{})
	if !ok {
		return []CertificateConfigPublicKey{}
	}

	if len(a) == 0 {
		return []CertificateConfigPublicKey{}
	}

	items := make([]CertificateConfigPublicKey, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenCertificateConfigPublicKey(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandCertificateConfigPublicKey expands an instance of CertificateConfigPublicKey into a JSON
// request object.
func expandCertificateConfigPublicKey(c *Client, f *CertificateConfigPublicKey, res *Certificate) (map[string]interface{}, error) {
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

// flattenCertificateConfigPublicKey flattens an instance of CertificateConfigPublicKey from a JSON
// response object.
func flattenCertificateConfigPublicKey(c *Client, i interface{}, res *Certificate) *CertificateConfigPublicKey {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &CertificateConfigPublicKey{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyCertificateConfigPublicKey
	}
	r.Key = dcl.FlattenString(m["key"])
	r.Format = flattenCertificateConfigPublicKeyFormatEnum(m["format"])

	return r
}

// expandCertificateRevocationDetailsMap expands the contents of CertificateRevocationDetails into a JSON
// request object.
func expandCertificateRevocationDetailsMap(c *Client, f map[string]CertificateRevocationDetails, res *Certificate) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandCertificateRevocationDetails(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandCertificateRevocationDetailsSlice expands the contents of CertificateRevocationDetails into a JSON
// request object.
func expandCertificateRevocationDetailsSlice(c *Client, f []CertificateRevocationDetails, res *Certificate) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandCertificateRevocationDetails(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenCertificateRevocationDetailsMap flattens the contents of CertificateRevocationDetails from a JSON
// response object.
func flattenCertificateRevocationDetailsMap(c *Client, i interface{}, res *Certificate) map[string]CertificateRevocationDetails {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]CertificateRevocationDetails{}
	}

	if len(a) == 0 {
		return map[string]CertificateRevocationDetails{}
	}

	items := make(map[string]CertificateRevocationDetails)
	for k, item := range a {
		items[k] = *flattenCertificateRevocationDetails(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenCertificateRevocationDetailsSlice flattens the contents of CertificateRevocationDetails from a JSON
// response object.
func flattenCertificateRevocationDetailsSlice(c *Client, i interface{}, res *Certificate) []CertificateRevocationDetails {
	a, ok := i.([]interface{})
	if !ok {
		return []CertificateRevocationDetails{}
	}

	if len(a) == 0 {
		return []CertificateRevocationDetails{}
	}

	items := make([]CertificateRevocationDetails, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenCertificateRevocationDetails(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandCertificateRevocationDetails expands an instance of CertificateRevocationDetails into a JSON
// request object.
func expandCertificateRevocationDetails(c *Client, f *CertificateRevocationDetails, res *Certificate) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.RevocationState; !dcl.IsEmptyValueIndirect(v) {
		m["revocationState"] = v
	}
	if v := f.RevocationTime; !dcl.IsEmptyValueIndirect(v) {
		m["revocationTime"] = v
	}

	return m, nil
}

// flattenCertificateRevocationDetails flattens an instance of CertificateRevocationDetails from a JSON
// response object.
func flattenCertificateRevocationDetails(c *Client, i interface{}, res *Certificate) *CertificateRevocationDetails {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &CertificateRevocationDetails{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyCertificateRevocationDetails
	}
	r.RevocationState = flattenCertificateRevocationDetailsRevocationStateEnum(m["revocationState"])
	r.RevocationTime = dcl.FlattenString(m["revocationTime"])

	return r
}

// expandCertificateCertificateDescriptionMap expands the contents of CertificateCertificateDescription into a JSON
// request object.
func expandCertificateCertificateDescriptionMap(c *Client, f map[string]CertificateCertificateDescription, res *Certificate) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandCertificateCertificateDescription(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandCertificateCertificateDescriptionSlice expands the contents of CertificateCertificateDescription into a JSON
// request object.
func expandCertificateCertificateDescriptionSlice(c *Client, f []CertificateCertificateDescription, res *Certificate) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandCertificateCertificateDescription(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenCertificateCertificateDescriptionMap flattens the contents of CertificateCertificateDescription from a JSON
// response object.
func flattenCertificateCertificateDescriptionMap(c *Client, i interface{}, res *Certificate) map[string]CertificateCertificateDescription {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]CertificateCertificateDescription{}
	}

	if len(a) == 0 {
		return map[string]CertificateCertificateDescription{}
	}

	items := make(map[string]CertificateCertificateDescription)
	for k, item := range a {
		items[k] = *flattenCertificateCertificateDescription(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenCertificateCertificateDescriptionSlice flattens the contents of CertificateCertificateDescription from a JSON
// response object.
func flattenCertificateCertificateDescriptionSlice(c *Client, i interface{}, res *Certificate) []CertificateCertificateDescription {
	a, ok := i.([]interface{})
	if !ok {
		return []CertificateCertificateDescription{}
	}

	if len(a) == 0 {
		return []CertificateCertificateDescription{}
	}

	items := make([]CertificateCertificateDescription, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenCertificateCertificateDescription(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandCertificateCertificateDescription expands an instance of CertificateCertificateDescription into a JSON
// request object.
func expandCertificateCertificateDescription(c *Client, f *CertificateCertificateDescription, res *Certificate) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v, err := expandCertificateCertificateDescriptionSubjectDescription(c, f.SubjectDescription, res); err != nil {
		return nil, fmt.Errorf("error expanding SubjectDescription into subjectDescription: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["subjectDescription"] = v
	}
	if v, err := expandCertificateCertificateDescriptionX509Description(c, f.X509Description, res); err != nil {
		return nil, fmt.Errorf("error expanding X509Description into x509Description: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["x509Description"] = v
	}
	if v, err := expandCertificateCertificateDescriptionPublicKey(c, f.PublicKey, res); err != nil {
		return nil, fmt.Errorf("error expanding PublicKey into publicKey: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["publicKey"] = v
	}
	if v, err := expandCertificateCertificateDescriptionSubjectKeyId(c, f.SubjectKeyId, res); err != nil {
		return nil, fmt.Errorf("error expanding SubjectKeyId into subjectKeyId: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["subjectKeyId"] = v
	}
	if v, err := expandCertificateCertificateDescriptionAuthorityKeyId(c, f.AuthorityKeyId, res); err != nil {
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
	if v, err := expandCertificateCertificateDescriptionCertFingerprint(c, f.CertFingerprint, res); err != nil {
		return nil, fmt.Errorf("error expanding CertFingerprint into certFingerprint: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["certFingerprint"] = v
	}

	return m, nil
}

// flattenCertificateCertificateDescription flattens an instance of CertificateCertificateDescription from a JSON
// response object.
func flattenCertificateCertificateDescription(c *Client, i interface{}, res *Certificate) *CertificateCertificateDescription {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &CertificateCertificateDescription{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyCertificateCertificateDescription
	}
	r.SubjectDescription = flattenCertificateCertificateDescriptionSubjectDescription(c, m["subjectDescription"], res)
	r.X509Description = flattenCertificateCertificateDescriptionX509Description(c, m["x509Description"], res)
	r.PublicKey = flattenCertificateCertificateDescriptionPublicKey(c, m["publicKey"], res)
	r.SubjectKeyId = flattenCertificateCertificateDescriptionSubjectKeyId(c, m["subjectKeyId"], res)
	r.AuthorityKeyId = flattenCertificateCertificateDescriptionAuthorityKeyId(c, m["authorityKeyId"], res)
	r.CrlDistributionPoints = dcl.FlattenStringSlice(m["crlDistributionPoints"])
	r.AiaIssuingCertificateUrls = dcl.FlattenStringSlice(m["aiaIssuingCertificateUrls"])
	r.CertFingerprint = flattenCertificateCertificateDescriptionCertFingerprint(c, m["certFingerprint"], res)

	return r
}

// expandCertificateCertificateDescriptionSubjectDescriptionMap expands the contents of CertificateCertificateDescriptionSubjectDescription into a JSON
// request object.
func expandCertificateCertificateDescriptionSubjectDescriptionMap(c *Client, f map[string]CertificateCertificateDescriptionSubjectDescription, res *Certificate) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandCertificateCertificateDescriptionSubjectDescription(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandCertificateCertificateDescriptionSubjectDescriptionSlice expands the contents of CertificateCertificateDescriptionSubjectDescription into a JSON
// request object.
func expandCertificateCertificateDescriptionSubjectDescriptionSlice(c *Client, f []CertificateCertificateDescriptionSubjectDescription, res *Certificate) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandCertificateCertificateDescriptionSubjectDescription(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenCertificateCertificateDescriptionSubjectDescriptionMap flattens the contents of CertificateCertificateDescriptionSubjectDescription from a JSON
// response object.
func flattenCertificateCertificateDescriptionSubjectDescriptionMap(c *Client, i interface{}, res *Certificate) map[string]CertificateCertificateDescriptionSubjectDescription {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]CertificateCertificateDescriptionSubjectDescription{}
	}

	if len(a) == 0 {
		return map[string]CertificateCertificateDescriptionSubjectDescription{}
	}

	items := make(map[string]CertificateCertificateDescriptionSubjectDescription)
	for k, item := range a {
		items[k] = *flattenCertificateCertificateDescriptionSubjectDescription(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenCertificateCertificateDescriptionSubjectDescriptionSlice flattens the contents of CertificateCertificateDescriptionSubjectDescription from a JSON
// response object.
func flattenCertificateCertificateDescriptionSubjectDescriptionSlice(c *Client, i interface{}, res *Certificate) []CertificateCertificateDescriptionSubjectDescription {
	a, ok := i.([]interface{})
	if !ok {
		return []CertificateCertificateDescriptionSubjectDescription{}
	}

	if len(a) == 0 {
		return []CertificateCertificateDescriptionSubjectDescription{}
	}

	items := make([]CertificateCertificateDescriptionSubjectDescription, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenCertificateCertificateDescriptionSubjectDescription(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandCertificateCertificateDescriptionSubjectDescription expands an instance of CertificateCertificateDescriptionSubjectDescription into a JSON
// request object.
func expandCertificateCertificateDescriptionSubjectDescription(c *Client, f *CertificateCertificateDescriptionSubjectDescription, res *Certificate) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v, err := expandCertificateCertificateDescriptionSubjectDescriptionSubject(c, f.Subject, res); err != nil {
		return nil, fmt.Errorf("error expanding Subject into subject: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["subject"] = v
	}
	if v, err := expandCertificateCertificateDescriptionSubjectDescriptionSubjectAltName(c, f.SubjectAltName, res); err != nil {
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

// flattenCertificateCertificateDescriptionSubjectDescription flattens an instance of CertificateCertificateDescriptionSubjectDescription from a JSON
// response object.
func flattenCertificateCertificateDescriptionSubjectDescription(c *Client, i interface{}, res *Certificate) *CertificateCertificateDescriptionSubjectDescription {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &CertificateCertificateDescriptionSubjectDescription{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyCertificateCertificateDescriptionSubjectDescription
	}
	r.Subject = flattenCertificateCertificateDescriptionSubjectDescriptionSubject(c, m["subject"], res)
	r.SubjectAltName = flattenCertificateCertificateDescriptionSubjectDescriptionSubjectAltName(c, m["subjectAltName"], res)
	r.HexSerialNumber = dcl.FlattenString(m["hexSerialNumber"])
	r.Lifetime = dcl.FlattenString(m["lifetime"])
	r.NotBeforeTime = dcl.FlattenString(m["notBeforeTime"])
	r.NotAfterTime = dcl.FlattenString(m["notAfterTime"])

	return r
}

// expandCertificateCertificateDescriptionSubjectDescriptionSubjectMap expands the contents of CertificateCertificateDescriptionSubjectDescriptionSubject into a JSON
// request object.
func expandCertificateCertificateDescriptionSubjectDescriptionSubjectMap(c *Client, f map[string]CertificateCertificateDescriptionSubjectDescriptionSubject, res *Certificate) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandCertificateCertificateDescriptionSubjectDescriptionSubject(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandCertificateCertificateDescriptionSubjectDescriptionSubjectSlice expands the contents of CertificateCertificateDescriptionSubjectDescriptionSubject into a JSON
// request object.
func expandCertificateCertificateDescriptionSubjectDescriptionSubjectSlice(c *Client, f []CertificateCertificateDescriptionSubjectDescriptionSubject, res *Certificate) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandCertificateCertificateDescriptionSubjectDescriptionSubject(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenCertificateCertificateDescriptionSubjectDescriptionSubjectMap flattens the contents of CertificateCertificateDescriptionSubjectDescriptionSubject from a JSON
// response object.
func flattenCertificateCertificateDescriptionSubjectDescriptionSubjectMap(c *Client, i interface{}, res *Certificate) map[string]CertificateCertificateDescriptionSubjectDescriptionSubject {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]CertificateCertificateDescriptionSubjectDescriptionSubject{}
	}

	if len(a) == 0 {
		return map[string]CertificateCertificateDescriptionSubjectDescriptionSubject{}
	}

	items := make(map[string]CertificateCertificateDescriptionSubjectDescriptionSubject)
	for k, item := range a {
		items[k] = *flattenCertificateCertificateDescriptionSubjectDescriptionSubject(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenCertificateCertificateDescriptionSubjectDescriptionSubjectSlice flattens the contents of CertificateCertificateDescriptionSubjectDescriptionSubject from a JSON
// response object.
func flattenCertificateCertificateDescriptionSubjectDescriptionSubjectSlice(c *Client, i interface{}, res *Certificate) []CertificateCertificateDescriptionSubjectDescriptionSubject {
	a, ok := i.([]interface{})
	if !ok {
		return []CertificateCertificateDescriptionSubjectDescriptionSubject{}
	}

	if len(a) == 0 {
		return []CertificateCertificateDescriptionSubjectDescriptionSubject{}
	}

	items := make([]CertificateCertificateDescriptionSubjectDescriptionSubject, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenCertificateCertificateDescriptionSubjectDescriptionSubject(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandCertificateCertificateDescriptionSubjectDescriptionSubject expands an instance of CertificateCertificateDescriptionSubjectDescriptionSubject into a JSON
// request object.
func expandCertificateCertificateDescriptionSubjectDescriptionSubject(c *Client, f *CertificateCertificateDescriptionSubjectDescriptionSubject, res *Certificate) (map[string]interface{}, error) {
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

// flattenCertificateCertificateDescriptionSubjectDescriptionSubject flattens an instance of CertificateCertificateDescriptionSubjectDescriptionSubject from a JSON
// response object.
func flattenCertificateCertificateDescriptionSubjectDescriptionSubject(c *Client, i interface{}, res *Certificate) *CertificateCertificateDescriptionSubjectDescriptionSubject {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &CertificateCertificateDescriptionSubjectDescriptionSubject{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyCertificateCertificateDescriptionSubjectDescriptionSubject
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

// expandCertificateCertificateDescriptionSubjectDescriptionSubjectAltNameMap expands the contents of CertificateCertificateDescriptionSubjectDescriptionSubjectAltName into a JSON
// request object.
func expandCertificateCertificateDescriptionSubjectDescriptionSubjectAltNameMap(c *Client, f map[string]CertificateCertificateDescriptionSubjectDescriptionSubjectAltName, res *Certificate) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandCertificateCertificateDescriptionSubjectDescriptionSubjectAltName(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandCertificateCertificateDescriptionSubjectDescriptionSubjectAltNameSlice expands the contents of CertificateCertificateDescriptionSubjectDescriptionSubjectAltName into a JSON
// request object.
func expandCertificateCertificateDescriptionSubjectDescriptionSubjectAltNameSlice(c *Client, f []CertificateCertificateDescriptionSubjectDescriptionSubjectAltName, res *Certificate) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandCertificateCertificateDescriptionSubjectDescriptionSubjectAltName(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenCertificateCertificateDescriptionSubjectDescriptionSubjectAltNameMap flattens the contents of CertificateCertificateDescriptionSubjectDescriptionSubjectAltName from a JSON
// response object.
func flattenCertificateCertificateDescriptionSubjectDescriptionSubjectAltNameMap(c *Client, i interface{}, res *Certificate) map[string]CertificateCertificateDescriptionSubjectDescriptionSubjectAltName {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]CertificateCertificateDescriptionSubjectDescriptionSubjectAltName{}
	}

	if len(a) == 0 {
		return map[string]CertificateCertificateDescriptionSubjectDescriptionSubjectAltName{}
	}

	items := make(map[string]CertificateCertificateDescriptionSubjectDescriptionSubjectAltName)
	for k, item := range a {
		items[k] = *flattenCertificateCertificateDescriptionSubjectDescriptionSubjectAltName(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenCertificateCertificateDescriptionSubjectDescriptionSubjectAltNameSlice flattens the contents of CertificateCertificateDescriptionSubjectDescriptionSubjectAltName from a JSON
// response object.
func flattenCertificateCertificateDescriptionSubjectDescriptionSubjectAltNameSlice(c *Client, i interface{}, res *Certificate) []CertificateCertificateDescriptionSubjectDescriptionSubjectAltName {
	a, ok := i.([]interface{})
	if !ok {
		return []CertificateCertificateDescriptionSubjectDescriptionSubjectAltName{}
	}

	if len(a) == 0 {
		return []CertificateCertificateDescriptionSubjectDescriptionSubjectAltName{}
	}

	items := make([]CertificateCertificateDescriptionSubjectDescriptionSubjectAltName, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenCertificateCertificateDescriptionSubjectDescriptionSubjectAltName(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandCertificateCertificateDescriptionSubjectDescriptionSubjectAltName expands an instance of CertificateCertificateDescriptionSubjectDescriptionSubjectAltName into a JSON
// request object.
func expandCertificateCertificateDescriptionSubjectDescriptionSubjectAltName(c *Client, f *CertificateCertificateDescriptionSubjectDescriptionSubjectAltName, res *Certificate) (map[string]interface{}, error) {
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
	if v, err := expandCertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSansSlice(c, f.CustomSans, res); err != nil {
		return nil, fmt.Errorf("error expanding CustomSans into customSans: %w", err)
	} else if v != nil {
		m["customSans"] = v
	}

	return m, nil
}

// flattenCertificateCertificateDescriptionSubjectDescriptionSubjectAltName flattens an instance of CertificateCertificateDescriptionSubjectDescriptionSubjectAltName from a JSON
// response object.
func flattenCertificateCertificateDescriptionSubjectDescriptionSubjectAltName(c *Client, i interface{}, res *Certificate) *CertificateCertificateDescriptionSubjectDescriptionSubjectAltName {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &CertificateCertificateDescriptionSubjectDescriptionSubjectAltName{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyCertificateCertificateDescriptionSubjectDescriptionSubjectAltName
	}
	r.DnsNames = dcl.FlattenStringSlice(m["dnsNames"])
	r.Uris = dcl.FlattenStringSlice(m["uris"])
	r.EmailAddresses = dcl.FlattenStringSlice(m["emailAddresses"])
	r.IPAddresses = dcl.FlattenStringSlice(m["ipAddresses"])
	r.CustomSans = flattenCertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSansSlice(c, m["customSans"], res)

	return r
}

// expandCertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSansMap expands the contents of CertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSans into a JSON
// request object.
func expandCertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSansMap(c *Client, f map[string]CertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSans, res *Certificate) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandCertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSans(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandCertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSansSlice expands the contents of CertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSans into a JSON
// request object.
func expandCertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSansSlice(c *Client, f []CertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSans, res *Certificate) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandCertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSans(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenCertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSansMap flattens the contents of CertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSans from a JSON
// response object.
func flattenCertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSansMap(c *Client, i interface{}, res *Certificate) map[string]CertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSans {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]CertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSans{}
	}

	if len(a) == 0 {
		return map[string]CertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSans{}
	}

	items := make(map[string]CertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSans)
	for k, item := range a {
		items[k] = *flattenCertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSans(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenCertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSansSlice flattens the contents of CertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSans from a JSON
// response object.
func flattenCertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSansSlice(c *Client, i interface{}, res *Certificate) []CertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSans {
	a, ok := i.([]interface{})
	if !ok {
		return []CertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSans{}
	}

	if len(a) == 0 {
		return []CertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSans{}
	}

	items := make([]CertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSans, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenCertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSans(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandCertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSans expands an instance of CertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSans into a JSON
// request object.
func expandCertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSans(c *Client, f *CertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSans, res *Certificate) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v, err := expandCertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSansObjectId(c, f.ObjectId, res); err != nil {
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

// flattenCertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSans flattens an instance of CertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSans from a JSON
// response object.
func flattenCertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSans(c *Client, i interface{}, res *Certificate) *CertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSans {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &CertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSans{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyCertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSans
	}
	r.ObjectId = flattenCertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSansObjectId(c, m["objectId"], res)
	r.Critical = dcl.FlattenBool(m["critical"])
	r.Value = dcl.FlattenString(m["value"])

	return r
}

// expandCertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSansObjectIdMap expands the contents of CertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSansObjectId into a JSON
// request object.
func expandCertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSansObjectIdMap(c *Client, f map[string]CertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSansObjectId, res *Certificate) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandCertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSansObjectId(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandCertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSansObjectIdSlice expands the contents of CertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSansObjectId into a JSON
// request object.
func expandCertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSansObjectIdSlice(c *Client, f []CertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSansObjectId, res *Certificate) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandCertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSansObjectId(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenCertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSansObjectIdMap flattens the contents of CertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSansObjectId from a JSON
// response object.
func flattenCertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSansObjectIdMap(c *Client, i interface{}, res *Certificate) map[string]CertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSansObjectId {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]CertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSansObjectId{}
	}

	if len(a) == 0 {
		return map[string]CertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSansObjectId{}
	}

	items := make(map[string]CertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSansObjectId)
	for k, item := range a {
		items[k] = *flattenCertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSansObjectId(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenCertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSansObjectIdSlice flattens the contents of CertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSansObjectId from a JSON
// response object.
func flattenCertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSansObjectIdSlice(c *Client, i interface{}, res *Certificate) []CertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSansObjectId {
	a, ok := i.([]interface{})
	if !ok {
		return []CertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSansObjectId{}
	}

	if len(a) == 0 {
		return []CertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSansObjectId{}
	}

	items := make([]CertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSansObjectId, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenCertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSansObjectId(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandCertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSansObjectId expands an instance of CertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSansObjectId into a JSON
// request object.
func expandCertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSansObjectId(c *Client, f *CertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSansObjectId, res *Certificate) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.ObjectIdPath; v != nil {
		m["objectIdPath"] = v
	}

	return m, nil
}

// flattenCertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSansObjectId flattens an instance of CertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSansObjectId from a JSON
// response object.
func flattenCertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSansObjectId(c *Client, i interface{}, res *Certificate) *CertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSansObjectId {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &CertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSansObjectId{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyCertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSansObjectId
	}
	r.ObjectIdPath = dcl.FlattenIntSlice(m["objectIdPath"])

	return r
}

// expandCertificateCertificateDescriptionX509DescriptionMap expands the contents of CertificateCertificateDescriptionX509Description into a JSON
// request object.
func expandCertificateCertificateDescriptionX509DescriptionMap(c *Client, f map[string]CertificateCertificateDescriptionX509Description, res *Certificate) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandCertificateCertificateDescriptionX509Description(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandCertificateCertificateDescriptionX509DescriptionSlice expands the contents of CertificateCertificateDescriptionX509Description into a JSON
// request object.
func expandCertificateCertificateDescriptionX509DescriptionSlice(c *Client, f []CertificateCertificateDescriptionX509Description, res *Certificate) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandCertificateCertificateDescriptionX509Description(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenCertificateCertificateDescriptionX509DescriptionMap flattens the contents of CertificateCertificateDescriptionX509Description from a JSON
// response object.
func flattenCertificateCertificateDescriptionX509DescriptionMap(c *Client, i interface{}, res *Certificate) map[string]CertificateCertificateDescriptionX509Description {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]CertificateCertificateDescriptionX509Description{}
	}

	if len(a) == 0 {
		return map[string]CertificateCertificateDescriptionX509Description{}
	}

	items := make(map[string]CertificateCertificateDescriptionX509Description)
	for k, item := range a {
		items[k] = *flattenCertificateCertificateDescriptionX509Description(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenCertificateCertificateDescriptionX509DescriptionSlice flattens the contents of CertificateCertificateDescriptionX509Description from a JSON
// response object.
func flattenCertificateCertificateDescriptionX509DescriptionSlice(c *Client, i interface{}, res *Certificate) []CertificateCertificateDescriptionX509Description {
	a, ok := i.([]interface{})
	if !ok {
		return []CertificateCertificateDescriptionX509Description{}
	}

	if len(a) == 0 {
		return []CertificateCertificateDescriptionX509Description{}
	}

	items := make([]CertificateCertificateDescriptionX509Description, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenCertificateCertificateDescriptionX509Description(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandCertificateCertificateDescriptionX509Description expands an instance of CertificateCertificateDescriptionX509Description into a JSON
// request object.
func expandCertificateCertificateDescriptionX509Description(c *Client, f *CertificateCertificateDescriptionX509Description, res *Certificate) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v, err := expandCertificateCertificateDescriptionX509DescriptionKeyUsage(c, f.KeyUsage, res); err != nil {
		return nil, fmt.Errorf("error expanding KeyUsage into keyUsage: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["keyUsage"] = v
	}
	if v, err := expandCertificateCertificateDescriptionX509DescriptionCaOptions(c, f.CaOptions, res); err != nil {
		return nil, fmt.Errorf("error expanding CaOptions into caOptions: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["caOptions"] = v
	}
	if v, err := expandCertificateCertificateDescriptionX509DescriptionPolicyIdsSlice(c, f.PolicyIds, res); err != nil {
		return nil, fmt.Errorf("error expanding PolicyIds into policyIds: %w", err)
	} else if v != nil {
		m["policyIds"] = v
	}
	if v := f.AiaOcspServers; v != nil {
		m["aiaOcspServers"] = v
	}
	if v, err := expandCertificateCertificateDescriptionX509DescriptionAdditionalExtensionsSlice(c, f.AdditionalExtensions, res); err != nil {
		return nil, fmt.Errorf("error expanding AdditionalExtensions into additionalExtensions: %w", err)
	} else if v != nil {
		m["additionalExtensions"] = v
	}

	return m, nil
}

// flattenCertificateCertificateDescriptionX509Description flattens an instance of CertificateCertificateDescriptionX509Description from a JSON
// response object.
func flattenCertificateCertificateDescriptionX509Description(c *Client, i interface{}, res *Certificate) *CertificateCertificateDescriptionX509Description {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &CertificateCertificateDescriptionX509Description{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyCertificateCertificateDescriptionX509Description
	}
	r.KeyUsage = flattenCertificateCertificateDescriptionX509DescriptionKeyUsage(c, m["keyUsage"], res)
	r.CaOptions = flattenCertificateCertificateDescriptionX509DescriptionCaOptions(c, m["caOptions"], res)
	r.PolicyIds = flattenCertificateCertificateDescriptionX509DescriptionPolicyIdsSlice(c, m["policyIds"], res)
	r.AiaOcspServers = dcl.FlattenStringSlice(m["aiaOcspServers"])
	r.AdditionalExtensions = flattenCertificateCertificateDescriptionX509DescriptionAdditionalExtensionsSlice(c, m["additionalExtensions"], res)

	return r
}

// expandCertificateCertificateDescriptionX509DescriptionKeyUsageMap expands the contents of CertificateCertificateDescriptionX509DescriptionKeyUsage into a JSON
// request object.
func expandCertificateCertificateDescriptionX509DescriptionKeyUsageMap(c *Client, f map[string]CertificateCertificateDescriptionX509DescriptionKeyUsage, res *Certificate) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandCertificateCertificateDescriptionX509DescriptionKeyUsage(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandCertificateCertificateDescriptionX509DescriptionKeyUsageSlice expands the contents of CertificateCertificateDescriptionX509DescriptionKeyUsage into a JSON
// request object.
func expandCertificateCertificateDescriptionX509DescriptionKeyUsageSlice(c *Client, f []CertificateCertificateDescriptionX509DescriptionKeyUsage, res *Certificate) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandCertificateCertificateDescriptionX509DescriptionKeyUsage(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenCertificateCertificateDescriptionX509DescriptionKeyUsageMap flattens the contents of CertificateCertificateDescriptionX509DescriptionKeyUsage from a JSON
// response object.
func flattenCertificateCertificateDescriptionX509DescriptionKeyUsageMap(c *Client, i interface{}, res *Certificate) map[string]CertificateCertificateDescriptionX509DescriptionKeyUsage {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]CertificateCertificateDescriptionX509DescriptionKeyUsage{}
	}

	if len(a) == 0 {
		return map[string]CertificateCertificateDescriptionX509DescriptionKeyUsage{}
	}

	items := make(map[string]CertificateCertificateDescriptionX509DescriptionKeyUsage)
	for k, item := range a {
		items[k] = *flattenCertificateCertificateDescriptionX509DescriptionKeyUsage(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenCertificateCertificateDescriptionX509DescriptionKeyUsageSlice flattens the contents of CertificateCertificateDescriptionX509DescriptionKeyUsage from a JSON
// response object.
func flattenCertificateCertificateDescriptionX509DescriptionKeyUsageSlice(c *Client, i interface{}, res *Certificate) []CertificateCertificateDescriptionX509DescriptionKeyUsage {
	a, ok := i.([]interface{})
	if !ok {
		return []CertificateCertificateDescriptionX509DescriptionKeyUsage{}
	}

	if len(a) == 0 {
		return []CertificateCertificateDescriptionX509DescriptionKeyUsage{}
	}

	items := make([]CertificateCertificateDescriptionX509DescriptionKeyUsage, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenCertificateCertificateDescriptionX509DescriptionKeyUsage(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandCertificateCertificateDescriptionX509DescriptionKeyUsage expands an instance of CertificateCertificateDescriptionX509DescriptionKeyUsage into a JSON
// request object.
func expandCertificateCertificateDescriptionX509DescriptionKeyUsage(c *Client, f *CertificateCertificateDescriptionX509DescriptionKeyUsage, res *Certificate) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v, err := expandCertificateCertificateDescriptionX509DescriptionKeyUsageBaseKeyUsage(c, f.BaseKeyUsage, res); err != nil {
		return nil, fmt.Errorf("error expanding BaseKeyUsage into baseKeyUsage: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["baseKeyUsage"] = v
	}
	if v, err := expandCertificateCertificateDescriptionX509DescriptionKeyUsageExtendedKeyUsage(c, f.ExtendedKeyUsage, res); err != nil {
		return nil, fmt.Errorf("error expanding ExtendedKeyUsage into extendedKeyUsage: %w", err)
	} else if !dcl.IsEmptyValueIndirect(v) {
		m["extendedKeyUsage"] = v
	}
	if v, err := expandCertificateCertificateDescriptionX509DescriptionKeyUsageUnknownExtendedKeyUsagesSlice(c, f.UnknownExtendedKeyUsages, res); err != nil {
		return nil, fmt.Errorf("error expanding UnknownExtendedKeyUsages into unknownExtendedKeyUsages: %w", err)
	} else if v != nil {
		m["unknownExtendedKeyUsages"] = v
	}

	return m, nil
}

// flattenCertificateCertificateDescriptionX509DescriptionKeyUsage flattens an instance of CertificateCertificateDescriptionX509DescriptionKeyUsage from a JSON
// response object.
func flattenCertificateCertificateDescriptionX509DescriptionKeyUsage(c *Client, i interface{}, res *Certificate) *CertificateCertificateDescriptionX509DescriptionKeyUsage {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &CertificateCertificateDescriptionX509DescriptionKeyUsage{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyCertificateCertificateDescriptionX509DescriptionKeyUsage
	}
	r.BaseKeyUsage = flattenCertificateCertificateDescriptionX509DescriptionKeyUsageBaseKeyUsage(c, m["baseKeyUsage"], res)
	r.ExtendedKeyUsage = flattenCertificateCertificateDescriptionX509DescriptionKeyUsageExtendedKeyUsage(c, m["extendedKeyUsage"], res)
	r.UnknownExtendedKeyUsages = flattenCertificateCertificateDescriptionX509DescriptionKeyUsageUnknownExtendedKeyUsagesSlice(c, m["unknownExtendedKeyUsages"], res)

	return r
}

// expandCertificateCertificateDescriptionX509DescriptionKeyUsageBaseKeyUsageMap expands the contents of CertificateCertificateDescriptionX509DescriptionKeyUsageBaseKeyUsage into a JSON
// request object.
func expandCertificateCertificateDescriptionX509DescriptionKeyUsageBaseKeyUsageMap(c *Client, f map[string]CertificateCertificateDescriptionX509DescriptionKeyUsageBaseKeyUsage, res *Certificate) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandCertificateCertificateDescriptionX509DescriptionKeyUsageBaseKeyUsage(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandCertificateCertificateDescriptionX509DescriptionKeyUsageBaseKeyUsageSlice expands the contents of CertificateCertificateDescriptionX509DescriptionKeyUsageBaseKeyUsage into a JSON
// request object.
func expandCertificateCertificateDescriptionX509DescriptionKeyUsageBaseKeyUsageSlice(c *Client, f []CertificateCertificateDescriptionX509DescriptionKeyUsageBaseKeyUsage, res *Certificate) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandCertificateCertificateDescriptionX509DescriptionKeyUsageBaseKeyUsage(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenCertificateCertificateDescriptionX509DescriptionKeyUsageBaseKeyUsageMap flattens the contents of CertificateCertificateDescriptionX509DescriptionKeyUsageBaseKeyUsage from a JSON
// response object.
func flattenCertificateCertificateDescriptionX509DescriptionKeyUsageBaseKeyUsageMap(c *Client, i interface{}, res *Certificate) map[string]CertificateCertificateDescriptionX509DescriptionKeyUsageBaseKeyUsage {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]CertificateCertificateDescriptionX509DescriptionKeyUsageBaseKeyUsage{}
	}

	if len(a) == 0 {
		return map[string]CertificateCertificateDescriptionX509DescriptionKeyUsageBaseKeyUsage{}
	}

	items := make(map[string]CertificateCertificateDescriptionX509DescriptionKeyUsageBaseKeyUsage)
	for k, item := range a {
		items[k] = *flattenCertificateCertificateDescriptionX509DescriptionKeyUsageBaseKeyUsage(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenCertificateCertificateDescriptionX509DescriptionKeyUsageBaseKeyUsageSlice flattens the contents of CertificateCertificateDescriptionX509DescriptionKeyUsageBaseKeyUsage from a JSON
// response object.
func flattenCertificateCertificateDescriptionX509DescriptionKeyUsageBaseKeyUsageSlice(c *Client, i interface{}, res *Certificate) []CertificateCertificateDescriptionX509DescriptionKeyUsageBaseKeyUsage {
	a, ok := i.([]interface{})
	if !ok {
		return []CertificateCertificateDescriptionX509DescriptionKeyUsageBaseKeyUsage{}
	}

	if len(a) == 0 {
		return []CertificateCertificateDescriptionX509DescriptionKeyUsageBaseKeyUsage{}
	}

	items := make([]CertificateCertificateDescriptionX509DescriptionKeyUsageBaseKeyUsage, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenCertificateCertificateDescriptionX509DescriptionKeyUsageBaseKeyUsage(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandCertificateCertificateDescriptionX509DescriptionKeyUsageBaseKeyUsage expands an instance of CertificateCertificateDescriptionX509DescriptionKeyUsageBaseKeyUsage into a JSON
// request object.
func expandCertificateCertificateDescriptionX509DescriptionKeyUsageBaseKeyUsage(c *Client, f *CertificateCertificateDescriptionX509DescriptionKeyUsageBaseKeyUsage, res *Certificate) (map[string]interface{}, error) {
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

// flattenCertificateCertificateDescriptionX509DescriptionKeyUsageBaseKeyUsage flattens an instance of CertificateCertificateDescriptionX509DescriptionKeyUsageBaseKeyUsage from a JSON
// response object.
func flattenCertificateCertificateDescriptionX509DescriptionKeyUsageBaseKeyUsage(c *Client, i interface{}, res *Certificate) *CertificateCertificateDescriptionX509DescriptionKeyUsageBaseKeyUsage {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &CertificateCertificateDescriptionX509DescriptionKeyUsageBaseKeyUsage{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyCertificateCertificateDescriptionX509DescriptionKeyUsageBaseKeyUsage
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

// expandCertificateCertificateDescriptionX509DescriptionKeyUsageExtendedKeyUsageMap expands the contents of CertificateCertificateDescriptionX509DescriptionKeyUsageExtendedKeyUsage into a JSON
// request object.
func expandCertificateCertificateDescriptionX509DescriptionKeyUsageExtendedKeyUsageMap(c *Client, f map[string]CertificateCertificateDescriptionX509DescriptionKeyUsageExtendedKeyUsage, res *Certificate) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandCertificateCertificateDescriptionX509DescriptionKeyUsageExtendedKeyUsage(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandCertificateCertificateDescriptionX509DescriptionKeyUsageExtendedKeyUsageSlice expands the contents of CertificateCertificateDescriptionX509DescriptionKeyUsageExtendedKeyUsage into a JSON
// request object.
func expandCertificateCertificateDescriptionX509DescriptionKeyUsageExtendedKeyUsageSlice(c *Client, f []CertificateCertificateDescriptionX509DescriptionKeyUsageExtendedKeyUsage, res *Certificate) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandCertificateCertificateDescriptionX509DescriptionKeyUsageExtendedKeyUsage(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenCertificateCertificateDescriptionX509DescriptionKeyUsageExtendedKeyUsageMap flattens the contents of CertificateCertificateDescriptionX509DescriptionKeyUsageExtendedKeyUsage from a JSON
// response object.
func flattenCertificateCertificateDescriptionX509DescriptionKeyUsageExtendedKeyUsageMap(c *Client, i interface{}, res *Certificate) map[string]CertificateCertificateDescriptionX509DescriptionKeyUsageExtendedKeyUsage {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]CertificateCertificateDescriptionX509DescriptionKeyUsageExtendedKeyUsage{}
	}

	if len(a) == 0 {
		return map[string]CertificateCertificateDescriptionX509DescriptionKeyUsageExtendedKeyUsage{}
	}

	items := make(map[string]CertificateCertificateDescriptionX509DescriptionKeyUsageExtendedKeyUsage)
	for k, item := range a {
		items[k] = *flattenCertificateCertificateDescriptionX509DescriptionKeyUsageExtendedKeyUsage(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenCertificateCertificateDescriptionX509DescriptionKeyUsageExtendedKeyUsageSlice flattens the contents of CertificateCertificateDescriptionX509DescriptionKeyUsageExtendedKeyUsage from a JSON
// response object.
func flattenCertificateCertificateDescriptionX509DescriptionKeyUsageExtendedKeyUsageSlice(c *Client, i interface{}, res *Certificate) []CertificateCertificateDescriptionX509DescriptionKeyUsageExtendedKeyUsage {
	a, ok := i.([]interface{})
	if !ok {
		return []CertificateCertificateDescriptionX509DescriptionKeyUsageExtendedKeyUsage{}
	}

	if len(a) == 0 {
		return []CertificateCertificateDescriptionX509DescriptionKeyUsageExtendedKeyUsage{}
	}

	items := make([]CertificateCertificateDescriptionX509DescriptionKeyUsageExtendedKeyUsage, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenCertificateCertificateDescriptionX509DescriptionKeyUsageExtendedKeyUsage(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandCertificateCertificateDescriptionX509DescriptionKeyUsageExtendedKeyUsage expands an instance of CertificateCertificateDescriptionX509DescriptionKeyUsageExtendedKeyUsage into a JSON
// request object.
func expandCertificateCertificateDescriptionX509DescriptionKeyUsageExtendedKeyUsage(c *Client, f *CertificateCertificateDescriptionX509DescriptionKeyUsageExtendedKeyUsage, res *Certificate) (map[string]interface{}, error) {
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

// flattenCertificateCertificateDescriptionX509DescriptionKeyUsageExtendedKeyUsage flattens an instance of CertificateCertificateDescriptionX509DescriptionKeyUsageExtendedKeyUsage from a JSON
// response object.
func flattenCertificateCertificateDescriptionX509DescriptionKeyUsageExtendedKeyUsage(c *Client, i interface{}, res *Certificate) *CertificateCertificateDescriptionX509DescriptionKeyUsageExtendedKeyUsage {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &CertificateCertificateDescriptionX509DescriptionKeyUsageExtendedKeyUsage{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyCertificateCertificateDescriptionX509DescriptionKeyUsageExtendedKeyUsage
	}
	r.ServerAuth = dcl.FlattenBool(m["serverAuth"])
	r.ClientAuth = dcl.FlattenBool(m["clientAuth"])
	r.CodeSigning = dcl.FlattenBool(m["codeSigning"])
	r.EmailProtection = dcl.FlattenBool(m["emailProtection"])
	r.TimeStamping = dcl.FlattenBool(m["timeStamping"])
	r.OcspSigning = dcl.FlattenBool(m["ocspSigning"])

	return r
}

// expandCertificateCertificateDescriptionX509DescriptionKeyUsageUnknownExtendedKeyUsagesMap expands the contents of CertificateCertificateDescriptionX509DescriptionKeyUsageUnknownExtendedKeyUsages into a JSON
// request object.
func expandCertificateCertificateDescriptionX509DescriptionKeyUsageUnknownExtendedKeyUsagesMap(c *Client, f map[string]CertificateCertificateDescriptionX509DescriptionKeyUsageUnknownExtendedKeyUsages, res *Certificate) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandCertificateCertificateDescriptionX509DescriptionKeyUsageUnknownExtendedKeyUsages(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandCertificateCertificateDescriptionX509DescriptionKeyUsageUnknownExtendedKeyUsagesSlice expands the contents of CertificateCertificateDescriptionX509DescriptionKeyUsageUnknownExtendedKeyUsages into a JSON
// request object.
func expandCertificateCertificateDescriptionX509DescriptionKeyUsageUnknownExtendedKeyUsagesSlice(c *Client, f []CertificateCertificateDescriptionX509DescriptionKeyUsageUnknownExtendedKeyUsages, res *Certificate) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandCertificateCertificateDescriptionX509DescriptionKeyUsageUnknownExtendedKeyUsages(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenCertificateCertificateDescriptionX509DescriptionKeyUsageUnknownExtendedKeyUsagesMap flattens the contents of CertificateCertificateDescriptionX509DescriptionKeyUsageUnknownExtendedKeyUsages from a JSON
// response object.
func flattenCertificateCertificateDescriptionX509DescriptionKeyUsageUnknownExtendedKeyUsagesMap(c *Client, i interface{}, res *Certificate) map[string]CertificateCertificateDescriptionX509DescriptionKeyUsageUnknownExtendedKeyUsages {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]CertificateCertificateDescriptionX509DescriptionKeyUsageUnknownExtendedKeyUsages{}
	}

	if len(a) == 0 {
		return map[string]CertificateCertificateDescriptionX509DescriptionKeyUsageUnknownExtendedKeyUsages{}
	}

	items := make(map[string]CertificateCertificateDescriptionX509DescriptionKeyUsageUnknownExtendedKeyUsages)
	for k, item := range a {
		items[k] = *flattenCertificateCertificateDescriptionX509DescriptionKeyUsageUnknownExtendedKeyUsages(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenCertificateCertificateDescriptionX509DescriptionKeyUsageUnknownExtendedKeyUsagesSlice flattens the contents of CertificateCertificateDescriptionX509DescriptionKeyUsageUnknownExtendedKeyUsages from a JSON
// response object.
func flattenCertificateCertificateDescriptionX509DescriptionKeyUsageUnknownExtendedKeyUsagesSlice(c *Client, i interface{}, res *Certificate) []CertificateCertificateDescriptionX509DescriptionKeyUsageUnknownExtendedKeyUsages {
	a, ok := i.([]interface{})
	if !ok {
		return []CertificateCertificateDescriptionX509DescriptionKeyUsageUnknownExtendedKeyUsages{}
	}

	if len(a) == 0 {
		return []CertificateCertificateDescriptionX509DescriptionKeyUsageUnknownExtendedKeyUsages{}
	}

	items := make([]CertificateCertificateDescriptionX509DescriptionKeyUsageUnknownExtendedKeyUsages, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenCertificateCertificateDescriptionX509DescriptionKeyUsageUnknownExtendedKeyUsages(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandCertificateCertificateDescriptionX509DescriptionKeyUsageUnknownExtendedKeyUsages expands an instance of CertificateCertificateDescriptionX509DescriptionKeyUsageUnknownExtendedKeyUsages into a JSON
// request object.
func expandCertificateCertificateDescriptionX509DescriptionKeyUsageUnknownExtendedKeyUsages(c *Client, f *CertificateCertificateDescriptionX509DescriptionKeyUsageUnknownExtendedKeyUsages, res *Certificate) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.ObjectIdPath; v != nil {
		m["objectIdPath"] = v
	}

	return m, nil
}

// flattenCertificateCertificateDescriptionX509DescriptionKeyUsageUnknownExtendedKeyUsages flattens an instance of CertificateCertificateDescriptionX509DescriptionKeyUsageUnknownExtendedKeyUsages from a JSON
// response object.
func flattenCertificateCertificateDescriptionX509DescriptionKeyUsageUnknownExtendedKeyUsages(c *Client, i interface{}, res *Certificate) *CertificateCertificateDescriptionX509DescriptionKeyUsageUnknownExtendedKeyUsages {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &CertificateCertificateDescriptionX509DescriptionKeyUsageUnknownExtendedKeyUsages{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyCertificateCertificateDescriptionX509DescriptionKeyUsageUnknownExtendedKeyUsages
	}
	r.ObjectIdPath = dcl.FlattenIntSlice(m["objectIdPath"])

	return r
}

// expandCertificateCertificateDescriptionX509DescriptionCaOptionsMap expands the contents of CertificateCertificateDescriptionX509DescriptionCaOptions into a JSON
// request object.
func expandCertificateCertificateDescriptionX509DescriptionCaOptionsMap(c *Client, f map[string]CertificateCertificateDescriptionX509DescriptionCaOptions, res *Certificate) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandCertificateCertificateDescriptionX509DescriptionCaOptions(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandCertificateCertificateDescriptionX509DescriptionCaOptionsSlice expands the contents of CertificateCertificateDescriptionX509DescriptionCaOptions into a JSON
// request object.
func expandCertificateCertificateDescriptionX509DescriptionCaOptionsSlice(c *Client, f []CertificateCertificateDescriptionX509DescriptionCaOptions, res *Certificate) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandCertificateCertificateDescriptionX509DescriptionCaOptions(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenCertificateCertificateDescriptionX509DescriptionCaOptionsMap flattens the contents of CertificateCertificateDescriptionX509DescriptionCaOptions from a JSON
// response object.
func flattenCertificateCertificateDescriptionX509DescriptionCaOptionsMap(c *Client, i interface{}, res *Certificate) map[string]CertificateCertificateDescriptionX509DescriptionCaOptions {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]CertificateCertificateDescriptionX509DescriptionCaOptions{}
	}

	if len(a) == 0 {
		return map[string]CertificateCertificateDescriptionX509DescriptionCaOptions{}
	}

	items := make(map[string]CertificateCertificateDescriptionX509DescriptionCaOptions)
	for k, item := range a {
		items[k] = *flattenCertificateCertificateDescriptionX509DescriptionCaOptions(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenCertificateCertificateDescriptionX509DescriptionCaOptionsSlice flattens the contents of CertificateCertificateDescriptionX509DescriptionCaOptions from a JSON
// response object.
func flattenCertificateCertificateDescriptionX509DescriptionCaOptionsSlice(c *Client, i interface{}, res *Certificate) []CertificateCertificateDescriptionX509DescriptionCaOptions {
	a, ok := i.([]interface{})
	if !ok {
		return []CertificateCertificateDescriptionX509DescriptionCaOptions{}
	}

	if len(a) == 0 {
		return []CertificateCertificateDescriptionX509DescriptionCaOptions{}
	}

	items := make([]CertificateCertificateDescriptionX509DescriptionCaOptions, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenCertificateCertificateDescriptionX509DescriptionCaOptions(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandCertificateCertificateDescriptionX509DescriptionCaOptions expands an instance of CertificateCertificateDescriptionX509DescriptionCaOptions into a JSON
// request object.
func expandCertificateCertificateDescriptionX509DescriptionCaOptions(c *Client, f *CertificateCertificateDescriptionX509DescriptionCaOptions, res *Certificate) (map[string]interface{}, error) {
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

// flattenCertificateCertificateDescriptionX509DescriptionCaOptions flattens an instance of CertificateCertificateDescriptionX509DescriptionCaOptions from a JSON
// response object.
func flattenCertificateCertificateDescriptionX509DescriptionCaOptions(c *Client, i interface{}, res *Certificate) *CertificateCertificateDescriptionX509DescriptionCaOptions {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &CertificateCertificateDescriptionX509DescriptionCaOptions{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyCertificateCertificateDescriptionX509DescriptionCaOptions
	}
	r.IsCa = dcl.FlattenBool(m["isCa"])
	r.MaxIssuerPathLength = dcl.FlattenInteger(m["maxIssuerPathLength"])

	return r
}

// expandCertificateCertificateDescriptionX509DescriptionPolicyIdsMap expands the contents of CertificateCertificateDescriptionX509DescriptionPolicyIds into a JSON
// request object.
func expandCertificateCertificateDescriptionX509DescriptionPolicyIdsMap(c *Client, f map[string]CertificateCertificateDescriptionX509DescriptionPolicyIds, res *Certificate) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandCertificateCertificateDescriptionX509DescriptionPolicyIds(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandCertificateCertificateDescriptionX509DescriptionPolicyIdsSlice expands the contents of CertificateCertificateDescriptionX509DescriptionPolicyIds into a JSON
// request object.
func expandCertificateCertificateDescriptionX509DescriptionPolicyIdsSlice(c *Client, f []CertificateCertificateDescriptionX509DescriptionPolicyIds, res *Certificate) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandCertificateCertificateDescriptionX509DescriptionPolicyIds(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenCertificateCertificateDescriptionX509DescriptionPolicyIdsMap flattens the contents of CertificateCertificateDescriptionX509DescriptionPolicyIds from a JSON
// response object.
func flattenCertificateCertificateDescriptionX509DescriptionPolicyIdsMap(c *Client, i interface{}, res *Certificate) map[string]CertificateCertificateDescriptionX509DescriptionPolicyIds {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]CertificateCertificateDescriptionX509DescriptionPolicyIds{}
	}

	if len(a) == 0 {
		return map[string]CertificateCertificateDescriptionX509DescriptionPolicyIds{}
	}

	items := make(map[string]CertificateCertificateDescriptionX509DescriptionPolicyIds)
	for k, item := range a {
		items[k] = *flattenCertificateCertificateDescriptionX509DescriptionPolicyIds(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenCertificateCertificateDescriptionX509DescriptionPolicyIdsSlice flattens the contents of CertificateCertificateDescriptionX509DescriptionPolicyIds from a JSON
// response object.
func flattenCertificateCertificateDescriptionX509DescriptionPolicyIdsSlice(c *Client, i interface{}, res *Certificate) []CertificateCertificateDescriptionX509DescriptionPolicyIds {
	a, ok := i.([]interface{})
	if !ok {
		return []CertificateCertificateDescriptionX509DescriptionPolicyIds{}
	}

	if len(a) == 0 {
		return []CertificateCertificateDescriptionX509DescriptionPolicyIds{}
	}

	items := make([]CertificateCertificateDescriptionX509DescriptionPolicyIds, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenCertificateCertificateDescriptionX509DescriptionPolicyIds(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandCertificateCertificateDescriptionX509DescriptionPolicyIds expands an instance of CertificateCertificateDescriptionX509DescriptionPolicyIds into a JSON
// request object.
func expandCertificateCertificateDescriptionX509DescriptionPolicyIds(c *Client, f *CertificateCertificateDescriptionX509DescriptionPolicyIds, res *Certificate) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.ObjectIdPath; v != nil {
		m["objectIdPath"] = v
	}

	return m, nil
}

// flattenCertificateCertificateDescriptionX509DescriptionPolicyIds flattens an instance of CertificateCertificateDescriptionX509DescriptionPolicyIds from a JSON
// response object.
func flattenCertificateCertificateDescriptionX509DescriptionPolicyIds(c *Client, i interface{}, res *Certificate) *CertificateCertificateDescriptionX509DescriptionPolicyIds {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &CertificateCertificateDescriptionX509DescriptionPolicyIds{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyCertificateCertificateDescriptionX509DescriptionPolicyIds
	}
	r.ObjectIdPath = dcl.FlattenIntSlice(m["objectIdPath"])

	return r
}

// expandCertificateCertificateDescriptionX509DescriptionAdditionalExtensionsMap expands the contents of CertificateCertificateDescriptionX509DescriptionAdditionalExtensions into a JSON
// request object.
func expandCertificateCertificateDescriptionX509DescriptionAdditionalExtensionsMap(c *Client, f map[string]CertificateCertificateDescriptionX509DescriptionAdditionalExtensions, res *Certificate) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandCertificateCertificateDescriptionX509DescriptionAdditionalExtensions(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandCertificateCertificateDescriptionX509DescriptionAdditionalExtensionsSlice expands the contents of CertificateCertificateDescriptionX509DescriptionAdditionalExtensions into a JSON
// request object.
func expandCertificateCertificateDescriptionX509DescriptionAdditionalExtensionsSlice(c *Client, f []CertificateCertificateDescriptionX509DescriptionAdditionalExtensions, res *Certificate) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandCertificateCertificateDescriptionX509DescriptionAdditionalExtensions(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenCertificateCertificateDescriptionX509DescriptionAdditionalExtensionsMap flattens the contents of CertificateCertificateDescriptionX509DescriptionAdditionalExtensions from a JSON
// response object.
func flattenCertificateCertificateDescriptionX509DescriptionAdditionalExtensionsMap(c *Client, i interface{}, res *Certificate) map[string]CertificateCertificateDescriptionX509DescriptionAdditionalExtensions {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]CertificateCertificateDescriptionX509DescriptionAdditionalExtensions{}
	}

	if len(a) == 0 {
		return map[string]CertificateCertificateDescriptionX509DescriptionAdditionalExtensions{}
	}

	items := make(map[string]CertificateCertificateDescriptionX509DescriptionAdditionalExtensions)
	for k, item := range a {
		items[k] = *flattenCertificateCertificateDescriptionX509DescriptionAdditionalExtensions(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenCertificateCertificateDescriptionX509DescriptionAdditionalExtensionsSlice flattens the contents of CertificateCertificateDescriptionX509DescriptionAdditionalExtensions from a JSON
// response object.
func flattenCertificateCertificateDescriptionX509DescriptionAdditionalExtensionsSlice(c *Client, i interface{}, res *Certificate) []CertificateCertificateDescriptionX509DescriptionAdditionalExtensions {
	a, ok := i.([]interface{})
	if !ok {
		return []CertificateCertificateDescriptionX509DescriptionAdditionalExtensions{}
	}

	if len(a) == 0 {
		return []CertificateCertificateDescriptionX509DescriptionAdditionalExtensions{}
	}

	items := make([]CertificateCertificateDescriptionX509DescriptionAdditionalExtensions, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenCertificateCertificateDescriptionX509DescriptionAdditionalExtensions(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandCertificateCertificateDescriptionX509DescriptionAdditionalExtensions expands an instance of CertificateCertificateDescriptionX509DescriptionAdditionalExtensions into a JSON
// request object.
func expandCertificateCertificateDescriptionX509DescriptionAdditionalExtensions(c *Client, f *CertificateCertificateDescriptionX509DescriptionAdditionalExtensions, res *Certificate) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v, err := expandCertificateCertificateDescriptionX509DescriptionAdditionalExtensionsObjectId(c, f.ObjectId, res); err != nil {
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

// flattenCertificateCertificateDescriptionX509DescriptionAdditionalExtensions flattens an instance of CertificateCertificateDescriptionX509DescriptionAdditionalExtensions from a JSON
// response object.
func flattenCertificateCertificateDescriptionX509DescriptionAdditionalExtensions(c *Client, i interface{}, res *Certificate) *CertificateCertificateDescriptionX509DescriptionAdditionalExtensions {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &CertificateCertificateDescriptionX509DescriptionAdditionalExtensions{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyCertificateCertificateDescriptionX509DescriptionAdditionalExtensions
	}
	r.ObjectId = flattenCertificateCertificateDescriptionX509DescriptionAdditionalExtensionsObjectId(c, m["objectId"], res)
	r.Critical = dcl.FlattenBool(m["critical"])
	r.Value = dcl.FlattenString(m["value"])

	return r
}

// expandCertificateCertificateDescriptionX509DescriptionAdditionalExtensionsObjectIdMap expands the contents of CertificateCertificateDescriptionX509DescriptionAdditionalExtensionsObjectId into a JSON
// request object.
func expandCertificateCertificateDescriptionX509DescriptionAdditionalExtensionsObjectIdMap(c *Client, f map[string]CertificateCertificateDescriptionX509DescriptionAdditionalExtensionsObjectId, res *Certificate) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandCertificateCertificateDescriptionX509DescriptionAdditionalExtensionsObjectId(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandCertificateCertificateDescriptionX509DescriptionAdditionalExtensionsObjectIdSlice expands the contents of CertificateCertificateDescriptionX509DescriptionAdditionalExtensionsObjectId into a JSON
// request object.
func expandCertificateCertificateDescriptionX509DescriptionAdditionalExtensionsObjectIdSlice(c *Client, f []CertificateCertificateDescriptionX509DescriptionAdditionalExtensionsObjectId, res *Certificate) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandCertificateCertificateDescriptionX509DescriptionAdditionalExtensionsObjectId(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenCertificateCertificateDescriptionX509DescriptionAdditionalExtensionsObjectIdMap flattens the contents of CertificateCertificateDescriptionX509DescriptionAdditionalExtensionsObjectId from a JSON
// response object.
func flattenCertificateCertificateDescriptionX509DescriptionAdditionalExtensionsObjectIdMap(c *Client, i interface{}, res *Certificate) map[string]CertificateCertificateDescriptionX509DescriptionAdditionalExtensionsObjectId {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]CertificateCertificateDescriptionX509DescriptionAdditionalExtensionsObjectId{}
	}

	if len(a) == 0 {
		return map[string]CertificateCertificateDescriptionX509DescriptionAdditionalExtensionsObjectId{}
	}

	items := make(map[string]CertificateCertificateDescriptionX509DescriptionAdditionalExtensionsObjectId)
	for k, item := range a {
		items[k] = *flattenCertificateCertificateDescriptionX509DescriptionAdditionalExtensionsObjectId(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenCertificateCertificateDescriptionX509DescriptionAdditionalExtensionsObjectIdSlice flattens the contents of CertificateCertificateDescriptionX509DescriptionAdditionalExtensionsObjectId from a JSON
// response object.
func flattenCertificateCertificateDescriptionX509DescriptionAdditionalExtensionsObjectIdSlice(c *Client, i interface{}, res *Certificate) []CertificateCertificateDescriptionX509DescriptionAdditionalExtensionsObjectId {
	a, ok := i.([]interface{})
	if !ok {
		return []CertificateCertificateDescriptionX509DescriptionAdditionalExtensionsObjectId{}
	}

	if len(a) == 0 {
		return []CertificateCertificateDescriptionX509DescriptionAdditionalExtensionsObjectId{}
	}

	items := make([]CertificateCertificateDescriptionX509DescriptionAdditionalExtensionsObjectId, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenCertificateCertificateDescriptionX509DescriptionAdditionalExtensionsObjectId(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandCertificateCertificateDescriptionX509DescriptionAdditionalExtensionsObjectId expands an instance of CertificateCertificateDescriptionX509DescriptionAdditionalExtensionsObjectId into a JSON
// request object.
func expandCertificateCertificateDescriptionX509DescriptionAdditionalExtensionsObjectId(c *Client, f *CertificateCertificateDescriptionX509DescriptionAdditionalExtensionsObjectId, res *Certificate) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.ObjectIdPath; v != nil {
		m["objectIdPath"] = v
	}

	return m, nil
}

// flattenCertificateCertificateDescriptionX509DescriptionAdditionalExtensionsObjectId flattens an instance of CertificateCertificateDescriptionX509DescriptionAdditionalExtensionsObjectId from a JSON
// response object.
func flattenCertificateCertificateDescriptionX509DescriptionAdditionalExtensionsObjectId(c *Client, i interface{}, res *Certificate) *CertificateCertificateDescriptionX509DescriptionAdditionalExtensionsObjectId {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &CertificateCertificateDescriptionX509DescriptionAdditionalExtensionsObjectId{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyCertificateCertificateDescriptionX509DescriptionAdditionalExtensionsObjectId
	}
	r.ObjectIdPath = dcl.FlattenIntSlice(m["objectIdPath"])

	return r
}

// expandCertificateCertificateDescriptionPublicKeyMap expands the contents of CertificateCertificateDescriptionPublicKey into a JSON
// request object.
func expandCertificateCertificateDescriptionPublicKeyMap(c *Client, f map[string]CertificateCertificateDescriptionPublicKey, res *Certificate) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandCertificateCertificateDescriptionPublicKey(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandCertificateCertificateDescriptionPublicKeySlice expands the contents of CertificateCertificateDescriptionPublicKey into a JSON
// request object.
func expandCertificateCertificateDescriptionPublicKeySlice(c *Client, f []CertificateCertificateDescriptionPublicKey, res *Certificate) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandCertificateCertificateDescriptionPublicKey(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenCertificateCertificateDescriptionPublicKeyMap flattens the contents of CertificateCertificateDescriptionPublicKey from a JSON
// response object.
func flattenCertificateCertificateDescriptionPublicKeyMap(c *Client, i interface{}, res *Certificate) map[string]CertificateCertificateDescriptionPublicKey {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]CertificateCertificateDescriptionPublicKey{}
	}

	if len(a) == 0 {
		return map[string]CertificateCertificateDescriptionPublicKey{}
	}

	items := make(map[string]CertificateCertificateDescriptionPublicKey)
	for k, item := range a {
		items[k] = *flattenCertificateCertificateDescriptionPublicKey(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenCertificateCertificateDescriptionPublicKeySlice flattens the contents of CertificateCertificateDescriptionPublicKey from a JSON
// response object.
func flattenCertificateCertificateDescriptionPublicKeySlice(c *Client, i interface{}, res *Certificate) []CertificateCertificateDescriptionPublicKey {
	a, ok := i.([]interface{})
	if !ok {
		return []CertificateCertificateDescriptionPublicKey{}
	}

	if len(a) == 0 {
		return []CertificateCertificateDescriptionPublicKey{}
	}

	items := make([]CertificateCertificateDescriptionPublicKey, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenCertificateCertificateDescriptionPublicKey(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandCertificateCertificateDescriptionPublicKey expands an instance of CertificateCertificateDescriptionPublicKey into a JSON
// request object.
func expandCertificateCertificateDescriptionPublicKey(c *Client, f *CertificateCertificateDescriptionPublicKey, res *Certificate) (map[string]interface{}, error) {
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

// flattenCertificateCertificateDescriptionPublicKey flattens an instance of CertificateCertificateDescriptionPublicKey from a JSON
// response object.
func flattenCertificateCertificateDescriptionPublicKey(c *Client, i interface{}, res *Certificate) *CertificateCertificateDescriptionPublicKey {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &CertificateCertificateDescriptionPublicKey{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyCertificateCertificateDescriptionPublicKey
	}
	r.Key = dcl.FlattenString(m["key"])
	r.Format = flattenCertificateCertificateDescriptionPublicKeyFormatEnum(m["format"])

	return r
}

// expandCertificateCertificateDescriptionSubjectKeyIdMap expands the contents of CertificateCertificateDescriptionSubjectKeyId into a JSON
// request object.
func expandCertificateCertificateDescriptionSubjectKeyIdMap(c *Client, f map[string]CertificateCertificateDescriptionSubjectKeyId, res *Certificate) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandCertificateCertificateDescriptionSubjectKeyId(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandCertificateCertificateDescriptionSubjectKeyIdSlice expands the contents of CertificateCertificateDescriptionSubjectKeyId into a JSON
// request object.
func expandCertificateCertificateDescriptionSubjectKeyIdSlice(c *Client, f []CertificateCertificateDescriptionSubjectKeyId, res *Certificate) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandCertificateCertificateDescriptionSubjectKeyId(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenCertificateCertificateDescriptionSubjectKeyIdMap flattens the contents of CertificateCertificateDescriptionSubjectKeyId from a JSON
// response object.
func flattenCertificateCertificateDescriptionSubjectKeyIdMap(c *Client, i interface{}, res *Certificate) map[string]CertificateCertificateDescriptionSubjectKeyId {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]CertificateCertificateDescriptionSubjectKeyId{}
	}

	if len(a) == 0 {
		return map[string]CertificateCertificateDescriptionSubjectKeyId{}
	}

	items := make(map[string]CertificateCertificateDescriptionSubjectKeyId)
	for k, item := range a {
		items[k] = *flattenCertificateCertificateDescriptionSubjectKeyId(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenCertificateCertificateDescriptionSubjectKeyIdSlice flattens the contents of CertificateCertificateDescriptionSubjectKeyId from a JSON
// response object.
func flattenCertificateCertificateDescriptionSubjectKeyIdSlice(c *Client, i interface{}, res *Certificate) []CertificateCertificateDescriptionSubjectKeyId {
	a, ok := i.([]interface{})
	if !ok {
		return []CertificateCertificateDescriptionSubjectKeyId{}
	}

	if len(a) == 0 {
		return []CertificateCertificateDescriptionSubjectKeyId{}
	}

	items := make([]CertificateCertificateDescriptionSubjectKeyId, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenCertificateCertificateDescriptionSubjectKeyId(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandCertificateCertificateDescriptionSubjectKeyId expands an instance of CertificateCertificateDescriptionSubjectKeyId into a JSON
// request object.
func expandCertificateCertificateDescriptionSubjectKeyId(c *Client, f *CertificateCertificateDescriptionSubjectKeyId, res *Certificate) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.KeyId; !dcl.IsEmptyValueIndirect(v) {
		m["keyId"] = v
	}

	return m, nil
}

// flattenCertificateCertificateDescriptionSubjectKeyId flattens an instance of CertificateCertificateDescriptionSubjectKeyId from a JSON
// response object.
func flattenCertificateCertificateDescriptionSubjectKeyId(c *Client, i interface{}, res *Certificate) *CertificateCertificateDescriptionSubjectKeyId {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &CertificateCertificateDescriptionSubjectKeyId{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyCertificateCertificateDescriptionSubjectKeyId
	}
	r.KeyId = dcl.FlattenString(m["keyId"])

	return r
}

// expandCertificateCertificateDescriptionAuthorityKeyIdMap expands the contents of CertificateCertificateDescriptionAuthorityKeyId into a JSON
// request object.
func expandCertificateCertificateDescriptionAuthorityKeyIdMap(c *Client, f map[string]CertificateCertificateDescriptionAuthorityKeyId, res *Certificate) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandCertificateCertificateDescriptionAuthorityKeyId(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandCertificateCertificateDescriptionAuthorityKeyIdSlice expands the contents of CertificateCertificateDescriptionAuthorityKeyId into a JSON
// request object.
func expandCertificateCertificateDescriptionAuthorityKeyIdSlice(c *Client, f []CertificateCertificateDescriptionAuthorityKeyId, res *Certificate) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandCertificateCertificateDescriptionAuthorityKeyId(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenCertificateCertificateDescriptionAuthorityKeyIdMap flattens the contents of CertificateCertificateDescriptionAuthorityKeyId from a JSON
// response object.
func flattenCertificateCertificateDescriptionAuthorityKeyIdMap(c *Client, i interface{}, res *Certificate) map[string]CertificateCertificateDescriptionAuthorityKeyId {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]CertificateCertificateDescriptionAuthorityKeyId{}
	}

	if len(a) == 0 {
		return map[string]CertificateCertificateDescriptionAuthorityKeyId{}
	}

	items := make(map[string]CertificateCertificateDescriptionAuthorityKeyId)
	for k, item := range a {
		items[k] = *flattenCertificateCertificateDescriptionAuthorityKeyId(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenCertificateCertificateDescriptionAuthorityKeyIdSlice flattens the contents of CertificateCertificateDescriptionAuthorityKeyId from a JSON
// response object.
func flattenCertificateCertificateDescriptionAuthorityKeyIdSlice(c *Client, i interface{}, res *Certificate) []CertificateCertificateDescriptionAuthorityKeyId {
	a, ok := i.([]interface{})
	if !ok {
		return []CertificateCertificateDescriptionAuthorityKeyId{}
	}

	if len(a) == 0 {
		return []CertificateCertificateDescriptionAuthorityKeyId{}
	}

	items := make([]CertificateCertificateDescriptionAuthorityKeyId, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenCertificateCertificateDescriptionAuthorityKeyId(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandCertificateCertificateDescriptionAuthorityKeyId expands an instance of CertificateCertificateDescriptionAuthorityKeyId into a JSON
// request object.
func expandCertificateCertificateDescriptionAuthorityKeyId(c *Client, f *CertificateCertificateDescriptionAuthorityKeyId, res *Certificate) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.KeyId; !dcl.IsEmptyValueIndirect(v) {
		m["keyId"] = v
	}

	return m, nil
}

// flattenCertificateCertificateDescriptionAuthorityKeyId flattens an instance of CertificateCertificateDescriptionAuthorityKeyId from a JSON
// response object.
func flattenCertificateCertificateDescriptionAuthorityKeyId(c *Client, i interface{}, res *Certificate) *CertificateCertificateDescriptionAuthorityKeyId {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &CertificateCertificateDescriptionAuthorityKeyId{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyCertificateCertificateDescriptionAuthorityKeyId
	}
	r.KeyId = dcl.FlattenString(m["keyId"])

	return r
}

// expandCertificateCertificateDescriptionCertFingerprintMap expands the contents of CertificateCertificateDescriptionCertFingerprint into a JSON
// request object.
func expandCertificateCertificateDescriptionCertFingerprintMap(c *Client, f map[string]CertificateCertificateDescriptionCertFingerprint, res *Certificate) (map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := make(map[string]interface{})
	for k, item := range f {
		i, err := expandCertificateCertificateDescriptionCertFingerprint(c, &item, res)
		if err != nil {
			return nil, err
		}
		if i != nil {
			items[k] = i
		}
	}

	return items, nil
}

// expandCertificateCertificateDescriptionCertFingerprintSlice expands the contents of CertificateCertificateDescriptionCertFingerprint into a JSON
// request object.
func expandCertificateCertificateDescriptionCertFingerprintSlice(c *Client, f []CertificateCertificateDescriptionCertFingerprint, res *Certificate) ([]map[string]interface{}, error) {
	if f == nil {
		return nil, nil
	}

	items := []map[string]interface{}{}
	for _, item := range f {
		i, err := expandCertificateCertificateDescriptionCertFingerprint(c, &item, res)
		if err != nil {
			return nil, err
		}

		items = append(items, i)
	}

	return items, nil
}

// flattenCertificateCertificateDescriptionCertFingerprintMap flattens the contents of CertificateCertificateDescriptionCertFingerprint from a JSON
// response object.
func flattenCertificateCertificateDescriptionCertFingerprintMap(c *Client, i interface{}, res *Certificate) map[string]CertificateCertificateDescriptionCertFingerprint {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]CertificateCertificateDescriptionCertFingerprint{}
	}

	if len(a) == 0 {
		return map[string]CertificateCertificateDescriptionCertFingerprint{}
	}

	items := make(map[string]CertificateCertificateDescriptionCertFingerprint)
	for k, item := range a {
		items[k] = *flattenCertificateCertificateDescriptionCertFingerprint(c, item.(map[string]interface{}), res)
	}

	return items
}

// flattenCertificateCertificateDescriptionCertFingerprintSlice flattens the contents of CertificateCertificateDescriptionCertFingerprint from a JSON
// response object.
func flattenCertificateCertificateDescriptionCertFingerprintSlice(c *Client, i interface{}, res *Certificate) []CertificateCertificateDescriptionCertFingerprint {
	a, ok := i.([]interface{})
	if !ok {
		return []CertificateCertificateDescriptionCertFingerprint{}
	}

	if len(a) == 0 {
		return []CertificateCertificateDescriptionCertFingerprint{}
	}

	items := make([]CertificateCertificateDescriptionCertFingerprint, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenCertificateCertificateDescriptionCertFingerprint(c, item.(map[string]interface{}), res))
	}

	return items
}

// expandCertificateCertificateDescriptionCertFingerprint expands an instance of CertificateCertificateDescriptionCertFingerprint into a JSON
// request object.
func expandCertificateCertificateDescriptionCertFingerprint(c *Client, f *CertificateCertificateDescriptionCertFingerprint, res *Certificate) (map[string]interface{}, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]interface{})
	if v := f.Sha256Hash; !dcl.IsEmptyValueIndirect(v) {
		m["sha256Hash"] = v
	}

	return m, nil
}

// flattenCertificateCertificateDescriptionCertFingerprint flattens an instance of CertificateCertificateDescriptionCertFingerprint from a JSON
// response object.
func flattenCertificateCertificateDescriptionCertFingerprint(c *Client, i interface{}, res *Certificate) *CertificateCertificateDescriptionCertFingerprint {
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	r := &CertificateCertificateDescriptionCertFingerprint{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyCertificateCertificateDescriptionCertFingerprint
	}
	r.Sha256Hash = dcl.FlattenString(m["sha256Hash"])

	return r
}

// flattenCertificateConfigPublicKeyFormatEnumMap flattens the contents of CertificateConfigPublicKeyFormatEnum from a JSON
// response object.
func flattenCertificateConfigPublicKeyFormatEnumMap(c *Client, i interface{}, res *Certificate) map[string]CertificateConfigPublicKeyFormatEnum {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]CertificateConfigPublicKeyFormatEnum{}
	}

	if len(a) == 0 {
		return map[string]CertificateConfigPublicKeyFormatEnum{}
	}

	items := make(map[string]CertificateConfigPublicKeyFormatEnum)
	for k, item := range a {
		items[k] = *flattenCertificateConfigPublicKeyFormatEnum(item.(interface{}))
	}

	return items
}

// flattenCertificateConfigPublicKeyFormatEnumSlice flattens the contents of CertificateConfigPublicKeyFormatEnum from a JSON
// response object.
func flattenCertificateConfigPublicKeyFormatEnumSlice(c *Client, i interface{}, res *Certificate) []CertificateConfigPublicKeyFormatEnum {
	a, ok := i.([]interface{})
	if !ok {
		return []CertificateConfigPublicKeyFormatEnum{}
	}

	if len(a) == 0 {
		return []CertificateConfigPublicKeyFormatEnum{}
	}

	items := make([]CertificateConfigPublicKeyFormatEnum, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenCertificateConfigPublicKeyFormatEnum(item.(interface{})))
	}

	return items
}

// flattenCertificateConfigPublicKeyFormatEnum asserts that an interface is a string, and returns a
// pointer to a *CertificateConfigPublicKeyFormatEnum with the same value as that string.
func flattenCertificateConfigPublicKeyFormatEnum(i interface{}) *CertificateConfigPublicKeyFormatEnum {
	s, ok := i.(string)
	if !ok {
		return nil
	}

	return CertificateConfigPublicKeyFormatEnumRef(s)
}

// flattenCertificateSubjectModeEnumMap flattens the contents of CertificateSubjectModeEnum from a JSON
// response object.
func flattenCertificateSubjectModeEnumMap(c *Client, i interface{}, res *Certificate) map[string]CertificateSubjectModeEnum {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]CertificateSubjectModeEnum{}
	}

	if len(a) == 0 {
		return map[string]CertificateSubjectModeEnum{}
	}

	items := make(map[string]CertificateSubjectModeEnum)
	for k, item := range a {
		items[k] = *flattenCertificateSubjectModeEnum(item.(interface{}))
	}

	return items
}

// flattenCertificateSubjectModeEnumSlice flattens the contents of CertificateSubjectModeEnum from a JSON
// response object.
func flattenCertificateSubjectModeEnumSlice(c *Client, i interface{}, res *Certificate) []CertificateSubjectModeEnum {
	a, ok := i.([]interface{})
	if !ok {
		return []CertificateSubjectModeEnum{}
	}

	if len(a) == 0 {
		return []CertificateSubjectModeEnum{}
	}

	items := make([]CertificateSubjectModeEnum, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenCertificateSubjectModeEnum(item.(interface{})))
	}

	return items
}

// flattenCertificateSubjectModeEnum asserts that an interface is a string, and returns a
// pointer to a *CertificateSubjectModeEnum with the same value as that string.
func flattenCertificateSubjectModeEnum(i interface{}) *CertificateSubjectModeEnum {
	s, ok := i.(string)
	if !ok {
		return nil
	}

	return CertificateSubjectModeEnumRef(s)
}

// flattenCertificateRevocationDetailsRevocationStateEnumMap flattens the contents of CertificateRevocationDetailsRevocationStateEnum from a JSON
// response object.
func flattenCertificateRevocationDetailsRevocationStateEnumMap(c *Client, i interface{}, res *Certificate) map[string]CertificateRevocationDetailsRevocationStateEnum {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]CertificateRevocationDetailsRevocationStateEnum{}
	}

	if len(a) == 0 {
		return map[string]CertificateRevocationDetailsRevocationStateEnum{}
	}

	items := make(map[string]CertificateRevocationDetailsRevocationStateEnum)
	for k, item := range a {
		items[k] = *flattenCertificateRevocationDetailsRevocationStateEnum(item.(interface{}))
	}

	return items
}

// flattenCertificateRevocationDetailsRevocationStateEnumSlice flattens the contents of CertificateRevocationDetailsRevocationStateEnum from a JSON
// response object.
func flattenCertificateRevocationDetailsRevocationStateEnumSlice(c *Client, i interface{}, res *Certificate) []CertificateRevocationDetailsRevocationStateEnum {
	a, ok := i.([]interface{})
	if !ok {
		return []CertificateRevocationDetailsRevocationStateEnum{}
	}

	if len(a) == 0 {
		return []CertificateRevocationDetailsRevocationStateEnum{}
	}

	items := make([]CertificateRevocationDetailsRevocationStateEnum, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenCertificateRevocationDetailsRevocationStateEnum(item.(interface{})))
	}

	return items
}

// flattenCertificateRevocationDetailsRevocationStateEnum asserts that an interface is a string, and returns a
// pointer to a *CertificateRevocationDetailsRevocationStateEnum with the same value as that string.
func flattenCertificateRevocationDetailsRevocationStateEnum(i interface{}) *CertificateRevocationDetailsRevocationStateEnum {
	s, ok := i.(string)
	if !ok {
		return nil
	}

	return CertificateRevocationDetailsRevocationStateEnumRef(s)
}

// flattenCertificateCertificateDescriptionPublicKeyFormatEnumMap flattens the contents of CertificateCertificateDescriptionPublicKeyFormatEnum from a JSON
// response object.
func flattenCertificateCertificateDescriptionPublicKeyFormatEnumMap(c *Client, i interface{}, res *Certificate) map[string]CertificateCertificateDescriptionPublicKeyFormatEnum {
	a, ok := i.(map[string]interface{})
	if !ok {
		return map[string]CertificateCertificateDescriptionPublicKeyFormatEnum{}
	}

	if len(a) == 0 {
		return map[string]CertificateCertificateDescriptionPublicKeyFormatEnum{}
	}

	items := make(map[string]CertificateCertificateDescriptionPublicKeyFormatEnum)
	for k, item := range a {
		items[k] = *flattenCertificateCertificateDescriptionPublicKeyFormatEnum(item.(interface{}))
	}

	return items
}

// flattenCertificateCertificateDescriptionPublicKeyFormatEnumSlice flattens the contents of CertificateCertificateDescriptionPublicKeyFormatEnum from a JSON
// response object.
func flattenCertificateCertificateDescriptionPublicKeyFormatEnumSlice(c *Client, i interface{}, res *Certificate) []CertificateCertificateDescriptionPublicKeyFormatEnum {
	a, ok := i.([]interface{})
	if !ok {
		return []CertificateCertificateDescriptionPublicKeyFormatEnum{}
	}

	if len(a) == 0 {
		return []CertificateCertificateDescriptionPublicKeyFormatEnum{}
	}

	items := make([]CertificateCertificateDescriptionPublicKeyFormatEnum, 0, len(a))
	for _, item := range a {
		items = append(items, *flattenCertificateCertificateDescriptionPublicKeyFormatEnum(item.(interface{})))
	}

	return items
}

// flattenCertificateCertificateDescriptionPublicKeyFormatEnum asserts that an interface is a string, and returns a
// pointer to a *CertificateCertificateDescriptionPublicKeyFormatEnum with the same value as that string.
func flattenCertificateCertificateDescriptionPublicKeyFormatEnum(i interface{}) *CertificateCertificateDescriptionPublicKeyFormatEnum {
	s, ok := i.(string)
	if !ok {
		return nil
	}

	return CertificateCertificateDescriptionPublicKeyFormatEnumRef(s)
}

// This function returns a matcher that checks whether a serialized resource matches this resource
// in its parameters (as defined by the fields in a Get, which definitionally define resource
// identity).  This is useful in extracting the element from a List call.
func (r *Certificate) matcher(c *Client) func([]byte) bool {
	return func(b []byte) bool {
		cr, err := unmarshalCertificate(b, c, r)
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

type certificateDiff struct {
	// The diff should include one or the other of RequiresRecreate or UpdateOp.
	RequiresRecreate bool
	UpdateOp         certificateApiOperation
	FieldName        string // used for error logging
}

func convertFieldDiffsToCertificateDiffs(config *dcl.Config, fds []*dcl.FieldDiff, opts []dcl.ApplyOption) ([]certificateDiff, error) {
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
	var diffs []certificateDiff
	// For each operation name, create a certificateDiff which contains the operation.
	for opName, fieldDiffs := range opNamesToFieldDiffs {
		// Use the first field diff's field name for logging required recreate error.
		diff := certificateDiff{FieldName: fieldDiffs[0].FieldName}
		if opName == "Recreate" {
			diff.RequiresRecreate = true
		} else {
			apiOp, err := convertOpNameToCertificateApiOperation(opName, fieldDiffs, opts...)
			if err != nil {
				return diffs, err
			}
			diff.UpdateOp = apiOp
		}
		diffs = append(diffs, diff)
	}
	return diffs, nil
}

func convertOpNameToCertificateApiOperation(opName string, fieldDiffs []*dcl.FieldDiff, opts ...dcl.ApplyOption) (certificateApiOperation, error) {
	switch opName {

	case "updateCertificateUpdateCertificateOperation":
		return &updateCertificateUpdateCertificateOperation{FieldDiffs: fieldDiffs}, nil

	default:
		return nil, fmt.Errorf("no such operation with name: %v", opName)
	}
}

func extractCertificateFields(r *Certificate) error {
	vConfig := r.Config
	if vConfig == nil {
		// note: explicitly not the empty object.
		vConfig = &CertificateConfig{}
	}
	if err := extractCertificateConfigFields(r, vConfig); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vConfig) {
		r.Config = vConfig
	}
	vRevocationDetails := r.RevocationDetails
	if vRevocationDetails == nil {
		// note: explicitly not the empty object.
		vRevocationDetails = &CertificateRevocationDetails{}
	}
	if err := extractCertificateRevocationDetailsFields(r, vRevocationDetails); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vRevocationDetails) {
		r.RevocationDetails = vRevocationDetails
	}
	vCertificateDescription := r.CertificateDescription
	if vCertificateDescription == nil {
		// note: explicitly not the empty object.
		vCertificateDescription = &CertificateCertificateDescription{}
	}
	if err := extractCertificateCertificateDescriptionFields(r, vCertificateDescription); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vCertificateDescription) {
		r.CertificateDescription = vCertificateDescription
	}
	return nil
}
func extractCertificateConfigFields(r *Certificate, o *CertificateConfig) error {
	vSubjectConfig := o.SubjectConfig
	if vSubjectConfig == nil {
		// note: explicitly not the empty object.
		vSubjectConfig = &CertificateConfigSubjectConfig{}
	}
	if err := extractCertificateConfigSubjectConfigFields(r, vSubjectConfig); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vSubjectConfig) {
		o.SubjectConfig = vSubjectConfig
	}
	vX509Config := o.X509Config
	if vX509Config == nil {
		// note: explicitly not the empty object.
		vX509Config = &CertificateConfigX509Config{}
	}
	if err := extractCertificateConfigX509ConfigFields(r, vX509Config); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vX509Config) {
		o.X509Config = vX509Config
	}
	vPublicKey := o.PublicKey
	if vPublicKey == nil {
		// note: explicitly not the empty object.
		vPublicKey = &CertificateConfigPublicKey{}
	}
	if err := extractCertificateConfigPublicKeyFields(r, vPublicKey); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vPublicKey) {
		o.PublicKey = vPublicKey
	}
	return nil
}
func extractCertificateConfigSubjectConfigFields(r *Certificate, o *CertificateConfigSubjectConfig) error {
	vSubject := o.Subject
	if vSubject == nil {
		// note: explicitly not the empty object.
		vSubject = &CertificateConfigSubjectConfigSubject{}
	}
	if err := extractCertificateConfigSubjectConfigSubjectFields(r, vSubject); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vSubject) {
		o.Subject = vSubject
	}
	vSubjectAltName := o.SubjectAltName
	if vSubjectAltName == nil {
		// note: explicitly not the empty object.
		vSubjectAltName = &CertificateConfigSubjectConfigSubjectAltName{}
	}
	if err := extractCertificateConfigSubjectConfigSubjectAltNameFields(r, vSubjectAltName); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vSubjectAltName) {
		o.SubjectAltName = vSubjectAltName
	}
	return nil
}
func extractCertificateConfigSubjectConfigSubjectFields(r *Certificate, o *CertificateConfigSubjectConfigSubject) error {
	return nil
}
func extractCertificateConfigSubjectConfigSubjectAltNameFields(r *Certificate, o *CertificateConfigSubjectConfigSubjectAltName) error {
	return nil
}
func extractCertificateConfigX509ConfigFields(r *Certificate, o *CertificateConfigX509Config) error {
	vKeyUsage := o.KeyUsage
	if vKeyUsage == nil {
		// note: explicitly not the empty object.
		vKeyUsage = &CertificateConfigX509ConfigKeyUsage{}
	}
	if err := extractCertificateConfigX509ConfigKeyUsageFields(r, vKeyUsage); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vKeyUsage) {
		o.KeyUsage = vKeyUsage
	}
	vCaOptions := o.CaOptions
	if vCaOptions == nil {
		// note: explicitly not the empty object.
		vCaOptions = &CertificateConfigX509ConfigCaOptions{}
	}
	if err := extractCertificateConfigX509ConfigCaOptionsFields(r, vCaOptions); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vCaOptions) {
		o.CaOptions = vCaOptions
	}
	return nil
}
func extractCertificateConfigX509ConfigKeyUsageFields(r *Certificate, o *CertificateConfigX509ConfigKeyUsage) error {
	vBaseKeyUsage := o.BaseKeyUsage
	if vBaseKeyUsage == nil {
		// note: explicitly not the empty object.
		vBaseKeyUsage = &CertificateConfigX509ConfigKeyUsageBaseKeyUsage{}
	}
	if err := extractCertificateConfigX509ConfigKeyUsageBaseKeyUsageFields(r, vBaseKeyUsage); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vBaseKeyUsage) {
		o.BaseKeyUsage = vBaseKeyUsage
	}
	vExtendedKeyUsage := o.ExtendedKeyUsage
	if vExtendedKeyUsage == nil {
		// note: explicitly not the empty object.
		vExtendedKeyUsage = &CertificateConfigX509ConfigKeyUsageExtendedKeyUsage{}
	}
	if err := extractCertificateConfigX509ConfigKeyUsageExtendedKeyUsageFields(r, vExtendedKeyUsage); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vExtendedKeyUsage) {
		o.ExtendedKeyUsage = vExtendedKeyUsage
	}
	return nil
}
func extractCertificateConfigX509ConfigKeyUsageBaseKeyUsageFields(r *Certificate, o *CertificateConfigX509ConfigKeyUsageBaseKeyUsage) error {
	return nil
}
func extractCertificateConfigX509ConfigKeyUsageExtendedKeyUsageFields(r *Certificate, o *CertificateConfigX509ConfigKeyUsageExtendedKeyUsage) error {
	return nil
}
func extractCertificateConfigX509ConfigKeyUsageUnknownExtendedKeyUsagesFields(r *Certificate, o *CertificateConfigX509ConfigKeyUsageUnknownExtendedKeyUsages) error {
	return nil
}
func extractCertificateConfigX509ConfigCaOptionsFields(r *Certificate, o *CertificateConfigX509ConfigCaOptions) error {
	return nil
}
func extractCertificateConfigX509ConfigPolicyIdsFields(r *Certificate, o *CertificateConfigX509ConfigPolicyIds) error {
	return nil
}
func extractCertificateConfigX509ConfigAdditionalExtensionsFields(r *Certificate, o *CertificateConfigX509ConfigAdditionalExtensions) error {
	vObjectId := o.ObjectId
	if vObjectId == nil {
		// note: explicitly not the empty object.
		vObjectId = &CertificateConfigX509ConfigAdditionalExtensionsObjectId{}
	}
	if err := extractCertificateConfigX509ConfigAdditionalExtensionsObjectIdFields(r, vObjectId); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vObjectId) {
		o.ObjectId = vObjectId
	}
	return nil
}
func extractCertificateConfigX509ConfigAdditionalExtensionsObjectIdFields(r *Certificate, o *CertificateConfigX509ConfigAdditionalExtensionsObjectId) error {
	return nil
}
func extractCertificateConfigPublicKeyFields(r *Certificate, o *CertificateConfigPublicKey) error {
	return nil
}
func extractCertificateRevocationDetailsFields(r *Certificate, o *CertificateRevocationDetails) error {
	return nil
}
func extractCertificateCertificateDescriptionFields(r *Certificate, o *CertificateCertificateDescription) error {
	vSubjectDescription := o.SubjectDescription
	if vSubjectDescription == nil {
		// note: explicitly not the empty object.
		vSubjectDescription = &CertificateCertificateDescriptionSubjectDescription{}
	}
	if err := extractCertificateCertificateDescriptionSubjectDescriptionFields(r, vSubjectDescription); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vSubjectDescription) {
		o.SubjectDescription = vSubjectDescription
	}
	vX509Description := o.X509Description
	if vX509Description == nil {
		// note: explicitly not the empty object.
		vX509Description = &CertificateCertificateDescriptionX509Description{}
	}
	if err := extractCertificateCertificateDescriptionX509DescriptionFields(r, vX509Description); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vX509Description) {
		o.X509Description = vX509Description
	}
	vPublicKey := o.PublicKey
	if vPublicKey == nil {
		// note: explicitly not the empty object.
		vPublicKey = &CertificateCertificateDescriptionPublicKey{}
	}
	if err := extractCertificateCertificateDescriptionPublicKeyFields(r, vPublicKey); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vPublicKey) {
		o.PublicKey = vPublicKey
	}
	vSubjectKeyId := o.SubjectKeyId
	if vSubjectKeyId == nil {
		// note: explicitly not the empty object.
		vSubjectKeyId = &CertificateCertificateDescriptionSubjectKeyId{}
	}
	if err := extractCertificateCertificateDescriptionSubjectKeyIdFields(r, vSubjectKeyId); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vSubjectKeyId) {
		o.SubjectKeyId = vSubjectKeyId
	}
	vAuthorityKeyId := o.AuthorityKeyId
	if vAuthorityKeyId == nil {
		// note: explicitly not the empty object.
		vAuthorityKeyId = &CertificateCertificateDescriptionAuthorityKeyId{}
	}
	if err := extractCertificateCertificateDescriptionAuthorityKeyIdFields(r, vAuthorityKeyId); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vAuthorityKeyId) {
		o.AuthorityKeyId = vAuthorityKeyId
	}
	vCertFingerprint := o.CertFingerprint
	if vCertFingerprint == nil {
		// note: explicitly not the empty object.
		vCertFingerprint = &CertificateCertificateDescriptionCertFingerprint{}
	}
	if err := extractCertificateCertificateDescriptionCertFingerprintFields(r, vCertFingerprint); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vCertFingerprint) {
		o.CertFingerprint = vCertFingerprint
	}
	return nil
}
func extractCertificateCertificateDescriptionSubjectDescriptionFields(r *Certificate, o *CertificateCertificateDescriptionSubjectDescription) error {
	vSubject := o.Subject
	if vSubject == nil {
		// note: explicitly not the empty object.
		vSubject = &CertificateCertificateDescriptionSubjectDescriptionSubject{}
	}
	if err := extractCertificateCertificateDescriptionSubjectDescriptionSubjectFields(r, vSubject); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vSubject) {
		o.Subject = vSubject
	}
	vSubjectAltName := o.SubjectAltName
	if vSubjectAltName == nil {
		// note: explicitly not the empty object.
		vSubjectAltName = &CertificateCertificateDescriptionSubjectDescriptionSubjectAltName{}
	}
	if err := extractCertificateCertificateDescriptionSubjectDescriptionSubjectAltNameFields(r, vSubjectAltName); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vSubjectAltName) {
		o.SubjectAltName = vSubjectAltName
	}
	return nil
}
func extractCertificateCertificateDescriptionSubjectDescriptionSubjectFields(r *Certificate, o *CertificateCertificateDescriptionSubjectDescriptionSubject) error {
	return nil
}
func extractCertificateCertificateDescriptionSubjectDescriptionSubjectAltNameFields(r *Certificate, o *CertificateCertificateDescriptionSubjectDescriptionSubjectAltName) error {
	return nil
}
func extractCertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSansFields(r *Certificate, o *CertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSans) error {
	vObjectId := o.ObjectId
	if vObjectId == nil {
		// note: explicitly not the empty object.
		vObjectId = &CertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSansObjectId{}
	}
	if err := extractCertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSansObjectIdFields(r, vObjectId); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vObjectId) {
		o.ObjectId = vObjectId
	}
	return nil
}
func extractCertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSansObjectIdFields(r *Certificate, o *CertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSansObjectId) error {
	return nil
}
func extractCertificateCertificateDescriptionX509DescriptionFields(r *Certificate, o *CertificateCertificateDescriptionX509Description) error {
	vKeyUsage := o.KeyUsage
	if vKeyUsage == nil {
		// note: explicitly not the empty object.
		vKeyUsage = &CertificateCertificateDescriptionX509DescriptionKeyUsage{}
	}
	if err := extractCertificateCertificateDescriptionX509DescriptionKeyUsageFields(r, vKeyUsage); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vKeyUsage) {
		o.KeyUsage = vKeyUsage
	}
	vCaOptions := o.CaOptions
	if vCaOptions == nil {
		// note: explicitly not the empty object.
		vCaOptions = &CertificateCertificateDescriptionX509DescriptionCaOptions{}
	}
	if err := extractCertificateCertificateDescriptionX509DescriptionCaOptionsFields(r, vCaOptions); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vCaOptions) {
		o.CaOptions = vCaOptions
	}
	return nil
}
func extractCertificateCertificateDescriptionX509DescriptionKeyUsageFields(r *Certificate, o *CertificateCertificateDescriptionX509DescriptionKeyUsage) error {
	vBaseKeyUsage := o.BaseKeyUsage
	if vBaseKeyUsage == nil {
		// note: explicitly not the empty object.
		vBaseKeyUsage = &CertificateCertificateDescriptionX509DescriptionKeyUsageBaseKeyUsage{}
	}
	if err := extractCertificateCertificateDescriptionX509DescriptionKeyUsageBaseKeyUsageFields(r, vBaseKeyUsage); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vBaseKeyUsage) {
		o.BaseKeyUsage = vBaseKeyUsage
	}
	vExtendedKeyUsage := o.ExtendedKeyUsage
	if vExtendedKeyUsage == nil {
		// note: explicitly not the empty object.
		vExtendedKeyUsage = &CertificateCertificateDescriptionX509DescriptionKeyUsageExtendedKeyUsage{}
	}
	if err := extractCertificateCertificateDescriptionX509DescriptionKeyUsageExtendedKeyUsageFields(r, vExtendedKeyUsage); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vExtendedKeyUsage) {
		o.ExtendedKeyUsage = vExtendedKeyUsage
	}
	return nil
}
func extractCertificateCertificateDescriptionX509DescriptionKeyUsageBaseKeyUsageFields(r *Certificate, o *CertificateCertificateDescriptionX509DescriptionKeyUsageBaseKeyUsage) error {
	return nil
}
func extractCertificateCertificateDescriptionX509DescriptionKeyUsageExtendedKeyUsageFields(r *Certificate, o *CertificateCertificateDescriptionX509DescriptionKeyUsageExtendedKeyUsage) error {
	return nil
}
func extractCertificateCertificateDescriptionX509DescriptionKeyUsageUnknownExtendedKeyUsagesFields(r *Certificate, o *CertificateCertificateDescriptionX509DescriptionKeyUsageUnknownExtendedKeyUsages) error {
	return nil
}
func extractCertificateCertificateDescriptionX509DescriptionCaOptionsFields(r *Certificate, o *CertificateCertificateDescriptionX509DescriptionCaOptions) error {
	return nil
}
func extractCertificateCertificateDescriptionX509DescriptionPolicyIdsFields(r *Certificate, o *CertificateCertificateDescriptionX509DescriptionPolicyIds) error {
	return nil
}
func extractCertificateCertificateDescriptionX509DescriptionAdditionalExtensionsFields(r *Certificate, o *CertificateCertificateDescriptionX509DescriptionAdditionalExtensions) error {
	vObjectId := o.ObjectId
	if vObjectId == nil {
		// note: explicitly not the empty object.
		vObjectId = &CertificateCertificateDescriptionX509DescriptionAdditionalExtensionsObjectId{}
	}
	if err := extractCertificateCertificateDescriptionX509DescriptionAdditionalExtensionsObjectIdFields(r, vObjectId); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vObjectId) {
		o.ObjectId = vObjectId
	}
	return nil
}
func extractCertificateCertificateDescriptionX509DescriptionAdditionalExtensionsObjectIdFields(r *Certificate, o *CertificateCertificateDescriptionX509DescriptionAdditionalExtensionsObjectId) error {
	return nil
}
func extractCertificateCertificateDescriptionPublicKeyFields(r *Certificate, o *CertificateCertificateDescriptionPublicKey) error {
	return nil
}
func extractCertificateCertificateDescriptionSubjectKeyIdFields(r *Certificate, o *CertificateCertificateDescriptionSubjectKeyId) error {
	return nil
}
func extractCertificateCertificateDescriptionAuthorityKeyIdFields(r *Certificate, o *CertificateCertificateDescriptionAuthorityKeyId) error {
	return nil
}
func extractCertificateCertificateDescriptionCertFingerprintFields(r *Certificate, o *CertificateCertificateDescriptionCertFingerprint) error {
	return nil
}

func postReadExtractCertificateFields(r *Certificate) error {
	vConfig := r.Config
	if vConfig == nil {
		// note: explicitly not the empty object.
		vConfig = &CertificateConfig{}
	}
	if err := postReadExtractCertificateConfigFields(r, vConfig); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vConfig) {
		r.Config = vConfig
	}
	vRevocationDetails := r.RevocationDetails
	if vRevocationDetails == nil {
		// note: explicitly not the empty object.
		vRevocationDetails = &CertificateRevocationDetails{}
	}
	if err := postReadExtractCertificateRevocationDetailsFields(r, vRevocationDetails); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vRevocationDetails) {
		r.RevocationDetails = vRevocationDetails
	}
	vCertificateDescription := r.CertificateDescription
	if vCertificateDescription == nil {
		// note: explicitly not the empty object.
		vCertificateDescription = &CertificateCertificateDescription{}
	}
	if err := postReadExtractCertificateCertificateDescriptionFields(r, vCertificateDescription); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vCertificateDescription) {
		r.CertificateDescription = vCertificateDescription
	}
	return nil
}
func postReadExtractCertificateConfigFields(r *Certificate, o *CertificateConfig) error {
	vSubjectConfig := o.SubjectConfig
	if vSubjectConfig == nil {
		// note: explicitly not the empty object.
		vSubjectConfig = &CertificateConfigSubjectConfig{}
	}
	if err := extractCertificateConfigSubjectConfigFields(r, vSubjectConfig); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vSubjectConfig) {
		o.SubjectConfig = vSubjectConfig
	}
	vX509Config := o.X509Config
	if vX509Config == nil {
		// note: explicitly not the empty object.
		vX509Config = &CertificateConfigX509Config{}
	}
	if err := extractCertificateConfigX509ConfigFields(r, vX509Config); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vX509Config) {
		o.X509Config = vX509Config
	}
	vPublicKey := o.PublicKey
	if vPublicKey == nil {
		// note: explicitly not the empty object.
		vPublicKey = &CertificateConfigPublicKey{}
	}
	if err := extractCertificateConfigPublicKeyFields(r, vPublicKey); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vPublicKey) {
		o.PublicKey = vPublicKey
	}
	return nil
}
func postReadExtractCertificateConfigSubjectConfigFields(r *Certificate, o *CertificateConfigSubjectConfig) error {
	vSubject := o.Subject
	if vSubject == nil {
		// note: explicitly not the empty object.
		vSubject = &CertificateConfigSubjectConfigSubject{}
	}
	if err := extractCertificateConfigSubjectConfigSubjectFields(r, vSubject); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vSubject) {
		o.Subject = vSubject
	}
	vSubjectAltName := o.SubjectAltName
	if vSubjectAltName == nil {
		// note: explicitly not the empty object.
		vSubjectAltName = &CertificateConfigSubjectConfigSubjectAltName{}
	}
	if err := extractCertificateConfigSubjectConfigSubjectAltNameFields(r, vSubjectAltName); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vSubjectAltName) {
		o.SubjectAltName = vSubjectAltName
	}
	return nil
}
func postReadExtractCertificateConfigSubjectConfigSubjectFields(r *Certificate, o *CertificateConfigSubjectConfigSubject) error {
	return nil
}
func postReadExtractCertificateConfigSubjectConfigSubjectAltNameFields(r *Certificate, o *CertificateConfigSubjectConfigSubjectAltName) error {
	return nil
}
func postReadExtractCertificateConfigX509ConfigFields(r *Certificate, o *CertificateConfigX509Config) error {
	vKeyUsage := o.KeyUsage
	if vKeyUsage == nil {
		// note: explicitly not the empty object.
		vKeyUsage = &CertificateConfigX509ConfigKeyUsage{}
	}
	if err := extractCertificateConfigX509ConfigKeyUsageFields(r, vKeyUsage); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vKeyUsage) {
		o.KeyUsage = vKeyUsage
	}
	vCaOptions := o.CaOptions
	if vCaOptions == nil {
		// note: explicitly not the empty object.
		vCaOptions = &CertificateConfigX509ConfigCaOptions{}
	}
	if err := extractCertificateConfigX509ConfigCaOptionsFields(r, vCaOptions); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vCaOptions) {
		o.CaOptions = vCaOptions
	}
	return nil
}
func postReadExtractCertificateConfigX509ConfigKeyUsageFields(r *Certificate, o *CertificateConfigX509ConfigKeyUsage) error {
	vBaseKeyUsage := o.BaseKeyUsage
	if vBaseKeyUsage == nil {
		// note: explicitly not the empty object.
		vBaseKeyUsage = &CertificateConfigX509ConfigKeyUsageBaseKeyUsage{}
	}
	if err := extractCertificateConfigX509ConfigKeyUsageBaseKeyUsageFields(r, vBaseKeyUsage); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vBaseKeyUsage) {
		o.BaseKeyUsage = vBaseKeyUsage
	}
	vExtendedKeyUsage := o.ExtendedKeyUsage
	if vExtendedKeyUsage == nil {
		// note: explicitly not the empty object.
		vExtendedKeyUsage = &CertificateConfigX509ConfigKeyUsageExtendedKeyUsage{}
	}
	if err := extractCertificateConfigX509ConfigKeyUsageExtendedKeyUsageFields(r, vExtendedKeyUsage); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vExtendedKeyUsage) {
		o.ExtendedKeyUsage = vExtendedKeyUsage
	}
	return nil
}
func postReadExtractCertificateConfigX509ConfigKeyUsageBaseKeyUsageFields(r *Certificate, o *CertificateConfigX509ConfigKeyUsageBaseKeyUsage) error {
	return nil
}
func postReadExtractCertificateConfigX509ConfigKeyUsageExtendedKeyUsageFields(r *Certificate, o *CertificateConfigX509ConfigKeyUsageExtendedKeyUsage) error {
	return nil
}
func postReadExtractCertificateConfigX509ConfigKeyUsageUnknownExtendedKeyUsagesFields(r *Certificate, o *CertificateConfigX509ConfigKeyUsageUnknownExtendedKeyUsages) error {
	return nil
}
func postReadExtractCertificateConfigX509ConfigCaOptionsFields(r *Certificate, o *CertificateConfigX509ConfigCaOptions) error {
	return nil
}
func postReadExtractCertificateConfigX509ConfigPolicyIdsFields(r *Certificate, o *CertificateConfigX509ConfigPolicyIds) error {
	return nil
}
func postReadExtractCertificateConfigX509ConfigAdditionalExtensionsFields(r *Certificate, o *CertificateConfigX509ConfigAdditionalExtensions) error {
	vObjectId := o.ObjectId
	if vObjectId == nil {
		// note: explicitly not the empty object.
		vObjectId = &CertificateConfigX509ConfigAdditionalExtensionsObjectId{}
	}
	if err := extractCertificateConfigX509ConfigAdditionalExtensionsObjectIdFields(r, vObjectId); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vObjectId) {
		o.ObjectId = vObjectId
	}
	return nil
}
func postReadExtractCertificateConfigX509ConfigAdditionalExtensionsObjectIdFields(r *Certificate, o *CertificateConfigX509ConfigAdditionalExtensionsObjectId) error {
	return nil
}
func postReadExtractCertificateConfigPublicKeyFields(r *Certificate, o *CertificateConfigPublicKey) error {
	return nil
}
func postReadExtractCertificateRevocationDetailsFields(r *Certificate, o *CertificateRevocationDetails) error {
	return nil
}
func postReadExtractCertificateCertificateDescriptionFields(r *Certificate, o *CertificateCertificateDescription) error {
	vSubjectDescription := o.SubjectDescription
	if vSubjectDescription == nil {
		// note: explicitly not the empty object.
		vSubjectDescription = &CertificateCertificateDescriptionSubjectDescription{}
	}
	if err := extractCertificateCertificateDescriptionSubjectDescriptionFields(r, vSubjectDescription); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vSubjectDescription) {
		o.SubjectDescription = vSubjectDescription
	}
	vX509Description := o.X509Description
	if vX509Description == nil {
		// note: explicitly not the empty object.
		vX509Description = &CertificateCertificateDescriptionX509Description{}
	}
	if err := extractCertificateCertificateDescriptionX509DescriptionFields(r, vX509Description); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vX509Description) {
		o.X509Description = vX509Description
	}
	vPublicKey := o.PublicKey
	if vPublicKey == nil {
		// note: explicitly not the empty object.
		vPublicKey = &CertificateCertificateDescriptionPublicKey{}
	}
	if err := extractCertificateCertificateDescriptionPublicKeyFields(r, vPublicKey); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vPublicKey) {
		o.PublicKey = vPublicKey
	}
	vSubjectKeyId := o.SubjectKeyId
	if vSubjectKeyId == nil {
		// note: explicitly not the empty object.
		vSubjectKeyId = &CertificateCertificateDescriptionSubjectKeyId{}
	}
	if err := extractCertificateCertificateDescriptionSubjectKeyIdFields(r, vSubjectKeyId); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vSubjectKeyId) {
		o.SubjectKeyId = vSubjectKeyId
	}
	vAuthorityKeyId := o.AuthorityKeyId
	if vAuthorityKeyId == nil {
		// note: explicitly not the empty object.
		vAuthorityKeyId = &CertificateCertificateDescriptionAuthorityKeyId{}
	}
	if err := extractCertificateCertificateDescriptionAuthorityKeyIdFields(r, vAuthorityKeyId); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vAuthorityKeyId) {
		o.AuthorityKeyId = vAuthorityKeyId
	}
	vCertFingerprint := o.CertFingerprint
	if vCertFingerprint == nil {
		// note: explicitly not the empty object.
		vCertFingerprint = &CertificateCertificateDescriptionCertFingerprint{}
	}
	if err := extractCertificateCertificateDescriptionCertFingerprintFields(r, vCertFingerprint); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vCertFingerprint) {
		o.CertFingerprint = vCertFingerprint
	}
	return nil
}
func postReadExtractCertificateCertificateDescriptionSubjectDescriptionFields(r *Certificate, o *CertificateCertificateDescriptionSubjectDescription) error {
	vSubject := o.Subject
	if vSubject == nil {
		// note: explicitly not the empty object.
		vSubject = &CertificateCertificateDescriptionSubjectDescriptionSubject{}
	}
	if err := extractCertificateCertificateDescriptionSubjectDescriptionSubjectFields(r, vSubject); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vSubject) {
		o.Subject = vSubject
	}
	vSubjectAltName := o.SubjectAltName
	if vSubjectAltName == nil {
		// note: explicitly not the empty object.
		vSubjectAltName = &CertificateCertificateDescriptionSubjectDescriptionSubjectAltName{}
	}
	if err := extractCertificateCertificateDescriptionSubjectDescriptionSubjectAltNameFields(r, vSubjectAltName); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vSubjectAltName) {
		o.SubjectAltName = vSubjectAltName
	}
	return nil
}
func postReadExtractCertificateCertificateDescriptionSubjectDescriptionSubjectFields(r *Certificate, o *CertificateCertificateDescriptionSubjectDescriptionSubject) error {
	return nil
}
func postReadExtractCertificateCertificateDescriptionSubjectDescriptionSubjectAltNameFields(r *Certificate, o *CertificateCertificateDescriptionSubjectDescriptionSubjectAltName) error {
	return nil
}
func postReadExtractCertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSansFields(r *Certificate, o *CertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSans) error {
	vObjectId := o.ObjectId
	if vObjectId == nil {
		// note: explicitly not the empty object.
		vObjectId = &CertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSansObjectId{}
	}
	if err := extractCertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSansObjectIdFields(r, vObjectId); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vObjectId) {
		o.ObjectId = vObjectId
	}
	return nil
}
func postReadExtractCertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSansObjectIdFields(r *Certificate, o *CertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSansObjectId) error {
	return nil
}
func postReadExtractCertificateCertificateDescriptionX509DescriptionFields(r *Certificate, o *CertificateCertificateDescriptionX509Description) error {
	vKeyUsage := o.KeyUsage
	if vKeyUsage == nil {
		// note: explicitly not the empty object.
		vKeyUsage = &CertificateCertificateDescriptionX509DescriptionKeyUsage{}
	}
	if err := extractCertificateCertificateDescriptionX509DescriptionKeyUsageFields(r, vKeyUsage); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vKeyUsage) {
		o.KeyUsage = vKeyUsage
	}
	vCaOptions := o.CaOptions
	if vCaOptions == nil {
		// note: explicitly not the empty object.
		vCaOptions = &CertificateCertificateDescriptionX509DescriptionCaOptions{}
	}
	if err := extractCertificateCertificateDescriptionX509DescriptionCaOptionsFields(r, vCaOptions); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vCaOptions) {
		o.CaOptions = vCaOptions
	}
	return nil
}
func postReadExtractCertificateCertificateDescriptionX509DescriptionKeyUsageFields(r *Certificate, o *CertificateCertificateDescriptionX509DescriptionKeyUsage) error {
	vBaseKeyUsage := o.BaseKeyUsage
	if vBaseKeyUsage == nil {
		// note: explicitly not the empty object.
		vBaseKeyUsage = &CertificateCertificateDescriptionX509DescriptionKeyUsageBaseKeyUsage{}
	}
	if err := extractCertificateCertificateDescriptionX509DescriptionKeyUsageBaseKeyUsageFields(r, vBaseKeyUsage); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vBaseKeyUsage) {
		o.BaseKeyUsage = vBaseKeyUsage
	}
	vExtendedKeyUsage := o.ExtendedKeyUsage
	if vExtendedKeyUsage == nil {
		// note: explicitly not the empty object.
		vExtendedKeyUsage = &CertificateCertificateDescriptionX509DescriptionKeyUsageExtendedKeyUsage{}
	}
	if err := extractCertificateCertificateDescriptionX509DescriptionKeyUsageExtendedKeyUsageFields(r, vExtendedKeyUsage); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vExtendedKeyUsage) {
		o.ExtendedKeyUsage = vExtendedKeyUsage
	}
	return nil
}
func postReadExtractCertificateCertificateDescriptionX509DescriptionKeyUsageBaseKeyUsageFields(r *Certificate, o *CertificateCertificateDescriptionX509DescriptionKeyUsageBaseKeyUsage) error {
	return nil
}
func postReadExtractCertificateCertificateDescriptionX509DescriptionKeyUsageExtendedKeyUsageFields(r *Certificate, o *CertificateCertificateDescriptionX509DescriptionKeyUsageExtendedKeyUsage) error {
	return nil
}
func postReadExtractCertificateCertificateDescriptionX509DescriptionKeyUsageUnknownExtendedKeyUsagesFields(r *Certificate, o *CertificateCertificateDescriptionX509DescriptionKeyUsageUnknownExtendedKeyUsages) error {
	return nil
}
func postReadExtractCertificateCertificateDescriptionX509DescriptionCaOptionsFields(r *Certificate, o *CertificateCertificateDescriptionX509DescriptionCaOptions) error {
	return nil
}
func postReadExtractCertificateCertificateDescriptionX509DescriptionPolicyIdsFields(r *Certificate, o *CertificateCertificateDescriptionX509DescriptionPolicyIds) error {
	return nil
}
func postReadExtractCertificateCertificateDescriptionX509DescriptionAdditionalExtensionsFields(r *Certificate, o *CertificateCertificateDescriptionX509DescriptionAdditionalExtensions) error {
	vObjectId := o.ObjectId
	if vObjectId == nil {
		// note: explicitly not the empty object.
		vObjectId = &CertificateCertificateDescriptionX509DescriptionAdditionalExtensionsObjectId{}
	}
	if err := extractCertificateCertificateDescriptionX509DescriptionAdditionalExtensionsObjectIdFields(r, vObjectId); err != nil {
		return err
	}
	if !dcl.IsEmptyValueIndirect(vObjectId) {
		o.ObjectId = vObjectId
	}
	return nil
}
func postReadExtractCertificateCertificateDescriptionX509DescriptionAdditionalExtensionsObjectIdFields(r *Certificate, o *CertificateCertificateDescriptionX509DescriptionAdditionalExtensionsObjectId) error {
	return nil
}
func postReadExtractCertificateCertificateDescriptionPublicKeyFields(r *Certificate, o *CertificateCertificateDescriptionPublicKey) error {
	return nil
}
func postReadExtractCertificateCertificateDescriptionSubjectKeyIdFields(r *Certificate, o *CertificateCertificateDescriptionSubjectKeyId) error {
	return nil
}
func postReadExtractCertificateCertificateDescriptionAuthorityKeyIdFields(r *Certificate, o *CertificateCertificateDescriptionAuthorityKeyId) error {
	return nil
}
func postReadExtractCertificateCertificateDescriptionCertFingerprintFields(r *Certificate, o *CertificateCertificateDescriptionCertFingerprint) error {
	return nil
}
