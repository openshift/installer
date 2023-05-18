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
package containeraws

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"time"

	"google.golang.org/api/googleapi"
	"github.com/GoogleCloudPlatform/declarative-resource-client-library/dcl"
)

type Cluster struct {
	Name                   *string                        `json:"name"`
	Description            *string                        `json:"description"`
	Networking             *ClusterNetworking             `json:"networking"`
	AwsRegion              *string                        `json:"awsRegion"`
	ControlPlane           *ClusterControlPlane           `json:"controlPlane"`
	Authorization          *ClusterAuthorization          `json:"authorization"`
	State                  *ClusterStateEnum              `json:"state"`
	Endpoint               *string                        `json:"endpoint"`
	Uid                    *string                        `json:"uid"`
	Reconciling            *bool                          `json:"reconciling"`
	CreateTime             *string                        `json:"createTime"`
	UpdateTime             *string                        `json:"updateTime"`
	Etag                   *string                        `json:"etag"`
	Annotations            map[string]string              `json:"annotations"`
	WorkloadIdentityConfig *ClusterWorkloadIdentityConfig `json:"workloadIdentityConfig"`
	Project                *string                        `json:"project"`
	Location               *string                        `json:"location"`
	Fleet                  *ClusterFleet                  `json:"fleet"`
}

func (r *Cluster) String() string {
	return dcl.SprintResource(r)
}

// The enum ClusterControlPlaneRootVolumeVolumeTypeEnum.
type ClusterControlPlaneRootVolumeVolumeTypeEnum string

// ClusterControlPlaneRootVolumeVolumeTypeEnumRef returns a *ClusterControlPlaneRootVolumeVolumeTypeEnum with the value of string s
// If the empty string is provided, nil is returned.
func ClusterControlPlaneRootVolumeVolumeTypeEnumRef(s string) *ClusterControlPlaneRootVolumeVolumeTypeEnum {
	v := ClusterControlPlaneRootVolumeVolumeTypeEnum(s)
	return &v
}

func (v ClusterControlPlaneRootVolumeVolumeTypeEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"VOLUME_TYPE_UNSPECIFIED", "GP2", "GP3"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "ClusterControlPlaneRootVolumeVolumeTypeEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum ClusterControlPlaneMainVolumeVolumeTypeEnum.
type ClusterControlPlaneMainVolumeVolumeTypeEnum string

// ClusterControlPlaneMainVolumeVolumeTypeEnumRef returns a *ClusterControlPlaneMainVolumeVolumeTypeEnum with the value of string s
// If the empty string is provided, nil is returned.
func ClusterControlPlaneMainVolumeVolumeTypeEnumRef(s string) *ClusterControlPlaneMainVolumeVolumeTypeEnum {
	v := ClusterControlPlaneMainVolumeVolumeTypeEnum(s)
	return &v
}

func (v ClusterControlPlaneMainVolumeVolumeTypeEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"VOLUME_TYPE_UNSPECIFIED", "GP2", "GP3"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "ClusterControlPlaneMainVolumeVolumeTypeEnum",
		Value: string(v),
		Valid: []string{},
	}
}

// The enum ClusterStateEnum.
type ClusterStateEnum string

// ClusterStateEnumRef returns a *ClusterStateEnum with the value of string s
// If the empty string is provided, nil is returned.
func ClusterStateEnumRef(s string) *ClusterStateEnum {
	v := ClusterStateEnum(s)
	return &v
}

func (v ClusterStateEnum) Validate() error {
	if string(v) == "" {
		// Empty enum is okay.
		return nil
	}
	for _, s := range []string{"STATE_UNSPECIFIED", "PROVISIONING", "RUNNING", "RECONCILING", "STOPPING", "ERROR", "DEGRADED"} {
		if string(v) == s {
			return nil
		}
	}
	return &dcl.EnumInvalidError{
		Enum:  "ClusterStateEnum",
		Value: string(v),
		Valid: []string{},
	}
}

type ClusterNetworking struct {
	empty                    bool     `json:"-"`
	VPCId                    *string  `json:"vpcId"`
	PodAddressCidrBlocks     []string `json:"podAddressCidrBlocks"`
	ServiceAddressCidrBlocks []string `json:"serviceAddressCidrBlocks"`
}

type jsonClusterNetworking ClusterNetworking

func (r *ClusterNetworking) UnmarshalJSON(data []byte) error {
	var res jsonClusterNetworking
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyClusterNetworking
	} else {

		r.VPCId = res.VPCId

		r.PodAddressCidrBlocks = res.PodAddressCidrBlocks

		r.ServiceAddressCidrBlocks = res.ServiceAddressCidrBlocks

	}
	return nil
}

// This object is used to assert a desired state where this ClusterNetworking is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyClusterNetworking *ClusterNetworking = &ClusterNetworking{empty: true}

func (r *ClusterNetworking) Empty() bool {
	return r.empty
}

func (r *ClusterNetworking) String() string {
	return dcl.SprintResource(r)
}

func (r *ClusterNetworking) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type ClusterControlPlane struct {
	empty                     bool                                          `json:"-"`
	Version                   *string                                       `json:"version"`
	InstanceType              *string                                       `json:"instanceType"`
	SshConfig                 *ClusterControlPlaneSshConfig                 `json:"sshConfig"`
	SubnetIds                 []string                                      `json:"subnetIds"`
	ConfigEncryption          *ClusterControlPlaneConfigEncryption          `json:"configEncryption"`
	SecurityGroupIds          []string                                      `json:"securityGroupIds"`
	IamInstanceProfile        *string                                       `json:"iamInstanceProfile"`
	RootVolume                *ClusterControlPlaneRootVolume                `json:"rootVolume"`
	MainVolume                *ClusterControlPlaneMainVolume                `json:"mainVolume"`
	DatabaseEncryption        *ClusterControlPlaneDatabaseEncryption        `json:"databaseEncryption"`
	Tags                      map[string]string                             `json:"tags"`
	AwsServicesAuthentication *ClusterControlPlaneAwsServicesAuthentication `json:"awsServicesAuthentication"`
	ProxyConfig               *ClusterControlPlaneProxyConfig               `json:"proxyConfig"`
}

type jsonClusterControlPlane ClusterControlPlane

func (r *ClusterControlPlane) UnmarshalJSON(data []byte) error {
	var res jsonClusterControlPlane
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyClusterControlPlane
	} else {

		r.Version = res.Version

		r.InstanceType = res.InstanceType

		r.SshConfig = res.SshConfig

		r.SubnetIds = res.SubnetIds

		r.ConfigEncryption = res.ConfigEncryption

		r.SecurityGroupIds = res.SecurityGroupIds

		r.IamInstanceProfile = res.IamInstanceProfile

		r.RootVolume = res.RootVolume

		r.MainVolume = res.MainVolume

		r.DatabaseEncryption = res.DatabaseEncryption

		r.Tags = res.Tags

		r.AwsServicesAuthentication = res.AwsServicesAuthentication

		r.ProxyConfig = res.ProxyConfig

	}
	return nil
}

// This object is used to assert a desired state where this ClusterControlPlane is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyClusterControlPlane *ClusterControlPlane = &ClusterControlPlane{empty: true}

func (r *ClusterControlPlane) Empty() bool {
	return r.empty
}

func (r *ClusterControlPlane) String() string {
	return dcl.SprintResource(r)
}

func (r *ClusterControlPlane) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type ClusterControlPlaneSshConfig struct {
	empty      bool    `json:"-"`
	Ec2KeyPair *string `json:"ec2KeyPair"`
}

type jsonClusterControlPlaneSshConfig ClusterControlPlaneSshConfig

func (r *ClusterControlPlaneSshConfig) UnmarshalJSON(data []byte) error {
	var res jsonClusterControlPlaneSshConfig
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyClusterControlPlaneSshConfig
	} else {

		r.Ec2KeyPair = res.Ec2KeyPair

	}
	return nil
}

// This object is used to assert a desired state where this ClusterControlPlaneSshConfig is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyClusterControlPlaneSshConfig *ClusterControlPlaneSshConfig = &ClusterControlPlaneSshConfig{empty: true}

func (r *ClusterControlPlaneSshConfig) Empty() bool {
	return r.empty
}

func (r *ClusterControlPlaneSshConfig) String() string {
	return dcl.SprintResource(r)
}

func (r *ClusterControlPlaneSshConfig) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type ClusterControlPlaneConfigEncryption struct {
	empty     bool    `json:"-"`
	KmsKeyArn *string `json:"kmsKeyArn"`
}

type jsonClusterControlPlaneConfigEncryption ClusterControlPlaneConfigEncryption

func (r *ClusterControlPlaneConfigEncryption) UnmarshalJSON(data []byte) error {
	var res jsonClusterControlPlaneConfigEncryption
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyClusterControlPlaneConfigEncryption
	} else {

		r.KmsKeyArn = res.KmsKeyArn

	}
	return nil
}

// This object is used to assert a desired state where this ClusterControlPlaneConfigEncryption is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyClusterControlPlaneConfigEncryption *ClusterControlPlaneConfigEncryption = &ClusterControlPlaneConfigEncryption{empty: true}

func (r *ClusterControlPlaneConfigEncryption) Empty() bool {
	return r.empty
}

func (r *ClusterControlPlaneConfigEncryption) String() string {
	return dcl.SprintResource(r)
}

func (r *ClusterControlPlaneConfigEncryption) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type ClusterControlPlaneRootVolume struct {
	empty      bool                                         `json:"-"`
	SizeGib    *int64                                       `json:"sizeGib"`
	VolumeType *ClusterControlPlaneRootVolumeVolumeTypeEnum `json:"volumeType"`
	Iops       *int64                                       `json:"iops"`
	KmsKeyArn  *string                                      `json:"kmsKeyArn"`
}

type jsonClusterControlPlaneRootVolume ClusterControlPlaneRootVolume

func (r *ClusterControlPlaneRootVolume) UnmarshalJSON(data []byte) error {
	var res jsonClusterControlPlaneRootVolume
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyClusterControlPlaneRootVolume
	} else {

		r.SizeGib = res.SizeGib

		r.VolumeType = res.VolumeType

		r.Iops = res.Iops

		r.KmsKeyArn = res.KmsKeyArn

	}
	return nil
}

// This object is used to assert a desired state where this ClusterControlPlaneRootVolume is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyClusterControlPlaneRootVolume *ClusterControlPlaneRootVolume = &ClusterControlPlaneRootVolume{empty: true}

func (r *ClusterControlPlaneRootVolume) Empty() bool {
	return r.empty
}

func (r *ClusterControlPlaneRootVolume) String() string {
	return dcl.SprintResource(r)
}

func (r *ClusterControlPlaneRootVolume) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type ClusterControlPlaneMainVolume struct {
	empty      bool                                         `json:"-"`
	SizeGib    *int64                                       `json:"sizeGib"`
	VolumeType *ClusterControlPlaneMainVolumeVolumeTypeEnum `json:"volumeType"`
	Iops       *int64                                       `json:"iops"`
	KmsKeyArn  *string                                      `json:"kmsKeyArn"`
}

type jsonClusterControlPlaneMainVolume ClusterControlPlaneMainVolume

func (r *ClusterControlPlaneMainVolume) UnmarshalJSON(data []byte) error {
	var res jsonClusterControlPlaneMainVolume
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyClusterControlPlaneMainVolume
	} else {

		r.SizeGib = res.SizeGib

		r.VolumeType = res.VolumeType

		r.Iops = res.Iops

		r.KmsKeyArn = res.KmsKeyArn

	}
	return nil
}

// This object is used to assert a desired state where this ClusterControlPlaneMainVolume is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyClusterControlPlaneMainVolume *ClusterControlPlaneMainVolume = &ClusterControlPlaneMainVolume{empty: true}

func (r *ClusterControlPlaneMainVolume) Empty() bool {
	return r.empty
}

func (r *ClusterControlPlaneMainVolume) String() string {
	return dcl.SprintResource(r)
}

func (r *ClusterControlPlaneMainVolume) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type ClusterControlPlaneDatabaseEncryption struct {
	empty     bool    `json:"-"`
	KmsKeyArn *string `json:"kmsKeyArn"`
}

type jsonClusterControlPlaneDatabaseEncryption ClusterControlPlaneDatabaseEncryption

func (r *ClusterControlPlaneDatabaseEncryption) UnmarshalJSON(data []byte) error {
	var res jsonClusterControlPlaneDatabaseEncryption
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyClusterControlPlaneDatabaseEncryption
	} else {

		r.KmsKeyArn = res.KmsKeyArn

	}
	return nil
}

// This object is used to assert a desired state where this ClusterControlPlaneDatabaseEncryption is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyClusterControlPlaneDatabaseEncryption *ClusterControlPlaneDatabaseEncryption = &ClusterControlPlaneDatabaseEncryption{empty: true}

func (r *ClusterControlPlaneDatabaseEncryption) Empty() bool {
	return r.empty
}

func (r *ClusterControlPlaneDatabaseEncryption) String() string {
	return dcl.SprintResource(r)
}

func (r *ClusterControlPlaneDatabaseEncryption) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type ClusterControlPlaneAwsServicesAuthentication struct {
	empty           bool    `json:"-"`
	RoleArn         *string `json:"roleArn"`
	RoleSessionName *string `json:"roleSessionName"`
}

type jsonClusterControlPlaneAwsServicesAuthentication ClusterControlPlaneAwsServicesAuthentication

func (r *ClusterControlPlaneAwsServicesAuthentication) UnmarshalJSON(data []byte) error {
	var res jsonClusterControlPlaneAwsServicesAuthentication
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyClusterControlPlaneAwsServicesAuthentication
	} else {

		r.RoleArn = res.RoleArn

		r.RoleSessionName = res.RoleSessionName

	}
	return nil
}

// This object is used to assert a desired state where this ClusterControlPlaneAwsServicesAuthentication is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyClusterControlPlaneAwsServicesAuthentication *ClusterControlPlaneAwsServicesAuthentication = &ClusterControlPlaneAwsServicesAuthentication{empty: true}

func (r *ClusterControlPlaneAwsServicesAuthentication) Empty() bool {
	return r.empty
}

func (r *ClusterControlPlaneAwsServicesAuthentication) String() string {
	return dcl.SprintResource(r)
}

func (r *ClusterControlPlaneAwsServicesAuthentication) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type ClusterControlPlaneProxyConfig struct {
	empty         bool    `json:"-"`
	SecretArn     *string `json:"secretArn"`
	SecretVersion *string `json:"secretVersion"`
}

type jsonClusterControlPlaneProxyConfig ClusterControlPlaneProxyConfig

func (r *ClusterControlPlaneProxyConfig) UnmarshalJSON(data []byte) error {
	var res jsonClusterControlPlaneProxyConfig
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyClusterControlPlaneProxyConfig
	} else {

		r.SecretArn = res.SecretArn

		r.SecretVersion = res.SecretVersion

	}
	return nil
}

// This object is used to assert a desired state where this ClusterControlPlaneProxyConfig is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyClusterControlPlaneProxyConfig *ClusterControlPlaneProxyConfig = &ClusterControlPlaneProxyConfig{empty: true}

func (r *ClusterControlPlaneProxyConfig) Empty() bool {
	return r.empty
}

func (r *ClusterControlPlaneProxyConfig) String() string {
	return dcl.SprintResource(r)
}

func (r *ClusterControlPlaneProxyConfig) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type ClusterAuthorization struct {
	empty      bool                             `json:"-"`
	AdminUsers []ClusterAuthorizationAdminUsers `json:"adminUsers"`
}

type jsonClusterAuthorization ClusterAuthorization

func (r *ClusterAuthorization) UnmarshalJSON(data []byte) error {
	var res jsonClusterAuthorization
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyClusterAuthorization
	} else {

		r.AdminUsers = res.AdminUsers

	}
	return nil
}

// This object is used to assert a desired state where this ClusterAuthorization is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyClusterAuthorization *ClusterAuthorization = &ClusterAuthorization{empty: true}

func (r *ClusterAuthorization) Empty() bool {
	return r.empty
}

func (r *ClusterAuthorization) String() string {
	return dcl.SprintResource(r)
}

func (r *ClusterAuthorization) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type ClusterAuthorizationAdminUsers struct {
	empty    bool    `json:"-"`
	Username *string `json:"username"`
}

type jsonClusterAuthorizationAdminUsers ClusterAuthorizationAdminUsers

func (r *ClusterAuthorizationAdminUsers) UnmarshalJSON(data []byte) error {
	var res jsonClusterAuthorizationAdminUsers
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyClusterAuthorizationAdminUsers
	} else {

		r.Username = res.Username

	}
	return nil
}

// This object is used to assert a desired state where this ClusterAuthorizationAdminUsers is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyClusterAuthorizationAdminUsers *ClusterAuthorizationAdminUsers = &ClusterAuthorizationAdminUsers{empty: true}

func (r *ClusterAuthorizationAdminUsers) Empty() bool {
	return r.empty
}

func (r *ClusterAuthorizationAdminUsers) String() string {
	return dcl.SprintResource(r)
}

func (r *ClusterAuthorizationAdminUsers) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type ClusterWorkloadIdentityConfig struct {
	empty            bool    `json:"-"`
	IssuerUri        *string `json:"issuerUri"`
	WorkloadPool     *string `json:"workloadPool"`
	IdentityProvider *string `json:"identityProvider"`
}

type jsonClusterWorkloadIdentityConfig ClusterWorkloadIdentityConfig

func (r *ClusterWorkloadIdentityConfig) UnmarshalJSON(data []byte) error {
	var res jsonClusterWorkloadIdentityConfig
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyClusterWorkloadIdentityConfig
	} else {

		r.IssuerUri = res.IssuerUri

		r.WorkloadPool = res.WorkloadPool

		r.IdentityProvider = res.IdentityProvider

	}
	return nil
}

// This object is used to assert a desired state where this ClusterWorkloadIdentityConfig is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyClusterWorkloadIdentityConfig *ClusterWorkloadIdentityConfig = &ClusterWorkloadIdentityConfig{empty: true}

func (r *ClusterWorkloadIdentityConfig) Empty() bool {
	return r.empty
}

func (r *ClusterWorkloadIdentityConfig) String() string {
	return dcl.SprintResource(r)
}

func (r *ClusterWorkloadIdentityConfig) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

type ClusterFleet struct {
	empty      bool    `json:"-"`
	Project    *string `json:"project"`
	Membership *string `json:"membership"`
}

type jsonClusterFleet ClusterFleet

func (r *ClusterFleet) UnmarshalJSON(data []byte) error {
	var res jsonClusterFleet
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	var m map[string]interface{}
	json.Unmarshal(data, &m)

	if len(m) == 0 {
		*r = *EmptyClusterFleet
	} else {

		r.Project = res.Project

		r.Membership = res.Membership

	}
	return nil
}

// This object is used to assert a desired state where this ClusterFleet is
// empty. Go lacks global const objects, but this object should be treated
// as one. Modifying this object will have undesirable results.
var EmptyClusterFleet *ClusterFleet = &ClusterFleet{empty: true}

func (r *ClusterFleet) Empty() bool {
	return r.empty
}

func (r *ClusterFleet) String() string {
	return dcl.SprintResource(r)
}

func (r *ClusterFleet) HashCode() string {
	// Placeholder for a more complex hash method that handles ordering, etc
	// Hash resource body for easy comparison later
	hash := sha256.New().Sum([]byte(r.String()))
	return fmt.Sprintf("%x", hash)
}

// Describe returns a simple description of this resource to ensure that automated tools
// can identify it.
func (r *Cluster) Describe() dcl.ServiceTypeVersion {
	return dcl.ServiceTypeVersion{
		Service: "container_aws",
		Type:    "Cluster",
		Version: "containeraws",
	}
}

func (r *Cluster) ID() (string, error) {
	if err := extractClusterFields(r); err != nil {
		return "", err
	}
	nr := r.urlNormalized()
	params := map[string]interface{}{
		"name":                     dcl.ValueOrEmptyString(nr.Name),
		"description":              dcl.ValueOrEmptyString(nr.Description),
		"networking":               dcl.ValueOrEmptyString(nr.Networking),
		"aws_region":               dcl.ValueOrEmptyString(nr.AwsRegion),
		"control_plane":            dcl.ValueOrEmptyString(nr.ControlPlane),
		"authorization":            dcl.ValueOrEmptyString(nr.Authorization),
		"state":                    dcl.ValueOrEmptyString(nr.State),
		"endpoint":                 dcl.ValueOrEmptyString(nr.Endpoint),
		"uid":                      dcl.ValueOrEmptyString(nr.Uid),
		"reconciling":              dcl.ValueOrEmptyString(nr.Reconciling),
		"create_time":              dcl.ValueOrEmptyString(nr.CreateTime),
		"update_time":              dcl.ValueOrEmptyString(nr.UpdateTime),
		"etag":                     dcl.ValueOrEmptyString(nr.Etag),
		"annotations":              dcl.ValueOrEmptyString(nr.Annotations),
		"workload_identity_config": dcl.ValueOrEmptyString(nr.WorkloadIdentityConfig),
		"project":                  dcl.ValueOrEmptyString(nr.Project),
		"location":                 dcl.ValueOrEmptyString(nr.Location),
		"fleet":                    dcl.ValueOrEmptyString(nr.Fleet),
	}
	return dcl.Nprintf("projects/{{project}}/locations/{{location}}/awsClusters/{{name}}", params), nil
}

const ClusterMaxPage = -1

type ClusterList struct {
	Items []*Cluster

	nextToken string

	pageSize int32

	resource *Cluster
}

func (l *ClusterList) HasNext() bool {
	return l.nextToken != ""
}

func (l *ClusterList) Next(ctx context.Context, c *Client) error {
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	if !l.HasNext() {
		return fmt.Errorf("no next page")
	}
	items, token, err := c.listCluster(ctx, l.resource, l.nextToken, l.pageSize)
	if err != nil {
		return err
	}
	l.Items = items
	l.nextToken = token
	return err
}

func (c *Client) ListCluster(ctx context.Context, project, location string) (*ClusterList, error) {
	ctx = dcl.ContextWithRequestID(ctx)
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	return c.ListClusterWithMaxResults(ctx, project, location, ClusterMaxPage)

}

func (c *Client) ListClusterWithMaxResults(ctx context.Context, project, location string, pageSize int32) (*ClusterList, error) {
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	// Create a resource object so that we can use proper url normalization methods.
	r := &Cluster{
		Project:  &project,
		Location: &location,
	}
	items, token, err := c.listCluster(ctx, r, "", pageSize)
	if err != nil {
		return nil, err
	}
	return &ClusterList{
		Items:     items,
		nextToken: token,
		pageSize:  pageSize,
		resource:  r,
	}, nil
}

func (c *Client) GetCluster(ctx context.Context, r *Cluster) (*Cluster, error) {
	ctx = dcl.ContextWithRequestID(ctx)
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	// This is *purposefully* supressing errors.
	// This function is used with url-normalized values + not URL normalized values.
	// URL Normalized values will throw unintentional errors, since those values are not of the proper parent form.
	extractClusterFields(r)

	b, err := c.getClusterRaw(ctx, r)
	if err != nil {
		if dcl.IsNotFound(err) {
			return nil, &googleapi.Error{
				Code:    404,
				Message: err.Error(),
			}
		}
		return nil, err
	}
	result, err := unmarshalCluster(b, c, r)
	if err != nil {
		return nil, err
	}
	result.Project = r.Project
	result.Location = r.Location
	result.Name = r.Name

	c.Config.Logger.InfoWithContextf(ctx, "Retrieved raw result state: %v", result)
	c.Config.Logger.InfoWithContextf(ctx, "Canonicalizing with specified state: %v", r)
	result, err = canonicalizeClusterNewState(c, result, r)
	if err != nil {
		return nil, err
	}
	if err := postReadExtractClusterFields(result); err != nil {
		return result, err
	}
	c.Config.Logger.InfoWithContextf(ctx, "Created result state: %v", result)

	return result, nil
}

func (c *Client) DeleteCluster(ctx context.Context, r *Cluster) error {
	ctx = dcl.ContextWithRequestID(ctx)
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	if r == nil {
		return fmt.Errorf("Cluster resource is nil")
	}
	c.Config.Logger.InfoWithContext(ctx, "Deleting Cluster...")
	deleteOp := deleteClusterOperation{}
	return deleteOp.do(ctx, r, c)
}

// DeleteAllCluster deletes all resources that the filter functions returns true on.
func (c *Client) DeleteAllCluster(ctx context.Context, project, location string, filter func(*Cluster) bool) error {
	listObj, err := c.ListCluster(ctx, project, location)
	if err != nil {
		return err
	}

	err = c.deleteAllCluster(ctx, filter, listObj.Items)
	if err != nil {
		return err
	}
	for listObj.HasNext() {
		err = listObj.Next(ctx, c)
		if err != nil {
			return nil
		}
		err = c.deleteAllCluster(ctx, filter, listObj.Items)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Client) ApplyCluster(ctx context.Context, rawDesired *Cluster, opts ...dcl.ApplyOption) (*Cluster, error) {
	ctx, cancel := context.WithTimeout(ctx, c.Config.TimeoutOr(0*time.Second))
	defer cancel()

	ctx = dcl.ContextWithRequestID(ctx)
	var resultNewState *Cluster
	err := dcl.Do(ctx, func(ctx context.Context) (*dcl.RetryDetails, error) {
		newState, err := applyClusterHelper(c, ctx, rawDesired, opts...)
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

func applyClusterHelper(c *Client, ctx context.Context, rawDesired *Cluster, opts ...dcl.ApplyOption) (*Cluster, error) {
	c.Config.Logger.InfoWithContext(ctx, "Beginning ApplyCluster...")
	c.Config.Logger.InfoWithContextf(ctx, "User specified desired state: %v", rawDesired)

	// 1.1: Validation of user-specified fields in desired state.
	if err := rawDesired.validate(); err != nil {
		return nil, err
	}

	if err := extractClusterFields(rawDesired); err != nil {
		return nil, err
	}

	initial, desired, fieldDiffs, err := c.clusterDiffsForRawDesired(ctx, rawDesired, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create a diff: %w", err)
	}

	diffs, err := convertFieldDiffsToClusterDiffs(c.Config, fieldDiffs, opts)
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
	var ops []clusterApiOperation
	if create {
		ops = append(ops, &createClusterOperation{})
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
	return applyClusterDiff(c, ctx, desired, rawDesired, ops, opts...)
}

func applyClusterDiff(c *Client, ctx context.Context, desired *Cluster, rawDesired *Cluster, ops []clusterApiOperation, opts ...dcl.ApplyOption) (*Cluster, error) {
	// 3.1, 3.2a Retrieval of raw new state & canonicalization with desired state
	c.Config.Logger.InfoWithContext(ctx, "Retrieving raw new state...")
	rawNew, err := c.GetCluster(ctx, desired)
	if err != nil {
		return nil, err
	}
	// Get additional values from the first response.
	// These values should be merged into the newState above.
	if len(ops) > 0 {
		lastOp := ops[len(ops)-1]
		if o, ok := lastOp.(*createClusterOperation); ok {
			if r, hasR := o.FirstResponse(); hasR {

				c.Config.Logger.InfoWithContext(ctx, "Retrieving raw new state from operation...")

				fullResp, err := unmarshalMapCluster(r, c, rawDesired)
				if err != nil {
					return nil, err
				}

				rawNew, err = canonicalizeClusterNewState(c, rawNew, fullResp)
				if err != nil {
					return nil, err
				}
			}
		}
	}

	c.Config.Logger.InfoWithContextf(ctx, "Canonicalizing with raw desired state: %v", rawDesired)
	// 3.2b Canonicalization of raw new state using raw desired state
	newState, err := canonicalizeClusterNewState(c, rawNew, rawDesired)
	if err != nil {
		return rawNew, err
	}

	c.Config.Logger.InfoWithContextf(ctx, "Created canonical new state: %v", newState)
	// 3.3 Comparison of the new state and raw desired state.
	// TODO(magic-modules-eng): EVENTUALLY_CONSISTENT_UPDATE
	newDesired, err := canonicalizeClusterDesiredState(rawDesired, newState)
	if err != nil {
		return newState, err
	}

	if err := postReadExtractClusterFields(newState); err != nil {
		return newState, err
	}

	// Need to ensure any transformations made here match acceptably in differ.
	if err := postReadExtractClusterFields(newDesired); err != nil {
		return newState, err
	}

	c.Config.Logger.InfoWithContextf(ctx, "Diffing using canonicalized desired state: %v", newDesired)
	newDiffs, err := diffCluster(c, newDesired, newState)
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
