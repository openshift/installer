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
// Package privateca contains methods and objects for handling privateca GCP resources.
package privateca

import (
	"bytes"
	"context"
	"fmt"

	"github.com/GoogleCloudPlatform/declarative-resource-client-library/dcl"
	"github.com/GoogleCloudPlatform/declarative-resource-client-library/dcl/operations"
)

// Waits for the state of the certificate authority to be "ENABLED".
func (r *CertificateAuthority) waitForCertificateAuthorityEnabled(ctx context.Context, c *Client) error {
	return dcl.Do(ctx, func(ctx context.Context) (*dcl.RetryDetails, error) {
		nr, err := c.GetCertificateAuthority(ctx, r)
		if err != nil {
			return nil, err
		}
		if *nr.State == *CertificateAuthorityStateEnumRef("ENABLED") || *nr.State == *CertificateAuthorityStateEnumRef("STAGED") {
			return nil, nil
		}
		return &dcl.RetryDetails{}, dcl.OperationNotDone{}
	}, c.Config.RetryProvider)
}

// Disables the certificate authority so that it can be deleted.
func (r *CertificateAuthority) disableCertificateAuthority(ctx context.Context, c *Client) error {
	if *r.State != *CertificateAuthorityStateEnumRef("ENABLED") {
		// Only enabled certificate authorities need to be disabled before deletion.
		return nil
	}
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"project":  dcl.ValueOrEmptyString(nr.Project),
		"location": dcl.ValueOrEmptyString(nr.Location),
		"ca_pool":  dcl.ValueOrEmptyString(nr.CaPool),
		"name":     dcl.ValueOrEmptyString(nr.Name),
	}
	u := dcl.URL("projects/{{project}}/locations/{{location}}/caPools/{{ca_pool}}/certificateAuthorities/{{name}}:disable", "https://privateca.googleapis.com/v1/", c.Config.BasePath, params)
	resp, err := dcl.SendRequest(ctx, c.Config, "POST", u, &bytes.Buffer{}, c.Config.RetryProvider)
	if err != nil {
		return err
	}

	// wait for certificate authority to be disabled.
	var o operations.StandardGCPOperation
	if err := dcl.ParseResponse(resp.Response, &o); err != nil {
		return err
	}
	if err := o.Wait(ctx, c.Config, "https://privateca.googleapis.com/v1beta1/", "GET"); err != nil {
		return err
	}
	return nil
}

func (r *Certificate) createURL(userBasePath string) (string, error) {
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"project":              dcl.ValueOrEmptyString(nr.Project),
		"location":             dcl.ValueOrEmptyString(nr.Location),
		"caPool":               dcl.ValueOrEmptyString(nr.CaPool),
		"name":                 dcl.ValueOrEmptyString(nr.Name),
		"certificateAuthority": dcl.ValueOrEmptyString(nr.CertificateAuthority),
	}
	basePath := dcl.URL("projects/{{project}}/locations/{{location}}/caPools/{{caPool}}/certificates", nr.basePath(), userBasePath, params)

	var err error
	if dcl.ValueOrEmptyString(nr.CertificateAuthority) != "" {
		basePath, err = dcl.AddQueryParams(basePath, map[string]string{"issuingCertificateAuthorityId": dcl.ValueOrEmptyString(nr.CertificateAuthority)})
		if err != nil {
			return "", err
		}
	}
	if dcl.ValueOrEmptyString(nr.Name) != "" {
		// Need to set name to nil or else it'll trigger error messages.
		basePath, err = dcl.AddQueryParams(basePath, map[string]string{"certificateId": dcl.ValueOrEmptyString(nr.Name)})
		if err != nil {
			return "", err
		}
	}
	return basePath, nil
}

func flattenCertificateConfigX509ConfigCAOptions(_ *Client, i interface{}, _ *Certificate) *CertificateConfigX509ConfigCaOptions {
	if i == nil {
		return nil
	}
	m, ok := i.(map[string]interface{})
	if !ok {
		return nil
	}

	result := &CertificateConfigX509ConfigCaOptions{}

	isCA, ok := m["isCa"].(bool)
	if ok {
		result.IsCa = dcl.Bool(isCA)
		if !isCA {
			result.NonCa = dcl.Bool(true)
		}
	}

	if _, ok := m["maxIssuerPathLength"]; ok {
		pathLen := dcl.FlattenInteger(m["maxIssuerPathLength"])
		result.MaxIssuerPathLength = pathLen
		if dcl.ValueOrEmptyInt64(pathLen) == 0 {
			result.ZeroMaxIssuerPathLength = dcl.Bool(true)
		}
	}

	return result
}

func expandCertificateConfigX509ConfigCAOptions(_ *Client, caOptions *CertificateConfigX509ConfigCaOptions, _ *Certificate) (map[string]interface{}, error) {
	if caOptions == nil {
		return nil, nil
	}

	m := make(map[string]interface{})
	isCA := dcl.ValueOrEmptyBool(caOptions.IsCa)
	nonCA := dcl.ValueOrEmptyBool(caOptions.NonCa)
	zeroPathLength := dcl.ValueOrEmptyBool(caOptions.ZeroMaxIssuerPathLength)
	maxIssuerPathLength := dcl.ValueOrEmptyInt64(caOptions.MaxIssuerPathLength)

	if !isCA && !nonCA {
		return nil, nil
	} else if isCA && nonCA {
		return nil, fmt.Errorf("is_ca and non_ca are mutually exclusive")
	} else if isCA || nonCA {
		m["isCa"] = isCA
	}

	if zeroPathLength && maxIssuerPathLength > 0 {
		return nil, fmt.Errorf("max_issuer_path_length and zero_max_issuer_path_length are mutually exclusive")
	}
	if maxIssuerPathLength > 0 || zeroPathLength {
		m["maxIssuerPathLength"] = maxIssuerPathLength
	}

	return m, nil
}

// base_key_usage has a custom flattener because the API does not return the object when all subfields are set to false.
func flattenCertificateTemplateBaseKeyUsage(_ *Client, i interface{}, res *CertificateTemplate) *CertificateTemplatePredefinedValuesKeyUsageBaseKeyUsage {
	m, ok := i.(map[string]interface{})
	if !ok {
		if res != nil && res.PredefinedValues != nil && res.PredefinedValues.KeyUsage != nil && res.PredefinedValues.KeyUsage.BaseKeyUsage != nil {
			baseKeyUsage := res.PredefinedValues.KeyUsage.BaseKeyUsage
			allFalse := true
			for _, booleanField := range []*bool{
				baseKeyUsage.DigitalSignature,
				baseKeyUsage.ContentCommitment,
				baseKeyUsage.KeyEncipherment,
				baseKeyUsage.DataEncipherment,
				baseKeyUsage.KeyAgreement,
				baseKeyUsage.CertSign,
				baseKeyUsage.CrlSign,
				baseKeyUsage.EncipherOnly,
				baseKeyUsage.DecipherOnly,
			} {
				if dcl.ValueOrEmptyBool(booleanField) {
					allFalse = false
				}
			}
			if allFalse {
				return dcl.Copy(baseKeyUsage).(*CertificateTemplatePredefinedValuesKeyUsageBaseKeyUsage)
			}
		}
		return nil
	}

	r := &CertificateTemplatePredefinedValuesKeyUsageBaseKeyUsage{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyCertificateTemplatePredefinedValuesKeyUsageBaseKeyUsage
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

// extended_key_usage has a custom flattener because the API does not return the object when all subfields are set to false.
func flattenCertificateTemplateExtendedKeyUsage(_ *Client, i interface{}, res *CertificateTemplate) *CertificateTemplatePredefinedValuesKeyUsageExtendedKeyUsage {
	m, ok := i.(map[string]interface{})
	if !ok {
		if res != nil && res.PredefinedValues != nil && res.PredefinedValues.KeyUsage != nil && res.PredefinedValues.KeyUsage.ExtendedKeyUsage != nil {
			extendedKeyUsage := res.PredefinedValues.KeyUsage.ExtendedKeyUsage
			allFalse := true
			for _, booleanField := range []*bool{
				extendedKeyUsage.ServerAuth,
				extendedKeyUsage.ClientAuth,
				extendedKeyUsage.CodeSigning,
				extendedKeyUsage.EmailProtection,
				extendedKeyUsage.TimeStamping,
				extendedKeyUsage.OcspSigning,
			} {
				if dcl.ValueOrEmptyBool(booleanField) {
					allFalse = false
				}
			}
			if allFalse {
				return dcl.Copy(extendedKeyUsage).(*CertificateTemplatePredefinedValuesKeyUsageExtendedKeyUsage)
			}
		}
		return nil
	}

	r := &CertificateTemplatePredefinedValuesKeyUsageExtendedKeyUsage{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyCertificateTemplatePredefinedValuesKeyUsageExtendedKeyUsage
	}
	r.ServerAuth = dcl.FlattenBool(m["serverAuth"])
	r.ClientAuth = dcl.FlattenBool(m["clientAuth"])
	r.CodeSigning = dcl.FlattenBool(m["codeSigning"])
	r.EmailProtection = dcl.FlattenBool(m["emailProtection"])
	r.TimeStamping = dcl.FlattenBool(m["timeStamping"])
	r.OcspSigning = dcl.FlattenBool(m["ocspSigning"])

	return r
}

func expandCaPoolIssuancePolicyBaselineValuesCAOptions(_ *Client, f *CaPoolIssuancePolicyBaselineValuesCaOptions, res *CaPool) (map[string]any, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]any)
	if v := f.IsCa; !dcl.IsEmptyValueIndirect(v) {
		m["isCa"] = v
	}

	maxIssuerPathLength := dcl.ValueOrEmptyInt64(f.MaxIssuerPathLength)
	zeroPathLength := dcl.ValueOrEmptyBool(f.ZeroMaxIssuerPathLength)
	if zeroPathLength && maxIssuerPathLength > 0 {
		return nil, fmt.Errorf("max_issuer_path_length and zero_max_issuer_path_length are mutually exclusive")
	}
	if maxIssuerPathLength > 0 || zeroPathLength {
		m["maxIssuerPathLength"] = maxIssuerPathLength
	}

	return m, nil
}

func flattenCaPoolIssuancePolicyBaselineValuesCAOptions(_ *Client, i any, res *CaPool) *CaPoolIssuancePolicyBaselineValuesCaOptions {
	m, ok := i.(map[string]any)
	if !ok {
		return nil
	}

	r := &CaPoolIssuancePolicyBaselineValuesCaOptions{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyCaPoolIssuancePolicyBaselineValuesCaOptions
	}

	isCA, ok := m["isCa"].(bool)
	if ok {
		r.IsCa = dcl.Bool(isCA)
	}

	if _, ok := m["maxIssuerPathLength"]; ok {
		pathLen := dcl.FlattenInteger(m["maxIssuerPathLength"])
		r.MaxIssuerPathLength = pathLen
		if dcl.ValueOrEmptyInt64(pathLen) == 0 {
			r.ZeroMaxIssuerPathLength = dcl.Bool(true)
		}
	}

	return r
}

func flattenCertificateAuthorityConfigX509ConfigCAOptions(_ *Client, i any, res *CertificateAuthority) *CertificateAuthorityConfigX509ConfigCaOptions {
	m, ok := i.(map[string]any)
	if !ok {
		return nil
	}

	r := &CertificateAuthorityConfigX509ConfigCaOptions{}

	if dcl.IsEmptyValueIndirect(i) {
		return EmptyCertificateAuthorityConfigX509ConfigCaOptions
	}

	isCA, ok := m["isCa"].(bool)
	if ok {
		r.IsCa = dcl.Bool(isCA)
	}

	if _, ok := m["maxIssuerPathLength"]; ok {
		pathLen := dcl.FlattenInteger(m["maxIssuerPathLength"])
		r.MaxIssuerPathLength = pathLen
		if dcl.ValueOrEmptyInt64(pathLen) == 0 {
			r.ZeroMaxIssuerPathLength = dcl.Bool(true)
		}
	}

	return r
}

func expandCertificateAuthorityConfigX509ConfigCAOptions(_ *Client, f *CertificateAuthorityConfigX509ConfigCaOptions, res *CertificateAuthority) (map[string]any, error) {
	if dcl.IsEmptyValueIndirect(f) {
		return nil, nil
	}

	m := make(map[string]any)
	if v := f.IsCa; !dcl.IsEmptyValueIndirect(v) {
		m["isCa"] = v
	}

	maxIssuerPathLength := dcl.ValueOrEmptyInt64(f.MaxIssuerPathLength)
	zeroPathLength := dcl.ValueOrEmptyBool(f.ZeroMaxIssuerPathLength)
	if zeroPathLength && maxIssuerPathLength > 0 {
		return nil, fmt.Errorf("max_issuer_path_length and zero_max_issuer_path_length are mutually exclusive")
	}
	if maxIssuerPathLength > 0 || zeroPathLength {
		m["maxIssuerPathLength"] = maxIssuerPathLength
	}

	return m, nil
}
