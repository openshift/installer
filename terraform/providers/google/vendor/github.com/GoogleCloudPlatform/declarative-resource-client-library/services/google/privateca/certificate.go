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

type Certificate struct {
	Name                       *string                            `json:"name"`
	PemCsr                     *string                            `json:"pemCsr"`
	Config                     *CertificateConfig                 `json:"config"`
	IssuerCertificateAuthority *string                            `json:"issuerCertificateAuthority"`
	Lifetime                   *string                            `json:"lifetime"`
	CertificateTemplate        *string                            `json:"certificateTemplate"`
	SubjectMode                *CertificateSubjectModeEnum        `json:"subjectMode"`
	RevocationDetails          *CertificateRevocationDetails      `json:"revocationDetails"`
	PemCertificate             *string                            `json:"pemCertificate"`
	CertificateDescription     *CertificateCertificateDescription `json:"certificateDescription"`
	PemCertificateChain        []string                           `json:"pemCertificateChain"`
	CreateTime                 *string                            `json:"createTime"`
	UpdateTime                 *string                            `json:"updateTime"`
	Labels                     map[string]string                  `json:"labels"`
	Project                    *string                            `json:"project"`
	Location                   *string                            `json:"location"`
	CaPool                     *string                            `json:"caPool"`
	CertificateAuthority       *string                            `json:"certificateAuthority"`
}

func (r *Certificate) String() string {
	return dcl.SprintResource(r)
}

// The enum CertificateConfigPublicKeyFormatEnum.
type CertificateConfigPublicKeyFormatEnum string

// CertificateConfigPublicKeyFormatEnumRef returns a *CertificateConfigPublicKeyFormatEnum with the value of string s
// If the empty string is provided, nil is returned.
func CertificateConfigPublicKeyFormatEnumRef(s string) *CertificateConfigPublicKeyFormatEnum {
	v := CertificateConfigPublicKeyFormatEnum(s)
	return &v
}

func (v CertificateConfigPublicKeyFormatEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"KEY_FORMAT_UNSPECIFIED", "PEM"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "CertificateConfigPublicKeyFormatEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum CertificateSubjectModeEnum.
type CertificateSubjectModeEnum string

// CertificateSubjectModeEnumRef returns a *CertificateSubjectModeEnum with the value of string s
// If the empty string is provided, nil is returned.
func CertificateSubjectModeEnumRef(s string) *CertificateSubjectModeEnum {
	v := CertificateSubjectModeEnum(s)
	return &v
}

func (v CertificateSubjectModeEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"SUBJECT_REQUEST_MODE_UNSPECIFIED", "DEFAULT", "REFLECTED_SPIFFE"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "CertificateSubjectModeEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum CertificateRevocationDetailsRevocationStateEnum.
type CertificateRevocationDetailsRevocationStateEnum string

// CertificateRevocationDetailsRevocationStateEnumRef returns a *CertificateRevocationDetailsRevocationStateEnum with the value of string s
// If the empty string is provided, nil is returned.
func CertificateRevocationDetailsRevocationStateEnumRef(s string) *CertificateRevocationDetailsRevocationStateEnum {
	v := CertificateRevocationDetailsRevocationStateEnum(s)
	return &v
}

func (v CertificateRevocationDetailsRevocationStateEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"REVOCATION_REASON_UNSPECIFIED", "KEY_COMPROMISE", "CERTIFICATE_AUTHORITY_COMPROMISE", "AFFILIATION_CHANGED", "SUPERSEDED", "CESSATION_OF_OPERATION", "CERTIFICATE_HOLD", "PRIVILEGE_WITHDRAWN", "ATTRIBUTE_AUTHORITY_COMPROMISE"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "CertificateRevocationDetailsRevocationStateEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum CertificateCertificateDescriptionPublicKeyFormatEnum.
type CertificateCertificateDescriptionPublicKeyFormatEnum string

// CertificateCertificateDescriptionPublicKeyFormatEnumRef returns a *CertificateCertificateDescriptionPublicKeyFormatEnum with the value of string s
// If the empty string is provided, nil is returned.
func CertificateCertificateDescriptionPublicKeyFormatEnumRef(s string) *CertificateCertificateDescriptionPublicKeyFormatEnum {
	v := CertificateCertificateDescriptionPublicKeyFormatEnum(s)
	return &v
}

func (v CertificateCertificateDescriptionPublicKeyFormatEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"KEY_FORMAT_UNSPECIFIED", "PEM"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "CertificateCertificateDescriptionPublicKeyFormatEnum",
		Value: string(v),
		Valid: []string{},
	}
}

type CertificateConfig struct {
	empty         bool                            `json:"-"`
	SubjectConfig *CertificateConfigSubjectConfig `json:"subjectConfig"`
	X509Config    *CertificateConfigX509Config    `json:"x509Config"`
	PublicKey     *CertificateConfigPublicKey     `json:"publicKey"`
}

type jsonCertificateConfig CertificateConfig

func (r *CertificateConfig) UnmarshalJSON(data []byte) error {
	var res jsonCertificateConfig
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyCertificateConfig
	} else {

		r.SubjectConfig = res.SubjectConfig

		r.X509Config = res.X509Config

		r.PublicKey = res.PublicKey

	}
	return nil
}

// This object is used to assert a desired state where this CertificateConfig is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyCertificateConfig *CertificateConfig = &CertificateConfig{empty: true}

func (r *CertificateConfig) Empty() bool {
	return r.empty
}

func (r *CertificateConfig) String() string {
	return dcl.SprintResource(r)
}

func (r *CertificateConfig) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type CertificateConfigSubjectConfig struct {
	empty          bool                                          `json:"-"`
	Subject        *CertificateConfigSubjectConfigSubject        `json:"subject"`
	SubjectAltName *CertificateConfigSubjectConfigSubjectAltName `json:"subjectAltName"`
}

type jsonCertificateConfigSubjectConfig CertificateConfigSubjectConfig

func (r *CertificateConfigSubjectConfig) UnmarshalJSON(data []byte) error {
	var res jsonCertificateConfigSubjectConfig
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyCertificateConfigSubjectConfig
	} else {

		r.Subject = res.Subject

		r.SubjectAltName = res.SubjectAltName

	}
	return nil
}

// This object is used to assert a desired state where this CertificateConfigSubjectConfig is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyCertificateConfigSubjectConfig *CertificateConfigSubjectConfig = &CertificateConfigSubjectConfig{empty: true}

func (r *CertificateConfigSubjectConfig) Empty() bool {
	return r.empty
}

func (r *CertificateConfigSubjectConfig) String() string {
	return dcl.SprintResource(r)
}

func (r *CertificateConfigSubjectConfig) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type CertificateConfigSubjectConfigSubject struct {
	empty              bool    `json:"-"`
	CommonName         *string `json:"commonName"`
	CountryCode        *string `json:"countryCode"`
	Organization       *string `json:"organization"`
	OrganizationalUnit *string `json:"organizationalUnit"`
	Locality           *string `json:"locality"`
	Province           *string `json:"province"`
	StreetAddress      *string `json:"streetAddress"`
	PostalCode         *string `json:"postalCode"`
}

type jsonCertificateConfigSubjectConfigSubject CertificateConfigSubjectConfigSubject

func (r *CertificateConfigSubjectConfigSubject) UnmarshalJSON(data []byte) error {
	var res jsonCertificateConfigSubjectConfigSubject
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyCertificateConfigSubjectConfigSubject
	} else {

		r.CommonName = res.CommonName

		r.CountryCode = res.CountryCode

		r.Organization = res.Organization

		r.OrganizationalUnit = res.OrganizationalUnit

		r.Locality = res.Locality

		r.Province = res.Province

		r.StreetAddress = res.StreetAddress

		r.PostalCode = res.PostalCode

	}
	return nil
}

// This object is used to assert a desired state where this CertificateConfigSubjectConfigSubject is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyCertificateConfigSubjectConfigSubject *CertificateConfigSubjectConfigSubject = &CertificateConfigSubjectConfigSubject{empty: true}

func (r *CertificateConfigSubjectConfigSubject) Empty() bool {
	return r.empty
}

func (r *CertificateConfigSubjectConfigSubject) String() string {
	return dcl.SprintResource(r)
}

func (r *CertificateConfigSubjectConfigSubject) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type CertificateConfigSubjectConfigSubjectAltName struct {
	empty          bool     `json:"-"`
	DnsNames       []string `json:"dnsNames"`
	Uris           []string `json:"uris"`
	EmailAddresses []string `json:"emailAddresses"`
	IPAddresses    []string `json:"ipAddresses"`
}

type jsonCertificateConfigSubjectConfigSubjectAltName CertificateConfigSubjectConfigSubjectAltName

func (r *CertificateConfigSubjectConfigSubjectAltName) UnmarshalJSON(data []byte) error {
	var res jsonCertificateConfigSubjectConfigSubjectAltName
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyCertificateConfigSubjectConfigSubjectAltName
	} else {

		r.DnsNames = res.DnsNames

		r.Uris = res.Uris

		r.EmailAddresses = res.EmailAddresses

		r.IPAddresses = res.IPAddresses

	}
	return nil
}

// This object is used to assert a desired state where this CertificateConfigSubjectConfigSubjectAltName is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyCertificateConfigSubjectConfigSubjectAltName *CertificateConfigSubjectConfigSubjectAltName = &CertificateConfigSubjectConfigSubjectAltName{empty: true}

func (r *CertificateConfigSubjectConfigSubjectAltName) Empty() bool {
	return r.empty
}

func (r *CertificateConfigSubjectConfigSubjectAltName) String() string {
	return dcl.SprintResource(r)
}

func (r *CertificateConfigSubjectConfigSubjectAltName) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type CertificateConfigX509Config struct {
	empty                bool                                              `json:"-"`
	KeyUsage             *CertificateConfigX509ConfigKeyUsage              `json:"keyUsage"`
	CaOptions            *CertificateConfigX509ConfigCaOptions             `json:"caOptions"`
	PolicyIds            []CertificateConfigX509ConfigPolicyIds            `json:"policyIds"`
	AiaOcspServers       []string                                          `json:"aiaOcspServers"`
	AdditionalExtensions []CertificateConfigX509ConfigAdditionalExtensions `json:"additionalExtensions"`
}

type jsonCertificateConfigX509Config CertificateConfigX509Config

func (r *CertificateConfigX509Config) UnmarshalJSON(data []byte) error {
	var res jsonCertificateConfigX509Config
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyCertificateConfigX509Config
	} else {

		r.KeyUsage = res.KeyUsage

		r.CaOptions = res.CaOptions

		r.PolicyIds = res.PolicyIds

		r.AiaOcspServers = res.AiaOcspServers

		r.AdditionalExtensions = res.AdditionalExtensions

	}
	return nil
}

// This object is used to assert a desired state where this CertificateConfigX509Config is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyCertificateConfigX509Config *CertificateConfigX509Config = &CertificateConfigX509Config{empty: true}

func (r *CertificateConfigX509Config) Empty() bool {
	return r.empty
}

func (r *CertificateConfigX509Config) String() string {
	return dcl.SprintResource(r)
}

func (r *CertificateConfigX509Config) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type CertificateConfigX509ConfigKeyUsage struct {
	empty                    bool                                                          `json:"-"`
	BaseKeyUsage             *CertificateConfigX509ConfigKeyUsageBaseKeyUsage              `json:"baseKeyUsage"`
	ExtendedKeyUsage         *CertificateConfigX509ConfigKeyUsageExtendedKeyUsage          `json:"extendedKeyUsage"`
	UnknownExtendedKeyUsages []CertificateConfigX509ConfigKeyUsageUnknownExtendedKeyUsages `json:"unknownExtendedKeyUsages"`
}

type jsonCertificateConfigX509ConfigKeyUsage CertificateConfigX509ConfigKeyUsage

func (r *CertificateConfigX509ConfigKeyUsage) UnmarshalJSON(data []byte) error {
	var res jsonCertificateConfigX509ConfigKeyUsage
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyCertificateConfigX509ConfigKeyUsage
	} else {

		r.BaseKeyUsage = res.BaseKeyUsage

		r.ExtendedKeyUsage = res.ExtendedKeyUsage

		r.UnknownExtendedKeyUsages = res.UnknownExtendedKeyUsages

	}
	return nil
}

// This object is used to assert a desired state where this CertificateConfigX509ConfigKeyUsage is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyCertificateConfigX509ConfigKeyUsage *CertificateConfigX509ConfigKeyUsage = &CertificateConfigX509ConfigKeyUsage{empty: true}

func (r *CertificateConfigX509ConfigKeyUsage) Empty() bool {
	return r.empty
}

func (r *CertificateConfigX509ConfigKeyUsage) String() string {
	return dcl.SprintResource(r)
}

func (r *CertificateConfigX509ConfigKeyUsage) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type CertificateConfigX509ConfigKeyUsageBaseKeyUsage struct {
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

type jsonCertificateConfigX509ConfigKeyUsageBaseKeyUsage CertificateConfigX509ConfigKeyUsageBaseKeyUsage

func (r *CertificateConfigX509ConfigKeyUsageBaseKeyUsage) UnmarshalJSON(data []byte) error {
	var res jsonCertificateConfigX509ConfigKeyUsageBaseKeyUsage
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyCertificateConfigX509ConfigKeyUsageBaseKeyUsage
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

// This object is used to assert a desired state where this CertificateConfigX509ConfigKeyUsageBaseKeyUsage is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyCertificateConfigX509ConfigKeyUsageBaseKeyUsage *CertificateConfigX509ConfigKeyUsageBaseKeyUsage = &CertificateConfigX509ConfigKeyUsageBaseKeyUsage{empty: true}

func (r *CertificateConfigX509ConfigKeyUsageBaseKeyUsage) Empty() bool {
	return r.empty
}

func (r *CertificateConfigX509ConfigKeyUsageBaseKeyUsage) String() string {
	return dcl.SprintResource(r)
}

func (r *CertificateConfigX509ConfigKeyUsageBaseKeyUsage) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type CertificateConfigX509ConfigKeyUsageExtendedKeyUsage struct {
	empty           bool  `json:"-"`
	ServerAuth      *bool `json:"serverAuth"`
	ClientAuth      *bool `json:"clientAuth"`
	CodeSigning     *bool `json:"codeSigning"`
	EmailProtection *bool `json:"emailProtection"`
	TimeStamping    *bool `json:"timeStamping"`
	OcspSigning     *bool `json:"ocspSigning"`
}

type jsonCertificateConfigX509ConfigKeyUsageExtendedKeyUsage CertificateConfigX509ConfigKeyUsageExtendedKeyUsage

func (r *CertificateConfigX509ConfigKeyUsageExtendedKeyUsage) UnmarshalJSON(data []byte) error {
	var res jsonCertificateConfigX509ConfigKeyUsageExtendedKeyUsage
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyCertificateConfigX509ConfigKeyUsageExtendedKeyUsage
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

// This object is used to assert a desired state where this CertificateConfigX509ConfigKeyUsageExtendedKeyUsage is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyCertificateConfigX509ConfigKeyUsageExtendedKeyUsage *CertificateConfigX509ConfigKeyUsageExtendedKeyUsage = &CertificateConfigX509ConfigKeyUsageExtendedKeyUsage{empty: true}

func (r *CertificateConfigX509ConfigKeyUsageExtendedKeyUsage) Empty() bool {
	return r.empty
}

func (r *CertificateConfigX509ConfigKeyUsageExtendedKeyUsage) String() string {
	return dcl.SprintResource(r)
}

func (r *CertificateConfigX509ConfigKeyUsageExtendedKeyUsage) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type CertificateConfigX509ConfigKeyUsageUnknownExtendedKeyUsages struct {
	empty        bool    `json:"-"`
	ObjectIdPath []int64 `json:"objectIdPath"`
}

type jsonCertificateConfigX509ConfigKeyUsageUnknownExtendedKeyUsages CertificateConfigX509ConfigKeyUsageUnknownExtendedKeyUsages

func (r *CertificateConfigX509ConfigKeyUsageUnknownExtendedKeyUsages) UnmarshalJSON(data []byte) error {
	var res jsonCertificateConfigX509ConfigKeyUsageUnknownExtendedKeyUsages
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyCertificateConfigX509ConfigKeyUsageUnknownExtendedKeyUsages
	} else {

		r.ObjectIdPath = res.ObjectIdPath

	}
	return nil
}

// This object is used to assert a desired state where this CertificateConfigX509ConfigKeyUsageUnknownExtendedKeyUsages is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyCertificateConfigX509ConfigKeyUsageUnknownExtendedKeyUsages *CertificateConfigX509ConfigKeyUsageUnknownExtendedKeyUsages = &CertificateConfigX509ConfigKeyUsageUnknownExtendedKeyUsages{empty: true}

func (r *CertificateConfigX509ConfigKeyUsageUnknownExtendedKeyUsages) Empty() bool {
	return r.empty
}

func (r *CertificateConfigX509ConfigKeyUsageUnknownExtendedKeyUsages) String() string {
	return dcl.SprintResource(r)
}

func (r *CertificateConfigX509ConfigKeyUsageUnknownExtendedKeyUsages) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type CertificateConfigX509ConfigCaOptions struct {
	empty                   bool   `json:"-"`
	IsCa                    *bool  `json:"isCa"`
	NonCa                   *bool  `json:"nonCa"`
	MaxIssuerPathLength     *int64 `json:"maxIssuerPathLength"`
	ZeroMaxIssuerPathLength *bool  `json:"zeroMaxIssuerPathLength"`
}

type jsonCertificateConfigX509ConfigCaOptions CertificateConfigX509ConfigCaOptions

func (r *CertificateConfigX509ConfigCaOptions) UnmarshalJSON(data []byte) error {
	var res jsonCertificateConfigX509ConfigCaOptions
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyCertificateConfigX509ConfigCaOptions
	} else {

		r.IsCa = res.IsCa

		r.NonCa = res.NonCa

		r.MaxIssuerPathLength = res.MaxIssuerPathLength

		r.ZeroMaxIssuerPathLength = res.ZeroMaxIssuerPathLength

	}
	return nil
}

// This object is used to assert a desired state where this CertificateConfigX509ConfigCaOptions is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyCertificateConfigX509ConfigCaOptions *CertificateConfigX509ConfigCaOptions = &CertificateConfigX509ConfigCaOptions{empty: true}

func (r *CertificateConfigX509ConfigCaOptions) Empty() bool {
	return r.empty
}

func (r *CertificateConfigX509ConfigCaOptions) String() string {
	return dcl.SprintResource(r)
}

func (r *CertificateConfigX509ConfigCaOptions) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type CertificateConfigX509ConfigPolicyIds struct {
	empty        bool    `json:"-"`
	ObjectIdPath []int64 `json:"objectIdPath"`
}

type jsonCertificateConfigX509ConfigPolicyIds CertificateConfigX509ConfigPolicyIds

func (r *CertificateConfigX509ConfigPolicyIds) UnmarshalJSON(data []byte) error {
	var res jsonCertificateConfigX509ConfigPolicyIds
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyCertificateConfigX509ConfigPolicyIds
	} else {

		r.ObjectIdPath = res.ObjectIdPath

	}
	return nil
}

// This object is used to assert a desired state where this CertificateConfigX509ConfigPolicyIds is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyCertificateConfigX509ConfigPolicyIds *CertificateConfigX509ConfigPolicyIds = &CertificateConfigX509ConfigPolicyIds{empty: true}

func (r *CertificateConfigX509ConfigPolicyIds) Empty() bool {
	return r.empty
}

func (r *CertificateConfigX509ConfigPolicyIds) String() string {
	return dcl.SprintResource(r)
}

func (r *CertificateConfigX509ConfigPolicyIds) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type CertificateConfigX509ConfigAdditionalExtensions struct {
	empty    bool                                                     `json:"-"`
	ObjectId *CertificateConfigX509ConfigAdditionalExtensionsObjectId `json:"objectId"`
	Critical *bool                                                    `json:"critical"`
	Value    *string                                                  `json:"value"`
}

type jsonCertificateConfigX509ConfigAdditionalExtensions CertificateConfigX509ConfigAdditionalExtensions

func (r *CertificateConfigX509ConfigAdditionalExtensions) UnmarshalJSON(data []byte) error {
	var res jsonCertificateConfigX509ConfigAdditionalExtensions
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyCertificateConfigX509ConfigAdditionalExtensions
	} else {

		r.ObjectId = res.ObjectId

		r.Critical = res.Critical

		r.Value = res.Value

	}
	return nil
}

// This object is used to assert a desired state where this CertificateConfigX509ConfigAdditionalExtensions is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyCertificateConfigX509ConfigAdditionalExtensions *CertificateConfigX509ConfigAdditionalExtensions = &CertificateConfigX509ConfigAdditionalExtensions{empty: true}

func (r *CertificateConfigX509ConfigAdditionalExtensions) Empty() bool {
	return r.empty
}

func (r *CertificateConfigX509ConfigAdditionalExtensions) String() string {
	return dcl.SprintResource(r)
}

func (r *CertificateConfigX509ConfigAdditionalExtensions) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type CertificateConfigX509ConfigAdditionalExtensionsObjectId struct {
	empty        bool    `json:"-"`
	ObjectIdPath []int64 `json:"objectIdPath"`
}

type jsonCertificateConfigX509ConfigAdditionalExtensionsObjectId CertificateConfigX509ConfigAdditionalExtensionsObjectId

func (r *CertificateConfigX509ConfigAdditionalExtensionsObjectId) UnmarshalJSON(data []byte) error {
	var res jsonCertificateConfigX509ConfigAdditionalExtensionsObjectId
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyCertificateConfigX509ConfigAdditionalExtensionsObjectId
	} else {

		r.ObjectIdPath = res.ObjectIdPath

	}
	return nil
}

// This object is used to assert a desired state where this CertificateConfigX509ConfigAdditionalExtensionsObjectId is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyCertificateConfigX509ConfigAdditionalExtensionsObjectId *CertificateConfigX509ConfigAdditionalExtensionsObjectId = &CertificateConfigX509ConfigAdditionalExtensionsObjectId{empty: true}

func (r *CertificateConfigX509ConfigAdditionalExtensionsObjectId) Empty() bool {
	return r.empty
}

func (r *CertificateConfigX509ConfigAdditionalExtensionsObjectId) String() string {
	return dcl.SprintResource(r)
}

func (r *CertificateConfigX509ConfigAdditionalExtensionsObjectId) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type CertificateConfigPublicKey struct {
	empty  bool                                  `json:"-"`
	Key    *string                               `json:"key"`
	Format *CertificateConfigPublicKeyFormatEnum `json:"format"`
}

type jsonCertificateConfigPublicKey CertificateConfigPublicKey

func (r *CertificateConfigPublicKey) UnmarshalJSON(data []byte) error {
	var res jsonCertificateConfigPublicKey
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyCertificateConfigPublicKey
	} else {

		r.Key = res.Key

		r.Format = res.Format

	}
	return nil
}

// This object is used to assert a desired state where this CertificateConfigPublicKey is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyCertificateConfigPublicKey *CertificateConfigPublicKey = &CertificateConfigPublicKey{empty: true}

func (r *CertificateConfigPublicKey) Empty() bool {
	return r.empty
}

func (r *CertificateConfigPublicKey) String() string {
	return dcl.SprintResource(r)
}

func (r *CertificateConfigPublicKey) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type CertificateRevocationDetails struct {
	empty           bool                                             `json:"-"`
	RevocationState *CertificateRevocationDetailsRevocationStateEnum `json:"revocationState"`
	RevocationTime  *string                                          `json:"revocationTime"`
}

type jsonCertificateRevocationDetails CertificateRevocationDetails

func (r *CertificateRevocationDetails) UnmarshalJSON(data []byte) error {
	var res jsonCertificateRevocationDetails
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyCertificateRevocationDetails
	} else {

		r.RevocationState = res.RevocationState

		r.RevocationTime = res.RevocationTime

	}
	return nil
}

// This object is used to assert a desired state where this CertificateRevocationDetails is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyCertificateRevocationDetails *CertificateRevocationDetails = &CertificateRevocationDetails{empty: true}

func (r *CertificateRevocationDetails) Empty() bool {
	return r.empty
}

func (r *CertificateRevocationDetails) String() string {
	return dcl.SprintResource(r)
}

func (r *CertificateRevocationDetails) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type CertificateCertificateDescription struct {
	empty                     bool                                                 `json:"-"`
	SubjectDescription        *CertificateCertificateDescriptionSubjectDescription `json:"subjectDescription"`
	X509Description           *CertificateCertificateDescriptionX509Description    `json:"x509Description"`
	PublicKey                 *CertificateCertificateDescriptionPublicKey          `json:"publicKey"`
	SubjectKeyId              *CertificateCertificateDescriptionSubjectKeyId       `json:"subjectKeyId"`
	AuthorityKeyId            *CertificateCertificateDescriptionAuthorityKeyId     `json:"authorityKeyId"`
	CrlDistributionPoints     []string                                             `json:"crlDistributionPoints"`
	AiaIssuingCertificateUrls []string                                             `json:"aiaIssuingCertificateUrls"`
	CertFingerprint           *CertificateCertificateDescriptionCertFingerprint    `json:"certFingerprint"`
}

type jsonCertificateCertificateDescription CertificateCertificateDescription

func (r *CertificateCertificateDescription) UnmarshalJSON(data []byte) error {
	var res jsonCertificateCertificateDescription
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyCertificateCertificateDescription
	} else {

		r.SubjectDescription = res.SubjectDescription

		r.X509Description = res.X509Description

		r.PublicKey = res.PublicKey

		r.SubjectKeyId = res.SubjectKeyId

		r.AuthorityKeyId = res.AuthorityKeyId

		r.CrlDistributionPoints = res.CrlDistributionPoints

		r.AiaIssuingCertificateUrls = res.AiaIssuingCertificateUrls

		r.CertFingerprint = res.CertFingerprint

	}
	return nil
}

// This object is used to assert a desired state where this CertificateCertificateDescription is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyCertificateCertificateDescription *CertificateCertificateDescription = &CertificateCertificateDescription{empty: true}

func (r *CertificateCertificateDescription) Empty() bool {
	return r.empty
}

func (r *CertificateCertificateDescription) String() string {
	return dcl.SprintResource(r)
}

func (r *CertificateCertificateDescription) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type CertificateCertificateDescriptionSubjectDescription struct {
	empty           bool                                                               `json:"-"`
	Subject         *CertificateCertificateDescriptionSubjectDescriptionSubject        `json:"subject"`
	SubjectAltName  *CertificateCertificateDescriptionSubjectDescriptionSubjectAltName `json:"subjectAltName"`
	HexSerialNumber *string                                                            `json:"hexSerialNumber"`
	Lifetime        *string                                                            `json:"lifetime"`
	NotBeforeTime   *string                                                            `json:"notBeforeTime"`
	NotAfterTime    *string                                                            `json:"notAfterTime"`
}

type jsonCertificateCertificateDescriptionSubjectDescription CertificateCertificateDescriptionSubjectDescription

func (r *CertificateCertificateDescriptionSubjectDescription) UnmarshalJSON(data []byte) error {
	var res jsonCertificateCertificateDescriptionSubjectDescription
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyCertificateCertificateDescriptionSubjectDescription
	} else {

		r.Subject = res.Subject

		r.SubjectAltName = res.SubjectAltName

		r.HexSerialNumber = res.HexSerialNumber

		r.Lifetime = res.Lifetime

		r.NotBeforeTime = res.NotBeforeTime

		r.NotAfterTime = res.NotAfterTime

	}
	return nil
}

// This object is used to assert a desired state where this CertificateCertificateDescriptionSubjectDescription is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyCertificateCertificateDescriptionSubjectDescription *CertificateCertificateDescriptionSubjectDescription = &CertificateCertificateDescriptionSubjectDescription{empty: true}

func (r *CertificateCertificateDescriptionSubjectDescription) Empty() bool {
	return r.empty
}

func (r *CertificateCertificateDescriptionSubjectDescription) String() string {
	return dcl.SprintResource(r)
}

func (r *CertificateCertificateDescriptionSubjectDescription) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type CertificateCertificateDescriptionSubjectDescriptionSubject struct {
	empty              bool    `json:"-"`
	CommonName         *string `json:"commonName"`
	CountryCode        *string `json:"countryCode"`
	Organization       *string `json:"organization"`
	OrganizationalUnit *string `json:"organizationalUnit"`
	Locality           *string `json:"locality"`
	Province           *string `json:"province"`
	StreetAddress      *string `json:"streetAddress"`
	PostalCode         *string `json:"postalCode"`
}

type jsonCertificateCertificateDescriptionSubjectDescriptionSubject CertificateCertificateDescriptionSubjectDescriptionSubject

func (r *CertificateCertificateDescriptionSubjectDescriptionSubject) UnmarshalJSON(data []byte) error {
	var res jsonCertificateCertificateDescriptionSubjectDescriptionSubject
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyCertificateCertificateDescriptionSubjectDescriptionSubject
	} else {

		r.CommonName = res.CommonName

		r.CountryCode = res.CountryCode

		r.Organization = res.Organization

		r.OrganizationalUnit = res.OrganizationalUnit

		r.Locality = res.Locality

		r.Province = res.Province

		r.StreetAddress = res.StreetAddress

		r.PostalCode = res.PostalCode

	}
	return nil
}

// This object is used to assert a desired state where this CertificateCertificateDescriptionSubjectDescriptionSubject is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyCertificateCertificateDescriptionSubjectDescriptionSubject *CertificateCertificateDescriptionSubjectDescriptionSubject = &CertificateCertificateDescriptionSubjectDescriptionSubject{empty: true}

func (r *CertificateCertificateDescriptionSubjectDescriptionSubject) Empty() bool {
	return r.empty
}

func (r *CertificateCertificateDescriptionSubjectDescriptionSubject) String() string {
	return dcl.SprintResource(r)
}

func (r *CertificateCertificateDescriptionSubjectDescriptionSubject) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type CertificateCertificateDescriptionSubjectDescriptionSubjectAltName struct {
	empty          bool                                                                          `json:"-"`
	DnsNames       []string                                                                      `json:"dnsNames"`
	Uris           []string                                                                      `json:"uris"`
	EmailAddresses []string                                                                      `json:"emailAddresses"`
	IPAddresses    []string                                                                      `json:"ipAddresses"`
	CustomSans     []CertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSans `json:"customSans"`
}

type jsonCertificateCertificateDescriptionSubjectDescriptionSubjectAltName CertificateCertificateDescriptionSubjectDescriptionSubjectAltName

func (r *CertificateCertificateDescriptionSubjectDescriptionSubjectAltName) UnmarshalJSON(data []byte) error {
	var res jsonCertificateCertificateDescriptionSubjectDescriptionSubjectAltName
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyCertificateCertificateDescriptionSubjectDescriptionSubjectAltName
	} else {

		r.DnsNames = res.DnsNames

		r.Uris = res.Uris

		r.EmailAddresses = res.EmailAddresses

		r.IPAddresses = res.IPAddresses

		r.CustomSans = res.CustomSans

	}
	return nil
}

// This object is used to assert a desired state where this CertificateCertificateDescriptionSubjectDescriptionSubjectAltName is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyCertificateCertificateDescriptionSubjectDescriptionSubjectAltName *CertificateCertificateDescriptionSubjectDescriptionSubjectAltName = &CertificateCertificateDescriptionSubjectDescriptionSubjectAltName{empty: true}

func (r *CertificateCertificateDescriptionSubjectDescriptionSubjectAltName) Empty() bool {
	return r.empty
}

func (r *CertificateCertificateDescriptionSubjectDescriptionSubjectAltName) String() string {
	return dcl.SprintResource(r)
}

func (r *CertificateCertificateDescriptionSubjectDescriptionSubjectAltName) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type CertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSans struct {
	empty    bool                                                                                 `json:"-"`
	ObjectId *CertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSansObjectId `json:"objectId"`
	Critical *bool                                                                                `json:"critical"`
	Value    *string                                                                              `json:"value"`
}

type jsonCertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSans CertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSans

func (r *CertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSans) UnmarshalJSON(data []byte) error {
	var res jsonCertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSans
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyCertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSans
	} else {

		r.ObjectId = res.ObjectId

		r.Critical = res.Critical

		r.Value = res.Value

	}
	return nil
}

// This object is used to assert a desired state where this CertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSans is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyCertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSans *CertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSans = &CertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSans{empty: true}

func (r *CertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSans) Empty() bool {
	return r.empty
}

func (r *CertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSans) String() string {
	return dcl.SprintResource(r)
}

func (r *CertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSans) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type CertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSansObjectId struct {
	empty        bool    `json:"-"`
	ObjectIdPath []int64 `json:"objectIdPath"`
}

type jsonCertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSansObjectId CertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSansObjectId

func (r *CertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSansObjectId) UnmarshalJSON(data []byte) error {
	var res jsonCertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSansObjectId
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyCertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSansObjectId
	} else {

		r.ObjectIdPath = res.ObjectIdPath

	}
	return nil
}

// This object is used to assert a desired state where this CertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSansObjectId is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyCertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSansObjectId *CertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSansObjectId = &CertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSansObjectId{empty: true}

func (r *CertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSansObjectId) Empty() bool {
	return r.empty
}

func (r *CertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSansObjectId) String() string {
	return dcl.SprintResource(r)
}

func (r *CertificateCertificateDescriptionSubjectDescriptionSubjectAltNameCustomSansObjectId) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type CertificateCertificateDescriptionX509Description struct {
	empty                bool                                                                   `json:"-"`
	KeyUsage             *CertificateCertificateDescriptionX509DescriptionKeyUsage              `json:"keyUsage"`
	CaOptions            *CertificateCertificateDescriptionX509DescriptionCaOptions             `json:"caOptions"`
	PolicyIds            []CertificateCertificateDescriptionX509DescriptionPolicyIds            `json:"policyIds"`
	AiaOcspServers       []string                                                               `json:"aiaOcspServers"`
	AdditionalExtensions []CertificateCertificateDescriptionX509DescriptionAdditionalExtensions `json:"additionalExtensions"`
}

type jsonCertificateCertificateDescriptionX509Description CertificateCertificateDescriptionX509Description

func (r *CertificateCertificateDescriptionX509Description) UnmarshalJSON(data []byte) error {
	var res jsonCertificateCertificateDescriptionX509Description
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyCertificateCertificateDescriptionX509Description
	} else {

		r.KeyUsage = res.KeyUsage

		r.CaOptions = res.CaOptions

		r.PolicyIds = res.PolicyIds

		r.AiaOcspServers = res.AiaOcspServers

		r.AdditionalExtensions = res.AdditionalExtensions

	}
	return nil
}

// This object is used to assert a desired state where this CertificateCertificateDescriptionX509Description is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyCertificateCertificateDescriptionX509Description *CertificateCertificateDescriptionX509Description = &CertificateCertificateDescriptionX509Description{empty: true}

func (r *CertificateCertificateDescriptionX509Description) Empty() bool {
	return r.empty
}

func (r *CertificateCertificateDescriptionX509Description) String() string {
	return dcl.SprintResource(r)
}

func (r *CertificateCertificateDescriptionX509Description) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type CertificateCertificateDescriptionX509DescriptionKeyUsage struct {
	empty                    bool                                                                               `json:"-"`
	BaseKeyUsage             *CertificateCertificateDescriptionX509DescriptionKeyUsageBaseKeyUsage              `json:"baseKeyUsage"`
	ExtendedKeyUsage         *CertificateCertificateDescriptionX509DescriptionKeyUsageExtendedKeyUsage          `json:"extendedKeyUsage"`
	UnknownExtendedKeyUsages []CertificateCertificateDescriptionX509DescriptionKeyUsageUnknownExtendedKeyUsages `json:"unknownExtendedKeyUsages"`
}

type jsonCertificateCertificateDescriptionX509DescriptionKeyUsage CertificateCertificateDescriptionX509DescriptionKeyUsage

func (r *CertificateCertificateDescriptionX509DescriptionKeyUsage) UnmarshalJSON(data []byte) error {
	var res jsonCertificateCertificateDescriptionX509DescriptionKeyUsage
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyCertificateCertificateDescriptionX509DescriptionKeyUsage
	} else {

		r.BaseKeyUsage = res.BaseKeyUsage

		r.ExtendedKeyUsage = res.ExtendedKeyUsage

		r.UnknownExtendedKeyUsages = res.UnknownExtendedKeyUsages

	}
	return nil
}

// This object is used to assert a desired state where this CertificateCertificateDescriptionX509DescriptionKeyUsage is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyCertificateCertificateDescriptionX509DescriptionKeyUsage *CertificateCertificateDescriptionX509DescriptionKeyUsage = &CertificateCertificateDescriptionX509DescriptionKeyUsage{empty: true}

func (r *CertificateCertificateDescriptionX509DescriptionKeyUsage) Empty() bool {
	return r.empty
}

func (r *CertificateCertificateDescriptionX509DescriptionKeyUsage) String() string {
	return dcl.SprintResource(r)
}

func (r *CertificateCertificateDescriptionX509DescriptionKeyUsage) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type CertificateCertificateDescriptionX509DescriptionKeyUsageBaseKeyUsage struct {
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

type jsonCertificateCertificateDescriptionX509DescriptionKeyUsageBaseKeyUsage CertificateCertificateDescriptionX509DescriptionKeyUsageBaseKeyUsage

func (r *CertificateCertificateDescriptionX509DescriptionKeyUsageBaseKeyUsage) UnmarshalJSON(data []byte) error {
	var res jsonCertificateCertificateDescriptionX509DescriptionKeyUsageBaseKeyUsage
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyCertificateCertificateDescriptionX509DescriptionKeyUsageBaseKeyUsage
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

// This object is used to assert a desired state where this CertificateCertificateDescriptionX509DescriptionKeyUsageBaseKeyUsage is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyCertificateCertificateDescriptionX509DescriptionKeyUsageBaseKeyUsage *CertificateCertificateDescriptionX509DescriptionKeyUsageBaseKeyUsage = &CertificateCertificateDescriptionX509DescriptionKeyUsageBaseKeyUsage{empty: true}

func (r *CertificateCertificateDescriptionX509DescriptionKeyUsageBaseKeyUsage) Empty() bool {
	return r.empty
}

func (r *CertificateCertificateDescriptionX509DescriptionKeyUsageBaseKeyUsage) String() string {
	return dcl.SprintResource(r)
}

func (r *CertificateCertificateDescriptionX509DescriptionKeyUsageBaseKeyUsage) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type CertificateCertificateDescriptionX509DescriptionKeyUsageExtendedKeyUsage struct {
	empty           bool  `json:"-"`
	ServerAuth      *bool `json:"serverAuth"`
	ClientAuth      *bool `json:"clientAuth"`
	CodeSigning     *bool `json:"codeSigning"`
	EmailProtection *bool `json:"emailProtection"`
	TimeStamping    *bool `json:"timeStamping"`
	OcspSigning     *bool `json:"ocspSigning"`
}

type jsonCertificateCertificateDescriptionX509DescriptionKeyUsageExtendedKeyUsage CertificateCertificateDescriptionX509DescriptionKeyUsageExtendedKeyUsage

func (r *CertificateCertificateDescriptionX509DescriptionKeyUsageExtendedKeyUsage) UnmarshalJSON(data []byte) error {
	var res jsonCertificateCertificateDescriptionX509DescriptionKeyUsageExtendedKeyUsage
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyCertificateCertificateDescriptionX509DescriptionKeyUsageExtendedKeyUsage
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

// This object is used to assert a desired state where this CertificateCertificateDescriptionX509DescriptionKeyUsageExtendedKeyUsage is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyCertificateCertificateDescriptionX509DescriptionKeyUsageExtendedKeyUsage *CertificateCertificateDescriptionX509DescriptionKeyUsageExtendedKeyUsage = &CertificateCertificateDescriptionX509DescriptionKeyUsageExtendedKeyUsage{empty: true}

func (r *CertificateCertificateDescriptionX509DescriptionKeyUsageExtendedKeyUsage) Empty() bool {
	return r.empty
}

func (r *CertificateCertificateDescriptionX509DescriptionKeyUsageExtendedKeyUsage) String() string {
	return dcl.SprintResource(r)
}

func (r *CertificateCertificateDescriptionX509DescriptionKeyUsageExtendedKeyUsage) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type CertificateCertificateDescriptionX509DescriptionKeyUsageUnknownExtendedKeyUsages struct {
	empty        bool    `json:"-"`
	ObjectIdPath []int64 `json:"objectIdPath"`
}

type jsonCertificateCertificateDescriptionX509DescriptionKeyUsageUnknownExtendedKeyUsages CertificateCertificateDescriptionX509DescriptionKeyUsageUnknownExtendedKeyUsages

func (r *CertificateCertificateDescriptionX509DescriptionKeyUsageUnknownExtendedKeyUsages) UnmarshalJSON(data []byte) error {
	var res jsonCertificateCertificateDescriptionX509DescriptionKeyUsageUnknownExtendedKeyUsages
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyCertificateCertificateDescriptionX509DescriptionKeyUsageUnknownExtendedKeyUsages
	} else {

		r.ObjectIdPath = res.ObjectIdPath

	}
	return nil
}

// This object is used to assert a desired state where this CertificateCertificateDescriptionX509DescriptionKeyUsageUnknownExtendedKeyUsages is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyCertificateCertificateDescriptionX509DescriptionKeyUsageUnknownExtendedKeyUsages *CertificateCertificateDescriptionX509DescriptionKeyUsageUnknownExtendedKeyUsages = &CertificateCertificateDescriptionX509DescriptionKeyUsageUnknownExtendedKeyUsages{empty: true}

func (r *CertificateCertificateDescriptionX509DescriptionKeyUsageUnknownExtendedKeyUsages) Empty() bool {
	return r.empty
}

func (r *CertificateCertificateDescriptionX509DescriptionKeyUsageUnknownExtendedKeyUsages) String() string {
	return dcl.SprintResource(r)
}

func (r *CertificateCertificateDescriptionX509DescriptionKeyUsageUnknownExtendedKeyUsages) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type CertificateCertificateDescriptionX509DescriptionCaOptions struct {
	empty               bool   `json:"-"`
	IsCa                *bool  `json:"isCa"`
	MaxIssuerPathLength *int64 `json:"maxIssuerPathLength"`
}

type jsonCertificateCertificateDescriptionX509DescriptionCaOptions CertificateCertificateDescriptionX509DescriptionCaOptions

func (r *CertificateCertificateDescriptionX509DescriptionCaOptions) UnmarshalJSON(data []byte) error {
	var res jsonCertificateCertificateDescriptionX509DescriptionCaOptions
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyCertificateCertificateDescriptionX509DescriptionCaOptions
	} else {

		r.IsCa = res.IsCa

		r.MaxIssuerPathLength = res.MaxIssuerPathLength

	}
	return nil
}

// This object is used to assert a desired state where this CertificateCertificateDescriptionX509DescriptionCaOptions is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyCertificateCertificateDescriptionX509DescriptionCaOptions *CertificateCertificateDescriptionX509DescriptionCaOptions = &CertificateCertificateDescriptionX509DescriptionCaOptions{empty: true}

func (r *CertificateCertificateDescriptionX509DescriptionCaOptions) Empty() bool {
	return r.empty
}

func (r *CertificateCertificateDescriptionX509DescriptionCaOptions) String() string {
	return dcl.SprintResource(r)
}

func (r *CertificateCertificateDescriptionX509DescriptionCaOptions) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type CertificateCertificateDescriptionX509DescriptionPolicyIds struct {
	empty        bool    `json:"-"`
	ObjectIdPath []int64 `json:"objectIdPath"`
}

type jsonCertificateCertificateDescriptionX509DescriptionPolicyIds CertificateCertificateDescriptionX509DescriptionPolicyIds

func (r *CertificateCertificateDescriptionX509DescriptionPolicyIds) UnmarshalJSON(data []byte) error {
	var res jsonCertificateCertificateDescriptionX509DescriptionPolicyIds
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyCertificateCertificateDescriptionX509DescriptionPolicyIds
	} else {

		r.ObjectIdPath = res.ObjectIdPath

	}
	return nil
}

// This object is used to assert a desired state where this CertificateCertificateDescriptionX509DescriptionPolicyIds is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyCertificateCertificateDescriptionX509DescriptionPolicyIds *CertificateCertificateDescriptionX509DescriptionPolicyIds = &CertificateCertificateDescriptionX509DescriptionPolicyIds{empty: true}

func (r *CertificateCertificateDescriptionX509DescriptionPolicyIds) Empty() bool {
	return r.empty
}

func (r *CertificateCertificateDescriptionX509DescriptionPolicyIds) String() string {
	return dcl.SprintResource(r)
}

func (r *CertificateCertificateDescriptionX509DescriptionPolicyIds) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type CertificateCertificateDescriptionX509DescriptionAdditionalExtensions struct {
	empty    bool                                                                          `json:"-"`
	ObjectId *CertificateCertificateDescriptionX509DescriptionAdditionalExtensionsObjectId `json:"objectId"`
	Critical *bool                                                                         `json:"critical"`
	Value    *string                                                                       `json:"value"`
}

type jsonCertificateCertificateDescriptionX509DescriptionAdditionalExtensions CertificateCertificateDescriptionX509DescriptionAdditionalExtensions

func (r *CertificateCertificateDescriptionX509DescriptionAdditionalExtensions) UnmarshalJSON(data []byte) error {
	var res jsonCertificateCertificateDescriptionX509DescriptionAdditionalExtensions
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyCertificateCertificateDescriptionX509DescriptionAdditionalExtensions
	} else {

		r.ObjectId = res.ObjectId

		r.Critical = res.Critical

		r.Value = res.Value

	}
	return nil
}

// This object is used to assert a desired state where this CertificateCertificateDescriptionX509DescriptionAdditionalExtensions is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyCertificateCertificateDescriptionX509DescriptionAdditionalExtensions *CertificateCertificateDescriptionX509DescriptionAdditionalExtensions = &CertificateCertificateDescriptionX509DescriptionAdditionalExtensions{empty: true}

func (r *CertificateCertificateDescriptionX509DescriptionAdditionalExtensions) Empty() bool {
	return r.empty
}

func (r *CertificateCertificateDescriptionX509DescriptionAdditionalExtensions) String() string {
	return dcl.SprintResource(r)
}

func (r *CertificateCertificateDescriptionX509DescriptionAdditionalExtensions) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type CertificateCertificateDescriptionX509DescriptionAdditionalExtensionsObjectId struct {
	empty        bool    `json:"-"`
	ObjectIdPath []int64 `json:"objectIdPath"`
}

type jsonCertificateCertificateDescriptionX509DescriptionAdditionalExtensionsObjectId CertificateCertificateDescriptionX509DescriptionAdditionalExtensionsObjectId

func (r *CertificateCertificateDescriptionX509DescriptionAdditionalExtensionsObjectId) UnmarshalJSON(data []byte) error {
	var res jsonCertificateCertificateDescriptionX509DescriptionAdditionalExtensionsObjectId
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyCertificateCertificateDescriptionX509DescriptionAdditionalExtensionsObjectId
	} else {

		r.ObjectIdPath = res.ObjectIdPath

	}
	return nil
}

// This object is used to assert a desired state where this CertificateCertificateDescriptionX509DescriptionAdditionalExtensionsObjectId is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyCertificateCertificateDescriptionX509DescriptionAdditionalExtensionsObjectId *CertificateCertificateDescriptionX509DescriptionAdditionalExtensionsObjectId = &CertificateCertificateDescriptionX509DescriptionAdditionalExtensionsObjectId{empty: true}

func (r *CertificateCertificateDescriptionX509DescriptionAdditionalExtensionsObjectId) Empty() bool {
	return r.empty
}

func (r *CertificateCertificateDescriptionX509DescriptionAdditionalExtensionsObjectId) String() string {
	return dcl.SprintResource(r)
}

func (r *CertificateCertificateDescriptionX509DescriptionAdditionalExtensionsObjectId) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type CertificateCertificateDescriptionPublicKey struct {
	empty  bool                                                  `json:"-"`
	Key    *string                                               `json:"key"`
	Format *CertificateCertificateDescriptionPublicKeyFormatEnum `json:"format"`
}

type jsonCertificateCertificateDescriptionPublicKey CertificateCertificateDescriptionPublicKey

func (r *CertificateCertificateDescriptionPublicKey) UnmarshalJSON(data []byte) error {
	var res jsonCertificateCertificateDescriptionPublicKey
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyCertificateCertificateDescriptionPublicKey
	} else {

		r.Key = res.Key

		r.Format = res.Format

	}
	return nil
}

// This object is used to assert a desired state where this CertificateCertificateDescriptionPublicKey is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyCertificateCertificateDescriptionPublicKey *CertificateCertificateDescriptionPublicKey = &CertificateCertificateDescriptionPublicKey{empty: true}

func (r *CertificateCertificateDescriptionPublicKey) Empty() bool {
	return r.empty
}

func (r *CertificateCertificateDescriptionPublicKey) String() string {
	return dcl.SprintResource(r)
}

func (r *CertificateCertificateDescriptionPublicKey) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type CertificateCertificateDescriptionSubjectKeyId struct {
	empty bool    `json:"-"`
	KeyId *string `json:"keyId"`
}

type jsonCertificateCertificateDescriptionSubjectKeyId CertificateCertificateDescriptionSubjectKeyId

func (r *CertificateCertificateDescriptionSubjectKeyId) UnmarshalJSON(data []byte) error {
	var res jsonCertificateCertificateDescriptionSubjectKeyId
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyCertificateCertificateDescriptionSubjectKeyId
	} else {

		r.KeyId = res.KeyId

	}
	return nil
}

// This object is used to assert a desired state where this CertificateCertificateDescriptionSubjectKeyId is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyCertificateCertificateDescriptionSubjectKeyId *CertificateCertificateDescriptionSubjectKeyId = &CertificateCertificateDescriptionSubjectKeyId{empty: true}

func (r *CertificateCertificateDescriptionSubjectKeyId) Empty() bool {
	return r.empty
}

func (r *CertificateCertificateDescriptionSubjectKeyId) String() string {
	return dcl.SprintResource(r)
}

func (r *CertificateCertificateDescriptionSubjectKeyId) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type CertificateCertificateDescriptionAuthorityKeyId struct {
	empty bool    `json:"-"`
	KeyId *string `json:"keyId"`
}

type jsonCertificateCertificateDescriptionAuthorityKeyId CertificateCertificateDescriptionAuthorityKeyId

func (r *CertificateCertificateDescriptionAuthorityKeyId) UnmarshalJSON(data []byte) error {
	var res jsonCertificateCertificateDescriptionAuthorityKeyId
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyCertificateCertificateDescriptionAuthorityKeyId
	} else {

		r.KeyId = res.KeyId

	}
	return nil
}

// This object is used to assert a desired state where this CertificateCertificateDescriptionAuthorityKeyId is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyCertificateCertificateDescriptionAuthorityKeyId *CertificateCertificateDescriptionAuthorityKeyId = &CertificateCertificateDescriptionAuthorityKeyId{empty: true}

func (r *CertificateCertificateDescriptionAuthorityKeyId) Empty() bool {
	return r.empty
}

func (r *CertificateCertificateDescriptionAuthorityKeyId) String() string {
	return dcl.SprintResource(r)
}

func (r *CertificateCertificateDescriptionAuthorityKeyId) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type CertificateCertificateDescriptionCertFingerprint struct {
	empty      bool    `json:"-"`
	Sha256Hash *string `json:"sha256Hash"`
}

type jsonCertificateCertificateDescriptionCertFingerprint CertificateCertificateDescriptionCertFingerprint

func (r *CertificateCertificateDescriptionCertFingerprint) UnmarshalJSON(data []byte) error {
	var res jsonCertificateCertificateDescriptionCertFingerprint
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyCertificateCertificateDescriptionCertFingerprint
	} else {

		r.Sha256Hash = res.Sha256Hash

	}
	return nil
}

// This object is used to assert a desired state where this CertificateCertificateDescriptionCertFingerprint is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyCertificateCertificateDescriptionCertFingerprint *CertificateCertificateDescriptionCertFingerprint = &CertificateCertificateDescriptionCertFingerprint{empty: true}

func (r *CertificateCertificateDescriptionCertFingerprint) Empty() bool {
	return r.empty
}

func (r *CertificateCertificateDescriptionCertFingerprint) String() string {
	return dcl.SprintResource(r)
}

func (r *CertificateCertificateDescriptionCertFingerprint) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

// Describe returns a simple description of this resource to ensure that automated tools
// can identify it.
func (r *Certificate) Describe() dcl.ServiceTypeVersion {
	return dcl.ServiceTypeVersion{
		Service: "privateca",
		Type:    "Certificate",
		Version: "privateca",
	}
}

func (r *Certificate) ID() (string, error) {
	if err := extractCertificateFields(r); err != nil {
		return "", err
	}
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"name":                         dcl.ValueOrEmptyString(nr.Name),
		"pem_csr":                      dcl.ValueOrEmptyString(nr.PemCsr),
		"config":                       dcl.ValueOrEmptyString(nr.Config),
		"issuer_certificate_authority": dcl.ValueOrEmptyString(nr.IssuerCertificateAuthority),
		"lifetime":                     dcl.ValueOrEmptyString(nr.Lifetime),
		"certificate_template":         dcl.ValueOrEmptyString(nr.CertificateTemplate),
		"subject_mode":                 dcl.ValueOrEmptyString(nr.SubjectMode),
		"revocation_details":           dcl.ValueOrEmptyString(nr.RevocationDetails),
		"pem_certificate":              dcl.ValueOrEmptyString(nr.PemCertificate),
		"certificate_description":      dcl.ValueOrEmptyString(nr.CertificateDescription),
		"pem_certificate_chain":        dcl.ValueOrEmptyString(nr.PemCertificateChain),
		"create_time":                  dcl.ValueOrEmptyString(nr.CreateTime),
		"update_time":                  dcl.ValueOrEmptyString(nr.UpdateTime),
		"labels":                       dcl.ValueOrEmptyString(nr.Labels),
		"project":                      dcl.ValueOrEmptyString(nr.Project),
		"location":                     dcl.ValueOrEmptyString(nr.Location),
		"ca_pool":                      dcl.ValueOrEmptyString(nr.CaPool),
		"certificate_authority":        dcl.ValueOrEmptyString(nr.CertificateAuthority),
	}
	return dcl.Nprintf("projects/{{project}}/locations/{{location}}/caPools/{{ca_pool}}/certificates/{{name}}", params), nil
}

const CertificateMaxPage = -1

type CertificateList struct {
	Items []*Certificate

	nextToken string

	pageSize int32

	resource *Certificate
}

func (l *CertificateList) HasNext() bool {
	return l.nextToken != ""
}

func (l *CertificateList) Next(ctx context.Context, c *Client) error {
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	if !l.HasNext() {
		return fmt.Errorf("no next page")
	}
	items, token, err := c.listCertificate(ctx, l.resource, l.nextToken, l.pageSize)
	if err != nil {
		return err
	}
	l.Items = items
	l.nextToken = token
	return err
}

func (c *Client) ListCertificate(ctx context.Context, project, location, caPool string) (*CertificateList, error) {
	ctx = dcl.ContextWithRequestID(ctx)
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	return c.ListCertificateWithMaxResults(ctx, project, location, caPool, CertificateMaxPage)

}

func (c *Client) ListCertificateWithMaxResults(ctx context.Context, project, location, caPool string, pageSize int32) (*CertificateList, error) {
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	// Create a resource object so that we can use proper url normalization methods.
	r := &Certificate{
		Project:  &project,
		Location: &location,
		CaPool:   &caPool,
	}
	items, token, err := c.listCertificate(ctx, r, "", pageSize)
	if err != nil {
		return nil, err
	}
	return &CertificateList{
		Items:     items,
		nextToken: token,
		pageSize:  pageSize,
		resource:  r,
	}, nil
}

func (c *Client) GetCertificate(ctx context.Context, r *Certificate) (*Certificate, error) {
	ctx = dcl.ContextWithRequestID(ctx)
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	// This is *purposefully* supressing errors.
	// This function is used with url-normalized values + not URL normalized values.
	// URL Normalized values will throw unintentional errors, since those values are not of the proper parent form.
	extractCertificateFields(r)

	b, err := c.getCertificateRaw(ctx, r)
	if err != nil {
		if dcl.IsNotFound(err) {
			return nil, &googleapi.Error{
				Code:    404,
				Message: err.Error(),
			}
		}
		return nil, err
	}
	result, err := unmarshalCertificate(b, c, r)
	if err != nil {
		return nil, err
	}
	result.Project = r.Project
	result.Location = r.Location
	result.CaPool = r.CaPool
	result.Name = r.Name

	c.Config.Logger.InfoWithContextf(ctx, "Retrieved raw result state: %v", result)
	c.Config.Logger.InfoWithContextf(ctx, "Canonicalizing with specified state: %v", r)
	result, err = canonicalizeCertificateNewState(c, result, r)
	if err != nil {
		return nil, err
	}
	if err := postReadExtractCertificateFields(result); err != nil {
		return result, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Created result state: %v", result)

	return result, nil
}

func (c *Client) DeleteCertificate(ctx context.Context, r *Certificate) error {
	ctx = dcl.ContextWithRequestID(ctx)
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	if r == nil {
		return fmt.Errorf("Certificate resource is nil")
	}
	c.Config.Logger.InfoWithContext(ctx, "Deleting Certificate...")
	deleteOp := deleteCertificateOperation{}
	return deleteOp.do(ctx, r, c)
}

// DeleteAllCertificate deletes all resources that the filter functions returns true on.
func (c *Client) DeleteAllCertificate(ctx context.Context, project, location, caPool string, filter func(*Certificate) bool) error {
	listObj, err := c.ListCertificate(ctx, project, location, caPool)
	if err != nil {
		return err
	}

	err = c.deleteAllCertificate(ctx, filter, listObj.Items)
	if err != nil {
		return err
	}
	for listObj.HasNext() {
		err = listObj.Next(ctx, c)
		if err != nil {
			return nil
		}
		err = c.deleteAllCertificate(ctx, filter, listObj.Items)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Client) ApplyCertificate(ctx context.Context, rawDesired *Certificate, opts ...dcl.ApplyOption) (*Certificate, error) {
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	ctx = dcl.ContextWithRequestID(ctx)
	var resultNewState *Certificate
	err := dcl.Do(ctx, func(ctx context.Context) (*dcl.RetryDetails, error) {
		newState, err := applyCertificateHelper(c, ctx, rawDesired, opts...)
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

func applyCertificateHelper(c *Client, ctx context.Context, rawDesired *Certificate, opts ...dcl.ApplyOption) (*Certificate, error) {
	c.Config.Logger.InfoWithContext(ctx, "Beginning ApplyCertificate...")
	c.Config.Logger.InfoWithContextf(ctx, "User specified desired state: %v", rawDesired)

	// 1.1: Validation of user-specified fields in desired state.
	if err := rawDesired.validate(); err != nil {
		return nil, err
	}

	if err := extractCertificateFields(rawDesired); err != nil {
		return nil, err
	}

	initial, desired, fieldDiffs, err := c.certificateDiffsForRawDesired(ctx, rawDesired, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create a diff: %w", err)
	}

	diffs, err := convertFieldDiffsToCertificateDiffs(c.Config, fieldDiffs, opts)
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
	var ops []certificateApiOperation
	if create {
		ops = append(ops, &createCertificateOperation{})
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
	return applyCertificateDiff(c, ctx, desired, rawDesired, ops, opts...)
}

func applyCertificateDiff(c *Client, ctx context.Context, desired *Certificate, rawDesired *Certificate, ops []certificateApiOperation, opts ...dcl.ApplyOption) (*Certificate, error) {
	// 3.1, 3.2a Retrieval of raw new state & canonicalization with desired state
	c.Config.Logger.InfoWithContext(ctx, "Retrieving raw new state...")
	rawNew, err := c.GetCertificate(ctx, desired)
	if err != nil {
		return nil, err
	}
	// Get additional values from the first response.
	// These values should be merged into the newState above.
	if len(ops) > 0 {
		lastOp := ops[len(ops)-1]
		if o, ok := lastOp.(*createCertificateOperation); ok {
			if r, hasR := o.FirstResponse(); hasR {

				c.Config.Logger.InfoWithContext(ctx, "Retrieving raw new state from operation...")

				fullResp, err := unmarshalMapCertificate(r, c, rawDesired)
				if err != nil {
					return nil, err
				}

				rawNew, err = canonicalizeCertificateNewState(c, rawNew, fullResp)
				if err != nil {
					return nil, err
				}
			}
		}
	}

	c.Config.Logger.InfoWithContextf(ctx, "Canonicalizing with raw desired state: %v", rawDesired)
	// 3.2b Canonicalization of raw new state using raw desired state
	newState, err := canonicalizeCertificateNewState(c, rawNew, rawDesired)
	if err != nil {
		return rawNew, err
	}

	c.Config.Logger.InfoWithContextf(ctx, "Created canonical new state: %v", newState)
	// 3.3 Comparison of the new state and raw desired state.
	// TODO(magic-modules-eng): EVENTUALLY_CONSISTENT_UPDATE
	newDesired, err := canonicalizeCertificateDesiredState(rawDesired, newState)
	if err != nil {
		return newState, err
	}

	if err := postReadExtractCertificateFields(newState); err != nil {
		return newState, err
	}

	// Need to ensure any transformations made here match acceptably in differ.
	if err := postReadExtractCertificateFields(newDesired); err != nil {
		return newState, err
	}

	c.Config.Logger.InfoWithContextf(ctx, "Diffing using canonicalized desired state: %v", newDesired)
	newDiffs, err := diffCertificate(c, newDesired, newState)
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
