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
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"time"

	"google.golang.org/api/googleapi"
	"github.com/GoogleCloudPlatform/declarative-resource-client-library/dcl"
)

type CaPool struct {
	Name              *string                  `json:"name"`
	Tier              *CaPoolTierEnum          `json:"tier"`
	IssuancePolicy    *CaPoolIssuancePolicy    `json:"issuancePolicy"`
	PublishingOptions *CaPoolPublishingOptions `json:"publishingOptions"`
	Labels            map[string]string        `json:"labels"`
	Project           *string                  `json:"project"`
	Location          *string                  `json:"location"`
}

func (r *CaPool) String() string {
	return dcl.SprintResource(r)
}

// The enum CaPoolTierEnum.
type CaPoolTierEnum string

// CaPoolTierEnumRef returns a *CaPoolTierEnum with the value of string s
// If the empty string is provided, nil is returned.
func CaPoolTierEnumRef(s string) *CaPoolTierEnum {
	v := CaPoolTierEnum(s)
	return &v
}

func (v CaPoolTierEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"TIER_UNSPECIFIED", "ENTERPRISE", "DEVOPS"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "CaPoolTierEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum CaPoolIssuancePolicyAllowedKeyTypesEllipticCurveSignatureAlgorithmEnum.
type CaPoolIssuancePolicyAllowedKeyTypesEllipticCurveSignatureAlgorithmEnum string

// CaPoolIssuancePolicyAllowedKeyTypesEllipticCurveSignatureAlgorithmEnumRef returns a *CaPoolIssuancePolicyAllowedKeyTypesEllipticCurveSignatureAlgorithmEnum with the value of string s
// If the empty string is provided, nil is returned.
func CaPoolIssuancePolicyAllowedKeyTypesEllipticCurveSignatureAlgorithmEnumRef(s string) *CaPoolIssuancePolicyAllowedKeyTypesEllipticCurveSignatureAlgorithmEnum {
	v := CaPoolIssuancePolicyAllowedKeyTypesEllipticCurveSignatureAlgorithmEnum(s)
	return &v
}

func (v CaPoolIssuancePolicyAllowedKeyTypesEllipticCurveSignatureAlgorithmEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"EC_SIGNATURE_ALGORITHM_UNSPECIFIED", "ECDSA_P256", "ECDSA_P384", "EDDSA_25519"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "CaPoolIssuancePolicyAllowedKeyTypesEllipticCurveSignatureAlgorithmEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum CaPoolIssuancePolicyPassthroughExtensionsKnownExtensionsEnum.
type CaPoolIssuancePolicyPassthroughExtensionsKnownExtensionsEnum string

// CaPoolIssuancePolicyPassthroughExtensionsKnownExtensionsEnumRef returns a *CaPoolIssuancePolicyPassthroughExtensionsKnownExtensionsEnum with the value of string s
// If the empty string is provided, nil is returned.
func CaPoolIssuancePolicyPassthroughExtensionsKnownExtensionsEnumRef(s string) *CaPoolIssuancePolicyPassthroughExtensionsKnownExtensionsEnum {
	v := CaPoolIssuancePolicyPassthroughExtensionsKnownExtensionsEnum(s)
	return &v
}

func (v CaPoolIssuancePolicyPassthroughExtensionsKnownExtensionsEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"KNOWN_CERTIFICATE_EXTENSION_UNSPECIFIED", "BASE_KEY_USAGE", "EXTENDED_KEY_USAGE", "CA_OPTIONS", "POLICY_IDS", "AIA_OCSP_SERVERS"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "CaPoolIssuancePolicyPassthroughExtensionsKnownExtensionsEnum",
		Value: string(v),
		Valid: []string{},
	}
}

type CaPoolIssuancePolicy struct {
	empty                 bool                                       `json:"-"`
	AllowedKeyTypes       []CaPoolIssuancePolicyAllowedKeyTypes      `json:"allowedKeyTypes"`
	MaximumLifetime       *string                                    `json:"maximumLifetime"`
	AllowedIssuanceModes  *CaPoolIssuancePolicyAllowedIssuanceModes  `json:"allowedIssuanceModes"`
	BaselineValues        *CaPoolIssuancePolicyBaselineValues        `json:"baselineValues"`
	IdentityConstraints   *CaPoolIssuancePolicyIdentityConstraints   `json:"identityConstraints"`
	PassthroughExtensions *CaPoolIssuancePolicyPassthroughExtensions `json:"passthroughExtensions"`
}

type jsonCaPoolIssuancePolicy CaPoolIssuancePolicy

func (r *CaPoolIssuancePolicy) UnmarshalJSON(data []byte) error {
	var res jsonCaPoolIssuancePolicy
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyCaPoolIssuancePolicy
	} else {

		r.AllowedKeyTypes = res.AllowedKeyTypes

		r.MaximumLifetime = res.MaximumLifetime

		r.AllowedIssuanceModes = res.AllowedIssuanceModes

		r.BaselineValues = res.BaselineValues

		r.IdentityConstraints = res.IdentityConstraints

		r.PassthroughExtensions = res.PassthroughExtensions

	}
	return nil
}

// This object is used to assert a desired state where this CaPoolIssuancePolicy is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyCaPoolIssuancePolicy *CaPoolIssuancePolicy = &CaPoolIssuancePolicy{empty: true}

func (r *CaPoolIssuancePolicy) Empty() bool {
	return r.empty
}

func (r *CaPoolIssuancePolicy) String() string {
	return dcl.SprintResource(r)
}

func (r *CaPoolIssuancePolicy) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type CaPoolIssuancePolicyAllowedKeyTypes struct {
	empty         bool                                              `json:"-"`
	Rsa           *CaPoolIssuancePolicyAllowedKeyTypesRsa           `json:"rsa"`
	EllipticCurve *CaPoolIssuancePolicyAllowedKeyTypesEllipticCurve `json:"ellipticCurve"`
}

type jsonCaPoolIssuancePolicyAllowedKeyTypes CaPoolIssuancePolicyAllowedKeyTypes

func (r *CaPoolIssuancePolicyAllowedKeyTypes) UnmarshalJSON(data []byte) error {
	var res jsonCaPoolIssuancePolicyAllowedKeyTypes
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyCaPoolIssuancePolicyAllowedKeyTypes
	} else {

		r.Rsa = res.Rsa

		r.EllipticCurve = res.EllipticCurve

	}
	return nil
}

// This object is used to assert a desired state where this CaPoolIssuancePolicyAllowedKeyTypes is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyCaPoolIssuancePolicyAllowedKeyTypes *CaPoolIssuancePolicyAllowedKeyTypes = &CaPoolIssuancePolicyAllowedKeyTypes{empty: true}

func (r *CaPoolIssuancePolicyAllowedKeyTypes) Empty() bool {
	return r.empty
}

func (r *CaPoolIssuancePolicyAllowedKeyTypes) String() string {
	return dcl.SprintResource(r)
}

func (r *CaPoolIssuancePolicyAllowedKeyTypes) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type CaPoolIssuancePolicyAllowedKeyTypesRsa struct {
	empty          bool   `json:"-"`
	MinModulusSize *int64 `json:"minModulusSize"`
	MaxModulusSize *int64 `json:"maxModulusSize"`
}

type jsonCaPoolIssuancePolicyAllowedKeyTypesRsa CaPoolIssuancePolicyAllowedKeyTypesRsa

func (r *CaPoolIssuancePolicyAllowedKeyTypesRsa) UnmarshalJSON(data []byte) error {
	var res jsonCaPoolIssuancePolicyAllowedKeyTypesRsa
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyCaPoolIssuancePolicyAllowedKeyTypesRsa
	} else {

		r.MinModulusSize = res.MinModulusSize

		r.MaxModulusSize = res.MaxModulusSize

	}
	return nil
}

// This object is used to assert a desired state where this CaPoolIssuancePolicyAllowedKeyTypesRsa is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyCaPoolIssuancePolicyAllowedKeyTypesRsa *CaPoolIssuancePolicyAllowedKeyTypesRsa = &CaPoolIssuancePolicyAllowedKeyTypesRsa{empty: true}

func (r *CaPoolIssuancePolicyAllowedKeyTypesRsa) Empty() bool {
	return r.empty
}

func (r *CaPoolIssuancePolicyAllowedKeyTypesRsa) String() string {
	return dcl.SprintResource(r)
}

func (r *CaPoolIssuancePolicyAllowedKeyTypesRsa) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type CaPoolIssuancePolicyAllowedKeyTypesEllipticCurve struct {
	empty              bool                                                                    `json:"-"`
	SignatureAlgorithm *CaPoolIssuancePolicyAllowedKeyTypesEllipticCurveSignatureAlgorithmEnum `json:"signatureAlgorithm"`
}

type jsonCaPoolIssuancePolicyAllowedKeyTypesEllipticCurve CaPoolIssuancePolicyAllowedKeyTypesEllipticCurve

func (r *CaPoolIssuancePolicyAllowedKeyTypesEllipticCurve) UnmarshalJSON(data []byte) error {
	var res jsonCaPoolIssuancePolicyAllowedKeyTypesEllipticCurve
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyCaPoolIssuancePolicyAllowedKeyTypesEllipticCurve
	} else {

		r.SignatureAlgorithm = res.SignatureAlgorithm

	}
	return nil
}

// This object is used to assert a desired state where this CaPoolIssuancePolicyAllowedKeyTypesEllipticCurve is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyCaPoolIssuancePolicyAllowedKeyTypesEllipticCurve *CaPoolIssuancePolicyAllowedKeyTypesEllipticCurve = &CaPoolIssuancePolicyAllowedKeyTypesEllipticCurve{empty: true}

func (r *CaPoolIssuancePolicyAllowedKeyTypesEllipticCurve) Empty() bool {
	return r.empty
}

func (r *CaPoolIssuancePolicyAllowedKeyTypesEllipticCurve) String() string {
	return dcl.SprintResource(r)
}

func (r *CaPoolIssuancePolicyAllowedKeyTypesEllipticCurve) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type CaPoolIssuancePolicyAllowedIssuanceModes struct {
	empty                    bool  `json:"-"`
	AllowCsrBasedIssuance    *bool `json:"allowCsrBasedIssuance"`
	AllowConfigBasedIssuance *bool `json:"allowConfigBasedIssuance"`
}

type jsonCaPoolIssuancePolicyAllowedIssuanceModes CaPoolIssuancePolicyAllowedIssuanceModes

func (r *CaPoolIssuancePolicyAllowedIssuanceModes) UnmarshalJSON(data []byte) error {
	var res jsonCaPoolIssuancePolicyAllowedIssuanceModes
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyCaPoolIssuancePolicyAllowedIssuanceModes
	} else {

		r.AllowCsrBasedIssuance = res.AllowCsrBasedIssuance

		r.AllowConfigBasedIssuance = res.AllowConfigBasedIssuance

	}
	return nil
}

// This object is used to assert a desired state where this CaPoolIssuancePolicyAllowedIssuanceModes is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyCaPoolIssuancePolicyAllowedIssuanceModes *CaPoolIssuancePolicyAllowedIssuanceModes = &CaPoolIssuancePolicyAllowedIssuanceModes{empty: true}

func (r *CaPoolIssuancePolicyAllowedIssuanceModes) Empty() bool {
	return r.empty
}

func (r *CaPoolIssuancePolicyAllowedIssuanceModes) String() string {
	return dcl.SprintResource(r)
}

func (r *CaPoolIssuancePolicyAllowedIssuanceModes) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type CaPoolIssuancePolicyBaselineValues struct {
	empty                bool                                                     `json:"-"`
	KeyUsage             *CaPoolIssuancePolicyBaselineValuesKeyUsage              `json:"keyUsage"`
	CaOptions            *CaPoolIssuancePolicyBaselineValuesCaOptions             `json:"caOptions"`
	PolicyIds            []CaPoolIssuancePolicyBaselineValuesPolicyIds            `json:"policyIds"`
	AiaOcspServers       []string                                                 `json:"aiaOcspServers"`
	AdditionalExtensions []CaPoolIssuancePolicyBaselineValuesAdditionalExtensions `json:"additionalExtensions"`
}

type jsonCaPoolIssuancePolicyBaselineValues CaPoolIssuancePolicyBaselineValues

func (r *CaPoolIssuancePolicyBaselineValues) UnmarshalJSON(data []byte) error {
	var res jsonCaPoolIssuancePolicyBaselineValues
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyCaPoolIssuancePolicyBaselineValues
	} else {

		r.KeyUsage = res.KeyUsage

		r.CaOptions = res.CaOptions

		r.PolicyIds = res.PolicyIds

		r.AiaOcspServers = res.AiaOcspServers

		r.AdditionalExtensions = res.AdditionalExtensions

	}
	return nil
}

// This object is used to assert a desired state where this CaPoolIssuancePolicyBaselineValues is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyCaPoolIssuancePolicyBaselineValues *CaPoolIssuancePolicyBaselineValues = &CaPoolIssuancePolicyBaselineValues{empty: true}

func (r *CaPoolIssuancePolicyBaselineValues) Empty() bool {
	return r.empty
}

func (r *CaPoolIssuancePolicyBaselineValues) String() string {
	return dcl.SprintResource(r)
}

func (r *CaPoolIssuancePolicyBaselineValues) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type CaPoolIssuancePolicyBaselineValuesKeyUsage struct {
	empty                    bool                                                                 `json:"-"`
	BaseKeyUsage             *CaPoolIssuancePolicyBaselineValuesKeyUsageBaseKeyUsage              `json:"baseKeyUsage"`
	ExtendedKeyUsage         *CaPoolIssuancePolicyBaselineValuesKeyUsageExtendedKeyUsage          `json:"extendedKeyUsage"`
	UnknownExtendedKeyUsages []CaPoolIssuancePolicyBaselineValuesKeyUsageUnknownExtendedKeyUsages `json:"unknownExtendedKeyUsages"`
}

type jsonCaPoolIssuancePolicyBaselineValuesKeyUsage CaPoolIssuancePolicyBaselineValuesKeyUsage

func (r *CaPoolIssuancePolicyBaselineValuesKeyUsage) UnmarshalJSON(data []byte) error {
	var res jsonCaPoolIssuancePolicyBaselineValuesKeyUsage
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyCaPoolIssuancePolicyBaselineValuesKeyUsage
	} else {

		r.BaseKeyUsage = res.BaseKeyUsage

		r.ExtendedKeyUsage = res.ExtendedKeyUsage

		r.UnknownExtendedKeyUsages = res.UnknownExtendedKeyUsages

	}
	return nil
}

// This object is used to assert a desired state where this CaPoolIssuancePolicyBaselineValuesKeyUsage is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyCaPoolIssuancePolicyBaselineValuesKeyUsage *CaPoolIssuancePolicyBaselineValuesKeyUsage = &CaPoolIssuancePolicyBaselineValuesKeyUsage{empty: true}

func (r *CaPoolIssuancePolicyBaselineValuesKeyUsage) Empty() bool {
	return r.empty
}

func (r *CaPoolIssuancePolicyBaselineValuesKeyUsage) String() string {
	return dcl.SprintResource(r)
}

func (r *CaPoolIssuancePolicyBaselineValuesKeyUsage) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type CaPoolIssuancePolicyBaselineValuesKeyUsageBaseKeyUsage struct {
	empty             bool  `json:"-"`
	DigitalSignature  *bool `json:"digitalSignature"`
	ContentCommitment *bool `json:"contentCommitment"`
	KeyEncipherment   *bool `json:"keyEncipherment"`
	DataEncipherment  *bool `json:"dataEncipherment"`
	KeyAgreement      *bool `json:"keyAgreement"`
	CertSign          *bool `json:"certSign"`
	CrlSign           *bool `json:"crlSign"`
	EncipherOnly      *bool `json:"encipherOnly"`
	DecipherOnly      *bool `json:"decipherOnly"`
}

type jsonCaPoolIssuancePolicyBaselineValuesKeyUsageBaseKeyUsage CaPoolIssuancePolicyBaselineValuesKeyUsageBaseKeyUsage

func (r *CaPoolIssuancePolicyBaselineValuesKeyUsageBaseKeyUsage) UnmarshalJSON(data []byte) error {
	var res jsonCaPoolIssuancePolicyBaselineValuesKeyUsageBaseKeyUsage
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyCaPoolIssuancePolicyBaselineValuesKeyUsageBaseKeyUsage
	} else {

		r.DigitalSignature = res.DigitalSignature

		r.ContentCommitment = res.ContentCommitment

		r.KeyEncipherment = res.KeyEncipherment

		r.DataEncipherment = res.DataEncipherment

		r.KeyAgreement = res.KeyAgreement

		r.CertSign = res.CertSign

		r.CrlSign = res.CrlSign

		r.EncipherOnly = res.EncipherOnly

		r.DecipherOnly = res.DecipherOnly

	}
	return nil
}

// This object is used to assert a desired state where this CaPoolIssuancePolicyBaselineValuesKeyUsageBaseKeyUsage is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyCaPoolIssuancePolicyBaselineValuesKeyUsageBaseKeyUsage *CaPoolIssuancePolicyBaselineValuesKeyUsageBaseKeyUsage = &CaPoolIssuancePolicyBaselineValuesKeyUsageBaseKeyUsage{empty: true}

func (r *CaPoolIssuancePolicyBaselineValuesKeyUsageBaseKeyUsage) Empty() bool {
	return r.empty
}

func (r *CaPoolIssuancePolicyBaselineValuesKeyUsageBaseKeyUsage) String() string {
	return dcl.SprintResource(r)
}

func (r *CaPoolIssuancePolicyBaselineValuesKeyUsageBaseKeyUsage) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type CaPoolIssuancePolicyBaselineValuesKeyUsageExtendedKeyUsage struct {
	empty           bool  `json:"-"`
	ServerAuth      *bool `json:"serverAuth"`
	ClientAuth      *bool `json:"clientAuth"`
	CodeSigning     *bool `json:"codeSigning"`
	EmailProtection *bool `json:"emailProtection"`
	TimeStamping    *bool `json:"timeStamping"`
	OcspSigning     *bool `json:"ocspSigning"`
}

type jsonCaPoolIssuancePolicyBaselineValuesKeyUsageExtendedKeyUsage CaPoolIssuancePolicyBaselineValuesKeyUsageExtendedKeyUsage

func (r *CaPoolIssuancePolicyBaselineValuesKeyUsageExtendedKeyUsage) UnmarshalJSON(data []byte) error {
	var res jsonCaPoolIssuancePolicyBaselineValuesKeyUsageExtendedKeyUsage
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyCaPoolIssuancePolicyBaselineValuesKeyUsageExtendedKeyUsage
	} else {

		r.ServerAuth = res.ServerAuth

		r.ClientAuth = res.ClientAuth

		r.CodeSigning = res.CodeSigning

		r.EmailProtection = res.EmailProtection

		r.TimeStamping = res.TimeStamping

		r.OcspSigning = res.OcspSigning

	}
	return nil
}

// This object is used to assert a desired state where this CaPoolIssuancePolicyBaselineValuesKeyUsageExtendedKeyUsage is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyCaPoolIssuancePolicyBaselineValuesKeyUsageExtendedKeyUsage *CaPoolIssuancePolicyBaselineValuesKeyUsageExtendedKeyUsage = &CaPoolIssuancePolicyBaselineValuesKeyUsageExtendedKeyUsage{empty: true}

func (r *CaPoolIssuancePolicyBaselineValuesKeyUsageExtendedKeyUsage) Empty() bool {
	return r.empty
}

func (r *CaPoolIssuancePolicyBaselineValuesKeyUsageExtendedKeyUsage) String() string {
	return dcl.SprintResource(r)
}

func (r *CaPoolIssuancePolicyBaselineValuesKeyUsageExtendedKeyUsage) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type CaPoolIssuancePolicyBaselineValuesKeyUsageUnknownExtendedKeyUsages struct {
	empty        bool    `json:"-"`
	ObjectIdPath []int64 `json:"objectIdPath"`
}

type jsonCaPoolIssuancePolicyBaselineValuesKeyUsageUnknownExtendedKeyUsages CaPoolIssuancePolicyBaselineValuesKeyUsageUnknownExtendedKeyUsages

func (r *CaPoolIssuancePolicyBaselineValuesKeyUsageUnknownExtendedKeyUsages) UnmarshalJSON(data []byte) error {
	var res jsonCaPoolIssuancePolicyBaselineValuesKeyUsageUnknownExtendedKeyUsages
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyCaPoolIssuancePolicyBaselineValuesKeyUsageUnknownExtendedKeyUsages
	} else {

		r.ObjectIdPath = res.ObjectIdPath

	}
	return nil
}

// This object is used to assert a desired state where this CaPoolIssuancePolicyBaselineValuesKeyUsageUnknownExtendedKeyUsages is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyCaPoolIssuancePolicyBaselineValuesKeyUsageUnknownExtendedKeyUsages *CaPoolIssuancePolicyBaselineValuesKeyUsageUnknownExtendedKeyUsages = &CaPoolIssuancePolicyBaselineValuesKeyUsageUnknownExtendedKeyUsages{empty: true}

func (r *CaPoolIssuancePolicyBaselineValuesKeyUsageUnknownExtendedKeyUsages) Empty() bool {
	return r.empty
}

func (r *CaPoolIssuancePolicyBaselineValuesKeyUsageUnknownExtendedKeyUsages) String() string {
	return dcl.SprintResource(r)
}

func (r *CaPoolIssuancePolicyBaselineValuesKeyUsageUnknownExtendedKeyUsages) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type CaPoolIssuancePolicyBaselineValuesCaOptions struct {
	empty                   bool   `json:"-"`
	IsCa                    *bool  `json:"isCa"`
	MaxIssuerPathLength     *int64 `json:"maxIssuerPathLength"`
	ZeroMaxIssuerPathLength *bool  `json:"zeroMaxIssuerPathLength"`
}

type jsonCaPoolIssuancePolicyBaselineValuesCaOptions CaPoolIssuancePolicyBaselineValuesCaOptions

func (r *CaPoolIssuancePolicyBaselineValuesCaOptions) UnmarshalJSON(data []byte) error {
	var res jsonCaPoolIssuancePolicyBaselineValuesCaOptions
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyCaPoolIssuancePolicyBaselineValuesCaOptions
	} else {

		r.IsCa = res.IsCa

		r.MaxIssuerPathLength = res.MaxIssuerPathLength

		r.ZeroMaxIssuerPathLength = res.ZeroMaxIssuerPathLength

	}
	return nil
}

// This object is used to assert a desired state where this CaPoolIssuancePolicyBaselineValuesCaOptions is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyCaPoolIssuancePolicyBaselineValuesCaOptions *CaPoolIssuancePolicyBaselineValuesCaOptions = &CaPoolIssuancePolicyBaselineValuesCaOptions{empty: true}

func (r *CaPoolIssuancePolicyBaselineValuesCaOptions) Empty() bool {
	return r.empty
}

func (r *CaPoolIssuancePolicyBaselineValuesCaOptions) String() string {
	return dcl.SprintResource(r)
}

func (r *CaPoolIssuancePolicyBaselineValuesCaOptions) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type CaPoolIssuancePolicyBaselineValuesPolicyIds struct {
	empty        bool    `json:"-"`
	ObjectIdPath []int64 `json:"objectIdPath"`
}

type jsonCaPoolIssuancePolicyBaselineValuesPolicyIds CaPoolIssuancePolicyBaselineValuesPolicyIds

func (r *CaPoolIssuancePolicyBaselineValuesPolicyIds) UnmarshalJSON(data []byte) error {
	var res jsonCaPoolIssuancePolicyBaselineValuesPolicyIds
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyCaPoolIssuancePolicyBaselineValuesPolicyIds
	} else {

		r.ObjectIdPath = res.ObjectIdPath

	}
	return nil
}

// This object is used to assert a desired state where this CaPoolIssuancePolicyBaselineValuesPolicyIds is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyCaPoolIssuancePolicyBaselineValuesPolicyIds *CaPoolIssuancePolicyBaselineValuesPolicyIds = &CaPoolIssuancePolicyBaselineValuesPolicyIds{empty: true}

func (r *CaPoolIssuancePolicyBaselineValuesPolicyIds) Empty() bool {
	return r.empty
}

func (r *CaPoolIssuancePolicyBaselineValuesPolicyIds) String() string {
	return dcl.SprintResource(r)
}

func (r *CaPoolIssuancePolicyBaselineValuesPolicyIds) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type CaPoolIssuancePolicyBaselineValuesAdditionalExtensions struct {
	empty    bool                                                            `json:"-"`
	ObjectId *CaPoolIssuancePolicyBaselineValuesAdditionalExtensionsObjectId `json:"objectId"`
	Critical *bool                                                           `json:"critical"`
	Value    *string                                                         `json:"value"`
}

type jsonCaPoolIssuancePolicyBaselineValuesAdditionalExtensions CaPoolIssuancePolicyBaselineValuesAdditionalExtensions

func (r *CaPoolIssuancePolicyBaselineValuesAdditionalExtensions) UnmarshalJSON(data []byte) error {
	var res jsonCaPoolIssuancePolicyBaselineValuesAdditionalExtensions
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyCaPoolIssuancePolicyBaselineValuesAdditionalExtensions
	} else {

		r.ObjectId = res.ObjectId

		r.Critical = res.Critical

		r.Value = res.Value

	}
	return nil
}

// This object is used to assert a desired state where this CaPoolIssuancePolicyBaselineValuesAdditionalExtensions is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyCaPoolIssuancePolicyBaselineValuesAdditionalExtensions *CaPoolIssuancePolicyBaselineValuesAdditionalExtensions = &CaPoolIssuancePolicyBaselineValuesAdditionalExtensions{empty: true}

func (r *CaPoolIssuancePolicyBaselineValuesAdditionalExtensions) Empty() bool {
	return r.empty
}

func (r *CaPoolIssuancePolicyBaselineValuesAdditionalExtensions) String() string {
	return dcl.SprintResource(r)
}

func (r *CaPoolIssuancePolicyBaselineValuesAdditionalExtensions) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type CaPoolIssuancePolicyBaselineValuesAdditionalExtensionsObjectId struct {
	empty        bool    `json:"-"`
	ObjectIdPath []int64 `json:"objectIdPath"`
}

type jsonCaPoolIssuancePolicyBaselineValuesAdditionalExtensionsObjectId CaPoolIssuancePolicyBaselineValuesAdditionalExtensionsObjectId

func (r *CaPoolIssuancePolicyBaselineValuesAdditionalExtensionsObjectId) UnmarshalJSON(data []byte) error {
	var res jsonCaPoolIssuancePolicyBaselineValuesAdditionalExtensionsObjectId
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyCaPoolIssuancePolicyBaselineValuesAdditionalExtensionsObjectId
	} else {

		r.ObjectIdPath = res.ObjectIdPath

	}
	return nil
}

// This object is used to assert a desired state where this CaPoolIssuancePolicyBaselineValuesAdditionalExtensionsObjectId is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyCaPoolIssuancePolicyBaselineValuesAdditionalExtensionsObjectId *CaPoolIssuancePolicyBaselineValuesAdditionalExtensionsObjectId = &CaPoolIssuancePolicyBaselineValuesAdditionalExtensionsObjectId{empty: true}

func (r *CaPoolIssuancePolicyBaselineValuesAdditionalExtensionsObjectId) Empty() bool {
	return r.empty
}

func (r *CaPoolIssuancePolicyBaselineValuesAdditionalExtensionsObjectId) String() string {
	return dcl.SprintResource(r)
}

func (r *CaPoolIssuancePolicyBaselineValuesAdditionalExtensionsObjectId) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type CaPoolIssuancePolicyIdentityConstraints struct {
	empty                           bool                                                  `json:"-"`
	CelExpression                   *CaPoolIssuancePolicyIdentityConstraintsCelExpression `json:"celExpression"`
	AllowSubjectPassthrough         *bool                                                 `json:"allowSubjectPassthrough"`
	AllowSubjectAltNamesPassthrough *bool                                                 `json:"allowSubjectAltNamesPassthrough"`
}

type jsonCaPoolIssuancePolicyIdentityConstraints CaPoolIssuancePolicyIdentityConstraints

func (r *CaPoolIssuancePolicyIdentityConstraints) UnmarshalJSON(data []byte) error {
	var res jsonCaPoolIssuancePolicyIdentityConstraints
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyCaPoolIssuancePolicyIdentityConstraints
	} else {

		r.CelExpression = res.CelExpression

		r.AllowSubjectPassthrough = res.AllowSubjectPassthrough

		r.AllowSubjectAltNamesPassthrough = res.AllowSubjectAltNamesPassthrough

	}
	return nil
}

// This object is used to assert a desired state where this CaPoolIssuancePolicyIdentityConstraints is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyCaPoolIssuancePolicyIdentityConstraints *CaPoolIssuancePolicyIdentityConstraints = &CaPoolIssuancePolicyIdentityConstraints{empty: true}

func (r *CaPoolIssuancePolicyIdentityConstraints) Empty() bool {
	return r.empty
}

func (r *CaPoolIssuancePolicyIdentityConstraints) String() string {
	return dcl.SprintResource(r)
}

func (r *CaPoolIssuancePolicyIdentityConstraints) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type CaPoolIssuancePolicyIdentityConstraintsCelExpression struct {
	empty       bool    `json:"-"`
	Expression  *string `json:"expression"`
	Title       *string `json:"title"`
	Description *string `json:"description"`
	Location    *string `json:"location"`
}

type jsonCaPoolIssuancePolicyIdentityConstraintsCelExpression CaPoolIssuancePolicyIdentityConstraintsCelExpression

func (r *CaPoolIssuancePolicyIdentityConstraintsCelExpression) UnmarshalJSON(data []byte) error {
	var res jsonCaPoolIssuancePolicyIdentityConstraintsCelExpression
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyCaPoolIssuancePolicyIdentityConstraintsCelExpression
	} else {

		r.Expression = res.Expression

		r.Title = res.Title

		r.Description = res.Description

		r.Location = res.Location

	}
	return nil
}

// This object is used to assert a desired state where this CaPoolIssuancePolicyIdentityConstraintsCelExpression is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyCaPoolIssuancePolicyIdentityConstraintsCelExpression *CaPoolIssuancePolicyIdentityConstraintsCelExpression = &CaPoolIssuancePolicyIdentityConstraintsCelExpression{empty: true}

func (r *CaPoolIssuancePolicyIdentityConstraintsCelExpression) Empty() bool {
	return r.empty
}

func (r *CaPoolIssuancePolicyIdentityConstraintsCelExpression) String() string {
	return dcl.SprintResource(r)
}

func (r *CaPoolIssuancePolicyIdentityConstraintsCelExpression) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type CaPoolIssuancePolicyPassthroughExtensions struct {
	empty                bool                                                            `json:"-"`
	KnownExtensions      []CaPoolIssuancePolicyPassthroughExtensionsKnownExtensionsEnum  `json:"knownExtensions"`
	AdditionalExtensions []CaPoolIssuancePolicyPassthroughExtensionsAdditionalExtensions `json:"additionalExtensions"`
}

type jsonCaPoolIssuancePolicyPassthroughExtensions CaPoolIssuancePolicyPassthroughExtensions

func (r *CaPoolIssuancePolicyPassthroughExtensions) UnmarshalJSON(data []byte) error {
	var res jsonCaPoolIssuancePolicyPassthroughExtensions
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyCaPoolIssuancePolicyPassthroughExtensions
	} else {

		r.KnownExtensions = res.KnownExtensions

		r.AdditionalExtensions = res.AdditionalExtensions

	}
	return nil
}

// This object is used to assert a desired state where this CaPoolIssuancePolicyPassthroughExtensions is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyCaPoolIssuancePolicyPassthroughExtensions *CaPoolIssuancePolicyPassthroughExtensions = &CaPoolIssuancePolicyPassthroughExtensions{empty: true}

func (r *CaPoolIssuancePolicyPassthroughExtensions) Empty() bool {
	return r.empty
}

func (r *CaPoolIssuancePolicyPassthroughExtensions) String() string {
	return dcl.SprintResource(r)
}

func (r *CaPoolIssuancePolicyPassthroughExtensions) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type CaPoolIssuancePolicyPassthroughExtensionsAdditionalExtensions struct {
	empty        bool    `json:"-"`
	ObjectIdPath []int64 `json:"objectIdPath"`
}

type jsonCaPoolIssuancePolicyPassthroughExtensionsAdditionalExtensions CaPoolIssuancePolicyPassthroughExtensionsAdditionalExtensions

func (r *CaPoolIssuancePolicyPassthroughExtensionsAdditionalExtensions) UnmarshalJSON(data []byte) error {
	var res jsonCaPoolIssuancePolicyPassthroughExtensionsAdditionalExtensions
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyCaPoolIssuancePolicyPassthroughExtensionsAdditionalExtensions
	} else {

		r.ObjectIdPath = res.ObjectIdPath

	}
	return nil
}

// This object is used to assert a desired state where this CaPoolIssuancePolicyPassthroughExtensionsAdditionalExtensions is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyCaPoolIssuancePolicyPassthroughExtensionsAdditionalExtensions *CaPoolIssuancePolicyPassthroughExtensionsAdditionalExtensions = &CaPoolIssuancePolicyPassthroughExtensionsAdditionalExtensions{empty: true}

func (r *CaPoolIssuancePolicyPassthroughExtensionsAdditionalExtensions) Empty() bool {
	return r.empty
}

func (r *CaPoolIssuancePolicyPassthroughExtensionsAdditionalExtensions) String() string {
	return dcl.SprintResource(r)
}

func (r *CaPoolIssuancePolicyPassthroughExtensionsAdditionalExtensions) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type CaPoolPublishingOptions struct {
	empty         bool  `json:"-"`
	PublishCaCert *bool `json:"publishCaCert"`
	PublishCrl    *bool `json:"publishCrl"`
}

type jsonCaPoolPublishingOptions CaPoolPublishingOptions

func (r *CaPoolPublishingOptions) UnmarshalJSON(data []byte) error {
	var res jsonCaPoolPublishingOptions
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyCaPoolPublishingOptions
	} else {

		r.PublishCaCert = res.PublishCaCert

		r.PublishCrl = res.PublishCrl

	}
	return nil
}

// This object is used to assert a desired state where this CaPoolPublishingOptions is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyCaPoolPublishingOptions *CaPoolPublishingOptions = &CaPoolPublishingOptions{empty: true}

func (r *CaPoolPublishingOptions) Empty() bool {
	return r.empty
}

func (r *CaPoolPublishingOptions) String() string {
	return dcl.SprintResource(r)
}

func (r *CaPoolPublishingOptions) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

// Describe returns a simple description of this resource to ensure that automated tools
// can identify it.
func (r *CaPool) Describe() dcl.ServiceTypeVersion {
	return dcl.ServiceTypeVersion{
		Service: "privateca",
		Type:    "CaPool",
		Version: "privateca",
	}
}

func (r *CaPool) ID() (string, error) {
	if err := extractCaPoolFields(r); err != nil {
		return "", err
	}
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"name":               dcl.ValueOrEmptyString(nr.Name),
		"tier":               dcl.ValueOrEmptyString(nr.Tier),
		"issuance_policy":    dcl.ValueOrEmptyString(nr.IssuancePolicy),
		"publishing_options": dcl.ValueOrEmptyString(nr.PublishingOptions),
		"labels":             dcl.ValueOrEmptyString(nr.Labels),
		"project":            dcl.ValueOrEmptyString(nr.Project),
		"location":           dcl.ValueOrEmptyString(nr.Location),
	}
	return dcl.Nprintf("projects/{{project}}/locations/{{location}}/caPools/{{name}}", params), nil
}

const CaPoolMaxPage = -1

type CaPoolList struct {
	Items []*CaPool

	nextToken string

	pageSize int32

	resource *CaPool
}

func (l *CaPoolList) HasNext() bool {
	return l.nextToken != ""
}

func (l *CaPoolList) Next(ctx context.Context, c *Client) error {
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	if !l.HasNext() {
		return fmt.Errorf("no next page")
	}
	items, token, err := c.listCaPool(ctx, l.resource, l.nextToken, l.pageSize)
	if err != nil {
		return err
	}
	l.Items = items
	l.nextToken = token
	return err
}

func (c *Client) ListCaPool(ctx context.Context, project, location string) (*CaPoolList, error) {
	ctx = dcl.ContextWithRequestID(ctx)
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	return c.ListCaPoolWithMaxResults(ctx, project, location, CaPoolMaxPage)

}

func (c *Client) ListCaPoolWithMaxResults(ctx context.Context, project, location string, pageSize int32) (*CaPoolList, error) {
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	// Create a resource object so that we can use proper url normalization methods.
	r := &CaPool{
		Project:  &project,
		Location: &location,
	}
	items, token, err := c.listCaPool(ctx, r, "", pageSize)
	if err != nil {
		return nil, err
	}
	return &CaPoolList{
		Items:     items,
		nextToken: token,
		pageSize:  pageSize,
		resource:  r,
	}, nil
}

func (c *Client) GetCaPool(ctx context.Context, r *CaPool) (*CaPool, error) {
	ctx = dcl.ContextWithRequestID(ctx)
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	// This is *purposefully* supressing errors.
	// This function is used with url-normalized values + not URL normalized values.
	// URL Normalized values will throw unintentional errors, since those values are not of the proper parent form.
	extractCaPoolFields(r)

	b, err := c.getCaPoolRaw(ctx, r)
	if err != nil {
		if dcl.IsNotFound(err) {
			return nil, &googleapi.Error{
				Code:    404,
				Message: err.Error(),
			}
		}
		return nil, err
	}
	result, err := unmarshalCaPool(b, c, r)
	if err != nil {
		return nil, err
	}
	result.Project = r.Project
	result.Location = r.Location
	result.Name = r.Name

	c.Config.Logger.InfoWithContextf(ctx, "Retrieved raw result state: %v", result)
	c.Config.Logger.InfoWithContextf(ctx, "Canonicalizing with specified state: %v", r)
	result, err = canonicalizeCaPoolNewState(c, result, r)
	if err != nil {
		return nil, err
	}
	if err := postReadExtractCaPoolFields(result); err != nil {
		return result, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Created result state: %v", result)

	return result, nil
}

func (c *Client) DeleteCaPool(ctx context.Context, r *CaPool) error {
	ctx = dcl.ContextWithRequestID(ctx)
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	if r == nil {
		return fmt.Errorf("CaPool resource is nil")
	}
	c.Config.Logger.InfoWithContext(ctx, "Deleting CaPool...")
	deleteOp := deleteCaPoolOperation{}
	return deleteOp.do(ctx, r, c)
}

// DeleteAllCaPool deletes all resources that the filter functions returns true on.
func (c *Client) DeleteAllCaPool(ctx context.Context, project, location string, filter func(*CaPool) bool) error {
	listObj, err := c.ListCaPool(ctx, project, location)
	if err != nil {
		return err
	}

	err = c.deleteAllCaPool(ctx, filter, listObj.Items)
	if err != nil {
		return err
	}
	for listObj.HasNext() {
		err = listObj.Next(ctx, c)
		if err != nil {
			return nil
		}
		err = c.deleteAllCaPool(ctx, filter, listObj.Items)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Client) ApplyCaPool(ctx context.Context, rawDesired *CaPool, opts ...dcl.ApplyOption) (*CaPool, error) {
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	ctx = dcl.ContextWithRequestID(ctx)
	var resultNewState *CaPool
	err := dcl.Do(ctx, func(ctx context.Context) (*dcl.RetryDetails, error) {
		newState, err := applyCaPoolHelper(c, ctx, rawDesired, opts...)
		resultNewState = newState
		if err != nil {
			// If the error is 409, there is conflict in resource update.
			// Here we want to apply changes based on latest state.
			if dcl.IsConflictError(err) {
				return &dcl.RetryDetails{}, dcl.OperationNotDone{Err: err}
			}
			return nil, err
		}
		return nil, nil
	}, c.Config.RetryProvider)
	return resultNewState, err
}

func applyCaPoolHelper(c *Client, ctx context.Context, rawDesired *CaPool, opts ...dcl.ApplyOption) (*CaPool, error) {
	c.Config.Logger.InfoWithContext(ctx, "Beginning ApplyCaPool...")
	c.Config.Logger.InfoWithContextf(ctx, "User specified desired state: %v", rawDesired)

	// 1.1: Validation of user-specified fields in desired state.
	if err := rawDesired.validate(); err != nil {
		return nil, err
	}

	if err := extractCaPoolFields(rawDesired); err != nil {
		return nil, err
	}

	initial, desired, fieldDiffs, err := c.caPoolDiffsForRawDesired(ctx, rawDesired, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create a diff: %w", err)
	}

	diffs, err := convertFieldDiffsToCaPoolDiffs(c.Config, fieldDiffs, opts)
	if err != nil {
		return nil, err
	}

	// TODO(magic-modules-eng): 2.2 Feasibility check (all updates are feasible so far).

	// 2.3: Lifecycle Directive Check
	var create bool
	lp := dcl.FetchLifecycleParams(opts)
	if initial == nil {
		if dcl.HasLifecycleParam(lp, dcl.BlockCreation) {
			return nil, dcl.ApplyInfeasibleError{Message: fmt.Sprintf("Creation blocked by lifecycle params: %#v.", desired)}
		}
		create = true
	} else if dcl.HasLifecycleParam(lp, dcl.BlockAcquire) {
		return nil, dcl.ApplyInfeasibleError{
			Message: fmt.Sprintf("Resource already exists - apply blocked by lifecycle params: %#v.", initial),
		}
	} else {
		for _, d := range diffs {
			if d.RequiresRecreate {
				return nil, dcl.ApplyInfeasibleError{
					Message: fmt.Sprintf("infeasible update: (%v) would require recreation", d),
				}
			}
			if dcl.HasLifecycleParam(lp, dcl.BlockModification) {
				return nil, dcl.ApplyInfeasibleError{Message: fmt.Sprintf("Modification blocked, diff (%v) unresolvable.", d)}
			}
		}
	}

	// 2.4 Imperative Request Planning
	var ops []caPoolApiOperation
	if create {
		ops = append(ops, &createCaPoolOperation{})
	} else {
		for _, d := range diffs {
			ops = append(ops, d.UpdateOp)
		}
	}
	c.Config.Logger.InfoWithContextf(ctx, "Created plan: %#v", ops)

	// 2.5 Request Actuation
	for _, op := range ops {
		c.Config.Logger.InfoWithContextf(ctx, "Performing operation %T %+v", op, op)
		if err := op.do(ctx, desired, c); err != nil {
			c.Config.Logger.InfoWithContextf(ctx, "Failed operation %T %+v: %v", op, op, err)
			return nil, err
		}
		c.Config.Logger.InfoWithContextf(ctx, "Finished operation %T %+v", op, op)
	}
	return applyCaPoolDiff(c, ctx, desired, rawDesired, ops, opts...)
}

func applyCaPoolDiff(c *Client, ctx context.Context, desired *CaPool, rawDesired *CaPool, ops []caPoolApiOperation, opts ...dcl.ApplyOption) (*CaPool, error) {
	// 3.1, 3.2a Retrieval of raw new state & canonicalization with desired state
	c.Config.Logger.InfoWithContext(ctx, "Retrieving raw new state...")
	rawNew, err := c.GetCaPool(ctx, desired)
	if err != nil {
		return nil, err
	}
	// Get additional values from the first response.
	// These values should be merged into the newState above.
	if len(ops) > 0 {
		lastOp := ops[len(ops)-1]
		if o, ok := lastOp.(*createCaPoolOperation); ok {
			if r, hasR := o.FirstResponse(); hasR {

				c.Config.Logger.InfoWithContext(ctx, "Retrieving raw new state from operation...")

				fullResp, err := unmarshalMapCaPool(r, c, rawDesired)
				if err != nil {
					return nil, err
				}

				rawNew, err = canonicalizeCaPoolNewState(c, rawNew, fullResp)
				if err != nil {
					return nil, err
				}
			}
		}
	}

	c.Config.Logger.InfoWithContextf(ctx, "Canonicalizing with raw desired state: %v", rawDesired)
	// 3.2b Canonicalization of raw new state using raw desired state
	newState, err := canonicalizeCaPoolNewState(c, rawNew, rawDesired)
	if err != nil {
		return rawNew, err
	}

	c.Config.Logger.InfoWithContextf(ctx, "Created canonical new state: %v", newState)
	// 3.3 Comparison of the new state and raw desired state.
	// TODO(magic-modules-eng): EVENTUALLY_CONSISTENT_UPDATE
	newDesired, err := canonicalizeCaPoolDesiredState(rawDesired, newState)
	if err != nil {
		return newState, err
	}

	if err := postReadExtractCaPoolFields(newState); err != nil {
		return newState, err
	}

	// Need to ensure any transformations made here match acceptably in differ.
	if err := postReadExtractCaPoolFields(newDesired); err != nil {
		return newState, err
	}

	c.Config.Logger.InfoWithContextf(ctx, "Diffing using canonicalized desired state: %v", newDesired)
	newDiffs, err := diffCaPool(c, newDesired, newState)
	if err != nil {
		return newState, err
	}

	if len(newDiffs) == 0 {
		c.Config.Logger.InfoWithContext(ctx, "No diffs found. Apply was successful.")
	} else {
		c.Config.Logger.InfoWithContextf(ctx, "Found diffs: %v", newDiffs)
		diffMessages := make([]string, len(newDiffs))
		for i, d := range newDiffs {
			diffMessages[i] = fmt.Sprintf("%v", d)
		}
		return newState, dcl.DiffAfterApplyError{Diffs: diffMessages}
	}
	c.Config.Logger.InfoWithContext(ctx, "Done Apply.")
	return newState, nil
}
