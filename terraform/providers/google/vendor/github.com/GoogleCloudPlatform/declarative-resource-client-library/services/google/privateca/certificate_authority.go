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

type CertificateAuthority struct {
	Name                      *string                                         `json:"name"`
	Type                      *CertificateAuthorityTypeEnum                   `json:"type"`
	Config                    *CertificateAuthorityConfig                     `json:"config"`
	Lifetime                  *string                                         `json:"lifetime"`
	KeySpec                   *CertificateAuthorityKeySpec                    `json:"keySpec"`
	SubordinateConfig         *CertificateAuthoritySubordinateConfig          `json:"subordinateConfig"`
	Tier                      *CertificateAuthorityTierEnum                   `json:"tier"`
	State                     *CertificateAuthorityStateEnum                  `json:"state"`
	PemCaCertificates         []string                                        `json:"pemCaCertificates"`
	CaCertificateDescriptions []CertificateAuthorityCaCertificateDescriptions `json:"caCertificateDescriptions"`
	GcsBucket                 *string                                         `json:"gcsBucket"`
	AccessUrls                *CertificateAuthorityAccessUrls                 `json:"accessUrls"`
	CreateTime                *string                                         `json:"createTime"`
	UpdateTime                *string                                         `json:"updateTime"`
	DeleteTime                *string                                         `json:"deleteTime"`
	ExpireTime                *string                                         `json:"expireTime"`
	Labels                    map[string]string                               `json:"labels"`
	Project                   *string                                         `json:"project"`
	Location                  *string                                         `json:"location"`
	CaPool                    *string                                         `json:"caPool"`
}

func (r *CertificateAuthority) String() string {
	return dcl.SprintResource(r)
}

// The enum CertificateAuthorityTypeEnum.
type CertificateAuthorityTypeEnum string

// CertificateAuthorityTypeEnumRef returns a *CertificateAuthorityTypeEnum with the value of string s
// If the empty string is provided, nil is returned.
func CertificateAuthorityTypeEnumRef(s string) *CertificateAuthorityTypeEnum {
	v := CertificateAuthorityTypeEnum(s)
	return &v
}

func (v CertificateAuthorityTypeEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"SELF_SIGNED", "SUBORDINATE"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "CertificateAuthorityTypeEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum CertificateAuthorityConfigPublicKeyFormatEnum.
type CertificateAuthorityConfigPublicKeyFormatEnum string

// CertificateAuthorityConfigPublicKeyFormatEnumRef returns a *CertificateAuthorityConfigPublicKeyFormatEnum with the value of string s
// If the empty string is provided, nil is returned.
func CertificateAuthorityConfigPublicKeyFormatEnumRef(s string) *CertificateAuthorityConfigPublicKeyFormatEnum {
	v := CertificateAuthorityConfigPublicKeyFormatEnum(s)
	return &v
}

func (v CertificateAuthorityConfigPublicKeyFormatEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"PEM"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "CertificateAuthorityConfigPublicKeyFormatEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum CertificateAuthorityKeySpecAlgorithmEnum.
type CertificateAuthorityKeySpecAlgorithmEnum string

// CertificateAuthorityKeySpecAlgorithmEnumRef returns a *CertificateAuthorityKeySpecAlgorithmEnum with the value of string s
// If the empty string is provided, nil is returned.
func CertificateAuthorityKeySpecAlgorithmEnumRef(s string) *CertificateAuthorityKeySpecAlgorithmEnum {
	v := CertificateAuthorityKeySpecAlgorithmEnum(s)
	return &v
}

func (v CertificateAuthorityKeySpecAlgorithmEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"RSA_PSS_2048_SHA256", "RSA_PSS_3072_SHA256", "RSA_PSS_4096_SHA256", "RSA_PKCS1_2048_SHA256", "RSA_PKCS1_3072_SHA256", "RSA_PKCS1_4096_SHA256", "EC_P256_SHA256", "EC_P384_SHA384"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "CertificateAuthorityKeySpecAlgorithmEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum CertificateAuthorityTierEnum.
type CertificateAuthorityTierEnum string

// CertificateAuthorityTierEnumRef returns a *CertificateAuthorityTierEnum with the value of string s
// If the empty string is provided, nil is returned.
func CertificateAuthorityTierEnumRef(s string) *CertificateAuthorityTierEnum {
	v := CertificateAuthorityTierEnum(s)
	return &v
}

func (v CertificateAuthorityTierEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"ENTERPRISE", "DEVOPS"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "CertificateAuthorityTierEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum CertificateAuthorityStateEnum.
type CertificateAuthorityStateEnum string

// CertificateAuthorityStateEnumRef returns a *CertificateAuthorityStateEnum with the value of string s
// If the empty string is provided, nil is returned.
func CertificateAuthorityStateEnumRef(s string) *CertificateAuthorityStateEnum {
	v := CertificateAuthorityStateEnum(s)
	return &v
}

func (v CertificateAuthorityStateEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"ENABLED", "DISABLED", "STAGED", "AWAITING_USER_ACTIVATION", "DELETED"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "CertificateAuthorityStateEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum CertificateAuthorityCaCertificateDescriptionsPublicKeyFormatEnum.
type CertificateAuthorityCaCertificateDescriptionsPublicKeyFormatEnum string

// CertificateAuthorityCaCertificateDescriptionsPublicKeyFormatEnumRef returns a *CertificateAuthorityCaCertificateDescriptionsPublicKeyFormatEnum with the value of string s
// If the empty string is provided, nil is returned.
func CertificateAuthorityCaCertificateDescriptionsPublicKeyFormatEnumRef(s string) *CertificateAuthorityCaCertificateDescriptionsPublicKeyFormatEnum {
	v := CertificateAuthorityCaCertificateDescriptionsPublicKeyFormatEnum(s)
	return &v
}

func (v CertificateAuthorityCaCertificateDescriptionsPublicKeyFormatEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"PEM"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "CertificateAuthorityCaCertificateDescriptionsPublicKeyFormatEnum",
		Value: string(v),
		Valid: []string{},
	}
}

type CertificateAuthorityConfig struct {
	empty         bool                                     `json:"-"`
	SubjectConfig *CertificateAuthorityConfigSubjectConfig `json:"subjectConfig"`
	X509Config    *CertificateAuthorityConfigX509Config    `json:"x509Config"`
	PublicKey     *CertificateAuthorityConfigPublicKey     `json:"publicKey"`
}

type jsonCertificateAuthorityConfig CertificateAuthorityConfig

func (r *CertificateAuthorityConfig) UnmarshalJSON(data []byte) error {
	var res jsonCertificateAuthorityConfig
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyCertificateAuthorityConfig
	} else {

		r.SubjectConfig = res.SubjectConfig

		r.X509Config = res.X509Config

		r.PublicKey = res.PublicKey

	}
	return nil
}

// This object is used to assert a desired state where this CertificateAuthorityConfig is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyCertificateAuthorityConfig *CertificateAuthorityConfig = &CertificateAuthorityConfig{empty: true}

func (r *CertificateAuthorityConfig) Empty() bool {
	return r.empty
}

func (r *CertificateAuthorityConfig) String() string {
	return dcl.SprintResource(r)
}

func (r *CertificateAuthorityConfig) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type CertificateAuthorityConfigSubjectConfig struct {
	empty          bool                                                   `json:"-"`
	Subject        *CertificateAuthorityConfigSubjectConfigSubject        `json:"subject"`
	SubjectAltName *CertificateAuthorityConfigSubjectConfigSubjectAltName `json:"subjectAltName"`
}

type jsonCertificateAuthorityConfigSubjectConfig CertificateAuthorityConfigSubjectConfig

func (r *CertificateAuthorityConfigSubjectConfig) UnmarshalJSON(data []byte) error {
	var res jsonCertificateAuthorityConfigSubjectConfig
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyCertificateAuthorityConfigSubjectConfig
	} else {

		r.Subject = res.Subject

		r.SubjectAltName = res.SubjectAltName

	}
	return nil
}

// This object is used to assert a desired state where this CertificateAuthorityConfigSubjectConfig is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyCertificateAuthorityConfigSubjectConfig *CertificateAuthorityConfigSubjectConfig = &CertificateAuthorityConfigSubjectConfig{empty: true}

func (r *CertificateAuthorityConfigSubjectConfig) Empty() bool {
	return r.empty
}

func (r *CertificateAuthorityConfigSubjectConfig) String() string {
	return dcl.SprintResource(r)
}

func (r *CertificateAuthorityConfigSubjectConfig) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type CertificateAuthorityConfigSubjectConfigSubject struct {
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

type jsonCertificateAuthorityConfigSubjectConfigSubject CertificateAuthorityConfigSubjectConfigSubject

func (r *CertificateAuthorityConfigSubjectConfigSubject) UnmarshalJSON(data []byte) error {
	var res jsonCertificateAuthorityConfigSubjectConfigSubject
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyCertificateAuthorityConfigSubjectConfigSubject
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

// This object is used to assert a desired state where this CertificateAuthorityConfigSubjectConfigSubject is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyCertificateAuthorityConfigSubjectConfigSubject *CertificateAuthorityConfigSubjectConfigSubject = &CertificateAuthorityConfigSubjectConfigSubject{empty: true}

func (r *CertificateAuthorityConfigSubjectConfigSubject) Empty() bool {
	return r.empty
}

func (r *CertificateAuthorityConfigSubjectConfigSubject) String() string {
	return dcl.SprintResource(r)
}

func (r *CertificateAuthorityConfigSubjectConfigSubject) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type CertificateAuthorityConfigSubjectConfigSubjectAltName struct {
	empty          bool                                                              `json:"-"`
	DnsNames       []string                                                          `json:"dnsNames"`
	Uris           []string                                                          `json:"uris"`
	EmailAddresses []string                                                          `json:"emailAddresses"`
	IPAddresses    []string                                                          `json:"ipAddresses"`
	CustomSans     []CertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSans `json:"customSans"`
}

type jsonCertificateAuthorityConfigSubjectConfigSubjectAltName CertificateAuthorityConfigSubjectConfigSubjectAltName

func (r *CertificateAuthorityConfigSubjectConfigSubjectAltName) UnmarshalJSON(data []byte) error {
	var res jsonCertificateAuthorityConfigSubjectConfigSubjectAltName
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyCertificateAuthorityConfigSubjectConfigSubjectAltName
	} else {

		r.DnsNames = res.DnsNames

		r.Uris = res.Uris

		r.EmailAddresses = res.EmailAddresses

		r.IPAddresses = res.IPAddresses

		r.CustomSans = res.CustomSans

	}
	return nil
}

// This object is used to assert a desired state where this CertificateAuthorityConfigSubjectConfigSubjectAltName is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyCertificateAuthorityConfigSubjectConfigSubjectAltName *CertificateAuthorityConfigSubjectConfigSubjectAltName = &CertificateAuthorityConfigSubjectConfigSubjectAltName{empty: true}

func (r *CertificateAuthorityConfigSubjectConfigSubjectAltName) Empty() bool {
	return r.empty
}

func (r *CertificateAuthorityConfigSubjectConfigSubjectAltName) String() string {
	return dcl.SprintResource(r)
}

func (r *CertificateAuthorityConfigSubjectConfigSubjectAltName) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type CertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSans struct {
	empty    bool                                                                     `json:"-"`
	ObjectId *CertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSansObjectId `json:"objectId"`
	Critical *bool                                                                    `json:"critical"`
	Value    *string                                                                  `json:"value"`
}

type jsonCertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSans CertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSans

func (r *CertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSans) UnmarshalJSON(data []byte) error {
	var res jsonCertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSans
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyCertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSans
	} else {

		r.ObjectId = res.ObjectId

		r.Critical = res.Critical

		r.Value = res.Value

	}
	return nil
}

// This object is used to assert a desired state where this CertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSans is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyCertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSans *CertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSans = &CertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSans{empty: true}

func (r *CertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSans) Empty() bool {
	return r.empty
}

func (r *CertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSans) String() string {
	return dcl.SprintResource(r)
}

func (r *CertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSans) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type CertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSansObjectId struct {
	empty        bool    `json:"-"`
	ObjectIdPath []int64 `json:"objectIdPath"`
}

type jsonCertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSansObjectId CertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSansObjectId

func (r *CertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSansObjectId) UnmarshalJSON(data []byte) error {
	var res jsonCertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSansObjectId
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyCertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSansObjectId
	} else {

		r.ObjectIdPath = res.ObjectIdPath

	}
	return nil
}

// This object is used to assert a desired state where this CertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSansObjectId is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyCertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSansObjectId *CertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSansObjectId = &CertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSansObjectId{empty: true}

func (r *CertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSansObjectId) Empty() bool {
	return r.empty
}

func (r *CertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSansObjectId) String() string {
	return dcl.SprintResource(r)
}

func (r *CertificateAuthorityConfigSubjectConfigSubjectAltNameCustomSansObjectId) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type CertificateAuthorityConfigX509Config struct {
	empty                bool                                                       `json:"-"`
	KeyUsage             *CertificateAuthorityConfigX509ConfigKeyUsage              `json:"keyUsage"`
	CaOptions            *CertificateAuthorityConfigX509ConfigCaOptions             `json:"caOptions"`
	PolicyIds            []CertificateAuthorityConfigX509ConfigPolicyIds            `json:"policyIds"`
	AiaOcspServers       []string                                                   `json:"aiaOcspServers"`
	AdditionalExtensions []CertificateAuthorityConfigX509ConfigAdditionalExtensions `json:"additionalExtensions"`
}

type jsonCertificateAuthorityConfigX509Config CertificateAuthorityConfigX509Config

func (r *CertificateAuthorityConfigX509Config) UnmarshalJSON(data []byte) error {
	var res jsonCertificateAuthorityConfigX509Config
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyCertificateAuthorityConfigX509Config
	} else {

		r.KeyUsage = res.KeyUsage

		r.CaOptions = res.CaOptions

		r.PolicyIds = res.PolicyIds

		r.AiaOcspServers = res.AiaOcspServers

		r.AdditionalExtensions = res.AdditionalExtensions

	}
	return nil
}

// This object is used to assert a desired state where this CertificateAuthorityConfigX509Config is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyCertificateAuthorityConfigX509Config *CertificateAuthorityConfigX509Config = &CertificateAuthorityConfigX509Config{empty: true}

func (r *CertificateAuthorityConfigX509Config) Empty() bool {
	return r.empty
}

func (r *CertificateAuthorityConfigX509Config) String() string {
	return dcl.SprintResource(r)
}

func (r *CertificateAuthorityConfigX509Config) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type CertificateAuthorityConfigX509ConfigKeyUsage struct {
	empty                    bool                                                                   `json:"-"`
	BaseKeyUsage             *CertificateAuthorityConfigX509ConfigKeyUsageBaseKeyUsage              `json:"baseKeyUsage"`
	ExtendedKeyUsage         *CertificateAuthorityConfigX509ConfigKeyUsageExtendedKeyUsage          `json:"extendedKeyUsage"`
	UnknownExtendedKeyUsages []CertificateAuthorityConfigX509ConfigKeyUsageUnknownExtendedKeyUsages `json:"unknownExtendedKeyUsages"`
}

type jsonCertificateAuthorityConfigX509ConfigKeyUsage CertificateAuthorityConfigX509ConfigKeyUsage

func (r *CertificateAuthorityConfigX509ConfigKeyUsage) UnmarshalJSON(data []byte) error {
	var res jsonCertificateAuthorityConfigX509ConfigKeyUsage
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyCertificateAuthorityConfigX509ConfigKeyUsage
	} else {

		r.BaseKeyUsage = res.BaseKeyUsage

		r.ExtendedKeyUsage = res.ExtendedKeyUsage

		r.UnknownExtendedKeyUsages = res.UnknownExtendedKeyUsages

	}
	return nil
}

// This object is used to assert a desired state where this CertificateAuthorityConfigX509ConfigKeyUsage is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyCertificateAuthorityConfigX509ConfigKeyUsage *CertificateAuthorityConfigX509ConfigKeyUsage = &CertificateAuthorityConfigX509ConfigKeyUsage{empty: true}

func (r *CertificateAuthorityConfigX509ConfigKeyUsage) Empty() bool {
	return r.empty
}

func (r *CertificateAuthorityConfigX509ConfigKeyUsage) String() string {
	return dcl.SprintResource(r)
}

func (r *CertificateAuthorityConfigX509ConfigKeyUsage) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type CertificateAuthorityConfigX509ConfigKeyUsageBaseKeyUsage struct {
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

type jsonCertificateAuthorityConfigX509ConfigKeyUsageBaseKeyUsage CertificateAuthorityConfigX509ConfigKeyUsageBaseKeyUsage

func (r *CertificateAuthorityConfigX509ConfigKeyUsageBaseKeyUsage) UnmarshalJSON(data []byte) error {
	var res jsonCertificateAuthorityConfigX509ConfigKeyUsageBaseKeyUsage
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyCertificateAuthorityConfigX509ConfigKeyUsageBaseKeyUsage
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

// This object is used to assert a desired state where this CertificateAuthorityConfigX509ConfigKeyUsageBaseKeyUsage is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyCertificateAuthorityConfigX509ConfigKeyUsageBaseKeyUsage *CertificateAuthorityConfigX509ConfigKeyUsageBaseKeyUsage = &CertificateAuthorityConfigX509ConfigKeyUsageBaseKeyUsage{empty: true}

func (r *CertificateAuthorityConfigX509ConfigKeyUsageBaseKeyUsage) Empty() bool {
	return r.empty
}

func (r *CertificateAuthorityConfigX509ConfigKeyUsageBaseKeyUsage) String() string {
	return dcl.SprintResource(r)
}

func (r *CertificateAuthorityConfigX509ConfigKeyUsageBaseKeyUsage) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type CertificateAuthorityConfigX509ConfigKeyUsageExtendedKeyUsage struct {
	empty           bool  `json:"-"`
	ServerAuth      *bool `json:"serverAuth"`
	ClientAuth      *bool `json:"clientAuth"`
	CodeSigning     *bool `json:"codeSigning"`
	EmailProtection *bool `json:"emailProtection"`
	TimeStamping    *bool `json:"timeStamping"`
	OcspSigning     *bool `json:"ocspSigning"`
}

type jsonCertificateAuthorityConfigX509ConfigKeyUsageExtendedKeyUsage CertificateAuthorityConfigX509ConfigKeyUsageExtendedKeyUsage

func (r *CertificateAuthorityConfigX509ConfigKeyUsageExtendedKeyUsage) UnmarshalJSON(data []byte) error {
	var res jsonCertificateAuthorityConfigX509ConfigKeyUsageExtendedKeyUsage
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyCertificateAuthorityConfigX509ConfigKeyUsageExtendedKeyUsage
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

// This object is used to assert a desired state where this CertificateAuthorityConfigX509ConfigKeyUsageExtendedKeyUsage is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyCertificateAuthorityConfigX509ConfigKeyUsageExtendedKeyUsage *CertificateAuthorityConfigX509ConfigKeyUsageExtendedKeyUsage = &CertificateAuthorityConfigX509ConfigKeyUsageExtendedKeyUsage{empty: true}

func (r *CertificateAuthorityConfigX509ConfigKeyUsageExtendedKeyUsage) Empty() bool {
	return r.empty
}

func (r *CertificateAuthorityConfigX509ConfigKeyUsageExtendedKeyUsage) String() string {
	return dcl.SprintResource(r)
}

func (r *CertificateAuthorityConfigX509ConfigKeyUsageExtendedKeyUsage) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type CertificateAuthorityConfigX509ConfigKeyUsageUnknownExtendedKeyUsages struct {
	empty        bool    `json:"-"`
	ObjectIdPath []int64 `json:"objectIdPath"`
}

type jsonCertificateAuthorityConfigX509ConfigKeyUsageUnknownExtendedKeyUsages CertificateAuthorityConfigX509ConfigKeyUsageUnknownExtendedKeyUsages

func (r *CertificateAuthorityConfigX509ConfigKeyUsageUnknownExtendedKeyUsages) UnmarshalJSON(data []byte) error {
	var res jsonCertificateAuthorityConfigX509ConfigKeyUsageUnknownExtendedKeyUsages
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyCertificateAuthorityConfigX509ConfigKeyUsageUnknownExtendedKeyUsages
	} else {

		r.ObjectIdPath = res.ObjectIdPath

	}
	return nil
}

// This object is used to assert a desired state where this CertificateAuthorityConfigX509ConfigKeyUsageUnknownExtendedKeyUsages is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyCertificateAuthorityConfigX509ConfigKeyUsageUnknownExtendedKeyUsages *CertificateAuthorityConfigX509ConfigKeyUsageUnknownExtendedKeyUsages = &CertificateAuthorityConfigX509ConfigKeyUsageUnknownExtendedKeyUsages{empty: true}

func (r *CertificateAuthorityConfigX509ConfigKeyUsageUnknownExtendedKeyUsages) Empty() bool {
	return r.empty
}

func (r *CertificateAuthorityConfigX509ConfigKeyUsageUnknownExtendedKeyUsages) String() string {
	return dcl.SprintResource(r)
}

func (r *CertificateAuthorityConfigX509ConfigKeyUsageUnknownExtendedKeyUsages) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type CertificateAuthorityConfigX509ConfigCaOptions struct {
	empty                   bool   `json:"-"`
	IsCa                    *bool  `json:"isCa"`
	MaxIssuerPathLength     *int64 `json:"maxIssuerPathLength"`
	ZeroMaxIssuerPathLength *bool  `json:"zeroMaxIssuerPathLength"`
}

type jsonCertificateAuthorityConfigX509ConfigCaOptions CertificateAuthorityConfigX509ConfigCaOptions

func (r *CertificateAuthorityConfigX509ConfigCaOptions) UnmarshalJSON(data []byte) error {
	var res jsonCertificateAuthorityConfigX509ConfigCaOptions
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyCertificateAuthorityConfigX509ConfigCaOptions
	} else {

		r.IsCa = res.IsCa

		r.MaxIssuerPathLength = res.MaxIssuerPathLength

		r.ZeroMaxIssuerPathLength = res.ZeroMaxIssuerPathLength

	}
	return nil
}

// This object is used to assert a desired state where this CertificateAuthorityConfigX509ConfigCaOptions is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyCertificateAuthorityConfigX509ConfigCaOptions *CertificateAuthorityConfigX509ConfigCaOptions = &CertificateAuthorityConfigX509ConfigCaOptions{empty: true}

func (r *CertificateAuthorityConfigX509ConfigCaOptions) Empty() bool {
	return r.empty
}

func (r *CertificateAuthorityConfigX509ConfigCaOptions) String() string {
	return dcl.SprintResource(r)
}

func (r *CertificateAuthorityConfigX509ConfigCaOptions) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type CertificateAuthorityConfigX509ConfigPolicyIds struct {
	empty        bool    `json:"-"`
	ObjectIdPath []int64 `json:"objectIdPath"`
}

type jsonCertificateAuthorityConfigX509ConfigPolicyIds CertificateAuthorityConfigX509ConfigPolicyIds

func (r *CertificateAuthorityConfigX509ConfigPolicyIds) UnmarshalJSON(data []byte) error {
	var res jsonCertificateAuthorityConfigX509ConfigPolicyIds
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyCertificateAuthorityConfigX509ConfigPolicyIds
	} else {

		r.ObjectIdPath = res.ObjectIdPath

	}
	return nil
}

// This object is used to assert a desired state where this CertificateAuthorityConfigX509ConfigPolicyIds is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyCertificateAuthorityConfigX509ConfigPolicyIds *CertificateAuthorityConfigX509ConfigPolicyIds = &CertificateAuthorityConfigX509ConfigPolicyIds{empty: true}

func (r *CertificateAuthorityConfigX509ConfigPolicyIds) Empty() bool {
	return r.empty
}

func (r *CertificateAuthorityConfigX509ConfigPolicyIds) String() string {
	return dcl.SprintResource(r)
}

func (r *CertificateAuthorityConfigX509ConfigPolicyIds) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type CertificateAuthorityConfigX509ConfigAdditionalExtensions struct {
	empty    bool                                                              `json:"-"`
	ObjectId *CertificateAuthorityConfigX509ConfigAdditionalExtensionsObjectId `json:"objectId"`
	Critical *bool                                                             `json:"critical"`
	Value    *string                                                           `json:"value"`
}

type jsonCertificateAuthorityConfigX509ConfigAdditionalExtensions CertificateAuthorityConfigX509ConfigAdditionalExtensions

func (r *CertificateAuthorityConfigX509ConfigAdditionalExtensions) UnmarshalJSON(data []byte) error {
	var res jsonCertificateAuthorityConfigX509ConfigAdditionalExtensions
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyCertificateAuthorityConfigX509ConfigAdditionalExtensions
	} else {

		r.ObjectId = res.ObjectId

		r.Critical = res.Critical

		r.Value = res.Value

	}
	return nil
}

// This object is used to assert a desired state where this CertificateAuthorityConfigX509ConfigAdditionalExtensions is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyCertificateAuthorityConfigX509ConfigAdditionalExtensions *CertificateAuthorityConfigX509ConfigAdditionalExtensions = &CertificateAuthorityConfigX509ConfigAdditionalExtensions{empty: true}

func (r *CertificateAuthorityConfigX509ConfigAdditionalExtensions) Empty() bool {
	return r.empty
}

func (r *CertificateAuthorityConfigX509ConfigAdditionalExtensions) String() string {
	return dcl.SprintResource(r)
}

func (r *CertificateAuthorityConfigX509ConfigAdditionalExtensions) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type CertificateAuthorityConfigX509ConfigAdditionalExtensionsObjectId struct {
	empty        bool    `json:"-"`
	ObjectIdPath []int64 `json:"objectIdPath"`
}

type jsonCertificateAuthorityConfigX509ConfigAdditionalExtensionsObjectId CertificateAuthorityConfigX509ConfigAdditionalExtensionsObjectId

func (r *CertificateAuthorityConfigX509ConfigAdditionalExtensionsObjectId) UnmarshalJSON(data []byte) error {
	var res jsonCertificateAuthorityConfigX509ConfigAdditionalExtensionsObjectId
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyCertificateAuthorityConfigX509ConfigAdditionalExtensionsObjectId
	} else {

		r.ObjectIdPath = res.ObjectIdPath

	}
	return nil
}

// This object is used to assert a desired state where this CertificateAuthorityConfigX509ConfigAdditionalExtensionsObjectId is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyCertificateAuthorityConfigX509ConfigAdditionalExtensionsObjectId *CertificateAuthorityConfigX509ConfigAdditionalExtensionsObjectId = &CertificateAuthorityConfigX509ConfigAdditionalExtensionsObjectId{empty: true}

func (r *CertificateAuthorityConfigX509ConfigAdditionalExtensionsObjectId) Empty() bool {
	return r.empty
}

func (r *CertificateAuthorityConfigX509ConfigAdditionalExtensionsObjectId) String() string {
	return dcl.SprintResource(r)
}

func (r *CertificateAuthorityConfigX509ConfigAdditionalExtensionsObjectId) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type CertificateAuthorityConfigPublicKey struct {
	empty  bool                                           `json:"-"`
	Key    *string                                        `json:"key"`
	Format *CertificateAuthorityConfigPublicKeyFormatEnum `json:"format"`
}

type jsonCertificateAuthorityConfigPublicKey CertificateAuthorityConfigPublicKey

func (r *CertificateAuthorityConfigPublicKey) UnmarshalJSON(data []byte) error {
	var res jsonCertificateAuthorityConfigPublicKey
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyCertificateAuthorityConfigPublicKey
	} else {

		r.Key = res.Key

		r.Format = res.Format

	}
	return nil
}

// This object is used to assert a desired state where this CertificateAuthorityConfigPublicKey is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyCertificateAuthorityConfigPublicKey *CertificateAuthorityConfigPublicKey = &CertificateAuthorityConfigPublicKey{empty: true}

func (r *CertificateAuthorityConfigPublicKey) Empty() bool {
	return r.empty
}

func (r *CertificateAuthorityConfigPublicKey) String() string {
	return dcl.SprintResource(r)
}

func (r *CertificateAuthorityConfigPublicKey) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type CertificateAuthorityKeySpec struct {
	empty              bool                                      `json:"-"`
	CloudKmsKeyVersion *string                                   `json:"cloudKmsKeyVersion"`
	Algorithm          *CertificateAuthorityKeySpecAlgorithmEnum `json:"algorithm"`
}

type jsonCertificateAuthorityKeySpec CertificateAuthorityKeySpec

func (r *CertificateAuthorityKeySpec) UnmarshalJSON(data []byte) error {
	var res jsonCertificateAuthorityKeySpec
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyCertificateAuthorityKeySpec
	} else {

		r.CloudKmsKeyVersion = res.CloudKmsKeyVersion

		r.Algorithm = res.Algorithm

	}
	return nil
}

// This object is used to assert a desired state where this CertificateAuthorityKeySpec is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyCertificateAuthorityKeySpec *CertificateAuthorityKeySpec = &CertificateAuthorityKeySpec{empty: true}

func (r *CertificateAuthorityKeySpec) Empty() bool {
	return r.empty
}

func (r *CertificateAuthorityKeySpec) String() string {
	return dcl.SprintResource(r)
}

func (r *CertificateAuthorityKeySpec) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type CertificateAuthoritySubordinateConfig struct {
	empty                bool                                                 `json:"-"`
	CertificateAuthority *string                                              `json:"certificateAuthority"`
	PemIssuerChain       *CertificateAuthoritySubordinateConfigPemIssuerChain `json:"pemIssuerChain"`
}

type jsonCertificateAuthoritySubordinateConfig CertificateAuthoritySubordinateConfig

func (r *CertificateAuthoritySubordinateConfig) UnmarshalJSON(data []byte) error {
	var res jsonCertificateAuthoritySubordinateConfig
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyCertificateAuthoritySubordinateConfig
	} else {

		r.CertificateAuthority = res.CertificateAuthority

		r.PemIssuerChain = res.PemIssuerChain

	}
	return nil
}

// This object is used to assert a desired state where this CertificateAuthoritySubordinateConfig is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyCertificateAuthoritySubordinateConfig *CertificateAuthoritySubordinateConfig = &CertificateAuthoritySubordinateConfig{empty: true}

func (r *CertificateAuthoritySubordinateConfig) Empty() bool {
	return r.empty
}

func (r *CertificateAuthoritySubordinateConfig) String() string {
	return dcl.SprintResource(r)
}

func (r *CertificateAuthoritySubordinateConfig) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type CertificateAuthoritySubordinateConfigPemIssuerChain struct {
	empty           bool     `json:"-"`
	PemCertificates []string `json:"pemCertificates"`
}

type jsonCertificateAuthoritySubordinateConfigPemIssuerChain CertificateAuthoritySubordinateConfigPemIssuerChain

func (r *CertificateAuthoritySubordinateConfigPemIssuerChain) UnmarshalJSON(data []byte) error {
	var res jsonCertificateAuthoritySubordinateConfigPemIssuerChain
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyCertificateAuthoritySubordinateConfigPemIssuerChain
	} else {

		r.PemCertificates = res.PemCertificates

	}
	return nil
}

// This object is used to assert a desired state where this CertificateAuthoritySubordinateConfigPemIssuerChain is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyCertificateAuthoritySubordinateConfigPemIssuerChain *CertificateAuthoritySubordinateConfigPemIssuerChain = &CertificateAuthoritySubordinateConfigPemIssuerChain{empty: true}

func (r *CertificateAuthoritySubordinateConfigPemIssuerChain) Empty() bool {
	return r.empty
}

func (r *CertificateAuthoritySubordinateConfigPemIssuerChain) String() string {
	return dcl.SprintResource(r)
}

func (r *CertificateAuthoritySubordinateConfigPemIssuerChain) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type CertificateAuthorityCaCertificateDescriptions struct {
	empty                     bool                                                             `json:"-"`
	SubjectDescription        *CertificateAuthorityCaCertificateDescriptionsSubjectDescription `json:"subjectDescription"`
	X509Description           *CertificateAuthorityCaCertificateDescriptionsX509Description    `json:"x509Description"`
	PublicKey                 *CertificateAuthorityCaCertificateDescriptionsPublicKey          `json:"publicKey"`
	SubjectKeyId              *CertificateAuthorityCaCertificateDescriptionsSubjectKeyId       `json:"subjectKeyId"`
	AuthorityKeyId            *CertificateAuthorityCaCertificateDescriptionsAuthorityKeyId     `json:"authorityKeyId"`
	CrlDistributionPoints     []string                                                         `json:"crlDistributionPoints"`
	AiaIssuingCertificateUrls []string                                                         `json:"aiaIssuingCertificateUrls"`
	CertFingerprint           *CertificateAuthorityCaCertificateDescriptionsCertFingerprint    `json:"certFingerprint"`
}

type jsonCertificateAuthorityCaCertificateDescriptions CertificateAuthorityCaCertificateDescriptions

func (r *CertificateAuthorityCaCertificateDescriptions) UnmarshalJSON(data []byte) error {
	var res jsonCertificateAuthorityCaCertificateDescriptions
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyCertificateAuthorityCaCertificateDescriptions
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

// This object is used to assert a desired state where this CertificateAuthorityCaCertificateDescriptions is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyCertificateAuthorityCaCertificateDescriptions *CertificateAuthorityCaCertificateDescriptions = &CertificateAuthorityCaCertificateDescriptions{empty: true}

func (r *CertificateAuthorityCaCertificateDescriptions) Empty() bool {
	return r.empty
}

func (r *CertificateAuthorityCaCertificateDescriptions) String() string {
	return dcl.SprintResource(r)
}

func (r *CertificateAuthorityCaCertificateDescriptions) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type CertificateAuthorityCaCertificateDescriptionsSubjectDescription struct {
	empty           bool                                                                           `json:"-"`
	Subject         *CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubject        `json:"subject"`
	SubjectAltName  *CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltName `json:"subjectAltName"`
	HexSerialNumber *string                                                                        `json:"hexSerialNumber"`
	Lifetime        *string                                                                        `json:"lifetime"`
	NotBeforeTime   *string                                                                        `json:"notBeforeTime"`
	NotAfterTime    *string                                                                        `json:"notAfterTime"`
}

type jsonCertificateAuthorityCaCertificateDescriptionsSubjectDescription CertificateAuthorityCaCertificateDescriptionsSubjectDescription

func (r *CertificateAuthorityCaCertificateDescriptionsSubjectDescription) UnmarshalJSON(data []byte) error {
	var res jsonCertificateAuthorityCaCertificateDescriptionsSubjectDescription
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyCertificateAuthorityCaCertificateDescriptionsSubjectDescription
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

// This object is used to assert a desired state where this CertificateAuthorityCaCertificateDescriptionsSubjectDescription is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyCertificateAuthorityCaCertificateDescriptionsSubjectDescription *CertificateAuthorityCaCertificateDescriptionsSubjectDescription = &CertificateAuthorityCaCertificateDescriptionsSubjectDescription{empty: true}

func (r *CertificateAuthorityCaCertificateDescriptionsSubjectDescription) Empty() bool {
	return r.empty
}

func (r *CertificateAuthorityCaCertificateDescriptionsSubjectDescription) String() string {
	return dcl.SprintResource(r)
}

func (r *CertificateAuthorityCaCertificateDescriptionsSubjectDescription) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubject struct {
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

type jsonCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubject CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubject

func (r *CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubject) UnmarshalJSON(data []byte) error {
	var res jsonCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubject
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubject
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

// This object is used to assert a desired state where this CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubject is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubject *CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubject = &CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubject{empty: true}

func (r *CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubject) Empty() bool {
	return r.empty
}

func (r *CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubject) String() string {
	return dcl.SprintResource(r)
}

func (r *CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubject) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltName struct {
	empty          bool                                                                                      `json:"-"`
	DnsNames       []string                                                                                  `json:"dnsNames"`
	Uris           []string                                                                                  `json:"uris"`
	EmailAddresses []string                                                                                  `json:"emailAddresses"`
	IPAddresses    []string                                                                                  `json:"ipAddresses"`
	CustomSans     []CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSans `json:"customSans"`
}

type jsonCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltName CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltName

func (r *CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltName) UnmarshalJSON(data []byte) error {
	var res jsonCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltName
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltName
	} else {

		r.DnsNames = res.DnsNames

		r.Uris = res.Uris

		r.EmailAddresses = res.EmailAddresses

		r.IPAddresses = res.IPAddresses

		r.CustomSans = res.CustomSans

	}
	return nil
}

// This object is used to assert a desired state where this CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltName is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltName *CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltName = &CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltName{empty: true}

func (r *CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltName) Empty() bool {
	return r.empty
}

func (r *CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltName) String() string {
	return dcl.SprintResource(r)
}

func (r *CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltName) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSans struct {
	empty    bool                                                                                             `json:"-"`
	ObjectId *CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSansObjectId `json:"objectId"`
	Critical *bool                                                                                            `json:"critical"`
	Value    *string                                                                                          `json:"value"`
}

type jsonCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSans CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSans

func (r *CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSans) UnmarshalJSON(data []byte) error {
	var res jsonCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSans
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSans
	} else {

		r.ObjectId = res.ObjectId

		r.Critical = res.Critical

		r.Value = res.Value

	}
	return nil
}

// This object is used to assert a desired state where this CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSans is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSans *CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSans = &CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSans{empty: true}

func (r *CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSans) Empty() bool {
	return r.empty
}

func (r *CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSans) String() string {
	return dcl.SprintResource(r)
}

func (r *CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSans) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSansObjectId struct {
	empty        bool    `json:"-"`
	ObjectIdPath []int64 `json:"objectIdPath"`
}

type jsonCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSansObjectId CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSansObjectId

func (r *CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSansObjectId) UnmarshalJSON(data []byte) error {
	var res jsonCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSansObjectId
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSansObjectId
	} else {

		r.ObjectIdPath = res.ObjectIdPath

	}
	return nil
}

// This object is used to assert a desired state where this CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSansObjectId is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyCertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSansObjectId *CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSansObjectId = &CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSansObjectId{empty: true}

func (r *CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSansObjectId) Empty() bool {
	return r.empty
}

func (r *CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSansObjectId) String() string {
	return dcl.SprintResource(r)
}

func (r *CertificateAuthorityCaCertificateDescriptionsSubjectDescriptionSubjectAltNameCustomSansObjectId) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type CertificateAuthorityCaCertificateDescriptionsX509Description struct {
	empty                bool                                                                               `json:"-"`
	KeyUsage             *CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsage              `json:"keyUsage"`
	CaOptions            *CertificateAuthorityCaCertificateDescriptionsX509DescriptionCaOptions             `json:"caOptions"`
	PolicyIds            []CertificateAuthorityCaCertificateDescriptionsX509DescriptionPolicyIds            `json:"policyIds"`
	AiaOcspServers       []string                                                                           `json:"aiaOcspServers"`
	AdditionalExtensions []CertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensions `json:"additionalExtensions"`
}

type jsonCertificateAuthorityCaCertificateDescriptionsX509Description CertificateAuthorityCaCertificateDescriptionsX509Description

func (r *CertificateAuthorityCaCertificateDescriptionsX509Description) UnmarshalJSON(data []byte) error {
	var res jsonCertificateAuthorityCaCertificateDescriptionsX509Description
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyCertificateAuthorityCaCertificateDescriptionsX509Description
	} else {

		r.KeyUsage = res.KeyUsage

		r.CaOptions = res.CaOptions

		r.PolicyIds = res.PolicyIds

		r.AiaOcspServers = res.AiaOcspServers

		r.AdditionalExtensions = res.AdditionalExtensions

	}
	return nil
}

// This object is used to assert a desired state where this CertificateAuthorityCaCertificateDescriptionsX509Description is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyCertificateAuthorityCaCertificateDescriptionsX509Description *CertificateAuthorityCaCertificateDescriptionsX509Description = &CertificateAuthorityCaCertificateDescriptionsX509Description{empty: true}

func (r *CertificateAuthorityCaCertificateDescriptionsX509Description) Empty() bool {
	return r.empty
}

func (r *CertificateAuthorityCaCertificateDescriptionsX509Description) String() string {
	return dcl.SprintResource(r)
}

func (r *CertificateAuthorityCaCertificateDescriptionsX509Description) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsage struct {
	empty                    bool                                                                                           `json:"-"`
	BaseKeyUsage             *CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageBaseKeyUsage              `json:"baseKeyUsage"`
	ExtendedKeyUsage         *CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageExtendedKeyUsage          `json:"extendedKeyUsage"`
	UnknownExtendedKeyUsages []CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageUnknownExtendedKeyUsages `json:"unknownExtendedKeyUsages"`
}

type jsonCertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsage CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsage

func (r *CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsage) UnmarshalJSON(data []byte) error {
	var res jsonCertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsage
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyCertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsage
	} else {

		r.BaseKeyUsage = res.BaseKeyUsage

		r.ExtendedKeyUsage = res.ExtendedKeyUsage

		r.UnknownExtendedKeyUsages = res.UnknownExtendedKeyUsages

	}
	return nil
}

// This object is used to assert a desired state where this CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsage is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyCertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsage *CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsage = &CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsage{empty: true}

func (r *CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsage) Empty() bool {
	return r.empty
}

func (r *CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsage) String() string {
	return dcl.SprintResource(r)
}

func (r *CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsage) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageBaseKeyUsage struct {
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

type jsonCertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageBaseKeyUsage CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageBaseKeyUsage

func (r *CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageBaseKeyUsage) UnmarshalJSON(data []byte) error {
	var res jsonCertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageBaseKeyUsage
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyCertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageBaseKeyUsage
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

// This object is used to assert a desired state where this CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageBaseKeyUsage is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyCertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageBaseKeyUsage *CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageBaseKeyUsage = &CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageBaseKeyUsage{empty: true}

func (r *CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageBaseKeyUsage) Empty() bool {
	return r.empty
}

func (r *CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageBaseKeyUsage) String() string {
	return dcl.SprintResource(r)
}

func (r *CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageBaseKeyUsage) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageExtendedKeyUsage struct {
	empty           bool  `json:"-"`
	ServerAuth      *bool `json:"serverAuth"`
	ClientAuth      *bool `json:"clientAuth"`
	CodeSigning     *bool `json:"codeSigning"`
	EmailProtection *bool `json:"emailProtection"`
	TimeStamping    *bool `json:"timeStamping"`
	OcspSigning     *bool `json:"ocspSigning"`
}

type jsonCertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageExtendedKeyUsage CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageExtendedKeyUsage

func (r *CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageExtendedKeyUsage) UnmarshalJSON(data []byte) error {
	var res jsonCertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageExtendedKeyUsage
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyCertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageExtendedKeyUsage
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

// This object is used to assert a desired state where this CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageExtendedKeyUsage is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyCertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageExtendedKeyUsage *CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageExtendedKeyUsage = &CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageExtendedKeyUsage{empty: true}

func (r *CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageExtendedKeyUsage) Empty() bool {
	return r.empty
}

func (r *CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageExtendedKeyUsage) String() string {
	return dcl.SprintResource(r)
}

func (r *CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageExtendedKeyUsage) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageUnknownExtendedKeyUsages struct {
	empty        bool    `json:"-"`
	ObjectIdPath []int64 `json:"objectIdPath"`
}

type jsonCertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageUnknownExtendedKeyUsages CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageUnknownExtendedKeyUsages

func (r *CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageUnknownExtendedKeyUsages) UnmarshalJSON(data []byte) error {
	var res jsonCertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageUnknownExtendedKeyUsages
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyCertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageUnknownExtendedKeyUsages
	} else {

		r.ObjectIdPath = res.ObjectIdPath

	}
	return nil
}

// This object is used to assert a desired state where this CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageUnknownExtendedKeyUsages is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyCertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageUnknownExtendedKeyUsages *CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageUnknownExtendedKeyUsages = &CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageUnknownExtendedKeyUsages{empty: true}

func (r *CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageUnknownExtendedKeyUsages) Empty() bool {
	return r.empty
}

func (r *CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageUnknownExtendedKeyUsages) String() string {
	return dcl.SprintResource(r)
}

func (r *CertificateAuthorityCaCertificateDescriptionsX509DescriptionKeyUsageUnknownExtendedKeyUsages) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type CertificateAuthorityCaCertificateDescriptionsX509DescriptionCaOptions struct {
	empty               bool   `json:"-"`
	IsCa                *bool  `json:"isCa"`
	MaxIssuerPathLength *int64 `json:"maxIssuerPathLength"`
}

type jsonCertificateAuthorityCaCertificateDescriptionsX509DescriptionCaOptions CertificateAuthorityCaCertificateDescriptionsX509DescriptionCaOptions

func (r *CertificateAuthorityCaCertificateDescriptionsX509DescriptionCaOptions) UnmarshalJSON(data []byte) error {
	var res jsonCertificateAuthorityCaCertificateDescriptionsX509DescriptionCaOptions
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyCertificateAuthorityCaCertificateDescriptionsX509DescriptionCaOptions
	} else {

		r.IsCa = res.IsCa

		r.MaxIssuerPathLength = res.MaxIssuerPathLength

	}
	return nil
}

// This object is used to assert a desired state where this CertificateAuthorityCaCertificateDescriptionsX509DescriptionCaOptions is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyCertificateAuthorityCaCertificateDescriptionsX509DescriptionCaOptions *CertificateAuthorityCaCertificateDescriptionsX509DescriptionCaOptions = &CertificateAuthorityCaCertificateDescriptionsX509DescriptionCaOptions{empty: true}

func (r *CertificateAuthorityCaCertificateDescriptionsX509DescriptionCaOptions) Empty() bool {
	return r.empty
}

func (r *CertificateAuthorityCaCertificateDescriptionsX509DescriptionCaOptions) String() string {
	return dcl.SprintResource(r)
}

func (r *CertificateAuthorityCaCertificateDescriptionsX509DescriptionCaOptions) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type CertificateAuthorityCaCertificateDescriptionsX509DescriptionPolicyIds struct {
	empty        bool    `json:"-"`
	ObjectIdPath []int64 `json:"objectIdPath"`
}

type jsonCertificateAuthorityCaCertificateDescriptionsX509DescriptionPolicyIds CertificateAuthorityCaCertificateDescriptionsX509DescriptionPolicyIds

func (r *CertificateAuthorityCaCertificateDescriptionsX509DescriptionPolicyIds) UnmarshalJSON(data []byte) error {
	var res jsonCertificateAuthorityCaCertificateDescriptionsX509DescriptionPolicyIds
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyCertificateAuthorityCaCertificateDescriptionsX509DescriptionPolicyIds
	} else {

		r.ObjectIdPath = res.ObjectIdPath

	}
	return nil
}

// This object is used to assert a desired state where this CertificateAuthorityCaCertificateDescriptionsX509DescriptionPolicyIds is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyCertificateAuthorityCaCertificateDescriptionsX509DescriptionPolicyIds *CertificateAuthorityCaCertificateDescriptionsX509DescriptionPolicyIds = &CertificateAuthorityCaCertificateDescriptionsX509DescriptionPolicyIds{empty: true}

func (r *CertificateAuthorityCaCertificateDescriptionsX509DescriptionPolicyIds) Empty() bool {
	return r.empty
}

func (r *CertificateAuthorityCaCertificateDescriptionsX509DescriptionPolicyIds) String() string {
	return dcl.SprintResource(r)
}

func (r *CertificateAuthorityCaCertificateDescriptionsX509DescriptionPolicyIds) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type CertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensions struct {
	empty    bool                                                                                      `json:"-"`
	ObjectId *CertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensionsObjectId `json:"objectId"`
	Critical *bool                                                                                     `json:"critical"`
	Value    *string                                                                                   `json:"value"`
}

type jsonCertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensions CertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensions

func (r *CertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensions) UnmarshalJSON(data []byte) error {
	var res jsonCertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensions
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyCertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensions
	} else {

		r.ObjectId = res.ObjectId

		r.Critical = res.Critical

		r.Value = res.Value

	}
	return nil
}

// This object is used to assert a desired state where this CertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensions is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyCertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensions *CertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensions = &CertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensions{empty: true}

func (r *CertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensions) Empty() bool {
	return r.empty
}

func (r *CertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensions) String() string {
	return dcl.SprintResource(r)
}

func (r *CertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensions) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type CertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensionsObjectId struct {
	empty        bool    `json:"-"`
	ObjectIdPath []int64 `json:"objectIdPath"`
}

type jsonCertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensionsObjectId CertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensionsObjectId

func (r *CertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensionsObjectId) UnmarshalJSON(data []byte) error {
	var res jsonCertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensionsObjectId
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyCertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensionsObjectId
	} else {

		r.ObjectIdPath = res.ObjectIdPath

	}
	return nil
}

// This object is used to assert a desired state where this CertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensionsObjectId is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyCertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensionsObjectId *CertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensionsObjectId = &CertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensionsObjectId{empty: true}

func (r *CertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensionsObjectId) Empty() bool {
	return r.empty
}

func (r *CertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensionsObjectId) String() string {
	return dcl.SprintResource(r)
}

func (r *CertificateAuthorityCaCertificateDescriptionsX509DescriptionAdditionalExtensionsObjectId) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type CertificateAuthorityCaCertificateDescriptionsPublicKey struct {
	empty  bool                                                              `json:"-"`
	Key    *string                                                           `json:"key"`
	Format *CertificateAuthorityCaCertificateDescriptionsPublicKeyFormatEnum `json:"format"`
}

type jsonCertificateAuthorityCaCertificateDescriptionsPublicKey CertificateAuthorityCaCertificateDescriptionsPublicKey

func (r *CertificateAuthorityCaCertificateDescriptionsPublicKey) UnmarshalJSON(data []byte) error {
	var res jsonCertificateAuthorityCaCertificateDescriptionsPublicKey
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyCertificateAuthorityCaCertificateDescriptionsPublicKey
	} else {

		r.Key = res.Key

		r.Format = res.Format

	}
	return nil
}

// This object is used to assert a desired state where this CertificateAuthorityCaCertificateDescriptionsPublicKey is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyCertificateAuthorityCaCertificateDescriptionsPublicKey *CertificateAuthorityCaCertificateDescriptionsPublicKey = &CertificateAuthorityCaCertificateDescriptionsPublicKey{empty: true}

func (r *CertificateAuthorityCaCertificateDescriptionsPublicKey) Empty() bool {
	return r.empty
}

func (r *CertificateAuthorityCaCertificateDescriptionsPublicKey) String() string {
	return dcl.SprintResource(r)
}

func (r *CertificateAuthorityCaCertificateDescriptionsPublicKey) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type CertificateAuthorityCaCertificateDescriptionsSubjectKeyId struct {
	empty bool    `json:"-"`
	KeyId *string `json:"keyId"`
}

type jsonCertificateAuthorityCaCertificateDescriptionsSubjectKeyId CertificateAuthorityCaCertificateDescriptionsSubjectKeyId

func (r *CertificateAuthorityCaCertificateDescriptionsSubjectKeyId) UnmarshalJSON(data []byte) error {
	var res jsonCertificateAuthorityCaCertificateDescriptionsSubjectKeyId
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyCertificateAuthorityCaCertificateDescriptionsSubjectKeyId
	} else {

		r.KeyId = res.KeyId

	}
	return nil
}

// This object is used to assert a desired state where this CertificateAuthorityCaCertificateDescriptionsSubjectKeyId is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyCertificateAuthorityCaCertificateDescriptionsSubjectKeyId *CertificateAuthorityCaCertificateDescriptionsSubjectKeyId = &CertificateAuthorityCaCertificateDescriptionsSubjectKeyId{empty: true}

func (r *CertificateAuthorityCaCertificateDescriptionsSubjectKeyId) Empty() bool {
	return r.empty
}

func (r *CertificateAuthorityCaCertificateDescriptionsSubjectKeyId) String() string {
	return dcl.SprintResource(r)
}

func (r *CertificateAuthorityCaCertificateDescriptionsSubjectKeyId) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type CertificateAuthorityCaCertificateDescriptionsAuthorityKeyId struct {
	empty bool    `json:"-"`
	KeyId *string `json:"keyId"`
}

type jsonCertificateAuthorityCaCertificateDescriptionsAuthorityKeyId CertificateAuthorityCaCertificateDescriptionsAuthorityKeyId

func (r *CertificateAuthorityCaCertificateDescriptionsAuthorityKeyId) UnmarshalJSON(data []byte) error {
	var res jsonCertificateAuthorityCaCertificateDescriptionsAuthorityKeyId
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyCertificateAuthorityCaCertificateDescriptionsAuthorityKeyId
	} else {

		r.KeyId = res.KeyId

	}
	return nil
}

// This object is used to assert a desired state where this CertificateAuthorityCaCertificateDescriptionsAuthorityKeyId is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyCertificateAuthorityCaCertificateDescriptionsAuthorityKeyId *CertificateAuthorityCaCertificateDescriptionsAuthorityKeyId = &CertificateAuthorityCaCertificateDescriptionsAuthorityKeyId{empty: true}

func (r *CertificateAuthorityCaCertificateDescriptionsAuthorityKeyId) Empty() bool {
	return r.empty
}

func (r *CertificateAuthorityCaCertificateDescriptionsAuthorityKeyId) String() string {
	return dcl.SprintResource(r)
}

func (r *CertificateAuthorityCaCertificateDescriptionsAuthorityKeyId) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type CertificateAuthorityCaCertificateDescriptionsCertFingerprint struct {
	empty      bool    `json:"-"`
	Sha256Hash *string `json:"sha256Hash"`
}

type jsonCertificateAuthorityCaCertificateDescriptionsCertFingerprint CertificateAuthorityCaCertificateDescriptionsCertFingerprint

func (r *CertificateAuthorityCaCertificateDescriptionsCertFingerprint) UnmarshalJSON(data []byte) error {
	var res jsonCertificateAuthorityCaCertificateDescriptionsCertFingerprint
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyCertificateAuthorityCaCertificateDescriptionsCertFingerprint
	} else {

		r.Sha256Hash = res.Sha256Hash

	}
	return nil
}

// This object is used to assert a desired state where this CertificateAuthorityCaCertificateDescriptionsCertFingerprint is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyCertificateAuthorityCaCertificateDescriptionsCertFingerprint *CertificateAuthorityCaCertificateDescriptionsCertFingerprint = &CertificateAuthorityCaCertificateDescriptionsCertFingerprint{empty: true}

func (r *CertificateAuthorityCaCertificateDescriptionsCertFingerprint) Empty() bool {
	return r.empty
}

func (r *CertificateAuthorityCaCertificateDescriptionsCertFingerprint) String() string {
	return dcl.SprintResource(r)
}

func (r *CertificateAuthorityCaCertificateDescriptionsCertFingerprint) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type CertificateAuthorityAccessUrls struct {
	empty                  bool     `json:"-"`
	CaCertificateAccessUrl *string  `json:"caCertificateAccessUrl"`
	CrlAccessUrls          []string `json:"crlAccessUrls"`
}

type jsonCertificateAuthorityAccessUrls CertificateAuthorityAccessUrls

func (r *CertificateAuthorityAccessUrls) UnmarshalJSON(data []byte) error {
	var res jsonCertificateAuthorityAccessUrls
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyCertificateAuthorityAccessUrls
	} else {

		r.CaCertificateAccessUrl = res.CaCertificateAccessUrl

		r.CrlAccessUrls = res.CrlAccessUrls

	}
	return nil
}

// This object is used to assert a desired state where this CertificateAuthorityAccessUrls is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyCertificateAuthorityAccessUrls *CertificateAuthorityAccessUrls = &CertificateAuthorityAccessUrls{empty: true}

func (r *CertificateAuthorityAccessUrls) Empty() bool {
	return r.empty
}

func (r *CertificateAuthorityAccessUrls) String() string {
	return dcl.SprintResource(r)
}

func (r *CertificateAuthorityAccessUrls) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.Sum256([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

// Describe returns a simple description of this resource to ensure that automated tools
// can identify it.
func (r *CertificateAuthority) Describe() dcl.ServiceTypeVersion {
	return dcl.ServiceTypeVersion{
		Service: "privateca",
		Type:    "CertificateAuthority",
		Version: "privateca",
	}
}

func (r *CertificateAuthority) ID() (string, error) {
	if err := extractCertificateAuthorityFields(r); err != nil {
		return "", err
	}
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"name":                        dcl.ValueOrEmptyString(nr.Name),
		"type":                        dcl.ValueOrEmptyString(nr.Type),
		"config":                      dcl.ValueOrEmptyString(nr.Config),
		"lifetime":                    dcl.ValueOrEmptyString(nr.Lifetime),
		"key_spec":                    dcl.ValueOrEmptyString(nr.KeySpec),
		"subordinate_config":          dcl.ValueOrEmptyString(nr.SubordinateConfig),
		"tier":                        dcl.ValueOrEmptyString(nr.Tier),
		"state":                       dcl.ValueOrEmptyString(nr.State),
		"pem_ca_certificates":         dcl.ValueOrEmptyString(nr.PemCaCertificates),
		"ca_certificate_descriptions": dcl.ValueOrEmptyString(nr.CaCertificateDescriptions),
		"gcs_bucket":                  dcl.ValueOrEmptyString(nr.GcsBucket),
		"access_urls":                 dcl.ValueOrEmptyString(nr.AccessUrls),
		"create_time":                 dcl.ValueOrEmptyString(nr.CreateTime),
		"update_time":                 dcl.ValueOrEmptyString(nr.UpdateTime),
		"delete_time":                 dcl.ValueOrEmptyString(nr.DeleteTime),
		"expire_time":                 dcl.ValueOrEmptyString(nr.ExpireTime),
		"labels":                      dcl.ValueOrEmptyString(nr.Labels),
		"project":                     dcl.ValueOrEmptyString(nr.Project),
		"location":                    dcl.ValueOrEmptyString(nr.Location),
		"ca_pool":                     dcl.ValueOrEmptyString(nr.CaPool),
	}
	return dcl.Nprintf("projects/{{project}}/locations/{{location}}/caPools/{{ca_pool}}/certificateAuthorities/{{name}}", params), nil
}

const CertificateAuthorityMaxPage = -1

type CertificateAuthorityList struct {
	Items []*CertificateAuthority

	nextToken string

	pageSize int32

	resource *CertificateAuthority
}

func (l *CertificateAuthorityList) HasNext() bool {
	return l.nextToken != ""
}

func (l *CertificateAuthorityList) Next(ctx context.Context, c *Client) error {
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	if !l.HasNext() {
		return fmt.Errorf("no next page")
	}
	items, token, err := c.listCertificateAuthority(ctx, l.resource, l.nextToken, l.pageSize)
	if err != nil {
		return err
	}
	l.Items = items
	l.nextToken = token
	return err
}

func (c *Client) ListCertificateAuthority(ctx context.Context, project, location, caPool string) (*CertificateAuthorityList, error) {
	ctx = dcl.ContextWithRequestID(ctx)
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	return c.ListCertificateAuthorityWithMaxResults(ctx, project, location, caPool, CertificateAuthorityMaxPage)

}

func (c *Client) ListCertificateAuthorityWithMaxResults(ctx context.Context, project, location, caPool string, pageSize int32) (*CertificateAuthorityList, error) {
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	// Create a resource object so that we can use proper url normalization methods.
	r := &CertificateAuthority{
		Project:  &project,
		Location: &location,
		CaPool:   &caPool,
	}
	items, token, err := c.listCertificateAuthority(ctx, r, "", pageSize)
	if err != nil {
		return nil, err
	}
	return &CertificateAuthorityList{
		Items:     items,
		nextToken: token,
		pageSize:  pageSize,
		resource:  r,
	}, nil
}

func (c *Client) GetCertificateAuthority(ctx context.Context, r *CertificateAuthority) (*CertificateAuthority, error) {
	ctx = dcl.ContextWithRequestID(ctx)
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	// This is *purposefully* supressing errors.
	// This function is used with url-normalized values + not URL normalized values.
	// URL Normalized values will throw unintentional errors, since those values are not of the proper parent form.
	extractCertificateAuthorityFields(r)

	b, err := c.getCertificateAuthorityRaw(ctx, r)
	if err != nil {
		if dcl.IsNotFound(err) {
			return nil, &googleapi.Error{
				Code:    404,
				Message: err.Error(),
			}
		}
		return nil, err
	}
	result, err := unmarshalCertificateAuthority(b, c, r)
	if err != nil {
		return nil, err
	}
	result.Project = r.Project
	result.Location = r.Location
	result.CaPool = r.CaPool
	result.Name = r.Name

	c.Config.Logger.InfoWithContextf(ctx, "Retrieved raw result state: %v", result)
	c.Config.Logger.InfoWithContextf(ctx, "Canonicalizing with specified state: %v", r)
	result, err = canonicalizeCertificateAuthorityNewState(c, result, r)
	if err != nil {
		return nil, err
	}
	if err := postReadExtractCertificateAuthorityFields(result); err != nil {
		return result, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Created result state: %v", result)

	return result, nil
}

func (c *Client) DeleteCertificateAuthority(ctx context.Context, r *CertificateAuthority) error {
	ctx = dcl.ContextWithRequestID(ctx)
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	if r == nil {
		return fmt.Errorf("CertificateAuthority resource is nil")
	}
	c.Config.Logger.InfoWithContext(ctx, "Deleting CertificateAuthority...")
	deleteOp := deleteCertificateAuthorityOperation{}
	return deleteOp.do(ctx, r, c)
}

// DeleteAllCertificateAuthority deletes all resources that the filter functions returns true on.
func (c *Client) DeleteAllCertificateAuthority(ctx context.Context, project, location, caPool string, filter func(*CertificateAuthority) bool) error {
	listObj, err := c.ListCertificateAuthority(ctx, project, location, caPool)
	if err != nil {
		return err
	}

	err = c.deleteAllCertificateAuthority(ctx, filter, listObj.Items)
	if err != nil {
		return err
	}
	for listObj.HasNext() {
		err = listObj.Next(ctx, c)
		if err != nil {
			return nil
		}
		err = c.deleteAllCertificateAuthority(ctx, filter, listObj.Items)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Client) ApplyCertificateAuthority(ctx context.Context, rawDesired *CertificateAuthority, opts ...dcl.ApplyOption) (*CertificateAuthority, error) {
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	ctx = dcl.ContextWithRequestID(ctx)
	var resultNewState *CertificateAuthority
	err := dcl.Do(ctx, func(ctx context.Context) (*dcl.RetryDetails, error) {
		newState, err := applyCertificateAuthorityHelper(c, ctx, rawDesired, opts...)
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

func applyCertificateAuthorityHelper(c *Client, ctx context.Context, rawDesired *CertificateAuthority, opts ...dcl.ApplyOption) (*CertificateAuthority, error) {
	c.Config.Logger.InfoWithContext(ctx, "Beginning ApplyCertificateAuthority...")
	c.Config.Logger.InfoWithContextf(ctx, "User specified desired state: %v", rawDesired)

	// 1.1: Validation of user-specified fields in desired state.
	if err := rawDesired.validate(); err != nil {
		return nil, err
	}

	if err := extractCertificateAuthorityFields(rawDesired); err != nil {
		return nil, err
	}

	initial, desired, fieldDiffs, err := c.certificateAuthorityDiffsForRawDesired(ctx, rawDesired, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create a diff: %w", err)
	}

	diffs, err := convertFieldDiffsToCertificateAuthorityDiffs(c.Config, fieldDiffs, opts)
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
	var ops []certificateAuthorityApiOperation
	if create {
		ops = append(ops, &createCertificateAuthorityOperation{})
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
	return applyCertificateAuthorityDiff(c, ctx, desired, rawDesired, ops, opts...)
}

func applyCertificateAuthorityDiff(c *Client, ctx context.Context, desired *CertificateAuthority, rawDesired *CertificateAuthority, ops []certificateAuthorityApiOperation, opts ...dcl.ApplyOption) (*CertificateAuthority, error) {
	// 3.1, 3.2a Retrieval of raw new state & canonicalization with desired state
	c.Config.Logger.InfoWithContext(ctx, "Retrieving raw new state...")
	rawNew, err := c.GetCertificateAuthority(ctx, desired)
	if err != nil {
		return nil, err
	}
	// Get additional values from the first response.
	// These values should be merged into the newState above.
	if len(ops) > 0 {
		lastOp := ops[len(ops)-1]
		if o, ok := lastOp.(*createCertificateAuthorityOperation); ok {
			if r, hasR := o.FirstResponse(); hasR {

				c.Config.Logger.InfoWithContext(ctx, "Retrieving raw new state from operation...")

				fullResp, err := unmarshalMapCertificateAuthority(r, c, rawDesired)
				if err != nil {
					return nil, err
				}

				rawNew, err = canonicalizeCertificateAuthorityNewState(c, rawNew, fullResp)
				if err != nil {
					return nil, err
				}
			}
		}
	}

	c.Config.Logger.InfoWithContextf(ctx, "Canonicalizing with raw desired state: %v", rawDesired)
	// 3.2b Canonicalization of raw new state using raw desired state
	newState, err := canonicalizeCertificateAuthorityNewState(c, rawNew, rawDesired)
	if err != nil {
		return rawNew, err
	}

	c.Config.Logger.InfoWithContextf(ctx, "Created canonical new state: %v", newState)
	// 3.3 Comparison of the new state and raw desired state.
	// TODO(magic-modules-eng): EVENTUALLY_CONSISTENT_UPDATE
	newDesired, err := canonicalizeCertificateAuthorityDesiredState(rawDesired, newState)
	if err != nil {
		return newState, err
	}

	if err := postReadExtractCertificateAuthorityFields(newState); err != nil {
		return newState, err
	}

	// Need to ensure any transformations made here match acceptably in differ.
	if err := postReadExtractCertificateAuthorityFields(newDesired); err != nil {
		return newState, err
	}

	c.Config.Logger.InfoWithContextf(ctx, "Diffing using canonicalized desired state: %v", newDesired)
	newDiffs, err := diffCertificateAuthority(c, newDesired, newState)
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
