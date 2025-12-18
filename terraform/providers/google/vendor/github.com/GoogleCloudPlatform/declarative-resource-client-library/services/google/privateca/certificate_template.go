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

type CertificateTemplate struct {
	Name                  *string                                   `json:"name"`
	PredefinedValues      *CertificateTemplatePredefinedValues      `json:"predefinedValues"`
	IdentityConstraints   *CertificateTemplateIdentityConstraints   `json:"identityConstraints"`
	PassthroughExtensions *CertificateTemplatePassthroughExtensions `json:"passthroughExtensions"`
	Description           *string                                   `json:"description"`
	CreateTime            *string                                   `json:"createTime"`
	UpdateTime            *string                                   `json:"updateTime"`
	Labels                map[string]string                         `json:"labels"`
	Project               *string                                   `json:"project"`
	Location              *string                                   `json:"location"`
}

func (r *CertificateTemplate) String() string {
	return dcl.SprintResource(r)
}

// The enum CertificateTemplatePassthroughExtensionsKnownExtensionsEnum.
type CertificateTemplatePassthroughExtensionsKnownExtensionsEnum string

// CertificateTemplatePassthroughExtensionsKnownExtensionsEnumRef returns a *CertificateTemplatePassthroughExtensionsKnownExtensionsEnum with the value of string s
// If the empty string is provided, nil is returned.
func CertificateTemplatePassthroughExtensionsKnownExtensionsEnumRef(s string) *CertificateTemplatePassthroughExtensionsKnownExtensionsEnum {
	v := CertificateTemplatePassthroughExtensionsKnownExtensionsEnum(s)
	return &v
}

func (v CertificateTemplatePassthroughExtensionsKnownExtensionsEnum) Validate() error {
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
		Enum:  "CertificateTemplatePassthroughExtensionsKnownExtensionsEnum",
		Value: string(v),
		Valid: []string{},
	}
}

type CertificateTemplatePredefinedValues struct {
	empty                bool                                                      `json:"-"`
	KeyUsage             *CertificateTemplatePredefinedValuesKeyUsage              `json:"keyUsage"`
	CaOptions            *CertificateTemplatePredefinedValuesCaOptions             `json:"caOptions"`
	PolicyIds            []CertificateTemplatePredefinedValuesPolicyIds            `json:"policyIds"`
	AiaOcspServers       []string                                                  `json:"aiaOcspServers"`
	AdditionalExtensions []CertificateTemplatePredefinedValuesAdditionalExtensions `json:"additionalExtensions"`
}

type jsonCertificateTemplatePredefinedValues CertificateTemplatePredefinedValues

func (r *CertificateTemplatePredefinedValues) UnmarshalJSON(data []byte) error {
	var res jsonCertificateTemplatePredefinedValues
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyCertificateTemplatePredefinedValues
	} else {

		r.KeyUsage = res.KeyUsage

		r.CaOptions = res.CaOptions

		r.PolicyIds = res.PolicyIds

		r.AiaOcspServers = res.AiaOcspServers

		r.AdditionalExtensions = res.AdditionalExtensions

	}
	return nil
}

// This object is used to assert a desired state where this CertificateTemplatePredefinedValues is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyCertificateTemplatePredefinedValues *CertificateTemplatePredefinedValues = &CertificateTemplatePredefinedValues{empty: true}

func (r *CertificateTemplatePredefinedValues) Empty() bool {
	return r.empty
}

func (r *CertificateTemplatePredefinedValues) String() string {
	return dcl.SprintResource(r)
}

func (r *CertificateTemplatePredefinedValues) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type CertificateTemplatePredefinedValuesKeyUsage struct {
	empty                    bool                                                                  `json:"-"`
	BaseKeyUsage             *CertificateTemplatePredefinedValuesKeyUsageBaseKeyUsage              `json:"baseKeyUsage"`
	ExtendedKeyUsage         *CertificateTemplatePredefinedValuesKeyUsageExtendedKeyUsage          `json:"extendedKeyUsage"`
	UnknownExtendedKeyUsages []CertificateTemplatePredefinedValuesKeyUsageUnknownExtendedKeyUsages `json:"unknownExtendedKeyUsages"`
}

type jsonCertificateTemplatePredefinedValuesKeyUsage CertificateTemplatePredefinedValuesKeyUsage

func (r *CertificateTemplatePredefinedValuesKeyUsage) UnmarshalJSON(data []byte) error {
	var res jsonCertificateTemplatePredefinedValuesKeyUsage
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyCertificateTemplatePredefinedValuesKeyUsage
	} else {

		r.BaseKeyUsage = res.BaseKeyUsage

		r.ExtendedKeyUsage = res.ExtendedKeyUsage

		r.UnknownExtendedKeyUsages = res.UnknownExtendedKeyUsages

	}
	return nil
}

// This object is used to assert a desired state where this CertificateTemplatePredefinedValuesKeyUsage is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyCertificateTemplatePredefinedValuesKeyUsage *CertificateTemplatePredefinedValuesKeyUsage = &CertificateTemplatePredefinedValuesKeyUsage{empty: true}

func (r *CertificateTemplatePredefinedValuesKeyUsage) Empty() bool {
	return r.empty
}

func (r *CertificateTemplatePredefinedValuesKeyUsage) String() string {
	return dcl.SprintResource(r)
}

func (r *CertificateTemplatePredefinedValuesKeyUsage) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type CertificateTemplatePredefinedValuesKeyUsageBaseKeyUsage struct {
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

type jsonCertificateTemplatePredefinedValuesKeyUsageBaseKeyUsage CertificateTemplatePredefinedValuesKeyUsageBaseKeyUsage

func (r *CertificateTemplatePredefinedValuesKeyUsageBaseKeyUsage) UnmarshalJSON(data []byte) error {
	var res jsonCertificateTemplatePredefinedValuesKeyUsageBaseKeyUsage
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyCertificateTemplatePredefinedValuesKeyUsageBaseKeyUsage
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

// This object is used to assert a desired state where this CertificateTemplatePredefinedValuesKeyUsageBaseKeyUsage is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyCertificateTemplatePredefinedValuesKeyUsageBaseKeyUsage *CertificateTemplatePredefinedValuesKeyUsageBaseKeyUsage = &CertificateTemplatePredefinedValuesKeyUsageBaseKeyUsage{empty: true}

func (r *CertificateTemplatePredefinedValuesKeyUsageBaseKeyUsage) Empty() bool {
	return r.empty
}

func (r *CertificateTemplatePredefinedValuesKeyUsageBaseKeyUsage) String() string {
	return dcl.SprintResource(r)
}

func (r *CertificateTemplatePredefinedValuesKeyUsageBaseKeyUsage) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type CertificateTemplatePredefinedValuesKeyUsageExtendedKeyUsage struct {
	empty           bool  `json:"-"`
	ServerAuth      *bool `json:"serverAuth"`
	ClientAuth      *bool `json:"clientAuth"`
	CodeSigning     *bool `json:"codeSigning"`
	EmailProtection *bool `json:"emailProtection"`
	TimeStamping    *bool `json:"timeStamping"`
	OcspSigning     *bool `json:"ocspSigning"`
}

type jsonCertificateTemplatePredefinedValuesKeyUsageExtendedKeyUsage CertificateTemplatePredefinedValuesKeyUsageExtendedKeyUsage

func (r *CertificateTemplatePredefinedValuesKeyUsageExtendedKeyUsage) UnmarshalJSON(data []byte) error {
	var res jsonCertificateTemplatePredefinedValuesKeyUsageExtendedKeyUsage
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyCertificateTemplatePredefinedValuesKeyUsageExtendedKeyUsage
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

// This object is used to assert a desired state where this CertificateTemplatePredefinedValuesKeyUsageExtendedKeyUsage is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyCertificateTemplatePredefinedValuesKeyUsageExtendedKeyUsage *CertificateTemplatePredefinedValuesKeyUsageExtendedKeyUsage = &CertificateTemplatePredefinedValuesKeyUsageExtendedKeyUsage{empty: true}

func (r *CertificateTemplatePredefinedValuesKeyUsageExtendedKeyUsage) Empty() bool {
	return r.empty
}

func (r *CertificateTemplatePredefinedValuesKeyUsageExtendedKeyUsage) String() string {
	return dcl.SprintResource(r)
}

func (r *CertificateTemplatePredefinedValuesKeyUsageExtendedKeyUsage) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type CertificateTemplatePredefinedValuesKeyUsageUnknownExtendedKeyUsages struct {
	empty        bool    `json:"-"`
	ObjectIdPath []int64 `json:"objectIdPath"`
}

type jsonCertificateTemplatePredefinedValuesKeyUsageUnknownExtendedKeyUsages CertificateTemplatePredefinedValuesKeyUsageUnknownExtendedKeyUsages

func (r *CertificateTemplatePredefinedValuesKeyUsageUnknownExtendedKeyUsages) UnmarshalJSON(data []byte) error {
	var res jsonCertificateTemplatePredefinedValuesKeyUsageUnknownExtendedKeyUsages
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyCertificateTemplatePredefinedValuesKeyUsageUnknownExtendedKeyUsages
	} else {

		r.ObjectIdPath = res.ObjectIdPath

	}
	return nil
}

// This object is used to assert a desired state where this CertificateTemplatePredefinedValuesKeyUsageUnknownExtendedKeyUsages is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyCertificateTemplatePredefinedValuesKeyUsageUnknownExtendedKeyUsages *CertificateTemplatePredefinedValuesKeyUsageUnknownExtendedKeyUsages = &CertificateTemplatePredefinedValuesKeyUsageUnknownExtendedKeyUsages{empty: true}

func (r *CertificateTemplatePredefinedValuesKeyUsageUnknownExtendedKeyUsages) Empty() bool {
	return r.empty
}

func (r *CertificateTemplatePredefinedValuesKeyUsageUnknownExtendedKeyUsages) String() string {
	return dcl.SprintResource(r)
}

func (r *CertificateTemplatePredefinedValuesKeyUsageUnknownExtendedKeyUsages) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type CertificateTemplatePredefinedValuesCaOptions struct {
	empty               bool   `json:"-"`
	IsCa                *bool  `json:"isCa"`
	MaxIssuerPathLength *int64 `json:"maxIssuerPathLength"`
}

type jsonCertificateTemplatePredefinedValuesCaOptions CertificateTemplatePredefinedValuesCaOptions

func (r *CertificateTemplatePredefinedValuesCaOptions) UnmarshalJSON(data []byte) error {
	var res jsonCertificateTemplatePredefinedValuesCaOptions
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyCertificateTemplatePredefinedValuesCaOptions
	} else {

		r.IsCa = res.IsCa

		r.MaxIssuerPathLength = res.MaxIssuerPathLength

	}
	return nil
}

// This object is used to assert a desired state where this CertificateTemplatePredefinedValuesCaOptions is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyCertificateTemplatePredefinedValuesCaOptions *CertificateTemplatePredefinedValuesCaOptions = &CertificateTemplatePredefinedValuesCaOptions{empty: true}

func (r *CertificateTemplatePredefinedValuesCaOptions) Empty() bool {
	return r.empty
}

func (r *CertificateTemplatePredefinedValuesCaOptions) String() string {
	return dcl.SprintResource(r)
}

func (r *CertificateTemplatePredefinedValuesCaOptions) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type CertificateTemplatePredefinedValuesPolicyIds struct {
	empty        bool    `json:"-"`
	ObjectIdPath []int64 `json:"objectIdPath"`
}

type jsonCertificateTemplatePredefinedValuesPolicyIds CertificateTemplatePredefinedValuesPolicyIds

func (r *CertificateTemplatePredefinedValuesPolicyIds) UnmarshalJSON(data []byte) error {
	var res jsonCertificateTemplatePredefinedValuesPolicyIds
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyCertificateTemplatePredefinedValuesPolicyIds
	} else {

		r.ObjectIdPath = res.ObjectIdPath

	}
	return nil
}

// This object is used to assert a desired state where this CertificateTemplatePredefinedValuesPolicyIds is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyCertificateTemplatePredefinedValuesPolicyIds *CertificateTemplatePredefinedValuesPolicyIds = &CertificateTemplatePredefinedValuesPolicyIds{empty: true}

func (r *CertificateTemplatePredefinedValuesPolicyIds) Empty() bool {
	return r.empty
}

func (r *CertificateTemplatePredefinedValuesPolicyIds) String() string {
	return dcl.SprintResource(r)
}

func (r *CertificateTemplatePredefinedValuesPolicyIds) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type CertificateTemplatePredefinedValuesAdditionalExtensions struct {
	empty    bool                                                             `json:"-"`
	ObjectId *CertificateTemplatePredefinedValuesAdditionalExtensionsObjectId `json:"objectId"`
	Critical *bool                                                            `json:"critical"`
	Value    *string                                                          `json:"value"`
}

type jsonCertificateTemplatePredefinedValuesAdditionalExtensions CertificateTemplatePredefinedValuesAdditionalExtensions

func (r *CertificateTemplatePredefinedValuesAdditionalExtensions) UnmarshalJSON(data []byte) error {
	var res jsonCertificateTemplatePredefinedValuesAdditionalExtensions
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyCertificateTemplatePredefinedValuesAdditionalExtensions
	} else {

		r.ObjectId = res.ObjectId

		r.Critical = res.Critical

		r.Value = res.Value

	}
	return nil
}

// This object is used to assert a desired state where this CertificateTemplatePredefinedValuesAdditionalExtensions is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyCertificateTemplatePredefinedValuesAdditionalExtensions *CertificateTemplatePredefinedValuesAdditionalExtensions = &CertificateTemplatePredefinedValuesAdditionalExtensions{empty: true}

func (r *CertificateTemplatePredefinedValuesAdditionalExtensions) Empty() bool {
	return r.empty
}

func (r *CertificateTemplatePredefinedValuesAdditionalExtensions) String() string {
	return dcl.SprintResource(r)
}

func (r *CertificateTemplatePredefinedValuesAdditionalExtensions) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type CertificateTemplatePredefinedValuesAdditionalExtensionsObjectId struct {
	empty        bool    `json:"-"`
	ObjectIdPath []int64 `json:"objectIdPath"`
}

type jsonCertificateTemplatePredefinedValuesAdditionalExtensionsObjectId CertificateTemplatePredefinedValuesAdditionalExtensionsObjectId

func (r *CertificateTemplatePredefinedValuesAdditionalExtensionsObjectId) UnmarshalJSON(data []byte) error {
	var res jsonCertificateTemplatePredefinedValuesAdditionalExtensionsObjectId
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyCertificateTemplatePredefinedValuesAdditionalExtensionsObjectId
	} else {

		r.ObjectIdPath = res.ObjectIdPath

	}
	return nil
}

// This object is used to assert a desired state where this CertificateTemplatePredefinedValuesAdditionalExtensionsObjectId is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyCertificateTemplatePredefinedValuesAdditionalExtensionsObjectId *CertificateTemplatePredefinedValuesAdditionalExtensionsObjectId = &CertificateTemplatePredefinedValuesAdditionalExtensionsObjectId{empty: true}

func (r *CertificateTemplatePredefinedValuesAdditionalExtensionsObjectId) Empty() bool {
	return r.empty
}

func (r *CertificateTemplatePredefinedValuesAdditionalExtensionsObjectId) String() string {
	return dcl.SprintResource(r)
}

func (r *CertificateTemplatePredefinedValuesAdditionalExtensionsObjectId) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type CertificateTemplateIdentityConstraints struct {
	empty                           bool                                                 `json:"-"`
	CelExpression                   *CertificateTemplateIdentityConstraintsCelExpression `json:"celExpression"`
	AllowSubjectPassthrough         *bool                                                `json:"allowSubjectPassthrough"`
	AllowSubjectAltNamesPassthrough *bool                                                `json:"allowSubjectAltNamesPassthrough"`
}

type jsonCertificateTemplateIdentityConstraints CertificateTemplateIdentityConstraints

func (r *CertificateTemplateIdentityConstraints) UnmarshalJSON(data []byte) error {
	var res jsonCertificateTemplateIdentityConstraints
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyCertificateTemplateIdentityConstraints
	} else {

		r.CelExpression = res.CelExpression

		r.AllowSubjectPassthrough = res.AllowSubjectPassthrough

		r.AllowSubjectAltNamesPassthrough = res.AllowSubjectAltNamesPassthrough

	}
	return nil
}

// This object is used to assert a desired state where this CertificateTemplateIdentityConstraints is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyCertificateTemplateIdentityConstraints *CertificateTemplateIdentityConstraints = &CertificateTemplateIdentityConstraints{empty: true}

func (r *CertificateTemplateIdentityConstraints) Empty() bool {
	return r.empty
}

func (r *CertificateTemplateIdentityConstraints) String() string {
	return dcl.SprintResource(r)
}

func (r *CertificateTemplateIdentityConstraints) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type CertificateTemplateIdentityConstraintsCelExpression struct {
	empty       bool    `json:"-"`
	Expression  *string `json:"expression"`
	Title       *string `json:"title"`
	Description *string `json:"description"`
	Location    *string `json:"location"`
}

type jsonCertificateTemplateIdentityConstraintsCelExpression CertificateTemplateIdentityConstraintsCelExpression

func (r *CertificateTemplateIdentityConstraintsCelExpression) UnmarshalJSON(data []byte) error {
	var res jsonCertificateTemplateIdentityConstraintsCelExpression
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyCertificateTemplateIdentityConstraintsCelExpression
	} else {

		r.Expression = res.Expression

		r.Title = res.Title

		r.Description = res.Description

		r.Location = res.Location

	}
	return nil
}

// This object is used to assert a desired state where this CertificateTemplateIdentityConstraintsCelExpression is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyCertificateTemplateIdentityConstraintsCelExpression *CertificateTemplateIdentityConstraintsCelExpression = &CertificateTemplateIdentityConstraintsCelExpression{empty: true}

func (r *CertificateTemplateIdentityConstraintsCelExpression) Empty() bool {
	return r.empty
}

func (r *CertificateTemplateIdentityConstraintsCelExpression) String() string {
	return dcl.SprintResource(r)
}

func (r *CertificateTemplateIdentityConstraintsCelExpression) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type CertificateTemplatePassthroughExtensions struct {
	empty                bool                                                           `json:"-"`
	KnownExtensions      []CertificateTemplatePassthroughExtensionsKnownExtensionsEnum  `json:"knownExtensions"`
	AdditionalExtensions []CertificateTemplatePassthroughExtensionsAdditionalExtensions `json:"additionalExtensions"`
}

type jsonCertificateTemplatePassthroughExtensions CertificateTemplatePassthroughExtensions

func (r *CertificateTemplatePassthroughExtensions) UnmarshalJSON(data []byte) error {
	var res jsonCertificateTemplatePassthroughExtensions
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyCertificateTemplatePassthroughExtensions
	} else {

		r.KnownExtensions = res.KnownExtensions

		r.AdditionalExtensions = res.AdditionalExtensions

	}
	return nil
}

// This object is used to assert a desired state where this CertificateTemplatePassthroughExtensions is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyCertificateTemplatePassthroughExtensions *CertificateTemplatePassthroughExtensions = &CertificateTemplatePassthroughExtensions{empty: true}

func (r *CertificateTemplatePassthroughExtensions) Empty() bool {
	return r.empty
}

func (r *CertificateTemplatePassthroughExtensions) String() string {
	return dcl.SprintResource(r)
}

func (r *CertificateTemplatePassthroughExtensions) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type CertificateTemplatePassthroughExtensionsAdditionalExtensions struct {
	empty        bool    `json:"-"`
	ObjectIdPath []int64 `json:"objectIdPath"`
}

type jsonCertificateTemplatePassthroughExtensionsAdditionalExtensions CertificateTemplatePassthroughExtensionsAdditionalExtensions

func (r *CertificateTemplatePassthroughExtensionsAdditionalExtensions) UnmarshalJSON(data []byte) error {
	var res jsonCertificateTemplatePassthroughExtensionsAdditionalExtensions
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyCertificateTemplatePassthroughExtensionsAdditionalExtensions
	} else {

		r.ObjectIdPath = res.ObjectIdPath

	}
	return nil
}

// This object is used to assert a desired state where this CertificateTemplatePassthroughExtensionsAdditionalExtensions is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyCertificateTemplatePassthroughExtensionsAdditionalExtensions *CertificateTemplatePassthroughExtensionsAdditionalExtensions = &CertificateTemplatePassthroughExtensionsAdditionalExtensions{empty: true}

func (r *CertificateTemplatePassthroughExtensionsAdditionalExtensions) Empty() bool {
	return r.empty
}

func (r *CertificateTemplatePassthroughExtensionsAdditionalExtensions) String() string {
	return dcl.SprintResource(r)
}

func (r *CertificateTemplatePassthroughExtensionsAdditionalExtensions) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

// Describe returns a simple description of this resource to ensure that automated tools
// can identify it.
func (r *CertificateTemplate) Describe() dcl.ServiceTypeVersion {
	return dcl.ServiceTypeVersion{
		Service: "privateca",
		Type:    "CertificateTemplate",
		Version: "privateca",
	}
}

func (r *CertificateTemplate) ID() (string, error) {
	if err := extractCertificateTemplateFields(r); err != nil {
		return "", err
	}
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"name":                   dcl.ValueOrEmptyString(nr.Name),
		"predefined_values":      dcl.ValueOrEmptyString(nr.PredefinedValues),
		"identity_constraints":   dcl.ValueOrEmptyString(nr.IdentityConstraints),
		"passthrough_extensions": dcl.ValueOrEmptyString(nr.PassthroughExtensions),
		"description":            dcl.ValueOrEmptyString(nr.Description),
		"create_time":            dcl.ValueOrEmptyString(nr.CreateTime),
		"update_time":            dcl.ValueOrEmptyString(nr.UpdateTime),
		"labels":                 dcl.ValueOrEmptyString(nr.Labels),
		"project":                dcl.ValueOrEmptyString(nr.Project),
		"location":               dcl.ValueOrEmptyString(nr.Location),
	}
	return dcl.Nprintf("projects/{{project}}/locations/{{location}}/certificateTemplates/{{name}}", params), nil
}

const CertificateTemplateMaxPage = -1

type CertificateTemplateList struct {
	Items []*CertificateTemplate

	nextToken string

	pageSize int32

	resource *CertificateTemplate
}

func (l *CertificateTemplateList) HasNext() bool {
	return l.nextToken != ""
}

func (l *CertificateTemplateList) Next(ctx context.Context, c *Client) error {
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	if !l.HasNext() {
		return fmt.Errorf("no next page")
	}
	items, token, err := c.listCertificateTemplate(ctx, l.resource, l.nextToken, l.pageSize)
	if err != nil {
		return err
	}
	l.Items = items
	l.nextToken = token
	return err
}

func (c *Client) ListCertificateTemplate(ctx context.Context, project, location string) (*CertificateTemplateList, error) {
	ctx = dcl.ContextWithRequestID(ctx)
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	return c.ListCertificateTemplateWithMaxResults(ctx, project, location, CertificateTemplateMaxPage)

}

func (c *Client) ListCertificateTemplateWithMaxResults(ctx context.Context, project, location string, pageSize int32) (*CertificateTemplateList, error) {
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	// Create a resource object so that we can use proper url normalization methods.
	r := &CertificateTemplate{
		Project:  &project,
		Location: &location,
	}
	items, token, err := c.listCertificateTemplate(ctx, r, "", pageSize)
	if err != nil {
		return nil, err
	}
	return &CertificateTemplateList{
		Items:     items,
		nextToken: token,
		pageSize:  pageSize,
		resource:  r,
	}, nil
}

func (c *Client) GetCertificateTemplate(ctx context.Context, r *CertificateTemplate) (*CertificateTemplate, error) {
	ctx = dcl.ContextWithRequestID(ctx)
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	// This is *purposefully* supressing errors.
	// This function is used with url-normalized values + not URL normalized values.
	// URL Normalized values will throw unintentional errors, since those values are not of the proper parent form.
	extractCertificateTemplateFields(r)

	b, err := c.getCertificateTemplateRaw(ctx, r)
	if err != nil {
		if dcl.IsNotFound(err) {
			return nil, &googleapi.Error{
				Code:    404,
				Message: err.Error(),
			}
		}
		return nil, err
	}
	result, err := unmarshalCertificateTemplate(b, c, r)
	if err != nil {
		return nil, err
	}
	result.Project = r.Project
	result.Location = r.Location
	result.Name = r.Name

	c.Config.Logger.InfoWithContextf(ctx, "Retrieved raw result state: %v", result)
	c.Config.Logger.InfoWithContextf(ctx, "Canonicalizing with specified state: %v", r)
	result, err = canonicalizeCertificateTemplateNewState(c, result, r)
	if err != nil {
		return nil, err
	}
	if err := postReadExtractCertificateTemplateFields(result); err != nil {
		return result, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Created result state: %v", result)

	return result, nil
}

func (c *Client) DeleteCertificateTemplate(ctx context.Context, r *CertificateTemplate) error {
	ctx = dcl.ContextWithRequestID(ctx)
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	if r == nil {
		return fmt.Errorf("CertificateTemplate resource is nil")
	}
	c.Config.Logger.InfoWithContext(ctx, "Deleting CertificateTemplate...")
	deleteOp := deleteCertificateTemplateOperation{}
	return deleteOp.do(ctx, r, c)
}

// DeleteAllCertificateTemplate deletes all resources that the filter functions returns true on.
func (c *Client) DeleteAllCertificateTemplate(ctx context.Context, project, location string, filter func(*CertificateTemplate) bool) error {
	listObj, err := c.ListCertificateTemplate(ctx, project, location)
	if err != nil {
		return err
	}

	err = c.deleteAllCertificateTemplate(ctx, filter, listObj.Items)
	if err != nil {
		return err
	}
	for listObj.HasNext() {
		err = listObj.Next(ctx, c)
		if err != nil {
			return nil
		}
		err = c.deleteAllCertificateTemplate(ctx, filter, listObj.Items)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Client) ApplyCertificateTemplate(ctx context.Context, rawDesired *CertificateTemplate, opts ...dcl.ApplyOption) (*CertificateTemplate, error) {
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	ctx = dcl.ContextWithRequestID(ctx)
	var resultNewState *CertificateTemplate
	err := dcl.Do(ctx, func(ctx context.Context) (*dcl.RetryDetails, error) {
		newState, err := applyCertificateTemplateHelper(c, ctx, rawDesired, opts...)
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

func applyCertificateTemplateHelper(c *Client, ctx context.Context, rawDesired *CertificateTemplate, opts ...dcl.ApplyOption) (*CertificateTemplate, error) {
	c.Config.Logger.InfoWithContext(ctx, "Beginning ApplyCertificateTemplate...")
	c.Config.Logger.InfoWithContextf(ctx, "User specified desired state: %v", rawDesired)

	// 1.1: Validation of user-specified fields in desired state.
	if err := rawDesired.validate(); err != nil {
		return nil, err
	}

	if err := extractCertificateTemplateFields(rawDesired); err != nil {
		return nil, err
	}

	initial, desired, fieldDiffs, err := c.certificateTemplateDiffsForRawDesired(ctx, rawDesired, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create a diff: %w", err)
	}

	diffs, err := convertFieldDiffsToCertificateTemplateDiffs(c.Config, fieldDiffs, opts)
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
	var ops []certificateTemplateApiOperation
	if create {
		ops = append(ops, &createCertificateTemplateOperation{})
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
	return applyCertificateTemplateDiff(c, ctx, desired, rawDesired, ops, opts...)
}

func applyCertificateTemplateDiff(c *Client, ctx context.Context, desired *CertificateTemplate, rawDesired *CertificateTemplate, ops []certificateTemplateApiOperation, opts ...dcl.ApplyOption) (*CertificateTemplate, error) {
	// 3.1, 3.2a Retrieval of raw new state & canonicalization with desired state
	c.Config.Logger.InfoWithContext(ctx, "Retrieving raw new state...")
	rawNew, err := c.GetCertificateTemplate(ctx, desired)
	if err != nil {
		return nil, err
	}
	// Get additional values from the first response.
	// These values should be merged into the newState above.
	if len(ops) > 0 {
		lastOp := ops[len(ops)-1]
		if o, ok := lastOp.(*createCertificateTemplateOperation); ok {
			if r, hasR := o.FirstResponse(); hasR {

				c.Config.Logger.InfoWithContext(ctx, "Retrieving raw new state from operation...")

				fullResp, err := unmarshalMapCertificateTemplate(r, c, rawDesired)
				if err != nil {
					return nil, err
				}

				rawNew, err = canonicalizeCertificateTemplateNewState(c, rawNew, fullResp)
				if err != nil {
					return nil, err
				}
			}
		}
	}

	c.Config.Logger.InfoWithContextf(ctx, "Canonicalizing with raw desired state: %v", rawDesired)
	// 3.2b Canonicalization of raw new state using raw desired state
	newState, err := canonicalizeCertificateTemplateNewState(c, rawNew, rawDesired)
	if err != nil {
		return rawNew, err
	}

	c.Config.Logger.InfoWithContextf(ctx, "Created canonical new state: %v", newState)
	// 3.3 Comparison of the new state and raw desired state.
	// TODO(magic-modules-eng): EVENTUALLY_CONSISTENT_UPDATE
	newDesired, err := canonicalizeCertificateTemplateDesiredState(rawDesired, newState)
	if err != nil {
		return newState, err
	}

	if err := postReadExtractCertificateTemplateFields(newState); err != nil {
		return newState, err
	}

	// Need to ensure any transformations made here match acceptably in differ.
	if err := postReadExtractCertificateTemplateFields(newDesired); err != nil {
		return newState, err
	}

	c.Config.Logger.InfoWithContextf(ctx, "Diffing using canonicalized desired state: %v", newDesired)
	newDiffs, err := diffCertificateTemplate(c, newDesired, newState)
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
